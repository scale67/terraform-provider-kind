# Terraform Kind Provider

This Terraform provider allows you to manage Kubernetes Kind (Kubernetes in Docker) clusters using Terraform.

## Features

- Create and manage Kind clusters
- Support for custom Kind configuration files (`kind-config.yaml`)
- Inline Kind configuration support
- Automatic cluster cleanup on destroy
- Kubeconfig and endpoint outputs
- Support for custom node images
- Wait for cluster readiness option

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
      version = "~> 1.0"
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

## Usage

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

## Resource: `kind_cluster`

### Arguments

- `name` (Required, String) - Name of the Kind cluster. Changing this forces a new resource to be created.
- `kind_config` (Optional, String) - Inline Kind configuration YAML. Conflicts with `config_path`.
- `config_path` (Optional, String) - Path to Kind configuration YAML file. Conflicts with `kind_config`.
- `node_image` (Optional, String) - Node Docker image to use for booting the cluster.
- `wait_for_ready` (Optional, Boolean) - Wait for the cluster to be ready before completing.

### Attributes

- `id` (String) - The cluster identifier.
- `kubeconfig` (String, Sensitive) - Kubeconfig for the created cluster.
- `endpoint` (String) - Kubernetes API server endpoint.

## Examples

Check the `examples/` directory for more comprehensive examples:

- `examples/basic/` - Basic cluster with inline configuration
- `examples/with-config-file/` - Cluster using external configuration file
- `examples/multi-node/` - Multi-node cluster with custom networking

## Development

### Building the Provider

```bash
go build -o terraform-provider-kind
```

### Testing

```bash
make test
```

### Installing for Local Development

```bash
make install
```

This will build and install the provider in your local Terraform plugins directory.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the Mozilla Public License v2.0 - see the LICENSE file for details.# gtm-infra-kind-provider
