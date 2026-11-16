## Why

The runner is enabled, but no workflow schedules repository CI and the shared
command still invokes a removed AL interface. Enable CI only for trunk pushes.

## What Changes

- Run the Forgejo workflow only for pushes to master.
- Check out the event revision with the existing pinned checkout mirror, then
  invoke the local CI action.
- Launch Bazel directly from tools/ci for build and test in every declared
  root/nested workspace, retaining complete failure collection.
- Exercise the actual command executable with an isolated Bazel process fixture.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `repository`: activate trunk-only CI and repair its executable boundary.

## Impact

.forgejo/workflows, tools/ci and the linked development-host documentation. No new dependencies,
runner deployment, credentials, protected-ref push or merge is required.
