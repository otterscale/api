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

// AcceleratorType specifies the accelerator hardware type.
// +enum
type AcceleratorType string

const (
	AcceleratorNvidia     AcceleratorType = "nvidia"
	AcceleratorAMD        AcceleratorType = "amd"
	AcceleratorIntelGaudi AcceleratorType = "intel-gaudi"
	AcceleratorGoogle     AcceleratorType = "google"
	AcceleratorCPU        AcceleratorType = "cpu"
)

// ModelServicePhase represents the current lifecycle phase of a ModelService.
// +enum
type ModelServicePhase string

const (
	// ModelServicePending indicates the serving resources are being created.
	ModelServicePending ModelServicePhase = "Pending"

	// ModelServiceRunning indicates the serving pods are running but not yet fully ready.
	ModelServiceRunning ModelServicePhase = "Running"

	// ModelServiceReady indicates all desired replicas are ready and serving traffic.
	ModelServiceReady ModelServicePhase = "Ready"

	// ModelServiceFailed indicates one or more components have failed.
	ModelServiceFailed ModelServicePhase = "Failed"
)

// EPPFailureMode determines behavior when the Endpoint Picker is unavailable.
// +enum
type EPPFailureMode string

const (
	EPPFailureModeFail EPPFailureMode = "FailClose"
	EPPFailureModeOpen EPPFailureMode = "FailOpen"
)

// ModelServiceSpec defines the desired state of a ModelService.
// It declares how to serve an OCI-packaged model using vLLM with optional
// Prefill/Decode disaggregation and Gateway API Inference Extension integration.
type ModelServiceSpec struct {
	// Model defines the OCI model artifact and serving identity.
	// +required
	Model ModelSpec `json:"model"`

	// Engine configures the inference engine (vLLM).
	// +required
	Engine EngineSpec `json:"engine"`

	// Accelerator configures the GPU/accelerator type for serving pods.
	// +required
	Accelerator AcceleratorSpec `json:"accelerator"`

	// Decode configures the decode (or unified) serving pods.
	// In non-disaggregated mode, these are the only serving pods.
	// +required
	Decode RoleSpec `json:"decode"`

	// Prefill optionally configures separate prefill pods for disaggregated serving.
	// When set, the serving architecture splits into Prefill (prompt processing)
	// and Decode (token generation) phases on separate pod groups.
	// +optional
	Prefill *RoleSpec `json:"prefill,omitempty"`

	// RoutingProxy configures the llm-d routing sidecar for disaggregated serving.
	// Required when Prefill is set; the proxy routes prefill requests between pods.
	// +optional
	RoutingProxy *RoutingProxySpec `json:"routingProxy,omitempty"`

	// InferencePool configures the Gateway API Inference Extension InferencePool.
	// When set, the operator creates and manages an InferencePool resource
	// with selector labels matching the serving pods.
	// +optional
	InferencePool *InferencePoolSpec `json:"inferencePool,omitempty"`

	// HTTPRoute optionally creates a Gateway API HTTPRoute pointing to the InferencePool.
	// Requires InferencePool to be set.
	// +optional
	HTTPRoute *HTTPRouteSpec `json:"httpRoute,omitempty"`

	// Monitoring configures observability features.
	// +optional
	Monitoring *MonitoringSpec `json:"monitoring,omitempty"`
}

// ModelSpec defines the OCI model artifact to serve.
// The model is mounted as a read-only Kubernetes image volume (requires K8s >= 1.35).
type ModelSpec struct {
	// Name is the model identifier used in OpenAI-compatible API requests
	// (e.g. "qwen/Qwen3-32B", "meta-llama/Llama-3-70B-Instruct").
	// This is passed to vLLM as --served-model-name.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +required
	Name string `json:"name"`

	// Image is the OCI reference for the model artifact.
	// It is used as a Kubernetes image volume source, leveraging node-level
	// container image caching for efficient model distribution.
	// Example: "registry.example.com/models/qwen3-32b:v1"
	// +kubebuilder:validation:MinLength=1
	// +required
	Image string `json:"image"`

	// MountPath is where the model artifact is mounted in serving containers.
	// The vLLM --model argument is automatically set to this path.
	// +kubebuilder:default="/models"
	// +optional
	MountPath string `json:"mountPath,omitempty"`

	// ImagePullSecrets for pulling the model OCI artifact from a private registry.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// EngineSpec configures the vLLM inference engine.
type EngineSpec struct {
	// Image is the vLLM container image.
	// +optional
	Image string `json:"image,omitempty"`

	// ImagePullPolicy for the engine container.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Args are additional vLLM command-line arguments.
	// The operator automatically sets --model, --port, --served-model-name,
	// --tensor-parallel-size, and --data-parallel-size based on the spec.
	// User-provided args are appended after the auto-generated ones.
	// +optional
	Args []string `json:"args,omitempty"`

	// Env are additional environment variables for the engine container.
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Port is the port the inference engine listens on.
	// When a routing proxy is enabled, this is the external port exposed by the proxy,
	// and vLLM listens on routingProxy.targetPort instead.
	// +kubebuilder:default=8000
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port int32 `json:"port,omitempty"`
}

// AcceleratorSpec configures the accelerator hardware for serving pods.
type AcceleratorSpec struct {
	// Type of accelerator hardware.
	// The operator uses this to determine the appropriate resource name
	// (e.g. nvidia.com/gpu, amd.com/gpu) and any accelerator-specific
	// environment variables.
	// +kubebuilder:validation:Enum=nvidia;amd;intel-gaudi;google;cpu
	// +required
	Type AcceleratorType `json:"type"`
}

// ParallelismSpec configures vLLM tensor and data parallelism.
type ParallelismSpec struct {
	// Tensor is the tensor-parallel-size: number of GPUs used to shard a single model.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +optional
	Tensor int32 `json:"tensor,omitempty"`

	// Data is the data-parallel-size: number of data-parallel replicas within a single pod.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +optional
	Data int32 `json:"data,omitempty"`

	// DataLocal is the data-parallel-size-local for disaggregated serving.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +optional
	DataLocal int32 `json:"dataLocal,omitempty"`
}

// RoleSpec configures a group of serving pods (decode or prefill).
type RoleSpec struct {
	// Replicas is the desired number of pod replicas.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Parallelism configures vLLM tensor/data parallelism for this role.
	// +optional
	Parallelism ParallelismSpec `json:"parallelism,omitempty"`

	// Resources for the vLLM container (CPU, memory).
	// GPU resources are automatically calculated from accelerator type
	// and parallelism (tensor * dataLocal) and should not be set manually.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// NodeSelector constrains pods to nodes with matching labels.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations for the serving pods.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Annotations for the serving pods.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// RoutingProxySpec configures the llm-d routing sidecar.
// The proxy runs as a native sidecar (restartable init container) that intercepts
// incoming requests and handles prefill/decode routing for disaggregated serving.
type RoutingProxySpec struct {
	// Image is the routing proxy container image.
	// +optional
	Image string `json:"image,omitempty"`

	// Connector specifies the KV-cache transfer protocol.
	// +kubebuilder:default="nixlv2"
	// +optional
	Connector string `json:"connector,omitempty"`

	// TargetPort is the port where vLLM actually listens when the proxy is enabled.
	// The proxy intercepts on engine.port and forwards to this port.
	// +kubebuilder:default=8200
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	TargetPort int32 `json:"targetPort,omitempty"`

	// ZapEncoder sets the Zap log encoding format (e.g. "json", "console").
	// +optional
	ZapEncoder string `json:"zapEncoder,omitempty"`

	// ZapLogLevel sets the Zap log level (e.g. "debug", "info", "error").
	// +optional
	ZapLogLevel string `json:"zapLogLevel,omitempty"`

	// SecureProxy enables TLS on the routing proxy.
	// +optional
	SecureProxy *bool `json:"secureProxy,omitempty"`

	// PrefillerUseTLS enables TLS for prefiller communication.
	// +optional
	PrefillerUseTLS *bool `json:"prefillerUseTLS,omitempty"`

	// CertPath is the path to TLS certificates for the routing proxy.
	// +optional
	CertPath string `json:"certPath,omitempty"`
}

// InferencePoolSpec configures the Gateway API Inference Extension InferencePool.
// The operator creates an InferencePool with selector labels matching the serving
// pods and target port matching the engine port. It also deploys and manages
// the Endpoint Picker (EPP) infrastructure: Deployment, Service, ConfigMap,
// ServiceAccount, RBAC, and optionally Istio DestinationRule.
type InferencePoolSpec struct {
	// EndpointPicker configures the Endpoint Picker extension deployment.
	// The operator creates and manages the EPP Deployment, Service, and
	// supporting resources alongside the InferencePool.
	// +required
	EndpointPicker EndpointPickerSpec `json:"endpointPicker"`
}

// EndpointPickerSpec configures the Endpoint Picker (EPP) extension deployment.
// The EPP service name is automatically derived from the ModelService name.
type EndpointPickerSpec struct {
	// Image is the EPP container image.
	// +optional
	Image string `json:"image,omitempty"`

	// Replicas is the number of EPP pod replicas.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Resources for the EPP container (CPU, memory).
	// CPU limits are intentionally left unset by default to allow bursting
	// during scheduling spikes.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Port of the EPP service (extProc gRPC port).
	// +kubebuilder:default=9002
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port int32 `json:"port,omitempty"`

	// FailureMode determines behavior when the EPP is unavailable.
	// +kubebuilder:validation:Enum=FailOpen;FailClose
	// +kubebuilder:default="FailOpen"
	// +optional
	FailureMode EPPFailureMode `json:"failureMode,omitempty"`
}

// HTTPRouteSpec configures a Gateway API HTTPRoute that routes traffic
// from a Gateway to the InferencePool backend.
type HTTPRouteSpec struct {
	// GatewayRef references the Gateway to attach the HTTPRoute to.
	// +required
	GatewayRef GatewayRef `json:"gatewayRef"`

	// Hostnames for the HTTPRoute. If empty, the route matches all hostnames.
	// +optional
	Hostnames []string `json:"hostnames,omitempty"`
}

// GatewayRef references a Gateway API Gateway.
type GatewayRef struct {
	// Name of the Gateway.
	// +kubebuilder:validation:MinLength=1
	// +required
	Name string `json:"name"`

	// Namespace of the Gateway. If empty, defaults to the ModelService namespace.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// MonitoringSpec configures observability features.
type MonitoringSpec struct {
	// PodMonitor configures Prometheus PodMonitor creation for serving pods.
	// +optional
	PodMonitor *PodMonitorSpec `json:"podMonitor,omitempty"`
}

// PodMonitorSpec configures Prometheus PodMonitor resources.
type PodMonitorSpec struct {
	// Enabled controls PodMonitor creation.
	// +kubebuilder:default=false
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// PortName to scrape metrics from. Must match a named port on the vLLM container.
	// +kubebuilder:default="http"
	// +optional
	PortName string `json:"portName,omitempty"`

	// Path is the HTTP endpoint to scrape metrics from.
	// +kubebuilder:default="/metrics"
	// +optional
	Path string `json:"path,omitempty"`

	// Interval between scrapes.
	// +kubebuilder:default="30s"
	// +optional
	Interval string `json:"interval,omitempty"`
}

// ModelServiceStatus defines the observed state of a ModelService.
type ModelServiceStatus struct {
	// ObservedGeneration is the most recent generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase is the high-level summary of the ModelService lifecycle.
	// +optional
	Phase ModelServicePhase `json:"phase,omitempty"`

	// DecodeReady is the number of ready decode replicas.
	// +optional
	DecodeReady int32 `json:"decodeReady,omitempty"`

	// DecodeReplicas is the desired number of decode replicas.
	// +optional
	DecodeReplicas int32 `json:"decodeReplicas,omitempty"`

	// PrefillReady is the number of ready prefill replicas.
	// +optional
	PrefillReady int32 `json:"prefillReady,omitempty"`

	// PrefillReplicas is the desired number of prefill replicas.
	// +optional
	PrefillReplicas int32 `json:"prefillReplicas,omitempty"`

	// Conditions store the status conditions of the ModelService.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Namespaced,categories={otterscale}
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Model",type=string,JSONPath=`.spec.model.name`
// +kubebuilder:printcolumn:name="Decode",type=string,JSONPath=`.status.decodeReady`,priority=1
// +kubebuilder:printcolumn:name="Prefill",type=string,JSONPath=`.status.prefillReady`,priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ModelService is the Schema for the modelservices API.
// A ModelService declares intent to serve an OCI-packaged model using vLLM,
// optionally with Prefill/Decode disaggregation and Gateway API integration.
// The model artifact is mounted via Kubernetes image volumes (requires K8s >= 1.35),
// eliminating the need for init containers or PVC-based model loading.
type ModelService struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// Spec defines the desired model serving configuration.
	// +required
	Spec ModelServiceSpec `json:"spec"`

	// Status represents the current state of the model serving deployment.
	// +optional
	Status ModelServiceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ModelServiceList contains a list of ModelService resources.
type ModelServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ModelService `json:"items"`
}
