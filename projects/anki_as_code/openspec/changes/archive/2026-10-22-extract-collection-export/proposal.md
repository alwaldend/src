## Why

Collection export is read-only and independent of reconciliation. Extract a prepared-connection library so export can be reviewed and reused independently.

## What Changes

- Construct Exporter from existing resource and state stores.
- Read native protobuf values directly through generated Anki service clients.
- Export readable identity-free TOML configs only, preserving existing naming while using direct filesystem writes.
- Add native E2E coverage without SQL, process ownership in export, or CLI lifecycle changes.

## Capabilities

### New Capabilities

- collection-export: Read a prepared collection and publish editable configurations without native writes.

### Modified Capabilities

None.

## Impact

Add internal/export, reuse existing config/resource/state and upstream protobuf dependencies, and document their APIs. CLI integration stays in #128. Both PRs target master; #128 retains the prerequisite until this extraction merges. No dependency or schema changes.
