package resources

import (
	"context"
	"errors"
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

// teamResource creates teams with management keys or imports existing identities.
type teamResource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// teamCreationConfirmedKey records successful team creation independently of its key.
const teamCreationConfirmedKey = "team_creation_confirmed"

// teamMissingKey requests explicit replacement while preserving credential cleanup.
const teamMissingKey = "team_missing"

// teamImportedKey distinguishes read-only adoption from confirmed creation.
const teamImportedKey = "team_imported"

// teamModel tracks durable ownership before each credential-dependent operation.
type teamModel struct {
	// ID is the immutable team name.
	ID types.String `tfsdk:"id"`
	// Name is the requested team name.
	Name types.String `tfsdk:"name"`
	// TeamID identifies the created team account.
	TeamID types.Int64 `tfsdk:"team_id"`
	// APIKey is generated management material, retained only in sensitive state.
	APIKey types.String `tfsdk:"api_key"`
	// KeyID identifies the management key for revocation.
	KeyID types.Int64 `tfsdk:"management_key_id"`
}

// Metadata names Terraform-managed teams.
func (r *teamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema exposes a zero-credit team and its sensitive generated management key.
func (r *teamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":                schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Team name used as resource identity."},
		"name":              schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "New team name. Creation transfers zero credits; renaming requires replacement."},
		"team_id":           schema.Int64Attribute{Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}, Description: "Team account ID."},
		"api_key":           schema.StringAttribute{Computed: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Generated team management credential, null for imported teams. Secure Terraform state; the personal key's context is unchanged."},
		"management_key_id": schema.Int64Attribute{Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}, Description: "Generated management key ID retained for cleanup, null for imported teams."},
	}}
}

// Configure receives the personal bootstrap client.
func (r *teamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// ValidateConfig rejects empty or path-like names before team creation.
func (r *teamResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data teamModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || data.Name.IsUnknown() {
		return
	}
	if strings.TrimSpace(data.Name.ValueString()) == "" || strings.ContainsAny(data.Name.ValueString(), "/\\\r\n") {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid team name", "Use a nonempty team name without path separators.")
	}
}

// Create persists a generated management key before making any team-context changes.
func (r *teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data teamModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keyID, key, err := r.providerData.Client.CreateAPIKey(ctx, data.Name.ValueString()+"_management", "")
	if err != nil {
		resp.Diagnostics.AddError("Cannot generate team management key", err.Error())
		return
	}
	data.ID = data.Name
	data.TeamID = types.Int64Null()
	data.KeyID = types.Int64Value(keyID)
	data.APIKey = types.StringValue(key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, teamCreationConfirmedKey, []byte("false"))...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := r.providerData.Client.WithKey(key)
	if err != nil {
		resp.Diagnostics.AddError("Cannot configure generated key", err.Error())
		return
	}
	if err := client.CreateTeam(ctx, data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Cannot create team", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, teamCreationConfirmedKey, []byte("true"))...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, contextID, err := r.providerData.Client.TeamContext(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot read created team", err.Error())
		return
	}
	data.TeamID = types.Int64Value(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if err := client.ActivateTeam(ctx, contextID, id); err != nil {
		resp.Diagnostics.AddError("Cannot activate generated team management key", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes team identity without retrieving or rotating credential values.
func (r *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data teamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	confirmed, diags := req.Private.GetKey(ctx, teamCreationConfirmedKey)
	resp.Diagnostics.Append(diags...)
	imported, diags := req.Private.GetKey(ctx, teamImportedKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || (string(confirmed) != "true" && string(imported) != "true") {
		// Failed creation owns only its generated key, even if this name exists.
		return
	}
	id, name := data.TeamID.ValueInt64(), data.Name.ValueString()
	var err error
	if id > 0 {
		name, _, err = r.providerData.Client.TeamContextByID(ctx, id)
	} else {
		id, _, err = r.providerData.Client.TeamContext(ctx, name)
	}
	if errors.Is(err, vastclient.ErrNotFound) {
		if string(imported) == "true" {
			resp.State.RemoveResource(ctx)
			return
		}
		// The independent key can survive team deletion; only apply may revoke it.
		resp.Diagnostics.Append(resp.Private.SetKey(ctx, teamMissingKey, []byte("true"))...)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read team", err.Error())
		return
	}
	data.Name = types.StringValue(name)
	data.TeamID = types.Int64Value(id)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, teamMissingKey, []byte("false"))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ModifyPlan replaces missing teams without abandoning their independently owned key.
func (r *teamResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	missing, diags := req.Private.GetKey(ctx, teamMissingKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || string(missing) != "true" {
		return
	}
	resp.RequiresReplace = append(resp.RequiresReplace, path.Root("id"))
	for _, name := range []string{"id", "api_key"} {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(name), types.StringUnknown())...)
	}
	for _, name := range []string{"team_id", "management_key_id"} {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(name), types.Int64Unknown())...)
	}
}

// Update refuses changes that should have been planned as replacement.
func (r *teamResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unsupported team update", "Team name changes require replacement.")
}

// Delete selects only this generated management context and revokes its key afterward.
func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data teamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	imported, diags := req.Private.GetKey(ctx, teamImportedKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if string(imported) == "true" {
		resp.Diagnostics.AddError("Imported teams cannot be destroyed", "Imported teams have no recovered management credential. Preserve the live team and detach it explicitly from Terraform state if ownership must be removed.")
		return
	}
	confirmed, diags := req.Private.GetKey(ctx, teamCreationConfirmedKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if string(confirmed) == "true" && data.TeamID.IsNull() {
		resp.Diagnostics.AddError("Missing owned team account ID", "Team creation was confirmed but account discovery did not finish. Refresh successfully before destroying this resource so cleanup can identify the created account.")
		return
	}
	if string(confirmed) == "true" && data.TeamID.ValueInt64() > 0 {
		_, contextID, err := r.providerData.Client.TeamContextByID(ctx, data.TeamID.ValueInt64())
		if err != nil && !errors.Is(err, vastclient.ErrNotFound) {
			resp.Diagnostics.AddError("Cannot read owned team before deletion", err.Error())
			return
		}
		if err == nil {
			if data.APIKey.IsNull() {
				resp.Diagnostics.AddError("Missing team management key", "Team creation was incomplete; recover a management credential before destruction.")
				return
			}
			client, err := r.providerData.Client.WithKey(data.APIKey.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Cannot configure team cleanup", err.Error())
				return
			}
			if err := client.ActivateTeam(ctx, contextID, data.TeamID.ValueInt64()); err != nil {
				resp.Diagnostics.AddError("Cannot select team for cleanup", err.Error())
				return
			}
			if err := client.DeleteTeam(ctx); err != nil {
				resp.Diagnostics.AddError("Cannot delete team", err.Error())
				return
			}
		}
	}
	if data.KeyID.ValueInt64() > 0 {
		if err := r.providerData.Client.DeleteAPIKey(ctx, data.KeyID.ValueInt64()); err != nil && !errors.Is(err, vastclient.ErrNotFound) {
			resp.Diagnostics.AddError("Cannot revoke management key", err.Error())
		}
	}
}

// ImportState adopts an owned team identity without generating or recovering keys.
func (r *teamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := resourceID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid team import ID", "Use the canonical positive decimal team account ID.")
		return
	}
	name, _, err := r.providerData.Client.TeamContextByID(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Cannot read imported team", err.Error())
		return
	}
	data := teamModel{
		ID:     types.StringValue(name),
		Name:   types.StringValue(name),
		TeamID: types.Int64Value(id),
		APIKey: types.StringNull(),
		KeyID:  types.Int64Null(),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, teamImportedKey, []byte("true"))...)
}

// NewTeam constructs an unconfigured team resource for provider registration.
func NewTeam() resource.Resource {
	return &teamResource{}
}
