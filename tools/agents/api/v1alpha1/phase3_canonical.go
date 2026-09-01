package v1alpha1

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func DecodeImpactPlan(content []byte) (ImpactPlan, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var plan ImpactPlan
	if err := decoder.Decode(&plan); err != nil {
		return ImpactPlan{}, fmt.Errorf("decode impact plan: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return ImpactPlan{}, err
	}
	if err := plan.Validate(); err != nil {
		return ImpactPlan{}, err
	}
	return plan, nil
}

func CanonicalImpactPlanJSON(plan ImpactPlan) ([]byte, error) {
	if err := plan.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := plan
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode impact plan: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeValidationSet(content []byte) (ValidationSet, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var set ValidationSet
	if err := decoder.Decode(&set); err != nil {
		return ValidationSet{}, fmt.Errorf("decode validation set: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return ValidationSet{}, err
	}
	if err := set.Validate(); err != nil {
		return ValidationSet{}, err
	}
	return set, nil
}

func CanonicalValidationSetJSON(set ValidationSet) ([]byte, error) {
	if err := set.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := set
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode validation set: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DecodeEvidenceAssertion(content []byte) (EvidenceAssertion, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var assertion EvidenceAssertion
	if err := decoder.Decode(&assertion); err != nil {
		return EvidenceAssertion{}, fmt.Errorf("decode evidence assertion: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return EvidenceAssertion{}, err
	}
	if err := assertion.Validate(); err != nil {
		return EvidenceAssertion{}, err
	}
	return assertion, nil
}

func CanonicalEvidenceAssertionJSON(assertion EvidenceAssertion) ([]byte, error) {
	if err := assertion.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := assertion
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode evidence assertion: %w", err)
	}
	withoutDigest.Digest = DigestOfCanonicalJSON(content)
	return json.Marshal(withoutDigest)
}

func DigestOfCanonicalJSON(content []byte) string {
	sum := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func CanonicalAdmissionJSON(request AdmissionRequest) ([]byte, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	content, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode admission request: %w", err)
	}
	return append(content, '\n'), nil
}

func DecodeAdmissionDecision(content []byte) (AdmissionDecision, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var decision AdmissionDecision
	if err := decoder.Decode(&decision); err != nil {
		return AdmissionDecision{}, fmt.Errorf("decode admission decision: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return AdmissionDecision{}, err
	}
	if err := decision.Validate(); err != nil {
		return AdmissionDecision{}, err
	}
	return decision, nil
}
