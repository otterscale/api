/*
Copyright 2026 The OtterScale Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ModuleSpec defines the desired state of an installed Module.
// A Module instantiates a ModuleTemplate by referencing it and optionally
// overriding the target namespace or Helm values.
type ModuleSpec struct {
	// TemplateRef is the name of the ModuleTemplate to instantiate.
	// This field is immutable after creation.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="templateRef is immutable"
	// +required
	TemplateRef string `json:"templateRef"`

	// Namespace overrides the default target namespace defined in the ModuleTemplate.
	// If not specified, the namespace from the ModuleTemplate is used.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?)$`
	// +optional
	Namespace *string `json:"namespace,omitempty"`

	// Values overrides the default Helm chart values for Helm-based modules.
	// Only applicable when the referenced ModuleTemplate uses a HelmChart template.
	// Ignored for Kustomization-based modules.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Values *runtime.RawExtension `json:"values,omitempty"`

	// ApprovedTemplateGeneration is the ModuleTemplate generation that has been
	// approved for deployment. When the referenced ModuleTemplate's generation
	// exceeds this value, the controller will not apply changes until this field
	// is updated to match or exceed the new generation.
	//
	// Leave unset (nil) to auto-approve all template changes (legacy behavior).
	// Set explicitly to enable manual upgrade approval.
	// +optional
	ApprovedTemplateGeneration *int64 `json:"approvedTemplateGeneration,omitempty"`
}

// ModuleStatus defines the observed state of a Module.
// It tracks the template generation lifecycle and reports the health of
// the underlying Helm release or Kustomization.
type ModuleStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// It corresponds to the Module's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// AppliedTemplateGeneration is the ModuleTemplate generation that was last
	// successfully applied. The controller uses this to detect whether a
	// template upgrade is pending.
	// +optional
	AppliedTemplateGeneration int64 `json:"appliedTemplateGeneration,omitempty"`

	// AvailableTemplateGeneration is the latest generation of the referenced
	// ModuleTemplate. When this exceeds AppliedTemplateGeneration, an upgrade
	// is available.
	// +optional
	AvailableTemplateGeneration int64 `json:"availableTemplateGeneration,omitempty"`

	// Namespace is the resolved target namespace where resources are deployed.
	// It reflects the effective namespace (Module override or ModuleTemplate default).
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// HelmRelease captures the observed state of the Helm release
	// when the Module is backed by a HelmChart template.
	// +optional
	HelmRelease *HelmReleaseStatus `json:"helmRelease,omitempty"`

	// Kustomization captures the observed state of the Kustomize-based
	// deployment when the Module is backed by a Kustomization template.
	// +optional
	Kustomization *KustomizationStatus `json:"kustomization,omitempty"`

	// Inventory tracks the Kubernetes resources managed by this Module.
	// Used for garbage collection (pruning) of resources that are no longer
	// part of the desired state.
	// +optional
	Inventory []InventoryEntry `json:"inventory,omitempty"`

	// Conditions store the status conditions of the Module (e.g., Ready, UpgradeAvailable).
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Template",type=string,JSONPath=`.spec.templateRef`
// +kubebuilder:printcolumn:name="Namespace",type=string,JSONPath=`.status.namespace`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Outdated",type=string,JSONPath=`.status.conditions[?(@.type=="UpgradeAvailable")].status`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Module is the Schema for the modules API.
// A Module represents an installed platform addon instantiated from a ModuleTemplate.
// The controller manages the underlying Helm release or Kustomization directly
// and reflects its status back to the Module.
type Module struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired behavior of the Module.
	// +required
	Spec ModuleSpec `json:"spec"`

	// Status represents the current information about the Module.
	// +optional
	Status ModuleStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ModuleList contains a list of Module resources.
type ModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Module `json:"items"`
}
