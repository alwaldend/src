package datasources

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/providerdata"
)

// templateDataSource resolves an existing template and its Docker option classification.
type templateDataSource struct {
	// providerData supplies the bounded shared client.
	providerData *providerdata.Data
}

// templateModel retains only public metadata and an environment-only classification.
type templateModel struct {
	// HashID selects the immutable template version.
	HashID types.String `tfsdk:"hash_id"`
	// Image is the resolved container image reference.
	Image types.String `tfsdk:"image"`
	// EnvironmentOnly excludes publishing, host networking, and unsupported flags.
	EnvironmentOnly types.Bool `tfsdk:"environment_only"`
}

// Metadata identifies exact template discovery under the provider prefix.
func (d *templateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

// Schema exposes safe metadata without storing template environment contents.
func (d *templateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Read one existing Vast.ai template by hash without creating resources.", Attributes: map[string]schema.Attribute{
		"hash_id":          schema.StringAttribute{Required: true, Description: "The template's 32-character hexadecimal content hash."},
		"image":            schema.StringAttribute{Computed: true, Description: "Template image including its configured tag."},
		"environment_only": schema.BoolAttribute{Computed: true, Description: "True only when Docker options contain literal environment assignments and no other options. Does not establish application trust."},
	}}
}

// Configure receives the provider's client without a separate credential context.
func (d *templateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// Read requires an exact hash match and omits environment values from Terraform state.
func (d *templateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data templateModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(data.HashID.ValueString()) {
		resp.Diagnostics.AddError("Invalid template hash", "hash_id must be a 32-character lowercase hexadecimal template hash.")
		return
	}
	remote, err := d.providerData.Client.Template(ctx, data.HashID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot read template", err.Error())
		return
	}
	data.Image = types.StringValue(remote.Image)
	data.EnvironmentOnly = types.BoolValue(remote.EnvironmentOnly)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// NewTemplate constructs a data source for provider registration.
func NewTemplate() datasource.DataSource { return &templateDataSource{} }
