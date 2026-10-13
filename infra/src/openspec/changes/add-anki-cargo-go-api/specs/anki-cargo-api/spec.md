## Purpose

Provide a Go API consuming the separately owned Anki server and preserving upstream collection semantics.

## ADDED Requirements

### Requirement: Server runtime inputs

The Go API build SHALL include the projects/anki_grpc executable as runtime data.
The reusable API SHALL live in projects/anki_grpc/pkg/backend. Go callers configure
its binary path and arguments and construct clients generated from upstream protobuf services.
The client MUST NOT duplicate native message schemas or provide a C ABI bridge.

#### Scenario: Go consumer build

- **WHEN** the Go API consumer is built with its server prerequisite
- **THEN** Bazel includes the projects/anki_grpc executable in its runtime inputs

### Requirement: Native Go API

Go callers SHALL invoke native Anki commands with upstream protobuf messages over gRPC HTTP/2 carried by stdin and stdout, without a network listener or C ABI. Go MUST own process startup, shutdown, cancellation, and crash reporting and reject commands after close. The server initializes Anki at startup; Go SHALL wait for gRPC readiness without an initialization RPC. Restart and recovery policy belongs to Go. Stdout MUST carry only protocol traffic; diagnostics MUST use stderr. Failed mutations MUST NOT be automatically retried. Cancelling a call MUST NOT promise rollback of a mutation already running in Anki.

The connection owner SHALL expose a grpc.ClientConnInterface for caller-created
generated service clients. Callers SHALL invoke typed methods without command
strings or manually supplied response objects. Native error translation and
lifecycle checks SHALL be shared by all generated clients through the connection.
The Rust server SHALL derive both original service routes and backend-prefixed
routes from Anki's descriptor helpers, preserving native dispatcher indices.

#### Scenario: Typed service clients

- **WHEN** a caller constructs Anki's generated collection, note, and note-type clients from the managed connection
- **THEN** it can open a collection and invoke typed note operations without handwritten command wrappers
- **AND** a generated collection-service method also declared by its backend service reaches the corresponding native operation

#### Scenario: Persist a note

- **WHEN** a caller opens a synthetic collection, creates a note through Anki, closes it, and reopens it
- **THEN** Anki returns the persisted fields and generated cards

#### Scenario: Failed command and closed process

- **WHEN** an invalid command, cancelled context, or closed backend is used
- **THEN** the API returns a contextual error without restarting the worker

#### Scenario: Call options and native failures

- **WHEN** a generated client supplies gRPC call options and Anki rejects its request
- **THEN** the options remain effective and the error retains its gRPC status, native message, and upstream error trailer

#### Scenario: Worker crash

- **WHEN** a live worker is killed or exits unexpectedly
- **THEN** Go reports failed calls, reaps the child, and returns its exit failure on close without replaying requests

#### Scenario: Startup failure

- **WHEN** the worker cannot start or complete the gRPC handshake within the caller's deadline
- **THEN** Go reports the failure and closes all owned transport and process resources

### Requirement: Explicit backend lifecycle

The Go API SHALL define Backend with Start(context.Context) error,
Connection() (grpc.ClientConnInterface, error), and Stop(context.Context) error.
BackendSTD SHALL implement it using stdin/stdout. New SHALL accept only a
BackendSTDOptions value with Binary, Arguments, and ShutdownTimeout, copy the
arguments, and perform no process startup. A zero ShutdownTimeout SHALL use five
seconds; Start SHALL reject a negative value before allocating worker resources.
Start SHALL create the child and wait for gRPC readiness; its context SHALL
control startup only; after successful startup its cancellation SHALL NOT stop
the child. The child SHALL inherit the parent's
environment and stderr. Each instance SHALL permit only one Start attempt.
Stop SHALL be safe before startup and after failure, retain its result on repeat,
reject new calls immediately, and reap every started child. The Stop context SHALL
bound graceful shutdown; cancellation or its deadline SHALL trigger forced process
termination and reaping while preserving the context cause. Connection SHALL return
an error until readiness, during shutdown, after Stop, and when the worker has exited.
Graceful shutdown SHALL expire after the configured ShutdownTimeout when the
caller's context permits a longer wait.

#### Scenario: Configure and start

- **WHEN** a caller constructs BackendSTD and then invokes Start through Backend
- **THEN** construction exposes no connection, Start forwards the recorded binary arguments, and generated clients can use the resulting connection
- **AND** a second Start fails without disrupting the existing worker

#### Scenario: Stop without a running worker

- **WHEN** Stop is called before Start or after failed startup
- **THEN** Stop releases any owned resources safely and a later Start cannot create a worker

#### Scenario: Startup context cancellation after success

- **WHEN** the Start context is cancelled after readiness
- **THEN** the worker remains available to service calls with their own contexts until Stop

#### Scenario: Shutdown deadline

- **WHEN** a worker does not exit gracefully before the Stop context deadline
- **THEN** Stop terminates and reaps it, returns an error preserving context.DeadlineExceeded, and subsequent Connection calls return an error

#### Scenario: Connection during a lifecycle transition

- **WHEN** Connection is called before Start completes or while Stop is waiting
- **THEN** it returns an error without exposing an unready handle or waiting for shutdown to complete

#### Scenario: Explicit shutdown cancellation

- **WHEN** a caller cancels Stop's context while the worker remains alive after EOF
- **THEN** Stop terminates and reaps the child and retains context.Canceled and the process exit failure on repeated calls

#### Scenario: Default shutdown grace expires

- **WHEN** a worker uses the default shutdown timeout, remains alive after EOF, and the Stop context has no earlier deadline
- **THEN** Stop terminates and reaps it after the five-second grace period and preserves the process exit failure

#### Scenario: Configured shutdown grace expires

- **WHEN** a caller sets a positive ShutdownTimeout and the worker remains alive after EOF
- **THEN** Stop applies that timeout instead of the default, reaps the child, and preserves its exit failure

#### Scenario: Invalid shutdown timeout

- **WHEN** a caller attempts Start with a negative ShutdownTimeout
- **THEN** Start reports the invalid configuration without allocating pipes or starting a child
