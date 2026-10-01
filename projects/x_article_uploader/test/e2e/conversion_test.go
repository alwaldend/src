package e2e

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
)

// postOutcome records the outcome a post is expected to produce. A post whose
// images all use an accepted media type must satisfy the mapping contract and
// may carry the continuing diagnostics its constructs produce; a post that
// references an unacceptable image must fail with the named diagnostic.
type postOutcome struct {
	name string
	// failing is true when conversion is expected to fail.
	failing bool
	// code is the failing diagnostic's code when failing is true.
	code string
	// continuingDiagnostics are diagnostics expected on a succeeding post.
	continuingDiagnostics []string
	// images is the number of image locators the artifact must carry.
	images int
}

// outcomes records the expected outcome per post directory name. The post set
// is discovered from the content tree; this table only says what each existing
// post produces, so a post added later is covered with the default outcome.
var outcomes = map[string]postOutcome{
	"diagrams-in-ac-era": {
		name:    "diagrams-in-ac-era",
		failing: true,
		code:    "image-media-type-rejected",
	},
	"dns-management-in-a-monorepo": {
		name:                  "dns-management-in-a-monorepo",
		continuingDiagnostics: []string{"inline-code-style-lost"},
	},
	"slop-without-a-clear-goal": {
		name: "slop-without-a-clear-goal",
	},
	"this-x-article-was-generated-from-markdown": {
		// The fixture post exercises representative vocabulary, so it
		// carries a raster banner and the inline-code style loss every
		// code-span-bearing post produces. Its diagram appears only in front
		// matter, which creates no image in the article body.
		name:                  "this-x-article-was-generated-from-markdown",
		continuingDiagnostics: []string{"inline-code-style-lost"},
		images:                0,
	},
}

const contentTree = "projects/alwaldend.com/content/blog"

// TestConvertEveryPost converts every post in the content tree and asserts the
// outcome its entry records.
func TestConvertEveryPost(t *testing.T) {
	root := workspaceRoot(t)
	posts, err := listPosts(filepath.Join(root, contentTree))
	if err != nil {
		t.Fatalf("list posts: %v", err)
	}
	if len(posts) == 0 {
		t.Fatal("no posts discovered in the content tree")
	}

	seen := map[string]bool{}
	for _, dir := range posts {
		name := filepath.Base(dir)
		seen[name] = true
		outcome, recorded := outcomes[name]
		if !recorded {
			// A post added later is expected to conform by default.
			outcome = postOutcome{name: name}
		}
		t.Run(name, func(t *testing.T) {
			raw := readPost(t, dir)
			converter := markdown.New(dir, filepath.ToSlash(filepath.Join(contentTree, name)))
			article, err := converter.Convert(raw)
			codes := diagnosticCodes(article)

			if outcome.failing {
				if err == nil {
					t.Fatalf("expected conversion to fail with %s, got success (diagnostics: %s)", outcome.code, joinCodes(codes))
				}
				if !hasDiagnostic(codes, outcome.code) {
					t.Fatalf("expected failing diagnostic %s, got %s", outcome.code, joinCodes(codes))
				}
				return
			}
			if err != nil {
				t.Fatalf("expected conversion to succeed, got: %v (diagnostics: %s)", err, joinCodes(codes))
			}
			for _, code := range outcome.continuingDiagnostics {
				if !hasDiagnostic(codes, code) {
					t.Errorf("expected continuing diagnostic %s, got %s", code, joinCodes(codes))
				}
			}
			if got := len(article.Locators); got != outcome.images {
				t.Errorf("expected %d image locators, got %d", outcome.images, got)
			}
			assertMappingContract(t, article)
		})
	}
	for name := range outcomes {
		if !seen[name] {
			t.Errorf("recorded outcome %q names a post that is not in the content tree", name)
		}
	}
}

// assertMappingContract checks the invariants every successful conversion must
// hold, independently of which post produced it.
func assertMappingContract(t *testing.T, article *markdown.Article) {
	t.Helper()
	if article.Document == nil || len(article.Document.Blocks) == 0 {
		t.Fatal("conversion produced no blocks")
	}
	if article.Title == "" {
		t.Error("conversion produced no parsed title")
	}

	allowed := map[string]bool{
		draftjs.BlockUnstyled:          true,
		draftjs.BlockHeaderOne:         true,
		draftjs.BlockHeaderTwo:         true,
		draftjs.BlockHeaderThree:       true,
		draftjs.BlockUnorderedListItem: true,
		draftjs.BlockOrderedListItem:   true,
		draftjs.BlockBlockquote:        true,
		draftjs.BlockAtomic:            true,
	}
	for _, block := range article.Document.Blocks {
		if !allowed[block.Type] {
			t.Errorf("block %s uses a type X does not expose: %s", block.Key, block.Type)
		}
		// Every range must select a span inside the block's final text,
		// measured in UTF-16 code units of that text.
		length := utf16Length(block.Text)
		for _, style := range block.InlineStyleRanges {
			if style.Offset < 0 || style.Offset+style.Length > length {
				t.Errorf("block %s has a style range outside its text (%d+%d > %d)", block.Key, style.Offset, style.Length, length)
			}
		}
		for _, entity := range block.EntityRanges {
			if entity.Offset < 0 || entity.Offset+entity.Length > length {
				t.Errorf("block %s has an entity range outside its text (%d+%d > %d)", block.Key, entity.Offset, entity.Length, length)
			}
			if entity.Key < 0 || entity.Key >= len(article.Document.Entities) {
				t.Errorf("block %s references entity key %d, which does not exist", block.Key, entity.Key)
			}
		}
		if block.Type == draftjs.BlockAtomic && len(block.EntityRanges) == 0 {
			t.Errorf("atomic block %s carries no entity range", block.Key)
		}
	}
}

// TestDraftArtifactIsRepeatable emits the draft artifact for every succeeding
// post to a task-owned path and asserts its document parses as content_state,
// its title matches the source front matter, and each image locator records the
// post package and the bytes' digest.
func TestDraftArtifactIsRepeatable(t *testing.T) {
	root := workspaceRoot(t)
	outDir := artifactOutputDir(root)
	// Under `bazel test` the workspace is not mounted for writing, so the
	// artifact must land in Bazel's undeclared-output directory; anywhere else
	// an explicit override or the task-owned workspace path is used. Asserting
	// the selected location keeps the advertised repeatable artifact from being
	// written into the disposable test sandbox.
	if undeclared := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); undeclared != "" && outDir != undeclared {
		t.Fatalf("artifact directory %q is not Bazel's undeclared-output directory %q", outDir, undeclared)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("create artifact directory: %v", err)
	}

	posts, err := listPosts(filepath.Join(root, contentTree))
	if err != nil {
		t.Fatalf("list posts: %v", err)
	}
	emitted := 0
	for _, dir := range posts {
		name := filepath.Base(dir)
		outcome, recorded := outcomes[name]
		if recorded && outcome.failing {
			continue
		}
		raw := readPost(t, dir)
		postPackage := filepath.ToSlash(filepath.Join(contentTree, name))
		converter := markdown.New(dir, postPackage)
		first, err := converter.Convert(raw)
		if err != nil {
			t.Fatalf("%s: convert: %v", name, err)
		}
		second, err := converter.Convert(raw)
		if err != nil {
			t.Fatalf("%s: reconvert: %v", name, err)
		}
		firstJSON, err := first.JSON()
		if err != nil {
			t.Fatalf("%s: encode artifact: %v", name, err)
		}
		secondJSON, err := second.JSON()
		if err != nil {
			t.Fatalf("%s: encode artifact again: %v", name, err)
		}
		if !bytes.Equal(firstJSON, secondJSON) {
			t.Errorf("%s: conversion is not deterministic", name)
		}

		// The emitted document must parse as a content_state payload.
		var decoded struct {
			Title    string `json:"title"`
			Document struct {
				Blocks   []draftjs.Block  `json:"blocks"`
				Entities []draftjs.Entity `json:"entities"`
			} `json:"content_state"`
			Locators []markdown.ImageLocator `json:"image_locators"`
		}
		if err := json.Unmarshal(firstJSON, &decoded); err != nil {
			t.Fatalf("%s: artifact does not parse: %v", name, err)
		}
		if decoded.Title != first.Title || decoded.Title == "" {
			t.Errorf("%s: artifact title %q does not match the parsed title %q", name, decoded.Title, first.Title)
		}
		if len(decoded.Document.Blocks) == 0 {
			t.Errorf("%s: artifact carries no blocks", name)
		}
		for _, locator := range decoded.Locators {
			if locator.PostPackage != postPackage {
				t.Errorf("%s: locator records package %q, expected %q", name, locator.PostPackage, postPackage)
			}
			if locator.Digest == "" {
				t.Errorf("%s: locator %q records no digest", name, locator.Path)
			}
		}
		path := filepath.Join(outDir, name+".json")
		if err := os.WriteFile(path, append(firstJSON, '\n'), 0o644); err != nil {
			t.Fatalf("%s: write artifact: %v", name, err)
		}
		// The artifact is the advertised verifiable output, so read it back and
		// confirm the bytes on disk are the ones conversion emitted.
		written, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s: read back artifact: %v", name, err)
		}
		if !bytes.Equal(written, append(firstJSON, '\n')) {
			t.Errorf("%s: artifact on disk differs from the emitted document", name)
		}
		emitted++
	}
	if emitted == 0 {
		t.Fatal("no artifacts were emitted")
	}
	t.Logf("emitted %d draft artifacts to %s", emitted, outDir)
}

// artifactOutputDir returns the directory the repeatable draft artifacts are
// written to. Bazel's undeclared-output directory comes first, because it is
// the only location a `bazel test` run preserves; an explicit override comes
// next; the task-owned workspace path is the default for a direct run.
func artifactOutputDir(workspaceRoot string) string {
	if outDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); outDir != "" {
		return outDir
	}
	if outDir := os.Getenv("X_ARTICLE_UPLOADER_OUT"); outDir != "" {
		return outDir
	}
	return filepath.Join(workspaceRoot, "out", "x_article_uploader", "artifacts")
}

// diagnosticCodes returns the codes of the diagnostics on an article, which
// conversion returns alongside a failing error.
func diagnosticCodes(article *markdown.Article) []string {
	if article == nil {
		return nil
	}
	codes := make([]string, 0, len(article.Diagnostics))
	for _, d := range article.Diagnostics {
		codes = append(codes, d.Code)
	}
	return codes
}

// utf16Length returns a string's length in UTF-16 code units, the unit DraftJS
// ranges are measured in.
func utf16Length(s string) int {
	n := 0
	for _, r := range s {
		if r > 0xFFFF {
			n += 2
			continue
		}
		n++
	}
	return n
}

// TestEmittedDocumentMatchesArticlesSchema asserts the emitted document uses
// the field names and value domains the Articles draft endpoint accepts. The
// endpoint's schema sets `additionalProperties: false` and names block ranges
// with underscores, so a canonical DraftJS spelling (`inlineStyleRanges`,
// `entityRanges`, `depth`) or a DraftJS `entityMap` would be rejected before
// the content is read. Asserting against the raw JSON keys, not the decoded Go
// struct, is what catches a renamed field: the decoder ignores unknown keys.
func TestEmittedDocumentMatchesArticlesSchema(t *testing.T) {
	source := "---\ntitle: Schema probe\n---\n\nA paragraph with a [link](https://example.com) and **bold** text.\n\n" +
		"## Heading\n\n- item\n  - nested\n\n1. one\n\n> quote\n\n---\n\n```go\ncode\n```\n"
	dir := t.TempDir()
	article, err := convertSource(t, dir, source)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	encoded, err := article.JSON()
	if err != nil {
		t.Fatalf("encode artifact: %v", err)
	}
	var raw struct {
		Document struct {
			Blocks   []map[string]json.RawMessage `json:"blocks"`
			Entities []map[string]json.RawMessage `json:"entities"`
		} `json:"content_state"`
	}
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatalf("decode artifact: %v", err)
	}
	if len(raw.Document.Blocks) == 0 {
		t.Fatal("document carries no blocks")
	}

	blockKeys := map[string]bool{
		"key": true, "type": true, "text": true,
		"inline_style_ranges": true, "entity_ranges": true, "data": true,
	}
	blockTypes := map[string]bool{
		"unstyled": true, "header-one": true, "header-two": true,
		"header-three": true, "unordered-list-item": true,
		"ordered-list-item": true, "blockquote": true, "atomic": true,
	}
	for i, block := range raw.Document.Blocks {
		for key := range block {
			if !blockKeys[key] {
				t.Errorf("block %d carries field %q, which the schema does not allow", i, key)
			}
		}
		if _, ok := block["inline_style_ranges"]; !ok {
			t.Errorf("block %d omits inline_style_ranges", i)
		}
		if _, ok := block["entity_ranges"]; !ok {
			t.Errorf("block %d omits entity_ranges", i)
		}
		var blockType string
		if err := json.Unmarshal(block["type"], &blockType); err != nil {
			t.Fatalf("block %d type: %v", i, err)
		}
		if !blockTypes[blockType] {
			t.Errorf("block %d uses type %q, outside the schema's enum", i, blockType)
		}
		if data, ok := block["data"]; ok {
			assertObjectKeys(t, "block "+itoa(i)+" data", data,
				map[string]bool{"cashtags": true, "hashtags": true, "mentions": true, "urls": true})
		}
	}

	valueKeys := map[string]bool{
		"caption": true, "entity_key": true, "markdown": true,
		"media_items": true, "post_id": true, "url": true,
	}
	entityTypes := map[string]bool{
		"post": true, "link": true, "image": true, "emoji": true,
		"markdown": true, "divider": true, "latex": true,
	}
	mutabilities := map[string]bool{"immutable": true, "mutable": true, "segmented": true}
	for i, entity := range raw.Document.Entities {
		if len(entity) != 2 || entity["key"] == nil || entity["value"] == nil {
			t.Fatalf("entity %d is not a {key, value} entry: %v", i, keysOf(entity))
		}
		var key string
		if err := json.Unmarshal(entity["key"], &key); err != nil {
			t.Errorf("entity %d key is not a string: %v", i, err)
		}
		var value struct {
			Type       string                     `json:"type"`
			Mutability string                     `json:"mutability"`
			Data       map[string]json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(entity["value"], &value); err != nil {
			t.Fatalf("entity %d value: %v", i, err)
		}
		if !entityTypes[value.Type] {
			t.Errorf("entity %d uses type %q, outside the schema's enum", i, value.Type)
		}
		if !mutabilities[value.Mutability] {
			t.Errorf("entity %d uses mutability %q, outside the schema's enum", i, value.Mutability)
		}
		for dataKey := range value.Data {
			if !valueKeys[dataKey] {
				t.Errorf("entity %d data carries field %q, which the schema does not allow", i, dataKey)
			}
		}
	}
}

// assertObjectKeys asserts an encoded JSON object carries only allowed keys.
func assertObjectKeys(t *testing.T, what string, encoded json.RawMessage, allowed map[string]bool) {
	t.Helper()
	var object map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &object); err != nil {
		t.Fatalf("%s is not an object: %v", what, err)
	}
	for key := range object {
		if !allowed[key] {
			t.Errorf("%s carries field %q, which the schema does not allow", what, key)
		}
	}
}

// keysOf returns a JSON object's field names for a diagnostic.
func keysOf(object map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	return keys
}

// itoa formats an index for a diagnostic.
func itoa(n int) string {
	return strconv.Itoa(n)
}
