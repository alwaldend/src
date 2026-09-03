package agent_system

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

func TestRunEmitsCompleteJSON(t *testing.T) {
	var stdout bytes.Buffer
	if err := Run([]string{
		"--workspace-root", "testdata/root",
		"--path", "projects/agents",
		"--revision", "abcdef0123456789",
		"--json",
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	var capsule v1alpha1.ContextCapsule
	if err := json.Unmarshal(stdout.Bytes(), &capsule); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	decoded, err := v1alpha1.DecodeContextCapsule(stdout.Bytes())
	if err != nil {
		t.Fatalf("capsule does not validate: %v", err)
	}
	if decoded.Provenance.Completeness != v1alpha1.CompletenessComplete {
		t.Fatalf("completeness = %s, want complete", decoded.Provenance.Completeness)
	}
	if !strings.Contains(decoded.Identity.SourceDigest, "abcdef0123456789") {
		t.Fatalf("source digest not bound: %s", decoded.Identity.SourceDigest)
	}
	if decoded.Component.Workspace == "" {
		t.Fatalf("applicable workspace is empty")
	}
	if len(decoded.Provenance.NextActions) == 0 {
		t.Fatalf("expected safe next discovery actions")
	}
}

func TestRunEmitsStructuredUnavailableOnMissingCatalogs(t *testing.T) {
	var stdout bytes.Buffer
	if err := Run([]string{
		"--workspace-root", "testdata/missing",
		"--path", ".",
		"--revision", "abcdef0123456789",
		"--json",
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	capsule, err := v1alpha1.DecodeContextCapsule(stdout.Bytes())
	if err != nil {
		t.Fatalf("capsule does not validate: %v", err)
	}
	if capsule.Provenance.Completeness != v1alpha1.CompletenessPartial {
		t.Fatalf("completeness = %s, want partial with structured unavailable", capsule.Provenance.Completeness)
	}
	if len(capsule.Provenance.Limitations) == 0 {
		t.Fatalf("expected structured unavailable limitations")
	}
}

func TestMarkdownRenderUsesSameData(t *testing.T) {
	var jsonOut, markdownOut bytes.Buffer
	if err := Run([]string{
		"--workspace-root", "testdata/root",
		"--path", "projects/agents",
		"--revision", "abcdef0123456789",
		"--json",
	}, &jsonOut); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{
		"--workspace-root", "testdata/root",
		"--path", "projects/agents",
		"--revision", "abcdef0123456789",
		"--markdown",
		"--json=false",
	}, &markdownOut); err != nil {
		t.Fatal(err)
	}
	var capsule v1alpha1.ContextCapsule
	if err := json.Unmarshal(jsonOut.Bytes(), &capsule); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(markdownOut.String(), capsule.ID) {
		t.Fatalf("markdown does not state the same capsule ID")
	}
}
