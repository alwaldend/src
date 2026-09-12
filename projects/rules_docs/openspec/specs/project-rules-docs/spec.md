# Rules docs

## Purpose

Package Markdown documentation and dependent documentation groups under an
archive prefix. This source baseline was observed on 2026-09-08 at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[documentation macro](../../../docs/defs.bzl), and
[module dependencies](../../../MODULE.bazel).

## Requirements

### Requirement: Provide package-relative documentation defaults

`docs_filegroup` SHALL default `srcs` to the package's `*.md` files and default
the archive prefix to `content/docs/` followed by the current package path.

#### Scenario: A package declares documentation without overrides

- **WHEN** `docs_filegroup` omits both `srcs` and `prefix`
- **THEN** the macro packages the package's Markdown files below its default documentation prefix.

### Requirement: Normalize bare child-package dependencies

The macro SHALL normalize relative dependency names containing no colon into
the named child package's `docs` target, while preserving explicit labels.

#### Scenario: A documentation group includes a child package

- **WHEN** package `parent` declares `deps = ["child"]`
- **THEN** the aggregate includes `//parent/child:docs`.

### Requirement: Keep generation support outside the packaging module

The documentation packaging module SHALL provide its rules without depending
on Gazelle or Go; consumers needing generation SHALL add the separate
`rules_docs_gazelle` module.

#### Scenario: A consumer only needs documentation packaging

- **WHEN** a consumer declares the `rules_docs` module dependency
- **THEN** its module supplies the packaging macro through rules_pkg without adding a Gazelle language implementation or Go toolchain dependency.

### Requirement: Allow preserving nested source paths

`docs_filegroup` SHALL provide an opt-in mode that keeps each source's path
relative to its package, so a package whose sources sit in subdirectories is
not forced to flatten identically named files onto one destination. The
default SHALL remain flattened so existing published paths do not change.

#### Scenario: A package holds identically named sources in subdirectories

- **WHEN** a package declares documentation whose sources include two files
  with the same basename in different subdirectories and opts into preserving
  paths
- **THEN** the macro packages both files at distinct package-relative
  destinations instead of failing analysis

#### Scenario: A package keeps the default behavior

- **WHEN** a package declares documentation without opting into preserving
  paths
- **THEN** each source keeps its flattened basename under the archive prefix,
  as before
