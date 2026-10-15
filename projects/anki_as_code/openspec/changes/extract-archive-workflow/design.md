## Context

Anki exposes ImportCollectionPackage on a closed backend and
ExportCollectionPackage on an open backend. Export consumes the open collection.
The existing Go backend owns process startup, shutdown, and reaping.

## Decisions

- Constructors accept an already-running gRPC connection with no open collection.
  The caller owns its process and connection, keeps them alive, and serializes
  collection use until Source.Close. Archive code never starts or stops a worker.
- Native filesystem RPCs finish before scratch cleanup: cancellation is checked
  between stages, while an already-started native call is awaited synchronously.
  Closing uses an uncancelled context. The caller owns transport/process failures.
- OpenSource accepts an archive or offline SQLite file; OpenPackage requires an
  archive. Raw files are copied as opaque bytes and journal sidecars are rejected.
- Native import restores into private paths. Native export includes media.
  Archive details and supported versions belong to Anki; arbitrary extra ZIP
  members and ZIP comments are outside Anki's collection-package contract.
- Publish consumes its source even on failure. Destination-side staging, mode,
  sync, optional digest guard, and rename are shared by both source formats.
- Cancellation does not guarantee rollback. Once native work starts, wait for its
  response before removing its files; do not detach cleanup from an active RPC.
- Keep FileHash and explicit mode/digest because #128 chooses replacement policy.
  Remove DatabasePath: consumers use the backend connection and Anki API.

### Failure cases recorded before implementation

- Native import/export loses media, note identity, review history, or custom state.
- Import is called with an open collection, or export is followed by double close.
- Failed or cancelled import/open/publication stops a borrowed worker or leaks scratch.
- Cancellation removes files while the native worker may still be writing them.
- A previously started worker cannot be reused after success, failure, or cleanup.
- A changed target is replaced, failed output changes an original, or staging leaks.
- A closed or consumed source publishes again; repeated cleanup fails.
- An offline source with WAL/journal is copied, or raw output cannot cross directories.
- Requested mode is lost or opening mutates the original source.

## Risks / Trade-offs

Digest verification and rename are separate operations. Callers serialize other
writers. Raw source publication handles the database only; media remains external
as with an offline .anki2 source. Package publication uses Anki's supported content.

## Migration

Deliver against master. #128 supplies its existing connection after merge; this PR does not publish #128 changes.
