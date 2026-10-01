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
media type, and content digest. It does not add a body image for the banner.

Every explicit Markdown image remains in the article body, including opening
images and repetitions of the selected banner. Front matter `images` selects
the banner independently: it neither inserts nor removes body images. Identical
banner and body bytes still share one media upload. Posts and older artifacts
without a banner remain supported.

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
- Fenced code and tables become `atomic` blocks backed by a mutable `markdown` entity
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
- Each body image becomes an `atomic` block with an unresolved `image` entity. Its
  source travels beside the document as a locator, because the entity's `data`
  object rejects additional properties. An image whose media type the upload
  endpoints reject is reported instead of being emitted as an entity publication
  could never resolve.

### Accepted image media types

The uploader supports static `image/jpeg`, `image/png`, `image/gif`, and
`image/webp` files no larger than 5,000,000 bytes, matching the documented
[simple-image upload formats and limit](https://docs.x.com/x-api/media/quickstart/best-practices).
Animated GIF, APNG, and WebP are refused explicitly; this tool does not implement
the separate animated-media upload contract. BMP, TIFF, PJPEG, and SVG are not
supported by this uploader. The initialize endpoint's broader MIME enumeration
does not establish support in the simple-upload path used here.

Conversion and draft preflight apply the same checks to banners and body images.
JPEG, PNG, and GIF are decoded; WebP receives RIFF container, chunk, frame-header,
and dimension validation, not full pixel decoding. A separate local limit of
32,000,000 pixels bounds decoder memory use; this is a tool resource guard,
not a claimed X limit. File extensions must match the detected format.
Authoring guidance points at this section rather than repeating the contract;
the implementation is owned by `internal/markdown/image.go`.

### Markdown payload budget

[X's schema](https://docs.x.com/x-api/articles/create-draft-article) documents a
10,000 weighted-length limit for Markdown entities per article,
but does not publish an Articles-specific weighting algorithm. Conversion uses
a conservative estimate with a local budget of 9,500: NFC-normalized text for
counting, the [published X codepoint weights](https://github.com/twitter/twitter-text/blob/master/config/v3.json),
separate weights for emoji components, and at least 23 units for each plausible
dotted domain fragment. Other characters are counted normally. The original
Markdown payload is preserved unchanged. This can reject some content X would
accept and is not a guarantee of equivalence to its backend validator. An
over-budget diagnostic reports the estimate and the local budget.

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
identifier within a resolution and, when its lifetime is known, across later
runs through a versioned cache under the task's ignored output directory.
Persistent reuse is scoped to a one-way fingerprint of the API endpoint, media
category, and all four OAuth credential fields. Changing accounts or rotating
credentials therefore causes a cache miss without an extra identity request.
Unscoped legacy entries are not trusted. The cache contains identifiers,
digests, expiries, and opaque scope fingerprints, never credential values.

Known lifetimes are anchored before the upload request, so upload and processing
time cannot extend them, and reserve 30 seconds for submission. Media that
expires while other files upload is refused before draft creation, including a
later duplicate of an earlier image; it is not uploaded again mid-resolution. Missing or
zero lifetime permits deduplication only within the current resolution; negative
or already unusable lifetimes fail. Banner and body images share the cache, so
identical bytes need only one upload within that resolution.
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

An upload without `processing_info`, or with a `succeeded` state, needs no
status request. Reported `pending` or `in_progress` processing triggers signed
read-only status checks after the server's `check_after_secs`, with a minimum
one-second delay and a two-minute default deadline. Failed, unknown, missing,
or inconsistent processing state, partial upload errors, and rejected status
requests stop the operation before the draft POST. Status polling never replays
an upload, and no rejected request is retried automatically.

The uploader creates a draft and reports its article identifier; it never calls
the publish endpoint, so a review in the X composer is always the step that makes
an article public. It does not retry automatically when a request is rejected.
HTTP redirects are returned as failures without following them, preventing a
redirect from issuing another media upload or draft request. The HTTP transport
cannot recreate a consumed POST body for an automatic replay.

### Failed requests and retry timing

Every non-2xx API error, including HTTP 503, retains and prints the response's
`X-Rate-Limit-{Limit,Remaining,Reset}`,
`X-User-Limit-24hour-{Limit,Remaining,Reset}`, and `Retry-After` headers.
Only those headers are copied; authorization, cookies, and other response
headers are excluded. The original operation, status, and body remain visible,
including the available body and underlying cause if reading the response fails.

Valid Unix reset values also show their UTC date and time. A
`Rate-limit retry not before` line appears only when known exhausted windows
(`Remaining: 0`) have unambiguous future resets. It uses the latest exhausted
window's reset and extends it with valid later `Retry-After` advice, either an
HTTP date or seconds after response receipt. Missing, invalid, ambiguous, or
elapsed reset evidence leaves the current retry time unknown. Positive remaining
counts and standalone `Retry-After` advice do not establish upload eligibility.
The displayed boundary does not guarantee service recovery or request success,
and it never triggers an automatic wait or retry.

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
