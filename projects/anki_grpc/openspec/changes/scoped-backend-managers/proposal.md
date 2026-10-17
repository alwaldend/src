## Why

Callers currently compose process and collection backends and coordinate their cleanup manually. Add reusable scoped resource managers so collection users can borrow prepared backends while the managers retain lifecycle ownership and transport selection remains replaceable.

## What Changes

- Add Manager.WithBackend and ManagerSTD to create, start, lend and stop a fresh stdin/stdout backend for each scope.
- Add ManagerCollection.WithCollection to borrow a transport manager, open a collection through BackendCollection, and close it before the transport scope ends.
- Add private withBackend as the shared startup/callback/cleanup implementation for both managers, with an explicit cleanup timeout.
- Preserve callback and cleanup errors, cancellation causes, native readiness, package restoration and direct database behavior.
- Keep the existing Backend, BackendSTD and BackendCollection APIs and family-first names. Add native E2E coverage and lifecycle examples without changing anki_as_code consumers in this PR.

## Capabilities

### New Capabilities

- backend-managers: Scoped backend lifecycle and transport-independent collection composition.

### Modified Capabilities

None.

## Impact

Changes belong to projects/anki_grpc/pkg/backend and the project's documentation/OpenSpec workspace. Reuse existing native backends, Go standard library and pinned dependencies. This additive master-based PR does not change source reservations, archive publication policy, or the other pending PRs.
