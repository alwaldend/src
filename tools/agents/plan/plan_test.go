package plan

import (
	"strings"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func testEnvelope() catalogv1alpha1.CatalogEnvelope {
	return catalogv1alpha1.CatalogEnvelope{
		Schema:            catalogv1alpha1.APIVersion + "/topology-catalog",
		Kind:              catalogv1alpha1.KindTopologyCatalog,
		ID:                "agent-system.topology",
		DerivationVersion: "1.0.0",
		ProducerRef:       "repository.topology-compiler",
		SourceRevision:    "0123456789abcdef0123456789abcdef01234567",
		Inputs:            []catalogv1alpha1.CatalogInput{},
		Bounds: catalogv1alpha1.CatalogBounds{
			Eligible: 1, Emitted: 1, Unavailable: 0,
			MaxItems: 100, MaxInputBytes: 1 << 20, MaxOutputBytes: 1 << 20,
		},
		Completeness: catalogv1alpha1.CompletenessComplete,
		Limitations:  []string{},
		Conflicts:    []catalogv1alpha1.CatalogConflict{},
	}
}

func testIndex() *Index {
	return &Index{
		Topology: &catalogv1alpha1.TopologyCatalog{
			CatalogEnvelope: testEnvelope(),
			Trees: []catalogv1alpha1.TopologyTree{{
				ID: "tools", Path: "tools",
				Boundary: catalogv1alpha1.TreeBoundaryTool,
			}},
			Components: []catalogv1alpha1.TopologyComponent{
				{
					ID: "agents", Path: "tools/agents",
					OwnerReadme: "tools/agents/README.md",
					BuildPath:   "tools/agents/BUILD.bazel",
					Title:       "Agents",
					Lifecycle:   "active",
					DocsState:   "owned",
				},
				{
					ID: "bazel-agent", Path: "projects/bazel_agent",
					OwnerReadme: "projects/bazel_agent/README.md",
					BuildPath:   "projects/bazel_agent/BUILD.bazel",
					Title:       "Bazel agent",
					Lifecycle:   "active",
					DocsState:   "owned",
				},
			},
			Workspaces: []catalogv1alpha1.TopologyWorkspace{
				{
					ID: "root", Path: ".", ModulePath: "MODULE.bazel",
					ModuleName: "com_alwaldend_src",
				},
				{
					ID: "agents.ws", Path: "tools/agents",
					ModulePath: "tools/agents/MODULE.bazel",
					ModuleName: "agents",
				},
			},
		},
		Capability: &catalogv1alpha1.CapabilityCatalog{
			CatalogEnvelope: testEnvelope(),
			Skills:          []catalogv1alpha1.CapabilitySkill{},
			Providers: []catalogv1alpha1.CapabilityProvider{{
				ID:             "bazel-agent",
				Owner:          "projects/bazel_agent",
				Kind:           "direct_binary",
				SourcePath:     "projects/bazel_agent/cmd/bazel_agent/main.go",
				Classification: "classified",
			}},
		},
		WorkspaceCheck: &catalogv1alpha1.WorkspaceCheckCatalog{
			CatalogEnvelope: testEnvelope(),
			Workspaces: []catalogv1alpha1.WorkspaceRecord{
				{
					ID: "root", Path: ".", ModulePath: "MODULE.bazel",
					ModuleName: "com_alwaldend_src",
				},
				{
					ID: "agents.ws", Path: "tools/agents",
					ModulePath: "tools/agents/MODULE.bazel",
					ModuleName: "agents",
				},
			},
		},
	}
}

func intent() v1alpha1.Reference {
	return v1alpha1.Reference{
		Kind: v1alpha1.ReferenceTask,
		ID:   "task/example",
	}
}

func mustPlan(
	t *testing.T,
	profile ImpactProfile,
	opts PlanOptions,
) v1alpha1.ImpactPlan {
	t.Helper()
	opts.Profile = profile
	opts.Intent = "example"
	plan, err := Plan(intent(), opts, testIndex())
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	return plan
}

func TestPlanDeterministicDigest(t *testing.T) {
	first := mustPlan(t, ProfileChangedFast, PlanOptions{
		AffectedPaths: []string{"tools/agents", "projects/bazel_agent"},
	})
	second := mustPlan(t, ProfileChangedFast, PlanOptions{
		AffectedPaths: []string{"tools/agents", "projects/bazel_agent"},
	})
	if first.Digest != second.Digest {
		t.Fatalf("identical inputs produced different digests: %s vs %s",
			first.Digest, second.Digest)
	}
	if !strings.HasPrefix(first.Digest, "sha256:") {
		t.Fatalf("digest is not sha256: %s", first.Digest)
	}
}

func TestProfileMappingChangedFastVsFullAudit(t *testing.T) {
	fast := mustPlan(t, ProfileChangedFast, PlanOptions{
		AffectedPaths: []string{"tools/agents"},
	})
	if fast.CostClass != v1alpha1.CostClassLow {
		t.Fatalf("changed/fast cost = %s, want low", fast.CostClass)
	}
	if len(fast.Targets) != 1 || fast.Targets[0].Label != "//tools/agents:all" {
		t.Fatalf("changed/fast targets = %+v", fast.Targets)
	}
	if len(fast.Checks) != 3 {
		t.Fatalf("changed/fast checks = %d, want 3", len(fast.Checks))
	}
	if fast.Cacheability != v1alpha1.Cacheable {
		t.Fatalf("changed/fast cacheability = %s, want cacheable",
			fast.Cacheability)
	}

	audit := mustPlan(t, ProfileFullAudit, PlanOptions{
		AffectedPaths: []string{"tools/agents"},
	})
	if audit.CostClass != v1alpha1.CostClassHigh {
		t.Fatalf("full/audit cost = %s, want high", audit.CostClass)
	}
	if len(audit.Targets) != 1 || audit.Targets[0].Label != "//..." {
		t.Fatalf("full/audit targets = %+v", audit.Targets)
	}
	if audit.Cacheability != v1alpha1.NotCacheable {
		t.Fatalf("full/audit cacheability = %s, want not cacheable",
			audit.Cacheability)
	}
	checkIDs := map[string]bool{}
	for _, check := range audit.Checks {
		checkIDs[check.Identifier] = true
	}
	if !checkIDs["full-repo-check"] {
		t.Fatalf("full/audit checks missing full-repo-check: %+v",
			audit.Checks)
	}
}

func TestWorkspaceProfileTargetsTrackedWorkspace(t *testing.T) {
	workspace := mustPlan(t, ProfileWorkspace, PlanOptions{
		AffectedPaths: []string{"tools/agents/plan/plan.go"},
	})
	if len(workspace.Targets) != 1 ||
		workspace.Targets[0].Label != "//tools/agents:all" {
		t.Fatalf("workspace targets = %+v", workspace.Targets)
	}
	if workspace.CostClass != v1alpha1.CostClassMedium {
		t.Fatalf("workspace cost = %s, want medium", workspace.CostClass)
	}
}

func TestFreshEvidenceBindsEvidenceRefs(t *testing.T) {
	fresh := mustPlan(t, ProfileFreshEvidence, PlanOptions{
		AffectedPaths: []string{"tools/agents"},
	})
	if len(fresh.EvidenceRefs) != 1 ||
		fresh.EvidenceRefs[0].Reference.ID != "task/example" {
		t.Fatalf("fresh evidence refs = %+v", fresh.EvidenceRefs)
	}
	freshDigest := fresh.Digest
	changed := mustPlan(t, ProfileChangedFast, PlanOptions{
		AffectedPaths: []string{"tools/agents"},
	})
	if fresh.Digest == changed.Digest {
		t.Fatalf("fresh/evidence digest equals changed/fast digest")
	}
	_ = freshDigest
}

func TestUnknownScopeProducesCoverageGap(t *testing.T) {
	plan, err := Plan(intent(), PlanOptions{
		Profile:       ProfileChangedFast,
		Intent:        "example",
		AffectedPaths: []string{"unknown/path"},
	}, testIndex())
	if err != nil {
		t.Fatalf("Plan() error = %v, want gap not error", err)
	}
	found := false
	for _, gap := range plan.CoverageGaps {
		if strings.Contains(gap, "unknown/path") {
			found = true
		}
	}
	if !found {
		t.Fatalf("coverage gaps = %+v, want unknown/path reason",
			plan.CoverageGaps)
	}
}

func TestRootModuleChangeExpandsScope(t *testing.T) {
	plan, err := Plan(intent(), PlanOptions{
		Profile:       ProfileChangedFast,
		Intent:        "example",
		AffectedPaths: []string{"MODULE.bazel"},
	}, testIndex())
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(plan.Targets) != 1 || plan.Targets[0].Label != "//..." {
		t.Fatalf("expanded targets = %+v", plan.Targets)
	}
	found := false
	for _, gap := range plan.CoverageGaps {
		if strings.Contains(gap, "expanded") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected expansion gap, got %+v", plan.CoverageGaps)
	}
}

func TestForbiddenEffectsDefault(t *testing.T) {
	plan := mustPlan(t, ProfileChangedFast, PlanOptions{
		AffectedPaths: []string{"tools/agents"},
	})
	for _, effect := range []v1alpha1.Effect{
		v1alpha1.EffectRemoteWrite,
		v1alpha1.EffectRemoteDestroy,
	} {
		found := false
		for _, forbidden := range plan.ForbiddenEffects {
			if forbidden == effect {
				found = true
			}
		}
		if !found {
			t.Fatalf("effect %s not forbidden by default", effect)
		}
	}
	for _, effect := range plan.Effects {
		if effect == v1alpha1.EffectRemoteWrite ||
			effect == v1alpha1.EffectRemoteDestroy {
			t.Fatalf("remote effect %s required without request", effect)
		}
	}
	if plan.Escalation != "" {
		t.Fatalf("escalation set without remote request: %s",
			plan.Escalation)
	}
}

func TestRequestedRemoteEffectsNoteEscalation(t *testing.T) {
	opts := PlanOptions{
		Profile:       ProfileChangedFast,
		Intent:        "example",
		AffectedPaths: []string{"tools/agents"},
		RequestRemote: true,
	}
	plan, err := Plan(intent(), opts, testIndex())
	if err != nil {
		t.Fatal(err)
	}
	if !containsEffect(plan.Effects, v1alpha1.EffectRemoteWrite) ||
		!containsEffect(plan.Effects, v1alpha1.EffectRemoteDestroy) {
		t.Fatalf("requested remote effects missing: %+v", plan.Effects)
	}
	if len(plan.ForbiddenEffects) != 0 {
		t.Fatalf("forbidden effects still set when requested: %+v",
			plan.ForbiddenEffects)
	}
	if !strings.Contains(plan.Escalation, "authority") {
		t.Fatalf("escalation does not note authority: %s", plan.Escalation)
	}
}

func TestEffectAndCheckCoverage(t *testing.T) {
	profiles := []ImpactProfile{
		ProfileChangedFast, ProfileWorkspace, ProfileFreshEvidence,
		ProfileFullAudit, ProfileDiagnose,
	}
	for _, profile := range profiles {
		plan := mustPlan(t, profile, PlanOptions{
			AffectedPaths: []string{"tools/agents"},
		})
		if len(plan.Effects) == 0 {
			t.Fatalf("%s: no effects", profile)
		}
		if !containsEffect(plan.Effects, v1alpha1.EffectSourceRead) {
			t.Fatalf("%s: missing source.read", profile)
		}
		if len(plan.Checks) == 0 {
			t.Fatalf("%s: no checks", profile)
		}
		if len(plan.Targets) == 0 {
			t.Fatalf("%s: no targets", profile)
		}
	}
	diagnose := mustPlan(t, ProfileDiagnose, PlanOptions{
		AffectedPaths: []string{"tools/agents"},
	})
	if containsEffect(diagnose.Effects, v1alpha1.EffectCodeExecute) {
		t.Fatalf("diagnose must not include code.execute: %+v",
			diagnose.Effects)
	}
}

func containsEffect(effects []v1alpha1.Effect, want v1alpha1.Effect) bool {
	for _, effect := range effects {
		if effect == want {
			return true
		}
	}
	return false
}
