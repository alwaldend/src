## Why

Repository SSH commands need the host client with AL's existing credential
injection and cleanup lifecycle.

## What Changes

- Add `tools/ssh:defs.bzl` with an `ssh_binary` macro using `al_binary_run`.
- Provide a default SSH target and document consumer-owned injection settings.

## Capabilities

### New Capabilities

- `ssh-command`: Run host SSH through AL with configurable injection and arguments.

### Modified Capabilities

None.

## Impact

Only the new `tools/ssh` package. Reuses AL without adding dependencies or
changing host configuration.
