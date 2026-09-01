package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func writeFixture(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for path, content := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func fixtureRegistry() string {
	return `{
		"schema": "agents.alwaldend.com/phase1-registry/v1alpha1",
		"criteriaRevision": 3,
		"operationFiles": []
	}`
}

func validProjectREADME(name string) string {
	return "---\ntitle:  " + name + "\nstatuses:\n  - active\n---\n# " + name + "\n\n" +
		"Maintained " + name + " component.\n"
}

func TestTopologyCompileComplete(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"projects/README.md":                      "# Projects\n",
		"infra/README.md":                         "# Infra\n",
		"tools/README.md":                         "# Tools\n",
		"data/README.md":                          "# Data\n",
		"third_party/README.md":                   "# Third party\n",
		"users/README.md":                         "# Users\n",
		"projects/agents/README.md":               validProjectREADME("agents"),
		"projects/agents/BUILD.bazel":             "",
		"MODULE.bazel":                            `module(name = "com_alwaldend_src")`,
		"tools/agents/declarations/registry.json": fixtureRegistry(),
	})
	var stdout, stderr bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/topology.json",
		"--markdown", "out/topology.md",
	}, &stdout, &stderr); err != nil {
		t.Fatalf("run: %v; stderr=%s", err, stderr.String())
	}
	content, err := os.ReadFile(filepath.Join(root, "out/topology.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := catalogv1alpha1.DecodeTopologyStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if catalog.Completeness != catalogv1alpha1.CompletenessComplete {
		t.Fatalf("expected complete, got %s: %v", catalog.Completeness, catalog.Limitations)
	}
	if len(catalog.Components) != 1 || catalog.Components[0].ID != "agents" {
		t.Fatalf("unexpected components: %#v", catalog.Components)
	}
	if len(catalog.Workspaces) != 1 || catalog.Workspaces[0].ModuleName != "com_alwaldend_src" {
		t.Fatalf("unexpected workspaces: %#v", catalog.Workspaces)
	}
	if !catalog.HasBoundaryClass(catalogv1alpha1.TreeBoundaryProduct) {
		t.Fatal("missing product boundary")
	}
}

func TestTopologyCompileFailsOnMissingProjectBuild(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"projects/README.md":                      "# Projects\n",
		"infra/README.md":                         "# Infra\n",
		"tools/README.md":                         "# Tools\n",
		"data/README.md":                          "# Data\n",
		"third_party/README.md":                   "# Third party\n",
		"users/README.md":                         "# Users\n",
		"projects/agents/README.md":               validProjectREADME("agents"),
		"tools/agents/declarations/registry.json": fixtureRegistry(),
	})
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/topology.json",
		"--markdown", "out/topology.md",
		"--check",
	}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected --check to fail on missing BUILD")
	}
	if !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTopologyCompileNonDeterministicInputRejected(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"projects/README.md":                      "# Projects\n",
		"projects/agents/README.md":               validProjectREADME("agents"),
		"projects/agents/BUILD.bazel":             "",
		"MODULE.bazel":                            `module(name = "com_alwaldend_src")`,
		"tools/agents/declarations/registry.json": fixtureRegistry(),
	})
	var first, second bytes.Buffer
	args := func(output string) []string {
		return []string{
			"--workspace-root", root,
			"--source-revision", "0123456789abcdef0123456789abcdef01234567",
			"--output", output + "/topology.json",
			"--markdown", output + "/topology.md",
		}
	}
	if err := run(args("out1"), &first, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if err := run(args("out2"), &second, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	firstContent, _ := os.ReadFile(filepath.Join(root, "out1/topology.json"))
	secondContent, _ := os.ReadFile(filepath.Join(root, "out2/topology.json"))
	if string(firstContent) != string(secondContent) {
		t.Fatalf("outputs differ:\n%s\nvs\n%s", firstContent, secondContent)
	}
}

func TestTopologyCheckFailsOnStaleJSON(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"projects/README.md":                      "# Projects\n",
		"infra/README.md":                         "# Infra\n",
		"tools/README.md":                         "# Tools\n",
		"data/README.md":                          "# Data\n",
		"third_party/README.md":                   "# Third party\n",
		"users/README.md":                         "# Users\n",
		"projects/agents/README.md":               validProjectREADME("agents"),
		"projects/agents/BUILD.bazel":             "",
		"MODULE.bazel":                            `module(name = "com_alwaldend_src")`,
		"tools/agents/declarations/registry.json": fixtureRegistry(),
		"out/topology.json":                       "stale",
	})
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/topology.json",
		"--markdown", "out/topology.md",
		"--check",
	}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected --check to fail on stale JSON")
	}
	if !strings.Contains(err.Error(), "stale") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestTopologyCheckFailsOnStaleMarkdown ensures --check compares the tracked
// Markdown projection too, instead of silently repairing a stale render while
// reporting success.
func TestTopologyCheckFailsOnStaleMarkdown(t *testing.T) {
	root := writeFixture(t, map[string]string{
		"projects/README.md":                      "# Projects\n",
		"infra/README.md":                         "# Infra\n",
		"tools/README.md":                         "# Tools\n",
		"data/README.md":                          "# Data\n",
		"third_party/README.md":                   "# Third party\n",
		"users/README.md":                         "# Users\n",
		"projects/agents/README.md":               validProjectREADME("agents"),
		"projects/agents/BUILD.bazel":             "",
		"MODULE.bazel":                            `module(name = "com_alwaldend_src")`,
		"tools/agents/declarations/registry.json": fixtureRegistry(),
	})
	args := []string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/topology.json",
		"--markdown", "out/topology.md",
	}
	// Generate current artifacts once.
	if err := run(args, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	// Stale the Markdown only; JSON stays current.
	if err := os.WriteFile(
		filepath.Join(root, "out/topology.md"),
		[]byte("# stale render\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	err := run(append(args, "--check"), &stdout, &stderr)
	if err == nil {
		t.Fatal("expected --check to fail on stale Markdown")
	}
	if !strings.Contains(err.Error(), "Markdown is stale") {
		t.Fatalf("unexpected error: %v", err)
	}
}
