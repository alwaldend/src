## Why

Create the printers AppRole and grant only the credentials required by its DNS and firewall roots.

## What Changes

- Add the printers AppRole and exact-secret read policy; verify policy paths, AppRole membership and Terraform formatting.

## Capabilities

### New Capabilities

- `printers-access`: Printer credential access.

### Modified Capabilities

None.

## Impact

Affects infra/printers, infra/mikrotik and infra/vault. The user subsequently
authorized deployment and state adoption after plan review. The
[deployment receipt](../../../../printers/openspec/changes/move-printer1-infrastructure/deployment.md) records the completed operations and limits.
