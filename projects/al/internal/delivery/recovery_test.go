package delivery

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// partialPushFixture creates an actual pushed revision with an interrupted PR
// update and a provider object that still reports the original commit.
func partialPushFixture(t *testing.T) (integrationDeliveryFixture, *prepareReport) {
	t.Helper()
	fixture := newIntegrationDeliveryFixture(t)
	first := fixture.prepare(t)
	if _, err := fixture.delivery.publish(context.Background(), publishOptions{ReceiptFile: "out/delivery/prepare.json"}); err != nil {
		t.Fatalf("initial publication: %v", err)
	}
	writeTestFile(t, filepath.Join(fixture.work, "out/delivery/commit.md"), "Recover interrupted publication\n\nPreserve the exact pushed candidate.\n")
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "revision\n")
	prepared, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md", ReceiptFile: "out/delivery/prepare.json", Paths: []string{"feature.txt"}, RewriteOID: first.HeadOID,
	})
	if err != nil {
		t.Fatalf("prepare interrupted revision: %v", err)
	}
	fixture.forge.headOIDOverride = first.HeadOID
	fixture.forge.failUpdates = 1
	report, err := fixture.delivery.publish(context.Background(), publishOptions{ReceiptFile: "out/delivery/prepare.json"})
	if err == nil || report == nil || report.Status != "publication_verification_pending" || report.HeadOID != prepared.HeadOID || report.RemoteHeadOID != prepared.HeadOID || report.PullRequest == nil {
		t.Fatalf("partial push lost its identity: %#v, %v", report, err)
	}
	if runTestGit(t, fixture.remote, "rev-parse", "feature") != prepared.HeadOID {
		t.Fatal("fixture did not complete the candidate push")
	}
	fixture.forge.requireObservedUpdate = true
	fixture.forge.syncHeadOnUpdate = true
	return fixture, prepared
}

// recoveryLocalState captures independently observable local Git state.
func recoveryLocalState(t *testing.T, work string) map[string]string {
	t.Helper()
	return map[string]string{
		"head":         runTestGit(t, work, "rev-parse", "HEAD"),
		"branch":       runTestGit(t, work, "branch", "--show-current"),
		"status":       runTestGit(t, work, "status", "--porcelain"),
		"index":        runTestGit(t, work, "ls-files", "--stage"),
		"working_diff": runTestGit(t, work, "diff", "--binary"),
		"staged_diff":  runTestGit(t, work, "diff", "--cached", "--binary"),
	}
}

// TestRecoverPullRequestPreservesAdvancedLocalWork covers the real partial-push
// failure, an advanced base/local commit, staged and unstaged files, and retry.
func TestRecoverPullRequestPreservesAdvancedLocalWork(t *testing.T) {
	fixture, prepared := partialPushFixture(t)
	fixture.advanceBase(t)
	writeTestFile(t, filepath.Join(fixture.work, "local.txt"), "new local progress\n")
	runTestGit(t, fixture.work, "add", "local.txt")
	runTestGit(t, fixture.work, "commit", "-m", "Continue local work")
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "staged progress\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "unstaged progress\n")
	writeTestFile(t, filepath.Join(fixture.work, "untracked.txt"), "untracked progress\n")
	before := recoveryLocalState(t, fixture.work)
	receiptPath := filepath.Join(fixture.work, "out/delivery/prepare.json")
	receiptBefore, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatalf("read original receipt: %v", err)
	}
	callsBefore := fixture.forge.updateCalls
	report, err := fixture.delivery.recoverPullRequest(context.Background(), recoverPullRequestOptions{receiptFile: "out/delivery/prepare.json"})
	if err != nil || report == nil || report.Status != "pull_request_synchronized" || report.HeadOID != prepared.HeadOID || report.RemoteHeadOID != prepared.HeadOID || report.Pushed || report.PullRequest == nil || report.PullRequest.HeadRefOID != prepared.HeadOID {
		t.Fatalf("recover PR = %#v, %v", report, err)
	}
	if !reflect.DeepEqual(before, recoveryLocalState(t, fixture.work)) {
		t.Fatal("recovery changed later local work")
	}
	if got, err := os.ReadFile(filepath.Join(fixture.work, "untracked.txt")); err != nil || string(got) != "untracked progress\n" {
		t.Fatalf("untracked work changed: %q, %v", got, err)
	}
	receiptAfter, err := os.ReadFile(receiptPath)
	if err != nil || string(receiptAfter) != string(receiptBefore) {
		t.Fatalf("recovery changed its original receipt: %v", err)
	}
	if runTestGit(t, fixture.remote, "rev-parse", "feature") != prepared.HeadOID {
		t.Fatal("recovery changed the remote feature ref")
	}
	second, err := fixture.delivery.recoverPullRequest(context.Background(), recoverPullRequestOptions{receiptFile: "out/delivery/prepare.json"})
	if err != nil || second.Status != "pull_request_synchronized" || fixture.forge.updateCalls != callsBefore+1 {
		t.Fatalf("non-idempotent recovery: %#v, %v", second, err)
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		artifact, err := json.MarshalIndent(map[string]any{"prepared_head": prepared.HeadOID, "local_head": before["head"], "local_state_preserved": true, "receipt_preserved": true, "remote_ref_preserved": true, "mutation_count": fixture.forge.updateCalls - callsBefore, "result": report}, "", "  ")
		if err != nil {
			t.Fatalf("encode recovery evidence: %v", err)
		}
		if err := os.WriteFile(filepath.Join(directory, "partial-push-pr-recovery.json"), append(artifact, '\n'), 0o600); err != nil {
			t.Fatalf("write recovery evidence: %v", err)
		}
	}
}

// TestRecoverPullRequestRefusesChangedBindings verifies that explicit recovery
// cannot authorize a different PR, human edits, another head, or a different receipt.
func TestRecoverPullRequestRefusesChangedBindings(t *testing.T) {
	cases := map[string]func(*testing.T, integrationDeliveryFixture, *prepareReport){
		"human_title":   func(_ *testing.T, f integrationDeliveryFixture, _ *prepareReport) { f.forge.pull.Title = "Human title" },
		"human_body":    func(_ *testing.T, f integrationDeliveryFixture, _ *prepareReport) { f.forge.pull.Body = "Human body" },
		"draft_changed": func(_ *testing.T, f integrationDeliveryFixture, _ *prepareReport) { f.forge.pull.IsDraft = true },
		"author_changed": func(_ *testing.T, f integrationDeliveryFixture, _ *prepareReport) {
			f.forge.pull.AuthorLogin = "another-author"
		},
		"replacement_pr": func(_ *testing.T, f integrationDeliveryFixture, _ *prepareReport) { f.forge.pull.ID = "replacement-pr" },
		"unrelated_pr_head": func(_ *testing.T, f integrationDeliveryFixture, _ *prepareReport) {
			f.forge.headOIDOverride = testOID('e')
		},
		"remote_advanced": func(t *testing.T, f integrationDeliveryFixture, _ *prepareReport) {
			writeTestFile(t, filepath.Join(f.work, "follow-up.txt"), "another candidate\n")
			runTestGit(t, f.work, "add", "follow-up.txt")
			runTestGit(t, f.work, "commit", "-m", "Another candidate")
			runTestGit(t, f.work, "push", "origin", "HEAD:feature")
		},
		"remote_not_pushed": func(t *testing.T, f integrationDeliveryFixture, p *prepareReport) {
			runTestGit(t, f.work, "push", "origin", p.Receipt.ExpectedRemoteHead.OID+":refs/heads/feature", "--force-with-lease=refs/heads/feature:"+p.HeadOID)
		},
		"local_branch_changed": func(t *testing.T, f integrationDeliveryFixture, _ *prepareReport) {
			runTestGit(t, f.work, "switch", "-c", "other")
		},
		"receipt_tree_changed": func(t *testing.T, f integrationDeliveryFixture, p *prepareReport) {
			receipt := p.Receipt
			receipt.PreparedTreeOID = testOID('e')
			encoded, err := json.Marshal(receipt)
			if err != nil {
				t.Fatalf("encode changed receipt: %v", err)
			}
			writeTestFile(t, filepath.Join(f.work, "out/delivery/prepare.json"), string(encoded))
		},
		"repository_binding_changed": func(_ *testing.T, f integrationDeliveryFixture, _ *prepareReport) {
			f.delivery.remoteRepository.Name = "another-repository"
		},
		"remote_advanced_during_recovery": func(t *testing.T, f integrationDeliveryFixture, _ *prepareReport) {
			f.forge.pullRequestsCalls = 0
			f.forge.onPullRequestsCall = func(call int) {
				if call != 1 {
					return
				}
				runTestGit(t, f.seed, "fetch", "origin", "feature")
				runTestGit(t, f.seed, "switch", "-c", "concurrent", "FETCH_HEAD")
				writeTestFile(t, filepath.Join(f.seed, "concurrent.txt"), "concurrent progress\n")
				runTestGit(t, f.seed, "add", "concurrent.txt")
				runTestGit(t, f.seed, "commit", "-m", "Concurrent progress")
				runTestGit(t, f.seed, "push", "origin", "HEAD:feature")
			}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			fixture, prepared := partialPushFixture(t)
			mutate(t, fixture, prepared)
			before := recoveryLocalState(t, fixture.work)
			calls := fixture.forge.updateCalls
			if report, err := fixture.delivery.recoverPullRequest(context.Background(), recoverPullRequestOptions{receiptFile: "out/delivery/prepare.json"}); err == nil {
				t.Fatalf("accepted changed binding: %#v", report)
			}
			if fixture.forge.updateCalls != calls {
				t.Fatal("refused recovery attempted a metadata mutation")
			}
			if !reflect.DeepEqual(before, recoveryLocalState(t, fixture.work)) {
				t.Fatal("refused recovery changed local source")
			}
		})
	}
}

// TestRecoverPullRequestRetainsUnknownMutationEvidence preserves known identities
// when the provider mutates metadata but cannot verify its response.
func TestRecoverPullRequestRetainsUnknownMutationEvidence(t *testing.T) {
	fixture, prepared := partialPushFixture(t)
	fixture.forge.failAfterUpdate = true
	report, err := fixture.delivery.recoverPullRequest(context.Background(), recoverPullRequestOptions{receiptFile: "out/delivery/prepare.json"})
	if err == nil || report == nil || report.Status != "pull_request_verification_pending" || report.HeadOID != prepared.HeadOID || report.RemoteHeadOID != prepared.HeadOID || report.PullRequest == nil || report.PullRequest.URL == "" || report.Pushed {
		t.Fatalf("lost unknown recovery evidence: %#v, %v", report, err)
	}
	calls := fixture.forge.updateCalls
	fixture.forge.failAfterUpdate = false
	if _, err := fixture.delivery.recoverPullRequest(context.Background(), recoverPullRequestOptions{receiptFile: "out/delivery/prepare.json"}); err != nil {
		t.Fatalf("inspect completed mutation: %v", err)
	}
	if fixture.forge.updateCalls != calls {
		t.Fatal("retry repeated an already completed mutation")
	}
}
