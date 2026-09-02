package v1alpha1

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func validGoal() Goal {
	return Goal{
		APIVersion: APIVersion,
		Kind:       KindGoal,
		Metadata: ObjectMeta{
			Name: "portable-goal",
		},
		Spec: GoalSpec{
			Title:     "Portable goal",
			Scope:     "workspace",
			Retention: Retention{Policy: "ephemeral"},
			Relationships: Relationships{
				DependsOnGoalRefs:  []GoalReference{},
				SupersedesGoalRefs: []GoalReference{},
			},
		},
		Status: GoalStatus{
			LifecycleGeneration: 1,
			Outcome:             "open",
			Execution:           "active",
			CriteriaRevision:    1,
		},
	}
}

func TestPortableValidationTreatsServerMetadataAsOpaque(t *testing.T) {
	for _, resourceVersion := range []string{
		"",
		"0007",
		"backend:opaque/token==",
	} {
		t.Run(resourceVersion, func(t *testing.T) {
			goal := validGoal()
			goal.Metadata.ResourceVersion = resourceVersion
			if err := goal.Validate(); err != nil {
				t.Fatalf("opaque resourceVersion %q: %v", resourceVersion, err)
			}
		})
	}
}

func TestPortableRecordIDsRespectDNSLabelBounds(t *testing.T) {
	for _, value := range []string{
		strings.Repeat("a", 63),
		strings.Repeat("a", 32) + "." + strings.Repeat("b", 31),
	} {
		if err := ValidateRecordID("test", value); err != nil {
			t.Errorf("valid record ID %q: %v", value, err)
		}
	}
	if err := ValidateRecordID("test", strings.Repeat("a", 64)); err == nil {
		t.Fatal("64-byte DNS label was accepted")
	}
}

func TestKubernetesCompatibleMetadataKeysAndLabelValues(t *testing.T) {
	prefix := strings.Repeat("a", 63) + "." +
		strings.Repeat("b", 63) + "." +
		strings.Repeat("c", 63) + "." + strings.Repeat("d", 61)
	validKeys := []string{
		"name",
		"example.com/name",
		"example.com/Name_1.x",
		strings.Repeat("n", 63),
		prefix + "/x",
	}
	validValues := []string{"", "A", "a-b_C.1", strings.Repeat("v", 63)}
	for index, key := range validKeys {
		goal := validGoal()
		goal.Metadata.Labels = map[string]string{
			key: validValues[index%len(validValues)],
		}
		goal.Metadata.Annotations = map[string]string{
			key: "arbitrary annotation\n{\"value\": true}",
		}
		if err := goal.Validate(); err != nil {
			t.Errorf("valid metadata key %q: %v", key, err)
		}
	}

	invalidKeys := []string{
		"",
		" white-space",
		"/name",
		"prefix/",
		"a/b/c",
		"Example.com/name",
		"example.com/-name",
		"example.com/name_",
		strings.Repeat("n", 64),
		strings.Repeat("a", 254) + "/x",
	}
	for _, key := range invalidKeys {
		t.Run("key-"+key, func(t *testing.T) {
			goal := validGoal()
			goal.Metadata.Annotations = map[string]string{key: "value"}
			if err := goal.Validate(); err == nil {
				t.Fatalf("invalid annotation key %q was accepted", key)
			}
		})
	}

	invalidValues := []string{
		strings.Repeat("v", 64),
		"-leading",
		"trailing_",
		"has/slash",
		"has space",
		"unicode-☃",
	}
	for _, value := range invalidValues {
		t.Run("value-"+value, func(t *testing.T) {
			goal := validGoal()
			goal.Metadata.Labels = map[string]string{"example.com/key": value}
			if err := goal.Validate(); err == nil {
				t.Fatalf("invalid label value %q was accepted", value)
			}
		})
	}
}

func TestPortableDigestsIgnoreBackendMetadataAndLocalAnnotations(
	t *testing.T,
) {
	first := validGoal()
	first.Metadata.ResourceVersion = "1"
	first.Metadata.CreationTimestamp = "2026-08-30T12:00:00Z"
	first.Metadata.Annotations = map[string]string{
		LocalOwnerRootAnnotation: "out/task",
	}
	second := first
	second.Metadata = first.Metadata
	second.Metadata.ResourceVersion = "backend:opaque/token=="
	second.Metadata.CreationTimestamp = "2027-01-01T00:00:00Z"
	second.Metadata.Annotations = map[string]string{
		LocalOwnerRootAnnotation: "projects/goal",
	}
	firstDigest, err := GoalStateDigest(first)
	if err != nil {
		t.Fatal(err)
	}
	secondDigest, err := GoalStateDigest(second)
	if err != nil {
		t.Fatal(err)
	}
	if firstDigest != secondDigest {
		t.Fatalf("backend metadata changed portable goal digest")
	}

	criteria := GoalCriteria{
		APIVersion: APIVersion,
		Kind:       KindCriteria,
		Metadata: ObjectMeta{
			Name:              first.Metadata.Name,
			ResourceVersion:   "1",
			CreationTimestamp: "2026-08-30T12:00:00Z",
		},
		Spec: CriteriaSpec{
			GoalRef:  GoalReference{Name: first.Metadata.Name},
			Revision: 1,
			Items: []Criterion{{
				CriterionID:    "criterion-001",
				Revision:       1,
				Required:       true,
				Statement:      "The result is verified.",
				EvidenceMethod: "Inspect the evidence.",
			}},
		},
	}
	otherCriteria := criteria
	otherCriteria.Metadata = criteria.Metadata
	otherCriteria.Metadata.ResourceVersion = "opaque"
	otherCriteria.Metadata.CreationTimestamp = ""
	criteriaDigest, err := CriteriaDigest(criteria)
	if err != nil {
		t.Fatal(err)
	}
	otherDigest, err := CriteriaDigest(otherCriteria)
	if err != nil {
		t.Fatal(err)
	}
	if criteriaDigest != otherDigest {
		t.Fatalf("backend metadata changed portable criteria digest")
	}
}

func TestGoalSpecCanBeValidatedBeforeStatusExists(t *testing.T) {
	goal := validGoal()
	goal.Status = GoalStatus{}
	if err := goal.ValidateSpec(); err != nil {
		t.Fatalf("valid desired state: %v", err)
	}
	if err := goal.ValidateStatus(); err == nil {
		t.Fatal("empty observed status was accepted")
	}
	if err := goal.Validate(); err == nil {
		t.Fatal("complete validation accepted empty observed status")
	}
}

func TestMigrationStatusRequiresPathDigestAndMappingProvenance(t *testing.T) {
	goal := validGoal()
	goal.Status.Migration = MigrationStatus{
		SourceFormat:   "unversioned",
		SourcePath:     "out/task/legacy/legacy-goal",
		SourceDigest:   "sha256:" + strings.Repeat("ab", 32),
		MappingVersion: "v1",
		ExtractionMode: "extracted",
		MigratedAt:     "2026-09-01T12:00:00Z",
	}
	if err := goal.Validate(); err != nil {
		t.Fatalf("valid migration status was rejected: %v", err)
	}

	for _, test := range []struct {
		name   string
		mutate func(*MigrationStatus)
	}{
		{
			name: "absolute source path",
			mutate: func(status *MigrationStatus) {
				status.SourcePath = "/var/legacy/legacy-goal"
			},
		},
		{
			name: "empty source path",
			mutate: func(status *MigrationStatus) {
				status.SourcePath = ""
			},
		},
		{
			name: "parent escaping source path",
			mutate: func(status *MigrationStatus) {
				status.SourcePath = "../legacy-goal"
			},
		},
		{
			name: "missing mapping version",
			mutate: func(status *MigrationStatus) {
				status.MappingVersion = ""
			},
		},
		{
			name: "unknown mapping version",
			mutate: func(status *MigrationStatus) {
				status.MappingVersion = "v9"
			},
		},
		{
			name: "unknown extraction mode",
			mutate: func(status *MigrationStatus) {
				status.ExtractionMode = "guessed"
			},
		},
		{
			name: "missing source digest",
			mutate: func(status *MigrationStatus) {
				status.SourceDigest = ""
			},
		},
		{
			name: "malformed source digest",
			mutate: func(status *MigrationStatus) {
				status.SourceDigest = "md5:deadbeef"
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidate := validGoal()
			candidate.Status.Migration = goal.Status.Migration
			test.mutate(&candidate.Status.Migration)
			if err := candidate.Validate(); err == nil {
				t.Fatalf("invalid migration status was accepted: %+v", candidate.Status.Migration)
			}
		})
	}
}

func TestRelationshipFieldNamesAreStableBeforeAlphaRelease(t *testing.T) {
	type expectedField struct {
		name string
		tag  string
	}
	for _, expected := range []expectedField{
		{name: "ParentGoalRef", tag: "parentGoalRef,omitempty"},
		{name: "DependsOnGoalRefs", tag: "dependsOnGoalRefs"},
		{name: "SupersedesGoalRefs", tag: "supersedesGoalRefs"},
	} {
		field, ok := reflect.TypeOf(Relationships{}).FieldByName(expected.name)
		if !ok {
			t.Fatalf("missing Relationships.%s", expected.name)
		}
		if got := field.Tag.Get("json"); got != expected.tag {
			t.Errorf("%s JSON tag = %q, want %q", expected.name, got, expected.tag)
		}
		if got := field.Tag.Get("yaml"); got != expected.tag {
			t.Errorf("%s YAML tag = %q, want %q", expected.name, got, expected.tag)
		}
	}
}

func TestRelationshipsNormalizeWithoutMutatingInput(t *testing.T) {
	parent := GoalReference{Name: "parent"}
	relationships := Relationships{
		ParentGoalRef: &parent,
		DependsOnGoalRefs: []GoalReference{
			{Name: "z-last"},
			{Name: "a-first"},
		},
		SupersedesGoalRefs: nil,
	}
	normalized := relationships.Normalized()
	wantDependencies := []GoalReference{
		{Name: "a-first"},
		{Name: "z-last"},
	}
	if !reflect.DeepEqual(normalized.DependsOnGoalRefs, wantDependencies) {
		t.Fatalf(
			"normalized dependencies = %#v, want %#v",
			normalized.DependsOnGoalRefs,
			wantDependencies,
		)
	}
	if normalized.SupersedesGoalRefs == nil ||
		len(normalized.SupersedesGoalRefs) != 0 {
		t.Fatalf("nil supersession list was not normalized to empty")
	}
	if relationships.DependsOnGoalRefs[0].Name != "z-last" {
		t.Fatal("normalization mutated the input relationship order")
	}
	normalized.ParentGoalRef.Name = "changed"
	if relationships.ParentGoalRef.Name != "parent" {
		t.Fatal("normalization aliased the input parent reference")
	}
}

func TestRelationshipOrderAndNilDoNotChangeGoalStateDigest(t *testing.T) {
	first := validGoal()
	first.Spec.Relationships = Relationships{
		DependsOnGoalRefs: []GoalReference{
			{Name: "goal-z"},
			{Name: "goal-a"},
		},
		SupersedesGoalRefs: nil,
	}
	second := first
	second.Spec.Relationships = Relationships{
		DependsOnGoalRefs: []GoalReference{
			{Name: "goal-a"},
			{Name: "goal-z"},
		},
		SupersedesGoalRefs: []GoalReference{},
	}
	firstDigest, err := GoalStateDigest(first)
	if err != nil {
		t.Fatal(err)
	}
	secondDigest, err := GoalStateDigest(second)
	if err != nil {
		t.Fatal(err)
	}
	if firstDigest != secondDigest {
		t.Fatalf(
			"semantic-equivalent relationship order changed digest: %s != %s",
			firstDigest,
			secondDigest,
		)
	}
}

func TestRelationshipValidationIsBoundedAndDeterministic(t *testing.T) {
	goal := validGoal()
	goal.Spec.Relationships.DependsOnGoalRefs = make(
		[]GoalReference,
		MaxGoalRelationshipReferences,
	)
	for index := range goal.Spec.Relationships.DependsOnGoalRefs {
		goal.Spec.Relationships.DependsOnGoalRefs[index] = GoalReference{
			Name: fmt.Sprintf("dependency-%03d", index),
		}
	}
	if err := goal.ValidateSpec(); err != nil {
		t.Fatalf("relationship bound should be inclusive: %v", err)
	}
	goal.Spec.Relationships.DependsOnGoalRefs = append(
		goal.Spec.Relationships.DependsOnGoalRefs,
		GoalReference{Name: "too-many"},
	)
	if err := goal.ValidateSpec(); err == nil ||
		!strings.Contains(err.Error(), "cardinality") {
		t.Fatalf("over-bound relationships error = %v", err)
	}

	duplicateErrors := []string{}
	for _, references := range [][]GoalReference{
		{{Name: "duplicate"}, {Name: "other"}, {Name: "duplicate"}},
		{{Name: "duplicate"}, {Name: "duplicate"}, {Name: "other"}},
	} {
		candidate := validGoal()
		candidate.Spec.Relationships.DependsOnGoalRefs = references
		err := candidate.ValidateSpec()
		if err == nil {
			t.Fatal("duplicate relationship was accepted")
		}
		duplicateErrors = append(duplicateErrors, err.Error())
	}
	if duplicateErrors[0] != duplicateErrors[1] {
		t.Fatalf("duplicate errors depend on input order: %#v", duplicateErrors)
	}

	self := validGoal()
	self.Spec.Relationships.SupersedesGoalRefs = []GoalReference{{
		Name: self.Metadata.Name,
	}}
	if err := self.ValidateSpec(); err == nil ||
		!strings.Contains(err.Error(), "self-reference") {
		t.Fatalf("self-reference error = %v", err)
	}
}

func validAttempt(envelopeName string) GoalAttempt {
	return GoalAttempt{
		APIVersion: APIVersion,
		Kind:       KindAttempt,
		Metadata:   ObjectMeta{Name: envelopeName},
		Spec: AttemptSpec{
			GoalRef:             GoalReference{Name: "portable-goal"},
			GoalGeneration:      1,
			LifecycleGeneration: 1,
			CriteriaRevision:    1,
			CriteriaDigest:      "sha256:" + strings.Repeat("11", 32),
			GoalStateDigest:     "sha256:" + strings.Repeat("22", 32),
			WorkType:            "change",
		},
		Status: AttemptStatus{
			State:      "open",
			ObservedAt: "2026-08-30T12:00:00Z",
			Artifacts: ArtifactManifest{
				PlanDigest:   "sha256:" + strings.Repeat("33", 32),
				ResultDigest: "sha256:" + strings.Repeat("44", 32),
			},
		},
	}
}

func TestAttemptResumeFieldsValidateAndBound(t *testing.T) {
	attempt := validAttempt("resume-attempt")
	attempt.Spec.StableDefect = "The goal catalog omits resume state for open goals."
	attempt.Spec.Hypothesis = "A structured continuation packet is needed."
	attempt.Spec.Subject = "projects/goal"
	attempt.Spec.AffectedCriteria = []string{"criterion-001", "criterion-002"}
	attempt.Spec.RegressionRefs = []string{"goal-resume-regression"}
	attempt.Spec.PriorAttemptID = "attempt-1"
	attempt.Spec.DominantFailure = "Fresh agents cannot resume open goals."
	attempt.Spec.MeasurableDelta = "Catalog lists every open goal."
	attempt.Spec.NextAction = "Promote the continuation packet."
	attempt.Spec.Blocker = "Awaiting acceptance."
	attempt.Spec.ResumeCondition = "The catalog check passes."
	goal := validGoal()
	goal.Status.CriteriaRevision = 1
	goal.Metadata.Generation = 1
	if err := attempt.ValidateForGoal(goal); err != nil {
		t.Fatalf("valid resume fields rejected: %v", err)
	}

	for _, test := range []struct {
		name   string
		mutate func(*GoalAttempt)
	}{
		{
			name: "untrimmed stable defect",
			mutate: func(a *GoalAttempt) {
				a.Spec.StableDefect = " leading"
			},
		},
		{
			name: "oversized next action",
			mutate: func(a *GoalAttempt) {
				a.Spec.NextAction = strings.Repeat("x", MaxResumeFieldBytes+1)
			},
		},
		{
			name: "NUL in resume condition",
			mutate: func(a *GoalAttempt) {
				a.Spec.ResumeCondition = "a\x00b"
			},
		},
		{
			name: "invalid prior attempt ID",
			mutate: func(a *GoalAttempt) {
				a.Spec.PriorAttemptID = "Not-A-Valid-ID!"
			},
		},
		{
			name: "unsorted affected criteria",
			mutate: func(a *GoalAttempt) {
				a.Spec.AffectedCriteria = []string{"criterion-002", "criterion-001"}
			},
		},
		{
			name: "duplicate regression refs",
			mutate: func(a *GoalAttempt) {
				a.Spec.RegressionRefs = []string{"dup", "dup"}
			},
		},
		{
			name: "oversized affected criteria",
			mutate: func(a *GoalAttempt) {
				a.Spec.AffectedCriteria = make([]string, MaxAffectedCriteria+1)
				for i := range a.Spec.AffectedCriteria {
					a.Spec.AffectedCriteria[i] = fmt.Sprintf("c-%d", i)
				}
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			modified := attempt
			test.mutate(&modified)
			if err := modified.ValidateForGoal(goal); err == nil {
				t.Fatalf("invalid resume fields were accepted")
			}
		})
	}
}
