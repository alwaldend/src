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

## Deployment

- DNS setup: [infra/dns](../../infra/dns)
- Per-project DNS declaration: [dnsconfig.json](dnsconfig.json). It declares
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

Create `img` for an svg file

Usage:

```md
{{</* alwaldend/svg_file file=local_file.svg */>}}
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
