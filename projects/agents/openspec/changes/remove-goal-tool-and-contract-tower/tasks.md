## 1. Relocate skill discovery

- [x] 1.1 Add `.agents/BUILD.bazel` with the `skill_discovery_links`
      declaration, `discovery_dir = ".agents/skills"`, and a
      `skill_discovery` package group for `//.agents`.
- [x] 1.2 Repoint every skill `BUILD.bazel` visibility label from
      `//tools/agents:skill_discovery` to `//.agents:skill_discovery`.
- [x] 1.3 Confirm `//.agents:write_skill_links_test` passes and
      `//.agents:write_skill_links` reconciles 27 links with no diff.

## 2. Remove the control tower

- [x] 2.1 Delete `tools/agents` (catalog compilers, capsule, control status,
      evidence tool, registry, catalogs, and shared libraries).

## 3. Reduce the project to its skills

- [x] 3.1 Delete `docs/current-state.md`, `docs/architecture.md`, and
      `docs/roadmap.md`; the project is a collection of skills, not an
      architecture program.
- [x] 3.2 Drop the `docs` dependency from `//projects/agents:docs` and retarget
      the README, root README, and capability-spec links.
- [x] 3.3 Update the `openspec` and `bazel-rules-skill` skills to the
      `//.agents` targets and drop the goal-tool deprecation wording.
- [x] 3.4 Drop the removed catalog entries from `.gitattributes`.

## 4. Validate

- [x] 4.1 `bazel_agent bazel test //.agents:write_skill_links_test` passes; the
      updater reconciles 27 links with no diff.
- [x] 4.2 `bazel_agent bazel test //tools/repo_quality:repo_quality_test`
      passes.
- [x] 4.3 `bazel_agent bazel test //:buildifier_test` passes.
- [x] 4.4 `bazel_agent bazel build` of `//projects/agents/...`, `//.agents/...`,
      `//projects/hugo_landing/...`, `//infra/dns/...`, and
      `//projects/alwaldend.com/...` succeeds (474 targets).
- [x] 4.5 No surviving reference to `tools/agents` outside git history,
      archived OpenSpec records, and the intentional "removed" notes.
