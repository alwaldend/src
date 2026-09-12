## Why

This change adopts the existing records for `rules_template` through its
prepared DNS root as part of the
[coordinated migration](../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform),
preserving live records as ownership moves into the project's state.

## What Changes

- Verify or provision scoped DNS access through the owning Vault workflow.
- Adopt the existing records through the documented root-workspace adapter.
- Retain enabled source defaults after adoption and verify no-change reconciliation.
- Record project-specific adoption evidence and preservation of unrelated records.

## Capabilities

### Modified Capabilities

- `project-dns`: Persist enabled ownership after authorized adoption.

## Impact

The project's Terraform state, DNS activation default, and dedicated Vault
AppRole and credential access. The existing canonical DNS declarations and
reusable Bazel module boundary remain the source and packaging owners.
