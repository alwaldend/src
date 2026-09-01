// Phase 3A: advisory impact planning.
package v1alpha1

// ImpactProfile names a conservative validation profile bound to a plan.
type ImpactProfile string

const (
	ImpactProfileChangedFast   ImpactProfile = "changed/fast"
	ImpactProfileWorkspace     ImpactProfile = "workspace"
	ImpactProfileFreshEvidence ImpactProfile = "fresh/evidence"
	ImpactProfileFullAudit     ImpactProfile = "full/audit"
	ImpactProfileDiagnose      ImpactProfile = "diagnose"
)

// CostClass bounds a plan's expected resource envelope.
type CostClass string

const (
	CostClassLow    CostClass = "low"
	CostClassMedium CostClass = "medium"
	CostClassHigh   CostClass = "high"
)

// PlanTarget is one changed or reverse-affected target at narrowest scope.
type PlanTarget struct {
	Label           string `json:"label"`
	Path            string `json:"path,omitempty"`
	AffectedBy      string `json:"affectedBy,omitempty"`
	ReverseAffected bool   `json:"reverseAffected,omitempty"`
}

// PlanCheck is one conservative minimum check plus its coverage note.
type PlanCheck struct {
	Identifier string `json:"identifier"`
	Scope      string `json:"scope,omitempty"`
	Covers     string `json:"covers,omitempty"`
}

// ImpactPlan is the deterministic smallest-sufficient plan for one intent.
// It is advisory: goal criteria or agent judgment own semantic sufficiency;
// selection never owns authority.
type ImpactPlan struct {
	APIVersion        string           `json:"apiVersion"`
	Kind              string           `json:"kind"`
	ID                string           `json:"id"`
	IntentRef         Reference        `json:"intentRef"`
	Profile           ImpactProfile    `json:"profile"`
	Capabilities      []Reference      `json:"capabilities"`
	CapabilityReasons []string         `json:"capabilityReasons,omitempty"`
	Effects           []Effect         `json:"effects"`
	ForbiddenEffects  []Effect         `json:"forbiddenEffects,omitempty"`
	Targets           []PlanTarget     `json:"targets"`
	Checks            []PlanCheck      `json:"checks"`
	CoverageGaps      []string         `json:"coverageGaps,omitempty"`
	CostClass         CostClass        `json:"costClass"`
	ExpectedMaxCost   Budget           `json:"expectedMaxCost"`
	Cacheability      Cacheability     `json:"cacheability"`
	Ordering          []string         `json:"ordering,omitempty"`
	Concurrency       int              `json:"concurrency,omitempty"`
	Escalation        string           `json:"escalation,omitempty"`
	EvidenceRefs      []InputReference `json:"evidenceRefs,omitempty"`
	Digest            string           `json:"digest,omitempty"`
}
