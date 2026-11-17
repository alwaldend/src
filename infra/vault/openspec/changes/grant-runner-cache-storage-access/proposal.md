## Why

NAS owns per-AppRole accounts and needs to publish their credentials without
granting AppRoles the NAS administrator API key. XCP-ng needs the runner account
for its storage connection.

## What Changes

- Permit NAS identity lookup and writes only to the reserved per-AppRole NAS credential publication paths.
- Own the truenas_admins and truenas_user entity groups; NAS account provisioning consumes their membership.
- Reuse existing shared-read AppRole policies.
- Use XCP-ng's own NAS account and remove its cross-account credential policy.

## Capabilities

### New Capabilities

- `runner-cache-credentials`: Scoped NAS account publication and consumption.

### Modified Capabilities

None.

## Impact

Changes Vault source grants only. Shared and AppRole NAS provisioning use the NAS
operator identity; consumer AppRoles and VM provisioning receive no API administrator key.
