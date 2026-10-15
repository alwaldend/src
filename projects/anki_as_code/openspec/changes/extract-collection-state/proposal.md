## Why

Collection-backed state needs a small API independent of reconciliation. State
commands use the existing managed backend, rather than configuring their own
gRPC connection or owning a second implementation of collection lifecycle.

## What Changes

- Add typed State.Resources maps and internal/state Dump, Import, Update and Replace.
- Keep the state store independent of process, collection lifecycle and files.
- Compose BackendSTD and BackendCollection in App.Start, close them in App.Stop,
  and let WithState borrow their prepared connection across operations.
- Accept an existing .anki2 database through --input. Archive persistence and
  resource reconciliation remain outside this PR.
- Keep --state-key namespaces and batched path/ID imports with resource-first validation.
- Share TOML/JSON codecs using a typed Format with named TOML and JSON values.

## Capabilities

### New Capabilities

- collection-state: Explicit stored-state operations using a borrowed connection.

### Modified Capabilities

None.

## Impact

projects/anki_as_code schemas, state package, application wiring, state CLI,
codec callers and focused E2E/docs. Reuse the existing backend library and its
native collection lifecycle. No new dependencies, SQL, archive implementation
or backend lifecycle changes.
