## ADDED Requirements

### Requirement: Configurable server transport

The server SHALL expose CLI help and select stdio, TCP, or Unix-socket transport.
Stdio MUST be the default. TCP and Unix transports MUST require an explicit
listen address or socket path. Invalid options MUST fail before serving requests.
All transports MUST forward the same upstream Anki protobuf commands. The server
MUST initialize one shared backend at startup without a client initialization RPC
and allow explicit collection close and reopen.
Anki's library compilation MUST remain in third_party; projects/anki_grpc MUST
depend on the native target rather than declare an upstream Anki library build.

#### Scenario: Default transport

- **WHEN** the existing Go API starts the server without arguments
- **THEN** the server serves gRPC on stdin/stdout with upstream commands ready without an initialization RPC

#### Scenario: TCP listener

- **WHEN** the server starts with --transport tcp and an explicit --listen address
- **THEN** clients can immediately call Anki and persist a synthetic collection over TCP
- **AND** the listener accepts a later client after the first disconnects

#### Scenario: Unix listener

- **WHEN** the server starts with --transport unix and an explicit --listen path
- **THEN** clients can immediately call Anki and persist a synthetic collection over a Unix socket
- **AND** the listener accepts a later client after the first disconnects

#### Scenario: CLI errors

- **WHEN** an unsupported transport, missing argument, unknown option, or listen value with stdio is supplied
- **THEN** the server exits unsuccessfully with a useful error and no protocol output
