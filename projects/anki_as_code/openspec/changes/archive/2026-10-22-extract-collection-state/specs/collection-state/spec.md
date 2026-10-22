## Purpose

Manage collection-backed path bindings through a borrowed, already-open Anki
connection without owning collection or process lifecycle.

## ADDED Requirements

### Requirement: Borrow a prepared backend

The Go store SHALL borrow a prepared gRPC connection and SHALL NOT own backend
or collection lifecycle. Application wiring SHALL use BackendSTD and
BackendCollection, obtain the prepared connection through Connection, and lend
it to the state store. app.Options SHALL own collection Input. app.StateOptions
SHALL contain only namespace Key and SHALL be passed by value to WithState.
The CLI SHALL bind those fields directly without a duplicate options type; output format and
stream handling SHALL remain command concerns. The CLI SHALL NOT construct
or configure gRPC connections.

State commands SHALL require --input for an existing regular .anki2 database.
They SHALL reject missing sources, directories and unsupported extensions before
starting a worker. Archive persistence SHALL remain outside this change.
Application cleanup SHALL close the collection before the worker, use a bounded
context independent of operation cancellation, and preserve cleanup errors.

#### Scenario: Reuse a managed connection

- **WHEN** several Go state operations run against one prepared managed backend
- **THEN** they share its connection without starting, closing or replacing it

Application Start SHALL prepare the managed backend once. WithState SHALL use
its prepared connection without starting or stopping resources or using an input
path. It SHALL reject use before successful Start and after Stop. Repeated Start
SHALL NOT replace an initialized backend. Stop SHALL remain the explicit cleanup
boundary, including after failed startup. The CLI SHALL initialize after parsing
command inputs and SHALL call Stop even when command execution fails.

#### Scenario: Share a prepared application

- **WHEN** nested or sequential WithState operations run after one successful Start
- **THEN** they share the prepared backend until explicit Stop
- **AND** callback errors, cancellation and rejected repeated Start do not replace or close it

#### Scenario: Persist managed CLI operations

- **WHEN** a state command completes against an existing database
- **THEN** its managed collection and worker are stopped before command exit
- **AND** reopening the database through Anki preserves its state and native objects

#### Scenario: Reject an unavailable collection input

- **WHEN** a state command receives a missing database, directory or archive
- **THEN** it fails without creating a collection or starting a worker

### Requirement: Dump stored bindings

Dump SHALL return the selected stored protobuf State without writes or discovery.
An absent namespace SHALL return empty version-1 state.

#### Scenario: Inspect an uninitialized namespace

- **WHEN** dump reads a key without state
- **THEN** it returns empty version-1 state without writing that key

### Requirement: Import external objects in one batch

Import SHALL accept parsed configs and multiple path/native-ID string bindings.
It SHALL resolve resources before interpreting IDs, verify the requested native
objects and ownership, and save once after the complete batch succeeds. Note-type
members SHALL be bound automatically. Preset aliases SHALL reference existing
presets without editing them.

#### Scenario: Reject a partial batch

- **WHEN** any supplied binding is undeclared, missing or conflicts
- **THEN** no binding in the batch is written

#### Scenario: Resolve before interpreting IDs

- **WHEN** an undeclared path has a nonnumeric identifier
- **THEN** resource resolution fails before numeric conversion

### Requirement: Update or replace a whole state object

Update SHALL merge a version-1 external State while retaining omitted map keys
and replacing complete supplied entries. Replace SHALL install the whole object,
removing omitted bindings. Invalid input or merged ownership SHALL fail before
writing. Empty updates and identical operations SHALL perform no state write.
Input messages and native objects SHALL remain unchanged.

#### Scenario: Retain omitted entries during update

- **WHEN** update supplies valid new or replacement entries
- **THEN** those entries are installed while omitted entries remain

#### Scenario: Clear bindings through replacement

- **WHEN** replacement supplies empty version-1 state
- **THEN** all bindings are removed without changing native resources

#### Scenario: Repair invalid saved state

- **WHEN** a valid replacement is supplied for an undecodable or invalid stored object
- **THEN** it replaces the selected object without changing other namespaces or native objects
- **AND** invalid replacement input and native read failures still fail before writing

### Requirement: Typed namespaced state

State SHALL contain version 1 and Resources with separate typed deck, note,
note-type and preset maps. Note-type state SHALL own named field/template maps.
Paths SHALL be unique across entity maps; native IDs SHALL be unique within
an entity map except preset aliases. --state-key and Store.New SHALL select one
key, defaulting to anki_as_code/state/v1. Only that key SHALL be written.

#### Scenario: Independent namespaces

- **WHEN** one namespace is imported, updated or replaced
- **THEN** the other namespace and native objects remain unchanged

### Requirement: Typed text formats

The shared codec SHALL accept typed Format values codec.TOML and codec.JSON.
Empty format SHALL select TOML. State CLI input/output SHALL default to TOML
and support JSON through --format json. Update and replace SHALL accept a file
or stdin. Unsupported formats SHALL fail; native 64-bit IDs and nested bindings
SHALL round-trip without loss.

#### Scenario: Exchange both encodings

- **WHEN** state is dumped and supplied to update or replace in TOML or JSON
- **THEN** the complete protobuf object is preserved and status stays off stdout

### Requirement: Reject incomplete and unknown state commands

Missing or unknown state subcommands SHALL fail with nonzero exit status and
stderr output. They SHALL NOT print help as a successful operation or access
the backend. Explicit --help SHALL continue to succeed.

#### Scenario: Mistyped state operation

- **WHEN** a caller runs state with an unknown subcommand
- **THEN** the command fails before backend access and leaves stdout empty
