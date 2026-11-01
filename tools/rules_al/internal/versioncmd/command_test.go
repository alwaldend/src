package versioncmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

func TestNoArgumentsPrintsUsage(t *testing.T) {
	var output bytes.Buffer
	if err := execute(nil, time.Now, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "release-start") {
		t.Fatalf("usage = %q", output.String())
	}
}

func TestTopLevelHelpPrintsUsageAndSucceeds(t *testing.T) {
	for _, help := range []string{"--help", "-h"} {
		var output bytes.Buffer
		if err := execute([]string{help}, time.Now, &output); err != nil {
			t.Fatalf("execute(%q) error = %v", help, err)
		}
		if !strings.Contains(output.String(), "release-start") {
			t.Fatalf("execute(%q) usage = %q", help, output.String())
		}
		if !strings.Contains(output.String(), "--channel") ||
			!strings.Contains(output.String(), "--release") {
			t.Fatalf("execute(%q) omitted context options: %q", help, output.String())
		}
	}
}

func TestEverySubcommandHelpPrintsUsageAndSucceeds(t *testing.T) {
	commands := []string{
		"show",
		"nightly-tag",
		"release-start",
		"release-tag",
		"release-plan",
		"release-publish",
		"bazel-status",
		"bazel",
	}
	for _, command := range commands {
		for _, help := range []string{"--help", "-h"} {
			var output bytes.Buffer
			if err := execute([]string{command, help}, time.Now, &output); err != nil {
				t.Fatalf("execute(%q, %q) error = %v", command, help, err)
			}
			if !strings.Contains(output.String(), command) {
				t.Fatalf("execute(%q, %q) usage = %q", command, help, output.String())
			}
			if !strings.Contains(output.String(), "--channel") ||
				!strings.Contains(output.String(), "--release") {
				t.Fatalf("execute(%q, %q) omitted context options: %q", command, help, output.String())
			}
		}
	}
}

func TestSubcommandHelpAfterOptionsPrintsUsageAndSucceeds(t *testing.T) {
	tests := map[string][]string{
		"show":            {"--format=json"},
		"nightly-tag":     {"--date=2026-08-30", "--dry-run"},
		"release-start":   {"--date=2026-08-30", "--switch=false"},
		"release-tag":     {"--dry-run"},
		"release-plan":    {"--plan=out/plan.json"},
		"release-publish": {"--plan=out/plan.json", "--receipt=out/receipt.json"},
	}
	for command, options := range tests {
		args := append([]string{command}, options...)
		args = append(args, "--help")
		var output bytes.Buffer
		if err := execute(args, time.Now, &output); err != nil {
			t.Fatalf("execute(%q) error = %v", args, err)
		}
		if got := output.String(); !strings.Contains(got, command) || !strings.Contains(got, "--repository") {
			t.Fatalf("execute(%q) omitted command or repository flag: %q", args, got)
		}
	}
}

func TestUnknownCommandFails(t *testing.T) {
	if err := execute([]string{"unknown"}, time.Now, &bytes.Buffer{}); err == nil {
		t.Fatal("unknown command succeeded")
	}
}

func TestMutationRejectsInspectionContextOptions(t *testing.T) {
	for _, args := range [][]string{
		{"--channel", "release", "nightly-tag", "--dry-run"},
		{"--release", "2026.35", "release-start", "--dry-run"},
	} {
		if err := execute(args, time.Now, &bytes.Buffer{}); err == nil ||
			!strings.Contains(err.Error(), "cannot be used with ref mutations") {
			t.Fatalf("execute(%q) error = %v", args, err)
		}
	}
}

func TestInvalidGlobalContextFailsBeforeDispatch(t *testing.T) {
	if err := execute(
		[]string{"--channel", "invalid", "nightly-tag", "--dry-run"},
		time.Now,
		&bytes.Buffer{},
	); err == nil || !strings.Contains(err.Error(), "invalid channel") {
		t.Fatalf("execute(invalid channel) error = %v", err)
	}
}

func TestWorkspaceStatusArgumentPrecedesTargetSeparator(t *testing.T) {
	got := insertBeforeSeparator(
		[]string{"run", "--config=release", "//some:target", "--", "target-arg"},
		"--workspace_status_command=versioning bazel-status",
	)
	want := []string{
		"run",
		"--config=release",
		"//some:target",
		"--workspace_status_command=versioning bazel-status",
		"--",
		"target-arg",
	}
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("insertBeforeSeparator() = %q, want %q", got, want)
	}
}

func TestShellQuote(t *testing.T) {
	if got := shellQuote("path with ' quote"); got != `'path with '"'"' quote'` {
		t.Fatalf("shellQuote() = %q", got)
	}
}

// TestLaunchBazelPreservesChildExitCode checks the built CLI and nested runner together.
func TestLaunchBazelPreservesChildExitCode(t *testing.T) {
	directory := t.TempDir()
	fakeBazel := filepath.Join(directory, "bazel")
	if err := os.WriteFile(fakeBazel, []byte("#!/bin/sh\nexit 7\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	binary, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(binary, "version", "--repo", directory, "bazel", "--", "build", "//app")
	command.Env = append(os.Environ(), "PATH="+directory+string(os.PathListSeparator)+os.Getenv("PATH"))
	output, err := command.CombinedOutput()
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 7 {
		t.Fatalf("stamped AL Bazel exit = %v, want 7; output: %s", err, output)
	}
}

func TestReleasePlanWritesReviewedJSON(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, directory, "init", "--initial-branch=master")
	runTestGit(t, directory, "config", "user.name", "Versioning Test")
	runTestGit(t, directory, "config", "user.email", "versioning@example.invalid")
	runTestGit(t, directory, "commit", "--allow-empty", "-m", "base")
	runTestGit(t, directory, "tag", "v2026.35.0-nightly.20260830")
	var output bytes.Buffer
	if err := execute([]string{
		"--repo", directory,
		"release-plan",
		"--plan", filepath.Join(directory, "out", "plan.json"),
	}, func() time.Time {
		return time.Date(2026, time.August, 30, 0, 0, 0, 0, time.UTC)
	}, &output); err != nil {
		t.Fatalf("release-plan error = %v", err)
	}
	if !strings.Contains(output.String(), `"schema": "versioning/release_ref_plan/v1"`) {
		t.Fatalf("release-plan output missing schema: %q", output.String())
	}
}

func TestReleasePlanRequiresPlanPath(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"--repo", t.TempDir(), "release-plan"}, time.Now, &output); err == nil ||
		!strings.Contains(err.Error(), "--plan") {
		t.Fatalf("release-plan() error = %v, want --plan refusal", err)
	}
}

func TestReleasePublishRequiresPlanPath(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"--repo", t.TempDir(), "release-publish"}, time.Now, &output); err == nil ||
		!strings.Contains(err.Error(), "--plan") {
		t.Fatalf("release-publish() error = %v, want --plan refusal", err)
	}
}

func TestReleasePublishRejectsMismatchedAtomicFlag(t *testing.T) {
	directory := t.TempDir()
	plan := `{
		"schema": "versioning/release_ref_plan/v1",
		"kind": "ReleaseRefPlan",
		"version": "2026.35.0",
		"channel": "release",
		"commit": "` + strings.Repeat("a", 40) + `",
		"treeState": "clean",
		"refs": [
			{"ref": "refs/heads/releases/2026.35", "expected": "` + strings.Repeat("a", 40) + `", "operation": "create"},
			{"ref": "refs/tags/v2026.35.0", "expected": "` + strings.Repeat("a", 40) + `", "operation": "create"}
		],
		"atomic": true,
		"lease": "lease"
	}`
	planPath := filepath.Join(directory, "plan.json")
	if err := os.WriteFile(planPath, []byte(plan), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := execute([]string{
		"--repo", directory,
		"release-publish",
		"--plan", planPath,
		"--receipt", filepath.Join(directory, "receipt.json"),
		"--atomic=false",
	}, time.Now, &output); err == nil || !strings.Contains(err.Error(), "--atomic must match") {
		t.Fatalf("release-publish() error = %v, want atomic mismatch refusal", err)
	}
}

func runTestGit(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, args...)...)
	var stdout bytes.Buffer
	command.Stdout = &stdout
	if err := command.Run(); err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(stdout.String())
}

// execute drives the native command tree with a fixture clock and output stream.
func execute(args []string, now func() time.Time, stdout io.Writer) error {
	cmd := newCommand(now)
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(io.Discard)
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("execute version fixture: %w", err)
	}
	return nil
}
