package main_test

import (
	"context"
	"encoding/json"
	"encoding/pem"
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

// fixture serves isolated Forgejo and Vault contracts to the actual executable.
type fixture struct {
	// mu serializes concurrent fake service requests.
	mu sync.Mutex
	// server is a private TLS service with no external dependencies.
	server *httptest.Server
	// environment supplies only public fixture credentials to subprocesses.
	environment []string
	// secret is the latest KV v2 record, if present.
	secret map[string]any
	// version enforces create-only Vault writes.
	version int
	// creates counts registration POST requests.
	creates int
	// writes counts accepted Vault writes.
	writes int
	// present selects whether the remote runner exists.
	present bool
	// failure selects a deliberately failing service operation.
	failure string
	// collision places a same-name runner on the second list page.
	collision bool
	// forwarded counts requests to a redirect destination.
	forwarded int
}

// newFixture prepares a trusted loopback HTTPS service and hermetic environment.
func newFixture(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{}
	f.server = httptest.NewTLSServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.server.Close)
	certificate := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(certificate, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: f.server.Certificate().Raw}), 0o600); err != nil {
		t.Fatal(err)
	}
	f.environment = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + t.TempDir(),
		"SSL_CERT_FILE=" + certificate,
		"VAULT_CACERT=" + certificate,
		"VAULT_ADDR=" + f.server.URL,
		"VAULT_TOKEN=fixture-vault-token",
		"FORGEJO_HOST=" + f.server.URL,
		"FORGEJO_API_TOKEN=fixture-api-token",
		"FORGEJO_RUNNER_REPOSITORY=owner/repository",
		"FORGEJO_RUNNER_NAME=fixture-worker",
		"FORGEJO_RUNNER_VAULT_MOUNT=secrets",
		"FORGEJO_RUNNER_VAULT_PATH=fixture/runners/worker",
	}
	for _, variable := range os.Environ() {
		if strings.HasPrefix(variable, "RUNFILES_") || strings.HasPrefix(variable, "TEST_") {
			f.environment = append(f.environment, variable)
		}
	}
	return f
}

// TestTerraformCreationOnly verifies the real provisioner declaration against
// built-in Terraform resources and private dummy services, with no live backend.
func TestTerraformCreationOnly(t *testing.T) {
	f := newFixture(t)
	workspace := t.TempDir()
	tool := filepath.Join(workspace, "runner_registration")
	content, err := os.ReadFile(executable(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tool, content, 0o700); err != nil {
		t.Fatal(err)
	}
	source, err := runfiles.Rlocation(os.Getenv("BOOTSTRAP_SOURCE"))
	if err != nil {
		t.Fatal(err)
	}
	declaration, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "bootstrap.tf"), declaration, 0o600); err != nil {
		t.Fatal(err)
	}
	configuration := map[string]any{"locals": map[string]any{
		"runner_registration_tool": tool,
		"runner_registrations": map[string]any{"fixture": map[string]any{
			"repository": "owner/repository", "name": "fixture-worker",
			"mount": "secrets", "path": "fixture/runners/worker",
		}},
	}}
	raw, err := json.Marshal(configuration)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "fixture.tf.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	terraform, err := runfiles.Rlocation(os.Getenv("TERRAFORM"))
	if err != nil {
		t.Fatal(err)
	}
	run := func(arguments ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, terraform, append([]string{"--chdir", workspace}, arguments...)...)
		command.Env = f.environment
		command.Stdin = strings.NewReader("yes\n")
		output, err := command.CombinedOutput()
		for _, secret := range []string{"fixture-api-token", "fixture-vault-token", "fixture-runner-token"} {
			if strings.Contains(string(output), secret) {
				t.Fatal("credential disclosed by Terraform")
			}
		}
		if err != nil {
			t.Fatalf("isolated Terraform operation failed: %v; %s", err, output)
		}
	}
	run("init", "-backend=false", "-input=false", "-no-color")
	run("plan", "-input=false", "-no-color")
	if f.creates != 0 || f.writes != 0 {
		t.Fatal("plan invoked registration")
	}
	run("apply", "-no-color")
	if f.creates != 1 || f.writes != 1 {
		t.Fatal("initial apply did not create one registration")
	}
	// Changed invalid executable content proves ordinary applies do not rerun it.
	if err := os.WriteFile(tool, []byte("changed tool content\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("apply", "-no-color")
	if f.creates != 1 || f.writes != 1 {
		t.Fatal("ordinary apply invalidated the cached registration")
	}
	if err := os.WriteFile(tool, content, 0o700); err != nil {
		t.Fatal(err)
	}
	run("apply", "-no-color", "-replace=terraform_data.runners[\"fixture\"]")
	state, err := os.ReadFile(filepath.Join(workspace, "terraform.tfstate"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(state), "fixture-runner-token") || strings.Contains(string(state), "fixture-api-token") || strings.Contains(string(state), "fixture-vault-token") || f.creates != 1 || f.writes != 1 {
		t.Fatal("explicit replacement changed credentials or leaked them to state")
	}
	run("destroy", "-no-color")
	if !f.present || f.secret == nil {
		t.Fatal("Terraform removal deleted retained credentials")
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "terraform-bootstrap.json"), []byte("{\"plan_has_no_mutation\":true,\"creation_once\":true,\"tool_changes_do_not_invalidate\":true,\"replacement_reuses_vault\":true,\"no_credentials_in_state\":true,\"removal_retains_registration\":true,\"live_operations\":0}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// serve emulates only the service operations needed by the bootstrap contract.
func (f *fixture) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/unexpected" {
		f.forwarded++
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/v1/") {
		if f.failure == "vault-redirect" {
			w.Header().Set("Location", f.server.URL+"/unexpected")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		if r.Header.Get("X-Vault-Token") != "fixture-vault-token" || f.failure == "vault-denied" {
			http.Error(w, "fixture-vault-token", http.StatusForbidden)
			return
		}
		if r.URL.Path != "/v1/secrets/data/fixture/runners/worker" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			if f.secret == nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"errors":[]}`)
				return
			}
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"data": f.secret, "metadata": map[string]any{"version": f.version}}})
			return
		}
		var request struct {
			// Data carries the proposed secret record.
			Data map[string]any `json:"data"`
			// Options selects the required current record version.
			Options struct {
				// CAS is Vault's atomic version precondition.
				CAS int `json:"cas"`
			} `json:"options"`
		}
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.Options.CAS != f.version {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if f.failure == "vault-save" {
			http.Error(w, "fixture-runner-token", http.StatusInternalServerError)
			return
		}
		f.secret, f.version, f.writes = request.Data, f.version+1, f.writes+1
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"version": f.version}})
		return
	}
	if r.Header.Get("Authorization") != "token fixture-api-token" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	const collection = "/api/v1/repos/owner/repository/actions/runners"
	if r.URL.Path == collection && r.Method == http.MethodGet {
		if f.failure == "redirect" {
			w.Header().Set("Location", f.server.URL+"/unexpected")
			w.WriteHeader(http.StatusFound)
			return
		}
		if f.collision && r.URL.Query().Get("page") == "1" {
			w.Header().Set("Link", "<"+f.server.URL+collection+"?page=2>; rel=\"next\"")
			fmt.Fprint(w, `[{"id":9,"name":"other"}]`)
			return
		}
		if f.present || f.collision {
			fmt.Fprint(w, `[{"id":7,"name":"fixture-worker","uuid":"11111111-2222-4333-8444-555555555555"}]`)
		} else {
			fmt.Fprint(w, `[]`)
		}
		return
	}
	if r.URL.Path == collection && r.Method == http.MethodPost {
		f.creates++
		if f.failure == "forgejo-create" {
			f.present = true
			http.Error(w, "fixture-api-token fixture-runner-token", http.StatusInternalServerError)
			return
		}
		f.present = true
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":7,"uuid":"11111111-2222-4333-8444-555555555555","token":"fixture-runner-token"}`)
		return
	}
	if r.URL.Path == collection+"/7" && r.Method == http.MethodGet && f.present {
		fmt.Fprint(w, `{"id":7,"name":"fixture-worker","uuid":"11111111-2222-4333-8444-555555555555"}`)
		return
	}
	http.NotFound(w, r)
}

// executable resolves the packaged standalone command from Bazel runfiles.
func executable(t *testing.T) string {
	t.Helper()
	path, err := runfiles.Rlocation(os.Getenv("RUNNER_TOOL"))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// invoke executes the actual CLI and checks that no public dummy secret leaks.
func (f *fixture) invoke(t *testing.T, success bool) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, executable(t), "ensure")
	command.Env = f.environment
	output, err := command.CombinedOutput()
	text := string(output)
	for _, secret := range []string{"fixture-api-token", "fixture-vault-token", "fixture-runner-token"} {
		if strings.Contains(text, secret) {
			t.Fatal("credential disclosed by command")
		}
	}
	if (err == nil) != success {
		t.Fatalf("unexpected command status: %v; %s", err, text)
	}
	return text
}

// TestRegistrationAcceptance verifies creation, reuse, refusal and recovery boundaries.
func TestRegistrationAcceptance(t *testing.T) {
	t.Run("create-and-reuse", func(t *testing.T) {
		f := newFixture(t)
		f.invoke(t, true)
		if f.creates != 1 || f.writes != 1 || f.secret["token"] != "fixture-runner-token" {
			t.Fatal("registration was not persisted")
		}
		f.invoke(t, true)
		if f.creates != 1 || f.writes != 1 {
			t.Fatal("repeat changed registration")
		}
		f.present = false
		f.invoke(t, true)
		if f.creates != 1 {
			t.Fatal("cached credentials were invalidated")
		}
	})
	for _, failure := range []string{"vault-denied", "forgejo-create", "vault-save", "redirect", "vault-redirect"} {
		t.Run(failure, func(t *testing.T) {
			f := newFixture(t)
			f.failure = failure
			f.invoke(t, false)
			if f.forwarded != 0 {
				t.Fatal("credentials forwarded through a redirect")
			}
			creates, writes := f.creates, f.writes
			f.failure = ""
			if failure == "forgejo-create" || failure == "vault-save" {
				f.invoke(t, false)
				if f.creates != creates || f.writes != writes {
					t.Fatal("retry changed incomplete registration")
				}
			} else if f.creates != 0 || f.writes != 0 {
				t.Fatal("preflight failure mutated services")
			}
		})
	}
	t.Run("later-page-collision", func(t *testing.T) {
		f := newFixture(t)
		f.collision = true
		f.invoke(t, false)
		if f.creates != 0 || f.writes != 0 {
			t.Fatal("existing runner was modified")
		}
	})
	for _, invalid := range []string{
		"FORGEJO_API_TOKEN=",
		"VAULT_TOKEN=",
		"FORGEJO_HOST=http://localhost",
		"FORGEJO_RUNNER_VAULT_PATH=../another/path",
		"VAULT_SKIP_VERIFY=true",
		"SSL_CERT_FILE=/nonexistent-fixture-certificate",
	} {
		t.Run(strings.SplitN(invalid, "=", 2)[0], func(t *testing.T) {
			f := newFixture(t)
			f.environment = append(f.environment, invalid)
			f.invoke(t, false)
			if f.creates != 0 || f.writes != 0 {
				t.Fatal("invalid configuration or untrusted TLS mutated services")
			}
		})
	}
	t.Run("conflicting-record", func(t *testing.T) {
		f := newFixture(t)
		f.invoke(t, true)
		f.secret["repository"] = "another/repository"
		f.invoke(t, false)
		if f.creates != 1 || f.writes != 1 {
			t.Fatal("conflicting record was overwritten")
		}
	})
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "runner-registration.json"), []byte("{\"creation_and_reuse\":true,\"partial_failure_stops\":true,\"pagination_collision\":true,\"redacted_output\":true,\"live_operations\":0}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
