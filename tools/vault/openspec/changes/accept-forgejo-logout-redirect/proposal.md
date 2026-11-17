## Why

The live Forgejo 17 logout endpoint returns HTTP 303. The login wrapper rejects
it before verifying session invalidation, making successful operations fail.

## What Changes

Accept HTTP 303 after logout, retaining the original-cookie invalidation check.

## Capabilities

### New Capabilities

- `forgejo-logout-verification`: Verify session invalidation after supported logout responses.

### Modified Capabilities

None.

## Impact

Forgejo login cleanup and its existing protocol fixture only.
