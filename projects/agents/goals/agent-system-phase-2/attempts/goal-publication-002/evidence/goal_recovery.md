# Goal publication and recovery evidence

## Test identity and verdict

- Candidate: working tree `t3code/check-1` HEAD `f7757ac1` with the Phase 2A
  publication intent, `doctor`, `recover`, and fail-closed gates applied.
- Bound goal: `agent-system-phase-2` generation 2, lifecycle generation 4,
  criteria revision 2.
- Evidence method: run publication-boundary fault injection, doctor, recovery,
  and validation tests against exact record identities.

## Local Go test run

```text
GOWORK=off GOCACHE=$PWD/out/agent_system_phase_2/gocache \
GOPATH=$PWD/out/agent_system_phase_2/gopath \
XDG_RUNTIME_DIR=$PWD/out/agent_system_phase_2/runtime \
go test ./projects/goal/...
```

PASS: all three package test targets (api/v1alpha1, cmd/goal,
internal/fsstore). The fsstore suite includes boundary-complete
fault-injection tests:

- `TestExistingAttemptCheckpointPublishesGoalTokenFirst` and
  `TestNewAttemptCheckpointPublishesGoalTokenFirst` cover every canonical
  rename for new/open/immediately-closed attempt checkpoints. A pre-commit
  failure preserves version 1 with no intent; a post-commit failure returns
  the committed version 2 reference.
- `TestCheckpointREADMEFailureReportsCommittedVersion` and
  `TestCriteriaREADMEFailureReportsCommittedVersion` cover README projection
  failures. The record stays at version 2/3 with a pending intent; a retry
  fails closed with `goal publication is incomplete`; `Recover()` completes
  the intended state; the stale pre-commit token is then rejected.
- `TestSetRelationshipsReportsCommittedProjectionFailure` covers relationship
  publication and verifies the committed projection failure is reported
  through the stable error.
- `TestDoctorReportsStableForValidGoal`,
  `TestCheckpointLeavesRecoverableIntentOnPublishFailure`, and
  `TestRecoverDiscardsIntentWhenNothingPublished` cover doctor classification
  (`stable`, `partial-intent`, `discardable-intent`) and idempotent recovery.
  `Recover()` discards a discardable intent, restores the prior record, and
  replays a partial intent to the intended final record.

## Bazel run

```text
bazel_agent test //projects/goal/...
```

PASS: six of six tests, including the `rules_skill` skill-validation aspect
for the updated concurrency reference document.

## Command-surface check

`goal --help` lists `doctor` and `recover`; `goal doctor --goal-dir ...`
returns `{"publicationState":"stable"}` on a valid record. The `init`,
`checkpoint`, `validate`, `promote`, `render`, `migrate`, and existing CLI
tests remain green.

## Scope

Fault injection covers every canonical multi-file publication boundary
(goal.yaml, attempt directory, finalization, README projection) for
checkpoint, criteria updates, and relationship updates. `init`, `migrate`, and
`promote` publish a complete new record with one atomic directory rename, so
they expose no torn canonical record and are not part of the intent protocol.
This evidence does not cover `legacy-migration` provenance or Phase 2 catalog
ingestion.
