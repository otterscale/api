# OtterScale API

[![Lint](https://github.com/otterscale/api/actions/workflows/lint.yml/badge.svg)](https://github.com/otterscale/api/actions/workflows/lint.yml)
[![Release](https://img.shields.io/github/v/release/otterscale/api)](https://github.com/otterscale/api/releases/latest)
[![npm](https://img.shields.io/npm/v/@otterscale/api)](https://www.npmjs.com/package/@otterscale/api)
[![Go Reference](https://pkg.go.dev/badge/github.com/otterscale/api.svg)](https://pkg.go.dev/github.com/otterscale/api)
[![License](https://img.shields.io/github/license/otterscale/api)](LICENSE)

Shared type definitions — CRDs (Kubebuilder) + ConnectRPC services (Protobuf) — for the OtterScale multi-cluster Kubernetes platform.

## Quick Start

```bash
# Go
go get github.com/otterscale/api@latest

# TypeScript
npm install @otterscale/api

# Generate all (proto, CRDs, deepcopy, lint)
make all

# Individual targets
make proto       # buf generate (Go + TypeScript + OpenAPI)
make manifests   # controller-gen CRDs
make generate    # controller-gen deepcopy
make lint        # golangci-lint
```

## API Groups

### CRDs (Kubebuilder)

| Group                             | Kind             | Scope      | Purpose                                                 |
| --------------------------------- | ---------------- | ---------- | ------------------------------------------------------- |
| `fleet.otterscale.io/v1alpha1`    | `Cluster`        | Cluster    | Talos bare metal Kubernetes cluster                     |
| `fleet.otterscale.io/v1alpha1`    | `Machine`        | Cluster    | Single bare metal node in a Talos cluster               |
| `model.otterscale.io/v1alpha1`    | `ModelArtifact`  | Namespaced | Import, package, and push model to OCI registry         |
| `model.otterscale.io/v1alpha1`    | `ModelService`   | Namespaced | Serve OCI-packaged model with llm-d (Prefill/Decode)    |
| `module.otterscale.io/v1alpha1`   | `Module`         | Cluster    | Installed platform module from a template               |
| `module.otterscale.io/v1alpha1`   | `ModuleTemplate` | Cluster    | Reusable module blueprint (Helm chart / Kustomization)  |
| `tenant.otterscale.io/v1alpha1`   | `Workspace`      | Cluster    | Namespace isolation with RBAC, quotas, network policies |
| `workload.otterscale.io/v1alpha1` | `Application`    | Namespaced | Unified Deployment + Service + PVC abstraction          |

### Install CRDs

```bash
kubectl apply --server-side -f https://github.com/otterscale/api/releases/download/<version>/crds.yaml
```

> **Note:** CRDs that embed Kubernetes core types (e.g. `DeploymentSpec`) produce large schemas that exceed the `kubectl.kubernetes.io/last-applied-configuration` annotation limit. Use `--server-side` to avoid this.

### ConnectRPC Services (Protobuf)

| Package                  | Service           | Key RPCs                                                                   |
| ------------------------ | ----------------- | -------------------------------------------------------------------------- |
| `otterscale.link.v1`     | `LinkService`     | `Register`, `ListLinks`, `GetAgentManifest`                                |
| `otterscale.resource.v1` | `ResourceService` | `Discovery`, `Schema`, `List`, `Get`, `Create`, `Apply`, `Delete`, `Watch` |
| `otterscale.runtime.v1`  | `RuntimeService`  | `PodLog`, `ExecuteTTY`, `PortForward`, `Scale`, `Restart`                  |

### Generated Outputs (`make proto`)

| Plugin               | Output                 | Description                                   |
| -------------------- | ---------------------- | --------------------------------------------- |
| `protocolbuffers/go` | `*.pb.go`              | Go protobuf types (opaque API)                |
| `connectrpc/go`      | `*.connect.go`         | Go ConnectRPC clients/handlers                |
| `bufbuild/es`        | `ts/src/`              | TypeScript protobuf types (`@otterscale/api`) |
| `connect-openapi`    | `openapi/openapi.yaml` | OpenAPI spec                                  |

## Toolchain

| Tool                                                                  | Version | Installed via                                       |
| --------------------------------------------------------------------- | ------- | --------------------------------------------------- |
| [buf](https://buf.build)                                              | v1.66.0 | `make proto` (auto-downloads)                       |
| [controller-gen](https://github.com/kubernetes-sigs/controller-tools) | v0.20.1 | `make manifests` / `make generate` (auto-downloads) |
| [golangci-lint](https://golangci-lint.run)                            | v2.10.1 | `make lint` (auto-downloads)                        |

## Feature Gating

RPCs are annotated with `otterscale.api.feature` method options for runtime feature gating:

```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse) {
  option (otterscale.api.feature) = {name: "link-enabled"};
}
```
