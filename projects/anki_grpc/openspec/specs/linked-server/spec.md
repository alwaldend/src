# Linked server Specification

## Purpose

Allow applications to link and start the existing native Anki server while retaining its transport and protocol behavior.

## Requirements

### Requirement: Linked server entry point

A Go caller SHALL be able to run the native server through one linked entry point with ordinary CLI arguments. The entry point SHALL preserve stdio, TCP, Unix, help, invalid-argument errors, and native gRPC framing. It SHALL NOT add per-method foreign-function bindings.

#### Scenario: Linked stdio operation

- **WHEN** an application invokes the linked server with serve and stdio arguments
- **THEN** generated gRPC clients can use its stdin/stdout transport with the same native responses and error trailers

### Requirement: Standalone server compatibility

The standalone anki_grpc executable SHALL use the same server implementation and start serving only through its serve subcommand while retaining its transport options. Bare invocation and root help SHALL show commands without initializing the native backend; serve --help SHALL show transport options. Unknown commands SHALL fail before protocol startup.

#### Scenario: Standalone network clients

- **WHEN** existing clients start the standalone server with serve on TCP or Unix
- **THEN** readiness, collection persistence and reconnection continue to work

### Requirement: Configurable server transport

The serve subcommand SHALL select stdio, TCP or Unix transport, defaulting to
stdio. TCP and Unix SHALL require an explicit listen address or socket path.
Invalid options SHALL fail before protocol serving. Every transport SHALL
forward the same upstream Anki protobuf commands. Serving SHALL initialize one
shared backend without a client initialization RPC and allow explicit collection
close and reopen. Upstream Anki library compilation SHALL remain owned by
third_party; the project SHALL consume that library.

#### Scenario: Default transport

- **WHEN** the server starts with serve and no transport override
- **THEN** it serves gRPC on stdin/stdout with upstream commands ready without an
  initialization RPC

#### Scenario: TCP listener

- **WHEN** the server starts with serve --transport tcp and an explicit --listen address
- **THEN** clients can persist a native collection over TCP
- **AND** a later connection uses the same initialized backend after a client disconnects

#### Scenario: Unix listener

- **WHEN** the server starts with serve --transport unix and an explicit --listen path
- **THEN** clients can persist a native collection over a Unix socket
- **AND** a later connection uses the same initialized backend after a client disconnects

#### Scenario: CLI errors

- **WHEN** the caller supplies an unsupported transport, missing argument, unknown
  option or listen value with stdio
- **THEN** the server exits unsuccessfully with a useful diagnostic and no protocol output
