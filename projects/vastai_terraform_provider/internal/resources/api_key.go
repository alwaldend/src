package resources

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/providerdata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// apiKeyResource generates a project credential using a Terraform-managed team key.
type apiKeyResource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// apiKeyModel retains sensitive generated credentials and replacement-only inputs.
type apiKeyModel struct {
	// ID is the generated key ID.
	ID types.String `tfsdk:"id"`
	// Name identifies the project key.
	Name types.String `tfsdk:"name"`
	// TeamKey selects the independently generated team management context.
	TeamKey types.String `tfsdk:"team_api_key"`
	// Permissions is the key's JSON permission object.
	Permissions types.String `tfsdk:"permissions"`
	// Key is the one-time generated credential, retained in sensitive state.
	Key types.String `tfsdk:"key"`
}

// Metadata names generated project API keys.
func (r *apiKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

// Schema makes credential changes explicit replacements and marks both keys sensitive.
func (r *apiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Generated API key ID."},
		"name":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Project key name."},
		"team_api_key": schema.StringAttribute{Required: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Terraform-generated team management credential from vastai_team.api_key."},
		"permissions":  schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "JSON permission object for the generated project key."},
		"key":          schema.StringAttribute{Computed: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Generated project token. Store directly in Vault; never output it. Remains in Terraform state."},
	}}
}

// Configure receives the endpoint-owning bootstrap client.
func (r *apiKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// ValidateConfig rejects empty names and malformed permission objects.
func (r *apiKeyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data apiKeyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.Name.IsUnknown() && strings.TrimSpace(data.Name.ValueString()) == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid key name", "Provide a nonempty key name.")
	}
	if !data.Permissions.IsUnknown() {
		var value map[string]any
		if err := json.Unmarshal([]byte(data.Permissions.ValueString()), &value); err != nil || value == nil {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "Invalid key permissions", "Provide a JSON permission object.")
		}
	}
}

// Create generates the project token inside its selected team context.
func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data apiKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := r.providerData.Client.WithKey(data.TeamKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure team key", err.Error())
		return
	}
	id, key, err := client.CreateAPIKey(ctx, data.Name.ValueString(), data.Permissions.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot generate project API key", err.Error())
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(id, 10))
	data.Key = types.StringValue(key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read verifies key identity and preserves its one-time secret value from state.
func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data apiKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := r.providerData.Client.WithKey(data.TeamKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure team key", err.Error())
		return
	}
	id, err := resourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid API key ID", err.Error())
		return
	}
	err = client.ReadAPIKey(ctx, id)
	if errors.Is(err, vastclient.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read project API key", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update guards the Framework-required method for this immutable resource.
// Every configurable attribute requires replacement, so valid plans bypass it.
func (r *apiKeyResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unexpected API key update", "All configurable API-key attributes require replacement. Terraform must revoke and recreate this key instead of updating it in place.")
}

// Delete revokes the project key while its management context remains alive.
func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keyID, err := resourceID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid API key ID", err.Error())
		return
	}
	teamClient, err := r.providerData.Client.WithKey(state.TeamKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure team key", err.Error())
		return
	}
	err = teamClient.DeleteAPIKey(ctx, keyID)
	if errors.Is(err, vastclient.ErrNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot revoke project API key", err.Error())
	}
}

// NewAPIKey constructs an unconfigured API key resource for provider registration.
func NewAPIKey() resource.Resource {
	return &apiKeyResource{}
}
