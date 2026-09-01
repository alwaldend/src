package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func writeActionFiles(t *testing.T, files map[string]string) string {
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

const validOperationFile = `{
	"schema": "agents.alwaldend.com/operations/v1alpha1",
	"owner": "projects.goal",
	"provider": "goal.local-store",
	"definition": "projects/goal/cmd/goal/command.go",
	"operations": [
		{
			"id": "goal.checkpoint",
			"selector": "checkpoint",
			"classification": "classified",
			"effects": ["source.write", "task_state.write"],
			"inputs": ["goal revision"],
			"outputs": ["versioned goal state"],
			"information": ["public"],
			"credentialUse": "none",
			"networkUse": "none",
			"environmentSelector": "explicit goal directory",
			"authorityGate": "goal writer",
			"preflight": "lock and validate",
			"verification": "record validation",
			"cost": "bounded",
			"cacheability": "not_cacheable",
			"cancellation": "atomic boundaries"
		}
	]
}`

func actionFixtureRegistry() string {
	return `{
		"schema": "agents.alwaldend.com/phase1-registry/v1alpha1",
		"criteriaRevision": 3,
		"operationFiles": ["projects/goal/agent_operations.json"]
	}`
}

func TestActionCompileComplete(t *testing.T) {
	root := writeActionFiles(t, map[string]string{
		"tools/agents/declarations/registry.json": actionFixtureRegistry(),
		"projects/goal/agent_operations.json":     validOperationFile,
		"projects/goal/cmd/goal/command.go":       "package main\n",
	})
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/action.json",
		"--markdown", "out/action.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/action.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := catalogv1alpha1.DecodeActionStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if catalog.Completeness != catalogv1alpha1.CompletenessComplete {
		t.Fatalf("expected complete, got %s: %v", catalog.Completeness, catalog.Limitations)
	}
	if len(catalog.Providers) != 1 || len(catalog.Actions) != 1 {
		t.Fatalf("unexpected providers/actions: %#v %#v", catalog.Providers, catalog.Actions)
	}
	if catalog.Actions[0].ID != "goal.checkpoint" ||
		catalog.Actions[0].ProviderRef != "goal.local-store" {
		t.Fatalf("unexpected action: %#v", catalog.Actions[0])
	}
}

func TestActionCompileFailsOnMissingOperationFile(t *testing.T) {
	root := writeActionFiles(t, map[string]string{
		"tools/agents/declarations/registry.json": actionFixtureRegistry(),
	})
	var stdout bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/action.json",
		"--markdown", "out/action.md",
		"--check",
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("expected check failure, got %v", err)
	}
}

func TestActionCompileRejectsUnknownEffect(t *testing.T) {
	bad := strings.Replace(validOperationFile,
		`"effects": ["source.write", "task_state.write"]`,
		`"effects": ["arbitrary.write"]`, 1)
	root := writeActionFiles(t, map[string]string{
		"tools/agents/declarations/registry.json": actionFixtureRegistry(),
		"projects/goal/agent_operations.json":     bad,
		"projects/goal/cmd/goal/command.go":       "package main\n",
	})
	var stdout bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/action.json",
		"--markdown", "out/action.md",
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "unknown effect") {
		t.Fatalf("expected unknown effect rejection, got %v", err)
	}
}
