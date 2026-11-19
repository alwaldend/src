# Repository agent-system design result

## Outcome

The design and documentation phase is achieved. The repository now presents
one coherent agent system without moving owner-local facts into a second
mutable registry.

The implemented documentation surface has four roles:

- root `README.md` and `AGENTS.md` are one-hop routers and a compact operating
  kernel;
- `projects/agents/docs/current-state.md` is the revision-bound observed
  baseline;
- `projects/agents/docs/architecture.md` is the normative target system; and
- `projects/agents/docs/roadmap.md` orders future implementation separately
  from present guarantees.

The durable project goal is owned by `projects/agents/goals/repo-agent-system`.
The architecture reserves `tools/agents` for a future repository-internal
controller or executor; no such implementation is claimed in this phase.

## System decision

Each fact has one natural mutation authority. Repository-wide cohesion comes
from thin typed contracts, stable references, and deterministic bounded
projections that retain provenance, freshness, conflicts, unavailable state,
and truncation. A derived view may inform or enforce a transition, but it does
not become a second owner.

The abstraction tower links:

1. intent and authority;
2. topology, classification, and ownership;
3. policy and action admission;
4. capability and provider selection;
5. durable work and planning;
6. execution;
7. observation, evidence, and acceptance;
8. delivery, review, version, release, and deployment; and
9. reviewed learning and accretion.

The temporal work loop is `Orient -> Bind -> Plan -> Admit -> Act -> Prove ->
Decide`, followed by replan, delivery and verification, wait or escalation, or
safe stop. Reviewed learning closes the loop through a minimized regression
and one owner-approved maintained change.

## Public-information and scratch correction

There is no generic “sensitive operational detail” confidentiality class.
Credentials, other secrets, and personal information are prohibited;
non-secret, non-personal operational, generated, live, state, and log facts
are public and may be reported. Raw artifacts require inspection only because
they can embed prohibited content.

This supersedes the active plan's overbroad statement that secret-bearing
material may never enter logs. An unavoidable secret-bearing raw artifact may
temporarily exist under ignored, task-private `out/<task>/` storage with
restrictive permissions, minimum lifetime, and cleanup. It must never be
tracked, staged, committed, or promoted as durable evidence.

The closed audit attempt is retained unchanged as historical evidence. Its
old `operational_sensitive` proposal and scratch-context relative links are
not current policy or navigation. The maintained current-state document is
the link-safe synthesis, and raw audit evidence is intentionally excluded from
the public documentation projection while remaining canonical goal evidence.

## Implemented repository changes

- Added the current-state, architecture, and roadmap documents and Bazel docs
  packaging.
- Added the durable goal landing and canonical project goal record.
- Reworked root and agent landing pages into one-hop system entry points.
- Normalized all six top-level tree guides across source disclosure, Bazel
  visibility, build-consumer, and publication axes.
- Corrected stale infrastructure routes and made state-changing examples
  explicit.
- Corrected question-skill routing for inert quoted payloads.
- Corrected secret handling to the content-based public-repository boundary.
- Added non-obvious tool and delegation narration to the root agent policy.
- Enabled Hugo's embedded Markdown link render hook so source-compatible
  repository links render as site pages.
- Excluded only tool-generated, digest-bound child goal records from generic
  formatting; the maintained goals landing and BUILD file remain checked.

## Acceptance

All six current criteria pass:

- `current-state-audit`: revision-bound observed state and retained audit
  evidence identify entry points, authorities, execution surfaces, feedback
  loops, duplicated facts, and material seams.
- `abstraction-tower`: the architecture defines every layer, owner, typed
  transition, invariant, provenance edge, failure state, and budget boundary.
- `prioritized-roadmap`: Phases 0 through 5 are dependency ordered and include
  measurable acceptance, cost controls, migrations, and non-goals.
- `zero-context-discovery`: a clean root reaches current state, target
  architecture, roadmap, operating policy, and durable goals in one hop.
- `closed-learning-loop`: exact context, plans, authority, receipts,
  validation, semantic assertions, delivery, review, release, outcomes, and
  learning remain linked without centralizing ownership.
- `repository-validation`: source links, rendered links, labels, formatting,
  goal integrity, focused tests, package builds, skill validation, and the
  full documentation site pass for the frozen source candidate.

Future roadmap phases are intentionally not implemented or presented as
current guarantees. The rejected alternatives remain a central mutable
manifest, mega-skill, universal orchestrator, mandatory daemon or vector
store, transcript ingestion, and physical repository reorganization for its
own sake.

## Exact subject

The accepted non-goal source candidate is the 20-file SHA-256 manifest with
digest
`d6e29c5644ca9f834f3a70d2856f3069673d88a19f59ee38dee599a06c40c8ca`.
The goal tool owns the canonical checkpoint files and validates them after
closing this attempt.
