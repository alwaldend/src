package main

import (
	"bytes"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

// storageFixture runs the command against an isolated TLS JSON-RPC service.
type storageFixture struct {
	// objects contains only synthetic inventory objects.
	objects map[string]map[string]any
	// creates counts storage creation operations.
	creates int
	// rejectCreate simulates a lost response after server-side creation.
	rejectCreate bool
	// rejectRead simulates an unavailable inventory service.
	rejectRead bool
	// wrongID simulates an inconsistent creation result.
	wrongID bool
}

// newStorageFixture prepares the same public command boundary used by Terraform.
func newStorageFixture(t *testing.T) *storageFixture {
	t.Helper()
	f := &storageFixture{objects: map[string]map[string]any{
		"host": {"id": "host", "type": "host", "name_label": "host.example.test", "$pool": "pool"},
	}}
	server := httptest.NewTLSServer(websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			var request struct {
				// ID correlates the response with its request.
				ID int `json:"id"`
				// Method selects the synthetic API operation.
				Method string `json:"method"`
				// Params holds the operation's arguments.
				Params map[string]any `json:"params"`
			}
			if websocket.JSON.Receive(ws, &request) != nil {
				return
			}
			var result, rpcError any
			switch request.Method {
			case "session.signInWithToken":
				result = map[string]any{"id": "fixture-user"}
			case "xo.getAllObjects":
				result = f.objects
				if f.rejectRead {
					rpcError = map[string]string{"message": "fixture-private-password"}
				}
			case "sr.createSmb":
				f.creates++
				if request.Params["host"] != "host" || request.Params["server"] != "//nas.example.test/cache" || request.Params["user"] != "cache-user" || request.Params["password"] != "fixture-private-password" {
					t.Error("incorrect storage connection parameters")
				}
				f.objects["sr"] = map[string]any{
					"id": "sr", "type": "SR", "name_label": "Runner cache",
					"name_description": request.Params["nameDescription"],
					"SR_type":          "smb", "shared": true, "$pool": "pool",
				}
				result = "sr"
				if f.rejectCreate {
					rpcError = map[string]string{"message": "fixture-private-password"}
				}
				if f.wrongID {
					result = "wrong-sr"
				}
			default:
				t.Errorf("unexpected API operation %q", request.Method)
			}
			if websocket.JSON.Send(ws, map[string]any{"jsonrpc": "2.0", "id": request.ID, "result": result, "error": rpcError}) != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)
	t.Setenv("XOA_URL", "wss"+strings.TrimPrefix(server.URL, "https"))
	t.Setenv("XOA_TOKEN", "fixture-private-token")
	t.Setenv("XOA_INSECURE", "true") // Only the isolated fixture's self-signed TLS.
	t.Setenv("XO_CACHE_HOST", "host.example.test")
	t.Setenv("XO_CACHE_SR", "Runner cache")
	t.Setenv("XO_CACHE_SMB_SHARE", "//nas.example.test/cache")
	t.Setenv("XO_CACHE_SMB_USER", "cache-user")
	t.Setenv("XO_CACHE_SMB_PASSWORD", "fixture-private-password")
	return f
}

// storageCommand checks command output and errors for credential exclusion.
func storageCommand(t *testing.T, name string, wantError bool) string {
	t.Helper()
	var out bytes.Buffer
	err := run([]string{name}, &out)
	if (err != nil) != wantError {
		t.Fatalf("command error = %v, want error %t", err, wantError)
	}
	if strings.Contains(out.String(), "fixture-private-") || (err != nil && strings.Contains(err.Error(), "fixture-private-")) {
		t.Fatal("credential escaped the command boundary")
	}
	return out.String()
}

// TestCacheStorageLifecycle verifies creation, repetition and partial recovery.
func TestCacheStorageLifecycle(t *testing.T) {
	f := newStorageFixture(t)
	storageCommand(t, "ensure-cache-storage", false)
	storageCommand(t, "ensure-cache-storage", false)
	if f.creates != 1 {
		t.Fatal("repeat invocation created additional storage")
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "cache-storage.json"), []byte("{\"create_operations\":1,\"repeat_operations\":0,\"transport\":\"isolated TLS JSON-RPC\"}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestCacheStorageFailures verifies that unsafe inventory cannot cause a create.
func TestCacheStorageFailures(t *testing.T) {
	for _, scenario := range []string{"missing-input", "missing-host", "duplicate-host", "duplicate-sr", "conflicting-sr", "read-error", "wrong-postcondition", "lost-response"} {
		t.Run(scenario, func(t *testing.T) {
			f := newStorageFixture(t)
			switch scenario {
			case "missing-input":
				t.Setenv("XO_CACHE_SMB_PASSWORD", "")
			case "missing-host":
				delete(f.objects, "host")
			case "duplicate-host":
				f.objects["host2"] = f.objects["host"]
			case "conflicting-sr", "duplicate-sr":
				f.objects["sr"] = map[string]any{"id": "sr", "type": "SR", "name_label": "Runner cache", "SR_type": "nfs", "$pool": "pool"}
				if scenario == "duplicate-sr" {
					f.objects["sr2"] = f.objects["sr"]
				}
			case "read-error":
				f.rejectRead = true
			case "wrong-postcondition":
				f.wrongID = true
			case "lost-response":
				f.rejectCreate = true
			}
			storageCommand(t, "ensure-cache-storage", true)
			if scenario == "lost-response" {
				f.rejectCreate = false
				storageCommand(t, "ensure-cache-storage", false)
			}
			want := 0
			if scenario == "wrong-postcondition" || scenario == "lost-response" {
				want = 1
			}
			if f.creates != want {
				t.Fatal("unexpected storage mutation count")
			}
		})
	}
}
