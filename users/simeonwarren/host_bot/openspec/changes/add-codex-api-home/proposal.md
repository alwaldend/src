## Why

Provide a separate Codex home for API-key login without changing the existing home.

## What Changes

- Manage `~/.codex-api` with the shared Codex settings and instructions.
- Require API login in this home and document login and selection.

## Capabilities

### New Capabilities

- `codex-api-home`: Isolated API-authenticated Codex home.

### Modified Capabilities

None.

## Impact

Host Bot Ansible and documentation only; no deployment or new dependencies.
