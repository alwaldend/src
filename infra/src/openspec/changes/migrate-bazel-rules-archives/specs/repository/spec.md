## ADDED Requirements

### Requirement: Root-owned ruleset development

The repository SHALL build and test its first-party rulesets from the root
Bazel workspace and provide standalone archives with their public rule labels.

#### Scenario: A consumer loads a published archive

- **WHEN** a standalone workspace resolves one of the archives through a registry
- **THEN** it can load the module's existing public rules

#### Scenario: CI runs the repository checks

- **WHEN** repository CI builds and tests the source checkout
- **THEN** the root wildcard selects every migrated ruleset under its
  repository-relative package path

## MODIFIED Requirements

### Requirement: Shared repository CI entry point

Every repository CI workflow SHALL use the local `tools/ci/action` action
after checkout. The action SHALL build a launcher with
`bazel run --config=ci --script_path=… //tools/ci`, execute it without a shell,
and propagate failure. The entry point SHALL run one build and one test from
the root Bazel workspace against normal `//...` targets. Runner infrastructure
acceptance probes SHALL run separately.

#### Scenario: Run normal CI

- **WHEN** the action runs repository CI
- **THEN** the shared command builds and tests the checkout root with the CI profile
- **AND** Bazel applies its normal exclusions for manual and incompatible targets
- **AND** ignored directories are never launched as independent workspaces

#### Scenario: One build fails

- **WHEN** the build phase fails
- **THEN** the command still attempts the test phase and exits unsuccessfully
- **AND** Forgejo receives both phase results and the command attempts Bazel
  server shutdown
