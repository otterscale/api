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

// ClusterPhase represents the lifecycle phase of a Cluster.
// +kubebuilder:validation:Enum=Pending;Provisioning;Provisioned;Ready;Failed;Deleting
// +enum
type ClusterPhase string

const (
	ClusterPhasePending      ClusterPhase = "Pending"
	ClusterPhaseProvisioning ClusterPhase = "Provisioning"
	ClusterPhaseProvisioned  ClusterPhase = "Provisioned"
	ClusterPhaseReady        ClusterPhase = "Ready"
	ClusterPhaseFailed       ClusterPhase = "Failed"
	ClusterPhaseDeleting     ClusterPhase = "Deleting"
)

// Endpoint defines a host:port pair for the Kubernetes API server.
type Endpoint struct {
	// Host is the hostname or IP address of the API server.
	// +kubebuilder:validation:MinLength=1
	// +required
	Host string `json:"host"`

	// Port is the port number of the API server.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=6443
	// +optional
	Port int32 `json:"port,omitempty"`
}

// ImageSpec defines the OS image to be provisioned on BareMetalHost resources.
type ImageSpec struct {
	// URL is the location of the OS image (e.g., HTTP server hosting Talos raw image).
	// +kubebuilder:validation:MinLength=1
	// +required
	URL string `json:"url"`

	// Checksum is the image checksum or a URL pointing to a checksum file.
	// +kubebuilder:validation:MinLength=1
	// +required
	Checksum string `json:"checksum"`

	// ChecksumType is the algorithm used for the checksum (sha256, sha512, md5, or auto).
	// +kubebuilder:validation:Enum=sha256;sha512;md5;auto
	// +kubebuilder:default="sha256"
	// +optional
	ChecksumType string `json:"checksumType,omitempty"`

	// Format is the disk image format (raw, qcow2, or live-iso).
	// +kubebuilder:validation:Enum=raw;qcow2;live-iso
	// +kubebuilder:default="raw"
	// +optional
	Format string `json:"format,omitempty"`
}

// ClusterNetworkSpec configures cluster-level networking.
type ClusterNetworkSpec struct {
	// DNSDomain is the DNS domain used for internal cluster DNS (default: cluster.local).
	// +kubebuilder:default="cluster.local"
	// +optional
	DNSDomain string `json:"dnsDomain,omitempty"`

	// PodSubnets is the list of CIDR ranges for pod IP allocation.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=4
	// +kubebuilder:validation:items:MaxLength=43
	// +kubebuilder:validation:XValidation:rule="self.all(cidr, isCIDR(cidr))",message="each entry must be a valid CIDR notation"
	// +optional
	PodSubnets []string `json:"podSubnets,omitempty"`

	// ServiceSubnets is the list of CIDR ranges for service ClusterIP allocation.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=4
	// +kubebuilder:validation:items:MaxLength=43
	// +kubebuilder:validation:XValidation:rule="self.all(cidr, isCIDR(cidr))",message="each entry must be a valid CIDR notation"
	// +optional
	ServiceSubnets []string `json:"serviceSubnets,omitempty"`
}

// TalosConfigSpec defines how Talos machine configuration is produced.
type TalosConfigSpec struct {
	// GenerateType determines how the Talos machine config is produced.
	// "controlplane" auto-generates a control plane config from cluster parameters.
	// "worker" auto-generates a worker (join) config from cluster parameters.
	// "none" expects the user to supply a complete config in the Data field.
	// When omitted, the controller defaults based on the Machine role.
	// +kubebuilder:validation:Enum=controlplane;worker;none
	// +optional
	GenerateType string `json:"generateType,omitempty"`

	// Data is the raw Talos machine configuration YAML.
	// Required when GenerateType is "none"; ignored otherwise.
	// +optional
	Data string `json:"data,omitempty"`

	// ConfigPatches is a list of RFC 6902 JSON patches applied to the generated config.
	// Ignored when GenerateType is "none".
	// +optional
	ConfigPatches []ConfigPatch `json:"configPatches,omitempty"`
}

// ConfigPatch defines a single RFC 6902 JSON patch operation.
type ConfigPatch struct {
	// Op is the patch operation (add, remove, replace, move, copy, test).
	// +kubebuilder:validation:Enum=add;remove;replace;move;copy;test
	// +required
	Op string `json:"op"`

	// Path is the JSON pointer to the target location.
	// +kubebuilder:validation:MinLength=1
	// +required
	Path string `json:"path"`

	// Value is the value to use in the patch operation (required for add/replace/test).
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Value *runtime.RawExtension `json:"value,omitempty"`
}

// ClusterSpec defines the desired state of a Cluster.
type ClusterSpec struct {
	// ControlPlaneEndpoint is the API server endpoint (typically a VIP or load balancer).
	// +required
	ControlPlaneEndpoint Endpoint `json:"controlPlaneEndpoint"`

	// TalosVersion is the Talos OS version to use for config generation (e.g., "v1.9").
	// +kubebuilder:validation:Pattern=`^v\d+\.\d+(\.\d+)?$`
	// +required
	TalosVersion string `json:"talosVersion"`

	// KubernetesVersion is the target Kubernetes version (e.g., "v1.32.0").
	// +kubebuilder:validation:Pattern=`^v\d+\.\d+\.\d+$`
	// +required
	KubernetesVersion string `json:"kubernetesVersion"`

	// TalosImage specifies the Talos OS image written to each BareMetalHost.
	// +required
	TalosImage ImageSpec `json:"talosImage"`

	// ClusterNetwork configures pod, service, and DNS networking.
	// +optional
	ClusterNetwork *ClusterNetworkSpec `json:"clusterNetwork,omitempty"`

	// ControlPlaneConfig defines how Talos machine configs are produced for control plane nodes.
	// Individual Machines may override this via their own TalosConfig field.
	// +optional
	ControlPlaneConfig TalosConfigSpec `json:"controlPlaneConfig,omitzero"`

	// WorkerConfig defines how Talos machine configs are produced for worker nodes.
	// When omitted, defaults to generateType "worker" which auto-generates a join config.
	// Individual Machines may override this via their own TalosConfig field.
	// +optional
	WorkerConfig TalosConfigSpec `json:"workerConfig,omitzero"`
}

// ResourceReference is a lightweight reference to a Kubernetes resource managed by the operator.
type ResourceReference struct {
	// Name is the name of the referenced resource.
	// +required
	Name string `json:"name"`

	// Namespace is the namespace of the referenced resource.
	// Empty for cluster-scoped resources.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// ClusterStatus defines the observed state of a Cluster.
type ClusterStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase is the high-level summary of the cluster lifecycle.
	// +optional
	Phase ClusterPhase `json:"phase,omitempty"`

	// ControlPlaneReady is true when all control plane Machines are bootstrapped and running.
	// +optional
	ControlPlaneReady bool `json:"controlPlaneReady,omitempty"`

	// Initialized is true once the first control plane node has been bootstrapped
	// and etcd has been initialized.
	// +optional
	Initialized bool `json:"initialized,omitempty"`

	// SecretsRef references the Secret containing the Talos secrets bundle.
	// +optional
	SecretsRef *ResourceReference `json:"secretsRef,omitempty"`

	// TalosconfigRef references the Secret containing the talosconfig client configuration.
	// +optional
	TalosconfigRef *ResourceReference `json:"talosconfigRef,omitempty"`

	// ReadyWorkers is the count of worker Machines that are ready.
	// +optional
	ReadyWorkers int32 `json:"readyWorkers,omitempty"`

	// TotalWorkers is the total count of worker Machines.
	// +optional
	TotalWorkers int32 `json:"totalWorkers,omitempty"`

	// Conditions store the status conditions of the Cluster.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster,categories={otterscale}
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="CP Ready",type=boolean,JSONPath=`.status.controlPlaneReady`
// +kubebuilder:printcolumn:name="Workers",type=string,JSONPath=`.status.readyWorkers`,priority=1
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Cluster is the Schema for the clusters API.
// A Cluster represents a Talos Linux Kubernetes cluster provisioned on bare metal
// via Metal3 BareMetalHost resources.
type Cluster struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired behavior of the Cluster.
	// +required
	Spec ClusterSpec `json:"spec"`

	// Status represents the current information about the Cluster.
	// +optional
	Status ClusterStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Cluster resources.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Cluster `json:"items"`
}
