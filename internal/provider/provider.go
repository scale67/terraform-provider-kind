package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure KindProvider satisfies various provider interfaces.
var _ provider.Provider = &KindProvider{}

// KindProvider defines the provider implementation.
type KindProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KindProviderModel describes the provider data model.
type KindProviderModel struct {
	// Provider configuration would go here if needed
}

func (p *KindProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kind"
	resp.Version = p.version
}

func (p *KindProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Kind provider allows Terraform to manage Kubernetes Kind clusters.",
		Attributes: map[string]schema.Attribute{
			// Provider configuration attributes can be added here if needed
		},
	}
}

func (p *KindProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KindProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Provider configuration logic can go here
}

func (p *KindProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
	}
}

func (p *KindProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Data sources can be added here if needed
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KindProvider{
			version: version,
		}
	}
}