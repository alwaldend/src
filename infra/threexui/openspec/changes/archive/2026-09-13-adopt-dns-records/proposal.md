## Why

3x-ui DNS records need retained ownership in the owner's Terraform state. This change
records its completed adoption in the
[coordinated migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/)
while preserving existing records and their provider identities.

## What Changes

- Adopt existing declared records into the owner's `tf_setup` state.
- Keep DNS management enabled by default in the adopted source revision.
- Require no-change reconciliation and evidence that unrelated records and
  declared DNS views remain intact.

## Capabilities

### Modified Capabilities

- `owned-dns`: Retain adopted ownership through scoped DNS execution.

## Impact

The `//infra/threexui/tf_setup` wrappers, DNS activation default, Terraform state,
and scoped DNS grants for the existing `src_infra_threexui` AppRole. Adoption
completed on 2026-09-13 with preserved record values and provider identities;
the [evidence](design.md#verified-adoption) records the scoped validation boundary.
