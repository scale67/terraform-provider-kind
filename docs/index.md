🔥 Ready to build something cool with Kubernetes and Terraform? The [Scale67 Kind provider](https://registry.terraform.io/providers/scale67/kind/latest) is your gateway to spinning up lightweight Kubernetes clusters using [Kind (Kubernetes IN Docker)](https://kind.sigs.k8s.io/)—perfect for testing, development, and CI/CD pipelines.

### Why You’d Use the Scale67 Kind Provider

Here’s how it can supercharge your workflow:

- **Local Kubernetes Clusters**: Instantly create and manage Kubernetes clusters inside Docker containers—no cloud account needed.
- **Infrastructure as Code**: Define your cluster setup in Terraform, making it reproducible, version-controlled, and easy to share.
- **CI/CD Friendly**: Ideal for automated testing environments where you need ephemeral clusters that spin up fast and tear down cleanly.
- **Multi-node Support**: Simulate real-world cluster topologies with multiple nodes for more robust testing.
- **Custom Configs**: Inject custom Kind configurations, like networking tweaks or volume mounts, directly from your Terraform code.

### Example Use Case

Imagine you're building a microservices app and want to test how it behaves in a Kubernetes cluster before deploying to production. With this provider, you can:

```hcl
resource "kind_cluster" "dev" {
  name = "dev-cluster"
  config = file("${path.module}/kind-config.yaml")
}
```

Boom—your cluster is up and running locally, ready for Helm charts, manifests, or whatever magic you’re cooking.

---

If you're serious about DevOps or just want to experiment with Kubernetes without the cloud overhead, this provider is a game-changer. Want help writing a full Terraform module with it? Let’s build it together.
