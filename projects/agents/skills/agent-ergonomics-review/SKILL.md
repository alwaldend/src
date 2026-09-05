---
name: agent-ergonomics-review
description: >-
  Review a completed session from the consuming agent's perspective by
  measuring discovery cost, context size, avoidable reads or commands, routing
  accuracy, and resume quality. Use at task close for substantial work, failed
  or inefficient sessions, and repeated friction; do not use for routine
  single-turn responses.
---

# Review agent ergonomics

Assess the session as if a fresh agent had to reproduce or continue it. Do not
rewrite history or relabel a poor path as success. Separate what the task
required from avoidable friction, and propose only changes whose expected
benefit justifies their cost.

## Measure the observed path

Before closing the session, inspect the actual evidence rather than relying on
memory. Report bounded counts and identities, not transcripts or secret-bearing
output:

- Stable task identity and exact delivered subject.
- Selected and considered skills, catalogs, runtimes, and tools.
- Context, catalog, fixture, and result identities; large reads or transfers.
- Avoidable or repeated reads, commands, checks, and waiting.
- Routing failures, conflicts, missing providers, and wrong-skill selections.
- Invalid assumptions, unsafe-adjacent actions, and their correction.
- Verification latency and whether the cheapest sufficient check was used.
- Whether a fresh agent could discover and resume the work from bounded
  records without free-form archaeology.

Prefer existing goal, validation, or delivery receipts over new storage. If a
workspace learning record is warranted, keep it sanitized, bounded, stable-ID,
and task-local. It is a proposal for review, never an automatic edit to shared
instructions, schemas, catalogs, or runtime behavior.

## Turn friction into reviewable evidence

1. Name one stable defect or opportunity per issue, using a stable public
   identifier rather than prose that changes between sessions.
2. Attach the smallest reproducible evidence, metric, and evidence tier
   (`configured`, `routed`, `fixture-tested`, `live`, `stale`, or
   `unverified`). Declarations and configuration checks do not establish
   observed routing, successful behavior, or runtime health.
3. Distinguish a routing error, missing contract, stale catalog, avoidable
   read, repeated check, and context-size failure from an acceptable cost.
4. Propose the narrowest owner change: fixture, assertion, skill description,
   reference, catalog projection, checkpoint field, or optimization.
5. Define how to measure the improvement and when to retire proposal tracking
   or an obsolete workaround. Clean runs can show that a fix works; they do
   not justify removing useful diagnostics or regression coverage. Remove a
   safeguard only when evidence shows it is obsolete or its cost outweighs
   its benefit, while preserving required guarantees.
6. Route the proposal through goal checkpoint, review, or delivery as
   applicable; do not silently edit canonical state. Fix small, in-scope
   ergonomics problems immediately without asking; obtain explicit user
   authorization when a proposed remedy is large, uncertain, costly, or
   expands task scope.

## Preserve proportion and safety

A concise session needs no formal report. A failed or inefficient substantial
session should produce enough evidence to prevent recurrence without spending
more context than the issue justifies. Never include credentials, personal
information, unreviewed runtime values, or secret-bearing output in records or
proposals.
