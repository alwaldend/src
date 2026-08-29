package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type expectedCommand struct {
	command command
	result  commandResult
	err     error
}

type transcriptRunner struct {
	t        *testing.T
	expected []expectedCommand
	index    int
}

func (r *transcriptRunner) Run(
	ctx context.Context,
	request command,
) (commandResult, error) {
	r.t.Helper()
	if r.index >= len(r.expected) {
		r.t.Fatalf("unexpected command: %#v", request)
	}
	want := r.expected[r.index]
	r.index++
	if !reflect.DeepEqual(request, want.command) {
		r.t.Fatalf("command %d = %#v, want %#v", r.index, request, want.command)
	}
	return want.result, want.err
}

func (r *transcriptRunner) done() {
	r.t.Helper()
	if r.index != len(r.expected) {
		r.t.Fatalf("ran %d commands, want %d", r.index, len(r.expected))
	}
}

func githubCommand(args []string, stdin string) command {
	return command{
		Name:             "gh-test",
		Args:             args,
		Dir:              "/repo",
		Env:              githubEnvironment,
		UnsetEnv:         append(githubUnsetEnvironment, gitUnsetEnvironment...),
		UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes,
		Stdin:            stdin,
		OutputLimit:      githubOutputLimit,
	}
}

func githubPRJSON(title, body string) string {
	return githubPRJSONWithBase(title, body, "master")
}

func githubPRJSONWithBase(title, body, base string) string {
	return fmt.Sprintf(`{
        "id":"PR_node",
        "number":7,
        "url":"https://github.com/owner/repo/pull/7",
        "state":"OPEN",
        "title":%q,
        "body":%q,
		"author":{"id":"U_agent","is_bot":false,"login":"agent","name":"Agent"},
		"baseRefName":%q,
        "baseRefOid":"%s",
        "headRefName":"feature;not-shell",
        "headRefOid":"%s",
		"headRepository":{"id":"R_repo","name":"repo","nameWithOwner":"owner/repo"},
		"headRepositoryOwner":{"id":"U_owner","login":"owner"},
        "isCrossRepository":false,
		"isDraft":false,
		"updatedAt":"2026-08-30T00:00:00Z"
	}`, title, body, base, testOID('a'), testOID('b'))
}

func testOID(character byte) string {
	value := make([]byte, 40)
	for index := range value {
		value[index] = character
	}
	return string(value)
}

func TestGitHubPreflightUsesExplicitHostRepository(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubCommand([]string{
			"repo", "view", "github.com/owner/repo",
			"--json", "id,nameWithOwner,url",
		}, ""),
		result: commandResult{Stdout: `{
            "id":"R_node",
            "nameWithOwner":"owner/repo",
            "url":"https://github.com/owner/repo"
        }`},
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	if err := forge.Preflight(context.Background(), repository); err != nil {
		t.Fatalf("Preflight() error = %v", err)
	}
	runner.done()
}

func TestGitHubPullRequestsUsesStructuredExactArguments(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	branch := "feature;not-shell"
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubCommand([]string{
			"pr", "list",
			"--repo", "github.com/owner/repo",
			"--state", "all",
			"--head", branch,
			"--limit", "100",
			"--json", githubPullRequestFields,
		}, ""),
		result: commandResult{Stdout: "[" + githubPRJSON("Title", "Body") + "]"},
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	got, err := forge.PullRequests(context.Background(), repository, branch)
	if err != nil {
		t.Fatalf("PullRequests() error = %v", err)
	}
	if len(got) != 1 || got[0].HeadRefName != branch {
		t.Fatalf("PullRequests() = %#v", got)
	}
	runner.done()
}

func TestGitHubUpdateRefusesConcurrentMetadataEdit(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	expected := pullRequest{
		ID: "PR_node", Number: 7, URL: "https://github.com/owner/repo/pull/7",
		State: "OPEN", Title: "Old", Body: "Old body", AuthorLogin: "agent",
		BaseRefName: "master", BaseRefOID: testOID('a'),
		HeadRefName: "feature;not-shell", HeadRefOID: testOID('b'),
		HeadRepositoryOwner: "owner", HeadRepositoryName: "repo",
	}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubCommand([]string{
			"pr", "view", "7",
			"--repo", "github.com/owner/repo",
			"--json", githubPullRequestFields,
		}, ""),
		result: commandResult{Stdout: githubPRJSON("Human edit", "Old body")},
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.UpdatePullRequest(
		context.Background(),
		repository,
		expected,
		pullRequestInput{
			BaseRefName: "master", HeadRefName: "feature;not-shell",
			ExpectedHeadOID: testOID('b'), Title: "New", Body: "New body",
		},
	)
	if err == nil {
		t.Fatal("UpdatePullRequest() unexpectedly overwrote concurrent edit")
	}
	runner.done()
}

func TestGitHubCreateReconcilesAmbiguousCommandFailure(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	input := pullRequestInput{
		BaseRefName: "master", HeadRefName: "feature;not-shell",
		ExpectedHeadOID: testOID('b'), Title: "Title", Body: "Body",
	}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubCommand([]string{
				"pr", "create",
				"--repo", "github.com/owner/repo",
				"--base", "master",
				"--head", "feature;not-shell",
				"--title", "Title",
				"--body-file", "-",
			}, "Body"),
			err: errors.New("connection closed after request"),
		},
		{
			command: githubCommand([]string{
				"pr", "list",
				"--repo", "github.com/owner/repo",
				"--state", "all",
				"--head", "feature;not-shell",
				"--limit", "100",
				"--json", githubPullRequestFields,
			}, ""),
			result: commandResult{Stdout: "[" + githubPRJSON("Title", "Body") + "]"},
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	got, err := forge.CreatePullRequest(context.Background(), repository, input)
	if err != nil {
		t.Fatalf("CreatePullRequest() error = %v", err)
	}
	if got.Number != 7 {
		t.Fatalf("CreatePullRequest() = %#v", got)
	}
	runner.done()
}

func TestGitHubUpdatePassesAndVerifiesBase(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	expected := pullRequest{
		ID: "PR_node", Number: 7, URL: "https://github.com/owner/repo/pull/7",
		State: "OPEN", Title: "Old", Body: "Old body", AuthorLogin: "agent",
		BaseRefName: "master", BaseRefOID: testOID('a'),
		HeadRefName: "feature;not-shell", HeadRefOID: testOID('b'),
		HeadRepositoryOwner: "owner", HeadRepositoryName: "repo",
	}
	input := pullRequestInput{
		BaseRefName: "release", HeadRefName: "feature;not-shell",
		ExpectedHeadOID: testOID('b'), Title: "New", Body: "New body",
	}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubCommand([]string{
				"pr", "view", "7",
				"--repo", "github.com/owner/repo",
				"--json", githubPullRequestFields,
			}, ""),
			result: commandResult{Stdout: githubPRJSON("Old", "Old body")},
		},
		{
			command: githubCommand([]string{
				"pr", "edit", "7",
				"--repo", "github.com/owner/repo",
				"--base", "release",
				"--title", "New",
				"--body-file", "-",
			}, "New body"),
		},
		{
			command: githubCommand([]string{
				"pr", "view", "7",
				"--repo", "github.com/owner/repo",
				"--json", githubPullRequestFields,
			}, ""),
			result: commandResult{
				Stdout: githubPRJSONWithBase("New", "New body", "release"),
			},
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	updated, err := forge.UpdatePullRequest(
		context.Background(),
		repository,
		expected,
		input,
	)
	if err != nil {
		t.Fatalf("UpdatePullRequest() error = %v", err)
	}
	if updated.BaseRefName != "release" {
		t.Fatalf("UpdatePullRequest().BaseRefName = %q", updated.BaseRefName)
	}
	runner.done()
}

func TestGitHubPullRequestsRejectsNullAndMalformedRecords(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	for name, output := range map[string]string{
		"top-level null":              "null",
		"null element":                "[null]",
		"malformed unrelated element": "[{}]",
	} {
		name, output := name, output
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := &transcriptRunner{t: t, expected: []expectedCommand{{
				command: githubCommand([]string{
					"pr", "list",
					"--repo", "github.com/owner/repo",
					"--state", "all",
					"--head", "feature",
					"--limit", "100",
					"--json", githubPullRequestFields,
				}, ""),
				result: commandResult{Stdout: output},
			}}}
			forge := &githubForge{
				executable: "gh-test", runner: runner, directory: "/repo",
			}
			if _, err := forge.PullRequests(
				context.Background(),
				repository,
				"feature",
			); err == nil {
				t.Fatal("PullRequests() unexpectedly accepted malformed JSON")
			}
			runner.done()
		})
	}
}

func TestDecodeGitHubJSONRejectsNoncanonicalOutput(t *testing.T) {
	t.Parallel()
	type nested struct {
		Value string `json:"value"`
	}
	type response struct {
		ID     string `json:"id"`
		Nested nested `json:"nested"`
	}
	tests := map[string][]byte{
		"invalid UTF-8":     {'{', '"', 'i', 'd', '"', ':', '"', 0xff, '"', '}'},
		"null scalar":       []byte(`{"id":null,"nested":{"value":"x"}}`),
		"missing scalar":    []byte(`{"nested":{"value":"x"}}`),
		"missing nested":    []byte(`{"id":"one"}`),
		"missing child":     []byte(`{"id":"one","nested":{}}`),
		"duplicate":         []byte(`{"id":"one","id":"two","nested":{"value":"x"}}`),
		"case variants":     []byte(`{"id":"one","ID":"two","nested":{"value":"x"}}`),
		"wrong casing":      []byte(`{"ID":"one","nested":{"value":"x"}}`),
		"unknown nested":    []byte(`{"id":"one","nested":{"value":"x","extra":true}}`),
		"trailing document": []byte(`{"id":"one","nested":{"value":"x"}} {}`),
	}
	for name, input := range tests {
		name, input := name, input
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var decoded response
			if err := decodeGitHubJSON("test response", input, &decoded); err == nil {
				t.Fatal("decodeGitHubJSON() unexpectedly accepted ambiguous output")
			}
		})
	}
}

func TestDecodeGitHubJSONRequiresNullablePointerKey(t *testing.T) {
	t.Parallel()
	type response struct {
		Required string  `json:"required"`
		Nullable *string `json:"nullable"`
		Optional string  `json:"optional,omitempty"`
	}
	var missing response
	if err := decodeGitHubJSON(
		"test response",
		[]byte(`{"required":"present"}`),
		&missing,
	); err == nil {
		t.Fatal("decodeGitHubJSON() accepted an omitted nullable pointer key")
	}

	var explicitNull response
	if err := decodeGitHubJSON(
		"test response",
		[]byte(`{"required":"present","nullable":null}`),
		&explicitNull,
	); err != nil {
		t.Fatalf("decodeGitHubJSON() rejected explicit null: %v", err)
	}
}

func TestDecodeGitHubJSONRequiresFlattenedThreadFields(t *testing.T) {
	t.Parallel()
	type response struct {
		Node *struct {
			githubReviewThread
			Comments *githubConnection[githubReviewComment] `json:"comments"`
		} `json:"node"`
	}
	withoutLine := []byte(`{
		"node": {
			"id": "RT_node",
			"isResolved": false,
			"isOutdated": false,
			"path": "file.go",
			"originalLine": null,
			"startLine": null,
			"originalStartLine": null,
			"comments": {
				"nodes": [],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}
	}`)
	var missing response
	if err := decodeGitHubJSON(
		"GitHub review thread",
		withoutLine,
		&missing,
	); err == nil {
		t.Fatal("decodeGitHubJSON() accepted an omitted flattened line key")
	}

	withNullLine := []byte(strings.Replace(
		string(withoutLine),
		`"path": "file.go",`,
		`"path": "file.go", "line": null,`,
		1,
	))
	var decoded response
	if err := decodeGitHubJSON(
		"GitHub review thread",
		withNullLine,
		&decoded,
	); err != nil {
		t.Fatalf("decodeGitHubJSON() rejected explicit nullable thread fields: %v", err)
	}
}

func TestGitHubPullRequestsStrictlyRejectsAmbiguousKeys(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	canonical := githubPRJSON("Title", "Body")
	for name, output := range map[string]string{
		"single wrong-case key": "[" + strings.Replace(canonical, `"id"`, `"ID"`, 1) + "]",
		"unknown nested key": "[" + strings.Replace(
			canonical,
			`"author":{"id":"U_agent","is_bot":false,"login":"agent","name":"Agent"}`,
			`"author":{"id":"U_agent","is_bot":false,"login":"agent","name":"Agent","secret":"unexpected"}`,
			1,
		) + "]",
		"duplicate nested key": "[" + strings.Replace(
			canonical,
			`"author":{"id":"U_agent","is_bot":false,"login":"agent","name":"Agent"}`,
			`"author":{"id":"U_agent","is_bot":false,"login":"agent","login":"other","name":"Agent"}`,
			1,
		) + "]",
	} {
		name, output := name, output
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := &transcriptRunner{t: t, expected: []expectedCommand{{
				command: githubCommand([]string{
					"pr", "list", "--repo", "github.com/owner/repo",
					"--state", "all", "--head", "feature", "--limit", "100",
					"--json", githubPullRequestFields,
				}, ""),
				result: commandResult{Stdout: output},
			}}}
			forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
			if _, err := forge.PullRequests(
				context.Background(), repository, "feature",
			); err == nil {
				t.Fatal("PullRequests() unexpectedly accepted ambiguous output")
			}
			runner.done()
		})
	}
}

func TestGitHubPullRequestsRejectsOmittedCLINestedFields(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	canonical := githubPRJSON("Title", "Body")
	for name, output := range map[string]string{
		"author is_bot": strings.Replace(
			canonical, `"is_bot":false,`, "", 1,
		),
		"repository nameWithOwner": strings.Replace(
			canonical, `,"nameWithOwner":"owner/repo"`, "", 1,
		),
		"repository owner id": strings.Replace(
			canonical, `"headRepositoryOwner":{"id":"U_owner",`,
			`"headRepositoryOwner":{`, 1,
		),
	} {
		name, output := name, output
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runner := &transcriptRunner{t: t, expected: []expectedCommand{{
				command: githubCommand([]string{
					"pr", "list", "--repo", "github.com/owner/repo",
					"--state", "all", "--head", "feature", "--limit", "100",
					"--json", githubPullRequestFields,
				}, ""),
				result: commandResult{Stdout: "[" + output + "]"},
			}}}
			forge := &githubForge{
				executable: "gh-test", runner: runner, directory: "/repo",
			}
			if _, err := forge.PullRequests(
				context.Background(), repository, "feature",
			); err == nil {
				t.Fatal("PullRequests() accepted an omitted CLI nested field")
			}
			runner.done()
		})
	}
}

func TestGitHubPullRequestsRejectsOmittedBooleanField(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	output := strings.Replace(
		githubPRJSON("Title", "Body"),
		`"isDraft":false,`,
		"",
		1,
	)
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubCommand([]string{
			"pr", "list", "--repo", "github.com/owner/repo",
			"--state", "all", "--head", "feature", "--limit", "100",
			"--json", githubPullRequestFields,
		}, ""),
		result: commandResult{Stdout: "[" + output + "]"},
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	if _, err := forge.PullRequests(
		context.Background(), repository, "feature",
	); err == nil {
		t.Fatal("PullRequests() accepted an omitted isDraft field")
	}
	runner.done()
}

func TestGitHubPreflightStrictlyRejectsUnknownOutput(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{{
		command: githubCommand([]string{
			"repo", "view", "github.com/owner/repo",
			"--json", "id,nameWithOwner,url",
		}, ""),
		result: commandResult{Stdout: `{
			"id":"R_node",
			"nameWithOwner":"owner/repo",
			"url":"https://github.com/owner/repo",
			"unexpected":true
		}`},
	}}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	if err := forge.Preflight(context.Background(), repository); err == nil {
		t.Fatal("Preflight() unexpectedly accepted unknown output")
	}
	runner.done()
}

func TestGitHubUpdateReportsUnknownOutcomeWhenRereadFails(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	expected := pullRequest{
		ID: "PR_node", Number: 7, URL: "https://github.com/owner/repo/pull/7",
		State: "OPEN", Title: "Old", Body: "Old body", AuthorLogin: "agent",
		BaseRefName: "master", BaseRefOID: testOID('a'),
		HeadRefName: "feature;not-shell", HeadRefOID: testOID('b'),
		HeadRepositoryOwner: "owner", HeadRepositoryName: "repo",
		UpdatedAt: "2026-08-30T00:00:00Z",
	}
	input := pullRequestInput{
		BaseRefName: "master", HeadRefName: "feature;not-shell",
		ExpectedHeadOID: testOID('b'), Title: "New", Body: "New body",
	}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubCommand([]string{
				"pr", "view", "7", "--repo", "github.com/owner/repo",
				"--json", githubPullRequestFields,
			}, ""),
			result: commandResult{Stdout: githubPRJSON("Old", "Old body")},
		},
		{
			command: githubCommand([]string{
				"pr", "edit", "7", "--repo", "github.com/owner/repo",
				"--base", "master", "--title", "New", "--body-file", "-",
			}, "New body"),
		},
		{
			command: githubCommand([]string{
				"pr", "view", "7", "--repo", "github.com/owner/repo",
				"--json", githubPullRequestFields,
			}, ""),
			err: errors.New("connection closed"),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.UpdatePullRequest(
		context.Background(), repository, expected, input,
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") {
		t.Fatalf("UpdatePullRequest() error = %v", err)
	}
	runner.done()
}

func TestGitHubCreatePostverifiesExpectedHead(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	input := pullRequestInput{
		BaseRefName: "master", HeadRefName: "feature;not-shell",
		ExpectedHeadOID: testOID('c'), Title: "Title", Body: "Body",
	}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubCommand([]string{
				"pr", "create",
				"--repo", "github.com/owner/repo",
				"--base", "master",
				"--head", "feature;not-shell",
				"--title", "Title",
				"--body-file", "-",
			}, "Body"),
		},
		{
			command: githubCommand([]string{
				"pr", "list",
				"--repo", "github.com/owner/repo",
				"--state", "all",
				"--head", "feature;not-shell",
				"--limit", "100",
				"--json", githubPullRequestFields,
			}, ""),
			result: commandResult{Stdout: "[" + githubPRJSON("Title", "Body") + "]"},
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.CreatePullRequest(context.Background(), repository, input)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "head") {
		t.Fatalf("CreatePullRequest() error = %v", err)
	}
	runner.done()
}

func TestGitHubCreateRejectsInvalidUTF8BeforeMutation(t *testing.T) {
	t.Parallel()
	runner := &transcriptRunner{t: t}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.CreatePullRequest(
		context.Background(),
		remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"},
		pullRequestInput{
			BaseRefName: "master", HeadRefName: "feature",
			ExpectedHeadOID: testOID('b'), Title: string([]byte{0xff}), Body: "Body",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "UTF-8") {
		t.Fatalf("CreatePullRequest() error = %v", err)
	}
	runner.done()
}

func TestGitHubCreateReportsUnknownOutcomeWhenRereadFails(t *testing.T) {
	t.Parallel()
	repository := remoteRepository{Host: "github.com", Owner: "owner", Name: "repo"}
	runner := &transcriptRunner{t: t, expected: []expectedCommand{
		{
			command: githubCommand([]string{
				"pr", "create",
				"--repo", "github.com/owner/repo",
				"--base", "master",
				"--head", "feature;not-shell",
				"--title", "Title",
				"--body-file", "-",
			}, "Body"),
		},
		{
			command: githubCommand([]string{
				"pr", "list",
				"--repo", "github.com/owner/repo",
				"--state", "all",
				"--head", "feature;not-shell",
				"--limit", "100",
				"--json", githubPullRequestFields,
			}, ""),
			err: errors.New("connection closed"),
		},
	}}
	forge := &githubForge{executable: "gh-test", runner: runner, directory: "/repo"}
	_, err := forge.CreatePullRequest(
		context.Background(),
		repository,
		pullRequestInput{
			BaseRefName: "master", HeadRefName: "feature;not-shell",
			ExpectedHeadOID: testOID('b'), Title: "Title", Body: "Body",
		},
	)
	if err == nil || !strings.Contains(err.Error(), "outcome unknown") ||
		!strings.Contains(err.Error(), "re-inspect before retrying") {
		t.Fatalf("CreatePullRequest() error = %v", err)
	}
	runner.done()
}
