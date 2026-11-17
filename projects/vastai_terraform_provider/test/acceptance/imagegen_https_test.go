package acceptance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// TestComfyUIModule verifies both access modes with an Ansible-managed installation.
func TestComfyUIModule(t *testing.T) {
	for _, exposed := range []bool{false, true} {
		t.Run(fmt.Sprint(exposed), func(t *testing.T) {
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
			data = []byte(strings.Replace(string(data), `template_hash_id = "`+strings.Repeat("a", 32)+`"`, "", 1))
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}
			extra := `workspace_directory = "/srv/imagegen"`
			if exposed {
				extra += "\naccess_mode = \"exposed\"\nweb_password = \"fixtureCredentialOnly01234567890123456789\""
			}
			configureImagegenAccess(t, h, extra)
			h.run(t, true, "apply", "-auto-approve", "-no-color")
			h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
			if command := h.run(t, true, "output", "-raw", "forward"); !strings.Contains(command, ":7860:127.0.0.1:18188") {
				t.Fatal("wrong ComfyUI forwarding backend")
			}
			for _, required := range []string{"--disable-all-custom-nodes", "--offline", "localhost:8288:18288:/:API Wrapper", "/etc/supervisor/conf.d/comfyui.conf", "export WORKSPACE='/srv/imagegen'"} {
				if !strings.Contains(f.startup, required) {
					t.Fatalf("missing runtime restriction %s", required)
				}
			}
			if exposed && !strings.Contains(f.startup, "reverse_proxy 127.0.0.1:18188") {
				t.Fatal("wrong HTTPS backend")
			}
			if strings.Contains(f.startup, "--enable-manager") || strings.Contains(f.startup, "--enable-cors-header") {
				t.Fatal("unrestricted runtime arguments")
			}
			verifyManagedArguments(t, f.startup)

			var manifest map[string]any
			if err := json.Unmarshal(stagedJSON(t, f.startup, "/etc/imagegen-models.json"), &manifest); err != nil {
				t.Fatal(err)
			}
			for _, key := range []string{"downloads", "git_repos", "pip_packages", "write_files"} {
				if entries, ok := manifest[key].([]any); !ok || len(entries) != 0 {
					t.Fatalf("ComfyUI bootstrap still installs %s", key)
				}
			}
			if strings.Contains(f.startup, "/user/default/workflows/") {
				t.Fatal("workflow payload retained in rental startup")
			}

			h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
			h.run(t, true, "destroy", "-auto-approve", "-no-color")
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
		`[{url="https://user:password@example.invalid/a",filename="model.safetensors"}]`, // trufflehog:ignore synthetic credential tests URL rejection.
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

// verifyManagedArguments executes the rendered boot guard against disposable environment files.
func verifyManagedArguments(t *testing.T, startup string) {
	t.Helper()
	_, guard, found := strings.Cut(startup, "# Ansible owns reviewed runtime arguments after application installation.\n")
	if !found {
		t.Fatal("missing runtime argument ownership guard")
	}
	guard, _, found = strings.Cut(guard, "export PORTAL_CONFIG=")
	if !found {
		t.Fatal("missing runtime guard boundary")
	}
	for _, managed := range []bool{false, true} {
		directory := t.TempDir()
		environment := filepath.Join(directory, "environment")
		initial := "export COMFYUI_ARGS='--listen 127.0.0.1 --offline --whitelist-custom-nodes FixtureNodes'"
		if managed {
			initial += " # imagegen-ansible"
		}
		initial += "\n"
		if err := os.WriteFile(environment, []byte(initial), 0o600); err != nil {
			t.Fatal(err)
		}
		script := strings.ReplaceAll(guard, "/etc/environment", environment)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		command := exec.CommandContext(ctx, "/bin/bash", "-eu", "-c", script)
		output, err := command.CombinedOutput()
		cancel()
		if err != nil {
			t.Fatalf("execute public boot guard: %v: %s", err, output)
		}
		actual, err := os.ReadFile(environment)
		if err != nil {
			t.Fatal(err)
		}
		if managed && string(actual) != initial {
			t.Fatal("boot overwrote Ansible runtime arguments")
		}
		if !managed && !strings.Contains(string(actual), "--disable-all-custom-nodes --offline") {
			t.Fatal("unmanaged image arguments bypassed restrictive defaults")
		}
	}
}
