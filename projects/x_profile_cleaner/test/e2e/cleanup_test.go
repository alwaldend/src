package e2e

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// cleanerBinary identifies the Bazel-provided command under test.
var cleanerBinary = flag.String("cleaner", "", "cleaner executable runfile")

// fixtureCredentials contains synthetic OAuth fields, including escaped secrets.
var fixtureCredentials = map[string]string{
	"X_API_KEY":             "fixture-api-key",
	"X_API_SECRET":          "fixture/api+secret&=",
	"X_ACCESS_TOKEN":        "fixture-access-token",
	"X_ACCESS_TOKEN_SECRET": "fixture/token+secret&=",
}

// fixture owns the HTTP observation state for one isolated command scenario.
type fixture struct {
	// mu protects requests made by the server goroutine.
	mu sync.Mutex
	// mutations records completed HTTP mutation paths.
	mutations []string
	// mutationAttempts records every mutation request, including rejected ones.
	mutationAttempts []string
	// allowedMutations contains only owned post and repost removal routes.
	allowedMutations map[string]bool
	// reads records timeline pagination requests.
	reads int
	// requests counts all HTTP requests, including rejected authentication.
	requests int
	// nonces rejects reused OAuth signatures across reads and mutations.
	nonces map[string]bool
	// windowStart is the expected inclusive API discovery boundary, or empty.
	windowStart string
	// windowEnd is the expected exclusive API discovery boundary, or empty.
	windowEnd string
	// ignoreWindow tests local rejection of out-of-window server records.
	ignoreWindow bool
	// relativeWindow records rolling bounds for verification against evaluation time.
	relativeWindow bool
	// requiredFields contains metadata requested by a richer policy fixture.
	requiredFields []string
	// mutationBarrier requires both endpoint workers to enter before either responds.
	mutationBarrier chan string
	// holdPostUntilCancelled keeps a post request in flight while the other endpoint fails.
	holdPostUntilCancelled bool
	// releaseMutations releases requests after the concurrency assertion.
	releaseMutations chan struct{}
	// account is the authenticated account returned by the fixture.
	account string
	// failMutation makes DELETE return a rate-limit error.
	failMutation bool
	// failMutationAt rejects the numbered mutation attempt, counting from one.
	failMutationAt int
	// malformed makes an owned inventory record incomplete.
	malformed bool
	// repeatPage makes pagination return a repeated token.
	repeatPage bool
	// timelineErrors supplies X problem details on the second page.
	timelineErrors []map[string]any
	// timelineStatus optionally overrides the second page's HTTP status.
	timelineStatus int
	// quotaMutation makes the first mutation exhaust its endpoint allowance.
	quotaMutation bool
	// quotaTimeline makes the first page exhaust the timeline allowance.
	quotaTimeline bool
	// resetAt records the fixture's announced quota reset time.
	resetAt time.Time
	// independentBeforeReset records another endpoint dispatching without waiting.
	independentBeforeReset bool
	// firstMutationPath identifies the endpoint whose allowance was exhausted.
	firstMutationPath string
	// mutationResponse overrides confirmation with an intentionally invalid body.
	mutationResponse string
	// mutationStatus overrides the status of a known mutation target.
	mutationStatus int
	// pages is the timeline fixture split across two responses.
	pages [][]map[string]any
}

// newFixture constructs owned records, foreign records, and graph branches.
func newFixture(t *testing.T) (*fixture, *httptest.Server) {
	t.Helper()
	f := &fixture{account: "100", allowedMutations: map[string]bool{}, nonces: map[string]bool{}}
	old := time.Now().Add(-60 * 24 * time.Hour).UTC().Format(time.RFC3339)
	recent := time.Now().Add(-5 * 24 * time.Hour).UTC().Format(time.RFC3339)
	var records []map[string]any
	for _, id := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "13", "14", "15", "18", "19", "20", "21", "22", "23", "24", "25", "27", "28"} {
		records = append(records, map[string]any{"id": id, "author_id": "100", "created_at": old, "text": "fixture " + id, "public_metrics": map[string]int{"like_count": 1, "repost_count": 0, "reply_count": 0, "quote_count": 0, "bookmark_count": 0, "impression_count": 0}})
	}
	byID := map[string]map[string]any{}
	for _, p := range records {
		byID[p["id"].(string)] = p
	}
	byID["8"]["author_id"] = "900"
	byID["10"]["created_at"] = recent
	byID["11"]["created_at"] = recent
	byID["25"]["article"] = map[string]any{"id": "999"}
	for _, edge := range [][3]string{{"2", "1", "quoted"}, {"3", "1", "quoted"}, {"4", "2", "replied_to"}, {"5", "4", "replied_to"}, {"6", "4", "quoted"}, {"7", "8", "replied_to"}, {"9", "901", "retweeted"}, {"10", "902", "retweeted"}, {"14", "13", "quoted"}, {"15", "14", "replied_to"}, {"19", "18", "replied_to"}, {"20", "19", "replied_to"}, {"21", "18", "quoted"}, {"22", "19", "quoted"}, {"23", "1", "quoted"}, {"24", "23", "replied_to"}, {"27", "903", "retweeted"}, {"28", "904", "retweeted"}} {
		byID[edge[0]]["referenced_posts"] = []map[string]string{{"id": edge[1], "type": edge[2]}}
	}
	// Keep references and policy seeds on earlier versions of edited roots.
	byID["1"]["id"] = "101"
	byID["1"]["edit_history_post_ids"] = []string{"1", "101"}
	byID["18"]["id"] = "118"
	byID["18"]["edit_history_post_ids"] = []string{"18", "118"}
	// A reply can also quote another own post, creating converging paths.
	byID["5"]["referenced_posts"] = []map[string]string{{"id": "4", "type": "replied_to"}, {"id": "3", "type": "quoted"}}
	for _, record := range records {
		if record["author_id"] != "100" || record["article"] != nil {
			continue
		}
		path := "/2/tweets/" + record["id"].(string)
		refs, _ := record["referenced_posts"].([]map[string]string)
		for _, ref := range refs {
			if ref["type"] == "retweeted" {
				path = "/2/users/100/retweets/" + ref["id"]
			}
		}
		f.allowedMutations[path] = true
	}
	f.pages = [][]map[string]any{records[:12], records[12:]}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && f.mutationBarrier != nil {
			select {
			case f.mutationBarrier <- r.URL.Path:
			case <-r.Context().Done():
				return
			}
			select {
			case <-f.releaseMutations:
			case <-r.Context().Done():
				return
			}
			if f.holdPostUntilCancelled && strings.HasPrefix(r.URL.Path, "/2/tweets/") {
				<-r.Context().Done()
				return
			}
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		f.requests++
		signedAt, err := verifyAuthorization(r, f.nonces)
		if err != nil {
			t.Errorf("fixture authentication: %v", err)
			w.WriteHeader(401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/2/users/me":
			fmt.Fprintf(w, `{"data":{"id":%q}}`, f.account)
		case r.Method == http.MethodGet && r.URL.Path == "/2/users/100/tweets":
			f.reads++
			if f.relativeWindow && f.reads == 1 {
				f.windowStart, f.windowEnd = r.URL.Query().Get("start_time"), r.URL.Query().Get("end_time")
			}
			if r.URL.Query().Get("start_time") != f.windowStart || r.URL.Query().Get("end_time") != f.windowEnd {
				t.Error("timeline discovery bounds changed or were omitted")
			}
			if r.URL.Query().Get("exclude") != "" || r.URL.Query().Get("max_results") != "100" {
				t.Error("timeline excluded activities or did not request maximum page size")
			}
			if r.URL.Query().Get("post.fields") == "" {
				t.Error("timeline did not use the SDK's schema-defined post fields")
			}
			for _, field := range f.requiredFields {
				if !strings.Contains(","+r.URL.Query().Get("post.fields")+",", ","+field+",") {
					t.Errorf("timeline omitted required field %s", field)
				}
			}
			if f.quotaTimeline && f.reads == 1 {
				f.resetAt = time.Unix(time.Now().Unix()+2, 0)
				w.Header().Set("X-Rate-Limit-Remaining", "0")
				w.Header().Set("X-Rate-Limit-Reset", fmt.Sprint(f.resetAt.Unix()))
			} else if f.quotaTimeline && time.Now().Before(f.resetAt) {
				t.Error("timeline request dispatched before quota reset")
			} else if f.quotaTimeline && signedAt.Before(f.resetAt) {
				t.Error("timeline request was signed before waiting for quota")
			}
			page := 0
			if r.URL.Query().Get("pagination_token") != "" {
				page = 1
			}
			data := f.pages[page]
			if !f.ignoreWindow && (f.windowStart != "" || f.windowEnd != "") {
				filtered := []map[string]any{}
				for _, post := range data {
					stamp, ok := post["created_at"].(string)
					if !ok || (f.windowStart != "" && stamp < f.windowStart) || (f.windowEnd != "" && stamp >= f.windowEnd) {
						continue
					}
					filtered = append(filtered, post)
				}
				data = filtered
			}
			if f.malformed {
				data = []map[string]any{{"id": "1", "author_id": "100"}}
			}
			meta := map[string]any{}
			if page == 0 || f.repeatPage {
				meta["next_token"] = "next"
			}
			body := map[string]any{"data": data, "meta": meta}
			if page == 1 && len(f.timelineErrors) > 0 {
				body["errors"] = f.timelineErrors
				body["debug"] = "unrelated response field"
				if f.timelineStatus != 0 {
					w.WriteHeader(f.timelineStatus)
				}
			}
			if err := json.NewEncoder(w).Encode(body); err != nil {
				t.Errorf("encode fixture: %v", err)
			}
		case r.Method == http.MethodDelete:
			f.mutationAttempts = append(f.mutationAttempts, r.URL.Path)
			if !f.allowedMutations[r.URL.Path] {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if f.failMutation || (f.failMutationAt > 0 && len(f.mutationAttempts) == f.failMutationAt) {
				w.Header().Set("X-Rate-Limit-Reset", "2000000000")
				w.WriteHeader(429)
				fmt.Fprint(w, `{"detail":"fixture rate limit"}`)
				return
			}
			if f.mutationResponse != "" {
				fmt.Fprint(w, f.mutationResponse)
				return
			}
			if f.mutationStatus != 0 {
				w.WriteHeader(f.mutationStatus)
				fmt.Fprint(w, `{"detail":"fixture mutation status"}`)
				return
			}
			if f.quotaMutation && len(f.mutationAttempts) == 1 {
				f.firstMutationPath = r.URL.Path
				f.resetAt = time.Unix(time.Now().Unix()+2, 0)
				w.Header().Set("X-Rate-Limit-Remaining", "0")
				w.Header().Set("X-Rate-Limit-Reset", fmt.Sprint(f.resetAt.Unix()))
			} else if f.quotaMutation && time.Now().Before(f.resetAt) {
				if strings.Contains(f.firstMutationPath, "/retweets/") == strings.Contains(r.URL.Path, "/retweets/") {
					t.Error("mutation dispatched before its endpoint quota reset")
				} else {
					f.independentBeforeReset = true
				}
			} else if f.quotaMutation && signedAt.Before(f.resetAt) {
				t.Error("mutation was signed before waiting for quota")
			}
			f.mutations = append(f.mutations, r.URL.Path)
			if strings.Contains(r.URL.Path, "/retweets/") {
				fmt.Fprint(w, `{"data":{"retweeted":false}}`)
			} else {
				fmt.Fprint(w, `{"data":{"deleted":true}}`)
			}
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	t.Cleanup(server.Close)
	return f, server
}

// newWindowFixture places owned activities on both sides of chunk boundaries.
func newWindowFixture(t *testing.T) (*fixture, *httptest.Server) {
	t.Helper()
	f, server := newFixture(t)
	posts := map[string]map[string]any{}
	for _, page := range f.pages {
		for _, post := range page {
			posts[post["id"].(string)] = post
		}
	}
	for id, stamp := range map[string]string{"7": "2024-01-01T00:00:00Z", "9": "2024-01-15T00:00:00Z", "11": "2024-02-01T00:00:00Z", "27": "2023-12-31T23:59:59Z", "28": "2024-02-10T00:00:00Z"} {
		posts[id]["created_at"] = stamp
	}
	f.pages = [][]map[string]any{{posts["7"], posts["9"], posts["11"]}, {posts["27"], posts["28"]}}
	return f, server
}

// TestTimeWindow verifies API bounds, chunk boundaries, cache scope, and exact apply.
func TestTimeWindow(t *testing.T) {
	for _, scenario := range []struct {
		// name identifies a discovery-window contract.
		name string
		// declaration configures the window before data requests.
		declaration string
		// start is the expected API start, or empty for an open lower bound.
		start string
		// end is the expected API end, or empty for an open upper bound.
		end string
		// ignoreWindow makes the fixture return additional boundary records.
		ignoreWindow bool
		// paths contains every expected mutation in exact plan order.
		paths []string
	}{
		{name: "bounded", declaration: `set_time_window("2024-01-01T00:00:00Z","2024-02-01T00:00:00Z")`, start: "2024-01-01T00:00:00Z", end: "2024-02-01T00:00:00Z", paths: []string{"/2/users/100/retweets/901", "/2/tweets/7"}},
		{name: "start only", declaration: `set_time_window("2024-01-01T00:00:00Z")`, start: "2024-01-01T00:00:00Z", paths: []string{"/2/users/100/retweets/901", "/2/users/100/retweets/904", "/2/tweets/11", "/2/tweets/7"}},
		{name: "end only", declaration: `set_time_window(nil,"2024-02-01T00:00:00Z")`, end: "2024-02-01T00:00:00Z", paths: []string{"/2/users/100/retweets/901", "/2/users/100/retweets/903", "/2/tweets/7"}},
		{name: "unbounded", declaration: `set_time_window(nil,nil)`, paths: []string{"/2/users/100/retweets/901", "/2/users/100/retweets/903", "/2/users/100/retweets/904", "/2/tweets/11", "/2/tweets/7"}},
		{name: "previous adjacent chunk", declaration: `set_time_window("2023-12-01T00:00:00Z","2024-01-01T00:00:00Z")`, start: "2023-12-01T00:00:00Z", end: "2024-01-01T00:00:00Z", paths: []string{"/2/users/100/retweets/903"}},
		{name: "UTC normalization", declaration: `set_time_window("2024-01-01T01:00:00+01:00","2024-02-01T01:00:00+01:00")`, start: "2024-01-01T00:00:00Z", end: "2024-02-01T00:00:00Z", paths: []string{"/2/users/100/retweets/901", "/2/tweets/7"}},
		{name: "server returns extra records", declaration: `set_time_window("2024-01-01T00:00:00Z","2024-02-01T00:00:00Z")`, start: "2024-01-01T00:00:00Z", end: "2024-02-01T00:00:00Z", ignoreWindow: true, paths: []string{"/2/users/100/retweets/901", "/2/tweets/7"}},
		{name: "empty window", declaration: `set_time_window("2025-01-01T00:00:00Z","2025-02-01T00:00:00Z")`, start: "2025-01-01T00:00:00Z", end: "2025-02-01T00:00:00Z"},
		{name: "caught invalid change is atomic", declaration: `set_time_window("2024-01-01T00:00:00Z","2024-02-01T00:00:00Z"); assert(not pcall(function() set_time_window("2025-01-01T00:00:00Z","2024-01-01T00:00:00Z") end))`, start: "2024-01-01T00:00:00Z", end: "2024-02-01T00:00:00Z", paths: []string{"/2/users/100/retweets/901", "/2/tweets/7"}},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			f, server := newWindowFixture(t)
			f.windowStart, f.windowEnd, f.ignoreWindow = scenario.start, scenario.end, scenario.ignoreWindow
			policy := writePolicy(t, scenario.declaration+fmt.Sprintf(`
assert(#get_posts()+#get_reposts()==%d)
for _,p in ipairs(get_posts()) do assert(get_post(p.id).id==p.id) end
for _,p in ipairs(get_reposts()) do assert(get_repost(p.id).id==p.id) end
force_delete_post_when(function(p) return true end)
force_delete_repost_when(function(p) return true end)
`, len(scenario.paths)))
			plan := filepath.Join(t.TempDir(), "plan.json")
			out, err := runCLI(t, server, "", "plan", "--config", policy, "--out", plan)
			if err != nil || f.reads != 2 || len(f.mutationAttempts) != 0 {
				t.Fatalf("window planning: %v pages=%d attempts=%v\n%s", err, f.reads, f.mutationAttempts, out)
			}
			if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" && scenario.name == "bounded" {
				content, err := os.ReadFile(plan)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(output, "bounded-cleanup-plan.json"), content, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.Remove(policy); err != nil {
				t.Fatal(err)
			}
			out, err = runCLI(t, server, "", "apply", "--plan", plan)
			if err != nil || f.reads != 2 || !sameEndpointOrder(f.mutations, scenario.paths) {
				t.Fatalf("window apply: %v pages=%d mutations=%v; want %v\n%s", err, f.reads, f.mutations, scenario.paths, out)
			}
		})
	}
}

// TestRejectedTimeWindow keeps invalid or late declarations from publishing plans.
func TestRejectedTimeWindow(t *testing.T) {
	for _, scenario := range []struct {
		// name describes the invalid declaration.
		name string
		// policy supplies an invalid bound, late change, or missing preservation seed.
		policy string
		// start is the first window's expected API start, if any.
		start string
		// end is the first window's expected API end, if any.
		end string
		// reads counts timeline pages before the expected failure.
		reads int
	}{
		{name: "invalid timestamp", policy: `set_time_window("yesterday",nil)`},
		{name: "numeric timestamp", policy: `set_time_window(2024,nil)`},
		{name: "reversed bounds", policy: `set_time_window("2025-01-01T00:00:00Z","2024-01-01T00:00:00Z")`},
		{name: "equal bounds", policy: `set_time_window("2024-01-01T00:00:00Z","2024-01-01T00:00:00Z")`},
		{name: "unsupported precision", policy: `set_time_window("2024-01-01T00:00:00.5Z",nil)`},
		{name: "unsupported date", policy: `set_time_window("0001-01-01T00:00:00Z",nil)`},
		{name: "after query", policy: `get_posts(); set_time_window(nil,"2024-01-01T00:00:00Z")`, reads: 2},
		{name: "inside callback", policy: `delete_post_when(function(p) set_time_window(nil,"2024-01-01T00:00:00Z"); return true end)`, reads: 2},
		{name: "missing seed", policy: `set_time_window("2024-01-01T00:00:00Z","2024-02-01T00:00:00Z"); preserve_thread("11")`, start: "2024-01-01T00:00:00Z", end: "2024-02-01T00:00:00Z", reads: 2},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			f, server := newWindowFixture(t)
			f.windowStart, f.windowEnd = scenario.start, scenario.end
			path := filepath.Join(t.TempDir(), "plan.json")
			if err := os.WriteFile(path, []byte("existing plan"), 0o600); err != nil {
				t.Fatal(err)
			}
			out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, scenario.policy), "--out", path)
			if err == nil || f.reads != scenario.reads || len(f.mutationAttempts) != 0 {
				t.Fatalf("rejected window: %v pages=%d attempts=%v\n%s", err, f.reads, f.mutationAttempts, out)
			}
			content, err := os.ReadFile(path)
			if err != nil || string(content) != "existing plan" {
				t.Fatal("failed planning replaced the existing plan")
			}
		})
	}
}

// TestWindowGraphBoundary demonstrates the documented loss of an omitted connector.
func TestWindowGraphBoundary(t *testing.T) {
	for _, bounded := range []bool{false, true} {
		t.Run(fmt.Sprint(bounded), func(t *testing.T) {
			f, server := newFixture(t)
			stamps := map[string]string{"101": "2024-01-01T00:00:00Z", "2": "2024-01-02T00:00:00Z", "3": "2024-01-03T00:00:00Z", "4": "2024-01-04T00:00:00Z", "5": "2024-02-01T00:00:00Z", "6": "2024-01-05T00:00:00Z"}
			for _, page := range f.pages {
				for _, post := range page {
					stamp, ok := stamps[post["id"].(string)]
					if !ok {
						stamp = "2023-12-01T00:00:00Z"
					}
					post["created_at"] = stamp
				}
			}
			declaration := ""
			var paths []string
			if bounded {
				f.windowStart, f.windowEnd = "2024-01-01T00:00:00Z", "2024-02-01T00:00:00Z"
				declaration = `set_time_window("2024-01-01T00:00:00Z","2024-02-01T00:00:00Z"); assert(get_post("5")==nil); assert(#get_posts()==5);`
				paths = []string{"/2/tweets/2", "/2/tweets/4", "/2/tweets/6"}
			}
			out, err := runCLI(t, server, "", "apply", "--config", writePolicy(t, declaration+`delete_post({"2","4","6"}); preserve_thread("3")`), "--yes")
			if err != nil || f.reads != 2 || !sameEndpointOrder(f.mutations, paths) {
				t.Fatalf("window graph: %v pages=%d mutations=%v; want %v\n%s", err, f.reads, f.mutations, paths, out)
			}
		})
	}
}

// verifyAuthorization independently checks the actual HTTP method and query.
// It never calls the production signer or includes credential values in errors.
func verifyAuthorization(request *http.Request, nonces map[string]bool) (time.Time, error) {
	header, ok := strings.CutPrefix(request.Header.Get("Authorization"), "OAuth ")
	if !ok {
		return time.Time{}, fmt.Errorf("expected OAuth 1.0a user authentication")
	}
	parameters := request.URL.Query()
	oauth := url.Values{}
	for _, field := range strings.Split(header, ",") {
		name, raw, ok := strings.Cut(strings.TrimSpace(field), "=")
		if !ok || len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' || oauth.Has(name) {
			return time.Time{}, fmt.Errorf("malformed OAuth header")
		}
		value, err := url.QueryUnescape(raw[1 : len(raw)-1])
		if err != nil {
			return time.Time{}, fmt.Errorf("decode OAuth field: %w", err)
		}
		oauth.Set(name, value)
	}
	if oauth.Get("oauth_consumer_key") != fixtureCredentials["X_API_KEY"] || oauth.Get("oauth_token") != fixtureCredentials["X_ACCESS_TOKEN"] || oauth.Get("oauth_signature_method") != "HMAC-SHA1" || oauth.Get("oauth_version") != "1.0" {
		return time.Time{}, fmt.Errorf("unexpected OAuth user or signing method")
	}
	stamp, err := strconv.ParseInt(oauth.Get("oauth_timestamp"), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("decode OAuth timestamp: %w", err)
	}
	nonce := oauth.Get("oauth_nonce")
	if nonce == "" || nonces[nonce] {
		return time.Time{}, fmt.Errorf("missing or reused OAuth nonce")
	}
	supplied := oauth.Get("oauth_signature")
	oauth.Del("oauth_signature")
	for key, values := range oauth {
		for _, value := range values {
			parameters.Add(key, value)
		}
	}
	encode := func(value string) string { return strings.ReplaceAll(url.QueryEscape(value), "+", "%20") }
	var pairs []string
	for key, values := range parameters {
		for _, value := range values {
			pairs = append(pairs, encode(key)+"="+encode(value))
		}
	}
	sort.Strings(pairs)
	base := request.Method + "&" + encode("http://"+request.Host+request.URL.EscapedPath()) + "&" + encode(strings.Join(pairs, "&"))
	key := encode(fixtureCredentials["X_API_SECRET"]) + "&" + encode(fixtureCredentials["X_ACCESS_TOKEN_SECRET"])
	mac := hmac.New(sha1.New, []byte(key))
	if _, err := mac.Write([]byte(base)); err != nil {
		return time.Time{}, fmt.Errorf("compute fixture signature: %w", err)
	}
	if !hmac.Equal([]byte(supplied), []byte(base64.StdEncoding.EncodeToString(mac.Sum(nil)))) {
		return time.Time{}, fmt.Errorf("OAuth signature does not cover the request")
	}
	nonces[nonce] = true
	return time.Unix(stamp, 0), nil
}

// runCLI runs the real command with only synthetic OAuth credentials.
func runCLI(t *testing.T, server *httptest.Server, input string, args ...string) (string, error) {
	t.Helper()
	out, err := runCLIWithCredentials(t, server, input, fixtureCredentials, args...)
	if err != nil {
		return out, fmt.Errorf("execute cleaner: %w", err)
	}
	return out, nil
}

// runCLIWithCredentials isolates all X credential variables from the host.
func runCLIWithCredentials(t *testing.T, server *httptest.Server, input string, credentials map[string]string, args ...string) (string, error) {
	t.Helper()
	cmd := cliCommand(t, server, credentials, args...)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	assertPrivateDiagnostics(t, string(out))
	if err != nil {
		return string(out), fmt.Errorf("run cleaner command: %w", err)
	}
	return string(out), nil
}

// cliCommand prepares the real executable with isolated synthetic credentials.
func cliCommand(t *testing.T, server *httptest.Server, credentials map[string]string, args ...string) *exec.Cmd {
	t.Helper()
	binary, err := runfiles.Rlocation(*cleanerBinary)
	if err != nil {
		t.Fatalf("locate cleaner: %v", err)
	}
	args = append(args, "--api-base-url", server.URL+"/2")
	cmd := exec.Command(binary, args...)
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if _, ok := fixtureCredentials[name]; !ok && name != "X_USER_ACCESS_TOKEN" {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	for name, value := range credentials {
		cmd.Env = append(cmd.Env, name+"="+value)
	}
	cmd.Env = append(cmd.Env, "X_USER_ACCESS_TOKEN=fixture-obsolete-token")
	return cmd
}

// assertPrivateDiagnostics rejects disclosure of fixture credentials.
func assertPrivateDiagnostics(t *testing.T, output string) {
	t.Helper()
	for _, value := range fixtureCredentials {
		if strings.Contains(output, value) {
			t.Fatal("credential value appeared in command output")
		}
	}
	if strings.Contains(output, "fixture-obsolete-token") {
		t.Fatal("obsolete bearer token appeared in command output")
	}
}

// promptObserver signals when the CLI is waiting without supplying stdin data.
type promptObserver struct {
	// output records both streams, written by exec's shared output copier.
	output bytes.Buffer
	// ready closes when the confirmation prompt has been written.
	ready chan struct{}
	// once permits several writes containing the accumulated prompt.
	once sync.Once
}

// Write observes the prompt while retaining evidence until the process exits.
func (p *promptObserver) Write(content []byte) (int, error) {
	n, err := p.output.Write(content)
	if err != nil {
		return n, fmt.Errorf("record prompt output: %w", err)
	}
	if strings.Contains(p.output.String(), "Apply this plan?") {
		p.once.Do(func() { close(p.ready) })
	}
	return n, nil
}

// TestConfirmationCancellation interrupts open stdin without authorizing any action.
func TestConfirmationCancellation(t *testing.T) {
	f, server := newFixture(t)
	cmd := cliCommand(t, server, fixtureCredentials, "apply", "--config", writePolicy(t, `delete_post("7")`))
	input, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	observer := &promptObserver{ready: make(chan struct{})}
	cmd.Stdout, cmd.Stderr = observer, observer
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-observer.ready:
	case err := <-done:
		t.Fatalf("exited before confirmation: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("confirmation prompt was not reached")
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("interrupted confirmation reported success")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("interrupt left confirmation blocked on open stdin")
	}
	if len(f.mutationAttempts) != 0 {
		t.Fatalf("interrupted confirmation mutated X: %v", f.mutationAttempts)
	}
	assertPrivateDiagnostics(t, observer.output.String())
}

// TestOperatorProgress checks separate JSON and logs, quota visibility, and final counts.
func TestOperatorProgress(t *testing.T) {
	f, server := newFixture(t)
	policy := writePolicy(t, `delete_post({"7","11"}); delete_repost({"9","10"})`)
	cmd := cliCommand(t, server, fixtureCredentials, "plan", "--config", policy)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("plan: %v\n%s", err, stderr.String())
	}
	if !json.Valid(stdout.Bytes()) || len(f.mutationAttempts) != 0 {
		t.Fatal("planning mixed diagnostics into JSON or mutated X")
	}
	for _, message := range []string{"Identifying account", "Evaluating cleanup policy", "Fetching timeline", "Fetched timeline page", "Timeline collected", "Plan prepared"} {
		if !strings.Contains(stderr.String(), message) {
			t.Errorf("missing planning progress %q", message)
		}
	}
	if !strings.Contains(stderr.String(), "page=2") || !strings.Contains(stderr.String(), "actions=4") {
		t.Error("missing page or action counts")
	}
	plan := filepath.Join(t.TempDir(), "plan.json")
	if err := os.WriteFile(plan, stdout.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	planningLog := stderr.String()
	stdout.Reset()
	stderr.Reset()
	f.quotaMutation = true
	cmd = cliCommand(t, server, fixtureCredentials, "apply", "--plan", plan)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("apply: %v\n%s", err, stderr.String())
	}
	for _, message := range []string{"Loading saved plan", "Applying cleanup plan", "Applying action", "Action completed", "Waiting for API quota", "API quota available", "Cleanup completed"} {
		if !strings.Contains(stderr.String(), message) {
			t.Errorf("missing execution progress %q", message)
		}
		if strings.Contains(stdout.String(), message) {
			t.Errorf("progress %q appeared on stdout", message)
		}
	}
	if !strings.Contains(stderr.String(), "completed=4 total=4") || len(f.mutations) != 4 {
		t.Fatalf("incorrect completed progress: mutations=%v\n%s", f.mutations, stderr.String())
	}
	logs := planningLog + stderr.String()
	assertPrivateDiagnostics(t, logs)
	if strings.Contains(logs, "fixture 7") || strings.Contains(logs, "fixture 11") {
		t.Fatal("progress disclosed post contents")
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "operator-progress.log"), []byte(logs), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// TestAuthentication rejects incomplete credentials without requests or disclosure.
func TestAuthentication(t *testing.T) {
	for missing := range fixtureCredentials {
		t.Run(missing, func(t *testing.T) {
			f, server := newFixture(t)
			credentials := map[string]string{}
			for name, value := range fixtureCredentials {
				if name != missing {
					credentials[name] = value
				}
			}
			out, err := runCLIWithCredentials(t, server, "", credentials, "apply", "--config", writePolicy(t, `delete_post("7")`), "--yes")
			if err == nil || !strings.Contains(out, missing) || f.requests != 0 {
				t.Fatalf("incomplete credentials: err=%v requests=%d", err, f.requests)
			}
		})
	}
}

// writePolicy writes caller-owned policy input into the test's private directory.
func writePolicy(t *testing.T, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "policy.lua")
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatalf("write policy: %v", err)
	}
	return path
}

// TestPlanAndApply verifies graph semantics, precedence, pagination, and action routing.
func TestPlanAndApply(t *testing.T) {
	f, server := newFixture(t)
	policy := writePolicy(t, `
assert(os == nil and io == nil and debug == nil and package == nil and require == nil)
assert(get_post("8") == nil)
assert(get_post("7").is_reply)
assert(not get_post("101").is_reply and not get_post("101").is_quote)
assert(get_post("5").is_reply and get_post("5").is_quote)
assert(not get_repost("9").is_reply and not get_repost("9").is_quote)
assert(get_post("1").id == "101" and get_post("101").id == "101")
assert(#get_quotes("1") == 3)
assert(#get_quotes("101") == 3)
assert(#get_replies("2") == 1)
assert(get_replies("18")[1].id == "19")
local copy = get_post("7"); copy.id = "999"
assert(get_post("7").id == "7")
local function old(item) return item.age_days > 30 end
delete_post_when(old)
delete_repost_when(old)
local nodes, edges = {}, {}
preserve_thread({"1", "14"}, {
  include_post = function(p)
    assert(not nodes[p.id]); nodes[p.id] = true
    return p.id ~= "4"
  end,
  follow_edge = function(e)
    local key = e.from_id .. ":" .. e.to_id .. ":" .. e.kind
    assert(not edges[key]); edges[key] = true
    return e.to_id ~= "23"
  end,
})
preserve_post({"18"}, function(p) return p.id ~= "19" end)
force_delete_post({"18"})
delete_post("18")
preserve_repost({"27"})
preserve_repost_when(function(p) return p.id == "28" end)
force_delete_post_when(function(p) return p.id == "2" end)
force_delete_repost_when(function(p) return p.id == "28" end)
`)
	plan := filepath.Join(t.TempDir(), "plan.json")
	if out, err := runCLI(t, server, "", "plan", "--config", policy, "--out", plan); err != nil {
		t.Fatalf("plan: %v\n%s", err, out)
	}
	if len(f.mutationAttempts) != 0 || f.reads != 2 {
		t.Fatalf("planning effects: mutation attempts=%v pages=%d", f.mutationAttempts, f.reads)
	}
	content, err := os.ReadFile(plan)
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Actions []struct {
			Kind     string `json:"kind"`
			TargetID string `json:"target_id"`
		} `json:"actions"`
	}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatal(err)
	}
	var actual []string
	for _, action := range result.Actions {
		actual = append(actual, action.Kind+":"+action.TargetID)
	}
	sort.Strings(actual)
	want := []string{"delete_post:118", "delete_post:13", "delete_post:19", "delete_post:2", "delete_post:21", "delete_post:22", "delete_post:23", "delete_post:24", "delete_post:4", "delete_post:7", "unrepost:901", "unrepost:904"}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("actions: got %v; want %v\n%s", actual, want, content)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "cleanup-plan.json"), content, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Remove(policy); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, server, "", "apply", "--plan", plan); err != nil {
		t.Fatalf("apply: %v\n%s", err, out)
	}
	if f.reads != 2 || len(f.mutations) != len(want) {
		t.Fatalf("saved apply reread inventory or lost actions: reads=%d mutations=%v", f.reads, f.mutations)
	}
	wantPaths := []string{"/2/users/100/retweets/901", "/2/users/100/retweets/904", "/2/tweets/118", "/2/tweets/13", "/2/tweets/19", "/2/tweets/2", "/2/tweets/21", "/2/tweets/22", "/2/tweets/23", "/2/tweets/24", "/2/tweets/4", "/2/tweets/7"}
	if !sameEndpointOrder(f.mutations, wantPaths) {
		t.Fatalf("mutation paths: got %v; want %v", f.mutations, wantPaths)
	}
}

// TestDirectSelections verifies scalar/array equivalence and monotonic precedence.
func TestDirectSelections(t *testing.T) {
	for _, scenario := range []struct {
		// name identifies the acceptance case.
		name string
		// policy registers direct selections through the public Lua functions.
		policy string
		// paths lists the exact expected mutations in execution order.
		paths []string
	}{
		{name: "empty arrays", policy: `delete_post({}); force_delete_post({}); delete_repost({}); force_delete_repost({}); preserve_post({}); preserve_thread({}); preserve_repost({})`},
		{name: "quote in a reply preserves its own reply ancestors", policy: `delete_post({"101", "2", "3", "4", "5", "6"}); preserve_thread("3")`, paths: []string{"/2/tweets/101"}},
		{name: "reply parent traversal can be filtered", policy: `delete_post({"101", "2", "3", "4", "5", "6"}); preserve_thread("3", {follow_edge=function(e) return e.to_id ~= "2" end})`, paths: []string{"/2/tweets/101", "/2/tweets/2"}},
		{name: "post preservation keeps reply descendants only", policy: `delete_post({"101", "2", "3", "4", "5", "6"}); preserve_post("5")`, paths: []string{"/2/tweets/101", "/2/tweets/2", "/2/tweets/3", "/2/tweets/4", "/2/tweets/6"}},
		{name: "forced before ordinary", policy: `force_delete_post({"7"}); delete_post("7"); preserve_post("7"); force_delete_repost({"9"}); delete_repost("9"); preserve_repost("9")`, paths: []string{"/2/users/100/retweets/901", "/2/tweets/7"}},
		{name: "ordinary before forced", policy: `preserve_post({"7"}); delete_post({"7"}); force_delete_post("7"); preserve_repost("9"); delete_repost({"9"}); force_delete_repost("9")`, paths: []string{"/2/users/100/retweets/901", "/2/tweets/7"}},
		{name: "preserved arrays", policy: `delete_post({"7"}); preserve_post({"7"}); delete_repost({"9"}); preserve_repost({"9"})`},
		{name: "invalid array is atomic", policy: `assert(not pcall(function() force_delete_post({"7", "999999"}) end))`},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			f, server := newFixture(t)
			out, err := runCLI(t, server, "", "apply", "--config", writePolicy(t, scenario.policy), "-y")
			if err != nil {
				t.Fatalf("direct selections: %v\n%s", err, out)
			}
			if !sameEndpointOrder(f.mutations, scenario.paths) {
				t.Fatalf("mutation paths: got %v; want %v", f.mutations, scenario.paths)
			}
		})
	}
}

// TestSavedPlanBoundary verifies complete validation and exact, bounded execution.
func TestSavedPlanBoundary(t *testing.T) {
	for _, scenario := range []struct {
		// name identifies the application boundary being checked.
		name string
		// actions is the explicit saved action list, independent of a policy.
		actions string
		// failAt rejects this numbered mutation attempt, or zero for success.
		failAt int
		// attempts lists every mutation request expected from the command.
		attempts []string
		// completed is the number of confirmed actions reported to the user.
		completed int
		// wantError distinguishes rejected and partially applied plans.
		wantError bool
	}{
		{name: "exact subset and endpoint order", actions: `[{"kind":"delete_post","target_id":"11","forced":true},{"kind":"unrepost","target_id":"902","forced":false}]`, attempts: []string{"/2/tweets/11", "/2/users/100/retweets/902"}, completed: 2},
		{name: "invalid later kind", actions: `[{"kind":"delete_post","target_id":"7","forced":false},{"kind":"unknown","target_id":"2","forced":false}]`, wantError: true},
		{name: "invalid later target", actions: `[{"kind":"delete_post","target_id":"7","forced":false},{"kind":"delete_post","target_id":"2/other","forced":false}]`, wantError: true},
		{name: "duplicate action", actions: `[{"kind":"delete_post","target_id":"7","forced":false},{"kind":"delete_post","target_id":"7","forced":true}]`, wantError: true},
		{name: "stop after failure", actions: `[{"kind":"delete_post","target_id":"7","forced":false},{"kind":"delete_post","target_id":"11","forced":false},{"kind":"delete_post","target_id":"2","forced":false}]`, failAt: 2, attempts: []string{"/2/tweets/7", "/2/tweets/11"}, completed: 1, wantError: true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			f, server := newFixture(t)
			f.failMutationAt = scenario.failAt
			plan := filepath.Join(t.TempDir(), "plan.json")
			content := `{"version":1,"account_id":"100","evaluated_at":"2026-10-30T00:00:00Z","actions":` + scenario.actions + `}`
			if err := os.WriteFile(plan, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			out, err := runCLI(t, server, "", "apply", "--plan", plan)
			if (err != nil) != scenario.wantError {
				t.Fatalf("apply result: %v\n%s", err, out)
			}
			if f.reads != 0 || !sameEndpointOrder(f.mutationAttempts, scenario.attempts) {
				t.Fatalf("saved plan boundary: timeline reads=%d mutation attempts=%v; want %v", f.reads, f.mutationAttempts, scenario.attempts)
			}
			if len(f.mutations) != scenario.completed {
				t.Fatalf("completed mutations: got %v; want %d", f.mutations, scenario.completed)
			}
			if len(scenario.attempts) > 0 && !strings.Contains(out, fmt.Sprintf("Completed %d of", scenario.completed)) {
				t.Fatalf("missing confirmed progress: %s", out)
			}
		})
	}
}

// TestTimelineProblems verifies useful bounded diagnostics without partial plans or mutations.
func TestTimelineProblems(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusForbidden} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			f, server := newFixture(t)
			detail := "The requested resource is unavailable."
			for _, credential := range fixtureCredentials {
				detail += " " + credential + " " + url.QueryEscape(credential)
			}
			problem := map[string]any{"type": "https://api.x.com/2/problems/resource-not-found", "title": "Not Found Error", "detail": detail, "resource_type": "tweet", "resource_id": "999", "debug": "unrelated problem field"}
			f.timelineErrors = []map[string]any{problem, problem, problem, problem}
			f.timelineStatus = status
			plan := filepath.Join(t.TempDir(), "plan.json")
			out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, `delete_post_when(function(post) return true end)`), "--out", plan)
			if err == nil || f.reads != 2 || len(f.mutationAttempts) != 0 {
				t.Fatalf("partial planning: err=%v reads=%d mutations=%v", err, f.reads, f.mutationAttempts)
			}
			if _, err := os.Stat(plan); !os.IsNotExist(err) {
				t.Fatal("failed collection must not publish a plan")
			}
			for _, expected := range []string{fmt.Sprintf("HTTP %d", status), "page 2", "after 12 activities", "resource-not-found", "Not Found Error", "resource_id=\"999\"", "1 additional problem", "[redacted]"} {
				if !strings.Contains(out, expected) {
					t.Errorf("missing diagnostic %q", expected)
				}
			}
			for _, credential := range fixtureCredentials {
				if strings.Contains(out, credential) || strings.Contains(out, url.QueryEscape(credential)) {
					t.Fatal("API diagnostics disclosed credentials")
				}
			}
			if strings.Contains(out, "unrelated") || strings.Contains(out, "fixture 1") {
				t.Fatal("API diagnostics included unrelated response content")
			}
			if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" && status == http.StatusOK {
				if err := os.WriteFile(filepath.Join(output, "partial-response-diagnostic.log"), []byte(out), 0o600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// TestLazyInventory verifies policies request only the data they actually use.
func TestLazyInventory(t *testing.T) {
	for _, scenario := range []struct {
		// name identifies the lazy collection contract.
		name string
		// policy selects or queries only the intended inventory views.
		policy string
		// reads is the expected number of shared timeline pages.
		reads int
		// malformedPost removes a post timestamp to exercise selective validation.
		malformedPost bool
		// wantError requires the complete plan to fail despite caught Lua errors.
		wantError bool
	}{
		{name: "empty", policy: ``},
		{name: "empty arrays", policy: `preserve_post({}); preserve_thread({}); delete_repost({})`},
		{name: "unused time window", policy: `set_time_window(nil,"2024-02-01T00:00:00Z"); preserve_thread({}); delete_repost({})`},
		{name: "repost only", policy: `delete_repost_when(function(p) return p.age_days > 30 end)`, reads: 2, malformedPost: true},
		{name: "cross kind callback", policy: `local n=0; delete_repost_when(function(p) n=n+1; assert(#get_posts()==18); return false end); preserve_repost_when(function(p) assert(n==4); return false end)`, reads: 2},
		{name: "caught collection failure", policy: `assert(not pcall(get_posts))`, reads: 2, malformedPost: true, wantError: true},
		{name: "caught callback collection failure", policy: `delete_repost_when(function(p) pcall(get_posts); return true end)`, reads: 2, malformedPost: true, wantError: true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			f, server := newFixture(t)
			if scenario.malformedPost {
				delete(f.pages[0][0], "created_at")
			}
			out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, scenario.policy))
			if (err != nil) != scenario.wantError {
				t.Fatalf("lazy plan: %v\n%s", err, out)
			}
			if f.reads != scenario.reads || len(f.mutationAttempts) != 0 {
				t.Fatalf("lazy effects: pages=%d mutation attempts=%v", f.reads, f.mutationAttempts)
			}
		})
	}
}

// TestLuaAPIArguments rejects invalid calls before reading inventory or selecting items.
func TestLuaAPIArguments(t *testing.T) {
	for _, source := range []string{
		`delete_post("7", function() return false end)`,
		`delete_post_when(function() return true end, function() return false end)`,
		`force_delete_post("7", {}, false)`,
		`preserve_post("7", {}, false)`,
		`get_post("7", false)`, `get_reposts(false)`, `get_replies("7", false)`,
		`delete_post("7", {include_post=function() return false end})`,
		`delete_post("7", {name=""})`, `delete_post("7", {name=7})`,
		`preserve_repost("9", {follow_edge=function() return true end})`,
		`days_ago(-1)`, `days_ago("30")`, `days_ago(0/0)`,
		`days_ago(1/0)`, `days_ago(1000000)`, `days_ago(1,2)`,
	} {
		t.Run(source, func(t *testing.T) {
			f, server := newFixture(t)
			out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, source))
			if err == nil || f.reads != 0 || len(f.mutationAttempts) != 0 {
				t.Fatalf("invalid API call: error=%v reads=%d mutations=%v output=%s", err, f.reads, f.mutationAttempts, out)
			}
		})
	}
}

// TestLuaAPIDecisions verifies rule provenance, exact actions, and metadata isolation.
func TestLuaAPIDecisions(t *testing.T) {
	f, server := newFixture(t)
	path := filepath.Join(t.TempDir(), "explained-plan.json")
	source := `
local calls = 0
delete_post_when(function(p) calls=calls+1; return p.id=="7" or p.id=="118" or p.id=="19" end, {name="old posts"})
assert(calls==0)
preserve_post_when(function(p) assert(calls==18); return p.id=="118" end, {name="important discussion"})
force_delete_post("19", {name="explicit removal"})
delete_repost("9", {name="old repost"})
`
	out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, source), "--out", path, "--explain")
	if err != nil || len(f.mutationAttempts) != 0 {
		t.Fatalf("explain plan: %v\n%s", err, out)
	}
	for _, want := range []string{"important discussion", "explicit removal", "fixture 18", "preserved"} {
		if !strings.Contains(out, want) {
			t.Errorf("explanation missing %q: %s", want, out)
		}
	}
	document := readPlanDocument(t, path)
	decisions := document["decisions"].([]any)
	byID := map[string]map[string]any{}
	for _, value := range decisions {
		d := value.(map[string]any)
		byID[d["item_id"].(string)] = d
	}
	if len(decisions) != 22 || byID["118"]["outcome"] != "preserved" || byID["19"]["outcome"] != "force_delete" || byID["101"]["outcome"] != "unchanged" {
		t.Fatalf("unexpected decisions: %v", decisions)
	}
	reasons := byID["19"]["reasons"].([]any)
	if len(reasons) != 3 {
		t.Fatalf("missing matching rules: %v", reasons)
	}
	found := false
	for _, value := range reasons {
		r := value.(map[string]any)
		if r["name"] == "important discussion" && r["seed_id"] == "118" && r["rule_id"].(float64) > 0 {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing preservation source: %v", reasons)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "explained-cleanup-plan.json"), content, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	// Explanations are descriptive and must never become execution input.
	document["decisions"] = append(decisions, map[string]any{"item_id": "13", "outcome": "force_delete"})
	content, err = json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	if out, err := runCLI(t, server, "", "apply", "--plan", path); err != nil {
		t.Fatalf("apply explained plan: %v\n%s", err, out)
	}
	if !sameEndpointOrder(f.mutations, []string{"/2/users/100/retweets/901", "/2/tweets/19", "/2/tweets/7"}) || f.reads != 2 {
		t.Fatalf("explanations changed execution: %v, reads %d", f.mutations, f.reads)
	}
}

// readPlanDocument reads the public plan schema without importing implementation types.
func readPlanDocument(t *testing.T, path string) map[string]any {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(content, &document); err != nil {
		t.Fatal(err)
	}
	return document
}

// TestLuaAPITime checks rolling bounds against the plan's fixed evaluation time.
func TestLuaAPITime(t *testing.T) {
	f, server := newFixture(t)
	f.relativeWindow = true
	path := filepath.Join(t.TempDir(), "relative-plan.json")
	source := `local start=days_ago(90); set_time_window(start,days_ago(30)); delete_post_when(function(p) assert(days_ago(90)==start); return true end)`
	if out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, source), "--out", path); err != nil {
		t.Fatalf("relative window: %v\n%s", err, out)
	}
	document := readPlanDocument(t, path)
	now, err := time.Parse(time.RFC3339Nano, document["evaluated_at"].(string))
	if err != nil {
		t.Fatal(err)
	}
	if f.windowStart != now.Add(-90*24*time.Hour).UTC().Truncate(time.Second).Format(time.RFC3339) || f.windowEnd != now.Add(-30*24*time.Hour).UTC().Truncate(time.Second).Format(time.RFC3339) || f.reads != 2 {
		t.Fatalf("inconsistent rolling bounds: %s %s", f.windowStart, f.windowEnd)
	}
}

// TestLuaAPIMetadata verifies optional structured data and isolated nested Lua copies.
func TestLuaAPIMetadata(t *testing.T) {
	f, server := newFixture(t)
	f.requiredFields = []string{"lang", "attachments", "entities"}
	p := f.pages[0][0]
	p["lang"] = "en"
	p["attachments"] = map[string]any{"media_keys": []string{"3_123"}, "poll_ids": []string{"456"}}
	p["entities"] = map[string]any{"hashtags": []map[string]any{{"start": 0, "end": 6, "tag": "short"}}}
	p["note_post"] = map[string]any{"text": "long #retained", "entities": map[string]any{"hashtags": []map[string]any{{"start": 5, "end": 14, "tag": "retained"}}, "mentions": []map[string]any{{"start": 0, "end": 4, "id": "900", "username": "name"}}, "urls": []map[string]any{{"start": 0, "end": 4, "url": "https://t.co/test", "expanded_url": "https://example.org/article"}}}}
	source := `
local p=get_post("101")
assert(p.language=="en" and p.attachments.media_keys[1]=="3_123" and p.attachments.poll_ids[1]=="456")
assert(p.entities.hashtags[1].tag=="retained" and p.entities.mentions[1].id=="900")
assert(p.entities.urls[1].expanded_url=="https://example.org/article")
p.attachments.media_keys[1]="changed"; p.entities.hashtags[1].tag="changed"
assert(get_post("101").attachments.media_keys[1]=="3_123" and get_post("101").entities.hashtags[1].tag=="retained")
assert(get_post("7").language==nil and get_post("7").entities==nil and get_post("7").attachments==nil)
delete_post("101")
preserve_post_when(function(p) return p.language=="en" and p.entities.hashtags[1].tag=="retained" end)
`
	path := filepath.Join(t.TempDir(), "metadata-plan.json")
	if out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, source), "--out", path); err != nil {
		t.Fatalf("metadata policy: %v\n%s", err, out)
	}
	if len(readPlanDocument(t, path)["actions"].([]any)) != 0 || f.reads != 2 || len(f.mutationAttempts) != 0 {
		t.Fatal("metadata policy selected unexpected actions or performed extra requests")
	}
}

// TestLuaAPITraversal verifies seed-specific callbacks and explicit traversal direction.
func TestLuaAPITraversal(t *testing.T) {
	f, server := newFixture(t)
	source := `
local visits={}
preserve_thread_when(function(p) return p.id=="2" or p.id=="3" end, {
 include_post=function(p,ctx)
   assert(ctx.seed_id=="2" or ctx.seed_id=="3")
   visits[ctx.seed_id..":"..p.id]=true
   ctx.seed_id="changed"
   return true
 end,
 follow_edge=function(e,ctx)
   assert(ctx.seed_id=="2" or ctx.seed_id=="3")
   assert(e.direction=="child" or e.direction=="parent")
   return true
 end,
})
preserve_repost_when(function(p) return p.id=="9" end, function(p,ctx) assert(ctx.seed_id=="9" and visits["2:4"] and visits["3:4"]); return false end)
delete_post({"2","3","4","5","6"})
`
	if out, err := runCLI(t, server, "", "apply", "--config", writePolicy(t, source), "-y"); err != nil || len(f.mutationAttempts) != 0 {
		t.Fatalf("seed callbacks: %v\n%s", err, out)
	}
}

// TestLuaAPIOverlappingThreads exercises callback-free union traversal at timeline scale.
func TestLuaAPIOverlappingThreads(t *testing.T) {
	f, server := newFixture(t)
	var records []map[string]any
	for i := 1000; i < 4200; i++ {
		p := map[string]any{"id": strconv.Itoa(i), "author_id": "100", "created_at": "2024-01-01T00:00:00Z", "text": "thread"}
		if i > 1000 {
			p["referenced_posts"] = []map[string]string{{"type": "replied_to", "id": strconv.Itoa(i - 1)}}
		}
		records = append(records, p)
	}
	f.pages = [][]map[string]any{records[:1600], records[1600:]}
	path := filepath.Join(t.TempDir(), "overlapping-plan.json")
	if out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, `delete_post_when(function() return true end); preserve_thread_when(function() return true end)`), "--out", path, "--policy-timeout", "5s"); err != nil {
		t.Fatalf("overlapping threads: %v\n%s", err, out)
	}
	if len(readPlanDocument(t, path)["actions"].([]any)) != 0 || len(f.mutationAttempts) != 0 {
		t.Fatal("overlapping traversal changed preservation")
	}
}

// TestParallelDeletion requires concurrent requests from the independent endpoint buckets.
func TestParallelDeletion(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprintf("failure=%v", fail), func(t *testing.T) {
			f, server := newFixture(t)
			f.holdPostUntilCancelled = fail
			f.failMutation = fail
			f.mutationBarrier = make(chan string, 4)
			f.releaseMutations = make(chan struct{})
			plan := filepath.Join(t.TempDir(), "parallel-plan.json")
			content := `{"version":1,"account_id":"100","evaluated_at":"2026-10-30T00:00:00Z","actions":[{"kind":"unrepost","target_id":"901","forced":false},{"kind":"delete_post","target_id":"7","forced":false},{"kind":"unrepost","target_id":"902","forced":false},{"kind":"delete_post","target_id":"11","forced":false}]}`
			if err := os.WriteFile(plan, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			done := make(chan struct{})
			go func() {
				defer close(done)
				if out, err := runCLI(t, server, "", "apply", "--plan", plan); (err != nil) != fail {
					t.Errorf("parallel apply: %v\n%s", err, out)
				}
			}()
			paths := []string{}
			timer := time.NewTimer(3 * time.Second)
			defer timer.Stop()
		waiting:
			for len(paths) < 2 {
				select {
				case path := <-f.mutationBarrier:
					paths = append(paths, path)
				case <-timer.C:
					t.Error("independent deletion endpoints did not run concurrently")
					break waiting
				}
			}
			close(f.releaseMutations)
			<-done
			sort.Strings(paths)
			if !reflect.DeepEqual(paths, []string{"/2/tweets/7", "/2/users/100/retweets/901"}) || f.reads != 0 {
				t.Fatalf("parallel actions: paths=%v mutations=%v reads=%d", paths, f.mutations, f.reads)
			}
			if fail {
				if len(f.mutationAttempts) != 1 || len(f.mutations) != 0 {
					t.Fatalf("failure dispatched further actions: %v", f.mutationAttempts)
				}
			} else if len(f.mutations) != 4 {
				t.Fatalf("lost parallel actions: %v", f.mutations)
			}
		})
	}
}

// TestDeletionDependencies keeps same-target operations in their explicit plan order.
func TestDeletionDependencies(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		f, server := newFixture(t)
		f.allowedMutations["/2/users/100/retweets/7"] = true
		actions := []string{`{"kind":"unrepost","target_id":"7","forced":false}`, `{"kind":"delete_post","target_id":"7","forced":false}`}
		paths := []string{"/2/users/100/retweets/7", "/2/tweets/7"}
		if reverse {
			actions[0], actions[1] = actions[1], actions[0]
			paths[0], paths[1] = paths[1], paths[0]
		}
		plan := filepath.Join(t.TempDir(), "dependent-plan.json")
		content := `{"version":1,"account_id":"100","evaluated_at":"2026-10-30T00:00:00Z","actions":[` + strings.Join(actions, ",") + `]}`
		if err := os.WriteFile(plan, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		if out, err := runCLI(t, server, "", "apply", "--plan", plan); err != nil {
			t.Fatalf("dependent apply: %v\n%s", err, out)
		}
		if !reflect.DeepEqual(f.mutations, paths) {
			t.Fatalf("same-target order: got %v want %v", f.mutations, paths)
		}
	}
}

// sameEndpointOrder compares exact requests while allowing independent endpoints to interleave.
func sameEndpointOrder(actual, expected []string) bool {
	group := func(paths []string) map[bool][]string {
		groups := map[bool][]string{}
		for _, path := range paths {
			key := strings.Contains(path, "/retweets/")
			groups[key] = append(groups[key], path)
		}
		return groups
	}
	return reflect.DeepEqual(group(actual), group(expected))
}

// TestQuotaPacing verifies waiting before dispatch and independent endpoint limits.
func TestQuotaPacing(t *testing.T) {
	for _, independent := range []bool{false, true} {
		t.Run(fmt.Sprint(independent), func(t *testing.T) {
			f, server := newFixture(t)
			f.quotaMutation = true
			actions := `[{"kind":"delete_post","target_id":"7","forced":false},{"kind":"delete_post","target_id":"11","forced":false}]`
			if independent {
				actions = `[{"kind":"delete_post","target_id":"7","forced":false},{"kind":"unrepost","target_id":"901","forced":false}]`
			}
			plan := filepath.Join(t.TempDir(), "plan.json")
			content := `{"version":1,"account_id":"100","evaluated_at":"2026-10-30T00:00:00Z","actions":` + actions + `}`
			if err := os.WriteFile(plan, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			out, err := runCLI(t, server, "", "apply", "--plan", plan)
			if err != nil || len(f.mutations) != 2 {
				t.Fatalf("quota apply: %v mutations=%v\n%s", err, f.mutations, out)
			}
			if independent && !f.independentBeforeReset {
				t.Fatal("another endpoint was unnecessarily delayed")
			}
		})
	}
	t.Run("cancel timeline wait", func(t *testing.T) {
		f, server := newFixture(t)
		f.quotaTimeline = true
		out, err := runCLI(t, server, "", "plan", "--config", writePolicy(t, `delete_post_when(function(p) return true end)`), "--policy-timeout", "150ms")
		if err == nil || f.reads != 1 || len(f.mutationAttempts) != 0 {
			t.Fatalf("wait cancellation: %v pages=%d attempts=%v\n%s", err, f.reads, f.mutationAttempts, out)
		}
	})
}

// TestIdempotentApply skips absent targets on explicit replay without masking other failures.
func TestIdempotentApply(t *testing.T) {
	for _, status := range []int{404, 403, 429, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			f, server := newFixture(t)
			plan := filepath.Join(t.TempDir(), "plan.json")
			content := `{"version":1,"account_id":"100","evaluated_at":"2026-10-30T00:00:00Z","actions":[{"kind":"delete_post","target_id":"7","forced":false},{"kind":"unrepost","target_id":"901","forced":false}]}`
			if err := os.WriteFile(plan, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			if out, err := runCLI(t, server, "", "apply", "--plan", plan); err != nil {
				t.Fatalf("first apply: %v\n%s", err, out)
			}
			f.mutationStatus = status
			out, err := runCLI(t, server, "", "apply", "--plan", plan)
			if (err == nil) != (status == 404) || len(f.mutations) != 2 {
				t.Fatalf("repeated apply status=%d: %v mutations=%v\n%s", status, err, f.mutations, out)
			}
			if status == 404 && (!strings.Contains(out, "Completed 2 of 2") || !strings.Contains(out, "already absent")) {
				t.Fatalf("missing skip result: %s", out)
			}
		})
	}
}

// TestUnconfirmedMutation rejects missing or null confirmations before later actions.
func TestUnconfirmedMutation(t *testing.T) {
	for _, response := range []string{`{"data":{}}`, `{"data":{"retweeted":null}}`} {
		t.Run(response, func(t *testing.T) {
			f, server := newFixture(t)
			f.mutationResponse = response
			plan := filepath.Join(t.TempDir(), "plan.json")
			content := `{"version":1,"account_id":"100","evaluated_at":"2026-10-30T00:00:00Z","actions":[{"kind":"unrepost","target_id":"901","forced":false},{"kind":"unrepost","target_id":"902","forced":false}]}`
			if err := os.WriteFile(plan, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			out, err := runCLI(t, server, "", "apply", "--plan", plan)
			if err == nil || len(f.mutationAttempts) != 1 {
				t.Fatalf("unconfirmed result: %v attempts=%v\n%s", err, f.mutationAttempts, out)
			}
		})
	}
}

// TestInlineApply verifies confirmation, both yes spellings, and no-op defaults.
func TestInlineApply(t *testing.T) {
	for _, answer := range []string{"n\n", "", "yes\n", "-y", "--yes"} {
		t.Run(strings.TrimSpace(answer), func(t *testing.T) {
			f, server := newFixture(t)
			policy := writePolicy(t, `delete_post_when(function(p) return p.id == "7" end)`)
			args := []string{"apply", "--config", policy}
			input := answer
			if strings.HasPrefix(answer, "-") {
				args = append(args, answer)
				input = ""
			}
			out, err := runCLI(t, server, input, args...)
			if err != nil {
				t.Fatalf("inline apply: %v\n%s", err, out)
			}
			want := 0
			if answer == "yes\n" || strings.HasPrefix(answer, "-") {
				want = 1
			}
			if len(f.mutations) != want {
				t.Fatalf("confirmation: got mutations %v; want %d", f.mutations, want)
			}
		})
	}
}

// TestRejectedInputs verifies failures cannot reach the mutation boundary.
func TestRejectedInputs(t *testing.T) {
	for _, policy := range []string{
		`delete_post_when(function(p) error("fixture callback") end)`,
		`delete_post_when(function(p) return "true" end)`,
		`preserve_thread("999999")`,
		`preserve_post(1)`,
		`delete_post({"7", 18})`,
		`preserve_thread({[2]="1"})`,
		`preserve_post({id="7"})`,
		`delete_repost({"7"})`,
		`preserve_post("1", function(p) return 1 end)`,
		`preserve_thread("1", {follow_edge=function(e) error("fixture edge") end})`,
		`while true do end`,
	} {
		t.Run(policy[:min(30, len(policy))], func(t *testing.T) {
			f, server := newFixture(t)
			out, err := runCLI(t, server, "", "apply", "--config", writePolicy(t, policy), "--yes", "--policy-timeout", "100ms")
			if err == nil || len(f.mutations) != 0 {
				t.Fatalf("accepted invalid policy: %v\n%s", err, out)
			}
		})
	}
	for _, mode := range []string{"schema", "account", "rate-limit", "inventory", "pagination", "edit-id", "edit-conflict"} {
		t.Run(mode, func(t *testing.T) {
			f, server := newFixture(t)
			plan := filepath.Join(t.TempDir(), "plan.json")
			if mode == "inventory" || mode == "pagination" || strings.HasPrefix(mode, "edit-") {
				f.malformed = mode == "inventory"
				f.repeatPage = mode == "pagination"
				if mode == "edit-id" {
					f.pages[0][0]["edit_history_post_ids"] = []string{"invalid", "101"}
				}
				if mode == "edit-conflict" {
					f.pages[0][1]["edit_history_post_ids"] = []string{"1", "2"}
				}
				out, err := runCLI(t, server, "", "apply", "--config", writePolicy(t, `delete_post_when(function(p) return true end)`), "-y")
				if err == nil || len(f.mutations) != 0 {
					t.Fatalf("invalid inventory accepted: %v\n%s", err, out)
				}
				return
			}
			text := `{"version":1,"account_id":"100","evaluated_at":"2026-10-30T00:00:00Z","actions":[{"kind":"delete_post","target_id":"7","forced":false}]}`
			if mode == "schema" {
				text = strings.Replace(text, `"delete_post"`, `"unknown"`, 1)
			}
			if mode == "account" {
				f.account = "200"
			}
			if mode == "rate-limit" {
				f.failMutation = true
			}
			if err := os.WriteFile(plan, []byte(text), 0o600); err != nil {
				t.Fatal(err)
			}
			out, err := runCLI(t, server, "", "apply", "--plan", plan)
			if err == nil || len(f.mutations) != 0 {
				t.Fatalf("invalid apply accepted: %v\n%s", err, out)
			}
		})
	}
}
