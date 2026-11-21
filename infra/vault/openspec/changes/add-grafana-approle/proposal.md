## Why

The Grafana deployment requires its own Vault and OIDC identity.

## What Changes

Register shared Vault memberships for the Grafana-owned AppRole.

## Capabilities

### New Capabilities

- `grafana-membership`: Shared Vault group membership for Grafana deployments.

### Modified Capabilities

None.

## Impact

Adds shared group membership entries to the Vault root.
