package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/injector"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// fixtureProvider supplies a deliberately non-Vault resource over the public RPC.
type fixtureProvider struct {
	// UnimplementedPluginServiceServer supplies unused startup RPC defaults.
	al_proto.UnimplementedPluginServiceServer
	// UnimplementedResourceServiceServer supplies future resource RPC defaults.
	al_proto.UnimplementedResourceServiceServer
}

// PluginStart advertises types alongside optional generic startup output.
func (p *fixtureProvider) PluginStart(_ context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	mode := os.Getenv("RESOURCE_MODE")
	if mode == "plugin_named_injector" && (req.Name != "injector" || req.Plugin.Calls["ordinary_call"] == nil) {
		return nil, fmt.Errorf("lost plugin or selected call map key")
	}
	if mode == "startup_failure" {
		return nil, fmt.Errorf("synthetic provider startup failure")
	}
	if mode == "nil_startup_response" {
		return nil, nil
	}
	if mode == "empty_startup_response" {
		return &al_proto.PluginStartResponse{}, nil
	}
	response := &al_proto.PluginStartResponse{ResourceTypes: []string{"fixture"}}
	switch mode {
	case "collision":
		response.ResourceTypes = []string{"env"}
	case "empty_type":
		response.ResourceTypes = []string{""}
	case "no_types", "plugin_named_injector", "startup_exports":
		response.ResourceTypes = nil
	}
	if mode == "env_override" || mode == "startup_exports" {
		response.Result = &al_proto.ResourceResult{Env: map[string]string{"RESOURCE_VALUE": mode}}
	}
	if mode == "empty_startup_result" {
		response.Result = &al_proto.ResourceResult{}
	}
	if mode != "startup_result" {
		return response, nil
	}
	path := os.Getenv("RESOURCE_STARTUP_FILE")
	if err := os.WriteFile(path, []byte("synthetic startup file"), 0o600); err != nil {
		return nil, fmt.Errorf("write startup file: %w", err)
	}
	result, err := injector.ToProto(&injector.ResourceResult{
		Name: "fixture_provider", Env: map[string]string{"STARTUP_VALUE": "synthetic startup"}, Files: []string{path},
		Data: map[string]any{"number": json.Number("9007199254740993"), "nested": map[string]string{"value": "retained"}},
	})
	if err != nil {
		return nil, fmt.Errorf("encode startup result: %w", err)
	}
	response.Result = result
	return response, nil
}

// ResourceFetch validates cross-process dependency results and returns public data.
func (p *fixtureProvider) ResourceFetch(ctx context.Context, r *al_proto.ResourceFetchRequest) (*al_proto.ResourceResult, error) {
	if os.Getenv("RESOURCE_MODE") == "provider_exit" {
		os.Exit(42)
	}
	if os.Getenv("RESOURCE_MODE") == "fetch_failure" {
		return nil, fmt.Errorf("synthetic resource failure")
	}
	if os.Getenv("RESOURCE_MODE") == "resource_data" {
		input, err := al.FromPbJson(r.Data).Get()
		expected := map[string]any{"number": int64(9007199254740993), "nested": map[string]any{"items": []any{"request", true, nil}}, "empty": map[string]any{}}
		if err != nil || !reflect.DeepEqual(input, expected) {
			return nil, fmt.Errorf("lost configured provider data")
		}
	}
	if r.Type != "fixture" || len(r.Dependencies) != 1 || len(r.Dependencies[0].Files) != 1 {
		return nil, fmt.Errorf("lost resource dependencies")
	}
	data, err := os.ReadFile(r.Dependencies[0].Files[0])
	if err != nil || string(data) != "synthetic file" {
		return nil, fmt.Errorf("dependency file unavailable")
	}
	wire, err := injector.ToProto(&injector.ResourceResult{Name: r.Name, Data: map[string]any{"value": "synthetic provider", "number": json.Number("9007199254740993"), "nested": map[string]string{"value": "retained"}}})
	if err != nil {
		return nil, err
	}
	return wire, nil
}

// TestResourceHelper runs in a real provider or child process.
func TestResourceHelper(t *testing.T) {
	if os.Args[len(os.Args)-1] == "al-resource-process" {
		// Replace the Bazel test runner's signal handler with normal termination.
		signal.Reset(syscall.SIGTERM)
		if err := os.WriteFile(os.Getenv("RESOURCE_PROCESS_READY"), []byte("ready"), 0o600); err != nil {
			t.Fatal(err)
		}
		// Default SIGTERM termination must be accepted for an owned helper.
		select {}
	}
	switch os.Getenv("RESOURCE_ROLE") {
	case "provider":
		if path := os.Getenv("RESOURCE_STARTUP_FILE"); path != "" {
			defer os.Remove(path)
		}
		shutdown, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
		defer cancel()
		ctx := al.NewCmdCtx(shutdown, "")
		server := al_plugin.NewPluginServer(ctx, &fixtureProvider{})
		if err := server.Start(ctx.Ctx); err != nil {
			t.Fatal(err)
		}
		<-ctx.Ctx.Done()
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Stop(cleanup); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(os.Getenv("RESOURCE_STOPPED"), []byte("stopped"), 0o600); err != nil {
			t.Fatal(err)
		}
	case "child":
		signals := make(chan os.Signal, 1)
		if os.Getenv("RESOURCE_MODE") == "cancel" {
			signal.Notify(signals, syscall.SIGTERM)
			defer signal.Stop(signals)
		}
		data, err := os.ReadFile(os.Getenv("RESOURCE_FILE"))
		if err != nil || string(data) != "synthetic file" {
			t.Fatal("injected file unavailable")
		}
		if os.Getenv("RESOURCE_MODE") == "process" {
			deadline := time.Now().Add(5 * time.Second)
			for {
				if _, err := os.Stat(os.Getenv("RESOURCE_PROCESS_READY")); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("resource process did not become ready")
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
		if err := json.NewEncoder(os.Stdout).Encode(map[string]string{"value": os.Getenv("RESOURCE_VALUE"), "file": os.Getenv("RESOURCE_FILE"), "pid": os.Getenv("RESOURCE_PID")}); err != nil {
			t.Fatal(err)
		}
		if os.Getenv("RESOURCE_MODE") == "cancel" {
			if err := os.WriteFile(os.Getenv("RESOURCE_CHILD_READY"), []byte("ready"), 0o600); err != nil {
				t.Fatal(err)
			}
			<-signals
			if _, err := os.ReadFile(os.Getenv("RESOURCE_FILE")); err != nil {
				t.Fatal("file removed before child shutdown")
			}
			if _, err := os.Stat(os.Getenv("RESOURCE_STOPPED")); !os.IsNotExist(err) {
				t.Fatal("provider stopped before child shutdown")
			}
			if err := os.WriteFile(os.Getenv("RESOURCE_CHILD_STOPPED"), []byte("stopped"), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		os.Exit(0)
	}
}

// TestBuiltinResourcesAndProviders retains repeatable receipts for real process outcomes.
func TestBuiltinResourcesAndProviders(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	receipts := map[string]bool{}
	for _, mode := range []string{"success", "builtin_only", "builtin_data", "resource_data", "legacy_extra", "map_override", "map_equal", "map_ambiguous", "map_cycle", "map_missing_parent", "map_equal_sibling", "legacy_root_lists", "empty_key", "process", "unknown", "missing", "cycle", "collision", "duplicate", "empty_type", "no_types", "non_provider", "startup_failure", "invalid_builtin", "unknown_builtin_field", "plugin_named_injector", "env_override", "fetch_failure", "provider_exit", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			stopped := filepath.Join(dir, "provider-stopped")
			resources := map[string]any{
				"file":           map[string]any{"type": "file", "args": map[string]any{"value": "synthetic file"}},
				"custom":         map[string]any{"type": "fixture", "args": map[string]any{}, "deps": []any{"file"}},
				"RESOURCE_FILE":  map[string]any{"type": "env", "args": map[string]any{"value": "{{ index .Last.Files 0 }}"}, "deps": []any{"file"}},
				"RESOURCE_VALUE": map[string]any{"type": "env", "args": map[string]any{"value": "{{ .Last.Data.value }}:{{ .Last.Data.number }}:{{ .Last.Data.nested.value }}"}, "deps": []any{"custom"}},
			}
			if mode == "builtin_only" || mode == "builtin_data" || mode == "plugin_named_injector" {
				delete(resources, "custom")
				delete(resources, "RESOURCE_VALUE")
			}
			if mode == "builtin_data" || mode == "resource_data" {
				resources["file"].(map[string]any)["args"] = map[string]any{"value": "{{ .Data.content }}"}
				resources["file"].(map[string]any)["data"] = map[string]any{"content": "synthetic file"}
				if mode == "builtin_data" {
					resources["RESOURCE_VALUE"] = map[string]any{"type": "env", "args": map[string]any{"value": "{{ .Data.number }}:{{ index .Data.items 0 }}:{{ .Data.flag }}"}, "data": map[string]any{"number": int64(9007199254740993), "items": []any{"local", nil}, "flag": false}}
				} else {
					resources["custom"].(map[string]any)["data"] = map[string]any{"number": int64(9007199254740993), "nested": map[string]any{"items": []any{"request", true, nil}}, "empty": map[string]any{}}
					resources["RESOURCE_VALUE"].(map[string]any)["args"] = map[string]any{"value": "{{ .Last.Data.value }}:{{ .Last.Data.number }}:{{ .Last.Data.nested.value }}:{{ .Data.suffix }}"}
					resources["RESOURCE_VALUE"].(map[string]any)["data"] = map[string]any{"suffix": "local"}
				}
			}
			if mode == "legacy_extra" {
				resources["file"].(map[string]any)["args"] = map[string]any{"value": "synthetic file", "extra": map[string]any{"unused": "removed"}}
			}
			if mode == "process" {
				resources["helper"] = map[string]any{"type": "process", "args": map[string]any{"name": self, "args": []any{"-test.run=^TestResourceHelper$", "--", "al-resource-process"}}}
				resources["RESOURCE_PID"] = map[string]any{"type": "env", "args": map[string]any{"value": "{{ .Last.Data.pid }}"}, "deps": []any{"helper"}}
			}
			if mode == "unknown" {
				resources["custom"].(map[string]any)["type"] = "unknown"
			}
			if mode == "missing" {
				resources["custom"].(map[string]any)["deps"] = []any{"absent"}
			}
			if mode == "cycle" {
				resources["custom"].(map[string]any)["deps"] = []any{"custom"}
			}
			if mode == "invalid_builtin" {
				resources["RESOURCE_VALUE"].(map[string]any)["args"] = map[string]any{"value": true}
			}
			if mode == "unknown_builtin_field" {
				resources["RESOURCE_VALUE"].(map[string]any)["args"] = map[string]any{"unknown": "synthetic"}
			}
			providerName := "fixture_provider"
			provider := &al_proto.PluginConfig{Bin: self, Args: []string{"-test.run=^TestResourceHelper$"}, ResourceProvider: true, Env: map[string]string{"RESOURCE_ROLE": "provider", "RESOURCE_MODE": mode, "RESOURCE_STOPPED": stopped}}
			if mode == "non_provider" {
				provider.ResourceProvider = false
				provider.Labels = map[string]string{"fixture": "1"}
			}
			if mode == "builtin_only" || mode == "builtin_data" {
				provider.Bin = filepath.Join(dir, "must-not-start")
			}
			if mode == "plugin_named_injector" {
				providerName = "injector"
				provider.ResourceProvider = false
				provider.Labels = nil
			}
			cfg := &al_proto.Config{PluginConfigs: map[string]*al_proto.PluginConfig{providerName: provider}, Inject: map[string]*al_proto.InjectorCall{"resources": {Labels: map[string]string{"fixture": "1"}, Res: fixtureResources(t, resources)}}}
			if mode == "plugin_named_injector" {
				cfg.Plugins = map[string]*al_proto.PluginCall{"ordinary_call": {Plugin: "injector", Labels: map[string]string{"fixture": "1"}}}
			}
			if mode == "env_override" {
				other := *provider
				other.ResourceProvider = false
				other.Labels = map[string]string{"fixture": "1"}
				other.Env = map[string]string{"RESOURCE_ROLE": "provider", "RESOURCE_MODE": "startup_exports", "RESOURCE_STOPPED": stopped}
				cfg.PluginConfigs["startup_defaults"] = &other
			}
			if mode == "duplicate" {
				other := *provider
				cfg.PluginConfigs["other_provider"] = &other
			}
			if strings.HasPrefix(mode, "map_") {
				base := cfg.Inject["resources"]
				override := proto.CloneOf(base.Res["file"])
				cfg.Inject["a_override"] = &al_proto.InjectorCall{Labels: map[string]string{"fixture": "1"}, Res: map[string]*al_proto.ResourceConfig{"file": override}}
				switch mode {
				case "map_override", "map_equal_sibling":
					base.Res["file"].Args, err = al.ToPbJson(map[string]any{"value": "must be overridden"}).Get()
					if err != nil {
						t.Fatal(err)
					}
					cfg.Inject["middle"] = &al_proto.InjectorCall{Extends: "resources", Labels: map[string]string{"unselected": "1"}, Res: map[string]*al_proto.ResourceConfig{"must_not_fetch": {Type: "unknown"}}}
					cfg.Inject["a_override"].Extends = "middle"
					if mode == "map_equal_sibling" {
						cfg.Inject["sibling"] = &al_proto.InjectorCall{Labels: map[string]string{"fixture": "1"}, Res: map[string]*al_proto.ResourceConfig{"file": proto.CloneOf(base.Res["file"])}}
					}
				case "map_ambiguous":
					override.Args, err = al.ToPbJson(map[string]any{"value": "conflicting file"}).Get()
					if err != nil {
						t.Fatal(err)
					}
				case "map_cycle":
					base.Extends, cfg.Inject["a_override"].Extends = "a_override", "resources"
				case "map_missing_parent":
					base.Extends = "absent"
				}
			}
			if mode == "empty_key" {
				cfg.Inject[""] = cfg.Inject["resources"]
				delete(cfg.Inject, "resources")
			}
			raw, err := protojson.Marshal(cfg)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "legacy_root_lists" {
				var sections map[string]any
				if err := json.Unmarshal(raw, &sections); err != nil {
					t.Fatal(err)
				}
				sections["inject"] = []any{}
				raw, err = json.Marshal(sections)
				if err != nil {
					t.Fatal(err)
				}
			}
			path := filepath.Join(dir, "config.json")
			if err := os.WriteFile(path, raw, 0o600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, binary, "run", "--config", path, "--plugin_label", "fixture=1", "--", self, "-test.run=^TestResourceHelper$")
			cmd.Env = append(os.Environ(), "RESOURCE_ROLE=child", "TMPDIR="+dir, "RESOURCE_MODE="+mode, "RESOURCE_CHILD_READY="+filepath.Join(dir, "child-ready"), "RESOURCE_CHILD_STOPPED="+filepath.Join(dir, "child-stopped"), "RESOURCE_PROCESS_READY="+filepath.Join(dir, "process-ready"), "RESOURCE_STOPPED="+stopped)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if mode == "cancel" {
				cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGTERM) }
				result := make(chan error, 1)
				go func() { result <- cmd.Run() }()
				deadline := time.NewTimer(10 * time.Second)
				defer deadline.Stop()
				ticker := time.NewTicker(10 * time.Millisecond)
				defer ticker.Stop()
				ready := false
				for !ready {
					select {
					case <-ticker.C:
						if _, e := os.Stat(filepath.Join(dir, "child-ready")); e == nil {
							ready = true
						}
					case e := <-result:
						t.Fatalf("child exited before cancellation: %v\n%s", e, stderr.String())
					case <-deadline.C:
						t.Fatal("child did not become ready")
					}
				}
				cancel()
				err = <-result
				if _, err := os.Stat(filepath.Join(dir, "child-stopped")); err != nil {
					t.Fatalf("child did not complete graceful cleanup: %v\n%s", err, stderr.String())
				}
			} else {
				err = cmd.Run()
			}
			success := mode == "success" || mode == "map_override" || mode == "map_equal" || mode == "builtin_only" || mode == "builtin_data" || mode == "resource_data" || mode == "process" || mode == "plugin_named_injector" || mode == "env_override" || mode == "cancel"
			if success {
				if err != nil && mode != "cancel" {
					t.Fatalf("run: %v\n%s", err, stderr.String())
				}
				var result map[string]string
				if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				if (mode == "success" || mode == "env_override") && result["value"] != "synthetic provider:9007199254740993:retained" {
					t.Fatal("custom provider export missing")
				}
				if mode == "builtin_data" && result["value"] != "9007199254740993:local:false" {
					t.Fatal("built-in template data missing")
				}
				if mode == "resource_data" && result["value"] != "synthetic provider:9007199254740993:retained:local" {
					t.Fatal("local and dependency data were not independent")
				}
				if _, err := os.Stat(result["file"]); !os.IsNotExist(err) {
					t.Fatal("injected file not cleaned")
				}
				if mode == "process" {
					pid, err := strconv.Atoi(result["pid"])
					if err != nil || pid <= 0 {
						t.Fatal("resource process identity missing")
					}
					process, err := os.FindProcess(pid)
					if err != nil {
						t.Fatal(err)
					}
					defer process.Release()
					if err := process.Signal(syscall.Signal(0)); !errors.Is(err, syscall.ESRCH) && !errors.Is(err, os.ErrProcessDone) {
						t.Fatalf("resource process remains after shutdown: %v", err)
					}
				}
			} else if err == nil || stdout.Len() != 0 {
				t.Fatalf("invalid graph/provider started the child: %v", err)
			}
			expectedErrors := map[string]string{
				"unknown": "no resource provider", "missing": "missing dependency", "cycle": "dependency cycle",
				"collision": "duplicate resource type", "duplicate": "duplicate resource type", "empty_type": "empty advertised resource type",
				"no_types": "advertised no resource types", "non_provider": "without resource_provider",
				"startup_failure": "synthetic provider startup failure", "invalid_builtin": "invalid resource settings", "unknown_builtin_field": "invalid resource settings", "fetch_failure": "synthetic resource failure", "provider_exit": "fetch provider resource",
				"legacy_extra":  "invalid resource settings",
				"map_ambiguous": "conflicting resource", "map_equal_sibling": "conflicting resource", "map_cycle": "extension cycle", "map_missing_parent": "missing injector extension", "empty_key": "empty name", "legacy_root_lists": "could not unmarshal config",
			}
			if expected := expectedErrors[mode]; expected != "" && !strings.Contains(stderr.String(), expected) {
				t.Fatalf("unexpected failure for %s: %s", mode, stderr.String())
			}
			invalidMap := mode == "map_ambiguous" || mode == "map_equal_sibling" || mode == "map_cycle" || mode == "map_missing_parent" || mode == "empty_key" || mode == "legacy_root_lists"
			if mode != "builtin_only" && mode != "builtin_data" && mode != "missing" && mode != "cycle" && mode != "provider_exit" && !invalidMap {
				if _, err := os.Stat(stopped); err != nil {
					t.Fatalf("provider did not stop: %v\n%s", err, stderr.String())
				}
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range entries {
				if strings.HasPrefix(e.Name(), "al_file_") {
					t.Fatal("temporary file leaked")
				}
			}
			receipts[mode] = true
		})
	}
	writeResourceReceipts(t, "resource-provider-results.json", receipts)
}

// TestStartupResourceResults verifies optional results across an actual plugin process.
func TestStartupResourceResults(t *testing.T) {
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	receipts := map[string]bool{}
	for _, mode := range []string{"no_startup_result", "nil_startup_response", "empty_startup_response", "empty_startup_result", "startup_result"} {
		t.Run(mode, func(t *testing.T) {
			dir := t.TempDir()
			path, stopped := filepath.Join(dir, "startup-file"), filepath.Join(dir, "stopped")
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			cmdCtx := al.NewCmdCtx(ctx, "")
			defer cmdCtx.RequestShutdown()
			client, err := al_plugin.NewPluginClient(cmdCtx, nil, &al_proto.Config{}, "fixture_provider", &al_proto.PluginConfig{
				Bin: self, Args: []string{"-test.run=^TestResourceHelper$"},
				Env: map[string]string{"RESOURCE_ROLE": "provider", "RESOURCE_MODE": mode, "RESOURCE_STARTUP_FILE": path, "RESOURCE_STOPPED": stopped},
			})
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				cleanup, done := context.WithTimeout(context.Background(), 5*time.Second)
				defer done()
				if err := client.Stop(cleanup); err != nil {
					t.Error(err)
				}
			}()
			if err := client.Start(ctx); err != nil {
				t.Fatal(err)
			}
			response := client.StartResponse()
			result := response.GetResult()
			switch mode {
			case "no_startup_result", "nil_startup_response", "empty_startup_response":
				if result != nil {
					t.Fatal("absent startup output produced a result")
				}
			case "empty_startup_result":
				if result == nil || result.Name != "" || len(result.Env) != 0 || len(result.Files) != 0 || result.Data != nil {
					t.Fatal("explicitly empty startup result lost its presence")
				}
			default:
				value, err := injector.FromProto(result)
				if err != nil {
					t.Fatal(err)
				}
				data := value.Data.(map[string]any)
				if value.Name != "fixture_provider" || value.Env["STARTUP_VALUE"] != "synthetic startup" || len(value.Files) != 1 || value.Files[0] != path || data["number"] != int64(9007199254740993) || data["nested"].(map[string]any)["value"] != "retained" {
					t.Fatal("startup result fields changed in transport")
				}
				content, err := os.ReadFile(path)
				if err != nil || string(content) != "synthetic startup file" {
					t.Fatal("startup file unavailable before shutdown")
				}
			}
			types := response.GetResourceTypes()
			if mode == "nil_startup_response" || mode == "empty_startup_response" {
				if len(types) != 0 {
					t.Fatal("empty startup response advertised resources")
				}
			} else if len(types) != 1 || types[0] != "fixture" {
				t.Fatal("startup advertisements missing")
			}
			if err := client.Stop(ctx); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatal("startup file survived provider shutdown")
			}
			if _, err := os.Stat(stopped); err != nil {
				t.Fatal(err)
			}
			receipts[mode] = true
		})
	}
	writeResourceReceipts(t, "startup-resource-results.json", receipts)
}

// writeResourceReceipts retains repeatable outcomes from process fixtures.
func writeResourceReceipts(t *testing.T, name string, receipts map[string]bool) {
	t.Helper()
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		t.Fatal("missing receipt directory")
	}
	data, err := json.MarshalIndent(receipts, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, name), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// fixtureResources encodes named fixture definitions through the public config schema.
func fixtureResources(t *testing.T, resources map[string]any) map[string]*al_proto.ResourceConfig {
	t.Helper()
	result := map[string]*al_proto.ResourceConfig{}
	for name, raw := range resources {
		fields := raw.(map[string]any)
		args, err := al.ToPbJson(fields["args"]).Get()
		if err != nil {
			t.Fatal(err)
		}
		r := &al_proto.ResourceConfig{Type: fields["type"].(string), Args: args}
		if value, ok := fields["data"]; ok {
			r.Data, err = al.ToPbJson(value).Get()
			if err != nil {
				t.Fatal(err)
			}
		}
		if deps, ok := fields["deps"].([]any); ok {
			for _, dep := range deps {
				r.Deps = append(r.Deps, dep.(string))
			}
		}
		result[name] = r
	}
	return result
}
