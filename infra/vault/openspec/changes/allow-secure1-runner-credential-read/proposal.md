## Why

The development deployment must read the already-created secure1 worker identity
to enable its runner without creating another registration.

## What Changes

- Grant the development AppRole read access to that exact KV data path.
- Preserve existing SSH and DNS policies and grant no writes or sibling access.

## Capabilities

### New Capabilities

- `secure1-runner-access`: deployment-only read access to cached credentials.

### Modified Capabilities

None.

## Impact

One Vault policy and its development AppRole binding. No new dependencies.
