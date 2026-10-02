package e2e

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draft"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

// Processing failure modes, specified before implementation:
//   - a processing or failed upload is attached before it is ready;
//   - simple ready images incur unnecessary status requests;
//   - status polling ignores server delay, its overall deadline, or cancellation;
//   - malformed states, mismatched identifiers, or partial errors are accepted;
//   - status GETs omit query parameters from their OAuth signature;
//   - status rejection or redirects cause retries, lose quota evidence, or draft.
//
// Each case converts a real post and runs its publisher against a local service.
func TestMediaProcessingBeforeDraft(t *testing.T) {
	const (
		ready      = `{"data":{"id":"1900000000000000001","expires_after_secs":3600}}`
		pending    = `{"data":{"id":"1900000000000000001","expires_after_secs":3600,"processing_info":{"state":"pending","check_after_secs":1}}}`
		succeeded  = `{"data":{"id":"1900000000000000001","processing_info":{"state":"succeeded"}}}`
		inProgress = `{"data":{"id":"1900000000000000001","expires_after_secs":3600,"processing_info":{"state":"in_progress","check_after_secs":0}}}`
		failed     = `{"data":{"id":"1900000000000000001","processing_info":{"state":"failed"}}}`
	)
	for _, test := range []struct {
		name         string
		upload       string
		statuses     []string
		statusCode   int
		timeout      time.Duration
		hangStatus   bool
		wantStatuses int
		wantDraft    bool
		wantDeadline bool
		wantError    string
	}{
		{name: "static_needs_no_status", upload: ready, wantDraft: true},
		{name: "succeeded_needs_no_status", upload: `{"data":{"id":"1900000000000000001","expires_after_secs":3600,"processing_info":{"state":"succeeded"}}}`, wantDraft: true},
		{name: "pending_then_succeeded", upload: pending, statuses: []string{succeeded}, wantStatuses: 1, wantDraft: true},
		{name: "in_progress_then_pending_then_succeeded", upload: inProgress, statuses: []string{pending, succeeded}, wantStatuses: 2, wantDraft: true},
		{name: "status_may_omit_id", upload: pending, statuses: []string{`{"data":{"processing_info":{"state":"succeeded"}}}`}, wantStatuses: 1, wantDraft: true},
		{name: "upload_failed", upload: failed, wantError: "processing"},
		{name: "upload_unknown_state", upload: `{"data":{"id":"1900000000000000001","processing_info":{"state":"unexpected"}}}`, wantError: "processing"},
		{name: "upload_missing_state", upload: `{"data":{"id":"1900000000000000001","processing_info":{}}}`, wantError: "processing"},
		{name: "negative_check_delay", upload: `{"data":{"id":"1900000000000000001","processing_info":{"state":"pending","check_after_secs":-1}}}`, wantError: "processing"},
		{name: "upload_partial_error", upload: `{"data":{"id":"1900000000000000001","expires_after_secs":3600},"errors":[{"detail":"fixture upload rejected"}]}`, wantError: "fixture upload rejected"},
		{name: "status_failed", upload: pending, statuses: []string{failed}, wantStatuses: 1, wantError: "processing"},
		{name: "status_missing_processing", upload: pending, statuses: []string{ready}, wantStatuses: 1, wantError: "processing"},
		{name: "status_changed_id", upload: pending, statuses: []string{`{"data":{"id":"1900000000000000002","processing_info":{"state":"succeeded"}}}`}, wantStatuses: 1, wantError: "identifier"},
		{name: "status_partial_error", upload: pending, statuses: []string{`{"data":{"id":"1900000000000000001","processing_info":{"state":"succeeded"}},"errors":[{"detail":"fixture status rejected"}]}`}, wantStatuses: 1, wantError: "fixture status rejected"},
		{name: "status_rejection_not_retried", upload: pending, statuses: []string{`{"detail":"fixture status rejection"}`}, statusCode: 503, wantStatuses: 1, wantError: "503"},
		{name: "status_redirect_not_followed", upload: pending, statuses: []string{`{"detail":"fixture status redirect"}`}, statusCode: 307, wantStatuses: 1, wantError: "307"},
		{name: "server_delay_exceeds_budget", upload: pending, timeout: 40 * time.Millisecond, wantDeadline: true},
		{name: "processing_never_finishes", upload: inProgress, statuses: []string{inProgress}, timeout: 1500 * time.Millisecond, wantStatuses: 1, wantDeadline: true},
		{name: "status_request_deadline", upload: pending, hangStatus: true, timeout: 1200 * time.Millisecond, wantStatuses: 1, wantDeadline: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, postDir, source := bannerPost(t, "images: [banner.png]\n", "Article body.")
			artifactPath, _ := bannerArtifact(t, root, postDir, source)
			artifact, err := draft.ReadArtifact(artifactPath)
			if err != nil {
				t.Fatalf("read converted artifact: %v", err)
			}
			var mu sync.Mutex
			var requests []string
			var statusTimes []time.Time
			var uploadTime time.Time
			statusCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				requests = append(requests, r.Method+" "+r.URL.RequestURI())
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				body := ""
				switch {
				case r.Method == http.MethodPost && r.URL.Path == "/media/upload":
					mu.Lock()
					uploadTime = time.Now()
					mu.Unlock()
					body = test.upload
				case r.Method == http.MethodGet && r.URL.Path == "/media/upload":
					if r.URL.Query().Get("media_id") != "1900000000000000001" || r.URL.Query().Get("command") != "STATUS" {
						t.Error("status request has incorrect media identifier or command")
					}
					if err := verifyStatusSignature(r); err != nil {
						t.Errorf("status authorization: %v", err)
					}
					mu.Lock()
					statusTimes = append(statusTimes, time.Now())
					index := statusCount
					statusCount++
					mu.Unlock()
					if test.hangStatus {
						<-r.Context().Done()
						return
					}
					if index >= len(test.statuses) {
						t.Error("unexpected additional status request")
						body = failed
					} else {
						body = test.statuses[index]
					}
					if test.statusCode != 0 {
						w.Header().Set("X-Rate-Limit-Remaining", "0")
						w.Header().Set("X-Rate-Limit-Reset", "1791055658")
						w.Header().Set("Location", "/unexpected-redirect")
						w.WriteHeader(test.statusCode)
					}
				case r.Method == http.MethodPost && r.URL.Path == "/articles/draft":
					var payload xapi.DraftRequest
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.CoverMedia == nil || payload.CoverMedia.MediaID != "1900000000000000001" {
						t.Error("draft did not attach the ready media identifier")
					}
					body = `{"data":{"id":"fixture-article"}}`
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					body = `{"detail":"unexpected request"}`
				}
				if _, err := io.WriteString(w, body); err != nil {
					t.Errorf("write media fixture response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := xapi.New(xapi.Credentials{APIKey: "fixture-key", APISecret: "fixture-secret", AccessToken: "fixture-token", AccessTokenSecret: "fixture-token-secret"})
			client.BaseURL = server.URL
			client.MediaProcessingTimeout = 10 * time.Second
			if test.timeout != 0 {
				client.MediaProcessingTimeout = test.timeout
			}
			started := time.Now()
			_, err = draft.New(client, root, filepath.Join(root, "cache.json")).CreateDraft(artifact)
			elapsed := time.Since(started)
			message := ""
			if err != nil {
				message = err.Error()
			}
			if test.wantDraft && err != nil {
				t.Errorf("ready media failed: %v", err)
			}
			if !test.wantDraft && err == nil {
				t.Error("unready media created a draft")
			}
			if test.wantError != "" && !strings.Contains(message, test.wantError) {
				t.Errorf("error %q lacks %q", message, test.wantError)
			}
			if test.wantDeadline && (!errors.Is(err, context.DeadlineExceeded) || elapsed > test.timeout+time.Second) {
				t.Errorf("processing deadline not enforced: elapsed=%s error=%v", elapsed, err)
			}
			if test.statusCode != 0 {
				var apiError *xapi.APIError
				if !errors.As(err, &apiError) || apiError.Status != test.statusCode || apiError.Headers.Get("X-Rate-Limit-Remaining") != "0" {
					t.Error("status rejection lost API status or quota headers")
				}
			}
			mu.Lock()
			attempts := append([]string(nil), requests...)
			observedStatusTimes := append([]time.Time(nil), statusTimes...)
			observedUploadTime := uploadTime
			mu.Unlock()
			want := []string{"POST /media/upload"}
			for range test.wantStatuses {
				want = append(want, "GET /media/upload?command=STATUS&media_id=1900000000000000001")
			}
			if test.wantDraft {
				want = append(want, "POST /articles/draft")
			}
			if !reflect.DeepEqual(attempts, want) {
				t.Errorf("request sequence=%v, want %v", attempts, want)
			}
			previous := observedUploadTime
			for _, at := range observedStatusTimes {
				if at.Sub(previous) < 900*time.Millisecond {
					t.Error("status polled before the server delay or minimum polling interval")
				}
				previous = at
			}
			writeRateLimitEvidence(t, message, nil, attempts)
		})
	}
}

// Verify the service's expected signature independently from the client signer.
// The signed query is reconstructed from the actual HTTP request, so omitting
// command or media_id from the OAuth parameter string causes a fixture failure.
func verifyStatusSignature(request *http.Request) error {
	parameters := request.URL.Query()
	header := strings.TrimPrefix(request.Header.Get("Authorization"), "OAuth ")
	for _, field := range strings.Split(header, ", ") {
		name, rawValue, ok := strings.Cut(field, "=")
		if !ok {
			return fmt.Errorf("malformed OAuth field")
		}
		value, err := url.QueryUnescape(strings.Trim(rawValue, `"`))
		if err != nil {
			return fmt.Errorf("decode OAuth field: %w", err)
		}
		parameters.Set(name, value)
	}
	supplied := parameters.Get("oauth_signature")
	parameters.Del("oauth_signature")
	encode := func(value string) string { return strings.ReplaceAll(url.QueryEscape(value), "+", "%20") }
	keys := make([]string, 0, len(parameters))
	for key := range parameters {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var pairs []string
	for _, key := range keys {
		pairs = append(pairs, encode(key)+"="+encode(parameters.Get(key)))
	}
	base := "GET&" + encode("http://"+request.Host+request.URL.Path) + "&" + encode(strings.Join(pairs, "&"))
	mac := hmac.New(sha1.New, []byte("fixture-secret&fixture-token-secret"))
	if _, err := io.WriteString(mac, base); err != nil {
		return fmt.Errorf("sign fixture request: %w", err)
	}
	if !hmac.Equal([]byte(supplied), []byte(base64.StdEncoding.EncodeToString(mac.Sum(nil)))) {
		return fmt.Errorf("signature excludes or changes request parameters")
	}
	return nil
}
