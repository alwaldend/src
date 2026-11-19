## Why

OpenHands needs owner-local reconciliation for its declared DNS records.
This change records its completed activation within the
[coordinated migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md),
preserving existing provider records while provisioning verified missing declarations.

## What Changes

- Compare declarations with complete provider inventories through the scoped
  `tf_setup:dns` workflow and provision only verified missing records.
- Persist enabled DNS ownership in the reviewed adoption revision.
- Require unchanged existing records, an exact additions-only plan for missing
  records, and a no-change follow-up plan with DNS resolution checks.

## Capabilities

### Modified Capabilities

- `owned-dns`: Persist enabled owner-local reconciliation after preserving
  existing records and provisioning verified missing declarations.

## Impact

The OpenHands `tf_setup` DNS activation default and Terraform state, plus scoped
DNS credential grants for its existing `src_infra_openhands` Vault AppRole.
Canonical declarations, existing records, VM configuration, and non-DNS
authentication and backend paths retain their owning sources.
The [adoption evidence](design.md) records 13 creations, zero imports, and the
verification-only recovery after the initial post-apply inventory check halted.
