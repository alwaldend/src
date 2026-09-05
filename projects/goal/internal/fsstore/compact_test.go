package fsstore

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSummaryCheckpointResumesLatestProgressAndRejectsStaleWriter(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "compact-goal")
	evidence := filepath.Join(root, "out", "task", "first-check.md")
	writeTestFile(t, evidence, "# Evidence\n\nCandidate one passed its focused check.\n")
	first, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1",
		Summary: "The focused check passed; integration remains unverified.",
		Subject: "candidate-one", NextAction: "Run integration checks.",
		EvidenceFiles: []string{evidence},
	})
	if err != nil {
		t.Fatal(err)
	}
	view, err := store.ShowGoal(goalDir, 1)
	if err != nil {
		t.Fatal(err)
	}
	active := view.ActiveAttempt
	if active == nil || active.Subject != "candidate-one" ||
		active.NextAction != "Run integration checks." ||
		!strings.Contains(active.ResultMarkdown, "integration remains unverified") ||
		active.ResultTruncated || len(active.Evidence) != 1 {
		t.Fatalf("incomplete checkpoint view: %+v", active)
	}
	attemptID := active.AttemptID
	planPath := filepath.Join(goalDir, "attempts", attemptID, "plan.md")
	plan, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(plan), "Goal compact-goal") ||
		!strings.Contains(string(plan), "The result is verified.") {
		t.Fatalf("generated plan lacks objective/acceptance: %s", plan)
	}
	update := CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: first.ResourceVersion,
		Summary: "Integration revealed a defect; the candidate has been corrected.",
		Subject: "candidate-two", NextAction: "Rerun affected integration checks.",
	}
	second, err := store.Checkpoint(update)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Checkpoint(update); err == nil || !strings.Contains(err.Error(), "stale resourceVersion") {
		t.Fatalf("stale checkpoint error = %v", err)
	}
	fresh, err := NewStore(root)
	if err != nil {
		t.Fatal(err)
	}
	view, err = fresh.ShowGoal(goalDir, 1)
	if err != nil {
		t.Fatal(err)
	}
	active = view.ActiveAttempt
	if view.Goal.Metadata.ResourceVersion != second.ResourceVersion || view.Total != 1 ||
		active == nil || active.AttemptID != attemptID || active.Subject != "candidate-two" ||
		active.NextAction != "Rerun affected integration checks." ||
		!strings.Contains(active.ResultMarkdown, "candidate has been corrected") ||
		active.ResultDigest == "" || active.SourcePath == "" ||
		active.ObservedAt == "" || len(active.Evidence) != 1 ||
		view.Goal.Status.Outcome != "open" {
		t.Fatalf("fresh resume lost latest progress or accepted it: %+v", view)
	}
	_, _, attempts, err := store.loadAndValidate(goalDir)
	if err != nil || attempts[0].Metadata.Generation != 2 {
		t.Fatalf("continuation spec update did not advance generation: %+v, %v", attempts, err)
	}
	updatedPlan, err := os.ReadFile(planPath)
	if err != nil || string(updatedPlan) != string(plan) {
		t.Fatalf("initial plan changed: %s, %v", updatedPlan, err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: second.ResourceVersion,
		Subject: "candidate-three",
	}); err == nil || !strings.Contains(err.Error(), "prior evidence") {
		t.Fatalf("candidate changed while retaining old result: %v", err)
	}
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: second.ResourceVersion,
		NextAction: "Inspect the integration fixture first.",
	}); err != nil {
		t.Fatal(err)
	}
	view, err = fresh.ShowGoal(goalDir, 1)
	if err != nil || view.ActiveAttempt.NextAction != "Inspect the integration fixture first." {
		t.Fatalf("next-action-only checkpoint ignored: %+v, %v", view, err)
	}
}

func TestSummaryCheckpointBoundsAndConflictingInputs(t *testing.T) {
	for _, test := range []struct {
		name string
		edit func(*CheckpointOptions)
	}{
		{"blank", func(o *CheckpointOptions) { o.Summary = " \n" }},
		{"too long", func(o *CheckpointOptions) { o.Summary = strings.Repeat("x", maxCheckpointSummaryBytes+1) }},
		{"invalid UTF8", func(o *CheckpointOptions) { o.Summary = string([]byte{255}) }},
		{"NUL", func(o *CheckpointOptions) { o.Summary = "progress\x00" }},
		{"missing subject", func(o *CheckpointOptions) { o.Subject = "" }},
		{"missing next action", func(o *CheckpointOptions) { o.NextAction = "" }},
		{"result conflict", func(o *CheckpointOptions) { o.ResultFile = "result.md" }},
		{"plan conflict", func(o *CheckpointOptions) { o.PlanFile = "plan.md" }},
		{"criteria conflict", func(o *CheckpointOptions) { o.CriteriaFile = "criteria.yaml" }},
		{"plan only conflict", func(o *CheckpointOptions) { o.PlanOnly = true }},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, root := newTestStore(t)
			goalDir := initTestGoal(t, store, root, "invalid-summary")
			options := CheckpointOptions{
				GoalDir: goalDir, ExpectedResourceVersion: "1",
				Summary: "Current progress.", Subject: "candidate", NextAction: "Run checks.",
			}
			test.edit(&options)
			if _, err := store.Checkpoint(options); err == nil {
				t.Fatal("invalid summary checkpoint succeeded")
			}
			view, err := store.ShowGoal(goalDir, 1)
			if err != nil || view.Goal.Metadata.ResourceVersion != "1" || view.Total != 0 {
				t.Fatalf("invalid checkpoint changed state: %+v, %v", view, err)
			}
		})
	}
}

func TestShowBoundsLargeActiveResultAndPreservesUTF8(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "bounded-summary")
	result := filepath.Join(root, "out", "task", "large-result.md")
	content := strings.Repeat("x", maxCheckpointSummaryBytes-1) + "界 tail"
	writeTestFile(t, result, content)
	if _, err := store.Checkpoint(CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", ResultFile: result,
	}); err != nil {
		t.Fatal(err)
	}
	view, err := store.ShowGoal(goalDir, 1)
	if err != nil {
		t.Fatal(err)
	}
	active := view.ActiveAttempt
	if active == nil || !active.ResultTruncated || !utf8.ValidString(active.ResultMarkdown) ||
		len(active.ResultMarkdown) >= maxCheckpointSummaryBytes || active.ResultBytes != len(content) ||
		active.ResultDigest != digestBytes([]byte(content)) {
		t.Fatalf("invalid bounded result view: %+v", active)
	}
}

func TestSummaryRequiresExplicitCloseReviewAndRetainsEvidence(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "summary-close")
	evidence := filepath.Join(root, "out", "task", "acceptance.md")
	writeTestFile(t, evidence, "# Evidence\n\nCandidate final passed the acceptance check.\n")
	options := CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1",
		Summary: "All checks passed for candidate final.", Subject: "candidate-final",
		NextAction: "Deliver the reviewed candidate.", EvidenceFiles: []string{evidence},
		CloseAttempt: true, Outcome: "achieved",
	}
	if _, err := store.Checkpoint(options); err == nil || !strings.Contains(err.Error(), "--review-file") {
		t.Fatalf("summary bypassed close review: %v", err)
	}
	review := filepath.Join(root, "out", "task", "review.yaml")
	writeTestFile(t, review, `decision: accept
criteria:
  - criterionID: criterion-001
    criterionRevision: 1
    verdict: pass
    evidenceRefs:
      - evidence/acceptance.md
`)
	options.ReviewFile = review
	if _, err := store.Checkpoint(options); err != nil {
		t.Fatal(err)
	}
	goal, _, attempts, err := store.loadAndValidate(goalDir)
	if err != nil || goal.Status.Outcome != "achieved" || len(attempts) != 1 ||
		len(attempts[0].Status.Artifacts.Evidence) != 1 ||
		attempts[0].Status.Review.Criteria[0].EvidenceRefs[0] != "evidence/acceptance.md" {
		t.Fatalf("accepted summary lost guarded review: %+v, %+v, %v", goal, attempts, err)
	}
}

func TestSummaryContinuationRecoversAfterResultPublicationFailure(t *testing.T) {
	store, root := newTestStore(t)
	goalDir := initTestGoal(t, store, root, "summary-recovery")
	options := CheckpointOptions{
		GoalDir: goalDir, ExpectedResourceVersion: "1", Summary: "First checkpoint.",
		Subject: "candidate-one", NextAction: "Run checks.",
	}
	if _, err := store.Checkpoint(options); err != nil {
		t.Fatal(err)
	}
	store.beforeRename = func(path string) error {
		if filepath.Base(path) == "result.md" {
			return errors.New("injected result publication failure")
		}
		return nil
	}
	options.ExpectedResourceVersion = "2"
	options.Summary = "Second checkpoint after candidate correction."
	options.Subject = "candidate-two"
	options.NextAction = "Rerun checks."
	ref, err := store.Checkpoint(options)
	if err == nil || ref.ResourceVersion != "3" {
		t.Fatalf("missing committed publication failure: %+v, %v", ref, err)
	}
	store.beforeRename = nil
	if _, err := store.ShowGoal(goalDir, 1); err == nil {
		t.Fatal("show accepted an incomplete checkpoint")
	}
	if _, err := store.Recover(goalDir); err != nil {
		t.Fatal(err)
	}
	view, err := store.ShowGoal(goalDir, 1)
	if err != nil || view.ActiveAttempt == nil || view.ActiveAttempt.Subject != "candidate-two" ||
		!strings.Contains(view.ActiveAttempt.ResultMarkdown, "Second checkpoint") {
		t.Fatalf("recovered checkpoint differs: %+v, %v", view, err)
	}
}
