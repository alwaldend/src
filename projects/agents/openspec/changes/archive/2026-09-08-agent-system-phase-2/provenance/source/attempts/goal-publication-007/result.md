# goal-publication-007: Context capsule and runtime isolation

## Outcome

This attempt implements the `context-capsule` and `runtime-isolation`
criteria: a bounded offline context command and the fixed runtime control
kernel with cross-process task isolation, plus the registry/baseline binding
for the two new direct binaries.

## Implemented

### `tools/agents/api/v1alpha1` — context capsule types

- `ContextCapsule`, `CapsuleIdentity`, `TaskBinding`, `CapsuleProvenance`,
  `CapsuleComponent`, `CapsuleCapabilityRef`, `CapsuleCheckRef`,
  `CapsuleProviderRef`, `CapsuleDocument`, `CapsuleNextAction`.
- `CapsuleID` (stable per path/task), `CanonicalContextJSON` (canonical
  ordered encoding), `DecodeContextCapsule` (self-validation), and
  `RenderContextMarkdown` (human projection from the same data).
- Unit tests: canonical encoding determinism, decode validation, and
  JSON/Markdown same-data rendering.

### `tools/agents/cmd/agent_system` — bounded offline context command

- Reads the six checked catalogs (topology, policy, action, capability,
  workspace-check, goal) plus `AGENTS.md` and renders a capsule for a
  path/task with provenance, completeness, limitations, and next actions.
- Fully offline: zero provider discovery calls; deterministic `observedAt`
  (fixed base time) for reproducible tests.
- Structured unavailable states: missing catalogs/docs produce
  `Completeness: partial` with bounded limitations; a complete run emits
  `byteSize: 2613` for `projects/agents`.
- `--json` default, `--json=false` for Markdown; both share the same capsule
  ID.
- Bazel targets: `agent_system`, `agent_system_test`, `agent_system_check`,
  `agent_system_update`; fixtures under `testdata/root` and `testdata/missing`.

### `tools/agents/control` — control kernel and isolation

- `Kernel` with per-package states (`loading`, `ready`, `degraded`,
  `failed`, `timed-out`, `draining`, `disabled`), per-package deadlines,
  and a never-settling timeout translation; `Health()`.
- Cross-process task `Lock`/`Unlock` via O_EXCL lock file under
  `locks/<ns>/task.lock`; `ErrLocked`, `ErrNotLockOwner`.
- Lease: `PublishLease`/`VerifyLease` with `ErrLeaseExpired`,
  `ErrStaleRevision`.
- Candidates: expected-revision `PublishAsset` (compare-and-swap),
  `PublishAssetIfMissing`, `ReadAsset`, `PurgeAsset`, namespace isolation.
- `Snapshot`/`ReadSnapshot` persists `packages.json` under `assets/<ns>/` so
  offline readers see package status.
- `KernelOptions.ObserveNow` for deterministic tests; runtime ID auto-generated
  via crypto/rand otherwise.
- Unit tests covering load/unload transitions, health, deadlines, keyed locks,
  leases, namespace isolation, and stale-revision publication failures.

### `tools/agents/cmd/control_status` — offline status render

- Reads the kernel snapshot/assets from a control root and renders status
  (`--markdown` supported); no provider discovery.
- Unit tests + Bazel targets.

### Registry and baseline

- `tools/agents/declarations/registry.json`: `agent-system` and
  `control-status` direct-binaries entries; `criteriaRevision` bumped 3 -> 4
  (4-line minimal diff; original compact formatting preserved).
- `tools/agents/declarations/resource_baseline.json`: rebound to the new
  registry digest (`sha256:b6db…6ba8`), criteriaRevision 4, measured
  `discovery_calls: 0` and `context_bytes: 2613`, with unavailable
  `cold_duration_ms`, `warm_duration_ms`, and `reused_checks` kept as null
  plus reasons.
- `tools/agents/cmd/phase1_check` criteria revision checks bumped to 4.
- All checked catalogs regenerated.

## Verification

- `bazel_agent test //tools/agents/... //projects/goal/...` — 27/27 pass.
- Catalog drift gates (topology/policy/action/capability/workspace-check/
  goal/index `_check`) — 7/7 pass.
- `bazel_agent test //:buildifier_test` — pass.
- `phase1_check` — criteriaRevision 4, registry digest bound, valid counts,
  `valid: false` only for the pre-existing `codex-migration` registry gap.
- `git diff --check` — clean.

## Not in scope

- `goal-recovery`, `legacy-migration`, `bounded-catalogs`, `system-index`
  were accepted in prior attempts.
- The pre-existing `codex-migration` registry/discovery mismatch remains a
  tracked defect outside this attempt's declared slices and did not regress:
  HEAD's registry also has 21 skills while discovery links expose 22.
