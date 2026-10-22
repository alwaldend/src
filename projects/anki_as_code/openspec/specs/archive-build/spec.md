# archive-build Specification

## Purpose

Build a separate Anki archive through staged publication while preserving the
caller's original collection package and the supplied storage's guarantees.

## Requirements

### Requirement: Staged archive publication

The build API SHALL copy a regular `.colpkg` to a private directory on the
destination filesystem, invoke the caller's apply operation on that copy, then
publish it with the supplied filesystem's rename operation after success and a
cancellation check. Publication SHALL be atomic with the default OS filesystem
and other adapters that provide atomic replacement; a custom adapter's rename
contract determines its publication and failure guarantees.
The caller SHALL close its native session before returning and coordinate
exclusive source and destination access. The output parent MUST already exist.

Staging cleanup is attempted after success as well as failure. Application build uses its ordinary apply workflow on the staged copy; the builder owns no native lifecycle or configuration processing.

#### Scenario: Changed archive

- **WHEN** the injected operation successfully changes the staged collection
- **THEN** the destination contains those changes and the source bytes remain unchanged

#### Scenario: Failure or cancellation

- **WHEN** copying or applying fails, or cancellation is observed before publication
- **THEN** the source and any existing destination remain unchanged and staging cleanup is attempted
- **AND** errors retain their original causes and any cleanup failure

#### Scenario: Publication failure

- **WHEN** the supplied filesystem's rename operation fails
- **THEN** the error retains its cause and staging cleanup is attempted
- **AND** the source remains unchanged and destination preservation follows the adapter's rename contract
- **AND** a failed default OS publication preserves any existing destination

#### Scenario: Invalid source or source alias

- **WHEN** the source is not a regular collection package or the destination aliases it
- **THEN** the operation fails before calling apply

#### Scenario: Unchanged archive

- **WHEN** apply leaves the staged archive bytes unchanged
- **THEN** the published copy has exactly the original bytes

#### Scenario: Build independently of an open source

- **WHEN** a caller builds from a source currently open in a scope and later changes that source scope
- **THEN** build returns with its completed independent output, later changes affect only the source and the output retains the built state
