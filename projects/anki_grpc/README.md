---
title: Anki gRPC
description: Rust CLI exposing Anki's backend over stdio, TCP, or Unix sockets
statuses:
  - finished
languages:
  - rust
  - go
tags:
  - anki
  - grpc
---

![Anki gRPC architecture](/docs/projects/anki_grpc/assets/architecture.svg)

Anki gRPC is a Rust CLI that forwards upstream protobuf commands to Anki's
native backend. It initializes one backend at startup using Anki's default
options. Anki owns collection semantics; the caller owns process lifetime and
crash handling. Stdin/stdout is the default transport. TCP and Unix sockets are
also supported. Diagnostics use stderr.

## Build

From the repository root:

```sh
bazel_agent bazel build //projects/anki_grpc:anki_grpc //projects/anki_grpc:clippy
```

[The third-party Anki package](../../third_party/com_github_ankitects_anki/README.md)
owns and builds the upstream native library, source pins, schemas, and single
Cargo dependency lock. Its native build consumes no project implementation
sources or build rules. `projects/anki_grpc` builds our server executable and
depends on `//third_party/com_github_ankitects_anki:native`. Both use standard
`rules_rust` targets. Compilation is offline and sandboxed; supported native
builds target Linux x86_64.

## Usage

The built executable accepts:

```sh
anki_grpc --help
anki_grpc
anki_grpc --transport stdio
anki_grpc --transport tcp --listen 127.0.0.1:50051
anki_grpc --transport unix --listen ./anki.sock
```

Stdio serves one gRPC HTTP/2 connection on the process pipes, reserving stdout
for protocol traffic. TCP and Unix listeners accept connections until the
process exits. They require an explicit listen address or socket path; stdio
rejects `--listen`. Listener startup prints `Listening on tcp://ADDRESS` or
`Listening on unix://PATH` to stderr, including the allocated port when TCP
port zero is requested. An existing Unix socket path causes startup to fail;
the caller owns removal of its socket after the server exits.

Network transports use plaintext gRPC. All connections share one initialized
backend and its open collection. Close a collection through Anki before
reopening it or shutting down the process. Disconnecting a network client
leaves the backend available to later clients.

## Project layout

- `cmd/anki_grpc/main.rs`: Transport listeners, opaque protobuf codec, and Anki dispatch.
- `cmd/anki_grpc/cli.rs`: CLI help and transport argument validation.
- `cmd/anki_grpc/build.rs`: Command lookup generated from Anki's own descriptors.
- `cmd/anki_grpc/BUILD.bazel`: Executable, descriptor generator, and Clippy targets.
- `pkg/backend/`: Reusable Go owner of the server process and gRPC connection.
- `docs/architecture.mmd`: Client, server, and storage diagram source.
- `test/transport/`: CLI and network transport E2E using official Anki messages.
- `rustfmt.toml`: Rust 2021 formatting settings for the project.
- `releases/`: Standard binary metadata for the repository release registry.

## Protocol and ownership

RPC paths use upstream service and method names, such as
`/anki.collection.BackendCollectionService/OpenCollection` and
`/anki.notes.NotesService/GetNote`. Anki decodes
requests and encodes responses using its own schema bindings. Clients can call
upstream operations immediately after connecting; initialization belongs to the
server and requires no RPC. Operation failures include Anki's serialized
`BackendError` in the binary `anki-error-bin` trailer. The transport uses stable
Tonic 0.14.6.

Dispatch runs on blocking threads because Anki owns its own runtime and
operation synchronization. Worker panics terminate the process. A cancelled
RPC does not promise rollback of a mutation already running in Anki.

Anki implicitly forwards collection-service methods through its backend services.
The descriptor generator exposes both original service paths for generated clients
and backend-prefixed paths for existing callers, using Anki's own method indices.
Methods declared in both services must have matching request and response schemas.
TCP and Unix clients can use normal gRPC connections with the same generated clients.

## Go backend

Import `git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend`.
`Backend` exposes `Start(context.Context) error`,
`Connection() (grpc.ClientConnInterface, error)`, and `Stop(context.Context) error`.
`BackendSTD` implements it using the server's stdin/stdout. `backend.New` accepts
one `BackendSTDOptions` value with `Binary`, `Arguments`, and `ShutdownTimeout`.
It copies the arguments without starting a process. A zero shutdown timeout uses
five seconds; Start rejects negative timeouts. `Start`
launches the server and creates its gRPC connection, waiting for readiness. Its
context controls startup only: cancelling it after successful startup leaves the
worker running. The child inherits the parent
environment and stderr. The Bazel library includes the server executable as
runtime data. Failed startup releases its resources and reaps the child.
The backend owns a separate lifetime context for the command.

The backend exposes `grpc.ClientConnInterface` through `Connection`. Construct
Anki's generated service clients with it; each method accepts a typed request and
returns a typed response. There are no handwritten per-operation wrappers.

```go
var worker backend.Backend = backend.New(backend.BackendSTDOptions{
    Binary:          workerPath,
    Arguments:       []string{"--transport", "stdio"},
    ShutdownTimeout: 5 * time.Second,
})
if err := worker.Start(ctx); err != nil {
    return fmt.Errorf("start Anki: %w", err)
}
connection, err := worker.Connection()
if err != nil {
    return fmt.Errorf("get Anki connection: %w", err)
}
collections := collection.NewBackendCollectionServiceClient(connection)
items := notes.NewNotesServiceClient(connection)

_, err = collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
    CollectionPath:  "/path/to/collection.anki2",
    MediaFolderPath: "/path/to/collection.media",
    MediaDbPath:     "/path/to/collection.media.db2",
})
if err != nil {
    return fmt.Errorf("open collection: %w", err)
}
note, err := items.GetNote(ctx, &notes.NoteId{Nid: noteID})
```

Generated clients use the upstream Go packages, such as
`github.com/ankitects/anki/proto/anki/collection` and
`github.com/ankitects/anki/proto/anki/notes`. Use `BackendCollectionServiceClient`
for collection open/close and `NotesServiceClient`, `NotetypesServiceClient`,
`DecksServiceClient`, or `CardsServiceClient` for collection operations.

One connection interceptor rejects calls after shutdown and decodes Anki's error
trailers, preserving gRPC status, context causes, and caller call options. Failed
calls are never automatically replayed or sent to a replacement process.
Cancellation stops waiting without guaranteeing rollback of a native mutation.

Close an open collection through `collections.CloseCollection` before calling
`worker.Stop(ctx)`. Stop signals EOF on stdin while draining stdout, waits for exit,
and cancels the command's lifetime context when its shutdown context expires or
after the configured shutdown timeout. `exec.CommandContext` handles forced
termination; an earlier Stop context deadline takes precedence.
It completes reaping even after cancellation, returning the context cause and
process exit errors. It also reaps crashes;
unexpected exits remain errors on repeated Stop calls. Stop is safe before Start
and after startup failure. Each instance permits one Start attempt; use a new
instance for a new process. Connection returns an error before readiness, during
startup/shutdown, after Stop, or after worker exit. Construct service clients after
Start succeeds. Retained clients reject calls
once Stop begins. The backend owns the connection: shut it down through Stop.

### Collection backend

`BackendCollection` implements the same interface for collection lifetime. It
borrows an already-started backend and never starts or stops that backend:

```go
collectionBackend := backend.NewCollection(backend.BackendCollectionOptions{
    Backend: worker, // already started; no collection open
    Path:    "/path/to/collection.colpkg", // or collection.anki2
})
if err := collectionBackend.Start(ctx); err != nil {
    return fmt.Errorf("prepare collection: %w", err)
}
connection, err := collectionBackend.Connection()
if err != nil {
    return fmt.Errorf("get collection connection: %w", err)
}
items := notes.NewNotesServiceClient(connection)
```

The options object is copied. `Path` accepts `.anki2` and `.colpkg` extensions,
case-insensitively. Databases open directly; Anki can create a missing database.
Media uses the database stem plus `.media`, with a `.media.db2` index. Database
writes commit in place and Stop does not roll them back.

Packages restore through Anki into an owned temporary directory by default.
Set `RestoreDirectory` to use an empty persistent destination instead; it can be
created when missing and is retained on Stop. `RestoreDirectory` is invalid for a
database input. The source package is not modified. The backend API remains
Start, Connection and Stop; restore paths are supplied through options.

After work, call `collectionBackend.Stop(cleanupContext)` and then
`worker.Stop(cleanupContext)`, preserving both errors. A clean collection Stop
closes native storage and deletes only its owned temporary restore. Call Stop
even after Start fails. Startup leaves cleanup to the caller and adds no separate
shutdown timeout. The borrowed backend stays available for another collection.
Retained generated Anki clients reject after Stop. Start permits one attempt;
Stop is terminal and idempotent, including before Start and after failed preparation. Lifecycle-changing native
calls must use the collection connection, and callers must give the wrapper
exclusive collection ownership. Anki exposes unary RPCs; a small connection adapter
tracks active calls and observes CloseCollection and ExportCollectionPackage
while forwarding generated client calls over the borrowed connection. These operations close the collection
without closing the transport. Export also consumes it on a native failure, so
confirmed completion lets Stop clean up without attempting a second close.
Transport failures or cancellation without a native response retain the files.

Cancellation can stop waiting while native work continues. If native completion
cannot be established, Stop reports failure and retains restore files instead
of deleting files potentially still in use. Stop the borrowed backend first;
after it is reaped, retained restore files can be removed using the location
reported in the cleanup error. For programmatic cleanup, pass a caller-owned
`RestoreDirectory` through options and remove it only after backend shutdown.
Explicit destinations always remain caller-owned. Callers retain
backend shutdown responsibility even when collection startup fails.

Ordinary calls remain concurrent. Stop waits for registered Go calls within its
context before closing storage. Any unconfirmed native response invalidates
collection access and keeps cleanup uncertain, including after later calls succeed.
Anki's collection mutex protects running native work; it cannot acknowledge a
cancelled call still queued in Rust. Stop therefore retains the restore and reports
failure in that case; stop the borrowed backend before cleaning its temporary files.

### Application scopes

Applications compose the process and collection owners. Start the worker before
the collection, close the collection before stopping the worker, and attempt
cleanup after failed startup, callback errors, cancellation or panic unwinding.
Use bounded cleanup contexts independent of operation cancellation and preserve
both operation and cleanup errors.

The [Anki as code collection store](../anki_as_code/README.md#collection-lifecycle)
owns its callback scope, source reservations and archive publication. The backend
library supplies the explicit native owners and their existing readiness,
shutdown and restore-retention contracts.

Public backend names are `BackendSTD` and `BackendCollection`. The Go package
is `backend`; `anki_grpc` is the project and server name.

## Verification

```sh
bazel_agent bazel test //projects/anki_grpc/test/transport:transport_test
bazel_agent bazel test //projects/anki_grpc/pkg/backend:backend_test
```

The server E2E checks CLI help, invalid arguments, startup initialization,
persistence, and reconnection over both TCP and Unix sockets. It emits closed
SQLite collections and JSON verification reports as repeatable artifacts.

The Go backend E2E uses generated clients over default stdio and checks native
persistence, generated cards, original service routes, overlapping service methods,
caller error trailers, concurrent calls, cancellation, startup deadlines, shutdown,
and process failures. It also verifies constructor configuration, argument copying,
repeated startup, stop before start, invalid arguments, startup-context isolation,
connection errors during startup/shutdown, Stop deadlines, and OS-observed child
reaping. Shutdown cases cover deadlines, explicit cancellation, the default
five-second grace period, and a configured grace period. They also verify invalid
timeout rejection and native failure details. The forwarding fixture reaps its
inner Anki process before testing forced shutdown of the outer worker. It emits `collection.anki2`
and JSON verification reports.

The collection-backend E2E verifies prepared-backend borrowing, direct database
persistence and media layout, native package exports/restores, temporary cleanup,
persistent destinations, retained-client rejection, backend reuse and failed or
cancelled lifecycle calls. It emits a native database/package and verification report.
It also blocks a real native media download, cancels the client, and verifies
that shutdown waits for active callers and retains unconfirmed restores until
backend shutdown.

Application scope verification belongs to the Anki as code collection-store E2E.
It exercises real workers across callback success, failure, cancellation, panic
and crashes, checking process reaping, source reuse and archive preservation.
