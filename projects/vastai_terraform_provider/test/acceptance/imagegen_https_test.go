package acceptance

import (
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	_, encoded, found := strings.Cut(f.startup, "printf '%s' '")
	if !found {
		t.Fatal("model manifest was not staged")
	}
	encoded, _, _ = strings.Cut(encoded, "'")
	manifest, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode staged manifest: %v", err)
	}
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
	_, settingsScript, found := strings.Cut(f.startup, "> /etc/imagegen-models.json")
	if !found {
		t.Fatal("model manifest destination was not staged")
	}
	_, settingsEncoded, found := strings.Cut(settingsScript, "printf '%s' '")
	if !found || !strings.Contains(f.startup, "--ui-settings-file /etc/imagegen-forge.json") {
		t.Fatal("Forge model selection was not connected to startup")
	}
	settingsEncoded, _, _ = strings.Cut(settingsEncoded, "'")
	settings, err := base64.StdEncoding.DecodeString(settingsEncoded)
	if err != nil {
		t.Fatalf("decode staged Forge settings: %v", err)
	}
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
	for _, exposed := range []bool{false, true} {
		t.Run(map[bool]string{false: "ssh", true: "https"}[exposed], func(t *testing.T) {
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
			extra := `model_downloads = [
 {url="https://example.invalid/flash.safetensors",filename="flash.safetensors",kind="checkpoint"},
 {url="https://example.invalid/t5.safetensors",filename="t5.safetensors",kind="text_encoder"},
 {url="https://example.invalid/ae.safetensors",filename="ae.safetensors",kind="vae"}
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
			_, encoded, found := strings.Cut(f.startup, "printf '%s' '")
			if !found {
				t.Fatal("model manifest missing")
			}
			encoded, _, _ = strings.Cut(encoded, "'")
			manifest, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				t.Fatal(err)
			}
			for _, destination := range []string{"models/diffusion_models/flash.safetensors", "models/text_encoders/t5.safetensors", "models/vae/ae.safetensors"} {
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
			if !found || !strings.Contains(staged, "| base64 -d > '/workspace/ComfyUI/user/default/workflows/fixture.json'") {
				t.Fatal("workflow was not staged after native workspace initialization")
			}

			h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
			h.run(t, true, "destroy", "-auto-approve", "-no-color")
			if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
				evidence, err := json.Marshal(map[string]any{"runtime": "comfyui", "exposed": exposed, "native_model_paths": true, "live_operations": 0})
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(directory, "comfyui-"+map[bool]string{false: "ssh", true: "https"}[exposed]+".json"), evidence, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
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
