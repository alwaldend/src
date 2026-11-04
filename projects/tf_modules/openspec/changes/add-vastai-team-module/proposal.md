## Why

Vast.ai bootstrap must support another team without copying a fixed resource graph.

## What Changes

- Add a reusable module creating one zero-credit team and a map of project API keys.
- Return key values only through sensitive module outputs for direct secret-resource use.
- Keep Vault authentication and destination paths in the calling root.

## Capabilities

### New Capabilities

- `vastai-team-module`: Declarative team and per-project key composition.

## Impact

Uses the in-tree provider and existing Terraform tooling, with no new external dependency.
