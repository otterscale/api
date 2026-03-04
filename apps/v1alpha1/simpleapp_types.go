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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimpleAppSpec defines the desired state of a SimpleApp.
// It wraps a Deployment, optional Service, and optional PersistentVolumeClaim
// into a single declarative resource.
type SimpleAppSpec struct {
	// Deployment defines the pod template, replicas, and update strategy
	// for the application workload.
	// +required
	Deployment appsv1.DeploymentSpec `json:"deployment"`

	// Service defines the Service configuration.
	// If specified, a Service will be created to expose the application.
	// +optional
	Service *corev1.ServiceSpec `json:"service,omitempty"`

	// PersistentVolumeClaim defines the PersistentVolumeClaim configuration.
	// If specified, a PersistentVolumeClaim will be created for persistent storage.
	// +optional
	PersistentVolumeClaim *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"`
}

// ResourceReference is a lightweight reference to a Kubernetes resource managed by the operator.
type ResourceReference struct {
	// Name is the name of the referenced resource.
	// +required
	Name string `json:"name"`

	// Namespace is the namespace of the referenced resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// SimpleAppStatus defines the observed state of a SimpleApp.
// It contains references to the actual Kubernetes resources created by the controller.
type SimpleAppStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// It corresponds to the SimpleApp's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// DeploymentRef is a reference to the Deployment managed by this SimpleApp.
	// +optional
	DeploymentRef *ResourceReference `json:"deploymentRef,omitempty"`

	// ServiceRef is a reference to the Service managed by this SimpleApp.
	// +optional
	ServiceRef *ResourceReference `json:"serviceRef,omitempty"`

	// PersistentVolumeClaimRef is a reference to the PersistentVolumeClaim managed by this SimpleApp.
	// +optional
	PersistentVolumeClaimRef *ResourceReference `json:"persistentVolumeClaimRef,omitempty"`

	// Conditions store the status conditions of the SimpleApp (e.g., Ready, Progressing, Degraded).
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SimpleApp is the Schema for the simpleapps API.
// A SimpleApp provides a simplified abstraction for deploying an application
// with an optional Service and PersistentVolumeClaim.
type SimpleApp struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired behavior of the SimpleApp.
	// +required
	Spec SimpleAppSpec `json:"spec"`

	// Status represents the current information about the SimpleApp.
	// +optional
	Status SimpleAppStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// SimpleAppList contains a list of SimpleApp resources.
type SimpleAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []SimpleApp `json:"items"`
}
