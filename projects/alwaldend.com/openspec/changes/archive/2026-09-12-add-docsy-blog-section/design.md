## Context

The apex site and the project landings are separate Hugo builds. The apex site
packages `projects/alwaldend.com` content plus repository documentation; each
landing packages a project README with the shared shell from `hugo_landing`.
Both render Docsy 0.17.0 with the same Bootstrap and Font-Awesome modules, so
they already share a component model but build from different source archives.

## Decisions

### Reuse the Docsy blog section

Docsy 0.17.0 already ships `blog` layouts, a date-grouped list template, an RSS
output, and print templates. The site's `[outputs]` declaration already includes
`rss` and `print` for sections, so the blog reuses those outputs. No site-local
layout or template override is added; the section works by placing content in
the `blog` directory and adding a menu entry.

### One post per directory

A post is its own Bazel package containing `index.md` and a `BUILD.bazel` that
declares its `docs_filegroup`. The section package aggregates post packages
through `subpackages(include = ["*"], allow_empty = True)`, matching the
content-tree aggregation used by `content/docs/misc`. A post directory is a
page bundle, so post-local resources resolve beside the post that uses them.
`allow_empty` keeps the section buildable before the first post exists.

### Share the canvas styles

The background belongs to the shared visual shell, not to either site. The rules
live in `projects/hugo_landing/assets/scss/_shared_theme.scss`; the landings pick
them up through their own `_styles_project.scss`, and the apex site packages the
same file into its `assets/scss/` tree and imports it. One implementation, two
consumers, no copied rules.

The overrides are CSS custom properties on `:root` and both `[data-bs-theme]`
selectors. Bootstrap and Docsy components read `--bs-*` at runtime, so setting
them once covers the palette, and Docsy's light/dark menu keeps working while
both selections render the same canvas. Values reuse the Bootstrap dark palette
that Docsy already compiles rather than introducing a parallel set of colors.

## Verification

- Build `//projects/alwaldend.com:site` and confirm the compiled stylesheet sets
  a pure black body background and that the black rule is last among the
  competing `--bs-body-bg` declarations.
- Inspect the served page in a browser and confirm the computed body background
  is `rgb(0, 0, 0)` with light foreground text.
- Build `//projects:landing_sites` and confirm a landing renders the same
  canvas from the shared source.
- Run the site's generated-site test for internal links, fragments, duplicated
  IDs, and image alternatives.

## Out of scope

Live deployment to GitHub Pages, DNS changes, alternate article taxonomies, and
comments are not part of this change.
