package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

func TestRunPlanEmitsValidJSON(t *testing.T) {
	var stdout bytes.Buffer
	if err := run([]string{
		"plan",
		"--workspace-root", "testdata/root",
		"--profile", "changed/fast",
		"--intent", "cli-example",
		"--path", "tools/agents",
	}, &stdout); err != nil {
		t.Fatalf("run plan: %v", err)
	}
	var plan v1alpha1.ImpactPlan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	decoded, err := v1alpha1.DecodeImpactPlan(stdout.Bytes())
	if err != nil {
		t.Fatalf("plan does not validate: %v", err)
	}
	if decoded.Kind != "ImpactPlan" {
		t.Fatalf("kind = %s, want ImpactPlan", decoded.Kind)
	}
	if decoded.Profile != v1alpha1.ImpactProfileChangedFast {
		t.Fatalf("profile = %s, want changed/fast", decoded.Profile)
	}
	if decoded.IntentRef.ID != "cli-example" {
		t.Fatalf("intent = %s, want cli-example", decoded.IntentRef.ID)
	}
	if !strings.HasPrefix(decoded.Digest, "sha256:") {
		t.Fatalf("digest = %s, want sha256", decoded.Digest)
	}
}

func TestRunPlanRequiresIntentAndPath(t *testing.T) {
	for _, args := range [][]string{
		{"plan", "--workspace-root", "testdata/root", "--path", "tools/agents"},
		{"plan", "--workspace-root", "testdata/root", "--intent", "x"},
	} {
		if err := run(args, &bytes.Buffer{}); err == nil {
			t.Fatalf("run %v: expected error", args)
		}
	}
}
