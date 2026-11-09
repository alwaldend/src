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
	configureImagegenAccess(t, h, `model_downloads = [{url="https://example.invalid/model.safetensors",filename="model.safetensors",kind="checkpoint"}]`)
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
	if len(decoded.Downloads) != 1 || decoded.Downloads[0].Destination != "/workspace/stable-diffusion-webui-forge/models/Stable-diffusion/model.safetensors" {
		t.Fatalf("wrong model destination: %s", manifest)
	}
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
}
