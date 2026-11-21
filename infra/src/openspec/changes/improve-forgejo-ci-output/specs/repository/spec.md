## ADDED Requirements

### Requirement: Forgejo repository CI reporting

Repository CI running on Forgejo SHALL provide a collapsible log group for each
Bazel phase and a job summary with the result and elapsed duration of each phase
in each workspace.

#### Scenario: One phase fails

- **WHEN** a Bazel phase fails in Forgejo
- **THEN** its diagnostics remain available in the workspace and phase group
- **AND** the summary marks the phase as failed and lists the results of the
  remaining phases
- **AND** the action reports failure

#### Scenario: Bazel prints failed-test output

- **WHEN** Bazel includes a failed test's output in the Forgejo test log
- **THEN** the test's output is available in its own collapsible section
- **AND** the test phase shows the test result directly when expanded

#### Scenario: Bazel prints routine diagnostics

- **WHEN** Bazel prints setup or progress messages and warnings in Forgejo
- **THEN** the messages and their continuation lines are available in a
  collapsible diagnostics group within the phase
- **AND** errors and their continuation lines are visible when the phase is
  expanded

#### Scenario: The action builds and then runs the CI command

- **WHEN** the Forgejo action builds the CI command
- **THEN** the launcher diagnostics have a collapsible section within the
  launcher build group
- **AND** any launcher build failure is visible when the build group is expanded
- **AND** the group closes before execution of the compiled command

#### Scenario: Run the command locally

- **WHEN** the command runs outside Forgejo
- **THEN** it streams the phase logs as plain text

### Requirement: WASI SDK sysroot directory input

The LLVM toolchain SHALL consume the pinned WASI SDK sysroot for
`wasm32-wasip1` as a single directory artifact.

#### Scenario: LLVM generates the WASI system module map

- **WHEN** Bazel analyzes a WASI compilation that uses LLVM's system module map
- **THEN** the WASI sysroot label resolves to one directory input
