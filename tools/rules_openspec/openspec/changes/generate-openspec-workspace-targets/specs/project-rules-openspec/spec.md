## ADDED Requirements

### Requirement: Generate workspace targets with Gazelle

The OpenSpec Gazelle extension SHALL generate a workspace source target and an owner-local validation target for a directory named `openspec` that contains `config.yaml`.

#### Scenario: Generate a workspace package

- **WHEN** Gazelle visits an owner-local OpenSpec workspace
- **THEN** its BUILD file declares a source target for the configuration, specs, and change artifacts
- **AND** the validation target runs strict current and archive checks against that source

#### Scenario: Regenerate an existing package

- **WHEN** Gazelle revisits a workspace with generated targets and manual attributes
- **THEN** it preserves the manual attributes and unrelated rules
- **AND** another run with the same inputs leaves the BUILD file unchanged
