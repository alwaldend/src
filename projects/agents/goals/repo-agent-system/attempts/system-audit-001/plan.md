# System audit plan

## Bindings

- Goal: `repo-agent-system`
- Goal generation: `1`
- Lifecycle generation: `4`
- Criteria revision: `2`
- Checkpoint resource version: `4`
- Work type: `investigation`
- Intended attempt: `system-audit-001`

The goal tool binds the portable goal-state and criteria digests when it
publishes the attempt.

## Target uncertainty

The repository has capable individual mechanisms, but there is no verified
repository-wide account of how a zero-context agent discovers intent, maps
ownership, selects capabilities, executes safely, validates results, delivers
changes, and turns evidence into durable improvement. It is also unclear which
documents are authoritative and where the main duplication and navigation
costs lie.

## Work

1. Inventory root entry points and all repository-owned agent, goal, build,
   quality, documentation, delivery, metadata, and runtime-control surfaces.
2. Delegate disjoint read-only audits with immutable bindings and isolated
   reports under `out/repo_agent_system/audits/`.
3. Trace representative request-to-delivery and failure-to-learning paths.
4. Compare declared contracts with discoverability and actual wiring.
5. Synthesize one evidence report that ranks gaps by consequence, frequency,
   and avoidable agent/resource cost.
6. Adversarially test the synthesis against the strongest case that the
   existing distributed documentation is already sufficient.

## Review packet

- Inputs: tracked repository state on the current feature branch, root and
  nearest ownership guides, Bazel target metadata, canonical skills, durable
  goal records, and repository-facing documentation.
- Candidate identity: SHA-256-bound audit result and selected evidence copied
  into this attempt by the goal tool.
- Affected criterion: `current-state-audit` revision 1.
- Regression checks: validate the goal record and confirm no worker wrote to
  canonical goal state or maintained source.
- Acceptance: the synthesis covers every named surface, cites direct
  repository evidence, distinguishes fact from recommendation, and identifies
  the highest-leverage architectural defect.
