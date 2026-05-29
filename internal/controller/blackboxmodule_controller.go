/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	modulev1alpha1 "github.com/luweglarz/blackbox-module-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// BlackboxModuleReconciler reconciles a BlackboxModule object
type BlackboxModuleReconciler struct {
	client.Client
	Scheme             *runtime.Scheme
	ConfigMapNamespace string
	ConfigMapName      string
}

const blackboxModuleFinalizer = "module.monitoring.ruup.amadeus.net/finalizer"

// +kubebuilder:rbac:groups=module.monitoring.ruup.amadeus.net,resources=blackboxmodules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=module.monitoring.ruup.amadeus.net,resources=blackboxmodules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=module.monitoring.ruup.amadeus.net,resources=blackboxmodules/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BlackboxModule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *BlackboxModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the triggering BlackboxModule to handle finalizer
	var module modulev1alpha1.BlackboxModule
	if err := r.Get(ctx, req.NamespacedName, &module); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// Object deleted and already gone — still re-sync the config
			return r.syncConfig(ctx)
		}
		return ctrl.Result{}, err
	}

	// Handle finalizer
	if module.DeletionTimestamp != nil {
		if controllerutil.ContainsFinalizer(&module, blackboxModuleFinalizer) {
			// Re-sync config without this module (it won't appear in List since it's terminating)
			result, err := r.syncConfig(ctx)
			if err != nil {
				return result, err
			}
			controllerutil.RemoveFinalizer(&module, blackboxModuleFinalizer)
			if err := r.Update(ctx, &module); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(&module, blackboxModuleFinalizer) {
		controllerutil.AddFinalizer(&module, blackboxModuleFinalizer)
		if err := r.Update(ctx, &module); err != nil {
			return ctrl.Result{}, err
		}
	}

	return r.syncConfig(ctx)
}

func (r *BlackboxModuleReconciler) syncConfig(ctx context.Context) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// List all BlackboxModule resources
	var blackboxModules modulev1alpha1.BlackboxModuleList
	if err := r.List(ctx, &blackboxModules); err != nil {
		logger.Error(err, "unable to list BlackboxModules")
		r.setConditionForAll(ctx, &blackboxModules, metav1.ConditionFalse, "ListFailed", err.Error())
		return ctrl.Result{}, err
	}

	// Generate the new config
	newConfig := make(map[string]interface{})
	modules := make(map[string]interface{})

	for _, module := range blackboxModules.Items {
		if module.DeletionTimestamp != nil {
			continue
		}
		modules[module.Name] = module.Spec
	}
	newConfig["modules"] = modules

	newConfigYAML, err := yaml.Marshal(newConfig)
	if err != nil {
		logger.Error(err, "unable to marshal new blackbox config")
		r.setConditionForAll(ctx, &blackboxModules, metav1.ConditionFalse, "MarshalFailed", err.Error())
		return ctrl.Result{}, err
	}

	var configMap corev1.ConfigMap
	err = r.Get(ctx, types.NamespacedName{Name: r.ConfigMapName, Namespace: r.ConfigMapNamespace}, &configMap)
	if err != nil {
		logger.Error(err, "unable to get ConfigMap")
		r.setConditionForAll(ctx, &blackboxModules, metav1.ConditionFalse, "ConfigMapNotFound", err.Error())
		return ctrl.Result{}, err
	}

	// Compare and update if necessary
	currentConfigYAML, ok := configMap.Data["config.yml"]
	if !ok || !reflect.DeepEqual(string(newConfigYAML), currentConfigYAML) {
		logger.Info("Updating Blackbox Exporter ConfigMap")
		configMap.Data["config.yml"] = string(newConfigYAML)
		if err := r.Update(ctx, &configMap); err != nil {
			logger.Error(err, "unable to update ConfigMap")
			r.setConditionForAll(ctx, &blackboxModules, metav1.ConditionFalse, "UpdateFailed", err.Error())
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Blackbox Exporter configuration is already up to date")
	}

	r.setConditionForAll(ctx, &blackboxModules, metav1.ConditionTrue, "ConfigSynced", "Module successfully synced to ConfigMap")
	return ctrl.Result{}, nil
}

func (r *BlackboxModuleReconciler) setConditionForAll(ctx context.Context, modules *modulev1alpha1.BlackboxModuleList, status metav1.ConditionStatus, reason, message string) {
	logger := log.FromContext(ctx)
	for i := range modules.Items {
		module := &modules.Items[i]
		condition := metav1.Condition{
			Type:               "ConfigSynced",
			Status:             status,
			ObservedGeneration: module.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             reason,
			Message:            fmt.Sprintf("%s/%s: %s", module.Namespace, module.Name, message),
		}
		changed := false
		for j, c := range module.Status.Conditions {
			if c.Type == "ConfigSynced" {
				if c.Status != status || c.Reason != reason {
					module.Status.Conditions[j] = condition
					changed = true
				}
				break
			}
		}
		if !changed && len(module.Status.Conditions) == 0 || !containsCondition(module.Status.Conditions, "ConfigSynced") {
			module.Status.Conditions = append(module.Status.Conditions, condition)
			changed = true
		}
		if changed {
			if err := r.Status().Update(ctx, module); err != nil {
				logger.Error(err, "unable to update status", "module", module.Name)
			}
		}
	}
}

func containsCondition(conditions []metav1.Condition, condType string) bool {
	for _, c := range conditions {
		if c.Type == condType {
			return true
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlackboxModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulev1alpha1.BlackboxModule{}).
		Watches(&corev1.ConfigMap{}, handler.EnqueueRequestsFromMapFunc(
			func(ctx context.Context, obj client.Object) []reconcile.Request {
				if obj.GetName() != r.ConfigMapName || obj.GetNamespace() != r.ConfigMapNamespace {
					return nil
				}
				// Re-queue all BlackboxModules to trigger a full re-sync
				var modules modulev1alpha1.BlackboxModuleList
				if err := mgr.GetClient().List(ctx, &modules); err != nil {
					return nil
				}
				var requests []reconcile.Request
				for _, m := range modules.Items {
					requests = append(requests, reconcile.Request{
						NamespacedName: types.NamespacedName{Name: m.Name, Namespace: m.Namespace},
					})
				}
				return requests
			},
		)).
		Complete(r)
}
