package main

import (
	"bytes"
	"context"
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
	return nil, errors.New("unexpected review inspection")
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
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
}

func TestVerifyRefusesPullRequestMetadataEditAfterInitialInspection(
	t *testing.T,
) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	if _, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
		ValidatedHead: prepared.HeadOID,
		ReceiptFile:   "out/delivery/prepare.json",
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

func TestPublishRetryRepairsPushOnlyPartialUpdate(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	first := fixture.prepare(t)
	if _, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: first.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
		ValidatedHead: prepared.HeadOID,
		ReceiptFile:   "out/delivery/prepare.json",
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
			ValidatedHead: first.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	oldRemote := runTestGit(t, fixture.remote, "rev-parse", "feature")
	fixture.forge.pull.ID = "replacement-pr"
	fixture.forge.pull.Number++
	fixture.forge.pull.URL = "https://github.com/owner/repo/pull/2"

	_, err = fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
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
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
	runTestGit(t, fixture.work, "switch", "-c", "other-feature")
	fixture.forge.expectedHeadBranch = "other-feature"
	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
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
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
	runTestGit(t, fixture.seed, "switch", "-c", "feature")
	writeTestFile(t, filepath.Join(fixture.seed, "remote.txt"), "remote\n")
	runTestGit(t, fixture.seed, "add", "remote.txt")
	runTestGit(t, fixture.seed, "commit", "-m", "concurrent feature")
	runTestGit(t, fixture.seed, "push", "origin", "feature")
	remoteHead := runTestGit(t, fixture.seed, "rev-parse", "HEAD")
	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "user edit\n")
	writeTestFile(t, filepath.Join(fixture.work, "notes.txt"), "user notes\n")

	published, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "unvalidated\n")

	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
	writeTestFile(t, filepath.Join(fixture.work, "staged.txt"), "staged\n")
	runTestGit(t, fixture.work, "add", "staged.txt")

	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
		},
	)
	var revalidation *revalidationRequiredError
	if !errors.As(err, &revalidation) {
		t.Fatalf("publish() error = %v, want revalidation", err)
	}
	if report == nil || report.Status != "revalidation_required" ||
		report.HeadOID == prepared.HeadOID || report.Receipt == nil ||
		revalidation.Report != report ||
		report.Receipt.PreparedHeadOID != report.HeadOID ||
		report.Receipt.PreparedTreeOID != report.TreeOID ||
		report.Receipt.BaseOID == prepared.Receipt.BaseOID {
		t.Fatalf("revalidation report = %#v, error = %#v", report, revalidation)
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
		t.Fatal("publish() pushed feature before revalidation")
	}
	published, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: report.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
		},
	)
	if err != nil {
		t.Fatalf("publish exact rebased candidate error = %v", err)
	}
	if !published.Verified || published.HeadOID != report.HeadOID {
		t.Fatalf("published = %#v, want exact rebased candidate", published)
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
	prepared := fixture.prepare(t)
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	fixture.advanceBase(t)
	writeTestFile(t, filepath.Join(fixture.work, "base.txt"), "user edit\n")
	writeTestFile(t, filepath.Join(fixture.work, "notes.txt"), "user notes\n")

	_, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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

func TestBaseSnapshotOIDIsNotPullRequestIdentity(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	if _, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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

func TestPrepareFinalizationFailureRestoresEntryState(t *testing.T) {
	for _, advanceBase := range []bool{false, true} {
		t.Run(fmt.Sprintf("base_advanced_%t", advanceBase), func(t *testing.T) {
			fixture := newIntegrationDeliveryFixture(t)
			if advanceBase {
				fixture.advanceBase(t)
			}
			if err := os.MkdirAll(filepath.Join(fixture.work, "out", "delivery"), 0o700); err != nil {
				t.Fatalf("create message directory: %v", err)
			}
			writeTestFile(
				t,
				filepath.Join(fixture.work, "out", "delivery", "commit.md"),
				"Transactional preparation\n",
			)
			writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
			originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
			originalIndex := runTestGit(t, fixture.work, "write-tree")
			fixture.delivery.repository.runner = &failPreparationFinalizationRunner{
				delegate: &execRunner{},
			}

			_, err := fixture.delivery.prepare(context.Background(), prepareOptions{
				MessageFile: "out/delivery/commit.md",
				ReceiptFile: "out/delivery/prepare.json",
				Paths:       []string{"feature.txt"},
			})
			if err == nil || !strings.Contains(err.Error(), "synthetic finalization failure") {
				t.Fatalf("prepare() error = %v, want injected finalization failure", err)
			}
			if got := runTestGit(t, fixture.work, "rev-parse", "HEAD"); got != originalHead {
				t.Fatalf("HEAD = %s, want restored %s", got, originalHead)
			}
			if got := runTestGit(t, fixture.work, "write-tree"); got != originalIndex {
				t.Fatalf(
					"index tree = %s, want restored %s; prepare error: %v",
					got,
					originalIndex,
					err,
				)
			}
			if got := runTestGit(t, fixture.work, "status", "--porcelain"); got != "?? feature.txt" {
				t.Fatalf("restored status = %q, want original untracked feature", got)
			}
			if _, err := os.Lstat(filepath.Join(fixture.work, "base-advance.txt")); advanceBase && !os.IsNotExist(err) {
				t.Fatalf("rebased base file survived rollback: %v", err)
			}
		})
	}
}

func TestPrepareCancellationAfterIndexInstallRestoresEntryState(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create message directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Cancel after exact index installation\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	originalIndex := runTestGit(t, fixture.work, "write-tree")
	ctx, cancel := context.WithCancel(context.Background())
	runner := &cancelAfterIndexInstallRunner{
		delegate: &execRunner{},
		cancel:   cancel,
	}
	fixture.delivery.repository.runner = runner

	_, err := fixture.delivery.prepare(ctx, prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("prepare() error = %v, want context cancellation", err)
	}
	if !runner.fired {
		t.Fatal("cancellation boundary did not fire after index installation")
	}
	assertPreparationEntryState(
		t,
		fixture.work,
		originalHead,
		originalIndex,
		"?? feature.txt",
	)
}

func TestPrepareCancellationAfterCandidateInstallRestoresEntryState(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create message directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Cancel after exact candidate installation\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	originalIndex := runTestGit(t, fixture.work, "write-tree")
	ctx, cancel := context.WithCancel(context.Background())
	runner := &cancelAfterCandidateInstallRunner{
		delegate: &execRunner{},
		cancel:   cancel,
	}
	fixture.delivery.repository.runner = runner

	_, err := fixture.delivery.prepare(ctx, prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("prepare() error = %v, want context cancellation", err)
	}
	if !runner.fired {
		t.Fatal("cancellation boundary did not fire after candidate installation")
	}
	assertPreparationEntryState(
		t,
		fixture.work,
		originalHead,
		originalIndex,
		"?? feature.txt",
	)
}

func TestPreparePostCommitIndexLockCleanupFailureRestoresEntryState(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create message directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Fail after exact commit installation\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	originalIndex := runTestGit(t, fixture.work, "write-tree")
	originalIndexContents, err := os.ReadFile(
		filepath.Join(fixture.work, ".git", "index"),
	)
	if err != nil {
		t.Fatalf("read original index: %v", err)
	}
	cleanupFailures := 0
	fixture.delivery.repository.removeLockFile = func(path string) error {
		removeErr := os.Remove(path)
		if removeErr == nil && cleanupFailures == 0 &&
			filepath.Base(path) == "index.lock" {
			cleanupFailures++
			return fmt.Errorf("synthetic post-commit index-lock cleanup failure")
		}
		return removeErr
	}

	_, err = fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if err == nil || !strings.Contains(
		err.Error(),
		"synthetic post-commit index-lock cleanup failure",
	) {
		t.Fatalf("prepare() error = %v, want post-install cleanup failure", err)
	}
	if strings.Contains(err.Error(), "roll back failed preparation") {
		t.Fatalf("prepare() failed to recognize installed commit: %v", err)
	}
	if cleanupFailures != 1 {
		t.Fatalf("injected cleanup failures = %d, want 1", cleanupFailures)
	}
	restoredIndexContents, err := os.ReadFile(
		filepath.Join(fixture.work, ".git", "index"),
	)
	if err != nil {
		t.Fatalf("read restored index: %v", err)
	}
	if !bytes.Equal(restoredIndexContents, originalIndexContents) {
		t.Fatal("rollback did not restore the exact original index bytes")
	}
	assertPreparationEntryState(
		t,
		fixture.work,
		originalHead,
		originalIndex,
		"?? feature.txt",
	)
}

func TestPreparePostRebaseHEADLockCleanupFailureRestoresEntryState(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	fixture.advanceBase(t)
	if err := os.MkdirAll(
		filepath.Join(fixture.work, "out", "delivery"),
		0o700,
	); err != nil {
		t.Fatalf("create message directory: %v", err)
	}
	writeTestFile(
		t,
		filepath.Join(fixture.work, "out", "delivery", "commit.md"),
		"Fail after exact rebase installation\n",
	)
	writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "feature\n")
	originalHead := runTestGit(t, fixture.work, "rev-parse", "HEAD")
	originalIndex := runTestGit(t, fixture.work, "write-tree")
	originalIndexContents, err := os.ReadFile(
		filepath.Join(fixture.work, ".git", "index"),
	)
	if err != nil {
		t.Fatalf("read original index: %v", err)
	}
	headLockRemovals := 0
	cleanupFailures := 0
	fixture.delivery.repository.removeLockFile = func(path string) error {
		removeErr := os.Remove(path)
		if removeErr == nil && filepath.Base(path) == "HEAD.lock" {
			headLockRemovals++
			if headLockRemovals == 2 {
				cleanupFailures++
				return fmt.Errorf("synthetic post-rebase HEAD-lock cleanup failure")
			}
		}
		return removeErr
	}

	_, err = fixture.delivery.prepare(context.Background(), prepareOptions{
		MessageFile: "out/delivery/commit.md",
		ReceiptFile: "out/delivery/prepare.json",
		Paths:       []string{"feature.txt"},
	})
	if err == nil || !strings.Contains(
		err.Error(),
		"synthetic post-rebase HEAD-lock cleanup failure",
	) {
		t.Fatalf("prepare() error = %v, want post-install cleanup failure", err)
	}
	if strings.Contains(err.Error(), "roll back failed preparation") {
		t.Fatalf("prepare() failed to recognize installed rebase: %v", err)
	}
	if cleanupFailures != 1 {
		t.Fatalf("injected cleanup failures = %d, want 1", cleanupFailures)
	}
	restoredIndexContents, err := os.ReadFile(
		filepath.Join(fixture.work, ".git", "index"),
	)
	if err != nil {
		t.Fatalf("read restored index: %v", err)
	}
	if !bytes.Equal(restoredIndexContents, originalIndexContents) {
		t.Fatal("rollback did not restore the exact original index bytes")
	}
	assertPreparationEntryState(
		t,
		fixture.work,
		originalHead,
		originalIndex,
		"?? feature.txt",
	)
	if _, err := os.Lstat(
		filepath.Join(fixture.work, "base-advance.txt"),
	); !os.IsNotExist(err) {
		t.Fatalf("rebased base file survived exact rollback: %v", err)
	}
}

func assertPreparationEntryState(
	t *testing.T,
	directory string,
	wantHead string,
	wantIndex string,
	wantStatus string,
) {
	t.Helper()
	if got := runTestGit(t, directory, "rev-parse", "HEAD"); got != wantHead {
		t.Fatalf("HEAD = %s, want restored %s", got, wantHead)
	}
	if got := runTestGit(t, directory, "write-tree"); got != wantIndex {
		t.Fatalf("index tree = %s, want restored %s", got, wantIndex)
	}
	if got := runTestGit(t, directory, "status", "--porcelain"); got != wantStatus {
		t.Fatalf("restored status = %q, want %q", got, wantStatus)
	}
}

func TestPublishRechecksPreparedScopeAfterFinalForgeRead(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	fixture.forge.pullRequestsCalls = 0
	fixture.forge.onPullRequestsCall = func(call int) {
		if call == 2 {
			writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "raced edit\n")
		}
	}

	_, err := fixture.delivery.publish(context.Background(), publishOptions{
		ValidatedHead: prepared.HeadOID,
		ReceiptFile:   "out/delivery/prepare.json",
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
		prepared := fixture.prepare(t)
		runTestGit(t, fixture.work, "update-index", "--assume-unchanged", "feature.txt")
		writeTestFile(t, filepath.Join(fixture.work, "feature.txt"), "hidden edit\n")
		_, err := fixture.delivery.publish(context.Background(), publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
		prepared := fixture.prepare(t)
		if _, err := fixture.delivery.publish(context.Background(), publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   "out/delivery/prepare.json",
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
	prepared := fixture.prepare(t)
	if _, err := fixture.delivery.publish(context.Background(), publishOptions{
		ValidatedHead: prepared.HeadOID,
		ReceiptFile:   "out/delivery/prepare.json",
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

type failPreparationFinalizationRunner struct {
	delegate  commandRunner
	candidate string
	armed     bool
	failed    bool
}

type cancelAfterIndexInstallRunner struct {
	delegate commandRunner
	cancel   context.CancelFunc
	armed    bool
	fired    bool
}

func (r *cancelAfterIndexInstallRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if r.armed && !r.fired && !hasEnvironmentPrefix(
		request.Env,
		"GIT_INDEX_FILE=",
	) && containsString(request.Args, "status") {
		r.fired = true
		r.cancel()
	}
	result, err := r.delegate.Run(ctx, request)
	if err == nil && hasEnvironmentPrefix(request.Env, "GIT_INDEX_FILE=") &&
		containsString(request.Args, "add") {
		r.armed = true
	}
	return result, err
}

type cancelAfterCandidateInstallRunner struct {
	delegate  commandRunner
	cancel    context.CancelFunc
	candidate string
	fired     bool
}

func (r *cancelAfterCandidateInstallRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	result, err := r.delegate.Run(ctx, request)
	if err != nil {
		return result, err
	}
	if containsString(request.Args, "commit-tree") {
		r.candidate = strings.TrimSpace(result.Stdout)
	}
	if !r.fired && isObjectID(r.candidate) &&
		containsString(request.Args, "for-each-ref") &&
		strings.Contains(result.Stdout, r.candidate) {
		r.fired = true
		r.cancel()
	}
	return result, nil
}

func hasEnvironmentPrefix(environment []string, prefix string) bool {
	for _, value := range environment {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func (r *failPreparationFinalizationRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	if r.armed && !r.failed && containsSequence(request.Args, "rev-list", "--count") {
		r.failed = true
		return syntheticCommandFailure(request, 73, "synthetic finalization failure")
	}
	result, err := r.delegate.Run(ctx, request)
	if err == nil && containsString(request.Args, "commit-tree") {
		r.candidate = strings.TrimSpace(result.Stdout)
	}
	if err == nil && isObjectID(r.candidate) &&
		containsString(request.Args, "for-each-ref") &&
		strings.Contains(result.Stdout, r.candidate) {
		r.armed = true
	}
	return result, err
}
