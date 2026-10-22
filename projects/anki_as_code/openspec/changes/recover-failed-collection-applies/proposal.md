## Why

Deferred; implementation has not started. Review finding 3 identifies a
documented nontransactional apply: native mutations can persist while a later
failure prevents binding persistence. A created deck can then require manual
identity recovery before the intended apply can be retried.

## What Changes

- Provide actionable recovery information after partial execution or failed
  binding persistence, without reporting success or assuming rollback.
- Preserve object identity, scheduling, review history, media, and immutable
  input configuration during recovery.
- Define a supported failure-to-recovery-to-retry workflow, with explicit
  refusal when observed collection state cannot safely be reconciled.
- Accept through an E2E exercise that fails after native creation, recovers,
  retries, and proves no duplicate objects or lost collection data.

## Capabilities

### New Capabilities

- `collection-apply-recovery`: recover from partially executed collection plans.

### Modified Capabilities

None. Existing nontransactional semantics remain the baseline until this
deferred change is implemented.

## Impact

Planning only is authorized. Future work affects reconciliation, collection
state and CLI error/recovery presentation, plus native collection behavior
tests. Recovery mechanism selection remains open; this proposal does not
mandate transaction emulation, rollback, or a new dependency.

Baseline: `e0d24cb3a`, inspected 2026-10-22. Evidence:
[execution before binding persistence](../../../internal/reconcile/reconciler.go)
and [existing injected native failures](../../../internal/reconcile/execute_test.go).
Those tests establish expected partial state in source; they were not rerun
while preparing this deferred plan.
