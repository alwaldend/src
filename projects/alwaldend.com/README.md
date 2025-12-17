---
title: Alwaldend.com
description: Alwaldend.com site
sites:
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
- Hugo rules: [../../tools/hugo](../../tools/hugo)

## Features

- [Hugo](https://gohugo.io) site
- [Docsy](https://github.com/google/docsy), [Bootstrap](https://getbootstrap.com)

## Deployment

- Setup dns: [infra/alwaldend.com/dns](../../infra/alwaldend.com/dns)
- Push to the Github Pages repo:
  ```sh
  bazel run --config=release //projects/alwaldend.com:deploy
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

## TODO

- Setup proper deployment
