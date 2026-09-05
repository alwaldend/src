package fsstore

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	v1alpha1 "git.alwaldend.com/alwaldend/src/projects/goal/api/v1alpha1"
)

const (
	// APIVersion is the wire version emitted by this backend.
	APIVersion                   = v1alpha1.APIVersion
	goalAPIVersion               = APIVersion
	localGoalReferenceAnnotation = v1alpha1.LocalGoalReferenceAnnotation
	localOwnerRootAnnotation     = v1alpha1.LocalOwnerRootAnnotation
	maxCriteria                  = v1alpha1.MaxCriteria
	maxCriteriaRevisions         = 4096
	maxGoals                     = 4096
	maxAttempts                  = 1024
	maxPlans                     = 64
	maxEvidenceFiles             = 256
	maxManifestBytes             = 256 * 1024
	maxPlanResultBytes           = 2 * 1024 * 1024
	maxEvidenceFileBytes         = 16 * 1024 * 1024
	maxCheckpointSummaryBytes    = 8192
	defaultOutputLimit           = 20
	maximumOutputLimit           = 100
)

type (
	ObjectMeta         = v1alpha1.ObjectMeta
	LocalGoalReference = v1alpha1.GoalReference
	Retention          = v1alpha1.Retention
	Relationships      = v1alpha1.Relationships
	GoalSpec           = v1alpha1.GoalSpec
	PromotionStatus    = v1alpha1.PromotionStatus
	MigrationStatus    = v1alpha1.MigrationStatus
	GoalStatus         = v1alpha1.GoalStatus
	Criterion          = v1alpha1.Criterion
	PlanSummary        = v1alpha1.PlanSummary
	CriteriaSpec       = v1alpha1.CriteriaSpec
	AttemptSpec        = v1alpha1.AttemptSpec
	ArtifactDigest     = v1alpha1.ArtifactDigest
	ArtifactManifest   = v1alpha1.ArtifactManifest
	CriterionReview    = v1alpha1.CriterionReview
	CloseReview        = v1alpha1.CloseReview
	AttemptStatus      = v1alpha1.AttemptStatus
	SessionSpec        = v1alpha1.SessionSpec
	SessionStatus      = v1alpha1.SessionStatus
)

// These local names preserve the store's concise implementation while every
// wire field and backend-neutral rule is owned by api/v1alpha1.
type (
	GoalManifest     v1alpha1.Goal
	CriteriaManifest v1alpha1.GoalCriteria
	AttemptManifest  v1alpha1.GoalAttempt
	SessionBinding   v1alpha1.GoalSessionBinding
)

func (g GoalManifest) validate() error {
	if err := v1alpha1.Goal(g).Validate(); err != nil {
		return err
	}
	if err := validatePersistedMetadata(g.Metadata); err != nil {
		return err
	}
	ownerRoot, ok := g.Metadata.Annotations[v1alpha1.LocalOwnerRootAnnotation]
	if !ok {
		return fmt.Errorf(
			"metadata.annotations must contain %s",
			v1alpha1.LocalOwnerRootAnnotation,
		)
	}
	if err := validatePortablePath(
		"local owner root annotation",
		ownerRoot,
		false,
		true,
	); err != nil {
		return err
	}
	return validateTimestamp(
		"status.observedAt",
		g.Status.ObservedAt,
		true,
	)
}

func (c CriteriaManifest) validate(goal GoalManifest) error {
	if err := v1alpha1.GoalCriteria(c).ValidateForGoal(
		v1alpha1.Goal(goal),
	); err != nil {
		return err
	}
	return validatePersistedMetadata(c.Metadata)
}

func (c CriteriaManifest) validateSnapshot(
	goalID string,
	revision uint64,
) error {
	if err := v1alpha1.GoalCriteria(c).ValidateSnapshot(
		goalID,
		revision,
	); err != nil {
		return err
	}
	return validatePersistedMetadata(c.Metadata)
}

func (a AttemptManifest) validate(goal GoalManifest) error {
	if err := v1alpha1.GoalAttempt(a).ValidateForGoal(v1alpha1.Goal(goal)); err != nil {
		return fmt.Errorf("planID %q: %w", a.Spec.PlanID, err)
	}
	if err := validatePersistedMetadata(a.Metadata); err != nil {
		return err
	}
	return validateTimestamp(
		"attempt.status.observedAt",
		a.Status.ObservedAt,
		true,
	)
}

func (s SessionBinding) validate() error {
	if err := v1alpha1.GoalSessionBinding(s).Validate(); err != nil {
		return err
	}
	if err := validatePersistedMetadata(s.Metadata); err != nil {
		return err
	}
	goalReference, ok := s.Metadata.Annotations[v1alpha1.LocalGoalReferenceAnnotation]
	if !ok {
		return fmt.Errorf(
			"metadata.annotations must contain %s",
			v1alpha1.LocalGoalReferenceAnnotation,
		)
	}
	if err := validatePortablePath(
		"local goal reference annotation",
		goalReference,
		false,
		false,
	); err != nil {
		return err
	}
	return validateTimestamp(
		"status.observedAt",
		s.Status.ObservedAt,
		true,
	)
}

func validatePersistedMetadata(metadata ObjectMeta) error {
	if _, err := parseResourceVersion(metadata.ResourceVersion); err != nil {
		return fmt.Errorf("metadata.resourceVersion: %w", err)
	}
	if metadata.Generation == 0 {
		return fmt.Errorf("metadata.generation must be positive")
	}
	return validateTimestamp(
		"metadata.creationTimestamp",
		metadata.CreationTimestamp,
		true,
	)
}

func parseResourceVersion(value string) (uint64, error) {
	if value == "" || (len(value) > 1 && value[0] == '0') {
		return 0, fmt.Errorf("must be a canonical positive decimal string")
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return 0, fmt.Errorf("must be a canonical positive decimal string")
	}
	return parsed, nil
}

func incrementResourceVersion(value string) (string, error) {
	parsed, err := parseResourceVersion(value)
	if err != nil || parsed == ^uint64(0) {
		return "", fmt.Errorf("cannot increment resourceVersion")
	}
	return strconv.FormatUint(parsed+1, 10), nil
}

func validateRecordID(field string, value string) error {
	return v1alpha1.ValidateRecordID(field, value)
}

func validatePortablePath(
	field string,
	value string,
	allowParent bool,
	allowDot bool,
) error {
	if value == "" || filepath.IsAbs(value) || strings.Contains(value, "\\") {
		return fmt.Errorf(
			"%s must be a non-empty relative portable path",
			field,
		)
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if clean != value || clean == "." && !allowDot {
		return fmt.Errorf("%s must be normalized", field)
	}
	if !allowParent &&
		(clean == ".." || strings.HasPrefix(clean, "../")) {
		return fmt.Errorf("%s must stay under the workspace root", field)
	}
	return nil
}

func validateTimestamp(field string, value string, required bool) error {
	if value == "" && !required {
		return nil
	}
	if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
		return fmt.Errorf("%s is not RFC3339: %w", field, err)
	}
	return nil
}

func validDigest(value string) bool {
	return v1alpha1.ValidDigest(value)
}

func oneOf(value string, values ...string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}
