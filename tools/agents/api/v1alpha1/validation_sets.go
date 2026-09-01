// Phase 3C: candidate-bound reusable evidence.
package v1alpha1

// ValidationResult is one check's outcome inside a validation set.
type ValidationResult struct {
	CheckID    string `json:"checkId"`
	Status     string `json:"status"`
	DurationMS int64  `json:"durationMs,omitempty"`
	OutputRef  string `json:"outputRef,omitempty"`
}

// ValidationInputs bind the exact inputs a validation set ran against.
type ValidationInputs struct {
	BaseOID          string            `json:"baseOid,omitempty"`
	TreeOID          string            `json:"treeOid,omitempty"`
	CommitOID        string            `json:"commitOid,omitempty"`
	ProviderRefs     []Reference       `json:"providerRefs,omitempty"`
	ConfigDigests    map[string]string `json:"configDigests,omitempty"`
	ToolchainDigests map[string]string `json:"toolchainDigests,omitempty"`
	PolicyDigests    map[string]string `json:"policyDigests,omitempty"`
}

// ValidationSet is an immutable, candidate-bound, criterion-neutral record of
// one executed validation profile. It never contains credential argv.
type ValidationSet struct {
	APIVersion      string             `json:"apiVersion"`
	Kind            string             `json:"kind"`
	ID              string             `json:"id"`
	Profile         ImpactProfile      `json:"profile"`
	Candidate       Reference          `json:"candidate"`
	Inputs          ValidationInputs   `json:"inputs"`
	SanitizedArgs   []string           `json:"sanitizedArgs"`
	WorkingScope    string             `json:"workingScope"`
	Results         []ValidationResult `json:"results"`
	Coverage        []string           `json:"coverage,omitempty"`
	TotalDurationMS int64              `json:"totalDurationMs,omitempty"`
	OutputBounds    Budget             `json:"outputBounds"`
	Limitations     []string           `json:"limitations,omitempty"`
	RawLogRefs      []Reference        `json:"rawLogRefs,omitempty"`
	CleanPreState   bool               `json:"cleanPreState"`
	CleanPostState  bool               `json:"cleanPostState"`
	Digest          string             `json:"digest,omitempty"`
}

// EvidenceAssertion applies one or more validation sets to an exact criterion
// revision with a semantic verdict. It is created by the goal or task owner.
type EvidenceAssertion struct {
	APIVersion     string      `json:"apiVersion"`
	Kind           string      `json:"kind"`
	ID             string      `json:"id"`
	CriterionRef   Reference   `json:"criterionRef"`
	CriterionRev   string      `json:"criterionRev"`
	ValidationRefs []Reference `json:"validationRefs"`
	Verdict        string      `json:"verdict"`
	Digest         string      `json:"digest,omitempty"`
}
