# Normalize agent semantics and declare safety contracts

## Why

The original Phase 1 audit found that shared types alone did not establish a closed registered operation universe, isolated scratch, or explicit mutation ownership. The accepted work normalized agent semantics and made those safety boundaries inspectable.

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline.

## What Changes

- Owner-local declarations: Operation, lifecycle, and generated-artifact facts stay with component owners. A report-only checker derives a closed registered universe and names legacy or unavailable entries instead of becoming a hand-maintained central registry.
- Task and run isolation: Cordis scratch is scoped to task and run identities with bounded manifests. Bazel action and test temporary directories retain Bazel ownership.
- Explicit mutation and measured resource limits: The unnamed Terraform mutation alias was replaced by its explicit operation. The runner doctor reports source, pin, profile, platform, scratch, and staleness without changing host state. A criteria-bound baseline records measurements and unavailable fields.

## Capabilities

### New Capabilities

None. This archive records historical changes to an existing repository capability.

### Modified Capabilities

- `project-agents`: Normalize agent semantics and declare safety contracts.

## Impact

The historical implementation affected the repository agent system and the component owners cited in its accepted plan and evidence. This import changes work-record location only; old plans do not grant present authority to deploy, publish, or change host configuration.

See [migration status](migration.json), [original objective](provenance/source/goal.yaml), [accepted result](provenance/source/attempts/phase1-completion-001/result.md), and [full provenance](provenance/README.md).
