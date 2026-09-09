## Why

`rules_skills` can only package skills whose files already exist in the calling
repository's source tree, so every externally maintained skill must be
hand-copied and then drifts silently from its upstream. The pinned OpenSpec CLI
is a concrete case: its twelve `openspec-*` skills are vendored as plain files
and were already out of date relative to the pinned `1.11.0` release.

## What Changes

- Add an `http_archive`-based skill source so a tagged upstream archive
  provides skill files without checking them in first.
- Add a `skill_archives` rule that turns the skill directories inside an
  external archive into ordinary `skill_library`-compatible targets, so
  discovery, validation, and installation treat them like local skills.
- Add a `skill_materialization` macro that writes those external skill files
  into the calling repository as regular files with `write_source_files`,
  producing an updater and an up-to-date test.
- **BREAKING** for maintainers: vendored copies become generated outputs that
  must be refreshed with the updater target rather than edited by hand.

## Capabilities

### New Capabilities

- `skill-materialization`: materialize skills owned by an external archive into
  the consuming repository as synchronized regular files.

### Modified Capabilities

- `project-rules-skill`: extend the packaging and discovery contract so skill
  libraries may originate from an external archive and still participate in
  validation and source-tree discovery.

## Impact

- `projects/rules_skills`: new rules, macros, docs, tests, and spec deltas.
- Consumers: `third_party/openspec` declares the archive; its `SKILL.md` files
  become materialized outputs.
- `.agents/skills` discovery continues to link canonical roots and is
  unaffected in shape.
