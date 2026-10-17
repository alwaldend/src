## Purpose

Describe intended collection changes from editable configurations, captured native resources and stored path bindings without modifying any input or collection.

## ADDED Requirements

### Requirement: Planning uses supplied values without I/O

Planning SHALL accept desired configurations, a native snapshot and stored state, return a deterministic protobuf execution plan with a path-based report, and leave inputs unchanged. Callers MUST provide all three inputs and own their freshness.

#### Scenario: Backend already stopped

- **WHEN** the caller supplies captured native values and loaded configurations after closing its collection
- **THEN** planning returns the same plan on repeated calls without native or filesystem access

#### Scenario: Missing required input

- **WHEN** a required input is absent
- **THEN** planning returns an error and no partial plan

### Requirement: Planning preserves path ownership

Planning SHALL resolve explicit moves through stored bindings, report additions, updates, deletions and moves of managed resources, and reject requests that would damage unmanaged references. Untracked resources and resident filtered cards MUST remain unmanaged.

#### Scenario: Explicit note move

- **WHEN** a bound note is moved through a declaration
- **THEN** its report identifies the source and destination paths without deleting or recreating its native identity

#### Scenario: Untracked native resources

- **WHEN** native resources have no stored binding
- **THEN** omission from configs does not plan their deletion

#### Scenario: Empty deck targeted by an untracked template

- **WHEN** an untracked native note-type template targets a managed normal deck with no cards
- **THEN** planning rejects removal or conversion of that deck while leaving the unmanaged template unchanged

#### Scenario: Cross-note-type member move

- **WHEN** a move attempts to transfer a field or template to another note-type file
- **THEN** planning rejects the request without reusing a destination member's identity or ordinal

### Requirement: Native generation remains deferred

Planning SHALL compare declared content and existing normal-card destinations without applying native methods or simulating card generation and normalization.

#### Scenario: New note

- **WHEN** a valid note is declared without an existing binding
- **THEN** the plan reports the new note and marks generated-card effects as deferred

### Requirement: Planning prepares execution data

The plan SHALL contain native operation payloads, resolved member mappings, desired bindings and expected stored state. Internal runtime models SHALL remain private to the planning package. Execution SHALL NOT require config resolution or native resource discovery.

#### Scenario: Detached plan

- **WHEN** the caller serializes and decodes a prepared plan
- **THEN** it contains the same ordered operations and binding references without the planning package's runtime models

#### Scenario: New referenced objects

- **WHEN** a note and a template reference a note type and deck created in the same plan
- **THEN** prepared references use configuration paths and proposed bindings contain no fabricated native IDs

#### Scenario: Converged collection

- **WHEN** desired values already match the supplied snapshot and stored bindings
- **THEN** the report and all native operation lists are empty

### Requirement: Native member drift preserves identity

Planning SHALL match stable native member IDs when available and original stored member names otherwise, without identifying current members through stale ordinals. Explicit member renames SHALL retain their source identity. Conflicting member mappings SHALL fail before returning a plan.

#### Scenario: Native member reordering

- **WHEN** fields or templates are reordered in Anki after their state is recorded
- **THEN** the plan maps desired members to their current native ordinals without exchanging their identities

#### Scenario: Removed stable member

- **WHEN** a saved stable member ID is absent from the supplied native type
- **THEN** its former ordinal does not claim another surviving member

#### Scenario: Resident filtered cloze cards

- **WHEN** a tracked cloze note has resident filtered cards at any cloze ordinal
- **THEN** planning rejects replacement or removal of its shared template, while permitting an explicit rename that preserves its original native template mapping
