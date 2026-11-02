# Bazel Agent Specification

## Purpose

Describe the validated AL Bazel entry point and installation diagnostics used by repository agents.

Sources: [project README](../../../README.md),
[Bazel targets](../../../BUILD.bazel),
[command and diagnostic implementation](../../../internal/bazelcmd/main.go),
and [runner tests](../../../internal/bazelcmd/main_test.go).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08.

## Requirements

### Requirement: Validate and preserve Bazel arguments

`al bazel` SHALL require a recognized Bazel command and insert
`--config=agent` immediately after that command. Subsequent options, targets,
and arguments after a `run` separator SHALL pass through unchanged.

#### Scenario: Run a target with application arguments

- **WHEN** the caller invokes `al bazel run //path:target -- argument`
- **THEN** the runner constructs `bazel run --config=agent //path:target -- argument`.

### Requirement: Replace the runner with Bazel

The Bazel entry point SHALL resolve `bazel` from `PATH` and replace the runner
process with it, preserving the supplied environment without creating or
injecting a host temporary directory.

#### Scenario: Bazel receives a termination signal

- **WHEN** a validated invocation starts Bazel
- **THEN** Bazel owns the runner process and its resulting exit status directly.

### Requirement: Report bounded installation diagnostics

`doctor` SHALL emit a read-only JSON report containing runner and available
built-source identities, Bazelisk pins, platform, configured profile and rc
paths, and task-scratch classification without dumping environment values.
Unavailable built-source identity SHALL remain explicit.

#### Scenario: Inspect a workspace without a built runner artifact

- **WHEN** `doctor` cannot find the expected Bazel output for the runner
- **THEN** the source report records that it is unavailable rather than claiming
  that the installed runner matches current source.
