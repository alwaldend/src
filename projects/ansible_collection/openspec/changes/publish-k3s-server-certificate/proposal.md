## Why

K3s callers need a common way to distribute the public API server CA after
bootstrap.

## What Changes

- Add opt-in Vault publication to the shared K3s role.
- Select the destination and Vault credentials in the calling project.

## Capabilities

### New Capabilities

- `k3s-server-ca-publication` publishes the public K3s CA to Vault KV v2.
