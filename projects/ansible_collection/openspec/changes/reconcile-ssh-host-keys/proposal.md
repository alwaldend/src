## Why

Partial cloud-init host-key sets fail SSH hardening because expected keys are
missing. Secure1 currently compensates with a host-specific command.

## What Changes

- Generate missing standard host keys in the shared SSH role before discovery and signing.
- Preserve existing private keys and converge without changes on repeat runs.
- Remove secure1's special-case key-generation task.

## Capabilities

### New Capabilities

- `ssh-host-key-reconciliation`: Complete partial host-key sets before shared SSH configuration.

### Modified Capabilities

None.

## Impact

Shared SSH role and its development-host consumer. Reuse community.crypto,
already required by SSH hardening. No new dependency or live deployment.
