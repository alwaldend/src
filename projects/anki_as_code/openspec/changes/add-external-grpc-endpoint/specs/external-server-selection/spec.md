## Purpose

Allow collection commands to use an explicitly selected Anki gRPC server while
retaining the embedded server and existing collection workflows by default.

## ADDED Requirements

### Requirement: Optional server address

Collection commands SHALL accept optional `--address ADDRESS`, inherited by
export, plan, apply, build and all state commands. An omitted or empty address
SHALL retain the embedded server. A nonempty address SHALL select the external
server without falling back to an embedded worker. The address SHALL accept
TCP host:port, `tcp://` endpoints, `unix://` endpoints and absolute socket paths.
Invalid addresses and connection failures SHALL return errors; cancellation
SHALL stop waiting for an unavailable server.

#### Scenario: Embedded default

- **WHEN** a collection command omits `--address` or supplies an empty value
- **THEN** it uses the embedded server with existing source and output behavior

#### Scenario: External server selection

- **WHEN** a collection command supplies a TCP or Unix address
- **THEN** it operates through that server and starts no embedded worker

#### Scenario: Invalid or unavailable address

- **WHEN** the address is malformed or connection startup is cancelled
- **THEN** the command fails without starting an embedded worker or changing its source

### Requirement: Managed collections on an external server

External selection SHALL preserve source validation, exclusive collection use,
collection preparation and cleanup, state namespaces, archive publication and
build staging. The CLI SHALL close its collection and client connection after
the operation while leaving the external process running for later commands.
Documentation SHALL describe plaintext transport, exclusive server access and
shared filesystem paths, including temporary restore and staging directories.

#### Scenario: Reuse the running server

- **WHEN** sequential collection or state commands use the same external server
- **THEN** each opens and closes its own collection and the server remains usable

#### Scenario: External archive build

- **WHEN** a build selects an external server
- **THEN** it reconciles a staged copy, publishes the output and retains the base archive

#### Scenario: Managed archive publication

- **WHEN** a read-only, converged or failing operation closes a managed archive
- **THEN** publication follows the existing write-back cleanup policy, and repeated successful apply and plan operations converge semantically

### Requirement: Application transport ownership

CLI SHALL inject one application and pass one `CollectionOptions` value containing
Path, Address, Connection and StateKey into each collection operation. Managed
workflows SHALL require Path and optionally Address; borrowed workflows SHALL
require Connection and reject simultaneous Path or Address fields. Typed nil
connections SHALL count as absent. Empty addresses/namespaces on path operations
SHALL select embedded transport and the default namespace. App SHALL retain no
endpoint between calls. Build SHALL require a path and forward copied options
with its staged path into apply. Borrowed connections SHALL retain caller-owned
transport, lifetime and publication.

#### Scenario: Conflicting collection access

- **WHEN** a collection operation supplies a borrowed connection with a path or address, or omits both path and connection
- **THEN** it fails before native access and leaves caller-owned resources usable

#### Scenario: Build with a borrowed connection

- **WHEN** a build supplies an already-open connection instead of a base path
- **THEN** it fails before staging or publishing output

#### Scenario: Alternating endpoints on one App

- **WHEN** successive calls select TCP, Unix and embedded transports on one App
- **THEN** each uses only its supplied options and state namespace

### Requirement: Native API server selection

The existing `api <service> <method>` command SHALL accept `--address` to select
an external server. Exactly one nonempty collection path or address SHALL be
required. Path access SHALL use the embedded worker by default. APIRequest SHALL carry collection access in the same CollectionOptions value,
without a separate backend field. Native API access SHALL select exactly one path,
address or borrowed connection. Remote calls
SHALL use the existing APIRequest workflow, leave server process and collection
lifetime caller-owned, and close the owned client connection.

#### Scenario: Remote native calls

- **WHEN** native open, query and close RPCs select a TCP or Unix address
- **THEN** they call that server and sequential queries preserve its open collection

#### Scenario: Embedded native call

- **WHEN** a native API call supplies a collection and omits the address or uses an empty address
- **THEN** it uses an embedded managed collection

#### Scenario: Conflicting API access

- **WHEN** a native API call supplies both a nonempty collection and address
- **THEN** it fails before opening native resources
