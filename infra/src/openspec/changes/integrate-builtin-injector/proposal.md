## Why

The user requested AL's built-in injector and advertised resource providers after
merging PR #276. Repository config and fixtures must use the new Vault provider.

## What Changes

- Compose the standalone Vault provider through the shared config.
- Remove external injector registration from built-in-only fixtures.
- Keep operational target labels, selections and credential references unchanged.

## Capabilities

### New Capabilities

None. Runtime behavior and protocol are owned by projects/al/openspec/changes/builtin-resource-injection.

### Modified Capabilities

None. This integration change uses skip_specs.

## Impact

Shared AL configuration, command fixtures and current workflow documentation.
No live infrastructure operations or new dependencies.
