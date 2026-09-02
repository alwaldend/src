package v1alpha1

import (
	"fmt"
	"sort"
	"strings"
)

func oneOfEvidenceTier(value EvidenceTier) bool {
	switch value {
	case TierConfigured, TierRouted, TierFixtureTested, TierLive,
		TierStale, TierUnverified:
		return true
	}
	return false
}

func validateStableID(value string) error {
	if value == "" || len(value) > 128 {
		return fmt.Errorf("stable ID must contain 1 to 128 bytes")
	}
	if strings.ContainsAny(value, "\x00\r\n") {
		return fmt.Errorf("stable ID contains control characters")
	}
	return nil
}

func validateDefectSignature(value string) error {
	if value == "" || len(value) > 256 {
		return fmt.Errorf("defect signature must contain 1 to 256 bytes")
	}
	if strings.ContainsAny(value, "\x00\r\n") {
		return fmt.Errorf("defect signature contains control characters")
	}
	return nil
}

func validateSkillCase(skillCase SkillCase) error {
	if skillCase.APIVersion != APIVersion || skillCase.Kind != "SkillCase" {
		return fmt.Errorf("unsupported skill case type %q %q", skillCase.APIVersion, skillCase.Kind)
	}
	if err := validateStableID(skillCase.ID); err != nil {
		return fmt.Errorf("case ID: %w", err)
	}
	if skillCase.SkillID == "" || len(skillCase.SkillID) > 253 {
		return fmt.Errorf("case skill ID must contain 1 to 253 bytes")
	}
	if err := skillCase.Provenance.Validate(); err != nil {
		return fmt.Errorf("provenance: %w", err)
	}
	if err := skillCase.SourceRef.Validate(); err != nil {
		return fmt.Errorf("source: %w", err)
	}
	if skillCase.Metric == "" || len(skillCase.Metric) > 128 {
		return fmt.Errorf("case metric must contain 1 to 128 bytes")
	}
	if !oneOfEvidenceTier(skillCase.EvidenceTier) {
		return fmt.Errorf("unknown evidence tier %q", skillCase.EvidenceTier)
	}
	if skillCase.Digest != "" && !validDigest(skillCase.Digest) {
		return fmt.Errorf("malformed case digest")
	}
	return nil
}

func (skillCase SkillCase) Validate() error {
	return validateSkillCase(skillCase)
}

func (assertion RequirementAssertion) Validate() error {
	if assertion.APIVersion != APIVersion || assertion.Kind != "RequirementAssertion" {
		return fmt.Errorf(
			"unsupported requirement assertion type %q %q",
			assertion.APIVersion,
			assertion.Kind,
		)
	}
	if err := validateStableID(assertion.ID); err != nil {
		return fmt.Errorf("assertion ID: %w", err)
	}
	if err := assertion.CaseRef.Validate(); err != nil {
		return fmt.Errorf("case reference: %w", err)
	}
	if err := assertion.Subject.Validate(); err != nil {
		return fmt.Errorf("subject: %w", err)
	}
	if strings.TrimSpace(assertion.Revision) == "" || len(assertion.Revision) > 64 {
		return fmt.Errorf("assertion revision must contain 1 to 64 bytes")
	}
	if !oneOf(assertion.Verdict, "pass", "fail", "skipped", "blocked", "unverified") {
		return fmt.Errorf("unknown assertion verdict %q", assertion.Verdict)
	}
	if assertion.Metric == "" || len(assertion.Metric) > 128 {
		return fmt.Errorf("assertion metric must contain 1 to 128 bytes")
	}
	if !oneOfEvidenceTier(assertion.Tier) {
		return fmt.Errorf("unknown evidence tier %q", assertion.Tier)
	}
	if assertion.Digest != "" && !validDigest(assertion.Digest) {
		return fmt.Errorf("malformed assertion digest")
	}
	return nil
}

func (entry CoverageEntry) Validate() error {
	if entry.SkillID == "" || len(entry.SkillID) > 253 {
		return fmt.Errorf("coverage skill ID must contain 1 to 253 bytes")
	}
	if err := validateStableID(entry.CaseID); err != nil {
		return fmt.Errorf("coverage case ID: %w", err)
	}
	if !oneOf(
		string(entry.State),
		string(CoverageConfigured),
		string(CoverageRouted),
		string(CoverageFixtureTested),
		string(CoverageLive),
		string(CoverageStale),
		string(CoverageUnverified),
	) {
		return fmt.Errorf("unknown coverage state %q", entry.State)
	}
	if entry.Metric == "" || len(entry.Metric) > 128 {
		return fmt.Errorf("coverage metric must contain 1 to 128 bytes")
	}
	if !oneOfEvidenceTier(entry.EvidenceTier) {
		return fmt.Errorf("unknown evidence tier %q", entry.EvidenceTier)
	}
	if entry.EvidenceRef.Kind != "" {
		if err := entry.EvidenceRef.Validate(); err != nil {
			return fmt.Errorf("evidence reference: %w", err)
		}
	}
	return nil
}

func (matrix CoverageMatrix) Validate() error {
	if matrix.APIVersion != APIVersion || matrix.Kind != "CoverageMatrix" {
		return fmt.Errorf("unsupported coverage matrix type %q %q", matrix.APIVersion, matrix.Kind)
	}
	if err := validateStableID(matrix.ID); err != nil {
		return fmt.Errorf("matrix ID: %w", err)
	}
	if err := matrix.CatalogRef.Validate(); err != nil {
		return fmt.Errorf("catalog reference: %w", err)
	}
	if matrix.CatalogRef.Kind != ReferenceArtifact {
		return fmt.Errorf("coverage catalog must be an artifact reference")
	}
	if len(matrix.Entries) > 4096 {
		return fmt.Errorf("coverage matrix exceeds 4096 entries")
	}
	seen := map[string]bool{}
	for _, entry := range matrix.Entries {
		if err := entry.Validate(); err != nil {
			return err
		}
		identity := entry.SkillID + "\x00" + entry.CaseID
		if seen[identity] {
			return fmt.Errorf("duplicate normalized case %q for skill %q", entry.CaseID, entry.SkillID)
		}
		seen[identity] = true
	}
	if matrix.Total < len(matrix.Entries) || matrix.Total > 4096 {
		return fmt.Errorf("coverage total must include every emitted entry")
	}
	if matrix.Truncated && matrix.Total == len(matrix.Entries) {
		return fmt.Errorf("truncated matrix cannot claim total equals emitted")
	}
	if matrix.Digest != "" && !validDigest(matrix.Digest) {
		return fmt.Errorf("malformed matrix digest")
	}
	return nil
}

func (record FrictionRecord) Validate() error {
	if record.APIVersion != APIVersion || record.Kind != "FrictionRecord" {
		return fmt.Errorf("unsupported friction record type %q %q", record.APIVersion, record.Kind)
	}
	if err := validateStableID(record.ID); err != nil {
		return fmt.Errorf("friction ID: %w", err)
	}
	if err := record.TaskRef.Validate(); err != nil {
		return fmt.Errorf("task reference: %w", err)
	}
	if record.TaskRef.Kind != ReferenceTask && record.TaskRef.Kind != ReferenceGoal {
		return fmt.Errorf("friction task must be a task or goal reference")
	}
	if err := record.EvidenceRef.Validate(); err != nil {
		return fmt.Errorf("evidence reference: %w", err)
	}
	if record.EvidenceRef.Kind != ReferenceArtifact {
		return fmt.Errorf("friction evidence must be an artifact reference")
	}
	if !oneOfEvidenceTier(record.EvidenceTier) {
		return fmt.Errorf("unknown evidence tier %q", record.EvidenceTier)
	}
	if err := validateDefectSignature(record.DefectSignature); err != nil {
		return fmt.Errorf("defect signature: %w", err)
	}
	if record.AvoidableReads < 0 || record.AvoidableCommands < 0 || record.LatencyMS < 0 {
		return fmt.Errorf("friction metrics cannot be negative")
	}
	if len(record.SelectedRefs) > 64 || len(record.ConsideredRefs) > 64 ||
		len(record.Conflicts) > 64 || len(record.MissingProviders) > 64 ||
		len(record.FailedAssumptions) > 64 {
		return fmt.Errorf("friction record exceeds 64 bounded references or observations")
	}
	totalObservationBytes := 0
	for _, observation := range [][]string{
		record.Conflicts,
		record.MissingProviders,
		record.FailedAssumptions,
	} {
		for _, value := range observation {
			if len(value) > 512 {
				return fmt.Errorf(
					"friction observation exceeds 512 bytes per value",
				)
			}
			totalObservationBytes += len(value)
		}
	}
	if totalObservationBytes > 8192 {
		return fmt.Errorf(
			"friction record exceeds 8192 bytes of total observations",
		)
	}
	for _, reference := range record.SelectedRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("selected reference: %w", err)
		}
	}
	for _, reference := range record.ConsideredRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("considered reference: %w", err)
		}
	}
	return nil
}

func (proposal LearningProposal) Validate() error {
	if proposal.APIVersion != APIVersion || proposal.Kind != "LearningProposal" {
		return fmt.Errorf("unsupported learning proposal type %q %q", proposal.APIVersion, proposal.Kind)
	}
	if err := validateStableID(proposal.ID); err != nil {
		return fmt.Errorf("proposal ID: %w", err)
	}
	if proposal.Owner == "" || len(proposal.Owner) > 253 {
		return fmt.Errorf("proposal owner must contain 1 to 253 bytes")
	}
	if err := validateDefectSignature(proposal.DefectSignature); err != nil {
		return fmt.Errorf("defect signature: %w", err)
	}
	if err := proposal.Reproducer.Validate(); err != nil {
		return fmt.Errorf("reproducer: %w", err)
	}
	if err := proposal.RegressionRef.Validate(); err != nil {
		return fmt.Errorf("regression: %w", err)
	}
	if err := proposal.ContractRef.Validate(); err != nil {
		return fmt.Errorf("contract: %w", err)
	}
	if strings.TrimSpace(proposal.Fallback) == "" || len(proposal.Fallback) > 1024 {
		return fmt.Errorf("proposal fallback must contain 1 to 1024 bytes")
	}
	if err := proposal.ResourceBudget.Validate(); err != nil {
		return fmt.Errorf("resource budget: %w", err)
	}
	if len(proposal.ValidationRefs) == 0 {
		return fmt.Errorf("proposal requires at least one validation reference")
	}
	for _, reference := range proposal.ValidationRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("validation reference: %w", err)
		}
	}
	if strings.TrimSpace(proposal.RetirementRule) == "" || len(proposal.RetirementRule) > 1024 {
		return fmt.Errorf("proposal retirement rule must contain 1 to 1024 bytes")
	}
	if len(proposal.FrictionRefs) < 2 {
		return fmt.Errorf("repeated proposal requires at least two friction references")
	}
	for _, reference := range proposal.FrictionRefs {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("friction reference: %w", err)
		}
		if reference.Kind != ReferenceArtifact {
			return fmt.Errorf("friction reference must be an artifact reference")
		}
	}
	seen := make(map[string]bool, len(proposal.FrictionRefs))
	for _, reference := range proposal.FrictionRefs {
		identity := string(reference.Kind) + "\x00" + reference.ID +
			"\x00" + reference.Digest
		if seen[identity] {
			return fmt.Errorf("friction references must be distinct")
		}
		seen[identity] = true
	}
	if proposal.Digest != "" && !validDigest(proposal.Digest) {
		return fmt.Errorf("malformed proposal digest")
	}
	return nil
}

func (measurement ContextMeasurement) Validate() error {
	if measurement.APIVersion != APIVersion || measurement.Kind != "ContextMeasurement" {
		return fmt.Errorf("unsupported context measurement type %q %q", measurement.APIVersion, measurement.Kind)
	}
	if err := validateStableID(measurement.ID); err != nil {
		return fmt.Errorf("measurement ID: %w", err)
	}
	if measurement.Profile == "" || len(measurement.Profile) > 128 {
		return fmt.Errorf("measurement profile must contain 1 to 128 bytes")
	}
	if measurement.ColdStartBytes < 0 || measurement.SessionBytes < 0 ||
		measurement.ResumeSteps < 0 || measurement.RepeatedChecks < 0 ||
		measurement.UnsupportedClaims < 0 || measurement.UnsafeActions < 0 {
		return fmt.Errorf("measurement values cannot be negative")
	}
	if !oneOfEvidenceTier(measurement.EvidenceTier) {
		return fmt.Errorf("unknown evidence tier %q", measurement.EvidenceTier)
	}
	if err := measurement.ProducerRef.Validate(); err != nil {
		return fmt.Errorf("producer: %w", err)
	}
	if err := measurement.SubjectRef.Validate(); err != nil {
		return fmt.Errorf("subject: %w", err)
	}
	if measurement.Digest != "" && !validDigest(measurement.Digest) {
		return fmt.Errorf("malformed measurement digest")
	}
	return nil
}

func (budget DisclosureBudget) Validate() error {
	if budget.APIVersion != APIVersion || budget.Kind != "DisclosureBudget" {
		return fmt.Errorf("unsupported disclosure budget type %q %q", budget.APIVersion, budget.Kind)
	}
	if err := validateStableID(budget.ID); err != nil {
		return fmt.Errorf("budget ID: %w", err)
	}
	if budget.CatalogSummaryBytes <= 0 || budget.ContentDigestBytes <= 0 ||
		budget.MaxBodyBytes <= 0 || budget.MaxSchemaBytes <= 0 {
		return fmt.Errorf("disclosure budgets must be positive")
	}
	if budget.MaxResumeSteps < 0 || budget.MaxRepeatedChecks < 0 {
		return fmt.Errorf("resume and check budgets cannot be negative")
	}
	if len(budget.Provenance) == 0 {
		return fmt.Errorf("disclosure budget requires provenance")
	}
	for _, reference := range budget.Provenance {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("provenance: %w", err)
		}
	}
	if !oneOfEvidenceTier(budget.EvidenceTier) {
		return fmt.Errorf("unknown evidence tier %q", budget.EvidenceTier)
	}
	if budget.Digest != "" && !validDigest(budget.Digest) {
		return fmt.Errorf("malformed budget digest")
	}
	return nil
}

func normalizeEvidenceReferences(references []Reference) []Reference {
	normalized := append([]Reference(nil), references...)
	sort.Slice(normalized, func(left int, right int) bool {
		if normalized[left].Kind != normalized[right].Kind {
			return normalized[left].Kind < normalized[right].Kind
		}
		return normalized[left].ID < normalized[right].ID
	})
	return normalized
}
