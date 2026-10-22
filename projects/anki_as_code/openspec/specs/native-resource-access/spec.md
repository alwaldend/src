# native-resource-access Specification

## Purpose

Provide shared native resource access without collection lifecycle, configuration projection or reconciliation policy.

## Requirements

### Requirement: Borrow prepared native access

Resource access SHALL reuse a caller-owned prepared connection and provide generated native service clients. It SHALL NOT start, open, close or own the backend or collection, or introduce a duplicate CRUD protocol.

#### Scenario: Reuse a collection after reads

- **WHEN** a snapshot read completes or fails
- **THEN** the caller can continue using the same collection and close it independently

### Requirement: Read complete native resources

A snapshot SHALL contain upstream deck, note-type, note and card messages, including unmanaged objects and resident filtered cards, plus existing preset names. Decks SHALL be ordered by native title and other resources by native ID. Reads SHALL preserve native values without state access, config projection or writes.

#### Scenario: Read a resident filtered card

- **WHEN** a card currently resides in a filtered deck
- **THEN** its snapshot retains the native current and original deck references

### Requirement: Fresh detached results

Each snapshot SHALL perform fresh reads without a persistent cache and return caller-owned results. Read failures SHALL return contextual errors without a partial snapshot. Reads SHALL respect context cancellation. Sequential reads SHALL NOT promise an atomic database snapshot.

#### Scenario: Update between snapshots

- **WHEN** a caller changes a native note after an earlier snapshot
- **THEN** a later snapshot contains the changed value and the earlier snapshot remains unchanged
