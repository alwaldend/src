## ADDED Requirements

### Requirement: Validate namespaced Python imports

The repository quality suite SHALL reject from-import statements in tracked
first-party Python source using the existing pinned syntax linter. Module imports
and qualified member access SHALL be accepted. Downloaded external source SHALL
remain outside this first-party check. The Python skill SHALL describe the
convention and its owning validation target.

#### Scenario: Reject a member import

- **WHEN** tracked Python uses a normal, aliased, multiline, relative, or star from-import
- **THEN** repository quality fails with a namespaced-import diagnostic

#### Scenario: Accept namespaced imports and inert text

- **WHEN** source imports modules or includes from-import text only in comments or strings
- **THEN** the import check passes

#### Scenario: Preserve existing consumers

- **WHEN** existing first-party imports are migrated to the enforced convention
- **THEN** their module references and supported package exports remain usable
