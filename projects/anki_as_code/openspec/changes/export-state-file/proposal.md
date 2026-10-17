## Why

Export knows the native identity of every generated config, but onboarding currently requires manually binding each path. An optional state file makes whole-collection onboarding repeatable without changing the source collection.

## What Changes

- Add `col export --state-output FILE` to write version-1 TOML bindings alongside configs from the same snapshot.
- Include exported decks, notes, note types, member identities and referenced preset aliases.
- Preserve config-only export and keep installation explicit through `col state update` or `replace`.
- **BREAKING**: App `ColExport` accepts an explicit state-output argument, empty for config-only behavior.
- Reject state destinations that alias the source or collide with generated configs.

## Capabilities

### New Capabilities

- `export-state-file`: Optional complete bindings for exported configurations.

### Modified Capabilities

None.

## Impact

Collection export projection, App, CLI, native E2E coverage and project documentation. No new dependencies.
