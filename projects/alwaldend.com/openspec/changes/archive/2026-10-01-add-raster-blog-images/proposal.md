## Why

A blog post's images are currently SVG by default, because a diagram is rendered
with the Mermaid pipeline and referenced as `![alt](./diagram.svg)`. X Articles
cannot accept an SVG: its media upload endpoint accepts a fixed set of image
media types that excludes `image/svg+xml`, so a post whose images are SVG cannot
be republished through the API at all.

The section therefore needs a stated rule that a post intended for syndication
references raster images only, and a way to render a diagram to a raster image
so the rule costs an author nothing but a rule choice.

## What Changes

- State that a post's images use a media type suitable for syndication, so an
  SVG is not the default image reference for a new post.
- Reference a Mermaid diagram through its raster render beside `index.md`, with
  the `.mmd` source remaining authoritative.
- Keep the existing SVG renders for posts that are not syndicated, rather than
  rewriting archived posts.
- Update the blog authoring guidance to state the rule and name the raster
  render target.
- Provide `mermaid_site_webp` for the native site palette, and use it for the
  X uploader post's padded Markdown-to-draft diagram. Select that same image
  in the post's social metadata.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `blog-section`: a post's image references are held to the media types a
  syndication target can accept, so the section's output is not limited to
  readers of the site alone.

## Impact

- `projects/alwaldend.com/content/blog/**`: new posts keep a diagram's raster
  render beside `index.md` in place of an SVG reference.
- `projects/alwaldend.com/skills/alwaldend-blog/SKILL.md`: the authoring
  guidance states the rule.
- The raster render target itself is owned by `tools/mermaid`; this change
  consumes it and records the rule its consumers follow. `mermaid_webp` is that
  target, and the authoring guidance in
  `projects/alwaldend.com/skills/alwaldend-blog/SKILL.md` states the rule and
  names it.
- Historical posts that reference SVG keep doing so. The X uploader post stays
  a draft; its preview image changes without publishing it or calling X.
