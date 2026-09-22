# Complete Phase 5 fixture and optimization follow-up

## Why

The recorded Phase 5 follow-up accepted a bounded routing-coverage baseline, an inventory of deterministic writable fixtures and their gaps, a policy for isolated live comparisons, and an evidence chain for optimization adoption.

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline.

## What Changes

- Bounded coverage is an observed subset: CoverageMatrix binds normalized routing cases and the exact capability catalog digest. Emitted cases are distinguished from the complete skill universe and truncation remains explicit.
- Fixture inventory preserves missing coverage: Accepted evidence inventories Git, forge, runner, and Cordis fixtures. It explicitly records missing Terraform, Ansible, Vault/injector, and Bazel end-to-end writable fixtures; those gaps were not implemented by the acceptance review.
- Live evaluation stays separate: Live comparisons are manual or scheduled, outside wildcard tests, and bind subject model, skill, catalog, fixture, judge, and isolated authentication identities. Offline tests validate harness configuration without model calls.
- Optimization adoption follows evidence: The accepted example traces measured friction, predeclared threshold, validated proposal, owner-local change, regression, delivered revision, fallback, and retirement. The source result still names remaining work although the final criterion review accepted the documented inventory and policy; both facts are retained.

## Capabilities

### New Capabilities

None. This archive records historical changes to an existing repository capability.

### Modified Capabilities

- `project-agents`: Complete Phase 5 fixture and optimization follow-up.

## Impact

The historical implementation affected the repository agent system and the component owners cited in its accepted plan and evidence. This import changes work-record location only; old plans do not grant present authority to deploy, publish, or change host configuration.

See [migration status](migration.json), [original objective](provenance/source/goal.yaml), [accepted result](provenance/source/attempts/followup-baseline-001/result.md), and [full provenance](provenance/README.md).
