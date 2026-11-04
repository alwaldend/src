package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

// teamRoleResource manages an unassigned custom role in the authenticated team.
type teamRoleResource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// teamRoleModel retains both name-based and numeric API identifiers.
type teamRoleModel struct {
	// ID is the stable role name used for reads and deletion.
	ID types.String `tfsdk:"id"`
	// Name is the replacement-only custom role name.
	Name types.String `tfsdk:"name"`
	// RoleID is the server's numeric update identifier.
	RoleID types.Int64 `tfsdk:"role_id"`
	// Permissions is a JSON object, preferably produced with jsonencode.
	Permissions types.String `tfsdk:"permissions"`
	// EffectivePermissions records the last observed server-normalized permissions.
	EffectivePermissions types.String `tfsdk:"effective_permissions"`
}

// Metadata names the team role resource.
func (r *teamRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_role"
}

// Schema declares name identity, numeric update ID, and permission JSON.
func (r *teamRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":                    schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Role name used as Terraform identity."},
		"name":                  schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Custom team role name. Changing it requires replacement."},
		"role_id":               schema.Int64Attribute{Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}, Description: "Numeric role ID used for updates."},
		"permissions":           schema.StringAttribute{Required: true, Description: "Requested JSON permission object. Use jsonencode; Vast.ai adds baseline permissions."},
		"effective_permissions": schema.StringAttribute{Computed: true, Description: "Server-normalized permission object, including Vast.ai baseline grants. This resource does not assign members."},
	}}
}

// Configure receives the provider's team-scoped client.
func (r *teamRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// ValidateConfig rejects empty or control-bearing names and non-object permissions.
func (r *teamRoleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data teamRoleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.Name.IsUnknown() && (strings.TrimSpace(data.Name.ValueString()) == "" || strings.ContainsAny(data.Name.ValueString(), "\\\r\n") || data.Name.ValueString() == "." || data.Name.ValueString() == "..") {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid role name", "Provide a nonempty role name without backslashes or line breaks.")
	}
	if !data.Permissions.IsUnknown() {
		var value map[string]any
		if err := json.Unmarshal([]byte(data.Permissions.ValueString()), &value); err != nil || value == nil {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "Invalid permissions", "Provide a JSON object, preferably with jsonencode.")
		}
	}
}

// readRole refreshes identity and distinguishes API normalization from later drift.
func (r *teamRoleResource) readRole(ctx context.Context, data *teamRoleModel) error {
	role, err := r.providerData.Client.TeamRole(ctx, data.ID.ValueString())
	if err != nil {
		return fmt.Errorf("refresh role: %w", err)
	}
	data.Name = types.StringValue(role.Name)
	data.RoleID = types.Int64Value(role.ID)
	var prior, current any
	priorErr := json.Unmarshal([]byte(data.EffectivePermissions.ValueString()), &prior)
	if err := json.Unmarshal(role.Permissions, &current); err != nil {
		return fmt.Errorf("decode role permissions: %w", err)
	}
	effective, err := json.Marshal(current)
	if err != nil {
		return fmt.Errorf("encode effective role permissions: %w", err)
	}
	if data.Permissions.IsNull() || data.Permissions.IsUnknown() || (priorErr == nil && !reflect.DeepEqual(prior, current)) {
		data.Permissions = types.StringValue(string(effective))
	}
	data.EffectivePermissions = types.StringValue(string(effective))
	return nil
}

// Create records role ownership before fetching its numeric identifier.
func (r *teamRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data teamRoleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.providerData.Client.WriteTeamRole(ctx, 0, data.Name.ValueString(), data.Permissions.ValueString()); err != nil {
		resp.Diagnostics.AddError("Cannot create team role", err.Error())
		return
	}
	data.ID = data.Name
	data.RoleID = types.Int64Null()
	data.EffectivePermissions = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if err := r.readRole(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Cannot read created team role", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the role or removes state on a definitive not-found response.
func (r *teamRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data teamRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.readRole(ctx, &data)
	if errors.Is(err, vastclient.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read team role", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update changes permissions through the server's numeric role ID.
func (r *teamRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data teamRoleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.RoleID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Missing role ID", "Refresh the role to obtain its numeric ID before updating.")
		return
	}
	if err := r.providerData.Client.WriteTeamRole(ctx, data.RoleID.ValueInt64(), data.Name.ValueString(), data.Permissions.ValueString()); err != nil {
		resp.Diagnostics.AddError("Cannot update team role", err.Error())
		return
	}
	if err := r.readRole(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Cannot read updated team role", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes only this role; it never removes team members or the team itself.
func (r *teamRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data teamRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.providerData.Client.DeleteTeamRole(ctx, data.ID.ValueString()); err != nil && !errors.Is(err, vastclient.ErrNotFound) {
		resp.Diagnostics.AddError("Cannot delete team role", err.Error())
	}
}

// ImportState accepts the role name, then Read hydrates permissions and numeric ID.
func (r *teamRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if strings.TrimSpace(req.ID) == "" || strings.ContainsAny(req.ID, "\\\r\n") || req.ID == "." || req.ID == ".." {
		resp.Diagnostics.AddError("Invalid role import ID", "Use a nonempty role name without backslashes or line breaks.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// NewTeamRole constructs an unconfigured team role resource for provider registration.
func NewTeamRole() resource.Resource {
	return &teamRoleResource{}
}
