# goal-publication-007: Context capsule and runtime isolation

## Goal

Complete the `context-capsule` and `runtime-isolation` Phase 2 criteria with
the bounded offline context command and the fixed runtime control kernel plus
cross-process task isolation, then rebind the resource baseline to the
revised registry.

## Criterion

> A bounded offline context command returns the path-, label-, or
> task-relevant situation with provenance, structured unavailable states, safe
> next discovery actions, and a human rendering from the same data.
>
> The fixed runtime control kernel and cross-process task isolation remain
> available under optional package failure, timeout, stale revision, and
> untrusted-package fault fixtures.

## Plan

### Context capsule (`tools/agents/cmd/agent_system`)

- `v1alpha1.ContextCapsule` API types with canonical JSON encoding
  (`CanonicalContextJSON`), decode with validation (`DecodeContextCapsule`),
  and human rendering (`RenderContextMarkdown`) from the same data.
- Offline join of the six checked catalogs (topology, policy, action,
  capability, workspace-check, goal) plus the `AGENTS.md` authority document
  for a path/task.
- Structured unavailability: missing catalogs/docs yield
  `Completeness: partial` with bounded limitations instead of failure.
- Fully offline: zero network or provider discovery calls; deterministic
  observedAt; shared capsule ID across JSON and Markdown projections.
- Bazel targets: go binary, unit tests, checked `agent_system_check` and
  regenerate `agent_system_update`.

### Runtime isolation (`tools/agents/control`)

- `Kernel` with per-package states, per-package deadlines, and a never-settling
  timeout translation, plus `Health()`.
- Cross-process `Lock`/`Unlock` (O_EXCL lock, ErrLocked/ErrNotLockOwner),
  `PublishLease`/`VerifyLease` (ErrLeaseExpired/ErrStaleRevision),
  expected-revision `PublishAsset` (compare-and-swap), `PublishAssetIfMissing`,
  `ReadAsset`, `PurgeAsset`, namespace isolation, and `Snapshot`/`ReadSnapshot`
  (offline-readable package status under the control root).
- `KernelOptions.ObserveNow` for deterministic tests; runtime ID auto-generated
  via crypto/rand.
- `tools/agents/cmd/control_status` offline status render (`--markdown`).

### Registry and baseline

- `tools/agents/declarations/registry.json`: add `agent-system` and
  `control-status` direct binaries, bump `criteriaRevision` 3 -> 4, preserving
  the compact original formatting.
- `tools/agents/declarations/resource_baseline.json`: rebind to the new registry
  digest; record measured `discovery_calls: 0` and `context_bytes: 2613`
  (bounded capsule); keep unavailable metrics (`cold_duration_ms`,
  `warm_duration_ms`, `reused_checks`) as null with reasons, per
  `phase1_check` baseline validation.
- Regenerate all checked catalogs and re-run `phase1_check` to confirm the
  baseline binding.

## Verification

- `bazel_agent test //tools/agents/... //projects/goal/...` — all suites.
- `bazel_agent test //:buildifier_test` — BUILD hygiene (new packages).
- `phase1_check` report with criteriaRevision 4, matching registry digest.
- `git diff --check` — clean.
