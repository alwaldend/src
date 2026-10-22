## Why

Callers duplicate collection import, open/close and temporary-file lifecycle around
the Go backend. A collection-scoped Backend implementation can own those resources
while borrowing an already-started server.

## What Changes

- Add BackendCollection and BackendCollectionOptions in the existing backend package.
- Reuse Start, Connection and Stop for collection lifetime; leave wrapped backend lifetime to callers.
- Open database inputs directly and restore package inputs through Anki into a temporary directory by default.
- Allow an explicit caller-owned restore directory and report retained restore paths in cleanup errors.
- Preserve generated service clients, handle consuming package export, and document resource ownership.
- Deliver this library independently; update anki_as_code to consume it in its existing PRs.

## Capabilities

### New Capabilities

- collection-backend: Collection-scoped connections over a caller-owned prepared backend.

### Modified Capabilities

None.

## Impact

projects/anki_grpc/pkg/backend, its README and native E2E, and visibility for an
existing upstream service binding. Reuse current Anki schemas and dependencies.
No server/Rust build or scheduling changes.
