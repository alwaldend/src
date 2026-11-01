package protocol_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// bridgePath identifies the executable supplied through Bazel runfiles.
var bridgePath = flag.String("bridge", "", "Responses bridge runfile")

// codexPath enables opt-in acceptance with an installed Codex executable.
var codexPath = flag.String("codex-binary", "", "optional pinned Codex executable")

// startBridge launches the real command and returns its ephemeral loopback URL.
func startBridge(t *testing.T, upstream string) string {
	t.Helper()
	path, err := runfiles.Rlocation(*bridgePath)
	if err != nil {
		t.Fatalf("resolve bridge: %v", err)
	}
	cmd := exec.Command(path, "--listen", "127.0.0.1:0", "--upstream", upstream)
	var diagnostics bytes.Buffer
	cmd.Stderr = &diagnostics
	output, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open bridge output: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start bridge: %v", err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	ready := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(output)
		if scanner.Scan() {
			ready <- scanner.Text()
		} else {
			ready <- ""
		}
	}()
	select {
	case address := <-ready:
		if address == "" {
			t.Fatalf("bridge startup failed: %s", diagnostics.String())
		}
		return "http://" + address
	case <-time.After(10 * time.Second):
		t.Fatal("bridge startup timed out")
	}
	return ""
}

// post sends a Responses request and returns the complete response for assertions.
func post(t *testing.T, endpoint string, body any) (int, string) {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Post(endpoint+"/v1/responses", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("post response: %v", err)
	}
	defer response.Body.Close()
	result, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	return response.StatusCode, string(result)
}

// TestProtocol verifies namespace routing, streaming, history, and explicit rejection.
func TestProtocol(t *testing.T) {
	var requests atomic.Int32
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		tools := body["tools"].([]any)
		names := make(map[string]bool)
		for _, entry := range tools {
			tool := entry.(map[string]any)
			if tool["type"] != "function" {
				t.Errorf("backend received unsupported tool: %v", tool["type"])
			}
			name := tool["name"].(string)
			if names[name] {
				t.Errorf("duplicate backend name: %s", name)
			}
			names[name] = true
		}
		if len(tools) != 3 {
			t.Errorf("lost tools: %d", len(tools))
		}
		selected := tools[1].(map[string]any)["name"].(string)
		if body["tool_choice"].(map[string]any)["name"] != selected {
			t.Error("forced tool identity was lost")
		}
		if body["tool_choice"].(map[string]any)["namespace"] != nil {
			t.Error("namespace leaked to function-only backend")
		}
		call := body["input"].([]any)[0].(map[string]any)
		if call["name"] != selected || call["namespace"] != nil {
			t.Error("replayed call was not flattened")
		}
		item := map[string]any{"type": "function_call", "name": selected, "call_id": "call_1", "arguments": "{}"}
		if body["stream"] == true {
			w.Header().Set("Content-Type", "text/event-stream")
			for _, event := range []any{
				map[string]any{"type": "response.output_item.added", "item": item},
				map[string]any{"type": "response.function_call_arguments.done", "name": selected, "arguments": "{}"},
				map[string]any{"type": "response.completed", "response": map[string]any{"output": []any{item}}},
			} {
				data, _ := json.Marshal(event)
				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush()
			}
			fmt.Fprint(w, "data: [DONE]\n\n")
		} else {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"output": []any{item}})
		}
	}))
	defer backend.Close()
	endpoint := startBridge(t, backend.URL)
	for _, stream := range []bool{false, true} {
		body := map[string]any{
			"stream": stream,
			"tools": []any{
				map[string]any{"type": "function", "name": "read", "parameters": map[string]any{}},
				map[string]any{"type": "namespace", "name": "t3", "description": "Local tools", "tools": []any{map[string]any{"type": "function", "name": "read", "parameters": map[string]any{}}}},
				map[string]any{"type": "namespace", "name": "other", "tools": []any{map[string]any{"type": "function", "name": "read", "parameters": map[string]any{}}}},
			},
			"input":       []any{map[string]any{"type": "function_call", "namespace": "t3", "name": "read", "call_id": "call_0", "arguments": "{}"}},
			"tool_choice": map[string]any{"type": "function", "namespace": "t3", "name": "read"},
		}
		status, result := post(t, endpoint, body)
		if status != 200 {
			t.Fatalf("status %d: %s", status, result)
		}
		if !strings.Contains(result, `"namespace":"t3"`) || !strings.Contains(result, `"name":"read"`) {
			t.Fatalf("identity not restored: %s", result)
		}
		if stream && !strings.Contains(result, "data: [DONE]") {
			t.Fatal("stream did not complete")
		}
	}
	for _, tool := range []any{
		map[string]any{"type": "custom", "name": "unsupported"},
		map[string]any{"type": "namespace", "name": "x", "tools": []any{map[string]any{"type": "custom", "name": "unsupported"}}},
	} {
		status, _ := post(t, endpoint, map[string]any{"tools": []any{tool}})
		if status != 400 {
			t.Errorf("unsupported schema accepted: %d", status)
		}
	}
	if requests.Load() != 2 {
		t.Error("invalid request reached backend")
	}
	t.Log("JSON and SSE namespace round trips passed; unsupported schemas rejected")
}

// TestTransport verifies isolation between concurrent mappings and error paths.
func TestTransport(t *testing.T) {
	cancelled := make(chan struct{})
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			fmt.Fprint(w, `{"status":"ok"}`)
			return
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			return
		}
		if body["model"] == "error" {
			http.Error(w, "fixture model unavailable", 503)
			return
		}
		if body["model"] == "broken-stream" {
			w.Header().Set("Content-Type", "text/event-stream")
			fmt.Fprint(w, "data: invalid-json\n\n")
			return
		}
		if body["model"] == "cancel" {
			w.Header().Set("Content-Type", "text/event-stream")
			fmt.Fprint(w, ": waiting\n\n")
			w.(http.Flusher).Flush()
			<-r.Context().Done()
			close(cancelled)
			return
		}
		name := body["tools"].([]any)[0].(map[string]any)["name"]
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"output": []any{map[string]any{"type": "function_call", "name": name}}})
	}))
	defer backend.Close()
	endpoint := startBridge(t, backend.URL)
	var group sync.WaitGroup
	for i := range 8 {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			namespace := fmt.Sprintf("session_%d", index)
			status, body := post(t, endpoint, map[string]any{"tools": []any{map[string]any{"type": "namespace", "name": namespace, "tools": []any{map[string]any{"type": "function", "name": "read"}}}}})
			if status != 200 || !strings.Contains(body, `"namespace":"`+namespace+`"`) {
				t.Errorf("concurrent mapping crossed sessions: %d %s", status, body)
			}
		}(i)
	}
	group.Wait()
	status, body := post(t, endpoint, map[string]any{"model": "error"})
	if status != 503 || !strings.Contains(body, "fixture model unavailable") {
		t.Error("backend failure was hidden")
	}
	response, err := http.Post(endpoint+"/v1/responses", "application/json", strings.NewReader(`{"model":"broken-stream"}`))
	if err == nil {
		data, readErr := io.ReadAll(response.Body)
		response.Body.Close()
		if readErr == nil && response.StatusCode == 200 {
			t.Errorf("invalid stream accepted: %s", data)
		}
	}
	response, err = http.Get(endpoint + "/health")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 200 {
		t.Error("readiness did not reach backend")
	}
	response, err = http.Get(endpoint + "/other")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 404 {
		t.Error("unrelated upstream path exposed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/v1/responses", strings.NewReader(`{"model":"cancel"}`))
	if err != nil {
		t.Fatal(err)
	}
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	cancel()
	response.Body.Close()
	select {
	case <-cancelled:
	case <-time.After(5 * time.Second):
		t.Error("client cancellation did not reach backend")
	}
	t.Log("Concurrent identities, backend errors, malformed streams, health, and route restrictions passed")
}

// TestCodexMCP verifies actual Codex dispatch and the following inference request.
func TestCodexMCP(t *testing.T) {
	if *codexPath == "" {
		t.Skip("pass -codex-binary to exercise pinned Codex")
	}
	var called atomic.Bool
	mcp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
			return
		}
		var result any
		switch request["method"] {
		case "initialize":
			result = map[string]any{"protocolVersion": "2025-03-26", "serverInfo": map[string]any{"name": "local-fixture", "version": "1"}, "capabilities": map[string]any{"tools": map[string]any{}}}
		case "tools/list":
			result = map[string]any{"tools": []any{map[string]any{"name": "read_marker", "description": "Read the local acceptance marker", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}}}}
		case "tools/call":
			called.Store(true)
			result = map[string]any{"content": []any{map[string]any{"type": "text", "text": "local-tool-result"}}}
		default:
			if request["id"] == nil {
				w.WriteHeader(http.StatusAccepted)
				return
			}
			result = map[string]any{}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": request["id"], "result": result})
	}))
	defer mcp.Close()
	var turns atomic.Int32
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		var name string
		for _, entry := range body["tools"].([]any) {
			tool := entry.(map[string]any)
			if tool["type"] != "function" {
				t.Errorf("llama.cpp would discard tool type %v", tool["type"])
			}
			if strings.Contains(fmt.Sprint(tool["description"]), "acceptance marker") {
				name = tool["name"].(string)
			}
		}
		if name == "" {
			t.Error("MCP tool did not reach backend")
			w.WriteHeader(400)
			return
		}
		n := turns.Add(1)
		var output []any
		if n == 1 {
			output = []any{map[string]any{"type": "function_call", "id": "fc_1", "call_id": "call_1", "name": name, "arguments": "{}", "status": "completed"}}
		} else {
			for _, entry := range body["input"].([]any) {
				if call, ok := entry.(map[string]any); ok && call["type"] == "function_call" {
					if call["namespace"] != nil || call["name"] != name {
						t.Error("Codex replayed call identity was not flattened")
					}
				}
			}
			input, _ := json.Marshal(body["input"])
			if !called.Load() || !bytes.Contains(input, []byte("local-tool-result")) {
				t.Error("tool result missing from follow-up")
			}
			output = []any{map[string]any{"type": "message", "id": "msg_1", "role": "assistant", "status": "completed", "content": []any{map[string]any{"type": "output_text", "text": "LOCAL_ACCEPTANCE_OK", "annotations": []any{}}}}}
		}
		w.Header().Set("Content-Type", "text/event-stream")
		for _, item := range output {
			for _, typ := range []string{"response.output_item.added", "response.output_item.done"} {
				data, _ := json.Marshal(map[string]any{"type": typ, "output_index": 0, "item": item})
				fmt.Fprintf(w, "data: %s\n\n", data)
			}
		}
		data, _ := json.Marshal(map[string]any{"type": "response.completed", "response": map[string]any{"id": fmt.Sprintf("resp_%d", n), "object": "response", "status": "completed", "output": output}})
		fmt.Fprintf(w, "data: %s\n\n", data)
	}))
	defer backend.Close()
	endpoint := startBridge(t, backend.URL)
	home := t.TempDir()
	work := filepath.Join(home, "work")
	if err := os.Mkdir(work, 0o700); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(*codexPath, "exec", "--skip-git-repo-check", "--json", "-C", work,
		"-c", `model_provider="acceptance"`, "-c", `model="local-coder"`,
		"-c", fmt.Sprintf(`model_providers.acceptance={name="Local fixture",base_url=%q,wire_api="responses",requires_openai_auth=false}`, endpoint+"/v1"),
		"-c", fmt.Sprintf(`mcp_servers.acceptance={url=%q,tools.read_marker.approval_mode="approve"}`, mcp.URL),
		"-c", `features.remote_models=false`, "-c", `features.apps=false`, "-c", `features.plugins=false`,
		"-c", `web_search="disabled"`, "-c", `features.image_generation=false`,
		"-c", `approval_policy="never"`, "Read the marker using the acceptance MCP tool, then reply with its result.")
	cmd.Env = []string{"PATH=/usr/local/bin:/usr/bin:/bin", "HOME=" + home, "CODEX_HOME=" + home, "TMPDIR=" + home}
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Codex round trip: %v\n%s", err, output)
	}
	if turns.Load() != 2 || !called.Load() || !bytes.Contains(output, []byte("LOCAL_ACCEPTANCE_OK")) {
		t.Fatalf("incomplete round trip: turns=%d, called=%v\n%s", turns.Load(), called.Load(), output)
	}
	t.Log("Codex executed a namespaced local MCP tool and submitted its result; follow-up completed")
}
