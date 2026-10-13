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
    MediaDbPath:     "/path/to/collection.media.db",
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
