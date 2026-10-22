## Context

See `proposal.md`. The reviewed starting candidate is
`e0d24cb3a9c250fc92665395d953f50aab2e5a6a`. Existing task checkmarks identify
archive candidates but do not establish that historical requirements still
describe current behavior. Source and retained acceptance evidence determine
the final baseline; no live deployment is inferred from checked-in definitions.

## Goals / Non-Goals

**Goals:** Make remaining work discoverable, preserve useful history, and leave
baseline specifications consistent with current owners and implementations.

**Non-Goals:** Implement deferred findings, re-run historical deployments,
change infrastructure, or recreate abandoned interfaces to satisfy stale plans.

## Decisions

- Review candidates by owner, collecting pinned CLI status and specification
  instructions before any baseline write or archive move. A completed-looking
  checklist alone is insufficient; use implementation, tests, and recorded
  evidence to decide readiness.
- Resolve overlapping deltas against current implementation. Preserve original
  historical artifacts and record superseded or relocated requirements in an
  archive audit rather than installing obsolete or duplicate baseline contracts.
  Merge surviving scenarios instead of replacing unrelated baseline content.
- Keep unfinished changes active. A stale unchecked source-maintenance task
  may be reconciled only when current source and later retained evidence prove
  it complete; record that decision explicitly.
- Record deferred findings in the narrowest owner workspace with visible
  deferred status and unchecked future tasks. Do not sync their proposed
  requirements into current baselines.
- Reuse the existing OpenSpec structural checks and extend their owner/archive
  membership where necessary. No new validator or dependency is introduced.

## Risks / Trade-offs

- Historical deployment evidence may be incomplete → preserve that limitation
  and leave acceptance work active rather than contact live systems.
- Archival changes relative link depth → inspect and repair maintained links,
  keeping provenance and external references intact.
- Superseded requirements may span owners → retain one current contract at the
  owning project and link the audit to it.

## Execution State

The isolated task branch is `t3/review-repo-quality`. Inventory and transient
CLI snapshots live under ignored `out/reconcile-openspec-work/`. The completed
review archives 89 of the original 97 records, including three stale unchecked
tasks reconciled from source or publication evidence. Eight original records
remain unfinished; partial hosting/website progress is recorded without losing
its outstanding acceptance criteria. Six deferred plans retain 24 unchecked
implementation tasks.

The base advanced to `e3d79e4819692be608d2d3d45189b219a7c3d291` during delivery.
Its completed Anki filesystem change was reviewed and archived as the 90th
record, preserving its six historical artifacts. Current publication guarantees
now respect injected filesystem rename semantics. Independent reviews found no
remaining defects after correcting documentation target wiring and the affected
specification wording. See `archive-review.md` for decisions and evidence limits.

Preflight OpenSpec, quality, semantic lint and documentation checks passed.
The delivery receipt owns final exact-candidate gates and publication state;
this maintenance record is archived with its accepted source reconciliation.
