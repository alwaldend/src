package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/providerdata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// apiKeyResource generates credentials through its provider or a same-apply bootstrap issuer.
type apiKeyResource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// apiKeyModel retains sensitive generated credentials, issuer metadata, and key settings.
type apiKeyModel struct {
	// ID is the generated key ID.
	ID types.String `tfsdk:"id"`
	// Name identifies the project key.
	Name types.String `tfsdk:"name"`
	// TeamKey optionally selects a management key generated during the same apply.
	TeamKey types.String `tfsdk:"team_api_key"`
	// TeamID optionally binds only the newly generated key to an existing owned team.
	TeamID types.Int64 `tfsdk:"team_id"`
	// Permissions is the key's JSON permission object.
	Permissions types.String `tfsdk:"permissions"`
	// Key is the one-time generated credential, retained in sensitive state.
	Key types.String `tfsdk:"key"`
}

// issuerClient uses provider authentication unless same-apply bootstrap supplies an override.
func (r *apiKeyResource) issuerClient(key types.String) (*vastclient.Client, error) {
	if key.IsNull() {
		return r.providerData.Client, nil
	}
	client, err := r.providerData.Client.WithKey(key.ValueString())
	if err != nil {
		return nil, fmt.Errorf("configure bootstrap key issuer: %w", err)
	}
	return client, nil
}

// validateTeamPermissions prevents owner-context binding from widening a scoped key.
func (data apiKeyModel) validateTeamPermissions() error {
	if !data.TeamID.IsNull() && !data.TeamID.IsUnknown() && !data.Permissions.IsUnknown() && data.Permissions.ValueString() != "" {
		return errors.New("Existing-team owner binding requires empty permissions for explicit full access; create scoped keys with an already selected team issuer.")
	}
	return nil
}

// Metadata names generated project API keys.
func (r *apiKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

// Schema makes key settings replacement-only and keeps an optional bootstrap issuer sensitive.
func (r *apiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Generated API key ID."},
		"name":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Project key name."},
		"team_api_key": schema.StringAttribute{Optional: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()}, Description: "Optional management key generated during the same bootstrap apply. Otherwise lifecycle operations use the configured provider. A new explicit issuer replaces the key; removing the override preserves it after access verification."},
		"team_id":      schema.Int64Attribute{Optional: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}, Description: "Existing owned team to select on the newly generated key; never changes the issuer's context."},
		"permissions":  schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "JSON permission object, or an empty string for the account's default full access."},
		"key":          schema.StringAttribute{Computed: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Generated project token. Store directly in Vault; never output it. Remains in Terraform state."},
	}}
}

// Configure receives the configured provider client for ordinary key operations.
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
	if err := data.validateTeamPermissions(); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("permissions"), "Invalid team-key permissions", err.Error())
	}
	if !data.Name.IsUnknown() && strings.TrimSpace(data.Name.ValueString()) == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid key name", "Provide a nonempty key name.")
	}
	if !data.TeamID.IsNull() && !data.TeamID.IsUnknown() && data.TeamID.ValueInt64() <= 0 {
		resp.Diagnostics.AddAttributeError(path.Root("team_id"), "Invalid existing team", "The existing team ID must be positive.")
	}
	if !data.Permissions.IsUnknown() && data.Permissions.ValueString() != "" {
		var value map[string]any
		if err := json.Unmarshal([]byte(data.Permissions.ValueString()), &value); err != nil || value == nil {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "Invalid key permissions", "Provide a JSON permission object.")
		}
	}
}

// Create generates and optionally binds a key while retaining partial ownership.
func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data apiKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := data.validateTeamPermissions(); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("permissions"), "Invalid team-key permissions", err.Error())
		return
	}
	client, err := r.issuerClient(data.TeamKey)
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure key issuer", err.Error())
		return
	}
	var contextID int64
	if !data.TeamID.IsNull() {
		_, contextID, err = client.TeamContextByID(ctx, data.TeamID.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError("Cannot resolve owned team", err.Error())
			return
		}
	}
	id, key, err := client.CreateAPIKey(ctx, data.Name.ValueString(), data.Permissions.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot generate project API key", err.Error())
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(id, 10))
	data.Key = types.StringValue(key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() || data.TeamID.IsNull() {
		return
	}
	generated, err := client.WithKey(key)
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure generated key", err.Error())
		return
	}
	if err := generated.ActivateTeam(ctx, contextID, data.TeamID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Cannot bind generated key to owned team", err.Error())
	}
}

// Read verifies key identity and preserves its one-time secret value from state.
func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data apiKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := r.issuerClient(data.TeamKey)
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure key issuer", err.Error())
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

// Update removes issuer metadata only after verifying provider access to the existing key.
func (r *apiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planned, previous apiKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planned)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previous)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !planned.TeamKey.IsNull() || !planned.Name.Equal(previous.Name) || !planned.Permissions.Equal(previous.Permissions) || !planned.TeamID.Equal(previous.TeamID) {
		resp.Diagnostics.AddError("Unexpected API key update", "Key name, permissions, team binding, and a new explicit bootstrap issuer require replacement.")
		return
	}
	client, err := r.issuerClient(planned.TeamKey)
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure key issuer", err.Error())
		return
	}
	id, err := resourceID(previous.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid API key ID", err.Error())
		return
	}
	if err := client.ReadAPIKey(ctx, id); err != nil {
		resp.Diagnostics.AddError("Cannot verify key issuer", err.Error())
		return
	}
	planned.ID, planned.Key = previous.ID, previous.Key
	resp.Diagnostics.Append(resp.State.Set(ctx, &planned)...)
}

// Delete revokes the generated key through its configured issuer.
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
	teamClient, err := r.issuerClient(state.TeamKey)
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure key issuer", err.Error())
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
