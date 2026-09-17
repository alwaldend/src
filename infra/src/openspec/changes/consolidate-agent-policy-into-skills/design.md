## Context

The root `AGENTS.md` had grown to 255 lines and 2,102 words while much of its
content restated procedures already owned by repository skills. The
duplication had already drifted: `git-rebase-remote` required a replaced commit
to "remain reachable", while `repo-delivery` permitted rewritten ancestry when
progress survived.

## Goals / Non-Goals

- Goal: keep repository-wide constraints visible and route execution detail to
  one owning skill per procedure.
- Goal: make the root document cheap to read without weakening any authority or
  safety boundary.
- Non-goal: change runtime behavior, deployment, or any external service.
- Non-goal: rewrite skill procedures that are already correct.

## Decisions

### Route from one table, ordered by phase

Triggers had accumulated under authority, procedures, verification, and
delivery. One "When to load a skill" table owns routing; each row names the
trigger and the skill. Skills keep their own procedures, so a routing change
does not duplicate content.

### Make initial reads task-dependent

The prior text required reading `BUILD.bazel` and `include.MODULE.bazel` even
for prose reviews. The revision starts from applicable policy and the owning
README, and inspects build declarations only when implementation or validation
depends on them.

### One authoritative source, not literally one copy

The absolute "one copy of each fact" rule conflicted with legitimate generated
projections, fixtures, and summaries. The rule now forbids independently
maintained duplicates and treats source-identifying projections as
projections.

### A `repo-workspace` skill triggers before the first write

Worktree isolation and scratch placement were stated in root and partly in
`repo-delivery`. `repo-delivery` loads too late to establish isolation before a
write, so the procedure moves to a new `repo-workspace` skill whose trigger is
"before the first mutation or task-scratch write". Root keeps the prohibition
on losing unrelated or shared work.

### Reconcile the rebase drift

`git-rebase-remote` now says a replacement must preserve every previously
remote task-owned commit's reachable progress, and that an amend or rebase may
change ancestry. This matches `repo-delivery` and removes the ambiguity.

### Transfer facts whose owner lacked them

- Go and Bazel-native automation rules and the `genrule` prohibition's
  alternatives move into `repo-bazel`.
- Commit-subject and trailer conventions move into a new
  `repo-delivery/references/commits.md`, which the delivery skill links.
- Cache-specific no-op diagnosis already lived in `repo-bazel`.

### Rejected alternatives

- Keeping all procedures in root with tighter wording: preserves the drift risk
  and the read cost.
- Putting workspace isolation only in `repo-delivery`: loads after the first
  write, which is the defect the review identified.
- Creating a by-type "skills" aggregator: `project-layout` rejects aggregating
  unrelated content by type.

## Risks / Trade-offs

- Splitting routing from procedure can hide a constraint if a skill is not
  loaded. Mitigation: every transferred constraint has a named owner, and root
  keeps the repository-wide boundaries listed in the proposal.
- A new skill adds discovery and eval maintenance. Mitigation: it follows the
  existing `bazel-rules-skill` packaging and offline eval contract.

## Migration Plan

No runtime migration. The change is documentation and skill content; delivery
follows the normal repository procedure.

## Open Questions

None.
