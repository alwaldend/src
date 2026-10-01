---
title: X Article uploader
description: Compile a Markdown blog post into a DraftJS content_state and create an X Article draft
statuses:
  - active
languages:
  - go
tags:
  - x
  - article
  - markdown
---

This project compiles a blog post into an X Article draft. It has two
operations with different inputs, outputs, and trust requirements:

- **Conversion** compiles a post's Markdown body into the DraftJS
  `content_state` document the X Articles draft endpoint accepts. It is
  offline, deterministic, and credential-free; it performs no network access.
- **Draft creation** resolves the body images and article banner through the
  media upload endpoint, then creates an X Article draft. It is the only
  operation that performs network calls, and the only one that needs
  credentials. The uploader creates drafts only and never publishes.

Splitting the two lets a draft be produced and diffed offline before anything
reaches X, and keeping the tool to draft creation means no invocation of it can
make a post publicly visible.

## What conversion does

A source post is YAML front matter plus a Markdown body. The front matter is
metadata: its `title` is returned alongside the document because the draft
endpoint requires a title that `content_state` cannot carry, and it never
appears in the article text. A source without a non-empty title fails
conversion.

The first entry in front matter `images`, when present, supplies the article
banner. It must name a supported raster image relative to the post package;
remote URLs, missing files, and references outside the package fail offline
conversion. Later entries do not become additional banners. Conversion records
the selected image in an optional `banner_locator` with its path, post package,
media type, and content digest. It does not add a body image for the banner:
an existing Markdown image remains in the body, including when it references
the same file. Posts and older artifacts without a banner remain supported.

Markdown maps onto the vocabulary the API exposes:

- Headings clamp: `#`/`##` to `header-one`, `###` to `header-two`, and `####`
  and deeper to `header-three`, because X exposes three heading levels.
- Paragraphs, both list kinds, nested list items, and block quotes map onto
  their matching block types.
- Bold, italic, and strikethrough become `inline_style_ranges`; links become
  `link` entities. Offsets are UTF-16 code units of the block's final text,
  matching how the document is interpreted as a JavaScript string. A link or
  image destination resolves the Markdown escapes and character references the
  source may carry, so the recorded URL or path is the one the author wrote.
- Fenced code and tables become `atomic` blocks backed by a `markdown` entity
  that preserves the original Markdown, because X has no table block or entity.
- Footnotes become a trailing `Footnotes` heading with each definition as an
  ordered list item.
- An inline-code span keeps its text; because DraftJS has no code style, the
  lost styling is reported as a continuing diagnostic. A thematic break becomes
  an `atomic` block backed by a `divider` entity.
- An image wrapped in a link, such as `[![alt](image.png)](target)`, is still
  emitted as an image; because the image carries no text span, the link cannot
  become a `link` entity over it, so the lost destination is reported as a
  continuing diagnostic naming the target. When the link also carries text,
  that text keeps the destination and the image's lost link is still reported,
  so the outcome does not depend on whether unrelated sibling text happens to
  exist. A link or image title has no place in the entity `data` object, so its
  loss is reported the same way.
- Each image becomes an `atomic` block with an unresolved `image` entity. Its
  source travels beside the document as a locator, because the entity's `data`
  object rejects additional properties. An image whose media type the upload
  endpoints reject is reported instead of being emitted as an entity publication
  could never resolve.

### Accepted image media types

This project owns the set of image media types the X media upload endpoints
accept, because it is the component that talks to them. The accepted set is
`image/jpeg`, `image/gif`, `image/bmp`, `image/png`, `image/webp`,
`image/pjpeg`, and `image/tiff`; `image/svg+xml` is not among them. A post
intended for publication therefore must not reference an SVG. Authoring
guidance for blog posts and Mermaid diagrams points at this list rather than
repeating it, so a change in X's support is made in the uploader's
`internal/markdown/image.go` and read from here.

Conversion distinguishes two outcomes and never silently drops content: a
construct whose content survives with lost formatting is reported and
conversion continues, while one whose content would be lost fails conversion.
Both name the source position that produced them.

## What draft creation does

Draft creation reads the draft artifact conversion emitted. It obtains each
body or banner image's bytes from the locator the artifact records — resolved
against the recorded post package, refused when the path escapes the post directory or the
bytes' digest no longer matches — and uploads it. Containment is enforced by
resolving the name through descriptor-anchored directory handles that reject a
symlink escaping the root, so the checked object is the object read and a
swapped symlink cannot redirect it. For body images, the returned `media_id` lands
on the entity the locator names, not by position. Identical bytes reuse one
identifier within a run and across a later run through a digest-to-`media_id`
cache under the task's ignored output directory; each entry records the lifetime
the upload endpoint reported, so an identifier that X has since dropped is
uploaded again rather than reused. The cache holds identifiers, digests, and
expiries only, never credential material. Banner and body images use this same
cache, so identical bytes need only one upload even when both reference them.
All referenced files, content digests, and media types are verified before the
first upload; an invalid banner cannot consume an upload or draft request.
Recording the cache is best-effort: a cache that cannot be written costs a later
re-upload rather than failing draft
creation after the uploads already succeeded.

The media uploads happen separately before draft creation. The selected banner's
uploaded identifier becomes the draft request's top-level `cover_media`, with
`media_category: "tweet_image"` and `media_id`. Body image identifiers remain
in their DraftJS entities. This follows the
[X draft endpoint's request schema](https://docs.x.com/x-api/articles/create-draft-article).

The uploader creates a draft and reports its article identifier; it never calls
the publish endpoint, so a review in the X composer is always the step that makes
an article public. It does not retry automatically when a request is rejected.

## Credentials

Draft creation authenticates with OAuth 1.0a user context, whose token does not
expire on the fixed schedule an OAuth 2.0 user token does. It requires four
fields, supplied at run time through the repository's Vault-backed injection
flow:

- `X_API_KEY`
- `X_API_SECRET`
- `X_ACCESS_TOKEN`
- `X_ACCESS_TOKEN_SECRET`

No credential value is recorded in source, and the values are never printed.
The identity, policy, and credential path are owned by
`infra/vault/openspec/changes/add-x-article-uploader-credentials`. The credential
sits inside the identity's own AppRole subtree, the same shape every other
component uses, and `al.lua` names its path directly, as every other component's
injected reference does; the shared AppRole module's own-subtree policy grants
the read.

## Run

Convert a post to a draft artifact:

```sh
bazel_agent bazel run //projects/x_article_uploader/cmd/convert -- \
  --source projects/alwaldend.com/content/blog/<slug>/index.md \
  --out out/x_article_uploader/<slug>.json \
  --post-package projects/alwaldend.com/content/blog/<slug>
```

`--post-package` may be omitted, in which case the package is derived from the
post directory's path relative to `--workspace` (default `.`). Whether supplied
or derived, the package is opened through a descriptor-anchored handle on the
workspace and the post's source and images are read through that handle, so the recorded
package and the bytes' source are the same object: a symlinked post directory
that escapes the workspace, or a supplied package whose opened directory is not
the post directory, is refused instead of recorded; the package is accepted by
comparing the identity of the directory actually opened, not a pathname. The workspace is the caller-supplied
root, never a root rediscovered by walking the filesystem.

Draft creation is the only operation that needs credentials, so run it through
the wrapper that injects them: `//projects/x_article_uploader:draft` selects
the `x_api=1` injection declared in `al.lua`, authenticates as the component
AppRole, and places each credential field in its environment variable. The
wrapped command is the draft binary, so its own flags are passed after
`--`; the binary and the injector travel in the target's data.

Create a draft from a converted artifact:

```sh
bazel_agent bazel run //projects/x_article_uploader:draft -- \
  --artifact out/x_article_uploader/<slug>.json
```

Publishing is deliberately unavailable: there is no publish flag and no publish
endpoint call, so a draft becomes public only through a separate review step
outside this tool. The injector only supplies credentials; it does not select
the artifact or publish anything. Do not export `X_*` values by hand, which
would bypass the AppRole boundary the wrapper establishes.

## Checks

```sh
bazel_agent bazel test //projects/x_article_uploader/test/e2e:e2e_test
```

The end-to-end check discovers every post in the site content tree, asserts each
post's recorded outcome, emits a repeatable draft artifact, and exercises
publication against a recorded-response client. It needs neither credentials
nor network access.
