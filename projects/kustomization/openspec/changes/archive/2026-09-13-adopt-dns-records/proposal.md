## Why

The user authorized project-by-project DNS deployment while preserving existing
records. This change adopts this project's public DNS through the
[coordinated migration](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform).

## What Changes

- Import this project's existing records and retain enabled ownership afterward.
- Verify a no-change adoption plan and preservation of unrelated DNS records.

## Capabilities

### Modified Capabilities

- `project-dns`: Preserve enabled ownership after authorized adoption.

## Impact

This project's Terraform state, DNS activation default, and scoped Vault access.
Record values and provider identities must remain unchanged.
