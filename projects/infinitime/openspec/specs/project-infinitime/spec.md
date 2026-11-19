# InfiniTime Specification

## Purpose

Package a pinned InfiniTime firmware fork documented to add a Text watchface
and a Pomodoro application. This baseline records the repository integration
at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08;
it does not verify firmware execution on a watch.

Sources: [project README](../../../README.md),
[firmware dependency declarations](../../../include.MODULE.bazel),
[build target](../../../BUILD.bazel),
and [release packaging](../../../releases/BUILD.bazel).

## Requirements

### Requirement: Pinned firmware source

The root workspace SHALL declare the firmware fork as a commit-pinned Git
repository with recursive submodules and the repository-owned CMake toolchain
patch.

#### Scenario: Resolve firmware source

- **WHEN** a consumer references a firmware target
- **THEN** Bazel's repository declaration SHALL select the declared commit,
  initialize recursive submodules, and apply the CMake toolchain patch.

### Requirement: Isolated firmware dependency extensions

The firmware integration SHALL use separate pip and npm extension identities
so unrelated Python and npm consumers do not require firmware lockfiles, while
retaining the firmware's declared SDK and compiler archive integrity pins.

#### Scenario: Resolve an unrelated package

- **WHEN** an unrelated root-workspace Python or npm target resolves its own
  dependency extension
- **THEN** that extension SHALL not use the firmware's lockfile as an input.

### Requirement: DFU artifact packaging

The project build test SHALL target the pinned firmware's `dfu` target, and
release packaging SHALL include that artifact and the local patch files.

#### Scenario: Assemble the firmware release inputs

- **WHEN** the project release file target is assembled
- **THEN** its inputs SHALL include the firmware DFU output and the project's
  patch package.
