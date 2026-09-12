## Why

The user authorized project-by-project DNS deployment while preserving existing
records. This project is the first public-only adoption in the
[coordinated migration](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/).

## What Changes

- Enable this project's DNS ownership after matching and importing its live record.
- Verify a no-change plan and preserve the unrelated public DNS inventory.

## Capabilities

### Modified Capabilities

- `project-dns`: Persist enabled ownership after authorized adoption.

## Impact

This project's Terraform state, DNS activation default, and dedicated Vault
AppRole. Existing record values and provider identities remain unchanged.
