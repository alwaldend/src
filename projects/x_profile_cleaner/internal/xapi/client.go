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
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/oauth1"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
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

// Client exposes only explicit X mutations through the generated SDK.
type Client struct {
	// api owns schema-defined requests and response decoding.
	api *sdk.APIClient
	// logger reports mutation progress without response bodies or credentials.
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
	if request.Method != http.MethodDelete {
		return nil, fmt.Errorf("X read requests are disabled; use an export archive")
	}
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

// ValidatePlan admits the complete supported executor set before any mutation.
// Duplicate operations on the same endpoint and target are rejected upfront.
func ValidatePlan(value *workerpb.Plan) error {
	if err := plan.Validate(value); err != nil {
		return fmt.Errorf("validate executable plan: %w", err)
	}
	seen := map[string]bool{}
	for _, operation := range value.Operations {
		if operation.GetDeletePost() == nil && operation.GetUnrepost() == nil {
			return fmt.Errorf("unsupported execution operation %q", plan.OperationName(operation))
		}
		key := plan.OperationName(operation) + ":" + plan.TargetID(operation)
		if seen[key] {
			return fmt.Errorf("duplicate execution operation")
		}
		seen[key] = true
	}
	return nil
}

// Execute sends exactly one typed SDK mutation for the validated action.
// A successful response or HTTP 404 confirms completion; other errors propagate.
func (c *Client) Execute(ctx context.Context, account string, action *workerpb.Operation) error {
	if !model.ValidID(account) {
		return fmt.Errorf("invalid execution identifier")
	}
	if err := ValidatePlan(&workerpb.Plan{Operations: []*workerpb.Operation{action}}); err != nil {
		return fmt.Errorf("validate execution payload: %w", err)
	}
	switch plan.OperationName(action) {
	case "delete_post":
		result, _, err := c.api.PostsAPI.DeletePosts(c.requestContext(ctx, "delete"), plan.TargetID(action)).Execute()
		if err != nil {
			if failure := c.mutationError(ctx, action, err); failure != nil {
				return fmt.Errorf("perform X mutation: %w", failure)
			}
			return nil
		}
		if result == nil || result.Data == nil || !result.Data.Deleted || len(result.Errors) > 0 {
			return fmt.Errorf("post deletion returned errors or lacked confirmation")
		}
	case "unrepost":
		result, _, err := c.api.UsersAPI.UnrepostPost(c.requestContext(ctx, "unrepost"), account, plan.TargetID(action)).Execute()
		if err != nil {
			if failure := c.mutationError(ctx, action, err); failure != nil {
				return fmt.Errorf("perform X mutation: %w", failure)
			}
			return nil
		}
		if result == nil || result.Data == nil || result.Data.Retweeted || len(result.Errors) > 0 {
			return fmt.Errorf("repost removal returned errors or lacked confirmation")
		}
	default:
		return fmt.Errorf("unsupported action %q", plan.OperationName(action))
	}
	return nil
}

// mutationError treats only an HTTP 404 from an explicit mutation as already absent.
// Authentication, quota, transport, and server failures remain errors for the caller.
func (c *Client) mutationError(ctx context.Context, action *workerpb.Operation, err error) error {
	var response *responseError
	if errors.As(err, &response) && response.status == http.StatusNotFound {
		c.logger.InfoContext(ctx, "Skipping target already absent", "kind", plan.OperationName(action), "target", plan.TargetID(action))
		return nil
	}
	return fmt.Errorf("execute %s for %s: %w", plan.OperationName(action), plan.TargetID(action), err)
}
