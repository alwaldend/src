## 1. Project scaffold and declarations

- [x] 1.1 Add the `projects/x_article_uploader` README documenting the two operations, their inputs and outputs, and the credential boundary; verify it matches `specs/x-article-uploader-build/spec.md` (Project documentation)
- [x] 1.2 Add the project BUILD with a `docs` filegroup and declare the Go library and both `cmd/` binaries; verify `bazel_agent bazel build //projects/x_article_uploader:docs` succeeds
- [x] 1.3 Exclude `x_article_uploader` from `//projects:deploy_heads` in `projects/BUILD.bazel`, since this project ships no release target; verify `bazel_agent bazel query 'deps(//projects:deploy_heads, 1)'` no longer requires `//projects/x_article_uploader/releases:head`
- [x] 1.4 Add a `site/content` landing package mirroring a sibling project; verify `bazel_agent bazel build //projects/alwaldend.com:site` succeeds
- [x] 1.5 Add the project to `PROJECTS` in `projects/projects.bzl`; verify `bazel_agent bazel build //projects:docs` succeeds
- [x] 1.6 Confirm the owner workspace registration added in task 7.4 covers this project; verify `bazel_agent bazel test //infra/src/openspec/validation:validate_test` validates this change

## 2. Markdown parser dependency

- [x] 2.1 Declare `github.com/yuin/goldmark`, the parser the user approved on 2026-09-29, and its transitive modules as pinned external build inputs through the owning dependency workflow, exposing the modules the project imports through the owning `use_repo` declaration; verify the applicable lock and package-check targets pass
- [x] 2.2 Confirm the parser exposes the CommonMark and GFM constructs the posts use — tables, footnotes, and strikethrough — with a focused parser test; verify the test asserts each construct rather than assuming a default

## 3. Behavior coverage first

Author the end-to-end conversion check and the failure cases before the
converter exists, so the implementation is written against established
expected outcomes.

- [x] 3.1 Write the end-to-end conversion check that discovers every post under the site content tree and asserts, per post, the outcome recorded for it: the mapped blocks, inline ranges, and entity payloads plus the continuing diagnostics its constructs produce for a post whose images are all acceptable, and the named failing diagnostic at the reported source position for a post that references an unacceptable image; verify it runs against the current posts and fails only because the converter is absent
- [x] 3.2 Write the failure cases for conversion — a source without a non-empty title, an image whose media type the upload endpoints reject, an over-budget `markdown` payload, an artifact locator that escapes the post directory, and an unrepresentable construct — enumerating expected diagnostics and source positions; verify each case fails only because the converter is absent
- [x] 3.3 Write the end-to-end check that emits a repeatable draft artifact — the document, its parsed title, and its image locators — to a task-owned path for every post whose recorded outcome is the mapping contract, and assert the converted document parses as `content_state` with the title matching the source front matter and, for each image, a locator entry that records the post package and the bytes' digest; verify it fails only because the converter is absent

## 4. Converter library

- [x] 4.1 Implement the input model that separates YAML front matter from the Markdown body and reports the parsed title; verify the front matter never appears in emitted text and a missing title fails per case 3.2
- [x] 4.2 Implement block mapping for paragraphs, both list kinds, nested list items, and block quotes; verify case 3.1 asserts the resulting block type for each
- [x] 4.3 Implement heading mapping with clamping — `#`/`##` to `header-one`, `###` to `header-two`, `####` and deeper to `header-three`; verify case 3.1 covers every depth present in the posts
- [x] 4.4 Implement `inline_style_ranges` for bold, italic, and strikethrough with offsets measured against the block's final text; verify case 3.1 asserts exact offset and length for the posts' emphasis
- [x] 4.5 Implement `link` entities and their `entity_ranges`; verify case 3.1 asserts the entity URL and the selecting range
- [x] 4.6 Implement `atomic` blocks with `markdown` entities for fenced code and tables, preserving source Markdown; verify case 3.1 asserts the payload keeps the fence language and the table's rows and cell contents
- [x] 4.7 Implement the per-article `markdown` payload measurement with a conservative budget and a failing diagnostic; verify case 3.2 fails conversion on an over-budget document
- [x] 4.8 Implement footnote handling — bracketed in-text references and a trailing `Footnotes` heading with definitions as ordered list items; verify case 3.1 covers the post with a footnote and asserts no heading is added when none exist
- [x] 4.9 Implement image handling as `atomic` blocks with unresolved `image` entities carrying alt text as `caption`, record each image's source locator beside the document with the post package it resolves against and the bytes' digest, and report images whose media type the upload endpoints reject; verify case 3.1 asserts the recorded failing outcome for the post with SVG diagrams, case 3.2 asserts the rejection diagnostic, and case 3.3 asserts a locator per image
- [x] 4.10 Implement conversion diagnostics that report unsupported constructs with source position and fail instead of silently omitting content; verify case 3.2 asserts the position and the failure
- [x] 4.11 Implement inline-code spans so the span's literal text survives in the block's final text and the lost code styling is reported as a continuing diagnostic, and implement thematic breaks as `atomic` blocks backed by `divider` entities; verify case 3.1 asserts both constructs for the posts that contain them and that neither is dropped or left as literal paragraph characters
- [x] 4.12 Record each image locator's `image` entity by `entity_key`, so draft creation resolves an uploaded `media_id` by the entity the locator names rather than by position or caption; verify case 3.3 asserts each locator's recorded entity key and case 8.1 asserts a `media_id` lands on the named entity
- [x] 4.13 Verify conversion determinism by converting the same input twice and asserting byte-identical output
- [x] 4.14 Resolve the Markdown escapes and character references a text node carries for prose, a link label, and an image caption while keeping a raw span such as a code span literal — including inside a caption, which is flattened from its inline content and must be built node by node rather than resolved as a whole; verify the prose and caption cases assert the resolved characters and the preserved raw literal
- [x] 4.16 Emit the Articles wire shape rather than canonical DraftJS: block ranges as `inline_style_ranges` and `entity_ranges` with no `depth`, and `entities` as `{key, value}` entries, matching the endpoint's schema; verify the live draft endpoint accepts a converted artifact and case 3.1 asserts the emitted field names
- [x] 4.15 Keep one list item one block: accumulate a loose item's paragraphs into a single list-item block, collecting the item's block constructs in source order and emitting them after its text rather than emitting one block per paragraph, so a paragraph that resumes after a nested list still joins the same item and a later construct never overtakes an earlier one, and keep a heading inside a block quote as quoted text while reporting the lost level, tracking quote ancestry through a nested list; verify the loose-item case asserts one block per item and the quoted-heading cases assert a `blockquote` block carries the heading and the loss is reported for a direct and a list-nested heading

## 5. Converter command

- [x] 5.1 Implement the conversion entry point that reads a post and writes the emitted draft artifact, including the image locators beside the document; verify the end-to-end check 3.3 passes for every post whose recorded outcome is the mapping contract and that a post expected to fail produces its diagnostic instead of an artifact
- [x] 5.2 Verify the command performs no network access during conversion, by running it with network access unavailable and confirming success

## 6. Raster-only image policy

A published post's images have to be media types the upload endpoints accept,
and a Mermaid diagram has no raster form on its own. This group adds the render
path and the authoring rule, and proves them against the converter.

- [x] 6.1 Add the end-to-end check that renders a diagram, asserts the output is an accepted media type produced through the pinned render inputs, and converts a post referencing that render without an unacceptable-image diagnostic; verify it fails only because the render path is absent
- [x] 6.2 Add the raster render path under `tools/mermaid` for a WebP render that reuses the maintained render contract — pinned browser, theme, fonts, and paint-order pass — and encodes WebP through the pinned browser rather than a new image-conversion dependency; verify the check 6.1 passes for the render and its output is WebP
- [x] 6.3 Track the owner changes for the raster policy and the render rule: `projects/alwaldend.com/openspec/changes/add-raster-blog-images` owns the rule for a syndicated post and the authoring guidance, and `tools/mermaid/openspec/changes/add-raster-diagram-renders` owns the render rule; verify each validates strictly in its own workspace
- [x] 6.4 Extend the maintained Mermaid authoring guidance with the raster rule and the accepted media types, and update the blog authoring guidance so a post's images use an accepted media type with a diagram referenced through the raster render beside `index.md`; verify the guidance names the rule and the accepted types

## 7. Credential injection and owner changes

The draft command needs a credential the repository can inject, and the raster
policy and render rule live in workspaces this project does not own. This
group carries the uploader's non-secret half and tracks the owner changes.

- [x] 7.1 Add the uploader's `al.lua` with its `al_config` target, AppRole name, and injector plugin call selecting the OAuth 1.0a credential reference, injecting the API key, API secret, access token, and access token secret into their named environment variables, including the plugin's data and label wiring; verify the configuration target builds and names every reference without a value
- [x] 7.2 Carry out `infra/vault/openspec/changes/add-x-article-uploader-credentials`, which owns the concrete Terraform identity, least-privilege policy, and credential reference this injection reads; verify that change's tasks complete and its structural checks pass without authenticating to Vault
- [x] 7.3 Track `projects/alwaldend.com/openspec/changes/add-raster-blog-images` as the owner of the raster-image rule and authoring guidance, and `tools/mermaid/openspec/changes/add-raster-diagram-renders` as the owner of the raster render rule; verify each validates strictly in its own workspace
- [x] 7.4 Add this project's `//projects/x_article_uploader/openspec:source` and the new `//tools/mermaid/openspec:source` to `_WORKSPACE_SOURCES` in `infra/src/openspec/validation/BUILD.bazel`, since neither workspace was registered before; verify `bazel_agent bazel test //infra/src/openspec/validation:validate_test` includes both and passes

## 8. Draft command

- [x] 8.1 Author the publisher's behavior coverage before its implementation: its failure cases — a missing one of the four OAuth 1.0a fields, an unresolved image, an image whose locator digest does not match, a rejected request, and draft creation requested without a parsed title — enumerating expected diagnostics, and its success cases — a recorded-response check that the draft request body carries `title` and `content_state`, one that the request carries an `Authorization` header with an RFC 5849 `oauth_signature`, one that each uploaded image's `media_id` lands on the entity its locator names, and one that draft creation alone does not publish; verify each fails only because the publisher is absent
- [x] 8.2 Implement credential loading of the four OAuth 1.0a fields through the repository's injection flow, with no credential value in source; verify case 8.1 fails with a diagnostic naming the absent field's reference
- [x] 8.3 Implement image upload through the media endpoint keyed by content hash, resolving `media_items` with `media_category` and `media_id` and attaching each identifier to the entity its locator names; verify case 8.1 asserts an unresolved image is not sent as resolved, that a locator whose digest does not match is refused, and that identical bytes reuse one identifier within a run and across a later run through the publisher's cache
- [x] 8.4 Implement draft creation against `POST /2/articles/draft`, sending the parsed title from the converted draft artifact; verify case 8.1's recorded-response check passes and still fails when the artifact carries no title
- [x] 8.5 Remove the publish operation: no publish endpoint call, no publish flag, and no `Publish` method, so no invocation can publish an article; verify the draft command exposes only `--artifact` and that a live run creates a draft and reports that it is not published
- [x] 8.6 Implement failure reporting that names the failing operation, reports created identifiers on success, and does not auto-retry on authentication or request errors; verify case 8.1 covers a rejected request and reports the created identifier on success

## 9. Validation

- [x] 9.1 Run `bazel_agent bazel build //projects/x_article_uploader:all` and confirm all targets build
- [x] 9.2 Run `bazel_agent bazel test //projects/x_article_uploader:all` and confirm the end-to-end conversion check and every failure case pass
- [x] 9.3 Run `bazel_agent bazel run //tools/openspec -- validate --all --strict --no-interactive` with `OPENSPEC_PROJECT=projects/x_article_uploader` and confirm this change validates strictly
- [x] 9.4 Confirm the new project is excluded from `//projects:deploy_heads` and that `//projects:docs` and `//projects/alwaldend.com:site` still build
- [x] 9.5 Run the repository formatter and confirm no unrelated target changes

## 10. Banner metadata

- [x] 10.1 Write and observe failing HTTP E2E coverage for banner selection, shared and separate uploads, absent metadata, cache lifetime, and preflight failures
- [x] 10.2 Resolve the first `images` entry during conversion, attach `cover_media` during draft creation, and preserve body images
- [x] 10.3 Pass the full uploader E2E suite and inspect the regenerated real-post artifact without contacting X
- [x] 10.4 Validate and archive the completed owner change
