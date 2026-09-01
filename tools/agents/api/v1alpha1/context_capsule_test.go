package v1alpha1

import (
	"bytes"
	"strings"
	"testing"
)

func validCapsule() ContextCapsule {
	return ContextCapsule{
		APIVersion: APIVersion,
		Kind:       KindContextCapsule,
		Identity: CapsuleIdentity{
			Repository:    "alwaldend/src",
			WorkspaceRoot: ".",
			WorktreePath:  ".",
			Revision:      "abc123",
			DirtyInputs:   true,
			InputDigest:   "sha256:" + strings.Repeat("ab", 32),
			ByteSize:      1024,
		},
		Outcome: CapsuleOutcome{
			RequestedOutcome: "orient",
			Authority:        "repository.defaults",
			Budget:           Budget{Calls: 8, Bytes: 131072, DurationMS: 15000},
		},
		Provenance: CapsuleProvenance{
			Freshness:    "fresh",
			Completeness: CompletenessComplete,
		},
	}
}

func TestContextCapsuleRoundTripsDeterministically(t *testing.T) {
	capsule := validCapsule()
	first, err := CanonicalContextJSON(capsule)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeContextCapsule(first)
	if err != nil {
		t.Fatal(err)
	}
	second, err := CanonicalContextJSON(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatalf("canonical round trip changed bytes:\nfirst: %s\nsecond: %s", first, second)
	}
	withUnknown := bytes.Replace(first, []byte(`"kind":"context-capsule"`),
		[]byte(`"kind":"context-capsule","unknown":true`), 1)
	if _, err := DecodeContextCapsule(withUnknown); err == nil ||
		!strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("DecodeContextCapsule() error = %v, want unknown-field rejection", err)
	}
}

func TestContextCapsuleRejectsIncompleteIdentity(t *testing.T) {
	capsule := validCapsule()
	capsule.Identity.InputDigest = ""
	if _, err := CanonicalContextJSON(capsule); err == nil ||
		!strings.Contains(err.Error(), "incomplete") {
		t.Fatalf("CanonicalContextJSON() error = %v, want incomplete-identity rejection", err)
	}
}

func TestContextCapsuleRejectsEmptyProvenance(t *testing.T) {
	capsule := validCapsule()
	capsule.Provenance.Completeness = ""
	if _, err := CanonicalContextJSON(capsule); err == nil ||
		!strings.Contains(err.Error(), "completeness") {
		t.Fatalf("CanonicalContextJSON() error = %v, want completeness rejection", err)
	}
}

func TestContextCapsuleMarkdownUsesSameData(t *testing.T) {
	capsule := validCapsule()
	first, err := CanonicalContextJSON(capsule)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeContextCapsule(first)
	if err != nil {
		t.Fatal(err)
	}
	markdown := RenderContextMarkdown(decoded)
	if !strings.Contains(markdown, decoded.ID) ||
		!strings.Contains(markdown, decoded.Identity.Repository) ||
		!strings.Contains(markdown, string(decoded.Provenance.Completeness)) {
		t.Fatalf("Markdown render does not project the same data: %s", markdown)
	}
}
