## Why

Move printer1 DNS and network access to the printers project while preserving the deployed objects.

## What Changes

- Move canonical DNS and firewall configuration, wire separate backends and dedicated credentials, and verify formatting, provider schema and DNS ownership; document state transfer and acceptance.

## Capabilities

### New Capabilities

- `printer-access`: Printer network configuration.

### Modified Capabilities

None.

## Impact

Affects infra/printers, infra/mikrotik and infra/vault. The user subsequently
authorized deployment and state adoption after plan review. The
[deployment receipt](deployment.md) records the completed operations and limits.
