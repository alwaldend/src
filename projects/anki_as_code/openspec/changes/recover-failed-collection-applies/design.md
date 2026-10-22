## Context

Deferred; implementation has not started. Authority is planning only. See the
[proposal](proposal.md) and its native failure-test evidence. The reconciler
executes mutations before replacing bindings; its API explicitly permits
persistent database writes after a later failure. Recovery must coexist with
caller-owned collection lifetime and exclusive offline access.

## Goals / Non-Goals

**Goals:** Make the known partial state understandable and recoverable while
retaining object identity and validating a fresh retry against current state.

**Non-Goals:** Promise a transaction across native APIs, automatically undo
review history, rewrite source identities, or change archive publication.

## Decisions

- Base the candidate on explicit recovery information and the existing state
  operations. Record completed mutations and their exact native identities
  when known; a generic retry message cannot identify an unbound created deck.
- Compare structured failure output with a durable execution journal during
  the first deferred design task. Durable recording may improve interruption
  recovery but adds consistency and lifecycle responsibilities. Process-crash
  recovery is not implied by this plan's handled-failure contract.
- Require observed identity/state agreement before binding recovered objects.
  Name-only automatic adoption could claim an unrelated deck; unconditional
  retry could create duplicates. Refuse ambiguous mappings explicitly.
- Preserve the existing native failure injection as the starting acceptance
  fixture and extend the exercised workflow through recovery and retry before
  implementation. Do not add tests that merely restate the internal call order.
- Use `decision-review` on the concrete persistence and recovery mechanism
  before commitment; update these artifacts if it changes the proposed contract.

## Risks / Trade-offs

- Failure while writing recovery information can hide completed native work ->
  define which information is available in memory, in error output, and after
  process exit before selecting a persistence mechanism.
- A user edits the collection before recovery -> compare current identities
  and state and reject stale recovery rather than overwrite changes.
- Partial binding writes can falsely imply all planned operations succeeded ->
  distinguish recovered objects from completion of the original plan.

## Migration Plan

Next action after implementation is requested: describe handled failure modes,
choose the recovery artifact/command contract, and establish an E2E failing
apply on a disposable native collection. Implement the smallest supported
recovery path, then verify native IDs, bindings, review history, media, and
unchanged input configuration after retry. Save a repeatable sanitized report
and collection-state evidence. No existing user collection is modified by
this planning work.
