package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestStateSubcommandErrors rejects incomplete and unknown operations before backend access.
func TestStateSubcommandErrors(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	for _, args := range [][]string{
		{"state"},
		{"state", "--collection", "unused.anki2"},
		{"state", "nonsense"},
		{"state", "nonsense", "--collection", "unused.anki2"},
	} {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		command := exec.CommandContext(ctx, binary, args...)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		cancel()
		if err == nil || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("invalid state command %v: error=%v stdout=%q stderr=%q", args, err, stdout.String(), stderr.String())
		}
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, binary, "state", "--help")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil || !strings.Contains(stdout.String(), "Available Commands:") || stderr.Len() != 0 {
		t.Fatalf("explicit state help: error=%v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
}

// TestStateRequiredFlags verifies missing state collection inputs fail with clean output streams.
func TestStateRequiredFlags(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	for _, args := range [][]string{
		{"state", "import", "--text", "unused", "notes/example.note.anki.toml=1"},
		{"state", "update", "unused"},
	} {
		var stdout, stderr bytes.Buffer
		command := exec.CommandContext(ctx, binary, args...)
		command.Stdout, command.Stderr = &stdout, &stderr
		if err := command.Run(); err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "required flag") {
			t.Fatalf("required state flags: error=%v stdout=%q stderr=%q", err, stdout.String(), stderr.String())
		}
	}
}

// TestStateCommandSyntax verifies binding errors are reported before collection access.
func TestStateCommandSyntax(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	for _, arguments := range [][]string{
		{"bad"},
		{"=opaque"},
		{"notes/a.note.anki.toml=1", "notes/a.note.anki.toml=2"},
	} {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		args := append([]string{"state", "import", "--collection", "unused.anki2", "--text", "missing"}, arguments...)
		command := exec.CommandContext(ctx, binary, args...)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		cancel()
		if err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "binding") || strings.Contains(stderr.String(), "unknown command") {
			t.Fatalf("invalid binding %v: error=%v stdout=%q stderr=%q", arguments, err, stdout.String(), stderr.String())
		}
	}
	for _, identifier := range []string{"opaque-id", "", "0", "-1", "9223372036854775808"} {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		command := exec.CommandContext(ctx, binary, "state", "import", "--collection", "unused.anki2", "--text", "missing", "notes/a.note.anki.toml="+identifier)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		cancel()
		if err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "load import configs") {
			t.Fatalf("identifier %q rejected during parsing: error=%v stdout=%q stderr=%q", identifier, err, stdout.String(), stderr.String())
		}
	}
}

// TestStateInputFailures rejects unavailable or malformed sources without creating collections.
func TestStateInputFailures(t *testing.T) {
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	root := t.TempDir()
	archive := filepath.Join(root, "collection.colpkg")
	if err := os.WriteFile(archive, []byte("invalid collection archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{filepath.Join(root, "missing.anki2"), root, archive} {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		command := exec.CommandContext(ctx, binary, "state", "dump", "--collection", input)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		cancel()
		if err == nil || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("invalid collection input %q: error=%v stdout=%q stderr=%q", input, err, stdout.String(), stderr.String())
		}
	}
	if _, err := os.Stat(filepath.Join(root, "missing.anki2")); !os.IsNotExist(err) {
		t.Fatalf("state command created a missing database: %v", err)
	}
}
