## Purpose

Keep collection operations within one application-owned native scope that preserves source identity, publication policy and deterministic cleanup.

## ADDED Requirements

### Requirement: Application-owned native scope

Collection access SHALL prepare one native worker and one collection, lend a prepared connection, close the collection and then stop its worker before releasing source reservations. Normal package restores SHALL be removed; uncertain restores SHALL retain the backend's existing retention contract. Existing source aliases and archive publication rules SHALL remain unchanged.

#### Scenario: Complete collection work

- **WHEN** a collection callback succeeds and requests publication
- **THEN** its native changes are published to its source when appropriate, its prepared connection expires and its source becomes reusable

### Requirement: Cleanup on unsuccessful exits

Failed startup, callback error, cancellation and panic SHALL execute bounded cleanup independently of caller cancellation. Callback errors, cancellation causes and process-exit errors SHALL remain observable. A panic SHALL propagate after cleanup.

#### Scenario: Callback fails while the worker crashes

- **WHEN** the callback returns an error after the worker exits unexpectedly
- **THEN** cleanup reaps the worker, retains both error identities and releases the source reservation without publishing the archive

#### Scenario: Cancel or panic inside a scope

- **WHEN** a callback cancels its context or panics
- **THEN** cleanup runs before scope exit, the original cause remains observable and the collection source can be opened again
