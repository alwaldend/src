# goal-publication-005: Policy + Action catalogs and shared schemas

## Outcome

This attempt closes with `refine`. It implements the shared catalog schema
types for all six remaining catalog kinds and ships the Policy and Action
catalog compilers with checked JSON+Markdown artifacts, deterministic
provenance, bounds, conflicts, and drift gates.

## Implemented

- `PolicyCatalog` (independent axes with owning sources, conflict reporting
  without resolution).
- `ActionCatalog` (providers, actions with the closed effect set, aliases).
- `CapabilityCatalog`, `WorkspaceCheckCatalog`, `GoalCatalog`, and
  `AgentSystemIndex` schema types + validation + canonical JSON helpers.
- `policy_check` and `action_check` compilers with `--check` drift gates and
  Bazel `update` / `check` targets; registry generatedArtifacts entries;
  resource baseline registry digest rebound.

## Verification

- `bazel_agent test //tools/agents/...` — 10/10 pass.
- `bazel_agent test //:buildifier_test` — pass.
- `git diff --check` — clean.

## Not in scope

- Remaining `bounded-catalogs` slices: Capability, WorkspaceCheck, Goal, and
  Index compilers (next attempts).
- The pre-existing `phase1_check` `codex-migration` registry/discovery
  mismatch remains (see evidence).
