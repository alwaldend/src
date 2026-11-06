package acceptance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// existingTeamConfiguration wires the actual module to an explicitly selected issuer.
func existingTeamConfiguration(t *testing.T, endpoint string) string {
	t.Helper()
	source, err := runfiles.Rlocation(*teamModulePath)
	if err != nil {
		t.Fatalf("resolve module: %v", err)
	}
	return fmt.Sprintf(`terraform {
 required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" {
 alias = "issuer"
 endpoint = %q
 api_key = "fixture-credential"
}
module "existing" {
 source = %q
 providers = { vastai = vastai.issuer }
 name = "existing_team"
 api_keys = { rentals = { name = "existing_team_key", permissions = "" } }
}
output "team_id" { value = module.existing.team_id }
`, endpoint, filepath.Dir(source))
}

// TestExistingTeamKey verifies independent binding and cleanup through Terraform.
func TestExistingTeamKey(t *testing.T) {
	for _, mode := range []string{"success", "activation-rejected", "missing-team", "scoped-permissions", "issuer-migration", "issuer-update-denied", "issuer-override-replacement", "module", "module-activation-rejected", "module-invalid-import", "module-scoped-permissions"} {
		t.Run(mode, func(t *testing.T) {
			var mu sync.Mutex
			exists, selected := false, false
			calls := []string{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				calls = append(calls, r.Method+" "+r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				issuer := r.Header.Get("Authorization") == "Bearer fixture-credential"
				generated := r.Header.Get("Authorization") == "Bearer fixture-existing-team-key"
				switch r.Method + " " + r.URL.Path {
				case "GET /api/v0/auth/apikeys/":
					if !issuer {
						w.WriteHeader(403)
						return
					}
					keys := []map[string]any{}
					if mode != "missing-team" {
						keys = append(keys, map[string]any{"id": 99, "name": "owner", "team_id": 42, "team_name": "existing_team"})
					}
					if exists {
						keys = append(keys, map[string]any{"id": 100, "name": "existing_team_key"})
					}
					if err := json.NewEncoder(w).Encode(map[string]any{"apikeys": keys}); err != nil {
						panic(err)
					}
				case "POST /api/v0/auth/apikeys/":
					if !issuer || exists {
						w.WriteHeader(403)
						return
					}
					var body map[string]any
					if json.NewDecoder(r.Body).Decode(&body) != nil || body["permissions"] != nil {
						w.WriteHeader(400)
						return
					}
					exists = true
					fmt.Fprint(w, `{"id":100,"key":"fixture-existing-team-key"}`)
				case "GET /api/v0/users/current/":
					if !generated {
						w.WriteHeader(403)
						return
					}
					if selected {
						fmt.Fprint(w, `{"id":42,"is_team":true}`)
					} else {
						fmt.Fprint(w, `{"id":1,"is_team":false}`)
					}
				case "PUT /api/v0/users/save-context/":
					var body map[string]int64
					if !generated || strings.Contains(mode, "activation-rejected") || json.NewDecoder(r.Body).Decode(&body) != nil || body["context_key_id"] != 99 {
						w.WriteHeader(400)
						return
					}
					selected = true
					fmt.Fprint(w, `{}`)
				case "DELETE /api/v0/auth/apikeys/100/":
					if !issuer {
						w.WriteHeader(403)
						return
					}
					exists = false
					fmt.Fprint(w, `{}`)
				default:
					w.WriteHeader(403)
				}
			}))
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			for i, entry := range h.environment {
				if strings.HasPrefix(entry, "VAST_API_KEY=") {
					h.environment[i] = "VAST_API_KEY=fixture-unprivileged"
				}
			}
			configuration := fmt.Sprintf(`terraform {
 required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" {
 endpoint = %q
 api_key = "fixture-credential"
}
resource "vastai_api_key" "existing" {
 name = "existing_team_key"
 team_id = 42
 permissions = ""
}
`, srv.URL+"/api/v0")
			issuerOverride := strings.HasPrefix(mode, "issuer-")
			migration := issuerOverride && mode != "issuer-override-replacement"
			if issuerOverride {
				configuration = strings.Replace(configuration, `name = "existing_team_key"`, "name = \"existing_team_key\"\n team_api_key = \"fixture-credential\"", 1)
			}
			address := "vastai_api_key.existing"
			moduleMode := strings.HasPrefix(mode, "module")
			if moduleMode {
				configuration = existingTeamConfiguration(t, srv.URL+"/api/v0")
				address = `module.existing.vastai_api_key.project["rentals"]`
			}
			if mode == "scoped-permissions" || mode == "module-scoped-permissions" {
				configuration = strings.Replace(configuration, `permissions = ""`, `permissions = jsonencode({api={user_read={}}})`, 1)
			}
			if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(configuration), 0o600); err != nil {
				t.Fatal(err)
			}
			if moduleMode {
				h.run(t, true, "get", "-no-color")
				if mode == "module-invalid-import" {
					h.run(t, false, "import", "-no-color", "module.existing.vastai_team.team", "0")
					mu.Lock()
					defer mu.Unlock()
					if len(calls) != 0 || exists {
						t.Fatal("invalid import reached the API")
					}
					if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
						if err := os.WriteFile(filepath.Join(dir, "existing-team-key-"+mode+".json"), []byte(`{"result":"passed","import_rejected":true,"requests":0}`), 0o600); err != nil {
							t.Fatal(err)
						}
					}
					return
				}
				h.run(t, true, "import", "-no-color", "module.existing.vastai_team.team", "42")
				mu.Lock()
				calls = nil // Import is read-only; measure the following key apply independently.
				mu.Unlock()
			}
			output := h.run(t, mode == "success" || issuerOverride || mode == "module", "apply", "-auto-approve", "-no-color")
			if strings.Contains(output, "fixture-existing-team-key") {
				t.Fatal("credential leaked")
			}
			mu.Lock()
			created := exists
			requestCount := len(calls)
			mu.Unlock()
			if mode == "scoped-permissions" || mode == "module-scoped-permissions" {
				mu.Lock()
				for _, call := range calls {
					if mode == "scoped-permissions" || !strings.HasPrefix(call, "GET ") {
						t.Error("invalid key configuration caused an unexpected API request")
					}
				}
				mu.Unlock()
				if created {
					t.Fatal("scoped owner binding created a key")
				}
				if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
					data, err := json.Marshal(map[string]any{"result": "passed", "mode": mode, "configuration_rejected": true, "requests": requestCount, "mutations": 0})
					if err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(dir, "existing-team-key-"+mode+".json"), data, 0o600); err != nil {
						t.Fatal(err)
					}
				}
				return
			}
			if created != (mode != "missing-team") {
				t.Fatalf("unexpected partial-key ownership: %v\n%s", created, output)
			}
			if mode == "success" || issuerOverride || mode == "module" {
				h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
			}
			if mode == "issuer-override-replacement" {
				replacement := strings.Replace(configuration, `team_api_key = "fixture-credential"`, `team_api_key = "fixture-replacement"`, 1)
				if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(replacement), 0o600); err != nil {
					t.Fatal(err)
				}
				if plan := h.run(t, false, "plan", "-detailed-exitcode", "-no-color"); !strings.Contains(plan, "must be replaced") {
					t.Fatal("changing the bootstrap issuer no longer replaces its project key")
				}
				if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(configuration), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if strings.Contains(mode, "activation-rejected") && !strings.Contains(h.run(t, true, "state", "list"), address) {
				t.Fatal("failed activation lost generated key")
			}
			if migration {
				configuration = strings.Replace(configuration, " team_api_key = \"fixture-credential\"\n", "", 1)
				if mode == "issuer-update-denied" {
					configuration = strings.Replace(configuration, `api_key = "fixture-credential"`, `api_key = "fixture-unprivileged"`, 1)
				}
				if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(configuration), 0o600); err != nil {
					t.Fatal(err)
				}
				plan := h.run(t, false, "plan", "-detailed-exitcode", "-no-color")
				if strings.Contains(plan, "must be replaced") {
					t.Fatal("removing issuer metadata replaced the key")
				}
				h.run(t, mode != "issuer-update-denied", "apply", "-auto-approve", "-no-color")
				if mode == "issuer-update-denied" {
					if state := h.run(t, true, "state", "show", address); !strings.Contains(state, "team_api_key") {
						t.Fatal("failed issuer change abandoned the previous credential metadata")
					}
					configuration = strings.Replace(configuration, `api_key = "fixture-unprivileged"`, `api_key = "fixture-credential"`, 1)
					if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(configuration), 0o600); err != nil {
						t.Fatal(err)
					}
					h.run(t, true, "apply", "-auto-approve", "-no-color")
				}
				if state := h.run(t, true, "state", "show", address); strings.Contains(state, "team_api_key") {
					t.Fatal("redundant resource credential remains in state")
				}
				h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
				if mode == "issuer-migration" {
					configuration = existingTeamConfiguration(t, srv.URL+"/api/v0") + `
moved {
 from = vastai_api_key.existing
 to = module.existing.vastai_api_key.project["rentals"]
}
`
					if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(configuration), 0o600); err != nil {
						t.Fatal(err)
					}
					h.run(t, true, "get", "-no-color")
					h.run(t, true, "import", "-no-color", "module.existing.vastai_team.team", "42")
					plan := h.run(t, false, "plan", "-detailed-exitcode", "-no-color")
					if !strings.Contains(plan, "0 to add, 0 to change, 0 to destroy") {
						t.Fatalf("moving key into module would mutate it: %s", plan)
					}
					h.run(t, true, "apply", "-auto-approve", "-no-color")
					address = `module.existing.vastai_api_key.project["rentals"]`
					if !strings.Contains(h.run(t, true, "state", "list"), address) {
						t.Fatal("module move lost generated key")
					}
					// The temporary migration is complete; final configuration must stand alone.
					configuration = existingTeamConfiguration(t, srv.URL+"/api/v0")
					if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(configuration), 0o600); err != nil {
						t.Fatal(err)
					}
					h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
				}
				mu.Lock()
				mutations := 0
				for _, call := range calls {
					if strings.HasPrefix(call, "POST ") || strings.HasPrefix(call, "DELETE ") {
						mutations++
					}
				}
				mu.Unlock()
				if mutations != 1 {
					t.Fatal("issuer migration rotated the key")
				}
			}
			if moduleMode || mode == "issuer-migration" {
				h.run(t, true, "destroy", "-target="+address, "-auto-approve", "-no-color")
				if output := h.run(t, false, "destroy", "-auto-approve", "-no-color"); !strings.Contains(output, "Imported teams cannot be destroyed") {
					t.Fatal("imported team destruction was not refused")
				}
				// Discard only fixture ownership after proving the remote team is protected.
				h.run(t, true, "state", "rm", "module.existing.vastai_team.team")
			}
			h.run(t, true, "destroy", "-auto-approve", "-no-color")
			mu.Lock()
			defer mu.Unlock()
			if exists {
				t.Fatal("generated key survived destruction")
			}
			for _, call := range calls {
				if call == "POST /api/v0/team/" || call == "DELETE /api/v0/team/" {
					t.Fatal("existing team was mutated")
				}
			}
			if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
				data, err := json.Marshal(map[string]any{"result": "passed", "mode": mode, "key_revoked": true, "calls": calls})
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "existing-team-key-"+mode+".json"), data, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
