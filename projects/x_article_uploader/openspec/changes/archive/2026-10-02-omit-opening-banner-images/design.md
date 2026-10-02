## Context

The converter already reads and hashes the selected banner separately from the
body. Goldmark provides top-level paragraph and inline image nodes; the draft
publisher operates on the resulting artifact and needs no change.

## Goals / Non-Goals

The conversion change removes only the opening banner repetition. Source Markdown,
draft request structure, upload caching, and live X behavior are outside this change.

## Decisions

- Inspect the top-level AST prefix before emitting body blocks. Tracking whether
  any DraftJS blocks have been emitted would incorrectly include nested list or
  quotation images and deferred list content.
- Compare digests only after the existing anchored read and media validation.
  This handles equivalent relative references and identical image aliases; a
  path-only comparison would retain an identical cover under another filename.
- End omission at the first nonmatching image or substantive nonimage content.
  A global image filter would erase later intentional illustrations.
- Keep mixed text and linked image nodes intact. They carry content beyond the
  standalone image that the cover can represent.

## Risks / Trade-offs

- Identical files with different names can be omitted in the opening prefix.
  Later repetitions remain available by placing them after other content.
- Invalid opening images must keep their normal diagnostic behavior. Leave them
  in the conversion path instead of treating failed validation as a match.
- A post containing only its opening banner can have an empty body document;
  preserve the banner and do not invent placeholder body text.

## Validation

Write and observe failing HTTP E2E scenarios before implementation. Record red
and green evidence under `out/x-article-banner-dedup/`, including the serialized
artifact and local fixture-service requests. No real X request is authorized.

On 2026-10-02, before implementation, the new HTTP cases failed for seven opening
image scenarios and the real article retained an unwanted body locator. The
failure is recorded in `out/x-article-banner-dedup/red.log` (Bazel invocation
`47589c98-260b-45b7-b9c4-34dc1b7f264b`).

The implementation on base `0ae03f4d` passed the complete uploader E2E target and
conversion command tests (Bazel invocation
`13e26f51-f56e-4d5e-89ec-75ccd4dc98b9`). Strict change validation also passed.
Representative artifacts and fixture-service HTTP requests are preserved under
`out/x-article-banner-dedup/evidence/`: the opening image disappears from the
body while the cover remains, a later repeat retains its body entity and shares
one upload, and a partial image-only paragraph preserves the remaining image
order. The real article has 47 body blocks, no body image locators, and its
`pipeline.webp` banner locator.

Independent review approved the implementation and inspected the emitted HTTP
evidence with no correctness blockers. The README and body-image specification
were clarified to distinguish retained images from the omitted opening prefix.
The parent owns final candidate formatting, gates, and delivery. No live X calls
were made.
