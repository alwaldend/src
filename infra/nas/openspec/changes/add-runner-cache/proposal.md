## Why

Each AppRole needs an isolated NAS account and encrypted dataset namespace.
NAS should group child dataset declarations by AppRole and reuse dataset defaults.

## What Changes

- Create accounts from the canonical Vault truenas_admins and truenas_user groups; each AppRole storage project owns its private namespace and SMB shares.
- Mirror administrator/user membership and publish generated 64-character passwords through existing shared-read credential paths.
- Give each infra/nas/approles/<approle> project its own Terraform state and plan/apply targets, using the NAS operator identity. Keep namespaces and child datasets out of the shared NAS state.

## Capabilities

### New Capabilities

- `runner-cache-export`: AppRole NAS account/namespace access and reusable child storage.

### Modified Capabilities

None.

## Impact

Uses existing pinned TrueNAS, Vault and random providers. Existing adopted NAS
addresses remain intact. No JSON configuration or live deployment is introduced.
