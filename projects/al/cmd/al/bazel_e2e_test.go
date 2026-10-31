package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// bazelTestObservation records the public selection and subprocess boundary of one run.
type bazelTestObservation struct {
	// Scenario names the independent behavior under test.
	Scenario string `json:"scenario"`
	// Report contains the selection emitted before test execution.
	Report json.RawMessage `json:"report,omitempty"`
	// Diagnostic preserves failed operation diagnostics.
	Diagnostic string `json:"diagnostic,omitempty"`
	// Calls contains the argument vectors observed by the Bazel fixture.
	Calls []string `json:"calls"`
	// ExitCode records the executable's result.
	ExitCode int `json:"exit_code"`
}

// TestAffectedBazelHelper substitutes only the Bazel process boundary; Git and AL are real.
// The parent validation separately checks these query expressions against a real Bazel graph.
func TestAffectedBazelHelper(t *testing.T) {
	log := os.Getenv("AL_BAZEL_CALLS")
	if log == "" {
		return
	}
	separator := slices.Index(os.Args, "--")
	if separator < 0 || separator+1 >= len(os.Args) {
		os.Exit(90)
	}
	args := os.Args[separator+1:]
	data, err := json.Marshal(args)
	if err != nil {
		os.Exit(91)
	}
	file, err := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		os.Exit(92)
	}
	_, writeErr := file.Write(append(data, '\n'))
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil {
		os.Exit(93)
	}
	scenario := os.Getenv("AL_BAZEL_SCENARIO")
	if scenario == "cancellation" {
		if err := os.WriteFile(log+".ready", []byte("ready"), 0o600); err != nil {
			os.Exit(94)
		}
		time.Sleep(20 * time.Second)
		os.Exit(0)
	}
	if args[0] == "test" {
		if scenario == "test_failure" {
			os.Exit(3)
		}
		os.Exit(0)
	}
	if args[0] != "query" {
		os.Exit(95)
	}
	if scenario == "query_failure" {
		fmt.Println("//service:behavior_test") // Partial output is never a valid selection.
		os.Exit(3)
	}
	if !slices.Contains(args, "--output=label") {
		os.Exit(96)
	}
	if strings.HasPrefix(args[1], "kind(\"source file\"") {
		if scenario == "external_scope" {
			fmt.Println("@consumer//lib:source.txt")
			os.Exit(0)
		}
		if scenario == "manual_suite_wildcard" {
			os.Exit(0)
		}
		fmt.Println("//lib:source.txt\n//lib:source with spaces.txt")
		os.Exit(0)
	}
	if scenario == "manual_suite_wildcard" {
		if !strings.Contains(args[1], `except attr("tags", "[\[ ]manual[,\]]"`) {
			fmt.Println("//outside:outside_test")
		}
		os.Exit(0)
	}
	if scenario == "no_tests" {
		os.Exit(0)
	}
	if scenario == "manual_suite_explicit" {
		fmt.Println("//outside:outside_test\n//outside:manual_member_test")
		os.Exit(0)
	}
	if scenario == "manual_explicit" {
		fmt.Println("//service:manual_test")
		os.Exit(0)
	}
	fmt.Println("//service:behavior_test")
	if scenario == "similar_manual_tag" {
		fmt.Println("//service:not_manual_test")
	}
	if !strings.Contains(strings.Join(args, " "), "rdeps(") {
		fmt.Println("//other:other_test")
	}
	os.Exit(0)
}

// TestAffectedBazelCLI checks changed-file collection, safe selection, and execution failures.
func TestAffectedBazelCLI(t *testing.T) {
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	var observations []bazelTestObservation
	for _, scenario := range []string{
		"source", "committed", "staged_reverted", "spaces", "dry_run", "unchanged", "no_tests",
		"deleted", "renamed", "build_file", "bzl_file", "untracked", "outside_workspace", "external_scope",
		"manual_explicit", "manual_suite_wildcard", "manual_suite_explicit", "similar_manual_tag",
		"query_failure", "test_failure", "invalid_base", "option_base", "cancellation",
	} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			directory := t.TempDir()
			repository := filepath.Join(directory, "repository")
			workspace := filepath.Join(repository, "workspace")
			write := func(name, body string) {
				t.Helper()
				path := filepath.Join(repository, name)
				if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			write("workspace/MODULE.bazel", "module(name = \"fixture\")\n")
			write("workspace/lib/BUILD.bazel", "exports_files([\"source.txt\", \"source with spaces.txt\"])\n")
			write("workspace/lib/source.txt", "original\n")
			write("workspace/lib/source with spaces.txt", "original\n")
			write("workspace/rules.bzl", "# original\n")
			write("configuration.txt", "original\n")
			runGit := func(args ...string) string {
				t.Helper()
				command := exec.CommandContext(ctx, git, append([]string{"-C", repository}, args...)...)
				command.Env = []string{"PATH=" + os.Getenv("PATH"), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}
				out, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("Git fixture %v: %v; %s", args, err, out)
				}
				return strings.TrimSpace(string(out))
			}
			runGit("init", "--quiet", "--initial-branch=main")
			runGit("config", "user.name", "AL Fixture")
			runGit("config", "user.email", "fixture@example.invalid")
			runGit("config", "commit.gpgsign", "false")
			runGit("add", ".")
			runGit("commit", "--quiet", "-m", "Initialize test selection fixture")
			base := runGit("rev-parse", "HEAD")
			args := []string{"bazel", "test", "--base", base, "--workspace", workspace}
			switch scenario {
			case "unchanged":
			case "deleted", "renamed":
				if err := os.Remove(filepath.Join(workspace, "lib/source.txt")); err != nil {
					t.Fatal(err)
				}
				if scenario == "renamed" {
					write("workspace/lib/renamed.txt", "original\n")
				}
			case "build_file":
				write("workspace/lib/BUILD.bazel", "# changed graph\n")
			case "bzl_file":
				write("workspace/rules.bzl", "# changed graph\n")
			case "untracked":
				write("workspace/new.txt", "untracked\n")
			case "outside_workspace":
				write("configuration.txt", "changed\n")
			case "spaces":
				write("workspace/lib/source with spaces.txt", "changed\n")
			default:
				write("workspace/lib/source.txt", "changed\n")
			}
			switch scenario {
			case "committed":
				runGit("add", ".")
				runGit("commit", "--quiet", "-m", "Change library source")
			case "staged_reverted":
				runGit("add", ".")
				write("workspace/lib/source.txt", "original\n")
			case "manual_explicit":
				args = append(args, "//service:manual_test")
			case "manual_suite_wildcard":
				args = append(args, "//suite:all")
			case "manual_suite_explicit":
				args = append(args, "//suite:manual_suite")
			case "external_scope":
				args = append(args, "@consumer//...")
			case "invalid_base":
				args[3] = "missing-ref"
			case "option_base":
				args[3] = "--help"
			case "dry_run":
				args = append(args, "--dry-run")
			}
			stub := filepath.Join(directory, "bazel")
			if err := os.WriteFile(stub, []byte("#!/bin/sh\nexec \"$AL_BAZEL_HELPER\" -test.run=^TestAffectedBazelHelper$ -- \"$@\"\n"), 0o700); err != nil {
				t.Fatal(err)
			}
			log := filepath.Join(directory, "calls.jsonl")
			command := exec.CommandContext(ctx, binary, args...)
			command.Dir = repository
			command.Env = append(os.Environ(), "PATH="+directory+":"+os.Getenv("PATH"), "BUILD_WORKSPACE_DIRECTORY=", "BUILD_WORKING_DIRECTORY=", "AL_BAZEL_HELPER="+helper, "AL_BAZEL_CALLS="+log, "AL_BAZEL_SCENARIO="+scenario, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null")
			var stdout, stderr bytes.Buffer
			command.Stdout, command.Stderr = &stdout, &stderr
			if err := command.Start(); err != nil {
				t.Fatal(err)
			}
			if scenario == "cancellation" {
				deadline := time.Now().Add(5 * time.Second)
				for {
					if _, err := os.Stat(log + ".ready"); err == nil {
						break
					}
					if time.Now().After(deadline) {
						t.Fatal("Bazel query did not start before cancellation")
					}
					time.Sleep(10 * time.Millisecond)
				}
				if err := command.Process.Signal(syscall.SIGTERM); err != nil {
					t.Fatal(err)
				}
			}
			runErr := command.Wait()
			logData, readErr := os.ReadFile(log)
			if readErr != nil && !os.IsNotExist(readErr) {
				t.Fatal(readErr)
			}
			calls := strings.Split(strings.TrimSpace(string(logData)), "\n")
			if len(logData) == 0 {
				calls = nil
			}
			observation := bazelTestObservation{Scenario: scenario, Diagnostic: stderr.String(), Calls: calls, ExitCode: command.ProcessState.ExitCode()}
			if stdout.Len() > 0 {
				if !json.Valid(stdout.Bytes()) {
					t.Fatalf("invalid selection report: %s; %s", &stdout, &stderr)
				}
				observation.Report = append([]byte(nil), stdout.Bytes()...)
			}
			observations = append(observations, observation)
			wantFailure := slices.Contains([]string{"query_failure", "test_failure", "invalid_base", "option_base", "cancellation"}, scenario)
			if (runErr != nil) != wantFailure {
				t.Fatalf("unexpected result: %v; stdout=%s; stderr=%s", runErr, &stdout, &stderr)
			}
			testCalled := len(calls) > 0 && strings.HasPrefix(calls[len(calls)-1], `["test",`)
			wantTest := !slices.Contains([]string{"dry_run", "unchanged", "no_tests", "manual_suite_wildcard", "query_failure", "invalid_base", "option_base", "cancellation"}, scenario)
			if testCalled != wantTest {
				t.Fatalf("Bazel test called=%v, expected=%v: %v", testCalled, wantTest, calls)
			}
			if wantFailure && scenario != "test_failure" {
				if stdout.Len() != 0 {
					t.Fatalf("failed selection emitted a success report: %s", &stdout)
				}
				return
			}
			var report struct {
				// Selection identifies affected, full-scope, or unchanged selection.
				Selection string `json:"selection"`
				// Tests contains the exact labels passed to Bazel test.
				Tests []string `json:"tests"`
				// ChangedPaths preserves Git path names without whitespace splitting.
				ChangedPaths []string `json:"changed_paths"`
			}
			if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			full := slices.Contains([]string{"deleted", "renamed", "build_file", "bzl_file", "untracked", "outside_workspace", "external_scope", "manual_suite_wildcard"}, scenario)
			wantSelection := "affected"
			if full {
				wantSelection = "full_scope"
			} else if scenario == "unchanged" {
				wantSelection = "no_changes"
			}
			if report.Selection != wantSelection {
				t.Fatalf("selection=%q, expected=%q: %s", report.Selection, wantSelection, &stdout)
			}
			if slices.Contains(report.Tests, "//service:manual_test") != (scenario == "manual_explicit") {
				t.Fatalf("manual test selection did not respect explicit opt-in: %s", &stdout)
			}
			if full && scenario != "manual_suite_wildcard" && !slices.Contains(report.Tests, "//other:other_test") {
				t.Fatalf("full scope missed independent test: %s", &stdout)
			}
			if scenario == "manual_suite_explicit" && (!slices.Contains(report.Tests, "//outside:outside_test") || !slices.Contains(report.Tests, "//outside:manual_member_test")) {
				t.Fatalf("explicit suite lost its members: %s", &stdout)
			}
			if scenario == "similar_manual_tag" && !slices.Contains(report.Tests, "//service:not_manual_test") {
				t.Fatalf("similar tag was treated as manual: %s", &stdout)
			}
			if scenario == "spaces" && !slices.Contains(report.ChangedPaths, "workspace/lib/source with spaces.txt") {
				t.Fatalf("spaced path was lost: %s", &stdout)
			}
			if scenario == "unchanged" && len(calls) != 0 {
				t.Fatalf("unchanged source invoked Bazel: %v", calls)
			}
		})
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(observations, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "affected-bazel-tests.json"), append(data, '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
