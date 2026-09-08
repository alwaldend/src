# Phase 2 goal publication and recovery plan

## Bindings

- Goal: `agent-system-phase-2`
- Source base / merged PR 39: `1ca06a0fc9697e8fb212b32d99d9ad3b996ea76e`
- Goal generation: 2
- Lifecycle generation: 4
- Criteria revision: 2
- Expected checkpoint resource version: 5
- Intended attempt: `goal-publication-001`
- Work type: `change`

The goal tool binds the portable goal-state and criteria digests when it
publishes this attempt.

## Target defects

Goal publication uses a multi-file protocol whose intermediate state fails
validation but has no stable diagnostic or supported recovery operation. A
catalog therefore cannot yet advertise rich continuation safely. During goal
initialization, the repeatable `--criterion` flag also split commas inside one
criterion into separate criteria; this bounded CLI defect must be fixed rather
than worked around.

## Decision and smallest slice

Proceed with the Phase 2A recoverability gate before rich goal-catalog joins.
Keep Phase 2B work read-only and design-bound during this attempt. Reject a
database or daemon: the existing filesystem store can expose a durable
publication intent and deterministic doctor/recover behavior without changing
owner authority.

## Ready workstreams

1. Audit every existing publication boundary, failure return, validation
   symptom, and migration guarantee in `projects/goal`.
2. Independently audit the Phase 1 declaration universe and identify the
   minimum owner-local inputs and output contracts for Phase 2 catalogs.
3. Coordinator: repair comma-preserving criterion flags, integrate audit
   evidence, select the smallest recovery representation, implement it, and
   validate the exact candidate.

The two audits have disjoint read-only outputs under
`out/agent_system_phase_2/audits/`; only the coordinator writes canonical goal
state and maintained source.

## Acceptance for this attempt

- `--criterion` and migration criterion overrides preserve commas, with a
  command regression test.
- Publication boundaries and recoverability invariants are enumerated against
  current source and tests.
- The smallest recovery protocol either passes injected-boundary tests or the
  attempt closes early with decisive reset evidence.
- Legacy migration remains non-destructive and current fixed goal-store
  regressions pass.
- No catalog claims resumability for an invalid record.

## Fixed regressions

- `//projects/goal/...` focused tests and builds.
- `//tools/agents/...` shared-contract and Phase 1 report tests.
- Root Buildifier if Bazel metadata changes.
- Goal catalog validation and `git diff --check`.

## Strategy reset

Reset if supported recovery needs cross-process transactional storage, a
database, or inference from generated projections. In that case, narrow the
catalog to validated identity and coarse unavailable status and retain the
fault-boundary evidence for a separate protocol design.
