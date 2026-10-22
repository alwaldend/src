# collection-source-api Specification

## Purpose

Allow library callers to reuse open Anki collections across application operations while preserving application-owned lifecycle for ordinary path inputs.

## Requirements

### Requirement: Select a collection for each operation

The application SHALL accept CollectionOptions selecting a managed path or an
already open collection connection for collection export, planning, in-place
apply and state operations. A managed path can carry a client address; borrowed
connections SHALL reject simultaneous path or address fields. Typed nil
connections SHALL count as absent. A single application instance SHALL support
multiple independent sources without retaining an operation's source or namespace.
Empty or conflicting sources SHALL fail before native operations.

#### Scenario: Reuse multiple open collections

- **WHEN** a caller alternates operations between two open collections on one application
- **THEN** each operation uses its supplied collection and subsequent calls can reuse both collections

#### Scenario: Reject an invalid source

- **WHEN** the caller supplies neither a path nor a connection, or combines a connection with a path or address
- **THEN** the operation fails without opening or mutating a collection

### Requirement: Preserve collection ownership

Managed path operations SHALL delegate preparation, publication and cleanup to the collection-store capability. Borrowed operations SHALL leave connection lifetime and package publication to the caller, including after errors and cancellation. Direct database mutations SHALL retain their native persistence behavior.

#### Scenario: Complete a borrowed workflow

- **WHEN** a caller exports configs, installs state, plans and applies changes through one borrowed connection
- **THEN** the caller can continue using that connection and explicitly publish the resulting archive after the operations finish

#### Scenario: Fail a borrowed operation

- **WHEN** an operation on a borrowed connection fails or is cancelled
- **THEN** the application returns the error without closing the caller's collection or publishing its package

### Requirement: Keep command handling independent

CLI commands SHALL continue accepting their existing path flags and SHALL delegate opening, publication and cleanup to the application. The complete application interface SHALL live in the shared API package and use shared source and protobuf types without importing the concrete application or native backend owner. CLI SHALL consume that public interface, and executable composition SHALL check concrete App compatibility.

#### Scenario: Supply an application implementation

- **WHEN** a caller or CLI test implements the public application interface
- **THEN** it can supply that implementation without depending on the concrete App or native worker

#### Scenario: Run a path-based command

- **WHEN** a user invokes an existing collection command with its path flag
- **THEN** the command retains its existing outputs and persistence semantics while the application owns collection preparation and cleanup
