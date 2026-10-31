package release

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestReleaseObservationCLI verifies offline planning and independent local and
// remote observations through the built command, retaining every JSON report.
func TestReleaseObservationCLI(t *testing.T) {
	dir := t.TempDir()
	status, payload, bundle := filepath.Join(dir, "status"), filepath.Join(dir, "website.tar.gz"), filepath.Join(dir, "bundle")
	writeBundleFixture(t, status, "STABLE_VERSION 2026.44.0\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	body := "<html>observed release</html>"
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
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--add_file", payload, "--output_dir", bundle).CombinedOutput(); err != nil {
		t.Fatalf("generate observation fixture: %v\n%s", err, output)
	}
	helper, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	remote := filepath.Join(dir, "remote")
	if err := os.Mkdir(remote, 0o755); err != nil {
		t.Fatal(err)
	}
	base := []string{"--release_dir", bundle, "--environment", "fixture", "--ssh_host", "fixture.invalid", "--ssh_root", remote, "--ssh_path", helper, "--rsync_path", helper}
	sequence := 0
	observe := func(command string, fail bool, extra ...string) map[string]any {
		t.Helper()
		before := observationTree(t, remote)
		args := append(append([]string{command}, base...), extra...)
		cmd := releaseBundleCommand(args...)
		if command == "extract" || command == "activate" {
			cmd = websiteCommand(args...)
		}
		cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1", "TMPDIR="+dir)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		output, err := cmd.Output()
		if (err != nil) != fail {
			t.Fatalf("%s %v: %v\n%s\n%s", command, extra, err, output, &stderr)
		}
		if after := observationTree(t, remote); !reflect.DeepEqual(before, after) {
			t.Fatalf("%s changed the remote tree:\nbefore=%v\nafter=%v", command, before, after)
		}
		var report map[string]any
		if err := json.Unmarshal(output, &report); err != nil {
			t.Fatalf("decode %s: %v\n%s\n%s", command, err, output, &stderr)
		}
		sequence++
		if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
			writeBundleFixture(t, filepath.Join(destination, fmt.Sprintf("observation-%02d-%s.json", sequence, command)), string(output))
		}
		operations := report["operations"].([]any)
		if len(operations) != 1 {
			t.Fatalf("unexpected operations: %s", output)
		}
		return operations[0].(map[string]any)
	}
	assertState := func(operation map[string]any, local, remoteState, site string, active any) {
		t.Helper()
		files := operation["local"].([]any)
		observed := operation["remote"].(map[string]any)
		website := observed["site"].(map[string]any)
		if files[0].(map[string]any)["status"] != local || observed["status"] != remoteState || website["status"] != site || website["active"] != active {
			t.Fatalf("unexpected observation: %#v", operation)
		}
	}
	plan := observe("plan", false, "--ssh_path", "/bin/false", "--rsync_path", "/bin/false", "--ssh_identity", "/nonexistent/private-identity", "--ssh_user", "fixture-admin", "--publish_user", "fixture-content")
	if plan["selection"] != "explicit" || plan["project"] != "fixture" || plan["version"] != "2026.44.0" || plan["destination"].(map[string]any)["port"] != nil || len(plan["files"].([]any)) != 1 || len(plan["steps"].([]any)) != 1 {
		t.Fatalf("incorrect offline preview: %#v", plan)
	}
	if plan["destination"].(map[string]any)["ssh_user"] != "fixture-admin" || plan["destination"].(map[string]any)["publish_user"] != "fixture-content" {
		t.Fatalf("preview omitted configured execution accounts: %#v", plan)
	}
	for _, command := range []string{"extract", "activate"} {
		step := "extract"
		if command == "activate" {
			step = "link"
		}
		extra := []string{"--dry-run", "--ssh_path", "/bin/false", "--rsync_path", "/bin/false", "--ssh_identity", "/nonexistent/private-identity"}
		if command == "extract" {
			extra = append(extra, "--site_archive", "website.tar.gz")
		}
		preview := observe(command, false, extra...)
		if !reflect.DeepEqual(preview["steps"], []any{step}) || preview["publication_digest"] != plan["publication_digest"] {
			t.Fatalf("website preview changed selection: %#v", preview)
		}
		receipt := filepath.Join(dir, command+"-dry-run-receipt.json")
		args := append(append([]string{command}, base...), append(extra, "--receipt", receipt)...)
		if output, err := websiteCommand(args...).CombinedOutput(); err == nil || !bytes.Contains(output, []byte("--receipt")) {
			t.Fatalf("website dry-run accepted a receipt: %v\n%s", err, output)
		}
		if _, err := os.Stat(receipt); !os.IsNotExist(err) {
			t.Fatalf("website dry-run wrote a receipt: %v", err)
		}
	}
	assertState(observe("inspect", false), "verified", "absent", "absent", false)
	assertState(observe("inspect", true, "--ssh_root", filepath.Join(dir, "unprovisioned")), "verified", "unavailable", "unavailable", nil)
	blocked := filepath.Join(dir, "blocked")
	if err := os.MkdirAll(filepath.Join(blocked, "remote"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(blocked, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(blocked, 0o755) })
	assertState(observe("inspect", true, "--ssh_root", filepath.Join(blocked, "remote")), "verified", "unavailable", "unavailable", nil)
	if err := os.Chmod(blocked, 0o755); err != nil {
		t.Fatal(err)
	}
	// Authentication and malformed protocol responses must not become absence.
	assertState(observe("inspect", true, "--ssh_path", "/bin/false"), "verified", "unavailable", "unavailable", nil)
	assertState(observe("inspect", true, "--ssh_path", "/bin/echo"), "verified", "unavailable", "unavailable", nil)
	for name, script := range map[string]string{
		"timeout":   "#!/bin/sh\nexec sleep 10\n",
		"oversized": "#!/bin/sh\nexec yes observation-output\n",
	} {
		filename := filepath.Join(dir, name)
		writeBundleFixture(t, filename, script)
		if err := os.Chmod(filename, 0o755); err != nil {
			t.Fatal(err)
		}
		deadline := "100ms"
		if name == "oversized" {
			deadline = "10s"
		}
		started := time.Now()
		assertState(observe("inspect", true, "--ssh_path", filename, "--timeout", deadline), "verified", "unavailable", "unavailable", nil)
		if time.Since(started) > 2*time.Second {
			t.Fatalf("%s did not stop promptly", name)
		}
	}
	publish := func(group, command string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(remote, "staging"), 0o755); err != nil {
			t.Fatal(err)
		}
		prefix := []string{"release", command}
		if group == "website" {
			prefix = []string{"release", "website", command}
		}
		args := append(prefix, base...)
		if command == "extract" {
			args = append(args, "--site_archive", "website.tar.gz")
		}
		cmd := artifactCommand(args...)
		cmd.Env = append(os.Environ(), "RELEASE_TRANSPORT_HELPER=1", "TMPDIR="+dir)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s %s: %v\n%s", group, command, err, output)
		}
		if err := os.RemoveAll(filepath.Join(remote, "staging")); err != nil {
			t.Fatal(err)
		}
	}
	publish("release", "publish")
	assertState(observe("inspect", false), "verified", "matching", "absent", false)
	local := filepath.Join(bundle, "files/website.tar.gz")
	if err := os.Remove(local); err != nil {
		t.Fatal(err)
	}
	assertState(observe("inspect", false), "missing", "matching", "absent", false)
	observe("plan", false, "--ssh_path", "/bin/false")
	writeBundleFixture(t, local, "damaged local payload")
	assertState(observe("inspect", false), "corrupt", "matching", "absent", false)
	writeBundleFixture(t, local, archive.String())
	public := filepath.Join(remote, "projects/fixture/releases/2026.44.0")
	corrupt := bytes.Repeat([]byte{'x'}, archive.Len())
	writeBundleFixture(t, filepath.Join(public, "website.tar.gz"), string(corrupt))
	assertState(observe("inspect", false), "verified", "corrupt", "absent", false)
	writeBundleFixture(t, filepath.Join(public, "website.tar.gz"), archive.String())
	writeBundleFixture(t, filepath.Join(public, "extra"), "undeclared")
	assertState(observe("inspect", false), "verified", "corrupt", "absent", false)
	if err := os.Remove(filepath.Join(public, "extra")); err != nil {
		t.Fatal(err)
	}
	manifest, err := os.ReadFile(filepath.Join(public, "release.json"))
	if err != nil {
		t.Fatal(err)
	}
	writeBundleFixture(t, filepath.Join(public, "release.json"), "different publication binding\n")
	assertState(observe("inspect", false), "verified", "conflicting", "absent", false)
	writeBundleFixture(t, filepath.Join(public, "release.json"), string(manifest))
	publish("website", "extract")
	assertState(observe("inspect", false), "verified", "matching", "retained", false)
	publish("website", "activate")
	assertState(observe("inspect", false), "verified", "matching", "active", true)
	site := filepath.Join(remote, "sites/fixture/releases/2026.44.0")
	current := filepath.Join(remote, "sites/fixture/current")
	for _, target := range []string{"releases/2026.44.0\n", "releases/other\nversion"} {
		if err := os.Remove(current); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, current); err != nil {
			t.Fatal(err)
		}
		invalid := observe("inspect", false)
		assertState(invalid, "verified", "matching", "retained", nil)
		if invalid["remote"].(map[string]any)["site"].(map[string]any)["current"] != "unsafe" {
			t.Fatalf("control character in current link was not identified: %#v", invalid)
		}
	}
	if err := os.Remove(current); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(site, current); err != nil {
		t.Fatal(err)
	}
	unsafe := observe("inspect", false)
	assertState(unsafe, "verified", "matching", "retained", nil)
	if unsafe["remote"].(map[string]any)["site"].(map[string]any)["current"] != "unsafe" {
		t.Fatalf("absolute current link was not identified: %#v", unsafe)
	}
	if err := os.Remove(current); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(current, 0o755); err != nil {
		t.Fatal(err)
	}
	assertState(observe("inspect", false), "verified", "matching", "retained", nil)
	if err := os.Remove(current); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("releases/2026.44.0", current); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(site, site+"-held"); err != nil {
		t.Fatal(err)
	}
	dangling := observe("inspect", false)
	assertState(dangling, "verified", "matching", "absent", true)
	if dangling["remote"].(map[string]any)["site"].(map[string]any)["current"] != "dangling" {
		t.Fatalf("dangling selected link was not identified: %#v", dangling)
	}
	if err := os.Rename(site+"-held", site); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(site, "index.html")); err != nil {
		t.Fatal(err)
	}
	assertState(observe("inspect", false), "verified", "matching", "corrupt", true)
	writeBundleFixture(t, filepath.Join(site, "index.html"), body)
	writeBundleFixture(t, filepath.Join(site, ".release.json"), "different site binding\n")
	assertState(observe("inspect", false), "verified", "matching", "conflicting", true)
	writeBundleFixture(t, filepath.Join(site, ".release.json"), string(manifest))
	if err := os.Rename(public, public+"-held"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(public+"-held", public); err != nil {
		t.Fatal(err)
	}
	assertState(observe("inspect", false), "verified", "corrupt", "active", true)
	// Inspection accepts no mutation checkpoint option.
	if output, err := releaseBundleCommand(append(append([]string{"inspect"}, base...), "--receipt", filepath.Join(dir, "receipt"))...).CombinedOutput(); err == nil {
		t.Fatalf("inspection accepted a receipt: %s", output)
	}
}

// TestReleaseInspectionCancellationCLI interrupts verification of a sparse large
// payload after its file is open and hashing has begun, before SSH may start.
func TestReleaseInspectionCancellationCLI(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the controlled hash-start observation requires Linux process I/O counters")
	}
	dir := t.TempDir()
	status, payload, bundle := filepath.Join(dir, "status"), filepath.Join(dir, "payload"), filepath.Join(dir, "bundle")
	writeBundleFixture(t, status, "STABLE_VERSION 2026.44.0\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
	writeBundleFixture(t, payload, "fixture")
	if output, err := releaseBundleCommand("generate", "--project", "projects/fixture", "--version_file", status, "--add_file", payload, "--output_dir", bundle).CombinedOutput(); err != nil {
		t.Fatalf("generate cancellation fixture: %v\n%s", err, output)
	}
	local := filepath.Join(bundle, "files/payload")
	if err := os.Chmod(local, 0o644); err != nil {
		t.Fatal(err)
	}
	const size = int64(128 * 1024 * 1024 * 1024)
	if err := os.Truncate(local, size); err != nil {
		t.Fatal(err)
	}
	filename := filepath.Join(bundle, "release.json")
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	manifest["items"].([]any)[0].(map[string]any)["file"].(map[string]any)["size"] = strconv.FormatInt(size, 10)
	data, err = json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filename, 0o644); err != nil {
		t.Fatal(err)
	}
	writeBundleFixture(t, filename, string(data))
	marker, transport := filepath.Join(dir, "ssh-started"), filepath.Join(dir, "ssh")
	writeBundleFixture(t, transport, "#!/bin/sh\nprintf started > \"$RELEASE_INSPECTION_STARTED\"\nexit 1\n")
	if err := os.Chmod(transport, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := releaseBundleCommand("inspect", "--release_dir", bundle, "--environment", "fixture", "--ssh_host", "fixture.invalid", "--ssh_path", transport)
	cmd.Env = append(os.Environ(), "RELEASE_INSPECTION_STARTED="+marker)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()
	process := filepath.Join("/proc", strconv.Itoa(cmd.Process.Pid))
	started := false
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		counters, err := os.ReadFile(filepath.Join(process, "io"))
		if err != nil {
			t.Fatal(err)
		}
		var readBytes int64
		if _, err := fmt.Sscanf(string(counters), "rchar: %d", &readBytes); err != nil {
			t.Fatal(err)
		}
		if readBytes > 1024*1024 {
			entries, err := os.ReadDir(filepath.Join(process, "fd"))
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				target, _ := os.Readlink(filepath.Join(process, "fd", entry.Name()))
				started = started || target == local
			}
		}
		if started {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !started {
		t.Fatal("inspection did not begin hashing the sparse payload")
	}
	signalled := time.Now()
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() { finished <- cmd.Wait() }()
	select {
	case err := <-finished:
		if err == nil || !strings.Contains(stderr.String(), "context canceled") {
			t.Fatalf("cancelled verification did not report cancellation: %v\n%s\n%s", err, &stdout, &stderr)
		}
	case <-time.After(2 * time.Second):
		_ = cmd.Process.Kill()
		<-finished
		t.Fatal("cancelled local verification did not stop promptly")
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("cancelled verification started SSH: %v", err)
	}
	if strings.Contains(stdout.String(), `"corrupt"`) {
		t.Fatalf("cancellation was mislabeled corruption: %s", &stdout)
	}
	if destination := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); destination != "" {
		evidence := fmt.Sprintf("{\"local_hash_started\":true,\"payload_size\":%d,\"signal\":\"SIGTERM\",\"ssh_started\":false,\"exit_nonzero\":true,\"elapsed_ms\":%d}\n", size, time.Since(signalled).Milliseconds())
		writeBundleFixture(t, filepath.Join(destination, "inspection-cancellation.json"), evidence)
	}
}

// observationTree snapshots entry names, bytes, permissions, modification times,
// and symlink targets so a nominally read-only operation cannot alter the fixture.
func observationTree(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	if err := filepath.WalkDir(root, func(filename string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk observation tree: %w", walkErr)
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect observation entry: %w", err)
		}
		value := fmt.Sprintf("%v:%d", info.Mode(), info.ModTime().UnixNano())
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("read observation entry: %w", err)
			}
			value += fmt.Sprintf(":%x", sha256.Sum256(data))
		} else if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(filename)
			if err != nil {
				return fmt.Errorf("read observation link: %w", err)
			}
			value += ":" + target
		}
		result[filename] = value
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return result
}
