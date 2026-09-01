# Policy + Action catalog evidence

## Criterion

> Versioned bounded topology, policy, action, capability, workspace-check, and
> goal catalogs derive from owner-local facts with deterministic provenance,
> conflicts, bounds, and completeness checks.

## Claim (this attempt)

The `bounded-catalogs` criterion is advanced: the shared catalog schema types
now cover Policy, Action, Capability, WorkspaceCheck, Goal, and AgentSystemIndex;
the Policy and Action compilers are implemented, checked as repository
artifacts, and covered by deterministic and negative fixtures. The remaining
Capability, WorkspaceCheck, Goal, and Index compilers run in the next attempt.

## Implemented

### Shared schema (`tools/agents/catalog/v1alpha1`)

- `policy.go`: `PolicyCatalog` with `PolicyRecord`, independent `PolicyAxis`
  values (`known`/`unknown`/`conflict`), owning source path, precedence, and
  canonical JSON + strict decode + validation (non-null arrays, sorted
  set-like IDs, unknown-axis rejection).
- `action.go`: `ActionCatalog` with `ActionProvider`, `ActionRecord` (closed
  effect set, full operation fields), `ActionAlias`, canonical JSON +
  strict decode + validation (unique IDs, provider refs, known effects,
  known cacheability, non-empty gates).
- `system.go`: `CapabilityCatalog`, `WorkspaceCheckCatalog`, `GoalCatalog`
  (available/unavailable with validated identity + coarse status only),
  `AgentSystemIndex` (descriptor-only, no embedded payloads) — schema types
  with validation and canonical JSON helpers used by the next compilers.
- `render_policy.go` / `render_action.go`: human Markdown renders from the
  same validated data, stating the JSON digest.

### Policy compiler (`tools/agents/cmd/policy_check`)

- Reads the root `AGENTS.md`, `CODEOWNERS`, and six top-level boundary
  READMEs as the explicitly closed policy universe (registry has no policy
  authority).
- Derives eight independent axes (source disclosure, evidence handling,
  bazel visibility, build consumer, artifact publication, documentation
  publication, information, live-environment association) with owning source
  paths; disagreements become conflicts.
- Emits checked `tools/agents/catalogs/policy.{json,md}` with a `--check`
  drift gate; wired as `policy_update` / `policy_check_check` Bazel targets.

### Action compiler (`tools/agents/cmd/action_check`)

- Reads the registry `operationFiles` (four paths), their owner-local
  definition files, and the closed effect set.
- Emits `tools/agents/catalogs/action.{json,md}` with providers, actions, and
  aliases; wired as `action_update` / `action_check_check`.

### Registry and baseline

- `registry.json` `generatedArtifacts` registers the policy and action
  checked outputs (inline one-line style preserved).
- `resource_baseline.json` `registryDigest` rebound to the new registry
  digest because the registry file changed.

## Tests

- `v1alpha1/policy_action_test.go`: policy + action canonical round-trip,
  markdown-digest parity, unknown-axis rejection, unknown-effect rejection.
- `cmd/policy_check/main_test.go`: complete policy fixture; `--check` fails
  on a missing boundary README.
- `cmd/action_check/main_test.go`: complete action fixture; `--check` fails
  on a missing operation file; unknown effect rejection.
- Checked drift gates `policy_check_check` and `action_check_check` pass on
  the tracked artifacts.

## Validation run

- `bazel_agent test //tools/agents/...` — 10/10 tests pass (schema, policy,
  action, topology, phase1).
- `bazel_agent test //:buildifier_test` — pass.
- `git diff --check` — clean.

## Known pre-existing issue (not caused by this attempt)

- `phase1_check` reports `valid: false` with `missing: skill:codex-migration`
  because the registry skills list (21) does not include `codex-migration`
  while `.agents/skills` discovery has 22 links. This mismatch exists on
  HEAD's registry and is outside this attempt's scope.

## Next

- Capability, WorkspaceCheck, Goal, and Index compilers.
