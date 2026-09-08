# Rules skill

## Purpose

Package named Codex skill bundles, validate their metadata, and reconcile
source-tree discovery links. This source baseline was observed on 2026-09-08
at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[skill provider](../../../skill/internal/skill_library.bzl),
[validation aspect](../../../skill/internal/skill_validation.bzl),
[metadata validator](../../../main/go/validate.go), and
[discovery link rules](../../../skill/internal/skill_discovery_links.bzl).

## Requirements

### Requirement: Package one named skill with explicit logical paths

`skill_library` SHALL require a non-root Bazel package and exactly one
`SKILL.md`. Its `SkillInfo` SHALL derive the skill name from the package's final
segment and map files by paths relative to that package. Sources outside the
owning package and duplicate logical paths SHALL be rejected during analysis.

#### Scenario: A source is owned by another package

- **WHEN** a skill library includes a file whose Bazel owner is outside the skill's package
- **THEN** analysis rejects the file as outside the skill root.

### Requirement: Validate instruction and optional OpenAI metadata

The validation aspect SHALL run its declared executable on `SkillInfo` files
and emit the `skill_validation` output group. Validation SHALL check required
instruction frontmatter, agreement between the skill name and package name,
nonempty instructions, and optional `agents/openai.yaml` metadata.

#### Scenario: The frontmatter name differs from the package name

- **WHEN** the validation output group is built for a skill whose declared name does not match its package's final segment
- **THEN** validation fails with a name diagnostic.

### Requirement: Derive discovery links from canonical source skills

`skill_discovery_links` SHALL derive its expected link set from supplied
`SkillInfo` targets in the same repository, require source `SKILL.md` files,
and use direct relative links to the canonical skill roots.

#### Scenario: A declaration references a generated skill instruction file

- **WHEN** a supplied skill provider identifies a generated `SKILL.md`
- **THEN** discovery-link analysis rejects it instead of generating a source-tree link.

### Requirement: Check the exact discovery directory state

The generated local discovery test SHALL require exactly the declared skill
names and direct relative targets, rejecting missing, extra, stale, absolute,
or indirect links. The updater SHALL hold a sibling lock directory while
reconciling links.

#### Scenario: An undeclared skill link remains in the discovery directory

- **WHEN** the generated discovery test observes a name absent from the BUILD declaration
- **THEN** the test fails until the directory is reconciled.
