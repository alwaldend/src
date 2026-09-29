## 1. Project scaffold and declarations

- [ ] 1.1 Add the `projects/x_article_uploader` README documenting the two operations, their inputs and outputs, and the credential boundary; verify it matches `specs/x-article-uploader-build/spec.md` (Project documentation)
- [ ] 1.2 Add the project BUILD with a `docs` filegroup and declare the Go library and both `cmd/` binaries; verify `bazel_agent bazel build //projects/x_article_uploader:docs` succeeds
- [ ] 1.3 Exclude `x_article_uploader` from `//projects:deploy_heads` in `projects/BUILD.bazel`, since this project ships no release target; verify `bazel_agent bazel query 'deps(//projects:deploy_heads, 1)'` no longer requires `//projects/x_article_uploader/releases:head`
- [ ] 1.4 Add a `site/content` landing package mirroring a sibling project; verify `bazel_agent bazel build //projects/alwaldend.com:site` succeeds
- [ ] 1.5 Add the project to `PROJECTS` in `projects/projects.bzl`; verify `bazel_agent bazel build //projects:docs` succeeds
- [ ] 1.6 Confirm the owner workspace registration added in task 7.4 covers this project; verify `bazel_agent bazel test //infra/src/openspec/validation:validate_test` validates this change

## 2. Markdown parser dependency

- [ ] 2.1 Declare `github.com/yuin/goldmark`, the parser the user approved on 2026-09-29, and its transitive modules as pinned external build inputs through the owning dependency workflow, exposing the modules the project imports through the owning `use_repo` declaration; verify the applicable lock and package-check targets pass
- [ ] 2.2 Confirm the parser exposes the CommonMark and GFM constructs the posts use — tables, footnotes, and strikethrough — with a focused parser test; verify the test asserts each construct rather than assuming a default

## 3. Behavior coverage first

Author the end-to-end conversion check and the failure cases before the
converter exists, so the implementation is written against established
expected outcomes.

- [ ] 3.1 Write the end-to-end conversion check that discovers every post under the site content tree and asserts, per post, the outcome recorded for it: the mapped blocks, inline ranges, and entity payloads plus the continuing diagnostics its constructs produce for a post whose images are all acceptable, and the named failing diagnostic at the reported source position for a post that references an unacceptable image; verify it runs against the current posts and fails only because the converter is absent
- [ ] 3.2 Write the failure cases for conversion — a source without a non-empty title, an image whose media type the upload endpoints reject, an over-budget `markdown` payload, an artifact locator that escapes the post directory, and an unrepresentable construct — enumerating expected diagnostics and source positions; verify each case fails only because the converter is absent
- [ ] 3.3 Write the end-to-end check that emits a repeatable draft artifact — the document, its parsed title, and its image locators — to a task-owned path for every post whose recorded outcome is the mapping contract, and assert the converted document parses as `content_state` with the title matching the source front matter and, for each image, a locator entry that records the post package and the bytes' digest; verify it fails only because the converter is absent

## 4. Converter library

- [ ] 4.1 Implement the input model that separates YAML front matter from the Markdown body and reports the parsed title; verify the front matter never appears in emitted text and a missing title fails per case 3.2
- [ ] 4.2 Implement block mapping for paragraphs, both list kinds, nested list items, and block quotes; verify case 3.1 asserts the resulting block type for each
- [ ] 4.3 Implement heading mapping with clamping — `#`/`##` to `header-one`, `###` to `header-two`, `####` and deeper to `header-three`; verify case 3.1 covers every depth present in the posts
- [ ] 4.4 Implement `inline_style_ranges` for bold, italic, and strikethrough with offsets measured against the block's final text; verify case 3.1 asserts exact offset and length for the posts' emphasis
- [ ] 4.5 Implement `link` entities and their `entity_ranges`; verify case 3.1 asserts the entity URL and the selecting range
- [ ] 4.6 Implement `atomic` blocks with `markdown` entities for fenced code and tables, preserving source Markdown; verify case 3.1 asserts the payload keeps the fence language and the table's rows and cell contents
- [ ] 4.7 Implement the per-article `markdown` payload measurement with a conservative budget and a failing diagnostic; verify case 3.2 fails conversion on an over-budget document
- [ ] 4.8 Implement footnote handling — bracketed in-text references and a trailing `Footnotes` heading with definitions as ordered list items; verify case 3.1 covers the post with a footnote and asserts no heading is added when none exist
- [ ] 4.9 Implement image handling as `atomic` blocks with unresolved `image` entities carrying alt text as `caption`, record each image's source locator beside the document with the post package it resolves against and the bytes' digest, and report images whose media type the upload endpoints reject; verify case 3.1 asserts the recorded failing outcome for the post with SVG diagrams, case 3.2 asserts the rejection diagnostic, and case 3.3 asserts a locator per image
- [ ] 4.10 Implement conversion diagnostics that report unsupported constructs with source position and fail instead of silently omitting content; verify case 3.2 asserts the position and the failure
- [ ] 4.11 Implement inline-code spans so the span's literal text survives in the block's final text and the lost code styling is reported as a continuing diagnostic, and implement thematic breaks as `atomic` blocks backed by `divider` entities; verify case 3.1 asserts both constructs for the posts that contain them and that neither is dropped or left as literal paragraph characters
- [ ] 4.12 Record each image locator's `image` entity by `entity_key`, so publication resolves an uploaded `media_id` by the entity the locator names rather than by position or caption; verify case 3.3 asserts each locator's recorded entity key and case 8.1 asserts a `media_id` lands on the named entity
- [ ] 4.13 Verify conversion determinism by converting the same input twice and asserting byte-identical output

## 5. Converter command

- [ ] 5.1 Implement the conversion entry point that reads a post and writes the emitted draft artifact, including the image locators beside the document; verify the end-to-end check 3.3 passes for every post whose recorded outcome is the mapping contract and that a post expected to fail produces its diagnostic instead of an artifact
- [ ] 5.2 Verify the command performs no network access during conversion, by running it with network access unavailable and confirming success

## 6. Raster-only image policy

A published post's images have to be media types the upload endpoints accept,
and a Mermaid diagram has no raster form on its own. This group adds the render
path and the authoring rule, and proves them against the converter.

- [ ] 6.1 Add the end-to-end check that renders a diagram, asserts the output is an accepted media type produced through the pinned render inputs, and converts a post referencing that render without an unacceptable-image diagnostic; verify it fails only because the render path is absent
- [ ] 6.2 Add the raster render path under `tools/mermaid` for a WebP render that reuses the maintained render contract — pinned browser, theme, fonts, and paint-order pass — and encodes WebP through the pinned browser rather than a new image-conversion dependency; verify the check 6.1 passes for the render and its output is WebP
- [ ] 6.3 Track the owner changes for the raster policy and the render rule: `projects/alwaldend.com/openspec/changes/add-raster-blog-images` owns the rule for a syndicated post and the authoring guidance, and `tools/mermaid/openspec/changes/add-raster-diagram-renders` owns the render rule; verify each validates strictly in its own workspace
- [ ] 6.4 Extend the maintained Mermaid authoring guidance with the raster rule and the accepted media types, and update the blog authoring guidance so a post's images use an accepted media type with a diagram referenced through the raster render beside `index.md`; verify the guidance names the rule and the accepted types

## 7. Credential injection and owner changes

The publisher needs a credential the repository can inject, and the raster
policy and render rule live in workspaces this project does not own. This
group carries the uploader's non-secret half and tracks the owner changes.

- [ ] 7.1 Add the uploader's `al.lua` with its `al_config` target, AppRole name, and injector plugin call selecting the credential reference, including the plugin's data and label wiring; verify the configuration target builds and names the reference without a value
- [ ] 7.2 Carry out `infra/vault/openspec/changes/add-x-article-uploader-credentials`, which owns the concrete Terraform identity, least-privilege policy, and credential reference this injection reads; verify that change's tasks complete and its structural checks pass without authenticating to Vault
- [ ] 7.3 Track `projects/alwaldend.com/openspec/changes/add-raster-blog-images` as the owner of the raster-image rule and authoring guidance, and `tools/mermaid/openspec/changes/add-raster-diagram-renders` as the owner of the raster render rule; verify each validates strictly in its own workspace
- [ ] 7.4 Add this project's `//projects/x_article_uploader/openspec:source` and the new `//tools/mermaid/openspec:source` to `_WORKSPACE_SOURCES` in `infra/src/openspec/validation/BUILD.bazel`, since neither workspace was registered before; verify `bazel_agent bazel test //infra/src/openspec/validation:validate_test` includes both and passes

## 8. Publisher

- [ ] 8.1 Author the publisher's behavior coverage before its implementation: its failure cases — missing credentials, an unresolved image, an image whose locator digest does not match, a rejected request, and draft creation requested without a parsed title — enumerating expected diagnostics, and its success cases — a recorded-response check that the draft request body carries `title` and `content_state`, one that each uploaded image's `media_id` lands on the entity its locator names, and one that draft creation alone does not publish; verify each fails only because the publisher is absent
- [ ] 8.2 Implement credential loading through the repository's injection flow, with no credential value in source; verify case 8.1 fails with a diagnostic naming the reference
- [ ] 8.3 Implement image upload through the media endpoint keyed by content hash, resolving `media_items` with `media_category` and `media_id` and attaching each identifier to the entity its locator names; verify case 8.1 asserts an unresolved image is not sent as resolved, that a locator whose digest does not match is refused, and that identical bytes reuse one identifier within a run and across a later run through the publisher's cache
- [ ] 8.4 Implement draft creation against `POST /2/articles/draft`, sending the parsed title from the converted draft artifact; verify case 8.1's recorded-response check passes and still fails when the artifact carries no title
- [ ] 8.5 Implement publication against `POST /2/articles/{article_id}/publish` as a separate explicitly requested operation; verify case 8.1's recorded-response check confirms draft creation alone does not publish
- [ ] 8.6 Implement failure reporting that names the failing operation, reports created identifiers on success, and does not auto-retry on authentication or request errors; verify case 8.1 covers a rejected request and reports the created identifier on success

## 9. Validation

- [ ] 9.1 Run `bazel_agent bazel build //projects/x_article_uploader:all` and confirm all targets build
- [ ] 9.2 Run `bazel_agent bazel test //projects/x_article_uploader:all` and confirm the end-to-end conversion check and every failure case pass
- [ ] 9.3 Run `bazel_agent bazel run //tools/openspec -- validate --all --strict --no-interactive` with `OPENSPEC_PROJECT=projects/x_article_uploader` and confirm this change validates strictly
- [ ] 9.4 Confirm the new project is excluded from `//projects:deploy_heads` and that `//projects:docs` and `//projects/alwaldend.com:site` still build
- [ ] 9.5 Run the repository formatter and confirm no unrelated target changes
