package xapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/oauth1"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
)

// endpointKey carries a stable endpoint bucket independently of individual IDs.
type endpointKey struct{}

// responseError retains an HTTP failure status through the SDK's error wrapping.
type responseError struct {
	// status identifies the server response without interpreting transport failures.
	status int
	// message contains the bounded, credential-redacted diagnostic.
	message string
}

// Error exposes the same diagnostic as other API failures.
func (e *responseError) Error() string { return e.message }

// Client adapts the generated SDK to owned cleanup data and explicit mutations.
type Client struct {
	// api owns schema-defined requests and response decoding.
	api *sdk.APIClient
	// logger reports read progress without response bodies or credentials.
	logger *slog.Logger
}

// quotaWindow records one account's remaining allowance for one endpoint.
type quotaWindow struct {
	// limit is the server-reported allowance or its documented default.
	limit int
	// remaining counts requests which may be dispatched before reset.
	remaining int
	// reset is the earliest next window boundary.
	reset time.Time
}

// quotaTransport paces requests and bounds network responses before SDK decoding.
type quotaTransport struct {
	// base performs one network request after quota admission.
	base http.RoundTripper
	// credentials sign requests after quota waits; they never reach Lua or plans.
	credentials oauth1.Credentials
	// mu protects the endpoint-specific quota windows.
	mu sync.Mutex
	// windows are private to this client's fixed user credentials.
	windows map[string]*quotaWindow
	// logger reports quota waits and resumption through the command's handler.
	logger *slog.Logger
}

// New configures the shared SDK with private authentication and safe transport.
func New(base string, credentials oauth1.Credentials, logger *slog.Logger) (*Client, error) {
	parsed, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse API base URL: %w", err)
	}
	if parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || !strings.HasSuffix(strings.TrimRight(parsed.Path, "/"), "/2") {
		return nil, fmt.Errorf("API base URL must end in /2 and have a host without credentials, query, or fragment")
	}
	ip := net.ParseIP(parsed.Hostname())
	loopback := parsed.Hostname() == "localhost" || (ip != nil && ip.IsLoopback())
	if parsed.Scheme != "https" && !(parsed.Scheme == "http" && loopback) {
		return nil, fmt.Errorf("API base URL requires HTTPS except for loopback fixtures")
	}
	if err := credentials.Validate(); err != nil {
		return nil, fmt.Errorf("validate OAuth 1.0a credentials: %w", err)
	}
	configuration := sdk.NewConfiguration()
	configuration.Servers = sdk.ServerConfigurations{{URL: strings.TrimSuffix(strings.TrimRight(base, "/"), "/2")}}
	configuration.HTTPClient = &http.Client{
		Transport:     &quotaTransport{base: http.DefaultTransport, credentials: credentials, windows: map[string]*quotaWindow{}, logger: logger},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	return &Client{api: sdk.NewAPIClient(configuration), logger: logger}, nil
}

// requestContext attaches endpoint identity for admission before authentication.
func (c *Client) requestContext(ctx context.Context, endpoint string) context.Context {
	return context.WithValue(ctx, endpointKey{}, endpoint)
}

// reserve waits before dispatch, with independent quotas and prompt cancellation.
func (t *quotaTransport) reserve(ctx context.Context, endpoint string) error {
	waited := false
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("admit %s request: %w", endpoint, err)
		}
		t.mu.Lock()
		now := time.Now()
		w := t.windows[endpoint]
		if w == nil {
			limit := 50
			switch endpoint {
			case "identity":
				limit = 75
			case "timeline":
				limit = 900
			}
			w = &quotaWindow{limit: limit, remaining: limit, reset: now.Add(15 * time.Minute)}
			t.windows[endpoint] = w
		}
		if !now.Before(w.reset) {
			w.remaining = w.limit
			w.reset = now.Add(15 * time.Minute)
		}
		if w.remaining > 0 {
			w.remaining--
			t.mu.Unlock()
			if waited {
				t.logger.InfoContext(ctx, "API quota available", "endpoint", endpoint)
			}
			return nil
		}
		delay := time.Until(w.reset)
		reset := w.reset
		t.mu.Unlock()
		t.logger.InfoContext(ctx, "Waiting for API quota", "endpoint", endpoint, "wait", delay.Round(time.Second), "resume_at", reset.UTC())
		waited = true
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("wait for %s quota: %w", endpoint, ctx.Err())
		case <-timer.C:
		}
	}
}

// observe adopts valid server allowances while retaining conservative defaults.
func (t *quotaTransport) observe(endpoint string, header http.Header) {
	t.mu.Lock()
	defer t.mu.Unlock()
	w := t.windows[endpoint]
	if limit, err := strconv.Atoi(header.Get("X-Rate-Limit-Limit")); err == nil && limit > 0 {
		w.limit = limit
	}
	w.remaining = min(w.remaining, w.limit)
	if remaining, err := strconv.Atoi(header.Get("X-Rate-Limit-Remaining")); err == nil && remaining >= 0 {
		w.remaining = min(remaining, w.limit)
	}
	if reset, err := strconv.ParseInt(header.Get("X-Rate-Limit-Reset"), 10, 64); err == nil && reset > time.Now().Unix() {
		w.reset = time.Unix(reset, 0).Add(time.Second)
	}
}

// RoundTrip performs one admitted request with a separate 30-second network limit.
// API errors and unconfirmed mutations are never replayed or exposed as success.
func (t *quotaTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	endpoint, ok := request.Context().Value(endpointKey{}).(string)
	if !ok {
		return nil, fmt.Errorf("request has no endpoint quota identity")
	}
	if err := t.reserve(request.Context(), endpoint); err != nil {
		return nil, fmt.Errorf("admit request: %w", err)
	}
	ctx, cancel := context.WithTimeout(request.Context(), 30*time.Second)
	defer cancel()
	signed := request.Clone(ctx)
	authorization, err := t.credentials.Sign(signed.Method, signed.URL.String(), signed.URL.Query(), nil, time.Now())
	if err != nil {
		return nil, fmt.Errorf("sign %s request: %w", endpoint, err)
	}
	signed.Header.Set("Authorization", authorization)
	response, err := t.base.RoundTrip(signed)
	if err != nil {
		return nil, fmt.Errorf("send %s request: %w", endpoint, err)
	}
	t.observe(endpoint, response.Header)
	body, readErr := io.ReadAll(io.LimitReader(response.Body, (8<<20)+1))
	closeErr := response.Body.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read %s response: %w", endpoint, readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close %s response: %w", endpoint, closeErr)
	}
	if len(body) > 8<<20 {
		return nil, fmt.Errorf("%s response exceeds 8 MiB", endpoint)
	}
	problems, partial := t.problemSummary(body)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, &responseError{status: response.StatusCode, message: fmt.Sprintf("%s: HTTP %d; %s; retry-after=%q rate-limit-reset=%q (request not retried)", endpoint, response.StatusCode, problems, response.Header.Get("Retry-After"), response.Header.Get("X-Rate-Limit-Reset"))}
	}
	if partial {
		return nil, fmt.Errorf("%s: HTTP %d; %s (partial response; request not retried)", endpoint, response.StatusCode, problems)
	}
	if request.Method == http.MethodDelete {
		var confirmation struct {
			// Data retains presence and nullability lost by generated boolean fields.
			Data map[string]*bool `json:"data"`
		}
		if err := json.Unmarshal(body, &confirmation); err != nil {
			return nil, fmt.Errorf("decode mutation confirmation: %w", err)
		}
		key, want := "deleted", true
		if endpoint == "unrepost" {
			key, want = "retweeted", false
		}
		value := confirmation.Data[key]
		if value == nil || *value != want {
			return nil, fmt.Errorf("%s did not confirm the requested state", endpoint)
		}
	}
	response.Body = io.NopCloser(bytes.NewReader(body))
	return response, nil
}

// problemSummary projects bounded diagnostic fields without exposing response data
// or credential values. Inspecting the wire envelope also retains unknown problem
// types that the generated response union cannot decode.
func (t *quotaTransport) problemSummary(body []byte) (string, bool) {
	var envelope struct {
		// Errors retains each problem independently of its schema discriminator.
		Errors []json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return "no readable problem details", false
	}
	partial := len(envelope.Errors) > 0
	problems := envelope.Errors
	if !partial {
		problems = []json.RawMessage{body}
	}
	var summaries []string
	for _, problem := range problems[:min(len(problems), 3)] {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(problem, &fields); err != nil {
			summaries = append(summaries, "unrecognized problem")
			continue
		}
		var parts []string
		for _, key := range []string{"type", "title", "detail", "resource_type", "resource_id", "parameter", "reason", "code"} {
			var value string
			if err := json.Unmarshal(fields[key], &value); err != nil {
				var number json.Number
				if err := json.Unmarshal(fields[key], &number); err != nil {
					continue
				}
				value = number.String()
			}
			if value == "" {
				continue
			}
			for _, secret := range []string{t.credentials.APISecret, t.credentials.AccessTokenSecret, t.credentials.APIKey, t.credentials.AccessToken} {
				for _, spelling := range []string{secret, url.QueryEscape(secret), url.PathEscape(secret)} {
					if spelling != "" {
						value = strings.ReplaceAll(value, spelling, "[redacted]")
					}
				}
			}
			if runes := []rune(value); len(runes) > 160 {
				value = string(runes[:160]) + "…"
			}
			parts = append(parts, fmt.Sprintf("%s=%q", key, value))
		}
		if len(parts) == 0 {
			parts = append(parts, "no readable problem details")
		}
		summaries = append(summaries, strings.Join(parts, " "))
	}
	if len(problems) > 3 {
		summaries = append(summaries, fmt.Sprintf("%d additional problem(s)", len(problems)-3))
	}
	return strings.Join(summaries, "; "), partial
}

// Account identifies the authenticated user through the generated SDK.
func (c *Client) Account(ctx context.Context) (string, error) {
	c.logger.InfoContext(ctx, "Identifying account")
	result, _, err := c.api.UsersAPI.GetUsersMe(c.requestContext(ctx, "identity")).Execute()
	if err != nil {
		return "", fmt.Errorf("identify authenticated account: %w", err)
	}
	if result == nil || result.Data == nil || len(result.Errors) > 0 || !model.ValidID(result.Data.GetId()) {
		return "", fmt.Errorf("account lookup returned errors or an invalid identifier")
	}
	return result.Data.GetId(), nil
}

// Snapshot caches mixed timeline pages and normalizes only requested owned kinds.
type Snapshot struct {
	// AccountID binds this planning snapshot to the authenticated account.
	AccountID string
	// client provides schema-defined, quota-paced read operations.
	client *Client
	// posts is populated once from the available paginated timeline.
	posts []sdk.Post
	// loaded distinguishes an empty completed timeline from an unread one.
	loaded bool
	// window binds cached pages and kind views to the initial discovery bounds.
	window model.TimeWindow
	// items caches normalized views independently by activity kind.
	items map[model.ItemKind][]model.Item
}

// Inventory identifies the account without fetching timeline data yet.
func (c *Client) Inventory(ctx context.Context) (*Snapshot, error) {
	account, err := c.Account(ctx)
	if err != nil {
		return nil, fmt.Errorf("collect inventory identity: %w", err)
	}
	return &Snapshot{AccountID: account, client: c, items: map[model.ItemKind][]model.Item{}}, nil
}

// collect requests full 100-record pages once, without per-post HTTP requests.
func (s *Snapshot) collect(ctx context.Context, window model.TimeWindow) error {
	if err := window.Validate(); err != nil {
		return fmt.Errorf("validate discovery window: %w", err)
	}
	if s.loaded {
		if !s.window.Start.Equal(window.Start) || !s.window.End.Equal(window.End) {
			return fmt.Errorf("cannot change the discovery window of a loaded snapshot")
		}
		return nil
	}
	seen := map[string]sdk.Post{}
	tokens := map[string]bool{}
	token := ""
	pages := 0
	s.client.logger.InfoContext(ctx, "Fetching timeline", "account", s.AccountID)
	for {
		request := s.client.api.UsersAPI.GetUsersPosts(s.client.requestContext(ctx, "timeline"), s.AccountID).
			MaxResults(100).PostFields([]string{"id", "text", "author_id", "created_at", "referenced_posts", "edit_history_post_ids", "public_metrics", "note_post", "article", "article_title", "lang", "attachments", "entities"})
		if !window.Start.IsZero() {
			request = request.StartTime(window.Start)
		}
		if !window.End.IsZero() {
			request = request.EndTime(window.End)
		}
		if token != "" {
			request = request.PaginationToken(token)
		}
		page, _, err := request.Execute()
		if err != nil {
			return fmt.Errorf("collect timeline page %d after %d activities: %w", pages+1, len(s.posts), err)
		}
		if page == nil {
			return fmt.Errorf("timeline page %d returned no response after %d activities; planning aborted", pages+1, len(s.posts))
		}
		for _, post := range page.Data {
			id := post.GetId()
			if previous, exists := seen[id]; exists {
				if !reflect.DeepEqual(previous, post) {
					return fmt.Errorf("activity %s changed between timeline pages", id)
				}
				continue
			}
			seen[id] = post
			s.posts = append(s.posts, post)
		}
		pages++
		s.client.logger.InfoContext(ctx, "Fetched timeline page", "page", pages, "records", len(page.Data), "collected", len(s.posts))
		token = ""
		if page.Meta != nil {
			token = page.Meta.GetNextToken()
		}
		if token == "" {
			break
		}
		if tokens[token] {
			return fmt.Errorf("timeline repeated pagination token; planning aborted")
		}
		tokens[token] = true
	}
	s.loaded = true
	s.window = window
	s.client.logger.InfoContext(ctx, "Timeline collected", "pages", pages, "records", len(s.posts))
	return nil
}

// Items returns a complete cached kind view; irrelevant kinds are not normalized.
func (s *Snapshot) Items(ctx context.Context, kind model.ItemKind, window model.TimeWindow) ([]model.Item, error) {
	if err := s.collect(ctx, window); err != nil {
		return nil, fmt.Errorf("load %s inventory: %w", kind, err)
	}
	if items, ok := s.items[kind]; ok {
		return items, nil
	}
	items := []model.Item{}
	for _, post := range s.posts {
		postKind := model.ItemPost
		for _, ref := range post.ReferencedPosts {
			if ref.Type == "retweeted" || ref.Type == "reposted" {
				postKind = model.ItemRepost
			}
		}
		if postKind != kind {
			continue
		}
		if post.AuthorId == nil {
			return nil, fmt.Errorf("timeline omitted activity author; planning aborted")
		}
		if post.GetAuthorId() != s.AccountID || post.Article != nil || post.ArticleTitle != nil {
			continue
		}
		if post.CreatedAt != nil && !window.Contains(*post.CreatedAt) {
			continue
		}
		item, err := normalize(post)
		if err != nil {
			return nil, fmt.Errorf("normalize timeline activity: %w", err)
		}
		items = append(items, item)
	}
	s.items[kind] = items
	s.client.logger.InfoContext(ctx, "Inventory ready", "kind", kind, "activities", len(items))
	return items, nil
}

// normalize maps one owned SDK post without inventing missing metadata.
func normalize(post sdk.Post) (model.Item, error) {
	id := post.GetId()
	if !model.ValidID(id) || !model.ValidID(post.GetAuthorId()) || post.CreatedAt == nil || post.CreatedAt.IsZero() || post.Text == nil {
		return model.Item{}, fmt.Errorf("invalid or incomplete owned activity %q", id)
	}
	item := model.Item{ID: id, AuthorID: post.GetAuthorId(), EditHistoryIDs: post.EditHistoryPostIds, Kind: model.ItemPost, CreatedAt: *post.CreatedAt, Text: post.GetText()}
	if post.NotePost != nil {
		item.Text = post.NotePost.Text
	}
	item.Language = post.Lang
	if post.Attachments != nil {
		item.Attachments = &model.Attachments{MediaKeys: post.Attachments.MediaKeys, PollIDs: post.Attachments.PollIds, MediaSourceIDs: post.Attachments.MediaSourceTweetId}
	}
	var entities any = post.Entities
	if post.NotePost != nil {
		entities = post.NotePost.Entities.Get()
	}
	if err := projectMetadata(entities, &item.Entities); err != nil {
		return model.Item{}, fmt.Errorf("normalize activity %s entities: %w", id, err)
	}
	if post.PublicMetrics != nil {
		content, err := json.Marshal(post.PublicMetrics)
		if err != nil {
			return model.Item{}, fmt.Errorf("encode activity %s metrics: %w", id, err)
		}
		if err := json.Unmarshal(content, &item.Metrics); err != nil {
			return model.Item{}, fmt.Errorf("decode activity %s counters: %w", id, err)
		}
	}
	for _, ref := range post.ReferencedPosts {
		if !model.ValidID(ref.Id) {
			return model.Item{}, fmt.Errorf("invalid reference on activity %s", id)
		}
		kind := model.ReferenceKind(ref.Type)
		switch kind {
		case model.ReferenceQuote, model.ReferenceReply:
		case model.ReferenceRepost, "reposted":
			if item.SourceID != "" {
				return model.Item{}, fmt.Errorf("ambiguous repost source on activity %s", id)
			}
			kind = model.ReferenceRepost
			item.Kind = model.ItemRepost
			item.SourceID = ref.Id
		default:
			return model.Item{}, fmt.Errorf("unknown reference kind %q on activity %s", ref.Type, id)
		}
		item.References = append(item.References, model.Reference{ID: ref.Id, Kind: kind})
	}
	return item, nil
}

// projectMetadata maps the SDK's ordinary or long-post entity schema into owned data.
func projectMetadata(source any, destination any) error {
	content, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("encode SDK metadata: %w", err)
	}
	if err := json.Unmarshal(content, destination); err != nil {
		return fmt.Errorf("decode owned metadata: %w", err)
	}
	return nil
}

// Execute sends exactly one typed SDK mutation for the validated action.
// A successful response or HTTP 404 confirms completion; other errors propagate.
func (c *Client) Execute(ctx context.Context, account string, action model.Action) error {
	if !model.ValidID(account) || !model.ValidID(action.TargetID) {
		return fmt.Errorf("invalid execution identifier")
	}
	switch action.Kind {
	case model.ActionDeletePost:
		result, _, err := c.api.PostsAPI.DeletePosts(c.requestContext(ctx, "delete"), action.TargetID).Execute()
		if err != nil {
			return c.mutationError(ctx, action, err)
		}
		if result == nil || result.Data == nil || !result.Data.Deleted || len(result.Errors) > 0 {
			return fmt.Errorf("post deletion returned errors or lacked confirmation")
		}
	case model.ActionUnrepost:
		result, _, err := c.api.UsersAPI.UnrepostPost(c.requestContext(ctx, "unrepost"), account, action.TargetID).Execute()
		if err != nil {
			return c.mutationError(ctx, action, err)
		}
		if result == nil || result.Data == nil || result.Data.Retweeted || len(result.Errors) > 0 {
			return fmt.Errorf("repost removal returned errors or lacked confirmation")
		}
	default:
		return fmt.Errorf("unsupported action %q", action.Kind)
	}
	return nil
}

// mutationError treats only an HTTP 404 from an explicit mutation as already absent.
// Authentication, quota, transport, and server failures remain errors for the caller.
func (c *Client) mutationError(ctx context.Context, action model.Action, err error) error {
	var response *responseError
	if errors.As(err, &response) && response.status == http.StatusNotFound {
		c.logger.InfoContext(ctx, "Skipping target already absent", "kind", action.Kind, "target", action.TargetID)
		return nil
	}
	return fmt.Errorf("execute %s for %s: %w", action.Kind, action.TargetID, err)
}
