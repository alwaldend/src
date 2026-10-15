## Context

The existing Backend owns process/transport lifetime. Native collection/package
services already provide storage operations. The user explicitly chooses a
collection wrapper borrowing an already-started Backend; callers stop both.

## Goals / Non-Goals

Own collection readiness, native open/close and package restore files. Retain
typed generated clients and avoid manual archives, SQL and process management.
Do not start, stop or replace the borrowed backend, normalize data, or reconcile configs.

## Decisions

- NewCollection accepts BackendCollectionOptions: Backend, Path and RestoreDirectory.
  Path names an .anki2 database or .colpkg package. Resolve paths at Start.
- Start requires a prepared backend, imports packages before native OpenCollection,
  and exposes a connection only after successful open. Database paths are opened
  directly; missing database files can be created by Anki.
- An omitted package RestoreDirectory creates an owned temporary directory. An explicit
  empty destination is persistent and remains caller-owned. Use Anki's standard
  collection.anki2, collection.media and collection.media.db2 layout.
- Stop closes the collection and removes only an owned temporary restore after
  confirmed native completion. It never stops the borrowed backend. Stop is terminal
  and idempotent, including before Start; each wrapper permits one Start attempt.
- Lend a small connection adapter for Anki's unary methods. It checks readiness
  and observes native CloseCollection and ExportCollectionPackage; an export
  consumes the collection even when it returns a native failure. Confirm completion
  from the native error trailer, distinguishing it from uncertain transport errors.
  Register every unary call under the lifecycle mutex; Stop waits for the active
  callers within its context. Unconfirmed completion is sticky and invalidates
  access. Keep ordinary calls concurrent and delegate the rest of the transport
  interface unchanged. Anki close is not
  idempotent, so this tracking prevents a second close during Stop.
- Leave cleanup after startup failure to the caller's explicit Stop. Avoid a
  hidden shutdown pass with its own context and timeout.
- If native completion cannot be established after cancellation/failure, retain
  restore files instead of deleting files potentially still in use. Cleanup errors
  report the restore location; programmatic callers can own an explicit option
  destination and remove it after stopping their backend.
- Keep the library in an independent master-based PR. anki_as_code will borrow its
  connection, stop the collection before the backend and clean owned restores only
  after process reaping on failure. Its CLI workflows remain outside this PR.

Alternatives: making this wrapper own the process would hide ordering but conflicts
with the final ownership choice; leaving collection lifetime in anki_as_code
duplicates a reusable native API concern. Proceed with the caller-owned backend.

## Risks / Trade-offs

- Callers must stop collection then backend; document this order and exercise real
  native persistence and reuse of the underlying backend after collection Stop.
- Cancelled RPCs can continue inside Anki; preserve files when cleanup is uncertain
  and report its location for cleanup after the caller reaps the backend.
- Concurrent use of lifecycle-changing native methods violates exclusive collection
  ownership; serialize consuming operations with Stop and reject later clients.
- Database updates commit in place. Stop cannot roll back native mutations.

Startup rechecks terminal shutdown after acquiring the serialized operation
token, before directory preparation or RPCs. A successful Stop that wins this
token cannot be followed by resource acquisition. Native E2E controls the
scheduling point through its startup context and verifies backend reuse.

Ordinary collection-service methods dispatch through Anki's with_col helper,
which holds its collection mutex. CloseCollection takes the same mutex, protecting
work already running. It cannot acknowledge canceled native jobs still queued by
the gRPC server's spawn_blocking call. Track all Go calls and retain sticky
uncertainty for every unconfirmed response, even if a subsequent call succeeds.
Stop first rejects new callers, waits for active Go calls with its context, and
only closes/removes storage if no call is uncertain. Retain ordinary-call
concurrency so native progress and abort RPCs can still be used during work.
Native media E2E blocks real work and verifies both active-call shutdown deadlines
and retention after a cancelled call, then reaps the backend before file cleanup.
