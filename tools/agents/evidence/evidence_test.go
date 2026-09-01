package evidence

import (
	"encoding/json"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

func TestEmitValidationSetDeterministicDigest(t *testing.T) {
	opts := EmitOptions{
		Profile:       v1alpha1.ImpactProfileChangedFast,
		Candidate:     repositoryReference("alwaldend/src"),
		BaseOID:       "base-001",
		TreeOID:       "tree-001",
		CommitOID:     "commit-001",
		SanitizedArgs: []string{"--define", "foo=bar"},
		WorkingScope:  "projects/agents",
		ProviderRefs:  []v1alpha1.Reference{providerReference("bazel")},
		ConfigDigests: map[string]string{
			"config": v1alpha1.DigestOfCanonicalJSON([]byte("cfg")),
		},
		OutputBounds: v1alpha1.Budget{Calls: 50},
	}
	results := []v1alpha1.ValidationResult{
		{CheckID: "build", Status: "pass", DurationMS: 10},
	}
	first, err := EmitValidationSet(opts, results)
	if err != nil {
		t.Fatalf("EmitValidationSet() error = %v", err)
	}
	second, err := EmitValidationSet(opts, results)
	if err != nil {
		t.Fatalf("EmitValidationSet() error = %v", err)
	}
	if first.ID != second.ID {
		t.Errorf("IDs differ: %q vs %q", first.ID, second.ID)
	}
	if first.Digest != second.Digest {
		t.Errorf("digests differ: %q vs %q", first.Digest, second.Digest)
	}
	if !strings.HasPrefix(first.ID, "validation.projects.agents.") {
		t.Errorf("ID = %q, want validation.projects.agents prefix", first.ID)
	}
	if !strings.HasPrefix(first.Digest, "sha256:") ||
		len(first.Digest) != len("sha256:")+64 {
		t.Errorf("Digest = %q, want sha256 hex", first.Digest)
	}
	if first.Digest != second.Digest || first.ID != second.ID {
		t.Fatal("deterministic identity mismatch")
	}
	// The digest covers the emitted identity: re-derive the digest from the
	// digest-free payload (with the provisional ID) and check the ID suffix
	// derives from it.
	draft := first
	draft.ID = "validation." + scopeSlug(opts.WorkingScope)
	draft.Digest = ""
	content, err := json.Marshal(draft)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	payloadDigest := v1alpha1.DigestOfCanonicalJSON(content)
	wantSuffix := shortDigest(payloadDigest)
	if !strings.HasSuffix(first.ID, "."+wantSuffix) {
		t.Errorf("ID = %q, want suffix .%s", first.ID, wantSuffix)
	}
	// Re-encoding the emitted set must preserve its digest.
	plain := first
	plain.Digest = ""
	content, err = json.Marshal(plain)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if got := v1alpha1.DigestOfCanonicalJSON(content); got != first.Digest {
		t.Errorf("recomputed digest = %q, want %q", got, first.Digest)
	}
}

func TestEmitValidationSetRejectsCredentialArgs(t *testing.T) {
	opts := EmitOptions{
		Profile:       v1alpha1.ImpactProfileChangedFast,
		Candidate:     repositoryReference("alwaldend/src"),
		WorkingScope:  "projects/agents",
		SanitizedArgs: []string{"--url", "https://user:pass@example.com"},
	}
	_, err := EmitValidationSet(opts, passingResults())
	if err == nil {
		t.Fatal("EmitValidationSet() accepted credential-like argument")
	}
	if !strings.Contains(err.Error(), "resembles a credential") {
		t.Errorf("error = %v, want credential rejection", err)
	}
}

func TestNewAssertionBindsSetDigests(t *testing.T) {
	set := mustEmitSet(t, EmitOptions{
		Profile:      v1alpha1.ImpactProfileChangedFast,
		Candidate:    repositoryReference("alwaldend/src"),
		WorkingScope: "projects/agents",
	})
	assertion, err := NewAssertion(
		"assertion.agents.build",
		goalReference("agents/repo-health"),
		"rev-42",
		"satisfied",
		set,
	)
	if err != nil {
		t.Fatalf("NewAssertion() error = %v", err)
	}
	if assertion.Kind != "EvidenceAssertion" {
		t.Errorf("Kind = %q", assertion.Kind)
	}
	if len(assertion.ValidationRefs) != 1 {
		t.Fatalf(
			"ValidationRefs = %d, want 1",
			len(assertion.ValidationRefs),
		)
	}
	ref := assertion.ValidationRefs[0]
	if ref.ID != set.ID || ref.Digest != set.Digest {
		t.Errorf(
			"validation ref = %+v, want set %q digest %q",
			ref,
			set.ID,
			set.Digest,
		)
	}
	if err := Apply(assertion, set); err != nil {
		t.Errorf("Apply() error = %v", err)
	}
	// The assertion digest is deterministic for the same inputs.
	again, err := NewAssertion(
		"assertion.agents.build",
		goalReference("agents/repo-health"),
		"rev-42",
		"satisfied",
		set,
	)
	if err != nil {
		t.Fatalf("NewAssertion() error = %v", err)
	}
	if again.Digest != assertion.Digest {
		t.Errorf(
			"deterministic assertion digest: %q vs %q",
			again.Digest,
			assertion.Digest,
		)
	}
}

func TestOneSetSupportsMultipleAssertions(t *testing.T) {
	set := mustEmitSet(t, EmitOptions{
		Profile:      v1alpha1.ImpactProfileChangedFast,
		Candidate:    repositoryReference("alwaldend/src"),
		WorkingScope: "projects/agents",
	})
	wantDigest := set.Digest
	first, err := NewAssertion(
		"assertion.agents.build",
		goalReference("agents/repo-health"),
		"rev-1",
		"satisfied",
		set,
	)
	if err != nil {
		t.Fatalf("NewAssertion() error = %v", err)
	}
	second, err := NewAssertion(
		"assertion.agents.safety",
		goalReference("agents/safety"),
		"rev-2",
		"satisfied",
		set,
	)
	if err != nil {
		t.Fatalf("NewAssertion() error = %v", err)
	}
	if err := Apply(first, set); err != nil {
		t.Errorf("Apply(first) error = %v", err)
	}
	if err := Apply(second, set); err != nil {
		t.Errorf("Apply(second) error = %v", err)
	}
	// Immutability: applying an assertion with a forged digest must fail and
	// must not alter the set.
	forged := first
	forged.Digest = forgeDigest(first.Digest)
	tampered := v1alpha1.EvidenceAssertion{
		APIVersion:   forged.APIVersion,
		Kind:         forged.Kind,
		ID:           forged.ID,
		CriterionRef: forged.CriterionRef,
		CriterionRev: forged.CriterionRev,
		ValidationRefs: []v1alpha1.Reference{
			{Kind: v1alpha1.ReferenceArtifact, ID: set.ID, Digest: forged.Digest},
		},
		Verdict: forged.Verdict,
		Digest:  forged.Digest,
	}
	if err := Apply(tampered, set); err == nil {
		t.Error("Apply() accepted forged validation digest")
	}
	if set.Digest != wantDigest {
		t.Errorf(
			"Apply() mutated the immutable set: digest = %q, want %q",
			set.Digest,
			wantDigest,
		)
	}
}

func TestApplyRejectsUnboundValidationRef(t *testing.T) {
	set := mustEmitSet(t, EmitOptions{
		Profile:      v1alpha1.ImpactProfileChangedFast,
		Candidate:    repositoryReference("alwaldend/src"),
		WorkingScope: "projects/agents",
	})
	assertion, err := NewAssertion(
		"assertion.agents.build",
		goalReference("agents/repo-health"),
		"rev-1",
		"satisfied",
		set,
	)
	if err != nil {
		t.Fatalf("NewAssertion() error = %v", err)
	}
	unbound := assertion
	unbound.ValidationRefs = []v1alpha1.Reference{
		{
			Kind:   v1alpha1.ReferenceArtifact,
			ID:     "validation.other-set",
			Digest: v1alpha1.DigestOfCanonicalJSON([]byte("other")),
		},
	}
	if err := Apply(unbound, set); err == nil {
		t.Error("Apply() accepted unbound validation reference")
	} else if !strings.Contains(err.Error(), "not bound") {
		t.Errorf("error = %v, want unbound message", err)
	}
}

func TestApplyAcceptsProviderBoundInputRef(t *testing.T) {
	set := mustEmitSet(t, EmitOptions{
		Profile:      v1alpha1.ImpactProfileChangedFast,
		Candidate:    repositoryReference("alwaldend/src"),
		WorkingScope: "projects/agents",
		ProviderRefs: []v1alpha1.Reference{providerReference("bazel")},
	})
	assertion := v1alpha1.EvidenceAssertion{
		APIVersion:   v1alpha1.APIVersion,
		Kind:         "EvidenceAssertion",
		ID:           "assertion.agents.provider",
		CriterionRef: goalReference("agents/repo-health"),
		CriterionRev: "rev-1",
		ValidationRefs: []v1alpha1.Reference{
			{
				Kind:   v1alpha1.ReferenceArtifact,
				ID:     "bazel",
				Digest: v1alpha1.DigestOfCanonicalJSON([]byte("provider")),
			},
		},
		Verdict: "satisfied",
	}
	// A provider ref without a digest is not bound-by-input with a digest.
	if err := Apply(assertion, set); err == nil {
		t.Error("Apply() accepted provider digests that are not recorded")
	}
	// Record the matching provider digest on the set; then the ref binds.
	recording := set
	recording.Inputs.ProviderRefs = []v1alpha1.Reference{
		{
			Kind:   v1alpha1.ReferenceProvider,
			ID:     "bazel",
			Digest: v1alpha1.DigestOfCanonicalJSON([]byte("provider")),
		},
	}
	if err := Apply(assertion, recording); err != nil {
		t.Errorf("Apply() with matching provider digest error = %v", err)
	}
}

func TestApplicabilityMatrix(t *testing.T) {
	tests := []struct {
		name          string
		setInputs     v1alpha1.ValidationInputs
		currentInputs v1alpha1.ValidationInputs
		currentBase   string
		currentTree   string
		currentCommit string
		want          v1alpha1.EvidenceApplicability
	}{
		{
			name: "tree-bound survives message-only rewrite",
			setInputs: v1alpha1.ValidationInputs{
				BaseOID: "base-001", TreeOID: "tree-001",
			},
			currentBase: "base-001", currentTree: "tree-001",
			currentCommit: "commit-002",
			want:          v1alpha1.EvidenceApplicable,
		},
		{
			name: "commit-bound goes stale on rewrite",
			setInputs: v1alpha1.ValidationInputs{
				BaseOID: "base-001", CommitOID: "commit-001",
			},
			currentBase: "base-001", currentTree: "tree-002",
			currentCommit: "commit-002",
			want:          v1alpha1.EvidenceStale,
		},
		{
			name: "tree only stays applicable on changed commit",
			setInputs: v1alpha1.ValidationInputs{
				TreeOID: "tree-001",
			},
			currentBase: "base-001", currentTree: "tree-001",
			currentCommit: "commit-002",
			want:          v1alpha1.EvidenceApplicable,
		},
		{
			name: "fully bound stale on changed tree",
			setInputs: v1alpha1.ValidationInputs{
				BaseOID: "base-001", TreeOID: "tree-001",
				CommitOID: "commit-001",
			},
			currentBase: "base-001", currentTree: "tree-002",
			currentCommit: "commit-002",
			want:          v1alpha1.EvidenceStale,
		},
		{
			name: "unknown when inputs missing",
			setInputs: v1alpha1.ValidationInputs{
				ProviderRefs: []v1alpha1.Reference{providerReference("bazel")},
			},
			want: v1alpha1.EvidenceUnknown,
		},
		{
			name: "config change goes stale",
			setInputs: v1alpha1.ValidationInputs{
				BaseOID: "base-001", TreeOID: "tree-001",
				ConfigDigests: map[string]string{
					"config": v1alpha1.DigestOfCanonicalJSON([]byte("old")),
				},
			},
			currentInputs: v1alpha1.ValidationInputs{
				BaseOID: "base-001", TreeOID: "tree-001",
				ConfigDigests: map[string]string{
					"config": v1alpha1.DigestOfCanonicalJSON([]byte("new")),
				},
			},
			currentBase: "base-001", currentTree: "tree-001",
			currentCommit: "commit-001",
			want:          v1alpha1.EvidenceStale,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := mustEmitSet(t, EmitOptions{
				Profile:      v1alpha1.ImpactProfileChangedFast,
				Candidate:    repositoryReference("alwaldend/src"),
				WorkingScope: "projects/agents",
			})
			set.Inputs = test.setInputs
			inputs := test.currentInputs
			if inputs.BaseOID == "" && inputs.TreeOID == "" &&
				inputs.CommitOID == "" &&
				len(inputs.ConfigDigests) == 0 {
				inputs = test.setInputs
			}
			got := Applicability(
				set,
				inputs,
				test.currentBase,
				test.currentTree,
				test.currentCommit,
			)
			if got != test.want {
				t.Errorf(
					"Applicability() = %q, want %q",
					got,
					test.want,
				)
			}
		})
	}
}

func passingResults() []v1alpha1.ValidationResult {
	return []v1alpha1.ValidationResult{
		{CheckID: "build", Status: "pass", DurationMS: 10},
	}
}

func repositoryReference(id string) v1alpha1.Reference {
	return v1alpha1.Reference{Kind: v1alpha1.ReferenceRepository, ID: id}
}

func providerReference(id string) v1alpha1.Reference {
	return v1alpha1.Reference{Kind: v1alpha1.ReferenceProvider, ID: id}
}

func goalReference(id string) v1alpha1.Reference {
	return v1alpha1.Reference{Kind: v1alpha1.ReferenceGoal, ID: id}
}

func mustEmitSet(
	t *testing.T,
	opts EmitOptions,
) v1alpha1.ValidationSet {
	t.Helper()
	set, err := EmitValidationSet(opts, passingResults())
	if err != nil {
		t.Fatalf("EmitValidationSet() error = %v", err)
	}
	return set
}

func forgeDigest(original string) string {
	if original == "" {
		return original
	}
	runes := []rune(original)
	runes[len(runes)-1] = 'f'
	return string(runes)
}
