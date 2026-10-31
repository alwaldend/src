package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstalledPayload verifies that the copied executable can start on its own.
func TestInstalledPayload(t *testing.T) {
	if os.Getenv("INSTALL_BINARY_PAYLOAD") == "1" {
		t.Log("installed executable started")
	}
}

// TestInstallBinary exercises the built installer and the resulting executable.
func TestInstallBinary(t *testing.T) {
	installer := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"),
		"tools/file_installer/cmd/install_binary/install_binary_/install_binary")
	source, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(t.TempDir(), "with spaces", "bin")
	arguments := []string{"--source", source, "--name", "fixture", "--directory", directory}
	for range 2 {
		if output, err := exec.Command(installer, arguments...).CombinedOutput(); err != nil {
			t.Fatalf("install: %v: %s", err, output)
		}
	}
	installed := filepath.Join(directory, "fixture")
	content, err := os.ReadFile(installed)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(content) != sha256.Sum256(expected) {
		t.Fatal("installed bytes differ from source")
	}
	info, err := os.Stat(installed)
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("installed permissions: %v: %v", info, err)
	}
	command := exec.Command(installed, "-test.run=^TestInstalledPayload$", "-test.v")
	command.Env = append(os.Environ(), "INSTALL_BINARY_PAYLOAD=1")
	if output, err := command.CombinedOutput(); err != nil || !strings.Contains(string(output), "installed executable started") {
		t.Fatalf("run installed executable: %v: %s", err, output)
	}
	for _, invalid := range [][]string{
		{"--source", source + ".missing", "--name", "fixture", "--directory", directory},
		{"--source", source, "--name", "../escape", "--directory", directory},
		{"--source", directory, "--name", "fixture", "--directory", directory},
	} {
		if output, err := exec.Command(installer, invalid...).CombinedOutput(); err == nil {
			t.Fatalf("invalid install succeeded: %s", output)
		}
	}
	after, err := os.ReadFile(installed)
	if err != nil || sha256.Sum256(after) != sha256.Sum256(expected) {
		t.Fatalf("failed update changed prior installation: %v", err)
	}
	entries, err := os.ReadDir(directory)
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary installation residue: %v: %v", entries, err)
	}
	digest := sha256.Sum256(content)
	receipt, err := json.MarshalIndent(map[string]any{
		"sha256": hex.EncodeToString(digest[:]), "mode": "0755",
		"executed": true, "replaced": true, "failed_update_preserved": true,
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "installation.json"), receipt, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
