## Why

Blog authors lack guidance to select a representative social image, and an
existing post with an opening diagram has none selected. The uploader post's
light WebP and excessive padding do not match the requested dark presentation.

## What Changes

- Encourage `images` metadata for every post, reusing opening images when present.
- Select a generated raster of the AC era post's opening diagram for its cards.
- Render the uploader diagram with the dark site palette and modest margins.
- Show separate media uploads, cached media reuse, and final draft creation.
- Preserve article bodies, draft states, and the paused X upload.

## Capabilities

### Modified Capabilities

- `blog-section`: selected image recommendations and native dark raster previews.

## Impact

The blog skill, two blog bundles, the site macro documentation, and browser
checks change. The Mermaid owner adds native raster color selection. No new
dependency, site deployment, or X API request is needed.
