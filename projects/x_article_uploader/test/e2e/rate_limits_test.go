package e2e

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draft"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

// Rate-limit failure modes, specified before implementation:
//   - 429 or 503 loses either window's headers through request-error wrapping;
//   - Unix reset values remain unreadable or lose their original value;
//   - nonexhausted windows or Retry-After alone invent an eligibility time;
//   - two exhausted windows use the earlier reset, or incomplete, malformed,
//     or ambiguous duplicate values produce an unjustified combined time;
//   - elapsed resets or an overflowing Retry-After invent a future boundary,
//     or a valid later Retry-After fails to extend an established quota boundary;
//   - a truncated rejected response loses the original status and rate headers;
//   - a diagnostic copies authorization, cookies, or unapproved response headers;
//   - a failed upload proceeds to draft creation, or either operation retries.
// Each case converts a post, reads its artifact, and creates a draft through the
// real publisher and HTTP client. Only the remote response is substituted.

func rateLimitHeaders(endpointRemaining, endpointReset, dailyRemaining, dailyReset, retryAfter string) http.Header {
	headers := make(http.Header)
	for name, value := range map[string]string{
		"X-Rate-Limit-Limit":            "40000",
		"X-Rate-Limit-Remaining":        endpointRemaining,
		"X-Rate-Limit-Reset":            endpointReset,
		"X-User-Limit-24hour-Limit":     "10",
		"X-User-Limit-24hour-Remaining": dailyRemaining,
		"X-User-Limit-24hour-Reset":     dailyReset,
		"Retry-After":                   retryAfter,
	} {
		if value != "" {
			headers.Set(name, value)
		}
	}
	return headers
}

func TestFailedDraftPreservesRateLimitEvidence(t *testing.T) {
	const (
		earlyReset = "1700000300"
		lateReset  = "1700003600"
		earlyUTC   = "2023-11-14 22:18:20 UTC"
		lateUTC    = "2023-11-14 23:13:20 UTC"
		retryUTC   = "2023-11-14 22:23:20 UTC"
		retryLabel = "Rate-limit retry not before: "
	)
	duplicateRemaining := rateLimitHeaders("0", earlyReset, "2", lateReset, "")
	duplicateRemaining.Add("X-Rate-Limit-Remaining", "1")
	duplicateReset := rateLimitHeaders("0", earlyReset, "2", lateReset, "")
	duplicateReset.Add("X-Rate-Limit-Reset", lateReset)
	malformedLimit := rateLimitHeaders("0", earlyReset, "2", lateReset, "")
	malformedLimit.Set("X-Rate-Limit-Limit", "not-an-integer")

	for _, test := range []struct {
		name      string
		status    int
		headers   http.Header
		notBefore string
		upload    bool
		truncated bool
	}{
		{"endpoint_exhausted", 429, rateLimitHeaders("0", earlyReset, "9", lateReset, ""), earlyUTC, false, false},
		{"daily_exhausted", 429, rateLimitHeaders("39999", earlyReset, "0", lateReset, ""), lateUTC, false, false},
		{"both_exhausted_endpoint_later", 429, rateLimitHeaders("0", lateReset, "0", earlyReset, ""), lateUTC, false, false},
		{"both_exhausted_daily_later", 429, rateLimitHeaders("0", earlyReset, "0", lateReset, ""), lateUTC, false, false},
		{"unavailable_with_one_slot", 503, rateLimitHeaders("40000", earlyReset, "1", lateReset, ""), "", false, false},
		{"missing_reset", 429, rateLimitHeaders("0", "", "2", lateReset, ""), "", false, false},
		{"malformed_reset", 429, rateLimitHeaders("0", "not-a-time", "2", lateReset, ""), "", false, false},
		{"missing_remaining", 429, rateLimitHeaders("", earlyReset, "2", lateReset, ""), "", false, false},
		{"malformed_remaining", 429, rateLimitHeaders("unknown", earlyReset, "2", lateReset, ""), "", false, false},
		{"both_exhausted_missing_one_reset", 429, rateLimitHeaders("0", earlyReset, "0", "", ""), "", false, false},
		{"both_exhausted_malformed_one_reset", 429, rateLimitHeaders("0", earlyReset, "0", "not-a-time", ""), "", false, false},
		{"duplicate_remaining", 429, duplicateRemaining, "", false, false},
		{"duplicate_reset", 429, duplicateReset, "", false, false},
		{"malformed_limit_does_not_change_remaining", 429, malformedLimit, earlyUTC, false, false},
		{"elapsed_reset", 429, rateLimitHeaders("0", "1699999900", "2", lateReset, ""), "", false, false},
		{"reset_at_receipt", 429, rateLimitHeaders("0", "1700000000", "2", lateReset, ""), "", false, false},
		{"both_exhausted_one_elapsed", 429, rateLimitHeaders("0", earlyReset, "0", "1699999900", ""), "", false, false},
		{"retry_after_is_separate", 503, rateLimitHeaders("40000", earlyReset, "1", lateReset, "120"), "", false, false},
		{"retry_after_alone", 503, http.Header{"Retry-After": []string{"120"}}, "", false, false},
		{"retry_after_delay_extends_bound", 429, rateLimitHeaders("0", earlyReset, "2", lateReset, "600"), retryUTC, false, false},
		{"retry_after_date_extends_bound", 429, rateLimitHeaders("0", earlyReset, "2", lateReset, "Tue, 14 Nov 2023 22:23:20 GMT"), retryUTC, false, false},
		{"retry_after_shorter_than_reset", 429, rateLimitHeaders("0", earlyReset, "2", lateReset, "120"), earlyUTC, false, false},
		{"retry_after_malformed", 429, rateLimitHeaders("0", earlyReset, "2", lateReset, "not-a-delay"), earlyUTC, false, false},
		{"retry_after_overflows_integer", 429, rateLimitHeaders("0", earlyReset, "2", lateReset, "999999999999999999999999"), earlyUTC, false, false},
		{"retry_after_overflows_duration", 429, rateLimitHeaders("0", earlyReset, "2", lateReset, "9223372036854775807"), earlyUTC, false, false},
		{"truncated_error_body", 429, rateLimitHeaders("0", earlyReset, "2", lateReset, ""), earlyUTC, false, true},
		{"upload_failure_stops_before_draft", 503, rateLimitHeaders("0", earlyReset, "2", lateReset, ""), earlyUTC, true, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			metadata := ""
			wantPath := "/articles/draft"
			wantOperation := "articles/draft"
			if test.upload {
				metadata = "images: [banner.png]\n"
				wantPath = "/media/upload"
				wantOperation = "upload image"
			}
			root, postDir, source := bannerPost(t, metadata, "Article text.")
			artifactPath, _ := bannerArtifact(t, root, postDir, source)
			artifact, err := draft.ReadArtifact(artifactPath)
			if err != nil {
				t.Fatalf("read converted artifact: %v", err)
			}

			var mu sync.Mutex
			var requests []string
			const responseBody = `{"detail":"fixture rejection"}`
			privateHeaders := map[string]string{
				"Authorization":    "fixture-response-authorization",
				"Set-Cookie":       "fixture-response-cookie",
				"X-Internal-Token": "fixture-response-internal-token",
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				requests = append(requests, r.Method+" "+r.URL.Path)
				mu.Unlock()
				for name, values := range test.headers {
					for _, value := range values {
						w.Header().Add(name, value)
					}
				}
				for name, value := range privateHeaders {
					w.Header().Set(name, value)
				}
				w.Header().Set("Content-Type", "application/json")
				if test.truncated {
					w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)+10))
				}
				w.WriteHeader(test.status)
				if _, err := io.WriteString(w, responseBody); err != nil {
					t.Errorf("write rejected response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := &xapi.Client{
				BaseURL: server.URL,
				Credentials: xapi.Credentials{
					APIKey: "fixture-key", APISecret: "fixture-secret",
					AccessToken: "fixture-token", AccessTokenSecret: "fixture-token-secret",
				},
				HTTPClient: server.Client(),
				Now:        func() time.Time { return time.Unix(1700000000, 0) },
			}
			_, err = draft.New(client, root, filepath.Join(root, "cache.json")).CreateDraft(artifact)
			if err == nil {
				t.Fatal("rejected HTTP request returned success")
			}
			var apiError *xapi.APIError
			if !errors.As(err, &apiError) {
				t.Fatalf("wrapped error lost APIError: %T", err)
			}
			if apiError.Status != test.status || apiError.Body != responseBody || apiError.Operation != wantOperation {
				t.Errorf("API error lost the original status, body, or operation")
			}
			message := err.Error()
			if !strings.Contains(message, responseBody) || !strings.Contains(message, strconv.Itoa(test.status)) || !strings.Contains(message, wantOperation) {
				t.Error("printed error lost the original status, body, or operation")
			}

			// Decode through JSON so this preimplementation test also compiles
			// before APIError exposes the structured Headers field.
			encoded, marshalErr := json.Marshal(apiError)
			if marshalErr != nil {
				t.Fatalf("encode API error: %v", marshalErr)
			}
			var recorded struct {
				Headers http.Header
			}
			if decodeErr := json.Unmarshal(encoded, &recorded); decodeErr != nil {
				t.Fatalf("decode API error evidence: %v", decodeErr)
			}
			if !reflect.DeepEqual(recorded.Headers, test.headers) {
				t.Error("structured error headers differ from the response's allowlisted headers")
			}
			for name, values := range test.headers {
				if !strings.Contains(strings.ToLower(message), strings.ToLower(name)) {
					t.Errorf("printed error omitted %s", name)
				}
				for _, value := range values {
					if !strings.Contains(message, value) {
						t.Errorf("printed error omitted the original value of %s", name)
					}
					if strings.HasSuffix(name, "-Reset") {
						if seconds, parseErr := strconv.ParseInt(value, 10, 64); parseErr == nil {
							wantUTC := time.Unix(seconds, 0).UTC().Format("2006-01-02 15:04:05 UTC")
							if !strings.Contains(message, value+" ("+wantUTC+")") {
								t.Errorf("printed %s lacks the original reset with readable UTC time", name)
							}
						}
					}
				}
			}
			if test.notBefore == "" {
				if strings.Contains(message, retryLabel) {
					t.Error("printed an unjustified rate-limit retry time")
				}
			} else if !strings.Contains(message, retryLabel+test.notBefore) {
				t.Errorf("printed error lacks the rate-limit retry time %s", test.notBefore)
			}
			for name, value := range privateHeaders {
				if recorded.Headers.Get(name) != "" || strings.Contains(message, value) || strings.Contains(string(encoded), value) {
					t.Fatalf("error evidence exposed unapproved response header %s", name)
				}
			}
			mu.Lock()
			attempts := append([]string(nil), requests...)
			mu.Unlock()
			if !reflect.DeepEqual(attempts, []string{"POST " + wantPath}) {
				t.Errorf("rejection caused extra or unexpected requests: %v", attempts)
			}
			writeRateLimitEvidence(t, message, recorded.Headers, attempts)
		})
	}
}

// Redirect failure modes, specified before implementation:
//   - 307/308 replays a draft or media POST at the redirect destination;
//   - 301/302/303 sends an unintended GET after the original POST;
//   - following a media redirect allows an additional draft request;
//   - refusing a redirect discards its original status, body, or rate headers.
//
// Exercise the production client constructor, since a fixture-supplied client
// would not verify the redirect policy used by the real command.
func TestRedirectStopsAfterOneRequest(t *testing.T) {
	for _, status := range []int{301, 302, 303, 307, 308} {
		for _, upload := range []bool{false, true} {
			name := "draft"
			if upload {
				name = "media"
			}
			t.Run(name+"_"+strconv.Itoa(status), func(t *testing.T) {
				metadata := ""
				wantPath := "/articles/draft"
				wantOperation := "articles/draft"
				if upload {
					metadata = "images: [banner.png]\n"
					wantPath = "/media/upload"
					wantOperation = "upload image"
				}
				root, postDir, source := bannerPost(t, metadata, "Article text.")
				artifactPath, _ := bannerArtifact(t, root, postDir, source)
				artifact, err := draft.ReadArtifact(artifactPath)
				if err != nil {
					t.Fatalf("read converted artifact: %v", err)
				}

				var mu sync.Mutex
				var requests []string
				const responseBody = `{"detail":"fixture redirect"}`
				headers := rateLimitHeaders("39999", "1700000300", "9", "1700003600", "")
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mu.Lock()
					requests = append(requests, r.Method+" "+r.URL.Path)
					mu.Unlock()
					w.Header().Set("Content-Type", "application/json")
					body := `{"data":{"id":"fixture-draft"}}`
					if r.URL.Path == wantPath {
						for name, values := range headers {
							for _, value := range values {
								w.Header().Add(name, value)
							}
						}
						w.Header().Set("Location", "/redirected")
						w.WriteHeader(status)
						body = responseBody
					} else if upload && r.URL.Path == "/redirected" {
						body = `{"data":{"id":"fixture-media","expires_after_secs":3600}}`
					}
					if _, err := io.WriteString(w, body); err != nil {
						t.Errorf("write redirect fixture response: %v", err)
					}
				}))
				t.Cleanup(server.Close)
				client := xapi.New(xapi.Credentials{
					APIKey: "fixture-key", APISecret: "fixture-secret",
					AccessToken: "fixture-token", AccessTokenSecret: "fixture-token-secret",
				})
				client.BaseURL = server.URL
				client.Now = func() time.Time { return time.Unix(1700000000, 0) }
				_, err = draft.New(client, root, filepath.Join(root, "cache.json")).CreateDraft(artifact)
				var apiError *xapi.APIError
				message := ""
				if err != nil {
					message = err.Error()
				}
				if !errors.As(err, &apiError) {
					t.Errorf("redirect did not return the original API rejection: %v", err)
				} else if apiError.Status != status || apiError.Body != responseBody || apiError.Operation != wantOperation || !reflect.DeepEqual(apiError.Headers, headers) {
					t.Error("redirect rejection lost its original status, body, operation, or rate headers")
				}
				mu.Lock()
				attempts := append([]string(nil), requests...)
				mu.Unlock()
				if !reflect.DeepEqual(attempts, []string{"POST " + wantPath}) {
					t.Errorf("redirect caused extra or unexpected requests: %v", attempts)
				}
				writeRateLimitEvidence(t, message, headers, attempts)
			})
		}
	}
}

func writeRateLimitEvidence(t *testing.T, message string, headers http.Header, requests []string) {
	t.Helper()
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		return
	}
	encoded, err := json.MarshalIndent(map[string]any{
		"error": message, "headers": headers, "requests": requests,
	}, "", "  ")
	if err != nil {
		t.Fatalf("encode rate-limit evidence: %v", err)
	}
	name := strings.ReplaceAll(t.Name(), "/", "-") + ".json"
	if err := os.WriteFile(filepath.Join(output, name), append(encoded, '\n'), 0o644); err != nil {
		t.Fatalf("write rate-limit evidence %q: %v", name, err)
	}
}
