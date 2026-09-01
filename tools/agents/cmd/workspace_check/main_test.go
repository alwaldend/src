package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func writeWorkspaceFiles(t *testing.T, files map[string]string) string {
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

const workspaceFixtureModule = `module(
    name = "com_example_src",
    version = "0.0.0",
)
`

const workspaceFixtureNestedModule = `module(
    name = "rules_docs",
    version = "0.0.0",
)
`

const workspaceFixtureBazelIgnore = `projects/rules_docs
out
node_modules
`

const workspaceFixtureFullCheckScript = `package main

var repositoryWorkspaces = []workspace{
	{name: "root", path: "."},
	{name: "projects/rules_docs", path: "projects/rules_docs"},
}
`

func TestWorkspaceCheckCompileComplete(t *testing.T) {
	root := writeWorkspaceFiles(t, map[string]string{
		"MODULE.bazel":                         workspaceFixtureModule,
		".bazelignore":                         workspaceFixtureBazelIgnore,
		"projects/rules_docs/MODULE.bazel":     workspaceFixtureNestedModule,
		"projects/rules_docs/docs/BUILD.bazel": "# docs\n",
		"projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go": workspaceFixtureFullCheckScript,
	})
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/workspace-check.json",
		"--markdown", "out/workspace-check.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/workspace-check.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := catalogv1alpha1.DecodeWorkspaceCheckStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if catalog.Completeness != catalogv1alpha1.CompletenessComplete {
		t.Fatalf("expected complete, got %s: %v", catalog.Completeness, catalog.Limitations)
	}
	if len(catalog.Workspaces) != 2 {
		t.Fatalf("unexpected workspaces: %#v", catalog.Workspaces)
	}
	byID := map[string]catalogv1alpha1.WorkspaceRecord{}
	for _, workspace := range catalog.Workspaces {
		byID[workspace.ID] = workspace
	}
	rootRecord, ok := byID["root"]
	if !ok {
		t.Fatalf("missing root workspace: %#v", byID)
	}
	if rootRecord.Path != "." || rootRecord.ModulePath != "MODULE.bazel" ||
		rootRecord.ModuleName != "com_example_src" ||
		len(rootRecord.Phases) != 1 ||
		rootRecord.Phases[0].ID != "root.check" ||
		rootRecord.Phases[0].ProviderRef != "repository.bazel-operations" {
		t.Fatalf("unexpected root record: %#v", rootRecord)
	}
	if !rootRecord.Projections.FullCheck || rootRecord.Projections.BazelIgnore {
		t.Fatalf("unexpected root projections: %#v", rootRecord.Projections)
	}
	nested, ok := byID["projects.rules_docs"]
	if !ok {
		t.Fatalf("missing nested workspace: %#v", byID)
	}
	if nested.ModuleName != "rules_docs" || !nested.Projections.BazelIgnore ||
		!nested.Projections.RootOverride || !nested.Projections.DocsAggregation ||
		!nested.Projections.FullCheck {
		t.Fatalf("unexpected nested projections: %#v", nested)
	}
}

// TestWorkspaceCheckCompileCheckFailure verifies that --check fails when the
// compiler reports a completeness problem (here: a MODULE.bazel that cannot
// be parsed for its module name). The tracked JSON must exist so the check
// reaches the completeness gate.
func TestWorkspaceCheckCompileCheckFailure(t *testing.T) {
	root := writeWorkspaceFiles(t, map[string]string{
		"MODULE.bazel":             workspaceFixtureModule,
		".bazelignore":             workspaceFixtureBazelIgnore,
		"out/workspace-check.json": "{}\n",
		"projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go": workspaceFixtureFullCheckScript,
		"projects/broken/MODULE.bazel":                                          "no module name here\n",
	})
	var stdout bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/workspace-check.json",
		"--markdown", "out/workspace-check.md",
		"--check",
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("expected check failure, got %v", err)
	}
}
