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
)

// teamFixture serves only non-compute endpoints and rejects every rental request.
type teamFixture struct {
	// mu protects fixture state.
	mu sync.Mutex
	// key is the registered SSH public key.
	key string
	// variables contains environment fixture values.
	variables map[string]string
	// permissions contains the test role's permissions, or nil if absent.
	permissions any
	// calls records non-sensitive API operations.
	calls []string
	// failReads selects a malformed response to test state retention.
	failReads bool
}

// ServeHTTP implements documented non-compute response shapes.
func (f *teamFixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.Method+" "+r.URL.EscapedPath())
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer fixture-credential" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if f.failReads && r.Method == http.MethodGet {
		if r.URL.Path == "/api/v0/ssh/" {
			fmt.Fprint(w, `[{}]`)
			return
		}
		fmt.Fprint(w, `{}`)
		return
	}
	var body map[string]any
	if r.Method != http.MethodGet {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil && r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	switch r.URL.EscapedPath() {
	case "/api/v0/ssh/":
		if r.Method == http.MethodPost {
			f.key, _ = body["ssh_key"].(string)
			fmt.Fprint(w, `{"success":true,"key":{"id":7}}`)
			return
		}
		if f.key == "" {
			fmt.Fprint(w, `[]`)
			return
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{{"id": 7, "key": f.key, "deleted_at": nil}}); err != nil {
			panic(err)
		}
	case "/api/v0/ssh/7/":
		if r.Method == http.MethodPut {
			f.key, _ = body["ssh_key"].(string)
		} else if r.Method == http.MethodDelete {
			f.key = ""
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprint(w, `{"success":true}`)
	case "/api/v0/secrets/":
		if r.Method == http.MethodGet {
			if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "secrets": f.variables}); err != nil {
				panic(err)
			}
			return
		}
		name, _ := body["key"].(string)
		if r.Method == http.MethodDelete {
			delete(f.variables, name)
		} else {
			f.variables[name], _ = body["value"].(string)
		}
		fmt.Fprint(w, `{"success":true}`)
	case "/api/v0/team/roles/":
		if r.Method != http.MethodPost || body["name"] != "test_engineering/read" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		f.permissions = body["permissions"]
		fmt.Fprint(w, `{"msg":"success"}`)
	case "/api/v0/team/roles/test_engineering%2Fread/":
		if r.Method == http.MethodDelete {
			f.permissions = nil
			fmt.Fprint(w, `{"msg":"success"}`)
			return
		}
		if f.permissions == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"id": 9, "name": "test_engineering/read", "permissions": f.permissions}); err != nil {
			panic(err)
		}
	case "/api/v0/team/roles/9/":
		f.permissions = body["permissions"]
		fmt.Fprint(w, `"Successfully Updated Team Role For test_engineering/read"`)
	default:
		w.WriteHeader(http.StatusForbidden)
	}
}

// writeTeam writes a configuration containing exactly three non-compute resource types.
func (h *harness) writeTeam(t *testing.T, endpoint, marker string) {
	t.Helper()
	text := fmt.Sprintf(`terraform {
  required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" { endpoint = %q }
resource "vastai_ssh_key" "test_key" { public_key = "ssh-ed25519 fixture-%s test_provider" }
resource "vastai_environment_variable" "test_marker" {
  name = "TEST_PROVIDER"
  value = %q
}
resource "vastai_team_role" "test_reader" {
  name = "test_engineering/read"
  permissions = jsonencode({api = {user_read = {marker = %q}}})
}
`, endpoint, marker, marker, marker)
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestTeamResources verifies CRUD, import, drift, and malformed-read behavior without compute.
func TestTeamResources(t *testing.T) {
	f := &teamFixture{variables: map[string]string{}}
	srv := httptest.NewServer(f)
	defer srv.Close()
	h := newHarness(t, srv.URL+"/api/v0")
	h.writeTeam(t, srv.URL+"/api/v0", "initial")
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	h.writeTeam(t, srv.URL+"/api/v0", "updated")
	h.run(t, true, "apply", "-auto-approve", "-no-color")
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	for address, id := range map[string]string{"vastai_ssh_key.test_key": "7", "vastai_environment_variable.test_marker": "TEST_PROVIDER", "vastai_team_role.test_reader": "test_engineering/read"} {
		h.run(t, true, "state", "rm", address)
		h.run(t, true, "import", "-no-color", address, id)
	}
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	f.mu.Lock()
	f.failReads = true
	f.mu.Unlock()
	failedRefresh := h.run(t, false, "apply", "-refresh-only", "-auto-approve", "-no-color")
	if !strings.Contains(failedRefresh, "Cannot read SSH key") {
		t.Fatal("malformed SSH collection was interpreted as remote deletion")
	}
	if state := h.run(t, true, "state", "list"); len(strings.Fields(state)) != 3 {
		t.Fatal("malformed read abandoned managed resources")
	}
	f.mu.Lock()
	f.failReads = false
	f.mu.Unlock()
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	f.mu.Lock()
	if f.key != "" || len(f.variables) != 0 || f.permissions != nil {
		t.Error("destroy left test objects")
	}
	calls := append([]string(nil), f.calls...)
	f.mu.Unlock()
	for _, call := range calls {
		if strings.Contains(call, "asks/") || strings.Contains(call, "instances/") {
			t.Fatal("non-compute test used compute API")
		}
	}
	encoded, err := json.MarshalIndent(map[string]any{"result": "passed", "calls": calls}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "team-resources.json"), encoded, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestTeamBootstrap verifies team creation and key generation without exposing tokens.
func TestTeamBootstrap(t *testing.T) {
	var mu sync.Mutex
	teamExists, keyExists := false, false
	projectID, projectGenerations := int64(0), 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		owner := r.Header.Get("Authorization") == "Bearer fixture-team-management"
		if !owner && r.Header.Get("Authorization") != "Bearer fixture-credential" {
			w.WriteHeader(401)
			return
		}
		projectPath := fmt.Sprintf("/api/v0/auth/apikeys/%d/", projectID)
		switch r.Method + " " + r.URL.Path {
		case "POST /api/v0/team/":
			if !owner {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			var input map[string]any
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input["team_name"] != "test_provider" || input["transfer_credit"] != float64(0) {
				w.WriteHeader(400)
				return
			}
			teamExists = true
			fmt.Fprint(w, `{"success":true}`)
		case "GET /api/v0/auth/apikeys/":
			if !teamExists {
				fmt.Fprint(w, `{"apikeys":[]}`)
				return
			}
			fmt.Fprint(w, `{"apikeys":[{"id":99,"name":"owner","key_type":"team","team_id":42,"team_name":"test_provider","key":null}]}`)
		case "POST /api/v0/auth/apikeys/":
			if !owner {
				fmt.Fprint(w, `{"id":100,"key":"fixture-team-management"}`)
				return
			}
			if keyExists {
				w.WriteHeader(http.StatusConflict)
				return
			}
			keyExists = true
			projectGenerations++
			projectID = int64(122 + projectGenerations)
			fmt.Fprintf(w, `{"id":%d,"key":"fixture-project-key-%d"}`, projectID, projectGenerations)
		case "PUT /api/v0/users/save-context/":
			if !owner {
				w.WriteHeader(403)
				return
			}
			fmt.Fprint(w, `{}`)
		case "GET /api/v0/users/current":
			if !owner {
				w.WriteHeader(403)
				return
			}
			fmt.Fprint(w, `{"id":42,"is_team":true}`)
		case "GET " + projectPath:
			if !owner {
				w.WriteHeader(403)
				return
			}
			if !keyExists {
				w.WriteHeader(404)
				return
			}
			fmt.Fprintf(w, `{"id":%d}`, projectID)
		case "DELETE " + projectPath:
			if !owner {
				w.WriteHeader(403)
				return
			}
			keyExists = false
			fmt.Fprint(w, `{"success":true}`)
		case "DELETE /api/v0/team/":
			if !owner || keyExists {
				w.WriteHeader(403)
				return
			}
			teamExists = false
			fmt.Fprint(w, `{"success":true}`)
		case "DELETE /api/v0/auth/apikeys/100/":
			fmt.Fprint(w, `{"success":true}`)
		default:
			w.WriteHeader(403)
		}
	}))
	defer srv.Close()
	h := newHarness(t, srv.URL+"/api/v0")
	config := fmt.Sprintf(`terraform {
 required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" { endpoint = %q }
resource "vastai_team" "test_team" { name = "test_provider" }
resource "vastai_api_key" "test_key" {
  name = "test_project"
  team_api_key = vastai_team.test_team.api_key
  permissions = jsonencode({api={user_read={}}})
}
output "project_key" {
 value = vastai_api_key.test_key.key
 sensitive = true
}
`, srv.URL+"/api/v0")
	if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	output := h.run(t, true, "apply", "-auto-approve", "-no-color")
	if strings.Contains(output, "fixture-project-key") || strings.Contains(output, "fixture-team-management") {
		t.Fatal("generated credential leaked into CLI output")
	}
	h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	for _, change := range []struct {
		// before is the current immutable setting.
		before string
		// after is the replacement setting.
		after string
	}{
		{`name = "test_project"`, `name = "test_project_updated"`},
		{`permissions = jsonencode({api={user_read={}}})`, `permissions = jsonencode({api={team_read={}}})`},
	} {
		config = strings.Replace(config, change.before, change.after, 1)
		if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(config), 0o600); err != nil {
			t.Fatal(err)
		}
		if plan := h.run(t, false, "plan", "-detailed-exitcode", "-no-color"); !strings.Contains(plan, "must be replaced") {
			t.Fatal("immutable API-key setting did not require replacement")
		}
		output = h.run(t, true, "apply", "-auto-approve", "-no-color")
		if strings.Contains(output, "fixture-project-key") {
			t.Fatal("replacement credential leaked into CLI output")
		}
		h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
	}
	h.run(t, true, "destroy", "-auto-approve", "-no-color")
	mu.Lock()
	defer mu.Unlock()
	if teamExists || keyExists {
		t.Fatal("bootstrap destroy left test team or token")
	}
	if projectGenerations != 3 {
		t.Fatal("API-key replacements did not revoke and recreate the key")
	}
	if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "team-bootstrap.json"), []byte(`{"result":"passed","team_deleted":true,"project_token_deleted":true,"credit_transfer":0,"project_key_generations":3}`), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
