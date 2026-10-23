## ADDED Requirements

### Requirement: Explicit managed collection sessions

App SHALL provide ColOpen accepting ColOpenOptions with Path, optional Address
and WriteBack, returning its factory-created ankiclient.Client already started
and with its collection open. The Store SHALL return that same client without a
reservation-client wrapper. Connection
SHALL lend native access and Stop SHALL publish enabled archives and release the
owned worker. The Store SHALL retain source reservation and session-context
cancellation while the reusable client owns native preparation and cleanup.
Existing workflows SHALL accept the client's borrowed Connection without changing
their source or persistence contracts. Stop SHALL initiate cleanup once and let
callers bound waiting with a context. Session context cancellation SHALL initiate
owned cleanup under the existing bounded policy. Source reservations SHALL remain
held until Client.Done confirms final native cleanup. Completed reservations
SHALL permit immediate source reuse independently of observer scheduling.

#### Scenario: Borrow a managed session

- **WHEN** a caller opens a session and supplies its connection to multiple
  existing App operations
- **THEN** they reuse one prepared native collection
- **AND** the caller calls Stop to publish enabled archives and release
  its worker and source reservation

#### Scenario: Reopen immediately after terminal cleanup

- **WHEN** a caller successfully stops a managed session and immediately opens its source
- **THEN** the new session can reserve that source without waiting for a separate
  reservation observer to run
