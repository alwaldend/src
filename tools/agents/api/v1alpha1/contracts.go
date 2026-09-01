// Phase 3B: action admission at provider gateways.
package v1alpha1

// AdmissionRequest is the exact authority, subject, environment, and pre-state
// context one action admission is evaluated against.
type AdmissionRequest struct {
	APIVersion       string            `json:"apiVersion"`
	Kind             string            `json:"kind"`
	Operation        Reference         `json:"operation"`
	Authority        AuthorityEnvelope `json:"authority"`
	SubjectRefs      []Reference       `json:"subjectRefs,omitempty"`
	Environment      []string          `json:"environment,omitempty"`
	PreStateDigest   string            `json:"preStateDigest,omitempty"`
	ExpectedPreState string            `json:"expectedPreState,omitempty"`
	Budget           Budget            `json:"budget"`
	RemoteWrite      bool              `json:"remoteWrite,omitempty"`
	Destroy          bool              `json:"destroy,omitempty"`
	PrepareVerified  bool              `json:"prepareVerified,omitempty"`
}

// AdmissionDecision is one fail-closed admission verdict.
type AdmissionDecision struct {
	APIVersion   string   `json:"apiVersion"`
	Kind         string   `json:"kind"`
	Admitted     bool     `json:"admitted"`
	ReasonCode   string   `json:"reasonCode,omitempty"`
	ReasonDetail []string `json:"reasonDetail,omitempty"`
	EffectSet    []Effect `json:"effectSet,omitempty"`
	ValidatedAt  string   `json:"validatedAt,omitempty"`
}
