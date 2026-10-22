## Why

The 2026-10-22 repository review found 86 fully checked task lists among 97
non-archived OpenSpec changes. Completed and superseded work obscures the
remaining backlog, while several implemented deltas have not reached their
owning baseline specifications.

## What Changes

- Reconcile each candidate against its planning artifacts, source, and retained
  evidence before archiving; preserve unfinished or unverified work as active.
- Merge implemented requirements into their natural owner's baseline without
  restoring superseded interfaces or duplicating another owner's contract.
- Preserve historical artifacts, record per-change archive decisions, and
  repair maintained references affected by archive moves.
- Include newly archived owner workspaces in the existing structural checks.
- Record review findings 1–4 and 6 as explicitly deferred owner-local changes
  with unchecked implementation tasks and concrete acceptance conditions.

## Capabilities

### New Capabilities

None. Deferred proposals describe future work and do not change current behavior.

### Modified Capabilities

None. This maintenance applies already accepted component deltas and restores
the existing work-state and validation contracts; it adds no product behavior.

## Impact

Owner-local OpenSpec changes and specifications, maintained archive links, and
the repository OpenSpec validation declarations. No production implementation,
live infrastructure operation, new dependency, or deferred fix is included.
Acceptance requires an auditable archive inventory, preservation of outstanding
work, valid deferred plans, structural validation, and the repository delivery
checks on the exact candidate.
