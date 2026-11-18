## Why

Agent commits need consistent co-author attribution drawn from the operator's
managed Git configuration.

## What Changes

- Read an optional co-author identity during commit preparation.
- Preserve existing attribution and keep the disclaimer at the end of the
  trailer block.
- Document the same setting for Forgejo delivery.

## Capabilities

### New Capabilities

- `agent-commit-attribution` covers the configured co-author trailer.

### Modified Capabilities

## Impact

AL delivery and its Forgejo compatibility procedure.
