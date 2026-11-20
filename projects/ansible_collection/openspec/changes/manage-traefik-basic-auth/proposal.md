## Why

Install Basic Auth verifier files through the same role that owns the proxy
account, private directories, and service activation.

## What Changes

- Allow the shared Traefik role to install an optional verifier file.
- Include verifier changes in the existing restart and rollback workflow.

## Capabilities

### New Capabilities

- `traefik-basic-auth-verifier`: Manages the proxy's optional verifier file.

## Impact

The VictoriaMetrics playbook uses the shared role to install its Vault-backed
verifier.
