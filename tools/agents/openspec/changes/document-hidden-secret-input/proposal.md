## Why

Users need pasteable Bash commands that request credentials without displaying them or placing their literal values in history or process arguments.

## What Changes

- Document hidden interactive input, standard-input delivery and subshell cleanup in repo-secrets.
- Add an evaluation case for this guidance and verify the command with dummy input.

## Capabilities

### New Capabilities

None; documentation and evaluation only (`skip_specs: true`).

### Modified Capabilities

None.

## Impact

Canonical repo-secrets skill and its existing offline evaluation assets. No host changes or dependencies.
