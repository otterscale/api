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
)

// KustomizationTemplate defines how to deploy a module via Kustomize.
// The operator clones the source Git repository, builds the kustomization,
// and applies the resulting manifests using server-side apply.
type KustomizationTemplate struct {
	// URL is the Git repository URL containing the kustomization.
	// Supports HTTPS and SSH URLs.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^(https?://|git@|ssh://).+$`
	// +required
	URL string `json:"url"`

	// Ref specifies the Git reference to check out.
	// If not set, the default branch (usually main) is used.
	// +optional
	Ref *GitReference `json:"ref,omitempty"`

	// Path is the directory path within the repository where
	// kustomization.yaml is located. Defaults to the repository root.
	// +optional
	Path string `json:"path,omitempty"`

	// Interval at which the operator re-reconciles this kustomization.
	// +required
	Interval metav1.Duration `json:"interval"`

	// Prune enables garbage collection: resources that were previously
	// applied but are no longer present in the kustomization output
	// will be deleted from the cluster.
	// +optional
	Prune bool `json:"prune,omitempty"`

	// Force instructs the operator to recreate resources that have
	// immutable field changes, instead of failing the apply.
	// +optional
	Force bool `json:"force,omitempty"`

	// Timeout is the maximum duration for the build and apply operation.
	// Defaults to 5m if not specified.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// TargetNamespace overrides the namespace for all resources in the
	// kustomization output. If empty, each resource keeps its own namespace.
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?)$`
	// +optional
	TargetNamespace string `json:"targetNamespace,omitempty"`

	// Patches is a list of strategic merge or JSON6902 patches to apply
	// on top of the kustomization output before sending to the cluster.
	// +optional
	Patches []KustomizePatch `json:"patches,omitempty"`

	// SecretRef references a Secret containing credentials for the Git
	// repository. Supported keys: username + password (HTTPS), identity +
	// identity.pub + known_hosts (SSH).
	// +optional
	SecretRef *SecretReference `json:"secretRef,omitempty"`
}

// GitReference specifies a Git reference for source checkout.
// At most one field should be set; if none are set the default branch is used.
// +kubebuilder:validation:XValidation:rule="[has(self.branch) && self.branch != '', has(self.tag) && self.tag != '', has(self.commit) && self.commit != '', has(self.semver) && self.semver != ''].filter(x, x).size() <= 1",message="at most one of branch, tag, commit, or semver may be set"
type GitReference struct {
	// Branch is the Git branch to check out.
	// +optional
	Branch string `json:"branch,omitempty"`

	// Tag is the Git tag to check out.
	// +optional
	Tag string `json:"tag,omitempty"`

	// Commit is the Git commit SHA to check out.
	// +kubebuilder:validation:Pattern=`^[a-f0-9]{7,40}$`
	// +optional
	Commit string `json:"commit,omitempty"`

	// Semver is a semver range expression used to select the latest
	// matching Git tag.
	// +optional
	Semver string `json:"semver,omitempty"`
}

// KustomizePatch defines an inline strategic merge or JSON6902 patch.
type KustomizePatch struct {
	// Patch is the inline YAML patch content.
	// +kubebuilder:validation:MinLength=1
	// +required
	Patch string `json:"patch"`

	// Target selects which resources to apply the patch to.
	// If not set, the patch is applied to all matching resources.
	// +optional
	Target *PatchSelector `json:"target,omitempty"`
}

// PatchSelector selects Kubernetes resources by GVK, name, and namespace.
type PatchSelector struct {
	// Group is the API group of the target resource.
	// +optional
	Group string `json:"group,omitempty"`

	// Version is the API version of the target resource.
	// +optional
	Version string `json:"version,omitempty"`

	// Kind is the kind of the target resource.
	// +optional
	Kind string `json:"kind,omitempty"`

	// Namespace of the target resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name of the target resource.
	// +optional
	Name string `json:"name,omitempty"`

	// AnnotationSelector filters resources by annotations
	// (e.g. "config.kubernetes.io/managed-by=kustomize").
	// +optional
	AnnotationSelector string `json:"annotationSelector,omitempty"`

	// LabelSelector filters resources by labels (e.g. "app=nginx").
	// +optional
	LabelSelector string `json:"labelSelector,omitempty"`
}

// KustomizationStatus captures the observed state of a Kustomize-based
// module managed by the operator.
type KustomizationStatus struct {
	// LastAppliedRevision is the Git commit SHA that was last successfully
	// built and applied to the cluster.
	// +optional
	LastAppliedRevision string `json:"lastAppliedRevision,omitempty"`

	// LastAttemptedRevision is the Git commit SHA that was last attempted
	// (may differ from LastAppliedRevision on failure).
	// +optional
	LastAttemptedRevision string `json:"lastAttemptedRevision,omitempty"`
}

// SecretReference holds a reference to a Kubernetes Secret by name.
// The Secret must exist in the operator's namespace (for cluster-scoped
// Modules) or the resolved target namespace.
type SecretReference struct {
	// Name is the name of the Secret.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`
}

// InventoryEntry records the identity of a single Kubernetes resource
// that was applied as part of a Module reconciliation. Used for
// garbage-collecting resources that are no longer part of the desired state.
type InventoryEntry struct {
	// ID uniquely identifies the resource in the format
	// "<namespace>_<name>_<group>_<kind>".
	// +required
	ID string `json:"id"`

	// Version is the API version of the resource (e.g. "v1", "apps/v1").
	// +required
	Version string `json:"version"`
}
