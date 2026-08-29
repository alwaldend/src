package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
)

func receiptTestPullRequest() pullRequest {
	return pullRequest{
		ID:                  "PR_test",
		Number:              17,
		URL:                 "https://github.com/owner/repo/pull/17",
		State:               "OPEN",
		Title:               "Prior title",
		Body:                "Prior body\n",
		AuthorLogin:         "delivery-bot",
		BaseRefName:         "master",
		BaseRefOID:          testOID('a'),
		HeadRefName:         "feature",
		HeadRefOID:          testOID('d'),
		HeadRepositoryOwner: "owner",
		HeadRepositoryName:  "repo",
	}
}

func validPreparationReceipt(
	t *testing.T,
	pullRequest *pullRequest,
) preparationReceipt {
	t.Helper()
	expectation, err := newPullRequestExpectation(pullRequest)
	if err != nil {
		t.Fatalf("newPullRequestExpectation() error = %v", err)
	}
	remoteHead := refExpectation{Ref: "refs/heads/feature"}
	if pullRequest != nil {
		remoteHead.Present = true
		remoteHead.OID = pullRequest.HeadRefOID
	}
	return preparationReceipt{
		Schema:                preparationReceiptSchema,
		RevisionNonce:         strings.Repeat("1", receiptRevisionBytes*2),
		RepositoryFingerprint: "sha256:" + strings.Repeat("a", 64),
		RemoteName:            "origin",
		FetchEndpointDigest:   "sha256:" + strings.Repeat("b", 64),
		PushEndpointDigest:    "sha256:" + strings.Repeat("c", 64),
		RemoteRepository: remoteRepository{
			Host: "github.com", Owner: "owner", Name: "repo",
		},
		Forge:               "github",
		BaseRef:             "refs/heads/master",
		BaseOID:             testOID('a'),
		HeadRef:             "refs/heads/feature",
		PreparedHeadOID:     testOID('b'),
		PreparedTreeOID:     testOID('c'),
		ExpectedRemoteHead:  remoteHead,
		ExpectedPullRequest: expectation,
		Scope: aggregateScope{
			Mode:            scopeModePaths,
			AuthorizedPaths: []string{"owned"},
			AggregatePaths:  []string{"owned/file.txt"},
		},
	}
}

func TestEndpointDigestDoesNotSerializeCredentials(t *testing.T) {
	t.Parallel()
	secret := "top-secret-token"
	endpoint := "https://user:" + secret + "@github.com/owner/repo.git"
	contents, err := json.Marshal(map[string]string{
		"endpoint_digest": endpointDigest("fetch", endpoint),
	})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if strings.Contains(string(contents), secret) ||
		strings.Contains(string(contents), endpoint) {
		t.Fatalf("serialized endpoint digest leaked credentials: %s", contents)
	}
}

func TestParsePorcelainV2Z(t *testing.T) {
	t.Parallel()
	oidA := strings.Repeat("a", 40)
	oidB := strings.Repeat("b", 40)
	oidC := strings.Repeat("c", 40)
	value := strings.Join([]string{
		"1 M. N... 100644 100644 100644 " + oidA + " " + oidB + " staged.txt",
		"1 .M N... 100644 100644 100644 " + oidA + " " + oidB + " unstaged name.txt",
		"1 MM N... 100644 100644 100644 " + oidA + " " + oidB + " both.txt",
		"? nested/untracked.txt",
		"u UU N... 100644 100644 100644 100644 " + oidA + " " + oidB + " " + oidC + " conflict.txt",
	}, "\x00") + "\x00"
	got, err := parsePorcelainV2Z(value)
	if err != nil {
		t.Fatalf("parsePorcelainV2Z() error = %v", err)
	}
	want := repositoryStatus{
		Staged:    []string{"both.txt", "conflict.txt", "staged.txt"},
		Unstaged:  []string{"both.txt", "conflict.txt", "unstaged name.txt"},
		Untracked: []string{"nested/untracked.txt"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parsePorcelainV2Z() = %#v, want %#v", got, want)
	}
}

func TestParsePorcelainV2ZFailsClosed(t *testing.T) {
	t.Parallel()
	for _, value := range []string{
		"? missing-terminal-nul",
		"? ok\x00\x00? bad\x00",
		"2 R. N... fields\x00old\x00",
		"# branch.head feature\x00",
		"! ignored\x00",
		"1 XX N... 1 2 3 a b bad.txt\x00",
		"1 M. N... 10064x 100644 100644 " + strings.Repeat("a", 40) + " " + strings.Repeat("b", 40) + " bad-mode.txt\x00",
		"1 M. N... 100644 100644 100644 short " + strings.Repeat("b", 40) + " bad-oid.txt\x00",
	} {
		if _, err := parsePorcelainV2Z(value); err == nil {
			t.Errorf("parsePorcelainV2Z(%q) unexpectedly succeeded", value)
		}
	}
}

func TestAggregateScopeValidation(t *testing.T) {
	t.Parallel()
	scope, err := newAggregateScope(
		scopeModePaths,
		[]string{"owned"},
		[]string{"owned/a.txt", "owned/b.txt"},
	)
	if err != nil {
		t.Fatalf("newAggregateScope() error = %v", err)
	}
	if err := scope.validate(); err != nil {
		t.Fatalf("scope.validate() error = %v", err)
	}
	for _, invalid := range []aggregateScope{
		{
			Mode:            scopeModePaths,
			AuthorizedPaths: []string{"owned"},
			AggregatePaths:  []string{"outside.txt"},
		},
		{
			Mode:            scopeModeUseIndex,
			AuthorizedPaths: []string{"owned"},
			AggregatePaths:  []string{"owned/a.txt"},
		},
		{
			Mode:            scopeModePaths,
			AuthorizedPaths: []string{"owned", "owned"},
			AggregatePaths:  []string{"owned/a.txt"},
		},
	} {
		if err := invalid.validate(); err == nil {
			t.Errorf("aggregateScope.validate(%#v) unexpectedly succeeded", invalid)
		}
	}
}

func TestPreparationReceiptValidation(t *testing.T) {
	t.Parallel()
	receipt := validPreparationReceipt(t, nil)
	if err := receipt.validate(); err != nil {
		t.Fatalf("receipt.validate() error = %v", err)
	}
	contents, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if _, err := decodePreparationReceipt(contents); err != nil {
		t.Fatalf("decodePreparationReceipt(absent PR) error = %v", err)
	}
	receipt.Schema = "repo_delivery/preparation/v2"
	if err := receipt.validate(); err == nil {
		t.Fatal("receipt.validate() accepted an unknown version")
	}
	receipt = validPreparationReceipt(t, nil)
	receipt.ExpectedPullRequest.ID = "PR_unexpected"
	if err := receipt.validate(); err == nil {
		t.Fatal("receipt.validate() accepted metadata on explicit PR absence")
	}
}

func TestPullRequestExpectationRejectsIdentityReplacement(t *testing.T) {
	t.Parallel()
	prior := receiptTestPullRequest()
	receipt := validPreparationReceipt(t, &prior)
	desired := commitProjection{Title: "Desired title", Body: "Desired body"}

	tests := []struct {
		name       string
		mutate     func(*pullRequest)
		remoteHead string
		want       receiptPullRequestState
		wantError  bool
	}{
		{
			name:       "prior",
			remoteHead: prior.HeadRefOID,
			want:       receiptPullRequestPrior,
		},
		{
			name: "head pushed before metadata",
			mutate: func(value *pullRequest) {
				value.HeadRefOID = receipt.PreparedHeadOID
			},
			remoteHead: receipt.PreparedHeadOID,
			want:       receiptPullRequestHeadPushedWithPriorProjection,
		},
		{
			name: "fully desired",
			mutate: func(value *pullRequest) {
				value.HeadRefOID = receipt.PreparedHeadOID
				value.Title = desired.Title
				value.Body = desired.Body
			},
			remoteHead: receipt.PreparedHeadOID,
			want:       receiptPullRequestDesired,
		},
		{
			name: "replacement id",
			mutate: func(value *pullRequest) {
				value.ID = "PR_replacement"
				value.HeadRefOID = receipt.PreparedHeadOID
				value.Title = desired.Title
				value.Body = desired.Body
			},
			remoteHead: receipt.PreparedHeadOID,
			wantError:  true,
		},
		{
			name: "replacement number",
			mutate: func(value *pullRequest) {
				value.Number++
				value.HeadRefOID = receipt.PreparedHeadOID
				value.Title = desired.Title
				value.Body = desired.Body
			},
			remoteHead: receipt.PreparedHeadOID,
			wantError:  true,
		},
		{
			name: "closed original",
			mutate: func(value *pullRequest) {
				value.State = "CLOSED"
			},
			remoteHead: prior.HeadRefOID,
			wantError:  true,
		},
		{
			name: "human metadata edit",
			mutate: func(value *pullRequest) {
				value.HeadRefOID = receipt.PreparedHeadOID
				value.Body = "human edit"
			},
			remoteHead: receipt.PreparedHeadOID,
			wantError:  true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			current := prior
			if test.mutate != nil {
				test.mutate(&current)
			}
			got, err := receipt.pullRequestState(
				&current,
				test.remoteHead,
				desired,
			)
			if test.wantError {
				if err == nil {
					t.Fatalf("pullRequestState() = %q, want error", got)
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf(
					"pullRequestState() = %q, %v; want %q, nil",
					got,
					err,
					test.want,
				)
			}
		})
	}
	if _, err := receipt.pullRequestState(
		nil,
		prior.HeadRefOID,
		desired,
	); err == nil {
		t.Fatal("pullRequestState() accepted disappearance of the expected PR")
	}
}

func TestAbsentPullRequestExpectationAllowsOnlyExactRecovery(t *testing.T) {
	t.Parallel()
	receipt := validPreparationReceipt(t, nil)
	receipt.ExpectedRemoteHead = refExpectation{
		Ref:     receipt.HeadRef,
		Present: true,
		OID:     testOID('d'),
	}
	receipt.ExpectedPullRequest.PriorHeadOID = receipt.ExpectedRemoteHead.OID
	desired := commitProjection{Title: "Desired title", Body: "Desired body"}
	if got, err := receipt.pullRequestState(
		nil,
		receipt.ExpectedRemoteHead.OID,
		desired,
	); err != nil || got != receiptPullRequestAbsent {
		t.Fatalf("prior absence = %q, %v", got, err)
	}
	if got, err := receipt.pullRequestState(
		nil,
		receipt.PreparedHeadOID,
		desired,
	); err != nil || got != receiptPullRequestHeadPushed {
		t.Fatalf("pushed absence = %q, %v", got, err)
	}
	created := receiptTestPullRequest()
	created.ID = "PR_created"
	created.Number = 18
	created.URL = "https://github.com/owner/repo/pull/18"
	created.HeadRefOID = receipt.PreparedHeadOID
	created.Title = desired.Title
	created.Body = desired.Body
	if got, err := receipt.pullRequestState(
		&created,
		receipt.PreparedHeadOID,
		desired,
	); err != nil || got != receiptPullRequestCreatedDesired {
		t.Fatalf("created recovery = %q, %v", got, err)
	}
	priorHeadCreation := created
	priorHeadCreation.HeadRefOID = receipt.ExpectedRemoteHead.OID
	if _, err := receipt.pullRequestState(
		&priorHeadCreation,
		receipt.ExpectedRemoteHead.OID,
		desired,
	); err == nil {
		t.Fatal("pullRequestState() accepted a new PR before the prepared push")
	}
	for name, mutate := range map[string]func(*pullRequest){
		"wrong projection": func(value *pullRequest) { value.Body = "unexpected" },
		"wrong head":       func(value *pullRequest) { value.HeadRefOID = testOID('e') },
		"draft":            func(value *pullRequest) { value.IsDraft = true },
		"wrong owner": func(value *pullRequest) {
			value.HeadRepositoryOwner = "somebody-else"
		},
	} {
		name := name
		mutate := mutate
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			candidate := created
			mutate(&candidate)
			if _, err := receipt.pullRequestState(
				&candidate,
				receipt.PreparedHeadOID,
				desired,
			); err == nil {
				t.Fatal("pullRequestState() accepted non-exact recovery")
			}
		})
	}
}

func TestCreatedPullRequestBindingRejectsLaterReplacement(t *testing.T) {
	t.Parallel()
	receipt := validPreparationReceipt(t, nil)
	receipt.ExpectedRemoteHead = refExpectation{
		Ref:     receipt.HeadRef,
		Present: true,
		OID:     testOID('d'),
	}
	receipt.ExpectedPullRequest.PriorHeadOID = receipt.ExpectedRemoteHead.OID
	desired := commitProjection{Title: "Desired title", Body: "Desired body"}
	created := receiptTestPullRequest()
	created.ID = "PR_created"
	created.Number = 18
	created.URL = "https://github.com/owner/repo/pull/18"
	created.HeadRefOID = receipt.PreparedHeadOID
	created.Title = desired.Title
	created.Body = desired.Body

	bound, err := receipt.bindCreatedPullRequest(
		&created,
		receipt.PreparedHeadOID,
		desired,
	)
	if err != nil {
		t.Fatalf("bindCreatedPullRequest() error = %v", err)
	}
	if !bound.ExpectedPullRequest.Present ||
		bound.ExpectedPullRequest.ID != created.ID {
		t.Fatalf("bound expectation = %#v", bound.ExpectedPullRequest)
	}
	if bound.RevisionNonce == receipt.RevisionNonce {
		t.Fatal("bound receipt reused the prior revision nonce")
	}
	if got, err := bound.pullRequestState(
		&created,
		receipt.PreparedHeadOID,
		desired,
	); err != nil || got != receiptPullRequestDesired {
		t.Fatalf("bound created state = %q, %v", got, err)
	}

	derived := bound
	derived.BaseOID = testOID('e')
	derived.PreparedHeadOID = testOID('f')
	derived.PreparedTreeOID = testOID('1')
	derived.ExpectedRemoteHead = refExpectation{
		Ref:     derived.HeadRef,
		Present: true,
		OID:     created.HeadRefOID,
	}
	if _, err := derived.pullRequestState(
		&created,
		created.HeadRefOID,
		desired,
	); err != nil {
		t.Fatalf("derived bound state error = %v", err)
	}
	replacement := created
	replacement.ID = "PR_replacement"
	replacement.Number = 19
	replacement.URL = "https://github.com/owner/repo/pull/19"
	if _, err := derived.pullRequestState(
		&replacement,
		replacement.HeadRefOID,
		desired,
	); err == nil {
		t.Fatal("derived bound receipt accepted a replacement pull request")
	}
}

func TestDerivedReceiptPreservesPartialPullRequestRecovery(t *testing.T) {
	t.Parallel()
	prior := receiptTestPullRequest()
	original := validPreparationReceipt(t, &prior)
	desired := commitProjection{Title: "Desired title", Body: "Desired body"}
	partialHead := original.PreparedHeadOID
	derived := original
	derived.BaseOID = testOID('e')
	derived.PreparedHeadOID = testOID('f')
	derived.PreparedTreeOID = testOID('1')
	derived.ExpectedRemoteHead.OID = partialHead

	current := prior
	current.HeadRefOID = partialHead
	if got, err := derived.pullRequestState(
		&current,
		partialHead,
		desired,
	); err != nil || got != receiptPullRequestPrior {
		t.Fatalf("derived prior projection = %q, %v", got, err)
	}
	current.Title = desired.Title
	current.Body = desired.Body
	if got, err := derived.pullRequestState(
		&current,
		partialHead,
		desired,
	); err != nil || got != receiptPullRequestDesired {
		t.Fatalf("derived desired projection = %q, %v", got, err)
	}

	absent := validPreparationReceipt(t, nil)
	absent.ExpectedPullRequest.PriorHeadOID = testOID('d')
	absent.ExpectedRemoteHead = refExpectation{
		Ref:     absent.HeadRef,
		Present: true,
		OID:     absent.PreparedHeadOID,
	}
	absent.PreparedHeadOID = testOID('f')
	absent.PreparedTreeOID = testOID('1')
	created := receiptTestPullRequest()
	created.ID = "PR_created"
	created.Number = 18
	created.URL = "https://github.com/owner/repo/pull/18"
	created.HeadRefOID = absent.ExpectedRemoteHead.OID
	created.Title = desired.Title
	created.Body = desired.Body
	if got, err := absent.pullRequestState(
		&created,
		absent.ExpectedRemoteHead.OID,
		desired,
	); err != nil || got != receiptPullRequestCreatedDesired {
		t.Fatalf("derived created projection = %q, %v", got, err)
	}
}

func TestCanonicalPreparationReceiptJSON(t *testing.T) {
	t.Parallel()
	prior := receiptTestPullRequest()
	receipt := validPreparationReceipt(t, &prior)
	contents, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent() error = %v", err)
	}
	if _, err := decodePreparationReceipt(contents); err != nil {
		t.Fatalf("decodePreparationReceipt(valid) error = %v", err)
	}
	valid := string(contents)
	tests := map[string]string{
		"duplicate top-level": strings.Replace(
			valid,
			`"schema": "repo_delivery/preparation/v1",`,
			`"schema": "repo_delivery/preparation/v1", "schema": "repo_delivery/preparation/v1",`,
			1,
		),
		"case variant": strings.Replace(valid, `"schema":`, `"Schema":`, 1),
		"duplicate repository key": strings.Replace(
			valid,
			`"host": "github.com",`,
			`"host": "github.com", "host": "github.com",`,
			1,
		),
		"case variant remote expectation": strings.Replace(
			valid,
			`"expected_remote_head": {`,
			`"Expected_remote_head": {`,
			1,
		),
		"duplicate pull request key": strings.Replace(
			valid,
			`"identity_sha256": "`+receipt.ExpectedPullRequest.IdentityDigest+`",`,
			`"identity_sha256": "`+receipt.ExpectedPullRequest.IdentityDigest+`", "identity_sha256": "`+receipt.ExpectedPullRequest.IdentityDigest+`",`,
			1,
		),
		"duplicate scope key": strings.Replace(
			valid,
			`"mode": "paths",`,
			`"mode": "paths", "mode": "paths",`,
			1,
		),
		"null boolean": strings.Replace(
			valid,
			`"present": true`,
			`"present": null`,
			1,
		),
		"string number": strings.Replace(
			valid,
			`"number": 17`,
			`"number": "17"`,
			1,
		),
		"missing canonical key": strings.Replace(
			valid,
			`    "oid": "`+receipt.ExpectedRemoteHead.OID+`"`+"\n",
			"",
			1,
		),
		"multiple values": valid + "\n{}",
	}
	for name, malformed := range tests {
		name := name
		malformed := malformed
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if malformed == valid {
				t.Fatal("test mutation did not change the receipt")
			}
			if _, err := decodePreparationReceipt([]byte(malformed)); err == nil {
				t.Fatal("decodePreparationReceipt() accepted non-canonical JSON")
			}
		})
	}
}

func TestValidateTaskOutputPath(t *testing.T) {
	t.Parallel()
	for _, value := range []string{
		"out/task/receipt.json",
		"out/task/nested/reply.json",
	} {
		if err := validateTaskOutputPath(value); err != nil {
			t.Errorf("validateTaskOutputPath(%q) error = %v", value, err)
		}
	}
	for _, value := range []string{
		"",
		"out/task",
		"out//receipt.json",
		"out/task/../receipt.json",
		"./out/task/receipt.json",
		"elsewhere/task/receipt.json",
		"/out/task/receipt.json",
		"out/task\\receipt.json",
		"out/task/receipt.json\x00suffix",
	} {
		if err := validateTaskOutputPath(value); err == nil {
			t.Errorf("validateTaskOutputPath(%q) unexpectedly succeeded", value)
		}
	}
}

func TestInstallReceiptAtomicallySyncsDirectoryAfterRename(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	temporaryPath := filepath.Join(directory, ".receipt.tmp")
	absolute := filepath.Join(directory, "receipt.json")
	contents := []byte("receipt contents")
	if err := os.WriteFile(temporaryPath, contents, 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if err := os.WriteFile(absolute, []byte("prior receipt"), 0o600); err != nil {
		t.Fatalf("os.WriteFile(existing receipt) error = %v", err)
	}

	var syncedDirectory string
	installed, err := installReceiptAtomically(
		temporaryPath,
		absolute,
		func(path string) error {
			syncedDirectory = path
			installedContents, readErr := os.ReadFile(absolute)
			if readErr != nil {
				t.Fatalf("receipt was not installed before directory sync: %v", readErr)
			}
			if !reflect.DeepEqual(installedContents, contents) {
				t.Fatalf(
					"installed contents = %q, want %q",
					installedContents,
					contents,
				)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("installReceiptAtomically() error = %v", err)
	}
	if !installed {
		t.Fatal("installReceiptAtomically() did not report the rename")
	}
	if syncedDirectory != directory {
		t.Fatalf("synced directory = %q, want %q", syncedDirectory, directory)
	}
	if _, err := os.Lstat(temporaryPath); !os.IsNotExist(err) {
		t.Fatalf("temporary receipt still exists after rename: %v", err)
	}
}

func TestInstallReceiptAtomicallyReportsUnknownOutcomeAfterSyncFailure(
	t *testing.T,
) {
	t.Parallel()
	directory := t.TempDir()
	temporaryPath := filepath.Join(directory, ".receipt.tmp")
	absolute := filepath.Join(directory, "receipt.json")
	if err := os.WriteFile(temporaryPath, []byte("receipt"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	syncError := errors.New("forced directory sync failure")

	installed, err := installReceiptAtomically(
		temporaryPath,
		absolute,
		func(string) error { return syncError },
	)
	if !installed {
		t.Fatal("installReceiptAtomically() did not report the completed rename")
	}
	if !errors.Is(err, syncError) {
		t.Fatalf("installReceiptAtomically() error = %v, want %v", err, syncError)
	}
	for _, guidance := range []string{
		"outcome unknown",
		"atomic rename succeeded",
		"re-inspect the receipt and repository state before retrying",
	} {
		if !strings.Contains(err.Error(), guidance) {
			t.Errorf("error %q does not contain %q", err, guidance)
		}
	}
	if _, statErr := os.Lstat(absolute); statErr != nil {
		t.Fatalf("renamed receipt is absent after sync failure: %v", statErr)
	}
	installedContents, readErr := os.ReadFile(absolute)
	if readErr != nil || string(installedContents) != "receipt" {
		t.Fatalf(
			"read receipt after lost acknowledgement = %q, %v",
			installedContents,
			readErr,
		)
	}
	if _, statErr := os.Lstat(temporaryPath); !os.IsNotExist(statErr) {
		t.Fatalf("temporary path exists after completed rename: %v", statErr)
	}
}

func TestReceiptStableGuardRejectsAlreadyStaleVersion(t *testing.T) {
	directory := t.TempDir()
	absolute := filepath.Join(directory, "prepare.json")
	initial := validPreparationReceipt(t, nil)
	initialContents, err := encodePreparationReceipt(initial)
	if err != nil {
		t.Fatalf("encode initial receipt: %v", err)
	}
	if err := os.WriteFile(absolute, initialContents, 0o600); err != nil {
		t.Fatalf("write initial receipt: %v", err)
	}
	staleVersion, err := captureReceiptFileVersion(absolute)
	if err != nil {
		t.Fatalf("capture stale receipt version: %v", err)
	}
	winnerReceipt := initial
	winnerReceipt.RevisionNonce, err = newReceiptRevisionNonce()
	if err != nil {
		t.Fatalf("generate winner revision: %v", err)
	}
	winnerContents, err := encodePreparationReceipt(winnerReceipt)
	if err != nil {
		t.Fatalf("encode winner receipt: %v", err)
	}
	winnerTemporary := filepath.Join(directory, ".winner.tmp")
	if err := os.WriteFile(winnerTemporary, winnerContents, 0o600); err != nil {
		t.Fatalf("write winner temporary receipt: %v", err)
	}
	installed, err := installReceiptAfterStableComparison(
		winnerTemporary,
		absolute,
		staleVersion,
		func(string) error { return nil },
	)
	if err != nil || !installed {
		t.Fatalf("winning receipt install = %v, %v; want installed", installed, err)
	}

	staleReceipt := initial
	staleReceipt.RevisionNonce, err = newReceiptRevisionNonce()
	if err != nil {
		t.Fatalf("generate stale-writer revision: %v", err)
	}
	staleContents, err := encodePreparationReceipt(staleReceipt)
	if err != nil {
		t.Fatalf("encode stale-writer receipt: %v", err)
	}
	staleTemporary := filepath.Join(directory, ".stale.tmp")
	if err := os.WriteFile(staleTemporary, staleContents, 0o600); err != nil {
		t.Fatalf("write stale temporary receipt: %v", err)
	}
	installed, err = installReceiptAfterStableComparison(
		staleTemporary,
		absolute,
		staleVersion,
		func(string) error { return nil },
	)
	if installed || err == nil || !strings.Contains(err.Error(), "exact byte version") {
		t.Fatalf(
			"stale receipt install = %v, %v; want exact-version refusal",
			installed,
			err,
		)
	}
	installedContents, err := readStableReceiptFile(absolute)
	if err != nil {
		t.Fatalf("read installed receipt: %v", err)
	}
	installedReceipt, err := decodePreparationReceipt(installedContents)
	if err != nil {
		t.Fatalf("decode installed receipt: %v", err)
	}
	if installedReceipt.RevisionNonce != winnerReceipt.RevisionNonce {
		t.Fatalf(
			"installed revision = %q, want winner %q",
			installedReceipt.RevisionNonce,
			winnerReceipt.RevisionNonce,
		)
	}
	if _, err := os.Lstat(staleTemporary); err != nil {
		t.Fatalf("refused stale temporary receipt was unexpectedly removed: %v", err)
	}
}

func TestReceiptLockSerializesOverlappingTransition(t *testing.T) {
	directory := t.TempDir()
	absolute := filepath.Join(directory, "prepare.json")
	lockPath := absolute + ".lock"
	initial := validPreparationReceipt(t, nil)
	initialContents, err := encodePreparationReceipt(initial)
	if err != nil {
		t.Fatalf("encode initial receipt: %v", err)
	}
	if err := os.WriteFile(absolute, initialContents, 0o600); err != nil {
		t.Fatalf("write initial receipt: %v", err)
	}

	first, err := openReceiptLockFile(lockPath)
	if err != nil {
		t.Fatalf("open first receipt lock: %v", err)
	}
	if err := lockReceiptFile(context.Background(), first); err != nil {
		t.Fatalf("acquire first receipt lock: %v", err)
	}
	type lockResult struct {
		file *os.File
		err  error
	}
	secondStarted := make(chan struct{})
	secondResult := make(chan lockResult, 1)
	waitContext, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()
	go func() {
		second, openErr := openReceiptLockFile(lockPath)
		close(secondStarted)
		if openErr == nil {
			openErr = lockReceiptFile(waitContext, second)
		}
		secondResult <- lockResult{file: second, err: openErr}
	}()
	<-secondStarted
	select {
	case result := <-secondResult:
		if result.file != nil {
			_ = result.file.Close()
		}
		t.Fatalf(
			"overlapping receipt writer completed while lock was held: %v",
			result.err,
		)
	case <-time.After(50 * time.Millisecond):
	}

	expected, err := captureReceiptFileVersion(absolute)
	if err != nil {
		t.Fatalf("capture locked receipt version: %v", err)
	}
	updated := initial
	updated.RevisionNonce, err = newReceiptRevisionNonce()
	if err != nil {
		t.Fatalf("generate updated revision: %v", err)
	}
	updatedContents, err := encodePreparationReceipt(updated)
	if err != nil {
		t.Fatalf("encode updated receipt: %v", err)
	}
	temporary := filepath.Join(directory, ".updated.tmp")
	if err := os.WriteFile(temporary, updatedContents, 0o600); err != nil {
		t.Fatalf("write updated temporary receipt: %v", err)
	}
	installed, err := installReceiptAfterStableComparison(
		temporary,
		absolute,
		expected,
		func(string) error { return nil },
	)
	if err != nil || !installed {
		t.Fatalf("install locked receipt = %v, %v", installed, err)
	}
	if err := syscall.Flock(int(first.Fd()), syscall.LOCK_UN); err != nil {
		t.Fatalf("release first receipt lock: %v", err)
	}
	if err := first.Close(); err != nil {
		t.Fatalf("close first receipt lock: %v", err)
	}

	var second lockResult
	select {
	case second = <-secondResult:
	case <-waitContext.Done():
		t.Fatalf("second receipt writer did not acquire released lock: %v", waitContext.Err())
	}
	if second.err != nil || second.file == nil {
		t.Fatalf("second receipt lock = %#v, %v", second.file, second.err)
	}
	observed, err := captureReceiptFileVersion(absolute)
	if err != nil {
		t.Fatalf("capture receipt after overlapping transition: %v", err)
	}
	if !observed.equal(newReceiptFileVersion(updatedContents)) {
		t.Fatal("second writer did not observe the completed locked transition")
	}
	if err := syscall.Flock(int(second.file.Fd()), syscall.LOCK_UN); err != nil {
		t.Fatalf("release second receipt lock: %v", err)
	}
	if err := second.file.Close(); err != nil {
		t.Fatalf("close second receipt lock: %v", err)
	}
}

func TestReceiptComparisonDoesNotExcludeLockIgnoringWriter(t *testing.T) {
	directory := t.TempDir()
	absolute := filepath.Join(directory, "prepare.json")
	initial := validPreparationReceipt(t, nil)
	initialContents, err := encodePreparationReceipt(initial)
	if err != nil {
		t.Fatalf("encode initial receipt: %v", err)
	}
	if err := os.WriteFile(absolute, initialContents, 0o600); err != nil {
		t.Fatalf("write initial receipt: %v", err)
	}
	expected, err := captureReceiptFileVersion(absolute)
	if err != nil {
		t.Fatalf("capture expected receipt: %v", err)
	}
	current, err := captureReceiptFileVersion(absolute)
	if err != nil {
		t.Fatalf("capture current receipt: %v", err)
	}
	if err := requireReceiptFileVersion(expected, current); err != nil {
		t.Fatalf("stable pre-rename comparison: %v", err)
	}

	lockIgnoring := initial
	lockIgnoring.RevisionNonce, err = newReceiptRevisionNonce()
	if err != nil {
		t.Fatalf("generate lock-ignoring revision: %v", err)
	}
	lockIgnoringContents, err := encodePreparationReceipt(lockIgnoring)
	if err != nil {
		t.Fatalf("encode lock-ignoring receipt: %v", err)
	}
	if err := os.WriteFile(absolute, lockIgnoringContents, 0o600); err != nil {
		t.Fatalf("write after stable comparison: %v", err)
	}

	intended := initial
	intended.RevisionNonce, err = newReceiptRevisionNonce()
	if err != nil {
		t.Fatalf("generate intended revision: %v", err)
	}
	intendedContents, err := encodePreparationReceipt(intended)
	if err != nil {
		t.Fatalf("encode intended receipt: %v", err)
	}
	temporary := filepath.Join(directory, ".intended.tmp")
	if err := os.WriteFile(temporary, intendedContents, 0o600); err != nil {
		t.Fatalf("write intended temporary receipt: %v", err)
	}
	// This deliberately models a writer that ignores the advisory lock and
	// lands after the stable comparison but before rename. There is no portable
	// pathname-content CAS joining these two filesystem operations.
	installed, err := installReceiptAtomically(
		temporary,
		absolute,
		func(string) error { return nil },
	)
	if err != nil || !installed {
		t.Fatalf("install after lock-ignoring write = %v, %v", installed, err)
	}
	finalContents, err := readStableReceiptFile(absolute)
	if err != nil {
		t.Fatalf("read final receipt: %v", err)
	}
	if !bytes.Equal(finalContents, intendedContents) {
		t.Fatal("test did not expose the documented non-cooperating-writer window")
	}
}

func TestReceiptLockSerializesAndReleasesAfterCancellation(t *testing.T) {
	lockPath := filepath.Join(t.TempDir(), "prepare.json.lock")
	first, err := openReceiptLockFile(lockPath)
	if err != nil {
		t.Fatalf("open first receipt lock: %v", err)
	}
	if err := lockReceiptFile(context.Background(), first); err != nil {
		t.Fatalf("acquire first receipt lock: %v", err)
	}
	second, err := openReceiptLockFile(lockPath)
	if err != nil {
		t.Fatalf("open contending receipt lock: %v", err)
	}

	waitContext, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()
	err = lockReceiptFile(waitContext, second)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf(
			"contending receipt lock error = %v; want deadline refusal",
			err,
		)
	}
	if err := second.Close(); err != nil {
		t.Fatalf("close canceled receipt lock: %v", err)
	}
	if err := syscall.Flock(int(first.Fd()), syscall.LOCK_UN); err != nil {
		t.Fatalf("release first receipt lock: %v", err)
	}
	if err := first.Close(); err != nil {
		t.Fatalf("close first receipt lock: %v", err)
	}
	afterRelease, err := openReceiptLockFile(lockPath)
	if err != nil {
		t.Fatalf("open receipt lock after cancellation: %v", err)
	}
	if err := lockReceiptFile(context.Background(), afterRelease); err != nil {
		t.Fatalf("reacquire receipt lock after cancellation: %v", err)
	}
	if err := syscall.Flock(int(afterRelease.Fd()), syscall.LOCK_UN); err != nil {
		t.Fatalf("release reacquired receipt lock: %v", err)
	}
	if err := afterRelease.Close(); err != nil {
		t.Fatalf("close reacquired receipt lock: %v", err)
	}
	lockInfo, err := os.Lstat(lockPath)
	if err != nil {
		t.Fatalf("inspect persistent receipt lock identity: %v", err)
	}
	if !lockInfo.Mode().IsRegular() || lockInfo.Mode().Perm() != 0o600 {
		t.Fatalf("receipt lock mode = %v, want regular 0600", lockInfo.Mode())
	}
}

func TestReceiptLockRejectsSymlinkIdentity(t *testing.T) {
	directory := t.TempDir()
	lockPath := filepath.Join(directory, "prepare.json.lock")
	target := filepath.Join(directory, "lock-target")
	if err := os.WriteFile(target, nil, 0o600); err != nil {
		t.Fatalf("write lock target: %v", err)
	}
	if err := os.Symlink(target, lockPath); err != nil {
		t.Fatalf("create receipt lock symlink: %v", err)
	}
	file, err := openReceiptLockFile(lockPath)
	if file != nil || err == nil {
		t.Fatalf(
			"symlink lock open = %#v, %v; want no-follow refusal",
			file,
			err,
		)
	}
}

func TestPublishRefusesChangedReceiptBeforeRemoteMutation(t *testing.T) {
	fixture := newIntegrationDeliveryFixture(t)
	prepared := fixture.prepare(t)
	relative := "out/delivery/prepare.json"
	absolute := filepath.Join(fixture.work, filepath.FromSlash(relative))
	changed := prepared.Receipt
	var err error
	changed.RevisionNonce, err = newReceiptRevisionNonce()
	if err != nil {
		t.Fatalf("generate external revision: %v", err)
	}
	changedContents, err := encodePreparationReceipt(changed)
	if err != nil {
		t.Fatalf("encode external receipt revision: %v", err)
	}
	fixture.forge.pullRequestsCalls = 0
	fixture.forge.onPullRequestsCall = func(call int) {
		if call != 1 {
			return
		}
		if writeErr := os.WriteFile(absolute, changedContents, 0o600); writeErr != nil {
			t.Fatalf("install external receipt revision: %v", writeErr)
		}
	}

	report, err := fixture.delivery.publish(
		context.Background(),
		publishOptions{
			ValidatedHead: prepared.HeadOID,
			ReceiptFile:   relative,
		},
	)
	if report != nil || err == nil || !strings.Contains(
		err.Error(),
		"pre-push preparation receipt gate",
	) {
		t.Fatalf(
			"publish after receipt replacement = %#v, %v; want pre-push refusal",
			report,
			err,
		)
	}
	if _, err := runTestGitAllowFailure(
		fixture.remote,
		"rev-parse",
		"--verify",
		"refs/heads/feature",
	); err == nil {
		t.Fatal("publish mutated the remote before rejecting the stale receipt")
	}
}
