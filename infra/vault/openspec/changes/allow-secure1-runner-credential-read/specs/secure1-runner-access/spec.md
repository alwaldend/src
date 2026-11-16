## Purpose

Let the owning development deployment consume the existing secure1 Forgejo
identity while keeping credential access out of the worker's job environment.

## ADDED Requirements

### Requirement: Exact credential read

The development AppRole SHALL receive read access only to the secure1 runner
record's KV data path. This grant SHALL NOT include writes, listing or siblings.

#### Scenario: Deploy secure1

- **WHEN** development Ansible authenticates through its own AppRole
- **THEN** it can read the cached secure1 worker identity
- **AND** this policy provides no write or sibling permissions
