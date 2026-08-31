package main

import (
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
