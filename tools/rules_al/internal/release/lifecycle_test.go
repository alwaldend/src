package release

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestImmutableReleaseCLI exercises publication, activation, deletion, and reuse
// through the built command and the destination shell protocol in a local tree.
func TestImmutableReleaseCLI(t *testing.T) {
	dir := t.TempDir()
	status, payload := filepath.Join(dir, "status"), filepath.Join(dir, "website.tar.gz")
	writeBundleFixture(t, status, "STABLE_VERSION 2026.44.0\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	body := "<html>release fixture</html>"
	if err := tw.WriteHeader(&tar.Header{Name: "index.html", Mode: 0o644, Size: int64(len(body))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(body)); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	writeBundleFixture(t, payload, archive.String())
	bundle := filepath.Join(dir, "bundle")
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--add_file", payload, "--output_dir", bundle).CombinedOutput(); err != nil {
		t.Fatalf("generate: %v\n%s", err, output)
	}
	rival := filepath.Join(dir, "rival")
	writeBundleFixture(t, payload, "different artifact bytes")
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--add_file", payload, "--output_dir", rival).CombinedOutput(); err != nil {
		t.Fatalf("generate conflicting bundle: %v\n%s", err, output)
	}
	ociDeployment, ociBundle := filepath.Join(dir, "oci.json"), filepath.Join(dir, "oci-bundle")
	writeBundleFixture(t, ociDeployment, `{"oci":{"repository":"registry.invalid/fixture"}}`)
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--deployment", ociDeployment, "--add_file", payload, "--output_dir", ociBundle).CombinedOutput(); err != nil {
		t.Fatalf("generate OCI metadata: %v\n%s", err, output)
	}
	if output, err := releaseBundleCommand("deploy", "--release_dir", ociBundle, "--soras_path", "/bin/true").CombinedOutput(); err == nil || !bytes.Contains(output, []byte("OCI publication is disabled")) {
		t.Fatalf("mutable OCI publication was not refused: %v\n%s", err, output)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	remote := filepath.Join(dir, "remote")
	if err := os.MkdirAll(filepath.Join(remote, "staging"), 0o755); err != nil {
		t.Fatal(err)
	}
	base := []string{"--release_dir", bundle, "--environment", "fixture", "--ssh_host", "fixture.invalid", "--ssh_root", remote, "--ssh_path", helper, "--rsync_path", helper}
	mount := ""
	bin := filepath.Join(dir, "bin")
	writeBundleFixture(t, filepath.Join(bin, "mountpoint"), "#!/bin/sh\n[ \"$3\" = \"$RELEASE_FIXTURE_MOUNT\" ]\n")
	if err := os.Chmod(filepath.Join(bin, "mountpoint"), 0o755); err != nil {
		t.Fatal(err)
	}
	run := func(command string, fail bool, extra ...string) {
		t.Helper()
		args := append([]string{command}, base...)
		cmd := releaseBundleCommand(append(args, extra...)...)
		cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1", "TMPDIR="+dir, "PATH="+bin+":"+os.Getenv("PATH"), "RELEASE_FIXTURE_MOUNT="+mount)
		output, err := cmd.CombinedOutput()
		if (err != nil) != fail {
			t.Fatalf("%s %v: %v\n%s", command, extra, err, output)
		}
	}
	public := filepath.Join(remote, "projects/fixture/releases/2026.44.0")
	run("publish", true, "--rsync_path", "/bin/false")
	if _, err := os.Stat(public); !os.IsNotExist(err) {
		t.Fatalf("failed transfer exposed a release: %v", err)
	}
	receipt := filepath.Join(dir, "publication-receipt.json")
	run("publish", false, "--receipt", receipt)
	checkReceipt := func(status string, states ...string) {
		t.Helper()
		data, err := os.ReadFile(receipt)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(receipt)
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("receipt permissions: %v", err)
		}
		var observed struct {
			// Schema identifies artifact operation evidence.
			Schema string `json:"schema"`
			// Status reports whether every requested operation completed.
			Status string `json:"status"`
			// Operations records one resolved destination and its ordered steps.
			Operations []struct {
				// Project identifies the independent artifact owner.
				Project string `json:"project"`
				// Version is the immutable version selected by this attempt.
				Version string `json:"version"`
				// PublicationDigest binds the selected canonical publication manifest.
				PublicationDigest string `json:"publication_digest"`
				// Destination records the resolved environment and content root.
				Destination struct {
					// Environment names the explicitly selected deployment environment.
					Environment string `json:"environment"`
					// Host identifies the selected SSH destination.
					Host string `json:"host"`
					// Root identifies its content directory.
					Root string `json:"root"`
				} `json:"destination"`
				// Steps preserves completed and uncertain observations in order.
				Steps []struct {
					// Status records pending, unknown, or completed.
					Status string `json:"status"`
				} `json:"steps"`
			} `json:"operations"`
		}
		if err := json.Unmarshal(data, &observed); err != nil {
			t.Fatal(err)
		}
		if observed.Schema != "release/artifact_receipt/v1" || observed.Status != status || len(observed.Operations) != 1 {
			t.Fatalf("invalid receipt: %s", data)
		}
		operation := observed.Operations[0]
		manifest, err := os.ReadFile(filepath.Join(public, "release.json"))
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(manifest)
		if operation.Project != "fixture" || operation.Version != "2026.44.0" || operation.PublicationDigest != "sha256:"+hex.EncodeToString(digest[:]) || operation.Destination.Host != "fixture.invalid" || operation.Destination.Environment != "fixture" || operation.Destination.Root != remote || len(operation.Steps) != len(states) {
			t.Fatalf("receipt identity mismatch: %s", data)
		}
		for index, state := range states {
			if operation.Steps[index].Status != state {
				t.Fatalf("receipt step mismatch: %s", data)
			}
		}
		if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
			writeBundleFixture(t, filepath.Join(output, "artifact-"+status+".json"), string(data))
		}
	}
	checkReceipt("completed", "completed")
	run("publish", true, "--receipt", receipt, "--release_dir", filepath.Join(dir, "missing-bundle"))
	checkReceipt("incomplete", "completed")
	run("publish", true, "--receipt", receipt, "--steps", "upload,extract,link", "--site_archive", "website.tar.gz", "--rsync_path", "/bin/false")
	checkReceipt("incomplete", "completed", "unknown", "pending")
	// Cancel during archive retrieval, after upload has already been confirmed.
	blocker, ready := filepath.Join(bin, "blocking-rsync"), filepath.Join(dir, "transfer-ready")
	writeBundleFixture(t, blocker, "#!/bin/sh\nprintf ready > \"$RELEASE_TRANSFER_READY\"\nexec sleep 30\n")
	if err := os.Chmod(blocker, 0o755); err != nil {
		t.Fatal(err)
	}
	interrupted := releaseBundleCommand(append(append([]string{"publish"}, base...), "--receipt", receipt, "--steps", "upload,extract,link", "--site_archive", "website.tar.gz", "--rsync_path", blocker)...)
	interrupted.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1", "TMPDIR="+dir, "PATH="+bin+":"+os.Getenv("PATH"), "RELEASE_TRANSFER_READY="+ready)
	if err := interrupted.Start(); err != nil {
		t.Fatal(err)
	}
	defer interrupted.Process.Kill()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if _, err := os.Stat(ready); err != nil {
		t.Fatal(err)
	}
	checkReceipt("running", "completed", "unknown", "pending")
	if err := interrupted.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() { finished <- interrupted.Wait() }()
	select {
	case err := <-finished:
		if err == nil {
			t.Fatal("cancelled publication succeeded")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancelled publication did not stop")
	}
	checkReceipt("incomplete", "completed", "unknown", "pending")
	// A checkpoint failure must precede all remote effects.
	run("publish", true, "--receipt", filepath.Join(receipt, "child.json"), "--project", "unstarted")
	if _, err := os.Stat(filepath.Join(remote, "projects/unstarted")); !os.IsNotExist(err) {
		t.Fatalf("receipt failure started publication: %v", err)
	}
	// An identical retry needs neither local payloads nor a working transfer tool.
	if err := os.Rename(filepath.Join(bundle, "files"), filepath.Join(bundle, "held-files")); err != nil {
		t.Fatal(err)
	}
	run("publish", false, "--rsync_path", "/bin/false")
	if err := os.Rename(filepath.Join(bundle, "held-files"), filepath.Join(bundle, "files")); err != nil {
		t.Fatal(err)
	}
	manifest, err := os.ReadFile(filepath.Join(public, "release.json"))
	if err != nil {
		t.Fatal(err)
	}
	// Concurrent identical publishers must both succeed and preserve the bytes.
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); run("publish", false) }()
	}
	wg.Wait()
	base[1] = rival
	run("publish", true)
	base[1] = bundle
	run("publish", false, "--project", "other")
	writeBundleFixture(t, filepath.Join(bundle, "files/website.tar.gz"), "tampered")
	run("publish", false)
	writeBundleFixture(t, filepath.Join(bundle, "files/website.tar.gz"), archive.String())
	writeBundleFixture(t, filepath.Join(public, "website.tar.gz"), "corrupted destination")
	run("publish", true)
	run("delete", false)
	run("publish", false, "--steps", "upload,extract,link", "--site_archive", "website.tar.gz")
	run("delete", true)
	writeBundleFixture(t, status, "STABLE_VERSION 2026.44.1\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
	next := filepath.Join(dir, "next")
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--add_file", filepath.Join(bundle, "files/website.tar.gz"), "--output_dir", next).CombinedOutput(); err != nil {
		t.Fatalf("generate next version: %v\n%s", err, output)
	}
	base[1] = next
	run("publish", false, "--steps", "upload,extract,link", "--site_archive", "website.tar.gz")
	base[1] = bundle
	run("publish", false, "--steps", "link")
	secondSite := filepath.Join(remote, "sites/fixture/releases/2026.44.1")
	if _, err := os.Stat(filepath.Join(secondSite, "index.html")); err != nil {
		t.Fatalf("rollback removed the retained site: %v", err)
	}
	for _, tree := range []string{"sites", "projects"} {
		original, held := filepath.Join(remote, tree), filepath.Join(remote, tree+"-held")
		if err := os.Rename(original, held); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(held, original); err != nil {
			t.Fatal(err)
		}
		run("publish", true, "--steps", "link")
		base[1] = next
		run("delete", true)
		base[1] = bundle
		if _, err := os.Stat(filepath.Join(secondSite, "index.html")); err != nil {
			t.Fatalf("symlink refusal removed retained content: %v", err)
		}
		if err := os.Remove(original); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(held, original); err != nil {
			t.Fatal(err)
		}
	}
	base[1] = next
	for _, mounted := range []string{secondSite, filepath.Join(secondSite, "assets"), filepath.Join(remote, "sites/fixture/releases"), filepath.Join(remote, "projects/fixture/releases/2026.44.1")} {
		if err := os.MkdirAll(mounted, 0o755); err != nil {
			t.Fatal(err)
		}
		mount = mounted
		run("delete", true)
		if _, err := os.Stat(filepath.Join(secondSite, "index.html")); err != nil {
			t.Fatalf("mount refusal removed retained content: %v", err)
		}
	}
	for _, mounted := range []string{filepath.Join(secondSite, "index.html"), filepath.Join(remote, "projects/fixture/releases/2026.44.1/website.tar.gz")} {
		mount = mounted
		run("delete", true)
		for _, retained := range []string{filepath.Join(secondSite, "index.html"), filepath.Join(remote, "projects/fixture/releases/2026.44.1/release.json")} {
			if _, err := os.Stat(retained); err != nil {
				t.Fatalf("file mount refusal removed retained content: %v", err)
			}
		}
	}
	mount = ""
	run("delete", false)
	base[1] = bundle
	if err := os.Remove(filepath.Join(remote, "sites/fixture/current")); err != nil {
		t.Fatal(err)
	}
	run("delete", false)
	run("delete", false)
	for _, path := range []string{public, filepath.Join(remote, "sites/fixture/releases/2026.44.0")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("deleted path remains %s: %v", path, err)
		}
	}
	// No tombstone prevents ordinary publication of different bytes after deletion.
	base[1] = rival
	run("publish", false)
	replacement, err := os.ReadFile(filepath.Join(public, "website.tar.gz"))
	if err != nil || string(replacement) != "different artifact bytes" {
		t.Fatalf("republished content: %q, %v", replacement, err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "immutable-release.json"), manifest, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
