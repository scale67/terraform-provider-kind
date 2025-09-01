# Terraform Kind Provider

The [Scale67 Kind provider](https://registry.terraform.io/providers/scale67/kind/latest) allows you to manage Kubernetes Kind (Kubernetes in Docker) clusters using Terraform. This provider is perfect for local development, testing, and CI/CD pipelines where you need lightweight, ephemeral Kubernetes clusters.

## Features

- **Create and manage Kind clusters** - Spin up and tear down Kubernetes clusters with a single resource
- **Support for custom Kind configuration files** - Use external `kind-config.yaml` files for complex cluster setups
- **Inline Kind configuration support** - Define cluster configuration directly in your Terraform code
- **Automatic cluster cleanup** - Clusters are automatically destroyed when the Terraform resource is removed
- **Kubeconfig and endpoint outputs** - Get immediate access to cluster connection details
- **Support for custom node images** - Use specific Kubernetes versions for your clusters
- **Wait for cluster readiness** - Ensure clusters are fully ready before proceeding

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) >= 0.20.0
- [Docker](https://docs.docker.com/get-docker/) or Podman

## Installation

### Using the Provider

```hcl
terraform {
  required_providers {
    kind = {
      source = "scale67/kind"
      version = "~> 0.0.1"
    }
  }
}

provider "kind" {}
```

### Building from Source

1. Clone the repository
2. Build and install the provider:

```bash
make install
```

## Quick Start

### Basic Example

```hcl
resource "kind_cluster" "example" {
  name           = "my-cluster"
  wait_for_ready = true

  kind_config = <<YAML
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
YAML
}

output "kubeconfig" {
  value = kind_cluster.example.kubeconfig
  sensitive = true
}
```

### Using External Configuration File

```hcl
resource "kind_cluster" "example" {
  name           = "my-cluster"
  config_path    = "${path.module}/kind-config.yaml"
  wait_for_ready = true
  node_image     = "kindest/node:v1.28.0"
}
```

## Use Cases

- **Local Development**: Create isolated Kubernetes environments for testing applications
- **CI/CD Pipelines**: Spin up ephemeral clusters for automated testing
- **Learning Kubernetes**: Experiment with Kubernetes features without cloud costs
- **Multi-node Testing**: Test applications across different cluster topologies
- **Integration Testing**: Validate Helm charts and Kubernetes manifests

## Examples

Check the `examples/` directory for comprehensive examples:

- `examples/basic/` - Basic cluster with inline configuration
- `examples/with-config-file/` - Cluster using external configuration file
- `examples/multi-node/` - Multi-node cluster with custom networking
