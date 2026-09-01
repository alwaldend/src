package v1alpha1

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	MaxCriteria                   = 256
	MaxGoalRelationshipReferences = 256
	MaxTitleBytes                 = 200
	MaxStatementBytes             = 4096
	MaxPublicationFiles           = 128
)

var (
	recordIDPattern = regexp.MustCompile(
		`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?` +
			`(?:\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)*$`,
	)
	dnsLabelPattern = regexp.MustCompile(
		`^[a-z0-9](?:[-a-z0-9]*[a-z0-9])?$`,
	)
	qualifiedNamePattern = regexp.MustCompile(
		`^[A-Za-z0-9](?:[-_.A-Za-z0-9]*[A-Za-z0-9])?$`,
	)
)

// Validate checks the complete backend-neutral Goal resource. Use
// ValidateSpec when validating desired state before a backend has populated
// status. ResourceVersion is opaque, and API-server or filesystem ownership
// rules are intentionally not applied.
func (g Goal) Validate() error {
	if err := g.ValidateSpec(); err != nil {
		return err
	}
	return g.ValidateStatus()
}

// ValidateSpec checks the Goal envelope and desired state without requiring
// backend-populated status.
func (g Goal) ValidateSpec() error {
	if err := validateEnvelope(g.APIVersion, g.Kind, KindGoal, g.Metadata); err != nil {
		return err
	}
	if strings.TrimSpace(g.Spec.Title) != g.Spec.Title ||
		g.Spec.Title == "" || len(g.Spec.Title) > MaxTitleBytes {
		return fmt.Errorf(
			"spec.title must be trimmed and contain 1..%d bytes",
			MaxTitleBytes,
		)
	}
	if !oneOf(g.Spec.Scope, "workspace", "project") {
		return fmt.Errorf("spec.scope must be workspace or project")
	}
	if !oneOf(g.Spec.Retention.Policy, "ephemeral", "durable") {
		return fmt.Errorf("spec.retention.policy must be ephemeral or durable")
	}
	if g.Spec.Scope == "project" && g.Spec.Retention.Policy != "durable" {
		return fmt.Errorf("project goals require durable retention")
	}
	if err := g.Spec.Relationships.validate(g.Metadata.Name); err != nil {
		return err
	}
	return nil
}

// ValidateStatus checks observed state. It assumes ValidateSpec has already
// established the envelope and spec invariants when validating a complete
// resource.
func (g Goal) ValidateStatus() error {
	if g.Status.LifecycleGeneration == 0 || g.Status.CriteriaRevision == 0 {
		return fmt.Errorf("status generation/revision fields must be positive")
	}
	if !oneOf(g.Status.Outcome, "open", "achieved", "abandoned", "superseded") {
		return fmt.Errorf("invalid status.outcome %q", g.Status.Outcome)
	}
	if !oneOf(g.Status.Execution, "active", "paused", "waiting", "blocked") {
		return fmt.Errorf("invalid status.execution %q", g.Status.Execution)
	}
	if g.Status.ActiveAttemptID != "" {
		if err := ValidateRecordID(
			"status.activeAttemptID",
			g.Status.ActiveAttemptID,
		); err != nil {
			return err
		}
	}
	if g.Status.Outcome == "achieved" {
		if err := ValidateRecordID(
			"status.acceptedAttemptID",
			g.Status.AcceptedAttemptID,
		); err != nil {
			return err
		}
		if !ValidDigest(g.Status.AcceptedResultDigest) {
			return fmt.Errorf(
				"achieved goal requires status.acceptedResultDigest",
			)
		}
	} else if g.Status.AcceptedAttemptID != "" ||
		g.Status.AcceptedResultDigest != "" {
		return fmt.Errorf(
			"acceptance pointers are allowed only for achieved outcome",
		)
	}
	if g.Status.Outcome != "open" &&
		(g.Status.Execution != "paused" || g.Status.ActiveAttemptID != "") {
		return fmt.Errorf(
			"closed outcomes require paused execution and no active attempt",
		)
	}
	if err := g.Status.Promotion.validate(g.Spec.Scope); err != nil {
		return err
	}
	if err := g.Status.Migration.validate(); err != nil {
		return err
	}
	return validateOptionalTimestamp(
		"status.observedAt",
		g.Status.ObservedAt,
	)
}

// ValidateForGoal checks a current criteria resource against its Goal.
func (c GoalCriteria) ValidateForGoal(goal Goal) error {
	return c.ValidateSnapshot(goal.Metadata.Name, goal.Status.CriteriaRevision)
}

// ValidateSnapshot checks a retained criteria revision independently of the
// current Goal status pointer.
func (c GoalCriteria) ValidateSnapshot(goalID string, revision uint64) error {
	if err := validateEnvelope(
		c.APIVersion,
		c.Kind,
		KindCriteria,
		c.Metadata,
	); err != nil {
		return err
	}
	if c.Metadata.Name != goalID || c.Spec.GoalRef.Name != goalID ||
		c.Spec.Revision == 0 || c.Spec.Revision != revision {
		return fmt.Errorf("criteria document does not match goal metadata")
	}
	if len(c.Spec.Items) > MaxCriteria {
		return fmt.Errorf(
			"criteria cardinality exceeds %d",
			MaxCriteria,
		)
	}
	seen := map[string]bool{}
	for _, criterion := range c.Spec.Items {
		if err := ValidateRecordID(
			"criterionID",
			criterion.CriterionID,
		); err != nil {
			return err
		}
		if criterion.Revision == 0 {
			return fmt.Errorf(
				"criterion %q revision must be positive",
				criterion.CriterionID,
			)
		}
		if seen[criterion.CriterionID] {
			return fmt.Errorf(
				"duplicate criterionID %q",
				criterion.CriterionID,
			)
		}
		seen[criterion.CriterionID] = true
		if strings.TrimSpace(criterion.Statement) != criterion.Statement ||
			criterion.Statement == "" ||
			len(criterion.Statement) > MaxStatementBytes {
			return fmt.Errorf(
				"criterion %q has an invalid statement",
				criterion.CriterionID,
			)
		}
		if strings.TrimSpace(criterion.EvidenceMethod) !=
			criterion.EvidenceMethod || criterion.EvidenceMethod == "" ||
			len(criterion.EvidenceMethod) > MaxStatementBytes {
			return fmt.Errorf(
				"criterion %q has an invalid evidenceMethod",
				criterion.CriterionID,
			)
		}
	}
	return nil
}

// ValidateForGoal checks an attempt and its portable input bindings.
func (a GoalAttempt) ValidateForGoal(goal Goal) error {
	if err := validateEnvelope(
		a.APIVersion,
		a.Kind,
		KindAttempt,
		a.Metadata,
	); err != nil {
		return err
	}
	if a.Spec.GoalRef.Name != goal.Metadata.Name {
		return fmt.Errorf("attempt spec.goalRef does not match goal")
	}
	if a.Spec.GoalGeneration == 0 ||
		a.Spec.LifecycleGeneration == 0 ||
		a.Spec.CriteriaRevision == 0 ||
		!ValidDigest(a.Spec.CriteriaDigest) ||
		!ValidDigest(a.Spec.GoalStateDigest) ||
		a.Spec.CriteriaRevision > goal.Status.CriteriaRevision ||
		a.Spec.LifecycleGeneration > goal.Status.LifecycleGeneration ||
		a.Spec.GoalGeneration > goal.Metadata.Generation {
		return fmt.Errorf(
			"attempt input generation/revision fields are invalid",
		)
	}
	if !oneOf(
		a.Spec.WorkType,
		"investigation",
		"candidate",
		"change",
		"integration",
		"validation",
		"decision",
	) {
		return fmt.Errorf(
			"invalid attempt spec.workType %q",
			a.Spec.WorkType,
		)
	}
	if !oneOf(a.Status.State, "open", "closed") {
		return fmt.Errorf("attempt status.state must be open or closed")
	}
	if a.Status.State == "closed" {
		if err := validateRequiredTimestamp(
			"attempt.status.closedAt",
			a.Status.ClosedAt,
		); err != nil {
			return err
		}
		if err := a.Status.Review.Validate(a.Status.Artifacts); err != nil {
			return err
		}
	}
	if a.Status.State == "open" &&
		(a.Status.ClosedAt != "" || a.Status.Review.Decision != "" ||
			len(a.Status.Review.Criteria) != 0) {
		return fmt.Errorf("open attempt cannot have close status")
	}
	if err := a.Status.Artifacts.Validate(); err != nil {
		return err
	}
	return validateOptionalTimestamp(
		"attempt.status.observedAt",
		a.Status.ObservedAt,
	)
}

// Validate checks the portable session binding. Filesystem annotations are
// optional here and are required only by the local adapter.
func (s GoalSessionBinding) Validate() error {
	if err := validateEnvelope(
		s.APIVersion,
		s.Kind,
		KindSessionBinding,
		s.Metadata,
	); err != nil {
		return err
	}
	if err := ValidateRecordID(
		"spec.goalRef.name",
		s.Spec.GoalRef.Name,
	); err != nil {
		return err
	}
	if err := ValidateRecordID("status.goalID", s.Status.GoalID); err != nil {
		return err
	}
	if s.Status.AttachedGeneration == 0 ||
		s.Status.AttachedLifecycleGeneration == 0 ||
		s.Status.AttachedCriteriaRevision == 0 ||
		!ValidDigest(s.Status.AttachedCriteriaDigest) ||
		!ValidDigest(s.Status.AttachedGoalStateDigest) {
		return fmt.Errorf(
			"session portable attachment bindings are invalid",
		)
	}
	if s.Spec.GoalRef.Name != s.Status.GoalID {
		return fmt.Errorf("session spec.goalRef and status.goalID differ")
	}
	return validateOptionalTimestamp(
		"status.observedAt",
		s.Status.ObservedAt,
	)
}

// Validate checks a close review against the frozen artifact inventory.
func (r CloseReview) Validate(artifacts ArtifactManifest) error {
	if !oneOf(r.Decision, "accept", "refine", "reset") {
		return fmt.Errorf(
			"closed attempt review decision must be accept, refine, or reset",
		)
	}
	if len(r.Criteria) > MaxCriteria {
		return fmt.Errorf(
			"attempt review criterion cardinality exceeds %d",
			MaxCriteria,
		)
	}
	allowedRefs := map[string]bool{"plan.md": true, "result.md": true}
	for _, artifact := range artifacts.Evidence {
		allowedRefs[artifact.Path] = true
	}
	seen := map[string]bool{}
	previousCriterion := ""
	for _, review := range r.Criteria {
		if err := ValidateRecordID(
			"attempt review criterionID",
			review.CriterionID,
		); err != nil {
			return err
		}
		if review.CriterionRevision == 0 ||
			!oneOf(review.Verdict, "pass", "fail", "unverified") {
			return fmt.Errorf(
				"attempt review for %q has an invalid revision or verdict",
				review.CriterionID,
			)
		}
		if seen[review.CriterionID] {
			return fmt.Errorf(
				"attempt review contains duplicate criterionID %q",
				review.CriterionID,
			)
		}
		if review.CriterionID <= previousCriterion {
			return fmt.Errorf(
				"attempt review criteria must be unique and sorted",
			)
		}
		seen[review.CriterionID] = true
		previousCriterion = review.CriterionID
		if review.Verdict != "unverified" && len(review.EvidenceRefs) == 0 {
			return fmt.Errorf(
				"attempt review %s for %q requires frozen evidenceRefs",
				review.Verdict,
				review.CriterionID,
			)
		}
		previousReference := ""
		for _, reference := range review.EvidenceRefs {
			if !allowedRefs[reference] {
				return fmt.Errorf(
					"attempt review evidence reference %q is not a frozen artifact",
					reference,
				)
			}
			if reference <= previousReference {
				return fmt.Errorf(
					"attempt review evidence references must be unique and sorted",
				)
			}
			previousReference = reference
		}
	}
	return nil
}

// Validate checks a frozen artifact digest manifest.
func (a ArtifactManifest) Validate() error {
	if !ValidDigest(a.PlanDigest) || !ValidDigest(a.ResultDigest) {
		return fmt.Errorf("attempt artifact plan/result digests are invalid")
	}
	previous := ""
	for _, item := range a.Evidence {
		if !validEvidenceArtifactPath(item.Path) ||
			!ValidDigest(item.Digest) {
			return fmt.Errorf("attempt evidence artifact entry is invalid")
		}
		if item.Path <= previous {
			return fmt.Errorf(
				"attempt evidence artifact entries must be unique and sorted",
			)
		}
		previous = item.Path
	}
	return nil
}

// ValidateRecordID validates the stable, portable identifier subset used by
// goal resources. The subset is stricter than a Kubernetes DNS subdomain.
func ValidateRecordID(field string, value string) error {
	if len(value) > 64 || !recordIDPattern.MatchString(value) ||
		value == "." || value == ".." {
		return fmt.Errorf(
			"%s %q is not a portable record ID",
			field,
			value,
		)
	}
	for _, label := range strings.Split(value, ".") {
		if len(label) > 63 {
			return fmt.Errorf(
				"%s %q is not a portable record ID",
				field,
				value,
			)
		}
	}
	return nil
}

// ValidDigest reports whether value is a canonical SHA-256 digest string.
func ValidDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") ||
		len(value) != len("sha256:")+64 {
		return false
	}
	for _, character := range value[len("sha256:"):] {
		if !(character >= '0' && character <= '9') &&
			!(character >= 'a' && character <= 'f') {
			return false
		}
	}
	return true
}

func validateEnvelope(
	apiVersion string,
	kind string,
	wantKind string,
	metadata ObjectMeta,
) error {
	if apiVersion != APIVersion {
		return fmt.Errorf("unsupported apiVersion %q", apiVersion)
	}
	if kind != wantKind {
		return fmt.Errorf("kind must be %s", wantKind)
	}
	if err := ValidateRecordID("metadata.name", metadata.Name); err != nil {
		return err
	}
	if metadata.Namespace != "" &&
		(len(metadata.Namespace) > 63 ||
			!dnsLabelPattern.MatchString(metadata.Namespace)) {
		return fmt.Errorf(
			"metadata.namespace is not a Kubernetes DNS label",
		)
	}
	if err := validateOptionalTimestamp(
		"metadata.creationTimestamp",
		metadata.CreationTimestamp,
	); err != nil {
		return err
	}
	for key, value := range metadata.Labels {
		if err := validateQualifiedName("metadata.labels", key); err != nil {
			return err
		}
		if err := validateLabelValue(key, value); err != nil {
			return err
		}
	}
	for key := range metadata.Annotations {
		if err := validateQualifiedName(
			"metadata.annotations",
			key,
		); err != nil {
			return err
		}
	}
	return nil
}

func validateQualifiedName(field string, value string) error {
	if value == "" {
		return fmt.Errorf("%s contains an empty key", field)
	}
	parts := strings.Split(value, "/")
	if len(parts) > 2 {
		return fmt.Errorf("%s key %q is not a qualified name", field, value)
	}
	name := parts[len(parts)-1]
	if len(name) == 0 || len(name) > 63 ||
		!qualifiedNamePattern.MatchString(name) {
		return fmt.Errorf("%s key %q is not a qualified name", field, value)
	}
	if len(parts) == 2 && !validDNSSubdomain(parts[0]) {
		return fmt.Errorf("%s key %q is not a qualified name", field, value)
	}
	return nil
}

func validateLabelValue(key string, value string) error {
	if len(value) > 63 ||
		(value != "" && !qualifiedNamePattern.MatchString(value)) {
		return fmt.Errorf(
			"metadata.labels value for %q is not Kubernetes-compatible",
			key,
		)
	}
	return nil
}

func validDNSSubdomain(value string) bool {
	if value == "" || len(value) > 253 {
		return false
	}
	for _, label := range strings.Split(value, ".") {
		if len(label) == 0 || len(label) > 63 ||
			!dnsLabelPattern.MatchString(label) {
			return false
		}
	}
	return true
}

func (r Relationships) validate(self string) error {
	normalized := r.Normalized()
	referenceCount := len(normalized.DependsOnGoalRefs) +
		len(normalized.SupersedesGoalRefs)
	if normalized.ParentGoalRef != nil {
		referenceCount++
	}
	if referenceCount > MaxGoalRelationshipReferences {
		return fmt.Errorf(
			"spec.relationships cardinality exceeds %d",
			MaxGoalRelationshipReferences,
		)
	}
	if normalized.ParentGoalRef != nil {
		if err := ValidateRecordID(
			"spec.relationships.parentGoalRef.name",
			normalized.ParentGoalRef.Name,
		); err != nil {
			return err
		}
		if normalized.ParentGoalRef.Name == self {
			return fmt.Errorf("goal cannot be its own parent")
		}
	}
	if err := validateGoalReferenceList(
		"spec.relationships.dependsOnGoalRefs",
		self,
		normalized.DependsOnGoalRefs,
	); err != nil {
		return err
	}
	return validateGoalReferenceList(
		"spec.relationships.supersedesGoalRefs",
		self,
		normalized.SupersedesGoalRefs,
	)
}

func validateGoalReferenceList(
	field string,
	self string,
	values []GoalReference,
) error {
	previous := ""
	for _, value := range values {
		if err := ValidateRecordID(field+".name", value.Name); err != nil {
			return err
		}
		if value.Name == self {
			return fmt.Errorf("%s contains a self-reference", field)
		}
		if value.Name == previous {
			return fmt.Errorf(
				"%s contains duplicate goal reference %q",
				field,
				value.Name,
			)
		}
		previous = value.Name
	}
	return nil
}

func (p PromotionStatus) validate(scope string) error {
	empty := p.SourceScope == "" && p.SourceGeneration == 0 &&
		p.SourceLifecycleGeneration == 0 &&
		p.SourceCriteriaRevision == 0 && p.SourceCriteriaDigest == "" &&
		p.SourceStateDigest == "" && p.SourceDigest == "" &&
		p.PromotedAt == ""
	if empty {
		return nil
	}
	if scope != "project" || p.SourceScope != "workspace" ||
		p.SourceGeneration == 0 || p.SourceLifecycleGeneration == 0 ||
		p.SourceCriteriaRevision == 0 ||
		!ValidDigest(p.SourceCriteriaDigest) ||
		!ValidDigest(p.SourceStateDigest) || !ValidDigest(p.SourceDigest) {
		return fmt.Errorf(
			"status.promotion must describe a valid workspace source",
		)
	}
	return validateRequiredTimestamp(
		"status.promotion.promotedAt",
		p.PromotedAt,
	)
}

func (m MigrationStatus) validate() error {
	empty := m.SourceFormat == "" && m.SourceDigest == "" &&
		m.MigratedAt == ""
	if empty {
		return nil
	}
	if m.SourceFormat != "unversioned" || !ValidDigest(m.SourceDigest) {
		return fmt.Errorf(
			"status.migration must describe an unversioned source",
		)
	}
	return validateRequiredTimestamp(
		"status.migration.migratedAt",
		m.MigratedAt,
	)
}

func validateOptionalTimestamp(field string, value string) error {
	if value == "" {
		return nil
	}
	return validateRequiredTimestamp(field, value)
}

func validateRequiredTimestamp(field string, value string) error {
	if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
		return fmt.Errorf("%s is not RFC3339: %w", field, err)
	}
	return nil
}

func validEvidenceArtifactPath(value string) bool {
	if !strings.HasPrefix(value, "evidence/") {
		return false
	}
	name := strings.TrimPrefix(value, "evidence/")
	if name == "" || len(name) > 128 || name == "." || name == ".." {
		return false
	}
	for _, character := range name {
		if !(character >= 'a' && character <= 'z') &&
			!(character >= 'A' && character <= 'Z') &&
			!(character >= '0' && character <= '9') &&
			character != '.' && character != '_' && character != '-' {
			return false
		}
	}
	return true
}

func oneOf(value string, values ...string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}
