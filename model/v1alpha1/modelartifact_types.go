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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PackFormat specifies the OCI artifact packaging format.
// +enum
type PackFormat string

const (
	// PackFormatModelPack produces a CNCF ModelPack compliant OCI artifact.
	PackFormatModelPack PackFormat = "ModelPack"

	// PackFormatModelKit produces a KitOps native ModelKit OCI artifact.
	PackFormatModelKit PackFormat = "ModelKit"
)

// ArtifactPhase represents the current lifecycle phase of a ModelArtifact.
// +enum
type ArtifactPhase string

const (
	// PhasePending indicates the pipeline has not yet started.
	PhasePending ArtifactPhase = "Pending"
	// PhaseRunning indicates the import/pack/push Job is in progress.
	PhaseRunning ArtifactPhase = "Running"
	// PhaseSucceeded indicates the artifact was successfully pushed to the registry.
	PhaseSucceeded ArtifactPhase = "Succeeded"
	// PhaseFailed indicates the pipeline encountered an error.
	PhaseFailed ArtifactPhase = "Failed"
)

// ModelArtifactSpec defines the desired state of a ModelArtifact.
// It declares the model source, target OCI registry, packaging format,
// and temporary storage for the import/pack/push pipeline.
type ModelArtifactSpec struct {
	// Source defines where to fetch the model from.
	// +required
	Source ModelSource `json:"source"`

	// Target defines the OCI registry destination for the packaged artifact.
	// +required
	Target OCITarget `json:"target"`

	// Format specifies the OCI artifact packaging format.
	// +kubebuilder:validation:Enum=ModelPack;ModelKit
	// +kubebuilder:default=ModelPack
	// +optional
	Format PackFormat `json:"format,omitempty"`

	// Storage configures the temporary PVC used during the import/pack/push pipeline.
	// The PVC is automatically cleaned up after the job completes.
	// +required
	Storage StorageSpec `json:"storage"`
}

// ModelSource defines the origin of the model to be packaged.
// Exactly one source type must be specified.
// +kubebuilder:validation:XValidation:rule="has(self.huggingFace)",message="at least one source must be specified"
type ModelSource struct {
	// HuggingFace specifies a HuggingFace Hub repository as the model source.
	// +optional
	HuggingFace *HuggingFaceSource `json:"huggingFace,omitempty"`
}

// HuggingFaceSource configures model retrieval from HuggingFace Hub.
//
// SECURITY: Repository and Revision are passed to shell scripts. Only users who can
// create ModelArtifacts should have access; they already have equivalent privileges.
type HuggingFaceSource struct {
	// Repository is the HuggingFace model repository identifier (e.g. "microsoft/phi-4").
	// Must contain only alphanumerics, dots, underscores, hyphens, and slashes.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9][a-zA-Z0-9._/-]*$"
	// +required
	Repository string `json:"repository"`

	// Revision pins a specific branch, tag, or commit hash.
	// If not specified, the default branch is used.
	// Must contain only alphanumerics, dots, underscores, hyphens, and slashes.
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9._/-]*$"
	// +optional
	Revision string `json:"revision,omitempty"`

	// TokenSecretRef references a Secret containing the HuggingFace access token.
	// Required for private or gated repositories.
	// +optional
	TokenSecretRef *SecretKeySelector `json:"tokenSecretRef,omitempty"`
}

// OCITarget defines the destination OCI registry for the packaged artifact.
//
// SECURITY: Repository and Tag are passed to shell scripts. Only users who can
// create ModelArtifacts should have access; they already have equivalent privileges.
type OCITarget struct {
	// Repository is the full OCI registry path (e.g. "ghcr.io/myorg/models/phi-4").
	// Must contain only alphanumerics, dots, underscores, hyphens, and slashes.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9][a-zA-Z0-9._/-]*$"
	// +required
	Repository string `json:"repository"`

	// Tag is the image tag to push. Defaults to "latest" if not specified.
	// Must contain only alphanumerics, dots, underscores, hyphens, and slashes.
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9._/-]*$"
	// +optional
	Tag string `json:"tag,omitempty"`

	// CredentialsSecretRef references a Secret containing OCI registry credentials.
	// The Secret must contain "username" and "password" keys.
	// +optional
	CredentialsSecretRef *SecretReference `json:"credentialsSecretRef,omitempty"`

	// Insecure uses an unencrypted connection to the registry instead of TLS.
	// Only use for development or air-gapped environments.
	// +optional
	Insecure bool `json:"insecure,omitempty"`
}

// StorageSpec configures the temporary PVC for the import/pack/push pipeline.
type StorageSpec struct {
	// Size is the requested PVC storage capacity (e.g. "100Gi").
	// Should be at least 2x the expected model size to accommodate both
	// the downloaded files and the packed artifact.
	// +required
	Size resource.Quantity `json:"size"`

	// StorageClassName overrides the cluster default StorageClass.
	// If not specified, the cluster default StorageClass is used.
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
}

// SecretReference references a Secret by name. Used when the Secret structure
// is fixed by convention (e.g. "username" and "password" keys for OCI credentials).
type SecretReference struct {
	// Name is the name of the Secret in the same namespace as the ModelArtifact.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`
}

// SecretKeySelector references a specific key within a Secret.
type SecretKeySelector struct {
	// Name is the name of the Secret in the same namespace as the ModelArtifact.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`

	// Key is the key within the Secret data. If omitted, defaults to "token".
	// +optional
	Key string `json:"key,omitempty"`
}

// ResourceReference is a lightweight reference to a namespaced Kubernetes resource.
type ResourceReference struct {
	// Name is the name of the referenced resource.
	// +required
	Name string `json:"name"`

	// Namespace is the namespace of the referenced resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// ModelArtifactStatus defines the observed state of a ModelArtifact.
type ModelArtifactStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase is the high-level summary of the artifact lifecycle.
	// +optional
	Phase ArtifactPhase `json:"phase,omitempty"`

	// Digest is the OCI manifest digest of the pushed artifact (e.g. "sha256:abc123...").
	// Only populated when Phase is Succeeded.
	// +optional
	Digest string `json:"digest,omitempty"`

	// Reference is the full OCI reference of the pushed artifact including tag.
	// +optional
	Reference string `json:"reference,omitempty"`

	// JobRef references the most recently created Job for this artifact.
	// +optional
	JobRef *ResourceReference `json:"jobRef,omitempty"`

	// StartTime is the timestamp when the most recent job was created.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the timestamp when the most recent job completed (succeeded or failed).
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Conditions store the status conditions of the ModelArtifact.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Repository",type=string,JSONPath=`.spec.target.repository`
// +kubebuilder:printcolumn:name="Digest",type=string,JSONPath=`.status.digest`,priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ModelArtifact is the Schema for the modelartifacts API.
// A ModelArtifact declares intent to import a model from a source (e.g. HuggingFace),
// package it as an OCI artifact (ModelPack or ModelKit format), and push it to an
// OCI-compliant registry. The controller creates a Kubernetes Job to execute the
// import/pack/push pipeline and reports the resulting digest back to the status.
type ModelArtifact struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired model artifact.
	// +required
	Spec ModelArtifactSpec `json:"spec"`

	// Status represents the current state of the model artifact pipeline.
	// +optional
	Status ModelArtifactStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ModelArtifactList contains a list of ModelArtifact resources.
type ModelArtifactList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ModelArtifact `json:"items"`
}
