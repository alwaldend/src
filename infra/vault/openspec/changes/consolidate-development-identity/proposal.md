## Why

The development-host consolidation now includes a single DNS root. Its shared identity needs DNS credential access, while the inactive Actions worker no longer needs a registration-token read.

## What Changes

- Grant src_infra_dev read access to the existing dc1 and public DNS provider credentials.
- Remove its unused legacy runner-config read and preserve the required OpenHands read.
- Retain both legacy identities until their DNS states have been reconciled under separate authorization.

## Capabilities

### New Capabilities

- `development-identity`: One scoped identity for development deployment and DNS.

## Impact

Vault Terraform policy declarations and infra/dev deployment wiring. No live Vault change, credential movement or state mutation is part of source implementation. Linked owner change: infra/dev/openspec/changes/consolidate-development-hosts.
