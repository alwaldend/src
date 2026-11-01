package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrecommitRequiresBazelAgent(t *testing.T) {
	t.Parallel()
	binDirectory := t.TempDir()
	fallbackMarker := filepath.Join(t.TempDir(), "bazel-called")
	fakeBazel := filepath.Join(binDirectory, "bazel")
	if err := os.WriteFile(
		fakeBazel,
		[]byte("#!/bin/sh\n: > \"${BAZEL_FALLBACK_MARKER}\"\nexit 42\n"),
		0o755,
	); err != nil {
		t.Fatalf("os.WriteFile(fake bazel) error = %v", err)
	}

	command := exec.Command("/bin/sh")
	command.Stdin = bytes.NewReader(precommit)
	command.Env = []string{
		"BAZEL_FALLBACK_MARKER=" + fallbackMarker,
		"PATH=" + binDirectory,
	}
	output, err := command.CombinedOutput()
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 127 {
		t.Fatalf("pre-commit exit = %v, want 127; output:\n%s", err, output)
	}
	if !strings.Contains(string(output), "//projects/bazel_agent:install") {
		t.Fatalf("pre-commit output lacks bootstrap guidance:\n%s", output)
	}
	if _, err := os.Stat(fallbackMarker); err == nil {
		t.Fatal("pre-commit invoked the fallback bazel")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("os.Stat(fallback marker) error = %v", err)
	}
}

func TestInstallAndCheckHook(t *testing.T) {
	hookPath := filepath.Join(t.TempDir(), "hooks", "pre-commit")
	want := []byte("#!/usr/bin/env sh\nexit 0\n")

	if err := installHook(hookPath, want); err != nil {
		t.Fatalf("installHook() error = %v", err)
	}
	if err := checkHook(hookPath, want); err != nil {
		t.Fatalf("checkHook() error = %v", err)
	}

	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}
	if got := info.Mode().Perm(); got != 0o755 {
		t.Errorf("installed hook mode = %o, want 755", got)
	}
}

func TestInstallHookReplacesStaleHook(t *testing.T) {
	hookPath := filepath.Join(t.TempDir(), "pre-commit")
	if err := os.WriteFile(hookPath, []byte("stale"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	want := []byte("current")
	if err := installHook(hookPath, want); err != nil {
		t.Fatalf("installHook() error = %v", err)
	}
	if err := checkHook(hookPath, want); err != nil {
		t.Fatalf("checkHook() error = %v", err)
	}
}

func TestCheckHookFailures(t *testing.T) {
	tests := []struct {
		name      string
		prepare   func(t *testing.T, path string)
		wantError string
	}{
		{
			name:      "missing",
			prepare:   func(t *testing.T, path string) {},
			wantError: "is not installed",
		},
		{
			name: "stale",
			prepare: func(t *testing.T, path string) {
				if err := os.WriteFile(path, []byte("stale"), 0o755); err != nil {
					t.Fatalf("os.WriteFile() error = %v", err)
				}
			},
			wantError: "is out of date",
		},
		{
			name: "not executable",
			prepare: func(t *testing.T, path string) {
				if err := os.WriteFile(path, []byte("current"), 0o644); err != nil {
					t.Fatalf("os.WriteFile() error = %v", err)
				}
			},
			wantError: "is not executable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hookPath := filepath.Join(t.TempDir(), "pre-commit")
			test.prepare(t, hookPath)
			err := checkHook(hookPath, []byte("current"))
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("checkHook() error = %v, want %q", err, test.wantError)
			}
		})
	}
}
