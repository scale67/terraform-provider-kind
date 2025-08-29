# With Config File Example

This example demonstrates how to use an external `kind-config.yaml` file to configure your Kind cluster with Terraform.

## Quick Start

```bash
# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Apply the configuration
terraform apply -auto-approve
```

## Configuration

The example creates a multi-node cluster with:
- 1 control-plane node with ingress-ready labels
- 2 worker nodes (one with custom labels)
- Port mappings for HTTP (80) and HTTPS (443)

## Using kubectl

After the cluster is created, you can use kubectl commands:

```bash
# Check cluster status
kubectl get nodes

# View all pods
kubectl get pods --all-namespaces

# Check namespaces
kubectl get namespaces
```

## Troubleshooting kubectl Access

If you encounter issues with kubectl commands, the most common cause is **TLS certificate validation errors**. This happens when the Kubernetes API server port (6443) is mapped to `0.0.0.0` instead of using the default localhost binding.

### Solution

The `kind-config.yaml` file has been configured to avoid this issue by:
- Not mapping port 6443 to avoid TLS certificate conflicts
- Using the default kind port binding behavior
- Ensuring cluster names match between Terraform and kind config

### Helper Script

Use the included `use-cluster.sh` script to easily set up kubectl access:

```bash
./use-cluster.sh
```

This script automatically extracts the kubeconfig from Terraform output and sets up your environment.

## Cleanup

```bash
terraform destroy -auto-approve
```

## Files

- `main.tf` - Terraform configuration
- `kind-config.yaml` - Kind cluster configuration
- `use-cluster.sh` - Helper script for kubectl setup
- `README.md` - This file

