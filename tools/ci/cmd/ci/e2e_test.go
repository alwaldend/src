package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// ciCall records one actual Bazel process launched by the compiled CI command.
type ciCall struct {
	// Directory identifies the workspace relative to the fixture repository.
	Directory string `json:"directory"`
	// Arguments contains the literal Bazel command arguments.
	Arguments []string `json:"arguments"`
}

// ciObservation retains the subprocess behavior of one isolated scenario.
type ciObservation struct {
	// Scenario identifies the failure or successful execution being exercised.
	Scenario string `json:"scenario"`
	// Calls contains all observed build and test processes in execution order.
	Calls []ciCall `json:"calls"`
	// ExitCode records the actual command's success or failure.
	ExitCode int `json:"exit_code"`
}

// TestMain substitutes only the external Bazel process in executable tests.
func TestMain(m *testing.M) {
	if filepath.Base(os.Args[0]) == "bazel" {
		os.Exit(recordBazelCall())
	}
	os.Exit(m.Run())
}

// recordBazelCall logs a private fixture invocation and injects selected failures.
func recordBazelCall() int {
	directory, err := os.Getwd()
	if err != nil {
		return 90
	}
	relative, err := filepath.Rel(os.Getenv("CI_FIXTURE_ROOT"), directory)
	if err != nil {
		return 91
	}
	file, err := os.OpenFile(os.Getenv("CI_FIXTURE_LOG"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return 92
	}
	writeErr := json.NewEncoder(file).Encode(ciCall{filepath.ToSlash(relative), os.Args[1:]})
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil || len(os.Args) < 2 {
		return 93
	}
	scenario := os.Getenv("CI_FIXTURE_SCENARIO")
	if scenario == "build_failure" && relative == "." && os.Args[1] == "build" {
		return 7
	}
	if scenario == "test_failure" && filepath.ToSlash(relative) == "nested/one" && os.Args[1] == "test" {
		return 8
	}
	return 0
}

// TestCIExecutable verifies real command launching, complete coverage and failure propagation.
func TestCIExecutable(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("CI_BINARY"))
	if err != nil {
		t.Fatalf("resolve CI executable: %v", err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatalf("resolve Bazel fixture executable: %v", err)
	}
	var observations []ciObservation
	for _, scenario := range []string{"success", "build_failure", "test_failure", "missing_bazel", "invalid_workspace"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			for path, content := range map[string]string{
				"MODULE.bazel":            "",
				".bazelignore":            "cache\nnested/one\nnested/two\nnested/one\n",
				"nested/one/MODULE.bazel": "",
				"nested/two/MODULE.bazel": "",
			} {
				fullPath := filepath.Join(root, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			bin := filepath.Join(root, "bin")
			if err := os.Mkdir(bin, 0o700); err != nil {
				t.Fatal(err)
			}
			if scenario != "missing_bazel" {
				if err := os.Symlink(helper, filepath.Join(bin, "bazel")); err != nil {
					t.Fatal(err)
				}
			}
			if scenario == "invalid_workspace" {
				if err := os.WriteFile(filepath.Join(root, ".bazelignore"), []byte("../outside\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			log := filepath.Join(root, "calls.jsonl")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, binary)
			cmd.Env = []string{
				"BUILD_WORKSPACE_DIRECTORY=" + root,
				"PATH=" + bin,
				"CI_FIXTURE_ROOT=" + root,
				"CI_FIXTURE_LOG=" + log,
				"CI_FIXTURE_SCENARIO=" + scenario,
			}
			output, runErr := cmd.CombinedOutput()
			if (runErr == nil) != (scenario == "success") {
				t.Fatalf("CI execution result %v: %s", runErr, output)
			}
			var calls []ciCall
			if file, err := os.Open(log); err == nil {
				decoder := json.NewDecoder(file)
				for {
					var call ciCall
					if err := decoder.Decode(&call); err == io.EOF {
						break
					} else if err != nil {
						t.Fatal(err)
					}
					calls = append(calls, call)
				}
				if err := file.Close(); err != nil {
					t.Fatal(err)
				}
			} else if !os.IsNotExist(err) {
				t.Fatal(err)
			}
			var expected []ciCall
			if scenario != "missing_bazel" && scenario != "invalid_workspace" {
				for _, directory := range []string{".", "nested/one", "nested/two"} {
					for _, phase := range []string{"build", "test"} {
						expected = append(expected, ciCall{directory, []string{phase, "--config=ci", "//..."}})
					}
				}
			}
			if !reflect.DeepEqual(calls, expected) {
				t.Fatalf("calls = %#v, want %#v; output: %s", calls, expected, output)
			}
			observations = append(observations, ciObservation{scenario, calls, cmd.ProcessState.ExitCode()})
		})
	}
	data, err := json.MarshalIndent(observations, "", "  ")
	if err != nil {
		t.Fatalf("encode command observations: %v", err)
	}
	if err := os.WriteFile(filepath.Join(os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), "ci-execution.json"), data, 0o600); err != nil {
		t.Fatalf("write command observations: %v", err)
	}
}
