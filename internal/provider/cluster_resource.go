package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gopkg.in/yaml.v3"
	kindv1alpha4 "sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ClusterResource{}
var _ resource.ResourceWithImportState = &ClusterResource{}

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

// ClusterResource defines the resource implementation.
type ClusterResource struct{}

// ClusterResourceModel describes the resource data model.
type ClusterResourceModel struct {
	Name         types.String `tfsdk:"name"`
	KindConfig   types.String `tfsdk:"kind_config"`
	ConfigPath   types.String `tfsdk:"config_path"`
	NodeImage    types.String `tfsdk:"node_image"`
	WaitForReady types.Bool   `tfsdk:"wait_for_ready"`
	Kubeconfig   types.String `tfsdk:"kubeconfig"`
	Endpoint     types.String `tfsdk:"endpoint"`
	Id           types.String `tfsdk:"id"`
}

func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Kind Kubernetes cluster",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Kind cluster",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kind_config": schema.StringAttribute{
				MarkdownDescription: "Inline Kind configuration YAML",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config_path": schema.StringAttribute{
				MarkdownDescription: "Path to Kind configuration YAML file",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_image": schema.StringAttribute{
				MarkdownDescription: "Node Docker image to use for booting the cluster",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"wait_for_ready": schema.BoolAttribute{
				MarkdownDescription: "Wait for the cluster to be ready before completing",
				Optional:            true,
			},
			"kubeconfig": schema.StringAttribute{
				MarkdownDescription: "Kubeconfig for the created cluster",
				Computed:            true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Kubernetes API server endpoint",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Cluster identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Provider-level data can be retrieved from req.ProviderData if needed
}

func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either kind_config or config_path is provided
	if data.KindConfig.IsNull() && data.ConfigPath.IsNull() {
		resp.Diagnostics.AddError(
			"Configuration Required",
			"Either 'kind_config' or 'config_path' must be specified",
		)
		return
	}

	if !data.KindConfig.IsNull() && !data.ConfigPath.IsNull() {
		resp.Diagnostics.AddError(
			"Conflicting Configuration",
			"Only one of 'kind_config' or 'config_path' should be specified",
		)
		return
	}

	clusterName := data.Name.ValueString()
	tflog.Info(ctx, "Creating Kind cluster", map[string]interface{}{
		"cluster_name": clusterName,
	})

	// Build kind create command
	args := []string{"create", "cluster", "--name", clusterName}

	// Handle configuration
	var configFile string
	var shouldRemoveTempFile bool
	if !data.ConfigPath.IsNull() {
		configFile = data.ConfigPath.ValueString()
	} else if !data.KindConfig.IsNull() {
		// Create temporary config file
		tempFilePath := fmt.Sprintf("/tmp/kind-config-%s.yaml", clusterName)
		err := os.WriteFile(tempFilePath, []byte(data.KindConfig.ValueString()), 0644)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create temporary config file",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
		configFile = tempFilePath
		shouldRemoveTempFile = true
		defer func() {
			if shouldRemoveTempFile {
				os.Remove(tempFilePath)
			}
		}()
	}

	if configFile != "" {
		// Validate the config file
		if err := r.validateKindConfig(configFile); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Kind configuration",
				fmt.Sprintf("Error validating config: %s", err),
			)
			return
		}
		args = append(args, "--config", configFile)
	}

	// Add node image if specified
	if !data.NodeImage.IsNull() {
		args = append(args, "--image", data.NodeImage.ValueString())
	}

	// Add wait flag if specified
	if !data.WaitForReady.IsNull() && data.WaitForReady.ValueBool() {
		args = append(args, "--wait", "300s")
	}

	// Execute kind create command
	cmd := exec.CommandContext(ctx, "kind", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Kind cluster",
			fmt.Sprintf("Command: kind %s\nOutput: %s\nError: %s",
				strings.Join(args, " "), string(output), err),
		)
		return
	}

	// Get kubeconfig
	kubeconfig, err := r.getKubeconfig(ctx, clusterName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get kubeconfig",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	// Get cluster endpoint
	endpoint, err := r.getClusterEndpoint(ctx, clusterName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get cluster endpoint",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	// Update the model with computed values
	data.Id = types.StringValue(clusterName)
	data.Kubeconfig = types.StringValue(kubeconfig)
	data.Endpoint = types.StringValue(endpoint)

	tflog.Info(ctx, "Created Kind cluster", map[string]interface{}{
		"cluster_name": clusterName,
		"endpoint":     endpoint,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterName := data.Name.ValueString()

	// Check if cluster exists
	exists, err := r.clusterExists(ctx, clusterName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to check cluster existence",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	if !exists {
		// Cluster doesn't exist, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	// Refresh computed values
	kubeconfig, err := r.getKubeconfig(ctx, clusterName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get kubeconfig",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	endpoint, err := r.getClusterEndpoint(ctx, clusterName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get cluster endpoint",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	data.Kubeconfig = types.StringValue(kubeconfig)
	data.Endpoint = types.StringValue(endpoint)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Most updates require replacement, but we can handle wait_for_ready changes
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterName := data.Name.ValueString()
	tflog.Info(ctx, "Deleting Kind cluster", map[string]interface{}{
		"cluster_name": clusterName,
	})

	// Execute kind delete command
	cmd := exec.CommandContext(ctx, "kind", "delete", "cluster", "--name", clusterName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Log the error but don't fail deletion if cluster doesn't exist
		if !strings.Contains(string(output), "not found") {
			resp.Diagnostics.AddError(
				"Failed to delete Kind cluster",
				fmt.Sprintf("Command: kind delete cluster --name %s\nOutput: %s\nError: %s",
					clusterName, string(output), err),
			)
			return
		}
	}

	tflog.Info(ctx, "Deleted Kind cluster", map[string]interface{}{
		"cluster_name": clusterName,
	})
}

func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Helper functions

func (r *ClusterResource) validateKindConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config kindv1alpha4.Cluster
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Basic validation
	if config.Kind != "Cluster" {
		return fmt.Errorf("kind must be 'Cluster', got '%s'", config.Kind)
	}

	return nil
}

func (r *ClusterResource) clusterExists(ctx context.Context, name string) (bool, error) {
	cmd := exec.CommandContext(ctx, "kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list clusters: %w", err)
	}

	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cluster := range clusters {
		if strings.TrimSpace(cluster) == name {
			return true, nil
		}
	}
	return false, nil
}

func (r *ClusterResource) getKubeconfig(ctx context.Context, name string) (string, error) {
	cmd := exec.CommandContext(ctx, "kind", "get", "kubeconfig", "--name", name)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	return string(output), nil
}

func (r *ClusterResource) getClusterEndpoint(ctx context.Context, name string) (string, error) {
	cmd := exec.CommandContext(ctx, "kind", "get", "kubeconfig", "--name", name)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig for endpoint: %w", err)
	}

	// Parse kubeconfig to extract server endpoint
	var kubeconfig struct {
		Clusters []struct {
			Cluster struct {
				Server string `yaml:"server"`
			} `yaml:"cluster"`
		} `yaml:"clusters"`
	}

	if err := yaml.Unmarshal(output, &kubeconfig); err != nil {
		return "", fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	if len(kubeconfig.Clusters) == 0 {
		return "", fmt.Errorf("no clusters found in kubeconfig")
	}

	return kubeconfig.Clusters[0].Cluster.Server, nil
}
