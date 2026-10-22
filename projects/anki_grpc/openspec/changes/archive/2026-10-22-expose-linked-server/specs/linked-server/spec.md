## Purpose

Allow applications to link and start the existing native Anki server while retaining its transport and protocol behavior.

## ADDED Requirements

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
