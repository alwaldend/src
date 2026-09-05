---
title: Alwaldend.com
description: Alwaldend.com site
websites:
  - alwaldend.com
statuses:
  - in_progress
tags:
  - hugo
  - github_pages
---

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/alwaldend.com
- Github Pages repo: https://github.com/alwaldend/alwaldend.github.io
- Hugo rules: [../../projects/rules_hugo](../../projects/rules_hugo)

## Features

- [Hugo](https://gohugo.io) site
- [Docsy](https://github.com/google/docsy), [Bootstrap](https://getbootstrap.com)

## Local preview and validation

Run from the repository root:

```sh
bazel_agent bazel build //projects/alwaldend.com:site
bazel_agent bazel test //projects/alwaldend.com:site_test \
  //projects/alwaldend.com/test/site:site_test
bazel_agent bazel run //projects/alwaldend.com:site_serve
```

The preview serves the local build at http://127.0.0.1:1313. The generated-site
test checks internal links, image paths, fragment targets, duplicate IDs, and
selected page markup. It does not check external destination availability.

## Documentation links

Markdown links and images resolve relative to their source directory.
`README.md` and `_index.md` links resolve to generated pages; packaged resources
use their published URLs, including when embedded in print pages. Link files
that exist only in the repository with explicit GitHub URLs. Unknown internal
destinations remain unchanged and are reported by the generated-site test;
they are not silently redirected to GitHub.

Print pages scope IDs and their fragment and control references to each source
document, keeping anchors distinct when documents are combined.

## Deployment

- DNS setup: [infra/dns](../../infra/dns)
- Per-project DNS declaration: [dnsconfig.json](https://github.com/alwaldend/src/blob/master/projects/alwaldend.com/dnsconfig.json). It declares
  the GitHub Pages `pages` address and one unproxied CNAME per project landing
  subdomain, for example `android-launcher.alwaldend.com`. Hostnames use
  hyphens while docs paths keep underscores. Each landing site links to its
  project documentation, so Cloudflare redirects are unnecessary. The apex
  and `www` records stay managed centrally.
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
