## ADDED Requirements

### Requirement: Scan the repository checkout and Git history

Repository CI SHALL run credential detection over the current tracked checkout
and over the Git history reachable from the event revision.

#### Scenario: A credential is in a tracked LFS payload

- **WHEN** CI checks out and materializes a revision's LFS assets
- **THEN** the checkout scanner inspects those assets
- **AND** the Git scan independently inspects the reachable Git history
