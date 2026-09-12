## MODIFIED Requirements

### Requirement: Reproducible Bazel development

Repository agent build and test operations MUST use `bazel_agent bazel` in the
owning workspace with pinned dependencies and the repository's shared agent
configuration. Generated dependency and catalog files MUST be updated through
their owning generator.

Every nested module that declares `MODULE.bazel` MUST declare only
dependencies its own sources use, with a resolvable version or an override
that applies when that module is the root. Every nested module directory MUST
appear in the root `.bazelignore` so root target expansion does not cross the
workspace boundary. Each nested workspace MUST build and test standalone.

#### Scenario: Change a shared dependency

- **WHEN** an implementation adds an external build input
- **THEN** its owner records an immutable version and integrity information
- **AND** applicable generated locks and package checks validate the declared
  dependency through the pinned build workflow

#### Scenario: Build a nested module standalone

- **WHEN** a nested Bazel workspace is built or tested on its own
- **THEN** its module graph resolves without depending on root-only overrides
- **AND** the workspace builds and tests through its shared configuration

#### Scenario: Expand root targets across a nested boundary

- **WHEN** the root workspace expands targets beneath a nested module directory
- **THEN** the nested workspace is excluded by the root ignore list
- **AND** root expansion neither loads nor silently omits that module
