## Why

The initial OpenSpec integration collected every component's specifications
and migrated work in one repository-root directory. The requested ownership
is local: components keep their own requirements and history, and `infra/src`
describes evolution of the repository itself.

## What Changes

- Move 48 component baselines and nine migrated changes into their owners.
- Move the repository baseline and migration index to `infra/src/openspec`.
- Select owner workspaces in the pinned CLI, context routing and skill.
- Validate all 49 workspaces, including projects with standalone Bazel modules.

## Capabilities

### New Capabilities

### Modified Capabilities

- `repository`: local specification ownership and repository evolution scope.

## Impact

OpenSpec source paths, declared Bazel inputs, context routing, documentation
and migration navigation change. Original historical bytes, acceptance and
execution state remain preserved. The goal skill stays disabled and the goal
tool remains deprecated. No deployment or new external dependency is involved.
