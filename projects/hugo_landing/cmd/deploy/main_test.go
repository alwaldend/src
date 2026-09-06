package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func testGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	data, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, data)
	}
	return strings.TrimSpace(string(data))
}

func writeTestFile(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDeployBootstrapUnchangedAndReplace(t *testing.T) {
	root := t.TempDir()
	// Prevent the developer's Git settings from affecting test repositories.
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	remote := filepath.Join(root, "remote.git")
	testGit(t, root, "init", "--bare", remote)
	site := filepath.Join(root, "site")
	writeTestFile(t, filepath.Join(site, "index.html"), "first landing")
	writeTestFile(t, filepath.Join(site, "assets", "old.css"), "old")
	opts := options{scratch: filepath.Join(root, "scratch"), git: "git"}
	checkDeploy := func(want string) {
		t.Helper()
		got, err := deploy(opts, "sample_project", site, remote)
		status, revision, _ := strings.Cut(got, " ")
		if want == "would publish" {
			status = got
		} else if err == nil && revision != testGit(t, remote, "rev-parse", "refs/heads/pages") {
			t.Fatalf("receipt revision %q differs from published branch", revision)
		}
		if err != nil || status != want {
			t.Fatalf("deploy = %q, %v; want %q", got, err, want)
		}
		entries, err := os.ReadDir(opts.scratch)
		if err != nil || len(entries) != 0 {
			t.Fatalf("scratch not cleaned: %v, %v", entries, err)
		}
	}
	opts.dryRun = true
	checkDeploy("would publish")
	if got := testGit(t, remote, "for-each-ref", "--format=%(refname)"); got != "" {
		t.Fatalf("dry run created refs: %s", got)
	}
	opts.dryRun = false
	checkDeploy("published")
	first := testGit(t, remote, "rev-parse", "refs/heads/pages")
	// GitHub can make the first pushed branch the repository default.
	// Subsequent clones then already contain a local pages branch.
	testGit(t, remote, "symbolic-ref", "HEAD", "refs/heads/pages")
	if got := testGit(t, remote, "show", "pages:CNAME"); got != "sample-project.alwaldend.com" {
		t.Fatalf("CNAME = %q", got)
	}
	testGit(t, remote, "show", "pages:.nojekyll")
	checkDeploy("unchanged")
	if got := testGit(t, remote, "rev-parse", "refs/heads/pages"); got != first {
		t.Fatal("unchanged deployment created a commit")
	}
	if err := os.Remove(filepath.Join(site, "assets", "old.css")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(site, "assets", "new.css"), "new")
	writeTestFile(t, filepath.Join(site, "index.html"), "second landing")
	opts.dryRun = true
	checkDeploy("would publish")
	if got := testGit(t, remote, "rev-parse", "refs/heads/pages"); got != first {
		t.Fatal("dry run updated remote")
	}
	opts.dryRun = false
	checkDeploy("published")
	if got := testGit(t, remote, "rev-parse", "pages^"); got != first {
		t.Fatal("deployment did not preserve branch history")
	}
	if got := testGit(t, remote, "show", "pages:index.html"); got != "second landing" {
		t.Fatalf("index = %q", got)
	}
	files := testGit(t, remote, "ls-tree", "-r", "--name-only", "pages")
	if strings.Contains(files, "old.css") || !strings.Contains(files, "assets/new.css") {
		t.Fatalf("unexpected published files: %s", files)
	}
}

func TestValidateRejectsUnsafeSite(t *testing.T) {
	site := t.TempDir()
	writeTestFile(t, filepath.Join(site, "index.html"), "landing")
	if err := validateSite("../project", site); err == nil {
		t.Fatal("accepted path traversal in project name")
	}
	if err := os.Symlink("index.html", filepath.Join(site, "linked")); err != nil {
		t.Fatal(err)
	}
	if err := validateSite("project", site); err == nil {
		t.Fatal("accepted symbolic link")
	}
	if err := os.Remove(filepath.Join(site, "linked")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(site, ".git", "config"), "unsafe")
	if err := validateSite("project", site); err == nil {
		t.Fatal("accepted Git metadata")
	}
}

func TestNormalizeRunfilesAndProjectSelection(t *testing.T) {
	root := t.TempDir()
	site := filepath.Join(root, "artifact")
	writeTestFile(t, filepath.Join(site, "index.html"), "landing")
	link := filepath.Join(root, "runfiles")
	if err := os.Symlink(site, link); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	relative, err := filepath.Rel(cwd, link)
	if err != nil {
		t.Fatal(err)
	}
	all := sites{"project": relative, "other": filepath.Join(root, "missing")}
	resolved, err := normalizeSites(all, "project")
	if err != nil || len(resolved) != 1 || resolved["project"] != site {
		t.Fatalf("normalize sites = %v, %v", resolved, err)
	}
	if _, err := normalizeSites(all, "unknown"); err == nil {
		t.Fatal("accepted unknown project selector")
	}
	if err := os.Symlink("index.html", filepath.Join(site, "interior")); err != nil {
		t.Fatal(err)
	}
	if _, err := normalizeSites(all, "project"); err == nil {
		t.Fatal("root normalization allowed an interior symlink")
	}
}

func TestDeployPreservesExistingDefaultBranch(t *testing.T) {
	root := t.TempDir()
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	remote := filepath.Join(root, "remote.git")
	testGit(t, root, "init", "--bare", remote)
	seed := filepath.Join(root, "seed")
	testGit(t, root, "init", "--initial-branch=main", seed)
	writeTestFile(t, filepath.Join(seed, "README.md"), "repository metadata")
	testGit(t, seed, "add", ".")
	testGit(t, seed, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-m", "Initialize")
	testGit(t, seed, "push", remote, "main")
	testGit(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")
	before := testGit(t, remote, "rev-parse", "main")
	site := filepath.Join(root, "site")
	writeTestFile(t, filepath.Join(site, "index.html"), "landing")
	opts := options{scratch: filepath.Join(root, "scratch"), git: "git"}
	if result, err := deploy(opts, "project", site, remote); err != nil || !strings.HasPrefix(result, "published ") {
		t.Fatalf("deploy = %q, %v", result, err)
	}
	if got := testGit(t, remote, "rev-parse", "main"); got != before {
		t.Fatal("changed existing default branch")
	}
	if files := testGit(t, remote, "ls-tree", "-r", "--name-only", "pages"); strings.Contains(files, "README.md") {
		t.Fatal("copied default branch metadata into pages")
	}
}

func TestDeployFailureCleansOnlyOwnScratch(t *testing.T) {
	root := t.TempDir()
	site := filepath.Join(root, "site")
	writeTestFile(t, filepath.Join(site, "index.html"), "landing")
	opts := options{scratch: filepath.Join(root, "scratch"), git: "git"}
	sentinel := filepath.Join(opts.scratch, "unrelated")
	writeTestFile(t, sentinel, "keep")
	if _, err := deploy(opts, "project", site, filepath.Join(root, "missing.git")); err == nil {
		t.Fatal("accepted inaccessible repository")
	}
	entries, err := os.ReadDir(opts.scratch)
	if err != nil || len(entries) != 1 || entries[0].Name() != "unrelated" {
		t.Fatalf("incorrect scratch cleanup: %v, %v", entries, err)
	}
}
