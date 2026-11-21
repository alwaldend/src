## Why

Repository credential scanning covers Git history. Materialized checkout
files also need to be inspected, including Git LFS payloads.

## What Changes

- Run the tracked-checkout scan as a normal repository CI test.

## Capabilities

### Modified Capabilities

- `repository` adds credential scanning of the tracked checkout in CI.

## Impact

The same checkout scan target runs locally and in CI.
