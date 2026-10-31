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

#### Scenario: Concurrent generated client calls

- **WHEN** Go callers issue collection commands concurrently during synchronous host calls
- **THEN** the browser adapter SHALL serialize command entry into each Rust reactor
- **AND** each successful write SHALL retain its own data across backend replacement

### Requirement: Owned command buffers

The JavaScript adapter SHALL copy request and response bytes between independent
memories, refresh views after growth and free every owned buffer on ordinary
completion or error. Empty messages SHALL retain their meaning.

#### Scenario: Repeated large Unicode transfer

- **WHEN** callers repeatedly exchange large and empty protobuf messages
- **THEN** the exact response bytes are retained without unbounded per-call memory growth

### Requirement: Explicit lifecycle and failure handling

The client SHALL permit one startup attempt, reject calls before readiness and
after shutdown, and report startup failures. Start SHALL check its context before
and after initialization and wait for the factory to finish even if cancellation
occurs during initialization. A canceled startup SHALL NOT publish its backend.
Command cancellation SHALL not release active state or permit later mutation of
caller-owned response objects. Stop SHALL wait
for active work and backend destruction; Done SHALL observe full cleanup even
after a waiting caller's deadline. A Wasm trap SHALL make the instance terminal.

#### Scenario: Cancelled asynchronous startup

- **WHEN** startup is canceled before its factory finishes
- **THEN** Start SHALL remain pending until the factory finishes and return the context error
- **AND** the completed backend SHALL be closed without becoming available to callers

#### Scenario: Concurrent Stop during initialization

- **WHEN** Stop exhausts its waiting context while startup's factory remains unresolved
- **THEN** Stop SHALL return its context error while Start remains pending
- **AND** Done SHALL remain open until initialization and backend cleanup finish

#### Scenario: Terminal shutdown or trap

- **WHEN** a backend stops or traps
- **THEN** retained clients reject calls and repeated shutdown preserves the cleanup result

#### Scenario: Recoverable command error

- **WHEN** a command returns an ordinary native error or unknown route
- **THEN** its details are reported and a subsequent valid command can succeed

### Requirement: Shared browser working filesystem

The browser WASI host SHALL use the caller's Node-compatible synchronous
JavaScript filesystem directly for collection and media operations, with an
absolute provider root mounted at guest /. Filesystem operations SHALL NOT call
Go. NewClient SHALL accept only a factory; its separate argument SHALL supply
timezone callbacks. The caller SHALL own provider/root alignment and backing
storage. Initialization failure, close and trap SHALL attempt to close every
provider descriptor owned by the host without destroying borrowed storage.

#### Scenario: Go and Rust share provider storage

- **WHEN** the caller installs the Go filesystem facade before Go.run and selects the same root through NewOsFs/BasePathFs
- **THEN** Go and Rust SHALL observe the same provider files
- **AND** the facade SHALL preserve provider receivers, runtime stdout/stderr, fixed cwd / and absolute path resolution
- **AND** the provider and runtime globals SHALL remain stable while Go runs

#### Scenario: Host descriptor cleanup

- **WHEN** initialization fails, the reactor traps or the backend closes
- **THEN** every provider descriptor opened by that host SHALL receive one close attempt
- **AND** the host SHALL invalidate each guest mapping without retrying uncertain provider close failures
- **AND** the caller's stored files and filesystem provider SHALL remain usable

#### Scenario: Collection and media reopen

- **WHEN** generated clients create a collection, write a note and media, close it and stop the backend
- **THEN** a new backend using the same provider and root SHALL reopen the collection and read the same note
- **AND** Go SHALL read unchanged binary media from that provider
- **AND** Rust SHALL observe media files created by Go

#### Scenario: Recoverable filesystem failure

- **WHEN** a collection path cannot be opened because its parent is a file
- **THEN** the command SHALL return a native error
- **AND** the backend SHALL remain usable for a valid collection path

### Requirement: Filesystem progress and recoverable failures

The host SHALL preserve completed transfer counts when a later recognized
provider read or write fails. It SHALL return errno for zero-progress failures
and keep unknown programming exceptions terminal. Sync failures after bytes
are written SHALL remain observable without replaying those bytes.

#### Scenario: Failure after partial transfer

- **WHEN** a recognized provider error follows completed chunks or iovecs
- **THEN** the host SHALL return success with the completed byte count and corresponding position
- **AND** resuming the remaining transfer SHALL preserve exact bytes without duplication or omission

#### Scenario: Sync fails after write progress

- **WHEN** synchronization fails after a positive write count
- **THEN** the host SHALL return that completed count and stop the transfer
- **AND** the next write, sync, datasync or close SHALL report the retained error before further writes
- **AND** close SHALL attempt to release its provider descriptor even when reporting that error

#### Scenario: Provider close also fails after a deferred sync error

- **WHEN** provider close fails while a recognized sync error is retained
- **THEN** the close error SHALL take precedence over the sync error
- **AND** a recognized close errno SHALL leave the backend usable, while an unknown exception SHALL terminate it
- **AND** the guest descriptor SHALL be invalid and the host SHALL NOT retry provider close, regardless of whether the failed call released it

#### Scenario: Final symlink lookup

- **WHEN** an open without symlink-follow requests creation or truncation through a final symlink
- **THEN** the host SHALL reject the open without changing or creating its target
- **AND** explicit follow SHALL retain ordinary provider open/create behavior
- **AND** the host SHALL use the native no-follow open flag when available

#### Scenario: Recognized filesystem error and recovery

- **WHEN** the provider returns a recognized filesystem error such as EDQUOT or ESTALE
- **THEN** the host SHALL return the corresponding WASI errno without terminating the backend
- **AND** a later successful operation and cleanup SHALL remain possible

### Requirement: Browser timezone hosting

The browser host SHALL supply Go timezone calculations for its named browser
timezone through the Rust guest's timezone imports.

#### Scenario: Non-UTC browser timezone

- **WHEN** Chromium runs in a named timezone with distinct winter and summer offsets
- **THEN** the Go host SHALL return the correct offsets and Rust collection operations SHALL use the timezone callbacks

### Requirement: Independent real-browser verification

The reusable backend SHALL have Chromium E2E coverage using production
`wasi.js`, a caller-supplied JavaScript filesystem, the Go OS filesystem adapter
and generated clients without application SDK code. Node production-host
verification SHALL remain distinct from required Chromium acceptance.

#### Scenario: Retained acceptance report

- **WHEN** the browser suite completes or fails
- **THEN** it SHALL retain a JSON report identifying artifacts, completed cases and failures

### Requirement: Standalone generated classic MEMFS module

The browser build SHALL provide `emscripten_memfs_files`, containing a generated
`emscripten_memfs.js` ES module and all required runtime artifacts. Its default
asynchronous factory SHALL return an independent initialized module exposing
classic Emscripten `.FS` and its `.ERRNO_CODES` constants. The build SHALL use pinned, integrity-checked official
SDK inputs, declared repository runtimes and sandboxed network-free actions.
It SHALL add neither a handwritten backing filesystem nor IDBFS integration.

#### Scenario: Independent factory calls

- **WHEN** two modules are created from one imported factory
- **THEN** file, directory, descriptor and working-directory state is independent
- **AND** binary and Unicode create/read/write/rename/truncate operations retain their values

#### Scenario: Worker-compatible module loading

- **WHEN** the declared files are imported and initialized without Node or DOM globals
- **THEN** the factory exposes working classic MEMFS without undeclared runtime inputs
- **AND** genuine generated-artifact E2E emits a repeatable report identifying exact artifact hashes
