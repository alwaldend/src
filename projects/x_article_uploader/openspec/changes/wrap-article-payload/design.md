## Context

Conversion builds offline JSON for later media resolution and draft creation. The user requested payload for X fields, custom fields beside it, no serialized diagnostics, and toggleable failure on warnings/errors.

## Goals / Non-Goals

- Make payload the only X request projection, retaining verified local media references outside it.
- Default to failing command conversion on warnings, with an explicit warnings-as-errors switch.
- Preserve fatal error handling and local-only validation; no upload or deployment is authorized.

## Decisions

- Embed an explicitly JSON-nested payload in the emission and decoding models. Their existing typed document and JSON-map representations remain appropriate for conversion and media attachment respectively.
- Keep diagnostics in memory for reporting, never in serialized artifacts. The command prints them before deciding whether to write output.
- Errors always fail; warnings-as-errors defaults to true and can be disabled to produce content with acknowledged formatting loss.
- Require payload.content_state on read rather than silently interpreting legacy flat artifacts.
- Resolve the banner into payload.cover_media, and body media IDs into payload.content_state.entities.

## Risks / Trade-offs

The wire format change requires regenerating saved artifacts. Existing posts with inline-code style loss require the explicit warning override. The unrelated working-tree blog edits remain user-owned and excluded from this change.

## Validation

Failure cases were enumerated and command E2E written before implementation in out/article-payload/failure-cases.md. Existing media, banner, cache and schema E2E checks exercise the updated artifact boundary. Successful command tests retain repeatable JSON artifacts.
