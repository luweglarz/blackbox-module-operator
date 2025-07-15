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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	modulev1alpha1 "github.com/luweglarz/blackbox-module-operator/api/v1alpha1"
)

var _ = Describe("BlackboxModule Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "blackbox-exporter-config"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "monitoring", // TODO(user):Modify as needed
		}
		blackboxmodule := &modulev1alpha1.BlackboxModule{}

		BeforeEach(func() {
			// Create a namespace if it doesn't exist
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "monitoring",
				},
			}
			_ = k8sClient.Create(ctx, namespace)
			// Create a ConfigMap with initial data
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "blackbox-exporter-config",
					Namespace: "monitoring",
				},
				Data: map[string]string{
					"config.yml": "initial: config", // Provide initial data as needed
				},
			}
			_ = k8sClient.Create(ctx, configMap) // Ignore error if already exists
			By("creating the custom resource for the Kind BlackboxModule")
			err := k8sClient.Get(ctx, typeNamespacedName, blackboxmodule)
			if err != nil && errors.IsNotFound(err) {
				resource := &modulev1alpha1.BlackboxModule{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "monitoring",
					},
					Spec: modulev1alpha1.BlackboxModuleSpec{
						Prober: "http", // Set to a supported value
						// Add other required spec fields if needed
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// Cleanup BlackboxModule resource
			resource := &modulev1alpha1.BlackboxModule{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			// Cleanup ConfigMap
			cm := &corev1.ConfigMap{}
			err = k8sClient.Get(ctx, typeNamespacedName, cm)
			if err == nil {
				Expect(k8sClient.Delete(ctx, cm)).To(Succeed())
			}

			// Cleanup monitoring namespace
			ns := &corev1.Namespace{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: "monitoring"}, ns)
			if err == nil {
				Expect(k8sClient.Delete(ctx, ns)).To(Succeed())
			}
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &BlackboxModuleReconciler{
				Client:             k8sClient,
				Scheme:             k8sClient.Scheme(),
				ConfigMapNamespace: "monitoring",               // Ensure this matches the namespace where the ConfigMap is created
				ConfigMapName:      "blackbox-exporter-config", // Ensure this matches the ConfigMap name
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})

// Go
var _ = Describe("BlackboxModule Controller - Aggregation", func() {
	ctx := context.Background()
	namespaces := []string{"ns1", "ns2"}
	crNames := []string{"bbm1", "bbm2"}

	BeforeEach(func() {
		Expect(k8sClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "monitoring2"},
		})).To(Succeed())
		Expect(k8sClient.Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "blackbox-exporter-config",
				Namespace: "monitoring2",
			},
			Data: map[string]string{"config.yml": "initial: config"},
		})).To(Succeed())

		// Create CRs in other namespaces
		for _, ns := range namespaces {
			Expect(k8sClient.Create(ctx, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: ns},
			})).To(Succeed())
		}
		for i, ns := range namespaces {
			Expect(k8sClient.Create(ctx, &modulev1alpha1.BlackboxModule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      crNames[i],
					Namespace: ns,
				},
				Spec: modulev1alpha1.BlackboxModuleSpec{
					Prober: "http",
					HTTP: &modulev1alpha1.HTTPProbe{
						Method: "GET",
						Headers: map[string]string{
							"X-Test": "true",
						},
					},
				},
			})).To(Succeed())
		}
	})

	AfterEach(func() {
		for i, ns := range namespaces {
			resource := &modulev1alpha1.BlackboxModule{}
			_ = k8sClient.Get(ctx, types.NamespacedName{Name: crNames[i], Namespace: ns}, resource)
			_ = k8sClient.Delete(ctx, resource)
		}
	})

	It("should aggregate CRs from multiple namespaces into the monitoring ConfigMap", func() {
		reconciler := &BlackboxModuleReconciler{
			Client:             k8sClient,
			Scheme:             k8sClient.Scheme(),
			ConfigMapNamespace: "monitoring2",              // Ensure this matches the namespace where the ConfigMap is created
			ConfigMapName:      "blackbox-exporter-config", // Ensure this matches the ConfigMap name
		}
		for i, ns := range namespaces {
			_, err := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: types.NamespacedName{Name: crNames[i], Namespace: ns},
			})
			Expect(err).NotTo(HaveOccurred())
		}

		cm := &corev1.ConfigMap{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "blackbox-exporter-config", Namespace: "monitoring2"}, cm)).To(Succeed())
		Expect(cm.Data["config.yml"]).NotTo(Equal("initial: config"))
	})
})
