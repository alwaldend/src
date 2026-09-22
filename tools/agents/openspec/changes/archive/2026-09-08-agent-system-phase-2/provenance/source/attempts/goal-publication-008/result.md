# goal-publication-008: Phase 2 acceptance

## Outcome

Accept `agent-system-phase-2`: every required criterion passes against the
exact delivered candidate on branch `t3code/continue-agent-system-1` (head
`3cd288a9`, PR 43).

## Delivered candidate

- Commit `48e02734` — Phase 2 delivery (context capsule, runtime isolation,
  registry/baseline, catalogs, goal-publication-007).
- Commit `3cd288a9` — host_bot `review_model` configuration (auto-review uses
  deepseek v4).
- PR 43: open, base `master`, head `t3code/continue-agent-system-1`, mergeable
  (`MERGEABLE`), title "agents: add policy and action catalogs plus shared
  catalog schemas".
- PR review state: one automated Codex review comment (info template, no
  findings) against commit `ab9f639`; no new review findings or inline
  comments on the delivered commits.

## Acceptance evidence

### `goal-recovery` — pass

- Publication-boundary behavior covered by goal store tests
  (`projects/goal/internal/fsstore`, `publication.go`, `recover`/`doctor`).
- `goal doctor` reports `publicationState: stable` on the live record.

### `legacy-migration` — pass

- Migration fixtures and provenance tests in `projects/goal` (prior attempt).

### `bounded-catalogs` — pass

- Six checked catalogs (topology, policy, action, capability, workspace-check,
  goal) regenerate deterministically with provenance, bounds, conflicts, and
  completeness; drift gates 7/7 pass.

### `system-index` — pass

- `AgentSystemIndex` emits descriptor-only inventory with input digests and
  query routes; no embedded catalog bodies (validated by index compiler tests
  and `Validate()`).

### `context-capsule` — pass

- `agent_system` joins the six catalogs plus AGENTS.md; `Completeness: complete`
  with bounded `byteSize: 2613`; structured `partial` on missing inputs;
  JSON/Markdown parity; zero discovery calls. Smoke output verified against
  `projects/agents`.

### `runtime-isolation` — pass

- `tools/agents/control` kernel: per-package states, deadlines, cross-process
  O_EXCL locks, leases, expected-revision CAS publication, namespace
  isolation, offline snapshots; `control_status` render. Control/status tests
  pass.

### `resource-baseline` — pass

- Baseline bound to registry criteriaRevision 4 + digest
  `sha256:b6dbab…6ba8`; ceilings valid; observations measured
  (`discovery_calls: 0`, `context_bytes: 2613`) and unavailable with reasons;
  `phase1_check` reports the binding.

### `exact-candidate-validation` — pass

- `bazel_agent test //tools/agents/... //projects/goal/...` — 27/27 pass.
- Catalog drift gates 7/7 pass; `//:buildifier_test` pass.
- `git diff --check` clean; working tree clean; branch in sync with origin.
- `goal validate --goals-root projects/agents/goals` — valid (3 goals).
- `goal doctor` — stable.
- Goal record RV 18, `goal-publication-007` closed, acceptance attempt bound
  to criteria revision 2 (all 8 required criteria pass).

## Not in scope

- Indicative `phase1_check` `codex-migration` registry gap (pre-existing,
  unchanged, tracked outside this goal).
- PR review assignment to `simeonwarren` is a human review step; the
  automated review has no findings.
