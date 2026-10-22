# Explicit clients Specification

## Purpose

Provide reusable native process and collection owners with explicit lifecycle contracts while leaving application callback scopes to their consumers.

## Requirements

### Requirement: Explicit native ownership

The Go backend library SHALL expose process and collection owners through Start, Connection and Stop. The collection owner SHALL borrow a running process backend and SHALL NOT start or stop that process. Callers SHALL close the collection before stopping its process. The library SHALL provide no transport or collection manager API.

#### Scenario: Construct and close native owners

- **WHEN** a caller constructs a process and a collection backend
- **THEN** the caller starts them in dependency order, obtains the prepared connection and stops them in reverse order while retaining the existing error and restore contracts

### Requirement: Typed native commands

Go callers SHALL invoke Anki using upstream generated protobuf service clients
through the prepared gRPC connection. Stdio clients SHALL carry gRPC HTTP/2 on
the worker pipes, reserve stdout for protocol traffic and use stderr for
diagnostics. Original service routes and backend-prefixed routes SHALL retain
Anki's dispatcher indices. Shared connection handling SHALL preserve call
options, gRPC status, native messages and upstream error trailers. Failed calls
SHALL NOT restart the worker or replay requests; cancellation SHALL NOT promise
rollback of native mutations.

#### Scenario: Typed service clients

- **WHEN** a caller constructs generated collection, note and note-type clients
  from the prepared connection
- **THEN** it can open storage and invoke typed operations with original and
  backend-prefixed routes reaching the corresponding native operations

#### Scenario: Persist a note

- **WHEN** a caller creates a note, closes its collection and reopens it
- **THEN** Anki returns its persisted fields and generated cards

#### Scenario: Call options and native failures

- **WHEN** a generated client supplies gRPC call options and Anki rejects its request
- **THEN** the options remain effective and the error retains its status, native
  message and upstream error trailer

#### Scenario: Failed command and closed process

- **WHEN** an invalid command, cancelled context or stopped worker is used
- **THEN** the caller receives a contextual error without worker restart or replay

### Requirement: Explicit worker lifecycle

The stdio client SHALL accept an explicitly configured executable, arguments and
shutdown timeout without starting a process during construction. It SHALL copy
arguments and inherit the parent environment and stderr. The library SHALL NOT
supply an implicit server executable runtime input. A zero shutdown timeout
SHALL select five seconds; a negative timeout SHALL fail before resource
allocation. Each instance SHALL permit one Start attempt and expose a connection
only after gRPC readiness, without a client initialization RPC. Cancelling the
startup context after successful readiness SHALL NOT stop the worker.

Stop SHALL be terminal, safe before startup and after failure, reject new calls,
reap every started child and preserve its terminal result across repeated calls.
A caller's earlier deadline or cancellation SHALL bound graceful shutdown;
otherwise the configured grace period SHALL apply. On expiry, Stop SHALL
terminate and reap the worker while preserving applicable context errors and
unexpected process-exit errors. Connection SHALL fail before readiness, during
shutdown and after worker exit.

#### Scenario: Configure and start

- **WHEN** a caller constructs a stdio client and starts it
- **THEN** the copied binary arguments reach the worker, generated clients become
  usable after readiness and a second Start fails without replacing the worker

#### Scenario: Stop without a running worker

- **WHEN** Stop is called before Start or after failed startup
- **THEN** owned resources are released safely and later Start cannot create a worker

#### Scenario: Startup context cancellation after success

- **WHEN** the Start context is cancelled after readiness
- **THEN** the worker remains available for calls with their own contexts until Stop

#### Scenario: Worker crash

- **WHEN** the worker exits unexpectedly or is killed
- **THEN** calls fail, the child is reaped and repeated Stop retains its exit failure
  without replaying requests

#### Scenario: Startup failure

- **WHEN** the worker cannot start or finish its handshake within the startup deadline
- **THEN** startup reports the failure and releases owned transport and process resources

#### Scenario: Shutdown deadline

- **WHEN** a worker remains alive after EOF until the Stop context deadline
- **THEN** Stop terminates and reaps it, preserves the deadline error and rejects
  subsequent connection access

#### Scenario: Explicit shutdown cancellation

- **WHEN** the Stop context is cancelled while the worker remains alive after EOF
- **THEN** Stop terminates and reaps it and preserves cancellation and process-exit
  errors across repeated calls

#### Scenario: Default shutdown grace expires

- **WHEN** the worker uses the default timeout and remains alive after EOF without
  an earlier caller deadline
- **THEN** Stop terminates and reaps it after the five-second grace period and
  preserves its process-exit failure

#### Scenario: Configured shutdown grace expires

- **WHEN** the caller selects a positive shutdown timeout and the worker remains
  alive after EOF
- **THEN** Stop uses that grace period, reaps the worker and preserves its exit failure

#### Scenario: Invalid shutdown timeout

- **WHEN** Start is attempted with a negative shutdown timeout
- **THEN** it reports invalid configuration before allocating pipes or starting a child

#### Scenario: Connection during a lifecycle transition

- **WHEN** Connection is called before readiness or during shutdown
- **THEN** it returns an error without lending an unready connection
