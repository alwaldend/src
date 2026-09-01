// Package v1alpha1 defines the goal publication intent resource.
//
// The intent is a bounded, versioned replay contract written before a
// multi-file goal publication begins. It is backend persistence metadata, not
// an independent authority: it records the exact staged after-images and
// before/after digests so that an interrupted publication can be diagnosed
// and recovered by the owning store.
package v1alpha1

import "fmt"

const (
	KindGoalPublication = "GoalPublication"

	// Publication states reported by doctor.
	PublicationStable      = "stable"
	PublicationIncomplete  = "incomplete-recoverable"
	PublicationProjection  = "committed-projection-stale"
	PublicationConflict    = "conflict"
	PublicationDiscardable = "discardable-intent"
	PublicationStaged      = "staged-intent"
	PublicationPartial     = "partial-intent"
)

type GoalPublication struct {
	APIVersion string                `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta            `json:"metadata" yaml:"metadata"`
	Spec       GoalPublicationSpec   `json:"spec" yaml:"spec"`
	Status     GoalPublicationStatus `json:"status" yaml:"status"`
}

type GoalPublicationSpec struct {
	// GoalRef binds the intent to exactly one local Goal record.
	GoalRef GoalReference `json:"goalRef" yaml:"goalRef"`
	// OperationID is a stable, process-unique identifier.
	OperationID string `json:"operationID" yaml:"operationID"`
	// PriorResourceVersion and IntendedResourceVersion bind the intent to the
	// optimistic-concurrency checkpoint it publishes.
	PriorResourceVersion    string `json:"priorResourceVersion" yaml:"priorResourceVersion"`
	IntendedResourceVersion string `json:"intendedResourceVersion" yaml:"intendedResourceVersion"`
	// Files are the canonical files to publish, in deterministic order.
	Files []GoalPublicationFile `json:"files" yaml:"files"`
	// SnapshotDigests retains immutable criteria-snapshot content so an
	// interrupted intent can be validated without regenerating state.
	SnapshotDigests map[uint64]string `json:"snapshotDigests" yaml:"snapshotDigests"`
}

type GoalPublicationFile struct {
	// Path is the workspace-relative canonical path within the goal record.
	Path string `json:"path" yaml:"path"`
	// BeforeDigest is empty when the file was absent before publication.
	BeforeDigest string `json:"beforeDigest,omitempty" yaml:"beforeDigest,omitempty"`
	// AfterDigest is the exact intended content digest.
	AfterDigest string `json:"afterDigest" yaml:"afterDigest"`
	// StagedRelative is the staged after-image path under the intent staging
	// directory.
	StagedRelative string `json:"stagedRelative" yaml:"stagedRelative"`
}

type GoalPublicationStatus struct {
	// State reports doctor classification.
	State      string `json:"state" yaml:"state"`
	Message    string `json:"message,omitempty" yaml:"message,omitempty"`
	ObservedAt string `json:"observedAt" yaml:"observedAt"`
}

// Validate checks the portable publication intent envelope and bounds. Paths
// are backend-owned and validated by the filesystem adapter; the API validates
// the stable identity and resource-version binding only.
func (p GoalPublication) Validate() error {
	if err := validateEnvelope(
		p.APIVersion,
		p.Kind,
		KindGoalPublication,
		p.Metadata,
	); err != nil {
		return err
	}
	if err := ValidateRecordID(
		"spec.goalRef.name",
		p.Spec.GoalRef.Name,
	); err != nil {
		return err
	}
	if p.Spec.OperationID == "" || len(p.Spec.OperationID) > 64 {
		return fmt.Errorf("spec.operationID is invalid")
	}
	if p.Spec.PriorResourceVersion == "" ||
		p.Spec.IntendedResourceVersion == "" {
		return fmt.Errorf("spec resource versions are required")
	}
	if len(p.Spec.Files) == 0 || len(p.Spec.Files) > MaxPublicationFiles {
		return fmt.Errorf(
			"spec.files cardinality must be between 1 and %d",
			MaxPublicationFiles,
		)
	}
	seen := map[string]bool{}
	for _, file := range p.Spec.Files {
		if file.Path == "" || file.AfterDigest == "" ||
			!ValidDigest(file.AfterDigest) ||
			(file.BeforeDigest != "" && !ValidDigest(file.BeforeDigest)) {
			return fmt.Errorf("spec.files entry is invalid")
		}
		if file.StagedRelative == "" || seen[file.StagedRelative] {
			return fmt.Errorf("spec.files staged path must be unique and present")
		}
		seen[file.StagedRelative] = true
	}
	for revision, digest := range p.Spec.SnapshotDigests {
		if revision == 0 || !ValidDigest(digest) {
			return fmt.Errorf("spec.snapshotDigests entry is invalid")
		}
	}
	if !oneOf(
		p.Status.State,
		"",
		PublicationStable,
		PublicationIncomplete,
		PublicationProjection,
		PublicationConflict,
		PublicationDiscardable,
		PublicationStaged,
		PublicationPartial,
	) {
		return fmt.Errorf("status.state is invalid")
	}
	return validateOptionalTimestamp(
		"status.observedAt",
		p.Status.ObservedAt,
	)
}
