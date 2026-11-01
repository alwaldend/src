---
title: Alwaldend.com
description: Main website and project documentation
websites:
  - alwaldend.com
statuses:
  - in_progress
tags:
  - hugo
  - github_pages
---

Alwaldend.com is the repository's main website and documentation site. It
uses Hugo with Docsy to publish project documentation through the local and
public [download servers](../../infra/download/README.md).

## Historical www architecture

This historical diagram preserves the former Proxmox unikernel deployment; see Deployment for the current website.
The [Mermaid source](docs/legacy-www.mmd) is authoritative; Bazel renders its light and dark SVGs.

[![www architecture](docs/assets/legacy-www.svg)](docs/assets/legacy-www.svg)

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/alwaldend.com
- Former GitHub Pages repo: https://github.com/alwaldend/alwaldend.github.io
- Hugo rules: [../../tools/rules_hugo](../../tools/rules_hugo)

## Features

- [Hugo](https://gohugo.io) site
- [Docsy](https://github.com/docsy/docsy), [Bootstrap](https://getbootstrap.com)

The homepage, blog, documentation, and project pages share a flat monochrome
theme with white and charcoal color modes, subtle borders, and system fonts.
The navigation's theme menu supports light, dark, and system preferences.
Every page uses the centered S8 square monogram SVG favicon from
`//data/alwaldend:monogram`, with adaptive light/dark colors.
Prose links are underlined, and keyboard focus uses a visible outline.
Scrollbars are thin and muted; taxonomy links use small rounded labels shared
with the project index. Documentation sidebars use plain action rows,
uppercase group headings, compact wrapping taxonomy items, subtle hover and active
backgrounds spanning the full tree row, and compact chevrons. Tree labels stay
indented; chevrons expand branches independently of the row's page link. Headings,
search fields, code blocks, tables, and callouts follow the same typography
and surface styles. Syntax highlighting
uses CSS classes so code colors follow the active theme.

Docsy's `sidebar_root_enabled` scopes `/docs/<category>/<owner>` sections to
their owner's subtree. One Hugo cascade applies `sidebar_root_for: self`;
the sidebar root links up to the parent section. Full navigation remains on
the Docs overview. Compact/foldable settings control visibility, while scoped
roots also reduce the HTML sent to readers. All tags remain available.

Documentation Mermaid images use the same light/dark palette and follow the
theme menu, including system preference. They use clean, rounded nodes and
containers, smooth connectors, embedded Liberation Mono lettering, and Mermaid's
native Dagre layout. Normal graph connectors are 2.25px, while thick, dashed,
invisible, and semantic marker distinctions remain native. Caption plates use
muted borders; sequence lifelines are thinner and grayer than messages.
Mindmaps use the site's surfaces, text, and connector colors in both modes.
The `mermaid_site_svg` macro in `pkg/mermaid/defs.bzl`
selects the shared `//data/alwaldend:mermaid_theme` target and renders both
SVGs through the native Mermaid API without post-processing. The preset supplies
caption styling through `themeCSS` and places titles on container borders with
Mermaid's `subGraphTitleMargin` setting. Compact node/rank spacing and the
pinned renderer's `edgeSpacing` option separate nested title plates. Package both
outputs and reference the light `.svg`; Markdown and the SVG shortcode recognize
its `.dark.svg` sibling. The site theme selects the visible image, and print
always uses light. Historical blog diagrams retain their original files and
presentation. New raster diagrams can use `mermaid_site_webp` from the same
macro file to encode the native diagram with the site's palette and monospace
font for syndication and social previews. Set `color_scheme = "dark"` to select
the dark palette; the default remains light.
The homepage panel links to GitHub, GitLab, Blog, Docs, and the project index,
followed by permanently visible, indented project links in title order.
Each row shows its destination alongside its title; narrow screens place the
destination below the title. The panel has space below the fixed header.

The [Mermaid theme gallery](docs/mermaid-theme/README.md) shows flowchart,
sequence, state, class, entity-relationship, and mindmap examples in the current
site theme. Each source builds light and dark SVGs for visual comparison.

`pkg/mermaid/color-schemes.json` saves the Cobalt drafting, Warm editorial,
Teal surfaces, and Violet technical color schemes in both modes. Each has
named canvas, surface, ink, outline, accent, and muted-border colors. This is
an inactive palette catalog; the render macro still selects the neutral
palette and diagram typography from `data/alwaldend/assets/theme.json`.
Mermaid geometry and color-role mappings remain in its preset templates,
independently of these saved colors.

## Local preview and validation

Run from the repository root:

```sh
bazel_agent bazel build //projects/alwaldend.com:site
bazel_agent bazel test //projects/alwaldend.com:site_test
bazel_agent bazel test //projects/alwaldend.com/test/mermaid:mermaid_test \
  //projects/alwaldend.com/test/mermaid:readability_test \
  //projects/alwaldend.com/test/search:search_test
bazel_agent bazel run //projects/alwaldend.com:site_serve
```

The preview serves the local build at http://127.0.0.1:1313.

Build the site's dependency diagram with
`bazel_agent bazel build //projects/alwaldend.com:deps_dark_webp` for a dark
WebP, or `:deps_webp` for light. The companion
`:deps_mermaid` and `:deps_svg` targets expose Mermaid and light/dark SVG
outputs. This overview analyzes `:site` at depth eight, reaching content
packages and repository documentation owners. It omits files, external
targets, tool dependencies, and visibility metadata; it is not the complete
transitive graph. See the
[analyzer documentation](../bazel_deps_analyzer/README.md) for options and
the build-cycle reason these artifacts remain outside the site's own inputs.

Use `:deps_mindmap`, `:deps_mindmap_svg`, or `:deps_mindmap_dark_webp` for an
ownership inventory of the same depth-eight targets, grouped by package.
Its branches describe ownership; use the flowchart to trace dependency edges.
The mindmap reuses the existing JSON capture and shared site appearance.

The Mermaid browser check writes screenshots and a JSON report to its Bazel
test outputs. It exercises the theme menu, saved and system preferences,
mobile layout, documentation print output, and blog header and body images.
The readability check measures the six gallery types and semantic edge cases
at article and mobile widths, retaining screenshots and contrast measurements.
The search browser check verifies loading, retry, stale-query handling, result
safety, owner navigation, the complete tag catalogue, and mobile navigation.
These browser checks need no external network.

## Search

The build exports Docsy's complete searchable document corpus, precomputes its
Lunr index, then renders the site with fingerprinted search assets. Preview,
release, and staging use the same content flags for the index and rendered pages;
draft posts remain absent from release and staging search results.

Search starts when its field receives focus. A Web Worker loads the index and
runs queries, keeping that work off the browser's UI thread. Results retain
Docsy's fields, ranking, excerpts, exclusions, and configured result limit.
Loading and download failures are visible, and Enter retries a failed request.
No search service or CDN is required. The build compiler and worker both use
Lunr 2.3.9 from the same lockfile; its license ships at `/licenses/lunr.txt`.

## Blog

The site publishes dated articles from `content/blog`. Each post is its own
Bazel package with an `index.md` and a `BUILD.bazel` declaring its
`docs_filegroup`; the section package aggregates them. The section publishes
HTML, an RSS feed, and a print edition, and appears in the main navigation.
The section index and its posts carry a `github_subdir` cascade so the
per-page GitHub links point at the content sources.

A post stays unpublished while its front matter declares `draft: true`. Local
builds render drafts for review and the release build excludes them, so the
draft state alone withholds a post from the deployed site.

A post selects its opening images with front matter `images: [pipeline.webp]`.
The shared blog template displays every entry in order above the article body,
including in the print edition. Hugo's Open Graph and Twitter card templates
use the same resources for social previews. Do not repeat a header image at
the start of the Markdown body; later illustrations stay in the body.

Every post is encouraged to have a representative image. Reuse an opening
image when available, preferring a raster projection for social cards. Posts
without a suitable image can omit this metadata. Bundled resource parameters
provide descriptive alternative text and an optional visible caption:

```yaml
images:
  - pipeline.webp
resources:
  - src: pipeline.webp
    params:
      alt: The uploader pipeline
```

Set `params.caption` when an image needs a visible caption. Without `params.alt`,
the image uses the page title as alternative text. Header images use the shared
image renderer and scale to the available width.

Header image files use a 5:2 canvas, matching the current
[X Article cover display](https://abs.twimg.com/x-web/x-web/assets/article-detail-B7xvnd3U.js).
This is a display ratio, not an API upload restriction. Keep the complete
illustration inside the canvas without cropping or stretching. Mermaid header
renders set `aspect_ratio = [5, 2]`; arrange their content for a wide canvas
before adding the minimum padding needed to reach that ratio. Other body
illustrations keep their natural dimensions.

Agent workflow: [Add a blog post](https://github.com/alwaldend/src/blob/master/projects/alwaldend.com/skills/alwaldend-blog/SKILL.md).

## Projects

The site publishes each registered project's visitor-facing landing page at
`/projects/<name>/`, with the section index at `/projects/`. Each project owns
its landing content in `projects/<name>/site/content/`; this site packages
those directories into `content/projects/<name>/` from the registry in
[projects/projects.bzl](../projects.bzl). Landing content carries no layouts,
styles, or build rules, and participates in the shared `statuses`,
`languages`, and `tags` taxonomies.

Repository reference documentation, including each project README, stays under
`/docs/projects/<name>/`. Because every landing is part of this one build, a
content error in any landing fails the whole site build and the blog
deployment.

## Documentation links

Markdown links and images resolve relative to their source directory.
`README.md` and `_index.md` links resolve to generated pages; packaged resources
use their published URLs, including when embedded in print pages. Link files
that exist only in the repository with explicit GitHub URLs. Unknown internal
destinations remain unchanged rather than being silently redirected to GitHub.

Print pages scope IDs and their fragment and control references to each source
document, keeping anchors distinct when documents are combined.

During site packaging, plain documentation that starts with an H1 receives a
derived Hugo title. Missing ancestor section pages preserve its directory
hierarchy, including OpenSpec specifications. Existing front matter and leaf
bundles remain authoritative. The title is rendered once and its original
heading anchor remains available; canonical Markdown needs no site metadata.
Tree disclosure controls support pointer clicks and keyboard toggling.
The tag index wraps its items, and `/404.html` offers links back to the site.

## Deployment

Build a portable website release bundle from the repository root:

```sh
tools/versioning/cmd/versioning/versioning.sh bazel -- \
  build --config=release //projects/alwaldend.com:release
```

The output directory `bazel-bin/projects/alwaldend.com/release` contains
`release.json` and `files/website.tar.gz`. The manifest uses `STABLE_VERSION`
from the versioning tool, including `0.0.0-dev` for development builds.
The archive contains the complete site with `index.html` at its root and
excludes drafts, generated project release pages/manifests, and project binary
payloads. Release artifacts remain on the download server. The release target
requires release configuration. Building
it does not connect to a host or publish anything.

With repository Vault authentication available, publish the bundle to XCP-ng:

```sh
bazel_agent bazel run //infra/download:publish.local -- \
  --release_dir "$PWD/bazel-bin/projects/alwaldend.com/release" \
  --site_archive website.tar.gz --steps upload,extract,link
```

After local verification, use `//infra/download:sync.yandex` to replace public
content from local. The wrappers inject short-lived SSH credentials and strict
host trust. Publication uses the release tool; synchronization runs the download
component's Ansible playbook. See the [release tool](../../tools/release/README.md)
for identity/configuration overrides and independently selecting upload,
extraction, or linking. The manifest version selects both the public archive
directory and the extracted site release; linking selects the main website.
Use a new version for changed local content: duplicate artifact or extracted
site versions fail without changes. Public sync overwrites matching versions.
Selecting a new site removes older extracted versions and retains every
downloadable archive. Rollback re-extracts and selects an earlier archive.
Development `0.0.0-dev` is suitable for the first local
deployment; subsequent releases should use distinct versions from the
[versioning workflow](../../tools/versioning/README.md).

The bundle is intentionally separate from the generated release documentation:
including the website archive in its own documentation would create a build
cycle. The Pages commands below remain available during the staged migration.

- DNS setup: [infra/dns](../../infra/dns)
- This project's [DNS declaration](https://github.com/alwaldend/src/blob/master/projects/alwaldend.com/dnsconfig.json)
  owns the shared `pages` address. Project landing pages are published by this
  site under `/projects/<name>/`, so no project owns a CNAME or a dedicated
  hostname. The [project directory](../README.md) links to every landing page.
  The apex and `www` records stay managed centrally.
- Deploy to the Github Pages repo (the `pages` branch of
  `alwaldend/alwaldend.github.io`, which GitHub Pages serves):
  ```sh
  tools/versioning/cmd/versioning/versioning.sh bazel -- \
    run --config=release //projects/alwaldend.com:deploy
  ```
  The deploy script clones the `pages` branch, replaces its contents with the
  built site, writes `.nojekyll`, and pushes only when the output changed.

## Update PVE VMs

```sh
bazel run //projects/alwaldend.com/tf # Apply tf
bazel run //projects/alwaldend.com/tf:update_pve_disk # Update the disk
```

## Taxonomy

| Taxonomy   | Meaning              |
| :--------- | :------------------- |
| Categories | General category     |
| Languages  | Programming language |
| Sites      | Sites                |
| Statuses   | Project status       |
| Tags       | Generic tags         |

## Shortcodes

### `alwaldend/alert`

```md
{{</* alwaldend/alert */>}}
Alert body
{{</* /alwaldend/alert */>}}
```

### `alwaldend/label_link`

Create a link using a bazel label

Usage:

```md
{{%/* alwaldend/label_link "//tools/qt" */%}}
```

### `alwaldend/links`

Render common links

Usage:

```md
{{%/* alwaldend/links */%}}
```

### `alwaldend/svg_file`

Render a packaged SVG using its published URL. Set `alt` to describe the image;
the page title is the fallback.

Usage:

```md
{{</* alwaldend/svg_file file="local_file.svg" alt="Project architecture" */>}}
```

### `alwaldend/include_html`

Include a local html file

Usage:

```md
{{</* alwaldend/include_html "file.html" */>}}
```

### `alwaldend/docs_misc_books`

Render books

Usage:

```md
{{</* alwaldend/docs_misc_books */>}}
```
