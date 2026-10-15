## Purpose

Expose the pinned Anki protobuf operations through a discoverable command-line
interface without maintaining handwritten wrappers for each operation.

## ADDED Requirements

### Requirement: Schema-complete commands

The CLI SHALL expose each declared upstream unary service and method with its
exact protobuf name and show its request/response types in help.

#### Scenario: Discover operations

- **WHEN** help is requested for a service or method
- **THEN** the operation and message types are discoverable without starting Anki

### Requirement: Protobuf JSON invocation

The CLI SHALL accept ProtoJSON inline, from a file, or stdin and emit a ProtoJSON
response. Invalid input, native failures, and transport errors SHALL fail.

#### Scenario: Edit and inspect

- **WHEN** a valid native mutation is invoked and its object is read afterward
- **THEN** the response is JSON and the mutation persists according to Anki semantics

### Requirement: Explicit connection ownership

The CLI SHALL support a managed default stdio worker with optional one-shot local
collection opening, and existing TCP/Unix targets. It SHALL stop managed workers
on success or failure and SHALL leave an external server's lifetime unchanged.

#### Scenario: Existing server

- **WHEN** separate CLI invocations open, mutate, and close a collection on a target
- **THEN** they share that server's state without implicitly stopping the server

#### Scenario: Cancelled call

- **WHEN** the deadline expires or cancellation occurs
- **THEN** the invocation fails and owned worker resources are released
