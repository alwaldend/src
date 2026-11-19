## Why

XCP-ng DNS records need adoption into the owning Terraform state while
preserving existing records as part of the
[coordinated migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/).

## What Changes

- Match and import the existing records through the scoped DNS plan, show, and
  saved-plan apply commands in `//infra/xcp_ng/tf`.
- Retain enabled source defaults after adoption and verify a no-change plan.
- Preserve unrelated DNS records and existing non-DNS resources.

## Capabilities

### Modified Capabilities

- `owned-dns`: Retain enabled ownership after authorized record adoption and
  specify the DNS-only command boundary.

## Impact

The owner's `tf` state and DNS activation default, plus scoped DNS access
through its configured Vault AppRole. Adoption is complete; the owner design
records the verified scope and evidence.
