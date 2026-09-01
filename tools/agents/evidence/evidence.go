// Package evidence emits immutable, candidate-bound validation evidence.
//
// A ValidationSet is a deterministic record of one executed validation
// profile: its digest covers every field above the Digest field, so the same
// inputs always produce the same identity. Sets are never mutated after
// emission; EvidenceAssertions bind sets by reference and may share one set.
package evidence

import (
	"encoding/json"
	"fmt"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

// EmitOptions carries the caller-supplied inputs of a validation set.
type EmitOptions struct {
	Profile          v1alpha1.ImpactProfile
	Candidate        v1alpha1.Reference
	BaseOID          string
	TreeOID          string
	CommitOID        string
	SanitizedArgs    []string
	WorkingScope     string
	ProviderRefs     []v1alpha1.Reference
	ConfigDigests    map[string]string
	ToolchainDigests map[string]string
	PolicyDigests    map[string]string
	OutputBounds     v1alpha1.Budget
	RawLogs          []v1alpha1.Reference
	TotalDurationMS  int64
	CleanPreState    bool
	CleanPostState   bool
}

// EmitValidationSet builds an immutable validation set with a deterministic
// identity and digest, delegating admission checks to the shared v1alpha1
// validators (including credential-like argument rejection).
func EmitValidationSet(
	opts EmitOptions,
	results []v1alpha1.ValidationResult,
) (v1alpha1.ValidationSet, error) {
	set := v1alpha1.ValidationSet{
		APIVersion: v1alpha1.APIVersion,
		Kind:       "ValidationSet",
		ID:         "validation." + scopeSlug(opts.WorkingScope),
		Profile:    opts.Profile,
		Candidate:  opts.Candidate,
		Inputs: v1alpha1.ValidationInputs{
			BaseOID:          opts.BaseOID,
			TreeOID:          opts.TreeOID,
			CommitOID:        opts.CommitOID,
			ProviderRefs:     append([]v1alpha1.Reference(nil), opts.ProviderRefs...),
			ConfigDigests:    cloneStringMap(opts.ConfigDigests),
			ToolchainDigests: cloneStringMap(opts.ToolchainDigests),
			PolicyDigests:    cloneStringMap(opts.PolicyDigests),
		},
		SanitizedArgs:   append([]string(nil), opts.SanitizedArgs...),
		WorkingScope:    opts.WorkingScope,
		Results:         append([]v1alpha1.ValidationResult(nil), results...),
		OutputBounds:    opts.OutputBounds,
		RawLogRefs:      append([]v1alpha1.Reference(nil), opts.RawLogs...),
		TotalDurationMS: opts.TotalDurationMS,
		CleanPreState:   opts.CleanPreState,
		CleanPostState:  opts.CleanPostState,
	}
	// Resolve the final identity first, then let the canonical encoder compute
	// the content-addressed digest (over the digest-omitted bytes) and emit
	// the completed document. Decode that document to obtain the canonical
	// digest-bearing validation set.
	provisional, err := v1alpha1.CanonicalValidationSetJSON(set)
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	var provisionalSet v1alpha1.ValidationSet
	if err := json.Unmarshal(provisional, &provisionalSet); err != nil {
		return v1alpha1.ValidationSet{}, fmt.Errorf("decode provisional set: %w", err)
	}
	set.ID = "validation." + scopeSlug(opts.WorkingScope) + "." +
		shortDigest(provisionalSet.Digest)
	canonical, err := v1alpha1.CanonicalValidationSetJSON(set)
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	var finalSet v1alpha1.ValidationSet
	if err := json.Unmarshal(canonical, &finalSet); err != nil {
		return v1alpha1.ValidationSet{}, fmt.Errorf("decode final set: %w", err)
	}
	return finalSet, nil
}

// NewAssertion binds the given sets' digests into artifact references and
// returns a deterministic evidence assertion for one criterion revision.
func NewAssertion(
	id string,
	criterionRef v1alpha1.Reference,
	criterionRev string,
	verdict string,
	sets ...v1alpha1.ValidationSet,
) (v1alpha1.EvidenceAssertion, error) {
	if len(sets) == 0 {
		return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
			"an assertion must bind at least one validation set",
		)
	}
	refs := make([]v1alpha1.Reference, 0, len(sets))
	seen := map[string]bool{}
	for _, set := range sets {
		if err := set.Validate(); err != nil {
			return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
				"validation set %q: %w",
				set.ID,
				err,
			)
		}
		if set.Digest == "" {
			return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
				"validation set %q has no digest",
				set.ID,
			)
		}
		if seen[set.ID] {
			return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
				"duplicate validation set %q",
				set.ID,
			)
		}
		seen[set.ID] = true
		refs = append(refs, v1alpha1.Reference{
			Kind:   v1alpha1.ReferenceArtifact,
			ID:     set.ID,
			Digest: set.Digest,
		})
	}
	assertion := v1alpha1.EvidenceAssertion{
		APIVersion:     v1alpha1.APIVersion,
		Kind:           "EvidenceAssertion",
		ID:             id,
		CriterionRef:   criterionRef,
		CriterionRev:   criterionRev,
		ValidationRefs: refs,
		Verdict:        verdict,
	}
	canonical, err := v1alpha1.CanonicalEvidenceAssertionJSON(assertion)
	if err != nil {
		return v1alpha1.EvidenceAssertion{}, err
	}
	assertion.Digest = v1alpha1.DigestOfCanonicalJSON(canonical)
	return assertion, nil
}

// Apply checks that every validation reference of the assertion is bound in
// the set (by set ID, or through InputRefs) and that bound set digests match
// the set's own digest. The immutable set is never mutated; any mismatch is
// returned as an error and one set may support many assertions.
func Apply(
	assertion v1alpha1.EvidenceAssertion,
	set v1alpha1.ValidationSet,
) error {
	if err := assertion.Validate(); err != nil {
		return err
	}
	if err := set.Validate(); err != nil {
		return err
	}
	for _, ref := range assertion.ValidationRefs {
		switch {
		case ref.ID == set.ID:
			if ref.Digest != "" && ref.Digest != set.Digest {
				return fmt.Errorf(
					"validation reference %q digest %s does not match set "+
						"digest %s",
					ref.ID,
					ref.Digest,
					set.Digest,
				)
			}
		case boundViaInputRefs(ref, set.Inputs.ProviderRefs):
			if ref.Digest != "" && !inputRefDigestMatches(
				ref,
				set.Inputs.ProviderRefs,
			) {
				return fmt.Errorf(
					"validation reference %q is bound via inputs but its "+
						"digest %s does not match any provider digest",
					ref.ID,
					ref.Digest,
				)
			}
		default:
			return fmt.Errorf(
				"validation reference %q is not bound by validation set %q",
				ref.ID,
				set.ID,
			)
		}
	}
	return nil
}

// Applicability resolves the reuse decision for a set against current inputs.
//
// A message-only commit rewrite keeps the same tree, so tree-bound evidence
// stays applicable; commit- and base-bound evidence goes stale; fully bound
// evidence goes stale when any binding mismatches; and evidence whose inputs
// are missing is unknown because it cannot be tied to the current state.
func Applicability(
	set v1alpha1.ValidationSet,
	inputs v1alpha1.ValidationInputs,
	currentBaseOID string,
	currentTreeOID string,
	currentCommitOID string,
) v1alpha1.EvidenceApplicability {
	if !inputsEqual(set.Inputs, inputs) {
		return v1alpha1.EvidenceStale
	}
	if set.Inputs.BaseOID == "" && set.Inputs.TreeOID == "" &&
		set.Inputs.CommitOID == "" {
		// No base/tree/commit binding leaves nothing to reuse.
		return v1alpha1.EvidenceUnknown
	}
	// Any bound OID that no longer matches the current state makes the
	// evidence stale; a message-only commit rewrite keeps the tree (and the
	// base), so tree- or base-bound evidence stays applicable.
	if set.Inputs.BaseOID != "" && set.Inputs.BaseOID != currentBaseOID {
		return v1alpha1.EvidenceStale
	}
	if set.Inputs.TreeOID != "" && set.Inputs.TreeOID != currentTreeOID {
		return v1alpha1.EvidenceStale
	}
	if set.Inputs.CommitOID != "" && set.Inputs.CommitOID != currentCommitOID {
		return v1alpha1.EvidenceStale
	}
	return v1alpha1.EvidenceApplicable
}

func inputsEqual(
	bound v1alpha1.ValidationInputs,
	current v1alpha1.ValidationInputs,
) bool {
	if bound.BaseOID != current.BaseOID || bound.TreeOID != current.TreeOID ||
		bound.CommitOID != current.CommitOID {
		return false
	}
	if !stringMapEqual(bound.ConfigDigests, current.ConfigDigests) ||
		!stringMapEqual(bound.ToolchainDigests, current.ToolchainDigests) ||
		!stringMapEqual(bound.PolicyDigests, current.PolicyDigests) {
		return false
	}
	return true
}

func scopeSlug(workingScope string) string {
	slug := strings.TrimSpace(workingScope)
	for _, mapper := range []struct {
		from, to string
	}{
		{"/", "."},
		{"\\", "."},
		{"_", "-"},
		{" ", "-"},
	} {
		slug = strings.ReplaceAll(slug, mapper.from, mapper.to)
	}
	slug = strings.Trim(slug, ".-")
	if slug == "" {
		return "scope"
	}
	return slug
}

func shortDigest(digest string) string {
	trimmed := digest
	if index := strings.Index(trimmed, ":"); index >= 0 {
		trimmed = trimmed[index+1:]
	}
	if len(trimmed) < 12 {
		return trimmed
	}
	return trimmed[:12]
}

func cloneStringMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	clone := make(map[string]string, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}

func stringMapEqual(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}

func boundViaInputRefs(ref v1alpha1.Reference, providers []v1alpha1.Reference) bool {
	for _, provider := range providers {
		if provider.ID == ref.ID {
			return true
		}
	}
	return false
}

func inputRefDigestMatches(
	ref v1alpha1.Reference,
	providers []v1alpha1.Reference,
) bool {
	for _, provider := range providers {
		if provider.ID == ref.ID && provider.Digest != "" &&
			provider.Digest == ref.Digest {
			return true
		}
	}
	return false
}
