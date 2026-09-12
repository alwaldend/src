---
name: alwaldend-blog
description: >-
  Publish an article on alwaldend.com from supplied content, then review it for
  style and inconsistencies. Use when the user supplies blog content to publish
  or asks to add, draft, review, or prepare a post in content/blog; do not use
  to redesign the site or to write the article's substance.
---

# Publish a blog post

The user supplies the article content. Reproduce it without modification, keep
the post unpublished, then review it.

## Preserve the supplied content

- Treat the supplied text as the finished article. Copy it verbatim: keep
  wording, punctuation, capitalization, markdown, code blocks, links, and
  images exactly as given.
- Do not rewrite, reorder, summarize, extend, or "improve" the content, and do
  not fix what looks like a typo. Report defects instead of editing them.
- Carry over only what the post needs to build: front matter and the package
  declaration. Do not invent title, description, date, or author values.
- Ask the user for a missing title, date, or other required field only when the
  content itself does not supply it and the build cannot proceed without it.

## Create the post package

Each post is its own package so its bundle is independently owned. Create
`projects/alwaldend.com/content/blog/<slug>/`:

```sh
mkdir -p projects/alwaldend.com/content/blog/<slug>
```

Write the supplied content to `index.md` with front matter:

```markdown
---
title: <title>
linkTitle: <short title>
date: <YYYY-MM-DD>
description: <one-line description>
author: <author>
draft: true
---
```

- `title` is required. `linkTitle` is the shorter navigation and list label.
- `date` orders the post. The section lists posts grouped by year, so a wrong
  or missing date misplaces the article.
- `description` renders as the lead paragraph and appears in listing excerpts.
- `author` renders in the post byline. Omit it only when the user supplied no
  author.
- Keep `draft: true` from the template below. New posts are unpublished.

Put any images or other resources the post references beside `index.md` in the
same directory, and reference them by bare relative path, such as
`![diagram](diagram.png)`. Do not place them in `static/`.

Declare the package in `BUILD.bazel`:

```starlark
load("@rules_docs//docs:defs.bzl", "docs_filegroup")

docs_filegroup(
    name = "docs",
    srcs = glob(
        ["*"],
        exclude = [
            "BUILD",
            "BUILD.bazel",
        ],
    ),
    prefix = "content/blog/<slug>",
    visibility = ["//projects/alwaldend.com/content/blog:__pkg__"],
)
```

The section aggregates every post package automatically through the
`subpackages` lookup in `content/blog/BUILD.bazel`. Do not edit that file or
`content/BUILD.bazel` to register a post.

## Keep the post unpublished

New posts stay out of the published site until the user asks to release them.
Keep `draft: true` in the front matter and stop there; do not set a future
`date` as a substitute.

The site renders drafts for local preview and excludes them from the release
build, so a draft is withheld from the deployed site while remaining visible
at `bazel_agent bazel run //projects/alwaldend.com:site_serve`. Publishing is a
separate authorized step: report the post's draft state in the handoff and let
the user decide when to release it. Remove `draft: true` only when the user
asks to publish.

## Review the post

Review after the package builds. Do not change the content: report findings.

Number the review so the user can reference a finding by number. Group the
findings into numbered sections, and render every finding as a Markdown list
item labelled with its section number and position, such as `- 3.4`. Keep each
item to a single finding with its location and reason, use one label per
finding, and never reuse a label within a review. Do not present findings as
prose paragraphs or unnumbered bullets; the numbering must stay readable as a
list.

- Build and inspect the rendered page: `bazel_agent bazel build
//projects/alwaldend.com:site`, then read the generated post under
  `bazel-bin/projects/alwaldend.com/site.dest/blog/<slug>/`.
- Check inconsistencies against the repository's other pages: heading case and
  depth, list and code-fence style, link text, terminology, and any spelling
  that conflicts with the repository's established usage.
- Check the front matter matches the content, including title, description,
  date, and author.
- Check that internal links and images resolve. Markdown targets resolve
  relative to the source directory; `README.md` and `_index.md` map to published
  pages. The site preserves an unknown destination instead of failing the build,
  so verify each one by hand and report the ones that point at nothing a reader
  can reach. Files that exist only in the repository need explicit GitHub URLs.
- Confirm the site builds clean by running
  `bazel_agent bazel test //projects/alwaldend.com:site_test` and report any
  failure.

Report every finding to the user with the exact location and the reason, and
let the user decide whether to change the content. Never apply a content edit
that the user did not ask for.

## Deliver

Ordinary delivery follows the `repo-delivery` skill for verification, commit,
push, and pull request. Publishing the post itself remains a separate step that
the user authorizes.

Local preview is available with `bazel_agent bazel run
//projects/alwaldend.com:site_serve` at `http://127.0.0.1:1313`. Use the
`repo-hugo` skill for the site's Hugo toolchain, theme, and deployment
workflow.
