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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TerminalResources overrides the compute resources for the containers of a
// Terminal's Pod. Empty fields fall back to the controller's configured
// defaults.
type TerminalResources struct {
	// Terminal overrides the resources of the "terminal" container, where the
	// user's kubectl shell actually runs.
	// +optional
	Terminal corev1.ResourceRequirements `json:"terminal,omitzero"`

	// Proxy overrides the resources of the "proxy" sidecar container, which
	// holds the impersonation credential and forwards kubectl traffic.
	// +optional
	Proxy corev1.ResourceRequirements `json:"proxy,omitzero"`
}

// TerminalSpec defines the desired state of the Terminal.
type TerminalSpec struct {
	// Subject identifies which user this terminal session belongs to. It must
	// be a lowercase UUID, must match the identity of the user who created,
	// updated, or deleted this Terminal (enforced by the validating webhook),
	// and metadata.name must equal "term-" followed by the first 8 characters
	// of Subject.
	// +kubebuilder:validation:Pattern=`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="subject is immutable"
	// +required
	Subject string `json:"subject"`

	// Image overrides the terminal/proxy container image. Empty uses the
	// controller's configured default. Changing this after the Terminal's Pod
	// has already been created has no effect on the running Pod (Pod
	// container images are immutable in Kubernetes); delete the Terminal to
	// have it recreated with the new image.
	// +optional
	Image string `json:"image,omitempty"`

	// IdleTimeoutSeconds is how long to wait, after the last recorded
	// activity, before the controller deletes this Terminal. Zero uses the
	// controller's configured default.
	// +kubebuilder:validation:Minimum=0
	// +optional
	IdleTimeoutSeconds int64 `json:"idleTimeoutSeconds,omitempty"`

	// Resources overrides compute resources for the terminal/proxy
	// containers. Empty fields use the controller's configured defaults.
	// Changing this after the Terminal's Pod has already been created has no
	// effect on the running Pod; delete the Terminal to have it recreated.
	// +optional
	Resources TerminalResources `json:"resources,omitzero"`
}

// TerminalStatus defines the observed state of the Terminal.
type TerminalStatus struct {
	// ObservedGeneration is the most recent generation observed by the
	// controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase summarizes the current lifecycle state of the Terminal.
	// +kubebuilder:validation:Enum=Pending;Creating;Ready;Failed
	// +optional
	Phase string `json:"phase,omitempty"`

	// PodName is the name of the reconciled Pod, in the same namespace as
	// this Terminal.
	// +optional
	PodName string `json:"podName,omitempty"`

	// PodReady mirrors the Ready condition of the reconciled Pod.
	// +optional
	PodReady bool `json:"podReady,omitempty"`

	// LastActivityTime is patched by the caller (not the controller) every
	// time a tab attaches/execs into this session. The controller reads it
	// (falling back to metadata.creationTimestamp if it was never set) to
	// decide when to garbage-collect an idle Terminal.
	// +optional
	LastActivityTime *metav1.Time `json:"lastActivityTime,omitempty"`

	// Conditions store the status conditions of the Terminal.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Namespaced,shortName=term,categories={otterscale}
// +kubebuilder:printcolumn:name="Subject",type=string,JSONPath=`.spec.subject`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="LastActivity",type=date,JSONPath=`.status.lastActivityTime`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'term-' + self.spec.subject.substring(0,8)",message="metadata.name must be term-<first 8 characters of spec.subject>"

// Terminal is the Schema for the terminals API.
// A Terminal represents a single user's interactive kubectl session Pod,
// reconciled by the console operator.
type Terminal struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired behavior of the Terminal.
	// +required
	Spec TerminalSpec `json:"spec"`

	// Status represents the current information about the Terminal.
	// +optional
	Status TerminalStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// TerminalList contains a list of Terminal resources.
type TerminalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Terminal `json:"items"`
}
