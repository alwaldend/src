## ADDED Requirements

### Requirement: Shared repository CI entry point

Every repository CI workflow SHALL use the local `tools/ci/action` action
after checkout. The action SHALL launch `bazel run --config=ci //tools/ci`
without a shell and propagate its failure.
The entry point SHALL build and test all normal targets in the root and nested
Bazel workspaces, continuing to attempt both phases and remaining workspaces
after a failure and returning failure if any command fails. It SHALL discover
nested workspace boundaries from the root `.bazelignore` and `MODULE.bazel`
files. Runner infrastructure acceptance probes SHALL NOT be part of normal CI.

#### Scenario: Run normal CI

- **WHEN** a repository CI workflow runs
- **THEN** the shared command builds and tests each workspace with the CI profile
- **AND** normal Bazel exclusions for manual and incompatible targets apply

#### Scenario: One build fails

- **WHEN** one workspace's build phase fails
- **THEN** its test phase and the remaining workspaces are still attempted
- **AND** the entry point reports failure
