package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draft"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

// recordedClient records the requests publication made and returns canned
// responses, so the checks need neither credentials nor network access.
type recordedClient struct {
	uploads       []string
	uploadBytes   [][]byte
	draftTitle    string
	draftBody     any
	drafts        int
	uploadID      string
	uploadExpires int
	draftID       string
	failUpload    bool
}

func (c *recordedClient) MediaCacheScope() string {
	return "recorded-client-fixture"
}

func (c *recordedClient) UploadImage(name string, content []byte) (*xapi.MediaUpload, error) {
	if c.failUpload {
		return nil, &xapi.APIError{Operation: "upload image", Status: 400, Body: "rejected"}
	}
	c.uploads = append(c.uploads, name)
	c.uploadBytes = append(c.uploadBytes, content)
	return &xapi.MediaUpload{MediaID: c.uploadID, MediaCategory: "tweet_image", ExpiresAfterSecs: c.uploadExpires}, nil
}

func (c *recordedClient) CreateDraft(request xapi.DraftRequest) (*xapi.Draft, error) {
	c.drafts++
	c.draftTitle = request.Title
	c.draftBody = request.ContentState
	return &xapi.Draft{ID: c.draftID}, nil
}

// publishedPost builds an artifact referencing one raster image, writes the
// image beside a temporary post package, and returns the artifact path.
func publishedPost(t *testing.T, imageName string, imageBytes []byte) (root string, artifactPath string, locator markdown.ImageLocator, document map[string]any) {
	t.Helper()
	root = t.TempDir()
	postPackage := "content/blog/example"
	dir := filepath.Join(root, filepath.FromSlash(postPackage))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create post directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, imageName), imageBytes, 0o644); err != nil {
		t.Fatalf("write image: %v", err)
	}

	// Convert a source post that references the image so the locator's digest
	// is the one the converter recorded.
	source := "---\ntitle: Example\n---\n\n![diagram](./" + imageName + ")\n"
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte(source), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	converter := markdown.New(dir, postPackage)
	article, err := converter.Convert([]byte(source))
	if err != nil {
		t.Fatalf("convert example post: %v", err)
	}
	if len(article.Locators) != 1 {
		t.Fatalf("expected one locator, got %d", len(article.Locators))
	}
	encoded, err := article.JSON()
	if err != nil {
		t.Fatalf("encode artifact: %v", err)
	}
	artifactPath = filepath.Join(root, "artifact.json")
	if err := os.WriteFile(artifactPath, encoded, 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	var decoded struct {
		Document map[string]any `json:"content_state"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("decode artifact: %v", err)
	}
	return root, artifactPath, article.Locators[0], decoded.Document
}

// tinyPNG is a one-pixel PNG, so the converter treats the reference as an
// accepted raster image.
var tinyPNG = []byte{
	0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 'I', 'D', 'A', 'T',
	0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00, 0x05,
	0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00,
	0x00, 0x00, 'I', 'E', 'N', 'D', 0xae, 0x42, 0x60, 0x82,
}

// TestCreateDraftCarriesTitleAndResolvesMedia asserts the draft request body
// carries the parsed title and content_state, that the uploaded media_id lands
// on the entity its locator names, and that draft creation alone does not
// publish.
func TestCreateDraftCarriesTitleAndResolvesMedia(t *testing.T) {
	root, artifactPath, locator, _ := publishedPost(t, "diagram.png", tinyPNG)
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))
	if err := pub.LoadCache(); err != nil {
		t.Fatalf("load cache: %v", err)
	}

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	draft, err := pub.CreateDraft(artifact)
	if err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if draft.ID != "draft-1" {
		t.Errorf("expected draft id draft-1, got %s", draft.ID)
	}
	if client.draftTitle != "Example" {
		t.Errorf("draft request carried title %q, expected the parsed title", client.draftTitle)
	}
	if len(client.uploads) != 1 {
		t.Errorf("expected one upload, got %d", len(client.uploads))
	}
	if len(client.uploadBytes) != 1 || !bytes.Equal(client.uploadBytes[0], tinyPNG) {
		t.Errorf("upload did not carry the verified image bytes")
	}
	if _, ok := client.draftBody.(map[string]any); !ok {
		t.Errorf("draft request carried content_state of type %T", client.draftBody)
	}
	// The draft command mutates the document it read, so the attachment is asserted
	// on that same document rather than on an earlier decode.
	assertMediaAttached(t, artifact.Document, locator.EntityKey, "media-1")
}

// TestCreateDraftFailsWithoutTitle asserts the draft command refuses to send a
// placeholder title when the artifact carries none.
func TestCreateDraftFailsWithoutTitle(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	artifact.Title = ""
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected draft creation without a title to fail")
	}
	if client.drafts != 0 {
		t.Errorf("draft command sent %d draft request(s) despite the missing title", client.drafts)
	}
}

// TestLocatorDigestMismatchIsRefused asserts bytes that no longer match the
// recorded digest are not uploaded in place of the image.
func TestLocatorDigestMismatchIsRefused(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	// Rewrite the image with different bytes of the same name.
	if err := os.WriteFile(filepath.Join(root, "content", "blog", "example", "diagram.png"), []byte("not a png"), 0o644); err != nil {
		t.Fatalf("rewrite image: %v", err)
	}
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected a digest mismatch to fail publication")
	}
	if len(client.uploads) != 0 {
		t.Errorf("draft command uploaded %d replacement image(s) despite the mismatch", len(client.uploads))
	}
}

// TestLocatorEscapeIsRefused asserts a locator that escapes the post package is
// refused rather than resolved.
func TestLocatorEscapeIsRefused(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	artifact.Locators[0].Path = "../../../etc/passwd"
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected an escaping locator to fail publication")
	}
	if len(client.uploads) != 0 {
		t.Errorf("draft command uploaded %d image(s) despite the escaping locator", len(client.uploads))
	}
}

// TestLocatorSymlinkEscapeIsRefused asserts an apparently safe relative path
// that is a symlink out of the post package is refused, because containment is
// a property of the resolved target rather than of the recorded path.
func TestLocatorSymlinkEscapeIsRefused(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	postDir := filepath.Join(root, "content", "blog", "example")
	outside := filepath.Join(root, "outside.png")
	if err := os.WriteFile(outside, tinyPNG, 0o644); err != nil {
		t.Fatalf("write outside image: %v", err)
	}
	link := filepath.Join(postDir, "link.png")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatalf("create symlink: %v", err)
	}
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	// The digest matches the symlink's target, so only containment can refuse
	// it; a path-only check would upload the outside bytes.
	artifact.Locators[0].Path = "link.png"
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected a symlink out of the post package to fail publication")
	}
	if len(client.uploads) != 0 {
		t.Errorf("draft command uploaded %d image(s) through an escaping symlink", len(client.uploads))
	}
}

// TestLocatorPackageEscapeIsRefused asserts a recorded post package that leaves
// the workspace is refused rather than resolved.
func TestLocatorPackageEscapeIsRefused(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	artifact.Locators[0].PostPackage = "content/blog/../../../../etc"
	artifact.Locators[0].Path = "passwd"
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected a post package outside the workspace to fail publication")
	}
	if len(client.uploads) != 0 {
		t.Errorf("draft command uploaded %d image(s) through an escaping package", len(client.uploads))
	}
}

// TestRejectedUploadFailsWithoutDraft asserts a failed upload does not create a
// draft referencing a placeholder media id.
func TestRejectedUploadFailsWithoutDraft(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	client := &recordedClient{failUpload: true, draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected a rejected upload to fail publication")
	}
	if client.drafts != 0 {
		t.Errorf("draft command created %d draft(s) despite the failed upload", client.drafts)
	}
}

// TestUploadedMediaIsReusedAcrossRuns asserts identical bytes reuse one
// identifier within a run and across a later run through the cache.
func TestUploadedMediaIsReusedAcrossRuns(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	cachePath := filepath.Join(root, "cache.json")

	first := &recordedClient{uploadID: "media-1", uploadExpires: 3600, draftID: "draft-1"}
	firstPub := draft.New(first, root, cachePath)
	if err := firstPub.LoadCache(); err != nil {
		t.Fatalf("load cache: %v", err)
	}
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if _, err := firstPub.CreateDraft(artifact); err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if len(first.uploads) != 1 {
		t.Fatalf("expected one upload on the first run, got %d", len(first.uploads))
	}

	// A later run in a new process reads the digest mapping and reuses it.
	second := &recordedClient{uploadID: "media-2", uploadExpires: 3600, draftID: "draft-2"}
	secondPub := draft.New(second, root, cachePath)
	if err := secondPub.LoadCache(); err != nil {
		t.Fatalf("load cache: %v", err)
	}
	if _, err := secondPub.CreateDraft(artifact); err != nil {
		t.Fatalf("create draft again: %v", err)
	}
	if len(second.uploads) != 0 {
		t.Errorf("expected the later run to reuse the cached identifier, but it uploaded %d time(s)", len(second.uploads))
	}
}

// TestMissingCredentialNamesItsReference asserts a missing field is reported by
// its reference rather than by a partial credential.
func TestMissingCredentialNamesItsReference(t *testing.T) {
	env := map[string]string{
		draft.EnvAPIKey:      "key",
		draft.EnvAPISecret:   "secret",
		draft.EnvAccessToken: "token",
	}
	_, err := draft.LoadCredentials(func(name string) string { return env[name] })
	if err == nil {
		t.Fatal("expected loading without the access token secret to fail")
	}
	if !strings.Contains(err.Error(), draft.EnvAccessTokenSecret) {
		t.Errorf("expected the error to name %s, got: %v", draft.EnvAccessTokenSecret, err)
	}
}

// TestMediaUploadDecodesV2Envelope asserts the draft command reads the identifier
// from the v2 envelope the upload endpoint returns — `data.id` with a lifetime —
// rather than from a top-level `media_id`, and that the identifier lands on the
// entity. The recorded-response client cannot catch this: it supplies a typed
// MediaUpload and never exercises the HTTP decode, so this check drives the real
// client against a served response.
func TestMediaUploadDecodesV2Envelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/media/upload" {
			t.Errorf("upload request went to %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"id":"1900000000000000001","media_key":"3_1900","expires_after_secs":86400}}`)
	}))
	defer server.Close()

	client := &xapi.Client{
		BaseURL: server.URL,
		Credentials: xapi.Credentials{
			APIKey: "key", APISecret: "secret",
			AccessToken: "token", AccessTokenSecret: "token-secret",
		},
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1700000000, 0) },
	}
	upload, err := client.UploadImage("diagram.png", tinyPNG)
	if err != nil {
		t.Fatalf("upload image: %v", err)
	}
	if upload.MediaID != "1900000000000000001" {
		t.Errorf("decoded media id %q, expected the envelope's data.id", upload.MediaID)
	}
	if upload.ExpiresAfterSecs != 86400 {
		t.Errorf("decoded lifetime %d, expected 86400", upload.ExpiresAfterSecs)
	}

	// The identifier must land on the image entity, so an artifact referencing
	// the uploaded bytes carries data.id rather than an empty identifier.
	root, artifactPath, locator, _ := publishedPost(t, "diagram.png", tinyPNG)
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if err := pub.ResolveUploads(artifact); err != nil {
		t.Fatalf("resolve uploads: %v", err)
	}
	assertMediaAttached(t, artifact.Document, locator.EntityKey, "1900000000000000001")
}

// TestExpiredMediaIsUploadedAgain asserts a cached identifier whose recorded
// lifetime has passed is treated as a miss, while an unexpired one is reused.
func TestExpiredMediaIsUploadedAgain(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	cachePath := filepath.Join(root, "cache.json")

	base := time.Unix(1700000000, 0)
	first := &recordedClient{uploadID: "media-1", uploadExpires: 3600, draftID: "draft-1"}
	firstPub := draft.New(first, root, cachePath)
	firstPub.Now = func() time.Time { return base }
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if _, err := firstPub.CreateDraft(artifact); err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if len(first.uploads) != 1 {
		t.Fatalf("expected one upload on the first run, got %d", len(first.uploads))
	}

	// Before the lifetime elapses, the cached identifier is reused.
	fresh := &recordedClient{uploadID: "media-2", uploadExpires: 3600, draftID: "draft-2"}
	freshPub := draft.New(fresh, root, cachePath)
	freshPub.Now = func() time.Time { return base.Add(30 * time.Minute) }
	if err := freshPub.LoadCache(); err != nil {
		t.Fatalf("load cache: %v", err)
	}
	if _, err := freshPub.CreateDraft(artifact); err != nil {
		t.Fatalf("create draft before expiry: %v", err)
	}
	if len(fresh.uploads) != 0 {
		t.Errorf("expected the unexpired identifier to be reused, but it uploaded %d time(s)", len(fresh.uploads))
	}

	// After the lifetime elapses, the stale identifier must not be reused.
	expired := &recordedClient{uploadID: "media-3", uploadExpires: 3600, draftID: "draft-3"}
	expiredPub := draft.New(expired, root, cachePath)
	expiredPub.Now = func() time.Time { return base.Add(2 * time.Hour) }
	if err := expiredPub.LoadCache(); err != nil {
		t.Fatalf("load cache: %v", err)
	}
	if _, err := expiredPub.CreateDraft(artifact); err != nil {
		t.Fatalf("create draft after expiry: %v", err)
	}
	if len(expired.uploads) != 1 {
		t.Errorf("expected the expired identifier to be re-uploaded, but it uploaded %d time(s)", len(expired.uploads))
	}
}

// TestDraftAndPublishDecodeV2Envelope asserts the real client reads the draft
// identifier from `data.id` and the post identifier from `data.post_id`, which
// the recorded-response client cannot catch because it returns typed structs.
func TestDraftDecodesV2Envelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/articles/draft":
			io.WriteString(w, `{"data":{"id":"1815550000000000000","title":"Example"}}`)
		default:
			t.Errorf("unexpected request path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := &xapi.Client{
		BaseURL: server.URL,
		Credentials: xapi.Credentials{
			APIKey: "key", APISecret: "secret",
			AccessToken: "token", AccessTokenSecret: "token-secret",
		},
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1700000000, 0) },
	}
	draft, err := client.CreateDraft(xapi.DraftRequest{Title: "Example", ContentState: map[string]any{"blocks": []any{}, "entities": []any{}}})
	if err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if draft.ID != "1815550000000000000" {
		t.Errorf("decoded draft id %q, expected the envelope's data.id", draft.ID)
	}
}

// TestUnresolvedImageWithoutLocatorIsRefused asserts an image entity that no
// locator resolves is refused before a network request, rather than sending an
// entity without media_items for X to reject.
func TestUnresolvedImageWithoutLocatorIsRefused(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	// A reviewed or hand-edited artifact can retain the unresolved image entity
	// while dropping its locator.
	artifact.Locators = nil
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected an unresolved image with no locator to fail publication")
	}
	if client.drafts != 0 {
		t.Errorf("draft command created %d draft(s) with an unresolved image", client.drafts)
	}
	if len(client.uploads) != 0 {
		t.Errorf("draft command uploaded %d image(s) despite having no locator", len(client.uploads))
	}
}

// TestUnwritableCacheDoesNotBlockPublication asserts a cache that cannot be
// written costs a later re-upload rather than aborting draft creation after the
// uploads already succeeded.
func TestUnwritableCacheDoesNotBlockPublication(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	// Point the cache at an existing directory: writing a file over it fails.
	cacheDir := filepath.Join(root, "cache-dir")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatalf("create cache directory: %v", err)
	}
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, cacheDir)

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	draft, err := pub.CreateDraft(artifact)
	if err != nil {
		t.Fatalf("an unwritable cache must not fail draft creation: %v", err)
	}
	if draft.ID != "draft-1" || client.drafts != 1 {
		t.Errorf("draft was not created despite the unwritable cache: %+v", draft)
	}
	if len(client.uploads) != 1 {
		t.Errorf("expected the image to still be uploaded, got %d upload(s)", len(client.uploads))
	}
}

// TestDuplicateImageEntityKeysAreRefused asserts an artifact whose two image
// entities share one key is refused before uploading. One locator would
// otherwise resolve only the first entity, leaving the second sent without
// media_items.
func TestDuplicateImageEntityKeysAreRefused(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	entities, ok := artifact.Document["entities"].([]any)
	if !ok {
		t.Fatalf("artifact entities are not a list: %T", artifact.Document["entities"])
	}
	duplicated := make([]any, 0, len(entities)+1)
	for _, raw := range entities {
		entity, _ := raw.(map[string]any)
		duplicated = append(duplicated, entity)
		value, _ := entity["value"].(map[string]any)
		if value["type"] == "image" {
			duplicated = append(duplicated, map[string]any{
				"key":   entity["key"],
				"value": map[string]any{"type": "image", "mutability": value["mutability"], "data": map[string]any{}},
			})
		}
	}
	artifact.Document["entities"] = duplicated

	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected duplicate image entity keys to fail publication")
	}
	if client.drafts != 0 || len(client.uploads) != 0 {
		t.Errorf("duplicate entity keys reached the network: drafts=%d uploads=%d", client.drafts, len(client.uploads))
	}
}

// assertMediaAttached checks the entity a locator names carries the media id.
func assertMediaAttached(t *testing.T, document map[string]any, entityKey int, mediaID string) {
	t.Helper()
	entities, ok := document["entities"].([]any)
	if !ok {
		t.Fatalf("document carries no entities")
	}
	for _, raw := range entities {
		entity, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		key, ok := entity["key"].(string)
		if !ok || key != strconv.Itoa(entityKey) {
			continue
		}
		value, ok := entity["value"].(map[string]any)
		if !ok {
			t.Fatalf("entity %d carries no value", entityKey)
		}
		data, ok := value["data"].(map[string]any)
		if !ok {
			t.Fatalf("entity %d carries no data", entityKey)
		}
		items, ok := data["media_items"].([]any)
		if !ok || len(items) != 1 {
			t.Fatalf("entity %d carries media_items %v", entityKey, data["media_items"])
		}
		item, ok := items[0].(map[string]any)
		if !ok {
			t.Fatalf("entity %d media item is not an object", entityKey)
		}
		if item["media_id"] != mediaID {
			t.Errorf("entity %d carries media_id %v, expected %s", entityKey, item["media_id"], mediaID)
		}
		if item["media_category"] != "tweet_image" {
			t.Errorf("entity %d carries media_category %v, expected tweet_image", entityKey, item["media_category"])
		}
		return
	}
	t.Fatalf("no entity carries key %d", entityKey)
}

// TestMalformedMediaItemsAreRefused asserts an image entity whose media_items
// entry lacks a usable media_id is not treated as resolved, so an invalid
// entity is not sent to X.
func TestMalformedMediaItemsAreRefused(t *testing.T) {
	root, artifactPath, locator, document := publishedPost(t, "diagram.png", tinyPNG)
	client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))

	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	// A hand-edited artifact claims the image is resolved but records no usable
	// identifier, and drops the locator that would have uploaded it.
	entities, ok := artifact.Document["entities"].([]any)
	if !ok {
		t.Fatalf("artifact entities are not a list: %T", artifact.Document["entities"])
	}
	for _, raw := range entities {
		entity, _ := raw.(map[string]any)
		value, _ := entity["value"].(map[string]any)
		if value["type"] != "image" {
			continue
		}
		value["data"] = map[string]any{"media_items": []any{map[string]any{"media_category": "tweet_image", "media_id": ""}}}
	}
	artifact.Locators = nil
	_ = locator
	_ = document

	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected an image with an unusable media_items entry to fail publication")
	}
	if client.drafts != 0 || len(client.uploads) != 0 {
		t.Errorf("malformed media_items reached the network: drafts=%d uploads=%d", client.drafts, len(client.uploads))
	}
}

// TestMediaItemsRequireCategory asserts an image entity whose media_items entry
// carries a media_id but no valid media category is treated as unresolved, so
// the malformed entry is not sent to X.
func TestMediaItemsRequireCategory(t *testing.T) {
	for _, tc := range []struct {
		name string
		item map[string]any
	}{
		{"missing category", map[string]any{"media_id": "media-1"}},
		{"empty category", map[string]any{"media_id": "media-1", "media_category": ""}},
		{"wrong category", map[string]any{"media_id": "media-1", "media_category": "tweet_video"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
			client := &recordedClient{uploadID: "media-2", draftID: "draft-1"}
			pub := draft.New(client, root, filepath.Join(root, "cache.json"))
			artifact, err := draft.ReadArtifact(artifactPath)
			if err != nil {
				t.Fatalf("read artifact: %v", err)
			}
			entities, _ := artifact.Document["entities"].([]any)
			for _, raw := range entities {
				entity, _ := raw.(map[string]any)
				value, _ := entity["value"].(map[string]any)
				if value["type"] == "image" {
					value["data"] = map[string]any{"media_items": []any{tc.item}}
				}
			}
			artifact.Locators = nil
			if _, err := pub.CreateDraft(artifact); err == nil {
				t.Fatal("expected a media_items entry without a valid category to fail publication")
			}
			if client.drafts != 0 || len(client.uploads) != 0 {
				t.Errorf("malformed media_items reached the network: drafts=%d uploads=%d", client.drafts, len(client.uploads))
			}
		})
	}
}

// TestImageBytesReadsResolvedPath asserts the bytes come from the canonical
// target containment was checked against, not the symlink path that could be
// swapped after validation.
func TestImageBytesReadsResolvedPath(t *testing.T) {
	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	postDir := filepath.Join(root, "content", "blog", "example")
	realPath := filepath.Join(postDir, "real.png")
	if err := os.WriteFile(realPath, tinyPNG, 0o644); err != nil {
		t.Fatalf("write real image: %v", err)
	}
	linkPath := filepath.Join(postDir, "link.png")
	if err := os.Symlink("real.png", linkPath); err != nil {
		t.Fatalf("create symlink: %v", err)
	}
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	locator := artifact.Locators[0]
	locator.Path = "link.png"
	_, resolved, err := draft.ImageBytes(locator.ImageSource, root)
	if err != nil {
		t.Fatalf("resolve image bytes: %v", err)
	}
	wantResolved, err := filepath.EvalSymlinks(realPath)
	if err != nil {
		t.Fatalf("resolve real path: %v", err)
	}
	if resolved != wantResolved {
		t.Errorf("ImageBytes returned path %q, expected the canonical target %q", resolved, wantResolved)
	}
}

// TestNonNumericImageEntityKeyIsRefused asserts an image entity whose key is
// not a number is refused rather than silently skipped. A skipped entity is
// invisible to locator coverage, so a hand-edited artifact could otherwise send
// an unresolved image to X.
func TestNonNumericImageEntityKeyIsRefused(t *testing.T) {
	for _, tc := range []struct {
		name string
		key  any
	}{
		{"missing key", nil},
		{"non-numeric key", "not-a-number"},
		{"negative key", "-1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
			client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
			pub := draft.New(client, root, filepath.Join(root, "cache.json"))
			artifact, err := draft.ReadArtifact(artifactPath)
			if err != nil {
				t.Fatalf("read artifact: %v", err)
			}
			entities, _ := artifact.Document["entities"].([]any)
			for _, raw := range entities {
				entity, _ := raw.(map[string]any)
				value, _ := entity["value"].(map[string]any)
				if value["type"] != "image" {
					continue
				}
				if tc.key == nil {
					delete(entity, "key")
				} else {
					entity["key"] = tc.key
				}
			}
			// The locator no longer names any usable entity, so without the key
			// check the malformed image would be sent unresolved.
			artifact.Locators = nil
			if _, err := pub.CreateDraft(artifact); err == nil {
				t.Fatal("expected an image entity with a non-numeric key to fail publication")
			}
			if client.drafts != 0 || len(client.uploads) != 0 {
				t.Errorf("malformed entity reached the network: drafts=%d uploads=%d", client.drafts, len(client.uploads))
			}
		})
	}
}

// TestEveryMediaItemsEntryMustBeUsable asserts an entity whose media_items
// holds a usable entry next to a malformed one is treated as unresolved, so the
// malformed array is not sent to X merely because one entry looked valid.
func TestEveryMediaItemsEntryMustBeUsable(t *testing.T) {
	for _, tc := range []struct {
		name string
		tail any
	}{
		{"empty entry", map[string]any{}},
		{"non-object entry", "media-2"},
		{"wrong category", map[string]any{"media_id": "media-2", "media_category": "tweet_video"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
			client := &recordedClient{uploadID: "media-1", draftID: "draft-1"}
			pub := draft.New(client, root, filepath.Join(root, "cache.json"))
			artifact, err := draft.ReadArtifact(artifactPath)
			if err != nil {
				t.Fatalf("read artifact: %v", err)
			}
			entities, _ := artifact.Document["entities"].([]any)
			for _, raw := range entities {
				entity, _ := raw.(map[string]any)
				value, _ := entity["value"].(map[string]any)
				if value["type"] == "image" {
					value["data"] = map[string]any{"media_items": []any{
						map[string]any{"media_id": "media-1", "media_category": "tweet_image"},
						tc.tail,
					}}
				}
			}
			artifact.Locators = nil
			if _, err := pub.CreateDraft(artifact); err == nil {
				t.Fatal("expected a partially malformed media_items array to fail publication")
			}
			if client.drafts != 0 || len(client.uploads) != 0 {
				t.Errorf("malformed media_items reached the network: drafts=%d uploads=%d", client.drafts, len(client.uploads))
			}
		})
	}
}

// TestDraftRequiresIdentifier asserts a 2xx draft response that omits the draft
// identifier is reported as a failure rather than decoded into a successful
// zero-value result the operator cannot act on.
func TestDraftRequiresIdentifier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/articles/draft":
			io.WriteString(w, `{"data":{"title":"Example"}}`)
		default:
			t.Errorf("unexpected request path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := &xapi.Client{
		BaseURL: server.URL,
		Credentials: xapi.Credentials{
			APIKey: "key", APISecret: "secret",
			AccessToken: "token", AccessTokenSecret: "token-secret",
		},
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1700000000, 0) },
	}
	if _, err := client.CreateDraft(xapi.DraftRequest{Title: "Example", ContentState: map[string]any{"blocks": []any{}, "entities": []any{}}}); err == nil {
		t.Error("expected a draft response without data.id to fail")
	}
}

// TestMediaUploadRejectsMissingIdentifier asserts a 2xx media upload response
// that omits data.id fails publication rather than resolving the image to an
// empty identifier, matching the draft and publish envelope checks.
func TestMediaUploadRejectsMissingIdentifier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"media_category":"tweet_image"}}`)
	}))
	defer server.Close()

	root, artifactPath, _, _ := publishedPost(t, "diagram.png", tinyPNG)
	client := &xapi.Client{
		BaseURL: server.URL,
		Credentials: xapi.Credentials{
			APIKey: "key", APISecret: "secret",
			AccessToken: "token", AccessTokenSecret: "token-secret",
		},
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1700000000, 0) },
	}
	pub := draft.New(client, root, filepath.Join(root, "cache.json"))
	artifact, err := draft.ReadArtifact(artifactPath)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if _, err := pub.CreateDraft(artifact); err == nil {
		t.Fatal("expected a media upload response without data.id to fail publication")
	}
}

// TestMediaUploadRequestCarriesCategory asserts the media upload request sends
// the `media_category` field the endpoint requires. The upload endpoint rejects
// a request without it, so a client that omitted the field would fail at the
// first image upload even though every recorded-client check passed, because
// those checks substitute the client rather than exercising the request.
func TestMediaUploadRequestCarriesCategory(t *testing.T) {
	var (
		gotCategory string
		gotFile     string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/media/upload" {
			t.Errorf("unexpected request path %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart request: %v", err)
		}
		gotCategory = r.FormValue("media_category")
		file, header, err := r.FormFile("media")
		if err != nil {
			t.Fatalf("read media part: %v", err)
		}
		defer file.Close()
		gotFile = header.Filename
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"id":"1900000000000000001","media_category":"`+gotCategory+`","expires_after_secs":86400}}`)
	}))
	defer server.Close()

	client := &xapi.Client{
		BaseURL: server.URL,
		Credentials: xapi.Credentials{
			APIKey: "key", APISecret: "secret",
			AccessToken: "token", AccessTokenSecret: "token-secret",
		},
		HTTPClient: server.Client(),
		Now:        func() time.Time { return time.Unix(1700000000, 0) },
	}
	if _, err := client.UploadImage("diagram.png", tinyPNG); err != nil {
		t.Fatalf("upload image: %v", err)
	}
	if gotCategory != "tweet_image" {
		t.Errorf("media upload sent media_category %q, expected tweet_image", gotCategory)
	}
	if gotFile != "diagram.png" {
		t.Errorf("media upload sent file %q, expected diagram.png", gotFile)
	}
}
