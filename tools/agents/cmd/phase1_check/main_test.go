package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTerraformSelectorsIncludesEveryNamedEntry(t *testing.T) {
	content := []byte(`DEFAULT_TERRAFORM_BINARIES = {
    "apply": ["apply"],
    "plan": ["plan"],
}
`)
	got := terraformSelectors(content)
	if !got["apply"] || !got["plan"] || got[""] || len(got) != 2 {
		t.Fatalf("terraformSelectors() = %#v", got)
	}
}

func TestValidateOperationRejectsIncompleteOrUnknownClassification(t *testing.T) {
	valid := operation{
		ID: "terraform.plan", Selector: "plan", Classification: "classified",
		Effects: []string{"network.read"}, Information: []string{"public"},
		EnvironmentSelector: "selected-workspace",
		AuthorityGate:       "read", Preflight: "init", Verification: "exit-status",
		Cost: "bounded", Cacheability: "not_cacheable", Cancellation: "signal",
	}
	if err := validateOperation(valid); err != nil {
		t.Fatalf("valid operation rejected: %v", err)
	}
	valid.Classification = "deprecated"
	if err := validateOperation(valid); err != nil {
		t.Fatalf("deprecated compatibility operation rejected: %v", err)
	}
	valid.Classification = "unknown"
	if err := validateOperation(valid); err == nil {
		t.Fatal("unknown classification accepted")
	}
	valid.Classification = "classified"
	valid.Verification = ""
	if err := validateOperation(valid); err == nil {
		t.Fatal("incomplete operation accepted")
	}
}

func TestFrontmatterStatuses(t *testing.T) {
	content := []byte("---\ntitle: Example\nstatuses:\n  - active\ntags:\n  - x\n---\n")
	got := frontmatterStatuses(content)
	if len(got) != 1 || got[0] != "active" {
		t.Fatalf("frontmatterStatuses() = %#v", got)
	}
}

func TestTaskRegistryIdentifiersRejectTraversal(t *testing.T) {
	for _, value := range []string{"../escape", "bad id", "UPPER"} {
		if validateID(value) {
			t.Fatalf("validateID(%q) = true", value)
		}
	}
}

func TestOpenSpecAuthorityAndChangeInventory(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		"MODULE.bazel",
		"openspec/changes/current/.openspec.yaml",
		"openspec/changes/archive/2026-09-08-finished/.openspec.yaml",
		"openspec/changes/archive/2026-09-08-finished/provenance/source/goal.yaml",
		"openspec/changes/archive/2026-09-08-finished/provenance/openspec/changes/past/.openspec.yaml",
		"openspec/changes/archive/.openspec.yaml",
		"infra/src/openspec/changes/current/.openspec.yaml",
		"infra/src/openspec/changes/archive/2026-09-08-finished/.openspec.yaml",
		"projects/sample/openspec/changes/current/.openspec.yaml",
		"projects/sample/openspec/changes/archive/2026-09-08-finished/.openspec.yaml",
		"infra/sample/openspec/changes/current/.openspec.yaml",
		"infra/sample/openspec/changes/archive/2026-09-08-finished/.openspec.yaml",
		"projects/sample/openspec/changes/current/provenance/.openspec.yaml",
		"projects/sample/testdata/openspec/changes/fixture/.openspec.yaml",
		"infra/sample/testdata/openspec/changes/fixture/.openspec.yaml",
		"tools/sample/testdata/openspec/changes/fixture/.openspec.yaml",
		"out/task/openspec/changes/scratch/.openspec.yaml",
	} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	value := checker{
		root: root,
		registry: registry{Authorities: []authority{
			{ID: "repository.openspec", Kind: "openspec", Source: "infra/src/openspec"},
		}},
		report: report{Counts: map[string]int{}},
	}
	value.checkAuthorities()
	if err := value.checkWorkspacesAndChanges(); err != nil {
		t.Fatal(err)
	}
	if len(value.report.Unclassified) != 0 || len(value.report.Missing) != 0 ||
		value.report.Counts["authorities"] != 1 || value.report.Counts["workspaces"] != 1 ||
		value.report.Counts["openspecChanges"] != 6 {
		t.Fatalf("incorrect OpenSpec inventory: %+v", value.report)
	}
	if _, exists := value.report.Counts["goals"]; exists {
		t.Fatal("legacy goal state advertised as current inventory")
	}
}
