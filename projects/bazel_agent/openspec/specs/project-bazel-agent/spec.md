# Bazel Agent Specification

## Purpose

Describe the validated Bazel entry point and cached repository-control tool
runner used by repository agents.

Sources: [project README](../../../README.md),
[Bazel targets](../../../BUILD.bazel),
[command and diagnostic implementation](../../../cmd/bazel_agent/main.go),
and [cache regression tests](../../../cmd/bazel_agent/tool_test.go).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08.

## Requirements

### Requirement: Validate and preserve Bazel arguments

`bazel_agent bazel` SHALL require a recognized Bazel command and insert
`--config=agent` immediately after that command. Subsequent options, targets,
and arguments after a `run` separator SHALL pass through unchanged.

#### Scenario: Run a target with application arguments

- **WHEN** the caller invokes `bazel_agent bazel run //path:target -- argument`
- **THEN** the runner constructs `bazel run --config=agent //path:target -- argument`.

### Requirement: Replace the runner with Bazel

The Bazel entry point SHALL resolve `bazel` from `PATH` and replace the runner
process with it, preserving the supplied environment without creating or
injecting a host temporary directory.

#### Scenario: Bazel receives a termination signal

- **WHEN** a validated invocation starts Bazel
- **THEN** Bazel owns the runner process and its resulting exit status directly.

### Requirement: Cache only matching runnable tool inputs

Cached-tool execution SHALL derive its key from the runner, declared tool
sources, dependency and Bazel pins, applicable rc inputs, and platform. Cache
misses for one key SHALL share an exclusive lock and atomically install a
complete runnable output. Changed inputs during a build SHALL select a new key.

#### Scenario: A relevant source file changes during cache population

- **WHEN** inputs rehashed after the build differ from the requested cache key
- **THEN** the output is not published as a valid entry for the stale key.

### Requirement: Report bounded installation diagnostics

`doctor` SHALL emit a read-only JSON report containing runner and available
built-source identities, Bazelisk pins, platform, configured profile and rc
paths, and task-scratch classification without dumping environment values.
Unavailable built-source identity SHALL remain explicit.

#### Scenario: Inspect a workspace without a built runner artifact

- **WHEN** `doctor` cannot find the expected Bazel output for the runner
- **THEN** the source report records that it is unavailable rather than claiming
  that the installed runner matches current source.
