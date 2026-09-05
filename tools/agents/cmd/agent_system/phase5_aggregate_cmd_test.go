package agent_system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func aggregateReference(
	kind v1alpha1.ReferenceKind,
	id string,
) v1alpha1.Reference {
	return v1alpha1.Reference{Kind: kind, ID: id, Version: "workspace"}
}

func aggregateCase(skillID, caseID, metric string) v1alpha1.SkillCase {
	return v1alpha1.SkillCase{
		APIVersion:   v1alpha1.APIVersion,
		Kind:         "SkillCase",
		ID:           caseID,
		SkillID:      skillID,
		Provenance:   aggregateReference("repository", "repository/source"),
		Metric:       metric,
		EvidenceTier: v1alpha1.TierConfigured,
		SourceRef:    aggregateReference("artifact", "artifact/case"),
	}
}

func aggregateFixtureCatalog(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "capability.json")
	skills := []catalogv1alpha1.CapabilitySkill{{
		ID:                 "example",
		Owner:              "test",
		CanonicalPath:      "test/example",
		DiscoveryPath:      ".agents/skills/example",
		Layer:              "procedure",
		Activation:         "test",
		ContextCost:        "small",
		EvaluationMaturity: "checked",
	}}
	skills = append(skills, catalogv1alpha1.CapabilitySkill{
		ID:                 "other",
		Owner:              "test",
		CanonicalPath:      "test/other",
		DiscoveryPath:      ".agents/skills/other",
		Layer:              "procedure",
		Activation:         "test",
		ContextCost:        "small",
		EvaluationMaturity: "checked",
	})
	bounds := catalogv1alpha1.CatalogBounds{
		Eligible:       2,
		Emitted:        2,
		Unavailable:    0,
		MaxItems:       100,
		MaxInputBytes:  1048576,
		MaxOutputBytes: 1048576,
	}
	catalog := catalogv1alpha1.CapabilityCatalog{
		CatalogEnvelope: catalogv1alpha1.CatalogEnvelope{
			Schema:            "agents.alwaldend.com/catalog/v1alpha1/capability-catalog",
			Kind:              catalogv1alpha1.KindCapabilityCatalog,
			ID:                "test.capability",
			DerivationVersion: "1.0.0",
			ProducerRef:       "test",
			SourceRevision:    "0123456789abcdef0123456789abcdef01234567",
			Inputs:            []catalogv1alpha1.CatalogInput{},
			Bounds:            bounds,
			Completeness:      catalogv1alpha1.CompletenessComplete,
			Limitations:       []string{},
			Conflicts:         []catalogv1alpha1.CatalogConflict{},
		},
		Skills:    skills,
		Providers: []catalogv1alpha1.CapabilityProvider{},
	}
	content, err := catalogv1alpha1.CanonicalJSONCapability(catalog)
	if err != nil {
		t.Fatalf("encode fixture capability catalog: %v", err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func aggregateFixtureCases(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "skill-cases.json")
	cases := []v1alpha1.SkillCase{
		aggregateCase("example", "case/example", "routing/positive"),
	}
	content, err := json.Marshal(cases)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunAggregatePublishesMatrix(t *testing.T) {
	root := t.TempDir()
	catalog := aggregateFixtureCatalog(t, root)
	cases := aggregateFixtureCases(t, root)
	output := filepath.Join(root, "coverage.json")
	markdown := filepath.Join(root, "coverage.md")
	var stdout bytes.Buffer
	if err := runAggregate([]string{
		"--catalog", catalog,
		"--input", cases,
		"--output", output,
		"--markdown", markdown,
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	matrix, err := v1alpha1.DecodeCoverageMatrix(content)
	if err != nil {
		t.Fatalf("decode coverage matrix: %v", err)
	}
	if len(matrix.Entries) != 1 || matrix.Total != 1 ||
		matrix.TotalSkills != 2 || matrix.CoveredSkills != 1 ||
		matrix.Truncated || matrix.Digest == "" {
		t.Fatalf("unexpected matrix: %#v", matrix)
	}
	markdownContent, err := os.ReadFile(markdown)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(markdownContent), matrix.Digest) {
		t.Fatalf("markdown does not state matrix digest")
	}
	if !strings.Contains(string(markdownContent), "`configured`") ||
		strings.Contains(string(markdownContent), "fixture-tested") {
		t.Fatalf("markdown overstates configured evidence: %s", markdownContent)
	}
	var summary map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &summary); err != nil {
		t.Fatalf("decode summary: %v", err)
	}
	if summary["digest"] != matrix.Digest {
		t.Fatalf("summary digest = %v, want %s", summary["digest"], matrix.Digest)
	}
}

func TestCoverageInventoriesSeparateCaseAndSkillCounts(t *testing.T) {
	for _, command := range []string{"aggregate", "coverage"} {
		for _, count := range []int{1, 2, 3} {
			t.Run(fmt.Sprintf("%s/%d-cases", command, count), func(t *testing.T) {
				root := t.TempDir()
				catalog := aggregateFixtureCatalog(t, root)
				cases := make([]v1alpha1.SkillCase, count)
				entries := make([]v1alpha1.CoverageEntry, count)
				for i := range cases {
					cases[i] = aggregateCase("example", fmt.Sprintf("case/example-%d", i), "routing/positive")
					entries[i] = v1alpha1.CoverageEntry{
						SkillID: "example", CaseID: cases[i].ID,
						State:        v1alpha1.CoverageConfigured,
						Metric:       cases[i].Metric,
						EvidenceTier: v1alpha1.TierConfigured,
						EvidenceRef:  cases[i].SourceRef,
					}
				}
				var input any = cases
				if command == "coverage" {
					input = entries
				}
				inputPath := filepath.Join(root, "input.json")
				writeCoverageFixture(t, inputPath, input)
				output := filepath.Join(root, "coverage.json")
				var stdout bytes.Buffer
				if err := Run([]string{
					command, "--catalog", catalog,
					"--input", inputPath, "--output", output,
				}, &stdout); err != nil {
					t.Fatal(err)
				}
				content, err := os.ReadFile(output)
				if err != nil {
					t.Fatal(err)
				}
				matrix, err := v1alpha1.DecodeCoverageMatrix(content)
				if err != nil {
					t.Fatal(err)
				}
				if matrix.Total != count || len(matrix.Entries) != count ||
					matrix.TotalSkills != 2 || matrix.CoveredSkills != 1 || matrix.Truncated {
					t.Fatalf("case and skill counts conflated: %#v", matrix)
				}
			})
		}
	}
}

func TestCoverageInventoriesRejectUnsupportedEvidence(t *testing.T) {
	for _, command := range []string{"aggregate", "coverage"} {
		for _, tier := range []v1alpha1.EvidenceTier{
			v1alpha1.TierRouted, v1alpha1.TierFixtureTested, v1alpha1.TierLive,
		} {
			t.Run(command+"/"+string(tier), func(t *testing.T) {
				root := t.TempDir()
				catalog := aggregateFixtureCatalog(t, root)
				skillCase := aggregateCase("example", "case/example", "routing/positive")
				skillCase.EvidenceTier = tier
				var input any = []v1alpha1.SkillCase{skillCase}
				if command == "coverage" {
					input = []v1alpha1.CoverageEntry{{
						SkillID: "example", CaseID: skillCase.ID,
						State: v1alpha1.CoverageState(tier), Metric: skillCase.Metric,
						EvidenceTier: tier, EvidenceRef: skillCase.SourceRef,
					}}
				}
				inputPath := filepath.Join(root, "input.json")
				writeCoverageFixture(t, inputPath, input)
				output := filepath.Join(root, "coverage.json")
				var stdout bytes.Buffer
				err := Run([]string{
					command, "--catalog", catalog,
					"--input", inputPath, "--output", output,
				}, &stdout)
				if err == nil || !strings.Contains(err.Error(), "configured cases only") {
					t.Fatalf("error = %v, want unsupported evidence rejection", err)
				}
				if _, err := os.Stat(output); !os.IsNotExist(err) {
					t.Fatalf("rejected input produced output: %v", err)
				}
			})
		}
	}
}

func TestCoverageRejectsSkillOutsideCatalog(t *testing.T) {
	root := t.TempDir()
	catalog := aggregateFixtureCatalog(t, root)
	inputPath := filepath.Join(root, "input.json")
	writeCoverageFixture(t, inputPath, []v1alpha1.CoverageEntry{{
		SkillID: "unknown", CaseID: "case/unknown",
		State: v1alpha1.CoverageConfigured, Metric: "routing/positive",
		EvidenceTier: v1alpha1.TierConfigured,
	}})
	var stdout bytes.Buffer
	err := runCoverage([]string{
		"--catalog", catalog, "--input", inputPath,
		"--output", filepath.Join(root, "coverage.json"),
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "absent from the bound catalog") {
		t.Fatalf("error = %v, want unknown skill rejection", err)
	}
}

func writeCoverageFixture(t *testing.T, path string, value any) {
	t.Helper()
	content, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRunAggregateRejectsDuplicateCase(t *testing.T) {
	root := t.TempDir()
	catalog := aggregateFixtureCatalog(t, root)
	cases := aggregateFixtureCases(t, root)
	duplicatedCases := []v1alpha1.SkillCase{
		aggregateCase("example", "case/example", "routing/positive"),
		aggregateCase("example", "case/example", "routing/adjacent-negative"),
	}
	duplicated, err := json.Marshal(duplicatedCases)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cases, duplicated, 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err = runAggregate([]string{
		"--catalog", catalog,
		"--input", cases,
		"--output", filepath.Join(root, "coverage.json"),
	}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "duplicate skill case") {
		t.Fatalf("runAggregate() error = %v, want duplicate rejection", err)
	}
}
