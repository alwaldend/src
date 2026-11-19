## Why

Shared apex and mail records need adoption into the `infra/dns` Terraform state
while preserving their existing identities and behavior. This is the shared
owner's batch in the [coordinated migration](../2026-09-13-migrate-project-dns-to-terraform/).

## What Changes

- Match and import the canonical shared records through `//infra/dns/tf`.
- Retain enabled defaults after adoption and verify a no-change plan.
- Preserve both declared DNS views and unrelated provider records.

## Capabilities

### Modified Capabilities

- `infra-dns`: Retain enabled shared-record ownership after authorized adoption.

## Impact

The shared owner's Terraform state, DNS activation default, and scoped access
through the existing `src_infra_dns` AppRole. Adoption completed on 2026-09-13;
the [evidence](design.md#verified-adoption) records preserved provider identities
and enabled ownership.
