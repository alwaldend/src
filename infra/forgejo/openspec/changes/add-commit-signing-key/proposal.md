## Why

Forgejo enables commit signing without a configured instance key. Rebuilding
the service must preserve its signing identity through Vault.

## What Changes

- Generate a dedicated Ed25519 SSH signing pair once during setup Terraform.
- Cache the pair in the existing Forgejo AppRole namespace with create-only writes.
- Inject and install the pair before controlled Forgejo activation.
- Configure SSH signing for the existing merge and CRUD policies.
- Reject corrupt caches and implicit replacement of deployed private keys.

## Capabilities

### New Capabilities

- `commit-signing`: Vault-backed instance signing identity and deployment wiring.

### Modified Capabilities

None.

## Impact

Forgejo setup Terraform, AL injection, the shared Vault SSH key tool, and the
shared Forgejo role's optional signing-key installation. Uses existing pinned
Go dependencies and the host's existing OpenSSH tools. No live operations.
