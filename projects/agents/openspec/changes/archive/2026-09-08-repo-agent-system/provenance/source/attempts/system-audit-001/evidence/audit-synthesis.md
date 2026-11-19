# Repository agent-system audit synthesis

## Binding and scope

- Goal: `repo-agent-system`
- Attempt: `system-audit-001`
- Tracked revision:
  `1423dce5fab45ce5223caeb6a24791bf1a2cc3ff`
- Goal generation: `1`
- Lifecycle generation: `4`
- Criteria revision: `2`
- Criteria digest:
  `sha256:2ac2db1242f5d3358e433b3499da5a622d06bdec49bfa690dd34cf3205e28f34`
- Goal-state digest:
  `sha256:193be5b38881faebc349f9ae1d273e24fac5d5925a9a4402b24706394ffaeb3a`

This synthesis covers repository entry points, topology, ownership, skills,
goals, runtime control, build control, safety, evidence, delivery, release,
and learning. The eight source audits are bounded, read-only reports. Their
recommendations are evidence, not independently added acceptance criteria.

For scratch terminology, `repository` means tracked or committable source.
Short-lived secret-bearing temporary material may live in ignored,
task-private `out/<task>/` storage when access, lifetime, and cleanup are
controlled. It must never be tracked, staged, committed, logged, or promoted
as durable evidence.

## Direct answer

The repository has strong components but no cheap, explicit system interface.
A zero-context agent can reconstruct much of the situation, but only by
reading a large flat policy, knowing which distributed documents and tools to
inspect, joining incompatible records in its own context, and trusting prose
at critical transition points. That reconstruction cost repeats for every
agent and after every context loss.

The highest-leverage improvement is a thin, typed composition contract over
the existing owner-local authorities. It should produce bounded, versioned,
provenance-carrying projections for discovery and evidence. It must not become
a second mutation authority.

## What should be preserved

1. Bazel is already the reproducible executable and dependency graph.
2. Nearest `AGENTS.md` and `README.md` files put policy and explanation near
   their owners.
3. Canonical skills have exact discovery links, validation, and offline eval
   configuration.
4. Goal records have compare-and-swap versions, immutable closed attempts,
   digest-bound criteria, and generated projections.
5. Delivery uses exact refs, OIDs, leases, receipts, and post-mutation
   verification.
6. Cordis provides useful bounded repository-reading and runtime primitives.
7. Root policy already distinguishes information requests from mutation
   authority and prefers narrow verification.

These are suitable leaf authorities and execution providers. Replacing them
with a universal agent service would discard useful integrity boundaries.

## Evidence-backed baseline

- The root README is not a repository map; it contains external links and the
  license but no path into the agent system.
- One 239-line repository-wide `AGENTS.md` is the only checked-in agent guide.
  It combines constitution, topology, procedures, and exceptions in one
  always-loaded surface.
- The repository has 409 tracked READMEs. Metadata is widespread but not
  typed: 407 use frontmatter, 383 have descriptions, 239 have tags, and 17
  have statuses.
- All 28 project roots have a README and BUILD file, but lifecycle metadata is
  incomplete and vocabulary is inconsistent.
- The root documentation dependency closure is not a cheap topology API; a
  cold audit crossed into external dependencies and failed after more than
  100 seconds.
- `CODEOWNERS` uses the literal path pattern `-`, not the catch-all `*`, so it
  does not assign ordinary repository paths to its listed owner.
- Twenty skills are registered. Their packaging is coherent, but registration
  is manual, composition is prose-only, and only three of eighteen canonical
  repository-agent skills have live behavioral targets.
- The goal store is locally rigorous but has no repository-wide catalog,
  bounded continuation packet, or structured joins to validation, delivery,
  review findings, and regressions.
- Build and runtime providers expose useful primitives but no common effect,
  authority, cost, or evidence contract. Some runnable Bazel targets can read
  secrets or mutate remote state while looking like ordinary `bazel run`
  entry points.
- Delivery receipts bind exact candidates, but validation execution is still
  caller-asserted. Versioning, release snapshots, review outcomes, goals, and
  regression targets lack typed cross-system links.
- Workspace lists and other system facts are repeated across several files,
  making normal accretion require coordinated manual edits.

## Ranked systemic gaps

### 1. No one-hop situation model

There is no bounded answer to: where am I, what does the user authorize, which
policies apply, who owns the path, what goal is active, which capabilities are
available, what effects and costs do they have, and which evidence is current.
Agents repeatedly spend reads, queries, context, and inference rebuilding it.

### 2. Safety and authority are prose at execution seams

Different surfaces conflate checked-in disclosure, Bazel visibility,
production dependency eligibility, artifact or documentation publication,
and operational sensitivity. Runnable operations have no mandatory effect
classification. Unknown or ambiguous actions therefore fail open to agent
judgment rather than closed at the execution gateway.

### 3. The work-to-proof chain is not traversable

Intent, attempt, command execution, validation, delivery, release, review,
and regression data exist in separate schemas or Markdown. Exact subject
identity and criteria do not travel through the entire chain, so evidence is
difficult to reuse and continuation depends on prose and memory.

### 4. Discovery is distributed but not projected

Owner-local facts are the right mutation authorities, yet there is no checked
projection joining paths, workspaces, ownership, lifecycle, skills, targets,
effects, runtime providers, and goals. Existing documentation aggregation is
a packaging graph, not a semantic catalog.

### 5. Accretion lacks a promotion protocol

Real task friction can be saved in goal evidence, but the repository does not
define when an observation becomes a stable defect, contract, test, routing
case, skill, or global invariant. Repeated failures can remain prose, while
speculative lessons can be promoted without regression evidence.

### 6. Cost is not a first-class planning input

The repository advises narrow checks, but providers do not expose comparable
cost, cacheability, freshness, or evidence maturity. Agents cannot reliably
choose the cheapest action that can discharge an acceptance criterion or
reuse an exact prior receipt.

## Required abstraction tower

The maintained design should define these layers without moving their facts
into one database:

1. intent, outcome, and authority;
2. topology, ownership, and information classification;
3. policy and action admission;
4. capability and provider selection;
5. durable goals and work packets;
6. execution providers;
7. evidence and acceptance;
8. delivery, release, and publication; and
9. reviewed learning and regression-backed accretion.

Each layer needs one mutation authority, typed inputs and outputs, explicit
invariants, and named consumers. Derived views must expose their sources,
versions or digests, observation time, freshness, unavailable state, and any
truncation.

The minimum provenance spine is:

```text
ContextCapsule
  -> GoalAttempt / WorkPacket
  -> ActionReceipt
  -> EvidenceAssertion
  -> DeliveryReceipt
  -> LearningProposal
```

Desired state and observed state remain separate. Authority is never inferred
or expanded by a projection. Unknown effects fail closed; read-only discovery
may degrade visibly. Evidence binds the exact subject and criterion revision.

## Resource-economy rule

Agents should climb a cost ladder only when the current evidence cannot answer
the acceptance question:

1. checked, cached metadata and narrow source reads;
2. targeted static queries and cached subject-bound receipts;
3. narrow local builds and tests;
4. broad or uncached validation;
5. external reads and model-assisted evaluation;
6. authorized remote mutation or destructive work.

The context projection must be bounded and progressively disclosed. A compact
catalog should route to full owner-local documents; it should not eagerly load
all instructions, logs, or dependency closures.

## Accretion rule

```text
task-local observation
  -> exact attempt evidence
  -> recurring stable defect
  -> owner-local document or typed contract
  -> executable check or regression
  -> routing case or reusable skill
  -> global invariant only when genuinely cross-cutting
```

Raw transcripts, credentials, private runtime values, and accidental local
state never enter this ladder. Promotion is reviewed and must preserve a
minimal public reproducer plus a regression or validation contract.

## Decision review

Three credible designs were tested.

1. **Keep distributed documentation only.** This preserves locality and has
   no new machinery, but the audit falsified its sufficiency: entry points,
   cross-layer joins, effect classification, and the learning loop are absent.
2. **Build a central agent hub.** This offers convenient lookup, but creates a
   duplicate mutable authority, synchronization and privacy burden, and a
   daemon/database dependency before demand is established.
3. **Keep distributed authorities and derive a thin typed projection.** This
   preserves current integrity boundaries while closing the discovery and
   provenance gaps. Its main risk is a projection becoming stale or treated
   as authoritative; mandatory provenance, regeneration, drift tests, and
   fail-visible freshness address that risk.

Verdict: **revise and proceed with option 3**. The implementation plan must
explicitly reject a mega-skill, hand-maintained central manifest, automatic
transcript ingestion, and a mandatory daemon or vector database. Physical
co-location under a new directory is not required for logical cohesion.

## Audit evidence

| Report | SHA-256 |
| --- | --- |
| `build_control.md` | `587b866ab22838def523c7ed2e5df966dbb71f5cbe5f91dc73b81d6bcf01b8fc` |
| `delivery.md` | `718f816f686c7e512076166d616590bc031ece58742470ff97d611b8b6cc38bc` |
| `entrypoint.md` | `6d608485e1dc75f40232fb32bf617327cb2510179abdd5d526790f66c5b2421b` |
| `goals.md` | `c50538e8341027142a804924c37a73cffaaecdf87ac0028926197ed2f5139f51` |
| `runtime.md` | `42643b02b875146d8fd7c7580f5a1dfc8cd4c73250902ec263c491f9cc497409` |
| `safety.md` | `a963222be2325cc8536c83b4143658f4c3f6dad865ef60c0ba46a3c82e9a27e8` |
| `skills.md` | `5c01950cd0562306fc07c3a1b2220cc32ffce1c3c54892113957d3fc28641b57` |
| `topology.md` | `f15f07927226947daeaef1a0d122d575c1311bb16626689efa5805f0dfcaa5dc` |

## Criterion conclusion

`current-state-audit` passes for this tracked revision. It is an audit
snapshot, not a claim that the proposed architecture or roadmap has yet been
implemented or validated. The remaining five criteria stay unverified for the
next attempt.
