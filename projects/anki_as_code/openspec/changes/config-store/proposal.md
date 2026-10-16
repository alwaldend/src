## Why

Configuration operations should share a root-bound store, just as native resources and state have their own stores. Package-level functions currently require callers to pass the same filesystem path repeatedly.

## What Changes

- Replace package-level Load, Dump, and Export with methods on config.Store, constructed with config.New(root).
- Migrate existing CLI, application, exporter, and E2E callers to the store.
- Document the API and update operation diagrams without changing serialization or filesystem behavior.

## Capabilities

No spec-level behavior changes; this is an internal API refactor with skip_specs enabled.

## Impact

Only projects/anki_as_code changes. The independent PR targets master and needs no new dependency. The reconciliation branch will adopt the API after this prerequisite merges.
