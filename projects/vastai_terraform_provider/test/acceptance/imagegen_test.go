package acceptance

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// imagegenModulePath locates the actual reusable module source.
var imagegenModulePath = flag.String("imagegen-module", "", "Image generation module runfile")

// imagegenFixture adds template discovery and creation guards to the rental fixture.
type imagegenFixture struct {
	// rental supplies the existing rental lifecycle API.
	rental fixture
	// options supplies template Docker options without executing their contents.
	options string
	// response overrides template metadata for malformed-response cases.
	response string
	// offerID changes the search result to verify existing rentals remain stable.
	offerID int
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
		w.Header().Set("Content-Type", "application/json")
		if f.rental.mode == "empty-offers" {
			fmt.Fprint(w, `{"offers":[]}`)
			return
		}
		id := f.offerID
		if id == 0 {
			id = 17
		}
		fmt.Fprintf(w, `{"offers":[{"id":%d,"gpu_name":"RTX 4090","num_gpus":1,"gpu_ram":24576,"dph_total":0.25}]}`, id)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/api/v0/template/" {
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
		if json.Unmarshal(body, &input) != nil || input["template_hash_id"] != strings.Repeat("a", 32) || input["runtype"] != "ssh_direct" || input["env"] != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		startup, _ := input["onstart"].(string)
		for _, required := range []string{"CF_QUICK_TUNNELS=false", "SUPERVISOR_NO_AUTOSTART=", "cloudflared", "quicktunnel", "FORGE_PORT_LOCAL=17860", "FORGE_ARGS=--server-name___127.0.0.1"} {
			if !strings.Contains(startup, required) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
	}
	f.rental.ServeHTTP(w, r)
}

// writeImagegen configures a consumer of the checked-in module rather than duplicating it.
func writeImagegen(t *testing.T, h *harness, endpoint, hash string) {
	t.Helper()
	module, err := runfiles.Rlocation(*imagegenModulePath)
	if err != nil {
		t.Fatalf("resolve image generation module: %v", err)
	}
	configuration := fmt.Sprintf(`terraform {
  required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" { endpoint = %q }
module "imagegen" {
  source = %q
  template_hash_id = %q
  disk_gb = 20
}

output "forward" { value = module.imagegen.ssh_forward_command }
`, endpoint, filepath.Dir(module), hash)
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
				delete(instance.(map[string]any)["attributes"].(map[string]any), "template_hash_id")
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
	for _, mode := range []string{"empty-offers", "over-budget"} {
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
