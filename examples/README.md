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

## Cleanup

To destroy the clusters created by these examples:

```bash
terraform destroy
```

Or manually using Kind:

```bash
kind delete cluster --name <cluster-name>
```