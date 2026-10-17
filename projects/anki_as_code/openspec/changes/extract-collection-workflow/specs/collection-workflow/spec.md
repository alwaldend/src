## Purpose

Provide one native lifecycle for state operations on offline Anki databases and collection packages while preserving their unmanaged content.

## ADDED Requirements

### Requirement: Collection inputs and preparation

State commands SHALL accept an existing offline .anki2 database or .colpkg package through --collection. Missing and non-regular inputs SHALL fail before native startup. App.WithCollection SHALL compose the native backends to prepare and release the connection once per command, closing the collection before the server on success, failure and cancellation. State callers SHALL construct namespaced stores directly from the prepared connection; stores SHALL NOT own native lifecycle or cleanup.

#### Scenario: State operation on either input type

- **WHEN** a state command receives a supported collection input
- **THEN** it uses one prepared native connection and releases its resources after the operation

### Requirement: Native persistence boundaries

Database operations SHALL write in place without copying the database. Changed package operations SHALL publish through Anki's native atomic export including media. Read-only or unchanged package operations SHALL preserve the source archive bytes. Failures before package export SHALL preserve the source package; committed direct database writes SHALL persist even if a later operation fails.

#### Scenario: Failed package operation

- **WHEN** a package operation writes state and then returns an error or is cancelled before publication
- **THEN** the source package remains unchanged and normal restored temporary files are cleaned up

### Requirement: State command integration

Dump, import, update and replace SHALL use the common collection workflow, preserve native resources and other state namespaces, and report successful mutation status only after persistence succeeds. Import SHALL keep all supplied bindings in one operation. Dump SHALL write encoded state to stdout and mutation status SHALL use stderr.

#### Scenario: Repeated and rejected writes

- **WHEN** state writes are identical, invalid or conflict with existing ownership
- **THEN** they leave the source collection unchanged and do not partially publish the requested bindings
