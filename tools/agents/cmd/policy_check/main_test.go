package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func writeFiles(t *testing.T, files map[string]string) string {
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

func fixturePolicyRoot() map[string]string {
	return map[string]string{
		"AGENTS.md": "# Host Codex instructions\n\n" +
			"## Repository visibility\n\npublic source is fair game.\n" +
			"## Infrastructure safety\n\nsecrets never in reports.\n",
		"CODEOWNERS":            "* @alwaldend\n",
		"projects/README.md":    "# Projects\n",
		"infra/README.md":       "# Infra\n",
		"tools/README.md":       "# Tools\n",
		"data/README.md":        "# Data\n",
		"third_party/README.md": "# Third party\n",
		"users/README.md":       "# Users\n",
	}
}

func TestPolicyCompileComplete(t *testing.T) {
	root := writeFiles(t, fixturePolicyRoot())
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/policy.json",
		"--markdown", "out/policy.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/policy.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := catalogv1alpha1.DecodePolicyStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if catalog.Completeness != catalogv1alpha1.CompletenessComplete {
		t.Fatalf("expected complete, got %s: %v", catalog.Completeness, catalog.Limitations)
	}
	if len(catalog.Policies) != 1 || !strings.HasPrefix(catalog.Policies[0].ID, "policy.") {
		t.Fatalf("unexpected policies: %#v", catalog.Policies)
	}
	axisFound := false
	for _, axis := range catalog.Policies[0].Axes {
		if axis.Name == "sourceDisclosure" && axis.Value == "known" {
			axisFound = true
		}
	}
	if !axisFound {
		t.Fatalf("missing sourceDisclosure axis: %#v", catalog.Policies[0].Axes)
	}
}

func TestPolicyCompileFailsOnMissingBoundaryReadme(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"AGENTS.md":          "# Host Codex instructions\n",
		"CODEOWNERS":         "* @alwaldend\n",
		"projects/README.md": "# Projects\n",
	})
	var stdout bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/policy.json",
		"--markdown", "out/policy.md",
		"--check",
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("expected check failure, got %v", err)
	}
}
