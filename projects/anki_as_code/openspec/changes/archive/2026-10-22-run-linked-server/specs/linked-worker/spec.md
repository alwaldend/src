## Purpose

Run Anki from one installed CLI executable without extracting another executable or changing collection workflows.

## ADDED Requirements

### Requirement: Native CLI passthrough

anki_as_code server run SHALL forward its argument vector to the linked Rust CLI without interpreting native subcommands or flags. Rust SHALL own serve dispatch, help and invalid-argument errors. Diagnostics SHALL use stderr, and interrupt signals SHALL terminate a listening native server.

The server group SHALL show help without starting native execution. Its serve shortcut SHALL delegate to run with serve prepended and SHALL NOT duplicate the native argument parser or server implementation.

#### Scenario: Show native operations

- **WHEN** a user invokes anki_as_code server run without native arguments or with --help
- **THEN** the command succeeds with native root help listing serve and starts no server

#### Scenario: Interrupt a listening server

- **WHEN** a user interrupts anki_as_code server run serve or server serve while it is listening
- **THEN** the server process terminates without leaving an owned worker process

### Requirement: Self-contained offline worker

Offline CLI collection scopes SHALL start their current executable with server serve and communicate over stdio gRPC. They SHALL NOT materialize an executable or require another installed server or PATH lookup. Existing collection behavior SHALL remain unchanged.

#### Scenario: Installed execution without temporary storage

- **WHEN** the CLI alone is copied outside Bazel, PATH is empty, and temporary storage is unavailable
- **THEN** a direct-database operation succeeds using its linked server and emits a repeatable result

### Requirement: Explicit library worker configuration

Application and collection-store native operations SHALL use only their configured executable and argument vector. Missing configuration SHALL fail without opening the collection or invoking the callback. Production code SHALL NOT discover a default worker through Bazel runfiles or PATH, and the client library SHALL NOT supply an implicit standalone-server runtime input. Configuration-only application operations SHALL remain usable without native worker settings.

#### Scenario: Missing native worker settings

- **WHEN** a caller accesses an existing collection through an unconfigured store, even with a discoverable standalone server
- **THEN** access fails without invoking the callback or changing the collection, and a later explicitly configured scope can use that source
