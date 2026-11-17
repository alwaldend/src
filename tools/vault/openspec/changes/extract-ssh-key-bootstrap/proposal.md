## Why

The SSH key bootstrap under infra/forgejo only generates and validates a pair
in Vault. Its behavior is reusable and belongs to the Vault tools owner.

## What Changes

- Move the command to tools/vault/cmd/ssh_key with neutral inputs and descriptions.
- Preserve the existing Ed25519 record format and immutable cache behavior.
- Keep the actual Forgejo Terraform declaration and its acceptance fixture in infra.

## Capabilities

### New Capabilities

- `ssh-key-cache`: Create and reuse a dedicated Ed25519 SSH pair in Vault KV v2.

### Modified Capabilities

None.

## Impact

Vault tool packaging, Forgejo's consumer wiring and documentation. No new
dependencies, secret migration, live generation, Vault write, or deployment.
