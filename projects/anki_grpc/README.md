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
The serve subcommand blocks until shutdown and owns the process's protocol streams. The optional direct client also exposes unary commands through
owned C ABI handles and buffers, while retaining generated Anki messages.

The generic `anki_as_code server run` command forwards arguments to this library.
`anki_as_code server serve` delegates to `server run serve`, selecting the Rust
serve subcommand. Its collection workflows
launch that command in a child process, retaining independent termination and
crash isolation without extracting or finding another executable. The standalone
`anki_grpc` executable remains a thin adapter over the shared Rust library.
Both builds currently target Linux x86_64 and require the repository's native
Rust/C toolchains; linking introduces no new external dependencies.

## Offline WASI server

This mode embeds the Rust CLI in a native Go executable so applications can
avoid cgo and native Rust linkage. The Go executable still targets the host
OS and CPU; only its embedded Rust backend targets WebAssembly. Browser
execution is outside this mode's scope.

The independent boolean `//projects/anki_grpc:wasm` defaults to false and
selects the standalone WASM worker when enabled. Application builds also honor
trunk's `//projects/anki_grpc:cgo` boolean, which defaults to true. Build a
pure-Go application or worker with:

```sh
bazel_agent bazel build -c opt //projects/anki_grpc //projects/anki_as_code --//projects/anki_grpc:cgo=false --//projects/anki_grpc:wasm=true
```

For the application, generic subprocess and in-process selection both prefer
cgo and otherwise use WASM. With neither engine, explicit external
clients remain available and local hosting fails clearly. The command owns a
scoped factory and registers each compiled engine's qualified subprocess
scheme. Those factories select the named engine in the application worker;
independent executables retain their exact arguments and caller-selected engine.
The generic subprocess keeps the preferred default. Generic client and
application libraries remain pure Go.

Guest and host follow the selected Bazel compilation mode. `-c opt` optimizes
both; omitting it retains the normal debug build. The Rust targets use
rules_rust's built-in WASI platform attribute and embed their outputs directly.

`pkg/ankigrpc/wasm` is a separate pure-Go host using wazero v1.12.0. It owns building
and embedding the production guest and exposes `Run(arguments)`, matching the
cgo adapter. Commands handle no guest bytes. `internal/wasmhost` accepts guest
bytes internally for production and integration fixtures and serves the same
opaque protobuf gRPC protocol. The default application strategy supervises a
child process through `ankiclient`. Explicit `inprocess+wasm://` instead runs
the same guest inside the application's process.
The complete Rust CLI owns gRPC framing and Anki dispatch in one guest.
Cancelling command execution closes that instance without replay or replacement. The guest receives filesystem access for the existing
absolute-path protocol. This is process isolation, with the worker's filesystem
authority, rather than a restricted filesystem sandbox.

`wasm.New()` constructs an unstarted managed client. Start executes the embedded
Rust server in a goroutine and supplies a `net.Pipe` dialer to the shared
`ankiclient.ClientGRPC`. Its in-memory stdio carries the existing gRPC protocol;
no process, listening socket or cgo is needed. Startup context controls readiness
only. Each client owns separate backend state. Close collections before Stop;
Stop cancels the guest and releases its streams/runtime, while Done observes
final cleanup if the caller's Stop deadline expires. A failed instance never
reconnects or replays requests. This strategy retains protobuf and HTTP/2 framing.
Executable composition registers `inprocess+wasm://` whenever WASM is enabled;
generic `inprocess://` prefers cgo when available and otherwise selects WASM.

Both backends execute the same Rust CLI and gRPC server. The WASI guest runs
on one execution thread, using an executor to schedule async futures; native
builds retain Tokio. WebAssembly supports multithreading, but this integration
does not. The guest targets [`wasm32-wasip1`](https://doc.rust-lang.org/rustc/platform-support/wasm32-wasip1.html),
where Rust thread spawning is unsupported. The [threaded Rust target](https://doc.rust-lang.org/rustc/platform-support/wasm32-wasip1-threads.html)
requires host support for WASI thread spawning. Wazero v1.12.0 provides
experimental shared-memory and atomic instructions, but no built-in WASI
thread-spawn implementation; its [threading example](https://github.com/wazero/wazero/blob/v1.12.0/experimental/features_example_test.go)
creates execution instances from Go. Our host neither enables that experimental
feature nor implements thread spawning. Changing the Rust target alone would
therefore be insufficient.

The Go host performs blocking I/O in goroutines and supplies
checked byte-stream imports for stdio, TCP and Unix sockets, using bounded I/O
queues and blocking readiness notifications. The command owns stream progress
and closure; pumps send results through channels, and guest tasks share the
host's readiness notification. It implements no command parser, gRPC server,
protobuf dispatcher or Anki allocation/free API. One guest retains
the backend across all client connections. Help and invalid arguments run before
backend initialization; they do not require collection or cache access.

The guest supports offline collection, package, and local media operations.
Network sync, media sync, remote media URLs, AnkiHub, AnkiWeb, GitHub downloads,
and Anki's HTTP sync-server entry points return explicit unsupported errors.
The host's gRPC stdio/TCP/Unix transports remain available. The cgo backend
retains its full API. Callers must keep exclusive offline access to collection
and media paths; the WASI SQLite VFS does not coordinate locks with other
processes. The guest uses the host-supplied `TMPDIR`; backups and package
compression run synchronously without guest threads.
Chrono local-time lookups use the Go host timezone at the requested timestamp,
including historical offsets, DST gaps and overlaps. The host uses `TZ` or its
system timezone and includes standard Go timezone data for named-zone fallback.

Set `ANKI_WASM_CACHE_DIR` to a caller-owned directory to reuse wazero's compiled
code across worker processes. Fresh directories are created with private
permissions; protect an existing cache from changes by other users. The cache
is optional and contains no required runtime inputs. Without it, each worker
compiles its embedded guest at startup.

The application target `//projects/anki_as_code` with the same backend setting uses the same
worker model inside one installed executable. Its shared entry point hosts either
compiled adapter; generic strategies prefer cgo, while qualified schemes select
their engine explicitly. Both command names stay unchanged. The Go executable
retains the host OS/CPU platform; only the embedded Rust guest targets WASI. The following checks exercise the actual Anki guest,
its unsupported HTTP sync-server entry points, command failures, process
transports, and the pinned SQLite build:

```sh
bazel_agent bazel test //projects/anki_grpc/internal/wasmhost:backend_test \
  //projects/anki_grpc/test/wasm:inprocess_test \
  //projects/anki_grpc/test/wasm:offline_test \
  //projects/anki_grpc/test/wasm:http_server_test \
  //projects/anki_grpc/test/wasm:worker_test \
  //projects/anki_grpc/test/wasm:sqlite_test
```

The worker check copies its executable alone and verifies stdio/TCP/Unix gRPC,
error trailers, and shutdown. The SQLite check retains `wasi-sqlite.anki2` with
an 11 MiB payload and `sqlite-verification.txt` in its test outputs; native and
WASI SQLite independently reopen it. The application `wasm_workflow_test`
retains `large-field-roundtrip.colpkg` with a 17 MiB note field and checks its
identity, scheduling, review history and media through native Anki after export.

### Runtime comparison

This snapshot was checked on 21 October 2026. Approximate GitHub star counts
refer to the runtime repository, not its smaller Go wrapper. The cgo column
describes the standard in-process Go binding linked in that row.

| Runtime (stars)                                                          | Guest multithreading and Go integration                                                                                                                                                                                                                                                                                                                                               | Uses cgo?                                                                                                |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| [Wasmer](https://github.com/wasmerio/wasmer) (~21.2k)                    | Modern Wasmer supports [WASI/WASIX threading](https://docs.wasmer.io/runtime/wasix/). The published Go wrapper still depends on [Wasmer 2.0](https://github.com/wasmerio/wasmer-go/blob/master/Cargo.toml) and does not expose the modern engine's full feature set.                                                                                                                  | [Yes](https://github.com/wasmerio/wasmer-go#supported-platforms)                                         |
| [Wasmtime](https://github.com/bytecodealliance/wasmtime) (~18.7k)        | The engine implements [WASI thread spawning](https://docs.rs/wasmtime-wasi-threads/latest/wasmtime_wasi_threads/). The Go binding exposes [thread instructions](https://github.com/bytecodealliance/wasmtime-go/blob/main/config_feat_threads.go); running Rust-created WASI threads requires additional host integration.                                                            | [Yes](https://github.com/bytecodealliance/wasmtime-go#installation)                                      |
| [WasmEdge](https://github.com/WasmEdge/WasmEdge) (~10.8k)                | Core threads are supported; the [example uses host-created workers](https://github.com/WasmEdge/WasmEdge/blob/master/examples/capi/mandelbrot-set-in-threads/README.md) sharing guest memory. The [Go SDK exposes the threads feature](https://github.com/second-state/WasmEdge-go/blob/master/wasmedge/configure.go); a ready-made Rust/WASI thread-spawn path is not verified here. | [Yes](https://github.com/second-state/WasmEdge-go/blob/master/wasmedge/cgo.go)                           |
| [Wasm3](https://github.com/wasm3/wasm3) (~8.1k)                          | No guest multithreading; the [feature list](https://github.com/wasm3/wasm3#features) includes atomics only. Its linked [community Go wrapper](https://github.com/matiasinsaurralde/go-wasm3) is marked work in progress.                                                                                                                                                              | [Yes](https://github.com/matiasinsaurralde/go-wasm3/blob/master/wasm3.go)                                |
| [wazero](https://github.com/wazero/wazero) (~6.4k; our pinned v1.12.0)   | Experimental shared memory and atomics, without built-in WASI thread spawning. Its [example creates workers from Go](https://github.com/wazero/wazero/blob/v1.12.0/experimental/features_example_test.go); guest multithreading requires custom host support.                                                                                                                         | [No](https://github.com/wazero/wazero/blob/v1.12.0/README.md)                                            |
| [WAMR](https://github.com/wasm-micro-runtime/wasm-micro-runtime) (~6.2k) | Supports [pthread and WASI thread spawning](https://github.com/wasm-micro-runtime/wasm-micro-runtime/blob/main/doc/pthread_impls.md). Build the native library with `WAMR_BUILD_LIB_WASI_THREADS=1` for the WASI threads ABI.                                                                                                                                                         | [Yes](https://github.com/wasm-micro-runtime/wasm-micro-runtime/blob/main/language-bindings/go/README.md) |

Core thread instructions, guest thread spawning and Go API support are separate
capabilities. This comparison checks upstream source and documentation; it is
not an Anki compatibility test.

## Project layout

- `cmd/anki_grpc/main.rs`: Standalone executable adapter.
- `pkg/ankiserver/lib.rs`: Transport listeners, opaque protobuf codec, and Anki dispatch.
- `pkg/ankiserver/cli.rs`: Shared Rust CLI, help and transport validation.
- `internal/backend/`: Shared Rust backend dispatch and descriptor-generated command lookup.
- `pkg/ankigrpc/native/`: Optional native C ABI bridge, CLI hosting and client registration.
- `pkg/ankigrpc/`: Pure Go connection-string routing and client construction.
- `pkg/ankigrpc/wasm/`: Production WASI command build, embedding and argument-vector API.
- `internal/wasmhost/`: Pure-Go WASI execution, native byte streams and timezone callbacks.
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

Native dispatch runs on blocking threads because Anki owns its own runtime and
operation synchronization. WASI serializes dispatch inside its single-thread
command. Worker panics terminate the process; WASM traps fail that command. A cancelled
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

The executable creates a configured factory and registers its worker schemes:

```go
factory := ankigrpc.NewFactory(ankigrpc.SubprocessAddress(workerPath, "serve"))
factory.RegisterClientFactory("subprocess", ankigrpc.ToClientFactory(ankigrpc.NewSubprocessClient))
```

`NewSubprocessClient` returns the concrete stdio client and requires an executable
in the URL. Use `ankigrpc.SubprocessAddress(binary, arguments...)` to preserve
literal paths and argv. Arguments are passed exactly, including zero arguments.
The executable owns supported worker schemes, entrypoints, server arguments and
its configured default. NewFactory snapshots built-in and opt-in package registrations;
later registrations remain private to that factory. Explicit addresses override the
default for one construction only. Package-level registration remains available for
opt-in adapters, and package-level NewClient retains no default. Importing the
native package does not register subprocess workers.

`ClientFactory` is `func(string) (ankiclient.StartableClient, error)`.
Pass a concrete Factory's `NewClient` method value when injecting its registry and
default selection. `ToClientFactory` adapts a typed constructor with contextual
error wrapping; a nil constructor produces a nil function. Registration rejects
nil factories. The factory returns the constructed client directly. Scheme
dispatch and consumer injection use the same three-method contract. Registrations retain
constructors, never clients, and reject duplicate or invalid schemes. Unavailable
schemes fail without fallback. Native workers receive no injected platform flags.

`WithClient(ctx, factory, address, operation)` constructs through the supplied
factory (or `NewClient` when nil), starts the client, lends its ready connection
and stops it with a five-second budget independent of request cancellation. It
opens and reserves no collection. Callers own native collection open/close,
archive publication and source reservations.

Import the optional `ankigrpc/native` package for direct client registration and
native CLI hosting. It registers only `inprocess` and `inprocess+cgo`; each label
creates fresh state. Start initializes Anki synchronously, checking its context
before and after; cancellation waits for initialization and failure cleanup.
Command cancellation stops waiting without interrupting Rust, and Stop rejects new calls while draining
active work. Its deadline bounds waiting; the concrete client's Done method
observes final cleanup. Native
crashes affect the hosting process, and streaming and unsupported call options
fail explicitly. Native errors retain the transport's error-decoding contract.

Core and application imports remain pure Go. Managed collection paths use the
selected client regardless of transport. Client implementations manage their own
resource lifetimes; callers do not query backend ownership. Private schemes and
injected constructor functions work without global registration.

Routing identifies only the scheme and passes the selected address unchanged.
The transport constructor parses it once. Network clients retain decoded endpoints,
including raw and escaped IPv6 zones and relative Unix socket paths containing spaces.

`pkg/ankiclient.Client` exposes `Connection` and `Stop` for callers that use and
release an owned connection. `StartableClient` embeds Client and adds `Start`
for callers that control startup. Prepared application collection handles use
the shared Client interface.
Concrete clients own their transport readiness and cleanup; constructors perform
no startup I/O. Collection operations use generated native service clients.
Concrete owned transports expose Done for native work that can outlive a
cancelled Stop wait. Completion observation adds no public client interface.

`ClientSTD` implements StartableClient using the server's stdin/stdout. `ankiclient.NewSTD` accepts
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
worker, err := factory.NewClient("")
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

When using raw native OpenCollection as above, close that storage through
`collections.CloseCollection` before calling `worker.Stop(ctx)`. Native full-package
export consumes collection storage, so callers must avoid closing it twice.
The process transport signals EOF while draining stdout, waits for exit and
forces termination when its own shutdown budget expires. It completes reaping
and retains unexpected process-exit errors on repeated Stop calls.

Stop is safe before Start and after startup failure. Each client permits one
Start attempt; use a new instance for a new transport. Connection rejects access
before readiness, while stopping and after transport exit. Retained generated
clients reject calls once their transport stops. Shut down the owning client
through Stop. Callers requiring final native completion after a cleanup error
can observe Done on the concrete owned transport.

The client library supplies no server executable or implicit runtime dependency.
Callers package their configured executable; native tests declare their own fixtures.

### Existing gRPC server

`ankiclient.NewGRPC(endpoint)` validates the
endpoint and returns `(*ClientGRPC, error)` without contacting the server. It implements
`StartableClient` for a running server. Endpoints accept
TCP host:port, `tcp://host:port`, `unix:///absolute/path` and absolute socket paths.
Unix URI paths may also be relative to the client's working directory.
Start waits for readiness within its context. Stop closes only the client;
it does not stop the remote server or close its collection. Retained clients
reject calls after Stop. Native error trailers are decoded as for stdio.

`ankiclient.NewGRPCWithDialer(dial)` provides the same lifecycle and RPC handling
for a caller-supplied `net.Conn`. The WASM adapter supplies a single-use memory
stream and retains ownership of guest execution. Other dialers may reconnect;
neither client replays failed RPCs.

### Application scopes

Applications own collection paths, archive restore/publication and source
reservations. They use the client connection for generated native RPCs and stop
the transport after collection cleanup. The
[Anki as Code collection store](../anki_as_code/README.md#collection-lifecycle)
provides managed callback scopes and already-open application handles.

The Go package is `ankiclient`; `anki_grpc` is the project and server name.
Concrete client types include `ClientSTD` and `ClientGRPC`.

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

The collection-lifecycle E2E verifies one client reusing its transport, direct database
persistence and media layout, native package exports/restores, temporary cleanup,
persistent destinations, retained-connection rejection, raw access between collections and failed or
cancelled lifecycle calls. It emits a native database/package and verification report.
It also blocks a real native media download, cancels the client, and verifies
that shutdown waits for active callers and retains unconfirmed restores until
backend shutdown.

Application scope verification belongs to the Anki as code collection-store E2E.
It exercises real workers across callback success, failure, cancellation, panic
and crashes, checking process reaping, source reuse and archive preservation.
