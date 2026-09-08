package fsstore

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExistingAttemptCheckpointPublishesGoalTokenFirst(t *testing.T) {
	for _, test := range []struct {
		name              string
		failureTarget     string
		wantVersion       string
		wantResultChanged bool
		wantCommitted     bool
	}{
		{
			name:              "goal rename failure changes no attempt data",
			failureTarget:     "goal.yaml",
			wantVersion:       "2",
			wantResultChanged: false,
		},
		{
			name:              "later failure leaves the advanced token",
			failureTarget:     "attempt.yaml",
			wantVersion:       "3",
			wantResultChanged: true,
			wantCommitted:     true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, root := newTestStore(t)
			goalDir := initTestGoal(t, store, root, "ordered-checkpoint")
			if _, err := store.Checkpoint(CheckpointOptions{
				GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
			}); err != nil {
				t.Fatal(err)
			}

			attemptDir := filepath.Join(goalDir, "attempts", "attempt-1")
			resultPath := filepath.Join(attemptDir, "result.md")
			attemptPath := filepath.Join(attemptDir, "attempt.yaml")
			oldResult, err := os.ReadFile(resultPath)
			if err != nil {
				t.Fatal(err)
			}
			oldAttempt, err := os.ReadFile(attemptPath)
			if err != nil {
				t.Fatal(err)
			}
			newResultPath := filepath.Join(root, "out", "task", "new-result.md")
			newResult := []byte("# Result\n\nNew evidence-backed result.\n")
			if err := os.WriteFile(newResultPath, newResult, 0o644); err != nil {
				t.Fatal(err)
			}

			store.beforeRename = func(target string) error {
				if filepath.Base(target) == test.failureTarget {
					return errors.New("simulated process failure")
				}
				return nil
			}
			reference, err := store.Checkpoint(CheckpointOptions{
				GoalDir:                 goalDir,
				ExpectedResourceVersion: "2",
				AttemptID:               "attempt-1",
				ResultFile:              newResultPath,
			})
			if err == nil || !strings.Contains(err.Error(), "simulated process failure") {
				t.Fatalf("Checkpoint() error = %v, want simulated failure", err)
			}
			if test.wantCommitted {
				if reference.GoalID != "ordered-checkpoint" ||
					reference.GoalRef != "." || reference.ResourceVersion != test.wantVersion {
					t.Fatalf(
						"Checkpoint() reference = %+v, want committed version %s",
						reference,
						test.wantVersion,
					)
				}
				if !strings.Contains(
					err.Error(),
					"checkpoint committed at resourceVersion "+test.wantVersion,
				) {
					t.Fatalf("Checkpoint() error = %v, want committed-version report", err)
				}
			} else if reference != (GoalReference{}) {
				t.Fatalf("Checkpoint() reference = %+v, want empty pre-commit reference", reference)
			}
			store.beforeRename = nil

			goal, err := store.readGoalManifest(goalDir)
			if err != nil {
				t.Fatal(err)
			}
			if goal.Metadata.ResourceVersion != test.wantVersion {
				t.Fatalf(
					"goal resourceVersion = %s, want %s",
					goal.Metadata.ResourceVersion,
					test.wantVersion,
				)
			}
			result, err := os.ReadFile(resultPath)
			if err != nil {
				t.Fatal(err)
			}
			if changed := !bytes.Equal(result, oldResult); changed != test.wantResultChanged {
				t.Fatalf("result changed = %t, want %t", changed, test.wantResultChanged)
			}
			attempt, err := os.ReadFile(attemptPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(attempt, oldAttempt) {
				t.Fatal("attempt manifest changed despite failure before its rename")
			}
			if test.wantResultChanged && !bytes.Equal(result, newResult) {
				t.Fatalf("published result = %q, want %q", result, newResult)
			}
			if test.wantResultChanged {
				if err := store.ValidateGoal(goalDir); err == nil {
					t.Fatal("partially published existing attempt did not fail closed")
				}
			}
		})
	}
}

func TestNewAttemptCheckpointPublishesGoalTokenFirst(t *testing.T) {
	for _, test := range []struct {
		name                   string
		closeAttempt           bool
		failGoalWrite          bool
		failAttemptPublication bool
		failGoalFinalization   bool
		wantVersion            string
		wantCommitted          bool
		wantAttempt            bool
	}{
		{
			name:          "goal rename failure publishes nothing",
			failGoalWrite: true,
			wantVersion:   "1",
		},
		{
			name:                   "open attempt publication failure advances token",
			failAttemptPublication: true,
			wantVersion:            "2",
			wantCommitted:          true,
		},
		{
			name:                   "immediate close publication failure advances token",
			closeAttempt:           true,
			failAttemptPublication: true,
			wantVersion:            "2",
			wantCommitted:          true,
		},
		{
			name:                 "immediate close finalization failure advances token",
			closeAttempt:         true,
			failGoalFinalization: true,
			wantVersion:          "2",
			wantCommitted:        true,
			wantAttempt:          true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, root := newTestStore(t)
			goalID := "new-attempt-order"
			goalDir := initTestGoal(t, store, root, goalID)
			attemptID := "attempt-1"
			attemptDir := filepath.Join(goalDir, "attempts", attemptID)
			goalPath := filepath.Join(goalDir, "goal.yaml")
			options := CheckpointOptions{
				GoalDir:                 goalDir,
				ExpectedResourceVersion: "1",
				AttemptID:               attemptID,
			}
			if test.closeAttempt {
				options.CloseAttempt = true
				options.ReviewFile = writeRefineReview(t, root, "new-attempt-review.yaml")
			}
			goalWrites := 0
			store.beforeRename = func(target string) error {
				if target == goalPath {
					goalWrites++
				}
				if (test.failGoalWrite && target == goalPath && goalWrites == 1) ||
					(test.failAttemptPublication && target == attemptDir) ||
					(test.failGoalFinalization && target == goalPath && goalWrites == 2) {
					return errors.New("simulated process failure")
				}
				return nil
			}
			reference, err := store.Checkpoint(options)
			if err == nil || !strings.Contains(err.Error(), "simulated process failure") {
				t.Fatalf("Checkpoint() error = %v, want simulated failure", err)
			}
			if test.wantCommitted {
				if reference.GoalID != goalID || reference.GoalRef != "." ||
					reference.ResourceVersion != test.wantVersion {
					t.Fatalf(
						"Checkpoint() reference = %+v, want committed version %s",
						reference,
						test.wantVersion,
					)
				}
				if !strings.Contains(
					err.Error(),
					"checkpoint committed at resourceVersion "+test.wantVersion,
				) {
					t.Fatalf("Checkpoint() error = %v, want committed-version report", err)
				}
			} else if reference != (GoalReference{}) {
				t.Fatalf("Checkpoint() reference = %+v, want empty pre-commit reference", reference)
			}
			store.beforeRename = nil

			goal, err := store.readGoalManifest(goalDir)
			if err != nil {
				t.Fatal(err)
			}
			if goal.Metadata.ResourceVersion != test.wantVersion {
				t.Fatalf(
					"goal resourceVersion = %s, want %s",
					goal.Metadata.ResourceVersion,
					test.wantVersion,
				)
			}
			if exists := pathExists(attemptDir); exists != test.wantAttempt {
				t.Fatalf("canonical attempt exists = %t, want %t", exists, test.wantAttempt)
			}
			if test.wantCommitted {
				if err := store.ValidateGoal(goalDir); err == nil {
					t.Fatal("partially committed new attempt did not fail closed")
				}
			} else if err := store.ValidateGoal(goalDir); err != nil {
				t.Fatalf("pre-commit failure changed the valid record: %v", err)
			}
		})
	}
}

func TestCheckpointREADMEFailureReportsCommittedVersion(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "checkpoint-readme-failure")
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
	if err == nil || !strings.Contains(
		err.Error(),
		"checkpoint committed at resourceVersion 2",
	) {
		t.Fatalf("Checkpoint() error = %v, want committed-version report", err)
	}
	if reference.GoalID != "checkpoint-readme-failure" ||
		reference.GoalRef != "." || reference.ResourceVersion != "2" {
		t.Fatalf("Checkpoint() reference = %+v, want committed version 2", reference)
	}
	store.beforeRename = nil
	// A pending intent blocks the retry with the stable gate, then recover
	// completes the intended state and the stale token is rejected.
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		Execution:               "waiting",
	}); err == nil || !strings.Contains(err.Error(), "publication is incomplete") {
		t.Fatalf("retry over pending intent did not fail closed: %v", err)
	}
	recovered, err := store.Recover(goalDir)
	if err != nil {
		t.Fatalf("Recover() after README failure: %v", err)
	}
	if recovered.ResourceVersion != "2" {
		t.Fatalf("Recover() ResourceVersion = %q, want 2", recovered.ResourceVersion)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		Execution:               "waiting",
	}); err == nil || !strings.Contains(err.Error(), "stale resourceVersion") {
		t.Fatalf("retry with pre-commit resourceVersion was not stale: %v", err)
	}
}

func TestCriteriaREADMEFailureReportsCommittedVersion(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "criteria-readme-failure")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "1",
		Execution:               "paused",
	}); err != nil {
		t.Fatal(err)
	}
	desiredPath := filepath.Join(root, "out", "task", "criteria-update.yaml")
	if err := os.WriteFile(desiredPath, []byte(`items:
  - criterionID: criterion-001
    statement: The committed criteria are verified.
    evidenceMethod: Inspect the focused test result.
`), 0o644); err != nil {
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
		CriteriaFile:            desiredPath,
	})
	if err == nil || !strings.Contains(
		err.Error(),
		"criteria update committed at resourceVersion 3",
	) {
		t.Fatalf("Checkpoint() error = %v, want committed-version report", err)
	}
	if reference.GoalID != "criteria-readme-failure" ||
		reference.GoalRef != "." || reference.ResourceVersion != "3" {
		t.Fatalf("Checkpoint() reference = %+v, want committed version 3", reference)
	}
	store.beforeRename = nil
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir:                 goalDir,
		ExpectedResourceVersion: "2",
		CriteriaFile:            desiredPath,
	}); err == nil || !strings.Contains(err.Error(), "publication is incomplete") {
		t.Fatalf("retry over pending intent did not fail closed: %v", err)
	}
	recovered, err := store.Recover(goalDir)
	if err != nil {
		t.Fatalf("Recover() after criteria README failure: %v", err)
	}
	if recovered.ResourceVersion != "3" {
		t.Fatalf("Recover() ResourceVersion = %q, want 3", recovered.ResourceVersion)
	}
}

func TestValidationCleansOnlyRecognizedTemporaryResidue(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "temporary-residue")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
	}); err != nil {
		t.Fatal(err)
	}
	attemptDir := filepath.Join(goalDir, "attempts", "attempt-1")
	paths := []string{
		filepath.Join(goalDir, ".goal-write-root123"),
		filepath.Join(goalDir, "criteria-revisions", ".goal-write-snapshot123"),
		filepath.Join(goalDir, "criteria-revisions", ".goal-immutable-snapshot123"),
		filepath.Join(attemptDir, ".goal-write-attempt123"),
		filepath.Join(attemptDir, "evidence", ".goal-write-evidence123"),
	}
	for _, path := range paths {
		if err := os.WriteFile(path, []byte("interrupted temporary file"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	temporaryAttempt := filepath.Join(goalDir, "attempts", ".goal-attempt-staging123")
	if err := os.MkdirAll(filepath.Join(temporaryAttempt, "evidence"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(temporaryAttempt, "partial"),
		[]byte("partial attempt"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("ValidateGoal() with internal temporary residue: %v", err)
	}
	for _, path := range append(paths, temporaryAttempt) {
		if pathExists(path) {
			t.Errorf("temporary residue was not removed: %s", path)
		}
	}
}

func TestMissingREADMEIsValidAndRenderRecreatesIt(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "replaceable-readme")
	readmePath := filepath.Join(goalDir, "README.md")
	if err := os.Remove(readmePath); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("replaceable README absence invalidated the record: %v", err)
	}
	if err := store.Render(goalDir, "1", 20); err != nil {
		t.Fatalf("Render() did not recreate a missing README: %v", err)
	}
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "# Goal replaceable-readme") {
		t.Fatalf("unexpected rendered README: %q", content)
	}
}

func TestNextCriteriaRevisionRejectsGrowthBeyondBound(t *testing.T) {
	next, err := nextCriteriaRevision(maxCriteriaRevisions - 1)
	if err != nil {
		t.Fatal(err)
	}
	if next != maxCriteriaRevisions {
		t.Fatalf("next criteria revision = %d, want %d", next, maxCriteriaRevisions)
	}
	if _, err := nextCriteriaRevision(maxCriteriaRevisions); err == nil ||
		!strings.Contains(err.Error(), "cannot exceed") {
		t.Fatalf("maximum criteria revision was allowed to grow: %v", err)
	}
}

func TestAttemptRootAndEvidenceLayoutsAreExact(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "exact-attempt-layout")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
	}); err != nil {
		t.Fatal(err)
	}
	attemptDir := filepath.Join(goalDir, "attempts", "attempt-1")
	extra := filepath.Join(attemptDir, "notes.md")
	if err := os.WriteFile(extra, []byte("unexpected"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateGoal(goalDir); err == nil ||
		!strings.Contains(err.Error(), "must contain only") {
		t.Fatalf("unexpected attempt-root file was accepted: %v", err)
	}
	if err := os.Remove(extra); err != nil {
		t.Fatal(err)
	}
	nonMarkdown := filepath.Join(attemptDir, "evidence", "evidence.txt")
	if err := os.WriteFile(nonMarkdown, []byte("unexpected"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateGoal(goalDir); err == nil ||
		!strings.Contains(err.Error(), "invalid evidence entry") {
		t.Fatalf("non-Markdown evidence was accepted: %v", err)
	}
}

func TestAttemptWithoutEvidenceSurvivesMissingEmptyDirectory(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "empty-evidence-roundtrip")
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
	}); err != nil {
		t.Fatal(err)
	}
	// Git and Bazel file packaging do not preserve empty directories.
	evidenceDir := filepath.Join(goalDir, "attempts", "attempt-1", "evidence")
	if err := os.Remove(evidenceDir); err != nil {
		t.Fatal(err)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("missing empty evidence directory invalidated the record: %v", err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "2", AttemptID: "attempt-1",
		NextAction: "Continue verification after checkout.",
	}); err != nil {
		t.Fatalf("missing empty evidence directory prevented continuation: %v", err)
	}
	if err := store.ValidateGoal(goalDir); err != nil {
		t.Fatalf("continued attempt is invalid: %v", err)
	}
}

func TestAttemptStillRequiresPlanAndResult(t *testing.T) {
	for _, name := range []string{"plan.md", "result.md"} {
		t.Run(name, func(t *testing.T) {
			store, root := newTestStore(t)
			goalDir := initTestGoal(t, store, root, "missing-required-artifact")
			if _, err := store.Checkpoint(CheckpointOptions{
				GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
			}); err != nil {
				t.Fatal(err)
			}
			attemptDir := filepath.Join(goalDir, "attempts", "attempt-1")
			if err := os.Remove(filepath.Join(attemptDir, name)); err != nil {
				t.Fatal(err)
			}
			if err := store.ValidateGoal(goalDir); err == nil {
				t.Fatalf("attempt with missing %s was accepted", name)
			}
		})
	}
}

func TestAttemptStillRequiresDeclaredEvidence(t *testing.T) {
	for _, removeDirectory := range []bool{false, true} {
		name := "missing file"
		if removeDirectory {
			name = "missing directory"
		}
		t.Run(name, func(t *testing.T) {
			store, root := newTestStore(t)
			goalDir := initTestGoal(t, store, root, "missing-declared-evidence")
			evidence := filepath.Join(root, "out", "task", "evidence.md")
			writeTestFile(t, evidence, "# Evidence\n\nCandidate verification.\n")
			if _, err := store.Checkpoint(CheckpointOptions{
				GoalDir: goalDir, ExpectedResourceVersion: "1", AttemptID: "attempt-1",
				EvidenceFiles: []string{evidence},
			}); err != nil {
				t.Fatal(err)
			}
			evidenceDir := filepath.Join(goalDir, "attempts", "attempt-1", "evidence")
			if err := os.Remove(filepath.Join(evidenceDir, "evidence.md")); err != nil {
				t.Fatal(err)
			}
			if removeDirectory {
				if err := os.Remove(evidenceDir); err != nil {
					t.Fatal(err)
				}
			}
			if err := store.ValidateGoal(goalDir); err == nil {
				t.Fatal("attempt with missing declared evidence was accepted")
			}
		})
	}
}
