package catalogv1alpha1

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestVerifySelfDigestRejectsEmptyAndMismatchedDigest ensures the strict
// decoders enforce the content-addressed catalog contract: a catalog whose
// fields were changed while retaining an old syntactically valid digest must
// be rejected, and a missing digest must be rejected even when every field is
// otherwise valid.
func TestVerifySelfDigestRejectsEmptyAndMismatchedDigest(t *testing.T) {
	valid := PolicyCatalog{
		CatalogEnvelope: sampleEnvelopeFor(KindPolicyCatalog),
		Policies: []PolicyRecord{{
			ID:                "root",
			PathPrefix:        "/",
			Precedence:        0,
			AgentPolicySource: "AGENTS.md",
			Axes: []PolicyAxis{
				{Name: "sourceDisclosure", Value: "known", Source: "AGENTS.md"},
				{Name: "bazelVisibility", Value: "unknown", Source: "AGENTS.md"},
			},
		}},
	}
	content, err := CanonicalJSONPolicy(valid)
	if err != nil {
		t.Fatal(err)
	}

	// Tamper with a field, leaving the old digest in place.
	var tampered PolicyCatalog
	if err := DecodeStrict(content, &tampered); err != nil {
		t.Fatal(err)
	}
	tampered.Policies[0].Precedence++
	tamperedBytes, err := json.Marshal(tampered)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodePolicyStrict(tamperedBytes); err == nil ||
		!strings.Contains(err.Error(), "self-digest mismatch") {
		t.Fatalf("tampered catalog decode error = %v, want self-digest mismatch", err)
	}

	// Strip the digest entirely.
	var missingDigest PolicyCatalog
	if err := DecodeStrict(content, &missingDigest); err != nil {
		t.Fatal(err)
	}
	missingDigest.Digest = ""
	missingBytes, err := json.Marshal(missingDigest)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodePolicyStrict(missingBytes); err == nil ||
		!strings.Contains(err.Error(), "digest is required") {
		t.Fatalf("missing-digest decode error = %v, want digest-required", err)
	}

	// The unmodified round trip still decodes.
	if _, err := DecodePolicyStrict(content); err != nil {
		t.Fatalf("valid round trip failed: %v", err)
	}
}
