terraform {
  required_providers {
    kind = {
      source  = "scale67/kind"
      version = "~> 1.0"
    }
  }
}

provider "kind" {}

# Kind cluster using external configuration file
resource "kind_cluster" "example" {
  name           = "scale-67-config-file-cluster"
  config_path    = "${path.module}/kind-config.yaml"
  wait_for_ready = true
  node_image     = "kindest/node:v1.28.0"
}

# Output the kubeconfig
output "kubeconfig" {
  value     = kind_cluster.example.kubeconfig
  sensitive = true
}

# Output the cluster endpoint
output "endpoint" {
  value = kind_cluster.example.endpoint
}