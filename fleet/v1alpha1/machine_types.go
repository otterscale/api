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

// MachineRole defines the role a Machine plays in the cluster.
// +kubebuilder:validation:Enum=controlplane;worker
// +enum
type MachineRole string

const (
	MachineRoleControlPlane MachineRole = "controlplane"
	MachineRoleWorker       MachineRole = "worker"
)

// MachinePhase represents the lifecycle phase of a Machine.
// +kubebuilder:validation:Enum=Pending;Provisioning;Provisioned;Bootstrapping;Running;Failed;Deleting
// +enum
type MachinePhase string

const (
	MachinePhasePending       MachinePhase = "Pending"
	MachinePhaseProvisioning  MachinePhase = "Provisioning"
	MachinePhaseProvisioned   MachinePhase = "Provisioned"
	MachinePhaseBootstrapping MachinePhase = "Bootstrapping"
	MachinePhaseRunning       MachinePhase = "Running"
	MachinePhaseFailed        MachinePhase = "Failed"
	MachinePhaseDeleting      MachinePhase = "Deleting"
)

// BareMetalHostReference identifies a Metal3 BareMetalHost resource.
type BareMetalHostReference struct {
	// Name is the name of the BareMetalHost.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`

	// Namespace is the namespace of the BareMetalHost.
	// +kubebuilder:validation:MinLength=1
	// +required
	Namespace string `json:"namespace"`
}

// MachineAddress contains information about a machine's address.
type MachineAddress struct {
	// Type is the type of the address (InternalIP, ExternalIP, Hostname).
	// +required
	Type string `json:"type"`

	// Address is the actual address value.
	// +required
	Address string `json:"address"`
}

// MachineSpec defines the desired state of a Machine.
// +kubebuilder:validation:XValidation:rule="!(self.role == 'worker' && has(self.bootstrap) && self.bootstrap)",message="worker nodes cannot be bootstrap nodes"
type MachineSpec struct {
	// ClusterRef is the name of the Cluster this Machine belongs to.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="clusterRef is immutable"
	// +required
	ClusterRef string `json:"clusterRef"`

	// Role defines whether this Machine is a control plane or worker node.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="role is immutable"
	// +required
	Role MachineRole `json:"role"`

	// BareMetalHostRef references the Metal3 BareMetalHost that backs this Machine.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="bareMetalHostRef is immutable"
	// +required
	BareMetalHostRef BareMetalHostReference `json:"bareMetalHostRef"`

	// TalosConfig overrides the Cluster-level default config for this Machine.
	// When omitted, the fallback is role-aware: control plane Machines use
	// the Cluster's ControlPlaneConfig, worker Machines use the Cluster's WorkerConfig.
	// +optional
	TalosConfig *TalosConfigSpec `json:"talosConfig,omitempty"`

	// Bootstrap indicates whether this Machine is the initial bootstrap node
	// responsible for initializing etcd. Exactly one control plane Machine per
	// Cluster should have this set to true. Must not be set on worker Machines.
	// +optional
	Bootstrap bool `json:"bootstrap,omitempty"`
}

// MachineStatus defines the observed state of a Machine.
type MachineStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase is the high-level summary of the machine lifecycle.
	// +optional
	Phase MachinePhase `json:"phase,omitempty"`

	// Ready is true when the Machine is fully bootstrapped and the node is healthy.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// InfrastructureReady is true when the BareMetalHost has finished provisioning.
	// +optional
	InfrastructureReady bool `json:"infrastructureReady,omitempty"`

	// BootstrapReady is true when the Talos bootstrap has completed.
	// +optional
	BootstrapReady bool `json:"bootstrapReady,omitempty"`

	// BootstrapDataSecretRef references the Secret containing the Talos machine config.
	// +optional
	BootstrapDataSecretRef *ResourceReference `json:"bootstrapDataSecretRef,omitempty"`

	// Addresses contains the addresses reported by the BareMetalHost hardware inspection.
	// +optional
	Addresses []MachineAddress `json:"addresses,omitempty"`

	// Conditions store the status conditions of the Machine.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Cluster",type=string,JSONPath=`.spec.clusterRef`
// +kubebuilder:printcolumn:name="Role",type=string,JSONPath=`.spec.role`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Ready",type=boolean,JSONPath=`.status.ready`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Machine is the Schema for the machines API.
// A Machine represents a single bare metal node in a Talos cluster,
// backed by a Metal3 BareMetalHost.
type Machine struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired behavior of the Machine.
	// +required
	Spec MachineSpec `json:"spec"`

	// Status represents the current information about the Machine.
	// +optional
	Status MachineStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// MachineList contains a list of Machine resources.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Machine `json:"items"`
}
