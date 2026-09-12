## Why

Host Bot's prepared DNS root needs to adopt its existing records before it can
reconcile them. This change tracks the host's pending adoption within the
[coordinated migration](../../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform).

## What Changes

- Provision scoped DNS access for the dedicated host DNS AppRole.
- Match and import the host's existing records into its DNS-only `tf` root.
- Retain enabled ownership after adoption and verify a no-change plan, both
  DNS views, and preservation of unrelated records.

## Capabilities

### Modified Capabilities

- `owned-dns`: Retain enabled ownership after authorized adoption while
  preserving the preparatory disabled state.

## Impact

The host's DNS Terraform state, activation default, and dedicated
`src_users_simeonwarren_host_bot` AppRole. Ansible retains its existing
`user_simeonwarren` authority. Canonical declarations remain with the host;
record values and provider identities must remain unchanged during adoption.
