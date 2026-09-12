## Why

Alwaldend.com publishes project documentation and repository guides but has no
place for dated, authored articles. Docsy ships a blog section with its own
layouts, RSS feed, print output, and sidebar navigation, so the site can gain an
article destination by publishing content into that section instead of adding a
one-off page type.

The repository's Hugo sites should also share one visual canvas. The apex site
and the project landing sites are separate Hugo builds, so a shared background
needs a single owner rather than repeated rules. That canvas is owned by
`projects/hugo_landing/openspec/changes/shared-site-canvas`; this change wires
the apex site to consume it.

## What Changes

- Add a `blog` content section to `projects/alwaldend.com`, packaged through the
  existing content rules with one package per post, and reachable from the main
  navigation.
- Rely on the Docsy blog layouts, RSS output, and print output the site already
  enables for sections, without adding site-local template overrides.
- Document how to add a post in a skill owned by this project.
- Keep draft posts out of the release build while rendering them for local
  preview, so an unpublished post cannot reach the deployed site.
- Consume the shared black canvas from the reusable site shell, so the apex site
  and the project landing sites render the same background from one
  implementation.

## Capabilities

### New Capabilities

- `blog-section`: A packaged article section for alwaldend.com with Docsy blog
  layouts, per-post packages, section index navigation, and syndication and print
  outputs.

### Modified Capabilities

- `project-alwaldend-com`: The site packages a blog section and consumes the
  shared canvas styles; its build, link resolution, print-anchor, preview, and
  publication behavior are otherwise unchanged.

## Impact

- `projects/alwaldend.com/content/blog/` and its `BUILD.bazel` declare the
  section and aggregate the post packages into `site_source_archive`.
- `projects/alwaldend.com/hugo.toml` gains the section's menu entry.
- `projects/alwaldend.com/BUILD.bazel` passes `--buildDrafts` only to the
  preview configuration, so the release build excludes draft content.
- `projects/alwaldend.com/assets/` packages the shared canvas style file from
  `hugo_landing` into its own assets tree; the apex tree keeps only its own
  presentation rules.
- `projects/alwaldend.com/skills/alwaldend-blog/` documents the post workflow
  and is registered with the repository skill discovery.
- The generated-site test in `projects/alwaldend.com/test/site/` continues to
  validate internal links, fragments, duplicate IDs, and image alternatives.
