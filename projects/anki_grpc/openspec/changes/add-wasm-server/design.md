## Context

The native adapter, approved runtime/sysroot and guest dependency adaptations
landed in #182 and #183. Both command variants retain the same-executable
supervised worker and upstream protobuf/gRPC protocol. See proposal.md for
motivation and specs/wasm-worker/spec.md for the behavioral contract.

## Goals / Non-Goals

Execute the complete existing Rust CLI in WASI while keeping the native cgo
adapter separate. Preserve collection semantics, process supervision and
stdio/TCP/Unix transports. Upstream network services retain the approved
offline-first errors; browser execution remains outside this change.

## Decisions

- `pkg/ankigrpc/wasm` owns a single Rust command build and embedding and exposes
  `Run(arguments)`. Application commands never handle guest bytes. Built-in
  rules_rust platform selection targets WASI only for the guest; Go stays native
  and pure. Normal Bazel compilation mode controls host and guest optimization.
- Preserve trunk's independent cgo engine flag and add an opt-in WASM flag.
  The application command owns optional imports, scoped factories and worker
  schemes. Both generic subprocess and in-process selection prefer cgo and
  otherwise use WASM. Explicit qualified subprocess schemes select
  any compiled engine for same-executable workers; the executable owns their
  `--server-engine` selector. Independent workers retain their exact arguments.
  The generic subprocess parser exposes validated options so executable
  composition can adapt its own workers without duplicating URL parsing.
  With neither engine, external clients remain
  usable. Fixed native/WASM/combined/external-only fixtures verify all combinations
  independently of public build flags and keep interoperability reopen native.
- `internal/backend` owns the shared `anki_backend_dispatch` library and its
  command-ID generator for native direct clients and both Rust server platforms.
  `pkg/ankigrpc/native` and `pkg/ankigrpc/wasm` own optional engine integration;
  the WASM runtime remains in `internal/wasmhost` for production and fixture guests.
- The shared Rust server library owns CLI parsing, gRPC framing and native
  error trailers, delegating Anki commands to the shared dispatch library. Its CLI is a private module, avoiding a
  single-consumer library target. A deliberately reusable dispatcher interface
  lets the executed fixture use exactly the same command/server implementation.
- Native builds keep Tokio and blocking backend tasks. WASI uses a task list
  and one shared readiness waker, with no waiter registry or guest threads.
  Guest task wakes and host I/O signal the same retained host notification;
  the executor polls and waits without a separate guest wake flag.
  Every host notification polls the current tasks; this favors simple ownership
  for the embedded worker over selective scheduling for many concurrent clients.
  Do not enter futures::executor: Anki's local media implementation invokes its
  own block_on, which cannot nest inside LocalPool's executor guard.
- Go supplies checked byte-stream imports for stdio and native TCP/Unix sockets,
  plus the existing host timezone callbacks. The host contains no command
  parser, gRPC server, protobuf dispatcher or Anki allocation/free ABI. WASI
  Preview 1 cannot create native sockets, so this platform adapter remains.
  Bounded pumps and blocking readiness notifications keep both stream directions
  live under backpressure without idle polling. Closing a command joins its
  pumps and releases every retained peer/listener.
  The command goroutine owns handle and stream progress. Pumps return results
  through channels; one outstanding write needs no shared mutex or completion
  counter. Removing a handle before closure provides the sole close guard.
  Unix stdio uses private nonblocking descriptor duplicates registered with Go's
  poller so closure interrupts inherited-pipe reads/writes before joining pumps.
  The worker reserves these streams; duplicates share nonblocking flags with
  their original process descriptors.
- One Rust backend is initialized after argument validation and shared across
  clients. Help does not initialize Anki or open a protocol transport. Optional
  compilation-cache initialization errors are deferred until serving so CLI
  help and invalid arguments do not require writable cache inputs.
- A cancelled command context closes executing WASM, including CPU-bound code.
  The host never restarts a module or replays a request. Parent termination and
  bounded reaping stay owned by the existing native process supervisor.
- Filesystem access retains the worker's existing absolute-path authority.
  The stock WASI SDK supplies SQLite libc through the existing LLVM toolchain.
  Host timezone callbacks preserve historical offsets, gaps and overlaps.

## Alternatives and verdict

The optional `inprocess+wasm://` client owns one embedded command and a
`net.Pipe` connection. It supplies a single-use dialer to `ankiclient.ClientGRPC`,
which owns gRPC setup, readiness, connection lending and native error decoding.
The WASM client owns guest execution and stream cleanup. The existing Rust
gRPC server executes in a goroutine;
Go supplies the same checked host streams. Startup context controls readiness,
not the instance lifetime after startup. Stop rejects new calls, cancels the
guest, and joins runtime/stream cleanup; Done observes final cleanup if the
Stop deadline expires. Unexpected guest exit is retained and never replayed.
The executable registers this scheme only when WASM is enabled. Generic
`inprocess://` prefers native when cgo is available and otherwise selects WASM;
with neither engine it is unavailable. The default remains the supervised worker.
Keep the direct cgo client: a native socket-pair transport adds platform-specific
plumbing and HTTP/2 overhead without being required for this feature.

Proceed with full Rust command execution and a stream-only platform adapter.
The old backend ABI worked but duplicated protocol serving and buffer ownership
in Go. Running the native Tokio path unchanged requires thread spawning and
socket APIs absent from this WASI host integration. See the project's
[runtime threading documentation](../../../README.md#offline-wasi-server) for
the selected target and runtime capabilities; this is not a general WebAssembly
threading restriction. A new runtime or a private framed-RPC protocol would expand
scope without removing the shared-server requirement.

## Risks / Trade-offs

- Protocol deadlock or lost connection state: execute actual stdio/TCP/Unix
  clients, large protobuf exchanges, reconnects and native error trailers.
- Guest traps/cancellation or host resource leaks: interrupt CPU-bound commands,
  close peer streams, join pumps and never replay against a replacement guest.
- Guest memory ranges: validate every host I/O buffer and copy queued writes
  before the synchronous import returns.
- Nested executor use: exercise real local file URL media import and backups;
  passing command help alone cannot establish Anki storage behavior.
- Compilation cost: a complete CLI compiles more code than the retired parsing
  guest. Optional caller-owned compilation caching amortizes repeated starts.
  Cache directories remain caller-protected and newly created privately.

## Acceptance and migration

Replace the backend ABI and separate parsing command with one full WASI CLI.
Use the direct command check's pre-fix failure as evidence of the original gap.
Run native and WASI collection/package/media workflows and reopen exported
artifacts through native Anki. Review the complete candidate for correctness
and concrete simplifications before required delivery gates and publication.
Git and ignored delivery receipts own exact revisions and publication state.
