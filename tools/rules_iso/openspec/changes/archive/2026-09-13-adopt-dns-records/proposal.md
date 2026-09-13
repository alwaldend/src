## Why

This project's prepared Terraform root needs an adoption record for the
[coordinated DNS migration](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/).
Existing public DNS records must retain their values and provider identities
when their ownership moves into the project's state.

## What Changes

- Verify or provision the project's scoped DNS access before adoption.
- Import the existing records, persist enabled ownership, and require a no-change plan.
- Verify public resolution and unrelated record preservation.

## Capabilities

### Modified Capabilities

- `project-dns`: Retain enabled ownership and no-change reconciliation after adoption.

## Impact

The project's `tf` state, DNS activation default, and scoped Vault access,
operated through its documented nested-project adapter. Deployment is pending.
