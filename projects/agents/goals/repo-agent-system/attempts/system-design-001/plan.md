# Agent-system design and documentation plan

## Bindings

- Goal: `repo-agent-system`
- Goal generation: `1`
- Lifecycle generation: `4`
- Criteria revision: `2`
- Expected checkpoint resource version: `6`
- Intended attempt: `system-design-001`
- Work type: `change`
- Audit attempt: `system-audit-001`

The goal tool will bind the current goal-state and criteria digests when this
attempt is published.

## Target result

Make the repository legible as one synthetic agent system without creating a
second source of truth. A zero-context human or agent should be able to enter
at the root, identify the canonical authorities and abstraction tower, follow
the current work-to-learning loop, distinguish present guarantees from future
work, and select the next source or operation without loading the whole tree.

## Design hypothesis

The best balance of correctness, cohesion, and resource economy is to keep
facts with owner-local authorities and define one thin composition contract
plus bounded derived projections. A monolithic registry, mega-skill, daemon,
or transcript store would duplicate authority and add more synchronization
cost than it removes at this stage.

## Work

1. Turn the audit verdict into a canonical architecture: layers, authorities,
   inputs, outputs, invariants, joins, failure behavior, and provenance.
2. Publish the evidence-backed current-state snapshot separately from the
   normative architecture.
3. Publish a dependency-ordered roadmap with acceptance signals, resource
   controls, migration boundaries, and explicit non-goals.
4. Make the root README and root agent guide route directly to those documents
   while keeping the guide a compact policy kernel rather than duplicating
   mutable system state.
5. Clarify that secret-bearing temporary material may use ignored,
   task-private `out/<task>/` storage but may never enter tracked, staged, or
   committed source, logs, or durable evidence.
6. Wire new documents into the existing Bazel documentation graph.
7. Run link, format, goal, skill, package, and Buildifier checks against the
   exact candidate; inspect formatting of digest-bound generated goal records
   before deciding whether formatter exclusions are required.
8. Perform an adversarial design review and incorporate only evidence-backed,
   in-scope corrections.
9. Close the attempt with exact evidence and use the repository delivery
   workflow to commit, publish, review, and verify the final candidate.

## Acceptance mapping

- `current-state-audit`: retained from `system-audit-001` and represented by a
  maintained current-state document with source links and snapshot identity.
- `abstraction-tower`: canonical architecture covers every layer and join.
- `prioritized-roadmap`: roadmap is dependency ordered, measurable, and
  explicit about current versus proposed behavior.
- `zero-context-discovery`: root README and agent guide form one-hop routes.
- `closed-learning-loop`: request-to-delivery and
  observation-to-regression traces preserve authority and evidence identity.
- `repository-validation`: every changed source and generated record passes
  the narrow applicable checks against the final candidate.

## Fixed regressions and evidence

- No design document may become a parallel mutable authority.
- Future proposals must not be described as current guarantees.
- Ignored task-private secret scratch remains allowed under strict handling;
  tracked or committed secret material remains prohibited.
- Root entry points must link to canonical documents with repository-relative
  paths that resolve.
- The docs graph must include the new documents without requiring a broad
  dependency-closure query for navigation.
- Goal records must remain tool-generated and digest-valid.
- BUILD files must remain Buildifier-clean and affected packages build/test.

## Strategy reset rule

If the architecture requires a new central mutable registry, mandatory daemon,
or physical reorganization to express the joins, stop and revise the design.
If validation shows generated goal records conflict with generic formatting,
add the narrowest documented formatter boundary rather than editing generated
records or weakening repository-wide formatting.
