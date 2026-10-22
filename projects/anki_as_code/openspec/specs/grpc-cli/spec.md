# grpc-cli Specification

## Purpose

Expose upstream Anki services through a complete command tree while retaining protobuf schemas as the API authority.

## Requirements

### Requirement: Complete API commands

The CLI SHALL expose kebab-case service and method commands for every supported upstream Anki RPC, including implicit backend methods, without loading proto files or contacting a server for discovery.

#### Scenario: Discover a previously unlinked service

- **WHEN** a user requests `api tags --help` without a backend
- **THEN** the CLI lists the pinned TagsService methods without starting Anki

### Requirement: Explicit backend selection

API commands SHALL require a collection path or client address. A path SHALL
select managed collection lifetime and can include an address selecting its
transport, as defined by in-process-selection. Address-only access SHALL leave
an external server's process and collection lifetime caller-owned. Borrowed
connections SHALL follow collection-source-api ownership. A borrowed connection
combined with a path or address SHALL fail before native access. CLI address
selection SHALL use --address.

#### Scenario: Missing access mode

- **WHEN** an API command supplies neither a collection path nor an address
- **THEN** the command fails before opening storage or contacting a server

### Requirement: Protobuf payloads and persistence control

Requests and responses SHALL support JSON and TOML with independent format selection, stdin or an input file, stdout responses and stderr errors. Application collection scopes SHALL explicitly enable backend write-back by default. Offline archive write-back SHALL default to enabled; `--write-back` and `--no-write-back` SHALL update one setting in argument order.

#### Scenario: Disable publication

- **WHEN** the last persistence flag disables write-back
- **THEN** an offline archive source remains unchanged after a mutating RPC

### Requirement: Self-contained application bundle

Default offline API access without an explicit client address in a host application
built with its native server SHALL run that linked server by
launching the same executable with server serve, without searching PATH or Bazel
runfiles or extracting another executable. Server-free builds SHALL follow their
explicit external-client or WASI capability contracts.

#### Scenario: Run an installed bundle

- **WHEN** the native-server application executable is copied outside Bazel and run with an empty PATH
- **THEN** offline API calls succeed using its linked server

### Requirement: Cancellable request input

The command timeout SHALL cover request input, startup and invocation. Interrupt cancellation SHALL end a command waiting for incomplete stdin without opening native storage or emitting a response.

#### Scenario: Stalled request producer

- **WHEN** a producer supplies request bytes but keeps stdin open and the deadline expires or the user interrupts the command
- **THEN** the command exits with the context error, empty stdout and an unopened collection

### Requirement: Shared application operation

The shared application SHALL validate and invoke upstream requests through one generic API operation. It SHALL accept managed collection options with an optional client address, a borrowed open connection or an explicit client address. CLI handlers SHALL delegate native execution and managed-lifetime validation to App. Borrowed connections SHALL retain caller-owned cleanup and publication, with native side effects limited to the requested RPC.

#### Scenario: Reuse a caller-owned collection

- **WHEN** a Go caller invokes a raw API request using a borrowed collection connection
- **THEN** App returns the protobuf response without closing the connection or automatically publishing its source archive
