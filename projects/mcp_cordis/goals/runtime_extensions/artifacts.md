# Artifact log

[Back to durable goal](./)

## Thin-wrapper strategy reset, 2026-08-30

- [Attempt 10 plan](attempts/010.md): dependency-owned reload serialization
  with deterministic slow-evaluation and slow-activation regressions. Review
  verdict: focused tests pass; complete delivery evidence remains open.
- [Pinned HMR patch](../../patches/hmr@1.0.16.patch): standard pnpm patch that
  serializes the package's module reload task and drains newly stashed URLs.
  Review verdict: exact-package apply check and focused regressions pass.
- [Attempt 9 plan](attempts/009.md): atomic persistence plus official Cordis
  Include/HMR, with the custom source-marker and acknowledgement transaction
  removed. Review verdict: refine after reproducing an HMR overlap race.

## Published candidate review reset, 2026-08-30

- [Attempt 8 review corrections](attempts/008.md): fail-closed recursive
  fallback, UTF-8 byte offsets, and valid-prefix HTTP previews. Review verdict:
  accepted locally with focused and complete package evidence.

- [Attempt 7 current record](attempts/007.md): official Loader, Include, and
  HMR architecture using `cordis.yaml` and normal ESM modules. Review verdict:
  active replacement candidate; exact delivery gates remain open.

- [Attempt 6 historical record](attempts/006.md): rejected worker/version-store
  architecture, its package hashes, and focused evidence. Review verdict:
  superseded by the standard Cordis architecture in Attempt 7.
- [Attempt 5 candidate and validation packet](attempts/005.md): exact
  implementation, immutable hashes, and current test/build results. Review
  verdict: rejected by the completed independent review.
- [Attempt 5 plan](attempts/005.md): exact response to independent review's
  semantic and evidence gaps. Review verdict: proceed.
- [Current MCP Cordis README](../..): documents standard Cordis config,
  ordinary reusable modules, disposable `out` modules, HMR, and the bounded
  `ctx.exec()` contract. Review verdict: current working-tree interface.
- [Imported decision-review skill](https://github.com/alwaldend/src/blob/da2085f1807bfea1c7f3979730f6b7df0033fdce/projects/agents/skills/decision-review/SKILL.md):
  immutable PR 24 instruction payload, now packaged with current offline
  validation. Review verdict: Bazel validation passes.
- [Attempt 3 commit](https://github.com/alwaldend/src/commit/7cfef0719075ad372c3bb257ad216b35770356b2):
  immutable published candidate. Review verdict: rejected as final after
  review.
- [PR 32](https://github.com/alwaldend/src/pull/32): evolving delivery vehicle
  for the runtime-extension goal.
- [PR 24](https://github.com/alwaldend/src/pull/24): provenance for the newly
  requested `projects/agents` subtree. Review verdict: import its four scoped
  changes only, merged against the current base.
- [Attempt 4 plan](attempts/004.md): frozen hypotheses, boundaries, and gates
  for the import and initial review corrections. Review verdict: rejected as a
  final candidate by independent source review.

## Final local candidate, 2026-08-30

- [Attempt 3 evidence](evidence.md#attempt-3-verdict): OIDs, direct rebase
  parent, forced test/build/format results, and the historical delivery state.
  Review verdict at the time: accepted locally, then superseded after
  publication by three valid review findings.
- [Attempt 3 record](attempts/003.md): preserved implementation and
  user-facing interface evidence for that historical local candidate. Review
  verdict: accepted at the time, then superseded by review findings.

## Rebased Attempt 2, 2026-08-30

- [Attempt 2 evidence](evidence.md#attempt-2-verdict): exact task-only rebase
  evidence for commit `e3e74cb1` onto `7ad2704c`. Review verdict: accepted as
  rebase evidence, not as the final candidate because integrated lifecycle
  evidence was nondeterministic.
- [Attempt 2 record](attempts/002.md): preserved documentation and evidence
  for the rebased server and runtime interface. Review verdict: implementation
  retained for Attempt 3; final validation was pending.

## Attempt 1

- [Attempt 1 record](attempts/001.md): preserved runtime/interface artifact.
  Review verdict: rejected as a final candidate because one included package
  had an undeclared runtime dependency.
