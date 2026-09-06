package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var testGit = flag.String("git", "", "Git executable runfile")

func TestGenerate(t *testing.T) {
	files := []string{
		"MODULE.bazel",
		"projects/z/include.MODULE.bazel",
		"third_party/z/include.MODULE.bazel",
		"tools/z/include.MODULE.bazel",
		"tools/a/include.MODULE.bazel",
		"tools/include.MODULE.bazel",
		"tools/a/include.MODULE.bazel", // staged and untracked discovery must not duplicate an include
		"projects/nested/MODULE.bazel",
		"projects/nested/include.MODULE.bazel",
		"projects/nested/child/include.MODULE.bazel",
		"projects/legacy/WORKSPACE.bazel",
		"projects/legacy/include.MODULE.bazel",
		"projects/old/WORKSPACE",
		"projects/old/include.MODULE.bazel",
		"projects/nested_neighbor/include.MODULE.bazel",
		"ignored/include.MODULE.bazel",
	}
	want := header + `
include("//third_party/z:include.MODULE.bazel")
include("//tools:include.MODULE.bazel")
include("//tools/a:include.MODULE.bazel")
include("//tools/z:include.MODULE.bazel")
include("//projects/nested_neighbor:include.MODULE.bazel")
include("//projects/z:include.MODULE.bazel")
`
	if got := string(generate(files, "# ignored directory\nignored/\n")); got != want {
		t.Fatalf("generated module:\n%s\nwant:\n%s", got, want)
	}
	// Input enumeration order must not change the checked-in module.
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}
	if got := string(generate(files, "ignored")); got != want {
		t.Fatal("discovery order changed the module")
	}
}

func TestIncludeAdditionAndRemoval(t *testing.T) {
	before := string(generate(nil, ""))
	after := string(generate([]string{"tools/new/include.MODULE.bazel"}, ""))
	if before == after || !strings.Contains(after, `include("//tools/new:include.MODULE.bazel")`) {
		t.Fatal("new include was not discovered")
	}
	if got := string(generate(nil, "")); strings.Contains(got, "//tools/new:") {
		t.Fatal("removed include was retained")
	}
}

func TestUpdateDiscoveryAndFreshness(t *testing.T) {
	root := t.TempDir()
	git, err := runfiles.Rlocation(*testGit)
	if err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command(git, "--git-dir", filepath.Join(root, ".git"), "init", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
	write := func(name, content string) {
		t.Helper()
		filename := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("MODULE.bazel", "")
	write(".gitignore", "out/\n")
	write("out/include.MODULE.bazel", "")
	write("tools/new/include.MODULE.bazel", "")
	if err := update(root, *testGit, false); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "MODULE.bazel"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "//tools/new:") || strings.Contains(string(got), "//out:") {
		t.Fatalf("incorrect discovery: %s", got)
	}
	if err := update(root, *testGit, true); err != nil {
		t.Fatal(err)
	}
	write("tools/added/include.MODULE.bazel", "")
	if err := update(root, *testGit, true); err == nil {
		t.Fatal("freshness accepted a missing include")
	}
	if err := update(root, *testGit, false); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "tools/added/include.MODULE.bazel")); err != nil {
		t.Fatal(err)
	}
	if err := update(root, *testGit, true); err == nil {
		t.Fatal("freshness accepted a stale include")
	}
}
