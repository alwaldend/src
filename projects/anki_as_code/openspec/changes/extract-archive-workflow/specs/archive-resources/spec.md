## Purpose

Provide native private collection resources and guarded publication while Anki
owns package compatibility, collection semantics, and media handling.

## ADDED Requirements

### Requirement: Native private collection resources

Sources SHALL restore packages through Anki's collection-package API into private
paths and open that collection. Offline databases SHALL be copied without SQL;
WAL and journal sidecars SHALL be rejected. The original SHALL remain unchanged.

#### Scenario: Package round trip

- **WHEN** a native package is opened, edited through Anki, and published
- **THEN** native reimport preserves media, identities, review history, and custom state

### Requirement: Guarded publication

Publication SHALL stage beside the destination, sync the completed file, preserve
requested mode, and check an optional expected digest before replacement.

#### Scenario: Changed target or failed output

- **WHEN** the target digest differs or an operation fails
- **THEN** the target remains unchanged and temporary outputs are removed

### Requirement: Borrowed native connection

Sources SHALL borrow an existing connection without starting or stopping its
worker. The caller SHALL retain process ownership. Native filesystem operations
SHALL finish before scratch cleanup, including when cancellation is requested.
Publish SHALL consume the source on every attempt, while cleanup is repeatable.

#### Scenario: Failed or cancelled operation

- **WHEN** native import, opening, or publication fails or is cancelled
- **THEN** native work is awaited before its files are removed, and the worker remains usable

#### Scenario: Reused source

- **WHEN** publication is attempted after closing or publishing
- **THEN** it fails without creating output
