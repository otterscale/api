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

// ApplicationSpec defines the desired state of an Application.
// It wraps a Deployment, optional Service, and optional PersistentVolumeClaim
// into a single declarative resource.
type ApplicationSpec struct {
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

// ApplicationStatus defines the observed state of an Application.
// It contains references to the actual Kubernetes resources created by the controller.
type ApplicationStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// It corresponds to the Application's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// DeploymentRef is a reference to the Deployment managed by this Application.
	// +optional
	DeploymentRef *ResourceReference `json:"deploymentRef,omitempty"`

	// ServiceRef is a reference to the Service managed by this Application.
	// +optional
	ServiceRef *ResourceReference `json:"serviceRef,omitempty"`

	// PersistentVolumeClaimRef is a reference to the PersistentVolumeClaim managed by this Application.
	// +optional
	PersistentVolumeClaimRef *ResourceReference `json:"persistentVolumeClaimRef,omitempty"`

	// Conditions store the status conditions of the Application (e.g., Ready, Progressing, Degraded).
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:categories={otterscale}
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Application is the Schema for the applications API.
// An Application provides a unified abstraction for deploying a workload
// with an optional Service and PersistentVolumeClaim.
type Application struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired behavior of the Application.
	// +required
	Spec ApplicationSpec `json:"spec"`

	// Status represents the current information about the Application.
	// +optional
	Status ApplicationStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application resources.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Application `json:"items"`
}
