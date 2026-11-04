package resources

import (
	"context"
	"errors"
	"regexp"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/providerdata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// environmentName matches the API's canonical uppercase variable names.
var environmentName = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// environmentResource manages one encrypted account environment variable.
type environmentResource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// environmentModel maps the canonical name and sensitive value.
type environmentModel struct {
	// ID is the canonical variable name.
	ID types.String `tfsdk:"id"`
	// Name is the uppercase identifier used by the API.
	Name types.String `tfsdk:"name"`
	// Value is stored in Terraform state and hidden in normal CLI output.
	Value types.String `tfsdk:"value"`
}

// Metadata names the account environment variable resource.
func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_variable"
}

// Schema makes renames replacement-only and marks values sensitive.
func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":    schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Canonical environment variable name."},
		"name":  schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Uppercase environment identifier. Changing it replaces the variable."},
		"value": schema.StringAttribute{Required: true, Sensitive: true, Description: "Variable value. Remains in Terraform state; use a secure backend."},
	}}
}

// Configure receives the shared client.
func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// ValidateConfig requires a canonical name and nonempty known value.
func (r *environmentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.Name.IsUnknown() && !environmentName.MatchString(data.Name.ValueString()) {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid environment name", "Use an uppercase environment identifier.")
	}
	if !data.Value.IsUnknown() && !data.Value.IsNull() && data.Value.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("value"), "Empty environment value", "The API requires a nonempty value.")
	}
}

// Create writes the variable and records its stable name.
func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.providerData.Client.WriteEnvironmentVariable(ctx, data.Name.ValueString(), data.Value.ValueString(), true); err != nil {
		resp.Diagnostics.AddError("Cannot create environment variable", err.Error())
		return
	}
	data.ID = data.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes only this variable and treats a missing map entry as absence.
func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	value, err := r.providerData.Client.EnvironmentVariable(ctx, data.ID.ValueString())
	if errors.Is(err, vastclient.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil && !errors.Is(err, vastclient.ErrValueRedacted) {
		resp.Diagnostics.AddError("Cannot read environment variable", err.Error())
		return
	}
	data.Name = data.ID
	if !errors.Is(err, vastclient.ErrValueRedacted) {
		data.Value = types.StringValue(value)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update changes the value in place.
func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.providerData.Client.WriteEnvironmentVariable(ctx, data.Name.ValueString(), data.Value.ValueString(), false); err != nil {
		resp.Diagnostics.AddError("Cannot update environment variable", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete confirms absence or removes this single variable.
func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.providerData.Client.EnvironmentVariable(ctx, data.ID.ValueString())
	if errors.Is(err, vastclient.ErrNotFound) {
		return
	}
	if err != nil && !errors.Is(err, vastclient.ErrValueRedacted) {
		resp.Diagnostics.AddError("Cannot read environment variable before deletion", err.Error())
		return
	}
	if err := r.providerData.Client.DeleteEnvironmentVariable(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Cannot delete environment variable", err.Error())
	}
}

// ImportState accepts the name; masked values require configuration on the next apply.
func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !environmentName.MatchString(req.ID) {
		resp.Diagnostics.AddError("Invalid environment import ID", "Use the canonical uppercase variable name.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// NewEnvironmentVariable constructs an unconfigured environment variable resource for provider registration.
func NewEnvironmentVariable() resource.Resource {
	return &environmentResource{}
}
