// Package inventory_test exercises the inventory command against isolated owner trees.
package inventory_test

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// repositoryMarker selects an optional real-checkout freshness check.
var repositoryMarker = flag.String("repository-marker", "", "Repository module runfile")

// TestCheckoutFreshness verifies generated projections against the actual source checkout.
func TestCheckoutFreshness(t *testing.T) {
	if *repositoryMarker == "" {
		t.Skip("fixture-only invocation")
	}
	marker, err := runfiles.Rlocation(*repositoryMarker)
	if err != nil {
		t.Fatal(err)
	}
	marker, err = filepath.EvalSymlinks(marker)
	if err != nil {
		t.Fatal(err)
	}
	binary, err := runfiles.Rlocation("_main/tools/repo_quality/cmd/inventory/inventory_/inventory")
	if err != nil {
		t.Fatal(err)
	}
	output, err := exec.Command(binary, "--root", filepath.Dir(marker), "--check").CombinedOutput()
	if err != nil {
		t.Fatalf("inventory freshness: %v\n%s", err, output)
	}
}

// TestOwnerDiscovery verifies regeneration, stale-output rejection and nested module labels.
func TestOwnerDiscovery(t *testing.T) {
	binary, err := runfiles.Rlocation("_main/tools/repo_quality/cmd/inventory/inventory_/inventory")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	files := map[string]string{
		"projects/example/README.md":              "---\ntitle: Example\ndescription: Example description\nstatuses: [active]\n---\n",
		"projects/statuses.yaml":                  "active: Ongoing development.\nmaintenance: Existing scope maintained.\ncompleted: Intended scope complete.\narchived: No further work.\n",
		"projects/example/site/content/_index.md": "---\ntitle: Example\ndate: 2026-01-01\nstatuses: [finished]\ntags: [example]\n---\n\nVisitor content.\n",
		"projects/example/openspec/BUILD.bazel":   "filegroup(name = \"source\", srcs = [\"config.yaml\"])\n",
		"tools/module/MODULE.bazel":               "module(name = \"example_rules\")\n",
		"tools/module/openspec/BUILD.bazel":       "filegroup(name = \"source\", srcs = [\"config.yaml\"])\n",
		"tools/draft/openspec/config.yaml":        "schema: spec-driven\n",
		"projects/README.md":                      "# Projects\n\n<!-- BEGIN GENERATED PROJECT STATUSES -->\nold\n<!-- END GENERATED PROJECT STATUSES -->\n\n<!-- BEGIN GENERATED PROJECT INDEX -->\nold\n<!-- END GENERATED PROJECT INDEX -->\n",
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
	run := func(check bool) error {
		args := []string{"--root", root}
		if check {
			args = append(args, "--check")
		}
		command := exec.Command(binary, args...)
		output, err := command.CombinedOutput()
		t.Log(string(output))
		return err
	}
	if run(true) == nil {
		t.Fatal("accepted missing generated inventory")
	}
	if err := run(false); err != nil {
		t.Fatal(err)
	}
	if err := run(true); err != nil {
		t.Fatal(err)
	}
	landingPath := filepath.Join(root, "projects/example/site/content/_index.md")
	landing, err := os.ReadFile(landingPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, preserved := range []string{"date: 2026-01-01", "tags: [example]", "Visitor content."} {
		if !strings.Contains(string(landing), preserved) {
			t.Fatalf("lost landing content %q", preserved)
		}
	}
	if !strings.Contains(string(landing), "statuses:\n  - active\n") {
		t.Fatalf("landing status not derived from README: %s", landing)
	}
	if err := run(false); err != nil {
		t.Fatal(err)
	}
	repeated, err := os.ReadFile(landingPath)
	if err != nil || string(landing) != string(repeated) {
		t.Fatalf("generation is not repeatable: %v", err)
	}
	readmePath := filepath.Join(root, "projects/example/README.md")
	for _, invalid := range []string{"", "statuses: []\n", "statuses: [active, maintenance]\n", "statuses: [finished]\n", "statuses: active\n"} {
		if err := os.WriteFile(readmePath, []byte("---\ntitle: Example\n"+invalid+"---\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if run(true) == nil || run(false) == nil {
			t.Fatalf("accepted invalid lifecycle declaration %q", invalid)
		}
	}
	if err := os.WriteFile(readmePath, []byte("---\ntitle: Example\nstatuses: [maintenance]\n---\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if run(true) == nil {
		t.Fatal("accepted stale landing status")
	}
	unchanged, err := os.ReadFile(landingPath)
	if err != nil || string(unchanged) != string(landing) {
		t.Fatalf("freshness check changed landing: %v", err)
	}
	if err := run(false); err != nil {
		t.Fatal(err)
	}
	updated, err := os.ReadFile(landingPath)
	if err != nil || !strings.Contains(string(updated), "statuses:\n  - maintenance\n") {
		t.Fatalf("status change was not projected: %v\n%s", err, updated)
	}
	path := filepath.Join(root, "infra/src/openspec/validation/workspaces.bzl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, label := range []string{"//projects/example/openspec:source", "@example_rules//openspec:source"} {
		if !strings.Contains(string(data), label) {
			t.Errorf("missing %s in %s", label, data)
		}
	}
	if strings.Contains(string(data), "draft") {
		t.Fatal("draft is not a declared source export")
	}
	if err := os.WriteFile(filepath.Join(root, "projects/example/README.md"), []byte("---\ntitle: Renamed\ndescription: New description\nstatuses: [active]\n---\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if run(true) == nil {
		t.Fatal("accepted stale project index")
	}
	if err := run(false); err != nil {
		t.Fatal(err)
	}
	if err := run(true); err != nil {
		t.Fatal(err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		for _, name := range []string{"projects/README.md", "projects/example/site/content/_index.md"} {
			artifact, err := os.ReadFile(filepath.Join(root, name))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(output, filepath.Base(filepath.Dir(name))+"-"+filepath.Base(name)), artifact, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(output, "workspaces.bzl"), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
