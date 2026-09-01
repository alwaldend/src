# Phase 2A goal publication and recovery result

## Outcome

This attempt closes with `refine`. It implements and validates the Phase 2A
recoverability gate selected by `goal-publication-001` and claims
`goal-recovery` support for the goal tool publication machinery. It does not
claim `legacy-migration`, `bounded-catalogs`, `system-index`,
`context-capsule`, `runtime-isolation`, or `resource-baseline`; those remain
open Phase 2 criteria.

## Implemented

The filesystem store now writes a bounded `.goal-publication.yaml` intent
containing deterministic after-images and before/after digests before any
multi-file mutation, stages after-images under
`.goal-publication-stage/<operation>/`, and exposes `goal doctor` and
`goal recover`. `doctor` is a read-only classifier (`stable`,
`discardable-intent`, `staged-intent`, `partial-intent`, `conflict`,
`committed-projection-stale`); `recover` replays remaining after-images for a
staged or partial intent, refuses a conflict, and discards a discardable
intent. A pre-commit failure (before the first `goal.yaml` rename) discards
the intent and preserves the exact prior record; a post-commit failure returns
the committed Goal reference and leaves a pending intent that normal mutations
fail closed on with a stable `goal publication is incomplete` error.
`checkpoint`, criteria updates, `set-relationships`, `promote`, `migrate`,
`render`, and `show` all enter through the pending-intent gate.

## Test coverage

Boundary-complete injected-failure tests cover every canonical rename
(`goal.yaml` first, attempt directory, finalization, README projection) for
new and existing attempt checkpoints, criteria updates, and
`set-relationships`. `TestExistingAttemptCheckpointPublishesGoalTokenFirst`,
`TestNewAttemptCheckpointPublishesGoalTokenFirst`,
`TestCheckpointREADMEFailureReportsCommittedVersion`,
`TestCriteriaREADMEFailureReportsCommittedVersion`,
`TestSetRelationshipsReportsCommittedProjectionFailure`,
`TestDoctorReportsStableForValidGoal`,
`TestCheckpointLeavesRecoverableIntentOnPublishFailure`, and
`TestRecoverDiscardsIntentWhenNothingPublished` pass. The focused
`//projects/goal/...` suite passes six of six tests.

`init`, `migrate`, and `promote` publish a complete new record with one atomic
directory rename, so they cannot expose a torn canonical record and satisfy
the recovery criterion by construction. They are not fault-injected through
the intent protocol in this attempt.

## Not in scope

`legacy-migration`, `bounded-catalogs`, `system-index`, `context-capsule`,
`runtime-isolation`, and `resource-baseline` remain open. `goal-recovery`
evidence is limited to the goal-tool publication machinery and does not extend
to catalog ingestion; `legacy-migration` still needs source-reference and
mapping provenance plus stronger snapshot binding.
