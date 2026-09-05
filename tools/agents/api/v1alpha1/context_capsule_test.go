package v1alpha1

import (
	"bytes"
	"encoding/json"
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

func TestContextCapsulePreservesUnknownVersusCleanGitState(t *testing.T) {
	for _, unavailable := range []bool{false, true} {
		capsule := validCapsule()
		clean := false
		capsule.Identity.Git = &CapsuleGitObservation{
			Revision: strings.Repeat("a", 40), Dirty: &clean,
			ObservedAt: "2026-09-05T12:00:00Z",
		}
		if unavailable {
			capsule.Identity.Git.Dirty = nil
			capsule.Identity.Git.Unavailable = []string{"status unavailable"}
		}
		content, err := CanonicalContextJSON(capsule)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := DecodeContextCapsule(content)
		if err != nil {
			t.Fatal(err)
		}
		if (decoded.Identity.Git.Dirty == nil) != unavailable {
			t.Fatalf("unknown/false Git state conflated: %s", content)
		}
		want := "Observed Git dirty state: false"
		if unavailable {
			want = "Observed Git dirty state: unavailable"
		}
		if !strings.Contains(RenderContextMarkdown(decoded), want) {
			t.Fatal("Markdown did not preserve observed Git availability")
		}
	}
}

func TestContextCapsuleRejectsUnexplainedGitUnavailable(t *testing.T) {
	capsule := validCapsule()
	capsule.Identity.Git = &CapsuleGitObservation{ObservedAt: "2026-09-05T12:00:00Z"}
	if _, err := CanonicalContextJSON(capsule); err == nil {
		t.Fatal("accepted unknown Git state without reason")
	}
}

func TestContextCapsuleValidatesIdentitySourceClaims(t *testing.T) {
	for _, test := range []struct {
		name      string
		change    func(*CapsuleIdentity)
		wantError string
	}{
		{name: "matching-observations"},
		{name: "matching-clean-observation", change: func(identity *CapsuleIdentity) {
			identity.DirtyInputs = false
			*identity.Git.Dirty = false
		}},
		{name: "matching-sha256-observation", change: func(identity *CapsuleIdentity) {
			identity.Revision = strings.Repeat("b", 64)
			identity.Git.Revision = identity.Revision
		}},
		{name: "legacy-without-labels-or-git", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource, identity.DirtyInputsSource = "", ""
			identity.Git = nil
			identity.Revision = "legacy-revision"
			identity.DirtyInputs = false
		}},
		{name: "legacy-values-differ-from-independent-git", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource, identity.DirtyInputsSource = "", ""
			identity.Revision = "legacy-revision"
			identity.DirtyInputs = false
		}},
		{name: "legacy-empty-revision", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = ""
			identity.Revision = ""
		}},
		{name: "legacy-whitespace-revision", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = ""
			identity.Revision = " \t\n"
		}},
		{name: "caller-declarations-differ-from-git", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource, identity.DirtyInputsSource = "caller-declared", "caller-declared"
			identity.Revision = "declared-revision"
			identity.DirtyInputs = false
		}},
		{name: "caller-declarations-without-git", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource, identity.DirtyInputsSource = "caller-declared", "caller-declared"
			identity.Git = nil
		}},
		{name: "caller-nonblank-revision-with-spaces", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "caller-declared"
			identity.Revision = " declared-revision "
		}},
		{name: "caller-empty-revision", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "caller-declared"
			identity.Revision = ""
		}, wantError: "requires a nonblank revision"},
		{name: "caller-whitespace-revision", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "caller-declared"
			identity.Revision = " \t\r\n\u00a0"
		}, wantError: "requires a nonblank revision"},
		{name: "fallback-values-without-git", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource, identity.DirtyInputsSource = "input-digest", "conservative-default"
			identity.Revision = identity.InputDigest
			identity.Git = nil
		}},
		{name: "observed-head-with-unavailable-dirty", change: func(identity *CapsuleIdentity) {
			identity.DirtyInputsSource = "conservative-default"
			identity.Git.Dirty = nil
			identity.Git.Unavailable = []string{"dirty state unavailable"}
		}},
		{name: "observed-dirty-with-unavailable-head", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "input-digest"
			identity.Revision = identity.InputDigest
			identity.Git.Revision = ""
			identity.Git.Unavailable = []string{"HEAD unavailable"}
		}},
		{name: "unknown-revision-source", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "unknown-source"
		}, wantError: "unknown capsule revision source"},
		{name: "cross-field-revision-source", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "observed-git-status"
		}, wantError: "unknown capsule revision source"},
		{name: "unknown-dirty-source", change: func(identity *CapsuleIdentity) {
			identity.DirtyInputsSource = "unknown-source"
		}, wantError: "unknown capsule dirty-input source"},
		{name: "cross-field-dirty-source", change: func(identity *CapsuleIdentity) {
			identity.DirtyInputsSource = "observed-git-head"
		}, wantError: "unknown capsule dirty-input source"},
		{name: "observed-head-mismatch", change: func(identity *CapsuleIdentity) {
			identity.Revision = strings.Repeat("b", 40)
		}, wantError: "matching Git revision observation"},
		{name: "observed-head-without-git", change: func(identity *CapsuleIdentity) {
			identity.Git = nil
		}, wantError: "matching Git revision observation"},
		{name: "observed-head-unavailable", change: func(identity *CapsuleIdentity) {
			identity.Git.Revision = ""
			identity.Git.Unavailable = []string{"HEAD unavailable"}
		}, wantError: "matching Git revision observation"},
		{name: "observed-dirty-mismatch", change: func(identity *CapsuleIdentity) {
			identity.DirtyInputs = false
		}, wantError: "matching Git dirty observation"},
		{name: "observed-dirty-without-git", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "caller-declared"
			identity.Git = nil
		}, wantError: "matching Git dirty observation"},
		{name: "observed-dirty-unavailable", change: func(identity *CapsuleIdentity) {
			identity.Git.Dirty = nil
			identity.Git.Unavailable = []string{"dirty state unavailable"}
		}, wantError: "matching Git dirty observation"},
		{name: "input-digest-mismatch", change: func(identity *CapsuleIdentity) {
			identity.RevisionSource = "input-digest"
		}, wantError: "revision to match inputDigest"},
		{name: "conservative-default-claims-clean", change: func(identity *CapsuleIdentity) {
			identity.DirtyInputsSource = "conservative-default"
			identity.DirtyInputs = false
		}, wantError: "dirtyInputs to be true"},
	} {
		t.Run(test.name, func(t *testing.T) {
			capsule := validCapsule()
			dirty := true
			capsule.Identity.Revision = strings.Repeat("a", 40)
			capsule.Identity.RevisionSource = "observed-git-head"
			capsule.Identity.DirtyInputsSource = "observed-git-status"
			capsule.Identity.Git = &CapsuleGitObservation{
				Revision: capsule.Identity.Revision, Dirty: &dirty,
				ObservedAt: "2026-09-05T12:00:00Z",
			}
			if test.change != nil {
				test.change(&capsule.Identity)
			}
			// Marshal directly so invalid external claims reach the decoder;
			// the canonical encoder must independently reject them too.
			content, err := json.Marshal(capsule)
			if err != nil {
				t.Fatal(err)
			}
			_, decodeErr := DecodeContextCapsule(content)
			_, encodeErr := CanonicalContextJSON(capsule)
			for operation, err := range map[string]error{"decode": decodeErr, "encode": encodeErr} {
				if test.wantError == "" {
					if err != nil {
						t.Fatalf("%s rejected valid source claim: %v", operation, err)
					}
				} else if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("%s error = %v, want %q", operation, err, test.wantError)
				}
			}
		})
	}
}

func TestContextCapsuleValidatesEveryGitUnavailableReason(t *testing.T) {
	for _, fields := range []string{"available", "missing-head", "missing-dirty", "missing-both"} {
		for _, test := range []struct {
			name    string
			reasons []string
			blank   bool
		}{
			{name: "omitted"},
			{name: "empty-list", reasons: []string{}},
			{name: "blank", reasons: []string{""}, blank: true},
			{name: "whitespace", reasons: []string{" \t\r\n\u00a0"}, blank: true},
			{name: "valid-then-blank", reasons: []string{"Git command unavailable", ""}, blank: true},
			{name: "blank-then-valid", reasons: []string{" \t", "Git command unavailable"}, blank: true},
			{name: "valid", reasons: []string{"Git command unavailable"}},
			{name: "multiple-valid", reasons: []string{"HEAD unavailable", "dirty state unavailable"}},
		} {
			t.Run(fields+"/"+test.name, func(t *testing.T) {
				capsule := validCapsule()
				dirty := false
				capsule.Identity.Git = &CapsuleGitObservation{
					Revision: strings.Repeat("a", 40), Dirty: &dirty,
					ObservedAt: "2026-09-05T12:00:00Z", Unavailable: test.reasons,
				}
				if fields == "missing-head" || fields == "missing-both" {
					capsule.Identity.Git.Revision = ""
				}
				if fields == "missing-dirty" || fields == "missing-both" {
					capsule.Identity.Git.Dirty = nil
				}
				wantError := ""
				if test.blank {
					wantError = "unavailable reasons must not be blank"
				} else if fields != "available" && len(test.reasons) == 0 {
					wantError = "unavailable capsule Git fields require a reason"
				}
				content, err := json.Marshal(capsule)
				if err != nil {
					t.Fatal(err)
				}
				_, decodeErr := DecodeContextCapsule(content)
				_, encodeErr := CanonicalContextJSON(capsule)
				for operation, err := range map[string]error{"decode": decodeErr, "encode": encodeErr} {
					if wantError == "" {
						if err != nil {
							t.Fatalf("%s rejected valid reasons: %v", operation, err)
						}
					} else if err == nil || !strings.Contains(err.Error(), wantError) {
						t.Fatalf("%s error = %v, want %q", operation, err, wantError)
					}
				}
			})
		}
	}
}
