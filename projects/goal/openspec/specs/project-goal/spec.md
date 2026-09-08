# Goal Specification

## Purpose

Describe the deprecated goal CLI, portable resource API, and local filesystem
store retained for legacy record compatibility and recovery. New maintained
work uses its owning OpenSpec workspace and the `//tools/openspec` entry point.

Sources: [project README](../../../README.md),
[Bazel package](../../../BUILD.bazel),
[command implementation](../../../cmd/goal/command.go),
[store regression tests](../../../internal/fsstore/store_test.go),
and [publication recovery tests](../../../internal/fsstore/publication_test.go).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08; deprecation behavior includes this OpenSpec migration.

## Requirements

### Requirement: Advertise deprecation while preserving compatibility

The legacy goal command SHALL identify itself as deprecated in help and
invocation diagnostics, directing new work to
`bazel_agent bazel run //tools/openspec -- ...`. Existing CLI, API, and store
operations SHALL remain available for legacy compatibility and recovery.

#### Scenario: An existing consumer invokes the goal command

- **WHEN** a consumer runs a retained legacy goal operation
- **THEN** the command emits a deprecation diagnostic on stderr while preserving
  the operation's supported output and behavior.

### Requirement: Disable goal-skill discovery

The goal skill SHALL be absent from enabled repository skill discovery and
registration. Its retained historical source SHALL state that it is disabled
and SHALL disable implicit invocation, with OpenSpec providing the supported
procedure for new maintained work.

#### Scenario: An agent selects a procedure for new maintained work

- **WHEN** the agent reads the enabled repository skill catalog
- **THEN** it finds the OpenSpec procedure without an enabled goal-skill entry.

### Requirement: Preserve legacy resource integrity

The filesystem store SHALL continue validating supported
`goals.alwaldend.com/v1alpha1` resources, enforcing expected resource versions,
and checking SHA-256 bindings for canonical attempt Markdown and evidence.
Record validation SHALL NOT itself establish that acceptance criteria are met.

#### Scenario: Two legacy writers use the same resource version

- **WHEN** one writer commits a mutation and advances the stored version before
  the other writer acquires its goal lock
- **THEN** the second writer's stale expected version is rejected.

### Requirement: Recover interrupted legacy publication

Multi-file mutations SHALL stage exact after-images and record publication
intent. `doctor` SHALL classify incomplete publication, and `recover` SHALL
replay or discard the pending intent according to its stored state. The store
SHALL preserve atomic per-file replacement without claiming independent
cross-file transaction semantics.

#### Scenario: A legacy checkpoint stops after a partial publication

- **WHEN** canonical files have changed but the generated README publication
  fails with a recoverable intent still present
- **THEN** `doctor` identifies the partial intent, ordinary mutation remains
  blocked, and `recover` can complete the intended publication.
