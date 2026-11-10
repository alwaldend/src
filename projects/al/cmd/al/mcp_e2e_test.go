package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpTestTransport exposes transport disconnect independently of graceful client shutdown.
type mcpTestTransport struct {
	// transport starts the real AL command over pipes.
	transport mcp.CommandTransport
	// connection allows the fixture to simulate physical input closure.
	connection mcp.Connection
}

// Connect retains the SDK transport connection for disconnect verification.
func (transport *mcpTestTransport) Connect(ctx context.Context) (mcp.Connection, error) {
	connection, err := transport.transport.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect fixture transport: %w", err)
	}
	transport.connection = connection
	return connection, nil
}

// TestBazelMCPHelper substitutes the Bazel process while AL and MCP remain real.
func TestBazelMCPHelper(t *testing.T) {
	if os.Getenv("AL_MCP_HELPER") != "1" {
		return
	}
	separator := slices.Index(os.Args, "--")
	if separator < 0 {
		os.Exit(90)
	}
	args := os.Args[separator+1:]
	switch {
	case slices.Contains(args, "wait"):
		if err := os.WriteFile(os.Getenv("AL_MCP_READY"), []byte("ready"), 0o600); err != nil {
			os.Exit(91)
		}
		time.Sleep(30 * time.Second)
	case slices.Contains(args, "fail"):
		fmt.Fprintln(os.Stderr, "fixture operation failed")
		fmt.Println("partial result")
		os.Exit(3)
	case slices.Contains(args, "large"):
		fmt.Print(strings.Repeat("o", 4096))
		fmt.Fprint(os.Stderr, strings.Repeat("e", 4096))
	case slices.Contains(args, "escaped"):
		fmt.Print(strings.Repeat("\x00", 524288))
		fmt.Fprint(os.Stderr, strings.Repeat("\x00", 524288))
	default:
		cwd, err := os.Getwd()
		if err != nil {
			os.Exit(92)
		}
		if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"workspace": cwd, "args": args}); err != nil {
			os.Exit(93)
		}
	}
	os.Exit(0)
}

// TestBazelMCPCLI verifies discovery, execution boundaries, failures and cleanup over stdio.
func TestBazelMCPCLI(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "MODULE.bazel"), []byte("module(name = \"fixture\")\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var observations []map[string]any
	for _, scenario := range []string{"operations", "invalid", "failure", "truncation", "escaped_output", "timeout", "cancellation", "disconnect"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			ready := filepath.Join(t.TempDir(), "ready")
			deadline := "20s"
			if scenario == "timeout" {
				deadline = "1s"
			}
			outputLimit := "1024"
			if scenario == "escaped_output" {
				outputLimit = "524288"
			}
			command := exec.CommandContext(ctx, binary, "bazel", "mcp", "serve", "--workspace", workspace,
				"--bazel-path", helper, "--startup-option=-test.run=TestBazelMCPHelper", "--startup-option=--",
				"--timeout="+deadline, "--max-output-bytes="+outputLimit)
			command.Env = append(os.Environ(), "AL_MCP_HELPER=1", "AL_MCP_READY="+ready)
			client := mcp.NewClient(&mcp.Implementation{Name: "al-e2e", Version: "1"}, nil)
			transport := &mcpTestTransport{transport: mcp.CommandTransport{Command: command, TerminateDuration: 2 * time.Second}}
			session, err := client.Connect(ctx, transport, nil)
			if err != nil {
				t.Fatal(err)
			}
			closed := false
			defer func() {
				if !closed {
					if err := session.Close(); err != nil {
						t.Errorf("close server: %v", err)
					}
				}
			}()
			call := func(name string, arguments any) (*mcp.CallToolResult, map[string]any) {
				t.Helper()
				result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
				if err != nil {
					t.Fatal(err)
				}
				data, err := json.Marshal(result.StructuredContent)
				if err != nil {
					t.Fatal(err)
				}
				var report map[string]any
				if err := json.Unmarshal(data, &report); err != nil {
					t.Fatal(err)
				}
				observations = append(observations, map[string]any{"scenario": scenario, "tool": name, "result": result})
				return result, report
			}
			switch scenario {
			case "operations":
				catalog, err := session.ListTools(ctx, nil)
				if err != nil || len(catalog.Tools) != 6 {
					t.Fatalf("tool catalog: %v; %+v", err, catalog)
				}
				for _, operation := range []string{"query", "cquery", "aquery", "build", "test", "info"} {
					arguments := map[string]any{}
					if strings.HasSuffix(operation, "query") {
						arguments["expression"] = "deps(//:literal with spaces)"
					} else if operation != "info" {
						arguments["targets"] = []string{"//:target with spaces"}
					}
					result, report := call("bazel_"+operation, arguments)
					if result.IsError || report["exit_code"] != float64(0) || report["workspace"] != workspace {
						t.Fatalf("unexpected operation result: %+v", result)
					}
					var child struct {
						// Workspace identifies the process working directory.
						Workspace string `json:"workspace"`
						// Args preserves the child argument vector.
						Args []string `json:"args"`
					}
					if err := json.Unmarshal([]byte(report["stdout"].(string)), &child); err != nil {
						t.Fatal(err)
					}
					if child.Workspace != workspace || child.Args[0] != operation {
						t.Fatalf("incorrect execution boundary: %+v", child)
					}
					if operation == "build" && !slices.Contains(child.Args, "//:target with spaces") {
						t.Fatalf("target was split: %v", child.Args)
					}
				}
			case "invalid":
				for _, input := range []struct {
					// name selects the invalid tool call.
					name string
					// arguments exercises schema and semantic validation.
					arguments any
				}{
					{"bazel_build", map[string]any{"targets": []string{"--help"}}},
					{"bazel_test", map[string]any{"targets": []string{}}},
					{"bazel_query", map[string]any{"expression": ""}},
					{"bazel_query", map[string]any{"expression": "--help"}},
					{"bazel_query", map[string]any{"expression": "//...", "options": []string{"build"}}},
					{"bazel_query", map[string]any{"expression": "//...", "options": []string{"--"}}},
					{"bazel_query", map[string]any{"expression": 42}},
					{"bazel_run", map[string]any{}},
				} {
					result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: input.name, Arguments: input.arguments})
					if err == nil && !result.IsError {
						t.Errorf("invalid input accepted: %+v", input)
					}
				}
			case "failure":
				result, report := call("bazel_query", map[string]any{"expression": "fail"})
				if !result.IsError || report["exit_code"] != float64(3) || !strings.Contains(report["stderr"].(string), "failed") {
					t.Fatalf("failure lost: %+v", result)
				}
			case "truncation":
				result, report := call("bazel_query", map[string]any{"expression": "large"})
				if !result.IsError || report["stdout_truncated"] != true || report["stderr_truncated"] != true || len(report["stdout"].(string)) != 1024 || len(report["stderr"].(string)) != 1024 {
					t.Fatalf("output not bounded: %+v", result)
				}
			case "timeout":
				result, report := call("bazel_query", map[string]any{"expression": "wait"})
				if !result.IsError || report["timed_out"] != true {
					t.Fatalf("deadline not reported: %+v", result)
				}
			case "escaped_output":
				result, report := call("bazel_query", map[string]any{"expression": "escaped"})
				if result.IsError || len(report["stdout"].(string)) != 524288 || len(report["stderr"].(string)) != 524288 {
					t.Fatal("maximum escaped output did not fit the client frame limit")
				}
			case "cancellation", "disconnect":
				requestCtx, requestCancel := context.WithCancel(ctx)
				defer requestCancel()
				done := make(chan error, 1)
				go func() {
					_, err := session.CallTool(requestCtx, &mcp.CallToolParams{Name: "bazel_query", Arguments: map[string]any{"expression": "wait"}})
					done <- err
				}()
				for {
					if _, err := os.Stat(ready); err == nil {
						break
					}
					select {
					case err := <-done:
						t.Fatalf("request ended before readiness: %v", err)
					case <-ctx.Done():
						t.Fatal(ctx.Err())
					case <-time.After(5 * time.Millisecond):
					}
				}
				if scenario == "cancellation" {
					requestCancel()
					if err := <-done; err == nil {
						t.Fatal("cancelled request succeeded")
					}
					result, _ := call("bazel_info", map[string]any{})
					if result.IsError {
						t.Fatalf("server did not recover: %+v", result)
					}
				} else {
					if err := transport.connection.Close(); err != nil {
						t.Fatal(err)
					}
					closed = true
					if !command.ProcessState.Success() {
						t.Fatalf("disconnect failed: %s", command.ProcessState)
					}
				}
			}
		})
	}
	for _, args := range [][]string{
		{"--workspace", t.TempDir()},
		{"--workspace", workspace, "--timeout=0s"},
		{"--workspace", workspace, "--max-output-bytes=0"},
		{"--workspace", workspace, "--max-output-bytes=524289"},
	} {
		command := exec.Command(binary, append([]string{"bazel", "mcp", "serve"}, args...)...)
		if output, err := command.CombinedOutput(); err == nil {
			t.Errorf("invalid server options accepted: %v; %s", args, output)
		}
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(observations, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "bazel-mcp.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
