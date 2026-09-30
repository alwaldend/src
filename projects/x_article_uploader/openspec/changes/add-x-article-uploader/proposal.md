## Why

Blog posts live as Markdown under `projects/alwaldend.com/content/blog/*/index.md`,
but X Articles has no Markdown import. Its `POST /2/articles/draft` endpoint
accepts only a DraftJS `content_state` and requires an authenticated call.
Republishing a post by hand means re-creating its structure, inline styles,
links, and code blocks in X's editor.

A deterministic converter plus an authenticated draft command turns each post
into a reviewable draft, so the same Markdown that renders the site also
produces the article body. The uploader creates drafts only, so a separate
review step in X is what makes an article public.

## What Changes

- Add `projects/x_article_uploader`, a Go project that compiles a Markdown post
  into a DraftJS `content_state` document and creates an X Article draft.
- Split conversion from draft creation: conversion is offline and hermetic;
  draft creation is the only path that performs network calls and requires
  credentials.
- Create drafts only. The uploader calls no publish endpoint and exposes no
  publish operation, so a draft becomes public only through a separate review
  step in X.
- Treat a source post as front matter plus body: the parsed title is carried
  alongside the document, because draft creation requires a title that
  `content_state` does not hold.
- Upload referenced images as media and reference them by `media_id`, because
  the `image` entity sources media from `media_items` and `data.url` is
  documented for `link` entities. Report images whose media type the upload
  endpoints reject instead of emitting entities that can never resolve.
- Compile Markdown with `github.com/yuin/goldmark`, approved by the user on
  2026-09-29, because it exposes the CommonMark and GFM constructs the posts
  rely on and matches the parser family Hugo renders with.
- Adopt a raster-only image policy for posts from now on, and add a
  `mermaid_webp` rule that rasterizes a Mermaid diagram into a media type the
  upload endpoints accept, so a post with diagrams can be drafted instead of
  only reported as unconvertible.
- Add Bazel targets for the library, both binaries, and hermetic tests over the
  real blog posts, plus the OpenSpec workspace and project declarations for the
  new project.
- **BREAKING**: none. The project is additive and changes no existing behavior.

## Capabilities

### New Capabilities

- `markdown-to-draftjs`: Compile a blog post's Markdown body into a
  X-compatible DraftJS `content_state` document and report its parsed title,
  deterministically and without network access.
- `x-article-publication`: Create X Article drafts using the X API, including
  uploading referenced images and resolving their `media_id`. Drafts only; no
  publication.
- `x-article-uploader-build`: Bazel targets, tests, and project declarations
  that build the converter and draft command hermetically from pinned
  dependencies.

### Modified Capabilities

None. No capability in this workspace has changed requirements; this change
adds a new project.

Related owner changes carry the cross-owner requirements this change depends on:

- `projects/alwaldend.com/openspec/changes/add-raster-blog-images` owns the
  raster-image rule for a syndicated post and the authoring guidance.
- `tools/mermaid/openspec/changes/add-raster-diagram-renders` owns the raster
  render rule and its appearance and hermeticity requirements.
- `infra/vault/openspec/changes/add-x-article-uploader-credentials` owns the
  component identity, the least-privilege policy for the X credential, and the
  injected reference.

## Impact

- `projects/x_article_uploader/**`: new project source, BUILD files, README,
  release and landing declarations, and this OpenSpec workspace.
- `tools/mermaid/**`: a `mermaid_webp` rule and its renderer support, plus the
  skill and README guidance for the raster-only policy, tracked by
  `tools/mermaid/openspec/changes/add-raster-diagram-renders`.
- `projects/alwaldend.com/skills/alwaldend-blog/SKILL.md`: the authoring
  guidance states the raster-only image policy.
- `projects/alwaldend.com/content/blog/**`: new posts reference raster images
  only; a diagram is referenced through the `mermaid_webp` render. The rule is
  tracked by `projects/alwaldend.com/openspec/changes/add-raster-blog-images`.
- `projects/projects.bzl`: the new project joins the landing-page catalog.
- `infra/src/openspec/validation/BUILD.bazel`: the new owner workspace joins
  `_WORKSPACE_SOURCES`, following the documented step for a new owner workspace.
- `infra/vault/**`: the uploader's component identity and its least-privilege
  credential policy, tracked by
  `infra/vault/openspec/changes/add-x-article-uploader-credentials`.
- The uploader's own `al.lua`, its `al_config` target, and plugin data and
  labels: the non-secret half of the injection path the Vault change describes.
- Root dependency declarations: `github.com/yuin/goldmark` and its transitive
  modules become pinned external build inputs, exposed through the owning
  `use_repo` declaration.
- Existing posts are produced only for new posts going forward; the current SVG
  diagrams in `projects/alwaldend.com/content/blog/diagrams-in-ac-era` are
  reported as unconvertible rather than retrofitted.
- X API credentials are a run-time input handled through the repository's
  secret-injection flow; no secret value is recorded in source.
