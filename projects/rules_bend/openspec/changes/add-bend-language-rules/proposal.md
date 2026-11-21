## Why

Bazel projects need to run and check Bend 2 programs with declared sources,
imports, and compiler inputs.

## What Changes

- Add `bend_check` and `bend_binary` rules under `projects/rules_bend`.
- Package the rules for standalone Bazel consumers and resolve Bend through registered toolchains.
- Register the project in the documentation site.

## Capabilities

### New Capabilities

- `bazel-bend`: Run and check Bend 2 programs from Bazel targets.

### Modified Capabilities

None.

## Impact

The rules resolve Bend 2 compiler executables through Bazel toolchains. The
default Linux x86-64 toolchains use the pinned upstream release distribution.
