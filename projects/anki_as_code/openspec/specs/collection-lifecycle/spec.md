# collection-lifecycle Specification

## Purpose

Keep collection operations within one application-owned native scope that preserves source identity, publication policy and deterministic cleanup.

## Requirements

### Requirement: Application-owned native scope

Collection access SHALL prepare one address-selected client and one collection,
lend a prepared connection, close the collection and then stop its client before
releasing source reservations. Reservations SHALL remain held until asynchronous
native cleanup completes, including when bounded client shutdown returns first.
Stopping an external client SHALL leave its server process running. Normal package
restores SHALL be removed; uncertain restores SHALL retain the backend's recovery
contract. Source aliases SHALL share one reservation, and archive publication SHALL
follow the configured WriteBack policy.

Collection package import and export use Anki APIs; no handwritten package codec or separate archive resource abstraction is required.

#### Scenario: Complete collection work

- **WHEN** a collection callback succeeds with WriteBack enabled
- **THEN** its native changes are published to its source when appropriate, its prepared connection expires and its source becomes reusable

#### Scenario: Native packaging preserves collection content

- **WHEN** a native package is reconciled or rebuilt and imported again through Anki
- **THEN** supported media, identities, history and stored state survive the package round trip
- **AND** managed source publication follows the collection-store write-back policy

### Requirement: Cleanup on unsuccessful exits

Failed startup, callback error, cancellation and panic SHALL execute bounded cleanup independently of caller cancellation. Callback errors, cancellation causes and process-exit errors SHALL remain observable. A panic SHALL propagate after cleanup.

#### Scenario: Callback fails while the worker crashes

- **WHEN** the callback returns an error after the worker exits unexpectedly
- **THEN** cleanup reaps the worker, retains both error identities and releases the source reservation without publishing the archive

#### Scenario: Cancel or panic inside a scope

- **WHEN** a callback cancels its context or panics
- **THEN** cleanup runs before scope exit, the original cause remains observable and the collection source can be opened again after native cleanup completes

#### Scenario: Failed preparation releases resources

- **WHEN** native preparation fails or is cancelled
- **THEN** no prepared connection is lent, owned client cleanup is attempted and supplied databases are not removed
