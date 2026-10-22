# collection-state Specification

## Purpose

Manage collection-backed path bindings through a borrowed, already-open Anki
connection without owning collection or process lifecycle.

## Requirements

### Requirement: Borrow a prepared backend

The Go state store SHALL borrow a prepared gRPC connection and SHALL NOT own
backend or collection lifecycle. Application collection selection SHALL follow
collection-source-api; managed scopes SHALL follow collection-store and
collection-lifecycle. The CLI SHALL bind source and namespace options without
constructing or configuring gRPC connections; output formats and streams SHALL
remain command concerns.

State commands SHALL require --collection for an existing regular offline
.anki2 database or .colpkg package. Missing sources and unsuitable filesystem
entries SHALL fail before native startup. Managed cleanup SHALL close the
collection before its client, use a bounded context independent of operation
cancellation, and preserve cleanup errors.

#### Scenario: Reuse a managed connection

- **WHEN** several Go state operations run against one prepared managed collection
- **THEN** they share its connection without starting, closing or replacing it

#### Scenario: Share a prepared application

- **WHEN** sequential state operations use a caller-owned open connection through one App
- **THEN** they share that supplied connection while its lifetime remains caller-owned
- **AND** errors or cancellation do not replace or close it

#### Scenario: Persist managed CLI operations

- **WHEN** a state command completes against an existing database
- **THEN** its managed collection and client are stopped before command exit
- **AND** reopening the database through Anki preserves its state and native objects

#### Scenario: Reject an unavailable collection input

- **WHEN** a state command receives a missing database or a directory
- **THEN** it fails without creating a collection or starting a client

#### Scenario: State operation on either input type

- **WHEN** a state command receives a supported database or package input
- **THEN** it uses one prepared native connection and releases its owned resources after the operation

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

CLI import parsing SHALL pass opaque identifier strings, reject malformed bindings and duplicate paths, and leave native numeric validation to the state store after resource resolution.

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

#### Scenario: Repeated and rejected writes

- **WHEN** state writes are identical, invalid or conflict with existing ownership
- **THEN** they leave stored bindings unchanged and do not partially install the requested batch
- **AND** managed archive cleanup follows the collection-store publication policy

#### Scenario: Reject a conflicting update

- **WHEN** input is valid alone but conflicts with retained native or path ownership
- **THEN** the complete update fails before writing any state

### Requirement: Typed namespaced state

State SHALL contain version 1 and Resources with separate typed deck, note,
note-type and preset maps. Note-type state SHALL own named field/template maps.
Paths SHALL be unique across entity maps; native IDs SHALL be unique within
an entity map except preset aliases. The collection group's --state-key flag and state.New SHALL select one
key, defaulting to anki_as_code/state/v1. Only that key SHALL be written.

Synthetic member fragments and preset addresses are not config-path keys. Persisted state uses typed entity maps without a generalized kind discriminator.

#### Scenario: Independent namespaces

- **WHEN** one namespace is imported, updated or replaced
- **THEN** the other namespace and native objects remain unchanged

#### Scenario: Persistent bindings

- **WHEN** bindings are explicitly installed and the collection is closed, reopened or rebuilt into a colpkg
- **THEN** native identities remain bound to paths and subsequent planning converges

#### Scenario: Preserve nested identities through a move

- **WHEN** a bound note type or one of its members is renamed through reconciliation
- **THEN** the nested native identity and ordinal remain bound without synthetic state paths

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
