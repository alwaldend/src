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

// TestCLIWorkflow verifies command-level export, identity generation, planning, and rebuild.
func TestCLIWorkflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	binary := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), os.Getenv("ANKI_CLI"))
	root := t.TempDir()
	note := filepath.Join(root, "example.note.anki.toml")
	before := "note_id = 100\nguid = 'old'\nnote_type_id = 10\nnote_type = 'Basic'\n[fields]\nFront = 'preserve this field'\n"
	if err := os.WriteFile(note, []byte(before), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	command := exec.CommandContext(ctx, binary, "generate-id", note)
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("generate identities: %v: %s", err, stderr.String())
	}
	if stdout.Len() != 0 || !strings.Contains(stderr.String(), "Generated identities for 1 notes") {
		t.Fatalf("status must use stderr: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	after, err := os.ReadFile(note)
	if err != nil || strings.Contains(string(after), "note_id = 100\n") || !strings.Contains(string(after), "preserve this field") {
		t.Fatalf("identity operation did not preserve content: %q %v", after, err)
	}
	for _, args := range [][]string{
		{"export", "--output", "unused"},
		{"plan", "--text", "unused"},
		{"apply", "--text", "unused"},
		{"build", "--text", "unused", "--output", "unused"},
	} {
		stdout.Reset()
		stderr.Reset()
		command = exec.CommandContext(ctx, binary, args...)
		command.Stdout, command.Stderr = &stdout, &stderr
		if err := command.Run(); err == nil || stdout.Len() != 0 || !strings.Contains(stderr.String(), "required flag") {
			t.Fatalf("required flags for %s: error=%v stdout=%q stderr=%q", args[0], err, stdout.String(), stderr.String())
		}
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		if err := os.WriteFile(filepath.Join(artifacts, "cli-streams-and-flags.txt"), []byte("Built CLI generates identities without changing field content; status and errors use stderr; export/plan/apply/build enforce required flags.\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
