# Examples

This directory contains various examples for using the Terraform Kind provider.

## Basic Example

The `basic/` directory contains a simple example that creates a single-node Kind cluster with inline configuration.

```bash
cd basic
terraform init
terraform plan
terraform apply
```

## With Config File Example

The `with-config-file/` directory demonstrates how to use an external `kind-config.yaml` file to configure your cluster.

```bash
cd with-config-file
terraform init
terraform plan
terraform apply
```

## Multi-Node Example

The `multi-node/` directory shows how to create a multi-node cluster with custom networking configuration and multiple clusters.

```bash
cd multi-node
terraform init
terraform plan
terraform apply
```

## Common Configuration Options

### Port Mappings

To expose services running in your Kind cluster to your host machine:

```yaml
extraPortMappings:
- containerPort: 80
  hostPort: 80
  protocol: TCP
- containerPort: 443
  hostPort: 443
  protocol: TCP
```

### Node Labels

To add labels to nodes for scheduling purposes:

```yaml
kubeadmConfigPatches:
- |
  kind: InitConfiguration
  nodeRegistration:
    kubeletExtraArgs:
      node-labels: "ingress-ready=true"
```

### Custom Node Images

To use a specific Kubernetes version:

```hcl
resource "kind_cluster" "example" {
  name       = "my-cluster"
  node_image = "kindest/node:v1.28.0"
  # ... other configuration
}
```

## Using kubectl with Your Cluster

After creating a cluster with Terraform, you can use kubectl to interact with it. The provider automatically configures kubectl for you.

### Troubleshooting kubectl Access

If you encounter issues with kubectl commands after creating a cluster, the most common causes are:

1. **TLS Certificate Issues**: Avoid mapping port 6443 (Kubernetes API server) to `0.0.0.0` in your kind config, as this can cause TLS certificate validation errors.

2. **Name Mismatch**: Ensure the cluster name in your Terraform configuration matches the name in your `kind-config.yaml` file.

### Helper Script

The `with-config-file/` example includes a helper script `use-cluster.sh` that you can run to easily set up kubectl access:

```bash
cd with-config-file
./use-cluster.sh
```

This script:
- Extracts the kubeconfig from Terraform output
- Sets the `KUBECONFIG` environment variable
- Provides instructions for permanent setup

### Manual Setup

Alternatively, you can manually set up kubectl access:

```bash
# Get the kubeconfig content
terraform output -raw kubeconfig > kubeconfig.yaml

# Use the kubeconfig
export KUBECONFIG=./kubeconfig.yaml

# Or merge with your existing config
kubectl config view --flatten > ~/.kube/config
```

## Cleanup

To destroy the clusters created by these examples:

```bash
terraform destroy
```

Or manually using Kind:

```bash
kind delete cluster --name <cluster-name>
```