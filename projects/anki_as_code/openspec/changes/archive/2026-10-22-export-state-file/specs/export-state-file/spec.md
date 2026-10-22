## Purpose

Make initial collection onboarding reproducible by exporting native bindings together with editable configurations.

## ADDED Requirements

### Requirement: Optional snapshot bindings

Collection export SHALL accept `--state-output FILE` and write version-1 TOML state for every exported deck, note and note type, including native field/template identities and referenced preset aliases. Configs and bindings MUST describe the same snapshot and generated paths. Export SHALL leave source data and stored state unchanged. Omitting the flag SHALL preserve config-only behavior.

#### Scenario: Whole-collection onboarding

- **WHEN** a user exports a collection with a state destination and installs that file using `col state update`
- **THEN** the exported configs are bound to existing native objects without individual path/ID arguments
- **AND** documentation explains export, state installation, editing, planning and applying in that order

#### Scenario: Namespace selection

- **WHEN** export selects a state namespace with existing renamed paths or preset aliases
- **THEN** generated state uses the actual exported paths and aliases, including newly discovered objects

### Requirement: Safe state destination

Export SHALL reject state destinations aliasing the source collection or conflicting with configuration paths before writing output. State serialization SHALL finish before config writes, and the state file SHALL be published only after successful config export. Filesystem errors MAY leave partial configs and SHALL be reported.

#### Scenario: Rejected destination

- **WHEN** the state destination aliases the source or a config declaration
- **THEN** export fails without changing the source or exported files

#### Scenario: Failed config export

- **WHEN** writing exported configs fails
- **THEN** an existing state output file remains unchanged
