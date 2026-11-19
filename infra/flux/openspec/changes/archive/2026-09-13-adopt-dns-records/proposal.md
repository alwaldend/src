## Why

Flux has a preparatory DNS integration whose resources are disabled.
The coordinated DNS migration needs an owner-specific adoption that preserves
existing records while transferring their management to this owner's Terraform state.

## What Changes

- Adopt this owner's existing DNS records through `//infra/flux/tf_setup`.
- Keep DNS ownership enabled by default after adoption and require a no-change
  adoption plan that preserves owned and unrelated records.
- Retain an owner-specific adoption receipt before synchronizing the specification.

## Capabilities

### Modified Capabilities

- `owned-dns`: Persist enabled ownership after authorized adoption of existing records.

## Impact

The owner's Terraform state, `tf_setup/dns.tf` activation default, setup documentation,
and scoped DNS grants for the existing `src_infra_flux` AppRole are affected.
The [shared cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, and recovery.
Deployment is pending; these artifacts do not establish live ownership.
