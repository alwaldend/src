## Why

Harbor's owner-local Terraform root is prepared for DNS management but still
defaults to disabled. Its pending adoption must preserve existing records while
establishing durable ownership under the
[coordinated migration](../../../../../dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform).

## What Changes

- Match and import the owner's existing DNS records through its `tf_setup` root.
- Persist enabled DNS ownership in the reviewed adoption revision.
- Require a no-change adoption plan and verify resolution and unrelated records.

## Capabilities

### Modified Capabilities

- `owned-dns`: Retain enabled owner-local DNS management after adopting existing records.

## Impact

The Harbor `tf_setup` DNS activation default and Terraform state, scoped DNS
grants for the existing `src_infra_harbor` AppRole, and owner adoption evidence.
Canonical declarations and existing record identities and values are preserved.
