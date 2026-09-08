# goal-publication-005: bounded catalogs (Policy + Action) and shared schemas

## Goal

Complete the `bounded-catalogs` criterion for the Policy and Action slices and
add the shared catalog schema types for all remaining catalog kinds, following
the safe sequence from the catalog-inputs audit.

## Criterion

> Versioned bounded topology, policy, action, capability, workspace-check, and
> goal catalogs derive from owner-local facts with deterministic provenance,
> conflicts, bounds, and completeness checks.

## Plan

### Shared schema types (tools/agents/catalog/v1alpha1)

- `PolicyCatalog`: `policies[]` with pathPrefix, precedence, policy sources,
  owner boundary source, review source, and independent axis values
  (`known`/`unknown`/`conflict` with owning source path).
- `ActionCatalog`: `providers[]`, `actions[]`, `aliases[]` with the closed
  effect set and full operation fields.
- Also add `CapabilityCatalog`, `WorkspaceCheckCatalog`, `GoalCatalog`, and
  `AgentSystemIndex` schema types with validation and canonical JSON helpers
  (used by later milestones).

### Policy compiler (tools/agents/cmd/policy_check)

- Inputs: every tracked `AGENTS.md`, `CODEOWNERS`, top-level boundary READMEs.
- Explicitly closes the root-only AGENTS.md universe (registry has no policy
  authority), so completeness stays truthful.
- Independent axes: source disclosure, evidence handling, bazel visibility,
  build consumer, artifact publication, documentation publication,
  information, live-environment association. Each value carries its owning
  source path; disagreements become conflicts.
- Emits checked `tools/agents/catalogs/policy.{json,md}` with a `--check`
  drift gate and Bazel update/check targets.

### Action compiler (tools/agents/cmd/action_check)

- Inputs: registry `operationFiles` (four paths), their definition files, and
  the closed effect set.
- Emits `tools/agents/catalogs/action.{json,md}` with providers, actions, and
  aliases, plus checked gates and Bazel targets.

### Tests

- Compiler unit tests for complete and negative (missing file, stale JSON,
  conflict) fixtures.
- Schema validation tests for policy/action records.
- Parity render asserts the JSON digest appears in Markdown.

## Out of scope

- Capability, WorkspaceCheck, Goal, and Index compilers run in subsequent
  attempts.
