# WASM worker Specification

## Purpose

Allow applications to execute an embedded Anki WASM backend without cgo while
retaining the existing supervised worker process and upstream protobuf protocol.

## Requirements

### Requirement: Dependency-isolated server builds

The system SHALL provide separate native and WASM server packages. One public
command target SHALL compose optional engines through independent
`//projects/anki_grpc:cgo` (default true) and `//projects/anki_grpc:wasm`
(default false) booleans, without separate command directories or changing the
host platform. WASM-only executables MUST build with cgo disabled and MUST NOT
link the native Rust static library. Native-only builds SHALL retain their
existing behavior. Application builds with both engines SHALL prefer cgo for
generic strategies and retain qualified WASM strategies; builds with neither engine
SHALL retain explicit external clients.

#### Scenario: Native package migration

- **WHEN** the existing application starts its native worker after the package split
- **THEN** root/serve help, stdio/TCP/Unix transports, Anki RPC responses and native error trailers retain their existing behavior

#### Scenario: WASM executable build

- **WHEN** the WASM application target is built with cgo disabled
- **THEN** it embeds the guest and requires no host native Rust/C linkage or companion worker file

#### Scenario: Adapter-owned guest packaging

- **WHEN** a command imports the WASM server adapter
- **THEN** the adapter owns building and embedding its production guest and exposes `Run(arguments)` without requiring the command to handle module bytes
- **AND** integration fixtures can supply their own guest through the internal host without depending on the production guest

#### Scenario: Shared Rust server CLI

- **WHEN** either backend receives server arguments
- **THEN** the same Rust CLI owns help, subcommands, transport validation and gRPC serving, including duplicate-option rejection
- **AND** help and argument failures do not initialize the Anki backend or require its filesystem/cache inputs
- **AND** the WASI guest executes as a command while Go supplies byte streams without implementing gRPC or an Anki allocation/dispatch ABI

#### Scenario: Composed optional engines

- **WHEN** both cgo and WASM engines are enabled
- **THEN** generic subprocess and in-process selection use cgo, while qualified WASM strategies remain available
- **AND** each compiled engine has its qualified subprocess scheme registered on the executable's scoped factory and an explicit scheme selects that engine for same-executable workers

#### Scenario: External-only application

- **WHEN** both cgo and WASM engines are disabled
- **THEN** explicit external clients remain available without native linkage, while local server hosting and in-process schemes fail clearly

### Requirement: Worker process isolation

The application's default local strategy SHALL use a supervised child process
of the same executable and communicate over gRPC stdio. It SHALL prefer cgo
and otherwise use WASM; explicit subprocess schemes SHALL select their engine
for same-executable workers and preserve independent executables' arguments. Parent
cancellation and bounded shutdown SHALL retain forced termination and reaping
behavior for that strategy.

#### Scenario: Forced worker termination

- **WHEN** a worker cannot finish before the existing shutdown deadline
- **THEN** the parent terminates and reaps the child without hosting backend execution in the parent process

### Requirement: Explicit in-process WASM client

WASM-enabled applications SHALL offer `inprocess+wasm://` using the existing
Rust server and Go gRPC client over in-memory streams without a subprocess,
listening socket or cgo. Each client SHALL own a fresh instance. Generic
`inprocess://` SHALL prefer native when available and otherwise select WASM;
it SHALL be unavailable with neither engine. Startup cancellation SHALL bound
readiness without owning lifetime after success. Stop SHALL reject new calls,
cancel guest execution and release streams/runtime; Done SHALL observe final
cleanup when a caller stops waiting. Failures SHALL never restart the instance.

#### Scenario: Independent embedded collections

- **WHEN** two in-process WASM clients open distinct collections
- **THEN** backend state and shutdown remain independent, large protobuf messages and native error details retain their behavior

#### Scenario: In-process lifecycle failure

- **WHEN** startup is cancelled, compilation fails or Stop interrupts execution
- **THEN** the caller receives a bounded result and final cleanup can be observed without recreating the guest

### Requirement: Protocol and guest failure handling

The WASM worker SHALL accept existing upstream Anki RPC method names and protobuf
messages. It SHALL preserve native backend error details. Invalid guest memory,
unknown RPCs and guest traps MUST return errors without unsafe host memory access
or automatic request replay.

#### Scenario: Invalid invocation

- **WHEN** a request names an unknown method or a guest returns an invalid memory range
- **THEN** the caller receives a bounded error and the host performs no unchecked memory access

### Requirement: Explicit offline guest capabilities

The first WASM backend SHALL support offline collection, package, and media
operations. Network sync, media sync, HTTP-server operations, and other
network-dependent services SHALL fail with an explicit unsupported error.
The cgo backend SHALL retain its full API. Unsupported calls MUST NOT perform
partial work or fall back to native execution.

#### Scenario: Network operation on an offline guest

- **WHEN** a client requests collection sync, media sync, or HTTP-server startup
- **THEN** it receives an unsupported error identifying the offline WASM limitation, and subsequent local collection operations remain usable

### Requirement: Native collection semantics

The WASM worker SHALL use Anki's collection and package implementations for the
application's collection workflows, preserving identities, review history and
media. It MUST NOT silently fall back to cgo.

#### Scenario: Collection package round trip

- **WHEN** a modern collection package containing media and more than 10 MiB of collection data is opened, changed and exported through the WASM worker
- **THEN** native Anki can reopen the resulting artifact with the expected notes, identities, review history and media

### Requirement: Host local-time semantics

The WASM backend SHALL use the host timezone for local-time conversions at the
requested timestamp, preserving historical and daylight-saving offsets as well
as ambiguous and nonexistent local times. Scheduling SHALL preserve the host's
non-UTC collection offset instead of replacing it with UTC. Backup local-time
operations SHALL use the same corrected timezone implementation.

#### Scenario: Non-UTC scheduling

- **WHEN** a New York host computes collection timing through the WASM worker
- **THEN** the collection retains the applicable host offset

#### Scenario: Timezone transition

- **WHEN** the guest converts UTC or local civil timestamps across seasonal, half-hour or skipped-date transitions
- **THEN** its Chrono Local results reflect the host timezone's valid offsets and ambiguity
