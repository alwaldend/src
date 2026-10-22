## Context

Deferred; implementation has not started. Authority is planning only. The
[proposal](proposal.md) identifies the current README ordering. Documentation
spans setup, native and WASM engines, collection workflows, and generated CLI
reference links; moving it requires preserving those link and ownership
boundaries.

## Goals / Non-Goals

**Goals:** Put one verified default workflow first and make detailed engine
reference discoverable without duplicating generated command help.

**Non-Goals:** Change command behavior, choose a new default engine, shorten
documentation to a line target, or redefine repository prerequisites.

## Decisions

- Start with default-engine prerequisites and one disposable collection
  workflow; then link advanced engine material from a compact reference index.
  Reordering alone is possible but leaves a large monolithic README; moving
  cohesive sections under the owner's `docs/` keeps reference discoverable.
- Preserve warnings about exclusive collection access, native cancellation,
  and WASM limitations with the workflows they constrain. Do not make an
  apparently simpler quickstart less accurate.
- Keep generated CLI help generated, and link shared repository setup rather
  than copy it. The linked contributor-onboarding record owns host portability.
- Set `skip_specs: true` because this is a documentation-only change. A design
  is retained to explain cross-page moves and avoid broken reference paths.

## Risks / Trade-offs

- Moving sections can break existing anchors -> inspect inbound links and
  preserve stable destinations or update all owned references deliberately.
- A quickstart can accidentally imply support for an untested engine/host ->
  name and execute one actual supported configuration in acceptance.
- Examples can modify valuable collections -> use a disposable fixture and
  identify the resulting artifact so the procedure remains repeatable.

## Migration Plan

Next action after implementation is requested: map the current sections and
inbound links, select the default workflow, and prepare the documentation move.
Build and inspect the rendered pages, execute the quickstart on a disposable
collection, and preserve its repeatable result artifact. No source behavior
or end-user collection is changed during this deferred planning work.
