# external-server-selection Specification

## Purpose

Allow collection commands to use an explicitly selected Anki gRPC server while
retaining the embedded server and existing collection workflows by default.

## Requirements

### Requirement: Optional server address

Collection commands SHALL accept optional --address ADDRESS, inherited by
export, plan, apply, build and all state commands. An omitted or empty address
SHALL retain the configured client default. TCP host:port, tcp://, unix:// and
absolute socket-path addresses SHALL select the requested external server without
falling back to a bundled worker. Invalid addresses and connection failures SHALL
return errors; cancellation SHALL stop waiting for an unavailable server. Other
registered local strategies SHALL follow in-process-selection and client factory
contracts.

#### Scenario: Embedded default

- **WHEN** a standard bundled collection command omits --address or supplies an empty value
- **THEN** it uses its configured bundled server with existing source and output behavior

#### Scenario: External server selection

- **WHEN** a collection command supplies a TCP or Unix address
- **THEN** it operates through that server and starts no bundled worker

#### Scenario: Invalid or unavailable address

- **WHEN** the address is malformed or connection startup is cancelled
- **THEN** the command fails without starting a bundled worker or changing its source

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

CLI SHALL inject one application and pass CollectionOptions with each operation.
Source validity and borrowed lifetime SHALL follow collection-source-api. Path
operations SHALL select their address and namespace per call; empty values SHALL
select the injected client default and default state namespace. Build SHALL
require a path and forward copied options with its staged path into apply. The
application SHALL retain no selected endpoint between calls.

#### Scenario: Conflicting collection access

- **WHEN** a collection operation combines a borrowed connection with a path or address, or omits both path and connection
- **THEN** it fails before native access and leaves caller-owned resources usable

#### Scenario: Build with a borrowed connection

- **WHEN** a build supplies an already-open connection instead of a base path
- **THEN** it fails before staging or publishing output

#### Scenario: Alternating endpoints on one App

- **WHEN** successive calls select TCP, Unix and default transports on one application
- **THEN** each uses only its supplied options and state namespace

### Requirement: Native API server selection

The api command SHALL accept --address for an external server. Managed path and
address-only selection SHALL follow grpc-cli and in-process-selection. APIRequest
SHALL carry access in the same CollectionOptions value without a separate backend
field. Address-only remote calls SHALL leave server process and collection lifetime
caller-owned and close their owned client connection. Managed paths SHALL retain
application-owned collection lifetime even when an external address selects their
transport.

#### Scenario: Remote native calls

- **WHEN** native open, query and close RPCs select a TCP or Unix address without a managed collection path
- **THEN** they call that server and sequential queries preserve its open collection

#### Scenario: Embedded native call

- **WHEN** a native API call supplies a collection and omits the address or uses an empty address
- **THEN** it uses the configured default client for its managed collection
