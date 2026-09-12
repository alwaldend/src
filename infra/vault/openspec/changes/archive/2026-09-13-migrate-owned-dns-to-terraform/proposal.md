## Why

The DNS migration assigns record state to `infra/vault` while retaining its canonical
`dnsconfig.json` declarations and existing dc1, global views.

## What Changes

- Integrate the reusable DNS Terraform module into `tf_setup` with an explicit
  disabled-by-default adoption gate.
- Package the canonical declaration, pinned providers, and scoped Vault injection.
- Retain existing service, provisioning, and Ansible authority.

## Capabilities

### New Capabilities

- `owned-dns`: owner-local Terraform DNS configuration and offline checks.

### Modified Capabilities

None.

## Impact

The owner Terraform root, AL configuration, and Bazel packaging change. Shared
migration and cutover requirements remain in the [DNS change](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/proposal.md).
No live backend, Vault, or DNS operation is part of implementation.
