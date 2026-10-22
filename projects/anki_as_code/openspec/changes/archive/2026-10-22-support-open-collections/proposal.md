## Why

App's collection operations always open a path, preventing callers from reusing an already open collection. Shared Go API contracts should remain independent of native backend construction so callers and CLI can use the application interface without importing its implementation.

## What Changes

- Add a public `ankiascodeapi` package for the application interface and path and borrowed-connection collection sources.
- **BREAKING**: Accept a collection source in App's existing collection operations; keep archive build's base path and CLI flags unchanged.
- Reuse caller-owned connections without opening, closing, or publishing their collections. Preserve managed path lifecycle and publication behavior.
- **BREAKING**: Rename `api/collection` and its Go package to `ankiascodepb`, and rename the protobuf namespace to `anki_as_code`.
- Have CLI consume the public application interface and verify concrete App compatibility through command composition.

## Capabilities

### New Capabilities

- `collection-source-api`: Shared Go source contracts, per-call borrowing, and ownership rules.

### Modified Capabilities

None.

## Impact

Changes are confined to Anki as code's API, App, CLI, generated-target declarations, direct consumers, and documentation. Existing protobuf field numbers, message names, state keys, and serialized config fields remain unchanged. Go import paths and protobuf fully qualified type names change intentionally. No external dependencies are added.
