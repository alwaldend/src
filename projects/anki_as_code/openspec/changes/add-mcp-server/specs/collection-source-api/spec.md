## ADDED Requirements

### Requirement: Explicit managed collection sessions

App SHALL provide ColOpen accepting ColOpenOptions with Path, optional Address
and WriteBack, returning an already-open ankiclient.Client exposing Connection
and Stop. Store SHALL prepare one private, single-use collection scope with a raw
ankiclient.StartableClient transport. The application SHALL reuse the shared
client contracts instead of declaring a separate Collection interface.
The application scope SHALL own native preparation, archive publication and
cleanup. Existing workflows SHALL borrow its Connection without changing their
source or persistence contracts. Stop SHALL initiate cleanup once and allow
callers to bound waiting. Session cancellation SHALL initiate bounded cleanup.
Source reservations SHALL remain held until private completion tracking confirms
native release, and completed reservations SHALL permit immediate source reuse.

#### Scenario: Borrow a managed session

- **WHEN** a caller opens a session and supplies its connection to multiple
  existing App operations
- **THEN** they reuse one prepared native collection
- **AND** Stop publishes enabled archives and releases its worker and reservation

#### Scenario: Reopen after terminal cleanup

- **WHEN** a caller successfully stops a session and immediately opens its source
- **THEN** the source can be reserved without waiting for an observer to run

### Requirement: Own managed storage in the application

The collection scope SHALL open databases directly and restore packages through
Anki into owned temporary storage. Enabled write-back SHALL publish archives
independently of operation errors. Cleanup SHALL drain active calls within its
budget, avoid duplicate native close after consuming operations, and preserve
owned restore files when native completion cannot be confirmed. Optional
transport completion, observed through a private Done-only capability, MAY
establish release after Stop returns early; a raw
network disconnect SHALL NOT establish remote native completion.

#### Scenario: Cancel a native operation

- **WHEN** a managed native request returns cancellation without a completion
  response and cleanup cannot confirm release
- **THEN** the application reports incomplete cleanup and retains its restore and
  source reservation while the transport StartableClient remains a three-method interface

#### Scenario: Complete an explicit package export

- **WHEN** a native full-package export consumes the managed collection
- **THEN** scope cleanup avoids closing consumed storage twice and restores access
  internally only when needed for enabled source-package publication
