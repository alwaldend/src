package v1alpha1

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func decodePhase5Object(content []byte, name string, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode %s: %w", name, err)
	}
	return requireJSONEOF(decoder)
}

func CanonicalSkillCaseJSON(skillCase SkillCase) ([]byte, error) {
	if err := skillCase.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := skillCase
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode skill case: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeSkillCase(content []byte) (SkillCase, error) {
	var skillCase SkillCase
	if err := decodePhase5Object(content, "skill case", &skillCase); err != nil {
		return SkillCase{}, err
	}
	if err := skillCase.Validate(); err != nil {
		return SkillCase{}, err
	}
	return skillCase, nil
}

func CanonicalRequirementAssertionJSON(assertion RequirementAssertion) ([]byte, error) {
	if err := assertion.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := assertion
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode requirement assertion: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeRequirementAssertion(content []byte) (RequirementAssertion, error) {
	var assertion RequirementAssertion
	if err := decodePhase5Object(content, "requirement assertion", &assertion); err != nil {
		return RequirementAssertion{}, err
	}
	if err := assertion.Validate(); err != nil {
		return RequirementAssertion{}, err
	}
	return assertion, nil
}

func CanonicalCoverageMatrixJSON(matrix CoverageMatrix) ([]byte, error) {
	if err := matrix.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := matrix
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode coverage matrix: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeCoverageMatrix(content []byte) (CoverageMatrix, error) {
	var matrix CoverageMatrix
	if err := decodePhase5Object(content, "coverage matrix", &matrix); err != nil {
		return CoverageMatrix{}, err
	}
	if err := matrix.Validate(); err != nil {
		return CoverageMatrix{}, err
	}
	return matrix, nil
}

func CanonicalFrictionRecordJSON(record FrictionRecord) ([]byte, error) {
	if err := record.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(record)
}

func DecodeFrictionRecord(content []byte) (FrictionRecord, error) {
	var record FrictionRecord
	if err := decodePhase5Object(content, "friction record", &record); err != nil {
		return FrictionRecord{}, err
	}
	if err := record.Validate(); err != nil {
		return FrictionRecord{}, err
	}
	return record, nil
}

func CanonicalLearningProposalJSON(proposal LearningProposal) ([]byte, error) {
	if err := proposal.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := proposal
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode learning proposal: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeLearningProposal(content []byte) (LearningProposal, error) {
	var proposal LearningProposal
	if err := decodePhase5Object(content, "learning proposal", &proposal); err != nil {
		return LearningProposal{}, err
	}
	if err := proposal.Validate(); err != nil {
		return LearningProposal{}, err
	}
	return proposal, nil
}

func CanonicalContextMeasurementJSON(measurement ContextMeasurement) ([]byte, error) {
	if err := measurement.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := measurement
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode context measurement: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeContextMeasurement(content []byte) (ContextMeasurement, error) {
	var measurement ContextMeasurement
	if err := decodePhase5Object(content, "context measurement", &measurement); err != nil {
		return ContextMeasurement{}, err
	}
	if err := measurement.Validate(); err != nil {
		return ContextMeasurement{}, err
	}
	return measurement, nil
}

func CanonicalDisclosureBudgetJSON(budget DisclosureBudget) ([]byte, error) {
	if err := budget.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := budget
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode disclosure budget: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeDisclosureBudget(content []byte) (DisclosureBudget, error) {
	var budget DisclosureBudget
	if err := decodePhase5Object(content, "disclosure budget", &budget); err != nil {
		return DisclosureBudget{}, err
	}
	if err := budget.Validate(); err != nil {
		return DisclosureBudget{}, err
	}
	return budget, nil
}
