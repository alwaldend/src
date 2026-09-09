## Context

`tools/agents` owned seven catalog compilers, an offline context capsule, a
control-status reader, an evidence emitter, a phase-1 registry, and shared
admission/control/plan/evidence libraries. Every one of them was offline and
advisory. The only production consumer was the goal tool, which is being
removed in the sibling repository change.

`tools/agents/BUILD.bazel` also owned `skill_discovery_links`, the declaration
that generates the `.agents/skills/` links making every repository skill
discoverable. That responsibility is generic repository infrastructure and
must outlive the removal.

## Goals / Non-Goals

**Goals:**

- Delete the unused control tower without losing skill discovery.
- Keep canonical skill content, link targets, and the 27 existing symlinks
  byte-identical.
- Leave documentation that describes the removed interfaces accurately.

**Non-Goals:**

- Replacing the catalogs or capsule with another generated projection.
- Changing `rules_skills` behavior or the `skills_write` contract.
- Relocating or re-authoring any skill.

## Decisions

The discovery declaration moves to `.agents/BUILD.bazel`, one level above the
directory it populates.

The rule resolves skill roots and the discovery directory relative to the
workspace root through `BUILD_WORKSPACE_DIRECTORY`, not relative to its own
package, so the generated relative symlinks do not depend on where the
declaration lives. A `write_skills` run after the move reconciled 27 links
with no diff, and `//.agents:write_skills_test` passes.

`.agents/skills/BUILD.bazel` was rejected as the location. The generated
updater and checker scripts scan every entry in the discovery directory and
fail on a non-symlink, so a `BUILD.bazel` there makes the rule fail its own
check. Placing the package in `.agents/` and setting `discovery_dir =
".agents/skills"` keeps that directory symlink-only. Verified by running the
moved target's checker, which reported `unexpected skill discovery entry:
BUILD.bazel` before the move and passed after.

Skill `BUILD.bazel` files previously granted `//tools/agents:skill_discovery`.
They now grant `//.agents:skill_discovery`, because the owning target's package
group must match the package that consumes the skills.

The removed bytes are not copied or restated in this change. They remain in git
history at the removal commit's parent, so duplicating them or their digests
would add no guarantee.

## Risks / Trade-offs

- The removal drops the only checked-in path that projected capability, action,
  policy, topology, and workspace catalogs. Nothing consumed them, so no
  workflow loses an input, but the projections are gone if a future need
  appears.
- Skill discovery now depends on a dot-directory package. This is deliberate:
  it colocates the generator with its output, and `//.agents:write_skills`
  is the single entry point.
- Historical provenance under `projects/agents/openspec/changes/archive/`
  still references `tools/agents` paths. Those bytes are retained evidence and
  are intentionally not rewritten.
