# python-import-validation Specification

## Purpose

Enforce module namespace imports in tracked first-party Python source through repository quality.

## Requirements

### Requirement: Validate namespaced Python imports

The repository quality suite SHALL reject individual symbol imports in tracked
first-party Python source using the existing pinned syntax linter, except direct
public re-exports listed in explicit `__all__` in `__init__.py`. Wildcard imports
SHALL remain forbidden. Module imports
and qualified member access SHALL be accepted, including importing a verified
module with `from package.path import module`. Long module paths SHALL use this
form. Downloaded external source SHALL
remain outside this first-party check. The Python skill SHALL describe the
convention and its owning validation target.

#### Scenario: Reject a member import

- **WHEN** tracked Python imports individual members outside explicit package re-exports or uses a star from-import
- **THEN** repository quality fails with a namespaced-import diagnostic

#### Scenario: Accept namespaced imports and inert text

- **WHEN** source imports modules, including verified module from-imports, or includes from-import text only in comments or strings
- **THEN** the import check passes

#### Scenario: Preserve existing consumers

- **WHEN** existing first-party imports are migrated to the enforced convention
- **THEN** their module references and supported package exports remain usable

#### Scenario: Accept explicit package re-exports

- **WHEN** `__init__.py` imports direct symbols and lists their bound names in explicit `__all__`
- **THEN** the import check accepts those exports, including relative imports
