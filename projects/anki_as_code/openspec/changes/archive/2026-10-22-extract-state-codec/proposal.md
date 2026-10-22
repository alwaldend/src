## Why

PR #128 mixes collection reconciliation with independent state and codec package changes. Extract the shared changes so they can be reviewed and merged separately.

## What Changes

- Move the unchanged generic TOML/JSON protobuf codec from internal/config/codec to internal/codec and update existing callers.
- Accept Resources protobuf input data in Store.Preview; keep state moves, ownership validation and persisted writes behind Store with detached snapshots.
- Document package ownership and validate existing state/config/application/CLI workflows without changing their behavior.

## Capabilities

### New Capabilities

None. This is a behavior-preserving extraction with skip_specs enabled.

### Modified Capabilities

None.

## Impact

Only the Anki project's shared state helper, codec location, imports, Bazel wiring, documentation and this change record. No dependencies, schemas, collection lifecycle, reconciliation or CLI options change. Both this PR and #128 target master; #128 retains required shared changes until this prerequisite merges.
