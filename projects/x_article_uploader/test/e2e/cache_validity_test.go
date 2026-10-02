package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draft"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

// cacheService records only public fixture request paths and responses. OAuth
// headers are deliberately excluded from the repeatable evidence.
type cacheService struct {
	mu       sync.Mutex
	paths    []string
	uploads  int
	drafts   int
	lifetime string
	delay    time.Duration
}

func newCacheService(t *testing.T, lifetime string, delay time.Duration) (*cacheService, *xapi.Client) {
	t.Helper()
	service := &cacheService{lifetime: lifetime, delay: delay}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.mu.Lock()
		defer service.mu.Unlock()
		service.paths = append(service.paths, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/media/upload"):
			service.uploads++
			if _, err := io.Copy(io.Discard, r.Body); err != nil {
				t.Errorf("read media request: %v", err)
			}
			time.Sleep(service.delay)
			if _, err := fmt.Fprintf(w, `{"data":{"id":"fixture-media-%d"%s}}`, service.uploads, service.lifetime); err != nil {
				t.Errorf("write media response: %v", err)
			}
		case strings.HasSuffix(r.URL.Path, "/articles/draft"):
			service.drafts++
			if _, err := io.WriteString(w, `{"data":{"id":"fixture-draft"}}`); err != nil {
				t.Errorf("write draft response: %v", err)
			}
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	client := xapi.New(xapi.Credentials{
		APIKey: "fixture-cache-api-key", APISecret: "fixture-cache-api-secret",
		AccessToken: "fixture-cache-access-token", AccessTokenSecret: "fixture-cache-access-secret",
	})
	client.BaseURL = server.URL
	return service, client
}

func (s *cacheService) counts() (int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.uploads, s.drafts
}

func (s *cacheService) evidence(t *testing.T) {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		return
	}
	encoded, err := json.MarshalIndent(map[string]any{
		"requests": s.paths, "uploads": s.uploads, "drafts": s.drafts,
	}, "", "  ")
	if err != nil {
		t.Fatalf("encode cache evidence: %v", err)
	}
	name := strings.ReplaceAll(t.Name(), "/", "-") + ".json"
	if err := os.WriteFile(filepath.Join(output, name), append(encoded, '\n'), 0o644); err != nil {
		t.Fatalf("write cache evidence: %v", err)
	}
}

func createCacheDraft(t *testing.T, client draft.MediaClient, root, path, cache string) error {
	t.Helper()
	artifact, err := draft.ReadArtifact(path)
	if err != nil {
		t.Fatalf("read cache fixture: %v", err)
	}
	publisher := draft.New(client, root, cache)
	if err := publisher.LoadCache(); err != nil {
		t.Fatalf("load cache: %v", err)
	}
	_, err = publisher.CreateDraft(artifact)
	return err
}

func TestMediaCacheBoundToUploadContext(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*xapi.Client)
	}{
		{"api_key", func(c *xapi.Client) { c.Credentials.APIKey += "-rotated" }},
		{"api_secret", func(c *xapi.Client) { c.Credentials.APISecret += "-rotated" }},
		{"access_token", func(c *xapi.Client) { c.Credentials.AccessToken += "-rotated" }},
		{"access_secret", func(c *xapi.Client) { c.Credentials.AccessTokenSecret += "-rotated" }},
		{"endpoint", func(c *xapi.Client) { c.BaseURL += "/other-api" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, path, _, _ := publishedPost(t, "diagram.png", tinyPNG)
			cache := filepath.Join(root, "cache.json")
			service, first := newCacheService(t, `,"expires_after_secs":86400`, 0)
			defer service.evidence(t)
			second := *first
			test.change(&second)
			// Returning to the first context must retain its still-valid entry.
			for _, client := range []*xapi.Client{first, &second, first} {
				if err := createCacheDraft(t, client, root, path, cache); err != nil {
					t.Fatalf("create draft: %v", err)
				}
			}
			if uploads, drafts := service.counts(); uploads != 2 || drafts != 3 {
				t.Errorf("got %d uploads and %d drafts; want two scoped uploads and three drafts", uploads, drafts)
			}
			raw, err := os.ReadFile(cache)
			if err != nil {
				t.Fatalf("read persisted cache: %v", err)
			}
			for _, credential := range []string{
				first.Credentials.APIKey, first.Credentials.APISecret,
				first.Credentials.AccessToken, first.Credentials.AccessTokenSecret,
			} {
				if strings.Contains(string(raw), credential) {
					t.Error("persisted cache contains a plaintext credential fixture")
				}
			}
		})
	}
}

func TestMediaCacheRejectsLegacyEntries(t *testing.T) {
	root, path, locator, _ := publishedPost(t, "diagram.png", tinyPNG)
	cache := filepath.Join(root, "cache.json")
	legacy, err := json.Marshal(map[string]any{locator.Digest: map[string]any{
		"media_id": "unscoped-media", "expires_at": time.Now().Add(time.Hour).Unix(),
	}})
	if err != nil {
		t.Fatalf("encode legacy cache: %v", err)
	}
	if err := os.WriteFile(cache, legacy, 0o644); err != nil {
		t.Fatalf("write legacy cache: %v", err)
	}
	service, client := newCacheService(t, `,"expires_after_secs":86400`, 0)
	defer service.evidence(t)
	if err := createCacheDraft(t, client, root, path, cache); err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if uploads, drafts := service.counts(); uploads != 1 || drafts != 1 {
		t.Errorf("legacy cache produced %d uploads and %d drafts; want one fresh upload and draft", uploads, drafts)
	}
}

func TestMediaCacheUnknownLifetimeIsResolutionLocal(t *testing.T) {
	for _, lifetime := range []string{"", `,"expires_after_secs":0`} {
		name := "omitted"
		if lifetime != "" {
			name = "zero"
		}
		t.Run(name, func(t *testing.T) {
			root, postDir, source := bannerPost(t, "images: [banner.png]\n", "Text.\n\n![repeat](banner.png)")
			path, _ := bannerArtifact(t, root, postDir, source)
			service, client := newCacheService(t, lifetime, 0)
			defer service.evidence(t)
			publisher := draft.New(client, root, filepath.Join(root, "cache.json"))
			// Even a repeated operation in the same process must not retain an
			// identifier whose lifetime X never supplied.
			for run := 0; run < 2; run++ {
				artifact, err := draft.ReadArtifact(path)
				if err != nil {
					t.Fatalf("read artifact: %v", err)
				}
				if _, err := publisher.CreateDraft(artifact); err != nil {
					t.Fatalf("create draft: %v", err)
				}
			}
			if err := createCacheDraft(t, client, root, path, filepath.Join(root, "cache.json")); err != nil {
				t.Fatalf("create draft from new publisher: %v", err)
			}
			if uploads, drafts := service.counts(); uploads != 3 || drafts != 3 {
				t.Errorf("unknown lifetime produced %d uploads and %d drafts; want one upload per resolution", uploads, drafts)
			}
		})
	}
}

func TestMediaCacheLifetimeCannotOutliveResponse(t *testing.T) {
	for _, test := range []struct {
		name     string
		lifetime int
		delay    time.Duration
	}{
		{"negative", -1, 0},
		{"expiry_reserve", 30, 0},
		{"upload_elapsed", 31, 1100 * time.Millisecond},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, path, _, _ := publishedPost(t, "diagram.png", tinyPNG)
			service, client := newCacheService(t, fmt.Sprintf(`,"expires_after_secs":%d`, test.lifetime), test.delay)
			defer service.evidence(t)
			if err := createCacheDraft(t, client, root, path, filepath.Join(root, "cache.json")); err == nil {
				t.Error("created a draft without a safely positive media lifetime")
			}
			if uploads, drafts := service.counts(); uploads != 1 || drafts != 0 {
				t.Errorf("unusable lifetime produced %d uploads and %d drafts; want one upload and no draft", uploads, drafts)
			}
		})
	}
}

type expiringSequenceClient struct {
	recordedClient
	clock *time.Time
}

func (c *expiringSequenceClient) UploadImage(name string, content []byte) (*xapi.MediaUpload, error) {
	c.uploads = append(c.uploads, name)
	lifetime := 31
	if len(c.uploads) == 2 {
		*c.clock = c.clock.Add(2 * time.Second)
		lifetime = 3600
	}
	return &xapi.MediaUpload{MediaID: fmt.Sprintf("media-%d", len(c.uploads)), ExpiresAfterSecs: lifetime}, nil
}

func TestMediaCacheRechecksAllMediaBeforeDraft(t *testing.T) {
	for _, test := range []struct{ name, body string }{
		{"distinct_body", "![body](body.png)"},
		{"later_banner_duplicate", "![body](body.png)\n\n![banner again](banner.png)"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, postDir, source := bannerPost(t, "images: [banner.png]\n", test.body)
			path, artifactJSON := bannerArtifact(t, root, postDir, source)
			clock := time.Unix(1700000000, 0)
			client := &expiringSequenceClient{clock: &clock}
			publisher := draft.New(client, root, filepath.Join(root, "cache.json"))
			publisher.Now = func() time.Time { return clock }
			artifact, err := draft.ReadArtifact(path)
			if err != nil {
				t.Fatalf("read artifact: %v", err)
			}
			if _, err := publisher.CreateDraft(artifact); err == nil {
				t.Error("draft used a banner whose safe lifetime elapsed while uploading the body")
			}
			if len(client.uploads) != 2 || client.drafts != 0 {
				t.Errorf("got %d uploads and %d drafts, want two uploads and no draft", len(client.uploads), client.drafts)
			}
			exchanges := make([]bannerExchange, 0, len(client.uploads)+client.drafts)
			for range client.uploads {
				exchanges = append(exchanges, bannerExchange{Path: "/media/upload"})
			}
			for range client.drafts {
				exchanges = append(exchanges, bannerExchange{Path: "/articles/draft"})
			}
			writeBannerEvidence(t, artifactJSON, exchanges)
		})
	}
}
