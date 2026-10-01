package e2e

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
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
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

// Banner failure modes, specified before implementation:
//   - front matter is ignored, emitted into body blocks, or selects a later image;
//   - the cover is attached as a body entity instead of top-level cover_media;
//   - shared banner/body bytes are uploaded twice, or distinct bytes are conflated;
//   - a valid cache is ignored, or an expired identifier is reused;
//   - malformed metadata, missing/unsupported media, escapes, changed bytes, or
//     forged media types spend an upload or draft request before being refused;
//   - a valid banner hides an invalid body image or vice versa.
// These checks drive conversion, artifact serialization, publisher resolution,
// and the real HTTP client together. Only the remote service is substituted.

type bannerExchange struct {
	Path    string         `json:"path"`
	Digest  string         `json:"digest,omitempty"`
	MediaID string         `json:"media_id,omitempty"`
	Payload map[string]any `json:"payload,omitempty"`
}

type bannerService struct {
	mu        sync.Mutex
	exchanges []bannerExchange
}

func newBannerService(t *testing.T) (*bannerService, *xapi.Client) {
	t.Helper()
	service := &bannerService{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.mu.Lock()
		defer service.mu.Unlock()
		service.exchanges = append(service.exchanges, bannerExchange{Path: r.URL.Path})
		exchange := &service.exchanges[len(service.exchanges)-1]
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method %s", r.Method)
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		switch r.URL.Path {
		case "/media/upload":
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Errorf("parse media upload: %v", err)
				http.Error(w, "bad media", http.StatusBadRequest)
				return
			}
			defer r.MultipartForm.RemoveAll()
			file, _, err := r.FormFile("media")
			if err != nil {
				t.Errorf("read upload form: %v", err)
				http.Error(w, "missing media", http.StatusBadRequest)
				return
			}
			content, readErr := io.ReadAll(file)
			closeErr := file.Close()
			if readErr != nil || closeErr != nil {
				t.Errorf("read upload bytes: %v; close: %v", readErr, closeErr)
				http.Error(w, "bad media", http.StatusBadRequest)
				return
			}
			if category := r.FormValue("media_category"); category != "tweet_image" {
				t.Errorf("upload category = %q, want tweet_image", category)
			}
			exchange.Digest = fmt.Sprintf("%x", sha256.Sum256(content))
			exchange.MediaID = fmt.Sprintf("190000000000000%04d", len(service.exchanges))
			if _, err := fmt.Fprintf(w, `{"data":{"id":%q,"expires_after_secs":3600}}`, exchange.MediaID); err != nil {
				t.Errorf("write media response: %v", err)
			}
		case "/articles/draft":
			if err := json.NewDecoder(r.Body).Decode(&exchange.Payload); err != nil {
				t.Errorf("decode draft request: %v", err)
				http.Error(w, "bad draft", http.StatusBadRequest)
				return
			}
			if _, err := io.WriteString(w, `{"data":{"id":"article-draft-1"}}`); err != nil {
				t.Errorf("write draft response: %v", err)
			}
		default:
			t.Errorf("unexpected endpoint %s", r.URL.Path)
			http.NotFound(w, r)
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
	return service, client
}

func (s *bannerService) recorded() []bannerExchange {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]bannerExchange(nil), s.exchanges...)
}

func bannerPost(t *testing.T, metadata, body string) (string, string, []byte) {
	t.Helper()
	root := t.TempDir()
	postDir := filepath.Join(root, "post")
	if err := os.Mkdir(postDir, 0o755); err != nil {
		t.Fatalf("create post: %v", err)
	}
	var encoded bytes.Buffer
	pixel := image.NewRGBA(image.Rect(0, 0, 1, 1))
	pixel.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(&encoded, pixel); err != nil {
		t.Fatalf("encode body image: %v", err)
	}
	for name, content := range map[string][]byte{
		"banner.png": tinyPNG,
		"body.png":   encoded.Bytes(),
		"vector.svg": []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`),
	} {
		if err := os.WriteFile(filepath.Join(postDir, name), content, 0o644); err != nil {
			t.Fatalf("write image %s: %v", name, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "outside.png"), tinyPNG, 0o644); err != nil {
		t.Fatalf("write outside image: %v", err)
	}
	if err := os.Symlink(filepath.Join(root, "outside.png"), filepath.Join(postDir, "escape.png")); err != nil {
		t.Fatalf("create escaping symlink: %v", err)
	}
	source := []byte("---\ntitle: Banner example\n" + metadata + "---\n\n" + body + "\n")
	if err := os.WriteFile(filepath.Join(postDir, "index.md"), source, 0o644); err != nil {
		t.Fatalf("write post: %v", err)
	}
	return root, postDir, source
}

func bannerArtifact(t *testing.T, root, postDir string, source []byte) (string, map[string]any) {
	t.Helper()
	article, err := markdown.New(postDir, "post").Convert(source)
	if err != nil {
		t.Fatalf("convert banner post: %v", err)
	}
	encoded, err := article.JSON()
	if err != nil {
		t.Fatalf("serialize banner artifact: %v", err)
	}
	path := filepath.Join(root, "artifact.json")
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	var value map[string]any
	if err := json.Unmarshal(encoded, &value); err != nil {
		t.Fatalf("decode artifact: %v", err)
	}
	return path, value
}

func writeBannerEvidence(t *testing.T, artifact map[string]any, exchanges []bannerExchange) {
	t.Helper()
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		return
	}
	encoded, err := json.MarshalIndent(map[string]any{"artifact": artifact, "requests": exchanges}, "", "  ")
	if err != nil {
		t.Fatalf("encode banner evidence: %v", err)
	}
	name := strings.ReplaceAll(t.Name(), "/", "-") + ".json"
	if err := os.WriteFile(filepath.Join(output, name), append(encoded, '\n'), 0o644); err != nil {
		t.Fatalf("write banner evidence: %v", err)
	}
}

func TestBannerDraftHTTPPipeline(t *testing.T) {
	for _, test := range []struct {
		name       string
		metadata   string
		body       string
		banner     bool
		uploads    int
		bodyImages int
	}{
		{"banner_only", "images: [banner.png]\n", "Article text.", true, 1, 0},
		{"shared_body_image", "images: [banner.png]\n", "![diagram](banner.png)", true, 1, 1},
		{"distinct_body_image", "images: [banner.png]\n", "![diagram](body.png)", true, 2, 1},
		{"first_image_only", "images: [banner.png, missing.png]\n", "Article text.", true, 1, 0},
		{"no_banner", "", "![diagram](body.png)", false, 1, 1},
		{"empty_images", "images: []\n", "Article text.", false, 0, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, postDir, source := bannerPost(t, test.metadata, test.body)
			path, artifactJSON := bannerArtifact(t, root, postDir, source)
			service, client := newBannerService(t)
			artifact, err := draft.ReadArtifact(path)
			if err != nil {
				t.Fatalf("read artifact: %v", err)
			}
			result, err := draft.New(client, root, filepath.Join(root, "cache.json")).CreateDraft(artifact)
			if err != nil {
				t.Fatalf("create draft: %v", err)
			}
			if result.ID != "article-draft-1" {
				t.Errorf("draft ID = %q", result.ID)
			}
			exchanges := service.recorded()
			writeBannerEvidence(t, artifactJSON, exchanges)
			if len(exchanges) != test.uploads+1 {
				t.Fatalf("got %d requests, want %d uploads followed by one draft", len(exchanges), test.uploads)
			}
			mediaByDigest := map[string]string{}
			for _, exchange := range exchanges[:test.uploads] {
				if exchange.Path != "/media/upload" {
					t.Fatalf("request before media resolution: %s", exchange.Path)
				}
				mediaByDigest[exchange.Digest] = exchange.MediaID
			}
			request := exchanges[len(exchanges)-1]
			if request.Path != "/articles/draft" || request.Payload["title"] != "Banner example" {
				t.Fatalf("unexpected draft request: %#v", request)
			}
			cover, hasCover := request.Payload["cover_media"].(map[string]any)
			if test.banner {
				bannerID := mediaByDigest[fmt.Sprintf("%x", sha256.Sum256(tinyPNG))]
				if !hasCover || bannerID == "" || cover["media_id"] != bannerID || cover["media_category"] != "tweet_image" || len(cover) != 2 {
					t.Errorf("cover_media = %#v, want uploaded banner %q", cover, bannerID)
				}
			} else if _, present := request.Payload["cover_media"]; present {
				t.Errorf("cover_media must be omitted without a banner: %#v", request.Payload["cover_media"])
			}
			document := request.Payload["content_state"].(map[string]any)
			blocks := document["blocks"].([]any)
			if len(blocks) != 1 {
				t.Errorf("front matter added body blocks: got %d, want 1", len(blocks))
			}
			if len(artifact.Locators) != test.bodyImages {
				t.Fatalf("body image count = %d, want %d", len(artifact.Locators), test.bodyImages)
			}
			for _, locator := range artifact.Locators {
				assertMediaAttached(t, document, locator.EntityKey, mediaByDigest[locator.Digest])
			}
		})
	}
}

func TestBannerMediaCacheLifetime(t *testing.T) {
	root, postDir, source := bannerPost(t, "images: [banner.png]\n", "![diagram](banner.png)")
	path, artifactJSON := bannerArtifact(t, root, postDir, source)
	service, client := newBannerService(t)
	for run, elapsed := range []int64{0, 60, 3601} {
		artifact, err := draft.ReadArtifact(path)
		if err != nil {
			t.Fatalf("read artifact: %v", err)
		}
		publisher := draft.New(client, root, filepath.Join(root, "cache.json"))
		publisher.Now = func() time.Time { return time.Unix(1700000000+elapsed, 0) }
		if err := publisher.LoadCache(); err != nil {
			t.Fatalf("load cache: %v", err)
		}
		if _, err := publisher.CreateDraft(artifact); err != nil {
			t.Fatalf("create draft run %d: %v", run, err)
		}
	}
	exchanges := service.recorded()
	writeBannerEvidence(t, artifactJSON, exchanges)
	wantPaths := []string{"/media/upload", "/articles/draft", "/articles/draft", "/media/upload", "/articles/draft"}
	if len(exchanges) != len(wantPaths) {
		t.Fatalf("cache produced %d requests, want %d", len(exchanges), len(wantPaths))
	}
	for index, path := range wantPaths {
		if exchanges[index].Path != path {
			t.Errorf("request %d = %s, want %s", index, exchanges[index].Path, path)
		}
	}
	for _, pair := range [][2]int{{1, 0}, {2, 0}, {4, 3}} {
		cover, ok := exchanges[pair[0]].Payload["cover_media"].(map[string]any)
		if !ok || cover["media_id"] != exchanges[pair[1]].MediaID {
			t.Errorf("request %d used unexpected cover: %#v", pair[0], cover)
		}
	}
}

func TestInvalidBannerSourceIsRefused(t *testing.T) {
	for _, test := range []struct{ name, metadata string }{
		{"scalar", "images: banner.png\n"},
		{"mapping", "images: {path: banner.png}\n"},
		{"nonstring", "images: [17]\n"},
		{"empty_reference", "images: ['']\n"},
		{"missing", "images: [missing.png]\n"},
		{"unsupported", "images: [vector.svg]\n"},
		{"remote", "images: ['https://example.invalid/banner.png']\n"},
		{"escape", "images: [../outside.png]\n"},
		{"symlink_escape", "images: [escape.png]\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, postDir, source := bannerPost(t, test.metadata, "Article text.")
			if _, err := markdown.New(postDir, "post").Convert(source); err == nil {
				t.Fatal("conversion accepted invalid banner metadata or source")
			}
		})
	}
}

func TestInvalidBannerArtifactMakesNoRequests(t *testing.T) {
	for _, test := range []struct {
		name   string
		target string
		field  string
		value  any
	}{
		{"missing_banner", "banner", "path", "missing.png"},
		{"banner_escape", "banner", "path", "../outside.png"},
		{"banner_symlink_escape", "banner", "path", "escape.png"},
		{"banner_package_escape", "banner", "post_package", "../outside"},
		{"banner_digest", "banner", "digest", strings.Repeat("0", 64)},
		{"banner_unsupported_type", "banner", "media_type", "image/svg+xml"},
		{"banner_mismatched_type", "banner", "media_type", "image/jpeg"},
		{"banner_empty_type", "banner", "media_type", ""},
		{"missing_body", "body", "path", "missing.png"},
		{"body_digest", "body", "digest", strings.Repeat("0", 64)},
		{"malformed_banner", "artifact", "banner_locator", "banner.png"},
		{"incomplete_banner", "artifact", "banner_locator", map[string]any{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, postDir, source := bannerPost(t, "images: [banner.png]\n", "![diagram](body.png)")
			path, artifactJSON := bannerArtifact(t, root, postDir, source)
			var target map[string]any
			switch test.target {
			case "artifact":
				target = artifactJSON
			case "body":
				target = artifactJSON["image_locators"].([]any)[0].(map[string]any)
			case "banner":
				var ok bool
				target, ok = artifactJSON["banner_locator"].(map[string]any)
				if !ok {
					t.Fatal("conversion did not record banner_locator")
				}
			}
			target[test.field] = test.value
			encoded, err := json.Marshal(artifactJSON)
			if err != nil {
				t.Fatalf("encode tampered artifact: %v", err)
			}
			if err := os.WriteFile(path, encoded, 0o644); err != nil {
				t.Fatalf("write tampered artifact: %v", err)
			}
			service, client := newBannerService(t)
			artifact, err := draft.ReadArtifact(path)
			if err == nil {
				_, err = draft.New(client, root, filepath.Join(root, "cache.json")).CreateDraft(artifact)
			}
			if err == nil {
				t.Error("tampered artifact was accepted")
			}
			exchanges := service.recorded()
			writeBannerEvidence(t, artifactJSON, exchanges)
			if len(exchanges) != 0 {
				t.Errorf("invalid artifact made %d requests before refusal", len(exchanges))
			}
		})
	}
}
