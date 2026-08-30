package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
		"show":          {"--format=json"},
		"nightly-tag":   {"--date=2026-08-30", "--dry-run"},
		"release-start": {"--date=2026-08-30", "--switch=false"},
		"release-tag":   {"--dry-run"},
	}
	for command, options := range tests {
		args := append([]string{command}, options...)
		args = append(args, "--help")
		var output bytes.Buffer
		if err := execute(args, time.Now, &output); err != nil {
			t.Fatalf("execute(%q) error = %v", args, err)
		}
		if got, want := output.String(), commandUsage[command]; got != want {
			t.Fatalf("execute(%q) usage = %q, want %q", args, got, want)
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

func TestProcessExitCodePreservesBazelExitCode(t *testing.T) {
	err := &bazelExitError{err: errors.New("Bazel failed"), code: 7}
	if got := processExitCode(err); got != 7 {
		t.Fatalf("processExitCode() = %d, want 7", got)
	}
	if got := processExitCode(errors.New("ordinary failure")); got != 1 {
		t.Fatalf("processExitCode(ordinary) = %d, want 1", got)
	}
	if got := processExitCode(ordinaryExitError{}); got != 1 {
		t.Fatalf("processExitCode(non-Bazel exit error) = %d, want 1", got)
	}
}

type ordinaryExitError struct{}

func (ordinaryExitError) Error() string { return "ordinary subprocess failure" }

func (ordinaryExitError) ExitCode() int { return 9 }

func TestLaunchBazelPreservesChildExitCode(t *testing.T) {
	directory := t.TempDir()
	executable := filepath.Join(directory, "bazel_agent")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nexit 7\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))
	err := launchBazel(directory, "master", "auto", "", []string{"build", "//app"})
	if err == nil {
		t.Fatal("launchBazel() succeeded")
	}
	if got := processExitCode(err); got != 7 {
		t.Fatalf("processExitCode(launchBazel error) = %d, want 7: %v", got, err)
	}
}
