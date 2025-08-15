terraform {
  required_providers {
    kind = {
      source = "gtm-cloud-ai/kind"
      version = "~> 1.0"
    }
  }
}

provider "kind" {}

# Multi-node Kind cluster with custom networking
resource "kind_cluster" "multi_node" {
  name           = "multi-node-cluster"
  wait_for_ready = true

  kind_config = <<YAML
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  # WARNING: It is _strongly_ recommended that you keep this the default
  # (127.0.0.1) when using kind for local development. However, if you are
  # running in a remote environment (like a cloud VM), you may need to change
  # this to a publicly accessible address.
  apiServerAddress: "127.0.0.1"
  # By default the API server listens on a random open port.
  # You may choose a specific port but probably don't need to in most cases.
  # Using a random port makes it easier to spin up multiple clusters.
  apiServerPort: 6443
nodes:
# the control plane node config
- role: control-plane
  # configure the control plane to use a specific version
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  # expose ingress controller port mappings
  - containerPort: 80
    hostPort: 80
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    listenAddress: "0.0.0.0"
    protocol: TCP
# the worker nodes
- role: worker
- role: worker
- role: worker
YAML
}

# Separate cluster for testing
resource "kind_cluster" "test" {
  name           = "test-cluster"
  wait_for_ready = false

  kind_config = <<YAML
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
YAML
}

# Output configurations for both clusters
output "multi_node_kubeconfig" {
  value = kind_cluster.multi_node.kubeconfig
  sensitive = true
}

output "multi_node_endpoint" {
  value = kind_cluster.multi_node.endpoint
}

output "test_kubeconfig" {
  value = kind_cluster.test.kubeconfig
  sensitive = true
}

output "test_endpoint" {
  value = kind_cluster.test.endpoint
}