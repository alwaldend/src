package resources

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/providerdata"
)

const (
	// runtypeSSH selects the default SSH launch mode.
	runtypeSSH = "ssh"
	// runtypeSSHProxy selects SSH access through the proxy.
	runtypeSSHProxy = "ssh_proxy"
	// runtypeSSHDirect selects direct SSH access.
	runtypeSSHDirect = "ssh_direct"
	// runtypeJupyter selects the default Jupyter launch mode.
	runtypeJupyter = "jupyter"
	// runtypeJupyterProxy selects Jupyter access through the proxy.
	runtypeJupyterProxy = "jupyter_proxy"
	// runtypeJupyterDirect selects direct Jupyter access.
	runtypeJupyterDirect = "jupyter_direct"
	// apiRuntypeSSHDirect is the legacy direct-SSH API spelling.
	apiRuntypeSSHDirect = "ssh_direc"
	// apiRuntypeJupyterDirect is the legacy direct-Jupyter API spelling.
	apiRuntypeJupyterDirect = "jupyter_direc"
	// instanceStateRunning is the requested running power state.
	instanceStateRunning = "running"
	// instanceStateStopped is the requested stopped power state.
	instanceStateStopped = "stopped"
	// maximumInstanceLabelLength bounds an instance label in Unicode characters.
	maximumInstanceLabelLength = 1024
	// maximumStartupBytes is Vast's inline startup command limit, including encoding.
	maximumStartupBytes = 16384
)

// supportedRuntypes owns the canonical launch modes used by validation and schema descriptions.
var supportedRuntypes = [...]string{
	runtypeSSH, runtypeSSHProxy, runtypeSSHDirect,
	runtypeJupyter, runtypeJupyterProxy, runtypeJupyterDirect,
}

// instanceResource owns Terraform lifecycle callbacks for a Vast.ai rental.
type instanceResource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// instanceModel maps supported rental configuration and observed connection data.
type instanceModel struct {
	// TemplateHashID selects optional template defaults only at creation.
	TemplateHashID types.String `tfsdk:"template_hash_id"`
	// ID is the decimal contract ID.
	ID types.String `tfsdk:"id"`
	// OfferID is the marketplace offer used at creation, unavailable on import.
	OfferID types.Int64 `tfsdk:"offer_id"`
	// Image selects the Docker image.
	Image types.String `tfsdk:"image"`
	// DiskGB is the disk allocation in GB.
	DiskGB types.Float64 `tfsdk:"disk_gb"`
	// Label is the mutable rental label.
	Label types.String `tfsdk:"label"`
	// State is the desired power state, separate from observed Status.
	State types.String `tfsdk:"state"`
	// Runtype is the creation-only container launch mode.
	Runtype types.String `tfsdk:"runtype"`
	// DockerOptions contains creation-only environment and publishing overrides.
	DockerOptions types.String `tfsdk:"docker_options"`
	// MappedPorts contains observed TCP/UDP bindings keyed by container port/protocol.
	MappedPorts types.Map `tfsdk:"mapped_ports"`
	// Onstart is the creation-only startup command.
	Onstart types.String `tfsdk:"onstart"`
	// CompressOnstart opts into gzip transport for images with the required shell utilities.
	CompressOnstart types.Bool `tfsdk:"compress_onstart"`
	// Status is the last observed container status.
	Status types.String `tfsdk:"actual_status"`
	// SSHHost is the last observed SSH endpoint.
	SSHHost types.String `tfsdk:"ssh_host"`
	// SSHPort is the last observed SSH port.
	SSHPort types.Int64 `tfsdk:"ssh_port"`
	// PublicIP is the last observed public address.
	PublicIP types.String `tfsdk:"public_ip"`
}

// Metadata identifies the resource under the provider's prefix.
func (r *instanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema declares creation-only inputs, mutable intent, and observed outputs.
func (r *instanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An on-demand Vast.ai GPU instance. Apply records API acceptance; refresh observes boot and connection details.",
		Attributes: map[string]schema.Attribute{
			"template_hash_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Optional 32-character template hash. Creation-only; image and launch settings override template defaults, while template environment and ports are inherited.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(
					func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
						// Older image-only state has no selector; the empty default is equivalent.
						resp.RequiresReplace = !(req.StateValue.IsNull() && req.PlanValue.ValueString() == "")
					},
					"Changes to a nonempty template selector replace the instance.",
					"Changes to a nonempty template selector replace the instance.",
				)},
			},
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "Decimal instance contract ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"offer_id": schema.Int64Attribute{
				Required:      true,
				Description:   "Positive marketplace offer ID. Creation-only and unavailable when importing; see import reconciliation guidance.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"image": schema.StringAttribute{
				Required:      true,
				Description:   "Docker image to launch. Changing it replaces the instance.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"disk_gb": schema.Float64Attribute{
				Required:      true,
				Description:   "Positive local disk allocation in GB. Changing it replaces the instance.",
				PlanModifiers: []planmodifier.Float64{float64planmodifier.RequiresReplace()},
			},
			"label": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: fmt.Sprintf("Mutable instance label, up to %d characters.", maximumInstanceLabelLength),
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(instanceStateRunning),
				Description: fmt.Sprintf("Desired power state: %s or %s. Stopped instances still incur storage charges.", instanceStateRunning, instanceStateStopped),
			},
			"runtype": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Default:       stringdefault.StaticString(runtypeSSH),
				Description:   "Creation-only launch mode: ssh, ssh_proxy, ssh_direct, jupyter, jupyter_proxy, or jupyter_direct.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"docker_options": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
				Description: "Creation-only Docker environment and port options. Merged with template defaults; sensitive but stored in state. API-returned environment is never read.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIf(
					func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = !(req.StateValue.IsNull() && req.PlanValue.ValueString() == "")
					}, "Changed Docker options replace the instance.", "Changed Docker options replace the instance.")},
			},
			"mapped_ports": schema.MapAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Observed host ports keyed by container port/protocol, for example 7860/tcp; null before refresh. Distinct host ports for one container port are rejected.",
			},
			"onstart": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Sensitive:     true,
				Default:       stringdefault.StaticString(""),
				Description:   "Creation-only startup command. Sensitive in CLI output; still stored in Terraform state.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"compress_onstart": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Compress startup when smaller. Requires /bin/bash, base64, gzip and mktemp in the image. Creation-only; the original sensitive script remains in state.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplaceIf(
					func(_ context.Context, req planmodifier.BoolRequest, resp *boolplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = !(req.StateValue.IsNull() && !req.PlanValue.ValueBool())
					}, "Changed startup compression replaces the instance.", "Changed startup compression replaces the instance.")},
			},
			"actual_status": schema.StringAttribute{
				Computed:    true,
				Description: "Observed container status; null before refresh or while unavailable.",
			},
			"ssh_host": schema.StringAttribute{
				Computed:    true,
				Description: "SSH host; null until reported by the API.",
			},
			"ssh_port": schema.Int64Attribute{
				Computed:    true,
				Description: "SSH port; null until reported by the API.",
			},
			"public_ip": schema.StringAttribute{
				Computed:    true,
				Description: "Public IP address; null until reported by the API.",
			},
		},
	}
}

// Configure receives the configured provider client.
func (r *instanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// validateInstance rejects known invalid parameters before any API operation.
func validateInstance(data instanceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	if !data.TemplateHashID.IsUnknown() && !data.TemplateHashID.IsNull() && data.TemplateHashID.ValueString() != "" && !regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(data.TemplateHashID.ValueString()) {
		diagnostics.AddAttributeError(path.Root("template_hash_id"), "Invalid template hash", "template_hash_id must be empty or a 32-character lowercase hexadecimal template hash.")
	}
	if !data.OfferID.IsUnknown() && !data.OfferID.IsNull() && data.OfferID.ValueInt64() <= 0 {
		diagnostics.AddAttributeError(path.Root("offer_id"), "Invalid offer ID", "offer_id must be a positive integer.")
	}
	if !data.Image.IsUnknown() && !data.Image.IsNull() && strings.TrimSpace(data.Image.ValueString()) == "" {
		diagnostics.AddAttributeError(path.Root("image"), "Invalid image", "image must not be empty.")
	}
	if !data.DiskGB.IsUnknown() && !data.DiskGB.IsNull() && data.DiskGB.ValueFloat64() <= 0 {
		diagnostics.AddAttributeError(path.Root("disk_gb"), "Invalid disk allocation", "disk_gb must be positive.")
	}
	if !data.Label.IsUnknown() && utf8.RuneCountInString(data.Label.ValueString()) > maximumInstanceLabelLength {
		diagnostics.AddAttributeError(path.Root("label"), "Invalid label", fmt.Sprintf("label must be at most %d characters.", maximumInstanceLabelLength))
	}
	if !data.State.IsUnknown() && !data.State.IsNull() && data.State.ValueString() != instanceStateRunning && data.State.ValueString() != instanceStateStopped {
		diagnostics.AddAttributeError(path.Root("state"), "Invalid power state", fmt.Sprintf("state must be %s or %s.", instanceStateRunning, instanceStateStopped))
	}
	if !data.Runtype.IsUnknown() && !data.Runtype.IsNull() && !slices.Contains(supportedRuntypes[:], data.Runtype.ValueString()) {
		diagnostics.AddAttributeError(path.Root("runtype"), "Invalid launch mode", fmt.Sprintf("Use one of: %s.", strings.Join(supportedRuntypes[:], ", ")))
	}
	if !data.Onstart.IsUnknown() && !data.CompressOnstart.IsUnknown() {
		if _, err := instanceStartupCommand(data.Onstart.ValueString(), data.CompressOnstart.ValueBool()); err != nil {
			diagnostics.AddAttributeError(path.Root("onstart"), "Invalid startup command", err.Error())
		}
	}
	return diagnostics
}

// instanceStartupCommand prepares bounded API input without including script contents in errors.
func instanceStartupCommand(script string, compress bool) (string, error) {
	command := script
	if compress && script != "" {
		var compressed bytes.Buffer
		writer, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
		if err != nil {
			return "", fmt.Errorf("initialize startup compression: %w", err)
		}
		if _, err := writer.Write([]byte(script)); err != nil {
			return "", fmt.Errorf("compress startup: %w", err)
		}
		if err := writer.Close(); err != nil {
			return "", fmt.Errorf("finish startup compression: %w", err)
		}
		encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())
		wrapped := "# Compressed Vast.ai startup\nset -eu\numask 077\nstartup_file=$(mktemp)\ntrap 'rm -f \"$startup_file\"' EXIT\nprintf '%s' '" + encoded + "' | base64 -d | gzip -dc > \"$startup_file\"\n/bin/bash \"$startup_file\"\n"
		if len(wrapped) < len(script) {
			command = wrapped
		}
	}
	if len(command) > maximumStartupBytes {
		return "", fmt.Errorf("startup command exceeds %d bytes after the selected encoding; reduce the payload or enable compression for a compatible image", maximumStartupBytes)
	}
	return command, nil
}

// ValidateConfig reports invalid known values during Terraform validation and planning.
func (r *instanceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data instanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(validateInstance(data)...)
	}
}

// Create records the accepted contract before making any further remote requests.
func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data instanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(validateInstance(data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	startup, err := instanceStartupCommand(data.Onstart.ValueString(), data.CompressOnstart.ValueBool())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("onstart"), "Invalid startup command", err.Error())
		return
	}
	id, err := r.providerData.Client.Create(ctx, data.OfferID.ValueInt64(), vastclient.CreateInput{TemplateHashID: data.TemplateHashID.ValueString(), Image: data.Image.ValueString(), Disk: data.DiskGB.ValueFloat64(), Label: data.Label.ValueString(), State: data.State.ValueString(), Runtype: data.Runtype.ValueString(), Onstart: startup, DockerOptions: data.DockerOptions.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Cannot create instance", err.Error())
		return
	}
	data.ID = types.StringValue(strconv.FormatInt(id, 10))
	data.Status = types.StringNull()
	data.SSHHost = types.StringNull()
	data.SSHPort = types.Int64Null()
	data.PublicIP = types.StringNull()
	data.MappedPorts = types.MapNull(types.Int64Type)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes drift and removes state for HTTP 404 or an explicit null instance.
func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data instanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := resourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance ID", err.Error())
		return
	}
	remote, err := r.providerData.Client.Get(ctx, id)
	if errors.Is(err, vastclient.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot read instance", err.Error())
		return
	}
	if strings.TrimSpace(remote.Image) == "" {
		resp.Diagnostics.AddError("Incomplete instance response", "Vast.ai returned no valid image identifier. The previous instance state is retained.")
		return
	}
	if remote.Disk <= 0 {
		resp.Diagnostics.AddError("Incomplete instance response", "Vast.ai returned no positive disk allocation. The previous instance state is retained.")
		return
	}
	runtype, err := normalizeRuntype(remote.Runtype, data.Runtype.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot read instance launch mode", err.Error())
		return
	}
	bindings := map[string]int64{}
	for name, ports := range remote.Ports {
		if len(ports) == 0 {
			continue
		}
		var port int64
		for _, binding := range ports {
			parsed, err := strconv.ParseInt(binding.HostPort, 10, 64)
			if err != nil || parsed < 1 || parsed > 65535 {
				resp.Diagnostics.AddError("Invalid port binding", "Instance returned an invalid mapped host port; prior state is retained.")
				return
			}
			if port != 0 && port != parsed {
				resp.Diagnostics.AddError("Ambiguous port binding", "Instance returned distinct host ports for a container port; prior state is retained.")
				return
			}
			port = parsed
		}
		bindings[name] = port
	}
	mapped, diagnostics := types.MapValueFrom(ctx, types.Int64Type, bindings)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.MappedPorts = mapped
	if data.DockerOptions.IsNull() {
		data.DockerOptions = types.StringValue("")
	}
	data.Image = types.StringValue(remote.Image)
	if data.TemplateHashID.IsNull() {
		data.TemplateHashID = types.StringValue("")
	}
	data.DiskGB = types.Float64Value(remote.Disk)
	data.Runtype = types.StringValue(runtype)
	data.Label = types.StringValue(remote.Label)
	data.State = types.StringValue(remote.State)
	if remote.Onstart != nil {
		startup, err := instanceStartupCommand(data.Onstart.ValueString(), data.CompressOnstart.ValueBool())
		if data.Onstart.IsNull() || err != nil || startup != *remote.Onstart {
			data.Onstart = types.StringPointerValue(remote.Onstart)
		}
	} else if data.Onstart.IsNull() {
		data.Onstart = types.StringValue("")
	}
	if data.CompressOnstart.IsNull() {
		data.CompressOnstart = types.BoolValue(false)
	}
	data.Status = types.StringPointerValue(remote.Status)
	data.SSHHost = types.StringPointerValue(remote.SSHHost)
	data.SSHPort = types.Int64PointerValue(remote.SSHPort)
	data.PublicIP = types.StringPointerValue(remote.PublicIP)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// normalizeRuntype retains compatible state or selects a supported import mode.
// Jupyter modes take precedence over their accompanying SSH access modes.
func normalizeRuntype(value, prior string) (string, error) {
	modes := map[string]bool{}
	for _, mode := range strings.Fields(value) {
		switch mode {
		case apiRuntypeSSHDirect:
			mode = runtypeSSHDirect
		case apiRuntypeJupyterDirect:
			mode = runtypeJupyterDirect
		}
		if !slices.Contains(supportedRuntypes[:], mode) {
			return "", errors.New("instance response contains an unsupported launch mode")
		}
		modes[mode] = true
	}
	if len(modes) == 0 {
		return "", errors.New("instance response omitted its launch mode")
	}
	candidates := []string{runtypeSSH, runtypeSSHDirect, runtypeSSHProxy}
	if modes[runtypeJupyter] || modes[runtypeJupyterDirect] || modes[runtypeJupyterProxy] {
		candidates = []string{runtypeJupyter, runtypeJupyterDirect, runtypeJupyterProxy}
	}
	if prior == candidates[0] {
		return prior, nil
	}
	for _, mode := range candidates {
		if prior == mode && modes[mode] {
			return prior, nil
		}
	}
	for _, mode := range candidates {
		if modes[mode] {
			return mode, nil
		}
	}
	return "", errors.New("instance response omitted a supported launch mode")
}

// Update submits mutable intent while preserving last-observed connection fields.
func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, prior instanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	resp.Diagnostics.Append(validateInstance(data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := resourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance ID", err.Error())
		return
	}
	if err := r.providerData.Client.Update(ctx, id, data.Label.ValueString(), data.State.ValueString()); err != nil {
		resp.Diagnostics.AddError("Cannot update instance", err.Error())
		return
	}
	data.Status, data.SSHHost, data.SSHPort, data.PublicIP, data.MappedPorts = prior.Status, prior.SSHHost, prior.SSHPort, prior.PublicIP, prior.MappedPorts
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete treats an already absent instance as successful destruction.
func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data instanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := resourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance ID", err.Error())
		return
	}
	if err := r.providerData.Client.Delete(ctx, id); err != nil && !errors.Is(err, vastclient.ErrNotFound) {
		resp.Diagnostics.AddError("Cannot destroy instance", err.Error())
	}
}

// ImportState validates a contract ID; Read hydrates the available API attributes.
func (r *instanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, err := resourceID(req.ID); err != nil {
		resp.Diagnostics.AddError("Invalid import ID", err.Error())
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// NewInstance constructs an unconfigured instance resource for provider registration.
func NewInstance() resource.Resource {
	return &instanceResource{}
}
