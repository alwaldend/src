package acceptance

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// providerPath names the provider binary in the Bazel runfiles tree.
var providerPath = flag.String("provider", "", "Provider runfile")

// terraformPath names the Terraform CLI binary in the Bazel runfiles tree.
var terraformPath = flag.String("terraform", "", "Terraform runfile")

// fixture implements the documented API shapes and records lifecycle operations.
type fixture struct {
	// mu protects state accessed by the HTTP server and test process.
	mu sync.Mutex
	// exists indicates whether the rental exists.
	exists bool
	// label is the current instance label.
	label string
	// state is the requested power state.
	state string
	// runtype overrides the API launch-mode representation for normalization fixtures.
	runtype string
	// image overrides the default image for pinned-runtime fixtures.
	image string
	// mode selects a specific API failure.
	mode string
	// calls records non-sensitive HTTP methods and paths.
	calls []string
	// onstart is the fixture's creation command, echoed to exercise refresh.
	onstart string
}

// ServeHTTP validates authentication and serves a minimal Vast.ai contract.
func (f *fixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.Method+" "+r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer fixture-credential" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if f.mode == "unauthorized" || f.mode == "rate-limit" {
		code := http.StatusUnauthorized
		if f.mode == "rate-limit" {
			code = http.StatusTooManyRequests
		}
		w.WriteHeader(code)
		fmt.Fprint(w, `{"msg":"fixture-secret-response"}`)
		return
	}
	if f.mode == "redirect" {
		w.Header().Set("Location", "/credential-target")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if f.mode == "malformed" {
		fmt.Fprint(w, `{`)
		return
	}
	if f.mode == "rejected" {
		fmt.Fprint(w, `{"success":false,"msg":"fixture-secret-response"}`)
		return
	}
	if r.Method == http.MethodPost && r.URL.Path == "/api/v0/bundles/" {
		var query map[string]any
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil || query["limit"] != float64(1) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if f.mode == "empty-offers" {
			fmt.Fprint(w, `{"offers":[]}`)
			return
		}
		fmt.Fprint(w, `{"offers":[{"id":17,"gpu_name":"RTX 4090","num_gpus":1,"gpu_ram":24576,"dph_total":0.25,"reliability":0.99}]}`)
		return
	}
	if r.Method == http.MethodPut && r.URL.Path == "/api/v0/asks/17/" {
		var body map[string]any
		expectedImage := f.image
		if expectedImage == "" {
			expectedImage = "fixture/image:1"
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body["image"] != expectedImage || body["disk"] != float64(20) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if f.mode == "missing-id" {
			fmt.Fprint(w, `{"success":true}`)
			return
		}
		f.exists = true
		f.onstart, _ = body["onstart"].(string)
		f.label, _ = body["label"].(string)
		f.state, _ = body["target_state"].(string)
		fmt.Fprint(w, `{"success":true,"new_contract":91,"instance_api_key":"fixture-secret-response"}`)
		return
	}
	if r.URL.Path == "/api/v0/instances/91/" {
		if !f.exists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"success":false}`)
			return
		}
		switch r.Method {
		case http.MethodGet:
			if f.mode == "null-instance" {
				f.exists = false
				fmt.Fprint(w, `{"instances":null}`)
				return
			}
			if f.mode == "missing-instance" {
				fmt.Fprint(w, `{}`)
				return
			}
			if f.mode == "rejected-instance" {
				fmt.Fprint(w, `{"success":false,"instances":null}`)
				return
			}
			if f.mode == "mismatched-instance" {
				fmt.Fprint(w, `{"instances":{"id":92}}`)
				return
			}
			if f.mode == "read-failure" {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			runtype := f.runtype
			if runtype == "" {
				runtype = "ssh"
			}
			if f.mode == "unknown-runtype" {
				runtype = "args"
			} else if f.mode == "missing-runtype" {
				runtype = ""
			}
			image := f.image
			if image == "" {
				image = "fixture/image:1"
			}
			instance := map[string]any{
				"id": 91, "label": f.label, "intended_status": f.state, "actual_status": f.state,
				"image_uuid": image, "disk_space": 20, "image_runtype": runtype,
				"ssh_host": "203.0.113.1", "ssh_port": 22001, "public_ipaddr": "203.0.113.2",
				"ports": map[string]any{"7860/tcp": []any{map[string]any{"HostPort": "27860"}}},
			}
			if f.onstart != "" {
				instance["onstart"] = f.onstart
			}
			switch f.mode {
			case "missing-image":
				delete(instance, "image_uuid")
			case "empty-image":
				instance["image_uuid"] = " "
			case "null-image":
				instance["image_uuid"] = nil
			case "missing-disk":
				delete(instance, "disk_space")
			case "zero-disk":
				instance["disk_space"] = 0
			case "negative-disk":
				instance["disk_space"] = -1
			case "null-disk":
				instance["disk_space"] = nil
			}
			if err := json.NewEncoder(w).Encode(map[string]any{"instances": instance}); err != nil {
				panic(err)
			}
		case http.MethodPut:
			var body struct {
				// Label is the requested replacement label.
				Label string `json:"label"`
				// State is the requested power state.
				State string `json:"state"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			f.label, f.state = body.Label, body.State
			fmt.Fprint(w, `{"success":true}`)
		case http.MethodDelete:
			f.exists = false
			fmt.Fprint(w, `{"success":true}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// harness owns a private Terraform directory and controlled subprocess environment.
type harness struct {
	// directory holds the configuration, state, and local development override.
	directory string
	// terraform is the resolved CLI executable.
	terraform string
	// environment contains only runtime essentials and fixture credentials.
	environment []string
}

// newHarness resolves pinned tools and installs a task-local provider override.
func newHarness(t *testing.T, endpoint string) *harness {
	t.Helper()
	p, err := runfiles.Rlocation(*providerPath)
	if err != nil {
		t.Fatalf("resolve provider: %v", err)
	}
	tf, err := runfiles.Rlocation(*terraformPath)
	if err != nil {
		t.Fatalf("resolve Terraform: %v", err)
	}
	dir := t.TempDir()
	config := filepath.Join(dir, "terraform.rc")
	if err := os.WriteFile(config, []byte(fmt.Sprintf(`provider_installation {
  dev_overrides { "registry.terraform.io/alwaldend/vastai" = %q }
  direct {}
}
`, filepath.Dir(p))), 0o600); err != nil {
		t.Fatal(err)
	}
	h := &harness{directory: dir, terraform: tf, environment: []string{
		"PATH=" + os.Getenv("PATH"), "HOME=" + dir, "TMPDIR=.",
		"TF_CLI_CONFIG_FILE=" + config, "TF_IN_AUTOMATION=1", "CHECKPOINT_DISABLE=1",
		"VAST_API_KEY=fixture-credential",
	}}
	runfilesEnv, err := runfiles.Env()
	if err != nil {
		t.Fatalf("resolve subprocess runfiles: %v", err)
	}
	h.environment = append(h.environment, runfilesEnv...)
	h.write(t, endpoint, "initial", "running", 20)
	return h
}

// write replaces only the fixture's Terraform configuration.
func (h *harness) write(t *testing.T, endpoint, label, state string, disk int) {
	t.Helper()
	text := fmt.Sprintf(`terraform {
  required_providers {
    vastai = { source = "alwaldend/vastai" }
  }
}
provider "vastai" { endpoint = %q }
data "vastai_offers" "gpu" { query = jsonencode({limit=1}) }
resource "vastai_instance" "gpu" {
  offer_id = 17
  image = "fixture/image:1"
  disk_gb = %d
  label = %q
  state = %q
}
output "offer" { value = data.vastai_offers.gpu.offers[0].id }
output "connection" { value = vastai_instance.gpu.ssh_host }
`, endpoint, disk, label, state)
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
}

// run executes the real Terraform CLI and checks its expected outcome.
func (h *harness) run(t *testing.T, success bool, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, h.terraform, args...)
	cmd.Dir, cmd.Env = h.directory, h.environment
	output, err := cmd.CombinedOutput()
	if (err == nil) != success {
		t.Fatalf("terraform %v: %v\n%s", args, err, output)
	}
	if strings.Contains(string(output), "fixture-secret-response") {
		t.Fatal("API response leaked into diagnostics")
	}
	return string(output)
}

// TestLifecycle verifies Terraform plans, persisted state, updates, imports, and drift.
func TestLifecycle(t *testing.T) {
	f := &fixture{}
	srv := httptest.NewServer(f)
	defer srv.Close()
	h := newHarness(t, srv.URL+"/api/v0")
	h.run(t, true, "validate", "-no-color")
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	stateBytes, err := os.ReadFile(filepath.Join(h.directory, "terraform.tfstate"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(stateBytes), "fixture-secret-response") || strings.Contains(string(stateBytes), "fixture-credential") {
		t.Fatal("API credentials were persisted in Terraform state")
	}
	h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
	if output := h.run(t, true, "output", "-raw", "connection"); output != "203.0.113.1" {
		t.Fatalf("connection: %q", output)
	}
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	h.write(t, srv.URL+"/api/v0", "updated", "stopped", 20)
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	f.mu.Lock()
	if f.label != "updated" || f.state != "stopped" {
		t.Error("update not applied")
	}
	f.mu.Unlock()
	h.write(t, srv.URL+"/api/v0", "updated", "stopped", 25)
	if output := h.run(t, false, "plan", "-detailed-exitcode", "-no-color"); !strings.Contains(output, "must be replaced") {
		t.Fatal("disk change did not require replacement")
	}
	h.write(t, srv.URL+"/api/v0", "updated", "stopped", 20)
	h.run(t, true, "state", "rm", "vastai_instance.gpu")
	h.run(t, false, "import", "-no-color", "vastai_instance.gpu", "../91")
	h.run(t, true, "import", "-no-color", "vastai_instance.gpu", "91")
	configuration, err := os.ReadFile(filepath.Join(h.directory, "main.tf"))
	if err != nil {
		t.Fatal(err)
	}
	configuration = []byte(strings.Replace(string(configuration), "offer_id = 17", "offer_id = 17\n  lifecycle { ignore_changes = [offer_id] }", 1))
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), configuration, 0o600); err != nil {
		t.Fatal(err)
	}
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	// Repeated destroy handles remote disappearance without a failed API operation.
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	f.mu.Lock()
	f.exists = false
	f.mu.Unlock()
	h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
	if state := h.run(t, true, "state", "list"); strings.Contains(state, "vastai_instance.gpu") {
		t.Fatal("missing remote instance retained in state")
	}
	f.mu.Lock()
	calls := append([]string(nil), f.calls...)
	f.mu.Unlock()
	encoded, err := json.MarshalIndent(map[string]any{"result": "passed", "calls": calls}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "lifecycle.json"), encoded, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestFailures exercises API rejection and uncertain or malformed responses through Terraform.
func TestFailures(t *testing.T) {
	for _, mode := range []string{"unauthorized", "rate-limit", "redirect", "malformed", "rejected", "missing-id"} {
		t.Run(mode, func(t *testing.T) {
			f := &fixture{mode: mode}
			srv := httptest.NewServer(f)
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			h.run(t, false, "apply", "-auto-approve", "-no-color")
			f.mu.Lock()
			defer f.mu.Unlock()
			if len(f.calls) == 0 {
				t.Fatal("failure did not reach the API fixture")
			}
			creates := 0
			for _, call := range f.calls {
				if strings.Contains(call, "asks/") {
					creates++
				}
				if strings.Contains(call, "credential-target") {
					t.Error("followed redirect")
				}
			}
			if creates > 1 {
				t.Error("retried non-idempotent creation")
			}
		})
	}
}

// TestReadFailureRetainsState ensures a failed refresh cannot silently orphan a rental.
func TestReadFailureRetainsState(t *testing.T) {
	for _, mode := range []string{"read-failure", "missing-instance", "mismatched-instance", "rejected-instance", "unknown-runtype", "missing-runtype", "missing-image", "empty-image", "null-image", "missing-disk", "zero-disk", "negative-disk", "null-disk"} {
		t.Run(mode, func(t *testing.T) {
			f := &fixture{}
			srv := httptest.NewServer(f)
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			h.run(t, true, "apply", "-auto-approve", "-no-color")
			f.mu.Lock()
			f.mode = mode
			f.mu.Unlock()
			h.run(t, false, "apply", "-refresh-only", "-auto-approve", "-no-color")
			if state := h.run(t, true, "state", "list"); !strings.Contains(state, "vastai_instance.gpu") {
				t.Fatal("failed refresh lost ownership")
			}
			state := strings.Join(strings.Fields(h.run(t, true, "state", "show", "-no-color", "vastai_instance.gpu")), " ")
			if !strings.Contains(state, `image = "fixture/image:1"`) || !strings.Contains(state, "disk_gb = 20") {
				t.Fatal("failed refresh changed the retained image or disk allocation")
			}
			report, err := json.Marshal(map[string]any{"result": "passed", "failure_mode": mode, "retained_image": "fixture/image:1", "retained_disk_gb": 20})
			if err != nil {
				t.Fatal(err)
			}
			if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
				if err := os.WriteFile(filepath.Join(dir, "instance-read-failure-"+mode+".json"), report, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestRemoteDeletion verifies both API absence signals permit refresh and recreation.
func TestRemoteDeletion(t *testing.T) {
	for _, mode := range []string{"http-not-found", "null-instance"} {
		t.Run(mode, func(t *testing.T) {
			f := &fixture{}
			srv := httptest.NewServer(f)
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			h.run(t, true, "apply", "-auto-approve", "-no-color")
			f.mu.Lock()
			if mode == "http-not-found" {
				f.exists = false
			} else {
				f.mode = mode
			}
			f.mu.Unlock()
			h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
			if state := h.run(t, true, "state", "list"); strings.Contains(state, "vastai_instance.gpu") {
				t.Fatal("externally deleted instance retained in state")
			}
			plan := h.run(t, false, "plan", "-detailed-exitcode", "-no-color")
			if !strings.Contains(strings.Join(strings.Fields(plan), " "), "1 to add, 0 to change, 0 to destroy") {
				t.Fatal("missing instance did not produce a recreation plan")
			}
			h.run(t, true, "apply", "-auto-approve", "-no-color")
			f.mu.Lock()
			f.mode = ""
			f.mu.Unlock()
			h.run(t, true, "destroy", "-auto-approve", "-no-color")
			f.mu.Lock()
			calls := append([]string(nil), f.calls...)
			f.mu.Unlock()
			encoded, err := json.MarshalIndent(map[string]any{"result": "passed", "absence": mode, "calls": calls}, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
				if err := os.WriteFile(filepath.Join(output, "remote-deletion-"+mode+".json"), encoded, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestValidation checks known invalid configuration without issuing rental requests.
func TestValidation(t *testing.T) {
	for _, testCase := range []struct {
		// name identifies the invalid configuration.
		name string
		// old is the valid HCL fragment to replace.
		old string
		// replacement is the invalid HCL fragment.
		replacement string
		// expected is the required provider diagnostic.
		expected string
	}{
		{"offer", "offer_id = 17", "offer_id = 0", "Invalid offer ID"},
		{"disk", "disk_gb = 20", "disk_gb = -1", "Invalid disk allocation"},
		{"image", `image = "fixture/image:1"`, `image = " "`, "Invalid image"},
		{"state", `state = "running"`, `state = "paused"`, "Invalid power state"},
		{"runtype", `label = "initial"`, `runtype = "invalid"`, "Invalid launch mode"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			f := &fixture{}
			srv := httptest.NewServer(f)
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			file := filepath.Join(h.directory, "main.tf")
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file, []byte(strings.Replace(string(data), testCase.old, testCase.replacement, 1)), 0o600); err != nil {
				t.Fatal(err)
			}
			if output := h.run(t, false, "validate", "-no-color"); !strings.Contains(output, testCase.expected) {
				t.Fatalf("missing expected diagnostic %q", testCase.expected)
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			if len(f.calls) > 0 {
				t.Fatal("validation contacted the API")
			}
		})
	}
}

// TestSearchAndConfiguration verifies empty results and pre-request configuration errors.
func TestSearchAndConfiguration(t *testing.T) {
	for _, mode := range []string{"empty-offers", "invalid-query", "missing-key", "insecure-endpoint"} {
		t.Run(mode, func(t *testing.T) {
			f := &fixture{mode: mode}
			srv := httptest.NewServer(f)
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			file := filepath.Join(h.directory, "main.tf")
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			text := string(data)
			switch mode {
			case "empty-offers":
				text = strings.Replace(text, "data.vastai_offers.gpu.offers[0].id", "length(data.vastai_offers.gpu.offers)", 1)
			case "invalid-query":
				text = strings.Replace(text, "jsonencode({limit=1})", "jsonencode([])", 1)
			case "missing-key":
				h.environment = append(h.environment, "VAST_API_KEY=")
			case "insecure-endpoint":
				text = strings.Replace(text, srv.URL+"/api/v0", "http://example.invalid/api/v0", 1)
			}
			if err := os.WriteFile(file, []byte(text), 0o600); err != nil {
				t.Fatal(err)
			}
			output := h.run(t, mode == "empty-offers", "plan", "-no-color")
			expected := map[string]string{"empty-offers": "offer = 0", "invalid-query": "Invalid offer query", "missing-key": "Missing API key", "insecure-endpoint": "Invalid API configuration"}[mode]
			if !strings.Contains(strings.Join(strings.Fields(output), " "), expected) {
				t.Fatalf("missing diagnostic or empty result %q", expected)
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			for _, call := range f.calls {
				if strings.Contains(call, "asks/") {
					t.Fatal("search or configuration failure rented an instance")
				}
			}
		})
	}
}

// TestLaunchModes verifies equivalent API modes, canonical imports, and real drift.
func TestLaunchModes(t *testing.T) {
	for _, testCase := range []struct {
		// name identifies the API representation being exercised.
		name string
		// configured is the operator's supported launch-mode setting.
		configured string
		// remote is the launch-mode string returned by the API.
		remote string
		// imported is the deterministic supported mode selected without prior state.
		imported string
	}{
		{"ssh_default", "ssh", "ssh ssh_direc ssh_proxy", "ssh"},
		{"ssh_direct", "ssh_direct", "ssh_direc ssh_proxy", "ssh_direct"},
		{"ssh_proxy", "ssh_proxy", "ssh_direc ssh_proxy", "ssh_direct"},
		{"jupyter_default", "jupyter", "jupyter_direc ssh_direc ssh_proxy", "jupyter_direct"},
		{"jupyter_direct", "jupyter_direct", "jupyter_direc ssh_direc ssh_proxy", "jupyter_direct"},
		{"jupyter_proxy", "jupyter_proxy", "jupyter_proxy ssh_proxy", "jupyter_proxy"},
		{"ssh_alias", "ssh_direct", "ssh_direc", "ssh_direct"},
		{"jupyter_alias", "jupyter_direct", "jupyter_direc", "jupyter_direct"},
		{"reordered_ssh", "ssh_proxy", " ssh_proxy  ssh_direc\tssh ", "ssh"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			f := &fixture{runtype: testCase.remote}
			srv := httptest.NewServer(f)
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			file := filepath.Join(h.directory, "main.tf")
			original, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			configuration := strings.Replace(string(original), "offer_id = 17", fmt.Sprintf("offer_id = 17\n  runtype = %q\n  lifecycle { ignore_changes = [offer_id] }", testCase.configured), 1)
			if err := os.WriteFile(file, []byte(configuration), 0o600); err != nil {
				t.Fatal(err)
			}
			h.run(t, true, "apply", "-auto-approve", "-no-color")
			h.run(t, true, "apply", "-refresh-only", "-auto-approve", "-no-color")
			h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
			h.run(t, true, "state", "rm", "vastai_instance.gpu")
			h.run(t, true, "import", "-no-color", "vastai_instance.gpu", "91")
			configuration = strings.Replace(configuration, fmt.Sprintf("runtype = %q", testCase.configured), fmt.Sprintf("runtype = %q", testCase.imported), 1)
			if err := os.WriteFile(file, []byte(configuration), 0o600); err != nil {
				t.Fatal(err)
			}
			h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
			f.mu.Lock()
			f.runtype = "jupyter_proxy ssh_proxy"
			if strings.HasPrefix(testCase.imported, "jupyter") {
				f.runtype = "ssh_proxy"
			}
			f.mu.Unlock()
			if output := h.run(t, false, "plan", "-detailed-exitcode", "-no-color"); !strings.Contains(output, "must be replaced") {
				t.Fatal("a real launch-family change did not plan replacement")
			}
			h.run(t, true, "destroy", "-auto-approve", "-no-color")
			report, err := json.Marshal(map[string]any{"result": "passed", "configured": testCase.configured, "api_mode": testCase.remote, "imported": testCase.imported, "equivalent_plan": "unchanged", "real_drift": "replacement"})
			if err != nil {
				t.Fatal(err)
			}
			if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
				if err := os.WriteFile(filepath.Join(dir, "launch-mode-"+testCase.name+".json"), report, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
