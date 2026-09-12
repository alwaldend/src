## Why

The coordinated DNS migration transfers Forgejo Actions runner's existing record
into its owning Terraform state while preserving the live DNS record.

## What Changes

- Adopt this owner's existing DNS records through the scoped
  `//infra/forgejo_runner/tf_setup:dns` wrappers and existing setup backend.
- Keep DNS ownership enabled by default after adoption and require a no-change
  adoption plan that preserves owned and unrelated records.
- Retain an owner-specific adoption receipt before synchronizing the specification.

## Capabilities

### Modified Capabilities

- `owned-dns`: Persist enabled ownership and execute DNS operations through the
  dedicated `dns=1` flow scoped to `module.dns`.

## Impact

The owner's Terraform state, `tf_setup/dns.tf` activation default, scoped command
contract, and setup documentation are affected. DNS operations use the existing
`src_infra_forgejo_runner` AppRole and setup backend.
The [shared cutover procedure](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/cutover.md)
owns migration ordering, writer coordination, credential handling, and recovery.
