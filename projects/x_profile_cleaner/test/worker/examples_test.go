package worker_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/cleanup"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// exampleConfig identifies the checked-in case configuration supplied by Bazel.
var exampleConfig = flag.String("example-config", "", "example case configuration runfile")

// exampleFiles contains every maintained Lua example, independently of the cases.
var exampleFiles = flag.String("example-files", "", "space-separated example runfiles")

// exampleConfiguration owns inputs and expectations without Go-specific case data.
type exampleConfiguration struct {
	// AccountID selects the expected owner of every configured export.
	AccountID string `json:"account_id"`
	// EvaluatedAt fixes time predicates independently of the host clock.
	EvaluatedAt time.Time `json:"evaluated_at"`
	// Archives maps a source key to ZIP entry names and relative fixture runfiles.
	Archives map[string]map[string]string `json:"archives"`
	// Cases selects an example, input archive, entrypoint and complete expected plan.
	Cases []exampleCase `json:"cases"`
}

// exampleCase describes one independently evaluated maintained script.
type exampleCase struct {
	// Name supplies the test case and worker operation identifier.
	Name string `json:"name"`
	// Example is the basename of a Lua file in the examples target.
	Example string `json:"example"`
	// Archive selects a configured fixture export.
	Archive string `json:"archive"`
	// Entrypoint is either script or the CLI/browser cleanup DSL.
	Entrypoint string `json:"entrypoint"`
	// ExpectedPlan is strict protobuf JSON, including optional report presence.
	ExpectedPlan json.RawMessage `json:"expected_plan"`
}

// readExampleRunfile reads only Bazel-declared inputs, including manifest-only runs.
func readExampleRunfile(t *testing.T, name string) []byte {
	t.Helper()
	file, err := runfiles.Rlocation(name)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// TestConfiguredExamples executes every shipped example through the real worker.
// Configuration owns data, time and exact plans; unconfigured examples fail coverage.
func TestConfiguredExamples(t *testing.T) {
	decoder := json.NewDecoder(bytes.NewReader(readExampleRunfile(t, *exampleConfig)))
	decoder.DisallowUnknownFields()
	var config exampleConfiguration
	if err := decoder.Decode(&config); err != nil {
		t.Fatal(err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		t.Fatal("example configuration contains trailing content")
	}
	if config.AccountID == "" || config.EvaluatedAt.IsZero() || len(config.Cases) == 0 {
		t.Fatal("example configuration requires account, time and cases")
	}
	scripts := map[string]string{}
	for _, file := range strings.Fields(*exampleFiles) {
		name := path.Base(file)
		if scripts[name] != "" {
			t.Fatalf("duplicate example basename %s", name)
		}
		scripts[name] = file
	}
	covered, names := map[string]bool{}, map[string]bool{}
	for _, item := range config.Cases {
		if scripts[item.Example] == "" || item.Name == "" || names[item.Name] {
			t.Fatalf("invalid case %s: example %s must exist and names must be unique", item.Name, item.Example)
		}
		covered[item.Example], names[item.Name] = true, true
	}
	for name := range scripts {
		if !covered[name] {
			t.Errorf("example %s has no configured case", name)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
	socket, _ := startWorker(t, "5s", "5m")
	for name, files := range config.Archives {
		if !filepath.IsLocal(name) || filepath.Base(name) != name {
			t.Fatalf("invalid fixture archive name %s", name)
		}
		entries := map[string]string{}
		for entry, file := range files {
			if !filepath.IsLocal(file) {
				t.Fatalf("fixture path must be relative: %s", file)
			}
			entries[entry] = string(readExampleRunfile(t, path.Join(path.Dir(*exampleConfig), file)))
		}
		if err := os.WriteFile(filepath.Join(filepath.Dir(socket), name+".zip"), datasetZIP(t, entries), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := connection.Close(); err != nil {
			t.Error(err)
		}
	}()
	rpc := workerpb.NewArchiveWorkerClient(connection)
	for _, item := range config.Cases {
		t.Run(item.Name, func(t *testing.T) {
			if _, exists := config.Archives[item.Archive]; !exists {
				t.Fatalf("unknown fixture archive %s", item.Archive)
			}
			expected := &workerpb.Plan{}
			if err := protojson.Unmarshal(item.ExpectedPlan, expected); err != nil {
				t.Fatalf("invalid expected protobuf plan: %v", err)
			}
			source := string(readExampleRunfile(t, scripts[item.Example]))
			switch item.Entrypoint {
			case "script":
			case "cleanup":
				source = cleanup.Script(source)
			default:
				t.Fatalf("unknown example entrypoint %s", item.Entrypoint)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			ref := &workerpb.EvaluationRef{JobId: item.Name, AccountId: config.AccountID}
			request := &workerpb.Evaluation{OperationId: ref.JobId, AccountId: ref.AccountId, SourceName: "exports", ArchiveKey: item.Archive + ".zip", Script: source, EvaluatedAt: timestamppb.New(config.EvaluatedAt)}
			if _, err := rpc.SubmitEvaluation(ctx, request); err != nil {
				t.Fatal(err)
			}
			defer func() {
				release, stop := context.WithTimeout(context.Background(), 2*time.Second)
				defer stop()
				if _, err := rpc.DeleteEvaluation(release, ref); err != nil {
					t.Error(err)
				}
			}()
			result := waitStatus(t, ctx, rpc, ref)
			if result.State != workerpb.EvaluationState_SUCCEEDED || result.Plan == nil {
				t.Fatalf("example failed: %v", result.Failure)
			}
			if !proto.Equal(result.Plan, expected) {
				t.Errorf("plan mismatch\nactual: %s\nexpected: %s", protojson.Format(result.Plan), protojson.Format(expected))
			}
			data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(result.Plan)
			if err != nil {
				t.Fatal(err)
			}
			writeAsyncArtifact(t, "example-"+item.Name+".json", json.RawMessage(data))
		})
	}
}
