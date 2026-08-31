# Delivery, review, and release workflow audit

## Audit binding

- Goal: `repo-agent-system`
- Attempt: `system-audit-001`
- Goal resource version: `5`
- Goal generation: `1`
- Lifecycle generation: `4`
- Criteria revision: `2`
- Criteria digest:
  `sha256:2ac2db1242f5d3358e433b3499da5a622d06bdec49bfa690dd34cf3205e28f34`
- Goal-state digest:
  `sha256:193be5b38881faebc349f9ae1d273e24fac5d5925a9a4402b24706394ffaeb3a`
- Scope: repository delivery, remote rebase, review, validation evidence,
  versioning, release packaging/deployment, and durable-learning handoffs.
- Method: read-only inspection of canonical skills, implementations, tests,
  Bazel declarations, and goal contracts. No Git, forge, network, maintained
  file, or goal-state mutation was performed.

## Executive verdict

The repository has unusually careful mechanisms at the individual-command
level. `repo_delivery` binds exact repository, ref, head, tree, lease,
pull-request, and path-scope state; its review mutations use exact inventory
digests and one-use reply authority. Versioning validates and creates local
refs transactionally. Goal records already bind criteria and evidence to exact
subjects.

The main defect is therefore not weak components. It is that the seams between
them are mostly procedural prose or copied scalar values. The repository can
prove that the candidate being published is the candidate named by the caller,
but it cannot mechanically carry what validations ran, why a remote rewrite
was authorized, whether a remote review reached a terminal state, how feedback
became a regression, or how a global version became a released artifact.

The right design is **not one giant delivery orchestrator**. Keep semantic
judgment, validation selection, forge operations, version calculation, Bazel
execution, and goal retention with their current owners. Add a small family of
versioned, redacted, subject-bound receipts between those owners, plus generated
projections for agents and humans.

## Existing strengths to preserve

1. The preparation receipt already records a schema, repository fingerprint,
   endpoint digests, forge, base/head refs and OIDs, prepared tree, expected
   remote head, pull-request expectation, and aggregate scope
   ([receipt.go:58-75](../../../tools/repo_delivery/main/go/receipt.go#L58)).
2. `publish` rereads that receipt and requires the caller's validated head to
   equal both the receipt head and current local head
   ([delivery.go:1632-1717](../../../tools/repo_delivery/main/go/delivery.go#L1632)).
3. Review inspection exposes exact PR, review, thread, comment, and reviewer
   identities, including thread and top-level inventory digests
   ([forge.go:150-229](../../../tools/repo_delivery/main/go/forge.go#L150)).
4. The versioning workflow has a single global version authority and guarded,
   clean-tree local ref transactions
   ([versioning/SKILL.md:12-48](../../../tools/versioning/skills/versioning/SKILL.md#L12)).
5. Goal evidence already requires criterion revision and exact tested subject,
   treats command success as insufficient, and calls for operation receipts
   ([lifecycle-and-evidence.md:25-39](../../../projects/goal/skills/goal/references/lifecycle-and-evidence.md#L25)).
6. The full-repository checker runs all registered workspace/phase pairs,
   continues after failures, retains private logs, and reports coverage limits
   ([run_full_repo_check.go:115-215](../../../projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go#L115)).

## Authority handoffs

| Transition | Current authority | Carrier | Verdict |
|---|---|---|---|
| Request to publication | Root policy and `repo-delivery` skill | An implementation request implicitly expands to commit, push, PR maintenance, comments, and review waiting | Policy is readable to an agent, but the selected external-effect scope is not an explicit task/session binding. |
| Prepared candidate to validation | Agent chooses checks; Bazel and other tools execute them | Literal `head_oid`, terminal output, and optional ad hoc logs | Candidate identity is strong; execution provenance and coverage are not carried. |
| Validation to publish | Agent calls `publish` | `--validated-head` plus preparation receipt | Publish checks subject equality, not that checks ran or passed. |
| Divergent remote to GitHub delivery | Agent applies `git-rebase-remote` ownership judgment; `repo_delivery` enforces state | A copied old remote OID in `--replace-remote` | Fresh state is mechanically guarded, but rewrite authority and workflow ownership are split. |
| Publish to review | `repo_delivery` and forge adapter | PR/head IDs and review inventory digests | Exact mutation guards are strong; remote review execution state is absent from the current adapter projection. |
| Review finding to regression | Agent judges feedback; Bazel/tests/skills own durable correction | Prose reply, code/test edit, and final report | The procedure exists, but there is no traversable `finding -> defect -> fix head -> regression` edge. |
| Version to remote refs | `versioning` creates local refs; an unspecified delivery path must push them | Local branch/tag names printed as text | Broken handoff: current `repo_delivery` only publishes one-commit feature branches and explicitly does not follow tags. |
| Version to release artifact | `versioning` owns global version; `tools/release` owns snapshot manifests | No typed connection; release rules accept an arbitrary `release_name` | Exact-HEAD snapshots are valid, but formal version identity and snapshot identity are conflated by naming and docs. |
| Delivery/release to goal evidence | Delivery/version tools own ephemeral receipts; goal owns durable evidence | Manual Markdown copy or link | Possible today, but lossy and non-machine-traversable. |

## Ranked findings

### D1 — P0: validation is asserted, not execution-attested

**Evidence.** The delivery skill deliberately leaves validation selection to
the agent
([repo-delivery/SKILL.md:58-69](../../../projects/agents/skills/repo-delivery/SKILL.md#L58))
and requires all checks after `prepare`
([repo-delivery/SKILL.md:112-134](../../../projects/agents/skills/repo-delivery/SKILL.md#L112)).
Publication then accepts only the literal OID and preparation receipt
([command.go:333-346](../../../tools/repo_delivery/main/go/command.go#L333)).
The receipt has no validation-set or result fields, and `publish` checks only
that the asserted OID equals the receipt and current repository state. The
full-repository runner records command, exit code, duration, and log path, but
its result is an unversioned Markdown report with no candidate commit/tree,
configuration identity, or log digest
([run_full_repo_check.go:225-327](../../../projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go#L225)).

**Consequence.** After context loss, delegation, or review-driven change, an
agent must trust narration or repeat work. Delivery can prevent candidate
substitution but cannot distinguish a real successful validation from a copied
OID. This is a provenance and ergonomics gap, not a claim that the existing
receipt is weak or that a local receipt would be cryptographically trustworthy.

The invalidation model is also too coarse. A message-only amendment preserves
the tree but requires every check to run again
([repo_delivery/README.md:110-128](../../../tools/repo_delivery/README.md#L110));
agent test results are explicitly not cached
([project.bazelrc:55-60](../../../tools/bazelrc/project.bazelrc#L55)).
Some checks really are commit-, parent-, base-, or published-interface-sensitive,
but pure tree checks need not be discarded when their complete subject and
inputs are unchanged.

**Recommendation.** Add a versioned `ValidationSet` receipt emitted by a
controlled execution wrapper. It should bind:

- repository fingerprint and candidate `{base_oid, commit_oid, tree_oid}`;
- stable check ID, exact argv/cwd, validation profile, executor/tool version,
  and declared configuration/toolchain digest;
- subject class (`tree`, `commit`, `base+tree`, or `remote-state`) and target or
  query-expression digest;
- clean pre/post state, start/end time, exit status, bounded result summary,
  log digest/location, and explicit coverage limitations; and
- a policy/profile digest stating which receipts were required, while leaving
  semantic sufficiency and check selection with the agent or goal criteria.

`repo_delivery publish` should consume one or more receipts and reject subject
or policy mismatch. The receipt is an auditable same-user consistency record,
not an authorization token. Reuse a receipt only when every declared subject
and input is unchanged; explain reuse and invalidation in output.

**Acceptance signal.** From a fresh session, an agent can prove which checks
covered the exact candidate without rereading raw logs, and a message-only
rewrite reruns commit-sensitive gates while reusing still-valid tree evidence.

### D2 — P0: the formal release path stops at local refs

**Evidence.** Versioning creates a nightly tag locally and explicitly delegates
pushing to "the repository delivery workflow"
([versioning/SKILL.md:50-65](../../../tools/versioning/skills/versioning/SKILL.md#L50)).
`release-start` similarly creates a local release branch and patch-zero tag.
The delivery CLI, however, describes and exposes only one-commit feature-branch
and PR commands
([command.go:39-85](../../../tools/repo_delivery/main/go/command.go#L39)),
and its exact push explicitly uses `--no-follow-tags`
([git.go:4876-4893](../../../tools/repo_delivery/main/go/git.go#L4876)).
Version mutation commands return a tag or branch name as text rather than a
ref plan/receipt
([versioning command.go:114-176](../../../tools/versioning/cmd/versioning/command.go#L114)).

**Consequence.** The highest-risk release handoff has no canonical executable
path. An agent must invent direct Git commands or reinterpret feature-delivery
policy, losing exact leases, multi-ref consistency, provider routing, and
postcondition evidence. Local atomicity for the branch/tag pair does not extend
to the remote.

**Recommendation.** Keep versioning as the semantic authority, but have its dry
run and apply emit a versioned `ReleaseRefPlan` containing exact commit, channel,
new refs, expected local and remote states, and intended atomicity. Add a
provider-neutral guarded ref publisher, separate from PR delivery, that:

- consumes the exact plan and explicit release-effect authority;
- fetches every affected ref, uses explicit leases and an atomic multi-ref push
  when supported, and fails closed when the required atomicity is unavailable;
- never moves an existing release tag; and
- refetches and emits a redacted `ReleaseRefReceipt` with exact postconditions.

**Acceptance signal.** One dry-run JSON document can be reviewed and then used
unchanged to publish and verify a nightly tag or release branch/tag pair; no
undocumented direct `git push` is needed.

### D3 — P0: global versions and release snapshots have no typed join

**Evidence.** Versioning says its value is shared by the repository and every
first-party project
([versioning/SKILL.md:12-24](../../../tools/versioning/skills/versioning/SKILL.md#L12)).
The release system legitimately creates an exact-HEAD Git bundle
([al_git_repo.bzl:87-126](../../../tools/git/main/bzl/al_git_repo.bzl#L87)),
but `al_release` accepts an unconstrained mandatory `release_name`
([al_release.bzl:112-145](../../../tools/release/main/bzl/al_release.bzl#L112)).
Projects commonly set both the target and release name to mutable `head`
([ci_platform releases:23-35](../../../projects/ci_platform/releases/BUILD.bazel#L23)),
and OCI tag construction uses that name directly
([generator.go:236-256](../../../tools/release/main/go/generator.go#L236)).
The current release READMEs contain metadata only and do not explain whether
these are snapshots or formal releases
([tools/release/README.md](../../../tools/release/README.md)).

There is also a hidden changelog coupling: the bundle imports every tag merged
into HEAD, and generation stops at the first earlier tagged commit, irrespective
of the canonical version/channel
([generator.go:301-327](../../../tools/release/main/go/generator.go#L301)).

**Consequence.** The two mechanisms are not intrinsically conflicting—`head`
can be a useful continuous snapshot—but an agent cannot prove that a formal
artifact's version, stamped commit, bundle head, changelog boundary, OCI tag,
and deployed bytes all describe one release. An unrelated merged tag can also
silently change changelog history.

**Recommendation.** Preserve both systems and make the distinction explicit:

- define `snapshot` and `nightly|release` channels in one typed
  `{version, channel, commit, tree_state}` provider;
- retain `head` only as an explicitly mutable snapshot alias, never as the
  identity of a formal immutable release;
- require formal `al_release` targets to consume the version provider, validate
  bundle-head equality, and use version-aware changelog boundaries;
- bind artifact SHA-256 digests and stamped version/commit into a release
  manifest; and
- have deployment verify the local manifest subject and report the remote OCI
  digest/postcondition. The current deployer only iterates `oras push` calls and
  returns process success
  ([deployer.go:26-59](../../../tools/release/main/go/deployer.go#L26)).

**Acceptance signal.** Given an OCI digest or release page, an agent can trace
back to one global version, exact source commit/tree, validation set, artifact
digest, release ref receipt, and deployment observation. Snapshot consumers can
still intentionally follow `head` without being told it is an immutable release.

### D4 — P1: publication authority is implicit and not session-bound

**Evidence.** Root policy assigns staging, commits, pushes, PR maintenance,
review comments, and final reporting to `repo-delivery`
([AGENTS.md:83-85](../../../AGENTS.md#L83)). The skill says that once requested
implementation is ready it should commit and push without asking
([repo-delivery/SKILL.md:10-14](../../../projects/agents/skills/repo-delivery/SKILL.md#L10)).

**Consequence.** This favors autonomy and avoids repetitive confirmation, but a
request that sounds local can silently acquire network, history-rewrite, human
notification, and long-running review effects. Conversely, a cautious agent may
ask redundant questions because the effective authority is not available as
machine state.

**Recommendation.** Make the policy explicit in both human and agent entry
points and bind one effect scope to the task/session:
`inspect`, `local-change`, `commit`, `publish`, `review-mutate`, `release-refs`,
and `deploy`. A repository default may still authorize ordinary feature
publication; the important properties are visibility, an explicit opt-out,
and fail-closed checks before crossing to a stronger effect class. A new user
request may widen the scope, but a skill relationship must not do so silently.

**Acceptance signal.** The agent can show the current allowed effect classes
without asking the user again, and every mutation tool can reject an operation
outside that binding.

### D5 — P1: rebase execution and rewrite authority have conflicting owners

**Evidence.** The `git-rebase-remote` skill instructs the agent to inspect,
fetch, rebase, rerun checks, push with an exact lease, refetch, and verify
([git-rebase-remote/SKILL.md:16-85](../../../projects/agents/skills/git-rebase-remote/SKILL.md#L16)).
It is packaged only as a prose skill/eval, not an executable producer
([git-rebase-remote BUILD:1-25](../../../projects/agents/skills/git-rebase-remote/BUILD.bazel#L1)).
For supported GitHub delivery, `repo_delivery` says it owns rebase and push
mechanics, but `--replace-remote` depends on the other skill having established
ownership and user authorization
([repo-delivery/SKILL.md:58-69](../../../projects/agents/skills/repo-delivery/SKILL.md#L58),
[repo-delivery/SKILL.md:105-111](../../../projects/agents/skills/repo-delivery/SKILL.md#L105)).

State is not unguarded: `--replace-remote` must equal the freshly fetched
divergent remote OID
([delivery.go:729-765](../../../tools/repo_delivery/main/go/delivery.go#L729)),
and that OID is carried into the receipt and later lease
([receipt.go:350-376](../../../tools/repo_delivery/main/go/receipt.go#L350)).
The missing element is structured authority/provenance, not OID checking.

**Consequence.** An agent faces two overlapping rebase/push workflows and must
translate a safety-critical judgment into a naked scalar. Provider behavior
diverges, repeated fetch/rebase/validation work is likely, and resumption cannot
explain who authorized replacing which exact commit range.

**Recommendation.** On a supported adapter, make `repo_delivery` the sole Git
state/execution owner. Add a non-cryptographic `RewriteAuthorization` receipt
bound to repository, refs, observed remote OID, inspected commit range,
task/goal reference, ownership rationale, explicit effect scope, issuer, and
expiry. The tool must continue to validate all observable state and must never
infer ownership from the receipt alone. Retain `git-rebase-remote` as the
unsupported-provider/Forgejo path, preferably backed by the same hardened Git
core and emitting a corresponding synchronization receipt.

**Acceptance signal.** Routing selects exactly one rebase/push executor, and a
fresh agent can audit the exact rewrite authorization without redoing discovery
or trusting conversational memory.

### D6 — P1: review is exact-state but not a closed learning loop

**Evidence.** The delivery skill already requires valid feedback to become
focused regression coverage and forces exact-head revalidation
([repo-delivery/SKILL.md:38-52](../../../projects/agents/skills/repo-delivery/SKILL.md#L38),
[repo-delivery/SKILL.md:209-217](../../../projects/agents/skills/repo-delivery/SKILL.md#L209)).
Goal records already support fixed regressions and exact-subject evidence
([record-format-v1alpha1.md:124-156](../../../projects/goal/skills/goal/references/record-format-v1alpha1.md#L124)).

The current review projection contains reviews, comments, threads, and requested
reviewers, but no running/completed/failed/cancelled execution state
([forge.go:220-229](../../../tools/repo_delivery/main/go/forge.go#L220)). The
skill therefore has to fall back to bounded observation and an `unverifiable`
result when no external monitor exposes that state
([repo-delivery/SKILL.md:192-207](../../../projects/agents/skills/repo-delivery/SKILL.md#L192)).
The reply receipt records remote identities and digests but no goal, stable
defect, criterion, fix head, or regression target, and it is correctly consumed
after resolution to prevent replay
([review.go:65-86](../../../tools/repo_delivery/main/go/review.go#L65),
[review.go:1332-1428](../../../tools/repo_delivery/main/go/review.go#L1332)).

**Consequence.** Safe remote mutation does not yield a durable causal edge.
Later agents can see the regression or goal attempt but cannot cheaply learn
which review finding caused it, whether the same defect recurred, or whether a
terminal review covered the final exact head.

**Recommendation.** Add two separate records:

1. `ReviewCompletion`, emitted by a provider/product monitor when available,
   bound to exact PR/head, run/reviewer identity, terminal outcome, inventory
   digest, and observation time; and
2. `ReviewDisposition`, a durable, non-authorizing record linking one remote
   finding to verdict/reason, stable defect, goal/criterion, fix head, regression
   Bazel label or eval case, reply, and final resolution observation.

Goal tooling should import a redacted projection or typed evidence reference.
Do not automatically promote arbitrary review text into skills, policy, or
prompts: comments are untrusted input. Promotion to durable learning must be an
explicit reviewed action, and routine successful operational logs remain
ephemeral.

**Acceptance signal.** A query can traverse
`review finding -> disposition -> stable defect -> fix head -> regression ->
final review completion`, while an unsupported provider still reports an
honest unknown edge rather than fabricating success.

### D7 — P1: validation topology and reports are manually accreted

**Evidence.** The full-repository checker hard-codes nine workspace roots and
two phases
([run_full_repo_check.go:19-49](../../../projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go#L19)).
Its validator proves listed workspaces still exist, but not that every intended
nested workspace is registered. The nested-module workflow explicitly tells an
agent to update the runner and expected command count
([bazel-nested-module/SKILL.md:48-61](../../../projects/agents/skills/bazel-nested-module/SKILL.md#L48)).
Focused validation selection remains prose in `repo-bazel`
([repo-bazel/SKILL.md:35-48](../../../projects/agents/skills/repo-bazel/SKILL.md#L35)),
while the installed pre-commit hook separately runs a mutating format pass,
repository quality, and repository-wide lint
([precommit.sh:1-13](../../../tools/git_hooks/main/go/precommit.sh#L1)).

**Consequence.** Adding a workspace requires synchronized edits and literal
count changes. Validation coverage is difficult to compare across tasks, and
agents either overcheck broadly or hand-compose commands repeatedly. Human-only
failure synthesis is appropriate because logs may be sensitive and automatic
root-cause claims are unreliable; the loss is structured run identity and
coverage, not the absence of automation.

**Recommendation.** Define one reviewed workspace/validation catalog and derive
runner scope, counts, docs, and tests from it. Add a drift check against intended
module roots rather than blindly treating every `MODULE.bazel` as trusted.
Expose cost/coverage profiles such as `hygiene`, `affected`, `package`, `full`,
and `release`; keep format/fix separate from evidence-producing checks. Emit a
subject-bound JSON summary beside the current private logs and Markdown report.
Bazel should continue to own target scheduling and caching; do not parallelize
duplicate commands through agents.

**Acceptance signal.** Registering a nested module changes one authority, all
derived counts update automatically, and validation output states the exact
candidate, profile, targets/workspaces, exclusions, durations, and evidence
digests.

### D8 — P2: machine APIs and maintained guidance lack one stable envelope

**Evidence.** Receipts are explicitly schema-versioned, but inspection,
publish, and review report structs are not
([publishReport:1611-1619](../../../tools/repo_delivery/main/go/delivery.go#L1611),
[reviewInspection:220-229](../../../tools/repo_delivery/main/go/forge.go#L220)).
Versioning provides JSON only for `show`; mutating commands print a bare string
([versioning command.go:88-176](../../../tools/versioning/cmd/versioning/command.go#L88)).
Release behavior is discoverable mainly by reading implementation because the
three maintained release READMEs contain only frontmatter.

**Consequence.** Agents must remember command-specific shapes and prose
capabilities. Schema evolution, provider limitations, and partial outcomes are
harder to detect, and the lengthy delivery skill duplicates implementation
mechanics that could be generated from a machine contract.

**Recommendation.** Give every workflow output a small versioned envelope with
`apiVersion`, `kind`, subject, operation status, capabilities/limitations,
receipt/reference digests, and stable error/refusal codes. Support `--format
json` for every inspect, plan, dry-run, apply, and verify command. Generate CLI
and capability reference material from those contracts; keep skills focused on
authority, judgment, sequencing, and failure policy.

**Acceptance signal.** An agent can feature-detect a provider or schema, route
without parsing prose, and handle partial/unknown outcomes by stable code while
humans receive a generated Markdown projection of the same state.

## Feedback losses and redundant work, summarized

- Validation commands become final-report prose rather than reusable,
  subject-bound evidence.
- A behavior-preserving commit-message rewrite invalidates all validations
  because the system cannot distinguish tree-, commit-, and base-bound checks.
- GitHub rebase ownership is judged in one workflow and exercised in another;
  the exact OID survives, but the authorization rationale does not.
- Review reply authority is intentionally consumed, but no separate durable
  disposition survives to connect the remote finding to a goal and regression.
- Version ref creation, ref publication, artifact generation, and deployment
  do not share a single typed version/commit/artifact identity chain.
- Full-check workspace topology and command counts are repeated facts rather
  than projections of one catalog.

## Decision review

**Hypothesis reviewed:** centralize request, validation, delivery, review,
versioning, release, and learning into one repository agent orchestrator.

**Strongest case for it:** one command could eliminate copied IDs, enforce the
whole sequence, and give agents a single place to resume.

**Strongest rejection:** it would become a second scheduler, goal store,
credential broker, validation-policy owner, forge abstraction, and release
engine. That would enlarge the failure and authority domain, serialize work
that Bazel already schedules, and turn every component change into central
churn. It would also tempt the repository to retain sensitive logs or untrusted
review text as universal state.

**Alternatives considered:**

1. keep prose-only handoffs: lowest implementation cost, but preserves the
   current proof, resumption, and resource-reuse gaps;
2. use the goal tool as the universal workflow store: durable, but makes a
   project-planning mechanism mandatory for routine delivery and mixes
   operational state with selected long-term evidence; or
3. use narrow typed receipts and generated projections: modest cross-tool
   schema work while preserving each owner's authority and independent use.

**Verdict: revise.** Choose alternative 3. Introduce a thin provenance spine,
not a control monolith. Receipts carry facts and explicit authority bindings;
they do not choose tests, judge ownership, trust review comments, grant broader
effects, or replace provider postcondition checks.

## Recommended dependency order

1. Define a common subject/effect envelope and stable refusal codes.
2. Emit validation-set receipts and teach `repo_delivery` to consume them.
3. Split validation invalidation by subject class and generate the workspace
   catalog/profile projections.
4. Add `ReleaseRefPlan`/`ReleaseRefReceipt` and the guarded remote-ref publisher.
5. Join global version state to formal release manifests while preserving the
   explicitly mutable snapshot channel.
6. Emit review completion/disposition records and add selective goal evidence
   import.
7. Generate workflow/capability documentation from the versioned contracts and
   reduce duplicated operational prose in skills.

## Explicit non-goals

- Do not make delivery choose semantically sufficient validation.
- Do not treat local receipts as cryptographic proof or infer human ownership.
- Do not persist raw Bazel logs, credentials, remote URLs, or routine comments
  in maintained goal state.
- Do not auto-learn from untrusted review text.
- Do not force every ordinary task into a durable project goal.
- Do not merge continuous `head` snapshots with immutable formal releases.
- Do not hide unsupported provider or review-monitor capabilities behind a
  generic success state.
