## Why

Bazel executables need a single-command way to run under systemd resource
limits without generating a separate launch script.

## What Changes

- Add a Linux `tools/systemd_run` executable and an executable Bazel rule.
- Declare the binary, arguments, data, and systemd properties in BUILD targets.
- Retain direct CLI execution with repeatable systemd properties.
- Preserve command arguments, environment, working directory, signals, and exit status.

## Capabilities

### New Capabilities

- `scoped-execution`: Launch an executable in a resource-limited user scope.

### Modified Capabilities

None.

## Impact

New Go command using the existing Cobra pin. No new dependencies, persistent
services, or host configuration. ComfyUI usage is recorded in its own change.

## Module organization

Move the launcher into the standalone `tools/rules_systemd` Bzlmod module.
Expose `@rules_systemd//systemd:defs.bzl`, keep runtime implementation private,
and integrate module documentation and workspace discovery with the parent.
Reuse current dependency pins and preserve the CLI and executable contract.
