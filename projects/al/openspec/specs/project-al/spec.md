# Al Specification

## Purpose

Describe AL, the repository command runner that prepares plugin resources and
environment variables, executes a command, and releases those resources.

Sources: [project README](../../../README.md),
[lifecycle contract](../../../docs/README.md),
[command implementation](../../../cmd/al/cmd.go),
[lifecycle manager](../../../pkg/lifecycle/lifecycle.go),
[Bazel package](../../../BUILD.bazel), and
[toolchain registration](../../../include.MODULE.bazel).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08.

## Requirements

### Requirement: Prepare selected plugins before command execution

`al run` SHALL load its supplied configurations, start the plugins selected by
`--plugin_label`, and append their environment values to the child command's
environment before starting that command. It SHALL reject invocation without
a command argument.

#### Scenario: Execute a configured command

- **WHEN** configuration loading and selected-plugin startup succeed
- **THEN** the requested command receives its arguments, inherited environment,
  Bazel runfiles environment, and prepared plugin environment.

### Requirement: Roll back failed resource startup

The lifecycle manager SHALL cancel sibling starts after a startup failure and
attempt cleanup for every attempted resource, including the resource that
failed to start. It SHALL return startup and cleanup failures to its caller.

#### Scenario: A plugin fails during startup

- **WHEN** one selected plugin cannot start
- **THEN** AL does not execute the requested command and releases resources
  allocated during the unsuccessful startup.

### Requirement: Retain resources until the command stops

AL SHALL wait for its command to exit before stopping plugin resources. On
cancellation it SHALL request direct-child termination, allow up to ten seconds
for shutdown, and kill and reap the child if the grace period expires.

#### Scenario: Cancel a command using plugin resources

- **WHEN** the invocation context is cancelled while the command is running
- **THEN** the command can use its prepared plugin resources during graceful
  shutdown, and plugin cleanup follows command termination.

### Requirement: Report cleanup failure

Dependency-ordered lifecycle resources SHALL stop in reverse registration
order. AL SHALL include cleanup errors in its invocation result even when the
requested command itself succeeded.

#### Scenario: Cleanup fails after successful command execution

- **WHEN** the command exits successfully but a plugin shutdown returns an error
- **THEN** `al run` reports failure rather than treating cleanup as successful.
