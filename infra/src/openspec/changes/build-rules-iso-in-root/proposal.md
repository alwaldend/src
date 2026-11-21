## Why

The ruleset's nested workspace adds a separate Bazel build and test cycle to
repository CI.

## What Changes

- Build and test the ISO ruleset from the root workspace.
- Point the XCP-ng and TrueNAS repository declarations at the root package.

## Capabilities

### Modified Capabilities

- `repository` builds the ISO ruleset as part of root Bazel CI.

## Impact

The ruleset retains an independent archive for registry consumers.
