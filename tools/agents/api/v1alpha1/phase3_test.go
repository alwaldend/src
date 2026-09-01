package v1alpha1

import (
	"strings"
	"testing"
)

func validReference(kind ReferenceKind, id string) Reference {
	return Reference{Kind: kind, ID: id, Version: "v1"}
}

func TestImpactPlanCanonicalRoundTrip(t *testing.T) {
	plan := ImpactPlan{
		APIVersion: APIVersion,
		Kind:       "ImpactPlan",
		ID:         "plan/example",
		IntentRef:  validReference(ReferenceTask, "task/example"),
		Profile:    ImpactProfileChangedFast,
		Capabilities: []Reference{
			validReference(ReferenceProvider, "provider/bazel"),
		},
		CapabilityReasons: []string{"fast changed-scope validation"},
		Effects:           []Effect{EffectCodeExecute, EffectSourceRead},
		ForbiddenEffects:  []Effect{EffectRemoteWrite},
		Targets: []PlanTarget{
			{Label: "pkg/example:all", Path: "pkg/example", AffectedBy: "source change"},
		},
		Checks: []PlanCheck{
			{Identifier: "test/pkg/example/all", Scope: "package", Covers: "changed package"},
		},
		CoverageGaps: []string{"no external dependency validation"},
		CostClass:    CostClassLow,
		ExpectedMaxCost: Budget{
			Calls:       4,
			Bytes:       1 << 20,
			DurationMS:  30000,
			Concurrency: 1,
		},
		Cacheability: Cacheable,
		Ordering:     []string{"test", "build"},
		Concurrency:  1,
		Escalation:   "prompt for remote write",
	}
	first, err := CanonicalImpactPlanJSON(plan)
	if err != nil {
		t.Fatalf("CanonicalImpactPlanJSON() error = %v", err)
	}
	decoded, err := DecodeImpactPlan(first)
	if err != nil {
		t.Fatalf("DecodeImpactPlan() error = %v", err)
	}
	second, err := CanonicalImpactPlanJSON(decoded)
	if err != nil {
		t.Fatalf("CanonicalImpactPlanJSON(decoded) error = %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("canonical round trip changed bytes:\n%s\n%s", first, second)
	}
}

func TestImpactPlanRejectsMissingCapability(t *testing.T) {
	plan := ImpactPlan{
		APIVersion: APIVersion,
		Kind:       "ImpactPlan",
		ID:         "plan/example",
		IntentRef:  validReference(ReferenceTask, "task/example"),
		Profile:    ImpactProfileWorkspace,
		Effects:    []Effect{EffectSourceRead},
		Targets:    []PlanTarget{{Label: "pkg/example:all"}},
		Checks:     []PlanCheck{{Identifier: "test/pkg/example/all"}},
		CostClass:  CostClassMedium,
	}
	if _, err := CanonicalImpactPlanJSON(plan); err == nil ||
		!strings.Contains(err.Error(), "at least one capability") {
		t.Fatalf("CanonicalImpactPlanJSON() error = %v, want capability rejection", err)
	}
}

func TestImpactPlanRejectsForbiddenAndRequiredOverlap(t *testing.T) {
	plan := ImpactPlan{
		APIVersion: APIVersion,
		Kind:       "ImpactPlan",
		ID:         "plan/example",
		IntentRef:  validReference(ReferenceTask, "task/example"),
		Profile:    ImpactProfileFullAudit,
		Capabilities: []Reference{
			validReference(ReferenceProvider, "provider/bazel"),
		},
		Effects:          []Effect{EffectSourceRead},
		ForbiddenEffects: []Effect{EffectSourceRead},
		Targets:          []PlanTarget{{Label: "pkg/example:all"}},
		Checks:           []PlanCheck{{Identifier: "test/pkg/example/all"}},
		CostClass:        CostClassHigh,
	}
	_, err := CanonicalImpactPlanJSON(plan)
	if err == nil || !strings.Contains(err.Error(), "both required and forbidden") {
		t.Fatalf("CanonicalImpactPlanJSON() error = %v, want overlap rejection", err)
	}
}

func TestValidationSetCanonicalRoundTrip(t *testing.T) {
	set := ValidationSet{
		APIVersion: APIVersion,
		Kind:       "ValidationSet",
		ID:         "validation/example",
		Profile:    ImpactProfileChangedFast,
		Candidate:  validReference(ReferenceSubject, "subject/commit"),
		Inputs: ValidationInputs{
			CommitOID: strings.Repeat("ab", 20),
			ProviderRefs: []Reference{
				validReference(ReferenceProvider, "provider/bazel"),
			},
			ConfigDigests:    map[string]string{"config": "sha256:" + strings.Repeat("cd", 32)},
			ToolchainDigests: map[string]string{"toolchain": "sha256:" + strings.Repeat("ef", 32)},
			PolicyDigests:    map[string]string{"policy": "sha256:" + strings.Repeat("01", 32)},
		},
		SanitizedArgs: []string{"test", "//pkg/example:all"},
		WorkingScope:  "pkg/example",
		Results: []ValidationResult{
			{CheckID: "test/pkg/example/all", Status: "pass", DurationMS: 1200},
		},
		Coverage:        []string{"pkg/example"},
		TotalDurationMS: 1200,
		OutputBounds:    Budget{Calls: 4, Bytes: 1 << 20, DurationMS: 30000, Concurrency: 1},
		CleanPreState:   true,
		CleanPostState:  true,
	}
	first, err := CanonicalValidationSetJSON(set)
	if err != nil {
		t.Fatalf("CanonicalValidationSetJSON() error = %v", err)
	}
	decoded, err := DecodeValidationSet(first)
	if err != nil {
		t.Fatalf("DecodeValidationSet() error = %v", err)
	}
	second, err := CanonicalValidationSetJSON(decoded)
	if err != nil {
		t.Fatalf("CanonicalValidationSetJSON(decoded) error = %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("canonical round trip changed bytes")
	}
	digest := DigestOfCanonicalJSON(first)
	if !strings.HasPrefix(digest, "sha256:") || len(digest) != 71 {
		t.Fatalf("DigestOfCanonicalJSON() = %q, want sha256 digest", digest)
	}
}

func TestValidationSetRejectsCredentialLikeArg(t *testing.T) {
	set := ValidationSet{
		APIVersion: APIVersion,
		Kind:       "ValidationSet",
		ID:         "validation/example",
		Profile:    ImpactProfileWorkspace,
		Candidate:  validReference(ReferenceSubject, "subject/commit"),
		Inputs:     ValidationInputs{},
		SanitizedArgs: []string{
			"https://user:secret@example.com/repo",
		},
		WorkingScope: "pkg/example",
		Results: []ValidationResult{
			{CheckID: "test/pkg/example/all", Status: "pass"},
		},
		OutputBounds:   Budget{},
		CleanPreState:  true,
		CleanPostState: true,
	}
	if _, err := CanonicalValidationSetJSON(set); err == nil ||
		!strings.Contains(err.Error(), "credential") {
		t.Fatalf("CanonicalValidationSetJSON() error = %v, want credential rejection", err)
	}
}

func TestEvidenceAssertionCanonicalRoundTrip(t *testing.T) {
	assertion := EvidenceAssertion{
		APIVersion:   APIVersion,
		Kind:         "EvidenceAssertion",
		ID:           "assertion/example",
		CriterionRef: validReference(ReferenceGoal, "goal/agent-system-phase-3"),
		CriterionRev: "1",
		ValidationRefs: []Reference{
			validReference(ReferenceArtifact, "validation/example"),
		},
		Verdict: "satisfied",
	}
	first, err := CanonicalEvidenceAssertionJSON(assertion)
	if err != nil {
		t.Fatalf("CanonicalEvidenceAssertionJSON() error = %v", err)
	}
	decoded, err := DecodeEvidenceAssertion(first)
	if err != nil {
		t.Fatalf("DecodeEvidenceAssertion() error = %v", err)
	}
	second, err := CanonicalEvidenceAssertionJSON(decoded)
	if err != nil {
		t.Fatalf("CanonicalEvidenceAssertionJSON(decoded) error = %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("canonical round trip changed bytes")
	}
}

func TestAdmissionRequestRequiresPrepareForRemoteWrite(t *testing.T) {
	request := AdmissionRequest{
		APIVersion: APIVersion,
		Kind:       "AdmissionRequest",
		Operation:  validReference(ReferenceOperation, "operation/publish"),
		Authority: AuthorityEnvelope{
			ActorRef:    validReference(ReferenceActor, "actor/agent"),
			SubjectRefs: []Reference{validReference(ReferenceSubject, "subject/commit")},
			Effects:     []Effect{EffectRemoteWrite},
			Budget:      Budget{Calls: 1, Bytes: 1, DurationMS: 1, Concurrency: 1},
		},
		RemoteWrite: true,
		Destroy:     false,
	}
	if _, err := CanonicalAdmissionJSON(request); err == nil ||
		!strings.Contains(err.Error(), "prepare verification") {
		t.Fatalf("CanonicalAdmissionJSON() error = %v, want prepare verification rejection", err)
	}
	request.PrepareVerified = true
	if _, err := CanonicalAdmissionJSON(request); err != nil {
		t.Fatalf("CanonicalAdmissionJSON() with prepare verified error = %v", err)
	}
}
