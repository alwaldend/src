## Purpose

Provide a supported recovery path after partial collection mutation so users
can safely reconcile identities and retry without losing collection data.

## ADDED Requirements

### Requirement: Describe recovery after partial execution

An apply that fails after committing native changes SHALL report failure and
provide sufficient recovery information to identify the affected collection,
state namespace, completed mutations, and next supported action. It MUST NOT
claim that the native changes were rolled back when they remain present.

#### Scenario: A deck is created before a later operation fails

- **WHEN** deck creation succeeds and a later note-type operation fails
- **THEN** the apply reports the failure and identifies how to recover the
  created deck's binding before retrying
- **AND** the user need not guess which native object the failed apply created

#### Scenario: Binding persistence fails after native execution

- **WHEN** all planned native mutations succeed but binding persistence fails
- **THEN** recovery information distinguishes completed native work from the
  failed state write and explains the supported recovery action

### Requirement: Recover identities without damaging collection data

The supported recovery workflow SHALL preserve existing native identities,
scheduling, review history, media, and input configuration. It SHALL allow a
fresh plan and retry after unambiguous recovery, and SHALL refuse ambiguous or
stale recovery instead of silently claiming unrelated resources.

#### Scenario: Recover and retry a partially applied collection

- **WHEN** the user follows the supported recovery steps for an unchanged
  collection and reruns the intended apply
- **THEN** the apply completes without duplicating the previously created deck
- **AND** existing collection data and input configuration remain unchanged
  except for the intended plan's operations

#### Scenario: Collection state changed before recovery

- **WHEN** another mutation makes the recovery identity mapping ambiguous
- **THEN** recovery refuses the unsafe mapping and explains the discrepancy
- **AND** it does not silently adopt or overwrite unrelated collection objects
