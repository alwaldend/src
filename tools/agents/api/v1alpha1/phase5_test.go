package v1alpha1

import (
	"strings"
	"testing"
)

func phase5Reference(kind ReferenceKind, id string) Reference {
	return Reference{Kind: kind, ID: id, Version: "v1"}
}

func validSkillCase() SkillCase {
	return SkillCase{
		APIVersion:   APIVersion,
		Kind:         "SkillCase",
		ID:           "case/example",
		SkillID:      "example",
		Provenance:   phase5Reference(ReferenceRepository, "repository/source"),
		Metric:       "routing/accuracy",
		EvidenceTier: TierFixtureTested,
		SourceRef:    phase5Reference(ReferenceArtifact, "artifact/case"),
	}
}

func TestSkillCaseCanonicalRoundTrip(t *testing.T) {
	first, err := CanonicalSkillCaseJSON(validSkillCase())
	if err != nil {
		t.Fatalf("CanonicalSkillCaseJSON() error = %v", err)
	}
	decoded, err := DecodeSkillCase(first)
	if err != nil {
		t.Fatalf("DecodeSkillCase() error = %v", err)
	}
	second, err := CanonicalSkillCaseJSON(decoded)
	if err != nil {
		t.Fatalf("CanonicalSkillCaseJSON(decoded) error = %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("canonical round trip changed bytes:\n%s\n%s", first, second)
	}
}

func TestSkillCaseRejectsUnknownTier(t *testing.T) {
	skillCase := validSkillCase()
	skillCase.EvidenceTier = "best-effort"
	if _, err := CanonicalSkillCaseJSON(skillCase); err == nil ||
		!strings.Contains(err.Error(), "unknown evidence tier") {
		t.Fatalf("CanonicalSkillCaseJSON() error = %v, want tier rejection", err)
	}
}

func TestCoverageMatrixRejectsDuplicateNormalizedCase(t *testing.T) {
	entry := CoverageEntry{
		SkillID:      "example",
		CaseID:       "case/example",
		State:        CoverageFixtureTested,
		Metric:       "routing/accuracy",
		EvidenceTier: TierFixtureTested,
	}
	matrix := CoverageMatrix{
		APIVersion:    APIVersion,
		Kind:          "CoverageMatrix",
		ID:            "matrix/example",
		CatalogRef:    phase5Reference(ReferenceArtifact, "artifact/catalog"),
		Entries:       []CoverageEntry{entry, entry},
		Total:         2,
		TotalSkills:   2,
		CoveredSkills: 1,
	}
	if _, err := CanonicalCoverageMatrixJSON(matrix); err == nil ||
		!strings.Contains(err.Error(), "duplicate normalized case") {
		t.Fatalf("CanonicalCoverageMatrixJSON() error = %v, want duplicate rejection", err)
	}
}

func TestCoverageMatrixRejectsImpossibleTruncation(t *testing.T) {
	matrix := CoverageMatrix{
		APIVersion: APIVersion,
		Kind:       "CoverageMatrix",
		ID:         "matrix/example",
		CatalogRef: phase5Reference(ReferenceArtifact, "artifact/catalog"),
		Entries: []CoverageEntry{{
			SkillID:      "example",
			CaseID:       "case/example",
			State:        CoverageFixtureTested,
			Metric:       "routing/accuracy",
			EvidenceTier: TierFixtureTested,
		}},
		Total:         1,
		Truncated:     true,
		TotalSkills:   2,
		CoveredSkills: 1,
	}
	if _, err := CanonicalCoverageMatrixJSON(matrix); err == nil ||
		!strings.Contains(err.Error(), "truncated matrix") {
		t.Fatalf("CanonicalCoverageMatrixJSON() error = %v, want truncation rejection", err)
	}
}

func TestCoverageMatrixValidatesSkillCountsIndependently(t *testing.T) {
	matrix := CoverageMatrix{
		APIVersion: APIVersion, Kind: "CoverageMatrix", ID: "matrix/example",
		CatalogRef: phase5Reference(ReferenceArtifact, "artifact/catalog"),
		Entries: []CoverageEntry{{
			SkillID: "example", CaseID: "case/example", Metric: "routing/accuracy",
			State: CoverageConfigured, EvidenceTier: TierConfigured,
		}},
		Total: 1, TotalSkills: 2, CoveredSkills: 1,
	}
	if err := matrix.Validate(); err != nil {
		t.Fatalf("partial skill coverage with complete output: %v", err)
	}
	for _, counts := range [][2]int{{2, 0}, {2, 2}, {0, 1}, {-1, 1}, {4097, 1}} {
		candidate := matrix
		candidate.TotalSkills, candidate.CoveredSkills = counts[0], counts[1]
		if err := candidate.Validate(); err == nil ||
			!strings.Contains(err.Error(), "skill counts") {
			t.Fatalf("counts %v: error = %v, want rejection", counts, err)
		}
	}
	matrix.Total = 2
	if err := matrix.Validate(); err == nil ||
		!strings.Contains(err.Error(), "untruncated matrix") {
		t.Fatalf("missing output without truncation: error = %v", err)
	}
	matrix.Truncated = true
	if err := matrix.Validate(); err != nil {
		t.Fatalf("explicitly truncated output: %v", err)
	}
}

func TestLearningProposalRequiresRepeatedFrictionAndRetirement(t *testing.T) {
	proposal := LearningProposal{
		APIVersion:      APIVersion,
		Kind:            "LearningProposal",
		ID:              "proposal/example",
		Owner:           "projects/agents",
		DefectSignature: "repeated-routing-error",
		Reproducer:      phase5Reference(ReferenceArtifact, "artifact/reproducer"),
		RegressionRef:   phase5Reference(ReferenceArtifact, "artifact/regression"),
		ContractRef:     phase5Reference(ReferenceArtifact, "artifact/contract"),
		Fallback:        "keep current routing until fixture passes",
		ResourceBudget:  Budget{Calls: 1, Bytes: 1024, DurationMS: 1000, Concurrency: 1},
		ValidationRefs: []Reference{
			phase5Reference(ReferenceArtifact, "artifact/validation"),
		},
		FrictionRefs: []Reference{
			phase5Reference(ReferenceArtifact, "artifact/friction-one"),
		},
		RetirementRule: "remove after two successful routing cycles",
	}
	if _, err := CanonicalLearningProposalJSON(proposal); err == nil ||
		!strings.Contains(err.Error(), "at least two friction references") {
		t.Fatalf("CanonicalLearningProposalJSON() error = %v, want friction rejection", err)
	}
	proposal.FrictionRefs = append(proposal.FrictionRefs,
		phase5Reference(ReferenceArtifact, "artifact/friction-two"))
	if _, err := CanonicalLearningProposalJSON(proposal); err != nil {
		t.Fatalf("CanonicalLearningProposalJSON() error = %v", err)
	}
}

func TestDisclosureBudgetRequiresProvenance(t *testing.T) {
	budget := DisclosureBudget{
		APIVersion:          APIVersion,
		Kind:                "DisclosureBudget",
		ID:                  "budget/example",
		CatalogSummaryBytes: 1024,
		ContentDigestBytes:  256,
		MaxBodyBytes:        8192,
		MaxSchemaBytes:      4096,
		EvidenceTier:        TierFixtureTested,
	}
	if _, err := CanonicalDisclosureBudgetJSON(budget); err == nil ||
		!strings.Contains(err.Error(), "requires provenance") {
		t.Fatalf("CanonicalDisclosureBudgetJSON() error = %v, want provenance rejection", err)
	}
	budget.Provenance = []Reference{
		phase5Reference(ReferenceArtifact, "artifact/measurement"),
	}
	if _, err := CanonicalDisclosureBudgetJSON(budget); err != nil {
		t.Fatalf("CanonicalDisclosureBudgetJSON() error = %v", err)
	}
}
