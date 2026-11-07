## Why

Wi-Fi clients need a stable local DNS name and network access to the wired
printer at `192.168.1.221`. The recorded router policy isolates the wireless
subnet from wired devices.

## What Changes

- Declare `printer1.dc1.alwaldend.com` in both local and global DNS views.
- Add a separate Terraform firewall root permitting Wi-Fi clients to reach
  only this printer, before the existing terminal forwarding drop.
- Keep historical router exports and the existing DNS-only root unchanged.
- Document deployment credentials, verification, and the limits of routed access.
- Keep the firewall lookup inline, without a module introduced solely for mocks.

## Capabilities

### New Capabilities

- `wifi-printer-access`: Global and local printer DNS and a host-specific routed Wi-Fi exception.

### Modified Capabilities

None.

## Impact

The firewall root reuses the pinned RouterOS provider and owner AppRole, with
separate state and scoped access to the existing non-DNS ingress RouterOS
credential. Its paired Vault policy change is maintained in
`infra/vault/openspec/changes/allow-mikrotik-firewall-credential`. No new dependency,
account, ingress service, print application, multicast relay, or host configuration
is introduced.
