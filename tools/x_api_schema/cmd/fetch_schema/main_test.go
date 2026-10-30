package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// TestFetchCommand exercises source replacement and rejection through the CLI.
func TestFetchCommand(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("FETCH_SCHEMA_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	// validSchema is a minimal schema with an observable operation.
	const validSchema = `{"openapi":"3.0.0","info":{"title":"Fixture","version":"1"},"paths":{"/users":{"get":{"responses":{"200":{"description":"OK"}}}}}}`
	for _, scenario := range []struct {
		// name identifies the failure boundary.
		name string
		// status is the server response status.
		status int
		// body is the server response content.
		body string
		// success permits replacing the existing schema.
		success bool
	}{
		{"success", 200, validSchema, true},
		{"unauthorized", 401, "credential-value", false},
		{"redirect", 302, validSchema, false},
		{"invalid-json", 200, "<html>error</html>", false},
		{"not-openapi", 200, `{"error":"failure"}`, false},
		{"empty-paths", 200, `{"openapi":"3.0.0","info":{},"paths":{}}`, false},
		{"oversized", 200, strings.Repeat("x", 1025), false},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if !strings.HasPrefix(r.Header.Get("Authorization"), "OAuth ") {
					t.Error("missing OAuth 1.0a signature")
				}
				w.Header().Set("Location", "/redirected")
				w.WriteHeader(scenario.status)
				_, _ = w.Write([]byte(scenario.body))
			}))
			defer server.Close()
			out := filepath.Join(t.TempDir(), "openapi.json")
			if err := os.WriteFile(out, []byte("previous schema"), 0o644); err != nil {
				t.Fatal(err)
			}
			command := exec.Command(binary, "--url", server.URL, "--out", out, "--max-bytes", "1024")
			command.Env = append(os.Environ(), "X_API_KEY=fixture-key", "X_API_SECRET=credential-value", "X_ACCESS_TOKEN=fixture-token", "X_ACCESS_TOKEN_SECRET=fixture-secret")
			output, err := command.CombinedOutput()
			if (err == nil) != scenario.success {
				t.Fatalf("success=%v, error=%v, output=%s", scenario.success, err, output)
			}
			if strings.Contains(string(output), "credential-value") {
				t.Fatal("response body or credential leaked into diagnostics")
			}
			if requests != 1 {
				t.Fatalf("issued %d requests; want one without retries or redirects", requests)
			}
			got, err := os.ReadFile(out)
			if err != nil {
				t.Fatal(err)
			}
			want := "previous schema"
			if scenario.success {
				want = validSchema + "\n"
			}
			if string(got) != want {
				t.Fatal("unexpected source bytes")
			}
			if scenario.success {
				if err := os.WriteFile(filepath.Join(os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), "openapi.json"), got, 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
