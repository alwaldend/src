package main

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
	"golang.org/x/net/websocket"
)

// TestCLIHTTPS executes the packaged CLI and bridge against an isolated API fixture.
func TestCLIHTTPS(t *testing.T) {
	bridge, err := runfiles.Rlocation(os.Getenv("BRIDGE_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	cli, err := runfiles.Rlocation(os.Getenv("CLI_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	queued := make(chan struct{})
	var queueOnce sync.Once
	var mu sync.Mutex
	var requests []string
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || user != "fixture" || password != "fixture-password" {
			http.Error(w, "unauthorized", 401)
			return
		}
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/system_stats":
			fmt.Fprint(w, `{"system":{"comfyui_version":"fixture"},"devices":[]}`)
		case "/object_info":
			fmt.Fprint(w, `{"FixtureOutput":{"name":"FixtureOutput","input":{"required":{}},"output":[],"output_name":[],"output_node":true}}`)
		case "/prompt":
			queueOnce.Do(func() { close(queued) })
			fmt.Fprint(w, `{"prompt_id":"fixture-job","number":1,"node_errors":{}}`)
		case "/history/fixture-job":
			fmt.Fprint(w, `{"fixture-job":{"prompt":[1,"fixture-job",{},{}],"outputs":{"1":{"images":[{"filename":"fixture.png","subfolder":"","type":"output"}]}},"status":{"status_str":"success","completed":true,"messages":[]}}}`)
		case "/queue":
			fmt.Fprint(w, `{"queue_running":[],"queue_pending":[]}`)
		case "/history":
			fmt.Fprint(w, `{}`)
		case "/view":
			w.Header().Set("Content-Type", "image/png")
			fmt.Fprint(w, "neutral-fixture-image")
		default:
			http.NotFound(w, r)
		}
	})
	mux.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		user, password, ok := ws.Request().BasicAuth()
		if !ok || user != "fixture" || password != "fixture-password" {
			return
		}
		select {
		case <-queued:
		case <-time.After(15 * time.Second):
			return
		}
		if err := websocket.Message.Send(ws, `{"type":"executed","data":{"node":"1","prompt_id":"fixture-job","output":{"images":[{"filename":"fixture.png","subfolder":"","type":"output"}]}}}`); err != nil {
			t.Error(err)
		}
		if err := websocket.Message.Send(ws, `{"type":"executing","data":{"node":null,"prompt_id":"fixture-job"}}`); err != nil {
			t.Error(err)
		}
	}))
	server := httptest.NewTLSServer(mux)
	defer server.Close()
	ca := filepath.Join(directory, "ca.pem")
	if err := os.WriteFile(ca, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw}), 0o600); err != nil {
		t.Fatal(err)
	}
	workflow := filepath.Join(directory, "workflow.json")
	if err := os.WriteFile(workflow, []byte(`{"1":{"class_type":"FixtureOutput","inputs":{}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	call := func(success bool, endpoint, password string, args ...string) string {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		argv := append([]string{"--cli", cli, "--url", endpoint, "--ca-file", ca, "--"}, args...)
		if endpoint == "" {
			self, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			argv = append([]string{"--cli", cli, "--endpoint-command", self, "--endpoint-arg=-test.run=^TestEndpointFixture$", "--ca-file", ca, "--"}, args...)
		}
		cmd := exec.CommandContext(ctx, bridge, argv...)
		cmd.Env = append(os.Environ(), "IMAGEGEN_WEB_USERNAME=fixture", "IMAGEGEN_WEB_PASSWORD="+password, "COMFY_REPOSITORY_PROFILE="+filepath.Join(directory, "profile"), "COMFY_WHERE=cloud", "ENDPOINT_FIXTURE_URL="+server.URL)
		output, err := cmd.CombinedOutput()
		if (err == nil) != success {
			t.Fatalf("CLI success=%t: %v\n%s", success, err, output)
		}
		if strings.Contains(string(output), "fixture-password") {
			t.Fatal("credential disclosed")
		}
		return string(output)
	}
	call(true, server.URL, "fixture-password", "--json", "system-stats")
	call(true, "", "fixture-password", "--json", "system-stats")
	call(false, server.URL, "", "--json", "system-stats")
	call(false, server.URL, "incorrect", "--json", "system-stats")
	call(false, server.URL, "fixture-password", "--json", "system-stats", "--where", "cloud")
	call(false, server.URL, "fixture-password", "run", "--host", "example.invalid")
	call(false, strings.Replace(server.URL, "https://127.0.0.1:", "https://localhost:", 1), "fixture-password", "--json", "system-stats")
	call(true, server.URL, "fixture-password", "run", "--workflow", workflow, "--wait", "--timeout", "10", "--json")
	call(true, server.URL, "fixture-password", "--json", "jobs", "status", "fixture-job")
	outputs := filepath.Join(directory, "outputs")
	call(true, server.URL, "fixture-password", "--json", "download", "fixture-job", "-o", outputs)
	files, err := filepath.Glob(filepath.Join(outputs, "*.png"))
	if err != nil || len(files) != 1 {
		t.Fatalf("output files: %v %v", files, err)
	}
	image, err := os.ReadFile(files[0])
	if err != nil || string(image) != "neutral-fixture-image" {
		t.Fatalf("output transfer: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, request := range []string{"GET /system_stats", "POST /prompt", "GET /history/fixture-job", "GET /view"} {
		found := false
		for _, observed := range requests {
			if observed == request {
				found = true
			}
		}
		if !found {
			t.Errorf("missing API operation %s", request)
		}
	}
	receipt, err := json.MarshalIndent(map[string]any{"cli_version": "1.25.0", "authenticated_https": true, "workflow_completed": true, "output_downloaded": true, "requests": requests}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "cli-https-result.json"), receipt, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestEndpointFixture emulates the owning runner's initialization and raw output.
func TestEndpointFixture(t *testing.T) {
	endpoint := os.Getenv("ENDPOINT_FIXTURE_URL")
	if endpoint == "" {
		t.Skip("subprocess fixture only")
	}
	if os.Getenv("IMAGEGEN_WEB_PASSWORD") != "" || os.Getenv("IMAGEGEN_WEB_USERNAME") != "" {
		fmt.Fprintln(os.Stderr, "endpoint subprocess inherited web credentials")
		os.Exit(1)
	}
	fmt.Printf("Initializing the backend...\nTerraform has been successfully initialized!\n%s", endpoint)
	os.Exit(0)
}
