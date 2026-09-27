package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Exercise the built CLI and inspect its portable output independently of the generator.
func TestReleaseBundleCLI(t *testing.T) {
	dir := t.TempDir()
	status := filepath.Join(dir, "status.txt")
	payload := filepath.Join(dir, "website.tar.gz")
	body := []byte("archive fixture bytes\n")
	writeBundleFixture(t, status, "STABLE_VERSION 2026.39.2\nSTABLE_VERSION_CHANNEL release\n")
	writeBundleFixture(t, payload, string(body))
	output := filepath.Join(dir, "release")
	cmd := releaseBundleCommand("generate", "--project", "projects/alwaldend.com", "--version_file", status,
		"--add_file", payload, "--output_dir", output)
	if result, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generate bundle: %v\n%s", err, result)
	}
	manifest, err := os.ReadFile(filepath.Join(output, "release.json"))
	if err != nil {
		t.Fatal(err)
	}
	var release struct {
		Name    string `json:"name"`
		Project struct {
			Subdir string `json:"subdir"`
		} `json:"project"`
		Items []struct {
			File struct {
				Name   string                           `json:"name"`
				Hashes []struct{ Algo, Content string } `json:"hashes"`
			} `json:"file"`
		} `json:"items"`
	}
	if err := json.Unmarshal(manifest, &release); err != nil {
		t.Fatal(err)
	}
	if release.Name != "2026.39.2" || release.Project.Subdir != "projects/alwaldend.com" || len(release.Items) != 1 {
		t.Fatalf("incorrect release identity: %s", manifest)
	}
	file := release.Items[0].File
	got, err := os.ReadFile(filepath.Join(output, "files", file.Name))
	if err != nil || !bytes.Equal(got, body) {
		t.Fatalf("payload mismatch: %v", err)
	}
	want := fmt.Sprintf("%x", sha256.Sum256(got))
	found := false
	for _, hash := range file.Hashes {
		if hash.Algo == "SHA-256" && hash.Content == want {
			found = true
		}
	}
	if !found {
		t.Fatal("manifest does not describe the bundled bytes")
	}
	if artifactDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifactDir != "" {
		if err := os.WriteFile(filepath.Join(artifactDir, "release.json"), manifest, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestReleaseBundleFailuresCLI(t *testing.T) {
	for _, name := range []string{"missing_version", "duplicate_version", "unsafe_version", "missing_project", "missing_file", "duplicate_filename", "no_files"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			status := filepath.Join(dir, "status.txt")
			payload := filepath.Join(dir, "archive.zip")
			writeBundleFixture(t, status, "STABLE_VERSION 0.0.0-dev\n")
			writeBundleFixture(t, payload, "payload")
			project := "projects/example"
			files := []string{payload}
			switch name {
			case "missing_version":
				writeBundleFixture(t, status, "BUILD_USER fixture\n")
			case "duplicate_version":
				writeBundleFixture(t, status, "STABLE_VERSION 1.0.0\nSTABLE_VERSION 2.0.0\n")
			case "unsafe_version":
				writeBundleFixture(t, status, "STABLE_VERSION ../outside\n")
			case "missing_project":
				project = ""
			case "missing_file":
				files = []string{filepath.Join(dir, "missing.zip")}
			case "duplicate_filename":
				other := filepath.Join(dir, "other", "archive.zip")
				writeBundleFixture(t, other, "different payload")
				files = append(files, other)
			case "no_files":
				files = nil
			}
			output := filepath.Join(dir, "release")
			args := []string{"generate", "--project", project, "--version_file", status, "--output_dir", output}
			for _, file := range files {
				args = append(args, "--add_file", file)
			}
			if result, err := releaseBundleCommand(args...).CombinedOutput(); err == nil {
				t.Fatalf("invalid input accepted: %s", result)
			}
			if _, err := os.Stat(filepath.Join(output, "release.json")); !os.IsNotExist(err) {
				t.Fatalf("invalid bundle appears complete: %v", err)
			}
		})
	}
}

func releaseBundleCommand(args ...string) *exec.Cmd {
	return exec.Command(filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"),
		"tools/release/main/go/go_/go"), args...)
}

func writeBundleFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
