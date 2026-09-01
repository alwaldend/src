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

func fixtureCatalogFile(kind string) (string, string) {
	content := `{
	"schema": "agents.alwaldend.com/catalog/v1alpha1/` + kind + `",
	"kind": "` + kind + `",
	"id": "agent-system.` + strings.ReplaceAll(kind, "-catalog", "") + `",
	"derivationVersion": "1.0.0",
	"producerRef": "repository.test",
	"sourceRevision": "0123456789abcdef0123456789abcdef01234567",
	"inputs": [
		{"path": "input/one.yaml", "role": "source", "digest": "sha256:` + strings.Repeat("a", 64) + `"},
		{"path": "input/one.yaml", "role": "source", "digest": "sha256:` + strings.Repeat("a", 64) + `"},
		{"path": "input/two.yaml", "role": "source", "digest": "sha256:` + strings.Repeat("b", 64) + `"}
	],
	"bounds": {
		"eligible": 1, "emitted": 1, "unavailable": 0,
		"maxItems": 1000, "maxInputBytes": 33554432, "maxOutputBytes": 1048576
	},
	"completeness": "complete",
	"limitations": [],
	"conflicts": [],
	"digest": "sha256:` + strings.Repeat("c", 64) + `"
}`
	return "agent-system." + strings.ReplaceAll(kind, "-catalog", ""), content
}

var fixtureKinds = []string{
	"topology-catalog",
	"policy-catalog",
	"action-catalog",
	"capability-catalog",
	"workspace-check-catalog",
	"goal-catalog",
}

func fixtureIndexFiles() map[string]string {
	files := map[string]string{}
	paths := map[string]string{
		"topology-catalog":        "tools/agents/catalogs/topology.json",
		"policy-catalog":          "tools/agents/catalogs/policy.json",
		"action-catalog":          "tools/agents/catalogs/action.json",
		"capability-catalog":      "tools/agents/catalogs/capability.json",
		"workspace-check-catalog": "tools/agents/catalogs/workspace-check.json",
		"goal-catalog":            "tools/agents/catalogs/goal.json",
	}
	for _, kind := range fixtureKinds {
		_, content := fixtureCatalogFile(kind)
		path := paths[kind]
		files[path] = content
	}
	return files
}

func TestIndexCompileComplete(t *testing.T) {
	root := writeFiles(t, fixtureIndexFiles())
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/index.json",
		"--markdown", "out/index.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/index.json"))
	if err != nil {
		t.Fatal(err)
	}
	index, err := catalogv1alpha1.DecodeIndexStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if index.Completeness != catalogv1alpha1.CompletenessComplete {
		t.Fatalf("expected complete, got %s: %v",
			index.Completeness, index.Limitations)
	}
	if len(index.Catalogs) != 6 {
		t.Fatalf("unexpected catalogs: %#v", index.Catalogs)
	}
	if index.Catalogs[0].ID != "agent-system.action" {
		t.Fatalf("unexpected first catalog: %#v", index.Catalogs[0])
	}
	descriptor := index.Catalogs[len(index.Catalogs)-1]
	if len(descriptor.InputDigests) != 2 {
		t.Fatalf("expected deduped input digests, got %#v", descriptor.InputDigests)
	}
}

func TestIndexCompileMissingCatalogs(t *testing.T) {
	_, topologyContent := fixtureCatalogFile("topology-catalog")
	_, policyContent := fixtureCatalogFile("policy-catalog")
	root := writeFiles(t, map[string]string{
		"tools/agents/catalogs/topology.json": topologyContent,
		"tools/agents/catalogs/policy.json":   policyContent,
	})
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/index.json",
		"--markdown", "out/index.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/index.json"))
	if err != nil {
		t.Fatal(err)
	}
	index, err := catalogv1alpha1.DecodeIndexStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if index.Completeness != catalogv1alpha1.CompletenessComplete {
		t.Fatalf("expected complete index, got %s: %v",
			index.Completeness, index.Limitations)
	}
	if len(index.Catalogs) != 6 {
		t.Fatalf("unexpected catalogs: %#v", index.Catalogs)
	}
	for _, descriptor := range index.Catalogs {
		if descriptor.Completeness != "unavailable" &&
			descriptor.ID != "agent-system.topology" &&
			descriptor.ID != "agent-system.policy" {
			t.Fatalf("unexpected descriptor state: %#v", descriptor)
		}
	}
}

func TestIndexCheckFailsOnMissingTrackedFile(t *testing.T) {
	root := writeFiles(t, map[string]string{})
	var stdout bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/index.json",
		"--markdown", "out/index.md",
		"--check",
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("expected check failure, got %v", err)
	}
}
