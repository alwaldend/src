## Context

See `proposal.md` for motivation and `specs/` for the behavior contract.

The target API is X API v2.168. Its Articles surface is two endpoints,
`POST /2/articles/draft` and `POST /2/articles/{article_id}/publish`; this
project uses only the draft endpoint. The draft body is `title` plus a
`content_state` DraftJS document (`blocks`, `entities`), with an optional
`cover_media`. The spec at `api.x.com/2/openapi.json` fixes the vocabulary this
project must target, and the draft shape is not canonical DraftJS: blocks name
their ranges `inline_style_ranges` and `entity_ranges`, carry no `depth`, and
`entities` is a mapping of `{key, value}` entries, because the schema sets
`additionalProperties: false`:

- block `type`: `unstyled`, `header-one`, `header-two`, `header-three`,
  `unordered-list-item`, `ordered-list-item`, `blockquote`, `atomic`;
- inline style: `bold`, `italic`, `strikethrough`;
- entity `type`: `post`, `link`, `image`, `emoji`, `markdown`, `divider`,
  `latex`, each with `mutability` `immutable`, `mutable`, or `segmented`;
- `data.markdown` is documented for `markdown` entities with a 10,000
  weighted-length budget per article;
- the entity `data` object accepts `caption`, `entity_key`, `markdown`,
  `media_items`, `post_id`, and `url`, with `additionalProperties: false`;
  `url` is documented for `link` and `media_items` for `image`, and
  `media_items` entries require `media_category` and `media_id`;
- the media upload endpoints accept a fixed `media_type` set whose image
  members are `image/jpeg`, `image/gif`, `image/bmp`, `image/png`,
  `image/webp`, `image/pjpeg`, and `image/tiff` — `image/svg+xml` is not among
  them.

There is no Markdown import endpoint: `content_state` is the only body the API
accepts, so a converter is required. Entity and inline ranges are character
offsets into the enclosing block's final `text`, which is the constraint that
makes the source document's parse tree — not its rendered HTML — the right input.

The blog posts live at `projects/alwaldend.com/content/blog/<slug>/index.md`.

## Goals / Non-Goals

**Goals:**

- Compile a blog post's Markdown into a valid `content_state` document offline.
- Keep conversion deterministic and testable without credentials or network.
- Resolve and upload referenced images, then create an X Article draft.
- Give a new post a raster form for its images, including a Mermaid diagram,
  because the upload endpoints accept no SVG.
- Build under the repository's pinned, hermetic Bazel workflow.

**Non-Goals:**

- Rendering the article inside Hugo, or adding any site output format.
- Authoring or editing posts, drafts, or titles on X outside this pipeline.
- Publishing an article: the uploader creates drafts only, so publication and
  any tracking or retraction of a published article are out of scope.
- Supporting `post`, `emoji`, and `latex` entities: the posts do not use them,
  and they carry composer-specific identifiers this project cannot synthesize.

## Decisions

### Parse Markdown rather than scrape rendered HTML

Convert from the Markdown source using `github.com/yuin/goldmark`, a
CommonMark/GFM parser that exposes a node tree with inline text and source
positions, instead of parsing `.Content` HTML inside Hugo templates. The user
approved this dependency on 2026-09-29, choosing it by name from the parser
options presented for this project.

The proposal named the exact package and its maintenance and supply-chain
costs before that approval: it is a widely used pure-Go parser with no cgo and
no runtime services, pinned to `v1.8.6` in the root Go manifest and lockfile and
exposed through the owning `use_repo` declaration, so the import graph stays
hermetic and reproducible. Its maintenance surface is a single module
(`github.com/yuin/goldmark`) plus its own extension packages; upgrading it is an
ordinary dependency-pin bump through the owning dependency workflow, and it
carries no new transitive dependency needing review.

Alternatives considered for the parser: `rsc.io/markdown`, already resolvable in
the repository as `@tools_io_rsc_markdown//:go_default_library`. It is rejected
because it cannot parse the constructs these posts use: its `Table` option is
marked `TODO` in its options struct and it has no footnote support at all, while
`slop-without-a-clear-goal/index.md` contains an 11-row GFM table and a
footnote. It is also declared inside the isolated `go_deps` extension of
`tools/pkgsite`, so consuming it from a project crosses a tool boundary. Using
the host Go toolchain to invoke a converter was rejected too: the repository
pins Go through `//tools/go:go`, and conversion has to build in the hermetic
Bazel graph.

Alternatives considered for rendering inside Hugo: render hooks plus template
string manipulation. That was measured against Hugo 0.165.0 and rejected. Only seven hooks exist
(`blockquote`, `codeblock`, `heading`, `image`, `link`, `passthrough`,
`table`); there is no paragraph or list-item hook, and no hook for `strong` or
`em`. Paragraphs and list items are the bulk of an article, so they could only
be produced by regexing rendered HTML. That approach also lacks what range
computation needs: the rendered text is already entity-escaped, table cells
arrive as rendered HTML, RE2 rejects the backreferences that paired-tag matching
would want, and nested lists break naive item extraction. Template hooks do
expose `.Position` with an `Offset`, `LineNumber`, and `ColumnNumber`, so hooks
are a viable input for the constructs they cover — but they cannot cover the
document, and a converter that takes Markdown directly has one input rather than
two.

### Parse front matter as metadata, not as body

A blog post is YAML front matter plus a Markdown body, and the draft endpoint
requires a `title` that `content_state` cannot carry. The input model is
therefore a parsed post: front matter is extracted as metadata, its `title` is
returned alongside the document for draft creation, and the body alone is
converted. Converting the file as one Markdown document would emit the front
matter as article text and leave draft creation without its required title.

### Raster-only authoring policy for posts from now on

The media endpoints accept a fixed set of image media types that excludes SVG,
and all six image references in the current posts are SVG, so those images
cannot be uploaded as-is. The user's decision is that this applies to new posts
only: a post published from now on references raster images, and a Mermaid
diagram is rendered to a raster image instead of being referenced as an SVG. The
current SVG diagrams in `diagrams-in-ac-era` are not retrofitted.

Two things carry the policy. The blog authoring guidance gains a stated
requirement that a post's images use a media type the upload endpoints accept,
and `tools/mermaid` gains a `mermaid_webp` rule so a diagram has a hermetic way
to become a raster image beside `index.md`. A diagram is the only source that
cannot be raster on its own: a hand-authored image is chosen by the post's
author, while a `.mmd` source needs a render step.

`mermaid_webp` follows the existing render contract rather than inventing a
second one. It reuses the pinned Chrome, font, theme, and paint-order inputs, so
a WebP render is as hermetic as the SVG render and carries the same appearance.
The output is a checked-in blog asset like every other post image, copied into
the post package with `write_source_files` so the `.mmd` stays authoritative and
the freshness test guards the copy. WebP is one of the accepted media types and
is already an LFS-tracked extension, so the rendered asset enters the repository
through the existing binary-asset path.

The converter still reports an image whose media type the upload endpoints
reject, with its source position, rather than emitting an entity that
draft creation can never resolve. That report is now a guardrail for a policy
violation in a source post rather than the expected outcome of every diagram
post; a post that follows the policy converts without such a report.

### Per-post expected outcome in the end-to-end check

The end-to-end check converts every post in the content tree, but a post's being
in the tree does not make it convertible. `diagrams-in-ac-era` references six
SVGs, which the upload endpoints cannot accept, so asserting that every post
maps cleanly would contradict the converter's own rule that an unacceptable
image is reported rather than silently emitted. The check therefore carries an
expected outcome per post: a policy-conforming post must map to the documented
block, range, and entity contract, and a post with an unacceptable image must
fail with the named diagnostic at the reported source position.

The outcome is keyed by post so the check still proves the converter over real
content, and the set of posts is discovered from the content tree rather than
listed by hand, so a post added later is covered without editing the check. A
new post is expected to conform by default; only a post that fails the
raster-only policy carries a failing expectation.

### Separate the library, the converter command, and the draft command

The parse-and-map logic is a library under `internal/`; conversion and draft
creation are separate entry points under `cmd/`. Conversion writes a
`content_state` document and performs no network access; draft creation consumes
that document and is the only network path. It calls only the draft endpoint and
exposes no publish operation, so no invocation can make an article public.

This keeps the hermetic, credential-free half independently testable, and lets
a draft be produced and diffed before anything reaches X. The draft command also
takes a pre-built document file so a draft can be re-run against a reviewed
document without re-converting.

Uploads need one piece of state that neither half may rewrite: the mapping from
an uploaded image's content digest to the `media_id` the API returned. A digest
alone cannot recover it, and the reviewed artifact must stay unchanged, so the
draft command owns a separate cache file under the task's ignored output directory
that records the mapping and is read on the next run. The cache holds identifiers
and digests only, never credential material; a missing or unreadable cache costs a
re-upload rather than failing draft creation, and a cache entry whose recorded
digest no longer matches the artifact's is ignored. Because that directory is
disposable, the reuse guarantee is scoped to a task workspace rather than to any
host or any later machine.

### Map headings by clamping, not by dropping

`#` and `##` become `header-one`, `###` becomes `header-two`, and `####` and
deeper become `header-three`, because X exposes exactly three heading levels and
a dropped heading loses structure. Every current post opens at `##`, so their
section headings become `header-one`.

### Preserve code and tables as `markdown` entities

Code fences and GFM tables become `atomic` blocks backed by a `markdown` entity
whose payload is the construct's original Markdown source. X has no table block
or entity, and the `markdown` entity is documented with a fenced-code example
plus a weighted budget, so preserving source Markdown is the intended use.
Because the budget is host-enforced and not fully specified, the converter
measures the total payload and fails before sending an over-budget document
rather than letting the API reject it.

### Inline code loses its styling; thematic breaks become dividers

DraftJS's inline style vocabulary is exactly `bold`, `italic`, and
`strikethrough`, so a Markdown inline-code span has no code style to map onto.
Keeping the span's literal text and reporting the lost monospace styling as a
continuing diagnostic is the honest mapping: dropping the span would lose
content, and inventing a style the API does not accept would be rejected.
`dns-management-in-a-monorepo/index.md` contains inline-code spans, so the path
is exercised by a real post.

A Markdown thematic break maps onto the `divider` entity, which X exposes for
exactly that construct, as an `atomic` block. Leaving it as paragraph text would
render `---` literally in the article, and dropping it would silently remove the
author's section separation. `slop-without-a-clear-goal/index.md` contains one.

The conversion-diagnostics contract distinguishes these two outcomes: a
construct whose content would be lost fails conversion, while one whose content
survives with lost formatting is reported and conversion continues. Both are
reported with the source position that produced them, so an author sees every
loss rather than only the fatal one.

### Emit footnotes as a trailing section

DraftJS has no footnote construct and no superscript inline style. Each in-text
`[^label]` reference becomes plain bracketed text, and the definitions become a
trailing `Footnotes` heading with each definition as an ordered list item,
rather than being dropped. `slop-without-a-clear-goal/index.md` uses one
footnote, so this path is exercised by a real post.

### Images: convert to a skeleton, resolve at draft creation

The converter emits an `atomic` block plus an `image` entity carrying the alt
text as `caption` and leaving the media unresolved. The draft command uploads
the image bytes through the media endpoint, sets `media_items` to the returned
`media_id` with category `tweet_image`, and only then creates the draft.

`media_items` requires a `media_id`, and `url` is documented for `link`, so a
document URL cannot source an image and an upload is unavoidable for an image
that is to render as one. Uploading at conversion time would break conversion's
offline and deterministic contract, so the resolution step belongs to draft
creation. The document format therefore distinguishes an unresolved image from
a resolved one, and draft creation refuses to send an unresolved image. Uploads
are keyed by content hash so re-running draft creation does not re-upload
identical bytes.

The API entity cannot carry the image's source, because its `data` object
rejects additional properties and is not the place for a local path. The
artifact therefore carries the locator as its own metadata, outside
`content_state`: each unresolved image gets an entry naming the image's source
path relative to the post's directory, the post package that base resolves
against, and a digest of the bytes as they were at conversion time. Recording
the base lets the artifact travel — the draft command resolves the path against the
recorded package rather than against wherever the artifact file happens to sit —
and recording the digest gives the draft command the value it tests the current
bytes against. A locator also records the `image` entity it resolves, by the
entity's `entity_key`, so the returned `media_id` is attached to the entity the
locator names: an artifact can carry several images, captions and
byte digests can repeat, and position or caption alone is not a stable
identifier. A locator that escapes the post's directory, or an image whose
digest no longer matches, is a failure rather than a silent substitution.

Containment is a property of the object actually read, not of the pathname that
was checked: a check that canonicalizes a path and then re-opens it can be
defeated by a writer that swaps the canonical target for a symlink in between,
so the bytes then come from outside the validated package even though the digest
matches. Both the converter's read and the draft command's read therefore resolve
the name through `os.Root`, which walks each component through descriptors and
refuses a symlink that escapes the directory it is anchored at. The check and
the read then operate on the same descriptor, so a swapped symlink cannot
redirect the read.

All six image references in the current posts are `.svg`
(`diagrams-in-ac-era/index.md`), and SVG is not among the media types the upload
endpoints accept, so these images cannot be uploaded as-is. Since only new posts
need a draft, the existing post is left as it is and its images are reported
with their source positions when it is converted. A new post follows the
raster-only policy in the decision above, so its images are uploadable without a
special case.

### OAuth 1.0a user context, signed with the standard library

The Articles and media endpoints accept either an OAuth 2.0 user token or an
OAuth 1.0a user context. The user's decision is OAuth 1.0a, because an OAuth 2.0
user token expires in about two hours and its refresh token is single-use and
rotating: a missed rotation invalidates the token family, which turns a
scheduled draft creation into a coordination problem with a credential store the
component may only read. An OAuth 1.0a user-context token does not expire on a
fixed schedule; it stays valid until it is revoked or the account is suspended,
so a stored credential keeps working without a rotation step. An app-only
bearer token was rejected as well: it carries no user context, so the Articles
and media endpoints reject it.

OAuth 1.0a needs four fields — API key, API secret, access token, and access
token secret — so the injected reference is a four-variable tuple rather than
the single `X_API_TOKEN` the first sketch assumed; the environment variable
names follow that sketch's prefix. Request signing is HMAC-SHA1 over a
normalized parameter string with the RFC 5849 percent-encoding rules, which the
standard library provides (`crypto/hmac`, `crypto/sha1`, `encoding/base64`,
`net/url`), so no signing dependency is added. The nonce and timestamp are the
only non-deterministic inputs; they are generated per request, while conversion
stays deterministic because signing belongs to draft creation alone.

### Credentials through the repository injection flow

No credential value is recorded in source. The draft command reads the four
credential fields at run time through the same Vault-backed injection flow the
other components use. Publication is the only operation that needs them.

Application-side loading is only half of that path, so this change also plans
the repository half. A run-time read needs a component identity to authenticate
it, a Vault path holding the credential, and a policy granting that identity
access to its own path, none of which exist for a new project. The uploader
therefore gets an `al.lua` and its `al_config` target that name the injected
reference and select the injector plugin, and the identity, policy, and
reference are carried by their owner in
`infra/vault/openspec/changes/add-x-article-uploader-credentials`. The credential is held at
`alwaldend.com/vault1/approles/src_projects_x_article_uploader/oauth1` -- its
path is one checked-in value that both the identity's checks and the injected
reference read, so the two cannot drift -- inside the identity's own AppRole
subtree, the same shape every other component uses, so the shared AppRole
module's own-subtree policy grants the read. The plan adds the wiring but does
not write a credential, create a role, or exercise a live read: those remain
separately authorized operations.

## Risks / Trade-offs

[The parser must expose the CommonMark and GFM semantics the posts rely on,
including tables and footnotes; a weaker parser would silently change the
delimiter rules] -> Require extensions for tables, footnotes, and strikethrough
explicitly, and assert those constructs in tests rather than assuming a default.

[The converter's parser version could diverge from the site's Markdown parser
version, so the same source might mean slightly different things] -> These are
independent concerns: the site renders HTML, this project produces DraftJS, and
the tests assert this project's mapping contract against the real posts. Version
drift is a future consideration, not a correctness bug, since only the DraftJS
output is normative here.

[`data.markdown`'s weighted-length rule is host-enforced and not fully
specified] -> Measure a documented conservative budget and fail with a
diagnostic instead of emitting a document the API rejects.

[Table-to-Markdown regeneration can lose cell formatting that the original
source had] -> tables are preserved as their original source slice where
available; a regenerated table is the fallback, and the tests assert cell
contents survive.

[Uploaded media and created drafts are external state that a failed run can
leave behind] -> draft creation reports the identifier it created and reuses an
upload through its digest-to-`media_id` cache; the cache lives in the task's
ignored output directory, so a re-run in the same task workspace reuses an
identifier while a lost cache costs a re-upload rather than a duplicate draft.

[All six image references in the existing posts are SVG, which the media
endpoints cannot accept] -> the reporter names each unacceptable image with its
source position instead of emitting an unresolvable entity, and new posts follow
the raster-only policy so they do not depend on that report.

[A new post could still reference an SVG and be rejected only at conversion
time] -> the authoring guidance states the raster-only requirement, and the
end-to-end check asserts that a policy-conforming post converts without an
unacceptable-image diagnostic, so the failure is caught before an upload attempt.

[`mermaid_webp` adds a second render path that could drift from the SVG render,
or rasterize a diagram whose text is unreadable at the chosen scale] -> it reuses
the existing renderer, theme, and pinned inputs, and its check compares the
render's geometry and label text against the same diagram's SVG render rather
than asserting only that a file was produced.

[Front matter carries the title that `content_state` cannot, so a naive
conversion would emit metadata as article text and leave the draft without its
required title] -> the input model parses front matter separately, reports the
title as document metadata, and fails conversion when no non-empty title exists.

## Open Questions

- The token's required scopes beyond the documented `tweet.read`,
  `tweet.write`, and `users.read` — media upload additionally needs
  `media.write` — can be confirmed against a live credential when draft creation
  is first exercised; the specs require the reference, not the exact scope list.
- Whether X accepts a draft with zero `atomic` blocks or requires at least one
  body block is unknown; the specs do not depend on it, and a live call settles
  it.

## Banner metadata and media reuse

The first front matter `images` entry selects the banner; subsequent entries
remain site metadata. Conversion records a separate optional banner locator,
with no additional body block. Banner and body locators share the existing
post-relative image source and digest contract. Draft creation preflights all
image sources before external effects, then resolves media using one cache.
Identical banner and body bytes reuse one unexpired media identifier.

The X [draft endpoint](https://docs.x.com/x-api/articles/create-draft-article)
and [OpenAPI schema](https://docs.x.com/openapi.json) define optional
`cover_media` with required `media_category` and `media_id` strings. The client
uses a typed request carrying that object separately from `content_state`.
Local HTTP E2E checks cover the exact wire requests, separate and shared media,
cache expiry, absent banners, malformed metadata, and refused image locators.
The tests were written and observed failing before implementation. Live X
acceptance remains unverified: the user's pause on live requests remains in
force and local validation does not spend a draft allowance.

## Validation evidence

The native light/dark raster, historical raster, post freshness, and site
browser checks passed (10 targets). Full uploader E2E and converter command
checks passed (2 targets), including the local HTTP banner scenarios. The
regenerated real-post artifact contains 39 blocks, one body image locator,
and a separate banner locator sharing its digest. The rendered dark diagram
was inspected for readable text, visible arrows, and modest canvas margins.
No live X request or site deployment was performed.
