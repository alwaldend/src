package resources

import (
	"context"
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

// sshKeyResource manages a registered public key without generating private material.
type sshKeyResource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// sshKeyModel maps the numeric key ID and public key text.
type sshKeyModel struct {
	// ID is the registered key identifier.
	ID types.String `tfsdk:"id"`
	// PublicKey is a single OpenSSH authorized-key line.
	PublicKey types.String `tfsdk:"public_key"`
}

// Metadata names the resource under the provider prefix.
func (r *sshKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

// Schema exposes public material only; private keys remain outside this provider.
func (r *sshKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Registered SSH key ID."},
		"public_key": schema.StringAttribute{Required: true, Description: "OpenSSH public key, including an optional test_ identifying comment. Vast.ai validates the key."},
	}}
}

// Configure receives the account client.
func (r *sshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// ValidateConfig rejects empty, multiline, or private key material.
func (r *sshKeyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data sshKeyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || data.PublicKey.IsUnknown() {
		return
	}
	text := data.PublicKey.ValueString()
	if strings.TrimSpace(text) == "" || strings.ContainsAny(text, "\r\n") || strings.Contains(text, "PRIVATE KEY") {
		resp.Diagnostics.AddAttributeError(path.Root("public_key"), "Invalid public key", "Provide one OpenSSH public key line; never private key material.")
	}
}

// Create registers a public key and persists the returned ID.
func (r *sshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data sshKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := r.providerData.Client.CreateSSHKey(ctx, data.PublicKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot create SSH key", err.Error())
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(id, 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the key or removes state when its ID is absent from the collection.
func (r *sshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data sshKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := resourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid SSH key ID", err.Error())
		return
	}
	key, err := r.providerData.Client.SSHKey(ctx, id)
	if errors.Is(err, vastclient.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read SSH key", err.Error())
		return
	}
	data.PublicKey = types.StringValue(key)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update changes key material without replacing its server ID.
func (r *sshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data sshKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := resourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid SSH key ID", err.Error())
		return
	}
	if err := r.providerData.Client.UpdateSSHKey(ctx, id, data.PublicKey.ValueString()); err != nil {
		resp.Diagnostics.AddError("Cannot update SSH key", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes only the registered key ID in state.
func (r *sshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data sshKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := resourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid SSH key ID", err.Error())
		return
	}
	if err := r.providerData.Client.DeleteSSHKey(ctx, id); err != nil && !errors.Is(err, vastclient.ErrNotFound) {
		resp.Diagnostics.AddError("Cannot delete SSH key", err.Error())
	}
}

// ImportState validates the numeric key ID before refresh hydrates its public key.
func (r *sshKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := resourceID(req.ID); err != nil {
		resp.Diagnostics.AddError("Invalid SSH key import ID", err.Error())
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// NewSSHKey constructs an unconfigured SSH public key resource for provider registration.
func NewSSHKey() resource.Resource {
	return &sshKeyResource{}
}
