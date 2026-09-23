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
uses Hugo with Docsy to publish project documentation on GitHub Pages.

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/alwaldend.com
- Github Pages repo: https://github.com/alwaldend/alwaldend.github.io
- Hugo rules: [../../tools/rules_hugo](../../tools/rules_hugo)

## Features

- [Hugo](https://gohugo.io) site
- [Docsy](https://github.com/google/docsy), [Bootstrap](https://getbootstrap.com)

The homepage, blog, documentation, and project pages share a flat monochrome
theme with white and charcoal color modes, subtle borders, and system fonts.
The navigation's theme menu supports light, dark, and system preferences.
Every page uses an "Al" SVG favicon outlined in Architects Daughter.
Prose links are underlined, and keyboard focus uses a visible outline.
Scrollbars are thin and muted; taxonomy links use small rounded labels shared
with the project index. Documentation sidebars use plain action rows,
uppercase group headings, compact wrapping taxonomy items, subtle hover and active
backgrounds spanning the full tree row, and compact chevrons. Tree labels stay
indented; chevrons expand branches independently of the row's page link. Headings,
search fields, code blocks, tables, and callouts follow the same typography
and surface styles. Syntax highlighting
uses CSS classes so code colors follow the active theme.
The homepage panel links to GitHub, GitLab, Blog, Docs, and the project index,
followed by permanently visible, indented project links in title order.
Each row shows its destination alongside its title; narrow screens place the
destination below the title. The panel has space below the fixed header.

## Local preview and validation

Run from the repository root:

```sh
bazel_agent bazel build //projects/alwaldend.com:site
bazel_agent bazel test //projects/alwaldend.com:site_test
bazel_agent bazel run //projects/alwaldend.com:site_serve
```

The preview serves the local build at http://127.0.0.1:1313.

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
