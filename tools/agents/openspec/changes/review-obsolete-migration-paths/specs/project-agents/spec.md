## ADDED Requirements

### Requirement: Review migrations for obsolete paths

Repository review policy SHALL require review-agent to trace superseded
execution paths through callers, configuration, and tests. A finding SHALL
identify a reachable path and its consequence for the requested change.

#### Scenario: A migration leaves a superseded path reachable

- **WHEN** a migration establishes one entry point for a workflow
- **AND** that workflow still invokes the superseded entry point
- **THEN** the reviewer reports the active caller and the resulting mismatch
- **AND** the reviewer preserves compatibility paths required by supported callers
