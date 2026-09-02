// Phase 5: measured learning and optimization.
package v1alpha1

// EvidenceTier distinguishes how behavior was observed.
type EvidenceTier string

const (
	TierConfigured    EvidenceTier = "configured"
	TierRouted        EvidenceTier = "routed"
	TierFixtureTested EvidenceTier = "fixture-tested"
	TierLive          EvidenceTier = "live"
	TierStale         EvidenceTier = "stale"
	TierUnverified    EvidenceTier = "unverified"
)

// CoverageState names why a behavior has not been normalized.
type CoverageState string

const (
	CoverageConfigured    CoverageState = "configured"
	CoverageRouted        CoverageState = "routed"
	CoverageFixtureTested CoverageState = "fixture-tested"
	CoverageLive          CoverageState = "live"
	CoverageStale         CoverageState = "stale"
	CoverageUnverified    CoverageState = "unverified"
)

// SkillCase is a stable, normalized skill-evaluation case. It names its
// provenance and intended metric without embedding prompt or model output.
type SkillCase struct {
	APIVersion   string       `json:"apiVersion"`
	Kind         string       `json:"kind"`
	ID           string       `json:"id"`
	SkillID      string       `json:"skillId"`
	Provenance   Reference    `json:"provenance"`
	Metric       string       `json:"metric"`
	EvidenceTier EvidenceTier `json:"evidenceTier"`
	SourceRef    Reference    `json:"sourceRef"`
	Digest       string       `json:"digest,omitempty"`
}

// RequirementAssertion binds one normalized case to an exact requirement
// revision and verdict.
type RequirementAssertion struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	ID         string       `json:"id"`
	CaseRef    Reference    `json:"caseRef"`
	Subject    Reference    `json:"subject"`
	Revision   string       `json:"revision"`
	Verdict    string       `json:"verdict"`
	Metric     string       `json:"metric"`
	Tier       EvidenceTier `json:"tier"`
	Digest     string       `json:"digest,omitempty"`
}

// CoverageEntry classifies one skill behavior by its strongest evidence.
type CoverageEntry struct {
	SkillID      string        `json:"skillId"`
	CaseID       string        `json:"caseId"`
	State        CoverageState `json:"state"`
	Metric       string        `json:"metric"`
	EvidenceTier EvidenceTier  `json:"evidenceTier"`
	EvidenceRef  Reference     `json:"evidenceRef,omitempty"`
}

// CoverageMatrix is the bounded published projection of normalized behavior.
type CoverageMatrix struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	ID         string          `json:"id"`
	CatalogRef Reference       `json:"catalogRef"`
	Entries    []CoverageEntry `json:"entries"`
	Total      int             `json:"total"`
	Truncated  bool            `json:"truncated"`
	Digest     string          `json:"digest,omitempty"`
}

// FrictionRecord is one bounded sanitized observation at task close.
type FrictionRecord struct {
	APIVersion        string       `json:"apiVersion"`
	Kind              string       `json:"kind"`
	ID                string       `json:"id"`
	TaskRef           Reference    `json:"taskRef"`
	SelectedRefs      []Reference  `json:"selectedRefs,omitempty"`
	ConsideredRefs    []Reference  `json:"consideredRefs,omitempty"`
	Conflicts         []string     `json:"conflicts,omitempty"`
	MissingProviders  []string     `json:"missingProviders,omitempty"`
	AvoidableReads    int          `json:"avoidableReads,omitempty"`
	AvoidableCommands int          `json:"avoidableCommands,omitempty"`
	FailedAssumptions []string     `json:"failedAssumptions,omitempty"`
	LatencyMS         int64        `json:"latencyMs,omitempty"`
	EvidenceRef       Reference    `json:"evidenceRef"`
	EvidenceTier      EvidenceTier `json:"evidenceTier"`
	DefectSignature   string       `json:"defectSignature"`
}

// LearningProposal is the only route from repeated friction to a reviewed
// change. It never authorizes an automatic edit.
type LearningProposal struct {
	APIVersion        string      `json:"apiVersion"`
	Kind              string      `json:"kind"`
	ID                string      `json:"id"`
	Owner             string      `json:"owner"`
	DefectSignature   string      `json:"defectSignature"`
	Reproducer        Reference   `json:"reproducer"`
	RegressionRef     Reference   `json:"regressionRef"`
	ContractRef       Reference   `json:"contractRef"`
	Fallback          string      `json:"fallback"`
	ResourceBudget    Budget      `json:"resourceBudget"`
	ValidationRefs    []Reference `json:"validationRefs"`
	DeliveredRevision string      `json:"deliveredRevision,omitempty"`
	RetirementRule    string      `json:"retirementRule"`
	FrictionRefs      []Reference `json:"frictionRefs"`
	Limitations       []string    `json:"limitations,omitempty"`
	Digest            string      `json:"digest,omitempty"`
}

// ContextMeasurement captures one bounded optimization observation.
type ContextMeasurement struct {
	APIVersion        string       `json:"apiVersion"`
	Kind              string       `json:"kind"`
	ID                string       `json:"id"`
	Profile           string       `json:"profile"`
	ColdStartBytes    int64        `json:"coldStartBytes"`
	SessionBytes      int64        `json:"sessionBytes"`
	ResumeSteps       int          `json:"resumeSteps"`
	RepeatedChecks    int          `json:"repeatedChecks"`
	UnsupportedClaims int          `json:"unsupportedClaims"`
	UnsafeActions     int          `json:"unsafeActions"`
	EvidenceTier      EvidenceTier `json:"evidenceTier"`
	ProducerRef       Reference    `json:"producerRef"`
	SubjectRef        Reference    `json:"subjectRef"`
	Digest            string       `json:"digest,omitempty"`
}

// DisclosureBudget binds a measured context budget to exact input identities.
type DisclosureBudget struct {
	APIVersion          string       `json:"apiVersion"`
	Kind                string       `json:"kind"`
	ID                  string       `json:"id"`
	CatalogSummaryBytes int64        `json:"catalogSummaryBytes"`
	ContentDigestBytes  int64        `json:"contentDigestBytes"`
	MaxBodyBytes        int64        `json:"maxBodyBytes"`
	MaxSchemaBytes      int64        `json:"maxSchemaBytes"`
	LazyLoad            bool         `json:"lazyLoad"`
	RangedResources     bool         `json:"rangedResources"`
	MaxResumeSteps      int          `json:"maxResumeSteps"`
	MaxRepeatedChecks   int          `json:"maxRepeatedChecks"`
	Provenance          []Reference  `json:"provenance"`
	EvidenceTier        EvidenceTier `json:"evidenceTier"`
	Digest              string       `json:"digest,omitempty"`
}
