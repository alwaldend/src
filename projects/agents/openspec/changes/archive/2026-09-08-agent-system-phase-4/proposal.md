# Join durable work, delivery, review, version, and release

## Why

The accepted Phase 4 change connected durable continuation, candidate validation, review evidence, and immutable release identities while preserving the owners of work, tests, delivery, and versioning.

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline.

## What Changes

- Structured continuation: The historical AttemptSpec and goal catalog carried the current attempt, stable defect, next action, and resume condition so fresh agents could find maintained open work. Native OpenSpec artifacts now replace that discovery workflow.
- Delivery consumes exact evidence: Delivery preparation and publication bind candidate ValidationSet identity rather than trusting a caller-supplied head. Review receipts retain work, delivery, and defect joins without becoming acceptance authorities.
- Guarded immutable releases: Versioning owns the version, channel, commit, and tree-state handoff. ReleaseRefPlan and receipts bind immutable artifact/ref names, exact remote leases, and verified publication; unsupported atomicity is an explicit refusal.

## Capabilities

### New Capabilities

None. This archive records historical changes to an existing repository capability.

### Modified Capabilities

- `project-agents`: Join durable work, delivery, review, version, and release.

## Impact

The historical implementation affected the repository agent system and the component owners cited in its accepted plan and evidence. This import changes work-record location only; old plans do not grant present authority to deploy, publish, or change host configuration.

See [migration status](migration.json), [original objective](provenance/source/goal.yaml), [accepted result](provenance/source/attempts/phase4-implementation/result.md), and [full provenance](provenance/README.md).
