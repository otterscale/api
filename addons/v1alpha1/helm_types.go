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

// HelmChartTemplate defines how to deploy a module via a Helm chart.
// The operator downloads the chart from the specified repository and
// manages the Helm release lifecycle directly using the Helm SDK.
type HelmChartTemplate struct {
	// RepoURL is the URL of the Helm chart repository.
	// Supports HTTP/HTTPS Helm repositories and OCI registries (oci://).
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^(https?://|oci://).+$`
	// +required
	RepoURL string `json:"repoURL"`

	// Chart is the name of the Helm chart within the repository.
	// +kubebuilder:validation:MinLength=1
	// +required
	Chart string `json:"chart"`

	// Version is the exact chart version to install (e.g. "1.2.3").
	// If empty, the latest version is used.
	// +optional
	Version string `json:"version,omitempty"`

	// Interval at which the operator re-reconciles this Helm release.
	// +required
	Interval metav1.Duration `json:"interval"`

	// Values holds the default Helm chart values as arbitrary JSON.
	// Module.Spec.Values can override these on a per-instance basis.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Values *runtime.RawExtension `json:"values,omitempty"`

	// ReleaseName overrides the Helm release name.
	// Defaults to the Module name if not specified.
	// +kubebuilder:validation:MaxLength=53
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	// +optional
	ReleaseName string `json:"releaseName,omitempty"`

	// CreateNamespace instructs the operator to create the target namespace
	// before installing the chart, if it does not already exist.
	// +optional
	CreateNamespace bool `json:"createNamespace,omitempty"`

	// Timeout is the maximum duration for any single Helm operation.
	// Defaults to 5m if not specified.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// MaxHistory limits the number of Helm release revisions saved.
	// Defaults to 10 if not specified.
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxHistory *int32 `json:"maxHistory,omitempty"`

	// Upgrade configures the Helm upgrade strategy.
	// +optional
	Upgrade *HelmUpgradeStrategy `json:"upgrade,omitempty"`

	// SecretRef references a Secret in the Module's resolved namespace
	// containing credentials for the Helm repository.
	// Supported keys: username, password (Basic Auth), caFile, certFile, keyFile (TLS).
	// +optional
	SecretRef *SecretReference `json:"secretRef,omitempty"`
}

// HelmUpgradeStrategy configures how Helm upgrades are performed.
type HelmUpgradeStrategy struct {
	// Force forces resource updates through a replacement strategy.
	// +optional
	Force bool `json:"force,omitempty"`

	// CleanupOnFail rolls back changes on upgrade failure.
	// +optional
	CleanupOnFail bool `json:"cleanupOnFail,omitempty"`

	// MaxRetries is the maximum number of retries before marking the release as failed.
	// Defaults to 0 (no retries).
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxRetries *int32 `json:"maxRetries,omitempty"`

	// EnableRollback triggers an automatic rollback when an upgrade fails.
	// +optional
	EnableRollback bool `json:"enableRollback,omitempty"`
}

// HelmReleaseStatus captures the observed state of a Helm release
// managed by the operator.
type HelmReleaseStatus struct {
	// ChartVersion is the version of the currently deployed Helm chart.
	// +optional
	ChartVersion string `json:"chartVersion,omitempty"`

	// Revision is the Helm release revision number.
	// +optional
	Revision int `json:"revision,omitempty"`

	// Status is the Helm release status (e.g. deployed, failed, pending-upgrade).
	// +kubebuilder:validation:Enum=unknown;deployed;uninstalled;superseded;failed;uninstalling;pending-install;pending-upgrade;pending-rollback
	// +optional
	Status string `json:"status,omitempty"`

	// ValuesChecksum is a SHA-256 hash of the rendered values used for
	// the most recent install or upgrade, enabling change detection.
	// +optional
	ValuesChecksum string `json:"valuesChecksum,omitempty"`
}
