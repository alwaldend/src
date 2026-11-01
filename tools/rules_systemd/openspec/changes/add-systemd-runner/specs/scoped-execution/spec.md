## Purpose

Run Bazel-built executables with explicit resource limits in transient user systemd scopes.

## ADDED Requirements

### Requirement: Scoped command execution

The executable rule SHALL accept a binary, configured arguments, data files,
and systemd properties. The runner SHALL also accept repeatable
`--property=NAME=VALUE` CLI options followed by an executable and its arguments. It SHALL execute the command through
`systemd-run --user --scope`, preserving literal arguments, environment,
working directory, standard streams, interrupt handling, and exit status.
It SHALL require at least one property and reject malformed properties or a
missing executable before launch. It SHALL NOT install a persistent service.

#### Scenario: Run a command with limits

- **WHEN** a declared executable target is run with memory and swap properties
- **THEN** only the launched command and its children enter that scope
- **AND** the command observes the requested kernel cgroup limits
- **AND** configured target arguments and caller arguments reach the executable unchanged
- **AND** application runfiles, runtime environment, and declared data remain available

#### Scenario: Command failure or interruption

- **WHEN** the launched command exits nonzero or receives an interrupt
- **THEN** the runner preserves the command's exit or signal outcome

#### Scenario: Invalid input or unavailable systemd

- **WHEN** the executable or properties are missing, a property is malformed, or systemd cannot launch the scope
- **THEN** the runner exits nonzero with a diagnostic
