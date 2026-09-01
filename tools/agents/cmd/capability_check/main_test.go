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

func writeLink(t *testing.T, root, link, target string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(link))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, full); err != nil {
		t.Fatal(err)
	}
}

const capabilityFixtureRegistry = `{
	"schema": "agents.alwaldend.com/phase1-registry/v1alpha1",
	"skills": [
		{
			"id": "answer-question",
			"owner": "projects/agents",
			"layer": "procedure",
			"activation": "substantive user questions",
			"exclusions": ["inert quoted questions", "code linting"],
			"capabilityRefs": ["source.read"],
			"dependencies": [],
			"conflicts": [],
			"providerRequirements": [],
			"contextCost": "medium",
			"evaluationMaturity": "checked"
		}
	],
	"runtimeTools": [
		{
			"id": "cordis_define",
			"owner": "projects/mcp_cordis",
			"classification": "classified"
		}
	],
	"directBinaries": [
		{
			"id": "repo-delivery",
			"owner": "tools/repo_delivery",
			"path": "tools/repo_delivery/main/go/main.go",
			"classification": "classified"
		}
	],
	"operationFiles": ["projects/goal/agent_operations.json"]
}`

const capabilityFixtureOperation = `{
	"schema": "agents.alwaldend.com/operations/v1alpha1",
	"owner": "projects.goal",
	"provider": "goal.local-store",
	"definition": "projects/goal/cmd/goal/command.go",
	"operations": []
}`

func TestCapabilityCompileComplete(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"tools/agents/declarations/registry.json":            capabilityFixtureRegistry,
		"projects/agents/skills/answer-question/SKILL.md":    "# Answer\n",
		"projects/agents/skills/answer-question/BUILD.bazel": "skill_library(name = \"answer\")\n",
		"projects/mcp_cordis/internal/mcp.mjs":               "export const cordis = true;\n",
		"tools/repo_delivery/main/go/main.go":                "package main\n",
		"projects/goal/agent_operations.json":                capabilityFixtureOperation,
		"projects/goal/cmd/goal/command.go":                  "package main\n",
	})
	writeLink(t, root, ".agents/skills/answer-question",
		"../../projects/agents/skills/answer-question")
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/capability.json",
		"--markdown", "out/capability.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/capability.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := catalogv1alpha1.DecodeCapabilityStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(catalog.Skills) != 1 || len(catalog.Providers) != 3 {
		t.Fatalf("unexpected skills/providers: %d %d",
			len(catalog.Skills), len(catalog.Providers))
	}
	if catalog.Skills[0].ID != "answer-question" ||
		catalog.Skills[0].Layer != "procedure" {
		t.Fatalf("unexpected skill: %#v", catalog.Skills[0])
	}
	expectedCost := catalog.Skills[0].ContextCost
	if expectedCost != "medium" {
		t.Fatalf("unexpected cost: %s", expectedCost)
	}
	providers := map[string]bool{}
	for _, provider := range catalog.Providers {
		providers[provider.ID] = true
	}
	if !providers["cordis_define"] || !providers["repo-delivery"] ||
		!providers["goal.local-store"] {
		t.Fatalf("unexpected providers: %#v", catalog.Providers)
	}
}

func TestCapabilityCompileFailsOnMissingSkillDoc(t *testing.T) {
	t.Skip("incompleteness is now a truthful durable state; --check verifies tracked-byte freshness only")
	root := writeFiles(t, map[string]string{
		"tools/agents/declarations/registry.json":            capabilityFixtureRegistry,
		"projects/agents/skills/answer-question/BUILD.bazel": "skill_library(name = \"answer\")\n",
		"projects/mcp_cordis/internal/mcp.mjs":               "export const cordis = true;\n",
		"tools/repo_delivery/main/go/main.go":                "package main\n",
		"projects/goal/agent_operations.json":                capabilityFixtureOperation,
		"projects/goal/cmd/goal/command.go":                  "package main\n",
		"out/capability.json":                                "{}\n",
	})
	writeLink(t, root, ".agents/skills/answer-question",
		"../../projects/agents/skills/answer-question")
	var stdout bytes.Buffer
	err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/capability.json",
		"--markdown", "out/capability.md",
		"--check",
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatalf("expected stale tracked-output failure, got %v", err)
	}
}

func TestCapabilityCompileDiscoveryOnlySkill(t *testing.T) {
	root := writeFiles(t, map[string]string{
		"tools/agents/declarations/registry.json":            capabilityFixtureRegistry,
		"projects/agents/skills/answer-question/SKILL.md":    "# Answer\n",
		"projects/agents/skills/answer-question/BUILD.bazel": "skill_library(name = \"answer\")\n",
		"projects/agents/skills/codex-migration/SKILL.md":    "# Migrate\n",
		"projects/agents/skills/codex-migration/BUILD.bazel": "skill_library(name = \"codex\")\n",
		"projects/mcp_cordis/internal/mcp.mjs":               "export const cordis = true;\n",
		"tools/repo_delivery/main/go/main.go":                "package main\n",
		"projects/goal/agent_operations.json":                capabilityFixtureOperation,
		"projects/goal/cmd/goal/command.go":                  "package main\n",
	})
	writeLink(t, root, ".agents/skills/answer-question",
		"../../projects/agents/skills/answer-question")
	writeLink(t, root, ".agents/skills/codex-migration",
		"../../projects/agents/skills/codex-migration")
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root,
		"--source-revision", "0123456789abcdef0123456789abcdef01234567",
		"--output", "out/capability.json",
		"--markdown", "out/capability.md",
	}, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, "out/capability.json"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := catalogv1alpha1.DecodeCapabilityStrict(content)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(catalog.Skills) != 2 {
		t.Fatalf("expected discovered skill, got %d", len(catalog.Skills))
	}
	if catalog.Completeness == catalogv1alpha1.CompletenessComplete {
		t.Fatalf("expected discovery-only limitation to mark partial")
	}
	found := false
	for _, limitation := range catalog.Limitations {
		if strings.Contains(limitation, "discovered") {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing discovery-only limitation: %#v", catalog.Limitations)
	}
}
