package gazelle

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// gazelleBinary is the runfiles path of the plugin-only executable.
var gazelleBinary = flag.String("gazelle-binary", "", "Gazelle executable runfiles path")

// TestGazelleWorkspace exercises generation and repeated merging through the plugin-only binary.
func TestGazelleWorkspace(t *testing.T) {
	binary, err := runfiles.Rlocation(*gazelleBinary)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	files := map[string]string{
		"MODULE.bazel":                          "module(name = \"example\")\nbazel_dep(name = \"rules_openspec\", version = \"0.0.0\", repo_name = \"spec_rules\")\n",
		"owner/openspec/config.yaml":            "schema: spec-driven\n",
		"owner/openspec/changes/nested/BUILD":   "filegroup(name = \"nested\")\n",
		"owner/openspec/changes/nested/spec.md": "# Nested\n",
		"manual/openspec/config.yaml":           "schema: spec-driven\n",
		"manual/openspec/BUILD":                 "load(\"@spec_rules//rules:defs.bzl\", \"openspec_workspace\", \"openspec_validate_test\")\nopenspec_workspace(name = \"source\", srcs = [\"config.yaml\"], visibility = [\"//visibility:public\"])\nopenspec_validate_test(name = \"validate_test\", source = \":custom\")\nfilegroup(name = \"custom\")\n",
		"collision/openspec/config.yaml":        "schema: spec-driven\n",
		"collision/openspec/BUILD.bazel":        "filegroup(name = \"source\")\n",
		"unconfigured/openspec/README.md":       "# Unconfigured\n",
		"other/config.yaml":                     "schema: spec-driven\n",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	run := func(args ...string) {
		t.Helper()
		command := exec.Command(binary, append([]string{"fix", "-repo_root", root}, args...)...)
		command.Dir = root
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("run Gazelle %v: %v\n%s", args, err, output)
		}
	}
	read := func(name string) []byte {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	const generated = "owner/openspec/BUILD.bazel"
	const manual = "manual/openspec/BUILD"

	run()
	initial := read(generated)
	for _, want := range []string{
		"@spec_rules//rules:defs.bzl",
		"openspec_workspace(",
		"name = \"source\"",
		"config = \"config.yaml\"",
		"\"specs/**/*.md\"",
		"\"changes/**\"",
		"allow_empty = True",
		"openspec_validate_test(",
		"name = \"validate_test\"",
		"source = \":source\"",
	} {
		if !bytes.Contains(initial, []byte(want)) {
			t.Errorf("generated BUILD lacks %q:\n%s", want, initial)
		}
	}
	if strings.Contains(string(initial), "changes/nested") {
		t.Errorf("generated BUILD explicitly refers to nested package:\n%s", initial)
	}
	custom := read(manual)
	for _, want := range []string{
		"srcs = [\"config.yaml\"]",
		"visibility = [\"//visibility:public\"]",
		"source = \":custom\"",
		"filegroup(name = \"custom\")",
	} {
		if !bytes.Contains(custom, []byte(want)) {
			t.Errorf("manual BUILD lost %q:\n%s", want, custom)
		}
	}
	for _, name := range []string{"unconfigured/openspec", "other"} {
		if _, err := os.Stat(filepath.Join(root, name, "BUILD.bazel")); !os.IsNotExist(err) {
			t.Errorf("unexpected package at %q: %v", name, err)
		}
	}
	if bytes.Contains(read("collision/openspec/BUILD.bazel"), []byte("openspec_validate_test")) {
		t.Error("generated a validation test for a colliding source target")
	}
	run("-mode=diff")
	if !bytes.Equal(read(generated), initial) || !bytes.Equal(read(manual), custom) {
		t.Error("Gazelle regeneration changed a BUILD file")
	}
	if err := os.Remove(filepath.Join(root, "owner/openspec/config.yaml")); err != nil {
		t.Fatal(err)
	}
	run()
	if !bytes.Equal(read(generated), initial) {
		t.Error("Gazelle removed a declaration after the configuration disappeared")
	}
	if err := os.WriteFile(filepath.Join(root, "MODULE.bazel"), []byte("module(name = \"example\")\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	fallback := filepath.Join(root, "fallback", "openspec")
	if err := os.MkdirAll(fallback, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fallback, "config.yaml"), []byte("schema: spec-driven\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("fallback/openspec")
	if output := read("fallback/openspec/BUILD.bazel"); !bytes.Contains(output, []byte("@rules_openspec//rules:defs.bzl")) {
		t.Errorf("generated BUILD lacks canonical load:\n%s", output)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "generated.BUILD.bazel"), initial, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
