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

Build an optimized worker with `bazel_agent bazel build -c opt //projects/anki_grpc`.
Optimized builds strip the final Rust executable's debugging information and static
symbols; development builds retain them.

[The third-party Anki package](../../third_party/com_github_ankitects_anki/README.md)
owns and builds the upstream native library, source pins, schemas, and single
Cargo dependency lock. Its native build consumes no project implementation
sources or build rules. `projects/anki_grpc` builds our server executable and
depends on `//third_party/com_github_ankitects_anki:native`. Both use standard
`rules_rust` targets. Compilation is offline and sandboxed; supported native
builds target Linux x86_64.

## Usage

The built executable lists commands when invoked without arguments or with
`--help`. The serve subcommand starts the server:

```sh
anki_grpc --help
anki_grpc serve
anki_grpc serve --transport stdio
anki_grpc serve --transport tcp --listen 127.0.0.1:50051
anki_grpc serve --transport unix --listen ./anki.sock
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

## Linked server

`pkg/ankigrpc/native` exposes the same Rust entry point to Go through cgo and one C ABI
argument-vector function. `native.RunServer(arguments)` takes a full argument vector,
including the program name, executes the native CLI and reports native failures.
The serve subcommand blocks until shutdown and owns the process's protocol streams. There are no per-RPC FFI
bindings: clients continue to use ordinary gRPC and generated Anki messages.

The generic `anki_as_code server run` command forwards arguments to this library.
`anki_as_code server serve` delegates to `server run serve`, selecting the Rust
serve subcommand. Its collection workflows
launch that command in a child process, retaining independent termination and
crash isolation without extracting or finding another executable. The standalone
`anki_grpc` executable remains a thin adapter over the shared Rust library.
Both builds currently target Linux x86_64 and require the repository's native
Rust/C toolchains; linking introduces no new external dependencies.

## Backend build preparation

The public command targets retain their executable names. Build-selected command
source files opt into `pkg/ankigrpc/native` for native hosting and scheme registration. The independent boolean
`//projects/anki_grpc:cgo` defaults to true; set it to false to build the application
command without a linked engine. Core client and application libraries stay pure Go
even when cgo support is enabled for the executable. Optional native support lives
in `pkg/ankigrpc/native`; importing the generic parent never links it. WASM runtime
support remains deferred.

The repository also pins wazero v1.12.0 and the WASI SDK 34.0 sysroot for
storage build validation. The Rust WASI probe uses the built-in `platform` attribute, preserving the
native consumer platform. Run the standalone
SQLite acceptance fixture with:

```sh
bazel_agent bazel test //projects/anki_grpc/test/wasm:sqlite_test
```

It writes an 11 MiB payload, verifies rollback and invalid destinations, and
reopens the database through both native and WASI SQLite. Its repeatable test
outputs include `wasi-sqlite.anki2` and `sqlite-verification.txt`. This build
preparation includes no WASM Anki server or guest port.

## Project layout

- `cmd/anki_grpc/main.rs`: Standalone executable adapter.
- `pkg/ankiserver/lib.rs`: Transport listeners, opaque protobuf codec, and Anki dispatch.
- `pkg/ankiserver/cli.rs`: Platform-independent `server_cli` library for help and transport validation.
- `pkg/ankiserver/build.rs`: Command lookup generated from Anki's own descriptors.
- `pkg/ankigrpc/native/`: Optional native C ABI bridge, CLI hosting and client registration.
- `pkg/ankigrpc/`: Pure Go connection-string routing and client construction.
- `pkg/ankiclient/`: Reusable Go owner of the server process and gRPC connection.
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

## Go client

Import `git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc` for
`NewClient(address)`. It constructs an unstarted lifecycle owner. TCP and Unix
endpoints are built in; other schemes require explicit registration. Empty library
addresses are rejected because the library owns no default worker entrypoint.

The executable registers its worker schemes through the common API:

```go
factory := ankigrpc.ClientFactoryFunc[*ankiclient.ClientSTD](ankigrpc.NewSubprocessClient)
ankigrpc.RegisterClientFactory("subprocess", factory)
```

`NewSubprocessClient` returns the concrete stdio client and requires an executable
in the URL. Use `ankigrpc.SubprocessAddress(binary, arguments...)` to preserve
literal paths and argv. Arguments are passed exactly, including zero arguments.
The executable owns supported worker schemes, entrypoints, server arguments and
any empty-operation default supplied by its injected factory. Importing the
native package does not register subprocess workers.

`ClientFactory` exposes `NewClient(address string) (ankiclient.Client, error)`.
`ClientFactoryFunc[T]` adapts a constructor returning any concrete client T to that
runtime interface. Fixed constructors retain concrete results; scheme dispatch
and consumer injection use the shared lifecycle contract. Registrations retain
constructors, never clients, and reject duplicate or invalid schemes. Unavailable
schemes fail without fallback. Native workers receive no injected platform flags.

`WithClient(ctx, factory, address, operation)` constructs through the supplied
factory (or `NewClient` when nil), starts the client, lends its ready connection
and stops it with a five-second budget independent of request cancellation. It
opens and reserves no collection. Collection publication and reservations remain
with the application's collection store.

`pkg/ankiclient` supplies the shared lifecycle implementations and collection scopes.
`Client` describes ownership and cleanup; RPC access remains the existing gRPC
connection interface, which can be shared by generated service clients.
`Client` exposes `Start(context.Context) error`,
`Connection() (grpc.ClientConnInterface, error)`, and `Stop(context.Context) error`.
`ClientSTD` implements it using the server's stdin/stdout. `ankiclient.NewSTD` accepts
one `ClientSTDOptions` value with `Binary`, `Arguments`, and `ShutdownTimeout`.
It copies the arguments without starting a process. A zero shutdown timeout uses
five seconds; Start rejects negative timeouts. `Start`
launches the server and creates its gRPC connection, waiting for readiness. Its
context controls startup only: cancelling it after successful startup leaves the
worker running. The child inherits the parent
environment and stderr. The caller supplies the exact server executable. Callers declare the executable as a runtime input. Failed startup releases its
resources and reaps the child.
The process client owns a separate lifetime context for the command.

The client exposes `grpc.ClientConnInterface` through `Connection`. Construct
Anki's generated service clients with it; each method accepts a typed request and
returns a typed response. There are no handwritten per-operation wrappers.

```go
worker, err := ankigrpc.NewClient(ankigrpc.SubprocessAddress(workerPath, "serve", "--transport", "stdio"))
if err != nil {
    return fmt.Errorf("construct Anki client: %w", err)
}
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

The client library supplies no server executable or implicit runtime dependency.
Callers package their configured executable; native tests declare their own fixtures.

### Existing gRPC server

`ankiclient.NewGRPC(ankiclient.ClientGRPCOptions{Address: endpoint})` implements the
same Start, Connection and Stop interface for a running server. Endpoints accept
TCP host:port, `tcp://host:port`, `unix:///absolute/path` and absolute socket paths.
Unix URI paths may also be relative to the client's working directory.
Start waits for readiness within its context. Stop closes only the client;
it does not stop the remote server or close its collection. Retained clients
reject calls after Stop. Native error trailers are decoded as for stdio.

### Collection backend

`ClientCollection` implements the same interface for collection lifetime. It
borrows an already-started backend and never starts or stops that backend:

```go
collectionBackend := ankiclient.NewCollection(ankiclient.ClientCollectionOptions{
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
database input. Source packages remain unchanged by default. Set `WriteBack: true`
in constructor options to publish the source package with media on Stop, regardless
of earlier request outcomes. This does not block mutating RPCs or roll
back direct database writes. The backend API remains Start, Connection and Stop.

After work, call `collectionBackend.Stop(cleanupContext)` and then
`worker.Stop(cleanupContext)`, preserving both errors. A clean collection Stop
publishes writable packages or closes native database storage, then deletes only
its owned temporary restore. Publication failures retain the restored data and
report its path for recovery. Call Stop
even after Start fails. Startup leaves cleanup to the caller and adds no separate
shutdown timeout. The borrowed backend stays available for another collection.
Retained generated Anki clients reject after Stop. Start permits one attempt;
Stop is terminal and idempotent, including before Start and after failed preparation. Lifecycle-changing native
calls must use the collection connection, and callers must give the wrapper
exclusive collection ownership. Anki exposes unary RPCs; a small connection adapter
tracks active calls and observes CloseCollection and ExportCollectionPackage
while forwarding generated client calls over the borrowed connection. These operations close the collection
without closing the transport. Export also consumes it on a native failure, so
confirmed completion lets Stop avoid a second close. When writable package
storage was consumed, Stop reopens its restored database before source publication.
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
owns its callback scope and source reservations. The backend library owns archive
write-back, readiness, shutdown and restore retention. Operation errors do not
suppress enabled write-back; callers preserve operation and cleanup errors together.

Public client implementation names are `ClientSTD`, `ClientGRPC` and `ClientCollection`. The Go package
is `ankiclient`; `anki_grpc` is the project and server name.

## Verification

```sh
bazel_agent bazel test //projects/anki_grpc/test/transport:transport_test
bazel_agent bazel test //projects/anki_grpc/pkg/ankiclient:ankiclient_test
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
