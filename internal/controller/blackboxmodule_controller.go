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
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

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

// +kubebuilder:rbac:groups=module.monitoring.ruup.amadeus.net,resources=blackboxmodules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=module.monitoring.ruup.amadeus.net,resources=blackboxmodules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=module.monitoring.ruup.amadeus.net,resources=blackboxmodules/finalizers,verbs=update

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
	logger := log.FromContext(ctx)

	// List all BlackboxModule resources
	var blackboxModules modulev1alpha1.BlackboxModuleList
	if err := r.List(ctx, &blackboxModules); err != nil {
		logger.Error(err, "unable to list BlackboxModules")
		return ctrl.Result{}, err
	}

	// Generate the new config
	newConfig := make(map[string]interface{})
	modules := make(map[string]interface{})

	for _, module := range blackboxModules.Items {
		modules[module.Name] = module.Spec
	}
	newConfig["modules"] = modules

	newConfigYAML, err := yaml.Marshal(newConfig)
	if err != nil {
		logger.Error(err, "unable to marshal new blackbox config")
		return ctrl.Result{}, err
	}

	var configMap corev1.ConfigMap
	err = r.Get(ctx, types.NamespacedName{Name: r.ConfigMapName, Namespace: r.ConfigMapNamespace}, &configMap)
	if err != nil {
		logger.Error(err, "unable to get ConfigMap")
		return ctrl.Result{}, err
	}

	// Compare and update if necessary
	currentConfigYAML, ok := configMap.Data["config.yml"]
	if !ok || !reflect.DeepEqual(string(newConfigYAML), currentConfigYAML) {
		logger.Info("Updating Blackbox Exporter ConfigMap")
		configMap.Data["config.yml"] = string(newConfigYAML)
		if err := r.Update(ctx, &configMap); err != nil {
			logger.Error(err, "unable to update ConfigMap")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Blackbox Exporter configuration is already up to date")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlackboxModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulev1alpha1.BlackboxModule{}).
		Complete(r)
}
