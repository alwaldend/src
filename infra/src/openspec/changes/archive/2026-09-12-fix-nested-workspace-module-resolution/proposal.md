## Why

Every nested Bazel workspace under `projects/` failed to resolve its module
graph. `projects/rules_openspec/MODULE.bazel` declared six sibling
`bazel_dep`s without versions, and their `local_path_override` entries only
apply when that module is the root. A nested workspace that depends on
`rules_openspec` therefore failed with `bad bazel_dep ... with no version`,
so no nested module could be built or tested standalone.

`projects/rules_iso` and `projects/rules_openspec` also declared their own
`MODULE.bazel` without a matching `.bazelignore` entry, so root-workspace
target expansion crossed the workspace boundary and failed to load
`@rules_openspec_npm`.

## What Changes

- Remove the six unused, versionless sibling dependencies from
  `projects/rules_openspec/MODULE.bazel`, keeping the one dependency its
  BUILD files actually load.
- Add `projects/rules_iso` and `projects/rules_openspec` to `.bazelignore`.
- Regenerate the affected nested module locks through the owning tool.

## Capabilities

### New Capabilities

### Modified Capabilities

- `repository`: every nested Bazel workspace resolves and builds standalone.

## Impact

Nested module manifests, locks and the root ignore list change. No runtime
behavior, deployment or external dependency changes.
