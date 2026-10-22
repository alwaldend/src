# In-process native clients Specification

## Purpose

Allow generated Anki clients to execute native operations inside their Go process with explicit ownership and the established native error contract.

## Requirements

### Requirement: Direct native execution

The backend SHALL implement the existing start, connection and stop contract and execute upstream unary commands through a C ABI without starting a process or transport listener. Generated clients SHALL retain native persistence, concurrent access, original service routes and serialized native errors. Unsupported streaming and call options SHALL fail explicitly.

#### Scenario: Native collection workflow

- **WHEN** generated clients create a note, close and reopen its database through the direct backend
- **THEN** the note and generated cards persist and native failures retain their status and error trailers

### Requirement: Native lifetime

The backend SHALL permit one start attempt and reject calls once stopping begins. Cancellation SHALL stop waiting without freeing resources still used by a native call. Stop SHALL release native state after active calls complete; its deadline SHALL bound waiting without pretending unfinished work was terminated. Repeated Stop calls SHALL be safe.

#### Scenario: Cancellation during native work

- **WHEN** a native media request remains active after its caller cancels and Stop times out
- **THEN** retained clients reject calls, native resources remain valid, and later Stop completes after the native work finishes

#### Scenario: Shutdown overlaps native initialization

- **WHEN** synchronous native startup is still running and Stop reaches its deadline
- **THEN** Stop returns the deadline error, Done waits for startup cleanup, and startup does not expose a connection after shutdown

### Requirement: In-process address selection

Importing the optional native adapter SHALL register `inprocess://` and `inprocess+cgo://` constructors selecting fresh direct clients. In-process suffixes SHALL be opaque labels and SHALL NOT create shared state or network resources. Executable composition MAY supply the WASM-only generic selection described by wasm-worker. Worker registration, address routing and configured factory isolation SHALL follow the client-connection-platforms capability.

#### Scenario: Independent named instances

- **WHEN** callers start two native clients selected with the same in-process address
- **THEN** each has independent collection state and stopping either leaves the other usable

### Requirement: Explicit optional native imports

Direct native execution SHALL require opting into the native subpackage through an explicit import that registers only its direct client constructors. Generic library dependency isolation, unavailable-engine errors and custom factory composition SHALL follow client-connection-platforms.

#### Scenario: Generic caller without native support

- **WHEN** a caller imports only the generic constructor and requests in-process execution
- **THEN** construction fails clearly and the caller can build with pure Go enabled

#### Scenario: Registered client strategy

- **WHEN** an adapter registers a client factory under a custom scheme
- **THEN** managed collection scopes use its client without transport classification or backend ownership methods
