package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type integrationForge struct {
	// reviewError simulates an unavailable or concurrently changing review inventory.
	reviewError                      error
	name                             string
	remote                           string
	repository                       remoteRepository
	pull                             *pullRequest
	baseRefOIDOverride               string
	pullRequestsCalls                int
	mutateMetadataOnPullRequestsCall int
	failUpdates                      int
	onPullRequestsCall               func(int)
	expectedHeadBranch               string
}

func (f *integrationForge) Name() string {
	if f.name != "" {
		return f.name
	}
	return "integration"
}

func (f *integrationForge) Preflight(
	context.Context,
	remoteRepository,
) error {
	return nil
}

func (f *integrationForge) InspectReviews(
	context.Context,
	remoteRepository,
	pullRequest,
) (*reviewInspection, error) {
	if f.reviewError != nil {
		return nil, f.reviewError
	}
	return &reviewInspection{}, nil
}

func (f *integrationForge) ReplyToReviewThread(
	context.Context,
	remoteRepository,
	pullRequest,
	reviewThreadExpectation,
	string,
) (*reviewInspection, error) {
	return nil, errors.New("unexpected review-thread reply")
}

func (f *integrationForge) CommentOnPullRequest(
	context.Context,
	remoteRepository,
	pullRequest,
	topLevelCommentExpectation,
	string,
) (*reviewInspection, error) {
	return nil, errors.New("unexpected top-level pull-request comment")
}

func (f *integrationForge) ResolveReviewThread(
	context.Context,
	remoteRepository,
	pullRequest,
	reviewThreadExpectation,
) (*reviewInspection, error) {
	return nil, errors.New("unexpected review-thread resolution")
}

func (f *integrationForge) RequestReview(
	context.Context,
	remoteRepository,
	pullRequest,
	[]string,
) (*reviewInspection, error) {
	return nil, errors.New("unexpected review request")
}

func (f *integrationForge) refOID(
	ctx context.Context,
	branch string,
) (string, error) {
	command := exec.CommandContext(
		ctx,
		"git",
		"--git-dir",
		f.remote,
		"rev-parse",
		"refs/heads/"+branch,
	)
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (f *integrationForge) refreshed(
	ctx context.Context,
) (*pullRequest, error) {
	if f.pull == nil {
		return nil, nil
	}
	result := *f.pull
	var err error
	if f.baseRefOIDOverride == "" {
		result.BaseRefOID, err = f.refOID(ctx, result.BaseRefName)
		if err != nil {
			return nil, err
		}
	} else {
		result.BaseRefOID = f.baseRefOIDOverride
	}
	result.HeadRefOID, err = f.refOID(ctx, result.HeadRefName)
	if err != nil {
		return nil, err
	}
	f.pull = &result
	return &result, nil
}

func (f *integrationForge) PullRequests(
	ctx context.Context,
	repository remoteRepository,
	headBranch string,
) ([]pullRequest, error) {
	if !sameRemoteRepository(repository, f.repository) {
		return nil, fmt.Errorf("unexpected repository")
	}
	if f.expectedHeadBranch != "" && headBranch != f.expectedHeadBranch {
		return nil, fmt.Errorf(
			"unexpected head branch %q, want %q",
			headBranch,
			f.expectedHeadBranch,
		)
	}
	pull, err := f.refreshed(ctx)
	if err != nil {
		return nil, err
	}
	f.pullRequestsCalls++
	if f.onPullRequestsCall != nil {
		f.onPullRequestsCall(f.pullRequestsCalls)
	}
	if pull == nil {
		return nil, nil
	}
	if f.pullRequestsCalls == f.mutateMetadataOnPullRequestsCall {
		pull.Title = "Concurrent human title edit"
		f.pull = pull
	}
	return []pullRequest{*pull}, nil
}

func (f *integrationForge) CreatePullRequest(
	ctx context.Context,
	repository remoteRepository,
	input pullRequestInput,
) (*pullRequest, error) {
	if f.pull != nil || !sameRemoteRepository(repository, f.repository) {
		return nil, fmt.Errorf("unexpected pull-request creation")
	}
	remoteHead, err := f.refOID(ctx, input.HeadRefName)
	if err != nil || input.BaseRefName == "" || input.HeadRefName == "" ||
		!isObjectID(input.ExpectedHeadOID) || remoteHead != input.ExpectedHeadOID {
		return nil, fmt.Errorf("pull-request creation did not bind its branch and expected head")
	}
	f.pull = &pullRequest{
		ID:                  "integration-pr",
		Number:              1,
		URL:                 "https://github.com/owner/repo/pull/1",
		State:               "OPEN",
		Title:               input.Title,
		Body:                input.Body,
		BaseRefName:         input.BaseRefName,
		HeadRefName:         input.HeadRefName,
		HeadRepositoryOwner: f.repository.Owner,
		HeadRepositoryName:  f.repository.Name,
	}
	return f.refreshed(ctx)
}

func (f *integrationForge) UpdatePullRequest(
	ctx context.Context,
	repository remoteRepository,
	pull pullRequest,
	input pullRequestInput,
) (*pullRequest, error) {
	if f.pull == nil || pull.ID != f.pull.ID ||
		!sameRemoteRepository(repository, f.repository) {
		return nil, fmt.Errorf("unexpected pull-request update")
	}
	if f.failUpdates > 0 {
		f.failUpdates--
		return nil, fmt.Errorf("injected pull-request update failure")
	}
	remoteHead, err := f.refOID(ctx, input.HeadRefName)
	if err != nil || input.BaseRefName == "" || input.HeadRefName == "" ||
		!isObjectID(input.ExpectedHeadOID) || remoteHead != input.ExpectedHeadOID {
		return nil, fmt.Errorf("pull-request update did not bind its branch and expected head")
	}
	f.pull.Title = input.Title
	f.pull.Body = input.Body
	f.pull.BaseRefName = input.BaseRefName
	f.pull.HeadRefName = input.HeadRefName
	return f.refreshed(ctx)
}

type integrationDeliveryFixture struct {
	delivery *delivery
	forge    *integrationForge
	remote   string
	seed     string
	work     string
}

type integrationHookRunner struct {
	delegate commandRunner
	hook     func(command)
	fired    bool
}

type integrationArgumentHookRunner struct {
	delegate commandRunner
	argument string
	hook     func(command)
	fired    bool
}

func (r *integrationArgumentHookRunner) Run(
	ctx context.Context,
	invocation command,
) (commandResult, error) {
	if !r.fired {
		for _, argument := range invocation.Args {
			if argument == r.argument {
				r.fired = true
				r.hook(invocation)
				break
			}
		}
	}
	return r.delegate.Run(ctx, invocation)
}

func (r *integrationHookRunner) Run(
	ctx context.Context,
	invocation command,
) (commandResult, error) {
	if !r.fired {
		for _, argument := range invocation.Args {
			if argument == "commit-tree" {
				r.fired = true
				r.hook(invocation)
				break
			}
		}
	}
	return r.delegate.Run(ctx, invocation)
}

func newIntegrationDeliveryFixture(
	t *testing.T,
) integrationDeliveryFixture {
	t.Helper()
	root := t.TempDir()
	remote := filepath.Join(root, "remote.git")
	seed := filepath.Join(root, "seed")
	work := filepath.Join(root, "work")

	runTestGit(t, root, "init", "--bare", remote)
	runTestGit(t, root, "init", "--initial-branch=master", seed)
	configureTestRepository(t, seed)
	writeTestFile(t, filepath.Join(seed, ".gitignore"), "/out/\n")
	writeTestFile(t, filepath.Join(seed, "base.txt"), "base\n")
	runTestGit(t, seed, "add", ".gitignore", "base.txt")
	runTestGit(t, seed, "commit", "-m", "base")
	runTestGit(t, seed, "remote", "add", "origin", remote)
	runTestGit(t, seed, "push", "origin", "master")

	runTestGit(t, root, "clone", remote, work)
	configureTestRepository(t, work)
	runTestGit(t, work, "switch", "-c", "feature")
	repository, err := openGitRepository(
		context.Background(),
		work,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("openGitRepository() error = %v", err)
	}
	identity := remoteRepository{
		Host: "github.com", Owner: "owner", Name: "repo",
	}
	forge := &integrationForge{
		remote: remote, repository: identity, expectedHeadBranch: "feature",
	}
	return integrationDeliveryFixture{
		delivery: &delivery{
			repository:       repository,
			remote:           "origin",
			fetchEndpoint:    remote,
			pushEndpoint:     remote,
			base:             "master",
			remoteRepository: identity,
			forge:            forge,
		},
		forge: forge, remote: remote, seed: seed, work: work,
	}
}

func (f integrationDeliveryFixture) prepare(
	t *testing.T,
) *prepareReport {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(f.work, "out", "delivery"), 0o700); err != nil {
		t.Fatalf("create message directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(f.work, "out", "delivery", "commit.md"),
		"Add delivery fixture\n\nExercise the complete delivery flow.\n",
	)
	writeTestFile(t, filepath.Join(f.work, "feature.txt"), "feature\n")
	report, err := f.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	return report
}

func (f integrationDeliveryFixture) advanceBase(t *testing.T) {
	t.Helper()
	writeTestFile(t, filepath.Join(f.seed, "base-advance.txt"), "advance\n")
	runTestGit(t, f.seed, "add", "base-advance.txt")
	runTestGit(t, f.seed, "commit", "-m", "advance base")
	runTestGit(t, f.seed, "push", "origin", "master")
}

func TestDeliveryPublishesLargePreparationReceipt(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	for i := range 6000 {
		path := filepath.Join(fixture.work, "submissions", strings.Repeat("a", 180), fmt.Sprintf("%04d", i), "submission.toml")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create submission directory: %v", err)
		}
		writeTestFile(t, path, fmt.Sprintf("id = %d\n", i))
	}
	if err := os.MkdirAll(filepath.Join(fixture.work, "out", "delivery"), 0o700); err != nil {
		t.Fatalf("create receipt directory: %v", err)
	}
	writeTestFile(t, filepath.Join(fixture.work, "out", "delivery", "commit.md"), "Add large submission collection\n\nExercise guarded publication of a large path inventory.\n")
	prepared, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"submissions"},
	})
	if err != nil {
		t.Fatalf("prepare large change: %v", err)
	}
	receiptPath := filepath.Join(fixture.work, "out", "delivery", "prepare.json")
	info, err := os.Stat(receiptPath)
	if err != nil {
		t.Fatalf("inspect large receipt: %v", err)
	}
	if info.Size() <= 256*1024 {
		t.Fatalf("receipt size = %d, want more than 256 KiB", info.Size())
	}
	published, err := fixture.delivery.publish(context.Background(), publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	})
	if err != nil {
		t.Fatalf("publish large change: %v", err)
	}
	if !published.Verified || published.PullRequest == nil || published.PullRequest.HeadRefOID != prepared.HeadOID {
		t.Fatalf("large change publication was not verified")
	}
	if artifacts := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifacts != "" {
		report := fmt.Sprintf("paths=6000\nreceipt_bytes=%d\nverified=%t\n", info.Size(), published.Verified)
		if err := os.WriteFile(filepath.Join(artifacts, "large-publication.txt"), []byte(report), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeliveryRejectsOversizedPreparationReceipt(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	receiptPath := filepath.Join(fixture.work, "out", "delivery", "prepare.json")
	if err := os.WriteFile(receiptPath, []byte(strings.Repeat(" ", 4*1024*1024+1)), 0o600); err != nil {
		t.Fatalf("write oversized receipt: %v", err)
	}
	_, err := fixture.delivery.publish(context.Background(), publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("publish oversized receipt error = %v, want size refusal", err)
	}
	if fixture.forge.pull != nil {
		t.Fatal("oversized receipt created a pull request")
	}
}

func TestPrepareStagesExplicitTrackedDeletion(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery output directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Delete owned file\n\nExercise explicit deletion staging.\n",
	)
	if err := os.Remove(filepath.Join(fixture.work, "base.txt")); err != nil {
		t.Fatalf("remove tracked file: %v", err)
	}
	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			Paths:       []string{"base.txt"},
		},
	)
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(fixture.work, "base.txt")); !os.IsNotExist(err) {
		t.Fatalf("deleted path stat error = %v, want not-exist", err)
	}
	if output := runTestGit(
		t,
		fixture.work,
		"ls-tree",
		"--name-only",
		prepared.HeadOID,
		"--",
		"base.txt",
	); output != "" {
		t.Fatalf("prepared tree still contains deleted path: %q", output)
	}
}

func TestPrepareStagesPartialDeletionUnderExplicitDirectory(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "owned"),
		0o700,
	); err != nil {
		t.Fatalf("create owned directory: %v", err)
	}
	keepPath := filepath.Join(fixture.work, "owned", "keep.txt")
	deletePath := filepath.Join(fixture.work, "owned", "delete.txt")
	writeTestFile(t, keepPath, "before\n")
	writeTestFile(t, deletePath, "delete me\n")
	runTestGit(t, fixture.work, "add", "owned")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Owned directory",
		"-m",
		commitDisclaimer,
	)
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	writeTestFile(t, keepPath, "after\n")
	if err := os.Remove(deletePath); err != nil {
		t.Fatalf("remove tracked child: %v", err)
	}
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery output directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Update owned directory\n\nExercise partial deletion staging.\n",
	)
	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			Paths:       []string{"owned"},
			RewriteOID:  originalHead,
		},
	)
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	if contents, err := os.ReadFile(keepPath); err != nil {
		t.Fatalf("read retained child: %v", err)
	} else if string(contents) != "after\n" {
		t.Fatalf("retained child contents = %q", contents)
	}
	if output := runTestGit(
		t,
		fixture.work,
		"ls-tree",
		"--name-only",
		prepared.HeadOID,
		"--",
		"owned/delete.txt",
	); output != "" {
		t.Fatalf("prepared tree still contains deleted child: %q", output)
	}
}

func TestPrepareStagesSymlinkReplacedByExplicitDirectory(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	linkPath := filepath.Join(fixture.work, "owned")
	if err := os.Symlink("base.txt", linkPath); err != nil {
		t.Fatalf("create tracked symlink: %v", err)
	}
	runTestGit(t, fixture.work, "add", "owned")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Owned symlink",
		"-m",
		commitDisclaimer,
	)
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	if err := os.Remove(linkPath); err != nil {
		t.Fatalf("remove tracked symlink: %v", err)
	}
	if err := os.Mkdir(linkPath, 0o700); err != nil {
		t.Fatalf("create replacement directory: %v", err)
	}
	writeTestFile(t, filepath.Join(linkPath, "child.txt"), "child\n")
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery output directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Replace owned symlink\n\nExercise directory type-change staging.\n",
	)
	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			Paths:       []string{"owned"},
			RewriteOID:  originalHead,
		},
	)
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	if output := runTestGit(
		t,
		fixture.work,
		"ls-tree",
		"--name-only",
		prepared.HeadOID,
		"--",
		"owned/child.txt",
	); output != "owned/child.txt" {
		t.Fatalf("prepared tree child = %q, want owned/child.txt", output)
	}
}

func TestDeliveryPreparePublishVerify(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	receiptPath := filepath.Join(
		fixture.work,
		"out",
		"delivery",
		"prepare.json",
	)
	receiptInfo, err := os.Lstat(receiptPath)
	if err != nil {
		t.Fatalf("inspect preparation receipt: %v", err)
	}
	if !receiptInfo.Mode().IsRegular() || receiptInfo.Mode().Perm() != 0o600 {
		t.Fatalf("receipt mode = %v, want regular 0600", receiptInfo.Mode())
	}
	receiptContents, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatalf("read preparation receipt: %v", err)
	}
	if strings.Contains(string(receiptContents), fixture.remote) {
		t.Fatalf("receipt exposed a remote endpoint: %s", receiptContents)
	}
	if prepared.Receipt.PreparedHeadOID != prepared.HeadOID ||
		prepared.Receipt.PreparedTreeOID != prepared.TreeOID ||
		prepared.Receipt.Scope.Mode != scopeModePaths ||
		!containsString(prepared.Receipt.Scope.AggregatePaths, "feature.txt") {
		t.Fatalf("prepared receipt = %#v", prepared.Receipt)
	}
	published, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	if !published.Verified || published.PullRequest == nil ||
		published.PullRequest.HeadRefOID != prepared.HeadOID {
		t.Fatalf("publish() = %#v", published)
	}
	if !hasFinalLine(published.PullRequest.Body, pullRequestDisclaimer) {
		t.Fatalf("pull request body lacks disclaimer: %q", published.PullRequest.Body)
	}
	receipt, err := fixture.delivery.readReceipt(context.Background(), "out/delivery/prepare.json")
	if err != nil {
		t.Fatalf("read published preparation receipt: %v", err)
	}
	verified, err := fixture.delivery.verify(context.Background(), &receipt)
	if err != nil || !verified.Verified {
		t.Fatalf("verify publication = %#v, %v", verified, err)
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		artifact, err := json.MarshalIndent(map[string]any{
			"prepared_head":  prepared.HeadOID,
			"published_head": published.HeadOID,
			"remote_head":    published.RemoteHeadOID,
			"verified":       verified.Verified,
		}, "", "  ")
		if err != nil {
			t.Fatalf("encode publication artifact: %v", err)
		}
		if err := os.WriteFile(filepath.Join(directory, "source-publication.json"), append(artifact, '\n'), 0o600); err != nil {
			t.Fatalf("write publication artifact: %v", err)
		}
	}
}

// TestPublicationRetainsEvidenceAfterReviewFailure exercises a completed push
// followed by a failed review read and verification-only recovery.
func TestPublicationRetainsEvidenceAfterReviewFailure(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	fixture.forge.reviewError = errGitHubReviewEpochAdvanced
	publication, err := fixture.delivery.publish(context.Background(), publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	})
	if !errors.Is(err, errGitHubReviewEpochAdvanced) || publication == nil {
		t.Fatalf("review failure lost publication: %#v, %v", publication, err)
	}
	if publication.Status != "publication_verification_pending" ||
		publication.Verified || publication.HeadOID != prepared.HeadOID ||
		publication.RemoteHeadOID != prepared.HeadOID ||
		publication.PullRequest == nil || publication.PullRequest.URL == "" {
		t.Fatalf("incomplete publication evidence = %#v", publication)
	}
	remoteBefore := runTestGit(t, fixture.remote, "rev-parse", "refs/heads/feature")
	if remoteBefore != prepared.HeadOID {
		t.Fatalf("remote HEAD = %s, want published candidate %s", remoteBefore, prepared.HeadOID)
	}
	receipt, err := fixture.delivery.readReceipt(context.Background(), "out/delivery/prepare.json")
	if err != nil {
		t.Fatalf("read published preparation receipt: %v", err)
	}
	fixture.forge.reviewError = nil
	verified, err := fixture.delivery.verify(context.Background(), &receipt)
	if err != nil || !verified.Verified {
		t.Fatalf("verification-only recovery = %#v, %v", verified, err)
	}
	if remoteAfter := runTestGit(t, fixture.remote, "rev-parse", "refs/heads/feature"); remoteAfter != remoteBefore {
		t.Fatal("verification changed the remote feature ref")
	}
	if verified.Inspection.PullRequest == nil ||
		verified.Inspection.PullRequest.ID != publication.PullRequest.ID ||
		verified.Inspection.PullRequest.HeadRefOID != prepared.HeadOID {
		t.Fatalf("verification did not retain the published pull request: %#v", verified)
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		artifact, err := json.MarshalIndent(map[string]any{
			"initial_status":       publication.Status,
			"head_oid":             prepared.HeadOID,
			"remote_head_oid":      remoteBefore,
			"pull_request_url":     publication.PullRequest.URL,
			"verified_after_retry": verified.Verified,
		}, "", "  ")
		if err != nil {
			t.Fatalf("encode publication recovery artifact: %v", err)
		}
		if err := os.WriteFile(filepath.Join(directory, "publication-review-recovery.json"), append(artifact, '\n'), 0o600); err != nil {
			t.Fatalf("write publication recovery artifact: %v", err)
		}
	}
}

func TestDeliveryPreparePublishLargeCollection(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	directory := filepath.Join(fixture.work, "notes")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	for index := 0; index < 2200; index++ {
		name := fmt.Sprintf("%04d-%s.toml", index, strings.Repeat("field", 24))
		writeTestFile(t, filepath.Join(directory, name), "note_id = 1\n")
	}
	messageDirectory := filepath.Join(fixture.work, "out", "delivery")
	if err := os.MkdirAll(messageDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(messageDirectory, "commit.md"), "Export a large collection\n")
	prepared, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"notes"},
	})
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(messageDirectory, "prepare.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() <= 256*1024 || len(prepared.Receipt.Scope.AggregatePaths) != 2200 {
		t.Fatalf("large receipt: %d bytes and %d paths", info.Size(), len(prepared.Receipt.Scope.AggregatePaths))
	}
	published, err := fixture.delivery.publish(context.Background(), publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !published.Verified || published.PullRequest == nil || published.PullRequest.HeadRefOID != prepared.HeadOID {
		t.Fatal("large collection publication was not verified")
	}
	oversize := filepath.Join(messageDirectory, "oversize.json")
	if err := os.WriteFile(oversize, make([]byte, 4*1024*1024+1), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readStableReceiptFileWithLimit(oversize, preparationReceiptFileLimit); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("oversized receipt was not refused: %v", err)
	}
	if artifactDirectory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifactDirectory != "" {
		verification := fmt.Sprintf("files=2200\nreceipt_bytes=%d\npublication_verified=true\noversize_refused=true\n", info.Size())
		if err := os.WriteFile(filepath.Join(artifactDirectory, "large-collection-verification.txt"), []byte(verification), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPrepareConsolidatesExactOwnedLinearRange(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery output directory: %v", err)
	}
	messagePath := filepath.Join(
		fixture.work,
		"out",
		"delivery",
		"commit.md",
	)
	writeTestFile(
		t,
		messagePath,
		"Consolidate delivery fixture\n\nExercise guarded range ownership.\n",
	)
	featurePath := filepath.Join(fixture.work, "feature.txt")
	writeTestFile(t, featurePath, "first\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Owned feature",
		"-m",
		commitDisclaimer,
	)
	message, err := withCommitDisclaimer(
		"Consolidate delivery fixture\n\nExercise guarded range ownership.\n",
	)
	if err != nil {
		t.Fatalf("withCommitDisclaimer() error = %v", err)
	}
	projection, err := messageProjection(message)
	if err != nil {
		t.Fatalf("messageProjection() error = %v", err)
	}
	writeTestFile(t, featurePath, "second\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(t, fixture.work, "commit", "-m", "Follow-up correction")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	runTestGit(t, fixture.work, "push", "origin", "feature")
	fixture.forge.pull = &pullRequest{
		ID:                  "integration-pr",
		Number:              1,
		URL:                 "https://github.com/owner/repo/pull/1",
		State:               "OPEN",
		Title:               projection.Title,
		Body:                projection.Body,
		AuthorLogin:         "task-bot",
		BaseRefName:         "master",
		HeadRefName:         "feature",
		HeadRepositoryOwner: "owner",
		HeadRepositoryName:  "repo",
	}
	writeTestFile(t, featurePath, "pending\n")

	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile:    "out/delivery/commit.md",
			ReceiptFile:    "out/delivery/prepare.json",
			Paths:          []string{"feature.txt"},
			ConsolidateOID: originalHead,
		},
	)
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	if prepared.Inspection.UniqueCommitCount != 2 {
		t.Fatalf(
			"inspected feature count = %d, want 2",
			prepared.Inspection.UniqueCommitCount,
		)
	}
	if prepared.HeadOID == originalHead {
		t.Fatal("consolidation did not replace the original feature head")
	}
	if count := runTestGit(
		t,
		fixture.work,
		"rev-list",
		"--count",
		prepared.Receipt.BaseOID+".."+prepared.HeadOID,
	); count != "1" {
		t.Fatalf("final feature commit count = %s, want 1", count)
	}
	if parent := runTestGit(
		t,
		fixture.work,
		"show",
		"--no-patch",
		"--format=%P",
		prepared.HeadOID,
	); parent != prepared.Receipt.BaseOID {
		t.Fatalf("consolidated parent = %s, want %s", parent, prepared.Receipt.BaseOID)
	}
	if contents, err := os.ReadFile(featurePath); err != nil {
		t.Fatalf("read consolidated feature: %v", err)
	} else if string(contents) != "pending\n" {
		t.Fatalf("consolidated contents = %q", contents)
	}
	if !prepared.Receipt.ExpectedRemoteHead.Present ||
		prepared.Receipt.ExpectedRemoteHead.OID != originalHead {
		t.Fatalf(
			"receipt remote expectation = %#v, want %s",
			prepared.Receipt.ExpectedRemoteHead,
			originalHead,
		)
	}
}

func TestPrepareConsolidatesCleanExactOwnedRange(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery output directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Consolidate clean delivery fixture\n",
	)
	featurePath := filepath.Join(fixture.work, "feature.txt")
	writeTestFile(t, featurePath, "first\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Owned feature",
		"-m",
		commitDisclaimer,
	)
	writeTestFile(t, featurePath, "second\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(t, fixture.work, "commit", "-m", "Follow-up correction")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	originalTree := runTestGit(t, fixture.work, "rev-parse", "HEAD^{tree}")

	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile:    "out/delivery/commit.md",
			ReceiptFile:    "out/delivery/prepare.json",
			Paths:          []string{"feature.txt"},
			ConsolidateOID: originalHead,
		},
	)
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	if prepared.TreeOID != originalTree {
		t.Fatalf("consolidated tree = %s, want %s", prepared.TreeOID, originalTree)
	}
	if count := runTestGit(
		t,
		fixture.work,
		"rev-list",
		"--count",
		prepared.Receipt.BaseOID+".."+prepared.HeadOID,
	); count != "1" {
		t.Fatalf("final feature commit count = %s, want 1", count)
	}
}

func TestPrepareRefusesMultiCommitRangeWithoutExactConsolidation(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery output directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Consolidate delivery fixture\n",
	)
	featurePath := filepath.Join(fixture.work, "feature.txt")
	writeTestFile(t, featurePath, "first\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Owned feature",
		"-m",
		commitDisclaimer,
	)
	writeTestFile(t, featurePath, "second\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(t, fixture.work, "commit", "-m", "Follow-up correction")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	writeTestFile(t, featurePath, "pending\n")

	for name, consolidateOID := range map[string]string{
		"absent": "",
		"stale":  testOID('a'),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := fixture.delivery.prepare(
				context.Background(),
				prepareOptions{
					MessageFile:    "out/delivery/commit.md",
					ReceiptFile:    "out/delivery/prepare.json",
					Paths:          []string{"feature.txt"},
					ConsolidateOID: consolidateOID,
				},
			)
			if err == nil {
				t.Fatal("prepare() unexpectedly succeeded")
			}
			if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != originalHead {
				t.Fatalf("HEAD = %s after refusal, want %s", got, originalHead)
			}
		})
	}
}

func TestVerifyRefusesPullRequestMetadataEditAfterInitialInspection(
	t *testing.T,
) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	if _, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	); err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	fixture.forge.pullRequestsCalls = 0
	fixture.forge.mutateMetadataOnPullRequestsCall = 2

	report, err := fixture.delivery.verify(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "metadata changed") {
		t.Fatalf("verify() = %#v, error = %v", report, err)
	}
}

func TestPublishReceiptSupportsIdempotentRetry(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	options := publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	}
	first, err := fixture.delivery.publish(context.Background(), options)
	if err != nil {
		t.Fatalf("first publish() error = %v", err)
	}
	second, err := fixture.delivery.publish(context.Background(), options)
	if err != nil {
		t.Fatalf("idempotent publish() error = %v", err)
	}
	if !first.Verified || !second.Verified ||
		first.HeadOID != prepared.HeadOID || second.HeadOID != prepared.HeadOID ||
		first.PullRequest == nil || second.PullRequest == nil ||
		first.PullRequest.ID != second.PullRequest.ID {
		t.Fatalf("publish reports = %#v, %#v", first, second)
	}
}

func TestPublishUsesReceiptBaseWithoutExplicitBase(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	fixture.delivery.base = ""
	report, err := fixture.delivery.publish(context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	if !report.Verified || report.PullRequest == nil {
		t.Fatalf("publish report = %#v", report)
	}
}

func TestPublishRetryRepairsPushOnlyPartialUpdate(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	first := fixture.prepare(t)
	if _, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	); err != nil {
		t.Fatalf("initial publish() error = %v", err)
	}
	priorTitle := fixture.forge.pull.Title
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Revise delivery fixture\n\nExercise partial publication recovery.\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "revised\n")
	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			Paths:       []string{"feature.txt"},
			RewriteOID:  first.HeadOID,
		},
	)
	if err != nil {
		t.Fatalf("second prepare() error = %v", err)
	}
	fixture.forge.failUpdates = 1
	options := publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	}
	if _, err := fixture.delivery.publish(context.Background(), options); err == nil {
		t.Fatal("partial publish unexpectedly succeeded")
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != prepared.HeadOID {
		t.Fatalf("remote head = %s, want pushed candidate %s", got, prepared.HeadOID)
	}
	if fixture.forge.pull.Title != priorTitle {
		t.Fatalf("failed update changed title to %q", fixture.forge.pull.Title)
	}
	recovered, err := fixture.delivery.publish(context.Background(), options)
	if err != nil {
		t.Fatalf("recovery publish() error = %v", err)
	}
	if !recovered.Verified || recovered.PullRequest == nil ||
		recovered.PullRequest.Title == priorTitle ||
		recovered.PullRequest.HeadRefOID != prepared.HeadOID {
		t.Fatalf("recovered publish = %#v", recovered)
	}
}

func TestPublishReceiptRejectsReplacementPullRequest(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	first := fixture.prepare(t)
	if _, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	); err != nil {
		t.Fatalf("initial publish() error = %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Revise before replacement\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "revised\n")
	_, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			Paths:       []string{"feature.txt"},
			RewriteOID:  first.HeadOID,
		},
	)
	if err != nil {
		t.Fatalf("second prepare() error = %v", err)
	}
	oldRemote := runTestGit(t, fixture.remote, "rev-parse", "feature")
	fixture.forge.pull.ID = "replacement-pr"
	fixture.forge.pull.Number++
	fixture.forge.pull.URL = "https://github.com/owner/repo/pull/2"

	_, err = fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "identity") {
		t.Fatalf("publish() error = %v, want replacement identity refusal", err)
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != oldRemote {
		t.Fatalf("remote head = %s, want unchanged %s", got, oldRemote)
	}
}

func TestPublishRejectsUnknownReceiptFields(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	receiptPath := filepath.Join(
		fixture.work,
		"out",
		"delivery",
		"prepare.json",
	)
	contents, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatalf("read preparation receipt: %v", err)
	}
	tampered := strings.Replace(
		string(contents),
		"{",
		"{\n  \"unknown_field\": true,",
		1,
	)
	writeTestFile(t, receiptPath, tampered)
	_, err = fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || (!strings.Contains(err.Error(), "unknown field") &&
		!strings.Contains(err.Error(), "non-canonical key")) {
		t.Fatalf("publish() error = %v, want strict receipt refusal", err)
	}
	command := exec.Command(
		"git",
		"--git-dir",
		fixture.remote,
		"show-ref",
		"--verify",
		"--quiet",
		"refs/heads/feature",
	)
	if err := command.Run(); err == nil {
		t.Fatal("publish() mutated the remote with an invalid receipt")
	}
}

func TestPublishReceiptCannotTransferToAnotherBranch(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	runTestGit(t, fixture.work, "switch", "-c", "other-feature")
	fixture.forge.expectedHeadBranch = "other-feature"
	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "head ref") {
		t.Fatalf("publish() error = %v, want receipt branch refusal", err)
	}
	command := exec.Command(
		"git",
		"--git-dir",
		fixture.remote,
		"show-ref",
		"--verify",
		"--quiet",
		"refs/heads/other-feature",
	)
	if err := command.Run(); err == nil {
		t.Fatal("transferred receipt published another branch")
	}
}

func TestPublishReceiptCannotTransferToAnotherClone(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	transfer := filepath.Join(filepath.Dir(fixture.work), "transfer")
	runTestGit(
		t,
		filepath.Dir(fixture.work),
		"clone",
		"--no-local",
		fixture.work,
		transfer,
	)
	configureTestRepository(t, transfer)
	runTestGit(t, transfer, "remote", "set-url", "origin", fixture.remote)
	if err := os.MkdirAll(
		filepath.Join(transfer, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create transferred receipt directory: %v", err)
	}
	receiptContents, err := os.ReadFile(filepath.Join(
		fixture.work,
		"out",
		"delivery",
		"prepare.json",
	))
	if err != nil {
		t.Fatalf("read source receipt: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(transfer, "out", "delivery", "prepare.json"),
		string(receiptContents),
	)
	repository, err := openGitRepository(
		context.Background(),
		transfer,
		"git",
		func(string) string { return "" },
		&execRunner{},
	)
	if err != nil {
		t.Fatalf("open transferred repository: %v", err)
	}
	transferredDelivery := &delivery{
		repository:       repository,
		remote:           "origin",
		fetchEndpoint:    fixture.remote,
		pushEndpoint:     fixture.remote,
		base:             "master",
		remoteRepository: fixture.forge.repository,
		forge:            fixture.forge,
	}
	_, err = transferredDelivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "repository fingerprint") {
		t.Fatalf("publish() error = %v, want cross-clone refusal", err)
	}
	command := exec.Command(
		"git",
		"--git-dir",
		fixture.remote,
		"show-ref",
		"--verify",
		"--quiet",
		"refs/heads/feature",
	)
	if err := command.Run(); err == nil {
		t.Fatal("transferred receipt mutated the remote")
	}
}

func TestPreparationReceiptBindsDeliveryDestination(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	report, err := fixture.delivery.inspect(context.Background())
	if err != nil {
		t.Fatalf("inspect() error = %v", err)
	}
	tests := []struct {
		name   string
		mutate func(*delivery)
		want   string
	}{
		{
			name: "remote name",
			mutate: func(candidate *delivery) {
				candidate.remote = "upstream"
			},
			want: "remote name",
		},
		{
			name: "fetch endpoint",
			mutate: func(candidate *delivery) {
				candidate.fetchEndpoint += "-different"
			},
			want: "fetch endpoint",
		},
		{
			name: "push endpoint",
			mutate: func(candidate *delivery) {
				candidate.pushEndpoint += "-different"
			},
			want: "push endpoint",
		},
		{
			name: "remote repository",
			mutate: func(candidate *delivery) {
				candidate.remoteRepository.Name = "different"
			},
			want: "remote repository",
		},
		{
			name: "forge",
			mutate: func(candidate *delivery) {
				forgeCopy := *fixture.forge
				forgeCopy.name = "other-forge"
				candidate.forge = &forgeCopy
			},
			want: "forge adapter",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := *fixture.delivery
			test.mutate(&candidate)
			err := candidate.validateReceiptContext(
				context.Background(),
				prepared.Receipt,
				report,
			)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf(
					"validateReceiptContext() error = %v, want %q",
					err,
					test.want,
				)
			}
		})
	}
}

func TestPublishReceiptRejectsChangedRemoteHeadLease(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	runTestGit(t, fixture.seed, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(fixture.seed, "remote.txt"), "remote\n")
	runTestGit(t, fixture.seed, "add", "remote.txt")
	runTestGit(t, fixture.seed, "commit", "-m", "concurrent feature")
	runTestGit(t, fixture.seed, "push", "origin", "feature")
	remoteHead := runTestGit(t, fixture.seed, "rev-parse", "HEAD")
	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "receipt lease") {
		t.Fatalf("publish() error = %v, want changed-lease refusal", err)
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != remoteHead {
		t.Fatalf("remote head = %s, want concurrent %s", got, remoteHead)
	}
}

func TestPublishReceiptRejectsNonDescendantBaseRewrite(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	receiptPath := filepath.Join(
		fixture.work,
		"out",
		"delivery",
		"prepare.json",
	)
	receiptBefore, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatalf("read receipt before base rewrite: %v", err)
	}
	runTestGit(t, fixture.seed, "switch", "--orphan", "rewritten")
	runTestGit(t, fixture.seed, "commit", "--allow-empty", "-m", "rewrite base")
	runTestGit(t, fixture.seed, "push", "--force", "origin", "HEAD:master")
	_, err = fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "not a descendant") {
		t.Fatalf("publish() error = %v, want rewritten-base refusal", err)
	}
	if head := runTestGit(t, fixture.work, "rev-parse", "HEAD"); head != prepared.HeadOID {
		t.Fatalf("HEAD = %s, want unchanged %s", head, prepared.HeadOID)
	}
	receiptAfter, readErr := os.ReadFile(receiptPath)
	if readErr != nil || string(receiptAfter) != string(receiptBefore) {
		t.Fatalf("receipt changed after rewritten-base refusal: %v", readErr)
	}
}

func TestPublishAndVerifyPreserveUnstagedAndUntrackedChanges(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "user edit\n")
	writeTestFile(t, filepath.Join(fixture.work, "notes.txt"), "user notes\n")

	published, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	if !published.Verified {
		t.Fatalf("publish() = %#v, want verified", published)
	}
	if _, err := fixture.delivery.verify(context.Background(), nil); err != nil {
		t.Fatalf("verify() with unrelated dirty paths error = %v", err)
	}

	baseContents, err := os.ReadFile(filepath.Join(fixture.work, "base.txt"))
	if err != nil {
		t.Fatalf("read preserved tracked path: %v", err)
	}
	if string(baseContents) != "user edit\n" {
		t.Fatalf("tracked path contents = %q, want preserved edit", baseContents)
	}
	noteContents, err := os.ReadFile(filepath.Join(fixture.work, "notes.txt"))
	if err != nil {
		t.Fatalf("read preserved untracked path: %v", err)
	}
	if string(noteContents) != "user notes\n" {
		t.Fatalf("untracked path contents = %q, want preserved notes", noteContents)
	}
	status, err := fixture.delivery.repository.status(context.Background())
	if err != nil {
		t.Fatalf("status() error = %v", err)
	}
	if len(status.Staged) != 0 {
		t.Fatalf("publication staged unrelated paths: %v", status.Staged)
	}
	if !containsString(status.Unstaged, "base.txt") ||
		!containsString(status.Untracked, "notes.txt") {
		t.Fatalf("status() = %#v, want both preserved dirty paths", status)
	}
	if got := runTestGit(
		t,
		fixture.remote,
		"show",
		"feature:base.txt",
	); got != "base" {
		t.Fatalf("published base.txt = %q, want committed contents", got)
	}
	command := exec.Command(
		"git",
		"--git-dir",
		fixture.remote,
		"cat-file",
		"-e",
		"feature:notes.txt",
	)
	if err := command.Run(); err == nil {
		t.Fatal("publication included the unrelated untracked path")
	}
}

func TestPublishRefusesDirtyPreparedScope(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "unvalidated\n")

	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "prepared-scope") {
		t.Fatalf("publish() error = %v, want dirty prepared-scope refusal", err)
	}
	command := exec.Command(
		"git",
		"--git-dir",
		fixture.remote,
		"show-ref",
		"--verify",
		"--quiet",
		"refs/heads/feature",
	)
	if err := command.Run(); err == nil {
		t.Fatal("publish() pushed an older tree after task-scope content changed")
	}
}

func TestPrepareRefusesRemoteOnlyFeatureWork(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	runTestGit(t, fixture.seed, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(fixture.seed, "remote.txt"), "remote work\n")
	runTestGit(t, fixture.seed, "add", "remote.txt")
	runTestGit(t, fixture.seed, "commit", "-m", "remote feature")
	runTestGit(t, fixture.seed, "push", "origin", "feature")
	remoteHead := runTestGit(t, fixture.seed, "rev-parse", "HEAD")
	localHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Local candidate\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "local.txt"), "local work\n")

	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"local.txt"},
	})
	if err == nil || !strings.Contains(err.Error(), "absent from") {
		t.Fatalf("prepare() error = %v, want remote-only-work refusal", err)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != localHead {
		t.Fatalf("local HEAD = %s, want unchanged %s", got, localHead)
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != remoteHead {
		t.Fatalf("remote feature = %s, want unchanged %s", got, remoteHead)
	}

	_, err = fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile:      "out/delivery/commit.md",
		ReceiptFile:      "out/delivery/prepare.json",
		Paths:            []string{"local.txt"},
		ReplaceRemoteOID: testOID('a'),
	})
	if err == nil || !strings.Contains(err.Error(), "differs from") {
		t.Fatalf("prepare() error = %v, want exact-remote mismatch refusal", err)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != localHead {
		t.Fatalf("local HEAD = %s, want unchanged %s", got, localHead)
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != remoteHead {
		t.Fatalf("remote feature = %s, want unchanged %s", got, remoteHead)
	}
}

func TestPrepareAuthorizesExactTaskOwnedRebasedRemote(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "first\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Add feature",
		"-m",
		"First version.\n\n"+commitDisclaimer,
	)
	runTestGit(t, fixture.work, "push", "origin", "feature")
	remoteHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")

	fixture.advanceBase(t)
	runTestGit(t, fixture.work, "fetch", "origin", "master")
	runTestGit(t, fixture.work, "rebase", "origin/master")
	rebasedHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	if rebasedHead == remoteHead {
		t.Fatal("rebase did not replace the local task-owned commit OID")
	}

	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Update feature\n\nFinish the task-owned change.\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "second\n")

	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile:      "out/delivery/commit.md",
			ReceiptFile:      "out/delivery/prepare.json",
			Paths:            []string{"feature.txt"},
			RewriteOID:       rebasedHead,
			ReplaceRemoteOID: remoteHead,
		},
	)
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	if !prepared.Inspection.RemoteHeadDiverged ||
		prepared.Inspection.RemoteHeadOID != remoteHead {
		t.Fatalf("prepared inspection = %#v, want divergent remote %s", prepared.Inspection, remoteHead)
	}
	if !prepared.Receipt.ExpectedRemoteHead.Present ||
		prepared.Receipt.ExpectedRemoteHead.OID != remoteHead {
		t.Fatalf("receipt remote expectation = %#v, want %s", prepared.Receipt.ExpectedRemoteHead, remoteHead)
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != remoteHead {
		t.Fatalf("prepare mutated remote feature = %s, want %s", got, remoteHead)
	}
}

func TestPrepareRefusesUnnecessaryRemoteReplacementAuthorization(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "first\n")
	runTestGit(t, fixture.work, "add", "feature.txt")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Add feature",
		"-m",
		"First version.\n\n"+commitDisclaimer,
	)
	runTestGit(t, fixture.work, "push", "origin", "feature")
	localHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Update feature\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "second\n")

	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile:      "out/delivery/commit.md",
		ReceiptFile:      "out/delivery/prepare.json",
		Paths:            []string{"feature.txt"},
		RewriteOID:       localHead,
		ReplaceRemoteOID: localHead,
	})
	if err == nil || !strings.Contains(err.Error(), "unnecessary") {
		t.Fatalf("prepare() error = %v, want unnecessary-authorization refusal", err)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != localHead {
		t.Fatalf("local HEAD = %s, want unchanged %s", got, localHead)
	}
	if got := runTestGit(t, fixture.remote, "rev-parse", "feature"); got != localHead {
		t.Fatalf("remote feature = %s, want unchanged %s", got, localHead)
	}
}

func TestPrepareRefusesNonDescendantBaseBeforeCommit(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Candidate against stable base\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	runner := &integrationArgumentHookRunner{
		delegate: &execRunner{},
		argument: "add",
		hook: func(command) {
			runTestGit(t, fixture.seed, "switch", "--orphan", "rewritten")
			runTestGit(t, fixture.seed, "commit", "--allow-empty", "-m", "rewrite base")
			runTestGit(t, fixture.seed, "push", "--force", "origin", "HEAD:master")
		},
	}
	fixture.delivery.repository.runner = runner

	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if err == nil || !strings.Contains(err.Error(), "non-descendant") {
		t.Fatalf("prepare() error = %v, want non-descendant-base refusal", err)
	}
	if !runner.fired {
		t.Fatal("base rewrite hook did not run")
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != originalHead {
		t.Fatalf("HEAD = %s, want unchanged %s", got, originalHead)
	}
}

func TestPublishRequiresCleanIndex(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	writeTestFile(t, filepath.Join(fixture.work, "staged.txt"), "staged\n")
	runTestGit(t, fixture.work, "add", "staged.txt")

	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "clean index") {
		t.Fatalf("publish() error = %v, want clean-index refusal", err)
	}
	command := exec.Command(
		"git",
		"--git-dir",
		fixture.remote,
		"show-ref",
		"--verify",
		"--quiet",
		"refs/heads/feature",
	)
	if err := command.Run(); err == nil {
		t.Fatal("publish() pushed feature with a dirty index")
	}
}

func TestPrepareRejectsStaleRewriteWhenFeatureRangeIsEmpty(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(filepath.Join(fixture.work, "out", "delivery"), 0o700); err != nil {
		t.Fatalf("create message directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Add delivery fixture\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
		RewriteOID:  testOID('a'),
	})
	if err == nil || !strings.Contains(err.Error(), "feature range has no commit") {
		t.Fatalf("prepare() error = %v, want stale rewrite refusal", err)
	}
}

func TestPrepareRequiresMessageAndReceiptInSameTaskDirectory(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Candidate\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")

	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/other/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if err == nil || !strings.Contains(err.Error(), "share one") {
		t.Fatalf("prepare() error = %v, want task-directory refusal", err)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != originalHead {
		t.Fatalf("HEAD = %s, want unchanged %s", got, originalHead)
	}
}

func TestPublishStopsBeforePushWhenBaseAdvanceChangesHead(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	fixture.advanceBase(t)

	report, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	var checksRequired *checksRequiredError
	if !errors.As(err, &checksRequired) {
		t.Fatalf("publish() error = %v, want a pause for caller checks", err)
	}
	if report == nil || report.Status != "checks_required" ||
		report.HeadOID == prepared.HeadOID || report.Receipt == nil ||
		checksRequired.Report != report ||
		report.Receipt.PreparedHeadOID != report.HeadOID ||
		report.Receipt.PreparedTreeOID != report.TreeOID ||
		report.Receipt.BaseOID == prepared.Receipt.BaseOID {
		t.Fatalf("checks-required report = %#v, error = %#v", report, checksRequired)
	}
	command := exec.Command(
		"git",
		"--git-dir",
		fixture.remote,
		"show-ref",
		"--verify",
		"--quiet",
		"refs/heads/feature",
	)
	if err := command.Run(); err == nil {
		t.Fatal("publish() pushed feature before caller checks")
	}
	published, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err != nil {
		t.Fatalf("publish exact rebased candidate error = %v", err)
	}
	if !published.Verified || published.HeadOID != report.HeadOID {
		t.Fatalf("published = %#v, want exact rebased candidate", published)
	}
}

// TestRebasePreservesRemoteUntilPublication observes both existing and absent remote refs.
func TestRebasePreservesRemoteUntilPublication(t *testing.T) {
	for _, published := range []bool{false, true} {
		t.Run(fmt.Sprintf("published_%t", published), func(t *testing.T) {
			fixture := newIntegrationDeliveryFixture(t)
			prepared := fixture.prepare(t)
			oldRemote := ""
			if published {
				runTestGit(t, fixture.work, "push", "origin", "feature")
				oldRemote = prepared.HeadOID
			}
			fixture.advanceBase(t)
			before := runTestGit(t, fixture.work, "ls-remote", "origin", "refs/heads/feature")

			report, err := fixture.delivery.rebase(context.Background(), rebaseOptions{})
			if err != nil {
				t.Fatalf("rebase() error = %v", err)
			}
			if report == nil || report.Status != "rebased" || report.Pushed ||
				report.HeadOID == prepared.HeadOID || report.RemoteHeadOID != oldRemote ||
				report.RemoteHeadLease != oldRemote || report.BaseOID == prepared.Receipt.BaseOID {
				t.Fatalf("rebase report = %#v", report)
			}
			after := runTestGit(t, fixture.work, "ls-remote", "origin", "refs/heads/feature")
			if after != before || fixture.forge.pull != nil {
				t.Fatalf("rebase changed publication: before=%q after=%q PR=%#v", before, after, fixture.forge.pull)
			}
			for _, name := range []string{"feature.txt", "base-advance.txt"} {
				if _, err := os.Stat(filepath.Join(fixture.work, name)); err != nil {
					t.Fatalf("rebased content %s missing: %v", name, err)
				}
			}
			if artifactDirectory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifactDirectory != "" {
				writeTestFile(t, filepath.Join(artifactDirectory, fmt.Sprintf("rebase-local-%t.json", published)),
					fmt.Sprintf("{\"head\":\"%s\",\"base\":\"%s\",\"remote_before\":\"%s\",\"remote_after\":\"%s\",\"pushed\":false}\n", report.HeadOID, report.BaseOID, oldRemote, report.RemoteHeadOID))
			}
		})
	}
}

// TestRebaseAlreadyBasedPreservesRemote confirms that a no-op does not publish a local branch.
func TestRebaseAlreadyBasedPreservesRemote(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	report, err := fixture.delivery.rebase(context.Background(), rebaseOptions{})
	if err != nil {
		t.Fatalf("rebase() error = %v", err)
	}
	if report.Status != "already_based" || report.Pushed || report.HeadOID != prepared.HeadOID {
		t.Fatalf("already based report = %#v", report)
	}
	if remote := runTestGit(t, fixture.work, "ls-remote", "origin", "refs/heads/feature"); remote != "" {
		t.Fatalf("rebase published an already based branch: %s", remote)
	}
}

func TestRebaseRequiresCleanWorktree(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	fixture.advanceBase(t)
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "user edit\n")
	writeTestFile(t, filepath.Join(fixture.work, "notes.txt"), "user notes\n")

	_, err := fixture.delivery.rebase(context.Background(), rebaseOptions{})
	if err == nil || !strings.Contains(err.Error(), "clean worktree and index") {
		t.Fatalf("rebase() error = %v, want dirty-worktree refusal", err)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != originalHead {
		t.Fatalf("HEAD = %s, want unchanged %s", got, originalHead)
	}
	baseContents, readErr := os.ReadFile(filepath.Join(fixture.work, "base.txt"))
	if readErr != nil || string(baseContents) != "user edit\n" {
		t.Fatalf("tracked edit was not preserved: %q, %v", baseContents, readErr)
	}
	noteContents, readErr := os.ReadFile(filepath.Join(fixture.work, "notes.txt"))
	if readErr != nil || string(noteContents) != "user notes\n" {
		t.Fatalf("untracked file was not preserved: %q, %v", noteContents, readErr)
	}
}

func TestPublishRebaseDropsOnlyPathsIdenticalInAdvancedBase(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Add overlapping feature\n\nKeep the non-upstream portion.\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "shared.txt"), "shared\n")
	writeTestFile(t, filepath.Join(fixture.work, "unique.txt"), "unique\n")
	_, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			Paths:       []string{"shared.txt", "unique.txt"},
		},
	)
	if err != nil {
		t.Fatalf("prepare() error = %v", err)
	}
	writeTestFile(t, filepath.Join(fixture.seed, "shared.txt"), "shared\n")
	runTestGit(t, fixture.seed, "add", "shared.txt")
	runTestGit(t, fixture.seed, "commit", "-m", "Add shared base file")
	runTestGit(t, fixture.seed, "push", "origin", "master")

	report, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	var checksRequired *checksRequiredError
	if !errors.As(err, &checksRequired) {
		t.Fatalf("publish() error = %v, want a pause for caller checks", err)
	}
	if report == nil || report.Receipt == nil || !reflect.DeepEqual(
		report.Receipt.Scope.AggregatePaths,
		[]string{"unique.txt"},
	) {
		t.Fatalf("rebased aggregate scope = %#v", report)
	}
}

func TestPreparePathConstrainsExistingAggregateDiff(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	writeTestFile(t, filepath.Join(fixture.work, "outside.txt"), "outside\n")
	runTestGit(t, fixture.work, "add", "outside.txt")
	runTestGit(
		t,
		fixture.work,
		"commit",
		"-m",
		"Existing feature",
		"-m",
		"Existing body.\n\n"+commitDisclaimer,
	)
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Rewrite feature\n\nKeep the aggregate scope exact.\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "allowed.txt"), "allowed\n")
	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"allowed.txt"},
		RewriteOID:  originalHead,
	})
	if err == nil || !strings.Contains(err.Error(), "outside.txt") {
		t.Fatalf("prepare() error = %v, want aggregate-scope refusal", err)
	}
	if head := runTestGit(t, fixture.work, "rev-parse", "HEAD"); head != originalHead {
		t.Fatalf("HEAD = %s, want unchanged %s", head, originalHead)
	}
}

func TestPrepareChecksCandidateTreeWhitespaceOnly(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Add whitespace candidate\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "bad trailing space \n")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if err == nil || !strings.Contains(err.Error(), "diff --check") {
		t.Fatalf("prepare() error = %v, want immutable whitespace refusal", err)
	}
	const diagnostic = "feature.txt:1: trailing whitespace."
	if count := strings.Count(err.Error(), diagnostic); count != 1 {
		t.Fatalf(
			"prepare() error contains %q %d times, want exactly once: %v",
			diagnostic,
			count,
			err,
		)
	}
	if head := runTestGit(t, fixture.work, "rev-parse", "HEAD"); head != originalHead {
		t.Fatalf("HEAD = %s, want unchanged %s", head, originalHead)
	}
}

func TestMessageOnlyDirtyRequiredRebaseRefusesBeforeAmend(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	receiptPath := filepath.Join(
		fixture.work,
		"out",
		"delivery",
		"prepare.json",
	)
	receiptBefore, err := os.ReadFile(receiptPath)
	if err != nil {
		t.Fatalf("read receipt before message-only: %v", err)
	}
	fixture.advanceBase(t)
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "user edit\n")
	writeTestFile(t, filepath.Join(fixture.work, "notes.txt"), "notes\n")
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Refresh delivery message\n\nKeep the same tree.\n",
	)
	_, err = fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		MessageOnly: true,
		RewriteOID:  prepared.HeadOID,
	})
	if err == nil || !strings.Contains(err.Error(), "required rebase") {
		t.Fatalf("message-only prepare error = %v, want dirty-rebase refusal", err)
	}
	if head := runTestGit(t, fixture.work, "rev-parse", "HEAD"); head != prepared.HeadOID {
		t.Fatalf("HEAD = %s, want unchanged %s", head, prepared.HeadOID)
	}
	receiptAfter, readErr := os.ReadFile(receiptPath)
	if readErr != nil || string(receiptAfter) != string(receiptBefore) {
		t.Fatalf("receipt changed after refusal: error=%v", readErr)
	}
}

func TestMessageOnlyRetainsPreAmendSnapshot(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "user edit\n")
	writeTestFile(t, filepath.Join(fixture.work, "notes.txt"), "notes\n")
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Refresh message with stable snapshot\n",
	)
	runner := &integrationHookRunner{
		delegate: &execRunner{},
		hook: func(command) {
			fixture.advanceBase(t)
		},
	}
	fixture.delivery.repository.runner = runner
	refreshed, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			MessageOnly: true,
			RewriteOID:  prepared.HeadOID,
		},
	)
	if err != nil {
		t.Fatalf("message-only prepare across later base advance error = %v", err)
	}
	if !runner.fired || refreshed.HeadOID == prepared.HeadOID ||
		refreshed.Receipt.BaseOID != prepared.Receipt.BaseOID {
		t.Fatalf("message-only preparation = %#v, hook fired = %v", refreshed, runner.fired)
	}
}

func TestRewriteEvidenceAcceptsLocalProjectionWhenRemoteCannotProject(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	projection, err := fixture.delivery.repository.projection(
		context.Background(),
		prepared.HeadOID,
	)
	if err != nil {
		t.Fatalf("project local aggregate: %v", err)
	}
	runTestGit(t, fixture.seed, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(fixture.seed, "legacy.txt"), "legacy\n")
	runTestGit(t, fixture.seed, "add", "legacy.txt")
	runTestGit(t, fixture.seed, "commit", "-m", "Legacy remote tail")
	runTestGit(t, fixture.seed, "push", "origin", "feature")
	fixture.forge.pull = &pullRequest{
		ID:                  "integration-pr",
		Number:              1,
		URL:                 "https://github.com/owner/repo/pull/1",
		State:               "OPEN",
		Title:               projection.Title,
		Body:                projection.Body,
		AuthorLogin:         "task-bot",
		BaseRefName:         "master",
		HeadRefName:         "feature",
		HeadRepositoryOwner: "owner",
		HeadRepositoryName:  "repo",
	}
	report, err := fixture.delivery.inspect(context.Background())
	if err != nil {
		t.Fatalf("inspect divergent rewrite fixture: %v", err)
	}
	if report.RemoteHeadOID == "" || report.RemoteHeadOID == prepared.HeadOID {
		t.Fatalf("remote head = %q, want divergent legacy commit", report.RemoteHeadOID)
	}
	if err := fixture.delivery.requireRewriteEvidence(
		context.Background(),
		report,
		prepared.HeadOID,
	); err != nil {
		t.Fatalf("requireRewriteEvidence() error = %v", err)
	}
	report.PullRequest.Body = "Human-edited body.\n"
	if err := fixture.delivery.requireRewriteEvidence(
		context.Background(),
		report,
		prepared.HeadOID,
	); err == nil || !strings.Contains(err.Error(), "preserve possible human edits") {
		t.Fatalf("human-edit error = %v", err)
	}
}

func TestPreparedIndexReceiptFreezesExactAggregatePaths(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create delivery scratch: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Commit prepared index\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "a.txt"), "a\n")
	writeTestFile(t, filepath.Join(fixture.work, "b.txt"), "b\n")
	runTestGit(t, fixture.work, "add", "a.txt", "b.txt")
	prepared, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			UseIndex:    true,
		},
	)
	if err != nil {
		t.Fatalf("prepare --use-index error = %v", err)
	}
	if prepared.Receipt.Scope.Mode != scopeModeUseIndex ||
		!reflect.DeepEqual(
			prepared.Receipt.Scope.AuthorizedPaths,
			prepared.Receipt.Scope.AggregatePaths,
		) || !reflect.DeepEqual(
		prepared.Receipt.Scope.AggregatePaths,
		[]string{"a.txt", "b.txt"},
	) {
		t.Fatalf("prepared-index scope = %#v", prepared.Receipt.Scope)
	}
}

func TestMessageOnlyReceiptFreezesExactAggregatePaths(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Refresh aggregate message\n\nPreserve the exact tree.\n",
	)
	refreshed, err := fixture.delivery.prepare(
		context.Background(),
		prepareOptions{
			MessageFile: "out/delivery/commit.md",
			ReceiptFile: "out/delivery/prepare.json",
			MessageOnly: true,
			RewriteOID:  prepared.HeadOID,
		},
	)
	if err != nil {
		t.Fatalf("prepare --message-only error = %v", err)
	}
	if refreshed.TreeOID != prepared.TreeOID ||
		refreshed.HeadOID == prepared.HeadOID ||
		refreshed.Receipt.Scope.Mode != scopeModeMessageOnly ||
		!reflect.DeepEqual(
			refreshed.Receipt.Scope.AuthorizedPaths,
			refreshed.Receipt.Scope.AggregatePaths,
		) || !reflect.DeepEqual(
		refreshed.Receipt.Scope.AggregatePaths,
		[]string{"feature.txt"},
	) {
		t.Fatalf("message-only preparation = %#v", refreshed)
	}
}

func TestPublishRequiresCleanWorktreeBeforeRebase(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	fixture.advanceBase(t)
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "user edit\n")
	writeTestFile(t, filepath.Join(fixture.work, "notes.txt"), "user notes\n")

	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "before rebase") {
		t.Fatalf("publish() error = %v, want dirty-rebase refusal", err)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != originalHead {
		t.Fatalf("HEAD = %s, want unchanged %s", got, originalHead)
	}
	baseContents, readErr := os.ReadFile(filepath.Join(fixture.work, "base.txt"))
	if readErr != nil || string(baseContents) != "user edit\n" {
		t.Fatalf("tracked edit was not preserved: %q, %v", baseContents, readErr)
	}
	noteContents, readErr := os.ReadFile(filepath.Join(fixture.work, "notes.txt"))
	if readErr != nil || string(noteContents) != "user notes\n" {
		t.Fatalf("untracked file was not preserved: %q, %v", noteContents, readErr)
	}
}

// TestDeliveryInspectionUsesSelectedBase distinguishes release progress from unrelated trunk progress.
func TestDeliveryInspectionUsesSelectedBase(t *testing.T) {
	for _, published := range []bool{false, true} {
		t.Run(fmt.Sprintf("published_%t", published), func(t *testing.T) {
			fixture := newIntegrationDeliveryFixture(t)
			const releaseBase = "releases/2026.45"
			runTestGit(t, fixture.seed, "push", "origin", "HEAD:refs/heads/"+releaseBase)
			fixture.delivery.base = releaseBase
			prepared := fixture.prepare(t)
			if published {
				if _, err := fixture.delivery.publish(context.Background(), publishOptions{
					ReceiptFile: "out/delivery/prepare.json",
				}); err != nil {
					t.Fatalf("publish release-based feature: %v", err)
				}
				// The existing pull request supplies the base when the flag is absent.
				fixture.delivery.base = ""
			}
			fixture.advanceBase(t)
			runTestGit(t, fixture.work, "fetch", "origin", "+refs/heads/master:refs/remotes/origin/master")
			report, err := fixture.delivery.inspect(context.Background())
			if err != nil {
				t.Fatalf("inspect after unrelated master advance: %v", err)
			}
			if artifactDirectory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifactDirectory != "" {
				var contents strings.Builder
				if err := writeJSON(&contents, report); err != nil {
					t.Fatalf("encode selected-base inspection: %v", err)
				}
				writeTestFile(t, filepath.Join(artifactDirectory, fmt.Sprintf("selected-base-%t.json", published)), contents.String())
			}
			if err := ensureNoRefusals(report); err != nil {
				t.Fatalf("unrelated master advance refused selected release base: %v", err)
			}
			if report.Base != releaseBase || report.BaseOID != prepared.Receipt.BaseOID ||
				report.BaseStale || report.NeedsRebase || report.LocalHeadOID != prepared.HeadOID {
				t.Fatalf("inspection after unrelated master advance = %#v", report)
			}

			runTestGit(t, fixture.seed, "push", "origin", "HEAD:refs/heads/"+releaseBase)
			advanced, err := fixture.delivery.inspect(context.Background())
			if err != nil {
				t.Fatalf("inspect after selected release base advance: %v", err)
			}
			if advanced.Base != releaseBase || advanced.BaseOID == prepared.Receipt.BaseOID ||
				!advanced.NeedsRebase || advanced.LocalHeadOID != prepared.HeadOID {
				t.Fatalf("inspection failed to retain selected-base ancestry guard: %#v", advanced)
			}
			if artifactDirectory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifactDirectory != "" {
				var contents strings.Builder
				if err := writeJSON(&contents, advanced); err != nil {
					t.Fatalf("encode advancing-base inspection: %v", err)
				}
				writeTestFile(t, filepath.Join(artifactDirectory, fmt.Sprintf("selected-base-advanced-%t.json", published)), contents.String())
			}
		})
	}
}

func TestBaseSnapshotOIDIsNotPullRequestIdentity(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	if _, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		},
	); err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	before, err := fixture.delivery.inspect(context.Background())
	if err != nil {
		t.Fatalf("inspect() before base advance error = %v", err)
	}
	staleBaseOID := before.PullRequest.BaseRefOID
	fixture.advanceBase(t)

	if err := fixture.delivery.refreshPullRequestSafety(
		context.Background(),
		before,
	); err != nil {
		t.Fatalf("base snapshot advance failed concurrency check: %v", err)
	}
	fixture.forge.baseRefOIDOverride = staleBaseOID
	after, err := fixture.delivery.inspect(context.Background())
	if err != nil {
		t.Fatalf("inspect() with stale provider base snapshot error = %v", err)
	}
	if after.PullRequest.BaseRefOID == after.BaseOID {
		t.Fatalf("test setup did not produce a lagging base snapshot: %#v", after)
	}
	if err := ensureNoRefusals(after); err != nil {
		t.Fatalf("lagging provider base snapshot caused refusal: %v", err)
	}
	if !after.NeedsRebase {
		t.Fatalf("inspect() = %#v, want advancing base to require rebase", after)
	}
}

func TestInspectRejectsBranchChangeDuringCapturedOIDChecks(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	runTestGit(t, fixture.work, "branch", "intruder", prepared.HeadOID)
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("find git: %v", err)
	}
	wrapper := filepath.Join(fixture.work, "out", "delivery", "git-wrapper")
	writeTestFile(t, wrapper, fmt.Sprintf(`#!/bin/sh
for argument in "$@"; do
    if [ "$argument" = rev-list ]; then
        %q symbolic-ref HEAD refs/heads/intruder || exit $?
        break
    fi
done
exec %q "$@"
`, realGit, realGit))
	if err := os.Chmod(wrapper, 0o700); err != nil {
		t.Fatalf("make Git wrapper executable: %v", err)
	}
	fixture.delivery.repository.executable = wrapper

	_, err = fixture.delivery.inspect(context.Background())
	if err == nil || !strings.Contains(err.Error(), "branch or HEAD changed") {
		t.Fatalf("inspect() error = %v, want final branch/HEAD refusal", err)
	}
}

func TestPrepareRejectsStatusHiddenIndexFlagsBeforeStaging(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(filepath.Join(fixture.work, "out", "delivery"), 0o700); err != nil {
		t.Fatalf("create message directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Reject hidden edit\n",
	)
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	runTestGit(t, fixture.work, "update-index", "--assume-unchanged", "base.txt")
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "hidden edit\n")
	if got := runTestGit(t, fixture.work, "status", "--porcelain"); got != "" {
		t.Fatalf("test edit was not hidden by assume-unchanged: %q", got)
	}

	_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"base.txt"},
	})
	if err == nil || !strings.Contains(err.Error(), "non-default index flags") {
		t.Fatalf("prepare() error = %v, want index-flag refusal", err)
	}
	if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != originalHead {
		t.Fatalf("HEAD = %s, want unchanged %s", got, originalHead)
	}
	if got := runTestGit(t, fixture.work, "ls-files", "-v", "base.txt"); !strings.HasPrefix(got, "h ") {
		t.Fatalf("assume-unchanged flag was altered: %q", got)
	}
	if _, err := os.Lstat(filepath.Join(fixture.work, "out", "delivery", "prepare.json")); !os.IsNotExist(err) {
		t.Fatalf("prepare created a receipt after refusal: %v", err)
	}
}

func TestPublishRechecksPreparedScopeAfterFinalForgeRead(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	fixture.forge.pullRequestsCalls = 0
	fixture.forge.onPullRequestsCall = func(call int) {
		if call == 2 {
			writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "raced edit\n")
		}
	}

	_, err := fixture.delivery.publish(context.Background(), publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	})
	if err == nil || !strings.Contains(err.Error(), "prepared-scope") {
		t.Fatalf("publish() error = %v, want final prepared-scope refusal", err)
	}
	if _, err := runTestGitAllowFailure(
		fixture.remote,
		"rev-parse",
		"--verify",
		"refs/heads/feature",
	); err == nil {
		t.Fatal("publish pushed after the final forge-read race")
	}
}

func TestPublishAndVerifyRejectHiddenPreparedScopeEdits(t *testing.T) {
	t.Run("publish", func(t *testing.T) {
		fixture := newIntegrationDeliveryFixture(t)
		fixture.prepare(t)
		runTestGit(t, fixture.work, "update-index", "--assume-unchanged", "feature.txt")
		writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "hidden edit\n")
		_, err := fixture.delivery.publish(context.Background(), publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		})
		if err == nil || !strings.Contains(err.Error(), "index flag gate") {
			t.Fatalf("publish() error = %v, want hidden-index refusal", err)
		}
		if _, err := runTestGitAllowFailure(
			fixture.remote,
			"rev-parse",
			"--verify",
			"refs/heads/feature",
		); err == nil {
			t.Fatal("publish pushed a status-hidden prepared-scope edit")
		}
	})

	t.Run("verify", func(t *testing.T) {
		fixture := newIntegrationDeliveryFixture(t)
		fixture.prepare(t)
		if _, err := fixture.delivery.publish(context.Background(), publishOptions{
			ReceiptFile: "out/delivery/prepare.json",
		}); err != nil {
			t.Fatalf("publish() error = %v", err)
		}
		receipt, err := fixture.delivery.readReceipt(
			context.Background(),
			"out/delivery/prepare.json",
		)
		if err != nil {
			t.Fatalf("readReceipt() error = %v", err)
		}
		runTestGit(t, fixture.work, "update-index", "--assume-unchanged", "feature.txt")
		writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "hidden edit\n")
		report, err := fixture.delivery.verify(context.Background(), &receipt)
		if err == nil || !strings.Contains(err.Error(), "index flag gate") {
			t.Fatalf("verify() = %#v, error = %v, want hidden-index refusal", report, err)
		}
	})
}

func TestVerifyRechecksPreparedScopeAfterFinalForgeRead(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	if _, err := fixture.delivery.publish(context.Background(), publishOptions{
		ReceiptFile: "out/delivery/prepare.json",
	}); err != nil {
		t.Fatalf("publish() error = %v", err)
	}
	receipt, err := fixture.delivery.readReceipt(
		context.Background(),
		"out/delivery/prepare.json",
	)
	if err != nil {
		t.Fatalf("readReceipt() error = %v", err)
	}
	fixture.forge.pullRequestsCalls = 0
	fixture.forge.onPullRequestsCall = func(call int) {
		if call == 2 {
			writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "raced edit\n")
		}
	}

	report, err := fixture.delivery.verify(context.Background(), &receipt)
	if err == nil || !strings.Contains(err.Error(), "prepared-scope") {
		t.Fatalf("verify() = %#v, error = %v, want final scope refusal", report, err)
	}
}

func TestVerifyRejectsUnboundPullRequestAfterExpectedAbsence(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	runTestGit(
		t,
		fixture.work,
		"push",
		"origin",
		prepared.HeadOID+":refs/heads/feature",
	)
	projection, err := fixture.delivery.repository.projection(
		context.Background(),
		prepared.HeadOID,
	)
	if err != nil {
		t.Fatalf("projection() error = %v", err)
	}
	fixture.forge.pull = &pullRequest{
		ID:                  "appeared-pr",
		Number:              1,
		URL:                 "https://github.com/owner/repo/pull/1",
		State:               "OPEN",
		Title:               projection.Title,
		Body:                projection.Body,
		BaseRefName:         "master",
		HeadRefName:         "feature",
		HeadRepositoryOwner: fixture.forge.repository.Owner,
		HeadRepositoryName:  fixture.forge.repository.Name,
	}

	report, err := fixture.delivery.verify(
		context.Background(),
		&prepared.Receipt,
	)
	if err == nil || !strings.Contains(err.Error(), "exact desired state") {
		t.Fatalf("verify() = %#v, error = %v, want unbound-identity refusal", report, err)
	}
}

func hasEnvironmentPrefix(environment []string, prefix string) bool {
	for _, value := range environment {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

// TestGitHubPreflightFailureStopsCLI checks API failures independently of Git transport.
func TestGitHubPreflightFailureStopsCLI(t *testing.T) {
	for _, transport := range []string{"ssh", "gh"} {
		for _, status := range []string{"401", "403"} {
			t.Run(transport+"_"+status, func(t *testing.T) {
				fixture := newIntegrationDeliveryFixture(t)
				runTestGit(t, fixture.work, "remote", "set-url", "origin", "git@github.com:owner/repo.git")
				head := runTestGit(t, fixture.work, "rev-parse", "HEAD")
				state := runTestGit(t, fixture.work, "status", "--porcelain=v2")
				refs := runTestGit(t, fixture.work, "show-ref")
				helper := filepath.Join(t.TempDir(), "gh")
				writeTestFile(t, helper, `#!/bin/sh
printf 'call\n' >> "$0.calls"
[ "$*" = 'repo view github.com/owner/repo --json id,nameWithOwner,url' ] || exit 2
printf 'gh: API rejected request (HTTP `+status+`)\nAuthorization: Bearer fixture-secret\n' >&2
exit 1
`)
				if err := os.Chmod(helper, 0o700); err != nil {
					t.Fatalf("chmod forge fixture: %v", err)
				}
				var output strings.Builder
				command := NewCommand(context.Background(), os.Getenv, &output)
				command.SetArgs([]string{
					"inspect", "--repository", fixture.work, "--base", "master",
					"--forge-cli", helper, "--git-transport", transport,
				})
				err := command.Execute()
				var commandErr *commandError
				if !errors.As(err, &commandErr) || commandErr.Result.ExitCode != 1 {
					t.Fatalf("inspect lost preflight failure: %v", err)
				}
				if strings.Contains(err.Error()+output.String(), "fixture-secret") {
					t.Fatal("inspect exposed a credential from forge stderr")
				}
				if got := strings.Contains(err.Error(), "API authentication"); got != (status == "401") {
					t.Fatalf("authentication advice for HTTP %s: %v", status, err)
				}
				calls, readErr := os.ReadFile(helper + ".calls")
				if readErr != nil || string(calls) != "call\n" {
					t.Fatalf("preflight was retried or skipped: calls=%q error=%v", calls, readErr)
				}
				if runTestGit(t, fixture.work, "rev-parse", "HEAD") != head ||
					runTestGit(t, fixture.work, "status", "--porcelain=v2") != state ||
					runTestGit(t, fixture.work, "show-ref") != refs {
					t.Fatal("failed preflight changed Git state")
				}
				t.Logf("Verified API failure: transport=%s status=%s calls=1 unchanged_git=true redacted=true", transport, status)
			})
		}
	}
}

// TestGHTransportPublication exercises helper execution and guarded publication using real local Git.
func TestGHTransportPublication(t *testing.T) {
	for _, failure := range []bool{false, true} {
		t.Run(fmt.Sprintf("authentication_failure_%t", failure), func(t *testing.T) {
			fixture := newIntegrationDeliveryFixture(t)
			helper := filepath.Join(t.TempDir(), "gh helper ' literal$")
			script := "#!/bin/sh\n[ \"$1 $2 $3\" = 'auth git-credential get' ] || exit 2\ncat >/dev/null\nprintf 'username=fixture\\npassword=fixture-only\\n\\n'\n"
			if failure {
				script = "#!/bin/sh\nexit 1\n"
			}
			writeTestFile(t, helper, script)
			if err := os.Chmod(helper, 0o700); err != nil {
				t.Fatalf("chmod helper: %v", err)
			}
			audit := &ghCredentialAuditRunner{delegate: &execRunner{}}
			fixture.delivery.repository.runner = audit
			fixture.delivery.repository.gitTransport = "gh"
			fixture.delivery.repository.credentialExecutable = helper
			t.Setenv("GH_DEBUG", "api")
			if failure {
				_, err := fixture.delivery.inspect(context.Background())
				if err == nil {
					t.Fatal("authentication failure must stop inspection")
				}
				if got := runTestGit(t, fixture.remote, "for-each-ref", "refs/heads/feature"); got != "" {
					t.Fatal("authentication failure published a feature ref")
				}
				return
			}
			prepared := fixture.prepare(t)
			published, err := fixture.delivery.publish(context.Background(), publishOptions{ReceiptFile: "out/delivery/prepare.json"})
			if err != nil {
				t.Fatalf("publish with gh credentials: %v", err)
			}
			if !published.Verified || audit.calls == 0 {
				t.Fatal("publication or credential protocol was not verified")
			}
			if got := runTestGit(t, fixture.remote, "rev-parse", "refs/heads/feature"); got != prepared.HeadOID {
				t.Fatal("published head differs from prepared candidate")
			}
			receipt, err := fixture.delivery.readReceipt(context.Background(), "out/delivery/prepare.json")
			if err != nil || receipt.GitTransport != "gh" {
				t.Fatalf("receipt does not bind gh transport: %v", err)
			}
			// The test log is a repeatable receipt for the isolated fixture's exact published candidate.
			t.Logf("Verified gh transport publication: head=%s credential_calls=%d", prepared.HeadOID, audit.calls)
		})
	}
}

// ghCredentialAuditRunner checks credential isolation before delegating real Git transport operations.
type ghCredentialAuditRunner struct {
	// delegate executes real local Git against the isolated fixture remote.
	delegate commandRunner
	// calls counts successful credential protocol exchanges.
	calls int
}

// Run validates the isolated Git credential protocol without exposing its dummy password.
func (r *ghCredentialAuditRunner) Run(ctx context.Context, request command) (commandResult, error) {
	if containsString(request.Env, "GIT_TERMINAL_PROMPT=0") {
		environment := environmentMap(mergeEnvironment(os.Environ(), request.Env, request.UnsetEnv, request.UnsetEnvPrefixes))
		if _, present := environment["GH_DEBUG"]; present {
			return commandResult{}, fmt.Errorf("transport retained credential diagnostics")
		}
		config, err := commandConfiguration(environment)
		if err != nil {
			return commandResult{}, fmt.Errorf("read transport configuration: %w", err)
		}
		for key, want := range map[string]string{"http.sslVerify": "true", "http.followRedirects": "false", "http.proxy": "", "credential.helper": "", "credential.interactive": "false", "protocol.https.allow": "always"} {
			if config[key] != want {
				return commandResult{}, fmt.Errorf("unsafe transport configuration for %s", key)
			}
		}
		if config["credential.https://github.com.helper"] == "" {
			return commandResult{}, fmt.Errorf("missing host-scoped gh credential helper")
		}
		credentialRequest := request
		credentialRequest.Args = []string{"credential", "fill"}
		credentialRequest.Stdin = "protocol=https\nhost=github.com\n\n"
		result, err := r.delegate.Run(ctx, credentialRequest)
		if err != nil {
			return commandResult{}, fmt.Errorf("fixture credential helper failed: %w", err)
		}
		if !strings.Contains(result.Stdout, "password=fixture-only") {
			return commandResult{}, fmt.Errorf("credential helper returned an unexpected fixture response")
		}
		r.calls++
	}
	result, err := r.delegate.Run(ctx, request)
	if err != nil {
		return result, fmt.Errorf("run fixture Git command: %w", err)
	}
	return result, nil
}

// TestReceiptRefusesTransportChange preserves the prepared authentication boundary.
func TestReceiptRefusesTransportChange(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.prepare(t)
	fixture.delivery.repository.gitTransport = "gh"
	fixture.delivery.repository.credentialExecutable = "gh"
	_, err := fixture.delivery.publish(context.Background(), publishOptions{ReceiptFile: "out/delivery/prepare.json"})
	if err == nil || !strings.Contains(err.Error(), "transport") {
		t.Fatalf("changed transport was not refused: %v", err)
	}
	if got := runTestGit(t, fixture.remote, "for-each-ref", "refs/heads/feature"); got != "" {
		t.Fatal("changed transport published a feature ref")
	}
}

// TestGHTransportEndpoints checks explicit selection and canonical repository boundaries before network access.
func TestGHTransportEndpoints(t *testing.T) {
	for _, endpoint := range []string{"git@github.com:owner/repository.git", "ssh://git@github.com/owner/repository.git", "https://github.com/owner/repository.git"} {
		t.Run(endpoint, func(t *testing.T) {
			repositoryPath := newDeliveryConstructorRepository(t)
			runTestGit(t, repositoryPath, "remote", "add", "origin", endpoint)
			runner := &constructorIsolationRunner{delegate: &execRunner{}}
			d, err := newDelivery(context.Background(), deliveryConfig{Repository: repositoryPath, ForgeCLI: "gh-test", gitTransport: "gh"}, func(string) string { return "" }, runner)
			if err != nil {
				t.Fatalf("configure gh transport: %v", err)
			}
			bound, err := d.repository.bindGitRemote(endpoint)
			if err != nil || bound.endpoint != "https://github.com/owner/repository.git" {
				t.Fatalf("incorrect HTTPS endpoint: %v", err)
			}
			if d.fetchEndpoint != endpoint {
				t.Fatal("captured raw endpoint changed")
			}
		})
	}
	for _, endpoint := range []string{"http://github.com/owner/repository.git", "https://user:secret@github.com/owner/repository.git", "https://github.com:443/owner/repository.git", "https://github.com/owner/%72epository.git", "https://github.com/owner/repository.git?query=1", "https://attacker.invalid/owner/repository.git", "git@attacker.invalid:owner/repository.git"} {
		t.Run(endpoint, func(t *testing.T) {
			repositoryPath := newDeliveryConstructorRepository(t)
			runTestGit(t, repositoryPath, "remote", "add", "origin", endpoint)
			runner := &constructorIsolationRunner{delegate: &execRunner{}}
			_, err := newDelivery(context.Background(), deliveryConfig{Repository: repositoryPath, ForgeCLI: "gh-test", gitTransport: "gh"}, func(string) string { return "" }, runner)
			if err == nil || runner.forgeCalls != 0 || runner.gitNetworkCalls != 0 {
				t.Fatal("unsafe endpoint was not refused before network access")
			}
		})
	}
	repositoryPath := newDeliveryConstructorRepository(t)
	runTestGit(t, repositoryPath, "remote", "add", "origin", "git@github.com:owner/repository.git")
	runner := &constructorIsolationRunner{delegate: &execRunner{}}
	_, err := newDelivery(context.Background(), deliveryConfig{Repository: repositoryPath, ForgeCLI: "gh-test", gitTransport: "automatic"}, func(string) string { return "" }, runner)
	if err == nil || runner.forgeCalls != 0 {
		t.Fatal("invalid transport selection was not refused")
	}
}

// TestGHTransportForgeExecutableDirectories checks preflight and helper execution from different directories.
func TestGHTransportForgeExecutableDirectories(t *testing.T) {
	for _, selection := range []string{"relative", "path"} {
		t.Run(selection, func(t *testing.T) {
			fixture := newIntegrationDeliveryFixture(t)
			relativeHelper := filepath.Join("out", "delivery", "gh fixture ' literal$")
			helper := filepath.Join(fixture.work, relativeHelper)
			script := `#!/bin/sh
case "$1 $2" in
 'repo view') printf '%s\n' '{"id":"fixture-repo","nameWithOwner":"owner/repo","url":"https://github.com/owner/repo"}' ;;
 'auth git-credential')
   [ "$3" = get ] || exit 2
   cat >/dev/null
   printf 'username=fixture\npassword=fixture-only\n\n'
   ;;
 *) exit 3 ;;
esac
`
			if err := os.MkdirAll(filepath.Dir(helper), 0o700); err != nil {
				t.Fatalf("create helper directory: %v", err)
			}
			writeTestFile(t, helper, script)
			if err := os.Chmod(helper, 0o700); err != nil {
				t.Fatalf("chmod helper: %v", err)
			}
			executable := "." + string(filepath.Separator) + relativeHelper
			if selection == "path" {
				executable = filepath.Base(helper)
				t.Setenv("PATH", filepath.Dir(helper)+string(os.PathListSeparator)+os.Getenv("PATH"))
			}
			runTestGit(t, fixture.work, "remote", "set-url", "origin", "git@github.com:owner/repo.git")
			d, err := newDelivery(context.Background(), deliveryConfig{Repository: fixture.work, ForgeCLI: executable, gitTransport: "gh", Base: "master"}, func(string) string { return "" }, &execRunner{})
			if err != nil {
				t.Fatalf("forge executable preflight failed: %v", err)
			}
			audit := &ghCredentialAuditRunner{delegate: &execRunner{}}
			d.repository.runner = audit
			d.forge = fixture.forge
			// Keep the same real credential protocol while publishing to a local fixture remote.
			d.fetchEndpoint, d.pushEndpoint = fixture.remote, fixture.remote
			fixture.delivery = d
			prepared := fixture.prepare(t)
			published, err := d.publish(context.Background(), publishOptions{ReceiptFile: "out/delivery/prepare.json"})
			if err != nil {
				t.Fatalf("publish after successful forge preflight: %v", err)
			}
			if !published.Verified || audit.calls == 0 {
				t.Fatal("credential helper and publication were not verified")
			}
			if got := runTestGit(t, fixture.remote, "rev-parse", "refs/heads/feature"); got != prepared.HeadOID {
				t.Fatal("published head differs from prepared candidate")
			}
			t.Logf("Verified %s forge executable: head=%s credential_calls=%d", selection, prepared.HeadOID, audit.calls)
		})
	}
}
