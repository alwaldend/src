package draft

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/xapi"
)

// mediaCacheEntry records an uploaded media identifier and a conservative reuse
// deadline. A zero ExpiresAt is usable only in the current image resolution and
// is never read from or written to persistent storage.
type mediaCacheEntry struct {
	MediaID   string `json:"media_id"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

// mediaCacheFile separates content digests by API endpoint, media category, and
// credential context. Legacy caches without this scope are deliberately misses.
type mediaCacheFile struct {
	Version int                                   `json:"version"`
	Scopes  map[string]map[string]mediaCacheEntry `json:"scopes"`
}

// mediaCacheScope lets a client opt into persistent reuse with a nonsecret,
// stable namespace. Clients without one get resolution-local reuse only.
type mediaCacheScope interface {
	MediaCacheScope() string
}

// Leave time for sending the resulting draft, instead of consuming a media
// identifier at the very end of the lifetime the API reported.
const mediaExpiryReserve = 30 * time.Second

// MediaClient is the subset of the API client draft creation uses, so the
// checks can substitute a recorded response.
type MediaClient interface {
	UploadImage(name string, content []byte) (*xapi.MediaUpload, error)
	CreateDraft(request xapi.DraftRequest) (*xapi.Draft, error)
}

// Publisher resolves a draft artifact's images and creates an unpublished draft.
type Publisher struct {
	Client MediaClient
	// WorkspaceRoot is the repository root the post package resolves against.
	WorkspaceRoot string
	// CachePath is the scoped digest-to-media_id cache. An unreadable cache
	// costs a re-upload rather than failing publication.
	CachePath string
	// Cache holds identifiers reused within this run, keyed by content digest.
	cache       map[string]mediaCacheEntry
	cacheScope  string
	cacheScopes map[string]map[string]mediaCacheEntry
	// Now supplies the current time for expiry checks; checks set it for a
	// stable run.
	Now func() time.Time
}

// New returns a draft creator.
func New(client MediaClient, workspaceRoot, cachePath string) *Publisher {
	return &Publisher{
		Client:        client,
		WorkspaceRoot: workspaceRoot,
		CachePath:     cachePath,
		cache:         map[string]mediaCacheEntry{},
		cacheScopes:   map[string]map[string]mediaCacheEntry{},
		Now:           time.Now,
	}
}

// LoadCache reads the digest-to-media_id mapping from the cache file. A missing,
// unreadable, or superseded cache is not an error: it costs a re-upload rather
// than failing publication.
func (p *Publisher) LoadCache() error {
	p.cache = map[string]mediaCacheEntry{}
	p.cacheScopes = map[string]map[string]mediaCacheEntry{}
	p.cacheScope = p.scope()
	if p.cacheScope == "" {
		return nil
	}
	raw, err := os.ReadFile(p.CachePath)
	if err != nil {
		return nil
	}
	var decoded mediaCacheFile
	if err := json.Unmarshal(raw, &decoded); err != nil || decoded.Version != 1 || decoded.Scopes == nil {
		// Older caches do not identify the account or endpoint that owns the
		// uploaded media. They cannot be trusted across credential changes.
		return nil
	}
	p.cacheScopes = decoded.Scopes
	for digest, entry := range decoded.Scopes[p.cacheScope] {
		if entry.ExpiresAt > 0 && p.unexpired(entry) {
			p.cache[digest] = entry
		}
	}
	return nil
}

func (p *Publisher) scope() string {
	if client, ok := p.Client.(mediaCacheScope); ok {
		return client.MediaCacheScope()
	}
	return ""
}

// beginResolution invalidates unknown lifetimes between operations and reloads
// the correct namespace if the caller changed its client or credentials.
func (p *Publisher) beginResolution() {
	scope := p.scope()
	if scope == "" {
		p.cache = map[string]mediaCacheEntry{}
		p.cacheScope = ""
		return
	}
	if scope != p.cacheScope {
		_ = p.LoadCache()
	}
	for digest, entry := range p.cache {
		if entry.ExpiresAt <= 0 || !p.unexpired(entry) {
			delete(p.cache, digest)
		}
	}
}

// saveCache writes the digest-to-media_id mapping with each entry's expiry. A
// failure to record the cache costs a later re-upload rather than failing the
// publication, so every failure is swallowed: cache storage is an optimization,
// not a precondition.
func (p *Publisher) saveCache() {
	if p.CachePath == "" || p.cacheScope == "" {
		return
	}
	persisted := map[string]mediaCacheEntry{}
	for digest, entry := range p.cache {
		if entry.ExpiresAt > 0 && p.unexpired(entry) {
			persisted[digest] = entry
		}
	}
	p.cacheScopes[p.cacheScope] = persisted
	encoded, err := json.MarshalIndent(mediaCacheFile{Version: 1, Scopes: p.cacheScopes}, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p.CachePath), 0o755); err != nil {
		return
	}
	// The write can fail for reasons the earlier checks cannot see, such as the
	// cache path naming an existing directory or a read-only mount.
	_ = os.WriteFile(p.CachePath, append(encoded, '\n'), 0o644)
}

// ResolveUploads validates and uploads body and banner images, attaching body
// media identifiers to the entities their locators name. The shared cache keeps
// a later CreateDraft call from uploading these images again.
func (p *Publisher) ResolveUploads(artifact *Artifact) error {
	if _, err := p.resolveUploads(artifact); err != nil {
		return fmt.Errorf("resolve uploads: %w", err)
	}
	return nil
}

// resolveUploads preflights every source before spending any API calls. The
// verified bytes are retained for uploading, so a later filesystem change cannot
// substitute different bytes between verification and upload.
func (p *Publisher) resolveUploads(artifact *Artifact) (*xapi.CoverMedia, error) {
	if artifact.Document == nil {
		return nil, fmt.Errorf("resolve images: the artifact carries no content_state")
	}
	entities, err := entitiesOf(artifact.Document)
	if err != nil {
		return nil, fmt.Errorf("read image entities: %w", err)
	}
	// An unresolved image entity must be resolved by exactly one locator, and
	// every locator must name an image entity. Validating this before any
	// upload keeps a malformed artifact from spending API calls and then
	// failing, and keeps an entity without media_items from being sent.
	if err := verifyLocators(entities, artifact.Locators); err != nil {
		return nil, fmt.Errorf("verify image locators: %w", err)
	}
	// The banner is not a body entity. It contributes one source to the same
	// upload resolution, and its identifier is returned separately for the
	// draft request's cover_media field.
	sources := make([]markdown.ImageSource, 0, len(artifact.Locators)+1)
	if artifact.Banner != nil {
		sources = append(sources, *artifact.Banner)
	}
	for _, locator := range artifact.Locators {
		sources = append(sources, locator.ImageSource)
	}
	contents := make([][]byte, len(sources))
	for index, source := range sources {
		content, _, err := ImageBytes(source, p.WorkspaceRoot)
		if err != nil {
			return nil, fmt.Errorf("resolve image %q: %w", source.Path, err)
		}
		contents[index] = content
	}
	p.beginResolution()
	var cover *xapi.CoverMedia
	bodyOffset := 0
	for index, source := range sources {
		mediaID, err := p.uploadImage(source, contents[index])
		if err != nil {
			return nil, fmt.Errorf("resolve upload %q: %w", source.Path, err)
		}
		if index == 0 && artifact.Banner != nil {
			cover = &xapi.CoverMedia{MediaCategory: draftjs.MediaCategoryImage, MediaID: mediaID}
			bodyOffset = 1
			continue
		}
		locator := artifact.Locators[index-bodyOffset]
		if err := attachMedia(entities, locator.EntityKey, mediaID); err != nil {
			return nil, fmt.Errorf("attach uploaded image %q: %w", source.Path, err)
		}
	}
	// A later upload or processing wait may consume an earlier item's safe
	// lifetime. Refuse the draft instead of referencing that expired item or
	// silently spending another upload request.
	for _, source := range sources {
		if !p.unexpired(p.cache[source.Digest]) {
			return nil, fmt.Errorf("resolve image %q: media lifetime expired before draft creation", source.Path)
		}
	}
	p.saveCache()
	return cover, nil
}

// uploadImage reuses the same digest cache for banner and body sources.
func (p *Publisher) uploadImage(source markdown.ImageSource, content []byte) (string, error) {
	if entry, ok := p.cache[source.Digest]; ok {
		if p.unexpired(entry) {
			return entry.MediaID, nil
		}
		// beginResolution already removed old entries. Replacing one that
		// expires during this operation could leave an earlier banner or body
		// attachment pointing at the old identifier, so stop without retrying.
		return "", fmt.Errorf("image %q: media lifetime expired during image resolution", source.Path)
	}
	started := p.now()
	upload, err := p.Client.UploadImage(source.Path, content)
	if err != nil {
		return "", fmt.Errorf("upload image %q: %w", source.Path, err)
	}
	if upload == nil || upload.MediaID == "" {
		return "", fmt.Errorf("upload image %q: the API returned no media_id", source.Path)
	}
	entry, err := p.cacheEntry(upload, started)
	if err != nil {
		return "", fmt.Errorf("cache uploaded image %q: %w", source.Path, err)
	}
	p.cache[source.Digest] = entry
	return upload.MediaID, nil
}

// CreateDraft resolves images, then creates a draft. It does not publish.
func (p *Publisher) CreateDraft(artifact *Artifact) (*xapi.Draft, error) {
	if strings.TrimSpace(artifact.Title) == "" {
		return nil, fmt.Errorf("create draft: the artifact carries no parsed title")
	}
	cover, err := p.resolveUploads(artifact)
	if err != nil {
		return nil, fmt.Errorf("prepare draft media: %w", err)
	}
	draft, err := p.Client.CreateDraft(xapi.DraftRequest{
		Title: artifact.Title, ContentState: artifact.Document, CoverMedia: cover,
	})
	if err != nil {
		return nil, fmt.Errorf("create draft: %w", err)
	}
	return draft, nil
}

// unexpired reports whether a cached identifier may be reused. Unknown expiry
// is allowed only for freshly uploaded entries in this resolution; LoadCache
// refuses it and beginResolution drops it before any subsequent operation.
func (p *Publisher) unexpired(entry mediaCacheEntry) bool {
	if entry.MediaID == "" {
		return false
	}
	if entry.ExpiresAt == 0 {
		return true
	}
	return p.now().Unix() < entry.ExpiresAt
}

// cacheEntry anchors lifetime before upload, so network and processing delays
// cannot extend it. A custom client without UploadedAt uses the caller's start.
func (p *Publisher) cacheEntry(upload *xapi.MediaUpload, started time.Time) (mediaCacheEntry, error) {
	entry := mediaCacheEntry{MediaID: upload.MediaID}
	if upload.ExpiresAfterSecs < 0 || int64(upload.ExpiresAfterSecs) > int64((1<<63-1)/time.Second) {
		return mediaCacheEntry{}, fmt.Errorf("the API returned an invalid media lifetime")
	}
	if upload.ExpiresAfterSecs == 0 {
		return entry, nil
	}
	if !upload.UploadedAt.IsZero() {
		started = upload.UploadedAt
	}
	entry.ExpiresAt = started.Add(time.Duration(upload.ExpiresAfterSecs)*time.Second - mediaExpiryReserve).Unix()
	if !p.unexpired(entry) {
		return mediaCacheEntry{}, fmt.Errorf("the uploaded media lifetime is too short or has expired")
	}
	return entry, nil
}

// now returns the creator's clock, defaulting to the wall clock.
func (p *Publisher) now() time.Time {
	if p.Now == nil {
		return time.Now()
	}
	return p.Now()
}

// entityValue returns the entity's payload object. The Articles endpoint wraps
// each entity as {key, value}; the type and data live under value.
func entityValue(entity map[string]any) (map[string]any, bool) {
	value, ok := entity["value"].(map[string]any)
	return value, ok
}

// entitiesOf returns the document's entity slice as mutable maps.
func entitiesOf(document map[string]any) ([]map[string]any, error) {
	raw, ok := document["entities"]
	if !ok {
		return nil, fmt.Errorf("resolve images: the artifact carries no entities")
	}
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("resolve images: the artifact's entities are not a list")
	}
	entities := make([]map[string]any, 0, len(list))
	for _, item := range list {
		entity, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("resolve images: an entity is not an object")
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

// attachMedia sets the media_items of the entity the locator names. A locator
// that names no entity, or an entity that is not an image, is an error rather
// than a silent substitution.
func attachMedia(entities []map[string]any, entityKey int, mediaID string) error {
	for _, entity := range entities {
		key, ok := entityKeyOf(entity)
		if !ok || key != entityKey {
			continue
		}
		value, ok := entityValue(entity)
		if !ok {
			return fmt.Errorf("resolve image: entity %d carries no value object", entityKey)
		}
		entityType, _ := value["type"].(string)
		if entityType != "image" {
			return fmt.Errorf("resolve image: entity %d is a %s entity, not an image", entityKey, entityType)
		}
		data, ok := value["data"].(map[string]any)
		if !ok {
			data = map[string]any{}
			value["data"] = data
		}
		// The slice uses the same shape the artifact decodes to, so the emitted
		// document matches a freshly decoded one.
		data["media_items"] = []any{map[string]any{
			"media_category": "tweet_image",
			"media_id":       mediaID,
		}}
		return nil
	}
	return fmt.Errorf("resolve image: no entity carries key %d", entityKey)
}

// verifyLocators checks that an artifact's image entities and locators agree
// exactly before any upload begins. Every image entity still missing
// media_items must be named by one locator, and every locator must name one
// distinct image entity. This refuses an artifact that retained an unresolved
// image while dropping its locator, which would otherwise be sent without
// media_items, and an artifact whose locator names a missing or non-image
// entity, rather than discovering the mismatch after paying for uploads.
func verifyLocators(entities []map[string]any, locators []markdown.ImageLocator) error {
	imageKeys := map[int]bool{}
	unresolved := map[int]bool{}
	seenKeys := map[int]bool{}
	for _, entity := range entities {
		value, _ := entityValue(entity)
		entityType, _ := value["type"].(string)
		key, ok := entityKeyOf(entity)
		if !ok {
			// An image entity whose key cannot be read is invisible to the
			// locator coverage below, so a hand-edited artifact could send an
			// unresolved image. Refuse it rather than skipping it.
			if entityType == "image" {
				return fmt.Errorf("resolve images: an image entity carries no nonnegative integral key")
			}
			continue
		}
		// Two entities that share one numeric key collapse to one locator, so
		// only the first would be attached and the other would be sent without
		// media_items. Refuse the duplicate before any upload.
		if seenKeys[key] {
			return fmt.Errorf("resolve images: entity key %d is used by more than one entity", key)
		}
		seenKeys[key] = true
		if entityType != "image" {
			continue
		}
		imageKeys[key] = true
		if !hasMediaItems(entity) {
			unresolved[key] = true
		}
	}
	claimed := map[int]bool{}
	for _, locator := range locators {
		if !imageKeys[locator.EntityKey] {
			return fmt.Errorf("resolve images: a locator names entity %d, which is not an image entity", locator.EntityKey)
		}
		if claimed[locator.EntityKey] {
			return fmt.Errorf("resolve images: entity %d is named by more than one locator", locator.EntityKey)
		}
		claimed[locator.EntityKey] = true
	}
	for key := range unresolved {
		if !claimed[key] {
			return fmt.Errorf("resolve images: image entity %d has no locator, so it would be sent without media_items", key)
		}
	}
	return nil
}

// entityKeyOf reads an entity's key as a nonnegative integer. The Articles
// endpoint records the key as a decimal string; a missing, non-numeric,
// fractional, or negative key is unusable, so the entity cannot be named by a
// locator and cannot be trusted as resolved.
func entityKeyOf(entity map[string]any) (int, bool) {
	raw, ok := entity["key"].(string)
	if !ok {
		return 0, false
	}
	key, err := strconv.Atoi(raw)
	if err != nil || key < 0 {
		return 0, false
	}
	return key, true
}

// hasMediaItems reports whether an image entity already carries a usable
// uploaded media identifier. Every supplied entry is checked, not just the
// first and not just the array's length: an array with one valid entry beside
// an empty, non-object, or incorrectly categorized one is not a resolution, so
// treating it as one would send an invalid entity to X.
func hasMediaItems(entity map[string]any) bool {
	value, ok := entityValue(entity)
	if !ok {
		return false
	}
	data, ok := value["data"].(map[string]any)
	if !ok {
		return false
	}
	items, ok := data["media_items"].([]any)
	if !ok || len(items) == 0 {
		return false
	}
	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			return false
		}
		mediaID, _ := item["media_id"].(string)
		mediaCategory, _ := item["media_category"].(string)
		if mediaID == "" || mediaCategory != draftjs.MediaCategoryImage {
			return false
		}
	}
	return true
}
