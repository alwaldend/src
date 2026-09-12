## Why

Ingress DNS records need adoption into the owning Terraform state while
preserving existing records as part of the
[coordinated migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/).

## What Changes

- Match and import the existing records through `//infra/ingress/tf`.
- Retain enabled source defaults after adoption and verify a no-change plan.
- Preserve unrelated DNS records and existing non-DNS resources.

## Capabilities

### Modified Capabilities

- `owned-dns`: Retain adopted ownership through scoped DNS execution.

## Impact

The owner's `tf` state and DNS activation default, plus scoped DNS access
through its configured Vault AppRole. Adoption completed on 2026-09-13;
the [evidence](design.md#verified-adoption) records preserved provider identities.
