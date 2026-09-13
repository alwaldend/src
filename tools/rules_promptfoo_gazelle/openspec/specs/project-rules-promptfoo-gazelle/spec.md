# Rules Promptfoo Gazelle

## Purpose

Generate offline validation targets for conventional Promptfoo configurations
under a package's `evals/` directory. This source baseline was observed on
2026-09-08 at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[generator](../../../gazelle/gazelle.go), and
[generator tests](../../../gazelle/gazelle_test.go).

## Requirements

### Requirement: Generate only offline validation rules

The extension SHALL recognize `promptfooconfig` and
`promptfooconfig.<variant>` files with `.yaml`, `.yml`, or `.json` extensions
directly under `evals/` and generate `promptfoo_validate_test` rules in their
containing package. It SHALL neither create an `evals` subpackage nor generate
live `promptfoo_test` rules.

#### Scenario: A conventional default configuration is discovered

- **WHEN** a package contains `evals/promptfooconfig.yaml`
- **THEN** the extension proposes `eval_config_test` in the containing package.

### Requirement: Generate deterministic collision-resistant target names

Variant target names SHALL lowercase their variants, replace runs of
non-ASCII-alphanumeric characters with underscores, and trim leading and
trailing underscores. A nonempty variant that normalizes to an empty string
SHALL use `variant`. Names SHALL take the form `eval_<normalized>_config_test`,
inserting a stable path-derived `__<8 lowercase hex digits>` before
`_config_test` when normalized names collide.

#### Scenario: Two configuration variants sanitize to the same name

- **WHEN** distinct conventional paths yield the same sanitized target name
- **THEN** each colliding target is named `eval_<normalized>__<hash>_config_test` using its stable path-derived hash.

#### Scenario: A variant contains only punctuation

- **WHEN** a nonempty conventional configuration variant normalizes to an empty string
- **THEN** its normalized variant is `variant`, retaining the distinction from the default `eval_config_test` configuration.

### Requirement: Infer skill and data dependencies from package contents

Ordinary validation targets in a named package containing `SKILL.md` SHALL
stage `:skill`; the exact `no_skill` variant and root packages SHALL omit that
inferred skill label. Non-configuration regular files under `evals/` SHALL be
included in sorted `data` unless they belong to a nested Bazel package.

#### Scenario: A skill package includes a no-skill control

- **WHEN** the package contains both a conventional main config and `promptfooconfig.no_skill.yaml`
- **THEN** the main target infers `:skill` and the control target omits it.

### Requirement: Reconcile discovered targets without destructive inference

The extension SHALL reconcile `config`, `data`, and `skills` on matching
conventional validation targets while preserving manual attributes and Gazelle
keep directives. Missing configurations or scan errors SHALL NOT cause the
extension to delete existing targets.

#### Scenario: A previously configured file is removed

- **WHEN** Gazelle no longer discovers a configuration that an existing target references
- **THEN** the extension leaves that target for explicit maintenance.

### Requirement: Bound evaluation discovery to the workspace

Discovery SHALL refuse packages whose resolved paths escape the repository,
refuse a symlink used as the package's `evals` directory, and stop traversal at
Bazel package boundaries, including configured BUILD filenames.

#### Scenario: The evals directory is a symlink

- **WHEN** the package's `evals` path is a symbolic link
- **THEN** discovery returns a non-destructive no-op without generating or removing targets.
