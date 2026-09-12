## Why

The [owner DNS migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/proposal.md)
requires component identities for Terraform state and DNS credentials. Several
owners have no AppRole, and existing identities lack the relevant DNS reads.

## What Changes

- Add missing component AppRoles as cohesive modules below
  `infra/vault/tf/approles/<name>/`, reusing the existing AppRole module.
- Grant each DNS identity reads only for its required provider credentials.
- Reuse existing component identities and preserve their resource addresses;
  give host_bot a separate DNS identity without changing Ansible authentication.
- Package the local modules and document source validation and deployment order.

## Capabilities

### New Capabilities

### Modified Capabilities

- `infra-vault`: Add owner-scoped DNS identities and provider-secret access.

## Impact

Changes belong to `infra/vault/tf` and its OpenSpec workspace. Credential values
remain in Vault. Applying the Vault configuration and DNS cutover require
separate explicit operator authorization.
