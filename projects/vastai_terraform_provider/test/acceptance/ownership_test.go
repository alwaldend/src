package acceptance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// ownershipFixture keeps management keys alive after external team deletion.
type ownershipFixture struct {
	// mu protects observable API state.
	mu sync.Mutex
	// teamExists identifies the team whose ownership must be preserved.
	teamExists bool
	// teamName identifies this fixture's independent team.
	teamName string
	// accountID distinguishes teams sharing one personal bootstrap account.
	accountID int64
	// nextKey assigns a distinct credential to each attempted creation.
	nextKey int64
	// keys records live independent management keys.
	keys map[int64]bool
	// contexts records the account selected on each generated key.
	contexts map[int64]int64
	// revokeFailures makes key revocation fail before a successful retry.
	revokeFailures int
	// malformedReads omits key identifiers to distinguish invalid data from absence.
	malformedReads bool
	// calls records non-sensitive operations for repeatable evidence.
	calls []string
}

// ServeHTTP models duplicate-name rejection and surviving management credentials.
func (f *ownershipFixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.Method+" "+r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	accountID := f.accountID
	if accountID == 0 {
		accountID = 42
	}
	personal := r.Header.Get("Authorization") == "Bearer fixture-credential"
	keyID, _ := strconv.ParseInt(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer fixture-management-"), 10, 64)
	if !personal && !f.keys[keyID] {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	switch r.Method + " " + r.URL.Path {
	case "POST /api/v0/auth/apikeys/":
		if !personal && f.contexts[keyID] != accountID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		id := f.nextKey
		f.nextKey++
		f.keys[id] = true
		fmt.Fprintf(w, `{"id":%d,"key":"fixture-management-%d"}`, id, id)
	case "POST /api/v0/team/":
		if personal || f.teamExists {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"msg":"Team name already exists"}`)
			return
		}
		f.teamExists = true
		fmt.Fprint(w, `{"success":true}`)
	case "GET /api/v0/auth/apikeys/":
		if f.malformedReads {
			fmt.Fprint(w, `{"apikeys":[{}]}`)
			return
		}
		if f.teamExists {
			name := f.teamName
			if name == "" {
				name = "test_provider"
			}
			fmt.Fprintf(w, `{"apikeys":[{"id":%d,"name":"owner","team_id":%d,"team_name":%q}]}`, accountID+57, accountID, name)
		} else {
			fmt.Fprint(w, `{"apikeys":[]}`)
		}
	case "PUT /api/v0/users/save-context/":
		if personal || !f.teamExists {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		f.contexts[keyID] = accountID
		fmt.Fprint(w, `{}`)
	case "GET /api/v0/users/current":
		fmt.Fprintf(w, `{"id":%d,"is_team":%t}`, f.contexts[keyID], f.contexts[keyID] == accountID)
	case "DELETE /api/v0/team/":
		if personal || !f.teamExists || f.contexts[keyID] != accountID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		f.teamExists = false
		fmt.Fprint(w, `{"success":true}`)
	default:
		if (!personal && f.contexts[keyID] != accountID) || !strings.HasPrefix(r.URL.Path, "/api/v0/auth/apikeys/") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		id, err := strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v0/auth/apikeys/"), "/"), 10, 64)
		if err != nil || !f.keys[id] {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, `{"id":%d}`, id)
			return
		}
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if f.revokeFailures > 0 {
			f.revokeFailures--
			w.WriteHeader(http.StatusForbidden)
			return
		}
		delete(f.keys, id)
		fmt.Fprint(w, `{"success":true}`)
	}
}

// TestTeamOwnership verifies failed-create cleanup and explicit missing-team replacement.
func TestTeamOwnership(t *testing.T) {
	for _, mode := range []string{"duplicate-name", "external-deletion", "revocation-failure", "malformed-context"} {
		t.Run(mode, func(t *testing.T) {
			f := &ownershipFixture{teamExists: mode == "duplicate-name", nextKey: 100, keys: map[int64]bool{}, contexts: map[int64]int64{}}
			srv := httptest.NewServer(f)
			defer srv.Close()
			h := newHarness(t, srv.URL+"/api/v0")
			config := fmt.Sprintf(`terraform {
  required_providers { vastai = { source = "alwaldend/vastai" } }
}
provider "vastai" { endpoint = %q }
resource "vastai_team" "test_team" { name = "test_provider" }
`, srv.URL+"/api/v0")
			if err := os.WriteFile(filepath.Join(h.directory, "main.tf"), []byte(config), 0o600); err != nil {
				t.Fatal(err)
			}
			h.run(t, mode != "duplicate-name", "apply", "-auto-approve", "-no-color")
			if mode == "malformed-context" {
				f.mu.Lock()
				f.malformedReads = true
				f.mu.Unlock()
			} else if mode != "duplicate-name" {
				f.mu.Lock()
				f.teamExists = false
				f.contexts[100] = 0
				if mode == "revocation-failure" {
					f.revokeFailures = 1
				}
				f.mu.Unlock()
			}
			h.run(t, mode != "malformed-context", "apply", "-refresh-only", "-auto-approve", "-no-color")
			if state := h.run(t, true, "state", "list"); !strings.Contains(state, "vastai_team.test_team") {
				t.Fatal("refresh abandoned management-key ownership")
			}
			f.mu.Lock()
			if !f.keys[100] {
				t.Error("refresh mutated the management key")
			}
			f.mu.Unlock()
			if mode == "malformed-context" {
				f.mu.Lock()
				f.malformedReads = false
				f.mu.Unlock()
				h.run(t, true, "plan", "-detailed-exitcode", "-no-color")
			} else if mode != "duplicate-name" {
				if mode == "revocation-failure" {
					h.run(t, false, "apply", "-auto-approve", "-no-color")
					if state := h.run(t, true, "state", "list"); !strings.Contains(state, "vastai_team.test_team") {
						t.Fatal("failed revocation abandoned the original key")
					}
				}
				h.run(t, true, "apply", "-auto-approve", "-no-color")
				f.mu.Lock()
				if !f.teamExists || f.keys[100] || len(f.keys) != 1 {
					t.Error("replacement failed to revoke the original management key before recreation")
				}
				f.mu.Unlock()
			}
			h.run(t, true, "destroy", "-auto-approve", "-no-color")
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.teamExists != (mode == "duplicate-name") || len(f.keys) != 0 {
				t.Fatal("cleanup deleted an unowned team or left a generated key")
			}
			report, err := json.MarshalIndent(map[string]any{"result": "passed", "scenario": mode, "calls": f.calls}, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); dir != "" {
				if err := os.WriteFile(filepath.Join(dir, "team-ownership-"+mode+".json"), report, 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
