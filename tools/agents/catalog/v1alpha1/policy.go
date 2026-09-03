package catalogv1alpha1

import (
	"encoding/json"
	"fmt"
)

const (
	// KindPolicyCatalog identifies the policy catalog kind.
	KindPolicyCatalog = "policy-catalog"
)

// PolicyAxis presents one independently derived policy axis. Every value
// carries its owning source path; no axis is inferred from another.
type PolicyAxis struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

// PolicyRecord is the portable policy projection for one path prefix.
type PolicyRecord struct {
	ID                string       `json:"id"`
	PathPrefix        string       `json:"pathPrefix"`
	Precedence        int          `json:"precedence"`
	AgentPolicySource string       `json:"agentPolicySource,omitempty"`
	OwnerBoundaryRef  string       `json:"ownerBoundaryRef,omitempty"`
	ReviewSource      string       `json:"reviewSource,omitempty"`
	Axes              []PolicyAxis `json:"axes"`
}

// PolicyCatalog is the bounded policy projection over tracked agent policy
// sources. Axes are independent and conflicts are reported, never resolved.
type PolicyCatalog struct {
	CatalogEnvelope
	Policies []PolicyRecord `json:"policies"`
}

// Validate checks policy-specific invariants and the shared envelope.
func (catalog PolicyCatalog) Validate() error {
	if catalog.Policies == nil {
		return fmt.Errorf("policies must be a non-null array")
	}
	itemIDs := make([]string, 0, len(catalog.Policies))
	seenPrefix := map[string]bool{}
	axisNames := map[string]bool{}
	for _, policy := range catalog.Policies {
		if policy.ID == "" || policy.PathPrefix == "" {
			return fmt.Errorf("policy record %q lacks identity or path prefix", policy.ID)
		}
		if seenPrefix[policy.PathPrefix] {
			return fmt.Errorf("duplicate path prefix %q", policy.PathPrefix)
		}
		seenPrefix[policy.PathPrefix] = true
		if policy.Precedence < 0 {
			return fmt.Errorf("policy %s has negative precedence", policy.ID)
		}
		if len(policy.Axes) == 0 {
			return fmt.Errorf("policy %s has no axes", policy.ID)
		}
		seen := map[string]bool{}
		for _, axis := range policy.Axes {
			if axis.Name == "" || axis.Value == "" || axis.Source == "" {
				return fmt.Errorf("policy %s has an incomplete axis", policy.ID)
			}
			if !oneOf(axis.Value, "known", "unknown", "conflict") {
				return fmt.Errorf("policy %s axis %s has unknown value %q",
					policy.ID, axis.Name, axis.Value)
			}
			if seen[axis.Name] {
				return fmt.Errorf("policy %s has duplicate axis %s", policy.ID, axis.Name)
			}
			seen[axis.Name] = true
			axisNames[axis.Name] = true
		}
		for _, axis := range policy.Axes {
			if axis.Value != "known" {
				continue
			}
			if axis.Source == "" || absPath(axis.Source) || escapes(axis.Source) {
				return fmt.Errorf("policy %s axis %s has malformed source %q",
					policy.ID, axis.Name, axis.Source)
			}
		}
		itemIDs = append(itemIDs, "policy."+policy.ID)
	}
	if len(axisNames) == 0 {
		return fmt.Errorf("policy catalog must define at least one axis")
	}
	return catalog.CatalogEnvelope.Validate(itemIDs)
}

// ItemIDs returns the sorted set-like identity list for a policy catalog.
func (catalog PolicyCatalog) ItemIDs() []string {
	ids := make([]string, 0, len(catalog.Policies))
	for _, policy := range catalog.Policies {
		ids = append(ids, "policy."+policy.ID)
	}
	return ids
}

// VerifySelfDigest recomputes the content-addressed digest over the canonical
// bytes of the complete policy document and compares it to the stored digest,
// rejecting an empty or mismatched digest.
func (catalog PolicyCatalog) VerifySelfDigest() error {
	withoutDigest := catalog
	withoutDigest.Digest = ""
	return verifySelfDigest(catalog.Digest, withoutDigest)
}

// CanonicalJSONPolicy encodes a complete policy catalog with a digest over
// the canonical bytes with the digest field omitted.
func CanonicalJSONPolicy(catalog PolicyCatalog) ([]byte, error) {
	if err := catalog.Validate(); err != nil {
		return nil, err
	}
	withoutDigest := catalog
	withoutDigest.Digest = ""
	content, err := json.Marshal(withoutDigest)
	if err != nil {
		return nil, fmt.Errorf("encode policy catalog: %w", err)
	}
	withoutDigest.Digest = digest(content)
	contentWithDigest, err := json.MarshalIndent(withoutDigest, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("encode policy catalog with digest: %w", err)
	}
	return append(contentWithDigest, '\n'), nil
}

// DecodePolicyStrict decodes and validates a complete policy catalog
// document with unknown-field and trailing-JSON rejection.
func DecodePolicyStrict(content []byte) (PolicyCatalog, error) {
	var catalog PolicyCatalog
	if err := DecodeStrict(content, &catalog); err != nil {
		return PolicyCatalog{}, err
	}
	if err := catalog.Validate(); err != nil {
		return PolicyCatalog{}, err
	}
	if err := catalog.VerifySelfDigest(); err != nil {
		return PolicyCatalog{}, err
	}
	return catalog, nil
}
