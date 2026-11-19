# Rules binary toolchain

## Purpose

Provide Bazel toolchains and runnable targets for packaged executable binaries,
including their runtime files. This source baseline was observed on 2026-09-08
at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[lock parsing](../../../main/bzl/binary_toolchain_lock.bzl),
[module extension](../../../main/bzl/binary_toolchain_extension.bzl),
[archive repositories](../../../main/bzl/binary_toolchain_repo.bzl),
and [runnable wrapper](../../../main/bzl/binary_toolchain_binary.bzl).

## Requirements

### Requirement: Select archives from the declared toolchain lock

The module extension SHALL resolve a requested toolchain by name and version
from its JSON lock and create repositories for the matching archives.

#### Scenario: A requested toolchain has no matching archive

- **WHEN** the requested toolchain is absent or has no archive matching its locked version
- **THEN** extension evaluation fails with a diagnostic instead of selecting another version.

### Requirement: Verify downloaded archive integrity

Archive repository download and download-and-extract actions SHALL compare the
reported integrity against the action's declared integrity and fail on mismatch.

#### Scenario: Downloaded content differs from its lock entry

- **WHEN** a download reports an integrity value different from the lock action
- **THEN** repository creation fails with the expected and observed integrity values.

### Requirement: Preserve runtime files in generated binary targets

Generated native binary and filegroup targets SHALL include the files selected
by each binary's `runtime_files` glob patterns, in addition to the executable.
Runnable toolchain wrappers SHALL merge the selected toolchain's runfiles and
declared data runfiles.

#### Scenario: A packaged executable needs an adjacent runtime asset

- **WHEN** an archive declares a runtime asset through `runtime_files`
- **THEN** the native binary carries that asset in its data, the filegroup exposes it, and a wrapper carries the toolchain runfiles into execution.

### Requirement: Forward wrapper arguments and environment

The runnable wrapper SHALL expand declared location and make-variable arguments,
append invocation arguments, and expose the declared `env` through
`RunEnvironmentInfo`.

#### Scenario: A caller supplies additional command arguments

- **WHEN** a generated wrapper is invoked with arguments after its configured arguments
- **THEN** the selected toolchain binary receives both argument sets in that order.
