## Why

API-created Forgejo runners already have a UUID and token. The shared role must
install that identity without requesting a second registration.

## What Changes

- Accept optional cached API credentials for initial runner-file installation.
- Validate identity fields, protect credentials and preserve existing files.
- Expose installation and service phases so callers can prepare prerequisites
  before a fresh worker starts accepting jobs.

## Capabilities

### New Capabilities

- `forgejo-runner-credentials`: import an existing worker identity.

### Modified Capabilities

None.

## Impact

The existing runner role and its local acceptance fixture. No new dependencies.
