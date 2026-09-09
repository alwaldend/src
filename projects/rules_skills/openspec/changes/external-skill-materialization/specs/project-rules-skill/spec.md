## REMOVED Requirements

### Requirement: Derive discovery links from canonical source skills

**Reason**: `skill_discovery_links` only installs symlinks to skills already in
the source tree, so it cannot express archive-backed skills or written payload
files. `skills_write` replaces it.

**Migration**: Replace `skill_discovery_links(...)` and its separate
`skill_discovery_links_updater` uses with a single `skills_write(...)`,
declaring source-tree skills under `symlinks` and archive skills under `files`.

### Requirement: Check the exact discovery directory state

**Reason**: The link-only check could not verify written payload files, and it
resolved the workspace by searching an unresolved `SKILL.md` symlink rather
than an explicit marker.

**Migration**: `skills_write` generates both the updater and a check that
requires exactly the declared names and kinds, compares written files
byte-for-byte against their declared sources, and locates the workspace through
a declared `workspace_marker` source file.

## MODIFIED Requirements

### Requirement: Package one named skill with explicit logical paths

`skill_library` SHALL require a non-root Bazel package and exactly one
`SKILL.md`. Its `SkillInfo` SHALL derive the skill name from the package's final
segment and map files by paths relative to that package. Sources outside the
owning package and duplicate logical paths SHALL be rejected during analysis.
`skill_archive` SHALL provide the same `SkillInfo` shape for a skill
directory owned by another repository, deriving the name from its target and
logical paths relative to the declared root.

#### Scenario: A source is owned by another package

- **WHEN** a skill library includes a file whose Bazel owner is outside the skill's package
- **THEN** analysis rejects the file as outside the skill root.

#### Scenario: An archive root lacks SKILL.md

- **WHEN** an archive skill root contains no SKILL.md
- **THEN** analysis fails naming the root and the paths it did find.

### Requirement: Validate instruction and optional OpenAI metadata

The validation aspect SHALL run its declared executable on `SkillInfo` files
and emit the `skill_validation` output group. Validation SHALL check required
instruction frontmatter, agreement between the skill name and the package name,
nonempty instructions, and optional `agents/openai.yaml` metadata.

#### Scenario: The frontmatter name differs from the package name

- **WHEN** the validation output group is built for a skill whose declared name does not match its package's final segment
- **THEN** validation fails with a name diagnostic.
