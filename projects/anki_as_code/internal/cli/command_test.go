package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"google.golang.org/protobuf/proto"
)

// recordingApplication substitutes command workflows without starting a backend.
type recordingApplication struct {
	// calls counts application invocations for syntax and batch assertions.
	calls int
	// ctx records the context forwarded by Cobra.
	ctx context.Context
	// operation identifies the dispatched workflow.
	operation string
	// collection records the source path forwarded by the command.
	collection string
	// paths records configuration and output arguments.
	paths []string
	// key records the namespace passed explicitly to the selected workflow.
	key string
	// bindings records opaque path-to-ID input without native interpretation.
	bindings map[string]string
	// state records decoded update or replacement input.
	state *collectionpb.State
	// configs records the decoded configuration set supplied for file export.
	configs *collectionpb.CollectionConfigSet
	// failure supplies a workflow error whose identity must survive command wrapping.
	failure error
}

// record captures a workflow invocation for command-boundary assertions.
func (application *recordingApplication) record(ctx context.Context, operation string, paths ...string) {
	application.calls++
	application.ctx, application.operation, application.paths = ctx, operation, paths
}

// ColExport records collection-to-directory dispatch and optional state output.
func (application *recordingApplication) ColExport(ctx context.Context, input string, output, stateOutput, key string) error {
	application.key = key
	application.collection = input
	application.record(ctx, "export", output)
	return application.failure
}

// ColPlan supplies a typed empty plan without native access.
func (application *recordingApplication) ColPlan(ctx context.Context, input string, text, key string) (*collectionpb.CollectionChangePlan, error) {
	application.key = key
	application.collection = input
	application.record(ctx, "plan", text)
	return &collectionpb.CollectionChangePlan{}, application.failure
}

// ColApply records in-place reconciliation dispatch.
func (application *recordingApplication) ColApply(ctx context.Context, input string, text, key string) error {
	application.key = key
	application.collection = input
	application.record(ctx, "apply", text)
	return application.failure
}

// ColBuild records archive creation dispatch.
func (application *recordingApplication) ColBuild(ctx context.Context, base string, text, output, key string) error {
	application.key = key
	application.collection = base
	application.record(ctx, "build", text, output)
	return application.failure
}

// ConfigDump returns recognizable output for stream assertions.
func (application *recordingApplication) ConfigDump(ctx context.Context, text, format string) ([]byte, error) {
	application.record(ctx, "config dump", text, format)
	return []byte("config output\n"), application.failure
}

// ConfigExport records a decoded configuration set without writing files.
func (application *recordingApplication) ConfigExport(ctx context.Context, output string, configs *collectionpb.CollectionConfigSet) error {
	application.configs = configs
	application.record(ctx, "config export", output)
	return application.failure
}

// ColStateDump supplies stored state for command-side encoding.
func (application *recordingApplication) ColStateDump(ctx context.Context, input string, key string) (*collectionpb.State, error) {
	application.key = key
	application.collection = input
	application.record(ctx, "state dump")
	return &collectionpb.State{Version: 1}, application.failure
}

// ColStateImport records the whole binding batch and config path without reading files.
func (application *recordingApplication) ColStateImport(ctx context.Context, input string, text, key string, bindings map[string]string) error {
	application.key = key
	application.collection = input
	application.bindings = bindings
	application.record(ctx, "state import", text)
	return application.failure
}

// ColStateUpdate records the decoded external state.
func (application *recordingApplication) ColStateUpdate(ctx context.Context, input string, key string, value *collectionpb.State) error {
	application.key = key
	application.collection = input
	application.state = value
	application.record(ctx, "state update")
	return application.failure
}

// ColStateReplace records the decoded complete replacement.
func (application *recordingApplication) ColStateReplace(ctx context.Context, input string, key string, value *collectionpb.State) error {
	application.key = key
	application.collection = input
	application.state = value
	application.record(ctx, "state replace")
	return application.failure
}

// TestCommandDispatch runs the complete command tree against an isolated application.
func TestCommandDispatch(t *testing.T) {
	for _, scenario := range []struct {
		// operation names the expected application method.
		operation string
		// args contains command arguments after the shared namespace flag.
		args []string
		// paths contains arguments expected at the application boundary.
		paths []string
		// stdout indicates a data-producing command instead of a mutation status.
		stdout bool
	}{
		{"export", []string{"col", "export", "--collection", "source.colpkg", "--output", "destination"}, []string{"source.colpkg", "destination"}, false},
		{"plan", []string{"col", "plan", "--collection", "source.anki2", "--text", "configs"}, []string{"source.anki2", "configs"}, true},
		{"apply", []string{"col", "apply", "--collection", "source.anki2", "--text", "configs"}, []string{"source.anki2", "configs"}, false},
		{"build", []string{"col", "build", "--base", "source.colpkg", "--text", "configs", "--output", "result.colpkg"}, []string{"source.colpkg", "configs", "result.colpkg"}, false},
		{"config dump", []string{"config", "dump", "--text", "configs"}, []string{"configs", "json"}, true},
		{"config dump", []string{"conf", "dump", "--text", "configs"}, []string{"configs", "json"}, true},
		{"config export", []string{"config", "export", "--output", "configs"}, []string{"configs"}, false},
		{"config export", []string{"conf", "export", "--output", "configs", "-"}, []string{"configs"}, false},
		{"export", []string{"collection", "export", "--collection", "source.colpkg", "--output", "destination"}, []string{"source.colpkg", "destination"}, false},
		{"state dump", []string{"col", "state", "dump", "--collection", "source.colpkg"}, []string{"source.colpkg"}, true},
		{"state import", []string{"col", "state", "import", "--collection", "source.colpkg", "--text", "unread-configs", "notes/a=b.note.anki.toml=opaque-id", "decks/New/deck.anki.toml="}, []string{"source.colpkg", "unread-configs"}, false},
		{"state update", []string{"col", "state", "update", "--collection", "source.colpkg"}, []string{"source.colpkg"}, false},
		{"state replace", []string{"col", "state", "replace", "--collection", "source.colpkg", "-"}, []string{"source.colpkg"}, false},
	} {
		t.Run(scenario.operation, func(t *testing.T) {
			for _, fail := range []bool{false, true} {
				fake := &recordingApplication{}
				failure := errors.New("application unavailable")
				if fail {
					fake.failure = failure
				}
				command := New(fake)
				arguments := scenario.args
				if arguments[0] == "col" || arguments[0] == "collection" {
					arguments = append([]string{arguments[0], "--state-key", "isolated"}, arguments[1:]...)
				}
				command.SetArgs(arguments)
				command.SetIn(strings.NewReader("version = '1'\n"))
				if scenario.operation == "config export" {
					command.SetIn(strings.NewReader(`{"config":{"version":"1"}}`))
				}
				var stdout, stderr bytes.Buffer
				command.SetOut(&stdout)
				command.SetErr(&stderr)
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				err := command.ExecuteContext(ctx)
				if fail && !errors.Is(err, failure) || !fail && err != nil {
					t.Fatalf("application failure=%v: %v", fail, err)
				}
				expectedPaths := scenario.paths
				if !strings.HasPrefix(scenario.operation, "config ") {
					if fake.collection != scenario.paths[0] {
						t.Fatalf("incorrect collection selector: %+v", fake.collection)
					}
					expectedPaths = scenario.paths[1:]
					if len(expectedPaths) == 0 {
						expectedPaths = nil
					}
				}
				if fake.calls != 1 || fake.operation != scenario.operation || fake.ctx != ctx || !reflect.DeepEqual(fake.paths, expectedPaths) || !strings.HasPrefix(scenario.operation, "config ") && fake.key != "isolated" {
					t.Fatalf("incorrect application dispatch: %+v", fake)
				}
				if fail {
					if stdout.Len() != 0 || stderr.Len() != 0 {
						t.Fatalf("failed operation emitted success output: %q %q", stdout.String(), stderr.String())
					}
					continue
				}
				if scenario.stdout != (stdout.Len() != 0) || scenario.stdout == (stderr.Len() != 0) {
					t.Fatalf("incorrect command streams: %q %q", stdout.String(), stderr.String())
				}
				if scenario.operation == "config export" && !proto.Equal(fake.configs, &collectionpb.CollectionConfigSet{Config: &collectionpb.CollectionConfig{Version: 1}}) {
					t.Fatalf("config input lost in command decoding: %v", fake.configs)
				}
				if scenario.operation == "state import" && !reflect.DeepEqual(fake.bindings, map[string]string{"notes/a=b.note.anki.toml": "opaque-id", "decks/New/deck.anki.toml": ""}) {
					t.Fatalf("import altered the binding batch: %v", fake.bindings)
				}
				if scenario.operation == "state update" || scenario.operation == "state replace" {
					if !proto.Equal(fake.state, &collectionpb.State{Version: 1}) {
						t.Fatalf("state input lost in command decoding: %v", fake.state)
					}
				}
				if scenario.operation == "state dump" {
					value := &collectionpb.State{}
					if err := codec.Decode(stdout.Bytes(), codec.TOML, value); err != nil || value.Version != 1 {
						t.Fatalf("state output is not TOML: %q %v", stdout.String(), err)
					}
				}
			}
		})
	}
}

// TestCommandRejections prevents malformed commands from reaching workflows.
func TestCommandRejections(t *testing.T) {
	for _, args := range [][]string{
		{"col"},
		{"collection"},
		{"export", "--collection", "unused", "--output", "unused"},
		{"state", "dump", "--collection", "unused"},
		{"config"},
		{"config", "unknown"},
		{"config", "dump"},
		{"config", "export"},
		{"config", "export", "--output", "unused"},
		{"config", "export", "--output", "unused", "missing.json"},
		{"config", "export", "--output", "unused", "--format", "invalid"},
		{"config", "export", "--output", "unused", "one", "two"},
		{"col", "export", "--collection", "unused"},
		{"col", "state", "dump"},
		{"col", "state", "import", "--collection", "unused", "--text", "unused", "bad"},
		{"col", "state", "import", "--collection", "unused", "--text", "unused", "a=1", "a=2"},
		{"col", "state", "update", "--collection", "unused"},
		{"col", "state", "replace", "--collection", "unused", "--format", "invalid"},
	} {
		fake := &recordingApplication{}
		command := New(fake)
		command.SetArgs(args)
		command.SetIn(strings.NewReader("invalid state input"))
		var output bytes.Buffer
		command.SetOut(&output)
		command.SetErr(&output)
		if err := command.ExecuteContext(t.Context()); err == nil || fake.calls != 0 || output.Len() != 0 {
			t.Fatalf("invalid command %v reached application or emitted data: calls=%d error=%v output=%q", args, fake.calls, err, output.String())
		}
	}
}

// TestStateCommandJSON verifies file decoding and JSON output remain command responsibilities.
func TestStateCommandJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte(`{"version":"1"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	fake := &recordingApplication{}
	for _, args := range [][]string{
		{"col", "state", "replace", "--collection", "unused", "--format", "json", "--state-key", "custom", path},
		{"col", "state", "dump", "--collection", "unused", "--format", "json"},
	} {
		command := New(fake)
		command.SetArgs(args)
		var stdout, stderr bytes.Buffer
		command.SetOut(&stdout)
		command.SetErr(&stderr)
		if err := command.ExecuteContext(t.Context()); err != nil {
			t.Fatal(err)
		}
		expectedKey := "custom"
		if args[2] == "dump" {
			expectedKey = state.DefaultKey
		}
		if fake.key != expectedKey {
			t.Fatalf("command namespace: got %q, want %q", fake.key, expectedKey)
		}
		if args[2] == "dump" {
			value := &collectionpb.State{}
			if err := codec.Decode(stdout.Bytes(), codec.JSON, value); err != nil || !proto.Equal(value, fake.state) {
				t.Fatalf("JSON state round trip: %q %v", stdout.String(), err)
			}
			if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
				if err := os.WriteFile(filepath.Join(destination, "isolated-state.json"), stdout.Bytes(), 0o600); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}

// failingWriter makes output failure identity observable at the command boundary.
type failingWriter struct {
	// failure is returned for every attempted write.
	failure error
}

// Write rejects output with the configured cause.
func (writer failingWriter) Write(data []byte) (int, error) {
	return 0, writer.failure
}

// TestCommandOutputFailure preserves write errors after successful workflow dispatch.
func TestCommandOutputFailure(t *testing.T) {
	failure := errors.New("output unavailable")
	for _, args := range [][]string{
		{"col", "export", "--collection", "unused", "--output", "unused"},
		{"col", "apply", "--collection", "unused", "--text", "unused"},
		{"col", "build", "--base", "unused", "--text", "unused", "--output", "unused"},
		{"config", "dump", "--text", "unused"},
		{"col", "state", "dump", "--collection", "unused"},
		{"col", "plan", "--collection", "unused", "--text", "unused"},
		{"col", "state", "update", "--collection", "unused"},
	} {
		fake := &recordingApplication{}
		command := New(fake)
		command.SetArgs(args)
		command.SetIn(strings.NewReader("version = '1'\n"))
		command.SetOut(failingWriter{failure})
		command.SetErr(failingWriter{failure})
		if err := command.ExecuteContext(t.Context()); !errors.Is(err, failure) || fake.calls != 1 {
			t.Fatalf("command lost output error: %v, calls=%d", err, fake.calls)
		}
	}
}
