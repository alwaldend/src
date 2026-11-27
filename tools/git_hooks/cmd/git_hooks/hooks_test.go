package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstallAndCheckHook verifies installed bytes and permissions.
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

// TestInstallHookReplacesStaleHook checks atomic replacement of an obsolete hook.
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

// TestCheckHookFailures checks diagnostics for missing, obsolete, and nonexecutable hooks.
func TestCheckHookFailures(t *testing.T) {
	tests := []struct {
		// name identifies the failure case.
		name string
		// prepare creates the input state for verification.
		prepare func(t *testing.T, path string)
		// wantError identifies the expected diagnostic.
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
