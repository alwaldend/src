package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type remoteRepository struct {
	Host  string `json:"host"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type providerReport struct {
	RemoteRepository           remoteRepository `json:"remote_repository"`
	ProviderHint               string           `json:"provider_hint"`
	AdapterAvailable           bool             `json:"adapter_available"`
	GitTransport               string           `json:"git_transport"`
	DeliveryTransportAvailable bool             `json:"delivery_transport_available"`
}

func providerForRepository(repository remoteRepository) (string, bool) {
	switch strings.ToLower(repository.Host) {
	case "github.com":
		return "github", true
	case "codeberg.org", "forgejo.alwaldend.com":
		return "forgejo", false
	default:
		return "unknown", false
	}
}

func (r remoteRepository) ghName() string {
	return r.Host + "/" + r.Owner + "/" + r.Name
}

func parseRemoteRepository(raw string) (remoteRepository, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return remoteRepository{}, fmt.Errorf("remote URL is empty")
	}
	if !utf8.ValidString(value) {
		return remoteRepository{}, fmt.Errorf("remote URL is not valid UTF-8")
	}

	var host string
	var repositoryPath string
	if strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err != nil {
			// url.Parse errors can contain the complete input, including URL
			// userinfo. Never carry that diagnostic across this boundary.
			return remoteRepository{}, fmt.Errorf("remote URL could not be parsed")
		}
		switch strings.ToLower(parsed.Scheme) {
		case "https", "ssh":
		case "git", "http":
			return remoteRepository{}, fmt.Errorf(
				"remote URL uses insecure transport %q; use HTTPS or SSH",
				parsed.Scheme,
			)
		default:
			return remoteRepository{}, fmt.Errorf(
				"remote URL uses an unsupported transport",
			)
		}
		if parsed.RawQuery != "" || parsed.Fragment != "" {
			return remoteRepository{}, fmt.Errorf(
				"remote URLs with queries or fragments are not supported",
			)
		}
		host = parsed.Hostname()
		if parsed.Port() != "" {
			return remoteRepository{}, fmt.Errorf(
				"remote URLs with explicit ports are not supported",
			)
		}
		if parsed.RawPath != "" || strings.Contains(parsed.EscapedPath(), "%") {
			return remoteRepository{}, fmt.Errorf(
				"remote URL path must use an unescaped canonical form",
			)
		}
		repositoryPath = parsed.Path
		if !strings.HasPrefix(repositoryPath, "/") {
			return remoteRepository{}, fmt.Errorf(
				"remote URL does not identify HOST/OWNER/REPOSITORY",
			)
		}
		repositoryPath = strings.TrimPrefix(repositoryPath, "/")
	} else {
		colon := strings.IndexByte(value, ':')
		at := -1
		if colon > 0 {
			at = strings.LastIndexByte(value[:colon], '@')
		}
		if colon <= 0 || at < 0 {
			return remoteRepository{}, fmt.Errorf("unsupported remote URL form")
		}
		if at == 0 || strings.Contains(value[colon+1:], ":") {
			return remoteRepository{}, fmt.Errorf("unsupported remote URL form")
		}
		host = value[at+1 : colon]
		repositoryPath = value[colon+1:]
		if strings.HasPrefix(repositoryPath, "/") {
			return remoteRepository{}, fmt.Errorf(
				"SCP-style remote paths must be repository-relative",
			)
		}
	}

	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" || strings.ContainsAny(host, "\\/:@ \t\r\n") {
		return remoteRepository{}, fmt.Errorf(
			"remote URL does not identify HOST/OWNER/REPOSITORY",
		)
	}
	if repositoryPath == "" || strings.Contains(repositoryPath, "\\") ||
		strings.HasPrefix(repositoryPath, "/") ||
		strings.HasSuffix(repositoryPath, "/") ||
		strings.Contains(repositoryPath, "//") {
		return remoteRepository{}, fmt.Errorf(
			"remote URL path is not canonical",
		)
	}
	repositoryPath = strings.TrimSuffix(repositoryPath, ".git")
	parts := strings.Split(repositoryPath, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" ||
		parts[0] == "." || parts[0] == ".." ||
		parts[1] == "." || parts[1] == ".." {
		return remoteRepository{}, fmt.Errorf(
			"remote URL does not identify HOST/OWNER/REPOSITORY",
		)
	}
	return remoteRepository{
		Host:  host,
		Owner: parts[0],
		Name:  parts[1],
	}, nil
}

func sameRemoteRepository(left, right remoteRepository) bool {
	return strings.EqualFold(left.Host, right.Host) &&
		strings.EqualFold(left.Owner, right.Owner) &&
		strings.EqualFold(left.Name, right.Name)
}

type pullRequest struct {
	ID                  string `json:"id"`
	Number              int    `json:"number"`
	URL                 string `json:"url"`
	State               string `json:"state"`
	Title               string `json:"title"`
	Body                string `json:"body"`
	AuthorLogin         string `json:"author_login,omitempty"`
	BaseRefName         string `json:"base_ref_name"`
	BaseRefOID          string `json:"base_ref_oid"`
	HeadRefName         string `json:"head_ref_name"`
	HeadRefOID          string `json:"head_ref_oid"`
	HeadRepositoryOwner string `json:"head_repository_owner"`
	HeadRepositoryName  string `json:"head_repository_name"`
	IsCrossRepository   bool   `json:"is_cross_repository"`
	IsDraft             bool   `json:"is_draft"`
	UpdatedAt           string `json:"updated_at,omitempty"`
}

type pullRequestInput struct {
	BaseRefName     string
	HeadRefName     string
	ExpectedHeadOID string
	Title           string
	Body            string
}

type reviewComment struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	Body         string `json:"body"`
	AuthorLogin  string `json:"author_login,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Path         string `json:"path,omitempty"`
	Line         *int   `json:"line,omitempty"`
	OriginalLine *int   `json:"original_line,omitempty"`
	CommitOID    string `json:"commit_oid,omitempty"`
}

type pullRequestReview struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Body        string `json:"body"`
	AuthorLogin string `json:"author_login,omitempty"`
	State       string `json:"state"`
	SubmittedAt string `json:"submitted_at,omitempty"`
	CommitOID   string `json:"commit_oid,omitempty"`
}

type reviewThread struct {
	ID                string          `json:"id"`
	IsResolved        bool            `json:"is_resolved"`
	IsOutdated        bool            `json:"is_outdated"`
	Path              string          `json:"path,omitempty"`
	Line              *int            `json:"line,omitempty"`
	OriginalLine      *int            `json:"original_line,omitempty"`
	StartLine         *int            `json:"start_line,omitempty"`
	OriginalStartLine *int            `json:"original_start_line,omitempty"`
	Comments          []reviewComment `json:"comments"`
	ExpectationDigest string          `json:"expectation_digest"`
}

type requestedReviewer struct {
	RequestID string `json:"request_id"`
	Type      string `json:"type"`
	ID        string `json:"id"`
	Login     string `json:"login"`
}

type reviewInspection struct {
	PullRequest                    pullRequest         `json:"pull_request"`
	PullRequestExpectationDigest   string              `json:"pull_request_expectation_digest"`
	Comments                       []reviewComment     `json:"top_level_comments"`
	ExpectedLastTopLevelComment    string              `json:"expected_last_top_level_comment"`
	ExpectedTopLevelCommentsDigest string              `json:"expected_top_level_comments_digest"`
	Reviews                        []pullRequestReview `json:"reviews"`
	Threads                        []reviewThread      `json:"review_threads"`
	RequestedReviewers             []requestedReviewer `json:"requested_reviewers"`
}

type reviewThreadExpectation struct {
	ThreadID                string
	ExpectedLastCommentID   string
	ExpectedDigest          string
	RequiredReplyBodyDigest string
	RequiredInventoryDigest string
	AuthorityExpiresAt      string
}

type topLevelCommentExpectation struct {
	ExpectedLastCommentID string
	ExpectedDigest        string
}

type forge interface {
	Name() string
	Preflight(context.Context, remoteRepository) error
	PullRequests(
		context.Context,
		remoteRepository,
		string,
	) ([]pullRequest, error)
	CreatePullRequest(
		context.Context,
		remoteRepository,
		pullRequestInput,
	) (*pullRequest, error)
	UpdatePullRequest(
		context.Context,
		remoteRepository,
		pullRequest,
		pullRequestInput,
	) (*pullRequest, error)
}

// reviewForge is an optional capability. A forge can support core delivery
// without pretending that its CLI exposes the review inventory and mutation
// operations needed by the review workflow.
type reviewForge interface {
	InspectReviews(
		context.Context,
		remoteRepository,
		pullRequest,
	) (*reviewInspection, error)
	ReplyToReviewThread(
		context.Context,
		remoteRepository,
		pullRequest,
		reviewThreadExpectation,
		string,
	) (*reviewInspection, error)
	CommentOnPullRequest(
		context.Context,
		remoteRepository,
		pullRequest,
		topLevelCommentExpectation,
		string,
	) (*reviewInspection, error)
	ResolveReviewThread(
		context.Context,
		remoteRepository,
		pullRequest,
		reviewThreadExpectation,
	) (*reviewInspection, error)
	RequestReview(
		context.Context,
		remoteRepository,
		pullRequest,
		[]string,
	) (*reviewInspection, error)
}

func resolveForge(
	name string,
	executable string,
	repository remoteRepository,
	runner commandRunner,
	directory string,
) (forge, error) {
	selected := strings.ToLower(strings.TrimSpace(name))
	if selected == "" || selected == "auto" {
		if strings.EqualFold(repository.Host, "github.com") {
			selected = "github"
		} else {
			return nil, fmt.Errorf(
				"no forge adapter recognizes host %q; select a supported adapter explicitly",
				repository.Host,
			)
		}
	}
	if selected != "github" {
		return nil, fmt.Errorf("forge adapter %q is not supported", selected)
	}
	if strings.TrimSpace(executable) == "" {
		executable = "gh"
	}
	return &githubForge{
		executable: executable,
		runner:     runner,
		directory:  directory,
	}, nil
}

func exactPullRequests(
	all []pullRequest,
	repository remoteRepository,
	branch string,
) ([]pullRequest, error) {
	for _, candidate := range all {
		if err := validatePullRequest(candidate); err != nil {
			return nil, fmt.Errorf(
				"forge returned a malformed pull request: %w",
				err,
			)
		}
	}
	if len(all) >= 100 {
		return nil, fmt.Errorf(
			"forge returned 100 pull requests for head %q; refusing a potentially truncated result",
			branch,
		)
	}
	result := make([]pullRequest, 0, len(all))
	for _, candidate := range all {
		if candidate.HeadRefName != branch {
			continue
		}
		if !strings.EqualFold(candidate.HeadRepositoryOwner, repository.Owner) {
			continue
		}
		if !strings.EqualFold(candidate.HeadRepositoryName, repository.Name) {
			if candidate.HeadRepositoryName == "" &&
				!candidate.IsCrossRepository {
				return nil, fmt.Errorf(
					"forge returned an ambiguous pull request for the exact owner and head",
				)
			}
			continue
		}
		if candidate.IsCrossRepository {
			return nil, fmt.Errorf(
				"pull request %d is cross-repository; v1 only supports same-repository heads",
				candidate.Number,
			)
		}
		result = append(result, candidate)
	}
	return result, nil
}

func selectOpenPullRequest(
	all []pullRequest,
) (*pullRequest, []string) {
	var open []pullRequest
	var closed []pullRequest
	for _, candidate := range all {
		if strings.EqualFold(candidate.State, "OPEN") {
			open = append(open, candidate)
		} else {
			closed = append(closed, candidate)
		}
	}
	var refusals []string
	if len(open) > 1 {
		refusals = append(
			refusals,
			fmt.Sprintf("found %d open pull requests for the exact head", len(open)),
		)
		return nil, refusals
	}
	if len(open) == 0 && len(closed) != 0 {
		refusals = append(
			refusals,
			"the exact head already has a closed or merged pull request",
		)
		return nil, refusals
	}
	if len(open) == 1 {
		selected := open[0]
		return &selected, refusals
	}
	return nil, refusals
}
