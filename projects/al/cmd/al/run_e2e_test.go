package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/encoding/protojson"
)

// runFixturePlugin returns synthetic environment exports through the real plugin RPC.
type runFixturePlugin struct {
	// UnimplementedPluginServiceServer supplies the forward-compatible RPC methods.
	al_proto.UnimplementedPluginServiceServer
	// exports contains only synthetic environment values for this fixture.
	exports map[string]string
	// ready records successful plugin initialization before child selection.
	ready string
}

// PluginStart records startup and returns the configured synthetic exports.
func (p *runFixturePlugin) PluginStart(context.Context, *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	if err := os.WriteFile(p.ready, []byte("started"), 0o600); err != nil {
		return nil, fmt.Errorf("record fixture startup: %w", err)
	}
	return &al_proto.PluginStartResponse{Result: &al_proto.ResourceResult{Env: p.exports}}, nil
}

// TestRunProcessFixture is a child executable and a plugin server for CLI scenarios.
func TestRunProcessFixture(t *testing.T) {
	switch os.Getenv("AL_RUN_PROCESS_FIXTURE") {
	case "child":
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(map[string]string{
			"input": string(input), "shared": os.Getenv("AL_RUN_SHARED"),
			"runfiles": os.Getenv("RUNFILES_DIR"),
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(os.Getenv("AL_RUN_CHILD_RECEIPT"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stdout.WriteString("out"); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stderr.WriteString("err"); err != nil {
			t.Fatal(err)
		}
		os.Exit(0)
	case "plugin":
		var exports map[string]string
		if err := json.Unmarshal([]byte(os.Getenv("AL_RUN_PLUGIN_EXPORTS")), &exports); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
		defer cancel()
		cmdCtx := al.NewCmdCtx(ctx, "")
		defer cmdCtx.RequestShutdown()
		cmdCtx.Logger = log.New(io.Discard, "", 0)
		server := al_plugin.NewPluginServer(cmdCtx, &runFixturePlugin{
			exports: exports, ready: os.Getenv("AL_RUN_PLUGIN_READY"),
		})
		if err := server.Start(ctx); err != nil {
			t.Fatal(err)
		}
		<-cmdCtx.Ctx.Done()
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		if err := server.Stop(stopCtx); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(os.Getenv("AL_RUN_PLUGIN_STOPPED"), []byte("cleaned"), 0o600); err != nil {
			t.Fatal(err)
		}
		os.Exit(0)
	}
}

// runCLIFixture holds a copied installation and synthetic child process inputs.
type runCLIFixture struct {
	// directory contains all mutable files for one isolated invocation.
	directory string
	// binary is the installed AL copy without adjacent Bazel runfiles by default.
	binary string
	// helper is this test's executable used as a child and plugin server.
	helper string
	// childReceipt records whether the intended child started and what it observed.
	childReceipt string
	// environment excludes ambient credentials and Bazel discovery variables.
	environment []string
}

// newRunCLIFixture creates a binary-only AL installation inside test scratch.
func newRunCLIFixture(t *testing.T) *runCLIFixture {
	t.Helper()
	source, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	binary := filepath.Join(directory, "al")
	if err := os.WriteFile(binary, data, 0o755); err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	receipt := filepath.Join(directory, "child.json")
	return &runCLIFixture{
		directory: directory, binary: binary, helper: helper, childReceipt: receipt,
		environment: []string{
			"PATH=" + os.Getenv("PATH"), "HOME=" + directory, "TMPDIR=" + directory,
			"AL_RUN_PROCESS_FIXTURE=child", "AL_RUN_CHILD_RECEIPT=" + receipt,
			"AL_RUN_SHARED=inherited-synthetic-value",
		},
	}
}

// invoke runs the installed CLI and retains command output even for a failing scenario.
func (f *runCLIFixture) invoke(t *testing.T, flags ...string) (string, string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	args := append([]string{"run"}, flags...)
	args = append(args, "--", f.helper, "-test.run=^TestRunProcessFixture$")
	cmd := exec.CommandContext(ctx, f.binary, args...)
	cmd.Dir = f.directory
	cmd.Env = f.environment
	cmd.Stdin = strings.NewReader("synthetic input")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	observation := map[string]any{
		"arguments": args, "stdout": stdout.String(), "stderr": stderr.String(),
		"exit_code": cmd.ProcessState.ExitCode(),
	}
	t.Cleanup(func() {
		files := map[string]any{}
		entries, readErr := os.ReadDir(f.directory)
		if readErr != nil {
			t.Error(readErr)
			return
		}
		for _, entry := range entries {
			if entry.Name() == "al" || !entry.Type().IsRegular() {
				continue
			}
			path := filepath.Join(f.directory, entry.Name())
			data, readErr := os.ReadFile(path)
			info, statErr := os.Stat(path)
			if readErr != nil || statErr != nil {
				t.Error(errors.Join(readErr, statErr))
				continue
			}
			files[entry.Name()] = map[string]any{"contents": string(data), "mode": info.Mode().Perm()}
		}
		observation["files"] = files
		if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
			data, marshalErr := json.MarshalIndent(observation, "", "  ")
			if marshalErr != nil {
				t.Error(marshalErr)
				return
			}
			name := strings.ReplaceAll(t.Name(), "/", "_") + ".json"
			if writeErr := os.WriteFile(filepath.Join(output, name), append(data, '\n'), 0o600); writeErr != nil {
				t.Error(writeErr)
			}
		}
	})
	return stdout.String(), stderr.String(), err
}

// config writes only synthetic plugin declarations for a real CLI invocation.
func (f *runCLIFixture) config(t *testing.T, plugins ...map[string]*al_proto.PluginConfig) string {
	t.Helper()
	declarations := map[string]*al_proto.PluginConfig{}
	for _, group := range plugins {
		for name, plugin := range group {
			declarations[name] = plugin
		}
	}
	data, err := protojson.Marshal(&al_proto.Config{PluginConfigs: declarations})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.directory, "config.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// plugin defines a synthetic RPC helper and independent startup/shutdown receipts.
func (f *runCLIFixture) plugin(t *testing.T, name, executable string, exports map[string]string) map[string]*al_proto.PluginConfig {
	t.Helper()
	data, err := json.Marshal(exports)
	if err != nil {
		t.Fatal(err)
	}
	return map[string]*al_proto.PluginConfig{name: {
		Bin: executable, Args: []string{"-test.run=^TestRunProcessFixture$"},
		Labels: map[string]string{"fixture": "run"},
		Env: map[string]string{
			"AL_RUN_PROCESS_FIXTURE": "plugin", "AL_RUN_PLUGIN_EXPORTS": string(data),
			"AL_RUN_PLUGIN_READY":   filepath.Join(f.directory, name+".ready"),
			"AL_RUN_PLUGIN_STOPPED": filepath.Join(f.directory, name+".stopped"),
		},
	}}
}

// TestRunStandaloneCLI distinguishes optional runfiles from required or invalid contexts.
func TestRunStandaloneCLI(t *testing.T) {
	for _, mode := range []string{"plain", "unselected", "absolute-plugin", "relative-plugin", "relative-plugin-with-runfiles", "adjacent-runfiles", "invalid-manifest", "invalid-directory"} {
		t.Run(mode, func(t *testing.T) {
			fixture := newRunCLIFixture(t)
			var flags []string
			wantFailure := false
			wantRunfiles := ""
			switch mode {
			case "unselected":
				flags = []string{"--config", fixture.config(t, fixture.plugin(t, "ignored", "fixture/missing", nil))}
			case "absolute-plugin", "relative-plugin", "relative-plugin-with-runfiles":
				executable := fixture.helper
				if mode != "absolute-plugin" {
					executable = "fixture/plugin"
				}
				if mode == "relative-plugin-with-runfiles" {
					wantRunfiles = filepath.Join(fixture.directory, "runfiles")
					if err := os.MkdirAll(filepath.Join(wantRunfiles, "fixture"), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.Symlink(fixture.helper, filepath.Join(wantRunfiles, executable)); err != nil {
						t.Fatal(err)
					}
					fixture.environment = append(fixture.environment, "RUNFILES_DIR="+wantRunfiles)
				}
				wantFailure = mode == "relative-plugin"
				flags = []string{"--config", fixture.config(t, fixture.plugin(t, "selected", executable, nil)), "--plugin_label", "fixture=run"}
			case "adjacent-runfiles":
				wantRunfiles = fixture.binary + ".runfiles"
				if err := os.Mkdir(wantRunfiles, 0o755); err != nil {
					t.Fatal(err)
				}
			case "invalid-manifest":
				fixture.environment = append(fixture.environment, "RUNFILES_MANIFEST_FILE="+filepath.Join(fixture.directory, "missing-manifest"))
				wantFailure = true
			case "invalid-directory":
				fixture.environment = append(fixture.environment, "RUNFILES_DIR="+filepath.Join(fixture.directory, "missing-directory"))
				wantFailure = true
			}
			stdout, stderr, err := fixture.invoke(t, flags...)
			if wantFailure {
				if err == nil {
					t.Fatal("invalid or unavailable required runfiles allowed child execution")
				}
				if _, statErr := os.Stat(fixture.childReceipt); !errors.Is(statErr, os.ErrNotExist) {
					t.Fatalf("child ran after invalid setup: %v", statErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("installed run failed: %v\n%s", err, stderr)
			}
			if stdout != "out" || !strings.Contains(stderr, "err") {
				t.Fatalf("child streams were not preserved: %q %q", stdout, stderr)
			}
			data, err := os.ReadFile(fixture.childReceipt)
			if err != nil {
				t.Fatal(err)
			}
			var observed map[string]string
			if err := json.Unmarshal(data, &observed); err != nil {
				t.Fatal(err)
			}
			if observed["input"] != "synthetic input" || observed["runfiles"] != wantRunfiles || observed["shared"] != "inherited-synthetic-value" {
				t.Fatalf("child context changed: %s", data)
			}
		})
	}
}

// TestRunStreamFiles checks creation, truncation, private new files, and input preservation.
func TestRunStreamFiles(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(fmt.Sprintf("existing-%t", existing), func(t *testing.T) {
			fixture := newRunCLIFixture(t)
			fixture.environment = append(fixture.environment, "RUNFILES_DIR="+fixture.directory)
			stdoutPath, stderrPath := filepath.Join(fixture.directory, "stdout"), filepath.Join(fixture.directory, "stderr")
			for _, path := range []string{stdoutPath, stderrPath} {
				if existing {
					if err := os.WriteFile(path, []byte("previous contents must disappear"), 0o644); err != nil {
						t.Fatal(err)
					}
				}
			}
			inputPath := filepath.Join(fixture.directory, "stdin")
			if err := os.WriteFile(inputPath, []byte("explicit input"), 0o600); err != nil {
				t.Fatal(err)
			}
			stdout, stderr, err := fixture.invoke(t, "--stdout", stdoutPath, "--stderr", stderrPath, "--stdin", inputPath)
			if err != nil {
				t.Fatalf("stream override failed: %v\n%s", err, stderr)
			}
			if stdout != "" || stderr != "" {
				t.Fatalf("child output escaped selected files: %q %q", stdout, stderr)
			}
			for path, want := range map[string]string{stdoutPath: "out", stderrPath: "err", inputPath: "explicit input"} {
				data, err := os.ReadFile(path)
				if err != nil || string(data) != want {
					t.Fatalf("stream file %s: %q, %v", filepath.Base(path), data, err)
				}
				info, err := os.Stat(path)
				if err != nil {
					t.Fatal(err)
				}
				if !existing && info.Mode().Perm() != 0o600 {
					t.Fatalf("new stream file is not private: %s %v", path, info.Mode())
				}
			}
		})
	}
}

// TestRunPluginEnvironment checks selected exports before starting the real child.
func TestRunPluginEnvironment(t *testing.T) {
	for _, mode := range []string{"equal", "conflict", "conflict-reversed"} {
		t.Run(mode, func(t *testing.T) {
			fixture := newRunCLIFixture(t)
			fixture.environment = append(fixture.environment, "RUNFILES_DIR="+fixture.directory)
			firstValue, secondValue := "synthetic-private-alpha", "synthetic-private-alpha"
			if mode != "equal" {
				secondValue = "synthetic-private-zeta"
			}
			first := fixture.plugin(t, "alpha", fixture.helper, map[string]string{"AL_RUN_SHARED": firstValue})
			second := fixture.plugin(t, "zeta", fixture.helper, map[string]string{"AL_RUN_SHARED": secondValue})
			plugins := []map[string]*al_proto.PluginConfig{second, first}
			if mode == "conflict-reversed" {
				plugins = []map[string]*al_proto.PluginConfig{first, second}
			}
			ignored := fixture.plugin(t, "ignored", "fixture/missing", map[string]string{"AL_RUN_SHARED": "ignored"})
			ignored["ignored"].Labels = map[string]string{"fixture": "unselected"}
			config := fixture.config(t, append(plugins, ignored)...)
			stdout, stderr, err := fixture.invoke(t, "--config", config, "--plugin_label", "fixture=run")
			if mode == "equal" {
				if err != nil {
					t.Fatalf("equal exports were rejected: %v\n%s", err, stderr)
				}
				data, readErr := os.ReadFile(fixture.childReceipt)
				if readErr != nil {
					t.Fatal(readErr)
				}
				var observed map[string]string
				if err := json.Unmarshal(data, &observed); err != nil {
					t.Fatal(err)
				}
				if observed["shared"] != firstValue {
					t.Fatalf("selected export did not override inherited value: %q", observed["shared"])
				}
			} else {
				if err == nil {
					t.Error("conflicting exports allowed child execution")
				}
				want := `conflicting environment variable "AL_RUN_SHARED" from plugins "alpha" and "zeta"`
				if !strings.Contains(stderr, want) {
					t.Errorf("missing deterministic conflict diagnostic: %s", stderr)
				}
				if _, statErr := os.Stat(fixture.childReceipt); !errors.Is(statErr, os.ErrNotExist) {
					t.Errorf("conflicting exports reached child: %v", statErr)
				}
			}
			for _, value := range []string{firstValue, secondValue} {
				if strings.Contains(stdout+stderr, value) {
					t.Error("plugin export value appeared in diagnostics")
				}
			}
			for _, name := range []string{"alpha", "zeta"} {
				for suffix, want := range map[string]string{".ready": "started", ".stopped": "cleaned"} {
					data, readErr := os.ReadFile(filepath.Join(fixture.directory, name+suffix))
					if readErr != nil || string(data) != want {
						t.Errorf("plugin %s%s lifecycle incomplete: %q %v", name, suffix, data, readErr)
					}
				}
			}
		})
	}
}
