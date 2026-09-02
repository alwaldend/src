// Package v1alpha1 defines the portable goal resource contract.
//
// The resources deliberately resemble Kubernetes API objects, but they do not
// implement Kubernetes runtime interfaces and are not a substitute for a CRD.
// Filesystem persistence policy belongs to internal/fsstore.
package v1alpha1

const (
	APIVersion = "goals.alwaldend.com/v1alpha1"

	KindGoal           = "Goal"
	KindCriteria       = "GoalCriteria"
	KindAttempt        = "GoalAttempt"
	KindSessionBinding = "GoalSessionBinding"

	// LocalGoalReferenceAnnotation stores a workspace-relative filesystem
	// reference for the local session adapter.
	LocalGoalReferenceAnnotation = "goals.alwaldend.com/local-goal-ref"

	// LocalOwnerRootAnnotation stores the workspace-relative owner root used by
	// the local filesystem adapter. It is not portable desired state.
	LocalOwnerRootAnnotation = "goals.alwaldend.com/local-owner-root"
)

type ObjectMeta struct {
	Name              string            `json:"name" yaml:"name"`
	Namespace         string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
	Generation        uint64            `json:"generation,omitempty" yaml:"generation,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// GoalReference identifies another Goal in the same catalog. References are
// deliberately name-only: filesystem paths, URLs, namespaces, and backend
// identifiers are not part of the portable relationship contract.
type GoalReference struct {
	Name string `json:"name" yaml:"name"`
}

type Retention struct {
	Policy string `json:"policy" yaml:"policy"`
}

type Relationships struct {
	ParentGoalRef      *GoalReference  `json:"parentGoalRef,omitempty" yaml:"parentGoalRef,omitempty"`
	DependsOnGoalRefs  []GoalReference `json:"dependsOnGoalRefs" yaml:"dependsOnGoalRefs"`
	SupersedesGoalRefs []GoalReference `json:"supersedesGoalRefs" yaml:"supersedesGoalRefs"`
}

type GoalSpec struct {
	Title         string        `json:"title" yaml:"title"`
	Scope         string        `json:"scope" yaml:"scope"`
	Retention     Retention     `json:"retention" yaml:"retention"`
	Relationships Relationships `json:"relationships" yaml:"relationships"`
}

type PromotionStatus struct {
	SourceScope               string `json:"sourceScope" yaml:"sourceScope"`
	SourceGeneration          uint64 `json:"sourceGeneration" yaml:"sourceGeneration"`
	SourceLifecycleGeneration uint64 `json:"sourceLifecycleGeneration" yaml:"sourceLifecycleGeneration"`
	SourceCriteriaRevision    uint64 `json:"sourceCriteriaRevision" yaml:"sourceCriteriaRevision"`
	SourceCriteriaDigest      string `json:"sourceCriteriaDigest" yaml:"sourceCriteriaDigest"`
	SourceStateDigest         string `json:"sourceStateDigest" yaml:"sourceStateDigest"`
	SourceDigest              string `json:"sourceDigest" yaml:"sourceDigest"`
	PromotedAt                string `json:"promotedAt" yaml:"promotedAt"`
}

type MigrationStatus struct {
	SourceFormat   string `json:"sourceFormat" yaml:"sourceFormat"`
	SourcePath     string `json:"sourcePath" yaml:"sourcePath"`
	SourceDigest   string `json:"sourceDigest" yaml:"sourceDigest"`
	MappingVersion string `json:"mappingVersion" yaml:"mappingVersion"`
	ExtractionMode string `json:"extractionMode" yaml:"extractionMode"`
	MigratedAt     string `json:"migratedAt" yaml:"migratedAt"`
}

type GoalStatus struct {
	LifecycleGeneration  uint64          `json:"lifecycleGeneration" yaml:"lifecycleGeneration"`
	Outcome              string          `json:"outcome" yaml:"outcome"`
	Execution            string          `json:"execution" yaml:"execution"`
	ActiveAttemptID      string          `json:"activeAttemptID" yaml:"activeAttemptID"`
	AcceptedAttemptID    string          `json:"acceptedAttemptID" yaml:"acceptedAttemptID"`
	AcceptedResultDigest string          `json:"acceptedResultDigest" yaml:"acceptedResultDigest"`
	CriteriaRevision     uint64          `json:"criteriaRevision" yaml:"criteriaRevision"`
	Promotion            PromotionStatus `json:"promotion" yaml:"promotion"`
	Migration            MigrationStatus `json:"migration" yaml:"migration"`
	ObservedAt           string          `json:"observedAt" yaml:"observedAt"`
}

type Goal struct {
	APIVersion string     `json:"apiVersion" yaml:"apiVersion"`
	Kind       string     `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       GoalSpec   `json:"spec" yaml:"spec"`
	Status     GoalStatus `json:"status" yaml:"status"`
}

type Criterion struct {
	CriterionID    string `json:"criterionID" yaml:"criterionID"`
	Revision       uint64 `json:"revision" yaml:"revision"`
	Required       bool   `json:"required" yaml:"required"`
	Statement      string `json:"statement" yaml:"statement"`
	EvidenceMethod string `json:"evidenceMethod" yaml:"evidenceMethod"`
}

type CriteriaSpec struct {
	GoalRef  GoalReference `json:"goalRef" yaml:"goalRef"`
	Revision uint64        `json:"revision" yaml:"revision"`
	Items    []Criterion   `json:"items" yaml:"items"`
}

type GoalCriteria struct {
	APIVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind       string       `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta   `json:"metadata" yaml:"metadata"`
	Spec       CriteriaSpec `json:"spec" yaml:"spec"`
}

type AttemptSpec struct {
	GoalRef             GoalReference `json:"goalRef" yaml:"goalRef"`
	GoalGeneration      uint64        `json:"goalGeneration" yaml:"goalGeneration"`
	LifecycleGeneration uint64        `json:"lifecycleGeneration" yaml:"lifecycleGeneration"`
	CriteriaRevision    uint64        `json:"criteriaRevision" yaml:"criteriaRevision"`
	CriteriaDigest      string        `json:"criteriaDigest" yaml:"criteriaDigest"`
	GoalStateDigest     string        `json:"goalStateDigest" yaml:"goalStateDigest"`
	WorkType            string        `json:"workType" yaml:"workType"`
	// StableDefect is the reproducible problem under investigation.
	StableDefect string `json:"stableDefect,omitempty" yaml:"stableDefect,omitempty"`
	// Hypothesis is the candidate explanation being tested.
	Hypothesis string `json:"hypothesis,omitempty" yaml:"hypothesis,omitempty"`
	// Subject names the exact system, artifact, or reference under test.
	Subject string `json:"subject,omitempty" yaml:"subject,omitempty"`
	// AffectedCriteria lists the criterion IDs this attempt exercises.
	AffectedCriteria []string `json:"affectedCriteria,omitempty" yaml:"affectedCriteria,omitempty"`
	// RegressionRefs records the reviewed regression set or fixtures.
	RegressionRefs []string `json:"regressionRefs,omitempty" yaml:"regressionRefs,omitempty"`
	// PriorAttemptID names an earlier attempt this one resumes or corrects.
	PriorAttemptID string `json:"priorAttemptID,omitempty" yaml:"priorAttemptID,omitempty"`
	// DominantFailure is the single most useful failure signal observed.
	DominantFailure string `json:"dominantFailure,omitempty" yaml:"dominantFailure,omitempty"`
	// MeasurableDelta is the measured difference this attempt produces.
	MeasurableDelta string `json:"measurableDelta,omitempty" yaml:"measurableDelta,omitempty"`
	// NextAction is the deterministic next step for a resuming agent.
	NextAction string `json:"nextAction,omitempty" yaml:"nextAction,omitempty"`
	// Blocker names the external condition preventing progress, if any.
	Blocker string `json:"blocker,omitempty" yaml:"blocker,omitempty"`
	// ResumeCondition states what must hold before resuming this attempt.
	ResumeCondition string `json:"resumeCondition,omitempty" yaml:"resumeCondition,omitempty"`
}

type ArtifactDigest struct {
	Path   string `json:"path" yaml:"path"`
	Digest string `json:"digest" yaml:"digest"`
}

type ArtifactManifest struct {
	PlanDigest   string           `json:"planDigest" yaml:"planDigest"`
	ResultDigest string           `json:"resultDigest" yaml:"resultDigest"`
	Evidence     []ArtifactDigest `json:"evidence" yaml:"evidence"`
}

type CriterionReview struct {
	CriterionID       string   `json:"criterionID" yaml:"criterionID"`
	CriterionRevision uint64   `json:"criterionRevision" yaml:"criterionRevision"`
	Verdict           string   `json:"verdict" yaml:"verdict"`
	EvidenceRefs      []string `json:"evidenceRefs" yaml:"evidenceRefs"`
}

type CloseReview struct {
	Decision string            `json:"decision" yaml:"decision"`
	Criteria []CriterionReview `json:"criteria" yaml:"criteria"`
}

type AttemptStatus struct {
	State      string           `json:"state" yaml:"state"`
	ClosedAt   string           `json:"closedAt" yaml:"closedAt"`
	Artifacts  ArtifactManifest `json:"artifacts" yaml:"artifacts"`
	Review     CloseReview      `json:"review" yaml:"review"`
	ObservedAt string           `json:"observedAt" yaml:"observedAt"`
}

type GoalAttempt struct {
	APIVersion string        `json:"apiVersion" yaml:"apiVersion"`
	Kind       string        `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta    `json:"metadata" yaml:"metadata"`
	Spec       AttemptSpec   `json:"spec" yaml:"spec"`
	Status     AttemptStatus `json:"status" yaml:"status"`
}

type SessionSpec struct {
	GoalRef GoalReference `json:"goalRef" yaml:"goalRef"`
}

type SessionStatus struct {
	GoalID                      string `json:"goalID" yaml:"goalID"`
	AttachedGeneration          uint64 `json:"attachedGeneration" yaml:"attachedGeneration"`
	AttachedLifecycleGeneration uint64 `json:"attachedLifecycleGeneration" yaml:"attachedLifecycleGeneration"`
	AttachedCriteriaRevision    uint64 `json:"attachedCriteriaRevision" yaml:"attachedCriteriaRevision"`
	AttachedCriteriaDigest      string `json:"attachedCriteriaDigest" yaml:"attachedCriteriaDigest"`
	AttachedGoalStateDigest     string `json:"attachedGoalStateDigest" yaml:"attachedGoalStateDigest"`
	ObservedAt                  string `json:"observedAt" yaml:"observedAt"`
}

type GoalSessionBinding struct {
	APIVersion string        `json:"apiVersion" yaml:"apiVersion"`
	Kind       string        `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta    `json:"metadata" yaml:"metadata"`
	Spec       SessionSpec   `json:"spec" yaml:"spec"`
	Status     SessionStatus `json:"status" yaml:"status"`
}
