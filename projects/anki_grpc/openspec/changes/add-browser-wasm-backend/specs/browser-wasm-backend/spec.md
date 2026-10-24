## Purpose

Expose Rust Anki commands to browser Go applications through a JavaScript API,
reusing generated protobuf clients and shared dispatch and lifetime handling.

## ADDED Requirements

### Requirement: Direct reusable command boundary

The browser adapter SHALL implement the generated Go clients' unary connection
interface using a JavaScript command API over the existing Rust dispatcher.
It SHALL reuse native direct-client protobuf, call-option and error handling.
Browser operation SHALL require neither a pipe/HTTP2 transport nor JSPI,
shared memory or a nested Wasm interpreter.

#### Scenario: Generated client call

- **WHEN** a caller starts the browser backend and invokes a generated Anki method
- **THEN** the shared dispatcher processes the same protobuf request and response
- **AND** native errors, unknown routes and supported call options retain their existing semantics

### Requirement: Owned command buffers

The JavaScript adapter SHALL copy request and response bytes between independent
memories, refresh views after growth and free every owned buffer on ordinary
completion or error. Empty messages SHALL retain their meaning.

#### Scenario: Repeated large Unicode transfer

- **WHEN** callers repeatedly exchange large and empty protobuf messages
- **THEN** the exact response bytes are retained without unbounded per-call memory growth

### Requirement: Explicit lifecycle and failure handling

The client SHALL permit one startup attempt, reject calls before readiness and
after shutdown, and report startup failures. Cancellation SHALL not release active
state or permit later mutation of caller-owned response objects. Stop SHALL wait
for active work and backend destruction; Done SHALL observe full cleanup even
after a waiting caller's deadline. A Wasm trap SHALL make the instance terminal.

#### Scenario: Cancelled asynchronous startup

- **WHEN** startup is canceled before its factory finishes
- **THEN** a later-created backend is closed and shutdown waits for that cleanup

#### Scenario: Terminal shutdown or trap

- **WHEN** a backend stops or traps
- **THEN** retained clients reject calls and repeated shutdown preserves the cleanup result

#### Scenario: Recoverable command error

- **WHEN** a command returns an ordinary native error or unknown route
- **THEN** its details are reported and a subsequent valid command can succeed

### Requirement: Shared injected collection filesystem

The collection client SHALL accept injected Afero storage for temporary
directories, files and cleanup, with OS storage as the native default.
The browser WASI host SHALL expose that same storage and virtual namespace.

#### Scenario: Package and media round trip

- **WHEN** Go imports a package with binary media and runs a collection operation
- **THEN** Rust accesses the same files and Go reads the completed package with media intact

#### Scenario: Native default

- **WHEN** a native caller omits injected storage
- **THEN** collection preparation and cleanup retain their OS filesystem behavior

### Requirement: WASI correctness and cleanup

The host SHALL preserve required filesystem offsets and binary content, validate
guest memory access and report unsupported operations. It SHALL retain the
collection client's uncertain-cleanup behavior and release owned descriptors.

#### Scenario: Failed collection operation

- **WHEN** an invalid package or filesystem error prevents completion
- **THEN** the error is reported and resources are cleaned up or retained under the existing uncertainty contract

#### Scenario: Native package verification

- **WHEN** the browser exports a completed package
- **THEN** native Anki can verify its contents, scheduling, review history and binary media
