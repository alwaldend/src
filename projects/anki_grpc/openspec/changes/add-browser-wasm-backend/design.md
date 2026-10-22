## Context

See proposal.md. The current browser backend runs the Rust WASI CLI over a Go
pipe. The repository already has a transport-independent protobuf dispatcher
and a native Go direct connection; these supply the replacement's reusable core.

## Goals / Non-Goals

**Goals:** Preserve generated clients, command semantics, error details and
explicit lifetime ownership while removing browser pipe/HTTP2/JSPI machinery.
Continue using the same Afero files for Go and Rust.

**Non-Goals:** Shared memory, another interpreter, new dependencies, rewriting
Anki operations, durable live SQLite storage or changing the standalone WASI CLI.

## Decisions

- Proceed with a JavaScript command boundary. Retaining the current pipe would
  reuse the complete artifact but also retain stream scheduling and framing.
  A per-method JavaScript API would duplicate schema bindings. One command
  method reuses the existing method-path/protobuf dispatcher instead.
- Extract native Go unary encoding, call-option handling, error decoding and
  lifecycle into `internal/directclient`. Native and browser clients supply
  a backend with Command and Close methods. Preserve native public interfaces,
  error decoding, concurrent progress queries and generated-client behavior.
  The native server implementation remains private, as on the current base.
- Share native backend ABI implementation under `internal/backend`. Build a
  browser WASI reactor with explicit initialization, command and destruction
  exports plus typed buffer helpers. Keep CLI startup in its native owner;
  the browser artifact does not link the gRPC server.
  Rust 1.98 omits the reactor startup object for this library build, so select
  `crt1-reactor.o` from that same pinned toolchain and export `_initialize`.
  Its constructor reference prevents per-command constructor/destructor
  wrappers. Remove the explicit linkage when the pinned compiler supplies it.
- `server.js` wraps an initialized reactor as `{command, close}`.
  Command accepts a method path and protobuf bytes, returning a status code and
  copied response bytes. Codes distinguish success, encoded native errors and
  unknown routes. JavaScript owns temporary Rust buffers through each call and
  uses typed descriptors rather than guessing struct-by-value ABI lowering.
- `wasi.js` instantiates and initializes the reactor using the existing
  synchronous Afero filesystem/timezone host. Go constructs that host and
  destroys Rust state before releasing descriptors and callbacks. There is no
  server loop, pipe, HTTP/2 handshake or JSPI suspension.
- Register initialization and calls before shutdown can race with them.
  Cancellation ends waiting without freeing active state or mutating a returned
  response later. Stop rejects new calls and Done observes complete cleanup,
  including handles produced by late asynchronous initialization.
- Synchronous Rust commands execute inside the worker. They cannot be preempted
  while blocking its event loop; cancellation is observed when execution can
  progress. Preserve the application's ordered queue. Host filesystem callbacks
  remain synchronous and memory-only during nested Go/JavaScript/Rust calls.
- Ordinary native errors and unknown commands remain recoverable. Wasm traps or
  broken bridge responses make that instance terminal; do not invoke unsafe
  cleanup exports on poisoned Rust state. Release the instance and host resources.
- Keep injected collection filesystem behavior and the existing uncertain-error
  cleanup contract. Reserve /tmp for working storage; the app persists completed
  packages and editable files after collection scopes finish.

## Risks / Trade-offs

- Separate Wasm memories require byte copies; performance is secondary to reuse.
- Memory growth invalidates views: acquire current memory views after allocations
  and exports, and copy results before freeing their Rust owners.
- Rust traps can leave state inconsistent: reject subsequent calls and retain the
  error while allowing a new independently owned backend.
- WASI reactor exports/toolchain behavior needs build verification: inspect the
  compiled imports/exports and execute real generated calls before app integration.
- Native lifecycle extraction can regress cancellation or options: retain native
  E2E coverage and validate async startup, Stop overlap and retained connections.

## Acceptance Cases Before Implementation

Exercise Unicode, empty and large protobufs, repeated allocation/free, unknown
routes, native error trailers, supported and unsupported call options, canceled
requests, startup rejection, cancellation during delayed startup, late-handle
cleanup, terminal traps and idempotent Stop. Emit hashes and a structured report.
A failed or canceled request must not later mutate the caller's reply.

Run real collection/package/media workflows without JSPI and retain the complete
browser save/reload, conflicting-writer and malformed-input coverage. Reopen the
browser-produced package in native Anki and verify identities, scheduling, review
history and binary media.

## Migration Plan

Replace the browser transport in the existing PR and package its reactor under
the existing server.wasm resource name. Update documentation and diagram to the
command boundary. The native and standalone WASI CLI paths remain independently
validated. Rollback is a feature-branch revert; no persistent schema migration or
live deployment is involved.
