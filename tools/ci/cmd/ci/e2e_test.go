package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// ciCall records one actual Bazel process launched by the compiled CI command.
type ciCall struct {
	// Directory identifies the working directory relative to the fixture root.
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
	if os.Args[1] == "shutdown" {
		if scenario == "shutdown_failure" {
			return 9
		}
		return 0
	}
	fmt.Fprintln(os.Stdout, "Starting local Bazel server...")
	fmt.Fprintln(os.Stderr, "(08:38:50) \x1b[32mINFO: \x1b[0mInvocation ID: fixture")
	fmt.Fprintln(os.Stderr, "(08:38:50) WARNING: option '--terminal_columns' was expanded twice")
	fmt.Fprintln(os.Stderr, "(13:08:42) INFO: From Generating JSON:")
	fmt.Fprintln(os.Stdout, "bazel-out/workflows.json")
	if scenario == "build_failure" && os.Args[1] == "build" && len(os.Args) == 4 {
		fmt.Fprintln(os.Stderr, "(13:08:43) \x1b[31mERROR: \x1b[0mfixture build failed")
		fmt.Fprintln(os.Stderr, "fixture compiler details")
		return 7
	}
	if scenario == "lint_failure" && len(os.Args) == 5 && os.Args[3] == "--config=lint" {
		fmt.Fprintln(os.Stderr, "ERROR: fixture lint failed")
		return 10
	}
	if scenario == "test_failure" && os.Args[1] == "test" {
		fmt.Fprint(os.Stdout, "\n==================== Test output for //fixture:failure:\nERROR: from the test\nfailed test diagnostics\n================================================================================\n")
		fmt.Fprintln(os.Stdout, "//fixture:failure FAILED in 0.1s")
		return 8
	}
	return 0
}

// TestCIExecutable verifies root-only command launching and failure propagation.
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
	for _, scenario := range []string{"success", "forgejo_success", "build_failure", "test_failure", "lint_failure", "shutdown_failure", "missing_bazel", "invalid_workspace", "unwritable_summary"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			if scenario != "invalid_workspace" {
				if err := os.WriteFile(filepath.Join(root, "MODULE.bazel"), nil, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if scenario == "forgejo_success" {
				if err := os.Mkdir(filepath.Join(root, "generated"), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "generated", "MODULE.bazel"), nil, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, ".bazelignore"), []byte("generated\n"), 0o600); err != nil {
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
			summary := filepath.Join(root, "summary.md")
			switch scenario {
			case "success":
				// A local invocation has neither the runner marker nor a summary file.
			case "test_failure":
				cmd.Env = append(cmd.Env, "FORGEJO_ACTIONS=true", "GITHUB_STEP_SUMMARY="+summary)
			case "unwritable_summary":
				cmd.Env = append(cmd.Env, "FORGEJO_ACTIONS=true", "FORGEJO_STEP_SUMMARY="+root)
			default:
				if err := os.WriteFile(summary, []byte("Existing content\n"), 0o600); err != nil {
					t.Fatal(err)
				}
				cmd.Env = append(cmd.Env, "FORGEJO_ACTIONS=true", "FORGEJO_STEP_SUMMARY="+summary)
			}
			output, runErr := cmd.CombinedOutput()
			success := scenario == "success" || scenario == "forgejo_success" || scenario == "unwritable_summary"
			if (runErr == nil) != success {
				t.Fatalf("CI execution result %v: %s", runErr, output)
			}
			verifyReporting(t, scenario, string(output), summary)
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
				expected = []ciCall{
					{".", []string{"build", "--config=ci", "//..."}},
					{".", []string{"test", "--config=ci", "//..."}},
					{".", []string{"build", "--config=ci", "--config=lint", "//..."}},
				}
				if scenario != "success" {
					expected = append(expected, ciCall{".", []string{"shutdown"}})
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

// verifyReporting checks the runner-facing log and summary at the executable boundary.
func verifyReporting(t *testing.T, scenario, output, summary string) {
	t.Helper()
	if scenario == "success" {
		if strings.Contains(output, "##[group]") || !strings.Contains(output, "Invocation ID: fixture") {
			t.Fatalf("local output should retain plain Bazel diagnostics: %s", output)
		}
		return
	}
	if scenario == "unwritable_summary" {
		if !strings.Contains(output, "could not write the job summary") {
			t.Fatalf("missing summary warning: %s", output)
		}
		return
	}
	data, err := os.ReadFile(summary)
	if err != nil {
		t.Fatal(err)
	}
	report := string(data)
	if scenario != "test_failure" && !strings.HasPrefix(report, "Existing content\n") {
		t.Fatalf("summary replaced existing content: %s", report)
	}
	if scenario == "invalid_workspace" {
		if strings.Contains(output, "##[group]") || !strings.Contains(report, "Repository workspace unavailable") {
			t.Fatalf("unexpected discovery failure report: %s\n%s", report, output)
		}
		return
	}
	if scenario == "shutdown_failure" && !strings.Contains(report, "Bazel server shutdown failed in Root.") {
		t.Fatalf("missing Bazel shutdown failure: %s", report)
	}
	buildStatus, testStatus, lintStatus := "Passed", "Passed", "Passed"
	if scenario == "missing_bazel" || scenario == "build_failure" {
		buildStatus = `\*\*Failed\*\*`
	}
	if scenario == "missing_bazel" || scenario == "test_failure" {
		testStatus = `\*\*Failed\*\*`
	}
	if scenario == "missing_bazel" || scenario == "lint_failure" {
		lintStatus = `\*\*Failed\*\*`
	}
	row := fmt.Sprintf(`(?m)^\| Root \| %s \([^)]+\) \| %s \([^)]+\) \| %s \([^)]+\) \|$`, buildStatus, testStatus, lintStatus)
	if !regexp.MustCompile(row).MatchString(report) {
		t.Fatalf("missing root phase results: %s", report)
	}
	for _, phase := range []string{"Build", "Test", "Lint"} {
		heading := "##[group]Root / " + phase + "\n"
		_, afterHeading, found := strings.Cut(output, heading)
		if scenario != "missing_bazel" {
			before, folded, opened := strings.Cut(afterHeading, "##[group]Bazel diagnostics\n")
			diagnostics, afterDiagnostics, closed := strings.Cut(folded, "\n##[endgroup]\n")
			if !opened || !closed || !strings.Contains(diagnostics, "Starting local Bazel server") ||
				!strings.Contains(diagnostics, "Invocation ID: fixture") ||
				!strings.Contains(diagnostics, "WARNING: option '--terminal_columns'") ||
				!strings.Contains(diagnostics, "INFO: From Generating JSON:\nbazel-out/workflows.json") ||
				strings.Contains(diagnostics, "fixture build failed") {
				t.Fatalf("unexpected folded diagnostics in %q: %s", heading, afterHeading)
			}
			afterHeading = before + afterDiagnostics
		}
		if scenario == "test_failure" && phase == "Test" {
			before, testLog, opened := strings.Cut(afterHeading, "##[group]Failed test //fixture:failure\n")
			diagnostics, afterTest, closed := strings.Cut(testLog, "\n##[endgroup]\n")
			if !opened || !closed || !strings.Contains(diagnostics, "ERROR: from the test\nfailed test diagnostics") ||
				!strings.HasPrefix(afterTest, "//fixture:failure FAILED in 0.1s\n") {
				t.Fatalf("missing collapsed failed-test diagnostics: %s", afterHeading)
			}
			afterHeading = before + diagnostics + afterTest
		}
		diagnostics, afterGroup, closed := strings.Cut(afterHeading, "\n##[endgroup]\n")
		if !found || !closed || strings.Contains(diagnostics, "##[group]") {
			t.Fatalf("missing or unbalanced group %q: %s", heading, output)
		}
		if scenario == "build_failure" && phase == "Build" &&
			!strings.Contains(diagnostics, "fixture build failed\nfixture compiler details") {
			t.Fatalf("missing visible build failure: %s", diagnostics)
		}
		if scenario == "lint_failure" && phase == "Lint" && !strings.Contains(diagnostics, "ERROR: fixture lint failed") {
			t.Fatalf("missing visible lint failure: %s", diagnostics)
		}
		output = afterGroup
	}
	if strings.Contains(output, "##[group]") || strings.Contains(output, "##[endgroup]") {
		t.Fatalf("unmatched group markers: %s", output)
	}
	if scenario == "build_failure" {
		artifact := filepath.Join(os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), "forgejo-job-summary.md")
		if err := os.WriteFile(artifact, data, 0o600); err != nil {
			t.Fatalf("write example job summary: %v", err)
		}
	}
}
