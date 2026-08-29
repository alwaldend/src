package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func testPullRequest() pullRequest {
	return pullRequest{
		ID:                  "PR_node",
		Number:              7,
		URL:                 "https://github.com/owner/repo/pull/7",
		State:               "OPEN",
		Title:               "Title",
		Body:                "Body",
		AuthorLogin:         "agent",
		BaseRefName:         "master",
		BaseRefOID:          testOID('a'),
		HeadRefName:         "feature;not-shell",
		HeadRefOID:          testOID('b'),
		HeadRepositoryOwner: "owner",
		HeadRepositoryName:  "repo",
		UpdatedAt:           "2026-08-30T00:00:00Z",
	}
}

func githubGraphQLTestCommand(
	t *testing.T,
	repository remoteRepository,
	query string,
	variables map[string]any,
) command {
	t.Helper()
	request, err := json.Marshal(struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}{Query: query, Variables: variables})
	if err != nil {
		t.Fatal(err)
	}
	return githubCommand([]string{
		"api", "--hostname", repository.Host,
		"graphql", "--input", "-",
	}, string(request))
}

func githubGraphQLPullRequestMap() map[string]any {
	return map[string]any{
		"id":                  "PR_node",
		"number":              7,
		"url":                 "https://github.com/owner/repo/pull/7",
		"state":               "OPEN",
		"title":               "Title",
		"body":                "Body",
		"author":              map[string]any{"login": "agent"},
		"baseRefName":         "master",
		"baseRefOid":          testOID('a'),
		"headRefName":         "feature;not-shell",
		"headRefOid":          testOID('b'),
		"headRepository":      map[string]any{"name": "repo"},
		"headRepositoryOwner": map[string]any{"login": "owner"},
		"isCrossRepository":   false,
		"isDraft":             false,
		"updatedAt":           "2026-08-30T00:00:00Z",
	}
}

func githubGraphQLResult(t *testing.T, data map[string]any) commandResult {
	t.Helper()
	encoded, err := json.Marshal(map[string]any{"data": data})
	if err != nil {
		t.Fatal(err)
	}
	return commandResult{Stdout: string(encoded)}
}

func githubPullRequestConnectionResult(
	t *testing.T,
	connection string,
	value any,
) commandResult {
	t.Helper()
	pullRequest := githubGraphQLPullRequestMap()
	pullRequest[connection] = value
	return githubGraphQLResult(t, map[string]any{
		"repository": map[string]any{"pullRequest": pullRequest},
	})
}

func githubTopComment(id string) map[string]any {
	return map[string]any{
		"id":        id,
		"url":       "https://github.com/owner/repo/pull/7#" + id,
		"body":      "Comment " + id,
		"author":    map[string]any{"login": "reviewer"},
		"createdAt": "2026-08-30T00:00:00Z",
		"updatedAt": "2026-08-30T00:00:00Z",
	}
}

func githubReviewNode(id string) map[string]any {
	return map[string]any{
		"id": id, "url": "https://github.com/owner/repo/pull/7#" + id,
		"body":   "Review " + id,
		"author": map[string]any{"login": "reviewer"},
		"state":  "COMMENTED", "submittedAt": "2026-08-30T00:00:00Z",
		"commit": map[string]any{"oid": testOID('b')},
	}
}

func githubReviewRequestNode(id string) map[string]any {
	return map[string]any{
		"id": "RR_" + id,
		"requestedReviewer": map[string]any{
			"__typename": "User", "id": "U_" + id, "login": id,
		},
	}
}

func githubThreadNode(id string) map[string]any {
	return map[string]any{
		"id": id, "isResolved": false, "isOutdated": false,
		"path": "file.go", "line": 12, "originalLine": 12,
		"startLine": nil, "originalStartLine": nil,
	}
}

func githubReviewCommentNode(id string) map[string]any {
	return map[string]any{
		"id": id, "url": "https://github.com/owner/repo/pull/7#" + id,
		"body":      "Comment " + id,
		"author":    map[string]any{"login": "reviewer"},
		"createdAt": "2026-08-30T00:00:00Z",
		"updatedAt": "2026-08-30T00:00:00Z",
		"path":      "file.go", "line": 12, "originalLine": 12,
		"commit": map[string]any{"oid": testOID('b')},
	}
}

func githubThreadCommentsResult(
	t *testing.T,
	threadID string,
	comments any,
) commandResult {
	t.Helper()
	thread := githubThreadNode(threadID)
	thread["comments"] = comments
	return githubGraphQLResult(t, map[string]any{
		"repository": map[string]any{
			"pullRequest": githubGraphQLPullRequestMap(),
		},
		"node": thread,
	})
}

func TestGitHubTopLevelCommentsPaginatesCompletely(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubGraphQLTestCommand(
				t,
				repository,
				githubCommentsQuery,
				githubVariables(repository, pullRequest, ""),
			),
			result: githubPullRequestConnectionResult(t, "comments", map[string]any{
				"nodes": []any{githubTopComment("C1")},
				"pageInfo": map[string]any{
					"hasNextPage": true,
					"endCursor":   "cursor-1",
				},
			}),
		},
		{
			command: githubGraphQLTestCommand(
				t,
				repository,
				githubCommentsQuery,
				githubVariables(repository, pullRequest, "cursor-1"),
			),
			result: githubPullRequestConnectionResult(t, "comments", map[string]any{
				"nodes": []any{githubTopComment("C2")},
				"pageInfo": map[string]any{
					"hasNextPage": false,
					"endCursor":   nil,
				},
			}),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	comments, err := forge.topLevelComments(
		context.Background(),
		repository,
		pullRequest,
	)
	if err != nil {
		t.Fatalf("topLevelComments() error = %v", err)
	}
	if len(comments) != 2 || comments[0].ID != "C1" || comments[1].ID != "C2" {
		t.Fatalf("topLevelComments() = %#v", comments)
	}
	runner.done()
}

func TestGitHubReviewsPaginateCompletely(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubGraphQLTestCommand(
				t, repository, githubReviewsQuery,
				githubVariables(repository, pullRequest, ""),
			),
			result: githubPullRequestConnectionResult(t, "reviews", map[string]any{
				"nodes": []any{githubReviewNode("R1")},
				"pageInfo": map[string]any{
					"hasNextPage": true, "endCursor": "reviews-1",
				},
			}),
		},
		{
			command: githubGraphQLTestCommand(
				t, repository, githubReviewsQuery,
				githubVariables(repository, pullRequest, "reviews-1"),
			),
			result: githubPullRequestConnectionResult(t, "reviews", map[string]any{
				"nodes": []any{githubReviewNode("R2")},
				"pageInfo": map[string]any{
					"hasNextPage": false, "endCursor": nil,
				},
			}),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	reviews, err := forge.pullRequestReviews(
		context.Background(), repository, pullRequest,
	)
	if err != nil {
		t.Fatalf("pullRequestReviews() error = %v", err)
	}
	if len(reviews) != 2 || reviews[0].ID != "R1" || reviews[1].ID != "R2" {
		t.Fatalf("pullRequestReviews() = %#v", reviews)
	}
	runner.done()
}

func TestGitHubReviewRequestsPaginateCompletely(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubGraphQLTestCommand(
				t, repository, githubReviewRequestsQuery,
				githubVariables(repository, pullRequest, ""),
			),
			result: githubPullRequestConnectionResult(
				t, "reviewRequests", map[string]any{
					"nodes": []any{githubReviewRequestNode("alice")},
					"pageInfo": map[string]any{
						"hasNextPage": true, "endCursor": "requests-1",
					},
				},
			),
		},
		{
			command: githubGraphQLTestCommand(
				t, repository, githubReviewRequestsQuery,
				githubVariables(repository, pullRequest, "requests-1"),
			),
			result: githubPullRequestConnectionResult(
				t, "reviewRequests", map[string]any{
					"nodes": []any{githubReviewRequestNode("bob")},
					"pageInfo": map[string]any{
						"hasNextPage": false, "endCursor": nil,
					},
				},
			),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	reviewers, err := forge.requestedReviewers(
		context.Background(), repository, pullRequest,
	)
	if err != nil {
		t.Fatalf("requestedReviewers() error = %v", err)
	}
	if len(reviewers) != 2 || reviewers[0].Login != "alice" ||
		reviewers[1].Login != "bob" {
		t.Fatalf("requestedReviewers() = %#v", reviewers)
	}
	runner.done()
}

func TestGitHubThreadCommentsPaginateCompletely(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	firstVariables := githubVariables(repository, pullRequest, "")
	firstVariables["threadId"] = "RT1"
	secondVariables := githubVariables(repository, pullRequest, "comments-1")
	secondVariables["threadId"] = "RT1"
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubGraphQLTestCommand(
				t, repository, githubThreadCommentsQuery, firstVariables,
			),
			result: githubThreadCommentsResult(t, "RT1", map[string]any{
				"nodes": []any{githubReviewCommentNode("RC1")},
				"pageInfo": map[string]any{
					"hasNextPage": true, "endCursor": "comments-1",
				},
			}),
		},
		{
			command: githubGraphQLTestCommand(
				t, repository, githubThreadCommentsQuery, secondVariables,
			),
			result: githubThreadCommentsResult(t, "RT1", map[string]any{
				"nodes": []any{githubReviewCommentNode("RC2")},
				"pageInfo": map[string]any{
					"hasNextPage": false, "endCursor": nil,
				},
			}),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	thread, err := forge.reviewThreadDetails(
		context.Background(), repository, pullRequest, "RT1",
	)
	if err != nil {
		t.Fatalf("reviewThreadDetails() error = %v", err)
	}
	if len(thread.Comments) != 2 || thread.Comments[0].ID != "RC1" ||
		thread.Comments[1].ID != "RC2" ||
		thread.ExpectationDigest != reviewThreadDigest(thread) {
		t.Fatalf("reviewThreadDetails() = %#v", thread)
	}
	runner.done()
}

func TestGitHubReviewThreadsPaginateCompletely(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	threadOneVariables := githubVariables(repository, pullRequest, "")
	threadOneVariables["threadId"] = "RT1"
	threadTwoVariables := githubVariables(repository, pullRequest, "")
	threadTwoVariables["threadId"] = "RT2"
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubGraphQLTestCommand(
				t, repository, githubThreadsQuery,
				githubVariables(repository, pullRequest, ""),
			),
			result: githubPullRequestConnectionResult(
				t, "reviewThreads", map[string]any{
					"nodes": []any{githubThreadNode("RT1")},
					"pageInfo": map[string]any{
						"hasNextPage": true, "endCursor": "threads-1",
					},
				},
			),
		},
		{
			command: githubGraphQLTestCommand(
				t, repository, githubThreadCommentsQuery, threadOneVariables,
			),
			result: githubThreadCommentsResult(t, "RT1", map[string]any{
				"nodes": []any{githubReviewCommentNode("RC1")},
				"pageInfo": map[string]any{
					"hasNextPage": false, "endCursor": nil,
				},
			}),
		},
		{
			command: githubGraphQLTestCommand(
				t, repository, githubThreadsQuery,
				githubVariables(repository, pullRequest, "threads-1"),
			),
			result: githubPullRequestConnectionResult(
				t, "reviewThreads", map[string]any{
					"nodes": []any{githubThreadNode("RT2")},
					"pageInfo": map[string]any{
						"hasNextPage": false, "endCursor": nil,
					},
				},
			),
		},
		{
			command: githubGraphQLTestCommand(
				t, repository, githubThreadCommentsQuery, threadTwoVariables,
			),
			result: githubThreadCommentsResult(t, "RT2", map[string]any{
				"nodes": []any{githubReviewCommentNode("RC2")},
				"pageInfo": map[string]any{
					"hasNextPage": false, "endCursor": nil,
				},
			}),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	threads, err := forge.reviewThreads(
		context.Background(), repository, pullRequest,
	)
	if err != nil {
		t.Fatalf("reviewThreads() error = %v", err)
	}
	if len(threads) != 2 || threads[0].ID != "RT1" ||
		threads[1].ID != "RT2" || threads[0].ExpectationDigest == "" ||
		threads[1].ExpectationDigest == "" {
		t.Fatalf("reviewThreads() = %#v", threads)
	}
	runner.done()
}

func TestGitHubTopLevelCommentsRefusesNullPullRequest(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t,
			repository,
			githubCommentsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: githubGraphQLResult(t, map[string]any{
			"repository": map[string]any{"pullRequest": nil},
		}),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.topLevelComments(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), "exact pull request") {
		t.Fatalf("topLevelComments() error = %v", err)
	}
	runner.done()
}

func TestGitHubTopLevelCommentsRefuseNullConnection(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t,
			repository,
			githubCommentsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: githubPullRequestConnectionResult(t, "comments", nil),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.topLevelComments(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), "null or missing connection") {
		t.Fatalf("topLevelComments() error = %v", err)
	}
	runner.done()
}

func TestGitHubReviewsRefuseMissingNodes(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t,
			repository,
			githubReviewsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: githubPullRequestConnectionResult(t, "reviews", map[string]any{
			"pageInfo": map[string]any{
				"hasNextPage": false, "endCursor": nil,
			},
		}),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.pullRequestReviews(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), `omits required key "nodes"`) {
		t.Fatalf("pullRequestReviews() error = %v", err)
	}
	runner.done()
}

func TestGitHubReviewsRefuseNullConnection(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t, repository, githubReviewsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: githubPullRequestConnectionResult(t, "reviews", nil),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.pullRequestReviews(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), "null or missing connection") {
		t.Fatalf("pullRequestReviews() error = %v", err)
	}
	runner.done()
}

func TestGitHubReviewRequestsRefuseMissingPageInfo(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t,
			repository,
			githubReviewRequestsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: githubPullRequestConnectionResult(
			t,
			"reviewRequests",
			map[string]any{"nodes": []any{}},
		),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.requestedReviewers(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), `omits required key "pageInfo"`) {
		t.Fatalf("requestedReviewers() error = %v", err)
	}
	runner.done()
}

func TestGitHubReviewRequestsRefuseNullConnection(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t, repository, githubReviewRequestsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: githubPullRequestConnectionResult(t, "reviewRequests", nil),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.requestedReviewers(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), "null or missing connection") {
		t.Fatalf("requestedReviewers() error = %v", err)
	}
	runner.done()
}

func TestGitHubReviewThreadsRefuseNullConnection(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t, repository, githubThreadsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: githubPullRequestConnectionResult(t, "reviewThreads", nil),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.reviewThreads(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), "null or missing connection") {
		t.Fatalf("reviewThreads() error = %v", err)
	}
	runner.done()
}

func TestGitHubThreadCommentsRefuseNullConnection(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	variables := githubVariables(repository, pullRequest, "")
	variables["threadId"] = "RT1"
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t, repository, githubThreadCommentsQuery, variables,
		),
		result: githubThreadCommentsResult(t, "RT1", nil),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.reviewThreadDetails(
		context.Background(), repository, pullRequest, "RT1",
	)
	if err == nil || !strings.Contains(err.Error(), "null or missing connection") {
		t.Fatalf("reviewThreadDetails() error = %v", err)
	}
	runner.done()
}

func TestNextGitHubCursorRefusesRepeatedCursor(t *testing.T) {
	t.Parallel()
	hasNextPage := true
	endCursor := "repeated"
	_, _, err := nextGithubCursor(
		"reviews",
		2,
		githubPageInfo{
			HasNextPage: &hasNextPage,
			EndCursor:   &endCursor,
		},
		map[string]bool{"repeated": true},
	)
	if err == nil || !strings.Contains(err.Error(), "invalid pagination cursor") {
		t.Fatalf("nextGithubCursor() error = %v", err)
	}
}

func TestNextGitHubCursorRefusesPageLimit(t *testing.T) {
	t.Parallel()
	hasNextPage := true
	endCursor := "next"
	_, _, err := nextGithubCursor(
		"review thread comments",
		githubPaginationLimit,
		githubPageInfo{
			HasNextPage: &hasNextPage,
			EndCursor:   &endCursor,
		},
		map[string]bool{},
	)
	if err == nil || !strings.Contains(err.Error(), "exceeded 1000 pages") {
		t.Fatalf("nextGithubCursor() error = %v", err)
	}
}

func TestGitHubThreadCommentsRefuseNullThread(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	variables := githubVariables(repository, pullRequest, "")
	variables["threadId"] = "RT_missing"
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t,
			repository,
			githubThreadCommentsQuery,
			variables,
		),
		result: githubGraphQLResult(t, map[string]any{
			"repository": map[string]any{
				"pullRequest": githubGraphQLPullRequestMap(),
			},
			"node": nil,
		}),
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.reviewThreadDetails(
		context.Background(),
		repository,
		pullRequest,
		"RT_missing",
	)
	if err == nil || !strings.Contains(err.Error(), "exact review thread") {
		t.Fatalf("reviewThreadDetails() error = %v", err)
	}
	runner.done()
}

func TestGitHubCommandsForceTTYOffByUnsettingEnvironment(t *testing.T) {
	t.Parallel()
	for _, value := range githubEnvironment {
		if strings.HasPrefix(value, "GH_FORCE_TTY=") {
			t.Fatalf("githubEnvironment contains %q", value)
		}
	}
	for _, name := range []string{
		"GH_FORCE_TTY",
		"GODEBUG",
		"SSLKEYLOGFILE",
		"CLICOLOR_FORCE",
	} {
		if !hasExactString(githubUnsetEnvironment, name) {
			t.Errorf("githubUnsetEnvironment does not unset %q: %#v", name, githubUnsetEnvironment)
		}
	}
}

func TestGitHubInventoryBudgetRefusesUnboundedAccumulation(t *testing.T) {
	t.Parallel()
	for name, test := range map[string]func(*githubInventoryBudget) error{
		"requests": func(budget *githubInventoryBudget) error {
			budget.Requests = githubInventoryRequestLimit
			return budget.consumeRequest("comments")
		},
		"nodes": func(budget *githubInventoryBudget) error {
			budget.Nodes = githubInventoryNodeLimit
			return budget.consumeNode("comments")
		},
		"bytes": func(budget *githubInventoryBudget) error {
			budget.Bytes = githubInventoryByteLimit
			return budget.consumeBytes("comments", 1)
		},
	} {
		if err := test(&githubInventoryBudget{}); err == nil ||
			!strings.Contains(err.Error(), "refusing an incomplete result") {
			t.Errorf("%s budget error = %v", name, err)
		}
	}
}

func TestGitHubInventoryBudgetChargesAllRawResponseBytes(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	comment := githubTopComment("C1")
	comment["body"] = ""
	result := githubPullRequestConnectionResult(t, "comments", map[string]any{
		"nodes": []any{comment},
		"pageInfo": map[string]any{
			"hasNextPage": false, "endCursor": nil,
		},
	})
	if len(result.Stdout) >= githubInventoryByteLimit {
		t.Fatal("test fixture unexpectedly exceeds the inventory limit")
	}
	budget := &githubInventoryBudget{
		Bytes: githubInventoryByteLimit - len(result.Stdout) + 1,
	}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubGraphQLTestCommand(
			t, repository, githubCommentsQuery,
			githubVariables(repository, pullRequest, ""),
		),
		result: result,
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.topLevelCommentsWithBudget(
		context.Background(), repository, pullRequest, budget,
	)
	if err == nil || !strings.Contains(err.Error(), "response bytes") {
		t.Fatalf("topLevelCommentsWithBudget() error = %v", err)
	}
	runner.done()
}

func TestGitHubNestedReviewRecordsRejectMalformedMetadata(t *testing.T) {
	t.Parallel()
	submittedAt := "2026-08-30T00:00:00Z"
	topComment := githubTopLevelComment{
		ID: "C1", URL: "https://github.com/owner/repo/pull/7#issuecomment-1",
		Body: "Comment", Author: githubIdentity{Login: "reviewer"},
		CreatedAt: "2026-08-30T00:00:00Z", UpdatedAt: "2026-08-30T00:00:00Z",
	}
	reviewComment := githubReviewComment{
		ID: "RC1", URL: "https://github.com/owner/repo/pull/7#discussion-1",
		Body: "Comment", Author: githubIdentity{Login: "reviewer"},
		CreatedAt: "2026-08-30T00:00:00Z", UpdatedAt: "2026-08-30T00:00:00Z",
		Path: "file.go", Commit: githubCommit{OID: testOID('b')},
	}
	review := githubReview{
		ID: "R1", URL: "https://github.com/owner/repo/pull/7#review-1",
		Body: "Review", Author: githubIdentity{Login: "reviewer"},
		State: "COMMENTED", SubmittedAt: &submittedAt,
		Commit: githubCommit{OID: testOID('b')},
	}
	thread := githubReviewThread{ID: "RT1", Path: "file.go"}

	invalidTopURL := topComment
	invalidTopURL.URL = "not a URL"
	missingTopAuthor := topComment
	missingTopAuthor.Author.Login = ""
	invalidTopTimestamp := topComment
	invalidTopTimestamp.UpdatedAt = "yesterday"
	invalidCommentCommit := reviewComment
	invalidCommentCommit.Commit.OID = "short"
	invalidCommentPath := reviewComment
	invalidCommentPath.Path = ""
	invalidReviewState := review
	invalidReviewState.State = "UNKNOWN"
	invalidReviewCommit := review
	invalidReviewCommit.Commit.OID = "short"
	invalidThreadPath := thread
	invalidThreadPath.Path = ""
	badLine := 0
	invalidThreadLine := thread
	invalidThreadLine.Line = &badLine

	tests := map[string]func() error{
		"top-level URL":       func() error { return validateGitHubTopLevelComment(invalidTopURL) },
		"top-level author":    func() error { return validateGitHubTopLevelComment(missingTopAuthor) },
		"top-level timestamp": func() error { return validateGitHubTopLevelComment(invalidTopTimestamp) },
		"comment commit":      func() error { return validateGitHubReviewComment(invalidCommentCommit) },
		"comment path":        func() error { return validateGitHubReviewComment(invalidCommentPath) },
		"review state":        func() error { return validateGitHubReview(invalidReviewState) },
		"review commit":       func() error { return validateGitHubReview(invalidReviewCommit) },
		"thread path":         func() error { return validateGitHubReviewThread(invalidThreadPath) },
		"thread line":         func() error { return validateGitHubReviewThread(invalidThreadLine) },
	}
	for name, validate := range tests {
		name, validate := name, validate
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if err := validate(); err == nil {
				t.Fatal("validator unexpectedly accepted malformed metadata")
			}
		})
	}
}

func TestGitHubPendingReviewAcceptsNullSubmittedAt(t *testing.T) {
	t.Parallel()
	var review githubReview
	output := []byte(`{"id":"R1","url":"https://github.com/owner/repo/pull/7#review-1","body":"draft","author":{"login":"reviewer"},"state":"PENDING","submittedAt":null,"commit":{"oid":"` + testOID('b') + `"}}`)
	if err := decodeGitHubJSON("pending review", output, &review); err != nil {
		t.Fatalf("decodeGitHubJSON() error = %v", err)
	}
	if err := validateGitHubReview(review); err != nil {
		t.Fatalf("validateGitHubReview() error = %v", err)
	}
	if review.pullRequestReview().SubmittedAt != "" {
		t.Fatalf("pending review submittedAt = %q", review.pullRequestReview().SubmittedAt)
	}
}

func TestGitHubGraphQLStrictEnvelopeAndData(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	query := "query Test { viewer { login } }"
	variables := map[string]any{}
	tests := map[string]struct {
		output  string
		wantErr bool
	}{
		"standard extensions": {
			output: `{"data":{"ok":true},"extensions":{"future":{"value":1}}}`,
		},
		"unknown data": {
			output:  `{"data":{"ok":true,"unexpected":1}}`,
			wantErr: true,
		},
		"duplicate in extension": {
			output:  `{"data":{"ok":true},"extensions":{"future":1,"Future":2}}`,
			wantErr: true,
		},
		"case-variant envelope": {
			output:  `{"Data":{"ok":true}}`,
			wantErr: true,
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := &transcriptRunner{t: t, expected: []expectedCommand{{
				command: githubGraphQLTestCommand(t, repository, query, variables),
				result:  commandResult{Stdout: test.output},
			}}}
			forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
			var output struct {
				OK bool `json:"ok"`
			}
			err := forge.graphQL(
				context.Background(), repository, query, variables, &output,
			)
			if test.wantErr && err == nil {
				t.Fatal("graphQL() unexpectedly accepted ambiguous output")
			}
			if !test.wantErr && (err != nil || !output.OK) {
				t.Fatalf("graphQL() output = %#v, error = %v", output, err)
			}
			runner.done()
		})
	}
}

func TestGitHubRequestedReviewerSupportsEverySchemaVariant(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input githubRequestedReviewer
		want  requestedReviewer
	}{
		{
			input: githubRequestedReviewer{Typename: "User", ID: "U1", Login: "alice"},
			want:  requestedReviewer{Type: "user", ID: "U1", Login: "alice"},
		},
		{
			input: githubRequestedReviewer{
				Typename: "Bot", ID: "B1", Login: "app[bot]",
			},
			want: requestedReviewer{Type: "bot", ID: "B1", Login: "app[bot]"},
		},
		{
			input: githubRequestedReviewer{
				Typename: "Mannequin", ID: "M1", Login: "imported",
			},
			want: requestedReviewer{
				Type: "mannequin", ID: "M1", Login: "imported",
			},
		},
		{
			input: githubRequestedReviewer{
				Typename: "Team", ID: "T1", Slug: "reviewers",
				Organization: githubIdentity{Login: "owner"},
			},
			want: requestedReviewer{
				Type: "team", ID: "T1", Login: "owner/reviewers",
			},
		},
		{
			input: githubRequestedReviewer{
				Typename: "EnterpriseTeam", ID: "ET1",
			},
			want: requestedReviewer{
				Type: "enterprise_team", ID: "ET1",
			},
		},
	}
	for _, test := range tests {
		got, err := test.input.requestedReviewer()
		if err != nil {
			t.Errorf("requestedReviewer(%q) error = %v", test.input.Typename, err)
		} else if got != test.want {
			t.Errorf("requestedReviewer(%q) = %#v, want %#v", test.input.Typename, got, test.want)
		}
	}
}

type reviewLifecycleRunner struct {
	t                           *testing.T
	threadID                    string
	comments                    []reviewComment
	topComments                 []reviewComment
	reviews                     []pullRequestReview
	resolved                    bool
	reviewers                   []string
	completeReviewImmediately   bool
	commentErr                  error
	commentResponse             string
	postMutationViewErr         error
	mutated                     bool
	nullReplyResponse           bool
	replyMutationCommentID      string
	pullRequestUpdatedAt        string
	advanceEpochAfterReply      bool
	epochAdvanced               bool
	staleReplyInventoryOnce     bool
	postMutationViewCount       int
	postMutationViewEpochs      []string
	replyMutationCount          int
	sawComment                  bool
	sawReply                    bool
	sawResolve                  bool
	sawRequest                  bool
	appendBeforeFinalThreadRead bool
	threadCommentsReadCount     int
}

type reviewInventoryResponse struct {
	Name  string
	Bytes int
}

type aggregateReviewInventoryRunner struct {
	t         *testing.T
	delegate  *reviewLifecycleRunner
	padding   string
	responses []reviewInventoryResponse
}

func (r *aggregateReviewInventoryRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	r.t.Helper()
	name := ""
	if len(request.Args) >= 2 && request.Args[0] == "pr" &&
		request.Args[1] == "view" {
		name = "pull-request coherence reads"
	} else if hasExactString(request.Args, "graphql") {
		var payload struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal([]byte(request.Stdin), &payload); err != nil {
			r.t.Fatalf("decode GraphQL request: %v", err)
		}
		switch payload.Query {
		case githubCommentsQuery:
			name = "comments"
		case githubReviewsQuery:
			name = "reviews"
		case githubThreadsQuery:
			name = "review threads"
		case githubThreadCommentsQuery:
			name = "review thread comments"
		case githubReviewRequestsQuery:
			name = "review requests"
		}
	}
	if name == "" {
		r.t.Fatalf("unexpected review inventory request: %#v", request)
	}
	result, err := r.delegate.Run(ctx, request)
	if err != nil {
		return result, err
	}
	result.Stdout += r.padding
	r.responses = append(r.responses, reviewInventoryResponse{
		Name:  name,
		Bytes: len(result.Stdout),
	})
	return result, nil
}

func newReviewLifecycleRunner(t *testing.T) *reviewLifecycleRunner {
	t.Helper()
	return &reviewLifecycleRunner{
		t:                    t,
		threadID:             "RT1",
		pullRequestUpdatedAt: "2026-08-30T00:00:00Z",
		topComments: []reviewComment{{
			ID:          "IC1",
			URL:         "https://github.com/owner/repo/pull/7#issuecomment-1",
			Body:        "Comment IC1",
			AuthorLogin: "reviewer",
			CreatedAt:   "2026-08-30T00:00:00Z",
			UpdatedAt:   "2026-08-30T00:00:00Z",
		}},
		comments: []reviewComment{{
			ID:          "RC_old",
			URL:         "https://github.com/owner/repo/pull/7#discussion-old",
			Body:        "Please fix this.",
			AuthorLogin: "reviewer",
			CreatedAt:   "2026-08-30T00:00:00Z",
			UpdatedAt:   "2026-08-30T00:00:00Z",
			Path:        "file.go",
			CommitOID:   testOID('b'),
		}},
		reviews: []pullRequestReview{{
			ID:          "R1",
			URL:         "https://github.com/owner/repo/pull/7#review-1",
			Body:        "Review",
			AuthorLogin: "reviewer",
			State:       "COMMENTED",
			SubmittedAt: "2026-08-30T00:00:00Z",
			CommitOID:   testOID('b'),
		}},
	}
}

func (r *reviewLifecycleRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	r.t.Helper()
	if request.Name != "gh-test" || request.Dir != "/repo" ||
		!reflect.DeepEqual(request.Env, githubEnvironment) ||
		!reflect.DeepEqual(
			request.UnsetEnv,
			append(githubUnsetEnvironment, gitUnsetEnvironment...),
		) || !reflect.DeepEqual(
		request.UnsetEnvPrefixes,
		gitUnsetEnvironmentPrefixes,
	) || request.OutputLimit != githubOutputLimit {
		r.t.Fatalf("unsafe GitHub command: %#v", request)
	}
	if len(request.Args) >= 2 && request.Args[0] == "pr" &&
		request.Args[1] == "view" {
		if r.mutated && r.postMutationViewErr != nil {
			return commandResult{}, r.postMutationViewErr
		}
		if r.mutated {
			r.postMutationViewCount++
			index := r.postMutationViewCount - 1
			if index < len(r.postMutationViewEpochs) {
				r.pullRequestUpdatedAt = r.postMutationViewEpochs[index]
			}
		}
		output := strings.Replace(
			githubPRJSON("Title", "Body"),
			`"updatedAt":"2026-08-30T00:00:00Z"`,
			fmt.Sprintf(`"updatedAt":%q`, r.pullRequestUpdatedAt),
			1,
		)
		return commandResult{Stdout: output}, nil
	}
	if len(request.Args) >= 1 && request.Args[0] == "api" {
		if hasExactString(request.Args, "graphql") {
			return r.graphQL(request)
		}
		if hasExactString(
			request.Args,
			"repos/owner/repo/issues/7/comments",
		) {
			return r.commentOnPullRequest(request)
		}
		return r.requestReview(request)
	}
	r.t.Fatalf("unexpected GitHub command: %#v", request)
	return commandResult{}, nil
}

func hasExactString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func (r *reviewLifecycleRunner) pullRequestMap() map[string]any {
	pullRequest := githubGraphQLPullRequestMap()
	pullRequest["updatedAt"] = r.pullRequestUpdatedAt
	return pullRequest
}

func (r *reviewLifecycleRunner) pullRequestConnectionResult(
	name string,
	connection any,
) commandResult {
	pullRequest := r.pullRequestMap()
	pullRequest[name] = connection
	return githubGraphQLResult(r.t, map[string]any{
		"repository": map[string]any{"pullRequest": pullRequest},
	})
}

func (r *reviewLifecycleRunner) graphQL(
	request command,
) (commandResult, error) {
	var payload struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}
	if err := json.Unmarshal([]byte(request.Stdin), &payload); err != nil {
		r.t.Fatalf("decode GraphQL request: %v", err)
	}
	switch payload.Query {
	case githubCommentsQuery:
		if r.mutated && r.advanceEpochAfterReply &&
			r.postMutationViewCount == 1 && !r.epochAdvanced {
			r.pullRequestUpdatedAt = "2026-08-30T00:01:00Z"
			r.reviews = append(r.reviews, pullRequestReview{
				ID: "R_reply", URL: "https://github.com/owner/repo/pull/7#review-reply",
				Body: "", AuthorLogin: "agent", State: "COMMENTED",
				SubmittedAt: "2026-08-30T00:01:00Z", CommitOID: testOID('b'),
			})
			r.epochAdvanced = true
		}
		nodes := make([]any, 0, len(r.topComments))
		for _, comment := range r.topComments {
			nodes = append(nodes, map[string]any{
				"id": comment.ID, "url": comment.URL, "body": comment.Body,
				"author":    map[string]any{"login": comment.AuthorLogin},
				"createdAt": comment.CreatedAt,
				"updatedAt": comment.UpdatedAt,
			})
		}
		return r.pullRequestConnectionResult("comments", map[string]any{
			"nodes": nodes,
			"pageInfo": map[string]any{
				"hasNextPage": false,
				"endCursor":   nil,
			},
		}), nil
	case githubReviewsQuery:
		nodes := make([]any, 0, len(r.reviews))
		for _, review := range r.reviews {
			nodes = append(nodes, map[string]any{
				"id": review.ID, "url": review.URL, "body": review.Body,
				"author":      map[string]any{"login": review.AuthorLogin},
				"state":       review.State,
				"submittedAt": review.SubmittedAt,
				"commit":      map[string]any{"oid": review.CommitOID},
			})
		}
		return r.pullRequestConnectionResult("reviews", map[string]any{
			"nodes": nodes,
			"pageInfo": map[string]any{
				"hasNextPage": false,
				"endCursor":   nil,
			},
		}), nil
	case githubThreadsQuery:
		return r.pullRequestConnectionResult("reviewThreads", map[string]any{
			"nodes": []any{r.threadMap()},
			"pageInfo": map[string]any{
				"hasNextPage": false,
				"endCursor":   nil,
			},
		}), nil
	case githubThreadCommentsQuery:
		r.requireThreadID(payload.Variables)
		r.threadCommentsReadCount++
		if r.appendBeforeFinalThreadRead && r.threadCommentsReadCount == 2 {
			r.comments = append(r.comments, reviewComment{
				ID: "RC_human_late", URL: "https://github.com/owner/repo/pull/7#late",
				Body: "late concurrent reply", AuthorLogin: "human",
				CreatedAt: "2026-08-30T00:02:00Z", UpdatedAt: "2026-08-30T00:02:00Z",
				Path: "file.go", CommitOID: testOID('b'),
			})
		}
		return githubGraphQLResult(r.t, map[string]any{
			"repository": map[string]any{
				"pullRequest": r.pullRequestMap(),
			},
			"node": r.threadCommentsMap(),
		}), nil
	case githubReviewRequestsQuery:
		nodes := make([]any, 0, len(r.reviewers))
		for _, reviewer := range r.reviewers {
			nodes = append(nodes, map[string]any{
				"id": "RR_" + reviewer,
				"requestedReviewer": map[string]any{
					"__typename": "User",
					"id":         "U_" + reviewer,
					"login":      reviewer,
				},
			})
		}
		return r.pullRequestConnectionResult("reviewRequests", map[string]any{
			"nodes": nodes,
			"pageInfo": map[string]any{
				"hasNextPage": false,
				"endCursor":   nil,
			},
		}), nil
	case githubReplyMutation:
		r.requireThreadID(payload.Variables)
		body, _ := payload.Variables["body"].(string)
		if hasExactString(request.Args, body) {
			r.t.Fatal("review reply body leaked into the process argument vector")
		}
		if !hasFinalLine(body, commentDisclaimer) {
			r.t.Fatalf("review reply lacks disclaimer: %q", body)
		}
		r.sawReply = true
		r.replyMutationCount++
		r.mutated = true
		r.comments = append(r.comments, reviewComment{
			ID:          "RC_reply",
			URL:         "https://github.com/owner/repo/pull/7#discussion-reply",
			Body:        body,
			AuthorLogin: "agent",
			CreatedAt:   "2026-08-30T00:01:00Z",
			UpdatedAt:   "2026-08-30T00:01:00Z",
			Path:        "file.go",
			CommitOID:   testOID('b'),
		})
		if r.nullReplyResponse {
			return githubGraphQLResult(r.t, map[string]any{
				"addPullRequestReviewThreadReply": nil,
			}), nil
		}
		mutationCommentID := r.replyMutationCommentID
		if mutationCommentID == "" {
			mutationCommentID = "RC_reply"
		}
		return githubGraphQLResult(r.t, map[string]any{
			"addPullRequestReviewThreadReply": map[string]any{
				"comment": map[string]any{"id": mutationCommentID},
			},
		}), nil
	case githubResolveMutation:
		r.requireThreadID(payload.Variables)
		r.sawResolve = true
		r.mutated = true
		r.resolved = true
		return githubGraphQLResult(r.t, map[string]any{
			"resolveReviewThread": map[string]any{
				"thread": map[string]any{
					"id": r.threadID, "isResolved": true,
				},
			},
		}), nil
	default:
		r.t.Fatalf("unexpected GraphQL query: %q", payload.Query)
		return commandResult{}, nil
	}
}

func (r *reviewLifecycleRunner) requireThreadID(variables map[string]any) {
	r.t.Helper()
	threadID, ok := variables["threadId"].(string)
	if !ok || threadID != r.threadID {
		r.t.Fatalf("GraphQL threadId = %#v, want %q", variables["threadId"], r.threadID)
	}
}

func (r *reviewLifecycleRunner) commentOnPullRequest(
	request command,
) (commandResult, error) {
	wantArgs := []string{
		"api", "--hostname", "github.com", "--method", "POST",
		"repos/owner/repo/issues/7/comments", "--input", "-",
	}
	if !reflect.DeepEqual(request.Args, wantArgs) {
		r.t.Fatalf("top-level comment args = %#v, want %#v", request.Args, wantArgs)
	}
	var body struct {
		Body string `json:"body"`
	}
	if err := json.Unmarshal([]byte(request.Stdin), &body); err != nil {
		r.t.Fatalf("decode top-level comment: %v", err)
	}
	if !hasFinalLine(body.Body, commentDisclaimer) {
		r.t.Fatalf("top-level comment lacks disclaimer: %q", body.Body)
	}
	r.sawComment = true
	r.mutated = true
	id := fmt.Sprintf("IC%d", len(r.topComments)+1)
	r.topComments = append(r.topComments, reviewComment{
		ID:          id,
		URL:         "https://github.com/owner/repo/pull/7#issuecomment-new",
		Body:        body.Body,
		AuthorLogin: "agent",
		CreatedAt:   "2026-08-30T00:02:00Z",
		UpdatedAt:   "2026-08-30T00:02:00Z",
	})
	if r.commentErr != nil {
		return commandResult{}, r.commentErr
	}
	if r.commentResponse != "" {
		return commandResult{Stdout: r.commentResponse}, nil
	}
	encoded, err := json.Marshal(map[string]any{"node_id": id})
	if err != nil {
		r.t.Fatal(err)
	}
	return commandResult{Stdout: string(encoded)}, nil
}

func (r *reviewLifecycleRunner) threadMap() map[string]any {
	line := 12
	return map[string]any{
		"id": r.threadID, "isResolved": r.resolved, "isOutdated": false,
		"path": "file.go", "line": line, "originalLine": line,
		"startLine": nil, "originalStartLine": nil,
	}
}

func (r *reviewLifecycleRunner) thread() reviewThread {
	line := 12
	thread := reviewThread{
		ID:           r.threadID,
		IsResolved:   r.resolved,
		Path:         "file.go",
		Line:         &line,
		OriginalLine: &line,
		Comments:     append([]reviewComment(nil), r.comments...),
	}
	thread.ExpectationDigest = reviewThreadDigest(thread)
	return thread
}

func (r *reviewLifecycleRunner) threadExpectation() reviewThreadExpectation {
	thread := r.thread()
	return reviewThreadExpectation{
		ThreadID:              thread.ID,
		ExpectedLastCommentID: thread.Comments[len(thread.Comments)-1].ID,
		ExpectedDigest:        thread.ExpectationDigest,
		RequiredReplyBodyDigest: reviewReplyBodyDigest(
			thread.Comments[len(thread.Comments)-1].Body,
		),
		RequiredInventoryDigest: reviewInventoryDigest(r.inspection()),
		AuthorityExpiresAt:      "2099-08-30T00:00:00Z",
	}
}

func (r *reviewLifecycleRunner) inspection() *reviewInspection {
	pullRequest := testPullRequest()
	pullRequest.UpdatedAt = r.pullRequestUpdatedAt
	requested := make([]requestedReviewer, 0, len(r.reviewers))
	for _, login := range r.reviewers {
		requested = append(requested, requestedReviewer{
			RequestID: "RR_" + login, Type: "user", ID: "U_" + login, Login: login,
		})
	}
	return &reviewInspection{
		PullRequest: pullRequest,
		PullRequestExpectationDigest: reviewPullRequestDigest(
			remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
			pullRequest,
		),
		Comments:           append([]reviewComment(nil), r.topComments...),
		Reviews:            append([]pullRequestReview(nil), r.reviews...),
		Threads:            []reviewThread{r.thread()},
		RequestedReviewers: requested,
	}
}

func (r *reviewLifecycleRunner) topLevelExpectation() topLevelCommentExpectation {
	return topLevelCommentExpectation{
		ExpectedLastCommentID: expectedLastTopLevelComment(r.topComments),
		ExpectedDigest:        reviewCommentsDigest(r.topComments),
	}
}

func (r *reviewLifecycleRunner) threadCommentsMap() map[string]any {
	result := r.threadMap()
	comments := r.comments
	if r.mutated && r.staleReplyInventoryOnce &&
		r.postMutationViewCount == 1 && len(comments) > 1 {
		comments = comments[:len(comments)-1]
	}
	nodes := make([]any, 0, len(comments))
	for _, comment := range comments {
		nodes = append(nodes, map[string]any{
			"id": comment.ID, "url": comment.URL, "body": comment.Body,
			"author":    map[string]any{"login": comment.AuthorLogin},
			"createdAt": comment.CreatedAt, "updatedAt": comment.UpdatedAt,
			"path": comment.Path, "line": comment.Line,
			"originalLine": comment.OriginalLine,
			"commit":       map[string]any{"oid": comment.CommitOID},
		})
	}
	result["comments"] = map[string]any{
		"nodes": nodes,
		"pageInfo": map[string]any{
			"hasNextPage": false,
			"endCursor":   nil,
		},
	}
	return result
}

func (r *reviewLifecycleRunner) requestReview(
	request command,
) (commandResult, error) {
	wantArgs := []string{
		"api", "--hostname", "github.com", "--method", "POST",
		"--silent",
		"repos/owner/repo/pulls/7/requested_reviewers", "--input", "-",
	}
	if !reflect.DeepEqual(request.Args, wantArgs) {
		r.t.Fatalf("review request args = %#v, want %#v", request.Args, wantArgs)
	}
	var body struct {
		Reviewers []string `json:"reviewers"`
	}
	if err := json.Unmarshal([]byte(request.Stdin), &body); err != nil {
		r.t.Fatalf("decode review request: %v", err)
	}
	r.sawRequest = true
	r.mutated = true
	if r.completeReviewImmediately {
		for _, reviewer := range body.Reviewers {
			r.reviews = append(r.reviews, pullRequestReview{
				ID:          "R_immediate_" + reviewer,
				URL:         "https://github.com/owner/repo/pull/7#review-immediate",
				Body:        "Reviewed immediately.",
				AuthorLogin: reviewer,
				State:       "APPROVED",
				SubmittedAt: "2026-08-30T00:03:00Z",
				CommitOID:   testOID('b'),
			})
		}
	} else {
		r.reviewers = append(r.reviewers, body.Reviewers...)
	}
	return commandResult{}, nil
}

func TestGitHubInspectReviewsReadsEveryReviewCategory(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.reviewers = []string{"alice"}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.InspectReviews(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
	)
	if err != nil {
		t.Fatalf("InspectReviews() error = %v", err)
	}
	if len(inspection.Comments) != 1 || len(inspection.Reviews) != 1 ||
		len(inspection.Threads) != 1 ||
		len(inspection.Threads[0].Comments) != 1 ||
		len(inspection.RequestedReviewers) != 1 ||
		inspection.ExpectedLastTopLevelComment != "IC1" ||
		inspection.ExpectedTopLevelCommentsDigest !=
			reviewCommentsDigest(inspection.Comments) ||
		inspection.Threads[0].ExpectationDigest == "" {
		t.Fatalf("InspectReviews() = %#v", inspection)
	}
}

func TestGitHubInspectReviewsSharesRawByteBudgetAcrossEveryRead(t *testing.T) {
	t.Parallel()
	lifecycle := newReviewLifecycleRunner(t)
	runner := &aggregateReviewInventoryRunner{
		t:        t,
		delegate: lifecycle,
		padding:  strings.Repeat(" ", githubInventoryByteLimit/7),
	}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.InspectReviews(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
	)
	if err == nil || !strings.Contains(
		err.Error(),
		"response bytes while reading pull-request coherence reads",
	) {
		t.Fatalf("InspectReviews() error = %v", err)
	}

	wantNames := []string{
		"pull-request coherence reads",
		"comments",
		"reviews",
		"review threads",
		"review thread comments",
		"review requests",
		"pull-request coherence reads",
	}
	gotNames := make([]string, 0, len(runner.responses))
	totalBytes := 0
	for _, response := range runner.responses {
		gotNames = append(gotNames, response.Name)
		totalBytes += response.Bytes
		if response.Bytes >= githubInventoryByteLimit {
			t.Fatalf(
				"individual %s response has %d bytes, limit %d",
				response.Name,
				response.Bytes,
				githubInventoryByteLimit,
			)
		}
	}
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("inventory reads = %#v, want %#v", gotNames, wantNames)
	}
	if totalBytes <= githubInventoryByteLimit {
		t.Fatalf(
			"aggregate response bytes = %d, want more than %d",
			totalBytes,
			githubInventoryByteLimit,
		)
	}
	for _, response := range runner.responses {
		if totalBytes-response.Bytes > githubInventoryByteLimit {
			t.Fatalf(
				"fixture exceeds the budget without %s: %d > %d",
				response.Name,
				totalBytes-response.Bytes,
				githubInventoryByteLimit,
			)
		}
	}
}

func TestGitHubInspectReviewsRefusesIncoherentUpdatedAt(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	changed := githubGraphQLPullRequestMap()
	changed["updatedAt"] = "2026-08-30T00:01:00Z"
	changed["comments"] = map[string]any{
		"nodes":    []any{},
		"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
	}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubCommand([]string{
				"pr", "view", "7", "--repo", "github.com/owner/repo",
				"--json", githubPullRequestFields,
			}, ""),
			result: commandResult{Stdout: githubPRJSON("Title", "Body")},
		},
		{
			command: githubGraphQLTestCommand(
				t, repository, githubCommentsQuery,
				githubVariables(repository, pullRequest, ""),
			),
			result: githubGraphQLResult(t, map[string]any{
				"repository": map[string]any{"pullRequest": changed},
			}),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.InspectReviews(context.Background(), repository, pullRequest)
	if err == nil || !strings.Contains(err.Error(), "epoch advanced") {
		t.Fatalf("InspectReviews() error = %v", err)
	}
	runner.done()
}

func TestGitHubTopLevelCommentUsesStdinAndVerifies(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.CommentOnPullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.topLevelExpectation(),
		"Explained why the top-level feedback does not apply.",
	)
	if err != nil {
		t.Fatalf("CommentOnPullRequest() error = %v", err)
	}
	last := inspection.Comments[len(inspection.Comments)-1]
	if !runner.sawComment || len(inspection.Comments) != 2 ||
		last.ID != "IC2" || !hasFinalLine(last.Body, commentDisclaimer) ||
		inspection.ExpectedLastTopLevelComment != "IC2" {
		t.Fatalf("CommentOnPullRequest() = %#v", inspection)
	}
}

func TestGitHubTopLevelCommentAcceptsNoCommentSentinel(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.topComments = nil
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.CommentOnPullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.topLevelExpectation(),
		"Recorded the review limitation.",
	)
	if err != nil {
		t.Fatalf("CommentOnPullRequest() error = %v", err)
	}
	if !runner.sawComment || len(inspection.Comments) != 1 {
		t.Fatalf("CommentOnPullRequest() = %#v", inspection)
	}
}

func TestGitHubTopLevelCommentRefusesStaleGuard(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	expectation := runner.topLevelExpectation()
	expectation.ExpectedLastCommentID = "stale"
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.CommentOnPullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		expectation,
		"Recorded the review limitation.",
	)
	if err == nil || !strings.Contains(err.Error(), "changed after inspection") ||
		runner.sawComment {
		t.Fatalf("CommentOnPullRequest() error = %v, saw = %v", err, runner.sawComment)
	}
}

func TestGitHubTopLevelCommentRefusesSameIDEdit(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	expectation := runner.topLevelExpectation()
	sensitiveEdit := "same ID, newly edited sensitive feedback"
	runner.topComments[0].Body = sensitiveEdit
	runner.topComments[0].UpdatedAt = "2026-08-30T00:01:00Z"
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.CommentOnPullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		expectation,
		"Recorded the review limitation.",
	)
	if err == nil || !strings.Contains(err.Error(), "changed after inspection") ||
		strings.Contains(err.Error(), sensitiveEdit) || runner.sawComment {
		t.Fatalf("CommentOnPullRequest() error = %v, saw = %v", err, runner.sawComment)
	}
}

func TestGitHubTopLevelCommentReconcilesAmbiguousFailure(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.commentErr = errors.New("connection closed after request")
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.CommentOnPullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.topLevelExpectation(),
		"Recorded the review limitation.",
	)
	if err != nil {
		t.Fatalf("CommentOnPullRequest() error = %v", err)
	}
	if !runner.sawComment || len(inspection.Comments) != 2 {
		t.Fatalf("CommentOnPullRequest() = %#v", inspection)
	}
}

func TestGitHubTopLevelCommentRejectsMalformedMutationResponse(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.commentResponse = `{"Node_ID":"IC2"}`
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.CommentOnPullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.topLevelExpectation(),
		"Recorded the review limitation.",
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "malformed") {
		t.Fatalf("CommentOnPullRequest() error = %v", err)
	}
}

func TestGitHubTopLevelCommentAllowsDocumentedRESTResponseEvolution(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.commentResponse = `{
		"node_id":"IC2",
		"url":"https://api.github.com/repos/owner/repo/issues/comments/2",
		"user":{"login":"agent","future_field":true}
	}`
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.CommentOnPullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.topLevelExpectation(),
		"Recorded the review limitation.",
	)
	if err != nil || inspection.Comments[len(inspection.Comments)-1].ID != "IC2" {
		t.Fatalf("CommentOnPullRequest() inspection = %#v, error = %v", inspection, err)
	}
}

func TestVerifyGitHubTopLevelCommentDetectsEditedPrefix(t *testing.T) {
	t.Parallel()
	before := &reviewInspection{Comments: []reviewComment{{
		ID: "IC1", Body: "Original", UpdatedAt: "before",
	}}}
	after := &reviewInspection{Comments: []reviewComment{
		{ID: "IC1", Body: "Edited", UpdatedAt: "after"},
		{ID: "IC2", Body: "Reply"},
	}}
	err := verifyGithubTopLevelComment(before, after, "Reply", "IC2")
	if err == nil || !strings.Contains(err.Error(), "concurrent") {
		t.Fatalf("verifyGithubTopLevelComment() error = %v", err)
	}
}

func TestGitHubReplyUsesStdinAddsDisclaimerAndVerifies(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.threadID = "RT_selected"
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
		"Implemented the requested fix.",
	)
	if err != nil {
		t.Fatalf("ReplyToReviewThread() error = %v", err)
	}
	last := inspection.Threads[0].Comments[1]
	if !runner.sawReply || last.ID != "RC_reply" ||
		!hasFinalLine(last.Body, commentDisclaimer) {
		t.Fatalf("ReplyToReviewThread() = %#v", inspection)
	}
}

func TestGitHubReplyRetriesWholeInventoryWhenUpdatedAtAdvances(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.advanceEpochAfterReply = true
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
		"Implemented the requested fix.",
	)
	if err != nil {
		t.Fatalf("ReplyToReviewThread() error = %v", err)
	}
	last := inspection.Threads[0].Comments[len(inspection.Threads[0].Comments)-1]
	if runner.replyMutationCount != 1 || last.ID != "RC_reply" ||
		inspection.PullRequest.UpdatedAt != "2026-08-30T00:01:00Z" ||
		len(inspection.Reviews) != 2 {
		t.Fatalf(
			"ReplyToReviewThread() = %#v, mutation count = %d",
			inspection,
			runner.replyMutationCount,
		)
	}
}

func TestGitHubReplyRejectsEpochRegressionAcrossRetryAttempts(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.advanceEpochAfterReply = true
	runner.postMutationViewEpochs = []string{
		"2026-08-30T00:00:00Z",
		"2026-08-30T00:00:00Z",
	}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
		"Implemented the requested fix.",
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "backward") ||
		runner.replyMutationCount != 1 {
		t.Fatalf(
			"ReplyToReviewThread() error = %v, mutation count = %d",
			err,
			runner.replyMutationCount,
		)
	}
}

func TestGitHubReplyRetriesCoherentStaleInventoryWithoutRemutating(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.staleReplyInventoryOnce = true
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
		"Implemented the requested fix.",
	)
	if err != nil {
		t.Fatalf("ReplyToReviewThread() error = %v", err)
	}
	last := inspection.Threads[0].Comments[len(inspection.Threads[0].Comments)-1]
	if runner.replyMutationCount != 1 || last.ID != "RC_reply" {
		t.Fatalf(
			"ReplyToReviewThread() = %#v, mutation count = %d",
			inspection,
			runner.replyMutationCount,
		)
	}
}

func TestGitHubInspectReviewsAcceptsEqualOrNewerInitialEpoch(t *testing.T) {
	t.Parallel()
	for name, updatedAt := range map[string]string{
		"equal": "2026-08-30T00:00:00Z",
		"newer": "2026-08-30T00:01:00Z",
	} {
		name, updatedAt := name, updatedAt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := newReviewLifecycleRunner(t)
			runner.pullRequestUpdatedAt = updatedAt
			forge := &githubForge{
				executable: "gh-test", runner: runner, directory: "/repo",
			}
			if _, err := forge.InspectReviews(
				context.Background(),
				remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
				testPullRequest(),
			); err != nil {
				t.Fatalf("InspectReviews() error = %v", err)
			}
		})
	}
}

func TestGitHubInspectReviewsRejectsBackwardInitialEpoch(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	expected := testPullRequest()
	expected.UpdatedAt = "2026-08-30T00:01:00Z"
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	if _, err := forge.InspectReviews(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		expected,
	); err == nil || !strings.Contains(err.Error(), "backward") {
		t.Fatalf("InspectReviews() error = %v", err)
	}
}

func TestGitHubReplyRequiresMutationReturnedCommentID(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.nullReplyResponse = true
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
		"Implemented the requested fix.",
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "comment ID") {
		t.Fatalf("ReplyToReviewThread() error = %v", err)
	}
}

func TestGitHubReplyRequiresMutationIDToMatchAppendedComment(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.replyMutationCommentID = "RC_different"
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
		"Implemented the requested fix.",
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "different created comment") {
		t.Fatalf("ReplyToReviewThread() error = %v", err)
	}
}

func TestGitHubReplyReportsUnknownOutcomeWhenRereadFails(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.postMutationViewErr = errors.New("connection closed")
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
		"Implemented the requested fix.",
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "re-inspect before retrying") {
		t.Fatalf("ReplyToReviewThread() error = %v", err)
	}
}

func TestGitHubReplyRefusesStaleLastComment(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	expectation := runner.threadExpectation()
	expectation.ExpectedLastCommentID = "stale"
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		expectation,
		"Implemented the requested fix.",
	)
	if err == nil || !strings.Contains(err.Error(), "changed after inspection") ||
		runner.sawReply {
		t.Fatalf("ReplyToReviewThread() error = %v, sawReply = %v", err, runner.sawReply)
	}
}

func TestGitHubReplyRefusesSameIDEdit(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	expectation := runner.threadExpectation()
	sensitiveEdit := "same ID, newly edited sensitive thread feedback"
	runner.comments[0].Body = sensitiveEdit
	runner.comments[0].UpdatedAt = "2026-08-30T00:01:00Z"
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ReplyToReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		expectation,
		"Implemented the requested fix.",
	)
	if err == nil || !strings.Contains(err.Error(), "changed after inspection") ||
		strings.Contains(err.Error(), sensitiveEdit) || runner.sawReply {
		t.Fatalf("ReplyToReviewThread() error = %v, saw = %v", err, runner.sawReply)
	}
}

func TestVerifyGitHubReplyDetectsEditedPrefix(t *testing.T) {
	t.Parallel()
	before := &reviewThread{Comments: []reviewComment{{
		ID: "RC1", Body: "Original", UpdatedAt: "before",
	}}}
	after := &reviewThread{Comments: []reviewComment{
		{ID: "RC1", Body: "Edited", UpdatedAt: "after"},
		{ID: "RC2", Body: "Reply"},
	}}
	err := verifyGithubReply(before, after, "Reply", "RC2")
	if err == nil || !strings.Contains(err.Error(), "concurrent") {
		t.Fatalf("verifyGithubReply() error = %v", err)
	}
}

func TestVerifyGitHubReplyDetectsThreadMetadataChange(t *testing.T) {
	t.Parallel()
	before := &reviewThread{
		ID: "RT1", Path: "before.go",
		Comments: []reviewComment{{ID: "RC1", Body: "Original"}},
	}
	after := &reviewThread{
		ID: "RT1", Path: "after.go",
		Comments: []reviewComment{
			{ID: "RC1", Body: "Original"},
			{ID: "RC2", Body: "Reply"},
		},
	}
	if err := verifyGithubReply(before, after, "Reply", "RC2"); err == nil ||
		!strings.Contains(err.Error(), "metadata") {
		t.Fatalf("verifyGithubReply() error = %v", err)
	}
}

func TestVerifyGitHubResolutionAllowsOnlyResolvedTransition(t *testing.T) {
	t.Parallel()
	before := &reviewThread{
		ID: "RT1", Path: "file.go",
		Comments: []reviewComment{{ID: "RC1", Body: "Reply"}},
	}
	after := *before
	after.IsResolved = true
	after.Path = "other.go"
	if err := verifyGithubResolution(before, &after); err == nil ||
		!strings.Contains(err.Error(), "metadata") {
		t.Fatalf("verifyGithubResolution() error = %v", err)
	}
}

func TestGitHubResolveVerifiesUnchangedThread(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.threadID = "RT_selected"
	runner.comments = append(runner.comments, reviewComment{
		ID:          "RC_reply",
		URL:         "https://github.com/owner/repo/pull/7#discussion-reply",
		Body:        "Implemented the requested fix.\n\n" + commentDisclaimer,
		AuthorLogin: "agent",
		CreatedAt:   "2026-08-30T00:01:00Z",
		UpdatedAt:   "2026-08-30T00:01:00Z",
		Path:        "file.go",
		CommitOID:   testOID('b'),
	})
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.ResolveReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
	)
	if err != nil {
		t.Fatalf("ResolveReviewThread() error = %v", err)
	}
	if !runner.sawResolve || !inspection.Threads[0].IsResolved {
		t.Fatalf("ResolveReviewThread() = %#v", inspection)
	}
}

func TestGitHubResolveRefusesAuthorityExpiredDuringInventory(t *testing.T) {
	runner := newReviewLifecycleRunner(t)
	runner.threadID = "RT_selected"
	runner.comments = append(runner.comments, reviewComment{
		ID: "RC_reply", URL: "https://github.com/owner/repo/pull/7#reply",
		Body: "Implemented.\n\n" + commentDisclaimer, AuthorLogin: "agent",
		CreatedAt: "2026-08-30T00:01:00Z", UpdatedAt: "2026-08-30T00:01:00Z",
		Path: "file.go", CommitOID: testOID('b'),
	})
	expectation := runner.threadExpectation()
	expectation.AuthorityExpiresAt = "2026-08-30T00:05:00Z"
	ctx := withReviewReplyReceiptClock(context.Background(), func() time.Time {
		return time.Date(2026, 8, 30, 0, 5, 1, 0, time.UTC)
	})
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ResolveReviewThread(
		ctx,
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		expectation,
	)
	if err == nil || runner.sawResolve {
		t.Fatalf("ResolveReviewThread() error = %v, sawResolve = %v", err, runner.sawResolve)
	}
}

func TestGitHubResolveRefusesLateTargetAppendWithUnchangedParent(t *testing.T) {
	runner := newReviewLifecycleRunner(t)
	runner.threadID = "RT_selected"
	runner.comments = append(runner.comments, reviewComment{
		ID: "RC_reply", URL: "https://github.com/owner/repo/pull/7#reply",
		Body: "Implemented.\n\n" + commentDisclaimer, AuthorLogin: "agent",
		CreatedAt: "2026-08-30T00:01:00Z", UpdatedAt: "2026-08-30T00:01:00Z",
		Path: "file.go", CommitOID: testOID('b'),
	})
	expectation := runner.threadExpectation()
	runner.appendBeforeFinalThreadRead = true
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ResolveReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		expectation,
	)
	if err == nil || runner.sawResolve || runner.pullRequestUpdatedAt != "2026-08-30T00:00:00Z" {
		t.Fatalf("ResolveReviewThread() error = %v, sawResolve = %v, updatedAt = %q", err, runner.sawResolve, runner.pullRequestUpdatedAt)
	}
}

func TestGitHubResolveRefusesAuthorityCrossingExpiryBeforeMutation(t *testing.T) {
	runner := newReviewLifecycleRunner(t)
	runner.threadID = "RT_selected"
	runner.comments = append(runner.comments, reviewComment{
		ID: "RC_reply", URL: "https://github.com/owner/repo/pull/7#reply",
		Body: "Implemented.\n\n" + commentDisclaimer, AuthorLogin: "agent",
		CreatedAt: "2026-08-30T00:01:00Z", UpdatedAt: "2026-08-30T00:01:00Z",
		Path: "file.go", CommitOID: testOID('b'),
	})
	expectation := runner.threadExpectation()
	expectation.AuthorityExpiresAt = "2026-08-30T00:05:00Z"
	clockCalls := 0
	ctx := withReviewReplyReceiptClock(context.Background(), func() time.Time {
		clockCalls++
		if clockCalls < 3 {
			return time.Date(2026, 8, 30, 0, 4, 59, 0, time.UTC)
		}
		return time.Date(2026, 8, 30, 0, 5, 0, 0, time.UTC)
	})
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ResolveReviewThread(
		ctx,
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		expectation,
	)
	if err == nil || runner.sawResolve || clockCalls != 3 {
		t.Fatalf(
			"ResolveReviewThread() error = %v, sawResolve = %v, clockCalls = %d",
			err,
			runner.sawResolve,
			clockCalls,
		)
	}
}

func TestGitHubResolveReportsUnknownOutcomeWhenRereadFails(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.comments = append(runner.comments, reviewComment{
		ID: "RC_reply", URL: "https://github.com/owner/repo/pull/7#discussion-reply",
		Body:        "Implemented the requested fix.\n\n" + commentDisclaimer,
		AuthorLogin: "agent", CreatedAt: "2026-08-30T00:01:00Z",
		UpdatedAt: "2026-08-30T00:01:00Z", Path: "file.go",
		CommitOID: testOID('b'),
	})
	runner.postMutationViewErr = errors.New("connection closed")
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ResolveReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "re-inspect before retrying") {
		t.Fatalf("ResolveReviewThread() error = %v", err)
	}
}

func TestGitHubResolveRefusesUnansweredFeedback(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.ResolveReviewThread(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		runner.threadExpectation(),
	)
	if err == nil || !strings.Contains(err.Error(), "exact disclaimer") ||
		runner.sawResolve {
		t.Fatalf("ResolveReviewThread() error = %v, saw = %v", err, runner.sawResolve)
	}
}

func TestGitHubRequestReviewUsesStdinAndVerifies(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.RequestReview(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		[]string{"alice"},
	)
	if err != nil {
		t.Fatalf("RequestReview() error = %v", err)
	}
	if !runner.sawRequest || len(inspection.RequestedReviewers) != 1 ||
		inspection.RequestedReviewers[0].Login != "alice" {
		t.Fatalf("RequestReview() = %#v", inspection)
	}
}

func TestGitHubRequestReviewReportsUnknownOutcomeWhenRereadFails(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.postMutationViewErr = errors.New("connection closed")
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.RequestReview(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		[]string{"alice"},
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "re-inspect before retrying") {
		t.Fatalf("RequestReview() error = %v", err)
	}
}

func TestGitHubRequestReviewAcceptsImmediateReviewOnExpectedHead(t *testing.T) {
	t.Parallel()
	runner := newReviewLifecycleRunner(t)
	runner.completeReviewImmediately = true
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	inspection, err := forge.RequestReview(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		testPullRequest(),
		[]string{"alice"},
	)
	if err != nil {
		t.Fatalf("RequestReview() error = %v", err)
	}
	if !runner.sawRequest || len(inspection.RequestedReviewers) != 0 ||
		!completedGithubReviewSince(
			&reviewInspection{Reviews: []pullRequestReview{runner.reviews[0]}},
			inspection,
			"alice",
			testOID('b'),
		) {
		t.Fatalf("RequestReview() = %#v", inspection)
	}
}

func TestReviewReplyCommandRequiresRaceGuards(t *testing.T) {
	t.Parallel()
	root := newRootCommand(
		context.Background(),
		func(string) string { return "" },
		io.Discard,
		&transcriptRunner{t: t},
	)
	root.SetArgs([]string{"review", "reply"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "required flag") {
		t.Fatalf("review reply error = %v", err)
	}
}

func TestReviewCommentCommandRequiresRaceGuards(t *testing.T) {
	t.Parallel()
	root := newRootCommand(
		context.Background(),
		func(string) string { return "" },
		io.Discard,
		&transcriptRunner{t: t},
	)
	root.SetArgs([]string{"review", "comment"})
	if err := root.Execute(); err == nil || !strings.Contains(err.Error(), "required flag") {
		t.Fatalf("review comment error = %v", err)
	}
}

func TestReviewMutationCommandsRequireInventoryDigests(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		missing string
	}{
		{
			name: "top-level comments",
			args: []string{
				"review", "comment",
				"--pull-request-id", "PR1",
				"--expected-head", testOID('b'),
				"--expected-pull-request-digest", strings.Repeat("a", 64),
				"--expected-last-top-level-comment", "IC1",
				"--body-file", "out/review/comment.md",
			},
			missing: "expected-top-level-comments-digest",
		},
		{
			name: "thread reply",
			args: []string{
				"review", "reply",
				"--pull-request-id", "PR1",
				"--expected-head", testOID('b'),
				"--expected-pull-request-digest", strings.Repeat("a", 64),
				"--thread-id", "RT1",
				"--expected-last-comment-id", "RC1",
				"--body-file", "out/review/reply.md",
				"--reply-receipt-file", "out/review/reply.json",
			},
			missing: "expected-thread-digest",
		},
		{
			name: "thread resolution",
			args: []string{
				"review", "resolve",
				"--pull-request-id", "PR1",
				"--expected-head", testOID('b'),
				"--expected-pull-request-digest", strings.Repeat("a", 64),
				"--thread-id", "RT1",
				"--expected-last-comment-id", "RC1",
				"--reply-receipt-file", "out/review/reply.json",
			},
			missing: "expected-thread-digest",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			root := newRootCommand(
				context.Background(),
				func(string) string { return "" },
				io.Discard,
				&transcriptRunner{t: t},
			)
			root.SetArgs(test.args)
			err := root.Execute()
			if err == nil || !strings.Contains(err.Error(), test.missing) {
				t.Fatalf("command error = %v", err)
			}
		})
	}
}

func providerGitCommand(args ...string) command {
	return command{
		Name:             "git-test",
		Args:             hardenedGitArguments(args),
		Dir:              "/repo",
		Env:              hardenedGitEnvironment(nil),
		UnsetEnv:         gitUnsetEnvironment,
		UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes,
	}
}

func missingProviderGitConfig(args ...string) expectedCommand {
	result := commandResult{ExitCode: 1}
	return expectedCommand{
		command: providerGitCommand(args...),
		result:  result,
		err: &commandError{
			Command: command{Name: "git-test"},
			Result:  result,
			Err:     errors.New("missing config value"),
		},
	}
}

func TestProviderCommandSanitizesRemoteWithoutForgePreflight(t *testing.T) {
	t.Parallel()
	secret := "credential-that-must-not-appear"
	remoteURL := "https://user:" + secret + "@github.com/owner/repo.git"
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: providerGitCommand("rev-parse", "--show-toplevel"),
			result:  commandResult{Stdout: "/repo\n"},
		},
		{
			command: providerGitCommand("rev-parse", "--is-shallow-repository"),
			result:  commandResult{Stdout: "false\n"},
		},
		{
			command: providerGitCommand("rev-parse", "--git-common-dir"),
			result:  commandResult{Stdout: "/repo/.git\n"},
		},
		missingProviderGitConfig(
			"config", "--local", "--no-includes", "--type=bool",
			"--get", "extensions.worktreeConfig",
		),
		{
			command: providerGitCommand(
				"config", "--local", "--no-includes", "--null",
				"--get-all", "remote.origin.url",
			),
			result: commandResult{Stdout: remoteURL + "\x00"},
		},
		missingProviderGitConfig(
			"config", "--local", "--no-includes", "--null",
			"--get-all", "remote.origin.pushurl",
		),
		missingProviderGitConfig(
			"config", "--local", "--no-includes", "--type=bool",
			"--get", "extensions.worktreeConfig",
		),
		{
			command: providerGitCommand(
				"config", "--local", "--no-includes", "--null",
				"--get-all", "remote.origin.url",
			),
			result: commandResult{Stdout: remoteURL + "\x00"},
		},
		missingProviderGitConfig(
			"config", "--local", "--no-includes", "--null",
			"--get-all", "remote.origin.pushurl",
		),
	}}
	var stdout bytes.Buffer
	root := newRootCommand(
		context.Background(),
		func(string) string { return "" },
		&stdout,
		runner,
	)
	root.SetArgs([]string{
		"--repository", "/repo",
		"--git-cli", "git-test",
		"provider",
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("provider command error = %v", err)
	}
	runner.done()
	if strings.Contains(stdout.String(), secret) ||
		strings.Contains(stdout.String(), "https://") {
		t.Fatalf("provider output exposed raw remote URL: %s", stdout.String())
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &fields); err != nil {
		t.Fatalf("decode provider output: %v", err)
	}
	if len(fields) != 5 || fields["remote_repository"] == nil ||
		fields["provider_hint"] == nil || fields["adapter_available"] == nil ||
		fields["git_transport"] == nil ||
		fields["delivery_transport_available"] == nil {
		t.Fatalf("provider fields = %#v", fields)
	}
	var report providerReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode provider report: %v", err)
	}
	if report.RemoteRepository != (remoteRepository{
		Host: "github.com", Owner: "owner", Name: "repo",
	}) || report.ProviderHint != "github" || !report.AdapterAvailable ||
		report.GitTransport != "https" || report.DeliveryTransportAvailable {
		t.Fatalf("provider report = %#v", report)
	}
}

func TestProviderHintsOnlyKnownHosts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		host      string
		provider  string
		available bool
	}{
		{host: "github.com", provider: "github", available: true},
		{host: "codeberg.org", provider: "forgejo", available: false},
		{
			host: "forgejo.alwaldend.com", provider: "forgejo",
			available: false,
		},
		{host: "git.example", provider: "unknown", available: false},
	}
	for _, test := range tests {
		provider, available := providerForRepository(remoteRepository{
			Host: test.host,
		})
		if provider != test.provider || available != test.available {
			t.Errorf(
				"providerForRepository(%q) = %q, %v, want %q, %v",
				test.host,
				provider,
				available,
				test.provider,
				test.available,
			)
		}
	}
}

func TestParseRemoteRepositoryRejectsInsecureTransports(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{
		"http://github.com/owner/repo.git",
		"HTTP://github.com/owner/repo.git",
		"git://github.com/owner/repo.git",
		"GIT://github.com/owner/repo.git",
	} {
		_, err := parseRemoteRepository(raw)
		if err == nil || !strings.Contains(err.Error(), "insecure transport") {
			t.Errorf("parseRemoteRepository(%q) error = %v", raw, err)
		}
	}
}

func TestParseRemoteRepositoryAcceptsSecureTransports(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{
		"https://github.com/owner/repo.git",
		"ssh://git@github.com/owner/repo.git",
		"git@github.com:owner/repo.git",
	} {
		got, err := parseRemoteRepository(raw)
		if err != nil {
			t.Errorf("parseRemoteRepository(%q) error = %v", raw, err)
			continue
		}
		if got.Owner != "owner" || got.Name != "repo" {
			t.Errorf("parseRemoteRepository(%q) = %#v", raw, got)
		}
	}
}

func TestParseRemoteRepositoryDoesNotEchoMalformedCredentials(t *testing.T) {
	t.Parallel()
	secret := "credential-that-must-not-appear"
	for _, raw := range []string{
		"https://user:" + secret + "@github.com:%zz/owner/repo.git",
		"https://user:" + secret + "@github.com:bad/owner/repo.git",
	} {
		_, err := parseRemoteRepository(raw)
		if err == nil || strings.Contains(err.Error(), secret) ||
			strings.Contains(err.Error(), raw) {
			t.Errorf("parseRemoteRepository() error = %v", err)
		}
	}
}

func TestParseRemoteRepositoryRejectsNonCanonicalSCPPaths(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{
		"git@github.com:/owner/repo.git",
		"git@github.com:owner//repo.git",
		"git@github.com:owner/../repo.git",
	} {
		if _, err := parseRemoteRepository(raw); err == nil {
			t.Errorf("parseRemoteRepository(%q) unexpectedly succeeded", raw)
		}
	}
}

func TestReviewPullRequestDigestBindsEditableContext(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	original := testPullRequest()
	want := reviewPullRequestDigest(repository, original)
	changes := []pullRequest{original, original, original, original}
	changes[0].Title = "Human edit"
	changes[1].BaseRefOID = testOID('c')
	changes[2].IsDraft = true
	changes[3].HeadRefOID = testOID('c')
	for _, changed := range changes {
		if got := reviewPullRequestDigest(repository, changed); got == want {
			t.Errorf("reviewPullRequestDigest() did not bind change: %#v", changed)
		}
	}
}

func TestReviewReplyReceiptJSONRejectsDuplicateAndNonCanonicalKeys(t *testing.T) {
	t.Parallel()
	for _, input := range []string{
		`{"schema":"one","schema":"two"}`,
		`{"Schema":"one"}`,
	} {
		if err := validateReviewReplyReceiptJSON([]byte(input)); err == nil {
			t.Errorf("validateReviewReplyReceiptJSON(%s) unexpectedly succeeded", input)
		}
	}
}

type reviewReceiptForge struct {
	repository         remoteRepository
	pullRequest        pullRequest
	thread             reviewThread
	comments           []reviewComment
	reviews            []pullRequestReview
	otherThreads       []reviewThread
	requestedReviewers []requestedReviewer
	sawResolve         bool
}

func (f *reviewReceiptForge) Name() string { return "github" }

func (f *reviewReceiptForge) Preflight(
	context.Context,
	remoteRepository,
) error {
	return nil
}

func (f *reviewReceiptForge) PullRequests(
	_ context.Context,
	repository remoteRepository,
	head string,
) ([]pullRequest, error) {
	if !sameRemoteRepository(repository, f.repository) ||
		head != f.pullRequest.HeadRefName {
		return nil, fmt.Errorf("unexpected review receipt target")
	}
	return []pullRequest{f.pullRequest}, nil
}

func (f *reviewReceiptForge) CreatePullRequest(
	context.Context,
	remoteRepository,
	pullRequestInput,
) (*pullRequest, error) {
	return nil, fmt.Errorf("unexpected create")
}

func (f *reviewReceiptForge) UpdatePullRequest(
	context.Context,
	remoteRepository,
	pullRequest,
	pullRequestInput,
) (*pullRequest, error) {
	return nil, fmt.Errorf("unexpected update")
}

func (f *reviewReceiptForge) inspection() *reviewInspection {
	thread := f.thread
	thread.Comments = append([]reviewComment(nil), f.thread.Comments...)
	thread.ExpectationDigest = reviewThreadDigest(thread)
	threads := []reviewThread{thread}
	threads = append(threads, f.otherThreads...)
	return &reviewInspection{
		PullRequest:                    f.pullRequest,
		PullRequestExpectationDigest:   reviewPullRequestDigest(f.repository, f.pullRequest),
		Comments:                       append([]reviewComment(nil), f.comments...),
		ExpectedLastTopLevelComment:    expectedLastTopLevelComment(f.comments),
		ExpectedTopLevelCommentsDigest: reviewCommentsDigest(f.comments),
		Reviews:                        append([]pullRequestReview(nil), f.reviews...),
		Threads:                        threads,
		RequestedReviewers:             append([]requestedReviewer(nil), f.requestedReviewers...),
	}
}

func (f *reviewReceiptForge) InspectReviews(
	_ context.Context,
	repository remoteRepository,
	expected pullRequest,
) (*reviewInspection, error) {
	if !sameRemoteRepository(repository, f.repository) ||
		!sameGithubReviewPullRequestContext(expected, f.pullRequest) {
		return nil, fmt.Errorf("unexpected review inspection target")
	}
	return f.inspection(), nil
}

func (f *reviewReceiptForge) ReplyToReviewThread(
	_ context.Context,
	repository remoteRepository,
	expected pullRequest,
	expectation reviewThreadExpectation,
	body string,
) (*reviewInspection, error) {
	if !sameRemoteRepository(repository, f.repository) ||
		!sameGithubReviewPullRequestContext(expected, f.pullRequest) ||
		expectation.ThreadID != f.thread.ID || f.thread.IsResolved ||
		len(f.thread.Comments) == 0 ||
		expectation.ExpectedLastCommentID !=
			f.thread.Comments[len(f.thread.Comments)-1].ID ||
		expectation.ExpectedDigest != reviewThreadDigest(f.thread) {
		return nil, fmt.Errorf("reply expectation changed")
	}
	f.thread.Comments = append(f.thread.Comments, reviewComment{
		ID:          "RC_tool_reply",
		URL:         "https://github.com/owner/repo/pull/7#discussion-tool-reply",
		Body:        body,
		AuthorLogin: "delivery-bot",
		CreatedAt:   "2026-08-30T00:01:00Z",
		UpdatedAt:   "2026-08-30T00:01:00Z",
		Path:        "file.go",
		CommitOID:   f.pullRequest.HeadRefOID,
	})
	f.pullRequest.UpdatedAt = "2026-08-30T00:01:00Z"
	f.thread.ExpectationDigest = reviewThreadDigest(f.thread)
	return f.inspection(), nil
}

func (f *reviewReceiptForge) CommentOnPullRequest(
	context.Context,
	remoteRepository,
	pullRequest,
	topLevelCommentExpectation,
	string,
) (*reviewInspection, error) {
	return nil, fmt.Errorf("unexpected comment")
}

func (f *reviewReceiptForge) ResolveReviewThread(
	_ context.Context,
	repository remoteRepository,
	expected pullRequest,
	expectation reviewThreadExpectation,
) (*reviewInspection, error) {
	if !sameRemoteRepository(repository, f.repository) ||
		!sameGithubReviewPullRequestContext(expected, f.pullRequest) ||
		expectation.ThreadID != f.thread.ID || f.thread.IsResolved ||
		len(f.thread.Comments) == 0 ||
		expectation.ExpectedLastCommentID !=
			f.thread.Comments[len(f.thread.Comments)-1].ID ||
		expectation.ExpectedDigest != reviewThreadDigest(f.thread) ||
		expectation.RequiredReplyBodyDigest != reviewReplyBodyDigest(
			f.thread.Comments[len(f.thread.Comments)-1].Body,
		) || expectation.RequiredInventoryDigest !=
		reviewInventoryDigest(f.inspection()) {
		return nil, fmt.Errorf("resolution expectation changed")
	}
	f.sawResolve = true
	f.thread.IsResolved = true
	f.pullRequest.UpdatedAt = "2026-08-30T00:02:00Z"
	f.thread.ExpectationDigest = reviewThreadDigest(f.thread)
	return f.inspection(), nil
}

func (f *reviewReceiptForge) RequestReview(
	context.Context,
	remoteRepository,
	pullRequest,
	[]string,
) (*reviewInspection, error) {
	return nil, fmt.Errorf("unexpected review request")
}

type reviewReceiptFixture struct {
	delivery    *delivery
	forge       *reviewReceiptForge
	ctx         context.Context
	clock       *reviewReceiptTestClock
	options     reviewTargetOptions
	expectation reviewThreadExpectation
	receipt     reviewReplyReceipt
	receiptFile string
	absolute    string
}

type reviewReceiptTestClock struct {
	now time.Time
}

func (c *reviewReceiptTestClock) time() time.Time { return c.now }

func newReviewReceiptFixture(t *testing.T) reviewReceiptFixture {
	t.Helper()
	clock := &reviewReceiptTestClock{
		now: time.Date(2026, time.August, 30, 4, 0, 0, 0, time.UTC),
	}
	ctx := withReviewReplyReceiptClock(context.Background(), clock.time)
	root := t.TempDir()
	repositoryPath := filepath.Join(root, "repository")
	runTestGit(t, root, "init", "--initial-branch=feature", repositoryPath)
	configureTestRepository(t, repositoryPath)
	writeTestFile(t, filepath.Join(repositoryPath, ".gitignore"), "out/\n")
	runTestGit(t, repositoryPath, "add", ".gitignore")
	runTestGit(t, repositoryPath, "commit", "-m", "test repository")
	receiptDirectory := filepath.Join(repositoryPath, "out", "review")
	if err := os.MkdirAll(receiptDirectory, 0o700); err != nil {
		t.Fatalf("create review receipt directory: %v", err)
	}

	remote := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	pullRequest := testPullRequest()
	pullRequest.HeadRefName = "feature"
	line := 12
	thread := reviewThread{
		ID:           "RT_receipt",
		Path:         "file.go",
		Line:         &line,
		OriginalLine: &line,
		Comments: []reviewComment{{
			ID:           "RC_feedback",
			URL:          "https://github.com/owner/repo/pull/7#discussion-feedback",
			Body:         "Please address this feedback.",
			AuthorLogin:  "reviewer",
			CreatedAt:    "2026-08-30T00:00:00Z",
			UpdatedAt:    "2026-08-30T00:00:00Z",
			Path:         "file.go",
			Line:         &line,
			OriginalLine: &line,
			CommitOID:    pullRequest.HeadRefOID,
		}},
	}
	thread.ExpectationDigest = reviewThreadDigest(thread)
	forge := &reviewReceiptForge{
		repository:  remote,
		pullRequest: pullRequest,
		thread:      thread,
	}
	repository := &gitRepository{
		directory:  repositoryPath,
		executable: "git",
		runner:     &execRunner{},
	}
	delivery := &delivery{
		repository:       repository,
		remoteRepository: remote,
		forge:            forge,
	}
	options := reviewTargetOptions{
		PullRequestID:             pullRequest.ID,
		ExpectedHeadOID:           pullRequest.HeadRefOID,
		ExpectedPullRequestDigest: reviewPullRequestDigest(remote, pullRequest),
	}
	priorDigest := thread.ExpectationDigest
	body, err := withCommentDisclaimer("Implemented the requested correction.")
	if err != nil {
		t.Fatal(err)
	}
	inspection, err := delivery.replyToReviewThread(
		ctx,
		options,
		reviewThreadExpectation{
			ThreadID:              thread.ID,
			ExpectedLastCommentID: thread.Comments[len(thread.Comments)-1].ID,
			ExpectedDigest:        priorDigest,
		},
		body,
	)
	if err != nil {
		t.Fatalf("replyToReviewThread() error = %v", err)
	}
	receipt, err := delivery.newReviewReplyReceipt(
		ctx,
		options,
		thread.ID,
		priorDigest,
		body,
		inspection,
	)
	if err != nil {
		t.Fatalf("newReviewReplyReceipt() error = %v", err)
	}
	receiptFile := "out/review/reply.json"
	if err := delivery.writeReviewReplyReceipt(ctx, receiptFile, receipt); err != nil {
		t.Fatalf("writeReviewReplyReceipt() error = %v", err)
	}
	resultThread := inspection.Threads[0]
	return reviewReceiptFixture{
		delivery: delivery,
		forge:    forge,
		ctx:      ctx,
		clock:    clock,
		options:  options,
		expectation: reviewThreadExpectation{
			ThreadID:              resultThread.ID,
			ExpectedLastCommentID: resultThread.Comments[len(resultThread.Comments)-1].ID,
			ExpectedDigest:        resultThread.ExpectationDigest,
		},
		receipt:     receipt,
		receiptFile: receiptFile,
		absolute:    filepath.Join(repositoryPath, filepath.FromSlash(receiptFile)),
	}
}

func TestReviewReplyReceiptFullReplyWriteReadResolveChain(t *testing.T) {
	fixture := newReviewReceiptFixture(t)
	read, err := fixture.delivery.readReviewReplyReceipt(
		context.Background(),
		fixture.receiptFile,
	)
	if err != nil {
		t.Fatalf("readReviewReplyReceipt() error = %v", err)
	}
	if !reflect.DeepEqual(read, fixture.receipt) {
		t.Fatalf("read receipt = %#v, want %#v", read, fixture.receipt)
	}
	info, err := os.Lstat(fixture.absolute)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		t.Fatalf("review reply receipt mode = %v", info.Mode())
	}
	inspection, err := fixture.delivery.resolveReviewThread(
		fixture.ctx,
		fixture.options,
		fixture.expectation,
		fixture.receiptFile,
	)
	if err != nil {
		t.Fatalf("resolveReviewThread() error = %v", err)
	}
	if !fixture.forge.sawResolve || !inspection.Threads[0].IsResolved {
		t.Fatalf("resolve result = %#v", inspection)
	}
	if _, err := os.Lstat(fixture.absolute); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("consumed review reply receipt still exists: %v", err)
	}
}

func TestReviewReplyReceiptAllowsUpdatedAtAdvanceWithIdenticalInventory(t *testing.T) {
	fixture := newReviewReceiptFixture(t)
	fixture.forge.pullRequest.UpdatedAt = "2026-08-30T00:02:00Z"
	inspection, err := fixture.delivery.resolveReviewThread(
		fixture.ctx,
		fixture.options,
		fixture.expectation,
		fixture.receiptFile,
	)
	if err != nil {
		t.Fatalf("resolveReviewThread() error = %v", err)
	}
	if !fixture.forge.sawResolve || !inspection.Threads[0].IsResolved {
		t.Fatalf("resolve result = %#v", inspection)
	}
}

func TestReviewReplyReceiptRejectsUpdatedAtRegression(t *testing.T) {
	fixture := newReviewReceiptFixture(t)
	fixture.forge.pullRequest.UpdatedAt = "2026-08-30T00:00:00Z"
	_, err := fixture.delivery.resolveReviewThread(
		fixture.ctx,
		fixture.options,
		fixture.expectation,
		fixture.receiptFile,
	)
	if err == nil || fixture.forge.sawResolve {
		t.Fatalf("resolveReviewThread() error = %v, sawResolve = %v", err, fixture.forge.sawResolve)
	}
}

func TestReviewReplyReceiptRejectsAnyReviewInventoryChange(t *testing.T) {
	for name, mutate := range map[string]func(*reviewReceiptForge){
		"top-level comment": func(forge *reviewReceiptForge) {
			forge.comments = append(forge.comments, reviewComment{ID: "IC_new", Body: "new"})
		},
		"review": func(forge *reviewReceiptForge) {
			forge.reviews = append(forge.reviews, pullRequestReview{ID: "R_new", State: "COMMENTED"})
		},
		"other thread": func(forge *reviewReceiptForge) {
			thread := reviewThread{
				ID: "RT_other", Path: "other.go",
				Comments: []reviewComment{{ID: "RC_other", Body: "new"}},
			}
			thread.ExpectationDigest = reviewThreadDigest(thread)
			forge.otherThreads = append(forge.otherThreads, thread)
		},
		"review request": func(forge *reviewReceiptForge) {
			forge.requestedReviewers = append(
				forge.requestedReviewers,
				requestedReviewer{RequestID: "RR_new", Type: "user", ID: "U_new", Login: "new"},
			)
		},
	} {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			fixture := newReviewReceiptFixture(t)
			mutate(fixture.forge)
			_, err := fixture.delivery.resolveReviewThread(
				fixture.ctx,
				fixture.options,
				fixture.expectation,
				fixture.receiptFile,
			)
			if err == nil || fixture.forge.sawResolve {
				t.Fatalf(
					"resolveReviewThread() error = %v, sawResolve = %v",
					err,
					fixture.forge.sawResolve,
				)
			}
			if _, statErr := os.Lstat(fixture.absolute); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("inventory mismatch did not consume receipt: %v", statErr)
			}
		})
	}
}

func TestReviewReplyReceiptRejectsTamperedInventoryDigestAndV3(t *testing.T) {
	t.Run("tampered digest", func(t *testing.T) {
		fixture := newReviewReceiptFixture(t)
		receipt := fixture.receipt
		receipt.ReviewInventoryDigest = strings.Repeat("a", 64)
		if err := fixture.delivery.writeReviewReplyReceipt(
			fixture.ctx,
			fixture.receiptFile,
			receipt,
		); err != nil {
			t.Fatal(err)
		}
		_, err := fixture.delivery.resolveReviewThread(
			fixture.ctx,
			fixture.options,
			fixture.expectation,
			fixture.receiptFile,
		)
		if err == nil || fixture.forge.sawResolve {
			t.Fatalf("resolveReviewThread() error = %v, sawResolve = %v", err, fixture.forge.sawResolve)
		}
		if _, statErr := os.Lstat(fixture.absolute); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("tampered digest did not consume receipt: %v", statErr)
		}
	})
	t.Run("missing inventory digest", func(t *testing.T) {
		fixture := newReviewReceiptFixture(t)
		contents, err := os.ReadFile(fixture.absolute)
		if err != nil {
			t.Fatal(err)
		}
		var document map[string]any
		if err := json.Unmarshal(contents, &document); err != nil {
			t.Fatal(err)
		}
		delete(document, "review_inventory_digest")
		contents, err = json.MarshalIndent(document, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		contents = append(contents, '\n')
		if err := os.WriteFile(fixture.absolute, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := fixture.delivery.readReviewReplyReceipt(
			fixture.ctx,
			fixture.receiptFile,
		); err == nil {
			t.Fatal("readReviewReplyReceipt() accepted missing inventory digest")
		}
	})
	t.Run("v3 schema", func(t *testing.T) {
		fixture := newReviewReceiptFixture(t)
		receipt := fixture.receipt
		receipt.Schema = "repo_delivery/review_reply/v3"
		if err := receipt.validate(); err == nil {
			t.Fatal("reviewReplyReceipt.validate() accepted v3")
		}
	})
}

func TestReviewInventoryDigestOrderingAndUniqueness(t *testing.T) {
	fixture := newReviewReceiptFixture(t)
	base := fixture.forge.inspection()
	base.Comments = []reviewComment{{ID: "IC_1"}, {ID: "IC_2"}}
	base.Reviews = []pullRequestReview{{ID: "R_1"}, {ID: "R_2"}}
	other := reviewThread{ID: "RT_other", Comments: []reviewComment{{ID: "RC_other"}}}
	other.ExpectationDigest = reviewThreadDigest(other)
	base.Threads = append(base.Threads, other)
	base.RequestedReviewers = []requestedReviewer{
		{RequestID: "RR_1", Type: "user", ID: "U_1"},
		{RequestID: "RR_2", Type: "team", ID: "T_1"},
	}
	want := reviewInventoryDigest(base)

	permuted := *base
	permuted.Threads = []reviewThread{base.Threads[1], base.Threads[0]}
	permuted.RequestedReviewers = []requestedReviewer{
		base.RequestedReviewers[1], base.RequestedReviewers[0],
	}
	if got := reviewInventoryDigest(&permuted); got != want {
		t.Fatalf("outer permutation digest = %q, want %q", got, want)
	}

	reordered := *base
	reordered.Comments = []reviewComment{base.Comments[1], base.Comments[0]}
	if got := reviewInventoryDigest(&reordered); got == want {
		t.Fatal("top-level comment reorder did not change digest")
	}
	reordered = *base
	reordered.Reviews = []pullRequestReview{base.Reviews[1], base.Reviews[0]}
	if got := reviewInventoryDigest(&reordered); got == want {
		t.Fatal("review reorder did not change digest")
	}

	duplicate := *base
	duplicate.Comments = []reviewComment{{ID: "IC_1"}, {ID: "IC_1"}}
	if err := validateReviewInventoryUnique(&duplicate); err == nil {
		t.Fatal("validateReviewInventoryUnique() accepted duplicate comments")
	}
}

func TestReviewReplyReceiptRefusesReplayAfterHumanUnresolve(t *testing.T) {
	fixture := newReviewReceiptFixture(t)
	if _, err := fixture.delivery.resolveReviewThread(
		fixture.ctx,
		fixture.options,
		fixture.expectation,
		fixture.receiptFile,
	); err != nil {
		t.Fatalf("first resolveReviewThread() error = %v", err)
	}
	// GitHub does not promise that a human unresolve advances a monotonic
	// thread epoch. Recreate the exact post-reply remote state, including its
	// UpdatedAt value, and prove that the consumed local receipt still prevents
	// an ABA replay.
	fixture.forge.thread.IsResolved = false
	fixture.forge.thread.ExpectationDigest = reviewThreadDigest(fixture.forge.thread)
	fixture.forge.pullRequest.UpdatedAt = fixture.receipt.PullRequestUpdatedAt
	fixture.forge.sawResolve = false
	_, err := fixture.delivery.resolveReviewThread(
		fixture.ctx,
		fixture.options,
		fixture.expectation,
		fixture.receiptFile,
	)
	if err == nil || fixture.forge.sawResolve {
		t.Fatalf("resolveReviewThread() error = %v, sawResolve = %v", err, fixture.forge.sawResolve)
	}
}

func TestReviewReplyReceiptRejectsExpiredFirstUse(t *testing.T) {
	fixture := newReviewReceiptFixture(t)
	// GitHub exposes no compare-and-swap or monotonic epoch for resolution.
	// Even if a human resolve and unresolve restore this exact observable state,
	// the short local authority window prevents an indefinitely stale first use.
	fixture.clock.now = fixture.clock.now.Add(reviewReplyAuthorityLimit)
	_, err := fixture.delivery.resolveReviewThread(
		fixture.ctx,
		fixture.options,
		fixture.expectation,
		fixture.receiptFile,
	)
	if err == nil || !strings.Contains(err.Error(), "expired") ||
		!strings.Contains(err.Error(), "fresh reasoned reply") ||
		fixture.forge.sawResolve {
		t.Fatalf("resolveReviewThread() error = %v, sawResolve = %v", err, fixture.forge.sawResolve)
	}
}

func TestReviewReplyReceiptAuthorityWindowIsCanonicalAndBounded(t *testing.T) {
	fixture := newReviewReceiptFixture(t)
	issuedAt, err := time.Parse(time.RFC3339Nano, fixture.receipt.IssuedAt)
	if err != nil {
		t.Fatal(err)
	}
	expiresAt, err := time.Parse(time.RFC3339Nano, fixture.receipt.ExpiresAt)
	if err != nil {
		t.Fatal(err)
	}
	if got := expiresAt.Sub(issuedAt); got != reviewReplyAuthorityLimit {
		t.Fatalf("review reply authority duration = %s", got)
	}
	for name, mutate := range map[string]func(*reviewReplyReceipt){
		"noncanonical issued_at": func(receipt *reviewReplyReceipt) {
			receipt.IssuedAt = "2026-08-30T07:00:00+03:00"
		},
		"reversed": func(receipt *reviewReplyReceipt) {
			receipt.ExpiresAt = receipt.IssuedAt
		},
		"too long": func(receipt *reviewReplyReceipt) {
			receipt.ExpiresAt = issuedAt.Add(reviewReplyAuthorityLimit + time.Nanosecond).
				Format(time.RFC3339Nano)
		},
	} {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			receipt := fixture.receipt
			mutate(&receipt)
			if err := receipt.validate(); err == nil {
				t.Fatal("reviewReplyReceipt.validate() unexpectedly succeeded")
			}
		})
	}
}

func TestReviewReplyReceiptRefusesCrossContextValues(t *testing.T) {
	for name, mutate := range map[string]func(*reviewReplyReceipt){
		"repository": func(receipt *reviewReplyReceipt) {
			receipt.RepositoryFingerprint = "sha256:" + strings.Repeat("a", 64)
		},
		"remote repository": func(receipt *reviewReplyReceipt) {
			receipt.RemoteRepository.Name = "other"
		},
		"pull request": func(receipt *reviewReplyReceipt) {
			receipt.PullRequestID = "PR_other"
		},
		"head": func(receipt *reviewReplyReceipt) {
			receipt.HeadOID = testOID('c')
		},
		"thread": func(receipt *reviewReplyReceipt) {
			receipt.ThreadID = "RT_other"
		},
	} {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			fixture := newReviewReceiptFixture(t)
			receipt := fixture.receipt
			mutate(&receipt)
			if err := fixture.delivery.writeReviewReplyReceipt(
				context.Background(), fixture.receiptFile, receipt,
			); err != nil {
				t.Fatalf("install mutated receipt: %v", err)
			}
			_, err := fixture.delivery.resolveReviewThread(
				fixture.ctx,
				fixture.options,
				fixture.expectation,
				fixture.receiptFile,
			)
			if err == nil || fixture.forge.sawResolve {
				t.Fatalf("resolveReviewThread() error = %v, sawResolve = %v", err, fixture.forge.sawResolve)
			}
		})
	}
}

func TestReviewReplyReceiptRefusesTamperingAndWrongMode(t *testing.T) {
	t.Run("tampered shape", func(t *testing.T) {
		fixture := newReviewReceiptFixture(t)
		contents, err := os.ReadFile(fixture.absolute)
		if err != nil {
			t.Fatal(err)
		}
		tampered := bytes.Replace(
			contents,
			[]byte("\n}"),
			[]byte(",\n  \"unexpected\": true\n}"),
			1,
		)
		if bytes.Equal(tampered, contents) {
			t.Fatal("test setup did not tamper with receipt")
		}
		if err := os.WriteFile(fixture.absolute, tampered, 0o600); err != nil {
			t.Fatal(err)
		}
		_, err = fixture.delivery.resolveReviewThread(
			fixture.ctx, fixture.options, fixture.expectation,
			fixture.receiptFile,
		)
		if err == nil || fixture.forge.sawResolve {
			t.Fatalf("resolveReviewThread() error = %v, sawResolve = %v", err, fixture.forge.sawResolve)
		}
	})
	t.Run("wrong mode", func(t *testing.T) {
		fixture := newReviewReceiptFixture(t)
		if err := os.Chmod(fixture.absolute, 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := fixture.delivery.resolveReviewThread(
			fixture.ctx, fixture.options, fixture.expectation,
			fixture.receiptFile,
		)
		if err == nil || fixture.forge.sawResolve {
			t.Fatalf("resolveReviewThread() error = %v, sawResolve = %v", err, fixture.forge.sawResolve)
		}
	})
}

func TestReviewReplyRequiresSubstantiveValidUTF8Body(t *testing.T) {
	t.Parallel()
	for _, body := range []string{
		"...",
		commentDisclaimer,
		string([]byte{0xff}),
	} {
		if _, err := withCommentDisclaimer(body); err == nil {
			t.Errorf("withCommentDisclaimer(%q) unexpectedly succeeded", body)
		}
	}
}

func TestReviewReplyArtifactPathsMustShareTaskDirectory(t *testing.T) {
	t.Parallel()
	delivery := &delivery{repository: &gitRepository{directory: "/repo"}}
	if err := delivery.requireSameReviewTaskDirectory(
		"out/task/reply.md",
		"/repo/out/other/reply.json",
	); err == nil || !strings.Contains(err.Error(), "share") {
		t.Fatalf("requireSameReviewTaskDirectory() error = %v", err)
	}
}

func TestCommentDisclaimerRejectsWrongArtifactType(t *testing.T) {
	t.Parallel()
	_, err := withCommentDisclaimer(fmt.Sprintf(
		"Done.\n\n%s",
		pullRequestDisclaimer,
	))
	if err == nil || !strings.Contains(err.Error(), "wrong disclaimer") {
		t.Fatalf("withCommentDisclaimer() error = %v", err)
	}
}

func TestCommentDisclaimerRejectsPaddedDisclaimer(t *testing.T) {
	t.Parallel()
	for _, body := range []string{
		"Done.\n\n " + commentDisclaimer,
		"Done.\n\n" + commentDisclaimer + " ",
		"Done.\n\n\t" + commentDisclaimer,
		"Done.\n\n" + commentDisclaimer + "\n ",
	} {
		_, err := withCommentDisclaimer(body)
		if err == nil || !strings.Contains(err.Error(), "wrong disclaimer") {
			t.Errorf("withCommentDisclaimer(%q) error = %v", body, err)
		}
	}
}

func TestGithubReviewersSatisfiedRejectsOldOrWrongHeadReview(t *testing.T) {
	t.Parallel()
	old := pullRequestReview{
		ID: "R1", AuthorLogin: "alice", State: "APPROVED",
		SubmittedAt: "2026-08-30T00:00:00Z", CommitOID: testOID('b'),
	}
	before := &reviewInspection{Reviews: []pullRequestReview{old}}
	for name, after := range map[string]*reviewInspection{
		"unchanged old review": {Reviews: []pullRequestReview{old}},
		"new wrong-head review": {Reviews: []pullRequestReview{{
			ID: "R2", AuthorLogin: "alice", State: "APPROVED",
			SubmittedAt: "2026-08-30T00:01:00Z", CommitOID: testOID('c'),
		}}},
	} {
		if githubReviewersSatisfied(
			before,
			after,
			[]string{"alice"},
			testOID('b'),
		) {
			t.Errorf("%s unexpectedly satisfied review request", name)
		}
	}
}

func TestCompletedGithubReviewSinceRejectsEditedCompletedReview(t *testing.T) {
	t.Parallel()
	beforeReview := pullRequestReview{
		ID: "R1", Body: "Original", AuthorLogin: "alice", State: "APPROVED",
		SubmittedAt: "2026-08-30T00:00:00Z", CommitOID: testOID('b'),
	}
	afterReview := beforeReview
	afterReview.Body = "Edited"
	if completedGithubReviewSince(
		&reviewInspection{Reviews: []pullRequestReview{beforeReview}},
		&reviewInspection{Reviews: []pullRequestReview{afterReview}},
		"alice",
		testOID('b'),
	) {
		t.Fatal("edited completed review unexpectedly satisfied a new request")
	}
}

func TestGithubReviewersPresentTreatsBotsConsistently(t *testing.T) {
	t.Parallel()
	inspection := &reviewInspection{RequestedReviewers: []requestedReviewer{{
		Type: "bot", ID: "B1", Login: "automation-bot",
	}}}
	if !githubReviewersPresent(inspection, []string{"automation-bot"}) {
		t.Fatal("requested bot was not recognized as present")
	}
}
