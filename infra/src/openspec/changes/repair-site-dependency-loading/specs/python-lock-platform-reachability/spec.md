## Purpose

Resolve native uv lockfiles according to dependency reachability and platform
markers so unrelated operating-system packages do not trigger invalid builds.

## ADDED Requirements

### Requirement: Preserve reachable platform dependencies

The uv locking configuration SHALL preserve platform-specific dependencies
for Bazel's native importer without patching rules_python. Autoscroll SHALL
resolve macOS and its complement separately, covering all platforms. Generated
package resolution markers SHALL retain PyObjC only for macOS. Existing package
identities and distribution hashes SHALL be preserved when applying this repair.

#### Scenario: Operating-system dependency

- **WHEN** a root dependency references PyObjC only on macOS
- **THEN** Linux and Windows imports omit PyObjC and macOS imports retain it

#### Scenario: Reproducible generation

- **WHEN** the documented native uv update target runs again with the same inputs
- **THEN** its checked-in lock output remains unchanged
- **AND** no custom importer patch or requirements-file export is needed

### Requirement: Preserve build consumers

The repair SHALL preserve the normal Linux autoscroll build and permit native
genquery loading of its dependency closure. The full site and its graph targets
SHALL remain the acceptance consumers for this repair.

#### Scenario: Site dependency capture

- **WHEN** the full site's native graph target builds on Linux
- **THEN** it produces its diagram without attempting a macOS source build

### Requirement: Portable locked Git acquisition

The pinned Bazel version SHALL preserve locked Git crate inputs when a shared
repository-content cache is reused after the original output base is removed.
The normal crate_universe recipe SHALL supply Anki's percent-encoding-iri
dependency without a custom source archive or copied build recipe.

#### Scenario: Reuse after retiring an output base

- **WHEN** a second output base reuses a clean cache from the first Git fetch
- **AND** the first output base is no longer present
- **THEN** the crate manifest, Rust source, and license files remain readable
- **AND** their contents match the locked commit
