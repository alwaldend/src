# Rules Promptfoo

## Purpose

Run skill evaluations and offline configuration validation using a pinned
Promptfoo CLI in Bazel. This source baseline was observed on 2026-09-08 at
revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`. The bundled real CLI is
intended for glibc-based Linux x86-64; Bazel currently constrains OS and CPU
without expressing a libc constraint.

Sources: [project description](../../../README.md),
[module dependencies](../../../MODULE.bazel),
[CLI declaration](../../../BUILD.bazel),
[runner](../../../promptfoo/private/promptfoo_test.bzl),
and [workspace staging](../../../promptfoo/private/workspace.bzl).

## Requirements

### Requirement: Separate offline validation from live evaluation

`promptfoo_validate_test` SHALL invoke configuration validation and reject both
`env_inherit` and `reuse_codex_login`. `promptfoo_test` SHALL be explicitly
tagged `manual`, `requires-network`, `no-cache`, `no-remote`, and `local`.

#### Scenario: An offline validation target requests host login reuse

- **WHEN** a validation target enables `reuse_codex_login` or declares inherited environment variables
- **THEN** Bazel analysis rejects that declaration.

### Requirement: Stage only the selected skill bundles

The runner SHALL physically copy selected `SkillInfo.files_by_path` bundles to
`workspace/.agents/skills/<name>/` and provide a separate empty judge workspace.
Workspace analysis SHALL reject duplicate skill names and unsafe logical paths.

#### Scenario: A target stages two skills with the same logical name

- **WHEN** two selected skill providers declare the same name
- **THEN** workspace analysis fails and identifies the duplicate name.

### Requirement: Isolate temporary runner state

The runner SHALL create private state below an absolute writable `TEST_TMPDIR`,
reject a temporary root beneath an ancestor containing `.agents`, and clean its
state on normal exit and handled termination signals.

#### Scenario: Test scratch is placed beneath an agent-enabled checkout

- **WHEN** an ancestor of the resolved `TEST_TMPDIR` contains `.agents`
- **THEN** the runner fails before staging its subject workspace.

### Requirement: Preserve explicit evaluation results

Live evaluation SHALL invoke Promptfoo with `--no-cache`, `--no-write`, and
`--no-share`, disable telemetry and update checks through the runner environment,
and request `results.json` in Bazel's undeclared outputs directory.

#### Scenario: A live evaluation produces its result artifact

- **WHEN** the evaluation completes and Promptfoo writes the requested results
- **THEN** `results.json` remains in Bazel's undeclared outputs while the isolated runner state is removed.

### Requirement: Resolve JavaScript dependencies from the pinned lock

The module SHALL resolve its CLI dependencies from the checked-in pnpm lock with
optional dependencies and lifecycle hooks disabled, explicitly declaring the
platform packages needed by its supported real CLI.

#### Scenario: Bazel resolves the Promptfoo npm repository

- **WHEN** the module's npm extension translates the checked-in lock
- **THEN** `no_optional` is enabled and `run_lifecycle_hooks` is disabled.
