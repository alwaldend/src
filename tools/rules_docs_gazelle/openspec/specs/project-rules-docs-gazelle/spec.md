# Rules docs Gazelle

## Purpose

Generate documentation packaging declarations in existing Bazel packages.
This source baseline was observed on 2026-09-08 at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[generator](../../../gazelle/gazelle.go), and
[generator tests](../../../gazelle/gazelle_test.go).

## Requirements

### Requirement: Generate only for an existing documented package

The extension SHALL generate `docs_filegroup(name = "docs")` only when a
directory already contains a BUILD file and a regular `README.md`.

#### Scenario: A README has no owning BUILD file

- **WHEN** Gazelle visits a directory containing `README.md` without a BUILD file
- **THEN** the documentation extension generates no rule for that directory.

### Requirement: Use Markdown defaults and ancestor visibility

A newly generated documentation rule SHALL use `glob(["*.md"])` and, when a
nearest ancestor Bazel package is found, grant visibility to that package only.

#### Scenario: A child package has an ancestor BUILD file

- **WHEN** Gazelle generates a new documentation rule below an existing ancestor package
- **THEN** the rule's visibility names that ancestor's `__pkg__` target rather than public visibility.

### Requirement: Preserve manually maintained documentation attributes

The extension SHALL preserve existing manually set rule attributes, including
`srcs`, `deps`, `prefix`, and `visibility`, while allowing missing attributes to
be populated by Gazelle.

#### Scenario: A package uses a custom archive prefix

- **WHEN** a documentation rule already declares a `prefix` and Gazelle runs again
- **THEN** the custom prefix remains in the rule.

### Requirement: Respect the apparent documentation repository name

The extension SHALL load `docs_filegroup` from the apparent Bzlmod name of
`rules_docs`, falling back to `rules_docs` when no mapping is supplied.

#### Scenario: The consumer renames the rules module

- **WHEN** Gazelle reports a different apparent repository name for `rules_docs`
- **THEN** the generated load uses that apparent repository name.
