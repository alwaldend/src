## Why

Requiring shared generated DNS pages to match every declaration change makes
independent DNS PRs edit the same files and creates avoidable merge conflicts.
The user requested explicit manual regeneration.

## What Changes

- Remove snapshot freshness from ordinary repository DNS validation.
- Keep declaration ownership/schema checks and renderer fixture coverage.
- Retain explicit `dump --write` and `dump --check` commands and document that
  checked-in pages can lag behind their authoritative declarations.

## Capabilities

### Modified Capabilities

- `infra-dns`: Refresh declaration pages manually without blocking DNS PRs.

## Impact

Only offline checks and documentation change. No DNS provider or live record
is modified, and no generated page refresh is included in this change.
