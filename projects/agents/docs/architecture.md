---
title: Repository agent-system architecture
description: >-
  Canonical composition contract for repository-wide agent work,
  evidence, delivery, and learning
tags:
  - agent
  - architecture
  - repository
---

# Repository agent-system architecture

## Status and purpose

This document defines the target composition contract for the repository as
one agent-driven system. It is normative about boundaries, invariants, and
ownership. A named interface or record is not a current capability merely
because it appears here. The [current-state snapshot](current-state.md)
describes implemented behavior; the [roadmap](roadmap.md) orders the work
needed to reach this design.

The objective is to let an agent form an accurate situation model, choose the
least costly safe action, prove the exact result, and make later work cheaper.
The design must remain useful with no daemon, network, credential, runtime
extension, or active durable goal.

## System thesis

The repository should have one composition protocol, not one mutable brain.

Each fact stays at its narrow natural owner. Typed, versioned references join
facts across owners. Deterministic bounded projections make those joins cheap
to read. An execution gateway enforces effects at the point where they occur.
No catalog, agent, skill, plan, or receipt may silently become an independent
authority.

This yields four structural rules:

1. **One mutation authority per fact.** A fact changes only at its owner.
2. **Many derived views.** Human maps, JSON catalogs, context capsules, and
   status views are replaceable projections with provenance.
3. **Typed transitions.** Cross-layer claims carry stable identities,
   versions or digests, completeness, and limitations.
4. **Effects are admitted at execution.** Selection and planning advise;
   explicit authority plus an executable contract controls mutation.

Physical co-location is not required for logical cohesion. Repository-agent
documents and goals belong under `projects/agents`. A future
repository-internal controller or executor may belong under `tools/agents`.
Component declarations remain with their component.

## Mutation authorities

| Fact                                    | Mutation authority                           | System projection                |
| --------------------------------------- | -------------------------------------------- | -------------------------------- |
| Requested outcome and granted authority | Current user/task interaction                | Task authority envelope          |
| Applicable agent policy                 | Nearest `AGENTS.md` chain                    | Policy slice with source digests |
| Component purpose and local boundary    | Nearest owner `README.md`                    | Topology record                  |
| Review accountability                   | `CODEOWNERS`                                 | Effective reviewer record        |
| Workspace and dependency structure      | `MODULE.bazel`, `BUILD`, and owning macros   | Workspace and action catalogs    |
| Reusable procedure and routing intent   | Canonical project-owned skill                | Capability record                |
| Durable outcome, criteria, and attempts | Owner-local goal resources                   | Goal catalog and resume packet   |
| Task-local scratch and worker ownership | Task manifest below ignored `out/<task>/`    | Bounded task-status record       |
| Desired runtime configuration           | Owning checked-in or task-local config       | Desired-state reference          |
| Observed runtime state                  | Live provider instance                       | Provider-health observation      |
| Source candidate                        | Git plus an explicit dirty-input declaration | Subject reference                |
| Check execution                         | Executing provider                           | Action receipt                   |
| Acceptance judgment                     | Goal criterion or task acceptance policy     | Evidence assertion               |
| Feature publication and review          | Git, forge, and repository delivery          | Delivery and review receipts     |
| Version and release identity            | Versioning and release owners                | Release identity and manifest    |
| Durable adopted learning                | Destination owner through normal review      | Updated contract and regression  |

The architecture owns only this composition model. It must not restate mutable
component values. If two authorities conflict, a projection reports the
conflict and its sources; it does not choose a winner. A relevant unresolved
conflict blocks effectful admission.

## Common artifact envelope

Cross-layer artifacts should share a small logical envelope. The roadmap owns
the eventual schema and migration; payloads remain subsystem-specific.

```text
ArtifactEnvelope
  apiVersion, kind, id
  producerRef
  authorityRefs[]?
  inputRefs[]: {kind, id, version?, digest, role}
  subjectRef?
  completeness: complete | partial | truncated | unknown
  limitations[]
  informationClass
  retentionClass
  digest
```

Observation or expiry times belong on live observations. Deterministic
catalog content must not contain nondeterministic generation timestamps.
Artifacts never contain credentials, secret-value hashes, personal
information, environment dumps, absolute checkout paths, unreviewed raw logs,
or full transcripts.

An authority reference is optional and repeatable. A deterministic catalog or
ordinary observation does not gain authority merely by using this envelope.
The artifact digest covers canonical schema-defined bytes without recursively
including its own digest field.

A `SubjectRef` identifies the thing a claim is about. Depending on the claim,
it can bind repository identity, base/commit/tree OIDs, the explicit
task-owned dirty input set, relevant configuration and toolchain digests, an
environment or remote resource selector, and a timed remote observation.
Claims must not use a weaker subject merely because it is easier to obtain.

## The abstraction tower

Each layer consumes contracts from the layer above, adds one kind of fact, and
produces a bounded contract for the layer below.

### 0. Intent, outcome, and authority

**Owner:** the user interaction; durable acceptance belongs to the goal when
one exists.

**Input:** the request, explicit constraints, and prior accepted context.

**Output:** an outcome, acceptance questions, scope, authority envelope, and
resource limits.

**Invariant:** a question, skill selection, plan, prior permission, or tool
availability never widens authority. Conditional authority is usable only
after its condition is evidenced. Missing authority is `unknown`, not
permission.

Repository policy may define an explicit implementation request to include
ordinary feature-branch commit, publication, and review through the delivery
workflow. That default must be visible in the authority envelope and never
extends to tags, releases, deployments, infrastructure changes, or
destruction.

Simple bounded tasks may keep this state in the active interaction. Iterative,
delegated, or resumable work uses a durable goal and attempt.

### 1. Topology, classification, and ownership

**Owners:** paths, nearest READMEs, `CODEOWNERS`, `MODULE.bazel`, BUILD files,
and explicit boundary declarations.

**Input:** workspace identity plus affected paths or labels.

**Output:** component and workspace identity, applicable policy sources,
review accountability, dependency role, and classification axes.

**Invariant:** keep these axes independent:

- checked-in source disclosure;
- log and evidence handling;
- Bazel target visibility;
- allowed build consumers;
- artifact and documentation publication; and
- secret or personal-information presence and live-environment association.

Checked-in content in this public repository is public information unless it
contains an accidental secret or personal information. Operational facts do
not become confidential merely because they are generated or live, but raw
artifacts require inspection because they can contain prohibited content.
Public information status does not make a target eligible for production
consumption or artifact publication.

### 2. Policy and action admission

**Owners:** applicable `AGENTS.md`, owner-local policy, and the executable
provider or owning build macro.

**Input:** authority envelope, topology slice, and requested operation.

**Output:** an admission decision with stable reason codes, required
preconditions, and an exact allowed effect set.

**Invariant:** every runnable operation declares atomic effects such as:

```text
source.read       source.write       task_state.write
host.write        history.write      code.execute
credential.consume
network.read      remote.write       remote.destroy
```

The executable/provider owns effects; a skill references them rather than
copying them. Unknown effects fail closed. Known read-only discovery may
degrade visibly. Remote writes and destruction require exact resource,
environment, candidate, pre-state, authority, and expiry bindings.

`history.write` is distinct from `source.write`: commits, local ref moves,
tags, and rewrites alter history even when worktree bytes do not change.
Operation scopes such as edit, commit, feature publication, review mutation,
release-ref publication, deployment, and destruction map to additive atomic
effects rather than substituting for them.

Review accountability, component ownership, task path/hunk or commit
ownership, goal-storage ownership, remote-ref ownership, external-resource
ownership, and user authority are separate relations. Unknown ownership fails
closed for history rewrite, credential expansion, publication, remote
mutation, and destruction.

### 3. Capability and provider selection

**Owners:** canonical skills for routing and procedure; executable and runtime
providers for actual capability contracts.

**Input:** intent, policy slice, available providers, effects, evidence
maturity, and cost.

**Output:** the smallest sufficient capability set, ordering, exclusions,
fallbacks, and an explanation of every choice.

**Invariant:** selection is advisory and never grants authority. Required
dependencies are acyclic. Conflicts, unavailable providers, unsupported
platforms, and unverified behavior are explicit. Full instructions and schemas
load only after compact metadata selects them.

### 4. Durable work and planning

**Owner:** the current task coordinator and, for durable work, owner-local goal
resources.

**Input:** intent, criteria, context slice, capability choices, prior evidence,
and relevant defects.

**Output:** a `WorkPacket` or goal attempt plus `CapabilityPlan` and
`ActionPlan`, with exact inputs, affected criteria, hypotheses where needed,
checks, omissions, fallbacks, fixed regressions, budgets, stop conditions, and
a strategy-reset condition.

**Invariant:** a plan does not authorize and does not prove. Delegated workers
receive immutable bindings and disjoint output ownership; only the coordinator
mutates canonical goal state. A repeated stable defect changes strategy rather
than accumulating indistinguishable retries.

### 5. Execution providers

**Owners:** Bazel and `bazel_agent`, native repository tools, Cordis or another
runtime provider, and explicit external adapters.

**Input:** admitted operation, exact subject and inputs, task/run identity,
deadline, quotas, and cancellation contract.

**Output:** immutable `ActionReceipt` records plus bounded task-local raw-log
pointers.

**Invariant:** executors remain thin and signal-transparent. Scratch is
isolated by workspace, task, run, and worker; generated execution receives
only declared filesystem, environment, process, network, and credential
capabilities. Desired and observed state remain separate. A control kernel
must stay responsive when optional data-plane extensions fail.

### 6. Observation, evidence, and acceptance

**Owners:** providers own observations; the task or goal acceptance policy owns
semantic sufficiency.

**Input:** receipts, subject, validation profile, declared policy, known
limitations, and—only for semantic judgment—a criterion revision.

**Output:** three distinct records:

- `EvidenceEvaluation` derives `applicable`, `stale`, `insufficient`,
  `unavailable`, or `unverifiable` with stable reasons;
- execution-owned, criterion-neutral `ValidationSet` aggregates exact candidate,
  plan/check identities, sanitized invocations, provider/configuration
  provenance, results, coverage, and limitations; and
- goal/task-owned `EvidenceAssertion` applies evidence to an exact criterion
  and subject with `satisfied`, `not_satisfied`, or `unknown`.

**Invariant:** receipts prove what ran and what was observed; they do not prove
that selected checks were semantically sufficient. One immutable validation
set may support non-goal delivery and multiple criterion assertions. Partial,
truncated, skipped, or unknown evidence cannot support complete coverage.
Reuse requires matching relevant input, policy, contract, configuration,
toolchain, and subject-class digests.

### 7. Delivery, review, version, and release

**Owners:** Git and the forge own remote facts; repository delivery owns
feature publication and review operations; versioning and release keep their
separate identities.

**Input:** exact candidate, authority, an applicable criterion-neutral
`ValidationSet`, expected remote state, and operation-specific policy.

**Output:** delivery, synchronization, review, version, release, and deployment
receipts linked by exact subject references.

**Invariant:** use the effectful protocol:

```text
inspect -> prepare -> immutable candidate -> validate -> authorize
        -> execute -> verify -> receipt
```

Delivery validates candidate-bound evidence but does not select semantic
checks. Feature publication does not imply tag or deployment authority.
Version identity and release snapshot identity are joined by a typed handoff,
not by assuming that any Git tag has the same meaning.

### 8. Reviewed learning and accretion

**Owner:** the component that will adopt the lesson through ordinary review
and delivery.

**Input:** exact task evidence, review dispositions, stable defect identity,
measured friction, and regression results.

**Output:** a `LearningProposal`, then—only after acceptance—an owner-local
document, contract, check, evaluation, skill, or global invariant.

**Invariant:** learning is never automatic mutation. It is minimized, public,
reviewed, regression-backed, measurable, reversible, and given a retirement
path. Runtime promotion creates a normal candidate; it does not write project
source or goal state directly.

## Provenance spine

The complete traversable chain is:

```text
ContextCapsule
  -> GoalAttempt / WorkPacket
  -> CapabilityPlan / ActionPlan
  -> AuthorityBinding + AdmissionDecision
  -> ActionReceipt(s)
  -> EvidenceEvaluation + ValidationSet + EvidenceAssertion(s)
  -> accepted candidate
  -> Delivery / Review / ReleaseRef / ReleaseManifest / Deployment receipts
  -> GoalOutcome
  -> LearningProposal / LearningAdoption
```

The logical records serve distinct purposes:

- `ContextCapsule` says what was known, applicable, authorized, available,
  stale, and omitted before planning.
- `WorkPacket` says what outcome and criteria a bounded unit will pursue.
- `CapabilityPlan` and `ActionPlan` say what was selected, omitted, ordered,
  budgeted, and expected to prove the result.
- `AuthorityBinding` and `AdmissionDecision` say which exact effects, subject,
  environment, policy, pre-state, and expiry were checked before execution.
- `ActionReceipt` says what provider executed or observed against which exact
  inputs and with what result and resource use.
- `ValidationSet` records criterion-neutral check execution and coverage;
  `EvidenceAssertion` makes the separate criterion judgment.
- Delivery, review, release-ref, manifest, and deployment receipts say which
  external transition occurred and which postcondition was observed.
- `GoalOutcome` records durable acceptance or the honest terminal state.
- `LearningProposal` requests owner review; `LearningAdoption` identifies the
  exact accepted change, regression, delivery, and measured effect.

Plans do not authorize. Receipts do not judge semantic sufficiency. Delivery
receipts do not prove that validation selection was complete. Learning
proposals do not edit policy.

Every edge binds the relevant subject and catalog-slice, policy, contract, and
input digests.

## Temporal work loop

The tower defines composition; this loop defines control over time:

```text
Orient -> Bind -> Plan -> Admit -> Act -> Prove -> Decide
                                              |-> replan
                                              |-> deliver -> verify -> close
                                              |-> wait or escalate
                                              `-> stop
                                                     |
                                                   Learn
```

Orientation obtains the smallest context capable of exposing uncertainty.
Binding freezes the intended outcome, exact subject, authority, criteria, and
budgets. Planning names omissions, fallback, stop conditions, and the cheapest
evidence capable of falsifying the current hypothesis. Admission checks the
exact plan against current policy and pre-state. Decision either accepts,
changes strategy, waits for named external state, escalates authority, or
stops honestly.

A declared transient failure permits at most one same-input retry. A
deterministic recurrence requires a changed hypothesis or strategy. A remote
operation with unknown postcondition is observed and reconciled before any
second write.

## Bounded context capsule

The target zero-context interface is one read-only, bounded projection. It may
be a repository tool and optionally an MCP surface, but it cannot require MCP
to function. It contains:

- repository, workspace, worktree, revision, and dirty-input identity;
- task/session, coordinator, worker, and run identity;
- requested outcome, authority, budgets, and durable goal binding if any;
- applicable instruction and owner-document paths with digests;
- component, workspace, review owner, and classification slice;
- candidate capabilities with effects, cost, dependencies, providers, and
  evidence maturity;
- provider boot/incarnation ID, catalog ETag, desired/observed revisions,
  action-contract hashes, observation time, expiry, and unavailable reasons;
- relevant checks and reusable evidence, not raw logs; and
- provenance, observation time, freshness, completeness, truncation, and safe
  next discovery actions.

The capsule is a join over owner-local catalogs and observations. It stores no
independent truth. A missing optional provider yields a structured unavailable
field; it does not make the whole capsule fail.

## Resource economy

An agent climbs this cost ladder only when existing evidence cannot answer the
acceptance question:

1. checked, cached metadata and narrow source reads;
2. targeted static queries and reusable subject-bound receipts;
3. narrow local builds and tests;
4. broad or uncached validation;
5. external reads and stochastic/model-assisted evaluation; and
6. authorized remote mutation or destructive work.

Budgets inherit from repository defaults to task to plan to operation/worker.
A child may narrow but never silently expand a hard bound. Reserve capacity for
cancellation cleanup and final verification. Track context/input bytes,
compute and target/workspace scope, elapsed time, output/log/artifact bytes,
retention, retries, worker fan-out, concurrency, and cacheability. Use measured
p50/p95 workloads to set numeric defaults.

Credentials, network access, history or remote writes, and destruction are
non-fungible authority scopes, never spendable budget. Compatible Bazel work
is batched. Context loads summaries and digests first, then full bodies,
schemas, source, or logs only on demand. A global catalog change triggers
cheap applicability reevaluation, not automatic reruns of unrelated evidence.

## Accretion protocol

```text
sanitized bounded friction event
  -> aggregation by stable signature and consequence
  -> LearningProposal with destination owner
  -> minimized regression or routing fixture
  -> destination-owner review
  -> one appropriate contract, policy, tool, test, or skill change
  -> normal validation and delivery
  -> LearningAdoption and regenerated evidence status
  -> recurrence and resource measurement
```

Raw transcripts, credentials, unreviewed runtime values, secret-bearing logs,
personal information, and accidental local state never enter the promotion
path. Ignored task-private `out/<task>/` may temporarily contain secret-bearing
material under restrictive access and cleanup, but that material is not
tracked or imported as durable evidence. A global invariant is an exceptional
owner-approved destination, not the automatic final rung of the sequence.

## Failure and degradation semantics

- `unknown` means the system lacks a fact; it is not false, safe, or allowed.
- `unavailable` names a provider or source that could not be observed.
- `stale` means a once-valid observation no longer matches its applicability
  inputs.
- `partial` and `truncated` bound what a result can prove.
- Desired state and observed state are always separate fields.
- Runtime observations identify boot/incarnation, catalog, desired and
  observed revisions, action contract, observation time, and expiry.
- Provider lifecycle is explicit: `loading`, `ready`, `degraded`, `failed`,
  `timed_out`, `draining`, or `disabled`.
- Shared mutable task, goal, and runtime state uses a cross-process lock or
  lease plus expected version/digest; atomic rename alone is not concurrency
  control.
- Nested work uses one absolute deadline with cancellation and cleanup margin.
- Failures carry phase, stable code/signature, subject, observed state, valid
  partial result, retry class, cost spent, redacted artifact references, and
  next safe action. Retry class is `never`, `same-input-once`,
  `after-state-change`, or `requires-human`.
- Read-only discovery returns the useful safe subset with explicit gaps.
- Effectful admission fails closed on unknown effects, authority, environment,
  subject, policy conflicts, or relevant stale state.
- Secret-bearing raw output stays in task-private scratch and only a
  bounded allowlisted summary may cross into ordinary evidence.

## Measures

Optimization is evaluated with representative zero-context tasks, not by file
count or architectural novelty. Track at least:

- task correctness and unsupported-claim rate;
- unauthorized or incorrectly scoped mutation rate;
- reads, tool calls, context bytes or tokens, and elapsed/resource cost before
  the first correct action;
- redundant checks and exact evidence reuse rate;
- steps and missing facts required to resume durable work;
- stale, partial, and conflict detection accuracy;
- provider startup and degraded-status latency; and
- the rate at which repeated defects gain a regression, owner, delivery
  identity, and eventual retirement.

## Non-goals

- No hand-maintained central manifest duplicating owner-local facts.
- No universal mega-skill or always-loaded copy of every procedure.
- No mandatory daemon, database, vector store, or runtime extension for basic
  orientation and safe degradation.
- No raw transcript ingestion, secret promotion, or automatic doctrine edits.
- No inference of authority from selection, planning, receipts, ownership, or
  tool availability.
- No durable goal requirement for a simple one-step task.
- No broad `//...` validation or live model evaluation by default.
- No physical repository reorganization merely to make the conceptual model
  look centralized.
