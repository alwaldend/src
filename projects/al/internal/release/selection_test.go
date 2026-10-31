package release

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestSelectedDestinationCLI verifies a symlink payload and environment-filtered
// publication through activation, retaining the published manifest as evidence.
func TestSelectedDestinationCLI(t *testing.T) {
	dir := t.TempDir()
	status, payload := filepath.Join(dir, "status"), filepath.Join(dir, "website.tar.gz")
	writeBundleFixture(t, status, "STABLE_VERSION 2026.44.0\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	body := []byte("<html>selected environment</html>")
	if err := tw.WriteHeader(&tar.Header{Name: "index.html", Mode: 0o644, Size: int64(len(body))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	writeBundleFixture(t, payload, archive.String())
	other := filepath.Join(dir, "other.txt")
	writeBundleFixture(t, other, "other environment")
	bundle := filepath.Join(dir, "bundle")
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--add_file", payload, "--add_file", other, "--output_dir", bundle).CombinedOutput(); err != nil {
		t.Fatalf("generate: %v\n%s", err, output)
	}
	manifestPath := filepath.Join(bundle, "release.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	remote := filepath.Join(dir, "remote")
	if err := os.MkdirAll(filepath.Join(remote, "staging"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, value := range manifest["items"].([]any) {
		item := value.(map[string]any)
		environment := "production"
		if item["file"].(map[string]any)["name"] == "other.txt" {
			environment = "development"
		}
		item["deployments"] = []any{map[string]any{"ssh": map[string]any{"environment": environment, "host": "fixture.invalid", "root": remote}}}
	}
	data, err = json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(manifestPath, 0o644); err != nil {
		t.Fatal(err)
	}
	writeBundleFixture(t, manifestPath, string(data))
	link := filepath.Join(bundle, "files", "website.tar.gz")
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(payload, link); err != nil {
		t.Fatal(err)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"publish", "--release_dir", bundle, "--environment", "production", "--ssh_path", helper, "--rsync_path", helper}
	preview := append([]string{"plan"}, args[1:]...)
	preview = append(preview, "--ssh_path", "/bin/false", "--rsync_path", "/bin/false")
	output, err := releaseBundleCommand(preview...).Output()
	if err != nil {
		t.Fatalf("preview selected destination: %v\n%s", err, output)
	}
	var plan map[string]any
	if err := json.Unmarshal(output, &plan); err != nil {
		t.Fatal(err)
	}
	operation := plan["operations"].([]any)[0].(map[string]any)
	files := operation["files"].([]any)
	if operation["selection"] != "environment" || len(files) != 1 || files[0].(map[string]any)["name"] != "website.tar.gz" {
		t.Fatalf("preview did not preserve environment selection: %s", output)
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		writeBundleFixture(t, filepath.Join(destination, "selected-plan.json"), string(output))
	}
	for _, invocation := range [][]string{
		append([]string{"release"}, args...),
		append(append([]string{"release", "website", "extract"}, args[1:]...), "--site_archive", "website.tar.gz"),
		append([]string{"release", "website", "activate"}, args[1:]...),
	} {
		cmd := artifactCommand(invocation...)
		cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1", "TMPDIR="+dir)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("selected destination operation %s: %v\n%s", invocation[1], err, output)
		}
	}
	selected, err := os.ReadFile(filepath.Join(remote, "sites/fixture/current/index.html"))
	if err != nil || !bytes.Equal(selected, body) {
		t.Fatalf("selected site %q: %v", selected, err)
	}
	if _, err := os.Stat(filepath.Join(remote, "projects/fixture/releases/2026.44.0/other.txt")); !os.IsNotExist(err) {
		t.Fatalf("unselected payload present: %v", err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		published, err := os.ReadFile(filepath.Join(remote, "projects/fixture/releases/2026.44.0/release.json"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "selected-publication.json"), published, 0o644); err != nil {
			t.Fatal(fmt.Errorf("record selected publication: %w", err))
		}
	}
}
