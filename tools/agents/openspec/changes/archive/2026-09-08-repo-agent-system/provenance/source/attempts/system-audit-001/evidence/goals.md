# Goal subsystem and durable-learning audit

## Audit binding

- Goal: `repo-agent-system`
- Goal resource version: `5`
- Goal generation: `1`
- Lifecycle generation: `4`
- Criteria revision: `2`
- Criteria digest:
  `sha256:2ac2db1242f5d3358e433b3499da5a622d06bdec49bfa690dd34cf3205e28f34`
- Goal-state digest:
  `sha256:193be5b38881faebc349f9ae1d273e24fac5d5925a9a4402b24706394ffaeb3a`
- Attempt: `system-audit-001`
- Method: read-only inspection of the goal skill, API, filesystem store, CLI,
  architecture, evaluations, and both goal records currently present outside
  ignored output. No canonical goal state or maintained source was changed.

## Verdict

The goal subsystem has a strong integrity kernel: stable IDs, exact generation
and digest bindings, immutable closed attempts, optimistic concurrency,
single-writer coordination, generated projections, and a pure relationship
analysis. Those are unusually good foundations.

Its main limitation is at the system boundary. The skill describes a rich
reasoning, evidence, retry, delivery, and learning protocol, while the durable
API records only a thin lifecycle and artifact index. At the same time, every
catalog operation is scoped to one explicitly supplied `goals/` directory.
Consequently a zero-context agent cannot discover the repository's goals as a
whole, cannot obtain a compact continuation packet, and cannot mechanically
answer the questions that matter most after a context loss: what exact defect
remains, what changed strategy, which checks are stale, what was actually
tested and delivered, and what durable lesson should alter future work.

The highest-leverage improvement is therefore not more narrative guidance. It
is a repository-wide discovery and continuation plane backed by a small set of
structured evidence and learning facts. Markdown should remain the rich human
narrative, but the facts used to resume, schedule, accept, and learn must be
queryable and validated.

## Ranked findings

### 1. P0 — Owner-local goals have no repository-wide discovery or identity

Facts:

- Project goals are intentionally stored below the narrowest owner at
  `<owner-root>/goals/<goal-id>`
  (`projects/goal/skills/goal/SKILL.md:32-35`).
- `list`, `graph`, and catalog validation each require one explicit
  `--goals-root` (`projects/goal/cmd/goal/command.go:123-177` and
  `:359-393`). A session binding explicitly does not create a repository-global
  catalog (`projects/goal/cmd/goal/README.md:56-68`).
- Relationship references are name-only and valid only in the same catalog
  (`projects/goal/api/v1alpha1/types.go:35-49` and
  `projects/goal/skills/goal/references/graph-organization.md:22-39`).
- The repository already has goals under two different owner roots:
  `projects/agents/goals/repo-agent-system` and the legacy
  `projects/mcp_cordis/goals/runtime_extensions`. The root README links neither
  the agent system nor any goal inventory (`README.md:6-14`).

Consequence:

A fresh agent must know paths in advance or perform an ad hoc repository scan.
Cross-project prerequisites and supersession cannot be represented, duplicate
names in different catalogs are ambiguous, and a repository-level graph or
ready-work queue cannot be derived. Durable knowledge is owner-local but not
system-visible.

Design delta:

1. Keep each goal canonical under its owner; do not centralize ownership.
2. Add a generated repository catalog containing a globally unambiguous key,
   preferably `{ownerRoot, name}` or a repository-relative goal URI. Treat the
   catalog as a projection, never a second mutation authority.
3. Register project goals through Bazel-owned metadata (`goal_record` or a
   similar provider) and generate registrations from the conventional layout
   to avoid hand-maintained indexes. Include only the current task's ignored
   workspace catalog through its session binding, not in the committed project
   catalog.
4. Add `goal catalog`, `goal resolve`, and `goal ready` projections that can
   report owner, schema, lifecycle, dependency state, active attempt, update
   time, next action, and truncation deterministically across owners.
5. Make the root agent/human entry points link to this catalog and the system
   model. Run duplicate-ID, broken-reference, schema, and owner-root checks in
   the normal repository validation graph.

Acceptance signal: from a clean root with no path knowledge, one bounded
command identifies every maintained open goal, resolves its canonical path,
and produces a complete cross-owner dependency analysis without scanning
private or ignored trees indiscriminately.

### 2. P0 — The skill's evidence contract is materially richer than the stored model

Facts:

- Before work, the skill requires a target uncertainty or stable defect,
  hypothesis, inputs, intended evidence, affected criteria, and fixed
  regressions; closeout requires the dominant failure, measurable movement,
  feedback bottlenecks, and a stable failure count
  (`projects/goal/skills/goal/SKILL.md:53-77,101-109` and
  `references/lifecycle-and-evidence.md:7-18,48-65`).
- `AttemptSpec` stores only input bindings and `workType`; `CloseReview` stores
  only a decision and criterion verdicts. `GoalStatus` has no blocker reason,
  waiting-on reference, resume condition, or next action
  (`projects/goal/api/v1alpha1/types.go:76-86,118-164`).
- A new attempt may use generated placeholder plan and result content
  (`projects/goal/internal/fsstore/checkpoint.go:300-345`). A pass may cite
  `plan.md` or `result.md`; validation checks that the file is frozen, not that
  it identifies the tested subject (`projects/goal/api/v1alpha1/validation.go:
  299-367`).
- The accepted digest is explicitly only the digest of `result.md`, not the
  external candidate or operation (`projects/goal/skills/goal/references/
  record-format-v1alpha1.md:148-156`).

Consequence:

The store proves that an agent froze some prose and made a structurally valid
claim. It cannot prove which source tree, deployment, PR, render, or operation
the claim describes, whether the fixed regressions ran, whether delivery
changed identity, or whether a recurring defect crossed the strategy-reset
threshold. Correctness and accretion remain dependent on the current agent's
memory and Markdown discipline.

Design delta:

- Extend the next record version with a minimal structured work packet:
  `targetRef` (stable defect or uncertainty), optional hypothesis,
  `affectedCriterionRefs`, `regressionCheckRefs`, immutable input refs,
  `subjectRef`, prior-attempt refs, and required capability/authority hints.
- Add structured evidence assertions that bind a criterion revision, exact
  subject identity, method, receipt, and observed postcondition. Preserve the
  Markdown artifact as explanation rather than parsing it as data.
- Add close-review fields for dominant failure, measurable delta, feedback
  cost, next action, and strategy decision. Derive recurrence counts by stable
  target ID and warn or reject a third unchanged strategy.
- Give inactive execution states a reason, `waitingOn`, and explicit resume
  condition. Require superseded outcomes to identify lineage.
- Integrate repository-delivery receipts as typed subject/delivery evidence.
  Achievement should bind both the explanatory result digest and the exact
  accepted/delivered subject identity.

Near-term, before a schema revision, add a strict protocol linter that refuses
placeholder plans/results at close, checks required plan/review fields, and
warns when acceptance has no exact subject or delivery receipt. This is cheap
feedback, not a substitute for the structured model.

Acceptance signal: given only canonical resources, a fresh agent can compute
the unresolved criterion set, stable defect recurrence, stale evidence,
highest-leverage next action, and exact accepted/delivered identity without
interpreting free-form prose.

### 3. P0 — A tolerated partial commit can brick the record, but there is no recovery surface

Facts:

- Checkpoint publication advances `goal.yaml` before publishing or finalizing
  attempt files. The architecture deliberately permits an interruption to
  leave an invalid record and says only to validate before resuming
  (`projects/goal/docs/architecture.md:60-82`).
- Every normal checkpoint first calls `loadAndValidate`; an invalid record is
  rejected before another mutation can repair it
  (`projects/goal/internal/fsstore/checkpoint.go:25-48`).
- The command set contains `validate` and `render` but no `doctor`, `recover`,
  or `repair` operation (`projects/goal/cmd/goal/command.go:73-85`). Direct YAML
  emulation is prohibited by the skill contract
  (`projects/goal/skills/goal/references/record-format-v1alpha1.md:82-92`).

Consequence:

Failing closed is correct, but the only supported state after the failure is an
error report. A transient filesystem or process failure can make the durable
coordination record unusable precisely when an agent most needs deterministic
recovery. Manual edits would bypass the store's own concurrency and digest
guarantees.

Design delta:

- Preferred end state: publish an immutable complete record revision and
  atomically switch one current-revision pointer. Old revisions then provide
  deterministic rollback and audit history without exposing a mixed state.
- Transitional end state: retain a write-ahead transaction manifest and staged
  content until a `goal recover` command can deterministically finish or roll
  back. `goal doctor` should classify pre-commit, committed-but-incomplete, and
  projection-only failures and emit a machine-readable recovery plan.
- Return stable JSON error codes plus committed version, affected paths, and
  allowed recovery operation. Never require parsing an English error to decide
  whether a retry is stale.

Acceptance signal: fault injection at every publication boundary followed by a
fresh process always produces either the old valid record or a deterministic,
idempotent recovery to the new valid record; no direct file edit is required.

### 4. P1 — The CLI exposes storage mechanics instead of an agent continuation protocol

Facts:

- `checkpoint` overloads attempt creation/update/close, evidence import,
  criteria replacement, outcome, and execution in one ten-flag surface
  (`projects/goal/cmd/goal/command.go:248-282`). Criteria replacement is
  mutually exclusive and requires a separately paused goal with no active
  attempt (`projects/goal/internal/fsstore/checkpoint.go:52-66`).
- `init` accepts only repeated criterion statements, generating generic IDs and
  a generic evidence method; it cannot atomically consume the structured
  criteria format (`projects/goal/cmd/goal/command.go:91-120` and
  `projects/goal/internal/fsstore/store.go:188-211`).
- The current first attempt already sits at Goal resource version 5,
  lifecycle generation 4, and criteria revision 2, illustrating the ceremony
  needed to initialize richer criteria before work
  (`projects/agents/goals/repo-agent-system/goal.yaml:3-25`).
- A mutation returns only `{goalID, goalRef, resourceVersion}`. If the CLI
  generates an attempt ID, that ID is absent from the receipt
  (`projects/goal/internal/fsstore/store.go:100-104` and
  `projects/goal/internal/fsstore/checkpoint.go:109-119`).
- `show` returns only attempt summaries. One shared `limit` truncates both
  criteria and attempts, while `Returned` and `Total` describe attempts only
  (`projects/goal/internal/fsstore/store.go:349-426`).

Consequence:

Agents must manufacture multiple YAML/Markdown files, parse a sparse receipt,
re-read state after mutations, remember session-root and ID flags, and open
sidecars individually to reconstruct context. This increases tool calls,
tokens, stale-version windows, and recovery mistakes.

Design delta:

1. Keep the atomic `checkpoint` primitive internally, but expose intent-level
   commands: `criteria init/replace`, `attempt open`, `evidence add`,
   `attempt close`, and `status set`.
2. Add `init --criteria-file` and a non-mutating `attempt scaffold` that emits
   validated plan, evidence, and review templates under task output.
3. Add `goal resume --goal/--session` returning one bounded continuation
   packet: exact versions/digests, all current required criteria, active plan,
   latest review, unresolved checks, stable defects, next action, and explicit
   per-section truncation metadata.
4. Return a versioned mutation receipt containing old/new versions, generated
   IDs, bound digests, changed fields, committed status, and valid next actions.
   Emit the same versioned envelope for errors, with stable codes such as
   `STALE_VERSION`, `COMMITTED_INCOMPLETE`, and `INVALID_TRANSITION`.
5. Allow one request-file input for complex mutations and explicit
   `name=source` evidence mapping. This avoids shell quoting and basename
   collisions while preserving reviewable intent.

Acceptance signal: an agent can initialize rich criteria, open an attempt, and
later resume it with one mutation plus one read; it never needs to discover a
generated ID or next resource version through a second ad hoc file read.

### 5. P1 — Durable failures do not automatically become reusable learning

Facts:

- The skill requires stable defect names, strategy changes after recurrence,
  preserved failures, and a feedback-cost audit
  (`projects/goal/skills/goal/SKILL.md:75-77,106-109`).
- The legacy MCP Cordis goal manually maintains exactly the useful learning
  shape: candidate, cause, strategy delta, and regression guard. One entry even
  records that a discovered failure changed `repo-delivery` policy
  (`projects/mcp_cordis/goals/runtime_extensions/failure_ledger.md:5-36`).
- None of these concepts exists in the v1alpha1 resource types; they can only
  be buried in result/evidence prose (`projects/goal/api/v1alpha1/types.go:
  118-156`).

Consequence:

The system records what happened but does not close the loop into a skill,
policy, test, tool, or architecture correction. Valuable lessons are found
only by a human or agent that already knows which historical record to read;
duplicated failures remain cheap to repeat.

Design delta:

- Add a small `GoalLearning`/`LearningProposal` resource (or equivalent
  canonical section) with stable ID, source attempt/evidence refs, scope,
  generalizable claim, proposed destination owner, guard/eval ref, and
  `proposed|adopted|rejected|superseded` status.
- Require owner review before promoting a local lesson into global AGENTS
  policy or a shared skill. This prevents every incidental observation from
  polluting repository-wide memory.
- Derive a repository learning queue: recurring defects without a guard,
  accepted learnings without a target change, and changed skills without a
  corresponding eval or deterministic regression.
- Link adopted learning back to the exact source evidence and forward to the
  enforcing test/skill/tool, creating an inspectable closed loop.

Acceptance signal: the MCP Cordis failure that tightened `repo-delivery` can be
traced mechanically from defect evidence to learning decision to skill change
to regression, and a future agent can query that chain from the root catalog.

### 6. P1 — The repository's existing durable goal cannot follow the documented migration path directly

Facts:

- The only pre-v1alpha1 durable project record is
  `projects/mcp_cordis/goals/runtime_extensions`; it is explicitly durable and
  contains 11 attempts, criteria, requirements, evidence, artifacts, and a
  failure ledger (`projects/mcp_cordis/goals/runtime_extensions/README.md:
  1-4,24-96`).
- Migration derives the new goal ID from the source directory basename,
  rejects an overriding ID that differs, and tells the caller to rename an
  invalid legacy directory (`projects/goal/internal/fsstore/migrate.go:
  326-341`). `runtime_extensions` is invalid under the hyphen/dot record-ID
  grammar.
- The import must target a non-overlapping directory, so a source already at
  the desired `<owner>/goals/<id>` boundary cannot be converted into that
  boundary in one supported operation
  (`projects/goal/skills/goal/references/promotion-and-migration.md:27-45`).

Consequence:

The repository's best evidence of durable goal practice remains invisible to
the new CLI and global graph. The user must rename/copy maintained history and
invent a cutover sequence, exactly the kind of protocol judgment the migration
tool is meant to own.

Design delta:

- Permit `migrate --new-goal-id` while preserving the immutable source
  basename/path and digest in provenance.
- Add a guided owner-local cutover that stages the converted goal, preserves
  the legacy tree under an explicit immutable legacy path or digest-bound
  snapshot, verifies links, then atomically publishes the canonical target.
- Allow an explicit mapping file for unambiguous criteria, attempts, stable
  defects, and receipts. Unmapped content remains an immutable legacy snapshot;
  it must not be guessed into structured fact.
- Make legacy records appear in repository catalog output as
  `Legacy/RequiresMigration` rather than disappearing.

Acceptance signal: the MCP Cordis record can be migrated without changing any
legacy byte, without an ad hoc manual rename/copy protocol, and without losing
its failure/strategy history or final delivery identity.

### 7. P2 — Durable evidence retention is advisory rather than verifiable

Facts:

- Project evidence is required to survive through repository-relative links,
  stable external references, or regeneration instructions
  (`projects/goal/skills/goal/SKILL.md:86-91`).
- Attempt manifests digest only copied Markdown files. External subject,
  repository links, regeneration commands, and operation receipts remain text
  (`projects/goal/skills/goal/references/record-format-v1alpha1.md:139-170`).
- Promotion rejects a limited class of absolute paths but explicitly leaves
  semantic privacy and credential review to callers
  (`projects/goal/cmd/goal/README.md:105-109`). There is no check that a
  relative evidence link exists, is tracked, is not ignored, or names a
  reproducible Bazel target.

Consequence:

A project goal may validate today while its acceptance-critical evidence is
ignored, moved, deleted, inaccessible, or impossible to regenerate. Its YAML
integrity can remain green after its epistemic basis has rotted.

Design delta:

- Add typed artifact references for `embedded`, `repository`, `external`, and
  `regenerable` evidence, each with digest/version and retention policy.
- Package project goal records as Bazel targets. Validate repository references
  for existence and tracked ownership; validate regeneration references as
  real targets; report external references as explicitly unverifiable rather
  than silently valid.
- Add `goal validate --retention --subject` and include it in project goal CI.
  Keep raw logs in ignored output when appropriate, but retain a deterministic
  receipt and regeneration recipe for acceptance-critical facts.
- Define cleanup semantics for ephemeral workspace goals and expose only a
  dry-run garbage-collection report by default.

Acceptance signal: removing or ignoring an acceptance-critical artifact makes
the owning project goal validation fail with the exact broken reference and
regeneration option.

### 8. P2 — Skill, docs, and CLI compatibility are not tested end to end

Facts:

- The goal eval target validates only Promptfoo configuration and performs no
  model call. The documentation explicitly says a longitudinal fresh-session
  fixture is still missing (`projects/goal/skills/goal/evals/README.md:7-20`).
- The project README, command README, architecture, record-format reference,
  and concurrency reference repeat mutable file, locking, atomicity, and
  Kubernetes-boundary details. The docs target merely groups Markdown
  (`projects/goal/skills/goal/BUILD.bazel:18-22`).
- The discovery symlink is a good single-authority pattern: canonical skill
  content stays under `projects/goal/skills/goal`, while `.agents/skills/goal`
  is only a relative link (`projects/goal/README.md:47-54`).

Consequence:

Store tests can remain green while the skill requests an operation the CLI
cannot express, an example goes stale, or a fresh agent fails to resume. The
most important behavior spans precisely the seam that current tests exclude.

Design delta:

- Generate JSON Schema/OpenAPI-like resource schemas, CLI reference, and
  version compatibility metadata from one API/command authority. Have the
  skill declare its minimum tool and record-format capabilities.
- Execute every documented CLI example in an isolated Bazel test.
- Add a deterministic longitudinal fixture covering init with rich criteria,
  attach, open attempt, isolated worker evidence, stale publication, strategy
  reset, exact-subject acceptance, delivery, context loss, and fresh-session
  resume. Most of this can test the protocol without a networked model.
- Add a small live model eval only for the decisions that cannot be reduced to
  state-machine tests. Keep it outside normal wildcard tests, but make its
  absence and last result visible in the system health report.
- Reduce duplicated prose: API/schema owns resource fields, store architecture
  owns transaction boundaries, CLI help owns flags, and the skill owns agent
  decisions. Other documents should link or use generated excerpts.

Acceptance signal: changing a flag, resource field, lifecycle invariant, or
skill workflow breaks a focused compatibility test before it can drift into a
fresh agent session.

## Recommended goal-system layer contract

The goal subsystem should participate in the repository's abstraction tower
without absorbing responsibilities that already have better owners:

| Layer | Goal-system responsibility | External authority linked by typed receipt |
| --- | --- | --- |
| Intent | Goal, criteria, assumptions, outcome | User request and repository policy |
| Topology | Owner-local identity and relationships; global generated catalog | Project ownership metadata |
| Capability | Attempt requirements and authority/cost hints | Skill/capability registry and security policy |
| Execution | Immutable work packet, coordinator/worker bindings, lifecycle | Agent harness and tool runner |
| Evidence | Criterion assertion against exact subject and check receipt | Bazel, reviewers, domain tools |
| Delivery | Accepted subject plus delivered identity | `repo-delivery` or stateful operation adapter |
| Learning | Stable defect, strategy delta, generalized learning proposal | Owning AGENTS rule, skill, test, or tool |
| Projection | Bounded resume packet, catalog, ready queue, human README | CLI/MCP and documentation surfaces |

The critical invariant is that each layer has one mutation authority and all
cross-layer statements are typed, digest/version-bound references. The goal
record coordinates and explains; it must not become a second scheduler, build
system, policy engine, or delivery adapter.

## Cheapest deterministic feedback, in order

1. Add `init --criteria-file`, generated templates, rich mutation receipts,
   structured error codes, and separate truncation counts.
2. Add a read-only `goal resume` packet and root `goal catalog` projection.
3. Add strict placeholder/subject/retention linting and executable doc examples.
4. Permit migration to a new valid ID and report legacy records in catalog
   health.
5. Add typed subject, regression, defect, next-action, and delivery-receipt
   fields while preserving Markdown narrative.
6. Add deterministic transaction recovery or revision-pointer publication.
7. Add learning proposals and the closed-loop repository learning report.

## Foundations to preserve

- Separate workspace and explicitly durable project goals.
- Stable goal IDs independent of chat/session identity.
- Exact generation, lifecycle, criteria, and portable digest bindings.
- One canonical writer; isolated workers never mutate goal state.
- Immutable closed attempts and preserved failed evidence.
- Opaque compare-and-swap resource versions and stale-publication refusal.
- Generated, bounded README projections distinct from canonical resources.
- Pure, complete graph analysis and explicit separation from live scheduling.
- Canonical project-owned skills with discovery symlinks rather than copied
  instruction trees.

These strengths should become the integrity substrate for the repository-wide
agent system. Replacing them with a central mutable plan document would lose
more than it gains; the needed improvement is a coherent discovery,
continuation, evidence, recovery, and learning interface above them.
