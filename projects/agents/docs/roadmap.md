---
title: Repository agent-system roadmap
description: >-
  Dependency-ordered plan for making repository work safer, cheaper,
  resumable, and regression-backed
tags:
  - agent
  - architecture
  - roadmap
---

# Repository agent-system roadmap

## Status

This is an intentional future plan, not a list of current guarantees. The
[current-state snapshot](current-state.md) is the baseline and the
[architecture](architecture.md) is the target contract. Phase 0 is the scope
of the current documentation change. Every later phase requires its own goal,
decision review, implementation, tests, and delivery.

The roadmap orders work by dependency and leverage rather than component. A
phase is complete only when its acceptance signals are measured against an
exact candidate. Later phases may prototype behind read-only or advisory
interfaces, but they must not become authoritative early.

## Prioritization rules

1. Correct false or ambiguous authority and safety semantics before adding
   automation that consumes them.
2. Put one typed fact at its natural owner before generating a system view.
3. Generate and check read-only projections before enforcing them.
4. Preserve cheap read/build/test paths while adding effect admission.
5. Bind exact execution evidence before using it to gate delivery.
6. Capture and measure friction before optimizing context or adding storage.
7. Promote learning only through a minimized regression and owner review.

Priority means:

- **P0:** a trust, authority, safety, or system-entry blocker;
- **P1:** a high-frequency correctness, continuation, or cost improvement; and
- **P2:** an optimization or consolidation justified by measurements.

Every qualitative resource gate must become numeric before enforcement. The
first implementation change in each phase records a revision/corpus-bound
baseline fixture and ceiling for its standing scenarios. At minimum measure
correctness, unsafe actions, calls, context bytes or tokens, cold/warm wall
time, target/workspace coverage, and reused checks. A claim names the scenario
digest, observation mode, baseline, ceiling, and final measurement; unavailable
measurements remain unavailable rather than estimated.

## Dependency spine

```text
semantic corrections
  -> shared identity, effect, authority, and evidence contracts
  -> owner-local typed declarations
  -> validated generated catalogs
  -> bounded context capsule
  -> advisory planning and admission
  -> candidate-bound execution evidence
  -> delivery, release, review, and goal joins
  -> regression-backed accretion and optimization
```

Runtime control-kernel and scratch-isolation work may proceed beside the
catalog work after the shared identity and effect vocabulary exists. Runtime
promotion cannot precede isolation and normal delivery.

## Phase 0: establish the system contract

**Status:** in scope for the current change.

### Deliverables

- Publish an evidence-backed current-state snapshot separately from the
  normative architecture.
- Publish the abstraction tower, mutation-authority model, provenance spine,
  degradation rules, cost ladder, and accretion protocol.
- Publish this dependency-ordered roadmap with explicit non-goals.
- Turn the root README and agent guide into one-hop routes to those documents.
- Clarify that source disclosure, publication eligibility, and the presence of
  secrets or personal information are separate axes.
- Clarify ignored task-private handling for temporary secret-bearing material
  and the sole raw-Bazel bootstrap exception.
- Align the affected repository-wide skills with those definitions.
- Store the audit, decision, plan, and results in a durable owner-local goal.
- Wire the documents into the existing Bazel documentation graph.
- Resolve source-compatible Markdown page links in the rendered documentation
  site and pin snapshot evidence that is not a packaged page to its revision.

### Acceptance signals

- A zero-context reader reaches current state, architecture, roadmap, policy,
  and durable evidence from the root in one hop.
- Every document labels current facts, target contracts, and proposed work
  without conflation.
- Repository-relative links resolve in source and render to valid site pages;
  revision-pinned evidence links resolve to public source; docs and affected
  skills package; goal digests validate; BUILD files are Buildifier-clean.
- The exact final candidate passes the narrow repository quality and delivery
  checks selected for these changes.

### Resource guard

This phase adds no daemon, schema implementation, generated catalog, runtime
dependency, live evaluation, or broad repository build. The audit reports are
public-source design evidence; operational details are allowed, and the
reports contain no credentials, other secrets, or personal information.

## Phase 1: normalize semantics and declare safety contracts

**Priority:** P0. **Depends on:** Phase 0 vocabulary and ownership model.

### 1A. Correct concrete trust defects

- Replace the literal `-` CODEOWNERS pattern with an intentional catch-all and
  add narrower ownership only where semantic responsibility differs.
- Audit actual documentation, build-consumer, and publication paths before
  changing tree guarantees. Narrow false blanket claims rather than inventing
  machinery to defend them.
- Remove or replace mutating unnamed/default entry points such as a root
  Terraform target whose default operation is `apply`.
- Resolve provider-specific ownership between repository delivery and remote
  rebase workflows. Distinguish edit, commit, feature publication, tag,
  release, deployment, and destructive authority.
- Replace shared `out/tmp` and worktree-global runtime scratch with explicit
  task/run namespaces. Keep Bazel-managed action/test temporary storage
  separate from host controller scratch.
- Add non-disclosing adapters for credentialed diagnostics; raw output that
  has not been checked for secrets or personal information remains in
  task-private scratch, while ordinary results expose reviewed structural
  summaries and safe operational facts.

### 1B. Define the shared contract vocabulary

Version a narrow, additive type family:

- repository, workspace, task/session, actor, provider, and subject references;
- authority and budget envelopes;
- independent path-policy axes;
- atomic operation effects;
- artifact provenance, completeness, information, and retention classes;
- stable unavailable, stale, partial, truncated, denied, and conflict reason
  codes; and
- evidence applicability states.

The types belong to a future repository-internal control package, likely under
`tools/agents`, while the architecture and durable goals remain under
`projects/agents`. They must not absorb component payload schemas.

### 1C. Add owner-local declarations

- Define an `AgentAction` contract in owning Bazel macros and runtime providers:
  effects, inputs and outputs, information classes, credential/network use,
  environment selector, authority gate, preflight, verification, cost,
  cacheability, cancellation, and owner.
- Add minimal skill metadata: stable ID, owner, layer, activation and
  exclusions, capability references, dependencies and conflicts, provider
  requirements, context cost, and evaluation maturity. Executable effects
  remain referenced from providers, not copied into skills.
- Define project lifecycle vocabulary and require one value at project roots.
- Define a task manifest for `out/<task>/` with task/run/worker ownership,
  classifications, budgets, retention, locks, and cleanup.
- Inventory generated artifacts, their formatter owner, checked output, source
  updater label, and intentionally curated exclusions.
- Add a read-only `bazel_agent doctor` contract for runner/source identity,
  Bazelisk path and pin, platform, rc/profile composition, task scratch, and
  stale host-install detection.

Completeness is always relative to a closed registered universe. The first
schema must enumerate its registration authorities: repository-owned
agent-reachable Bazel operations and supported direct binaries, registered
runtime tools, discoverable skills, maintained components, tracked workspaces,
and registered goals. Every default or unnamed alias resolves to an operation
record. Legacy, ignored, unavailable, observed-but-unregistered, and
`RequiresMigration` entries are reported explicitly; the system never claims
absence outside the named universe.

### Acceptance signals

- Schema fixtures round-trip deterministically and reject unknown effectful
  operations, malformed identities, widened authority, and incompatible
  information flows.
- Every operation in the registered agent-reachable universe can represent its
  present behavior; adding one through any supported definition surface
  without classification fails a cheap focused check that names the universe.
- Public information and secret or personal-information content are classified
  without contradiction.
- Two tasks in one worktree cannot share controller/runtime scratch
  accidentally; stale workers fail compare-and-swap publication.
- Safe explicit replacement labels exist before mutating aliases are removed.

### Resource guard

Start in report-only mode and inventory actual behavior before enforcement.
Do not wrap every Bazel action in a new process or parse full dependency
closures. Static declarations and direct-owner checks belong in normal
hygiene; broad graph audits remain explicit.

## Phase 2: generate a trustworthy system map and context capsule

**Priority:** P0/P1. **Depends on:** Phase 1 types and declarations.

### 2A. Make durable goal state diagnosable and recoverable

Before a goal catalog advertises rich continuation or evidence joins, define
stable committed/incomplete states and supported `doctor/recover` behavior.
Use immutable revisions plus an atomic current pointer, write-ahead design, or
an equivalent protocol. Inject interruption at every publication boundary;
recovery must yield the prior valid record or idempotently finish the new one.

Specify legacy migration with a new valid ID, retained source path/digest/raw
bytes, explicit field mapping, and preserved unmapped prose. No direct YAML
editing or inference from projection text is a recovery mechanism.

Until this gate passes, a generated goal catalog may expose only validated
identity and coarse status; an invalid record is `unavailable`, not resumable.

### 2B. Generate bounded catalogs

Create one versioned schema and compiler per concern, with owner-local facts as
inputs:

- `TopologyCatalog`: components, workspaces, lifecycle, docs, and source paths;
- `PolicyCatalog`: applicable policy sources and independent boundary axes;
- `ActionCatalog`: runnable effects and validation contracts;
- `CapabilityCatalog`: skills, providers, composition, and evidence maturity;
- `WorkspaceCheckCatalog`: registered workspaces and supported check phases;
- `GoalCatalog`: owner-local goal identities and, after 2A, continuation state;
  and
- `AgentSystemIndex`: only the catalog identities, versions, input digests,
  conflicts, and query routes.

Generate checked Markdown for people and portable JSON for tools. Every output
has stable ordering, schema/derivation version, source-relative paths, input
digests, bounds, and a check mode. An omitted or invalid eligible skill,
workspace, goal, or component must break the cheap completeness test.

Use tracked `MODULE.bazel` roots as workspace facts and reconcile the intended
ignore, Bzlmod override, docs, and full-check projections. Replace manual
counts with assertions over registered members. Absorb `readme_tree` into the
new topology projection or deprecate it after the replacement passes fidelity
and cost gates.

### 2C. Provide one bounded zero-context read

Implement an offline `agent_system status/context --json` interface that joins
only the slice relevant to a path, label, or task. It returns the architecture's
`ContextCapsule`, including safe next discovery actions. A human rendering
uses the same data.

The command must work without Cordis and may be exposed through MCP as an
adapter. It never owns facts or mutates source. Git, Bazel, Cordis, goals, and
optional connectors each have explicit unavailable states.

### 2D. Make runtime control and isolation independently safe

- Start a fixed health/status/control kernel before optional extension code.
- Load packages asynchronously with per-package deadlines and explicit
  loading, ready, degraded, failed, timed-out, draining, and disabled states.
- Publish desired and observed revision, runtime incarnation, contract hash,
  bounded error, queue/in-flight state, and catalog ETag.
- Add cross-process task namespace, lock, lease, and expected-revision checks.
- Run generated/untrusted packages in a worker or subprocess with sanitized
  environment, descriptor-bound filesystem access, explicit network/process
  capabilities, one deadline, output/resource quotas, revocable lease, and no
  maintained-project write authority. Keep in-process execution only for a
  reviewed trusted tier with an explicit contract.
- Keep healthy packages and native fallbacks available when one extension
  fails or never settles.

### Acceptance signals

- From a clean root, one bounded offline call answers: where am I, what policy
  and owner apply, what work is active, which capabilities exist, what effects
  and checks are relevant, what is authorized, what is degraded, and where the
  facts came from.
- Removing an optional provider produces a structured unavailable field rather
  than losing the whole capsule.
- Catalog generation performs no network or stateful operation and stays
  within a committed byte/time budget without a transitive `//:docs` query.
- A never-settling runtime package cannot block status; two simultaneous tasks
  cannot observe or mutate each other's scratch namespace.
- Runtime fault fixtures cannot impair the fixed kernel, another package,
  another task, or maintained source before dynamic effect admission exists.
- Every injected goal publication interruption returns the prior valid record
  or an idempotent recovery to the intended record.

### Resource guard

The index stores references and digests, not full skill bodies, schemas, logs,
or README content. No database, vector store, long-lived daemon, or mandatory
MCP is introduced. Regeneration reads bounded source authorities directly.

## Phase 3: add advisory planning, admission, and reusable evidence

**Priority:** P0/P1. **Depends on:** Phase 2 context and action catalogs.

### 3A. Resolve the smallest sufficient plan

Given intent, affected paths/labels, criteria, context, and prior receipts,
produce a deterministic `ImpactPlan` containing:

- selected capabilities/providers and explicit reasons;
- required and forbidden effects, authority, credentials, and network;
- changed and reverse-affected targets at the narrowest practical scope;
- conservative minimum checks plus coverage gaps;
- expected/maximum cost, cacheability, ordering, concurrency, and escalation;
  and
- reusable evidence with applicability decisions.

The plan binds canonical candidate, declared new task files, catalog/config
and profile digests, selected and omitted targets, coverage gaps, cost class,
and escalation. Named profiles are `changed/fast`, `workspace`,
`fresh/evidence`, `full/audit`, and `diagnose`. Unknown source, root module,
rc, toolchain, or generator changes conservatively expand or refuse scope.

The planner is advisory. Goal criteria or agent judgment owns semantic
sufficiency; selection never owns authority.

### 3B. Enforce effects at provider gateways

- Admit action contracts against exact authority, subject, environment,
  pre-state, and budgets.
- Preserve a very cheap path for declared read and hermetic compute work.
- Require prepare/validate/authorize/execute/verify for remote write or destroy.
- Reject unknown actions, relevant catalog conflicts, stale pre-state, wrong
  environment, weakened safety flags, and missing cancellation/cleanup budget.
- Consume the Phase 2 runtime-isolation substrate and enforce each admitted
  action's narrower filesystem, process, network, credential, and write
  contract at invocation.

### 3C. Emit candidate-bound evidence

Providers emit immutable action receipts. Execution owns a criterion-neutral
`ValidationSet` with exact candidate, validation profile, impact/check
identities, canonical sanitized arguments, working scope,
provider/config/toolchain/policy digests, clean pre/post state, results,
coverage, duration, output bounds, limitations, and task-local raw-log
digests/pointers. Sanitized arguments contain placeholders or secret
references, never credential-bearing argv.

The goal/task owner separately creates an `EvidenceAssertion` that applies one
or more validation sets to an exact criterion revision and semantic verdict.
Delivery consumes the validation set; goals consume both. The same immutable
set may support non-goal delivery and multiple criterion assertions.

Receipt applicability uses relevant dependency-slice digests. Tree-bound
checks may survive a message-only commit rewrite; commit-bound checks may not.
Changes to relevant base, tree, config, policy, contract, toolchain, generator,
environment, or remote generation make evidence stale. Global catalog churn
alone triggers reevaluation, not a rerun.

### 3D. Make exhaustive checks structured and resumable

Generate the full-repository workspace matrix from the checked inventory.
`full-repo-check` emits incremental versioned JSON plus generated Markdown,
supports workspace/phase selectors, binds exact inputs and profile, records
the target-universe count, rejects zero or unexpected reductions, and resumes
only against identical inputs. Interrupted or unexecuted work remains explicit
rather than being inferred from a final report.

### Acceptance signals

- Representative edit scenarios select the expected minimum checks and explain
  omissions; deliberately changed dependencies invalidate the right receipts.
- Unknown or mismatched effects fail before execution. Read-only discovery
  remains fast and degrades with useful typed gaps.
- Fault fixtures cover wrong environment, stale authority/receipt, secret
  output, symlink escape, redirect to private resources, infinite runtime code,
  cancellation, concurrent writers, and cleanup failure.
- A later agent can reuse exact applicable evidence without rereading raw logs
  or rerunning an unchanged check.
- Identical canonical inputs yield the same impact-plan digest. An interrupted
  broad run resumes without claiming unexecuted coverage.
- One validation set can support delivery and several goal assertions without
  mutation or loss of exact-subject identity.

### Resource guard

Do not persist raw BEP, environment dumps, or unrestricted subprocess output.
Batch compatible Bazel checks, deduplicate identical concurrent checks, and
cap receipts. Benchmark the persistent Bazel server's warm-query latency and
keep scratch, output-base, and loopback policy aligned with server lifetime.

## Phase 4: join durable work, delivery, review, version, and release

**Priority:** P1. **Depends on:** Phase 3 subject-bound receipts.

### Deliverables

- Extend the next goal record version with structured stable defect/uncertainty,
  hypothesis, subject, affected criteria, regression references, prior
  attempt, dominant failure, measurable delta, next action, blocker, and
  resume-condition fields. Keep rich explanation in Markdown.
- Generate a repository-wide goal catalog with globally unambiguous owner-root
  identities and a bounded continuation packet. Do not centralize goal
  mutation.
- Make repository delivery consume a successful candidate-bound validation
  set while leaving check selection with the task/goal policy.
- Replace raw remote-rewrite handoff with a typed, non-authorizing
  synchronization or rewrite-authorization receipt and make provider ownership
  unambiguous.
- Join review finding, disposition, stable defect, fix candidate, regression
  label, and delivery identity through durable references.
- Define a version-owned `{version, channel, commit, tree_state}` handoff to
  release. Keep mutable `head` snapshots explicit and make formal release/tag
  publication a separately authorized, verified operation.
- Add a versioning-owned `ReleaseRefPlan` and provider-neutral guarded
  publisher. It consumes the exact reviewed plan, requires distinct
  `release-refs` authority, fetches expected remote state, uses explicit leases
  and atomic multi-ref publication where required and supported, and emits a
  `ReleaseRefReceipt`. Unsupported atomicity or observation is an explicit
  refusal or unknown, never generic success.
- Bind formal release manifests to version subject, validation set, bundle
  head, artifact/changelog digests, release-ref receipt, and deployment
  observation.

### Acceptance signals

- A fresh agent can discover every maintained open goal and resume one without
  prior path knowledge or free-form archaeology.
- `prepare -> publish` cannot succeed with only a caller-asserted head; the
  supplied validation set must match the exact candidate and policy.
- A review outcome is traversable to the delivered fix and regression without
  making delivery the goal or test owner.
- Formal release identity, bundle head, version/channel, immutable artifact
  names, remote refs, and deployment manifest agree; unrelated tags do not
  truncate changelogs.
- The same reviewed release-ref plan can publish and verify a nightly tag or a
  release branch/tag pair; an existing immutable release tag never moves.
- Interrupted multi-file goal publication has a supported deterministic doctor
  or recovery result.

### Resource guard

Receipts are local audit evidence, not cryptographic attestation against a
same-user forger. Do not merge goal, delivery, versioning, release, or review
into one orchestrator. Each consumes typed references and retains its own
authority.

## Phase 5: close the learning loop and optimize from measurements

**Priority:** P1/P2. **Depends on:** stable context, plans, and receipts.

### 5A. Normalize behavioral evidence

- Give every skill evaluation case and requirement assertion a stable ID,
  provenance, metric, and evidence tier.
- Reject duplicate normalized cases and publish a coverage matrix distinguishing
  configured, routed, fixture-tested, live, stale, and unverified behavior.
- Test routing across the whole skill graph: positive, adjacent negative,
  inert payload, exclusion, conflict, and composition cases.
- Add deterministic writable fixtures for high-risk Git, Bazel, forge,
  Terraform, Ansible, secret, and runtime trajectories.
- Keep live stochastic comparisons manual or scheduled, never in ordinary
  wildcard tests; bind model, skill, catalog, fixture, and judge identities.

### 5B. Capture friction and promote regressions

At task close, retain a bounded sanitized record of selected/considered
capabilities, conflicts, missing providers, avoidable reads or commands,
failed assumptions, verification latency, and exact public evidence identity.
Aggregate by stable defect signature.

Repeated issues create a `LearningProposal`, never an automatic edit. Adoption
requires an owner, minimized public reproducer, regression, contract, fallback,
resource budget, validation, delivered revision, and deprecation/retirement
rule. Runtime promotion creates an isolated source candidate and enters the
ordinary layout, validation, review, and delivery path.

### 5C. Optimize progressive disclosure

- Refactor the largest skills into compact routing/invariant cores with
  conditional references after behavior baselines exist.
- Load catalog summaries and content digests before bodies and schemas; use
  ranged resources for large source or logs.
- Cache launch and projection artifacts by exact provenance; lazy-load optional
  runtime packages after measuring phase latency.
- Promote the advisory resolver to authoritative routing only if predeclared
  correctness and critical-boundary thresholds beat the current selector.
  Otherwise narrow or abandon it.

### Acceptance signals

- Every promoted repository lesson traces from attempts to a minimized
  regression, reviewed owner change, exact delivered digest, catalog update,
  and retirement path.
- Representative tasks reduce cold-start reads, context bytes, repeated checks,
  and resume steps without increasing unsupported claims or unsafe actions.
- All critical effect/authority boundary fixtures pass with zero misroutes
  before routing enforcement is enabled.
- Unchanged instruction bodies and schemas are not transferred twice; context
  and status remain within committed budgets.

### Resource guard

No transcripts, credentials, personal information, unreviewed runtime values,
or secret-bearing outputs enter learning records. Measurements precede
storage or runtime optimization. A new service, database, or model-assisted
index needs a separate decision review proving that bounded files and catalogs
are insufficient.

## Cross-phase verification matrix

| Property                              | First gate                  | Continuing regression                            |
| ------------------------------------- | --------------------------- | ------------------------------------------------ |
| One authority per fact                | Phase 0 review              | Catalog conflict and duplicate-authority checks  |
| Public/prohibited-content distinction | Phase 1 schema fixtures     | Candidate secret and personal-information checks |
| Complete owner-local inventory        | Phase 2 catalog check       | Normal repository hygiene                        |
| Zero-context orientation              | Phase 2 task scenarios      | Calls, bytes, latency, and correctness budget    |
| Effect-safe execution                 | Phase 3 fault fixtures      | Gateway negative tests                           |
| Reusable exact evidence               | Phase 3 applicability suite | Delivery and goal subject checks                 |
| Traversable delivery/release          | Phase 4 integration tests   | Exact remote/ref/postcondition verification      |
| Reviewed accretion                    | Phase 5 lineage tests       | Promotion and retirement audit                   |

## Explicit non-goals

- No central mutable repository brain or hand-maintained mega-manifest.
- No universal mega-skill, eager loading of every instruction, or automatic
  policy synthesis from transcripts.
- No mandatory daemon, MCP server, database, vector store, or network call for
  basic repository orientation.
- No physical relocation of components solely to mimic conceptual cohesion.
- No authority inferred from ownership, selection, planning, a receipt, a
  passing test, tool availability, or earlier unrelated permission.
- No goal ceremony for simple one-step work.
- No secret or personal payloads and no unbounded or unreviewed raw log/state
  payloads in durable evidence; bounded reviewed public facts are allowed.
- No broad build, full-repository test, or billable live evaluation by default.
- No enforcement rollout without a report-only period, adversarial fixtures,
  measured improvement, and a reversible fallback.
