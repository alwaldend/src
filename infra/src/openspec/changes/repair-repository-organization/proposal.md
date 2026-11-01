## Why

Repository ownership is documented, but duplicated tool pins, regular workflow dependencies and independently maintained inventories undermine those boundaries.

## What Changes

- Route legacy Helm and Drawio build rules through their third-party dependency owners.
- Export nested module specifications without depending on the OpenSpec CLI.
- Generate specification-validation membership and the project index from owner declarations; share workspace discovery between CI and the full repository checker.
- Add an explicit artifact-dependency check alongside native package visibility enforcement.

## Capabilities

### New Capabilities

- `repository-inventory`: Derived owner inventories and bounded artifact-dependency auditing.

### Modified Capabilities

None. Binary lock schema changes are owned by `tools/rules_binary_toolchain/openspec/changes/validate-binary-lock-schema`.

## Impact

Shared Bazel integration, documentation, CI discovery and validation. No new external dependencies, host configuration changes, deployment or merge.
