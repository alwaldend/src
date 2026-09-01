package catalogv1alpha1

import (
	"strings"
	"testing"
)

func sampleEnvelopeFor(kind string) CatalogEnvelope {
	return CatalogEnvelope{
		Schema:            APIVersion + "/" + kind,
		Kind:              kind,
		ID:                "agent-system.test",
		DerivationVersion: "1.0.0",
		ProducerRef:       "repository.test-compiler",
		SourceRevision:    "0123456789abcdef0123456789abcdef01234567",
		Inputs: []CatalogInput{{
			Path: "AGENTS.md", Role: "agent-policy", Digest: digest([]byte("policy")),
		}},
		Bounds: CatalogBounds{
			Eligible: 1, Emitted: 1, MaxItems: 100,
			MaxInputBytes: 1 << 20, MaxOutputBytes: 1 << 20,
		},
		Completeness: CompletenessComplete,
		Limitations:  []string{},
		Conflicts:    []CatalogConflict{},
	}
}

func TestPolicyCatalogCanonicalRoundTrip(t *testing.T) {
	catalog := PolicyCatalog{
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
	content, err := CanonicalJSONPolicy(catalog)
	if err != nil {
		t.Fatalf("canonical encode: %v", err)
	}
	decoded, err := DecodePolicyStrict(content)
	if err != nil {
		t.Fatalf("decode round-trip: %v", err)
	}
	markdown := RenderPolicyMarkdown(decoded)
	if !strings.Contains(markdown, "JSON digest: `"+decoded.Digest+"`") {
		t.Fatalf("markdown lacks digest: %s", markdown)
	}
	if !strings.Contains(markdown, "`sourceDisclosure`: known (from `AGENTS.md`)") {
		t.Fatalf("markdown lacks axis: %s", markdown)
	}
}

func TestPolicyCatalogRejectsUnknownAxisValue(t *testing.T) {
	catalog := PolicyCatalog{
		CatalogEnvelope: sampleEnvelopeFor(KindPolicyCatalog),
		Policies: []PolicyRecord{{
			ID: "root", PathPrefix: "/", Axes: []PolicyAxis{{
				Name: "sourceDisclosure", Value: "inferred", Source: "AGENTS.md",
			}},
		}},
	}
	if _, err := CanonicalJSONPolicy(catalog); err == nil ||
		!strings.Contains(err.Error(), "unknown value") {
		t.Fatalf("expected unknown axis rejection, got %v", err)
	}
}

func TestActionCatalogCanonicalRoundTrip(t *testing.T) {
	catalog := ActionCatalog{
		CatalogEnvelope: sampleEnvelopeFor(KindActionCatalog),
		Providers: []ActionProvider{{
			ID: "goal.local-store", Owner: "projects.goal",
			DefinitionPath: "projects/goal/cmd/goal/command.go",
		}},
		Actions: []ActionRecord{{
			ID: "goal.checkpoint", ProviderRef: "goal.local-store", Owner: "projects.goal",
			SourcePath: "projects/goal/agent_operations.json", Selector: "checkpoint",
			Classification: "classified", Effects: []string{"source.write"},
			Inputs: []string{"goal revision"}, Outputs: []string{"goal state"},
			Information: []string{"public"}, CredentialUse: "none",
			NetworkUse: "none", EnvironmentSelector: "explicit goal directory",
			AuthorityGate: "goal writer", Preflight: "lock and validate",
			Verification: "record validation", Cost: "bounded",
			Cacheability: "not_cacheable", Cancellation: "atomic boundaries",
		}},
		Aliases: []ActionAlias{},
	}
	content, err := CanonicalJSONAction(catalog)
	if err != nil {
		t.Fatalf("canonical encode: %v", err)
	}
	decoded, err := DecodeActionStrict(content)
	if err != nil {
		t.Fatalf("decode round-trip: %v", err)
	}
	markdown := RenderActionMarkdown(decoded)
	if !strings.Contains(markdown, "JSON digest: `"+decoded.Digest+"`") {
		t.Fatalf("markdown lacks digest: %s", markdown)
	}
	if !strings.Contains(markdown, "`goal.checkpoint` (goal.local-store.checkpoint)") {
		t.Fatalf("markdown lacks action: %s", markdown)
	}
}

func TestActionCatalogRejectsUnknownEffect(t *testing.T) {
	catalog := ActionCatalog{
		CatalogEnvelope: sampleEnvelopeFor(KindActionCatalog),
		Providers: []ActionProvider{{
			ID: "goal.local-store", Owner: "projects.goal",
			DefinitionPath: "projects/goal/cmd/goal/command.go",
		}},
		Actions: []ActionRecord{{
			ID: "goal.checkpoint", ProviderRef: "goal.local-store", Owner: "projects.goal",
			SourcePath: "projects/goal/agent_operations.json", Selector: "checkpoint",
			Classification: "classified", Effects: []string{"arbitrary.write"},
			Inputs: []string{"goal revision"}, Outputs: []string{"goal state"},
			Information: []string{"public"}, CredentialUse: "none",
			NetworkUse: "none", EnvironmentSelector: "explicit goal directory",
			AuthorityGate: "goal writer", Preflight: "lock and validate",
			Verification: "record validation", Cost: "bounded",
			Cacheability: "not_cacheable", Cancellation: "atomic boundaries",
		}},
		Aliases: []ActionAlias{},
	}
	if _, err := CanonicalJSONAction(catalog); err == nil ||
		!strings.Contains(err.Error(), "unknown effect") {
		t.Fatalf("expected unknown effect rejection, got %v", err)
	}
}
