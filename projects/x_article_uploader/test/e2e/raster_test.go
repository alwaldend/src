package e2e

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
	"git.alwaldend.com/alwaldend/src/tools/mermaid/test/svgdoc"
)

// rasterPostSource is a post whose only image is the diagram's raster render.
// It is a fixture rather than a published post: the policy it demonstrates is
// prospective, and the check reads the real rule output beside it.
const rasterPostSource = `---
title: Raster diagram policy
description: A post whose diagram is referenced as a raster render.
---

The pipeline below is a diagram.

![Pipeline](diagram.webp)

It converts like any other post.
`

// TestRasterRenderConvertsPost renders a diagram through the maintained raster
// rule and converts a post that references the render. It asserts the render
// is a WebP — the media type the X media upload endpoints accept — that it
// reproduces the SVG render's geometry rather than introducing a second
// appearance, and that the post converts without an unacceptable-image
// diagnostic.
func TestRasterRenderConvertsPost(t *testing.T) {
	render := rlocation(t, *rasterRender)
	rasterBytes, err := os.ReadFile(render)
	if err != nil {
		t.Fatalf("read raster render: %v", err)
	}
	// Chrome encodes the raster, so a valid WebP is the observable proof that
	// the pinned browser produced it rather than a host converter.
	width, height, err := svgdoc.WebPDimensions(rasterBytes)
	if err != nil {
		t.Fatalf("raster render is not a readable WebP: %v", err)
	}

	// The raster is a projection of the SVG render, so the two agree on the
	// drawing's size. A raster that no longer follows the maintained document
	// fails here even though it is still a WebP.
	svgBytes, err := os.ReadFile(rlocation(t, *rasterSource))
	if err != nil {
		t.Fatalf("read source render: %v", err)
	}
	viewWidth, viewHeight, err := svgdoc.CanvasSize(string(svgBytes))
	if err != nil {
		t.Fatalf("read source render canvas: %v", err)
	}
	const scale = 2
	if math.Abs(float64(width)-viewWidth*scale) > scale || math.Abs(float64(height)-viewHeight*scale) > scale {
		t.Errorf("raster render measures %dx%d, but the SVG geometry at scale %d is %.0fx%.0f",
			width, height, scale, viewWidth*scale, viewHeight*scale)
	}

	// Convert a post that references the render from its own directory, so the
	// converter resolves and type-checks the real rule output.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte(rasterPostSource), 0o644); err != nil {
		t.Fatalf("write fixture post: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "diagram.webp"), rasterBytes, 0o644); err != nil {
		t.Fatalf("write raster render: %v", err)
	}
	converter := markdown.New(dir, "content/blog/raster-post")
	converter.PostSource = "raster-post/index.md"
	article, err := converter.Convert([]byte(rasterPostSource))
	if err != nil {
		t.Fatalf("convert a post with a raster diagram: %v", err)
	}
	// The diagnostic is the failing one this policy exists to avoid, so its
	// absence is the property under test rather than a side effect.
	codes := diagnosticCodes(article)
	if hasDiagnostic(codes, "image-media-type-rejected") {
		t.Fatalf("raster diagram produced an unacceptable-image diagnostic: %s", joinCodes(codes))
	}
	if len(article.Locators) != 1 {
		t.Fatalf("post carries %d image locators, expected one", len(article.Locators))
	}
	locator := article.Locators[0]
	if locator.MediaType != "image/webp" {
		t.Errorf("locator records media type %q, expected image/webp", locator.MediaType)
	}
	if locator.Path != "diagram.webp" {
		t.Errorf("locator records path %q, expected diagram.webp", locator.Path)
	}
	if locator.Digest == "" {
		t.Error("locator records no digest for the raster render")
	}
}
