package catalogv1alpha1

import (
	"strings"
	"testing"
)

func sampleEnvelope() CatalogEnvelope {
	return CatalogEnvelope{
		Schema:            APIVersion + "/" + KindTopologyCatalog,
		Kind:              KindTopologyCatalog,
		ID:                "agent-system.topology",
		DerivationVersion: "1.0.0",
		ProducerRef:       "repository.topology-compiler",
		SourceRevision:    "0123456789abcdef0123456789abcdef01234567",
		Inputs: []CatalogInput{
			{Path: "projects", Role: "owner-root", Digest: digest([]byte("projects"))},
		},
		Bounds: CatalogBounds{
			Eligible: 1, Emitted: 1, Unavailable: 0,
			MaxItems: 100, MaxInputBytes: 1 << 20, MaxOutputBytes: 1 << 20,
		},
		Completeness: CompletenessComplete,
		Limitations:  []string{},
		Conflicts:    []CatalogConflict{},
	}
}

func TestEnvelopeValidation(t *testing.T) {
	envelope := sampleEnvelope()
	if err := envelope.Validate([]string{"component.demo"}); err != nil {
		t.Fatalf("valid envelope rejected: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(*CatalogEnvelope)
	}{
		{"wrong schema", func(e *CatalogEnvelope) { e.Schema = "wrong" }},
		{"bad kind", func(e *CatalogEnvelope) { e.Kind = "bad kind" }},
		{"bad derivation", func(e *CatalogEnvelope) { e.DerivationVersion = "1.0" }},
		{"bad revision", func(e *CatalogEnvelope) { e.SourceRevision = "" }},
		{"bad input path", func(e *CatalogEnvelope) {
			e.Inputs[0].Path = "../escape"
		}},
		{"bad input digest", func(e *CatalogEnvelope) {
			e.Inputs[0].Digest = "md5:abc"
		}},
		{"bad bounds", func(e *CatalogEnvelope) { e.Bounds.MaxItems = 0 }},
		{"negative emitted", func(e *CatalogEnvelope) { e.Bounds.Emitted = -1 }},
		{"emitted over max", func(e *CatalogEnvelope) { e.Bounds.Emitted = 200 }},
		{"truncated mismatch", func(e *CatalogEnvelope) {
			e.Bounds.Truncated = true
			e.Completeness = CompletenessComplete
		}},
		{"non-complete no limitation", func(e *CatalogEnvelope) {
			e.Completeness = CompletenessPartial
			e.Limitations = nil
		}},
		{"unknown completeness", func(e *CatalogEnvelope) {
			e.Completeness = "weird"
		}},
		{"bad conflict", func(e *CatalogEnvelope) {
			e.Conflicts = []CatalogConflict{{ID: "bad id", Code: "x"}}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value := sampleEnvelope()
			tc.mutate(&value)
			if err := value.Validate([]string{"component.demo"}); err == nil {
				t.Fatalf("%s: invalid envelope accepted", tc.name)
			}
		})
	}
}

func TestTopologyCatalogValidationRejectsDuplicateAndUnknown(t *testing.T) {
	catalog := TopologyCatalog{
		CatalogEnvelope: sampleEnvelope(),
		Trees:           []TopologyTree{},
		Components: []TopologyComponent{
			{
				ID: "agents", Path: "projects/agents", OwnerReadme: "projects/agents/README.md",
				BuildPath: "projects/agents/BUILD.bazel", Title: "Agents",
				Lifecycle: "active", DocsState: "owned",
			},
		},
		Workspaces: []TopologyWorkspace{},
	}
	if err := catalog.Validate(); err != nil {
		t.Fatalf("valid topology rejected: %v", err)
	}
	duplicate := TopologyCatalog{
		CatalogEnvelope: sampleEnvelope(),
		Trees:           []TopologyTree{},
		Components: []TopologyComponent{
			{ID: "agents", Lifecycle: "active", DocsState: "owned"},
			{ID: "agents", Path: "projects/agents", Lifecycle: "active", DocsState: "owned"},
		},
		Workspaces: []TopologyWorkspace{},
	}
	if err := duplicate.Validate(); err == nil {
		t.Fatal("duplicate component identity accepted")
	}
	unknownLifecycle := TopologyCatalog{
		CatalogEnvelope: sampleEnvelope(),
		Trees:           []TopologyTree{},
		Workspaces:      []TopologyWorkspace{},
	}
	unknownLifecycle.Components = []TopologyComponent{
		{ID: "agents", Lifecycle: "mystery", DocsState: "owned"},
	}
	if err := unknownLifecycle.Validate(); err == nil {
		t.Fatal("unknown lifecycle accepted")
	}
}

func TestCanonicalJSONTopologyDeterministicDigest(t *testing.T) {
	catalog := TopologyCatalog{
		CatalogEnvelope: sampleEnvelope(),
		Trees:           []TopologyTree{},
		Components: []TopologyComponent{
			{
				ID: "agents", Path: "projects/agents", OwnerReadme: "projects/agents/README.md",
				BuildPath: "projects/agents/BUILD.bazel", Title: "Agents",
				Lifecycle: "active", DocsState: "owned",
			},
		},
		Workspaces: []TopologyWorkspace{},
	}
	first, err := CanonicalJSONTopology(catalog)
	if err != nil {
		t.Fatalf("canonical encode: %v", err)
	}
	second, err := CanonicalJSONTopology(catalog)
	if err != nil {
		t.Fatalf("canonical encode: %v", err)
	}
	if string(first) != string(second) {
		t.Fatal("canonical JSON is not deterministic")
	}
	if !strings.Contains(string(first), `"digest":"sha256:`) {
		t.Fatalf("digest field missing: %s", first)
	}
}

func TestDecodeStrictRejectsUnknownAndTrailing(t *testing.T) {
	sample := `{"kind":"TopologyCatalog","extra":1}`
	if err := DecodeStrict([]byte(sample), &struct {
		Kind string `json:"kind"`
	}{}); err == nil {
		t.Fatal("unknown field accepted")
	}
	sample = `{"kind":"TopologyCatalog"} {}`
	if err := DecodeStrict([]byte(sample), &struct {
		Kind string `json:"kind"`
	}{}); err == nil {
		t.Fatal("trailing JSON accepted")
	}
}

func TestGoalCatalogRejectsInvalidContinuation(t *testing.T) {
	base := goalCatalogFixture()
	if base.Goals[0].Identity == nil || base.Goals[0].CoarseStatus == nil {
		t.Fatal("fixture goal must be available")
	}
	base.Goals[0].Continuation = &GoalContinuation{
		ActiveAttempt: "resume-attempt",
		NextAction:    "Continue.",
	}
	if err := base.Validate(); err != nil {
		t.Fatalf("valid continuation rejected: %v", err)
	}

	for _, test := range []struct {
		name   string
		mutate func(*GoalCatalog)
	}{
		{
			name: "continuation on closed goal",
			mutate: func(catalog *GoalCatalog) {
				catalog.Goals[0].CoarseStatus.Outcome = "achieved"
			},
		},
		{
			name: "continuation without active attempt",
			mutate: func(catalog *GoalCatalog) {
				catalog.Goals[0].Continuation.ActiveAttempt = ""
			},
		},
		{
			name: "continuation on unavailable goal",
			mutate: func(catalog *GoalCatalog) {
				catalog.Goals[0].Availability = "unavailable"
				catalog.Goals[0].Reason = "invalid record"
				catalog.Goals[0].Identity = nil
				catalog.Goals[0].CoarseStatus = nil
			},
		},
		{
			name: "unsorted affected criteria",
			mutate: func(catalog *GoalCatalog) {
				catalog.Goals[0].Continuation.AffectedCriteria = []string{
					"criterion-002", "criterion-001",
				}
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			catalog := base
			catalog.Goals = append([]GoalRecord(nil), base.Goals...)
			if test.name == "unsorted affected criteria" {
				catalog.Goals[0].Continuation = &GoalContinuation{
					ActiveAttempt:    "resume-attempt",
					AffectedCriteria: []string{"criterion-002", "criterion-001"},
				}
			} else {
				test.mutate(&catalog)
			}
			if err := catalog.Validate(); err == nil {
				t.Fatalf("invalid continuation was accepted")
			}
		})
	}
}

func goalCatalogFixture() GoalCatalog {
	return GoalCatalog{
		CatalogEnvelope: CatalogEnvelope{
			Schema:            APIVersion + "/" + KindGoalCatalog,
			Kind:              KindGoalCatalog,
			ID:                "agent-system.goal",
			DerivationVersion: "1.0.0",
			SourceRevision:    "0123456789abcdef0123456789abcdef01234567",
			Inputs: []CatalogInput{{
				Path:   "projects/agents/goals/example/goal.yaml",
				Role:   "goal-manifest",
				Digest: "sha256:" + strings.Repeat("1", 64),
			}},
			Bounds: CatalogBounds{
				Eligible: 1, Emitted: 1, Unavailable: 0,
				MaxItems: 1000, MaxInputBytes: 1 << 20, MaxOutputBytes: 1 << 20,
			},
			Completeness: CompletenessComplete,
			Limitations:  []string{},
			Conflicts:    []CatalogConflict{},
		},
		Goals: []GoalRecord{{
			CandidatePath: "projects/agents/goals/example",
			Availability:  "available",
			Identity: &GoalCoarseIdentity{
				Name:      "example",
				OwnerRoot: "projects/agents",
				Scope:     "project",
			},
			CoarseStatus: &GoalCoarseStatus{
				Outcome:   "open",
				Execution: "active",
			},
		}},
	}
}
