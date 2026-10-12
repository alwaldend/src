---
title: Anki gRPC
description: Rust CLI exposing Anki's backend over stdio, TCP, or Unix sockets
statuses:
  - finished
languages:
  - rust
tags:
  - anki
  - grpc
---

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
- `test/transport/`: CLI and network transport E2E using official Anki messages.
- `rustfmt.toml`: Rust 2021 formatting settings for the project.
- `releases/`: Standard binary metadata for the repository release registry.

## Protocol and ownership

RPC paths use upstream service and method names, such as
`/anki.collection.BackendCollectionService/OpenCollection`. Anki decodes
requests and encodes responses using its own schema bindings. Clients can call
upstream operations immediately after connecting; initialization belongs to the
server and requires no RPC. Operation failures include Anki's serialized
`BackendError` in the binary `anki-error-bin` trailer. The transport uses stable
Tonic 0.14.6.

Dispatch runs on blocking threads because Anki owns its own runtime and
operation synchronization. Worker panics terminate the process. A cancelled
RPC does not promise rollback of a mutation already running in Anki.

The Go consumer is delivered separately in
[PR #146](https://github.com/alwaldend/src/pull/146). It consumes this project's
binary as runtime data and owns pipes, the gRPC readiness deadline, process wait,
shutdown, and crash reporting. It uses generated upstream protobuf messages
with `grpc.ClientConn.Invoke`. For stdio shutdown, close the collection, then
signal EOF on stdin while draining stdout. TCP and Unix clients can use normal
gRPC connections with the same message types and RPC paths.

## Verification

```sh
bazel_agent bazel test //projects/anki_grpc/test/transport:transport_test
```

The server E2E checks CLI help, invalid arguments, startup initialization,
persistence, and reconnection over both TCP and Unix sockets. It emits closed
SQLite collections and JSON verification reports as repeatable artifacts.

The separately owned Go consumer E2E in PR #146 also exercises default stdio,
generated cards, concurrent calls, cancellation, startup deadlines, shutdown,
and process failures. Its test target is supplied by that PR.
