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

package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *BlackboxModule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-module-monitoring-ruup-amadeus-net-v1alpha1-blackboxmodule,mutating=false,failurePolicy=fail,sideEffects=None,groups=module.monitoring.ruup.amadeus.net,resources=blackboxmodules,verbs=create;update,versions=v1alpha1,name=vblackboxmodule.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &BlackboxModule{}

// ValidateCreate implements webhook.Validator
func (r *BlackboxModule) ValidateCreate() (admission.Warnings, error) {
	return r.validateBlackboxModule()
}

// ValidateUpdate implements webhook.Validator
func (r *BlackboxModule) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	return r.validateBlackboxModule()
}

// ValidateDelete implements webhook.Validator
func (r *BlackboxModule) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (r *BlackboxModule) validateBlackboxModule() (admission.Warnings, error) {
	switch r.Spec.Prober {
	case "http":
		if r.Spec.HTTP == nil {
			return nil, fmt.Errorf("http probe config is required when prober is 'http'")
		}
	case "tcp":
		if r.Spec.TCP == nil {
			return nil, fmt.Errorf("tcp probe config is required when prober is 'tcp'")
		}
	case "icmp":
		if r.Spec.ICMP == nil {
			return nil, fmt.Errorf("icmp probe config is required when prober is 'icmp'")
		}
	case "dns":
		if r.Spec.DNS == nil {
			return nil, fmt.Errorf("dns probe config is required when prober is 'dns'")
		}
	case "grpc":
		if r.Spec.GRPC == nil {
			return nil, fmt.Errorf("grpc probe config is required when prober is 'grpc'")
		}
	}

	// Warn if unrelated probe sections are set
	var warnings admission.Warnings
	if r.Spec.Prober != "http" && r.Spec.HTTP != nil {
		warnings = append(warnings, "http config is set but prober is not 'http'; it will be ignored")
	}
	if r.Spec.Prober != "tcp" && r.Spec.TCP != nil {
		warnings = append(warnings, "tcp config is set but prober is not 'tcp'; it will be ignored")
	}
	if r.Spec.Prober != "icmp" && r.Spec.ICMP != nil {
		warnings = append(warnings, "icmp config is set but prober is not 'icmp'; it will be ignored")
	}
	if r.Spec.Prober != "dns" && r.Spec.DNS != nil {
		warnings = append(warnings, "dns config is set but prober is not 'dns'; it will be ignored")
	}
	if r.Spec.Prober != "grpc" && r.Spec.GRPC != nil {
		warnings = append(warnings, "grpc config is set but prober is not 'grpc'; it will be ignored")
	}

	return warnings, nil
}
