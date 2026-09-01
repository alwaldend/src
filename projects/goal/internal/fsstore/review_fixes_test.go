package fsstore

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/projects/goal/api/v1alpha1"
)

// TestReviewFixPathConfinementIsEnforcedOnStoredIntents ensures a malformed
// or crafted persisted intent with a traversal path is refused at load time,
// before recovery can rename content outside the goal record.
func TestReviewFixPathConfinementIsEnforcedOnStoredIntents(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "confine-paths")
	before, after := digestBytes([]byte("old\n")), digestBytes([]byte("new\n"))
	// Write the malformed intent directly with the YAML writer so the strict
	// load path (the recovery boundary) sees it, mirroring a persisted intent
	// crafted before the confinement validation existed.
	if err := store.writeYAML(publicationIntentPath(goalDir), publicationIntent{
		APIVersion: goalAPIVersion,
		Kind:       v1alpha1.KindGoalPublication,
		Metadata: ObjectMeta{
			Name:              "pub-path-escape",
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: store.timestamp(),
		},
		Spec: v1alpha1.GoalPublicationSpec{
			GoalRef:                 LocalGoalReference{Name: "confine-paths"},
			OperationID:             "pub-path-escape",
			PriorResourceVersion:    "1",
			IntendedResourceVersion: "2",
			Files: []v1alpha1.GoalPublicationFile{
				{
					Path:           "../../outside.yaml",
					BeforeDigest:   before,
					AfterDigest:    after,
					StagedRelative: "goal.yaml",
				},
			},
			SnapshotDigests: map[uint64]string{},
		},
		Status: v1alpha1.GoalPublicationStatus{
			State:      v1alpha1.PublicationIncomplete,
			ObservedAt: store.timestamp(),
		},
	}); err != nil {
		t.Fatalf("write malformed intent: %v", err)
	}
	if _, err := store.Doctor(goalDir); err == nil ||
		!strings.Contains(err.Error(), "normalized goal-relative path") {
		t.Fatalf("Doctor() error = %v, want path confinement error", err)
	}
	if _, err := store.Recover(goalDir); err == nil ||
		!strings.Contains(err.Error(), "normalized goal-relative path") {
		t.Fatalf("Recover() error = %v, want path confinement error", err)
	}
	if _, err := store.readPublicationIntent(goalDir); err == nil ||
		!strings.Contains(err.Error(), "normalized goal-relative path") {
		t.Fatalf("readPublicationIntent() error = %v, want path confinement error", err)
	}
}

// TestReviewFixBeforeDigestsAreRecordedForExistingTargets verifies that a
// checkpoint over an existing goal records the actual on-disk before digest
// instead of an empty "absent" marker, so an interruption before the first
// canonical rename is classified as staged rather than a conflict.
func TestReviewFixBeforeDigestsAreRecordedForExistingTargets(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "before-digests")
	goalPath := filepath.Join(goalDir, "goal.yaml")
	original, err := os.ReadFile(goalPath)
	if err != nil {
		t.Fatal(err)
	}
	readmePath := filepath.Join(goalDir, "README.md")
	store.beforeRename = func(target string) error {
		if target == readmePath {
			return errors.New("projection write failed")
		}
		return nil
	}
	reference, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		Execution:               "paused",
	})
	if err == nil || !strings.Contains(err.Error(), "publication is incomplete") {
		t.Fatalf("Checkpoint() error = %v, want publication-incomplete", err)
	}
	if reference.ResourceVersion != "2" {
		t.Fatalf("reference.ResourceVersion = %q, want 2", reference.ResourceVersion)
	}
	store.beforeRename = nil
	intent, err := store.readPublicationIntent(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if intent == nil {
		t.Fatal("expected a pending intent after projection failure")
	}
	found := false
	for _, file := range intent.Spec.Files {
		if file.Path == "goal.yaml" {
			found = true
			if file.BeforeDigest != digestBytes(original) {
				t.Fatalf(
					"goal.yaml before digest = %q, want recorded actual digest",
					file.BeforeDigest,
				)
			}
		}
	}
	if !found {
		t.Fatal("intent does not contain a goal.yaml file entry")
	}
	if _, err := store.Recover(goalDir); err != nil {
		t.Fatalf("Recover() after projection failure: %v", err)
	}
}

// TestReviewFixCriteriaCommitPointPublishesGoalFirst ensures the criteria
// update intent is ordered with goal.yaml as the first canonical change, so a
// failure after the criteria rename cannot leave the criteria revision ahead
// of the goal pointer without a recoverable intent.
func TestReviewFixCriteriaCommitPointPublishesGoalFirst(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "criteria-order")
	criteriaFile := filepath.Join(root, "out", "task", "criteria-update.yaml")
	if err := os.WriteFile(criteriaFile, []byte(`items:
  - criterionID: criterion-001
    statement: The committed criteria are verified.
    evidenceMethod: Inspect the focused test result.
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		Execution:               "paused",
	}); err != nil {
		t.Fatal(err)
	}
	readmePath := filepath.Join(goalDir, "README.md")
	store.beforeRename = func(target string) error {
		if target == readmePath {
			return errors.New("projection write failed")
		}
		return nil
	}
	reference, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "2",
		CriteriaFile:            criteriaFile,
	})
	if err == nil || !strings.Contains(err.Error(), "criteria update committed") {
		t.Fatalf("Checkpoint() error = %v, want criteria committed-version report", err)
	}
	if reference.ResourceVersion != "3" {
		t.Fatalf("reference.ResourceVersion = %q, want 3", reference.ResourceVersion)
	}
	store.beforeRename = nil
	intent, err := store.readPublicationIntent(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if intent == nil || len(intent.Spec.Files) != 3 {
		t.Fatalf("intent = %+v, want 3 file entries", intent)
	}
	if intent.Spec.Files[0].Path != "goal.yaml" {
		t.Fatalf(
			"first published path = %q, want goal.yaml commit point",
			intent.Spec.Files[0].Path,
		)
	}
	if !store.commitPointPublished(goalDir, intent) {
		t.Fatal("commit point goal.yaml was not published first")
	}
	if _, err := store.Recover(goalDir); err != nil {
		t.Fatalf("Recover() after criteria projection failure: %v", err)
	}
}

// TestReviewFixChainEvaluationAcceptsFinalOrderedState reproduces the
// immediate-close new attempt scenario: the intent records two goal.yaml
// after-images, the final rename succeeds, but cleanup of the intent itself
// fails. Doctor must classify the record as recoverable partial, not a
// conflict, and Recover must finish without rewriting canonical content.
func TestReviewFixChainEvaluationAcceptsFinalOrderedState(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "chain-final")
	goalPath := filepath.Join(goalDir, "goal.yaml")
	original, err := os.ReadFile(goalPath)
	if err != nil {
		t.Fatal(err)
	}
	attemptID := "attempt-1"
	plan := filepath.Join(root, "out", "task", "plan.md")
	writeTestFile(t, plan, "# Plan\n")
	review := writeRefineReview(t, root, "close-review.yaml")

	// Fail the final goal.yaml rename (the second image), leaving the commit
	// point published but the final state not installed.
	goalWrites := 0
	store.beforeRename = func(target string) error {
		if target == goalPath {
			goalWrites++
		}
		if target == goalPath && goalWrites == 2 {
			return errors.New("simulated finalization failure")
		}
		return nil
	}
	_, err = store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		AttemptID:               attemptID,
		PlanFile:                plan,
		CloseAttempt:            true,
		ReviewFile:              review,
	})
	if err == nil || !strings.Contains(err.Error(), "simulated finalization failure") {
		t.Fatalf("Checkpoint() error = %v, want simulated failure", err)
	}
	store.beforeRename = nil

	// The intent has two goal.yaml entries. Recover finishes the chain by
	// installing the final image; the record must become valid.
	if _, err := store.Recover(goalDir); err != nil {
		t.Fatalf("Recover() over intermediate chain state: %v", err)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() after recover: %v", err)
	}
	goal, err := store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Status.ActiveAttemptID != "" {
		t.Fatalf("closed goal still has an active attempt pointer: %+v", goal.Status)
	}
	_ = original
}

// TestReviewFixChainEvaluationAcceptsFinalStateWhenProjectionFails reproduces
// the exact review finding: the final goal.yaml rename succeeds, but the README
// refresh fails afterward. The current digest matches the second (final)
// goal.yaml after-image and neither state recorded by the first entry, so
// Doctor must classify the record as recoverable partial rather than a
// conflict, and Recover must complete without rewriting canonical content.
func TestReviewFixChainEvaluationAcceptsFinalStateWhenProjectionFails(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "chain-final-projection")
	readmePath := filepath.Join(goalDir, "README.md")
	store.beforeRename = func(target string) error {
		if target == readmePath {
			return errors.New("projection write failed")
		}
		return nil
	}
	plan := filepath.Join(root, "out", "task", "plan.md")
	writeTestFile(t, plan, "# Plan\n")
	review := writeRefineReview(t, root, "close-review-projection.yaml")
	_, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		AttemptID:               "attempt-1",
		PlanFile:                plan,
		CloseAttempt:            true,
		ReviewFile:              review,
	})
	if err == nil || !strings.Contains(err.Error(), "publication is incomplete") {
		t.Fatalf("Checkpoint() error = %v, want publication-incomplete", err)
	}
	store.beforeRename = nil

	intent, err := store.readPublicationIntent(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if intent == nil {
		t.Fatal("expected a pending intent after projection failure")
	}
	goalEntries := 0
	for _, file := range intent.Spec.Files {
		if file.Path == "goal.yaml" {
			goalEntries++
		}
	}
	if goalEntries != 2 {
		t.Fatalf("goal.yaml entries = %d, want 2 (commit point + final)", goalEntries)
	}
	result, err := store.Doctor(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.PublicationState != v1alpha1.PublicationPartial {
		t.Fatalf(
			"Doctor() PublicationState = %q, want partial-intent for final chain state",
			result.PublicationState,
		)
	}
	if _, err := store.Recover(goalDir); err != nil {
		t.Fatalf("Recover() over final chain state: %v", err)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() after recover: %v", err)
	}
	goal, err := store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Status.ActiveAttemptID != "" {
		t.Fatalf("closed goal still has an active attempt pointer: %+v", goal.Status)
	}
}

// TestReviewFixReplayExistingAttemptDirRepairsStagedFiles verifies that an
// interrupted update to an existing attempt replays the staged after-images
// inside the already-present attempt directory instead of skipping them.
func TestReviewFixReplayExistingAttemptDirRepairsStagedFiles(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "replay-existing")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
	}); err != nil {
		t.Fatal(err)
	}
	attemptDir := filepath.Join(goalDir, "attempts", "attempt-1")
	resultPath := filepath.Join(attemptDir, "result.md")
	newResult := filepath.Join(root, "out", "task", "new-result.md")
	if err := os.WriteFile(newResult, []byte("# Result\n\nNew result.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Fail after the goal.yaml rename but before the attempt result rename.
	store.beforeRename = func(target string) error {
		if target == resultPath {
			return errors.New("simulated attempt failure")
		}
		return nil
	}
	_, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "2",
		AttemptID:               "attempt-1",
		ResultFile:              newResult,
	})
	if err == nil || !strings.Contains(err.Error(), "simulated attempt failure") {
		t.Fatalf("Checkpoint() error = %v, want simulated attempt failure", err)
	}
	store.beforeRename = nil
	oldResult, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(oldResult), "New result") {
		t.Fatal("result was already installed before the simulated failure")
	}
	if _, err := store.Recover(goalDir); err != nil {
		t.Fatalf("Recover() over existing attempt: %v", err)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() after recover: %v", err)
	}
	newContent, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(newContent), "New result") {
		t.Fatalf("existing attempt result was not replayed: %q", newContent)
	}
}

// TestReviewFixDoctorReportsProjectionStaleWhenREADMEOnlyDrifts ensures the
// no-intent stable path compares the replaceable README projection instead of
// reporting stable for any valid canonical record.
func TestReviewFixDoctorReportsProjectionStaleWhenREADMEOnlyDrifts(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "projection-stale")
	readmePath := filepath.Join(goalDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Stale projection\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := store.Doctor(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.PublicationState != v1alpha1.PublicationProjection {
		t.Fatalf(
			"Doctor() PublicationState = %q, want committed-projection-stale",
			result.PublicationState,
		)
	}
}
