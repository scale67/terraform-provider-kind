terraform {
  required_providers {
    kind = {
      source = "gtm-cloud-ai/kind"
      version = "~> 1.0"
    }
  }
}

provider "kind" {}

# Basic Kind cluster with inline configuration
resource "kind_cluster" "example" {
  name           = "example-cluster"
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

# Output the kubeconfig
output "kubeconfig" {
  value = kind_cluster.example.kubeconfig
  sensitive = true
}

# Output the cluster endpoint
output "endpoint" {
  value = kind_cluster.example.endpoint
}