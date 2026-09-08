# Generate a trustworthy system map and bounded context capsule

## Why

The original Phase 2 work joined owner-local facts into deterministic, bounded catalogs and context views, while making interrupted durable publication and optional runtime failures explicit.

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline.

## What Changes

- Bounded catalogs and descriptor-only index: Topology, policy, action, capability, workspace-check, and historical goal catalogs are derived from owner facts with input digests, conflicts, limits, and completeness. AgentSystemIndex contains catalog descriptors and query routes rather than copied catalog bodies.
- Explicit unavailable states: The offline context command returns relevant path, label, or task context with provenance and structured partial results when optional inputs are unavailable. Runtime isolation uses package deadlines, leases, namespaces, and expected-revision publication.
- Recoverable legacy publication: The historical goal store made multi-file publication interruption recoverable and preserved legacy source identities and raw bytes during imports. That compatibility implementation is historical; OpenSpec now owns maintained work.

## Capabilities

### New Capabilities

None. This archive records historical changes to an existing repository capability.

### Modified Capabilities

- `project-agents`: Generate a trustworthy system map and bounded context capsule.

## Impact

The historical implementation affected the repository agent system and the component owners cited in its accepted plan and evidence. This import changes work-record location only; old plans do not grant present authority to deploy, publish, or change host configuration.

See [migration status](migration.json), [original objective](provenance/source/goal.yaml), [accepted result](provenance/source/attempts/goal-publication-008/result.md), and [full provenance](provenance/README.md).
