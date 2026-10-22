# skill-materialization Specification

## Purpose

Materialize skills owned by an external archive into a consuming repository as
regular files that stay synchronized with the pinned upstream revision.

## Requirements

### Requirement: Package skills from an external archive

`skill_archives` SHALL derive one skill target per skill directory in a
declared external archive and SHALL expose them as `SkillInfo` providers whose
`root` identifies the skill directory within its owning package and whose
`files_by_path` keys are relative to that directory. It SHALL reject
an archive directory that lacks a `SKILL.md` and SHALL reject duplicate skill
names within one declaration.

#### Scenario: An archive skill directory lacks SKILL.md

- **WHEN** a declared archive skill directory contains no `SKILL.md`
- **THEN** analysis fails naming the offending directory.

#### Scenario: Two archive directories declare the same skill name

- **WHEN** two skill directories in one declaration share a final path segment
- **THEN** analysis fails instead of silently preferring one.

### Requirement: Materialize archive skills as synchronized source files

`skills_write` SHALL copy every file of each skill declared through `archives`
into the selected discovery directory in the calling repository as regular
files. The macro SHALL generate a runnable updater that reconciles the
destination files and a test that fails while any destination file is missing,
stale, or extra.

#### Scenario: An upstream archive revision changes a skill file

- **WHEN** the declared archive revision changes the content of a materialized skill file and the updater has not run
- **THEN** the generated exact-state test fails and reports the differing destination file.

#### Scenario: A destination file is written by hand

- **WHEN** a materialized destination file differs from the archive content that produced it
- **THEN** the generated test fails until the updater reconciles the file.

#### Scenario: The archive gains or loses a skill file

- **WHEN** a file exists in the destination directory but not in the declared archive skill
- **THEN** the updater removes it and the generated test detects the extra file before reconciliation.
