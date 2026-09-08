# Phase 2A goal publication and recovery plan

## Bindings

- Goal: `agent-system-phase-2`
- Goal generation: 2
- Lifecycle generation: 4
- Criteria revision: 2
- Expected checkpoint resource version: 8
- Prior attempt: `goal-publication-001` (`refine`)
- Work type: `change`

The goal tool binds the portable goal-state and criteria digests when it
publishes this attempt.

## Target defect

Goal publication uses a multi-file protocol whose intermediate state fails
validation but has no stable diagnostic or supported recovery operation. A
catalog therefore cannot yet advertise rich continuation safely. This attempt
implements the recoverability gate selected by `goal-publication-001`: a
bounded, versioned per-goal publication intent with exact staged after-images
and before/after digests, plus deterministic `doctor` and `recover` behavior.

## Decision and smallest slice

Proceed with the Phase 2A recoverability gate before rich goal-catalog joins.
Keep Phase 2B work read-only and design-bound during this attempt. Reject a
database or daemon: the existing filesystem store can expose a durable
publication intent and deterministic doctor/recover behavior without changing
owner authority.

## Ready workstreams

1. Coordinator: implement the publication intent, staging, doctor, recover,
   and fail-closed gates in `projects/goal`.
2. Coordinator: add boundary-complete injected-failure tests for every
   publication boundary and the committed/incomplete error states.
3. Coordinator: update the goal command, README, architecture, and
   concurrency reference documentation for the new protocol.

The audits from `goal-publication-001` already enumerate the publication
boundary matrix and catalog input contracts; this attempt is implementation
and validation of the recovery slice only.

## Acceptance for this attempt

- Publication has a stable committed state and a stable `incomplete`
  diagnosis: `doctor` classifies `stable`, `discardable-intent`,
  `staged-intent`, `partial-intent`, `conflict`, and
  `committed-projection-stale`.
- `recover` replays every remaining after-image for a staged or partial
  intent, refuses a conflict, and discards a discardable intent to restore the
  exact prior record.
- A pre-commit failure discards the intent and preserves the prior valid
  record; a post-commit failure returns the committed Goal reference and
  leaves a pending intent that normal mutations fail closed on.
- Boundary-complete injected-failure and doctor/recover tests pass against
  exact record identities; the legacy migration and goal-store regressions
  remain green.
- No catalog claims resumability for an invalid record.

## Fixed regressions

- `//projects/goal/...` focused tests and builds.
- Repository Buildifier if Bazel metadata changes.
- Goal catalog validation and `git diff --check`.

## Strategy reset

Reset if supported recovery needs cross-process transactional storage, a
database, inference from generated projections, or a second serialization
authority. In that case, narrow the catalog to validated identity and coarse
unavailable status and retain the fault-boundary evidence for a separate
protocol design.
