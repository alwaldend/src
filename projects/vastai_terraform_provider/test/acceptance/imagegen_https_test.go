package acceptance

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// stagedJSON reads a quoted heredoc without executing the public fixture startup.
func stagedJSON(t *testing.T, startup, path string) []byte {
	t.Helper()
	_, section, found := strings.Cut(startup, "cat > '"+path+"' <<'")
	if !found {
		t.Fatalf("missing staged JSON file %s", path)
	}
	delimiter, section, found := strings.Cut(section, "'\n")
	if !found {
		t.Fatal("missing JSON heredoc delimiter")
	}
	content, _, found := strings.Cut(section, "\n"+delimiter+"\n")
	if !found || !json.Valid([]byte(content)) {
		t.Fatal("staged heredoc is not complete JSON")
	}
	return []byte(content)
}

// TestConfiguredImagegenStartup verifies every public profile without live resources.
func TestConfiguredImagegenStartup(t *testing.T) {
	configurationPath, err := runfiles.Rlocation(*imagegenConfigurationPath)
	if err != nil {
		t.Fatalf("resolve profile configuration: %v", err)
	}
	sequencePath, err := runfiles.Rlocation(*promptSequencePath)
	if err != nil {
		t.Fatalf("resolve configured node: %v", err)
	}
	workflowPath, err := runfiles.Rlocation(*configuredWorkflowPath)
	if err != nil {
		t.Fatalf("resolve configured workflows: %v", err)
	}
	data, err := os.ReadFile(configurationPath)
	if err != nil {
		t.Fatalf("read public profile: %v", err)
	}
	var configuration struct {
		// Downloads owns the public publisher URLs and model destination kinds.
		Downloads map[string]struct {
			// URL identifies the immutable publisher file.
			URL string `json:"url"`
			// Kind selects the native model directory.
			Kind string `json:"kind"`
		} `json:"downloads"`
		// ModelProfiles supplies the bundles without another maintained list.
		ModelProfiles map[string]struct {
			// Workflows identifies the complete native graph catalog.
			Workflows []string `json:"workflows"`
			// StartupWorkflows optionally limits inline graphs; remaining files are uploaded through the API.
			StartupWorkflows *[]string `json:"startup_workflows"`
			// Downloads is the ordered list of filenames resolved from the registry.
			Downloads []string `json:"downloads"`
			// CheckpointDirectory distinguishes SDXL bundles from standalone denoisers.
			CheckpointDirectory string `json:"checkpoint_directory"`
		} `json:"model_profiles"`
	}
	if err := json.Unmarshal(data, &configuration); err != nil {
		t.Fatalf("decode public profile: %v", err)
	}
	gifPath, err := runfiles.Rlocation(*chainedGIFPath)
	if err != nil {
		t.Fatalf("resolve configured GIF node: %v", err)
	}
	var names []string
	for name := range configuration.ModelProfiles {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, profileName := range names {
		t.Run(profileName, func(t *testing.T) {
			profile := configuration.ModelProfiles[profileName]
			startupWorkflows := profile.Workflows
			if profile.StartupWorkflows != nil {
				startupWorkflows = *profile.StartupWorkflows
			}
			f := &imagegenFixture{noTemplate: true, exposed: true, comfyUI: true}
			f.rental.runtype = "ssh_direct"
			server := httptest.NewServer(f)
			defer server.Close()
			h := newHarness(t, server.URL+"/api/v0")
			writeImagegenProfile(t, h, server.URL+"/api/v0", strings.Repeat("a", 32), *comfyuiModulePath)
			path := filepath.Join(h.directory, "main.tf")
			consumer, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read profile consumer: %v", err)
			}
			text := strings.Replace(string(consumer), `template_hash_id = "`+strings.Repeat("a", 32)+`"`, "", 1)
			text += fmt.Sprintf(`
locals {
  configuration = jsondecode(file(%q))
  profile_name = %q
  selected_profile = local.configuration.model_profiles[local.profile_name]
  profile = merge(local.selected_profile, {
    downloads = [for name in local.selected_profile.downloads : merge(local.configuration.downloads[name], {filename=name})]
    custom_nodes = {for name in try(local.selected_profile.custom_nodes, []) : name => local.configuration.custom_nodes[name]}
  })
  configuration_directory = %q
  workflow_directory = %q
}
`, configurationPath, profileName, filepath.Dir(configurationPath), filepath.Dir(filepath.Dir(workflowPath)))
			if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
				t.Fatalf("write profile consumer: %v", err)
			}
			configureImagegenAccess(t, h, fmt.Sprintf(`
access_mode = "exposed"
web_username = local.configuration.web_username
web_password = "fixtureCredentialOnly01234567890123456789"
model_downloads = local.profile.downloads
checkpoint_directory = try(local.profile.checkpoint_directory, "models/diffusion_models")
custom_nodes = try(local.profile.custom_nodes, {})
python_requirements = try(local.profile.python_requirements, null) == null ? null : file("${local.configuration_directory}/${local.profile.python_requirements}")
custom_node_sources = merge(
  {prompt_sequence = file(%q)},
  contains(local.profile.workflows, "${local.profile_name}-image-to-image-gif") ? {chained_img2img_gif=file(%q)} : {},
)
workflows = {for name in try(local.profile.startup_workflows, local.profile.workflows) : name => file("${local.workflow_directory}/${local.profile_name}/${name}.json")}
`, sequencePath, gifPath))
			h.run(t, true, "apply", "-auto-approve", "-no-color")
			if f.startupBytes == 0 || f.startupBytes > 16384 {
				t.Fatal("configured bundle exceeded the API startup limit")
			}
			for _, name := range configuration.ModelProfiles[profileName].Workflows {
				expected, err := os.ReadFile(filepath.Join(filepath.Dir(filepath.Dir(workflowPath)), profileName, name+".json"))
				if err != nil {
					t.Fatalf("read generated workflow: %v", err)
				}
				var actualGraph, expectedGraph any
				if json.Unmarshal(expected, &expectedGraph) != nil {
					t.Fatalf("configured workflow %s is not valid JSON", name)
				}
				if !slices.Contains(startupWorkflows, name) {
					if strings.Contains(f.startup, "cat > '/workspace/ComfyUI/user/default/workflows/"+name+".json'") {
						t.Fatalf("API-staged workflow %s was embedded in startup", name)
					}
					continue
				}
				actual := stagedJSON(t, f.startup, "/workspace/ComfyUI/user/default/workflows/"+name+".json")
				if json.Unmarshal(actual, &actualGraph) != nil || !reflect.DeepEqual(actualGraph, expectedGraph) {
					t.Fatalf("configured workflow %s changed during startup packaging", name)
				}
			}
			var manifest struct {
				// Downloads retains the exact URL/destination pairs sent to the provisioner.
				Downloads []struct {
					// URL selects the publisher file rather than another profile's weights.
					URL string `json:"url"`
					// Destination names the final native ComfyUI model path.
					Destination string `json:"dest"`
				} `json:"downloads"`
			}
			if err := json.Unmarshal(stagedJSON(t, f.startup, "/etc/imagegen-models.json"), &manifest); err != nil {
				t.Fatalf("decode configured model manifest: %v", err)
			}
			if len(manifest.Downloads) != len(profile.Downloads) {
				t.Fatal("configured model asset count changed during packaging")
			}
			for index, filename := range profile.Downloads {
				asset, exists := configuration.Downloads[filename]
				if !exists {
					t.Fatalf("model asset %s is missing from the registry", filename)
				}
				directory := map[string]string{"checkpoint": profile.CheckpointDirectory, "text_encoder": "models/text_encoders", "vae": "models/vae", "ipadapter": "models/ipadapter", "clip_vision": "models/clip_vision", "lora": "models/loras", "model_patch": "models/model_patches"}[asset.Kind]
				if directory == "" && asset.Kind == "checkpoint" {
					directory = "models/diffusion_models"
				}
				if directory == "" || manifest.Downloads[index].URL != asset.URL || manifest.Downloads[index].Destination != "/workspace/ComfyUI/"+directory+"/"+filename {
					t.Fatalf("model asset %s has the wrong publisher or destination", filename)
				}
			}
			h.run(t, true, "destroy", "-auto-approve", "-no-color")
			if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
				evidence, err := json.Marshal(map[string]any{"profile": profileName, "model_files": profile.Downloads, "workflow_count": len(profile.Workflows), "startup_workflow_count": len(startupWorkflows), "startup_bytes": f.startupBytes, "limit_bytes": 16384, "workflow_contents_preserved": true, "model_destinations_preserved": true, "live_operations": 0})
				if err != nil {
					t.Fatalf("encode configured startup evidence: %v", err)
				}
				if err := os.WriteFile(filepath.Join(directory, "configured-imagegen-startup-"+profileName+".json"), evidence, 0o600); err != nil {
					t.Fatalf("write configured startup evidence: %v", err)
				}
			}
		})
	}
}

// configureImagegenAccess supplies public fixture credentials and explicit access configuration.
func configureImagegenAccess(t *testing.T, h *harness, extra string) {
	t.Helper()
	path := filepath.Join(h.directory, "main.tf")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture configuration: %v", err)
	}
	text := strings.Replace(string(data), "disk_gb = 20", "disk_gb = 20\n"+extra, 1)
	text += "\noutput \"url\" { value = module.imagegen.https_url }\n"
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatalf("write access configuration: %v", err)
	}
}

// TestImagegenHTTPS exercises exposed access, credential redaction, and mode replacement through Terraform.
func TestImagegenHTTPS(t *testing.T) {
	f := &imagegenFixture{exposed: true}
	f.rental.runtype = "ssh_direct"
	server := httptest.NewServer(f)
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
	const password = "fixtureCredentialOnly01234567890123456789"
	configureImagegenAccess(t, h, `access_mode = "exposed"
  web_username = "imagegen"
  web_password = "`+password+`"`)
	output := h.run(t, true, "apply", "-auto-approve", "-no-color")
	if strings.Contains(output, password) {
		t.Fatal("password appeared in normal Terraform output")
	}
	h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
	if url := h.run(t, true, "output", "-raw", "url"); url != "https://203.0.113.2:27860" {
		t.Fatalf("incorrect mapped HTTPS URL: %s", url)
	}
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
	if output := h.run(t, true, "plan", "-no-color"); !strings.Contains(output, "must be replaced") {
		t.Fatal("access mode change did not replace rental")
	}
	// Restore the original configuration so destruction checks the same access contract.
	configureImagegenAccess(t, h, `access_mode = "exposed"
  web_username = "imagegen"
  web_password = "`+password+`"`)
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		evidence, err := json.Marshal(map[string]any{"result": "passed", "mode": "exposed", "live_operations": 0, "credential_redacted": true})
		if err != nil {
			t.Fatalf("encode HTTPS evidence: %v", err)
		}
		if err := os.WriteFile(filepath.Join(directory, "imagegen-https.json"), evidence, 0o600); err != nil {
			t.Fatalf("write HTTPS evidence: %v", err)
		}
	}
}

// TestImagegenAccessRejection verifies invalid access and memory inputs fail before rental creation.
func TestImagegenAccessRejection(t *testing.T) {
	for name, extra := range map[string]string{
		"memory":              `min_gpu_ram_gb = 24`,
		"unknown-mode":        `access_mode = "public"`,
		"missing-credentials": `access_mode = "exposed"`,
		"weak-password": `access_mode = "exposed"
web_username = "imagegen"
web_password = "password"`,
		"shell-injection": `access_mode = "exposed"
web_username = "invalid;name"
web_password = "01234567890123456789012345678901"`,
	} {
		t.Run(name, func(t *testing.T) {
			f := &imagegenFixture{}
			server := httptest.NewServer(f)
			defer server.Close()
			h := newHarness(t, server.URL+"/api/v0")
			writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
			configureImagegenAccess(t, h, extra)
			h.run(t, false, "plan", "-no-color")
			if f.rental.exists {
				t.Fatal("invalid access created a rental")
			}
		})
	}
}

// TestImagegenWithoutTemplate verifies the pinned-runtime consumer creates no template lookup or inherited options.
func TestImagegenWithoutTemplate(t *testing.T) {
	f := &imagegenFixture{noTemplate: true}
	f.rental.runtype = "ssh_direct"
	server := httptest.NewServer(f)
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
	path := filepath.Join(h.directory, "main.tf")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read default-runtime configuration: %v", err)
	}
	text := strings.Replace(string(data), `template_hash_id = "`+strings.Repeat("a", 32)+`"`, "", 1)
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatalf("write default-runtime configuration: %v", err)
	}
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
}

// TestImagegenModelManifest exercises model-path configuration through a real Terraform apply.
func TestImagegenModelManifest(t *testing.T) {
	f := &imagegenFixture{}
	f.rental.runtype = "ssh_direct"
	server := httptest.NewServer(f)
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
	configureImagegenAccess(t, h, `model_downloads = [
  {url="https://example.invalid/model.safetensors",filename="model.safetensors",kind="checkpoint"},
  {url="https://example.invalid/encoder.safetensors",filename="encoder.safetensors",kind="text_encoder"},
  {url="https://example.invalid/vae.safetensors",filename="vae.safetensors",kind="vae"}
]
model_selection = {
  checkpoint = "model.safetensors"
  preset = "flux"
  modules = ["encoder.safetensors", "vae.safetensors"]
  steps = 40
  cfg_scale = 3
}`)
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	manifest := stagedJSON(t, f.startup, "/etc/imagegen-models.json")
	var decoded struct {
		// Downloads contains the staged model destinations.
		Downloads []struct {
			// URL identifies the public fixture file.
			URL string `json:"url"`
			// Destination is the staged Forge model path.
			Destination string `json:"dest"`
		} `json:"downloads"`
	}
	if err := json.Unmarshal(manifest, &decoded); err != nil {
		t.Fatalf("decode model configuration: %v", err)
	}
	if len(decoded.Downloads) != 3 || decoded.Downloads[0].Destination != "/workspace/stable-diffusion-webui-forge/models/Stable-diffusion/model.safetensors" || decoded.Downloads[1].Destination != "/workspace/stable-diffusion-webui-forge/models/text_encoder/encoder.safetensors" || decoded.Downloads[2].Destination != "/workspace/stable-diffusion-webui-forge/models/VAE/vae.safetensors" {
		t.Fatalf("wrong model destination: %s", manifest)
	}
	if !strings.Contains(f.startup, "--ui-settings-file /etc/imagegen-forge.json") {
		t.Fatal("Forge model selection was not connected to startup")
	}
	settings := stagedJSON(t, f.startup, "/etc/imagegen-forge.json")
	var selected map[string]any
	if err := json.Unmarshal(settings, &selected); err != nil {
		t.Fatalf("decode Forge settings: %v", err)
	}
	for _, key := range []string{"sd_model_checkpoint", "forge_checkpoint_flux"} {
		if selected[key] != "model.safetensors" {
			t.Fatalf("checkpoint missing from %s", key)
		}
	}
	for _, key := range []string{"forge_additional_modules", "forge_additional_modules_flux"} {
		modules, ok := selected[key].([]any)
		if !ok || len(modules) != 2 || modules[0] != decoded.Downloads[1].Destination || modules[1] != decoded.Downloads[2].Destination {
			t.Fatalf("external modules missing from %s", key)
		}
	}
	if selected["forge_preset"] != "flux" || selected["flux_t2i_step"] != float64(40) || selected["flux_i2i_step"] != float64(40) || selected["flux_t2i_cfg"] != float64(3) || selected["flux_i2i_cfg"] != float64(3) {
		t.Fatal("Forge preset defaults were not staged")
	}
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
}

// TestComfyUIModule verifies both access modes and native model paths through Terraform.
func TestComfyUIModule(t *testing.T) {
	for _, checkpointDirectory := range []string{"models/diffusion_models", "models/checkpoints"} {
		for _, exposed := range []bool{false, true} {
			t.Run(checkpointDirectory+"/"+map[bool]string{false: "ssh", true: "https"}[exposed], func(t *testing.T) {
				f := &imagegenFixture{noTemplate: true, exposed: exposed, comfyUI: true}
				f.rental.runtype = "ssh_direct"
				server := httptest.NewServer(f)
				defer server.Close()
				h := newHarness(t, server.URL+"/api/v0")
				writeImagegenProfile(t, h, server.URL+"/api/v0", strings.Repeat("a", 32), *comfyuiModulePath)
				path := filepath.Join(h.directory, "main.tf")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				configuration := strings.Replace(string(data), `template_hash_id = "`+strings.Repeat("a", 32)+`"`, "", 1)
				if err := os.WriteFile(path, []byte(configuration), 0o600); err != nil {
					t.Fatal(err)
				}
				extra := `checkpoint_directory = "` + checkpointDirectory + `"
model_downloads = [
 {url="https://example.invalid/flash.safetensors",filename="flash.safetensors",kind="checkpoint"},
 {url="https://example.invalid/t5.safetensors",filename="t5.safetensors",kind="text_encoder"},
 {url="https://example.invalid/ae.safetensors",filename="ae.safetensors",kind="vae"},
 {url="https://example.invalid/adapter.safetensors",filename="adapter.safetensors",kind="ipadapter"},
 {url="https://example.invalid/vision.safetensors",filename="vision.safetensors",kind="clip_vision"}
]
workflows = { fixture = jsonencode({version=0.4,nodes=[],links=[]}) }
custom_nodes = { FixtureNodes = {url="https://github.com/example/fixture",ref="aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"} }
python_requirements = "fixture-package==1.0 --hash=sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n${join("",[for i in range(1000):"# Public fixture dependency commentary.\n"])}"`
				if exposed {
					extra += "\naccess_mode = \"exposed\"\nweb_password = \"fixtureCredentialOnly01234567890123456789\""
				}
				configureImagegenAccess(t, h, extra)
				h.run(t, true, "apply", "-auto-approve", "-no-color")
				h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
				if command := h.run(t, true, "output", "-raw", "forward"); !strings.Contains(command, ":7860:127.0.0.1:18188") {
					t.Fatalf("ComfyUI forwarding targets wrong backend: %s", command)
				}
				manifest := stagedJSON(t, f.startup, "/etc/imagegen-models.json")
				for _, destination := range []string{checkpointDirectory + "/flash.safetensors", "models/text_encoders/t5.safetensors", "models/vae/ae.safetensors", "models/ipadapter/adapter.safetensors", "models/clip_vision/vision.safetensors"} {
					if !strings.Contains(string(manifest), "/workspace/ComfyUI/"+destination) {
						t.Fatalf("wrong ComfyUI destination: %s", manifest)
					}
				}
				for _, required := range []string{"--disable-all-custom-nodes", "--offline", "localhost:8288:18288:/:API Wrapper", "/etc/supervisor/conf.d/comfyui.conf"} {
					if !strings.Contains(f.startup, required) {
						t.Fatalf("missing ComfyUI restriction %s", required)
					}
				}
				if !strings.Contains(f.startup, "--whitelist-custom-nodes FixtureNodes") {
					t.Fatal("approved node pack was not selectively enabled")
				}
				var provisioning struct {
					// GitRepos contains the pinned native fetcher configuration.
					GitRepos []struct {
						// Ref fixes the reviewed node-pack revision.
						Ref string `json:"ref"`
						// Dest installs nodes under the native custom-node directory.
						Dest string `json:"dest"`
					} `json:"git_repos"`
					// PipPackages retains the requirements and integrity policy.
					PipPackages []struct {
						// Args requires hashes without unconstrained dependency resolution.
						Args string `json:"args"`
					} `json:"pip_packages"`
				}
				if err := json.Unmarshal(manifest, &provisioning); err != nil || len(provisioning.GitRepos) != 1 || provisioning.GitRepos[0].Ref != strings.Repeat("a", 40) || provisioning.GitRepos[0].Dest != "/workspace/ComfyUI/custom_nodes/FixtureNodes" || len(provisioning.PipPackages) != 1 || provisioning.PipPackages[0].Args != "--require-hashes --no-deps" {
					t.Fatal("native fetcher did not retain pinned nodes and locked dependencies")
				}
				if strings.Contains(f.startup, "for service in tunnel_manager syncthing pyworker api-wrapper") {
					t.Fatal("local API wrapper was unnecessarily disabled")
				}
				if strings.Contains(f.startup, "--enable-cors-header") || strings.Contains(f.startup, "--enable-manager") {
					t.Fatal("unrestricted runtime option retained")
				}
				if exposed && !strings.Contains(f.startup, "reverse_proxy 127.0.0.1:18188") {
					t.Fatal("proxy targets wrong backend")
				}

				_, staged, found := strings.Cut(f.startup, "mkdir -p \"$(dirname '/workspace/ComfyUI/user/default/workflows/fixture.json')\"")
				if !found {
					t.Fatal("workflow was not staged after native workspace initialization")
				}
				stagedJSON(t, staged, "/workspace/ComfyUI/user/default/workflows/fixture.json")

				h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
				h.run(t, true, "destroy", "-auto-approve", "-no-color")
				if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
					evidence, err := json.Marshal(map[string]any{"runtime": "comfyui", "exposed": exposed, "native_model_paths": true, "checkpoint_directory": checkpointDirectory, "live_operations": 0})
					if err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(directory, "comfyui-"+filepath.Base(checkpointDirectory)+"-"+map[bool]string{false: "ssh", true: "https"}[exposed]+".json"), evidence, 0o600); err != nil {
						t.Fatal(err)
					}
				}
			})
		}
	}
}

// TestRentalRejectionDiagnostic retains the public startup limit without echoing a response body.
func TestRentalRejectionDiagnostic(t *testing.T) {
	f := &imagegenFixture{}
	f.rental.mode = "startup-too-long"
	server := httptest.NewServer(f)
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
	output := h.run(t, false, "apply", "-auto-approve", "-no-color")
	output = strings.Join(strings.Fields(output), " ")
	if !strings.Contains(output, "startup arguments exceed the server limit of 32768") || strings.Contains(output, "fixture-private-content") || f.rental.exists {
		t.Fatal("rental rejection was not classified safely")
	}
}

// TestImagegenModelRejection rejects malformed model declarations before any paid request.
func TestImagegenModelRejection(t *testing.T) {
	for _, models := range []string{
		`[{url="https://example.invalid/a",filename="../model.safetensors"}]`,
		`[{url="https://user:password@example.invalid/a",filename="model.safetensors"}]`,
		`[{url="https://example.invalid/a",filename="model.safetensors",kind="unknown"}]`,
		`[{url="https://example.invalid/a",filename="model.safetensors"},{url="https://example.invalid/b",filename="model.safetensors"}]`,
	} {
		t.Run(models, func(t *testing.T) {
			f := &imagegenFixture{}
			server := httptest.NewServer(f)
			defer server.Close()
			h := newHarness(t, server.URL+"/api/v0")
			writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
			configureImagegenAccess(t, h, "model_downloads = "+models)
			h.run(t, false, "plan", "-no-color")
			if f.rental.exists {
				t.Fatal("invalid manifest created rental")
			}
		})
	}
}
