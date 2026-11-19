## Why

Forgejo's DNS adoption preserves existing records while transferring their
management to this owner's Terraform state through scoped DNS commands.

## What Changes

- Adopt this owner's existing DNS records through the `dns.plan`, `dns.show`,
  and saved-plan `dns.apply` targets in `//infra/forgejo/tf_setup`.
- Keep DNS ownership enabled by default after adoption and require a no-change
  adoption plan that preserves owned and unrelated records.
- Retain an owner-specific adoption receipt before synchronizing the specification.

## Capabilities

### Modified Capabilities

- `owned-dns`: Persist enabled ownership after authorized adoption of existing
  records and select the scoped DNS execution boundary.

## Impact

The owner's Terraform state, `tf_setup/dns.tf` activation default, setup documentation,
and scoped DNS execution through the existing `src_infra_dc1_forgejo1` AppRole are affected.
The [shared cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, and recovery.
The adoption evidence records a historical DNS snapshot; aggregate source
validation remains with repository delivery.
