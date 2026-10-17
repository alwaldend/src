## Why

Application workflows need shared source reservations and successful archive publication. Keep these policies in collection.Store while backend.ManagerCollection owns native lifetime and lends a prepared connection.

## What Changes

- Keep internal/collection.Store for source reservations and successful archive publication.
- Pass a prepared grpc.ClientConnInterface directly to each WithCollection callback.
- Construct resource/state stores and exporters where used; callers reuse local instances for multiple operations.
- Adapt the existing master application workflow without changing its CLI or implementing placeholder commands.
- Replace App.WithCollection with typed state workflows; keep command flags and serialization in the CLI.
- Use the merged backend managers for native lifetime, with an optional injected transport manager.
- Remove the Collection service container and namespace-indexed caches.
- Keep reconciliation and the pending CLI workflow changes in their existing PRs.

## Capabilities

### New Capabilities

- collection-store: Source reservations and publication around managed native collection access.

### Modified Capabilities

None.

## Impact

Changes belong to projects/anki_as_code, with a visibility entry in projects/anki_grpc for native collection tests. Reuse existing Anki, Go, protobuf and TOML dependencies. Keep the new PR against master and separate from the pending workflow implementations.
