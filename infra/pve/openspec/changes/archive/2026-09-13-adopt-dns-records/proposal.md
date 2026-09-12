## Why

Proxmox DNS records need adoption into the owning Terraform state while
preserving existing records as part of the
[coordinated migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/).

## What Changes

- Match and import existing records through the scoped `//infra/pve/tf:dns` targets.
- Retain enabled source defaults after adoption and verify a no-change plan.
- Preserve unrelated DNS records and existing non-DNS resources.

## Capabilities

### Modified Capabilities

- `owned-dns`: Retain enabled ownership after authorized record adoption.

## Impact

The owner's `tf` state and DNS activation default, with scoped DNS access
through its existing Vault AppRole. Adoption is complete; [verified evidence](design.md)
records the saved-plan imports and DNS reconciliation boundary.
