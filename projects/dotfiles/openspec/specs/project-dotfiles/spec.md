# Dotfiles Specification

## Purpose

Package personal configuration files with commands to compare and install the
packaged files. This baseline records checked-in behavior at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08; no host
configuration was installed during specification capture.

Sources: [project README](../../../README.md),
[archive and command targets](../../../BUILD.bazel),
[executable mappings](../../../bin/BUILD.bazel),
[home mappings](../../../home/BUILD.bazel),
and [Neovim mappings](../../../nvim/BUILD.bazel).

## Requirements

### Requirement: Installable configuration archive

The dotfiles archive SHALL package the generated Make installation interface
and the `bin`, `home`, and `nvim` configuration groups beneath a `dotfiles`
archive directory.

#### Scenario: Build the distributable archive

- **WHEN** the dotfiles archive target is assembled
- **THEN** its input SHALL be the Make installation package containing the
  three declared configuration groups.

### Requirement: Comparison and installation commands

The project SHALL expose Bazel `help`, `diff`, and `install` entry points that
invoke the matching operations in the shared Make installation package.

#### Scenario: Inspect available configuration differences

- **WHEN** a user selects the project's `diff` entry point
- **THEN** the entry point SHALL invoke the package's `diff` operation.

### Requirement: Selected configuration operations

The archive's Make interface SHALL support group-specific and file-specific
comparison and installation targets.

#### Scenario: Select Neovim configuration

- **WHEN** a user selects `diff/nvim` or `install/nvim` from the unpacked archive
- **THEN** the selected operation SHALL apply to the Neovim configuration group.
