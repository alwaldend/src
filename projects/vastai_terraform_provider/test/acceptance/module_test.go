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
	"strconv"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// teamModulePath locates the reusable module's declared source input.
var teamModulePath = flag.String("team-module", "", "Team module source runfile")

// moduleFixture exposes several teams through one personal-account API endpoint.
type moduleFixture struct {
	// teams holds independent account and key state behind the shared endpoint.
	teams []*ownershipFixture
}

// ServeHTTP routes generated credentials and aggregates personal owner contexts.
func (f *moduleFixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	personal := r.Header.Get("Authorization") == "Bearer fixture-credential"
	if personal && r.Method == http.MethodGet && r.URL.Path == "/api/v0/auth/apikeys/" {
		keys := []map[string]any{}
		for _, team := range f.teams {
			team.mu.Lock()
			if team.teamExists {
				// Owner-context metadata must not collide with generated fixture keys.
				keys = append(keys, map[string]any{"id": team.accountID + 1_000_000, "name": "owner", "team_id": team.accountID, "team_name": team.teamName})
			}
			for id := range team.keys {
				keys = append(keys, map[string]any{"id": id, "name": "test_generated"})
			}
			team.mu.Unlock()
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"apikeys": keys}); err != nil {
			panic(err)
		}
		return
	}
	if personal && r.Method == http.MethodPost && r.URL.Path == "/api/v0/auth/apikeys/" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		var input struct {
			// Name associates each generated management key with its requested team.
			Name string `json:"name"`
		}
		if err := json.Unmarshal(body, &input); err != nil {
			panic(err)
		}
		for _, team := range f.teams {
			if input.Name == team.teamName+"_management" {
				team.ServeHTTP(w, r)
				return
			}
		}
	} else {
		keyID, _ := strconv.ParseInt(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer fixture-management-"), 10, 64)
		if personal {
			keyID, _ = strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v0/auth/apikeys/"), "/"), 10, 64)
		}
		for _, team := range f.teams {
			team.mu.Lock()
			ownsKey := team.keys[keyID]
			team.mu.Unlock()
			if ownsKey {
				team.ServeHTTP(w, r)
				return
			}
		}
	}
	w.WriteHeader(http.StatusForbidden)
}

// TestTeamModule exercises independent module instances and multiple project keys.
func TestTeamModule(t *testing.T) {
	first := &ownershipFixture{teamName: "test_first", accountID: 42, nextKey: 100, keys: map[int64]bool{}, contexts: map[int64]int64{}}
	second := &ownershipFixture{teamName: "test_second", accountID: 43, nextKey: 200, keys: map[int64]bool{}, contexts: map[int64]int64{}}
	server := httptest.NewServer(&moduleFixture{teams: []*ownershipFixture{first, second}})
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	source, err := runfiles.Rlocation(*teamModulePath)
	if err != nil {
		t.Fatalf("resolve team module: %v", err)
	}
	config := fmt.Sprintf(`terraform {
  required_providers { vastai = { source = "alwaldend/vastai" } }
}


provider "vastai" { endpoint = %q }
locals {
  teams = {
  test_first = {
    first = { name = "test_first_project", permissions = jsonencode({api={user_read={}}}) }
    second = { name = "test_second_project", permissions = jsonencode({api={team_read={}}}) }
  }
  test_second = {
    third = { name = "test_third_project", permissions = jsonencode({api={user_read={}}}) }
  }
  }
}
module "teams" {
  for_each = local.teams
  source = %q
  name = each.key
  api_keys = each.value
}
`, server.URL+"/api/v0", filepath.Dir(source))
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	// get installs only the declared local modules, without provider registry access.
	h.run(t, true, "get", "-no-color")
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	for index, f := range []*ownershipFixture{first, second} {
		f.mu.Lock()
		expected := 3 - index
		if !f.teamExists || len(f.keys) != expected {
			t.Errorf("module %d created %d keys; expected %d", index, len(f.keys), expected)
		}
		f.mu.Unlock()
	}
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	for _, f := range []*ownershipFixture{first, second} {
		f.mu.Lock()
		if f.teamExists || len(f.keys) != 0 {
			t.Error("module cleanup left a team or key")
		}
		f.mu.Unlock()
	}
	report, err := json.Marshal(map[string]any{"result": "passed", "bootstrap_accounts": 1, "teams": 2, "project_keys": 3, "credit_transfer": 0, "cleanup": "complete"})
	if err != nil {
		t.Fatal(err)
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "team-module.json"), report, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestTeamModuleUpgrade preserves the singleton-team state contract from b2696dd.
func TestTeamModuleUpgrade(t *testing.T) {
	f := &ownershipFixture{teamName: "test_first", accountID: 42, nextKey: 100, keys: map[int64]bool{}, contexts: map[int64]int64{}}
	server := httptest.NewServer(&moduleFixture{teams: []*ownershipFixture{f}})
	defer server.Close()
	h := newHarness(t, server.URL+"/api/v0")
	legacy := filepath.Join(h.directory, "legacy")
	if err := os.Mkdir(legacy, 0o700); err != nil {
		t.Fatal(err)
	}
	// This frozen input models the pre-existing singleton address, not the new mode.
	legacySource := `terraform {
 required_providers { vastai = { source = "alwaldend/vastai" } }
}
variable "name" { type = string }
variable "api_keys" { type = map(object({ name = string, permissions = string })) }
resource "vastai_team" "team" { name = var.name }
resource "vastai_api_key" "project" {
 for_each = var.api_keys
 name = each.value.name
 permissions = each.value.permissions
 team_api_key = vastai_team.team.api_key
}
`
	if err := os.WriteFile(filepath.Join(legacy, "main.tf"), []byte(legacySource), 0o600); err != nil {
		t.Fatal(err)
	}
	configuration := `terraform {
 required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" { endpoint = %q }
module "team" {
 source = %q
 name = "test_first"
 api_keys = { reader = { name = "test_first_project", permissions = jsonencode({api={user_read={}}}) } }
}
`
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(fmt.Sprintf(configuration, server.URL+"/api/v0", legacy)), 0o600); err != nil {
		t.Fatal(err)
	}
	h.run(t, true, "get", "-no-color")
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	f.mu.Lock()
	previousCalls := len(f.calls)
	f.mu.Unlock()
	source, err := runfiles.Rlocation(*teamModulePath)
	if err != nil {
		t.Fatalf("resolve module: %v", err)
	}
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(fmt.Sprintf(configuration, server.URL+"/api/v0", filepath.Dir(source))), 0o600); err != nil {
		t.Fatal(err)
	}
	h.run(t, true, "get", "-no-color")
	plan := h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	if !strings.Contains(plan, "No changes.") {
		t.Fatalf("upgrade would change managed resources: %s", plan)
	}
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	state := h.run(t, true, "state", "list")
	if !strings.Contains(state, "module.team.vastai_team.team") || !strings.Contains(state, `module.team.vastai_api_key.project["reader"]`) {
		t.Fatal("upgrade lost managed addresses")
	}
	f.mu.Lock()
	for _, call := range f.calls[previousCalls:] {
		if !strings.HasPrefix(call, "GET ") {
			t.Errorf("upgrade mutated API resources: %s", call)
		}
	}
	if f.nextKey != 102 || len(f.keys) != 2 || !f.teamExists {
		t.Error("upgrade changed team or key identity")
	}
	f.mu.Unlock()
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.teamExists || len(f.keys) != 0 {
		t.Fatal("owned module resources survived cleanup")
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "team-module-upgrade.json"), []byte(`{"result":"passed","upgrade_api_mutations":0,"generated_keys":2,"team_identity_preserved":true,"cleanup":"complete"}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
