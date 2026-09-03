package main

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestSuggestedValidationCommands(t *testing.T) {
	t.Parallel()
	commands := suggestedValidationCommands([]string{
		"//tools/foo",
		"//tools/bar",
	})
	want := []string{
		"bazel_agent bazel test //tools/bar //tools/foo",
		"bazel_agent bazel build --config=lint //tools/bar //tools/foo",
	}
	if !reflect.DeepEqual(commands, want) {
		t.Fatalf("commands = %v, want %v", commands, want)
	}
}

func TestSuggestedValidationCommandsEmpty(t *testing.T) {
	t.Parallel()
	if commands := suggestedValidationCommands(nil); commands != nil {
		t.Fatalf("commands = %v, want nil", commands)
	}
}

func TestCommitAndPullRequestDisclaimers(t *testing.T) {
	t.Parallel()
	message, err := withCommitDisclaimer(
		"Add delivery automation\n\nExplain the aggregate change.\n",
	)
	if err != nil {
		t.Fatalf("withCommitDisclaimer() error = %v", err)
	}
	if !hasFinalLine(message, commitDisclaimer) {
		t.Fatalf("message does not end with commit disclaimer:\n%s", message)
	}
	body, err := pullRequestBody(
		"Explain the aggregate change.\n\n" + commitDisclaimer + "\n",
	)
	if err != nil {
		t.Fatalf("pullRequestBody() error = %v", err)
	}
	if !hasFinalLine(body, pullRequestDisclaimer) {
		t.Fatalf("body does not end with pull request disclaimer:\n%s", body)
	}
	if strings.Contains(body, commitDisclaimer) {
		t.Fatalf("body retained commit disclaimer:\n%s", body)
	}
}

func TestWithCommitDisclaimerAcceptsRegularConventions(t *testing.T) {
	t.Parallel()
	message, err := withCommitDisclaimer(
		"Add delivery automation\n\nExplain the aggregate change.\n\n" +
			"Goal-Ref: goal\nAttempt-ID: attempt\n",
	)
	if err != nil {
		t.Fatalf("withCommitDisclaimer() error = %v", err)
	}
	want := "Add delivery automation\n\nExplain the aggregate change.\n\n" +
		"Goal-Ref: goal\nAttempt-ID: attempt\n\n" + commitDisclaimer + "\n"
	if message != want {
		t.Fatalf("message =\n%s\nwant:\n%s", message, want)
	}
}

func TestWithCommitDisclaimerRejectsConventionalPrefix(t *testing.T) {
	t.Parallel()
	_, err := withCommitDisclaimer("feat: Add delivery automation")
	if err == nil || !strings.Contains(
		err.Error(),
		"artificial prefix",
	) {
		t.Fatalf(
			"withCommitDisclaimer() error = %v, want artificial prefix",
			err,
		)
	}
}

func TestWithCommitDisclaimerRejectsMissingSubjectSeparator(t *testing.T) {
	t.Parallel()
	_, err := withCommitDisclaimer("Add delivery automation\nBody")
	if err == nil || !strings.Contains(err.Error(), "blank line") {
		t.Fatalf(
			"withCommitDisclaimer() error = %v, want blank line",
			err,
		)
	}
}

func TestWithCommitDisclaimerRejectsUnpairedGoalTrailers(t *testing.T) {
	t.Parallel()
	_, err := withCommitDisclaimer(
		"Add delivery automation\n\nGoal-Ref: goal",
	)
	if err == nil || !strings.Contains(err.Error(), "Goal-Ref") {
		t.Fatalf(
			"withCommitDisclaimer() error = %v, want Goal-Ref pairing",
			err,
		)
	}
}

func TestWithCommitDisclaimerRejectsDisclaimerBeforeGoalTrailers(t *testing.T) {
	t.Parallel()
	_, err := withCommitDisclaimer(
		"Add delivery automation\n\n" +
			commitDisclaimer + "\n\n" +
			"Goal-Ref: goal\nAttempt-ID: attempt",
	)
	if err == nil || !strings.Contains(err.Error(), "final trailer") {
		t.Fatalf(
			"withCommitDisclaimer() error = %v, want final trailer",
			err,
		)
	}
}

func TestWithCommitDisclaimerRejectsWrongFinalDisclaimer(t *testing.T) {
	t.Parallel()
	_, err := withCommitDisclaimer(
		"Subject\n\nBody\n\n" + pullRequestDisclaimer,
	)
	if err == nil || !strings.Contains(err.Error(), "wrong disclaimer") {
		t.Fatalf("withCommitDisclaimer() error = %v, want wrong disclaimer", err)
	}
}

func TestParseRemoteRepository(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value string
		want  remoteRepository
	}{
		{
			name:  "scp",
			value: "git@github.com:alwaldend/src.git",
			want: remoteRepository{
				Host: "github.com", Owner: "alwaldend", Name: "src",
			},
		},
		{
			name:  "https",
			value: "https://github.example/alwaldend/src.git",
			want: remoteRepository{
				Host: "github.example", Owner: "alwaldend", Name: "src",
			},
		},
		{
			name:  "ssh",
			value: "ssh://git@github.com/alwaldend/src.git",
			want: remoteRepository{
				Host: "github.com", Owner: "alwaldend", Name: "src",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseRemoteRepository(test.value)
			if err != nil {
				t.Fatalf("parseRemoteRepository() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("parseRemoteRepository() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestParseRemoteRepositoryDoesNotEchoCredentialURL(t *testing.T) {
	t.Parallel()
	secret := "top-secret-token"
	_, err := parseRemoteRepository("https://user:" + secret + "@github.com/bad")
	if err == nil {
		t.Fatal("parseRemoteRepository() unexpectedly succeeded")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error leaked credential URL: %v", err)
	}
}

func TestParseRemoteRepositoryRejectsLocalTransport(t *testing.T) {
	t.Parallel()
	_, err := parseRemoteRepository("file://github.com/owner/repo.git")
	if err == nil || !strings.Contains(err.Error(), "unsupported transport") {
		t.Fatalf("parseRemoteRepository() error = %v, want transport refusal", err)
	}
}

func TestCanonicalSSHDeliveryEndpoints(t *testing.T) {
	t.Parallel()
	for _, endpoint := range []string{
		"git@github.com:owner/repository.git",
		"ssh://git@github.com/owner/repository.git",
	} {
		if err := requireCanonicalSSHEndpoint(endpoint); err != nil {
			t.Errorf("requireCanonicalSSHEndpoint(%q) error = %v", endpoint, err)
		}
	}
}

func TestCanonicalSSHDeliveryEndpointRejectsAmbiguousSyntax(t *testing.T) {
	t.Parallel()
	secret := "credential-that-must-not-appear"
	for _, endpoint := range []string{
		"https://user:" + secret + "@github.com/owner/repository.git",
		"-oProxyCommand=attacker@github.com:owner/repository.git",
		"git user@github.com:owner/repository.git",
		"git@@github.com:owner/repository.git",
		"git@-github.com:owner/repository.git",
		"ssh://-option@github.com/owner/repository.git",
		"ssh://git@-github.com/owner/repository.git",
		"ssh://github.com/owner/repository.git",
		"git@github.com:owner/repository name.git",
		"git@github.com:owner/repository;command.git",
		"git@github.com:owner/../repository.git",
	} {
		err := requireCanonicalSSHEndpoint(endpoint)
		if err == nil {
			t.Errorf("requireCanonicalSSHEndpoint(%q) unexpectedly succeeded", endpoint)
			continue
		}
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("endpoint refusal leaked credentials: %v", err)
		}
	}
}

func TestValidateRemoteNameRejectsConfigKeyAmbiguity(t *testing.T) {
	t.Parallel()
	for _, remote := range []string{
		"origin.extra",
		"origin/url",
		"origin\nurl",
		"-origin",
		"origin key",
	} {
		if err := validateRemoteName(remote); err == nil {
			t.Errorf("validateRemoteName(%q) unexpectedly succeeded", remote)
		}
	}
	for _, remote := range []string{"origin", "upstream-2", "review_branch"} {
		if err := validateRemoteName(remote); err != nil {
			t.Errorf("validateRemoteName(%q) error = %v", remote, err)
		}
	}
}

func TestExactPullRequestsIgnoresUnrelatedForkBeforeCrossRepoRefusal(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	all := []pullRequest{
		{
			ID: "fork", Number: 1,
			URL: "https://github.com/owner/repo/pull/1", State: "OPEN",
			BaseRefName: "master", BaseRefOID: testOID('a'),
			HeadRefName: "feature", HeadRefOID: testOID('c'),
			HeadRepositoryOwner: "fork", HeadRepositoryName: "repo",
			IsCrossRepository: true,
		},
		{
			ID: "exact", Number: 2, URL: "https://github.com/owner/repo/pull/2",
			HeadRefName: "feature", HeadRefOID: testOID('b'),
			HeadRepositoryOwner: "owner", HeadRepositoryName: "repo",
			BaseRefName: "master", BaseRefOID: testOID('a'), State: "OPEN",
		},
	}
	got, err := exactPullRequests(all, repository, "feature")
	if err != nil {
		t.Fatalf("exactPullRequests() error = %v", err)
	}
	if len(got) != 1 || got[0].Number != 2 {
		t.Fatalf("exactPullRequests() = %#v, want PR 2", got)
	}
}

func TestExactPullRequestsIgnoresOtherRepositoryFromSameOwner(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	got, err := exactPullRequests([]pullRequest{{
		ID: "other", Number: 1,
		URL: "https://github.com/owner/other/pull/1", State: "OPEN",
		BaseRefName: "master", BaseRefOID: testOID('a'),
		HeadRefName: "feature", HeadRefOID: testOID('b'),
		HeadRepositoryOwner: "owner", HeadRepositoryName: "other",
		IsCrossRepository: true,
	}}, repository, "feature")
	if err != nil {
		t.Fatalf("exactPullRequests() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("exactPullRequests() = %#v, want no exact PR", got)
	}
}

func TestSelectOpenPullRequestRefusesClosedOnly(t *testing.T) {
	t.Parallel()
	selected, refusals := selectOpenPullRequest([]pullRequest{{State: "MERGED"}})
	if selected != nil || len(refusals) != 1 {
		t.Fatalf(
			"selectOpenPullRequest() = %#v, %#v; want nil and one refusal",
			selected,
			refusals,
		)
	}
}

func TestMergeEnvironmentReplacesAndUnsets(t *testing.T) {
	t.Parallel()
	got := mergeEnvironment(
		[]string{"PATH=/bin", "GH_HOST=wrong", "GH_PAGER=less", "TOKEN=keep"},
		[]string{"GH_PAGER=cat", "NO_COLOR=1"},
		[]string{"GH_HOST"},
		nil,
	)
	want := []string{"PATH=/bin", "TOKEN=keep", "GH_PAGER=cat", "NO_COLOR=1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mergeEnvironment() = %q, want %q", got, want)
	}
}

func TestLimitedBufferFailsClosedOnTruncation(t *testing.T) {
	t.Parallel()
	buffer := limitedBuffer{limit: 4}
	if _, err := buffer.Write([]byte("abcdef")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if !buffer.truncated || buffer.String() != "abcd" {
		t.Fatalf(
			"limitedBuffer = %q, truncated %v; want abcd, true",
			buffer.String(),
			buffer.truncated,
		)
	}
}

func TestCommandErrorRedactsURLCredentials(t *testing.T) {
	t.Parallel()
	secret := "top-secret-token"
	err := (&commandError{
		Command: command{Name: "git"},
		Result: commandResult{
			Stderr:   "fatal: https://user:" + secret + "@github.com/owner/repo",
			ExitCode: 128,
		},
	}).Error()
	if strings.Contains(err, secret) || !strings.Contains(err, "<redacted>@github.com") {
		t.Fatalf("command error did not safely redact credentials: %q", err)
	}
}

func TestValidateTaskPathsRejectsRepositoryRoot(t *testing.T) {
	t.Parallel()
	for _, value := range []string{
		".", "..", "../outside", "/absolute", "out/task",
	} {
		if _, err := validateTaskPaths([]string{value}); err == nil {
			t.Errorf("validateTaskPaths(%q) unexpectedly succeeded", value)
		}
	}
}

func TestPreparePathFlagPreservesCommaInLiteralPath(t *testing.T) {
	t.Parallel()
	command := newPrepareCommand(
		context.Background(),
		&deliveryConfig{},
		func(string) string { return "" },
		io.Discard,
		nil,
	)
	if err := command.Flags().Parse([]string{
		"--message-file", "out/task/message.md",
		"--path", "owned/foo,bar",
	}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	paths, err := command.Flags().GetStringArray("path")
	if err != nil {
		t.Fatalf("GetStringArray() error = %v", err)
	}
	if !reflect.DeepEqual(paths, []string{"owned/foo,bar"}) {
		t.Fatalf("--path values = %q, want one literal comma path", paths)
	}
}

func TestPrepareReplaceRemoteFlagPreservesLiteralOID(t *testing.T) {
	t.Parallel()
	want := testOID('a')
	command := newPrepareCommand(
		context.Background(),
		&deliveryConfig{},
		func(string) string { return "" },
		io.Discard,
		nil,
	)
	if err := command.Flags().Parse([]string{
		"--message-file", "out/task/message.md",
		"--receipt-file", "out/task/prepare.json",
		"--replace-remote", want,
	}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	got, err := command.Flags().GetString("replace-remote")
	if err != nil {
		t.Fatalf("GetString() error = %v", err)
	}
	if got != want {
		t.Fatalf("--replace-remote = %q, want %q", got, want)
	}
}

func TestPrepareConsolidateFlagPreservesLiteralOID(t *testing.T) {
	t.Parallel()
	want := testOID('c')
	command := newPrepareCommand(
		context.Background(),
		&deliveryConfig{},
		func(string) string { return "" },
		io.Discard,
		nil,
	)
	if err := command.Flags().Parse([]string{
		"--message-file", "out/task/message.md",
		"--receipt-file", "out/task/prepare.json",
		"--consolidate", want,
	}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	got, err := command.Flags().GetString("consolidate")
	if err != nil {
		t.Fatalf("GetString() error = %v", err)
	}
	if got != want {
		t.Fatalf("--consolidate = %q, want %q", got, want)
	}
}

func TestRequireConsolidationEvidence(t *testing.T) {
	t.Parallel()
	parentOID := testOID('1')
	firstOID := testOID('2')
	headOID := testOID('3')
	base := inspection{
		LocalHeadOID:      headOID,
		UniqueCommitCount: 2,
		FeatureCommits: []featureCommit{
			{
				OID:             firstOID,
				Parents:         []string{parentOID},
				Title:           "First feature commit",
				AuthorName:      "Task Bot",
				AuthorEmail:     "task@example.com",
				CommitterName:   "Task Bot",
				CommitterEmail:  "task@example.com",
				SignatureStatus: "N",
				HasDisclaimer:   true,
			},
			{
				OID:             headOID,
				Parents:         []string{firstOID},
				AuthorName:      "Task Bot",
				AuthorEmail:     "task@example.com",
				CommitterName:   "Task Bot",
				CommitterEmail:  "task@example.com",
				SignatureStatus: "N",
			},
		},
	}
	evidence, err := (&delivery{}).requireConsolidationEvidence(
		context.Background(),
		&base,
		headOID,
		"Consolidated change\n\nDetails.\n\n"+commitDisclaimer+"\n",
	)
	if err != nil {
		t.Fatalf("requireConsolidationEvidence() error = %v", err)
	}
	if evidence.ParentOID != parentOID || evidence.AuthorOID != firstOID ||
		!reflect.DeepEqual(evidence.SignatureSource, []string{firstOID, headOID}) {
		t.Fatalf("consolidation evidence = %#v", evidence)
	}

	tests := []struct {
		name        string
		mutate      func(*inspection)
		expectedOID string
		want        string
	}{
		{
			name:        "stale exact head",
			expectedOID: testOID('4'),
			want:        "differs from the exact inspected local head",
		},
		{
			name: "missing ownership marker",
			mutate: func(report *inspection) {
				report.FeatureCommits[0].HasDisclaimer = false
			},
			expectedOID: headOID,
			want:        "lacks the required ownership marker; amend",
		},
		{
			name: "nonlinear chain",
			mutate: func(report *inspection) {
				report.FeatureCommits[1].Parents = []string{parentOID}
			},
			expectedOID: headOID,
			want:        "not one linear parent chain",
		},
		{
			name: "different identity",
			mutate: func(report *inspection) {
				report.FeatureCommits[1].CommitterEmail = "human@example.com"
			},
			expectedOID: headOID,
			want:        "different author or committer identity",
		},
		{
			name: "bad signature",
			mutate: func(report *inspection) {
				report.FeatureCommits[1].SignatureStatus = "B"
			},
			expectedOID: headOID,
			want:        "has a bad signature",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			report := base
			report.FeatureCommits = append(
				[]featureCommit(nil),
				base.FeatureCommits...,
			)
			for index := range report.FeatureCommits {
				report.FeatureCommits[index].Parents = append(
					[]string(nil),
					report.FeatureCommits[index].Parents...,
				)
			}
			if test.mutate != nil {
				test.mutate(&report)
			}
			_, err := (&delivery{}).requireConsolidationEvidence(
				context.Background(),
				&report,
				test.expectedOID,
				"Consolidated change\n\nDetails.\n\n"+commitDisclaimer+"\n",
			)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestRequireConsolidationEvidencePreservesRequestedPullRequest(t *testing.T) {
	t.Parallel()
	parentOID := testOID('1')
	firstOID := testOID('2')
	headOID := testOID('3')
	message := "Consolidated change\n\nCurrent aggregate description.\n\n" +
		commitDisclaimer + "\n"
	projection, err := messageProjection(message)
	if err != nil {
		t.Fatalf("messageProjection() error = %v", err)
	}
	report := inspection{
		LocalHeadOID:      headOID,
		UniqueCommitCount: 2,
		Base:              "master",
		Branch:            "feature",
		PullRequest: &pullRequest{
			Title:       projection.Title,
			Body:        projection.Body,
			BaseRefName: "master",
			HeadRefName: "feature",
		},
		FeatureCommits: []featureCommit{
			{
				OID:             firstOID,
				Parents:         []string{parentOID},
				AuthorName:      "Task Bot",
				AuthorEmail:     "task@example.com",
				CommitterName:   "Task Bot",
				CommitterEmail:  "task@example.com",
				SignatureStatus: "N",
				HasDisclaimer:   true,
			},
			{
				OID:             headOID,
				Parents:         []string{firstOID},
				AuthorName:      "Task Bot",
				AuthorEmail:     "task@example.com",
				CommitterName:   "Task Bot",
				CommitterEmail:  "task@example.com",
				SignatureStatus: "N",
			},
		},
	}
	if _, err := (&delivery{}).requireConsolidationEvidence(
		context.Background(),
		&report,
		headOID,
		message,
	); err != nil {
		t.Fatalf("requireConsolidationEvidence() error = %v", err)
	}
	report.PullRequest.Body = "Human-edited body.\n"
	if _, err := (&delivery{}).requireConsolidationEvidence(
		context.Background(),
		&report,
		headOID,
		message,
	); err == nil || !strings.Contains(err.Error(), "requested consolidation projection") {
		t.Fatalf("error = %v, want requested projection refusal", err)
	}
}

func TestEnsurePreparationRefusalsOnlyWaivesMultiCommitOwnership(t *testing.T) {
	t.Parallel()
	report := inspection{
		UniqueCommitCount: 2,
		Refusals:          []string{multiCommitRefusal(2)},
	}
	if err := ensurePreparationRefusals(&report, true); err != nil {
		t.Fatalf("ensurePreparationRefusals() error = %v", err)
	}
	report.Refusals = append(report.Refusals, "unrelated refusal")
	if err := ensurePreparationRefusals(&report, true); err == nil ||
		!strings.Contains(err.Error(), "unrelated refusal") {
		t.Fatalf("error = %v, want unrelated refusal", err)
	}
}

func TestRequireRemoteReplacementAuthorization(t *testing.T) {
	t.Parallel()
	remoteOID := testOID('a')
	otherOID := testOID('b')
	tests := []struct {
		name        string
		report      inspection
		expected    string
		wantErrPart string
	}{
		{name: "ordinary preparation"},
		{
			name:     "exact divergent remote",
			report:   inspection{RemoteHeadOID: remoteOID, RemoteHeadDiverged: true},
			expected: remoteOID,
		},
		{
			name:        "missing authorization",
			report:      inspection{RemoteHeadOID: remoteOID, RemoteHeadDiverged: true},
			wantErrPart: "requires --replace-remote=",
		},
		{
			name:        "mismatched authorization",
			report:      inspection{RemoteHeadOID: remoteOID, RemoteHeadDiverged: true},
			expected:    otherOID,
			wantErrPart: "differs from",
		},
		{
			name:        "unnecessary authorization",
			report:      inspection{RemoteHeadOID: remoteOID},
			expected:    remoteOID,
			wantErrPart: "unnecessary",
		},
		{
			name:        "absent remote",
			expected:    remoteOID,
			wantErrPart: "absent remote",
		},
		{
			name:        "malformed object ID",
			report:      inspection{RemoteHeadOID: remoteOID, RemoteHeadDiverged: true},
			expected:    "not-an-object-id",
			wantErrPart: "full Git object ID",
		},
		{
			name:        "inconsistent divergence report",
			report:      inspection{RemoteHeadDiverged: true},
			expected:    remoteOID,
			wantErrPart: "without an exact remote OID",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := requireRemoteReplacementAuthorization(&test.report, test.expected)
			if test.wantErrPart == "" {
				if err != nil {
					t.Fatalf("requireRemoteReplacementAuthorization() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErrPart) {
				t.Fatalf(
					"requireRemoteReplacementAuthorization() error = %v, want %q",
					err,
					test.wantErrPart,
				)
			}
		})
	}
}

func TestPushUsesExactAbsentLeaseForNewRef(t *testing.T) {
	t.Parallel()
	got := exactPushArguments(
		"repo-delivery-test",
		testOID('c'),
		"refs/heads/feature",
		"",
	)
	want := []string{
		"push",
		"--atomic",
		"--no-verify",
		"--no-follow-tags",
		"--no-push-option",
		"--recurse-submodules=no",
		"--force-with-lease=refs/heads/feature:",
		"--",
		"repo-delivery-test",
		testOID('c') + ":refs/heads/feature",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("exactPushArguments() = %q, want %q", got, want)
	}
}

func TestCommitHeaderHasSignatureRecognizesBothObjectFormats(t *testing.T) {
	for _, header := range []string{
		"tree " + strings.Repeat("a", 40) + "\ngpgsig signature",
		"tree " + strings.Repeat("a", 64) + "\ngpgsig-sha256 signature",
	} {
		if !commitHeaderHasSignature(header) {
			t.Fatalf("commitHeaderHasSignature(%q) = false", header)
		}
	}
	if commitHeaderHasSignature("tree " + strings.Repeat("a", 40) + "\ngpgsignature nope") {
		t.Fatal("commitHeaderHasSignature accepted a lookalike header")
	}
}
