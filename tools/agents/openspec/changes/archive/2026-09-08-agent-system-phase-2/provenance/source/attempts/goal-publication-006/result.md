# goal-publication-006: Capability, WorkspaceCheck, Goal, and Index catalogs

## Outcome

This attempt implements the remaining `bounded-catalogs` slices and the
`system-index` criterion: the Capability, WorkspaceCheck, Goal, and
AgentSystemIndex compilers with checked JSON+Markdown artifacts, deterministic
provenance, bounds, conflicts, completeness checks, and Bazel drift gates.

## Implemented

- `tools/agents/cmd/capability_check`: capability compiler over registry
  skills, runtime tools, direct binaries, operation providers, `.agents/skills`
  discovery links, and per-skill `SKILL.md` + `BUILD.bazel` inputs; emits
  `tools/agents/catalogs/capability.{json,md}` with skills and providers.
- `tools/agents/cmd/workspace_check`: workspace-check compiler over tracked
  `MODULE.bazel` roots, `.bazelignore`, docs aggregation, and the full-check
  script; emits `tools/agents/catalogs/workspace-check.{json,md}` with
  workspace records, projections, and phases.
- `tools/agents/cmd/goal_check`: goal compiler over the registered goals root;
  emits validated `identity` + `coarseStatus` only for fully valid records and
  `unavailable` with a bounded reason otherwise;
  `tools/agents/catalogs/goal.{json,md}`.
- `tools/agents/cmd/index_check`: index compiler over the six catalog
  descriptors; emits the descriptor-only `AgentSystemIndex`
  (`tools/agents/catalogs/index.{json,md}`) with digests and query routes, no
  embedded payloads.
- Shared `AgentSystemIndex` refactor: index conflicts now live on the shared
  catalog envelope so the index remains descriptor-only; catalog descriptors
  carry limitations.
- Registry `generatedArtifacts` entries for the four new checked JSON
  artifacts; resource baseline rebound to the new registry digest.
- Checked drift-gate targets: `capability_check_check`, `workspace_check_check`,
  `goal_check_check`, `index_check_check`; update targets regenerate artifacts.

## Verification

- `bazel_agent test //tools/agents/... //projects/goal/...` — 23/23 pass.
- `bazel_agent test //:buildifier_test` — pass (registry/build edits).
- `git diff --check` — clean.
- `goal validate --goals-root projects/agents/goals` — valid (3 goals).

## Not in scope

- `context-capsule`, `runtime-isolation`, `resource-baseline` run in
  subsequent attempts.
- The pre-existing `phase1_check` `codex-migration` registry/discovery
  mismatch remains a tracked defect outside this attempt's declared slices.
