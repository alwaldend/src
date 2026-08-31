# Phase 2 goal publication investigation result

## Outcome

This attempt closes with `refine`. It bound Phase 2 to merged PR 39, repaired
the goal CLI defect encountered while initializing the record, enumerated the
complete goal publication crash surface, and selected a dependency-safe
recovery and catalog sequence. It does not claim either `goal-recovery` or
`legacy-migration` acceptance.

## Implemented bounded repair

The repeatable `init --criterion` and `migrate --criterion` flags now preserve
commas inside each argument instead of treating them as hidden list
separators. Command-level regressions cover both paths. The full focused
`//projects/goal/...` suite passed: six of six tests.

## Recovery verdict

Current checkpoints detect torn records but cannot recover them because every
normal command performs full validation before it can retry. Criteria
replacement has the largest untested gap; new and existing attempt publication
also expose invalid canonical states after the first rename.

The next implementation should add a bounded, versioned per-goal publication
intent containing exact staged after-images and before/after digests. Under the
existing lock, `doctor` classifies stable, recoverable-incomplete,
projection-stale, and conflict states; `recover` idempotently rolls forward
nonconflicting files and refuses unknown content. This is smaller than an
immutable-generation/current-pointer redesign and preserves the no-daemon,
no-database guard.

Legacy migration keeps its non-destructive full-directory staging, but still
needs a portable source reference, explicit mapping provenance, imported-byte
digest rebinding, and atomic no-replace publication before its criterion can
pass.

## Catalog boundary

Static Topology, Policy, Action, Capability, WorkspaceCheck, and descriptor-
only AgentSystemIndex work can proceed beside recovery. GoalCatalog live
ingestion must wait for a pure read-only recovery-aware inspector. Until then,
invalid/interrupted eligible goals are explicit `unavailable` entries and no
catalog may advertise resumability.

## Next attempt

Implement the publication intent and pure doctor inspection first, beginning
with criteria replacement and then new/existing attempt checkpoints. Inject
failure at every intent and canonical rename, prove repeatable recovery, and
only then connect coarse validated goal status to GoalCatalog.
