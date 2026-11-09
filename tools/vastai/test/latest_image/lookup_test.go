package latest_image_test

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// binary identifies the built external-data adapter.
var binary = flag.String("binary", "", "Image lookup executable runfile")

// TestLookup exercises the external-data protocol against a bounded registry fixture.
func TestLookup(t *testing.T) {
	executable, err := runfiles.Rlocation(*binary)
	if err != nil {
		t.Fatalf("resolve adapter: %v", err)
	}
	for name, response := range map[string]string{
		"newest":       `{"results":[{"name":"neo-abcd123-2026-11-09-cuda-12.9"},{"name":"neo-deadbee-2026-11-02-cuda-12.9"}]}`,
		"wrong-family": `{"results":[{"name":"reforge-abcd123-2026-11-09-cuda-12.9"}]}`,
		"malformed":    `{"results":`,
	} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(response)) }))
			defer server.Close()
			input, _ := json.Marshal(map[string]string{"endpoint": server.URL})
			cmd := exec.Command(executable)
			cmd.Stdin = strings.NewReader(string(input))
			output, err := cmd.CombinedOutput()
			if name != "newest" {
				if err == nil {
					t.Fatal("invalid registry metadata accepted")
				}
				return
			}
			if err != nil {
				t.Fatalf("lookup protocol: %v: %s", err, output)
			}
			var result map[string]string
			if err := json.Unmarshal(output, &result); err != nil {
				t.Fatalf("decode result: %v", err)
			}
			if result["image"] != "vastai/sd-forge:neo-abcd123-2026-11-09-cuda-12.9" {
				t.Fatalf("incorrect image: %s", output)
			}
			if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
				if err := os.WriteFile(filepath.Join(directory, "latest-image.json"), output, 0o600); err != nil {
					t.Fatalf("write evidence: %v", err)
				}
			}
		})
	}
}
