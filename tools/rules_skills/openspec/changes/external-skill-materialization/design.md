## Context

`rules_skills` packaged only files already present in the calling package. The
pinned OpenSpec CLI is the first consumer that needs skills maintained upstream
at a pinned revision, and its vendored copies had already drifted. The
`skill_discovery_links` macro modeled only symlinks to in-repository skills and
resolved the workspace by dereferencing an unresolved `SKILL.md` symlink.

## Goals / Non-Goals

**Goals:**

- Package skills from a pinned external archive without checking them in first.
- Materialize archive skills into the consuming repository as regular files.
- Verify the source-tree discovery directory exactly, including payload bytes.
- Keep one macro for both symlinked and written discovery entries.

**Non-Goals:**

- Publishing `rules_skills` to a registry.
- Supporting native Windows checkouts.
- Fetching skills over the network at build time; the consumer owns the
  `http_archive` pin and integrity.

## Decisions

- `skill_archive` lives in the archive-owned `BUILD` file and takes the
  archive's own file set, so the rule derives `files_by_path` with the existing
  `SkillInfo` shape and needs no new provider. A macro over several roots would
  have to return many providers, which rules cannot expose per target.
- `skills_write` replaces `skill_discovery_links` rather than extending it: one
  updater and one check cover both symlink and written entries, and the name
  states what it does.
- Written entries are compared with `cmp` against their declared archive source
  instead of being regenerated, so drift and hand edits both fail.
- The check locates the workspace through an explicit `workspace_marker` source
  file instead of searching for an unresolved `SKILL.md` symlink, which no
  longer exists once skills are written.
- Validation passes the provider's skill name instead of the package segment so
  archive skills, whose package differs from the skill name, validate.
- `compatibility` joins the allowed `SKILL.md` frontmatter keys because
  upstream OpenSpec skills declare it.

## Risks / Trade-offs

- The check binds to the local checkout by resolving runfiles into the source
  tree, so it is tagged `local`/`no-sandbox` and cannot be remote-executed.
- Written entries duplicate archive bytes in the repository; the trade-off is
  reviewable, diffable skill content and no network at build time.
- `skills_write` cannot symlink archive skills, because a symlink into the
  external repository would not resolve for agents that read the checkout.
