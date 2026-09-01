package fsstore

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDoctorReportsStableForValidGoal confirms the default healthy state.
func TestDoctorReportsStableForValidGoal(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "doctor-stable")
	result, err := store.Doctor(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.PublicationState != "stable" {
		t.Fatalf("PublicationState = %q, want stable", result.PublicationState)
	}
}

// TestCheckpointLeavesRecoverableIntentOnPublishFailure injects a failure at
// the final README rename and verifies doctor classifies the record as
// incomplete-recoverable and recover completes the intended record.
func TestCheckpointLeavesRecoverableIntentOnPublishFailure(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "recoverable-checkpoint")
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
		t.Fatalf("Checkpoint() error = %v, want PublicationIncompleteError", err)
	}
	if reference.ResourceVersion != "2" {
		t.Fatalf("reference.ResourceVersion = %q, want 2", reference.ResourceVersion)
	}
	store.beforeRename = nil

	// Doctor must report a recoverable intent, not a generic invalid record.
	result, err := store.Doctor(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.PublicationState != "partial-intent" {
		t.Fatalf(
			"Doctor() PublicationState = %q, want partial-intent",
			result.PublicationState,
		)
	}
	if result.OperationID == "" {
		t.Fatal("Doctor() did not return the pending operation ID")
	}

	// A normal mutation must fail closed with the stable error.
	_, err = store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "2",
		Execution:               "active",
	})
	if err == nil || !strings.Contains(err.Error(), "publication is incomplete") {
		t.Fatalf("Checkpoint() over pending intent error = %v, want publication-incomplete", err)
	}

	// Recover must complete the intended state.
	recovered, err := store.Recover(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if recovered.Discarded {
		t.Fatal("Recover() discarded instead of replaying the intended state")
	}
	if recovered.ResourceVersion != "2" {
		t.Fatalf("Recover() ResourceVersion = %q, want 2", recovered.ResourceVersion)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() after recover: %v", err)
	}
	goal, err := store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Metadata.ResourceVersion != "2" || goal.Status.Execution != "paused" {
		t.Fatalf("recovered goal state = %+v, want version 2 paused", goal)
	}
	if _, pending := store.readPublicationIntent(goalDir); pending != nil {
		t.Fatal("publication intent was not removed after recovery")
	}
}

// TestRecoverDiscardsIntentWhenNothingPublished simulates an interruption
// before any canonical rename: checkpoint discards the intent automatically,
// and recover discards a stale intent left by a crash that removed staging.
func TestRecoverDiscardsIntentWhenNothingPublished(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "discard-intent")
	// A pre-commit failure (before the first goal.yaml rename) must discard
	// the intent and return the raw injected error, preserving the prior
	// record exactly.
	goalPath := filepath.Join(goalDir, "goal.yaml")
	store.beforeRename = func(target string) error {
		if target == goalPath {
			return errors.New("simulated process failure before first rename")
		}
		return nil
	}
	_, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		Execution:               "paused",
	})
	if err == nil || !strings.Contains(
		err.Error(),
		"simulated process failure before first rename",
	) {
		t.Fatalf("Checkpoint() error = %v, want raw simulated failure", err)
	}
	store.beforeRename = nil
	if intent, readErr := store.readPublicationIntent(goalDir); readErr != nil || intent != nil {
		t.Fatalf("pre-commit failure left an intent, intent=%v readErr=%v", intent, readErr)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() after pre-commit failure: %v", err)
	}
	goal, err := store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Metadata.ResourceVersion != "1" ||
		goal.Status.Execution != "active" ||
		goal.Status.ActiveAttemptID != "" {
		t.Fatalf("pre-commit failure changed the prior record: %+v", goal)
	}

	// Recover must discard an intent left behind by a crash that removed the
	// staged after-images before any canonical rename. The recorded before
	// digest matches the untouched goal.yaml so the intent is classified as
	// having published nothing.
	goalBytes, err := os.ReadFile(goalPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.beginPublication(
		goalDir,
		"discard-intent",
		"1",
		"2",
		[]publicationFileEntry{
			{
				Path:         "goal.yaml",
				BeforeDigest: digestBytes(goalBytes),
				Content:      []byte(string(goalBytes) + "\n"),
			},
		},
		nil,
	); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(goalDir, publicationStageDir)); err != nil {
		t.Fatal(err)
	}
	result, err := store.Doctor(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if result.PublicationState != "discardable-intent" {
		t.Fatalf(
			"Doctor() PublicationState = %q, want discardable-intent",
			result.PublicationState,
		)
	}
	recovered, err := store.Recover(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if !recovered.Discarded {
		t.Fatal("Recover() did not discard the intent")
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() after discard: %v", err)
	}
	goal, err = store.readGoalManifest(goalDir)
	if err != nil {
		t.Fatal(err)
	}
	if goal.Metadata.ResourceVersion != "1" {
		t.Fatalf("goal resourceVersion after discard = %q, want prior 1", goal.Metadata.ResourceVersion)
	}
}
