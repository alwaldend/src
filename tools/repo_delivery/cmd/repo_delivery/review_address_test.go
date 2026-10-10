package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// This workflow uses real Git preparation, validation records and publication,
// with a deterministic forge to inject failures without posting live comments.
type addressTestForge struct {
	*integrationForge
	reviews    *reviewReceiptForge
	replyCalls int
	replyErr   error
}

func (f *addressTestForge) PullRequests(ctx context.Context, repository remoteRepository, head string) ([]pullRequest, error) {
	values, err := f.integrationForge.PullRequests(ctx, repository, head)
	if err == nil && len(values) == 1 {
		values[0].UpdatedAt = f.reviews.pullRequest.UpdatedAt
	}
	return values, err
}

func (f *addressTestForge) InspectReviews(ctx context.Context, r remoteRepository, p pullRequest) (*reviewInspection, error) {
	return f.reviews.InspectReviews(ctx, r, p)
}

func (f *addressTestForge) ReplyToReviewThread(ctx context.Context, r remoteRepository, p pullRequest, e reviewThreadExpectation, b string) (*reviewInspection, error) {
	f.replyCalls++
	if f.replyErr != nil {
		return nil, f.replyErr
	}
	return f.reviews.ReplyToReviewThread(ctx, r, p, e, b)
}

func (f *addressTestForge) ResolveReviewThread(ctx context.Context, r remoteRepository, p pullRequest, e reviewThreadExpectation) (*reviewInspection, error) {
	return f.reviews.ResolveReviewThread(ctx, r, p, e)
}

func newAddressFixture(t *testing.T) (integrationDeliveryFixture, *addressTestForge, reviewAddressOptions) {
	t.Helper()
	fixture, _, _, _ := newValidationFixture(t)
	requireValidated(t, fixture)
	if _, err := fixture.delivery.continueCandidate(context.Background(), testValidationReceipt, true, false); err != nil {
		t.Fatal(err)
	}
	pull, err := fixture.forge.refreshed(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	pull.UpdatedAt = "2026-08-30T00:00:00Z"
	reviews := &reviewReceiptForge{repository: fixture.delivery.remoteRepository, pullRequest: *pull, thread: reviewThread{
		ID: "RT_address", Path: "file.go", Comments: []reviewComment{{ID: "RC_feedback", Body: "Please fix this.", Path: "file.go", CommitOID: pull.HeadRefOID}},
	}}
	forge := &addressTestForge{integrationForge: fixture.forge, reviews: reviews}
	fixture.delivery.forge = forge
	options := reviewAddressOptions{ReceiptFile: testValidationReceipt, ThreadID: reviews.thread.ID, BodyFile: "out/delivery/reply.md"}
	writeTestFile(t, filepath.Join(fixture.work, options.BodyFile), "Fixed the requested behavior and validated the published candidate.")
	return fixture, forge, options
}

func TestReviewAddressPublishedCandidateAndRecovery(t *testing.T) {
	fixture, forge, options := newAddressFixture(t)
	forge.reviews.beforeResolve = func() error { return &reviewResolutionReadError{err: errors.New("inventory epoch advanced")} }
	workflow := reviewAddressWorkflow{delivery: fixture.delivery, options: options}
	if _, err := workflow.run(context.Background()); err == nil || !strings.Contains(err.Error(), "was restored") {
		t.Fatalf("first attempt = %v", err)
	}
	if forge.replyCalls != 1 {
		t.Fatalf("reply calls = %d", forge.replyCalls)
	}
	forge.reviews.beforeResolve = nil
	report, err := workflow.run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if forge.replyCalls != 1 || !report.Threads[0].IsResolved {
		t.Fatalf("retry duplicated reply or failed resolution: %d %#v", forge.replyCalls, report)
	}
	if _, err := workflow.run(context.Background()); err == nil {
		t.Fatal("completed workflow must not repost")
	}
	if forge.replyCalls != 1 {
		t.Fatal("completed workflow posted again")
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		file, err := os.Create(filepath.Join(directory, "review-address.json"))
		if err != nil {
			t.Fatal(err)
		}
		if err := writeJSON(file, report); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestReviewAddressRefusesDuplicateAfterUnknownReplyOutcome(t *testing.T) {
	fixture, forge, options := newAddressFixture(t)
	forge.replyErr = errors.New("reply response lost; outcome unknown")
	workflow := reviewAddressWorkflow{delivery: fixture.delivery, options: options}
	if _, err := workflow.run(context.Background()); err == nil {
		t.Fatal("unknown reply accepted")
	}
	forge.replyErr = nil
	if _, err := workflow.run(context.Background()); err == nil {
		t.Fatal("unknown attempt replayed")
	}
	if forge.replyCalls != 1 {
		t.Fatalf("unknown reply retried %d times", forge.replyCalls)
	}
}

func TestReviewAddressRefusesChangedInputs(t *testing.T) {
	for _, change := range []string{"dirty candidate", "changed body", "missing body", "changed thread", "expired receipt", "altered validation log", "unpublished remote", "missing validation", "changed head"} {
		t.Run(change, func(t *testing.T) {
			fixture, forge, options := newAddressFixture(t)
			workflow := reviewAddressWorkflow{delivery: fixture.delivery, options: options}
			forge.reviews.beforeResolve = func() error { return &reviewResolutionReadError{err: errors.New("read failed")} }
			if change != "dirty candidate" && change != "altered validation log" && change != "unpublished remote" && change != "missing validation" && change != "changed head" {
				if _, err := workflow.run(context.Background()); err == nil {
					t.Fatal("expected read failure")
				}
			}
			switch change {
			case "dirty candidate":
				writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "changed")
			case "changed body":
				writeTestFile(t, filepath.Join(fixture.work, options.BodyFile), "A different reply.")
			case "missing body":
				if err := os.Remove(filepath.Join(fixture.work, options.BodyFile)); err != nil {
					t.Fatal(err)
				}
			case "changed thread":
				forge.reviews.thread.Comments = append(forge.reviews.thread.Comments, reviewComment{ID: "RC_late", Body: "Wait."})
			case "expired receipt":
				// Advance the authority clock beyond the fixed five-minute receipt lifetime.
				ctx := withReviewReplyReceiptClock(context.Background(), func() time.Time { return time.Now().Add(10 * time.Minute) })
				if _, err := workflow.run(ctx); err == nil {
					t.Fatal("expired receipt accepted")
				}
				if forge.replyCalls != 1 {
					t.Fatal("expired authority caused another reply")
				}
				return
			case "unpublished remote":
				runTestGit(t, fixture.seed, "push", "origin", ":feature")
			case "missing validation":
				if err := os.Remove(filepath.Join(fixture.work, validationStatePath(testValidationReceipt))); err != nil {
					t.Fatal(err)
				}
			case "changed head":
				writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "new candidate")
				runTestGit(t, fixture.work, "add", "feature.txt")
				runTestGit(t, fixture.work, "commit", "-m", "Changed candidate")
			case "altered validation log":
				writeTestFile(t, filepath.Join(fixture.work, validationStatePath(testValidationReceipt)+".check-01.log"), "altered")
			}
			forge.reviews.beforeResolve = nil
			before := forge.replyCalls
			if _, err := workflow.run(context.Background()); err == nil {
				t.Fatal("changed input accepted")
			}
			if forge.replyCalls != before || forge.reviews.sawResolve {
				t.Fatal("changed input caused a mutation")
			}
		})
	}
}
