## Why

Provision a dedicated RouterOS identity and store its generated password in Vault so printer configuration no longer borrows ingress credentials.

## What Changes

- Package and validate the accounts root and AL injection without live authentication; document bootstrap permissions and inspect sensitive-field wiring.

## Capabilities

### New Capabilities

- `printer-credentials`: Printer account provisioning.

### Modified Capabilities

None.

## Impact

Affects infra/printers, infra/mikrotik and infra/vault. The user subsequently
authorized deployment and state adoption after plan review. The
[deployment receipt](../../../../printers/openspec/changes/move-printer1-infrastructure/deployment.md) records the completed operations and limits.
