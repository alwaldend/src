# Rules skill Gazelle

## Purpose

Generate Bazel skill bundle declarations from named directories containing
`SKILL.md`. This source baseline was observed on 2026-09-08 at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[generator](../../../gazelle/gazelle.go), and
[generator tests](../../../gazelle/gazelle_test.go).

## Requirements

### Requirement: Generate skill libraries in named subpackages

The extension SHALL propose `skill_library(name = "skill")` in every named
subpackage containing a regular `SKILL.md`, including directories without an
existing BUILD file. A root-level `SKILL.md` SHALL NOT trigger generation.

#### Scenario: A new named skill has no BUILD file

- **WHEN** Gazelle visits a non-root directory containing `SKILL.md`
- **THEN** the extension supplies a skill rule that allows Gazelle to create the package's BUILD file.

### Requirement: Exclude build definitions and evaluation files from bundles

Newly generated skill source globs SHALL include `**` while excluding
`BUILD.bazel`, `BUILD`, every custom BUILD filename configured in Gazelle, and
`evals/**`, with duplicate exclusion names removed.

#### Scenario: The repository configures an additional BUILD filename

- **WHEN** Gazelle generates a skill library using that repository configuration
- **THEN** the configured filename is excluded from the generated source bundle along with the standard BUILD filenames and evaluation subtree.

### Requirement: Preserve manual skill declarations

The extension SHALL preserve existing attributes on a manually maintained
`skill_library(name = "skill")` and SHALL NOT delete that rule solely because
`SKILL.md` disappears.

#### Scenario: A manual skill rule has an explicit source list

- **WHEN** Gazelle runs on the package again
- **THEN** the existing source list remains instead of being replaced with the generated glob.

### Requirement: Respect the apparent skill repository name

The extension SHALL load `skill_library` from the apparent Bzlmod name of
`rules_skills`, falling back to `rules_skills` when no mapping is available.

#### Scenario: A consumer maps the skill module to another apparent name

- **WHEN** Gazelle resolves that module mapping
- **THEN** the generated load names the mapped repository.
