package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const githubPullRequestFields = "id,number,url,state,title,body,author," +
	"baseRefName,baseRefOid,headRefName,headRefOid,headRepository," +
	"headRepositoryOwner,isCrossRepository,isDraft,updatedAt"

const githubOutputLimit = 32 * 1024 * 1024

const githubPaginationLimit = 1000

const (
	githubInventoryRequestLimit = 512
	githubInventoryNodeLimit    = 20_000
	githubInventoryByteLimit    = 8 * 1024 * 1024
)

const githubGraphQLPullRequestFields = `
    id
    number
    url
    state
    title
    body
    author { login }
    baseRefName
    baseRefOid
    headRefName
    headRefOid
    headRepository { name }
    headRepositoryOwner { login }
    isCrossRepository
    isDraft
    updatedAt
`

var errGitHubReviewEpochAdvanced = errors.New(
	"GitHub pull-request review epoch advanced while reading inventory",
)

type githubReviewEpochAdvancedError struct {
	PullRequest pullRequest
}

func (e *githubReviewEpochAdvancedError) Error() string {
	return errGitHubReviewEpochAdvanced.Error()
}

func (e *githubReviewEpochAdvancedError) Unwrap() error {
	return errGitHubReviewEpochAdvanced
}

const githubCommentsQuery = `
query($owner:String!,$name:String!,$number:Int!,$cursor:String) {
  repository(owner:$owner,name:$name) {
    pullRequest(number:$number) {
` + githubGraphQLPullRequestFields + `
      comments(first:100,after:$cursor) {
        nodes {
          id
          url
          body
          author { login }
          createdAt
          updatedAt
        }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
}`

const githubReviewsQuery = `
query($owner:String!,$name:String!,$number:Int!,$cursor:String) {
  repository(owner:$owner,name:$name) {
    pullRequest(number:$number) {
` + githubGraphQLPullRequestFields + `
      reviews(first:100,after:$cursor) {
        nodes {
          id
          url
          body
          author { login }
          state
          submittedAt
          commit { oid }
        }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
}`

const githubThreadsQuery = `
query($owner:String!,$name:String!,$number:Int!,$cursor:String) {
  repository(owner:$owner,name:$name) {
    pullRequest(number:$number) {
` + githubGraphQLPullRequestFields + `
      reviewThreads(first:100,after:$cursor) {
        nodes {
          id
          isResolved
          isOutdated
          path
          line
          originalLine
          startLine
          originalStartLine
        }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
}`

const githubThreadCommentsQuery = `
query($owner:String!,$name:String!,$number:Int!,$threadId:ID!,$cursor:String) {
  repository(owner:$owner,name:$name) {
    pullRequest(number:$number) {
` + githubGraphQLPullRequestFields + `
    }
  }
  node(id:$threadId) {
    ... on PullRequestReviewThread {
      id
      isResolved
      isOutdated
      path
      line
      originalLine
      startLine
      originalStartLine
      comments(first:100,after:$cursor) {
        nodes {
          id
          url
          body
          author { login }
          createdAt
          updatedAt
          path
          line
          originalLine
          commit { oid }
        }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
}`

const githubReviewRequestsQuery = `
query($owner:String!,$name:String!,$number:Int!,$cursor:String) {
  repository(owner:$owner,name:$name) {
    pullRequest(number:$number) {
` + githubGraphQLPullRequestFields + `
      reviewRequests(first:100,after:$cursor) {
        nodes {
          id
          requestedReviewer {
            __typename
            ... on Node { id }
            ... on Actor { login }
            ... on Team { slug organization { login } }
          }
        }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
}`

const githubReplyMutation = `
mutation($threadId:ID!,$body:String!) {
  addPullRequestReviewThreadReply(input:{
    pullRequestReviewThreadId:$threadId,
    body:$body
  }) {
    comment { id }
  }
}`

const githubResolveMutation = `
mutation($threadId:ID!) {
  resolveReviewThread(input:{threadId:$threadId}) {
    thread { id isResolved }
  }
}`

var githubEnvironment = []string{
	"GH_PROMPT_DISABLED=1",
	"GH_PAGER=cat",
	"PAGER=cat",
	"NO_COLOR=1",
	"CLICOLOR=0",
	"GH_NO_UPDATE_NOTIFIER=1",
	"GH_SPINNER_DISABLED=1",
}

var githubUnsetEnvironment = []string{
	"GH_HOST",
	"GH_REPO",
	"GH_DEBUG",
	"GH_FORCE_TTY",
	"DEBUG",
	"GODEBUG",
	"SSLKEYLOGFILE",
	"CLICOLOR_FORCE",
}

type githubForge struct {
	executable string
	runner     commandRunner
	directory  string
}

func githubMutationOutcomeUnknown(
	action string,
	mutationErr error,
	reconciliationErr error,
) error {
	if mutationErr != nil {
		return fmt.Errorf(
			"%s outcome unknown: mutation may have succeeded; re-inspect before retrying: mutation reported %v; post-mutation reconciliation failed: %w",
			action,
			mutationErr,
			reconciliationErr,
		)
	}
	return fmt.Errorf(
		"%s outcome unknown: mutation may have succeeded; re-inspect before retrying: post-mutation reconciliation failed: %w",
		action,
		reconciliationErr,
	)
}

var jsonRawMessageType = reflect.TypeOf(json.RawMessage(nil))

type githubOpenObject map[string]json.RawMessage

var githubOpenObjectType = reflect.TypeOf(githubOpenObject(nil))

type githubJSONField struct {
	Type     reflect.Type
	Required bool
}

// decodeGitHubJSON treats forge CLI output as an untrusted protocol boundary.
// encoding/json otherwise accepts invalid UTF-8, overwrites duplicate keys,
// and matches struct keys case-insensitively. None of those behaviors is safe
// for state used as a mutation guard.
func decodeGitHubJSON(label string, contents []byte, output any) error {
	if !utf8.Valid(contents) {
		return fmt.Errorf("decode %s: output is not valid UTF-8", label)
	}
	value := reflect.ValueOf(output)
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() {
		return fmt.Errorf("decode %s: destination must be a nonnil pointer", label)
	}

	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.UseNumber()
	shape, err := readUniqueGitHubJSONValue(decoder)
	if err != nil {
		return fmt.Errorf("decode %s shape: %w", label, err)
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode %s shape: multiple JSON values", label)
		}
		return fmt.Errorf("decode %s shape: trailing data: %w", label, err)
	}
	if err := validateGitHubJSONShape(shape, value.Type().Elem()); err != nil {
		return fmt.Errorf("decode %s shape: %w", label, err)
	}

	typed := json.NewDecoder(bytes.NewReader(contents))
	typed.DisallowUnknownFields()
	if err := typed.Decode(output); err != nil {
		return fmt.Errorf("decode %s: %w", label, err)
	}
	var extra any
	if err := typed.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode %s: multiple JSON values", label)
		}
		return fmt.Errorf("finish decoding %s: %w", label, err)
	}
	return nil
}

func readUniqueGitHubJSONValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return token, nil
	}
	switch delimiter {
	case '{':
		result := map[string]any{}
		seen := map[string]string{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return nil, err
			}
			key, ok := keyToken.(string)
			if !ok {
				return nil, fmt.Errorf("object contains a non-string key")
			}
			folded := strings.ToLower(key)
			if prior, exists := seen[folded]; exists {
				if prior == key {
					return nil, fmt.Errorf("object contains a duplicate key")
				}
				return nil, fmt.Errorf("object contains case-variant keys")
			}
			seen[folded] = key
			child, err := readUniqueGitHubJSONValue(decoder)
			if err != nil {
				return nil, err
			}
			result[key] = child
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim('}') {
			return nil, fmt.Errorf("object has an invalid terminator")
		}
		return result, nil
	case '[':
		result := []any{}
		for decoder.More() {
			child, err := readUniqueGitHubJSONValue(decoder)
			if err != nil {
				return nil, err
			}
			result = append(result, child)
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim(']') {
			return nil, fmt.Errorf("array has an invalid terminator")
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unexpected JSON delimiter")
	}
}

func validateGitHubJSONShape(value any, target reflect.Type) error {
	for target.Kind() == reflect.Pointer {
		if value == nil {
			return nil
		}
		target = target.Elem()
	}
	if target == jsonRawMessageType {
		return nil
	}
	if value == nil {
		return fmt.Errorf("nonnullable field contains null")
	}
	switch target.Kind() {
	case reflect.Struct:
		object, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("value does not have the required object shape")
		}
		fields := map[string]githubJSONField{}
		if err := collectGitHubJSONFields(
			target,
			fields,
			map[reflect.Type]bool{},
		); err != nil {
			return err
		}
		for key, child := range object {
			field, ok := fields[key]
			if !ok {
				return fmt.Errorf("object contains a non-canonical or unknown key")
			}
			if err := validateGitHubJSONShape(child, field.Type); err != nil {
				return err
			}
		}
		missing := []string{}
		for key, field := range fields {
			if _, exists := object[key]; field.Required && !exists {
				missing = append(missing, key)
			}
		}
		if len(missing) != 0 {
			sort.Strings(missing)
			return fmt.Errorf(
				"object omits required key %q",
				missing[0],
			)
		}
	case reflect.Slice, reflect.Array:
		values, ok := value.([]any)
		if !ok {
			return fmt.Errorf("value does not have the required array shape")
		}
		for _, child := range values {
			if err := validateGitHubJSONShape(child, target.Elem()); err != nil {
				return err
			}
		}
	case reflect.Map:
		// Deliberate forward-compatibility boundary for explicitly declared raw
		// REST response objects. Duplicate and case-variant keys were still
		// rejected by readUniqueGitHubJSONValue.
		if target == githubOpenObjectType {
			return nil
		}
		return fmt.Errorf("decoder schema contains an unbounded map")
	case reflect.Interface:
		return fmt.Errorf("decoder schema contains an unbounded interface")
	}
	return nil
}

func collectGitHubJSONFields(
	target reflect.Type,
	result map[string]githubJSONField,
	visiting map[reflect.Type]bool,
) error {
	if visiting[target] {
		return fmt.Errorf("decoder schema contains recursive anonymous embedding")
	}
	visiting[target] = true
	defer delete(visiting, target)
	for index := 0; index < target.NumField(); index++ {
		field := target.Field(index)
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}
		tag := strings.Split(field.Tag.Get("json"), ",")
		name := tag[0]
		if name == "-" {
			continue
		}
		if field.Anonymous && name == "" {
			embedded := field.Type
			for embedded.Kind() == reflect.Pointer {
				embedded = embedded.Elem()
			}
			if embedded.Kind() == reflect.Struct {
				if err := collectGitHubJSONFields(
					embedded,
					result,
					visiting,
				); err != nil {
					return err
				}
				continue
			}
		}
		if name == "" {
			name = field.Name
		}
		omitEmpty := false
		for _, option := range tag[1:] {
			if option == "omitempty" {
				omitEmpty = true
				break
			}
		}
		if _, exists := result[name]; exists {
			return fmt.Errorf("decoder schema contains an ambiguous field")
		}
		result[name] = githubJSONField{
			Type:     field.Type,
			Required: !omitEmpty,
		}
	}
	return nil
}

func (g *githubForge) Name() string {
	return "github"
}

func (g *githubForge) run(
	ctx context.Context,
	args []string,
	stdin string,
) (commandResult, error) {
	return g.runner.Run(ctx, command{
		Name:             g.executable,
		Args:             args,
		Dir:              g.directory,
		Env:              githubEnvironment,
		UnsetEnv:         append(githubUnsetEnvironment, gitUnsetEnvironment...),
		UnsetEnvPrefixes: gitUnsetEnvironmentPrefixes,
		Stdin:            stdin,
		OutputLimit:      githubOutputLimit,
	})
}

func (g *githubForge) Preflight(
	ctx context.Context,
	repository remoteRepository,
) error {
	result, err := g.run(ctx, []string{
		"repo", "view", repository.ghName(),
		"--json", "id,nameWithOwner,url",
	}, "")
	if err != nil {
		return fmt.Errorf("GitHub CLI repository preflight failed: %w", err)
	}
	var metadata struct {
		ID            string `json:"id"`
		NameWithOwner string `json:"nameWithOwner"`
		URL           string `json:"url"`
	}
	if err := decodeGitHubJSON(
		"GitHub repository preflight",
		[]byte(result.Stdout),
		&metadata,
	); err != nil {
		return err
	}
	wantName := repository.Owner + "/" + repository.Name
	if metadata.ID == "" ||
		!strings.EqualFold(metadata.NameWithOwner, wantName) ||
		metadata.URL == "" {
		return fmt.Errorf("GitHub CLI repository preflight returned incomplete or mismatched metadata")
	}
	return nil
}

type githubIdentity struct {
	Login string `json:"login"`
}

type githubRepositoryIdentity struct {
	Name string `json:"name"`
}

type githubPullRequest struct {
	ID                  string                   `json:"id"`
	Number              int                      `json:"number"`
	URL                 string                   `json:"url"`
	State               string                   `json:"state"`
	Title               string                   `json:"title"`
	Body                string                   `json:"body"`
	Author              githubIdentity           `json:"author"`
	BaseRefName         string                   `json:"baseRefName"`
	BaseRefOID          string                   `json:"baseRefOid"`
	HeadRefName         string                   `json:"headRefName"`
	HeadRefOID          string                   `json:"headRefOid"`
	HeadRepository      githubRepositoryIdentity `json:"headRepository"`
	HeadRepositoryOwner githubIdentity           `json:"headRepositoryOwner"`
	IsCrossRepository   bool                     `json:"isCrossRepository"`
	IsDraft             bool                     `json:"isDraft"`
	UpdatedAt           string                   `json:"updatedAt"`
}

// githubCLIPullRequest models the exact nested object projection emitted by
// `gh pr list/view --json`. The CLI projection is richer than the GraphQL
// selection above and must remain a separate strict protocol shape.
type githubCLIIdentity struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type githubCLIRepository struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
}

type githubCLIRepositoryOwner struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

type githubCLIPullRequest struct {
	ID                  string                   `json:"id"`
	Number              int                      `json:"number"`
	URL                 string                   `json:"url"`
	State               string                   `json:"state"`
	Title               string                   `json:"title"`
	Body                string                   `json:"body"`
	Author              githubCLIIdentity        `json:"author"`
	BaseRefName         string                   `json:"baseRefName"`
	BaseRefOID          string                   `json:"baseRefOid"`
	HeadRefName         string                   `json:"headRefName"`
	HeadRefOID          string                   `json:"headRefOid"`
	HeadRepository      githubCLIRepository      `json:"headRepository"`
	HeadRepositoryOwner githubCLIRepositoryOwner `json:"headRepositoryOwner"`
	IsCrossRepository   bool                     `json:"isCrossRepository"`
	IsDraft             bool                     `json:"isDraft"`
	UpdatedAt           string                   `json:"updatedAt"`
}

func (p githubCLIPullRequest) pullRequest() pullRequest {
	return pullRequest{
		ID:                  p.ID,
		Number:              p.Number,
		URL:                 p.URL,
		State:               p.State,
		Title:               p.Title,
		Body:                p.Body,
		AuthorLogin:         p.Author.Login,
		BaseRefName:         p.BaseRefName,
		BaseRefOID:          p.BaseRefOID,
		HeadRefName:         p.HeadRefName,
		HeadRefOID:          p.HeadRefOID,
		HeadRepositoryOwner: p.HeadRepositoryOwner.Login,
		HeadRepositoryName:  p.HeadRepository.Name,
		IsCrossRepository:   p.IsCrossRepository,
		IsDraft:             p.IsDraft,
		UpdatedAt:           p.UpdatedAt,
	}
}

func validateGitHubCLIPullRequest(value githubCLIPullRequest) error {
	pullRequest := value.pullRequest()
	if err := validateGitHubPullRequestRecord(pullRequest); err != nil {
		return err
	}
	for name, id := range map[string]string{
		"author":                value.Author.ID,
		"head repository":       value.HeadRepository.ID,
		"head repository owner": value.HeadRepositoryOwner.ID,
	} {
		if err := validateOpaqueID("GitHub CLI "+name+" ID", id); err != nil {
			return fmt.Errorf("record contains malformed GitHub CLI %s identity", name)
		}
	}
	wantNameWithOwner := value.HeadRepositoryOwner.Login + "/" +
		value.HeadRepository.Name
	if value.HeadRepository.NameWithOwner == "" ||
		!strings.EqualFold(value.HeadRepository.NameWithOwner, wantNameWithOwner) ||
		strings.ContainsAny(value.Author.Login, " \t\r\n") ||
		strings.ContainsAny(value.HeadRepositoryOwner.Login, " \t\r\n") {
		return fmt.Errorf("record contains inconsistent GitHub CLI nested metadata")
	}
	return nil
}

func (p githubPullRequest) pullRequest() pullRequest {
	return pullRequest{
		ID:                  p.ID,
		Number:              p.Number,
		URL:                 p.URL,
		State:               p.State,
		Title:               p.Title,
		Body:                p.Body,
		AuthorLogin:         p.Author.Login,
		BaseRefName:         p.BaseRefName,
		BaseRefOID:          p.BaseRefOID,
		HeadRefName:         p.HeadRefName,
		HeadRefOID:          p.HeadRefOID,
		HeadRepositoryOwner: p.HeadRepositoryOwner.Login,
		HeadRepositoryName:  p.HeadRepository.Name,
		IsCrossRepository:   p.IsCrossRepository,
		IsDraft:             p.IsDraft,
		UpdatedAt:           p.UpdatedAt,
	}
}

func validGitHubURL(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	return err == nil && parsed.IsAbs() && parsed.Host != "" &&
		strings.EqualFold(parsed.Scheme, "https") && parsed.User == nil
}

func validGitHubTimestamp(value string) bool {
	if value == "" {
		return false
	}
	_, err := time.Parse(time.RFC3339Nano, value)
	return err == nil
}

func validateGitHubPullRequestRecord(value pullRequest) error {
	if err := validatePullRequest(value); err != nil {
		return err
	}
	if !validGitHubURL(value.URL) || strings.TrimSpace(value.Title) == "" ||
		value.AuthorLogin == "" || !validGitHubTimestamp(value.UpdatedAt) {
		return fmt.Errorf("record contains malformed GitHub metadata")
	}
	return nil
}

func (g *githubForge) PullRequests(
	ctx context.Context,
	repository remoteRepository,
	head string,
) ([]pullRequest, error) {
	result, err := g.run(ctx, []string{
		"pr", "list",
		"--repo", repository.ghName(),
		"--state", "all",
		"--head", head,
		"--limit", "100",
		"--json", githubPullRequestFields,
	}, "")
	if err != nil {
		return nil, fmt.Errorf("list GitHub pull requests: %w", err)
	}
	var decoded []*githubCLIPullRequest
	if err := decodeGitHubJSON(
		"GitHub pull requests",
		[]byte(result.Stdout),
		&decoded,
	); err != nil {
		return nil, err
	}
	if decoded == nil {
		return nil, fmt.Errorf("decode GitHub pull requests: top-level result is null")
	}
	pullRequests := make([]pullRequest, 0, len(decoded))
	for index, candidate := range decoded {
		if candidate == nil {
			return nil, fmt.Errorf(
				"decode GitHub pull requests: record %d is null",
				index,
			)
		}
		pullRequest := candidate.pullRequest()
		if err := validateGitHubCLIPullRequest(*candidate); err != nil {
			return nil, fmt.Errorf(
				"decode GitHub pull requests: malformed record %d: %w",
				index,
				err,
			)
		}
		pullRequests = append(pullRequests, pullRequest)
	}
	return pullRequests, nil
}

func (g *githubForge) viewPullRequest(
	ctx context.Context,
	repository remoteRepository,
	number int,
) (*pullRequest, error) {
	return g.viewPullRequestWithBudget(
		ctx,
		repository,
		number,
		nil,
		"",
	)
}

func (g *githubForge) viewPullRequestWithBudget(
	ctx context.Context,
	repository remoteRepository,
	number int,
	budget *githubInventoryBudget,
	connection string,
) (*pullRequest, error) {
	if budget != nil {
		if err := budget.consumeRequest(connection); err != nil {
			return nil, err
		}
	}
	result, err := g.run(ctx, []string{
		"pr", "view", strconv.Itoa(number),
		"--repo", repository.ghName(),
		"--json", githubPullRequestFields,
	}, "")
	if err != nil {
		return nil, fmt.Errorf("view GitHub pull request %d: %w", number, err)
	}
	if budget != nil {
		if err := budget.consumeBytes(connection, len(result.Stdout)); err != nil {
			return nil, err
		}
	}
	var decoded githubCLIPullRequest
	if err := decodeGitHubJSON(
		fmt.Sprintf("GitHub pull request %d", number),
		[]byte(result.Stdout),
		&decoded,
	); err != nil {
		return nil, err
	}
	pullRequest := decoded.pullRequest()
	if err := validateGitHubCLIPullRequest(decoded); err != nil {
		return nil, fmt.Errorf("decode GitHub pull request %d: %w", number, err)
	}
	return &pullRequest, nil
}

func (g *githubForge) CreatePullRequest(
	ctx context.Context,
	repository remoteRepository,
	input pullRequestInput,
) (*pullRequest, error) {
	if err := validatePullRequestInput(input); err != nil {
		return nil, err
	}
	_, createErr := g.run(ctx, []string{
		"pr", "create",
		"--repo", repository.ghName(),
		"--base", input.BaseRefName,
		"--head", input.HeadRefName,
		"--title", input.Title,
		"--body-file", "-",
	}, input.Body)

	pullRequests, lookupErr := g.PullRequests(ctx, repository, input.HeadRefName)
	if lookupErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"create GitHub pull request",
			createErr,
			lookupErr,
		)
	}
	exact, err := exactPullRequests(pullRequests, repository, input.HeadRefName)
	if err != nil {
		return nil, githubMutationOutcomeUnknown(
			"create GitHub pull request",
			createErr,
			err,
		)
	}
	pullRequest, refusals := selectOpenPullRequest(exact)
	if pullRequest == nil || len(refusals) != 0 {
		return nil, githubMutationOutcomeUnknown(
			"create GitHub pull request",
			createErr,
			fmt.Errorf(
				"no exact open pull request appeared: %s",
				strings.Join(refusals, "; "),
			),
		)
	}
	if !pullRequestMatchesInput(*pullRequest, input) {
		return nil, githubMutationOutcomeUnknown(
			"create GitHub pull request",
			createErr,
			fmt.Errorf(
				"created pull request %d does not match the requested projection or head",
				pullRequest.Number,
			),
		)
	}
	return pullRequest, nil
}

func (g *githubForge) UpdatePullRequest(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	input pullRequestInput,
) (*pullRequest, error) {
	if err := validatePullRequestInput(input); err != nil {
		return nil, err
	}
	current, err := g.viewPullRequest(ctx, repository, pullRequest.Number)
	if err != nil {
		return nil, err
	}
	if !samePullRequestMetadata(*current, pullRequest) {
		return nil, fmt.Errorf(
			"pull request %d changed before update; preserve possible human edits",
			pullRequest.Number,
		)
	}
	_, updateErr := g.run(ctx, []string{
		"pr", "edit", strconv.Itoa(pullRequest.Number),
		"--repo", repository.ghName(),
		"--base", input.BaseRefName,
		"--title", input.Title,
		"--body-file", "-",
	}, input.Body)
	refreshed, viewErr := g.viewPullRequest(ctx, repository, pullRequest.Number)
	if viewErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"update GitHub pull request",
			updateErr,
			viewErr,
		)
	}
	if !pullRequestMatchesInput(*refreshed, input) {
		return nil, githubMutationOutcomeUnknown(
			"update GitHub pull request",
			updateErr,
			fmt.Errorf(
				"updated pull request %d does not match the requested projection or head",
				pullRequest.Number,
			),
		)
	}
	expectedIdentity := pullRequest
	expectedIdentity.BaseRefName = input.BaseRefName
	expectedIdentity.HeadRefOID = input.ExpectedHeadOID
	if !samePullRequestIdentity(*refreshed, expectedIdentity) {
		return nil, githubMutationOutcomeUnknown(
			"update GitHub pull request",
			updateErr,
			fmt.Errorf(
				"pull request %d identity changed during update",
				pullRequest.Number,
			),
		)
	}
	return refreshed, nil
}

type githubPageInfo struct {
	HasNextPage *bool   `json:"hasNextPage"`
	EndCursor   *string `json:"endCursor"`
}

type githubConnection[T any] struct {
	Nodes    *[]T            `json:"nodes"`
	PageInfo *githubPageInfo `json:"pageInfo"`
}

type githubInventoryBudget struct {
	Requests int
	Nodes    int
	Bytes    int
}

func (b *githubInventoryBudget) consumeRequest(connection string) error {
	if b.Requests >= githubInventoryRequestLimit {
		return fmt.Errorf(
			"GitHub review inventory exceeded %d requests while reading %s; refusing an incomplete result",
			githubInventoryRequestLimit,
			connection,
		)
	}
	b.Requests++
	return nil
}

func (b *githubInventoryBudget) consumeNode(connection string) error {
	if b.Nodes >= githubInventoryNodeLimit {
		return fmt.Errorf(
			"GitHub review inventory exceeded %d nodes while reading %s; refusing an incomplete result",
			githubInventoryNodeLimit,
			connection,
		)
	}
	b.Nodes++
	return nil
}

func (b *githubInventoryBudget) consumeBytes(
	connection string,
	bytes int,
) error {
	if bytes < 0 || b.Bytes > githubInventoryByteLimit-bytes {
		return fmt.Errorf(
			"GitHub review inventory exceeded %d response bytes while reading %s; refusing an incomplete result",
			githubInventoryByteLimit,
			connection,
		)
	}
	b.Bytes += bytes
	return nil
}

func githubConnectionPage[T any](
	name string,
	connection *githubConnection[T],
) ([]T, githubPageInfo, error) {
	if connection == nil {
		return nil, githubPageInfo{}, fmt.Errorf(
			"GitHub %s returned a null or missing connection",
			name,
		)
	}
	if connection.Nodes == nil {
		return nil, githubPageInfo{}, fmt.Errorf(
			"GitHub %s returned null or missing nodes",
			name,
		)
	}
	if connection.PageInfo == nil ||
		connection.PageInfo.HasNextPage == nil {
		return nil, githubPageInfo{}, fmt.Errorf(
			"GitHub %s returned null or missing pageInfo",
			name,
		)
	}
	return *connection.Nodes, *connection.PageInfo, nil
}

type githubCommit struct {
	OID string `json:"oid"`
}

type githubReviewComment struct {
	ID           string         `json:"id"`
	URL          string         `json:"url"`
	Body         string         `json:"body"`
	Author       githubIdentity `json:"author"`
	CreatedAt    string         `json:"createdAt"`
	UpdatedAt    string         `json:"updatedAt"`
	Path         string         `json:"path"`
	Line         *int           `json:"line"`
	OriginalLine *int           `json:"originalLine"`
	Commit       githubCommit   `json:"commit"`
}

func (c githubReviewComment) reviewComment() reviewComment {
	return reviewComment{
		ID:           c.ID,
		URL:          c.URL,
		Body:         c.Body,
		AuthorLogin:  c.Author.Login,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		Path:         c.Path,
		Line:         c.Line,
		OriginalLine: c.OriginalLine,
		CommitOID:    c.Commit.OID,
	}
}

type githubTopLevelComment struct {
	ID        string         `json:"id"`
	URL       string         `json:"url"`
	Body      string         `json:"body"`
	Author    githubIdentity `json:"author"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
}

func validateGitHubCommentIdentity(
	id string,
	url string,
	author githubIdentity,
	createdAt string,
	updatedAt string,
) error {
	if err := validateOpaqueID("GitHub comment ID", id); err != nil {
		return fmt.Errorf("GitHub returned a malformed comment ID")
	}
	if !validGitHubURL(url) || author.Login == "" ||
		strings.ContainsAny(author.Login, " \t\r\n") {
		return fmt.Errorf("GitHub returned malformed comment identity metadata")
	}
	created, createdErr := time.Parse(time.RFC3339Nano, createdAt)
	updated, updatedErr := time.Parse(time.RFC3339Nano, updatedAt)
	if createdErr != nil || updatedErr != nil || updated.Before(created) {
		return fmt.Errorf("GitHub returned malformed comment timestamps")
	}
	return nil
}

func validateGitHubTopLevelComment(value githubTopLevelComment) error {
	return validateGitHubCommentIdentity(
		value.ID,
		value.URL,
		value.Author,
		value.CreatedAt,
		value.UpdatedAt,
	)
}

func validateGitHubReviewComment(value githubReviewComment) error {
	if err := validateGitHubCommentIdentity(
		value.ID,
		value.URL,
		value.Author,
		value.CreatedAt,
		value.UpdatedAt,
	); err != nil {
		return err
	}
	if value.Path == "" || strings.ContainsAny(value.Path, "\x00\r\n") ||
		!isObjectID(value.Commit.OID) {
		return fmt.Errorf("GitHub returned malformed review-comment metadata")
	}
	for _, line := range []*int{value.Line, value.OriginalLine} {
		if line != nil && *line <= 0 {
			return fmt.Errorf("GitHub returned malformed review-comment line metadata")
		}
	}
	return nil
}

func (c githubTopLevelComment) reviewComment() reviewComment {
	return reviewComment{
		ID:          c.ID,
		URL:         c.URL,
		Body:        c.Body,
		AuthorLogin: c.Author.Login,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type githubReview struct {
	ID          string         `json:"id"`
	URL         string         `json:"url"`
	Body        string         `json:"body"`
	Author      githubIdentity `json:"author"`
	State       string         `json:"state"`
	SubmittedAt *string        `json:"submittedAt"`
	Commit      githubCommit   `json:"commit"`
}

func validateGitHubReview(value githubReview) error {
	if err := validateOpaqueID("GitHub review ID", value.ID); err != nil {
		return fmt.Errorf("GitHub returned a malformed review ID")
	}
	if !validGitHubURL(value.URL) || value.Author.Login == "" ||
		strings.ContainsAny(value.Author.Login, " \t\r\n") ||
		!isObjectID(value.Commit.OID) {
		return fmt.Errorf("GitHub returned malformed review identity metadata")
	}
	switch strings.ToUpper(value.State) {
	case "PENDING":
		if value.SubmittedAt != nil {
			return fmt.Errorf("GitHub returned malformed pending-review metadata")
		}
	case "APPROVED", "CHANGES_REQUESTED", "COMMENTED", "DISMISSED":
		if value.SubmittedAt == nil || !validGitHubTimestamp(*value.SubmittedAt) {
			return fmt.Errorf("GitHub returned malformed review timestamp metadata")
		}
	default:
		return fmt.Errorf("GitHub returned unknown review state")
	}
	return nil
}

func (r githubReview) pullRequestReview() pullRequestReview {
	submittedAt := ""
	if r.SubmittedAt != nil {
		submittedAt = *r.SubmittedAt
	}
	return pullRequestReview{
		ID:          r.ID,
		URL:         r.URL,
		Body:        r.Body,
		AuthorLogin: r.Author.Login,
		State:       r.State,
		SubmittedAt: submittedAt,
		CommitOID:   r.Commit.OID,
	}
}

type githubReviewThread struct {
	ID                string `json:"id"`
	IsResolved        bool   `json:"isResolved"`
	IsOutdated        bool   `json:"isOutdated"`
	Path              string `json:"path"`
	Line              *int   `json:"line"`
	OriginalLine      *int   `json:"originalLine"`
	StartLine         *int   `json:"startLine"`
	OriginalStartLine *int   `json:"originalStartLine"`
}

func validateGitHubReviewThread(value githubReviewThread) error {
	if err := validateOpaqueID("GitHub review thread ID", value.ID); err != nil {
		return fmt.Errorf("GitHub returned a malformed review-thread ID")
	}
	if value.Path == "" || strings.ContainsAny(value.Path, "\x00\r\n") {
		return fmt.Errorf("GitHub returned malformed review-thread path metadata")
	}
	for _, line := range []*int{
		value.Line,
		value.OriginalLine,
		value.StartLine,
		value.OriginalStartLine,
	} {
		if line != nil && *line <= 0 {
			return fmt.Errorf("GitHub returned malformed review-thread line metadata")
		}
	}
	if value.Line != nil && value.StartLine != nil &&
		*value.StartLine > *value.Line {
		return fmt.Errorf("GitHub returned inverted review-thread line metadata")
	}
	if value.OriginalLine != nil && value.OriginalStartLine != nil &&
		*value.OriginalStartLine > *value.OriginalLine {
		return fmt.Errorf("GitHub returned inverted review-thread line metadata")
	}
	return nil
}

func (t githubReviewThread) reviewThread() reviewThread {
	return reviewThread{
		ID:                t.ID,
		IsResolved:        t.IsResolved,
		IsOutdated:        t.IsOutdated,
		Path:              t.Path,
		Line:              t.Line,
		OriginalLine:      t.OriginalLine,
		StartLine:         t.StartLine,
		OriginalStartLine: t.OriginalStartLine,
		Comments:          []reviewComment{},
	}
}

type githubRequestedReviewer struct {
	Typename     string         `json:"__typename"`
	ID           string         `json:"id"`
	Login        string         `json:"login,omitempty"`
	Slug         string         `json:"slug,omitempty"`
	Organization githubIdentity `json:"organization,omitempty"`
}

type githubReviewRequest struct {
	ID       string                  `json:"id"`
	Reviewer githubRequestedReviewer `json:"requestedReviewer"`
}

func (r githubRequestedReviewer) requestedReviewer() (requestedReviewer, error) {
	if err := validateOpaqueID("GitHub requested reviewer ID", r.ID); err != nil {
		return requestedReviewer{}, fmt.Errorf(
			"GitHub returned a requested reviewer with a malformed node ID",
		)
	}
	switch r.Typename {
	case "User", "Bot", "Mannequin":
		if r.Login == "" || strings.ContainsAny(r.Login, " \t\r\n") {
			return requestedReviewer{}, fmt.Errorf(
				"GitHub returned a requested %s without a login",
				strings.ToLower(r.Typename),
			)
		}
		return requestedReviewer{
			Type:  strings.ToLower(r.Typename),
			ID:    r.ID,
			Login: r.Login,
		}, nil
	case "Team":
		if r.Organization.Login == "" || r.Slug == "" ||
			strings.ContainsAny(r.Organization.Login, " \t\r\n") ||
			strings.ContainsAny(r.Slug, " \t\r\n") {
			return requestedReviewer{}, fmt.Errorf(
				"GitHub returned a requested team without an identity",
			)
		}
		return requestedReviewer{
			Type:  "team",
			ID:    r.ID,
			Login: r.Organization.Login + "/" + r.Slug,
		}, nil
	case "EnterpriseTeam":
		return requestedReviewer{
			Type: "enterprise_team",
			ID:   r.ID,
		}, nil
	default:
		if r.Typename == "" {
			return requestedReviewer{}, fmt.Errorf(
				"GitHub returned a requested reviewer without a type",
			)
		}
		return requestedReviewer{
			Type: strings.ToLower(r.Typename),
			ID:   r.ID,
		}, nil
	}
}

type githubGraphQLError struct {
	Message    string          `json:"message"`
	Type       string          `json:"type,omitempty"`
	Path       json.RawMessage `json:"path,omitempty"`
	Locations  json.RawMessage `json:"locations,omitempty"`
	Extensions json.RawMessage `json:"extensions,omitempty"`
}

func (g *githubForge) graphQL(
	ctx context.Context,
	repository remoteRepository,
	query string,
	variables map[string]any,
	output any,
) error {
	return g.graphQLWithBudget(
		ctx,
		repository,
		query,
		variables,
		output,
		nil,
		"",
	)
}

func (g *githubForge) inventoryGraphQL(
	ctx context.Context,
	repository remoteRepository,
	query string,
	variables map[string]any,
	output any,
	budget *githubInventoryBudget,
	connection string,
) error {
	if budget == nil {
		return fmt.Errorf("GitHub review inventory budget is absent")
	}
	return g.graphQLWithBudget(
		ctx,
		repository,
		query,
		variables,
		output,
		budget,
		connection,
	)
}

func (g *githubForge) graphQLWithBudget(
	ctx context.Context,
	repository remoteRepository,
	query string,
	variables map[string]any,
	output any,
	budget *githubInventoryBudget,
	connection string,
) error {
	request, err := json.Marshal(struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("encode GitHub GraphQL request: %w", err)
	}
	result, err := g.run(ctx, []string{
		"api", "--hostname", repository.Host,
		"graphql", "--input", "-",
	}, string(request))
	if err != nil {
		return fmt.Errorf("run GitHub GraphQL request: %w", err)
	}
	if budget != nil {
		if err := budget.consumeBytes(connection, len(result.Stdout)); err != nil {
			return err
		}
	}
	var response struct {
		Data       json.RawMessage      `json:"data"`
		Errors     []githubGraphQLError `json:"errors,omitempty"`
		Extensions json.RawMessage      `json:"extensions,omitempty"`
	}
	if err := decodeGitHubJSON(
		"GitHub GraphQL envelope",
		[]byte(result.Stdout),
		&response,
	); err != nil {
		return err
	}
	if len(response.Errors) != 0 {
		messages := make([]string, 0, len(response.Errors))
		for _, graphError := range response.Errors {
			message := strings.TrimSpace(graphError.Message)
			if message == "" {
				message = "unspecified GraphQL error"
			}
			messages = append(messages, redactCredentials(message))
		}
		return fmt.Errorf(
			"GitHub GraphQL returned errors: %s",
			strings.Join(messages, "; "),
		)
	}
	if len(response.Data) == 0 || string(response.Data) == "null" {
		return fmt.Errorf("GitHub GraphQL returned no data")
	}
	if err := decodeGitHubJSON(
		"GitHub GraphQL data",
		response.Data,
		output,
	); err != nil {
		return err
	}
	return nil
}

func githubVariables(
	repository remoteRepository,
	pullRequest pullRequest,
	cursor string,
) map[string]any {
	variables := map[string]any{
		"owner":  repository.Owner,
		"name":   repository.Name,
		"number": pullRequest.Number,
		"cursor": nil,
	}
	if cursor != "" {
		variables["cursor"] = cursor
	}
	return variables
}

func validateGithubReviewPullRequest(
	value *githubPullRequest,
	expected pullRequest,
) error {
	if value == nil {
		return fmt.Errorf("GitHub did not return the exact pull request")
	}
	actual := value.pullRequest()
	if err := validateGitHubPullRequestRecord(actual); err != nil {
		return fmt.Errorf("GitHub returned malformed pull-request identity: %w", err)
	}
	if !sameGithubReviewPullRequestContext(actual, expected) {
		return fmt.Errorf(
			"pull request identity, head, or review state changed while reading review state",
		)
	}
	if githubUpdatedAtAdvanced(expected.UpdatedAt, actual.UpdatedAt) {
		return &githubReviewEpochAdvancedError{PullRequest: actual}
	}
	if actual.UpdatedAt == "" || expected.UpdatedAt == "" ||
		actual.UpdatedAt != expected.UpdatedAt {
		return fmt.Errorf(
			"pull request identity, head, or review state changed while reading review state",
		)
	}
	return nil
}

func githubUpdatedAtAdvanced(before, after string) bool {
	beforeTime, beforeErr := time.Parse(time.RFC3339Nano, before)
	afterTime, afterErr := time.Parse(time.RFC3339Nano, after)
	return beforeErr == nil && afterErr == nil && afterTime.After(beforeTime)
}

func sameGithubReviewPullRequestContext(left, right pullRequest) bool {
	return samePullRequestMetadata(left, right) &&
		left.BaseRefOID == right.BaseRefOID
}

func nextGithubCursor(
	connection string,
	page int,
	info githubPageInfo,
	seen map[string]bool,
) (string, bool, error) {
	if info.HasNextPage == nil {
		return "", false, fmt.Errorf(
			"GitHub %s returned null or missing pageInfo",
			connection,
		)
	}
	if !*info.HasNextPage {
		return "", true, nil
	}
	if page >= githubPaginationLimit {
		return "", false, fmt.Errorf(
			"GitHub %s exceeded %d pages; refusing an incomplete result",
			connection,
			githubPaginationLimit,
		)
	}
	if info.EndCursor == nil || *info.EndCursor == "" ||
		seen[*info.EndCursor] {
		return "", false, fmt.Errorf(
			"GitHub %s returned an invalid pagination cursor; refusing an incomplete result",
			connection,
		)
	}
	seen[*info.EndCursor] = true
	return *info.EndCursor, false, nil
}

func rememberGithubNode(
	connection string,
	id string,
	seen map[string]bool,
) error {
	if id == "" {
		return fmt.Errorf("GitHub %s returned a node without an ID", connection)
	}
	if seen[id] {
		return fmt.Errorf(
			"GitHub %s returned duplicate node %q; refusing an inconsistent result",
			connection,
			id,
		)
	}
	seen[id] = true
	return nil
}

func (g *githubForge) topLevelComments(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
) ([]reviewComment, error) {
	return g.topLevelCommentsWithBudget(
		ctx,
		repository,
		pullRequest,
		&githubInventoryBudget{},
	)
}

func (g *githubForge) topLevelCommentsWithBudget(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	budget *githubInventoryBudget,
) ([]reviewComment, error) {
	comments := []reviewComment{}
	seenNodes := map[string]bool{}
	seenCursors := map[string]bool{}
	cursor := ""
	for page := 1; ; page++ {
		if err := budget.consumeRequest("comments"); err != nil {
			return nil, err
		}
		var data struct {
			Repository struct {
				PullRequest *struct {
					githubPullRequest
					Comments *githubConnection[githubTopLevelComment] `json:"comments"`
				} `json:"pullRequest"`
			} `json:"repository"`
		}
		err := g.inventoryGraphQL(
			ctx,
			repository,
			githubCommentsQuery,
			githubVariables(repository, pullRequest, cursor),
			&data,
			budget,
			"comments",
		)
		if err != nil {
			return nil, fmt.Errorf("read GitHub top-level comments: %w", err)
		}
		current := data.Repository.PullRequest
		if current == nil {
			return nil, fmt.Errorf("GitHub did not return the exact pull request")
		}
		if err := validateGithubReviewPullRequest(
			&current.githubPullRequest,
			pullRequest,
		); err != nil {
			return nil, err
		}
		nodes, pageInfo, err := githubConnectionPage(
			"comments",
			current.Comments,
		)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			if err := budget.consumeNode("comments"); err != nil {
				return nil, err
			}
			if err := validateGitHubTopLevelComment(node); err != nil {
				return nil, err
			}
			if err := rememberGithubNode("comments", node.ID, seenNodes); err != nil {
				return nil, err
			}
			comments = append(comments, node.reviewComment())
		}
		next, done, err := nextGithubCursor(
			"comments",
			page,
			pageInfo,
			seenCursors,
		)
		if err != nil {
			return nil, err
		}
		if done {
			return comments, nil
		}
		cursor = next
	}
}

func (g *githubForge) pullRequestReviews(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
) ([]pullRequestReview, error) {
	return g.pullRequestReviewsWithBudget(
		ctx,
		repository,
		pullRequest,
		&githubInventoryBudget{},
	)
}

func (g *githubForge) pullRequestReviewsWithBudget(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	budget *githubInventoryBudget,
) ([]pullRequestReview, error) {
	reviews := []pullRequestReview{}
	seenNodes := map[string]bool{}
	seenCursors := map[string]bool{}
	cursor := ""
	for page := 1; ; page++ {
		if err := budget.consumeRequest("reviews"); err != nil {
			return nil, err
		}
		var data struct {
			Repository struct {
				PullRequest *struct {
					githubPullRequest
					Reviews *githubConnection[githubReview] `json:"reviews"`
				} `json:"pullRequest"`
			} `json:"repository"`
		}
		err := g.inventoryGraphQL(
			ctx,
			repository,
			githubReviewsQuery,
			githubVariables(repository, pullRequest, cursor),
			&data,
			budget,
			"reviews",
		)
		if err != nil {
			return nil, fmt.Errorf("read GitHub reviews: %w", err)
		}
		current := data.Repository.PullRequest
		if current == nil {
			return nil, fmt.Errorf("GitHub did not return the exact pull request")
		}
		if err := validateGithubReviewPullRequest(
			&current.githubPullRequest,
			pullRequest,
		); err != nil {
			return nil, err
		}
		nodes, pageInfo, err := githubConnectionPage(
			"reviews",
			current.Reviews,
		)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			if err := budget.consumeNode("reviews"); err != nil {
				return nil, err
			}
			if err := validateGitHubReview(node); err != nil {
				return nil, err
			}
			if err := rememberGithubNode("reviews", node.ID, seenNodes); err != nil {
				return nil, err
			}
			reviews = append(reviews, node.pullRequestReview())
		}
		next, done, err := nextGithubCursor(
			"reviews",
			page,
			pageInfo,
			seenCursors,
		)
		if err != nil {
			return nil, err
		}
		if done {
			return reviews, nil
		}
		cursor = next
	}
}

func (g *githubForge) reviewThreads(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
) ([]reviewThread, error) {
	return g.reviewThreadsWithBudget(
		ctx,
		repository,
		pullRequest,
		&githubInventoryBudget{},
	)
}

func (g *githubForge) reviewThreadsWithBudget(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	budget *githubInventoryBudget,
) ([]reviewThread, error) {
	threads := []reviewThread{}
	seenNodes := map[string]bool{}
	seenCursors := map[string]bool{}
	cursor := ""
	for page := 1; ; page++ {
		if err := budget.consumeRequest("review threads"); err != nil {
			return nil, err
		}
		var data struct {
			Repository struct {
				PullRequest *struct {
					githubPullRequest
					Threads *githubConnection[githubReviewThread] `json:"reviewThreads"`
				} `json:"pullRequest"`
			} `json:"repository"`
		}
		err := g.inventoryGraphQL(
			ctx,
			repository,
			githubThreadsQuery,
			githubVariables(repository, pullRequest, cursor),
			&data,
			budget,
			"review threads",
		)
		if err != nil {
			return nil, fmt.Errorf("read GitHub review threads: %w", err)
		}
		current := data.Repository.PullRequest
		if current == nil {
			return nil, fmt.Errorf("GitHub did not return the exact pull request")
		}
		if err := validateGithubReviewPullRequest(
			&current.githubPullRequest,
			pullRequest,
		); err != nil {
			return nil, err
		}
		nodes, pageInfo, err := githubConnectionPage(
			"review threads",
			current.Threads,
		)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			if err := budget.consumeNode("review threads"); err != nil {
				return nil, err
			}
			if err := validateGitHubReviewThread(node); err != nil {
				return nil, err
			}
			if err := rememberGithubNode("review threads", node.ID, seenNodes); err != nil {
				return nil, err
			}
			thread, err := g.reviewThreadDetailsWithBudget(
				ctx,
				repository,
				pullRequest,
				node.ID,
				budget,
			)
			if err != nil {
				return nil, err
			}
			if len(thread.Comments) == 0 {
				return nil, fmt.Errorf(
					"GitHub review thread %q has no comments; refusing incomplete state",
					thread.ID,
				)
			}
			threads = append(threads, thread)
		}
		next, done, err := nextGithubCursor(
			"review threads",
			page,
			pageInfo,
			seenCursors,
		)
		if err != nil {
			return nil, err
		}
		if done {
			return threads, nil
		}
		cursor = next
	}
}

func (g *githubForge) reviewThreadDetails(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	threadID string,
) (reviewThread, error) {
	return g.reviewThreadDetailsWithBudget(
		ctx,
		repository,
		pullRequest,
		threadID,
		&githubInventoryBudget{},
	)
}

func (g *githubForge) reviewThreadDetailsWithBudget(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	threadID string,
	budget *githubInventoryBudget,
) (reviewThread, error) {
	comments := []reviewComment{}
	var snapshot *reviewThread
	seenNodes := map[string]bool{}
	seenCursors := map[string]bool{}
	cursor := ""
	for page := 1; ; page++ {
		if err := budget.consumeRequest("review thread comments"); err != nil {
			return reviewThread{}, err
		}
		variables := githubVariables(repository, pullRequest, cursor)
		variables["threadId"] = threadID
		var data struct {
			Repository struct {
				PullRequest *githubPullRequest `json:"pullRequest"`
			} `json:"repository"`
			Node *struct {
				githubReviewThread
				Comments *githubConnection[githubReviewComment] `json:"comments"`
			} `json:"node"`
		}
		err := g.inventoryGraphQL(
			ctx,
			repository,
			githubThreadCommentsQuery,
			variables,
			&data,
			budget,
			"review thread comments",
		)
		if err != nil {
			return reviewThread{}, fmt.Errorf(
				"read comments for GitHub review thread %q: %w",
				threadID,
				err,
			)
		}
		if err := validateGithubReviewPullRequest(
			data.Repository.PullRequest,
			pullRequest,
		); err != nil {
			return reviewThread{}, err
		}
		if data.Node == nil || data.Node.ID != threadID {
			return reviewThread{}, fmt.Errorf(
				"GitHub did not return exact review thread %q",
				threadID,
			)
		}
		if err := validateGitHubReviewThread(
			data.Node.githubReviewThread,
		); err != nil {
			return reviewThread{}, err
		}
		current := data.Node.githubReviewThread.reviewThread()
		if snapshot == nil {
			snapshot = &current
		} else if !sameGithubThreadMetadata(*snapshot, current) {
			return reviewThread{}, fmt.Errorf(
				"GitHub review thread %q changed during pagination",
				threadID,
			)
		}
		nodes, pageInfo, err := githubConnectionPage(
			"review thread comments",
			data.Node.Comments,
		)
		if err != nil {
			return reviewThread{}, err
		}
		for _, node := range nodes {
			if err := budget.consumeNode("review thread comments"); err != nil {
				return reviewThread{}, err
			}
			if err := validateGitHubReviewComment(node); err != nil {
				return reviewThread{}, err
			}
			if err := rememberGithubNode(
				"review thread comments",
				node.ID,
				seenNodes,
			); err != nil {
				return reviewThread{}, err
			}
			comments = append(comments, node.reviewComment())
		}
		next, done, err := nextGithubCursor(
			"review thread comments",
			page,
			pageInfo,
			seenCursors,
		)
		if err != nil {
			return reviewThread{}, err
		}
		if done {
			if snapshot == nil {
				return reviewThread{}, fmt.Errorf(
					"GitHub did not return exact review thread %q",
					threadID,
				)
			}
			snapshot.Comments = comments
			snapshot.ExpectationDigest = reviewThreadDigest(*snapshot)
			return *snapshot, nil
		}
		cursor = next
	}
}

func sameGithubThreadMetadata(left, right reviewThread) bool {
	return left.ID == right.ID && left.IsResolved == right.IsResolved &&
		left.IsOutdated == right.IsOutdated && left.Path == right.Path &&
		equalOptionalInt(left.Line, right.Line) &&
		equalOptionalInt(left.OriginalLine, right.OriginalLine) &&
		equalOptionalInt(left.StartLine, right.StartLine) &&
		equalOptionalInt(left.OriginalStartLine, right.OriginalStartLine)
}

func equalOptionalInt(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func (g *githubForge) requestedReviewers(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
) ([]requestedReviewer, error) {
	return g.requestedReviewersWithBudget(
		ctx,
		repository,
		pullRequest,
		&githubInventoryBudget{},
	)
}

func (g *githubForge) requestedReviewersWithBudget(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	budget *githubInventoryBudget,
) ([]requestedReviewer, error) {
	reviewers := []requestedReviewer{}
	seenNodes := map[string]bool{}
	seenCursors := map[string]bool{}
	cursor := ""
	for page := 1; ; page++ {
		if err := budget.consumeRequest("review requests"); err != nil {
			return nil, err
		}
		var data struct {
			Repository struct {
				PullRequest *struct {
					githubPullRequest
					ReviewRequests *githubConnection[githubReviewRequest] `json:"reviewRequests"`
				} `json:"pullRequest"`
			} `json:"repository"`
		}
		err := g.inventoryGraphQL(
			ctx,
			repository,
			githubReviewRequestsQuery,
			githubVariables(repository, pullRequest, cursor),
			&data,
			budget,
			"review requests",
		)
		if err != nil {
			return nil, fmt.Errorf("read GitHub review requests: %w", err)
		}
		current := data.Repository.PullRequest
		if current == nil {
			return nil, fmt.Errorf("GitHub did not return the exact pull request")
		}
		if err := validateGithubReviewPullRequest(
			&current.githubPullRequest,
			pullRequest,
		); err != nil {
			return nil, err
		}
		nodes, pageInfo, err := githubConnectionPage(
			"review requests",
			current.ReviewRequests,
		)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			if err := budget.consumeNode("review requests"); err != nil {
				return nil, err
			}
			if err := rememberGithubNode(
				"review requests",
				node.ID,
				seenNodes,
			); err != nil {
				return nil, err
			}
			reviewer, err := node.Reviewer.requestedReviewer()
			if err != nil {
				return nil, err
			}
			reviewer.RequestID = node.ID
			reviewers = append(reviewers, reviewer)
		}
		next, done, err := nextGithubCursor(
			"review requests",
			page,
			pageInfo,
			seenCursors,
		)
		if err != nil {
			return nil, err
		}
		if done {
			return reviewers, nil
		}
		cursor = next
	}
}

func (g *githubForge) InspectReviews(
	ctx context.Context,
	repository remoteRepository,
	expected pullRequest,
) (*reviewInspection, error) {
	return g.inspectReviews(
		ctx,
		repository,
		expected,
		&githubInventoryBudget{},
	)
}

func (g *githubForge) reconcileReplyInspection(
	ctx context.Context,
	repository remoteRepository,
	expected pullRequest,
	beforeThread *reviewThread,
	body string,
	mutationCommentID string,
) (*reviewInspection, error) {
	budget := &githubInventoryBudget{}
	floor := expected
	for attempt := 0; attempt < 3; attempt++ {
		inspection, err := g.inspectReviews(
			ctx,
			repository,
			floor,
			budget,
		)
		var advanced *githubReviewEpochAdvancedError
		if errors.As(err, &advanced) {
			floor = advanced.PullRequest
			continue
		}
		if err != nil {
			return nil, err
		}
		floor = inspection.PullRequest
		afterThread, verifyErr := findGithubThread(
			inspection,
			beforeThread.ID,
		)
		if verifyErr == nil {
			verifyErr = verifyGithubReply(
				beforeThread,
				afterThread,
				body,
				mutationCommentID,
			)
		}
		if verifyErr == nil {
			return inspection, nil
		}
		if mutationCommentID != "" && afterThread != nil &&
			sameGithubReviewThreadSnapshot(*beforeThread, *afterThread) {
			continue
		}
		return nil, verifyErr
	}
	return nil, fmt.Errorf(
		"GitHub review state did not stabilize on the exact appended reply",
	)
}

func sameGithubReviewThreadSnapshot(left, right reviewThread) bool {
	if !sameGithubThreadMetadata(left, right) ||
		len(left.Comments) != len(right.Comments) {
		return false
	}
	for index := range left.Comments {
		if !sameReviewComment(left.Comments[index], right.Comments[index]) {
			return false
		}
	}
	return true
}

func (g *githubForge) inspectReviews(
	ctx context.Context,
	repository remoteRepository,
	expected pullRequest,
	budget *githubInventoryBudget,
) (*reviewInspection, error) {
	current, err := g.viewPullRequestWithBudget(
		ctx,
		repository,
		expected.Number,
		budget,
		"pull-request coherence reads",
	)
	if err != nil {
		return nil, err
	}
	if !sameGithubReviewPullRequestContext(*current, expected) {
		return nil, fmt.Errorf(
			"pull request identity or head changed before reading review state",
		)
	}
	if current.UpdatedAt == "" {
		return nil, fmt.Errorf(
			"GitHub pull request omitted updatedAt; cannot read coherent review state",
		)
	}
	if current.UpdatedAt != expected.UpdatedAt &&
		!githubUpdatedAtAdvanced(expected.UpdatedAt, current.UpdatedAt) {
		return nil, fmt.Errorf(
			"pull request review epoch moved backward or is not comparable",
		)
	}
	comments, err := g.topLevelCommentsWithBudget(
		ctx,
		repository,
		*current,
		budget,
	)
	if err != nil {
		return nil, err
	}
	reviews, err := g.pullRequestReviewsWithBudget(
		ctx,
		repository,
		*current,
		budget,
	)
	if err != nil {
		return nil, err
	}
	threads, err := g.reviewThreadsWithBudget(
		ctx,
		repository,
		*current,
		budget,
	)
	if err != nil {
		return nil, err
	}
	requested, err := g.requestedReviewersWithBudget(
		ctx,
		repository,
		*current,
		budget,
	)
	if err != nil {
		return nil, err
	}
	refreshed, err := g.viewPullRequestWithBudget(
		ctx,
		repository,
		expected.Number,
		budget,
		"pull-request coherence reads",
	)
	if err != nil {
		return nil, err
	}
	if !sameGithubReviewPullRequestContext(*refreshed, *current) {
		return nil, fmt.Errorf(
			"pull request identity, head, or review state changed while reading review state",
		)
	}
	if githubUpdatedAtAdvanced(current.UpdatedAt, refreshed.UpdatedAt) {
		return nil, &githubReviewEpochAdvancedError{PullRequest: *refreshed}
	}
	if refreshed.UpdatedAt == "" || refreshed.UpdatedAt != current.UpdatedAt {
		return nil, fmt.Errorf(
			"pull request identity, head, or review state changed while reading review state",
		)
	}
	return &reviewInspection{
		PullRequest:                    *refreshed,
		PullRequestExpectationDigest:   reviewPullRequestDigest(repository, *refreshed),
		Comments:                       comments,
		ExpectedLastTopLevelComment:    expectedLastTopLevelComment(comments),
		ExpectedTopLevelCommentsDigest: reviewCommentsDigest(comments),
		Reviews:                        reviews,
		Threads:                        threads,
		RequestedReviewers:             requested,
	}, nil
}

func expectedGithubTopLevelComments(
	inspection *reviewInspection,
	expectation topLevelCommentExpectation,
) error {
	actual := expectedLastTopLevelComment(inspection.Comments)
	if actual != expectation.ExpectedLastCommentID ||
		reviewCommentsDigest(inspection.Comments) != expectation.ExpectedDigest {
		return fmt.Errorf("top-level comments changed after inspection")
	}
	return nil
}

func sameReviewComment(left, right reviewComment) bool {
	return left.ID == right.ID && left.URL == right.URL &&
		normalizeText(left.Body) == normalizeText(right.Body) &&
		left.AuthorLogin == right.AuthorLogin &&
		left.CreatedAt == right.CreatedAt && left.UpdatedAt == right.UpdatedAt &&
		left.Path == right.Path && equalOptionalInt(left.Line, right.Line) &&
		equalOptionalInt(left.OriginalLine, right.OriginalLine) &&
		left.CommitOID == right.CommitOID
}

func verifyGithubTopLevelComment(
	before *reviewInspection,
	after *reviewInspection,
	body string,
	mutationCommentID string,
) error {
	if len(after.Comments) != len(before.Comments)+1 {
		return fmt.Errorf("comment verification found concurrent top-level changes")
	}
	for index := range before.Comments {
		if !sameReviewComment(before.Comments[index], after.Comments[index]) {
			return fmt.Errorf("comment verification found concurrent top-level changes")
		}
	}
	created := after.Comments[len(after.Comments)-1]
	if mutationCommentID != "" && created.ID != mutationCommentID {
		return fmt.Errorf("comment verification found a different created comment")
	}
	if normalizeText(created.Body) != normalizeText(body) {
		return fmt.Errorf("comment verification found different comment text")
	}
	return nil
}

func (g *githubForge) CommentOnPullRequest(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	expectation topLevelCommentExpectation,
	body string,
) (*reviewInspection, error) {
	if err := validateTopLevelCommentExpectation(expectation); err != nil {
		return nil, err
	}
	body, err := withCommentDisclaimer(body)
	if err != nil {
		return nil, err
	}
	before, err := g.InspectReviews(ctx, repository, pullRequest)
	if err != nil {
		return nil, err
	}
	if err := expectedGithubTopLevelComments(before, expectation); err != nil {
		return nil, err
	}
	requestBody, err := json.Marshal(struct {
		Body string `json:"body"`
	}{Body: body})
	if err != nil {
		return nil, fmt.Errorf("encode GitHub top-level comment: %w", err)
	}
	result, mutationErr := g.run(ctx, []string{
		"api", "--hostname", repository.Host,
		"--method", "POST",
		fmt.Sprintf(
			"repos/%s/%s/issues/%d/comments",
			repository.Owner,
			repository.Name,
			pullRequest.Number,
		),
		"--input", "-",
	}, string(requestBody))
	mutationCommentID := ""
	responseMalformed := false
	if strings.TrimSpace(result.Stdout) != "" {
		// The REST endpoint returns a large evolving comment object. Preserve
		// forward compatibility for fields we do not use, while the strict map
		// decode still rejects invalid UTF-8 and duplicate/case-variant keys.
		var response githubOpenObject
		if decodeErr := decodeGitHubJSON(
			"GitHub top-level comment response",
			[]byte(result.Stdout),
			&response,
		); decodeErr != nil {
			responseMalformed = true
			if mutationErr == nil {
				mutationErr = fmt.Errorf(
					"decode GitHub top-level comment response: %w",
					decodeErr,
				)
			}
		} else {
			rawID, ok := response["node_id"]
			if !ok {
				responseMalformed = true
				if mutationErr == nil {
					mutationErr = fmt.Errorf(
						"GitHub top-level comment response omitted its canonical node ID",
					)
				}
			} else if decodeErr := decodeGitHubJSON(
				"GitHub top-level comment node ID",
				rawID,
				&mutationCommentID,
			); decodeErr != nil {
				responseMalformed = true
				if mutationErr == nil {
					mutationErr = decodeErr
				}
			}
		}
	}
	if mutationErr == nil && mutationCommentID == "" {
		responseMalformed = true
		mutationErr = fmt.Errorf(
			"GitHub top-level comment response omitted its node ID",
		)
	}
	after, inspectErr := g.InspectReviews(ctx, repository, before.PullRequest)
	if inspectErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"add GitHub top-level comment",
			mutationErr,
			inspectErr,
		)
	}
	verifyErr := verifyGithubTopLevelComment(
		before,
		after,
		body,
		mutationCommentID,
	)
	if verifyErr == nil && responseMalformed {
		verifyErr = fmt.Errorf(
			"top-level comment mutation response was malformed",
		)
	}
	if verifyErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"add GitHub top-level comment",
			mutationErr,
			verifyErr,
		)
	}
	return after, nil
}

func expectedGithubThread(
	inspection *reviewInspection,
	expectation reviewThreadExpectation,
) (*reviewThread, error) {
	var selected *reviewThread
	for index := range inspection.Threads {
		candidate := &inspection.Threads[index]
		if candidate.ID != expectation.ThreadID {
			continue
		}
		if selected != nil {
			return nil, fmt.Errorf(
				"GitHub returned duplicate review thread %q",
				expectation.ThreadID,
			)
		}
		selected = candidate
	}
	if selected == nil {
		return nil, fmt.Errorf(
			"review thread %q no longer exists",
			expectation.ThreadID,
		)
	}
	if selected.IsResolved {
		return nil, fmt.Errorf(
			"review thread %q is already resolved",
			expectation.ThreadID,
		)
	}
	if len(selected.Comments) == 0 ||
		selected.Comments[len(selected.Comments)-1].ID !=
			expectation.ExpectedLastCommentID ||
		reviewThreadDigest(*selected) != expectation.ExpectedDigest {
		return nil, fmt.Errorf(
			"review thread %q changed after inspection",
			expectation.ThreadID,
		)
	}
	return selected, nil
}

func findGithubThread(
	inspection *reviewInspection,
	threadID string,
) (*reviewThread, error) {
	var selected *reviewThread
	for index := range inspection.Threads {
		candidate := &inspection.Threads[index]
		if candidate.ID != threadID {
			continue
		}
		if selected != nil {
			return nil, fmt.Errorf(
				"GitHub returned duplicate review thread %q",
				threadID,
			)
		}
		selected = candidate
	}
	if selected == nil {
		return nil, fmt.Errorf("review thread %q no longer exists", threadID)
	}
	return selected, nil
}

func verifyGithubReply(
	before *reviewThread,
	after *reviewThread,
	body string,
	mutationCommentID string,
) error {
	if mutationCommentID == "" {
		return fmt.Errorf("reply mutation omitted the created comment ID")
	}
	if before.IsResolved || after.IsResolved ||
		!sameGithubThreadMetadata(*before, *after) {
		return fmt.Errorf("reply verification found concurrent thread metadata changes")
	}
	if len(after.Comments) != len(before.Comments)+1 {
		return fmt.Errorf("reply verification found concurrent thread changes")
	}
	for index := range before.Comments {
		if !sameReviewComment(before.Comments[index], after.Comments[index]) {
			return fmt.Errorf("reply verification found concurrent thread changes")
		}
	}
	created := after.Comments[len(after.Comments)-1]
	if created.ID != mutationCommentID {
		return fmt.Errorf("reply verification found a different created comment")
	}
	if normalizeText(created.Body) != normalizeText(body) {
		return fmt.Errorf("reply verification found different comment text")
	}
	return nil
}

func (g *githubForge) ReplyToReviewThread(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	expectation reviewThreadExpectation,
	body string,
) (*reviewInspection, error) {
	if err := validateThreadExpectation(expectation); err != nil {
		return nil, err
	}
	body, err := withCommentDisclaimer(body)
	if err != nil {
		return nil, err
	}
	before, err := g.InspectReviews(ctx, repository, pullRequest)
	if err != nil {
		return nil, err
	}
	beforeThread, err := expectedGithubThread(before, expectation)
	if err != nil {
		return nil, err
	}
	var mutationData struct {
		Reply *struct {
			Comment *struct {
				ID string `json:"id"`
			} `json:"comment"`
		} `json:"addPullRequestReviewThreadReply"`
	}
	mutationErr := g.graphQL(
		ctx,
		repository,
		githubReplyMutation,
		map[string]any{
			"threadId": expectation.ThreadID,
			"body":     body,
		},
		&mutationData,
	)
	mutationCommentID := ""
	if mutationData.Reply != nil && mutationData.Reply.Comment != nil {
		mutationCommentID = mutationData.Reply.Comment.ID
	}
	after, inspectErr := g.reconcileReplyInspection(
		ctx,
		repository,
		before.PullRequest,
		beforeThread,
		body,
		mutationCommentID,
	)
	if inspectErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"reply to GitHub review thread",
			mutationErr,
			inspectErr,
		)
	}
	if mutationErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"reply to GitHub review thread",
			mutationErr,
			fmt.Errorf(
				"reply mutation response did not establish unambiguous success",
			),
		)
	}
	return after, nil
}

func verifyUnchangedGithubComments(
	before *reviewThread,
	after *reviewThread,
) error {
	if len(before.Comments) != len(after.Comments) {
		return fmt.Errorf("resolution verification found concurrent thread changes")
	}
	for index := range before.Comments {
		if !sameReviewComment(before.Comments[index], after.Comments[index]) {
			return fmt.Errorf("resolution verification found concurrent thread changes")
		}
	}
	return nil
}

func verifyGithubResolution(
	before *reviewThread,
	after *reviewThread,
) error {
	if before.IsResolved || !after.IsResolved ||
		before.ID != after.ID ||
		before.IsOutdated != after.IsOutdated ||
		before.Path != after.Path ||
		!equalOptionalInt(before.Line, after.Line) ||
		!equalOptionalInt(before.OriginalLine, after.OriginalLine) ||
		!equalOptionalInt(before.StartLine, after.StartLine) ||
		!equalOptionalInt(before.OriginalStartLine, after.OriginalStartLine) {
		return fmt.Errorf("resolution verification found concurrent thread metadata changes")
	}
	return verifyUnchangedGithubComments(before, after)
}

func (g *githubForge) ResolveReviewThread(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	expectation reviewThreadExpectation,
) (*reviewInspection, error) {
	if err := validateResolutionThreadExpectation(expectation); err != nil {
		return nil, err
	}
	budget := &githubInventoryBudget{}
	before, err := g.inspectReviews(ctx, repository, pullRequest, budget)
	if err != nil {
		return nil, err
	}
	if err := validateReviewInventoryUnique(before); err != nil {
		return nil, err
	}
	if reviewInventoryDigest(before) != expectation.RequiredInventoryDigest {
		return nil, fmt.Errorf(
			"complete bounded review inventory projection changed after the tool-issued reply receipt",
		)
	}
	expiresAt, err := time.Parse(time.RFC3339Nano, expectation.AuthorityExpiresAt)
	if err != nil || !reviewReplyReceiptNow(ctx).Before(expiresAt) {
		return nil, fmt.Errorf("review reply resolution authority expired")
	}
	beforeThread, err := expectedGithubThread(before, expectation)
	if err != nil {
		return nil, err
	}
	if err := requireReasonedToolComment(
		beforeThread.Comments[len(beforeThread.Comments)-1].Body,
	); err != nil {
		return nil, err
	}
	if reviewReplyBodyDigest(
		beforeThread.Comments[len(beforeThread.Comments)-1].Body,
	) != expectation.RequiredReplyBodyDigest {
		return nil, fmt.Errorf(
			"last review-thread comment does not match the tool-issued reply receipt",
		)
	}
	latestThread, err := g.reviewThreadDetailsWithBudget(
		ctx,
		repository,
		before.PullRequest,
		expectation.ThreadID,
		budget,
	)
	if err != nil {
		return nil, fmt.Errorf("reread review thread before resolution: %w", err)
	}
	if !sameGithubReviewThreadSnapshot(*beforeThread, latestThread) {
		return nil, fmt.Errorf("review thread changed immediately before resolution")
	}
	if !reviewReplyReceiptNow(ctx).Before(expiresAt) {
		return nil, fmt.Errorf("review reply resolution authority expired")
	}
	var mutationData struct {
		Resolve *struct {
			Thread *struct {
				ID         string `json:"id"`
				IsResolved bool   `json:"isResolved"`
			} `json:"thread"`
		} `json:"resolveReviewThread"`
	}
	if !reviewReplyReceiptNow(ctx).Before(expiresAt) {
		return nil, fmt.Errorf("review reply resolution authority expired")
	}
	mutationErr := g.graphQL(
		ctx,
		repository,
		githubResolveMutation,
		map[string]any{"threadId": expectation.ThreadID},
		&mutationData,
	)
	after, inspectErr := g.InspectReviews(ctx, repository, before.PullRequest)
	if inspectErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"resolve GitHub review thread",
			mutationErr,
			inspectErr,
		)
	}
	afterThread, verifyErr := findGithubThread(after, expectation.ThreadID)
	if verifyErr == nil {
		verifyErr = verifyGithubResolution(beforeThread, afterThread)
	}
	if verifyErr == nil &&
		(mutationData.Resolve == nil || mutationData.Resolve.Thread == nil ||
			mutationData.Resolve.Thread.ID != expectation.ThreadID ||
			!mutationData.Resolve.Thread.IsResolved) {
		verifyErr = fmt.Errorf("resolution mutation returned missing or mismatched state")
	}
	if verifyErr == nil && mutationErr != nil {
		verifyErr = fmt.Errorf(
			"resolution mutation response did not establish unambiguous success",
		)
	}
	if verifyErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"resolve GitHub review thread",
			mutationErr,
			verifyErr,
		)
	}
	return after, nil
}

func validGithubLogin(login string) bool {
	if len(login) == 0 || len(login) > 39 || login[0] == '-' ||
		login[len(login)-1] == '-' {
		return false
	}
	for _, character := range login {
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') || character == '-' {
			continue
		}
		return false
	}
	return true
}

func githubReviewersPresent(
	inspection *reviewInspection,
	reviewers []string,
) bool {
	present := make(map[string]bool, len(inspection.RequestedReviewers))
	for _, reviewer := range inspection.RequestedReviewers {
		if reviewer.Type == "user" || reviewer.Type == "bot" {
			present[strings.ToLower(reviewer.Login)] = true
		}
	}
	for _, reviewer := range reviewers {
		if !present[strings.ToLower(reviewer)] {
			return false
		}
	}
	return true
}

func missingGithubReviewers(
	inspection *reviewInspection,
	reviewers []string,
) []string {
	missing := make([]string, 0, len(reviewers))
	for _, reviewer := range reviewers {
		if !githubReviewersPresent(inspection, []string{reviewer}) {
			missing = append(missing, reviewer)
		}
	}
	return missing
}

func samePullRequestReview(left, right pullRequestReview) bool {
	return left.ID == right.ID && left.URL == right.URL &&
		normalizeText(left.Body) == normalizeText(right.Body) &&
		left.AuthorLogin == right.AuthorLogin && left.State == right.State &&
		left.SubmittedAt == right.SubmittedAt && left.CommitOID == right.CommitOID
}

func githubReviewIsCompleted(state string) bool {
	switch strings.ToUpper(state) {
	case "APPROVED", "CHANGES_REQUESTED", "COMMENTED":
		return true
	default:
		return false
	}
}

func completedGithubReviewSince(
	before *reviewInspection,
	after *reviewInspection,
	reviewer string,
	headOID string,
) bool {
	for _, candidate := range after.Reviews {
		if !strings.EqualFold(candidate.AuthorLogin, reviewer) ||
			!githubReviewIsCompleted(candidate.State) ||
			candidate.SubmittedAt == "" || candidate.CommitOID != headOID {
			continue
		}
		priorFound := false
		for _, prior := range before.Reviews {
			if prior.ID != candidate.ID {
				continue
			}
			priorFound = true
			// An already-completed review on this exact head cannot become
			// evidence of a newly satisfied request merely because it was
			// edited or resubmitted in place.
			if strings.EqualFold(prior.AuthorLogin, reviewer) &&
				githubReviewIsCompleted(prior.State) &&
				prior.SubmittedAt != "" && prior.CommitOID == headOID {
				break
			}
			if strings.EqualFold(prior.AuthorLogin, reviewer) {
				return true
			}
			break
		}
		if !priorFound {
			return true
		}
	}
	return false
}

func githubReviewersSatisfied(
	before *reviewInspection,
	after *reviewInspection,
	reviewers []string,
	headOID string,
) bool {
	for _, reviewer := range reviewers {
		if githubReviewersPresent(after, []string{reviewer}) ||
			completedGithubReviewSince(before, after, reviewer, headOID) {
			continue
		}
		return false
	}
	return true
}

func (g *githubForge) RequestReview(
	ctx context.Context,
	repository remoteRepository,
	pullRequest pullRequest,
	reviewers []string,
) (*reviewInspection, error) {
	reviewers, err := normalizeReviewers(reviewers)
	if err != nil {
		return nil, err
	}
	for _, reviewer := range reviewers {
		if !validGithubLogin(reviewer) {
			return nil, fmt.Errorf(
				"GitHub reviewer %q is not a valid user login",
				reviewer,
			)
		}
	}
	before, err := g.InspectReviews(ctx, repository, pullRequest)
	if err != nil {
		return nil, err
	}
	missing := missingGithubReviewers(before, reviewers)
	if len(missing) == 0 {
		return before, nil
	}
	body, err := json.Marshal(struct {
		Reviewers []string `json:"reviewers"`
	}{Reviewers: missing})
	if err != nil {
		return nil, fmt.Errorf("encode GitHub review request: %w", err)
	}
	_, mutationErr := g.run(ctx, []string{
		"api", "--hostname", repository.Host,
		"--method", "POST",
		"--silent",
		fmt.Sprintf(
			"repos/%s/%s/pulls/%d/requested_reviewers",
			repository.Owner,
			repository.Name,
			pullRequest.Number,
		),
		"--input", "-",
	}, string(body))
	after, inspectErr := g.InspectReviews(ctx, repository, before.PullRequest)
	if inspectErr != nil {
		return nil, githubMutationOutcomeUnknown(
			"request GitHub review",
			mutationErr,
			inspectErr,
		)
	}
	if !githubReviewersSatisfied(
		before,
		after,
		reviewers,
		before.PullRequest.HeadRefOID,
	) {
		verifyErr := fmt.Errorf("requested reviewers are absent after mutation")
		return nil, githubMutationOutcomeUnknown(
			"request GitHub review",
			mutationErr,
			verifyErr,
		)
	}
	return after, nil
}

func validatePullRequest(pullRequest pullRequest) error {
	state := strings.ToUpper(pullRequest.State)
	if pullRequest.ID == "" || pullRequest.Number <= 0 ||
		pullRequest.URL == "" || pullRequest.BaseRefName == "" ||
		pullRequest.HeadRefName == "" ||
		pullRequest.HeadRepositoryOwner == "" ||
		pullRequest.HeadRepositoryName == "" {
		return fmt.Errorf("record contains an empty required field")
	}
	if state != "OPEN" && state != "CLOSED" && state != "MERGED" {
		return fmt.Errorf("record has unknown state %q", pullRequest.State)
	}
	if !isObjectID(pullRequest.BaseRefOID) {
		return fmt.Errorf("record has invalid base OID")
	}
	if !isObjectID(pullRequest.HeadRefOID) {
		return fmt.Errorf("record has invalid head OID")
	}
	return nil
}

func validatePullRequestInput(input pullRequestInput) error {
	for name, value := range map[string]string{
		"base ref": input.BaseRefName,
		"head ref": input.HeadRefName,
		"title":    input.Title,
		"body":     input.Body,
	} {
		if !utf8.ValidString(value) {
			return fmt.Errorf("pull-request %s is not valid UTF-8", name)
		}
	}
	if input.BaseRefName == "" || input.HeadRefName == "" ||
		strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("pull-request input contains an empty required field")
	}
	if !isObjectID(input.ExpectedHeadOID) {
		return fmt.Errorf("pull-request input has an invalid expected head OID")
	}
	return nil
}

func pullRequestMatchesInput(
	pullRequest pullRequest,
	input pullRequestInput,
) bool {
	return pullRequest.BaseRefName == input.BaseRefName &&
		pullRequest.HeadRefName == input.HeadRefName &&
		pullRequest.HeadRefOID == input.ExpectedHeadOID &&
		normalizeText(pullRequest.Title) == normalizeText(input.Title) &&
		normalizeText(pullRequest.Body) == normalizeText(input.Body)
}
