## Why

The user requested the existing HTTP proxy for secure1 CI. The development
deployment identity currently cannot read that proxy's Vault-owned password.

## What Changes

Grant `src_infra_dev` read access to the exact existing HTTP proxy secret path,
alongside its existing runner identity read. No writes, listing, proxy-account
creation or secret-value changes.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This wires an existing deployment identity to an existing secret.

## Impact

`approle_src_infra_dev.tf`; implementation and acceptance are linked from
`infra/dev/openspec/changes/configure-runner-proxy`. Applying policy requires
separate live-operation authorization.
