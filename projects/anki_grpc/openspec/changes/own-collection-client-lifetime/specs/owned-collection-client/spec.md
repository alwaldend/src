## Purpose

Provide one reusable client owning transport and optional native collection
storage, with explicit RPC access and observable final completion.

## ADDED Requirements

### Requirement: Own transport and optional storage through one Client

Factories SHALL return one Client exposing Start, Connection, OpenCollection,
CloseCollection, Stop and Done. Start SHALL prepare its transport once;
OpenCollection SHALL attach private collection state to that same client using
a copied CollectionOptions value. Connection SHALL reject access before startup
and after terminal shutdown begins. Stop SHALL close storage before stopping
the transport, including after failed collection preparation. Collection state
SHALL NOT implement another Client lifecycle.

#### Scenario: Stop a client with an open collection

- **WHEN** a caller starts a client, opens storage, uses its connection and stops it
- **THEN** native collection storage closes before transport shutdown
- **AND** retained connection access fails and repeated Stop remains safe

#### Scenario: Fail collection opening after starting the transport

- **WHEN** collection preparation fails on a started client
- **THEN** OpenCollection reports the failure without exposing prepared access
- **AND** CloseCollection or Stop can clean up its recorded storage state

### Requirement: Preserve publication and bounded cleanup

CollectionOptions SHALL contain archive WriteBack and terminal CloseTimeout
settings. Stop SHALL publish enabled archives independently of its caller's wait
context, under the configured publication budget (default ten minutes) and a
separate five-second transport budget. A cancelled wait SHALL NOT cancel owned
cleanup. Joining collection cleanup already running in another caller SHALL honor
the terminal cleanup budget and SHALL NOT prevent transport shutdown after that
budget expires. Done SHALL close only after terminal native transport cleanup finishes.
Start's context SHALL control startup only. Hosting applications SHALL retain
source-reservation and session-context policy without another returned client.

#### Scenario: Cancel waiting for cleanup

- **WHEN** a caller stops with an expired wait context
- **THEN** Stop reports cancellation while cleanup continues within its owned budgets
- **AND** Done permits callers to observe final native completion

#### Scenario: Stop while another caller is closing storage

- **WHEN** CloseCollection is already waiting for active native work with an
  unbounded context and terminal Stop reaches its collection cleanup deadline
- **THEN** Stop proceeds to transport shutdown within its separate transport budget
- **AND** its cleanup diagnostic preserves the collection deadline and Done still
  observes final native completion
