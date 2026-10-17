## Why

Archive builds need staged copying and publication independently of reconciliation.
Extract that policy from #128 so it can be reviewed and merged separately.

## What Changes

- Add `internal/rebuild.Builder` with an injected apply operation.
- Publish only a completed staged `.colpkg`, preserving the base and destination on failure.
- Verify the package independently with native Anki archive/state operations.
- Keep App and reconciliation integration in #128.

## Capabilities

### New Capabilities

- `archive-build`: Apply changes to a private archive copy and publish it atomically.

### Modified Capabilities

None.

## Impact

Adds a standard-library-only internal Go package, its native integration test,
documentation and narrow test visibility. No dependency additions or CLI changes.
