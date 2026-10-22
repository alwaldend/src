## Why

The collection store already owns each source scope and always selects the stdin/stdout worker. Backend manager interfaces add another ownership layer without production pooling or transport selection.

## What Changes

- Construct the worker and collection backends directly in the store.
- Keep bounded cleanup, source reservations, prepared connections and archive publication in the same application scope.
- Remove manager injection and retain worker configuration privately for process-boundary fixtures.
- Preserve native lifecycle coverage using real subprocesses and observable native/config artifacts.

## Capabilities

### New Capabilities

- `collection-lifecycle`: Application-owned native scopes with deterministic cleanup and source publication.

### Modified Capabilities

None.

## Impact

Changes internal collection code, its native E2E and current docs. The paired anki_grpc change removes the unused manager API. No dependencies or CLI syntax change.
