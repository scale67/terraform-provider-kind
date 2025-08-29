#!/bin/bash

# Helper script to use the kind cluster created by Terraform
# Run this script after running 'terraform apply' to set up kubectl access

echo "Setting up kubectl access for the kind cluster..."

# Get the kubeconfig from Terraform output
KUBECONFIG_CONTENT=$(terraform output -raw kubeconfig)

# Create a temporary kubeconfig file
TEMP_KUBECONFIG=$(mktemp)
echo "$KUBECONFIG_CONTENT" > "$TEMP_KUBECONFIG"

# Set KUBECONFIG environment variable for this session
export KUBECONFIG="$TEMP_KUBECONFIG"

echo "KUBECONFIG set to: $TEMP_KUBECONFIG"
echo "Current kubectl context: $(kubectl config current-context)"
echo ""
echo "You can now use kubectl commands. For example:"
echo "  kubectl get nodes"
echo "  kubectl get pods --all-namespaces"
echo ""
echo "Note: This script sets KUBECONFIG for the current shell session only."
echo "To make it permanent, add this line to your shell profile:"
echo "  export KUBECONFIG=\"$TEMP_KUBECONFIG\""
echo ""
echo "Or copy the kubeconfig to your default location:"
echo "  cp \"$TEMP_KUBECONFIG\" ~/.kube/config"

