package acceptance

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// imagegenModulePath locates the actual reusable module source.
var imagegenModulePath = flag.String("imagegen-module", "", "Image generation module runfile")

// comfyuiModulePath locates the native ComfyUI wrapper.
var comfyuiModulePath = flag.String("comfyui-module", "", "ComfyUI module runfile")

// imagegenConfigurationPath selects the public deployment profile for full-bundle acceptance.
var imagegenConfigurationPath = flag.String("imagegen-configuration", "", "Image generation configuration runfile")

// promptSequencePath selects the actual first-party node staged in the startup manifest.
var promptSequencePath = flag.String("prompt-sequence", "", "Prompt Sequence source runfile")

// chainedGIFPath selects the first-party GIF source for profiles that stage it.
var chainedGIFPath = flag.String("chained-gif", "", "Chained img2img GIF source runfile")

// configuredWorkflowPath anchors the generated workflow tree consumed by the real profile.
var configuredWorkflowPath = flag.String("configured-workflow", "", "Configured workflow runfile")

// imagegenFixture adds template discovery and creation guards to the rental fixture.
type imagegenFixture struct {
	// rental supplies the existing rental lifecycle API.
	rental fixture
	// options supplies template Docker options without executing their contents.
	options string
	// response overrides template metadata for malformed-response cases.
	response string
	// startup retains the public-fixture startup for manifest inspection.
	startup string
	// startupBytes records the encoded API payload size without retaining it in evidence.
	startupBytes int
	// noTemplate requires a rental without inherited template configuration.
	noTemplate bool
	// exposed selects the single-port HTTPS request contract.
	exposed bool
	// offerID changes the search result to verify existing rentals remain stable.
	offerID int

	// comfyUI selects the native ComfyUI service and backend contract.
	comfyUI bool
}

// ServeHTTP rejects incorrect template selection or networking before accepting a fixture rental.
func (f *imagegenFixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/api/v0/bundles/" {
		var query map[string]any
		if json.NewDecoder(r.Body).Decode(&query) != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, name := range []string{"datacenter", "verified", "rentable"} {
			filter, ok := query[name].(map[string]any)
			if !ok || filter["eq"] != true {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		memory, ok := query["gpu_ram"].(map[string]any)
		if !ok || memory["gte"] != float64(32768) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if f.rental.mode == "empty-offers" {
			fmt.Fprint(w, `{"offers":[]}`)
			return
		}
		id := f.offerID
		if id == 0 {
			id = 17
		}
		tier := `,"hosting_type":1,"verification":"verified"`
		if f.rental.mode == "uncertified" {
			tier = `,"hosting_type":0,"verification":"verified"`
		}
		if f.rental.mode == "missing-tier" {
			tier = ""
		}
		fmt.Fprintf(w, `{"offers":[{"id":%d,"gpu_name":"RTX 5090","num_gpus":1,"gpu_ram":32768,"dph_total":0.25%s}]}`, id, tier)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/api/v0/template/" {
		if f.noTemplate {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var filters map[string]map[string]string
		if r.Header.Get("Authorization") != "Bearer fixture-credential" || json.Unmarshal([]byte(r.URL.Query().Get("select_filters")), &filters) != nil || len(filters["hash_id"]["eq"]) != 32 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if f.response != "" {
			fmt.Fprint(w, f.response)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "templates_found": 1, "templates": []any{map[string]any{
			"hash_id": filters["hash_id"]["eq"], "image": "fixture/image", "tag": "1", "env": f.options,
		}}})
		return
	}
	if r.Method == http.MethodPut && r.URL.Path == "/api/v0/asks/17/" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var input map[string]any
		if json.Unmarshal(body, &input) != nil || (f.noTemplate && input["template_hash_id"] != nil) || (!f.noTemplate && input["template_hash_id"] != strings.Repeat("a", 32)) || input["runtype"] != "ssh_direct" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		startup, _ := input["onstart"].(string)
		f.startupBytes = len(startup)
		if f.rental.mode == "startup-too-long" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":"invalid_args","msg":"Invalid args: len(args) > 32768; fixture-private-content"}`)
			return
		}
		if strings.HasPrefix(startup, "# Compressed Vast.ai startup") {
			_, encoded, found := strings.Cut(startup, "printf '%s' '")
			if !found || len(startup) > 16384 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			encoded, _, _ = strings.Cut(encoded, "'")
			compressed, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			reader, err := gzip.NewReader(bytes.NewReader(compressed))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			decoded, err := io.ReadAll(reader)
			closeErr := reader.Close()
			if err != nil || closeErr != nil || !strings.Contains(startup, `/bin/bash "$startup_file"`) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			startup = string(decoded)
		}
		f.startup = startup
		options, _ := input["env"].(string)
		if f.exposed {
			if options != "-p 7860:7860" || !strings.Contains(startup, "basicauth") || !strings.Contains(startup, "tls /etc/instance.crt /etc/instance.key") || !strings.Contains(startup, "admin off") || strings.Contains(startup, "OPEN_BUTTON_TOKEN") {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if options != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		arguments := "--port 17860 --server-name 127.0.0.1"
		if f.comfyUI {
			arguments = "--listen 127.0.0.1 --port 18188"
		}
		for _, required := range []string{"CF_QUICK_TUNNELS=false", "INSTANCE_TEST=false", "--no-update-portal", "--no-update-vast", "unset CF_TUNNEL_TOKEN HOTFIX_SCRIPT BOOT_SCRIPT", arguments, "localhost:1111:11111:/:Instance Portal", "localhost:8080:18080:/:Jupyter", "localhost:6006:16006:/:Tensorboard"} {
			if !strings.Contains(startup, required) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		f.rental.image, _ = input["image"].(string)
		r.Body = io.NopCloser(bytes.NewReader(body))
	}
	f.rental.ServeHTTP(w, r)
}

// TestInstanceStartupCompression executes the API command and verifies refresh and drift through Terraform.
func TestInstanceStartupCompression(t *testing.T) {
	f := &fixture{}
	server := httptest.NewServer(f)
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	script := strings.Repeat("# Public neutral startup fixture.\n", 4000) + "printf '%s' 'literal $HOME and apostrophe: '\"'\"'' > result.txt\n"
	if err := os.WriteFile(filepath.Join(h.directory, "startup.sh"), []byte(script), 0o600); err != nil {
		t.Fatalf("write startup fixture: %v", err)
	}
	configuration, err := os.ReadFile(filepath.Join(h.directory, "main.tf"))
	if err != nil {
		t.Fatalf("read startup configuration: %v", err)
	}
	configuration = []byte(strings.Replace(string(configuration), "offer_id = 17", "offer_id = 17\n  onstart = file(\"${path.module}/startup.sh\")\n  compress_onstart = true", 1))
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), configuration, 0o600); err != nil {
		t.Fatalf("write startup configuration: %v", err)
	}
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	f.mu.Lock()
	command := f.onstart
	f.mu.Unlock()
	if len(command) > 16384 || len(command) >= len(script) {
		t.Fatal("API command did not fit the startup limit")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	process := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	process.Dir, process.Env = h.directory, h.environment
	if output, err := process.CombinedOutput(); err != nil {
		t.Fatalf("execute startup command: %v\n%s", err, output)
	}
	result, err := os.ReadFile(filepath.Join(h.directory, "result.txt"))
	if err != nil || string(result) != "literal $HOME and apostrophe: '" {
		t.Fatalf("startup literal content changed: %v", err)
	}
	h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	f.mu.Lock()
	f.onstart = "printf 'external change'"
	f.mu.Unlock()
	if output := h.run(t, true, "plan", "-no-color"); !strings.Contains(output, "must be replaced") {
		t.Fatal("external startup drift was not detected")
	}
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		evidence, err := json.Marshal(map[string]any{"original_bytes": len(script), "command_bytes": len(command), "execution": "passed", "refresh": "stable", "drift": "detected", "live_operations": 0})
		if err != nil {
			t.Fatalf("encode startup evidence: %v", err)
		}
		if err := os.WriteFile(filepath.Join(directory, "startup-compression.json"), evidence, 0o600); err != nil {
			t.Fatalf("write startup evidence: %v", err)
		}
	}
}

// TestInstanceStartupLimit rejects oversized plain and incompressible commands before creation.
func TestInstanceStartupLimit(t *testing.T) {
	var payload strings.Builder
	for index := range 2048 {
		fmt.Fprintf(&payload, "%x", sha256.Sum256([]byte(fmt.Sprint(index))))
	}
	for _, compress := range []bool{false, true} {
		t.Run(fmt.Sprint(compress), func(t *testing.T) {
			f := &fixture{}
			server := httptest.NewServer(f)
			defer server.Close()
			h := newHarness(t, server.URL+"/api/v0")
			if err := os.WriteFile(filepath.Join(h.directory, "startup.sh"), []byte(payload.String()), 0o600); err != nil {
				t.Fatalf("write oversized fixture: %v", err)
			}
			configuration, err := os.ReadFile(filepath.Join(h.directory, "main.tf"))
			if err != nil {
				t.Fatalf("read oversized configuration: %v", err)
			}
			configuration = []byte(strings.Replace(string(configuration), "offer_id = 17", fmt.Sprintf("offer_id = 17\n  onstart = file(\"${path.module}/startup.sh\")\n  compress_onstart = %t", compress), 1))
			if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), configuration, 0o600); err != nil {
				t.Fatalf("write oversized configuration: %v", err)
			}
			output := h.run(t, false, "apply", "-auto-approve", "-no-color")
			if !strings.Contains(output, "startup command exceeds 16384 bytes") || strings.Contains(output, payload.String()[:80]) || f.exists {
				t.Fatal("oversized startup did not fail safely before creation")
			}
		})
	}
}

// writeImagegen configures a consumer of the checked-in module rather than duplicating it.
func writeImagegen(t *testing.T, h *harness, endpoint, hash string) {
	t.Helper()
	writeImagegenProfile(t, h, endpoint, hash, *imagegenModulePath)
}

// writeImagegenProfile supplies the selected real module to a fixture consumer.
func writeImagegenProfile(t *testing.T, h *harness, endpoint, hash, modulePath string) {
	t.Helper()
	module, err := runfiles.Rlocation(modulePath)
	if err != nil {
		t.Fatalf("resolve image generation module: %v", err)
	}
	// Use a local package so sibling module sources retain their normal boundaries.
	if err := os.Symlink(filepath.Dir(filepath.Dir(module)), filepath.Join(h.directory, "modules")); err != nil && !os.IsExist(err) {
		t.Fatalf("link image generation modules: %v", err)
	}
	configuration := fmt.Sprintf(`terraform {
  required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" { endpoint = %q }
module "imagegen" {
  source = %q
  image = "fixture/image@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  template_hash_id = %q
  disk_gb = 20
}

output "forward" { value = module.imagegen.ssh_forward_command }
`, endpoint, "./modules/"+filepath.Base(filepath.Dir(module)), hash)
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(configuration), 0o600); err != nil {
		t.Fatalf("write module consumer: %v", err)
	}
	h.run(t, true, "get", "-no-color")
}

// TestImagegenModule verifies template selection, deferred connection output, and replacement without live capacity.
func TestImagegenModule(t *testing.T) {
	f := &imagegenFixture{options: `-e PROVISIONING_SCRIPT="https://example.invalid/model setup.sh" --env=AUTO_UPDATE=false`}
	f.rental.runtype = "ssh_direct"
	server := httptest.NewServer(f)
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
	h.run(t, true, "validate", "-no-color")
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	if output := h.run(t, true, "output", "-no-color"); strings.Contains(output, "ssh -N") {
		t.Fatalf("connection was invented before refresh: %s", output)
	}
	h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
	command := h.run(t, true, "output", "-raw", "forward")
	if command != "ssh -N -o ExitOnForwardFailure=yes -L 127.0.0.1:7860:127.0.0.1:17860 -p 22001 root@'203.0.113.1'" {
		t.Fatalf("incorrect forwarding command: %s", command)
	}
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	f.offerID = 18
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	f.offerID = 17
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("b", 32))
	if output := h.run(t, true, "plan", "-no-color"); !strings.Contains(output, "must be replaced") {
		t.Fatalf("template change did not replace the rental: %s", output)
	}
	writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.Marshal(map[string]any{"result": "passed", "forward_command": command, "rental_api": "loopback fixture", "operations": f.rental.calls, "live_operations": 0})
		if err != nil {
			t.Fatalf("encode image generation evidence: %v", err)
		}
		if err := os.WriteFile(filepath.Join(output, "imagegen-lifecycle.json"), data, 0o600); err != nil {
			t.Fatalf("write image generation evidence: %v", err)
		}
	}
}

// TestTemplateSchemaUpgrade verifies adding the optional selector does not replace preexisting image-only rentals.
func TestTemplateSchemaUpgrade(t *testing.T) {
	f := &fixture{}
	server := httptest.NewServer(f)
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
	path := filepath.Join(h.directory, "terraform.tfstate")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read legacy fixture state: %v", err)
	}
	var state map[string]any
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("decode legacy fixture state: %v", err)
	}
	for _, resource := range state["resources"].([]any) {
		object := resource.(map[string]any)
		if object["type"] == "vastai_instance" {
			for _, instance := range object["instances"].([]any) {
				for _, name := range []string{"template_hash_id", "docker_options", "mapped_ports"} {
					delete(instance.(map[string]any)["attributes"].(map[string]any), name)
				}
			}
		}
	}
	data, err = json.Marshal(state)
	if err != nil {
		t.Fatalf("encode legacy fixture state: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write legacy fixture state: %v", err)
	}
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
}

// TestImagegenTemplateRejection checks unsafe options and incomplete metadata before creation.
func TestImagegenTemplateRejection(t *testing.T) {
	for _, options := range []string{`-p 7860:7860`, `-p7860:7860`, `--publish=7860:7860`, `-P`, `--network=host`, `-e OK=value -p 8080:8080`, `-e "UNTERMINATED=value`, `--unknown value`, `-e MISSING_ASSIGNMENT`} {
		t.Run(options, func(t *testing.T) {
			f := &imagegenFixture{options: options}
			server := httptest.NewServer(f)
			defer server.Close()
			h := newHarness(t, server.URL+"/api/v0")
			writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
			if output := h.run(t, false, "plan", "-no-color"); !strings.Contains(output, "Template Docker options must contain only environment assignments") {
				t.Fatalf("missing safety diagnostic: %s", output)
			}
			if f.rental.exists {
				t.Fatal("unsafe template created a rental")
			}
		})
	}
	for _, response := range []string{`{"success":true,"templates":[]}`, `{"success":true}`, `{"success":true,"templates":[{"hash_id":"wrong","image":"fixture/image","env":""}]}`, `{"success":true,"templates":[{},{}]}`, `{"success":true,"templates":[{"hash_id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","image":"fixture/image"}]}`} {
		t.Run(response, func(t *testing.T) {
			f := &imagegenFixture{response: response}
			server := httptest.NewServer(f)
			defer server.Close()
			h := newHarness(t, server.URL+"/api/v0")
			writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
			h.run(t, false, "plan", "-no-color")
		})
	}
}

// TestImagegenOfferDefaults verifies empty and over-budget markets fail without renting.
func TestImagegenOfferDefaults(t *testing.T) {
	for _, mode := range []string{"empty-offers", "over-budget", "uncertified", "missing-tier"} {
		t.Run(mode, func(t *testing.T) {
			f := &imagegenFixture{}
			f.rental.mode = mode
			server := httptest.NewServer(f)
			defer server.Close()
			h := newHarness(t, server.URL+"/api/v0")
			writeImagegen(t, h, server.URL+"/api/v0", strings.Repeat("a", 32))
			if mode == "over-budget" {
				path := filepath.Join(h.directory, "main.tf")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(strings.Replace(string(data), "disk_gb = 20", "disk_gb = 20\n  max_hourly_price = 0.10", 1)), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			h.run(t, false, "plan", "-no-color")
			if f.rental.exists {
				t.Fatal("unavailable or over-budget market created a rental")
			}
		})
	}
}
