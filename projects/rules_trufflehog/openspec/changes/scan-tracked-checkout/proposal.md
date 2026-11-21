## Why

Git history contains pointers to LFS content. A filesystem scan of the tracked
checkout reaches the materialized payloads.

## What Changes

- Add an optional checkout scan to the shared scan rule.
- Fail the checkout scan if tracked LFS files have not been materialized.

## Capabilities

### New Capabilities

- `tracked-checkout-scan` provides scan-only credential detection over the
  tracked filesystem snapshot.

## Impact

Checkout scans use the pinned TruffleHog filesystem source.
