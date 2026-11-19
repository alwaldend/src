# Rules skills

## Purpose

Package named Codex skill bundles, validate their metadata, and reconcile
source-tree discovery links. This source baseline was observed on 2026-09-08
at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[skill provider](../../../skill/internal/skill_library.bzl),
[validation aspect](../../../skill/internal/skill_validation.bzl),
[metadata validator](../../../main/go/validate.go), and
[discovery write rules](../../../skill/internal/skills_write.bzl), and
[archive packaging](../../../skill/internal/skill_archive.bzl).

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

`skills_write` SHALL derive its expected entry set from supplied `SkillInfo`
targets, installing `symlinks` labels as direct relative links to canonical
source roots and `archives` labels as regular files copied from their declared
sources. It SHALL require source `SKILL.md` files for symlink entries and
reject duplicate names across both inputs.

#### Scenario: A symlink declaration references a generated skill instruction file

- **WHEN** a supplied symlink skill provider identifies a generated or external `SKILL.md`
- **THEN** analysis rejects it instead of generating a source-tree link.

#### Scenario: An archive skill is materialized

- **WHEN** an `archives` provider identifies a skill packaged from an external archive
- **THEN** the updater writes its declared files as regular files in the discovery directory.

### Requirement: Check the exact discovery directory state

The generated local discovery test SHALL require exactly the declared entry
names and kinds, rejecting missing, extra, stale, absolute, or indirect
entries, and SHALL reject a written entry whose bytes differ from its declared
source. The updater SHALL hold a sibling lock directory while reconciling
entries, and SHALL locate the workspace from a declared `workspace_marker`
source file.

#### Scenario: An undeclared skill entry remains in the discovery directory

- **WHEN** the generated discovery test observes a name absent from the BUILD declaration
- **THEN** the test fails until the directory is reconciled.

#### Scenario: A materialized file is edited by hand

- **WHEN** a written entry's file no longer matches the archive source that produced it
- **THEN** the test fails until the updater reconciles the file.

### Requirement: Keep generated actions out of rule source

Rules SHALL NOT embed templated executable content, such as bash assembled from
string literals or f-strings in a `.bzl` file. A rule that needs an executable
SHALL reference an ordinary checked-in script and pass it data as arguments or
as a declared data file that the script parses.

#### Scenario: A rule needs to run a script

- **WHEN** a rule must reconcile a source-tree directory
- **THEN** it references a checked-in script, supplies the entry and payload data as declared inputs, and contains no generated shell program text.

### Requirement: Generate stardoc for every public rule set

Every public Bazel rule set SHALL declare stardoc documentation for its public
`.bzl` entry points, and that generated documentation SHALL be included in the
repository's main documentation output.

#### Scenario: A public rule set gains a rule

- **WHEN** a rule set exposes a new public rule or macro
- **THEN** its stardoc target covers the declaring file and the main documentation build includes the generated pages.
