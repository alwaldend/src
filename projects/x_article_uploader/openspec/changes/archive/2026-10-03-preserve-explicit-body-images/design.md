## Context

The previous converter removed an opening image prefix when its validated
bytes match the selected banner. The site change
`projects/alwaldend.com/openspec/changes/render-blog-frontmatter-images` moves
opening-image presentation into Hugo, so conversion should follow the explicit
Markdown body without filtering it against metadata.

## Goals / Non-Goals

Preserve body images and their order while keeping the first `images` entry as
the separate article banner. This change does not prepend metadata images to
DraftJS, modify upload caching, change image validation, or call X.

## Decisions

- Remove the opening-image AST filter and its specialized collection helper.
  The existing Markdown traversal already emits image blocks and locators in
  source order; retaining the filter would continue deleting explicit content.
- Keep banner selection and digest-based upload reuse intact. A deliberately
  repeated image may appear in the body and cover while sharing one upload.
- Let the Hugo owner remove redundant opening Markdown references from current
  posts. The uploader does not need site presentation rules to interpret a body.

## Risks / Trade-offs

- A post explicitly repeating its cover in Markdown will show both references
  in the X draft. This preserves the author's source; current blog duplication
  is handled by the linked site change.
- Real-post expectations depend on the companion source cleanup. The X article
  keeps its metadata banner with zero body image locators once its opening
  Markdown diagram is removed. The Diagrams post still fails on its remaining
  unsupported SVG body images.

## Validation

Existing HTTP E2E expectations were changed before production code. On
2026-10-03, the red run failed in seven opening, repeated, and identical-alias
cases because body blocks and locators were missing. Evidence is in
`out/hugo-blog-images/uploader-preservation-red.log`, Bazel invocation
`01f09f13-76c2-4f21-b5bf-b9d838093cf0`, on base `42878ef`.

The combined site's browser and output tests, uploader HTTP E2E suite,
conversion command tests, and blog skill configuration check passed in
`out/hugo-blog-images/integration.log`. A freshly converted story has 47 blocks,
its `pipeline.webp` banner, and no body images because the redundant Markdown
reference was removed from the source. The explicit-image fixture cases now
preserve all opening/repeated images and reuse one media upload for identical
bytes. Independent review found no blockers. The initial MODIFIED delta was refused because it dropped the superseded
suppression scenario names. The delta now explicitly removes that requirement
and adds the independent banner/body conversion contract, retaining validation
guarantees. Strict validation passed for this corrected delta.
No live X requests were made for this change.
