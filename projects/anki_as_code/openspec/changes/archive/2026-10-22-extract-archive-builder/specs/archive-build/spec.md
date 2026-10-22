## Purpose

Build a separate Anki archive without exposing a partial destination or modifying
the caller's original collection package.

## ADDED Requirements

### Requirement: Staged archive publication

The build API SHALL copy a regular `.colpkg` to a private directory on the
destination filesystem, invoke the caller's apply operation on that copy, then
atomically replace the destination after success and a cancellation check.
The caller SHALL close its native session before returning and coordinate
exclusive source and destination access. The output parent MUST already exist.

#### Scenario: Changed archive

- **WHEN** the injected operation successfully changes the staged collection
- **THEN** the destination contains those changes and the source bytes remain unchanged

#### Scenario: Failure or cancellation

- **WHEN** copying, applying or publication fails, or cancellation is observed before publication
- **THEN** the source and any existing destination remain unchanged and staging cleanup is attempted
- **AND** errors retain their original causes and any cleanup failure

#### Scenario: Invalid source or source alias

- **WHEN** the source is not a regular collection package or the destination aliases it
- **THEN** the operation fails before calling apply

#### Scenario: Unchanged archive

- **WHEN** apply leaves the staged archive bytes unchanged
- **THEN** the published copy has exactly the original bytes
