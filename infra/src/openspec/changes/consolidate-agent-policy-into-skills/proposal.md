## Why

The root `AGENTS.md` grew to 255 lines and 2,102 words, and much of its
procedural detail already has a skill owner. Duplicated procedures drift, as
the reachability wording in `git-rebase-remote` and `repo-delivery` had
already done, and the length raises the read cost for every task. The file
should keep repository-wide constraints visible and route execution detail to
the skill that owns it.

## What Changes

- **BREAKING** Reduce the root `AGENTS.md` to repository-wide constraints and
  one compact "when to load a skill" routing table; move execution detail to
  its owning skill.
- Make initial reads task-dependent: start from applicable policy and the
  owning README, and inspect BUILD and MODULE files only when implementation or
  validation depends on them.
- Replace the absolute "one copy of each fact" rule with one authoritative
  source, distinguishing independently maintained duplicates from generated
  projections, fixtures, and source-identifying summaries.
- Simplify the retry rule to its principle; move cache-specific diagnosis to
  `repo-bazel`.
- Fix the generated-file ambiguity: exclude disposable build outputs, include
  required generator-maintained files updated through their owning workflow.
- Add a `repo-workspace` skill that owns worktree isolation and task-scratch
  placement, triggered before the first write, so that procedure is not loaded
  too late by `repo-delivery`.
- Transfer facts whose owner lacked them: Go and Bazel-native automation rules
  into `repo-bazel`, commit-message and trailer conventions into a
  `repo-delivery` reference, and the ownership-verdict and infrastructure
  approval boundaries explicitly into root.
- Reconcile existing drift so `git-rebase-remote` and `repo-delivery` agree
  that a task-owned rewrite may change ancestry while preserving every
  previously remote commit's reachable progress.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `repository`: The agent-policy requirement changes from a monolithic root
  document to repository-wide constraints plus one routing table with
  skill-owned procedures; the isolated-implementation requirement gains a
  named workspace skill and an explicit generated-versus-disposable delivery
  boundary.

## Impact

`AGENTS.md` shrinks to roughly 100-140 lines. New skill
`projects/agents/skills/repo-workspace` with its package, discovery link, and
offline eval; new reference
`tools/repo_delivery/skills/repo-delivery/references/commits.md`.
`repo-bazel`, `git-rebase-remote`, `project-layout`, and the delivery skill see
small reconciliations. No runtime capability, deployment, or external service
changes.
