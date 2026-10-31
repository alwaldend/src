## Why

Hardware diagnostics need a reusable authenticated SSH entry point instead of
ad hoc credential handling.

## What Changes

- Add a host SSH target using the existing Vault password injector and host SSH client.
- Document read-only invocation and host prerequisites.

## Capabilities

### New Capabilities

- `host-ssh`: Authenticated SSH access to the XCP-ng host through AL.

## Impact

Only repository command packaging changes; no host deployment or new dependency.
