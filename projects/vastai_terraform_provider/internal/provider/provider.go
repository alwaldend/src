package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/datasources"
	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/providerdata"
	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/resources"
)

// vastProvider implements the Framework provider contract.
type vastProvider struct {
	// version is the build version reported through the plugin protocol.
	version string
}

// providerModel holds optional provider configuration.
type providerModel struct {
	// APIKey overrides VAST_API_KEY when configured.
	APIKey types.String `tfsdk:"api_key"`
	// Endpoint overrides the default REST root for controlled API testing.
	Endpoint types.String `tfsdk:"endpoint"`
}

// New returns the provider factory consumed by the plugin server.
func New(version string) func() provider.Provider {
	return func() provider.Provider { return &vastProvider{version: version} }
}

// Metadata identifies the provider and its version to Terraform.
func (p *vastProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vastai"
	resp.Version = p.version
}

// Schema defines the credential and API-root configuration.
func (p *vastProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Manage Vast.ai compute and account resources.", Attributes: map[string]schema.Attribute{
		"api_key":  schema.StringAttribute{Optional: true, Sensitive: true, Description: "Vast.ai API key. Defaults to VAST_API_KEY; prefer environment injection to keep credentials out of configuration."},
		"endpoint": schema.StringAttribute{Optional: true, Description: "REST API root. Defaults to https://console.vast.ai/api/v0. HTTP is allowed only for literal loopback testing."},
	}}
}

// Configure validates configuration and shares a client with resources and data sources.
func (p *vastProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.APIKey.IsUnknown() || data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddError("Unknown provider configuration", "API key and endpoint must be known before configuring the provider.")
		return
	}
	key := os.Getenv("VAST_API_KEY")
	if !data.APIKey.IsNull() {
		key = data.APIKey.ValueString()
	}
	if strings.TrimSpace(key) == "" {
		resp.Diagnostics.AddError("Missing API key", "Set VAST_API_KEY or the sensitive api_key provider attribute.")
		return
	}
	endpoint := "https://console.vast.ai/api/v0"
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}
	client, err := vastclient.New(vastclient.Options{Endpoint: endpoint, APIKey: key})
	if err != nil {
		resp.Diagnostics.AddError("Invalid API configuration", err.Error())
		return
	}
	providerData := &providerdata.Data{Client: client}
	resp.ResourceData = providerData
	resp.DataSourceData = providerData
}

// Resources registers compute and non-compute account resources.
func (p *vastProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewInstance,
		resources.NewSSHKey,
		resources.NewEnvironmentVariable,
		resources.NewTeamRole,
		resources.NewTeam,
		resources.NewAPIKey,
	}
}

// DataSources registers marketplace offer and template discovery.
func (p *vastProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{datasources.NewOffers, datasources.NewTemplate}
}
