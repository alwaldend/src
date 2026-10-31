package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// TestDownloadSyncPlaybook exercises the packaged playbook through real loopback
// OpenSSH and rsync; all content and credentials belong to the isolated fixture.
func TestDownloadSyncPlaybook(t *testing.T) {
	if os.Getenv("DOWNLOAD_SYNC_MOUNT_NAMESPACE") != "1" {
		command := exec.Command("unshare", "--user", "--map-current-user", "--keep-caps", "--mount", "--propagation", "private", os.Args[0], "-test.run=^TestDownloadSyncPlaybook$", "-test.v")
		command.Env = append(os.Environ(), "DOWNLOAD_SYNC_MOUNT_NAMESPACE=1")
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("isolated mount-namespace acceptance: %v\n%s", err, output)
		}
		return
	}
	for _, tool := range []string{"ssh", "sshd", "ssh-keygen", "rsync", "prlimit", "mount", "umount", "setpriv"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Fatalf("acceptance fixture requires installed %s: %v", tool, err)
		}
	}
	dir := t.TempDir()
	login, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"host", "client", "ca"} {
		syncFixtureCommand(t, "ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", filepath.Join(dir, name))
	}
	syncFixtureCommand(t, "ssh-keygen", "-q", "-s", filepath.Join(dir, "ca"), "-I", "sync-fixture", "-n", login.Username, "-V", "-1m:+10m", filepath.Join(dir, "client.pub"))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, "sshd_config")
	syncFixtureFile(t, config, fmt.Sprintf("Port %d\nListenAddress 127.0.0.1\nHostKey %s/host\nPidFile %s/pid\nAuthorizedKeysFile none\nTrustedUserCAKeys %s/ca.pub\nPasswordAuthentication no\nKbdInteractiveAuthentication no\nUsePAM no\nStrictModes no\nUseDNS no\nLogLevel ERROR\n", port, dir, dir, dir))
	sshdPath, err := exec.LookPath("sshd")
	if err != nil {
		t.Fatal(err)
	}
	log, err := os.Create(filepath.Join(dir, "sshd.log"))
	if err != nil {
		t.Fatal(err)
	}
	// Keep mount capability in the fixture controller; SSH publishers must
	// remain unable to override ordinary filesystem permissions.
	server := exec.Command("sh", "-c", "umask 077; exec \"$@\"", "fixture-sshd", "setpriv", "--inh-caps=-all", "--ambient-caps=-all", "--", sshdPath, "-D", "-e", "-f", config)
	server.Stdout, server.Stderr = log, log
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = server.Process.Kill()
		_ = server.Wait()
		_ = log.Close()
	})
	pub, err := os.ReadFile(filepath.Join(dir, "host.pub"))
	if err != nil {
		t.Fatal(err)
	}
	known := filepath.Join(dir, "known_hosts")
	syncFixtureFile(t, known, fmt.Sprintf("[127.0.0.1]:%d %s", port, pub))
	// Supply the shared injector references expected by the current inventory.
	t.Setenv("DOWNLOAD_SSH_PRIVATE_KEY", filepath.Join(dir, "client"))
	t.Setenv("DOWNLOAD_SSH_CERTIFICATE", filepath.Join(dir, "client-cert.pub"))
	t.Setenv("DOWNLOAD_SSH_KNOWN_HOSTS", known)
	ready := false
	for deadline := time.Now().Add(10 * time.Second); time.Now().Before(deadline); {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			ready = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if !ready {
		body, _ := os.ReadFile(filepath.Join(dir, "sshd.log"))
		t.Fatalf("isolated sshd did not start: %s", body)
	}
	source, destination := filepath.Join(dir, "source store"), filepath.Join(dir, "destination store")
	for _, root := range []string{source, destination} {
		for _, name := range []string{"projects", "sites", "staging", "state"} {
			if err := os.MkdirAll(filepath.Join(root, name), 0o755); err != nil {
				t.Fatal(err)
			}
		}
	}
	artifact := "projects/sample/releases/v1/payload with ' spaces.bin"
	syncFixtureFile(t, filepath.Join(source, artifact), "source bytes\n")
	syncFixtureFile(t, filepath.Join(source, "sites/web.example/releases/v1/index.html"), "new website\n")
	syncFixtureFile(t, filepath.Join(source, "sites/web.example/releases/v1/assets/main.css"), "body{}\n")
	syncFixtureFile(t, filepath.Join(source, "sites/web.example/releases/v2/index.html"), "unselected website\n")
	syncFixtureFile(t, filepath.Join(source, "staging/private"), "must not transfer\n")
	syncFixtureFile(t, filepath.Join(source, "state/private"), "must not transfer\n")
	syncFixtureFile(t, filepath.Join(destination, "projects/sample/releases/old/keep.bin"), "retain me\n")
	syncFixtureFile(t, filepath.Join(destination, "state/private"), "destination state stays private\n")
	syncFixtureFile(t, filepath.Join(destination, "sites/web.example/releases/old/index.html"), "old website\n")
	syncFixtureLink(t, filepath.Join(source, "sites/web.example/current"), "releases/v1")
	syncFixtureLink(t, filepath.Join(destination, "sites/web.example/current"), "releases/old")
	// The real host root is administrator-owned; publishers can replace entries
	// only inside projects, sites, and staging, not those top-level directories.
	if err := os.Chmod(destination, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(destination, 0o755) })
	metadata := func(name, root string) string {
		file := filepath.Join(dir, name+".json")
		body, err := json.Marshal(map[string]any{"ssh": map[string]string{"host": "127.0.0.1", "user": login.Username, "root": root, "environment": name}})
		if err != nil {
			t.Fatal(err)
		}
		syncFixtureFile(t, file, string(body))
		return file
	}
	sourceMetadata := metadata("source", source)
	t.Setenv("RELEASE_TEST_IDENTITY", filepath.Join(dir, "client"))
	transport := []string{
		"--ssh_config", "/dev/null", "--ssh_identity", "${RELEASE_TEST_IDENTITY}", "--ssh_certificate", filepath.Join(dir, "client-cert.pub"),
		"--ssh_disable_agent", "--ssh_known_hosts", known, "--ssh_port", strconv.Itoa(port), "--timeout", "30s",
	}
	commonSSH := fmt.Sprintf("-F /dev/null -o BatchMode=yes -o IdentityAgent=none -o IdentitiesOnly=yes -o StrictHostKeyChecking=yes -o GlobalKnownHostsFile=/dev/null -o UserKnownHostsFile=%s -o CertificateFile=%s", known, filepath.Join(dir, "client-cert.pub"))
	inventory := map[string]any{"all": map[string]any{
		"vars": map[string]any{"ansible_host": "127.0.0.1", "ansible_port": port, "ansible_user": login.Username, "ansible_ssh_private_key_file": filepath.Join(dir, "client"), "ansible_ssh_common_args": commonSSH, "ansible_python_interpreter": "/usr/bin/python3", "ansible_pipelining": true, "ansible_ssh_args": "-o ControlMaster=no -o ControlPath=none", "download_sync_require_mount": false, "download_sync_become": false, "download_publisher": ""},
		"children": map[string]any{
			"download_local":  map[string]any{"hosts": map[string]any{"source": map[string]string{"download_root": source}}},
			"download_yandex": map[string]any{"hosts": map[string]any{"destination": map[string]string{"download_root": destination}}},
		},
	}}
	inventoryBytes, err := json.Marshal(inventory)
	if err != nil {
		t.Fatal(err)
	}
	inventoryPath := filepath.Join(dir, "inventory.json")
	syncFixtureFile(t, inventoryPath, string(inventoryBytes))
	runfiles := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"))
	args := []string{"--inventory", inventoryPath, "--extra-vars", "ansible_user=" + login.Username, filepath.Join(runfiles, "playbook_sync.yaml")}
	receipt := []string{}
	run := func(name string, success bool, extra ...string) string {
		t.Helper()
		command := exec.Command(filepath.Join(runfiles, "infra/download/ansible/ansible_bin.script.sh"), append(append([]string{}, args...), extra...)...)
		command.Dir = runfiles
		command.Env = append(os.Environ(), "TMPDIR="+dir, "BUILD_WORKSPACE_DIRECTORY=", "ANSIBLE_LOCAL_TEMP="+filepath.Join(dir, "ansible-local"), "ANSIBLE_CONFIG="+filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"), "ansible.cfg"))
		output, err := command.CombinedOutput()
		if (err == nil) != success {
			t.Fatalf("%s: success=%t, error=%v\n%s", name, success, err, output)
		}
		if strings.Contains(string(output), "unknown command") || strings.Contains(string(output), "unknown flag") {
			t.Fatalf("%s: command contract missing: %s", name, output)
		}
		receipt = append(receipt, name)
		return string(output)
	}
	syncFixturePublication(t, dir, source, sourceMetadata, transport)
	receipt = append(receipt, "local_publication_collision_and_concurrency", "mounted_release_deletion_refusal")
	before := syncFixtureTree(t, destination)
	if output := run("dry_run", true, "--check"); !strings.Contains(output, "*deleting") {
		t.Fatal("dry run did not report public-only content deletion")
	}
	if got := syncFixtureTree(t, destination); got != before {
		t.Fatal("dry run changed destination")
	}
	run("required_mount_missing", false, "--extra-vars", "download_sync_require_mount=true")
	wrongKnown := filepath.Join(dir, "wrong_known_hosts")
	syncFixtureFile(t, wrongKnown, "")
	trustOverride, err := json.Marshal(map[string]string{"ansible_ssh_common_args": strings.ReplaceAll(commonSSH, known, wrongKnown)})
	if err != nil {
		t.Fatal(err)
	}
	run("host_trust_refusal", false, "--extra-vars", string(trustOverride))
	syncFixtureLink(t, filepath.Join(source, "projects/sample/releases/v1/escape"), filepath.Join(source, "state/private"))
	run("source_symlink_refusal", false)
	if err := os.Remove(filepath.Join(source, "projects/sample/releases/v1/escape")); err != nil {
		t.Fatal(err)
	}
	syncFixtureLink(t, filepath.Join(source, "sites/web.example/current"), "../../state")
	run("current_escape_refusal", false)
	syncFixtureLink(t, filepath.Join(source, "sites/web.example/current"), "releases/v1")
	index := filepath.Join(source, "sites/web.example/releases/v1/index.html")
	if err := os.Rename(index, index+".saved"); err != nil {
		t.Fatal(err)
	}
	run("missing_index_refusal", false)
	if err := os.Rename(index+".saved", index); err != nil {
		t.Fatal(err)
	}
	if got := syncFixtureTree(t, destination); got != before {
		t.Fatal("preflight failure changed destination")
	}
	// Published files must never act as rsync policy and hide other content.
	syncFixtureFile(t, filepath.Join(source, "projects/sample/releases/v1/.rsync-filter"), "- *\n")
	run("initial_copy", true)
	if got, err := os.ReadFile(filepath.Join(destination, artifact)); err != nil || string(got) != "source bytes\n" {
		t.Fatalf("artifact copy: %q, %v", got, err)
	}
	current := filepath.Join(destination, "sites/web.example/current")
	if got, err := os.Readlink(current); err != nil || got != "releases/v1" {
		t.Fatalf("selected website: %q, %v", got, err)
	}
	for _, name := range []string{"staging/private"} {
		if _, err := os.Lstat(filepath.Join(destination, name)); !os.IsNotExist(err) {
			t.Fatalf("private source content copied: %s", name)
		}
	}
	if _, err := os.Stat(filepath.Join(destination, "projects/sample/releases/old/keep.bin")); !os.IsNotExist(err) {
		t.Fatalf("public-only artifact retained: %v", err)
	}
	if _, err := os.Stat(filepath.Join(destination, "sites/web.example/releases/v2/index.html")); err != nil {
		t.Fatalf("unselected site missing: %v", err)
	}
	first := syncFixtureTree(t, destination)
	selection, err := os.Lstat(current)
	if err != nil {
		t.Fatal(err)
	}
	run("unchanged_repeat", true)
	if got := syncFixtureTree(t, destination); got != first {
		t.Fatal("unchanged repeat altered content")
	}
	repeated, err := os.Lstat(current)
	if err != nil || !selection.ModTime().Equal(repeated.ModTime()) {
		t.Fatal("unchanged selection replaced current link")
	}
	syncFixtureFile(t, filepath.Join(source, artifact), "updated ordinary bytes\n")
	run("ordinary_file_update", true)
	if got, err := os.ReadFile(filepath.Join(destination, artifact)); err != nil || string(got) != "updated ordinary bytes\n" {
		t.Fatalf("artifact update: %q, %v", got, err)
	}
	syncFixtureFile(t, index, "conflicting immutable website\n")
	run("source_replaces_differing_site", true)
	if got, err := os.ReadFile(filepath.Join(destination, "sites/web.example/releases/v1/index.html")); err != nil || string(got) != "conflicting immutable website\n" {
		t.Fatalf("source website did not win: %q, %v", got, err)
	}
	syncFixtureFile(t, index, "new website\n")
	syncFixtureLink(t, filepath.Join(source, "sites/web.example/current"), "releases/v2")
	publicBeforeFailure := syncFixtureTree(t, destination)
	syncFixtureFile(t, filepath.Join(source, "projects/sample/releases/v1/large.bin"), strings.Repeat("x", 16384))
	// Restrict only destination rsync so writing the staged payload fails.
	// Ansible module execution and controller resources are unchanged.
	run("staged_write_failure_preserves_public_trees", false, "--extra-vars", `{"download_sync_rsync_path":"prlimit --fsize=4096: -- rsync"}`)
	if got := syncFixtureTree(t, destination); got != publicBeforeFailure {
		t.Fatal("failed staged transfer changed published trees")
	}
	if err := os.Remove(filepath.Join(source, "projects/sample/releases/v1/large.bin")); err != nil {
		t.Fatal(err)
	}
	syncFixtureFile(t, filepath.Join(source, artifact), "fail before selection\n")
	if err := os.Chmod(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(destination, "projects"), filepath.Join(destination, "saved-projects")); err != nil {
		t.Fatal(err)
	}
	syncFixtureLink(t, filepath.Join(destination, "projects"), filepath.Join(destination, "saved-projects"))
	run("destination_symlink_refusal", false)
	if err := os.Remove(filepath.Join(destination, "projects")); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(destination, "saved-projects"), filepath.Join(destination, "projects")); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(destination, 0o555); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(source, "sites/web.example/current")); err != nil {
		t.Fatal(err)
	}
	run("remove_missing_source_selection", true)
	if _, err := os.Lstat(current); !os.IsNotExist(err) {
		t.Fatal("missing source selection did not clear public selection")
	}
	run("mirror_retained_site_and_archives", true)
	for _, root := range []string{source, destination} {
		versions, err := os.ReadDir(filepath.Join(root, "sites/publication.example/releases"))
		if err != nil || len(versions) != 2 || versions[0].Name() != "2026.44.0" || versions[1].Name() != "2026.44.2" {
			t.Fatalf("retained extracted versions in %s: %v, %v", root, versions, err)
		}
		for _, version := range []string{"2026.44.0", "2026.44.2"} {
			if _, err := os.Stat(filepath.Join(root, "projects/publication.example/releases", version, "website.zip")); err != nil {
				t.Fatalf("retained archive %s in %s: %v", version, root, err)
			}
		}
	}
	for _, tree := range []string{"projects", "sites"} {
		if err := os.RemoveAll(filepath.Join(source, tree)); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(source, tree), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if got, err := os.ReadFile(filepath.Join(destination, "state/private")); err != nil || string(got) != "destination state stays private\n" {
		t.Fatal("private destination state changed")
	}
	run("empty_authoritative_source", true)
	for _, tree := range []string{"projects", "sites"} {
		entries, err := os.ReadDir(filepath.Join(destination, tree))
		if err != nil || len(entries) != 0 {
			t.Fatalf("public %s not empty: %v", tree, err)
		}
	}
	if err := os.Rename(source, source+".absent"); err != nil {
		t.Fatal(err)
	}
	run("missing_source_root", false)
	residue, err := filepath.Glob(filepath.Join(dir, "capture.*"))
	if err != nil || len(residue) != 0 {
		t.Fatalf("controller staging residue: %v, %v", residue, err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		body, err := json.MarshalIndent(map[string]any{"passed": receipt, "initial_mirror": json.RawMessage(first), "destination_tree": json.RawMessage(syncFixtureTree(t, destination)), "transport": "packaged Ansible playbook over real loopback OpenSSH and rsync with an ephemeral client certificate"}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "sync-receipt.json"), append(body, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// syncFixturePublication verifies local immutability through the real deploy CLI.
func syncFixturePublication(t *testing.T, dir, root, metadata string, transport []string) {
	t.Helper()
	bundle := filepath.Join(dir, "release")
	// Model sudo retaining a login directory that the content user cannot search.
	privateHome := filepath.Join(dir, "private-login-home")
	if err := os.Mkdir(privateHome, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(privateHome, 0o700) })
	quote := func(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'" }
	remotePrefix := "chmod 700 -- " + quote(privateHome) + "; cd -- " + quote(privateHome) + "; chmod 000 -- " + quote(privateHome) + "; "
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		t.Fatal(err)
	}
	sshWrapper := filepath.Join(dir, "ssh-private-home")
	syncFixtureFile(t, sshWrapper, "#!/usr/bin/env bash\nset -eu\nargs=(\"$@\")\nlast=$((${#args[@]} - 1))\nargs[last]="+quote(remotePrefix)+"\"${args[last]}\"\nexec "+quote(sshPath)+" \"${args[@]}\"\n")
	if err := os.Chmod(sshWrapper, 0o700); err != nil {
		t.Fatal(err)
	}
	al, err := runfiles.Rlocation(os.Getenv("AL_BINARY"))
	if err != nil {
		t.Fatal(err)
	}
	writeBundle := func(version, content string) {
		t.Helper()
		var buffer bytes.Buffer
		archive := zip.NewWriter(&buffer)
		index, err := archive.Create("index.html")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := index.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
		if err := archive.Close(); err != nil {
			t.Fatal(err)
		}
		payload, status := filepath.Join(dir, "website.zip"), filepath.Join(dir, "status")
		syncFixtureFile(t, payload, buffer.String())
		syncFixtureFile(t, status, "STABLE_VERSION "+version+"\nSTABLE_VERSION_CHANNEL release\nSTABLE_GIT_COMMIT 0123456789012345678901234567890123456789\nSTABLE_GIT_TREE_STATE clean\n")
		if err := os.RemoveAll(bundle); err != nil {
			t.Fatal(err)
		}
		syncFixtureCommand(t, al, "release", "generate", "--project", "projects/publication.example", "--version_file", status, "--add_file", payload, "--output_dir", bundle)
	}
	command := func(steps string) *exec.Cmd {
		options := []string{"release", "deploy", "--environment", "source", "--ssh_deployment", metadata, "--release_dir", bundle, "--project", "publication.example", "--site_archive", "website.zip"}
		if steps == "delete" {
			options[1] = "delete"
		} else {
			options = append(options, "--steps", steps)
		}
		options = append(options, transport...)
		if steps == "link" || steps == "delete" {
			options = append(options, "--ssh_path", sshWrapper)
		}
		cmd := exec.Command(al, options...)
		cmd.Env = append(os.Environ(), "TMPDIR="+dir, "BUILD_WORKSPACE_DIRECTORY=")
		return cmd
	}
	run := func(steps string, success bool) {
		t.Helper()
		output, err := command(steps).CombinedOutput()
		if (err == nil) != success {
			t.Fatalf("publication %s: success=%t, error=%v\n%s", steps, success, err, output)
		}
	}
	writeBundle("2026.44.0", "first local website\n")
	run("upload,extract,link", true)
	archives := filepath.Join(root, "projects/publication.example")
	archiveDirectory := filepath.Join(archives, "releases/2026.44.0")
	if info, err := os.Stat(archiveDirectory); err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("archive directory permissions: %v, %v", info, err)
	}
	firstArchive, err := os.ReadFile(filepath.Join(archiveDirectory, "website.zip"))
	if err != nil {
		t.Fatal(err)
	}
	before := syncFixtureTree(t, root)
	run("upload,extract,link", true)
	if got := syncFixtureTree(t, root); got != before {
		t.Fatal("identical retry changed local content")
	}
	writeBundle("2026.44.0", "changed duplicate website\n")
	run("upload,extract,link", false)
	run("extract", false)
	if got := syncFixtureTree(t, root); got != before {
		t.Fatal("conflicting publication changed local content")
	}
	writeBundle("2026.44.2", "concurrent website\n")
	first, second := command("upload"), command("upload")
	var firstOutput, secondOutput bytes.Buffer
	first.Stdout, first.Stderr = &firstOutput, &firstOutput
	second.Stdout, second.Stderr = &secondOutput, &secondOutput
	if err := first.Start(); err != nil {
		t.Fatal(err)
	}
	if err := second.Start(); err != nil {
		_ = first.Wait()
		t.Fatal(err)
	}
	firstErr, secondErr := first.Wait(), second.Wait()
	if firstErr != nil || secondErr != nil {
		t.Fatalf("identical concurrent publication: %v, %v\n%s\n%s", firstErr, secondErr, &firstOutput, &secondOutput)
	}
	run("link", false)
	run("extract", true)
	versionsRoot := filepath.Join(root, "sites/publication.example/releases")
	private := filepath.Join(root, "state/retention-sentinel")
	syncFixtureFile(t, private, "private state survives deletion\n")
	mountedVersion := filepath.Join(versionsRoot, "2026.44.2")
	syncFixtureCommand(t, "mount", "--bind", filepath.Join(root, "state"), mountedVersion)
	output, mountErr := command("delete").CombinedOutput()
	syncFixtureCommand(t, "umount", mountedVersion)
	if mountErr == nil || !strings.Contains(strings.ToLower(string(output)), "mount") {
		t.Fatalf("mounted deletion was not refused: %v, %s", mountErr, output)
	}
	if body, err := os.ReadFile(private); err != nil || string(body) != "private state survives deletion\n" {
		t.Fatalf("deletion changed mounted state: %q, %v", body, err)
	}
	if current, err := os.Readlink(filepath.Join(root, "sites/publication.example/current")); err != nil || current != "releases/2026.44.0" {
		t.Fatalf("mount refusal changed selection: %q, %v", current, err)
	}
	for _, mountedFile := range []string{filepath.Join(mountedVersion, "index.html"), filepath.Join(archives, "releases/2026.44.2/website.zip")} {
		syncFixtureCommand(t, "mount", "--bind", private, mountedFile)
		before := syncFixtureTree(t, root)
		output, mountErr := command("delete").CombinedOutput()
		after := syncFixtureTree(t, root)
		syncFixtureCommand(t, "umount", mountedFile)
		if mountErr == nil || !strings.Contains(strings.ToLower(string(output)), "mount") || before != after {
			t.Fatalf("file mount refusal changed release content: %v, %s", mountErr, output)
		}
	}
	allArchives := syncFixtureTree(t, archives)
	run("link", true)
	writeBundle("2026.44.0", "first local website\n")
	run("link", true)
	for _, version := range []string{"2026.44.0", "2026.44.2"} {
		if _, err := os.Stat(filepath.Join(versionsRoot, version, "index.html")); err != nil {
			t.Fatalf("selection removed retained site %s: %v", version, err)
		}
	}
	if err := os.Rename(versionsRoot, versionsRoot+".saved"); err != nil {
		t.Fatal(err)
	}
	syncFixtureLink(t, versionsRoot, versionsRoot+".saved")
	before = syncFixtureTree(t, root)
	run("link", false)
	if got := syncFixtureTree(t, root); got != before {
		t.Fatal("unsafe ancestry refusal changed content")
	}
	if err := os.Remove(versionsRoot); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(versionsRoot+".saved", versionsRoot); err != nil {
		t.Fatal(err)
	}
	if got := syncFixtureTree(t, archives); got != allArchives {
		t.Fatal("rollback changed downloadable archives")
	}
	original, err := os.ReadFile(filepath.Join(archiveDirectory, "website.zip"))
	if err != nil || !bytes.Equal(original, firstArchive) {
		t.Fatalf("original archive changed after rollback: %v", err)
	}
}

// syncFixtureCommand runs one bounded setup operation without exposing keys.
func syncFixtureCommand(t *testing.T, name string, args ...string) {
	t.Helper()
	if output, err := exec.Command(name, args...).CombinedOutput(); err != nil {
		t.Fatalf("fixture %s: %v: %s", name, err, output)
	}
}

// syncFixtureFile creates an ordinary content fixture and its parent directories.
func syncFixtureFile(t *testing.T, name, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

// syncFixtureLink changes only the explicitly named fixture symlink.
func syncFixtureLink(t *testing.T, name, target string) {
	t.Helper()
	if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if err := os.Symlink(target, name); err != nil {
		t.Fatal(err)
	}
}

// syncFixtureTree captures content hashes, permissions, and literal link targets.
func syncFixtureTree(t *testing.T, root string) string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(name string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk fixture %q: %w", name, walkErr)
		}
		rel, err := filepath.Rel(root, name)
		if err != nil {
			return fmt.Errorf("relativize fixture %q: %w", name, err)
		}
		if rel == "staging" {
			return filepath.SkipDir
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect fixture %q: %w", name, err)
		}
		value := info.Mode().String()
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(name)
			if err != nil {
				return fmt.Errorf("read fixture %q: %w", name, err)
			}
			value += fmt.Sprintf(" %x", sha256.Sum256(data))
		} else if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(name)
			if err != nil {
				return fmt.Errorf("read fixture link %q: %w", name, err)
			}
			value += " " + target
		}
		result[rel] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
