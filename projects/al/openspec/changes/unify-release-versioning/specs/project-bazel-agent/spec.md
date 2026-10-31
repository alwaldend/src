## REMOVED Requirements

### Requirement: Validate and preserve Bazel arguments

**Reason**: AL no longer provides generic Bazel forwarding; user rc files own agent defaults.
**Migration**: Invoke Bazel directly with its ordinary command and startup options.

### Requirement: Replace the runner with Bazel

**Reason**: Direct Bazel invocation already owns its process and exit status.
**Migration**: Replace generic al bazel forwarding with direct bazel calls and configure user preferences in a personal rc file. Use al bazel test only when affected-test selection is intended.

## ADDED Requirements

### Requirement: Select affected tests with Bazel's dependency graph

`al bazel test` SHALL use Git changes against an explicitly selected base to
select affected tests through Bazel's dependency graph within the selected
workspace and target scope. The comparison SHALL include committed branch
changes since the merge base and staged, unstaged, and nonignored untracked
files. Bazel SHALL own dependency interpretation, test execution, configuration,
and caching. AL SHALL NOT infer dependencies from file proximity, maintain a
parallel dependency graph, or require delivery preparation or publication state.

#### Scenario: Select a transitive test dependency

- **WHEN** a changed source file is a transitive dependency of a test in the selected scope
- **THEN** the selected test list includes that test according to Bazel's graph
- **AND** ordinary Bazel test execution determines whether cached results can be reused

#### Scenario: Keep workspace coverage explicit

- **WHEN** the operator selects one Bazel workspace and a target scope
- **THEN** selection and test execution are limited to that scope
- **AND** AL does not claim that tests in other standalone workspaces were covered

#### Scenario: Include local edits

- **WHEN** a relevant source change is staged, unstaged, or untracked and not ignored by Git
- **THEN** it participates in the changed-path input before any test execution

### Requirement: Do not silently narrow tests on incomplete graph evidence

Affected-test selection SHALL widen to the complete selected test scope when
deletions, build-configuration changes, or paths that cannot be mapped reliably
prevent safe narrowing. It SHALL report that wider selection and its reason.
A failed query or nonzero query with partial output SHALL fail the command
without treating those partial labels as a complete test selection. A successful
empty selection SHALL require complete successful selection evidence.

#### Scenario: Delete a source file

- **WHEN** a changed source file is absent from the current workspace graph
- **THEN** selection includes the complete selected test scope
- **AND** it does not assume that absence means the file had no dependents

#### Scenario: Change build configuration

- **WHEN** a changed build declaration or configuration can alter the selected dependency graph
- **THEN** AL selects the complete requested test scope and lets Bazel evaluate it

#### Scenario: Receive a partial query result

- **WHEN** a Bazel query emits labels and then exits unsuccessfully
- **THEN** AL fails selection without running that partial list or claiming that no tests were affected

### Requirement: Keep affected-test execution independent

Affected-test execution SHALL propagate test failures and cancellation and SHALL
NOT write delivery validation records or publish source. A dry-run selection
SHALL report the selected tests and selection scope without running tests.
The command SHALL NOT restore arbitrary Bazel argument forwarding.

#### Scenario: Inspect the selected tests

- **WHEN** the operator requests a dry run
- **THEN** AL reports changed paths, selected tests, scope, and whether selection widened
- **AND** it does not execute tests or change delivery receipts

#### Scenario: A selected test fails

- **WHEN** Bazel reports failed test execution
- **THEN** al bazel test fails and does not claim publication readiness

## MODIFIED Requirements

### Requirement: Report bounded installation diagnostics

`doctor` SHALL emit a read-only JSON report containing AL and available
built-source identities, Bazelisk pins, platform, declared rc paths, and task-scratch
classification without dumping environment values. It SHALL NOT claim that a
particular Bazel profile is active. Unavailable built-source identity SHALL
remain explicit.

#### Scenario: Inspect a workspace without a built runner artifact

- **WHEN** doctor cannot find the expected Bazel output for AL
- **THEN** the source report records that it is unavailable rather than claiming that installed AL matches current source

#### Scenario: Declare rc inputs

- **WHEN** the operator supplies rc-file paths
- **THEN** doctor reports the declared paths without claiming to inspect their contents or effective Bazel options
