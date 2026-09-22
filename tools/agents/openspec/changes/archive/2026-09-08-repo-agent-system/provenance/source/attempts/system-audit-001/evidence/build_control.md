# Build-control plane audit

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
- Scope: the repository-wide Bazel/build/query/format/lint/test control
  plane, not the behavior of individual product subprojects.
- Method: read-only source and history inspection. No Bazel build, test, or
  broad query was run. This report is the only file written by this audit.

## Executive finding

The repository has a strong execution substrate but not yet a coherent
agent control plane.

Bazel accurately models declared build dependencies. `bazel_agent` gives
agents a consistent low-level executor. `repo_quality` supplies pinned tools,
and `full-repo-check` supplies a deliberately exhaustive build/test backend.
What is missing is the layer that tells an agent, in machine-readable form:

- what an operation means;
- what it can read or change;
- what authority it requires;
- how much it is expected to cost;
- which checks a change requires and why;
- what was actually covered; and
- which exact candidate a result proves.

Today those answers are distributed among target names, macros, rc flags,
skills, `AGENTS.md`, generated-file comments, and agent judgment. Bazel query
can enumerate labels and dependencies, but it cannot distinguish a safe read
from a source, host, or remote mutation. The same undifferentiated agent
profile serves a cheap inner loop, final correctness evidence, diagnosis, and
a full-repository audit. This is the highest-leverage build-control defect.

The preferred system is a linked abstraction tower:

```text
repository invariants and execution context
    -> locally owned workspace/operation/check declarations
    -> generated validated catalog
    -> impact, authority, and cost planner
    -> thin Bazel executor
    -> structured candidate-bound evidence
    -> goal and delivery gates
    -> reviewed regression and policy accretion
```

The catalog should be a generated projection over component-owned facts, not
a new hand-maintained mega-manifest. Bazel remains the graph and execution
engine. `bazel_agent` remains the transparent leaf executor. Planning,
authorization, evidence, and learning belong above it.

## Current control plane

| Layer | Current authority | What works | System gap |
| --- | --- | --- | --- |
| Policy | `AGENTS.md` and Bazel-related skills | Strong narrow-first, hermeticity, safety, and evidence guidance | Prose is not projected into executable operation/check contracts |
| Executor | `projects/bazel_agent` | Thin pass-through, `syscall.Exec`, nearest-workspace lookup, pinned config selection | No operation semantics, provenance report, task context, or stale-install detection |
| Configuration | `.bazelrc` -> `tools/bazelrc/root.bazelrc` -> preset/project rc | Central flags, lockfile error, bounded build/test output, skill validation | One profile conflates iteration, evidence, and audit; local/absolute state leaks into actions |
| Discovery | `bazel_agent query` and root labels | Native dependency and target discovery | No effect, authority, owner, cost, generated-artifact, or check metadata |
| Quality | `tools/repo_quality`, root aliases, `--config=lint` | Broad pinned formatter coverage and semantic lint aspects | Names and entry points disagree; no one coherent check taxonomy or changed-file planner |
| Generation | Gazelle and heterogeneous update/write targets | Generators are explicit and generated files often name an updater | No complete catalog, uniform naming/effects, freshness gate, or nested-workspace orchestration |
| Exhaustive assurance | `full-repo-check` | All current workspaces, continue-after-failure, private logs, honest basic caveats | Hard-coded topology, build/test only, no plan/resume/JSON/candidate identity/coverage universe |
| Delivery evidence | Agent-selected commands and prose reports | Exact-candidate delivery policy exists | Validation selection and result normalization remain manual |

## Preserved strengths

Any redesign should preserve these properties:

1. `.bazeliskrc:1-2` pins a Bazel version and archive digest.
2. `projects/bazel_agent/cmd/bazel_agent/main.go:19-27,83-90` keeps the
   executor small, preserves Bazel argument precedence, and replaces itself
   with the child process.
3. `tools/bazelrc/project.bazelrc:38-60` makes agent output non-interactive,
   timestamps progress, bounds build/test output, and rejects lockfile drift.
4. `AGENTS.md:152-154` and
   `projects/agents/skills/full-repo-check/SKILL.md:48-50` correctly avoid raw
   Build Event Protocol artifacts because they can contain the client
   environment.
5. `tools/repo_quality/README.md:13-19` and
   `tools/repo_quality/BUILD.bazel:186-269` use pinned Bazel-acquired tools and
   provide broad whole-index formatting coverage.
6. `projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go:88-108`
   restricts report/log directories, and lines `192-215` continue through the
   matrix instead of hiding later failures.
7. `AGENTS.md:78-79,170-179,219-234` already teaches narrow discovery and
   verification before repository-wide work.

## Ranked system proposals

### P0.1: Establish a typed operation and target-effect contract

**Classification:** central control-plane architecture, with immediate leaf
safety corrections.

**Evidence**

- Root `//:tf` is declared at `BUILD.bazel:26-29`. The Terraform macro maps
  its empty/default label to `apply` and also emits `migrate`, `import`,
  `destroy`, `deploy`, `state`, and `force_unlock` at
  `tools/terraform/defs.bzl:28-45`.
- This is systemic rather than Terraform-specific. The default `dns` target
  runs `dnscontrol push` at `infra/dns/BUILD.bazel:29-45`, and default Ansible
  targets execute deployment playbooks, for example
  `infra/pve/ansible/BUILD.bazel:16-37`.
- Root `BUILD.bazel:46-90,170-193` also exposes host/source writers such as Git
  hook installation, formatting, skill-link reconciliation, Gazelle, and
  compile-command refresh alongside ordinary checks.
- `projects/bazel_agent/README.md:24-30` deliberately states that the runner
  performs no validation and allows later arguments to override its config.
- `AGENTS.md:193-200` supplies important prose safety rules, but no Bazel
  metadata lets a planner enforce or even discover those rules.

**Proposal**

Define a versioned repository operation descriptor, locally declared by the
owning macro or target and validated centrally. At minimum it should contain:

- stable operation ID and owner;
- effect: `read`, `source-write`, `host-write`, or `remote-write`;
- required capabilities such as `network`, `secrets`, `credentials`, or
  privileged host access;
- scope: target, package, workspace, repository, or external system;
- expected cost class and typical latency;
- determinism/cache class;
- safe commands (`build`, `test`, `run`, or update only);
- prerequisite and validation operation IDs; and
- generated inputs/outputs where applicable.

Use reserved, queryable Bazel tags/providers where Bazel owns the fact, and
generate one compact operation catalog for consumers that cannot inspect the
configured graph. A validation target should reject unclassified public or
agent-runnable operations. The planner must require explicit matching
authority before a write effect is executed.

Immediate corrections should remove implicit write defaults: generic labels
such as `:tf`, `:dns`, and deployment playbook labels should resolve to help,
status, preview, or fail-closed guidance; writes should have explicit labels
such as `:tf.apply` or `:dns.deploy` and still carry effect metadata.

**Acceptance signals**

- An agent can ask for JSON metadata for any root or public runnable label.
- A repository test fails on an unclassified runnable operation.
- No generic/default root operation performs a source, host, or remote write.
- A remote-write plan is rejected without explicit authority for the exact
  operation and scope.

### P0.2: Separate task scratch from Bazel action temporaries

**Classification:** immediate foundational correction; prerequisite for
reliable concurrent control-plane execution.

**Evidence**

- Repository policy requires task-specific `out/<task>/` storage at
  `AGENTS.md:22-34`.
- The mandatory runner instead hard-codes shared `out/tmp` at
  `projects/bazel_agent/cmd/bazel_agent/main.go:17,77-81`; its README and tests
  make that behavior contractual at `projects/bazel_agent/README.md:34-38`
  and `projects/bazel_agent/cmd/bazel_agent/main_test.go:147-159`.
- `tools/bazelrc/project.bazelrc:42-59` injects those inherited absolute
  `TEMP`, `TMP`, and `TMPDIR` values into repository rules, normal actions,
  host actions, and tests.
- `tools/repo_quality/tool_adapter.sh:5-16` independently shares one
  `out/repo_quality` cache/config/tmp tree across invocations.

**Consequence**

The declared task isolation cannot be satisfied through the required runner.
Concurrent agents share mutable temporary state. Absolute worktree paths enter
action and repository-rule environments, making action identities
worktree-specific and creating paths that may not exist in remote or sandboxed
execution. This spends cache capacity while weakening hermeticity.

**Proposal**

- Introduce a canonical task/run execution context with private
  `out/<task>/<run>/` paths and explicit ownership/retention.
- Use that context for controller processes, repository-updating tools, and
  private logs.
- Do not globally export host absolute temporary paths to Bazel actions or
  tests. Normal actions should use declared outputs and Bazel-managed temp;
  tests should use their test temp contract.
- Opt host scratch into only the small set of intentional no-sandbox,
  workspace-inspecting developer operations.
- Add concurrency, permissions, sandbox, and cross-worktree cache-portability
  tests.

### P0.3: Add an impact/cost planner and candidate-bound evidence protocol

**Classification:** central roadmap architecture.

**Evidence**

- `AGENTS.md:219-234` tells each agent to choose checks from changed files, but
  provides no machine planner.
- `projects/agents/skills/repo-delivery/SKILL.md:68-69` explicitly leaves
  validation selection to agent judgment.
- Whole-index quality discovery is intentionally based on Git-tracked files
  (`tools/repo_quality/README.md:13-15`); custom adapters use `git ls-files` at
  `go_module_quality.sh:61-64`, `lua_quality.sh:51-54`, and
  `qt_xml_quality.sh:47-50`. Untracked new task files are outside that stated
  guarantee during the inner loop.
- Full-check results record only workspace, phase, command, exit, duration,
  and log path at
  `projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go:42-49`.
  The Markdown report is written only after all commands complete at lines
  `225-327,351-358`; it records no candidate Git/tree identity, plan digest,
  target universe, or configuration identity.

**Proposal**

Provide one controller interface with machine and human views:

- `discover PATH|LABEL`: owner, workspace, operations, effects, and checks;
- `plan`: changed paths, declared task-owned new paths, affected targets and
  conservative reverse dependencies, selected profile, reason per check,
  estimated cost, authority needs, and escalation conditions;
- `run`: execute an immutable plan through the thin Bazel executor;
- `explain`: why a check was included or omitted; and
- `resume`: continue only input/config-identical completed work.

Plans and results should be schema-versioned and digest-bound. Evidence should
identify the exact Git commit/tree plus relevant dirty/task-owned inputs,
workspace/config/tool identities, start/end times, target counts, coverage
gaps, exit taxonomy, and private-log digests. Emit append-only events after
each operation and generate Markdown from the JSON record. Do not restore raw
BEP or environment dumps.

The planner must be conservative: root module/rc/toolchain changes, unknown
files, cross-workspace dependencies, or incomplete metadata escalate scope.
It should never infer permission to perform writes merely because an operation
is present in the plan.

**Acceptance signals**

- The same immutable inputs produce the same plan digest.
- Every selected check has a reason and every omitted in-scope check has an
  explainable policy result.
- Task-owned new files are included without inspecting or claiming unrelated
  untracked files.
- A result can be proven to apply to one exact candidate and validation
  profile without reading prose logs.

### P1.1: Define one validation taxonomy with staged profiles

**Classification:** central roadmap, followed by component migration.

**Evidence**

- Root `//:lint` aliases `//tools/repo_quality:check` at
  `BUILD.bazel:56-59`. That target is an alias for formatter checks at
  `tools/repo_quality/BUILD.bazel:285-322`.
- Semantic Buildifier, Ruff, and ShellCheck linting is instead an aspect
  activated by `build --config=lint` at
  `tools/bazelrc/root.bazelrc:9-12` and documented separately at
  `tools/repo_quality/README.md:27-37`.
- The Git hook performs a mutating whole-repository format, whole-index quality
  test, and root-wide lint in sequence at
  `tools/git_hooks/main/go/precommit.sh:11-13`, despite narrow-first guidance
  at `AGENTS.md:78-79,221-234` and the no-broad-churn invariant at line `217`.
- `full-repo-check` has only `build` and `test` phases
  (`run_full_repo_check.go:40`) although its skill advertises release-readiness
  use (`projects/agents/skills/full-repo-check/SKILL.md:3-8`). It does not
  explicitly run the semantic lint config or generator idempotence checks.
- Every agent test is forced fresh, and all commands keep going, at
  `tools/bazelrc/project.bazelrc:55-56,67-68`; every runner call is batch mode
  at `projects/bazel_agent/cmd/bazel_agent/main.go:19-26`.

**Proposal**

Define stable phases such as `hygiene`, `format-check`, `lint`, `generated`,
`build`, `test`, `security`, and `release`, then compose them into explicit
profiles:

- `changed/fast`: cached, fail-fast, quiet, narrow iteration;
- `workspace`: complete owning-workspace confidence;
- `fresh/evidence`: deliberate fresh tests against an exact candidate;
- `full/audit`: keep-going repository matrix and complete diagnostics; and
- `diagnose`: progressively expanded output for known failures.

Formatting checks should be non-mutating by default. Source mutation should
be a separate explicit operation over task-owned paths. Rename the current
root surfaces so `lint` means semantic lint and `format-check` means formatter
drift. A release-readiness claim must enumerate all required phases or narrow
its claim.

Do not simply remove `--batch` or freshness checks. Batch mode provides useful
process isolation, and fresh tests provide stronger final evidence. First
separate intent, then benchmark batch versus a managed server under realistic
concurrent-agent workloads. The most certain waste is applying audit semantics
to every iterative command.

### P1.2: Generate one workspace and generated-artifact inventory

**Classification:** central roadmap projection; local declarations remain
authoritative.

**Evidence**

- The full-check runner hard-codes root plus eight nested workspaces at
  `run_full_repo_check.go:19-38`. It only checks that listed workspaces exist at
  lines `67-85`; an unlisted new workspace is invisible.
- The literal topology is repeated as "eight"/"eighteen" in
  `projects/agents/skills/full-repo-check/SKILL.md:32-33,77-78`, its eval cases,
  and a unit test expecting 18 commands at
  `run_full_repo_check_test.go:72-76`.
- `projects/agents/skills/bazel-nested-module/SKILL.md:48-65` requires agents to
  update `.bazelignore`, parent dependencies, docs, project exclusions,
  full-check registration/counts, and Gazelle registration by hand.
- Root Gazelle is one mutating `fix` target at `BUILD.bazel:170-186`; root
  Gazelle ignores nested workspaces, while the nested-module/Gazelle skills
  require separate execution without providing one repository-wide generator
  plan.
- Update operations are heterogeneous: `.update` from
  `tools/helm/main/bzl/al_helm_chart_lock.bzl:30-33`, `-write` from
  `tools/bzl/main/bzl/al_genquery_write_to_source_file.bzl:73-76`, `write_*`
  in `projects/goal/docs/BUILD.bazel:11-26`, and `:gazelle` at the root.
- Generated headers instruct direct `bazel run`, for example
  `gazelle_python.yaml:1-4`, `requirements.txt:1-5`, and
  `tools/bazelrc/preset.bazelrc:1-3`, while agent policy requires
  `bazel_agent` at `AGENTS.md:145-151`.

**Proposal**

Discover tracked workspace roots mechanically and combine them with small,
locally owned declarations for non-derivable facts: owner, supported profiles,
optional configs, generators, and platform constraints. Validate the result
bidirectionally against tracked `MODULE.bazel` files, `.bazelignore`, parent
module overrides, docs aggregation, and full-check coverage. Generate counts,
docs, runner plans, and tests from the inventory rather than copying literals.

Give every maintained generated artifact a validated relation:

```text
source inputs -> generator operation -> generated outputs -> freshness check
```

Standardize update/check metadata even if legacy target names remain. Add a
non-mutating generated-freshness profile and require generator idempotence.
Generated comments should use a neutral canonical operation ID or explicitly
show both human and agent invocation routes.

**Acceptance signals**

- Adding a nested module requires one local declaration plus generated
  projections, not synchronized literal counts.
- Missing and extra workspace registrations both fail a focused meta-test.
- Every generated maintained file has one discoverable updater and freshness
  check.
- A second generator run is empty and recorded as evidence.

### P1.3: Upgrade `full-repo-check` into an exhaustive backend, not the planner

**Classification:** foundational component work under the central planner.

Keep its valuable continue-after-failure and private-log behavior, but add:

- `--plan`, workspace/phase selectors, failed-only and resume modes;
- per-operation timeouts and safe signal propagation;
- schema-versioned incremental JSON/events plus generated Markdown;
- candidate, plan, config, workspace-registry, and log digests;
- planned and observed target counts;
- explicit manual, incompatible, optional-config, skipped, and unsupported
  sets; and
- resource-budgeted parallelism across independent workspaces, with a
  conservative sequential default.

The current static report caveat at `run_full_repo_check.go:238-249` is honest
but not measurable. Empty or unexpectedly reduced target universes should be
coverage failures. Interruption should leave a valid partial record. Repeated
fatal startup failures should not be paid twice for build and test when the
second phase cannot add evidence.

### P1.4: Make the repository root a tiny safe control surface

**Classification:** architecture-aligned component migration.

`BUILD.bazel:1-15` eagerly loads npm, pip, Gazelle, compile-command, docs,
Python, skill, Buildifier, AL, Terraform, TOML/TXT, and Vault macros before a
lightweight root alias can be parsed. Lines `26-44,102-139,170-193` mix
privileged runners, dependency generation, documentation, quality, Gazelle,
and compile commands in one package.

Keep root BUILD and README entry points deliberately small: native aliases or
suites for `help`, `discover`, `plan`, `check`, and `doctor`, pointing to
implementations owned under the repository tooling/control plane. Heavyweight
and mutating operations should remain with their owners and appear through the
generated catalog. This reduces unrelated load/fetch failure coupling and
makes the root an accurate zero-context interface rather than a warehouse of
operations.

### P2.1: Make runner, platform, and rc provenance inspectable

**Classification:** foundational hardening.

- Root policy claims Bazelisk-managed Bazel at `AGENTS.md:147-151`, but the
  runner accepts any `bazel` on `PATH` at
  `projects/bazel_agent/cmd/bazel_agent/main.go:83-89`.
- The installed binary can silently become stale; manual reinstall is required
  after source changes at `projects/bazel_agent/README.md:48-56`.
- Repository Bazelisk packaging is Linux x86_64 only at
  `third_party/com_github_bazelbuild_bazelisk/include.MODULE.bazel:5-22` and
  `binary_toolchain.json:9-35`; the runner uses POSIX `syscall.Exec`, while
  `tools/bazelrc/project.bazelrc:1-2` hard-codes a Fedora/OpenJDK truststore.
- `tools/bazelrc/README.md:6-10` still documents direct raw Bazel invocation.

Add a cheap, structured `doctor/describe` contract reporting runner source and
contract versions, installed/source drift, workspace root, resolved Bazel and
Bazelisk identities, expected Bazel pin/digest, effective profile, scratch
context, rc provenance, and supported platform. Either declare and fail early
on the managed Linux/x86_64 platform boundary or provide real multi-platform
pins/configs. Generate or validate root and nested rc composition from one
typed profile definition.

Logical ownership of the repo-private runner should be decided with the
repository topology architecture. The evidence favors `tools/agents` because
`tools/README.md:10-18` owns repo-wide private tooling and
`projects/BUILD.bazel:24-50` excludes `bazel_agent` from project releases, but
the move itself is lower leverage than establishing the contracts above and
should preserve compatibility labels.

### P2.2: Govern quality exclusions and close the learning loop

**Classification:** later roadmap phase.

The formatter exclusions at `.gitattributes:106-175` contain useful prose
rationales, but most do not name a machine-readable owner or semantic
validator. Create a coverage projection mapping relevant file classes to one
formatter policy plus owning semantic checks, or to an explicit exclusion with
owner and reason. This gives an agent a cheap exact answer to "what checks this
file?"

Preserve sanitized duration, coverage, and failure fingerprints from structured
receipts. Repeated cost hotspots or failures should enter a reviewed promotion
path:

```text
receipt evidence -> minimized recurring cause -> regression/check/runbook
    -> owner review -> catalog/profile update
```

Raw logs stay private and ignored. Only stable, non-sensitive lessons become
maintained policy or tests.

## Decision review

**Outcome sought:** maximize agent accuracy, safe control, and feedback quality
while minimizing compute, elapsed time, context, and manual navigation.

**Strongest case for the current design:** Bazel already supplies a rigorous
dependency graph and cache; a thin pass-through avoids a second command
language; prose keeps nuanced policy readable; full checks intentionally favor
freshness and failure coverage. A central manifest or clever changed-target
algorithm could drift, hide Bazel semantics, or create false confidence.

**Strongest rejection:** the missing facts are not derivable from the Bazel
dependency graph. `:tf` invoking apply, `:lint` meaning formatter check, shared
absolute scratch entering actions, and a release-readiness tool omitting an
explicit lint phase are direct counterexamples. Literal workspace counts and
manual validation selection demonstrate that the current system is not
accretive. More prose cannot make these facts machine-enforceable.

**Alternatives considered**

1. **Keep the current distributed conventions.** Lowest implementation cost,
   but retains safety ambiguity, repeated navigation, and avoidable full-scope
   work.
2. **Turn `bazel_agent` into a validating all-in-one wrapper.** Centralizes
   behavior, but compromises its transparent leaf role, adds recursive-query
   and bootstrap problems, and risks diverging from Bazel.
3. **Add a generated typed control plane above the thin runner.** Preserves
   local ownership and Bazel semantics while making effects, costs, plans, and
   evidence explicit and testable.

**Verdict: revise.** Choose alternative 3. Keep `bazel_agent` thin; do not
immediately remove batch mode or reintroduce BEP. First repair scratch and
unsafe default labels, establish local typed declarations and generated
catalogs, then add planning/evidence profiles. Change process lifetime only
after concurrent-agent measurements show the benefit and isolation conditions.

## Recommended central roadmap order

| Phase | Central system work | Foundational component work | Exit signal |
| --- | --- | --- | --- |
| 0: Correct contradictions | Define task/run context and initial effect vocabulary | Stop global host temp injection; isolate quality scratch; make default write labels explicit/safe; correct `lint` naming | Existing policy is executable without contradiction |
| 1: Describe | Define locally owned operation, workspace, check, and generated-artifact schemas; generate a catalog | Add metadata to root/tool macros; add runner `doctor`; validate platform/rc/pins | Zero-context JSON discovery answers owner/effect/cost/check questions |
| 2: Plan | Implement deterministic changed/workspace/full/release planning with reasons, budgets, and authority gates | Make format non-mutating by default; add generated freshness; split agent rc profiles | Narrow plans are conservative, reproducible, and cheaper than full checks |
| 3: Prove | Define immutable plan/result/event schemas and delivery/goal evidence linkage | Upgrade full-check with JSON, coverage, selectors, interruption safety, and resume | Every claimed result identifies exact candidate and coverage |
| 4: Accrete | Define reviewed evidence-to-regression promotion and profile cost governance | Add exclusion ownership, duration/failure fingerprints, and contract regression tests | Repeated friction decreases without growing eager agent context |

## Explicit non-goals

- Replacing Bazel as the dependency and execution engine.
- Building a hand-maintained central catalog that duplicates component facts.
- Parsing or retaining raw secret-bearing BEP/environment output.
- Treating a changed-target heuristic as proof when metadata is incomplete.
- Automatically authorizing writes because a target is discoverable.
- Running every optional configuration or manual target in the cheap inner
  loop.
- Moving directories before the repository-wide ownership/topology decision is
  settled.

## Audit limitations

This was a bounded read-only design audit. No command latency, cache-hit rate,
remote-execution behavior, or concurrent-runner contention was measured. Those
measurements are acceptance work for the proposed profile/process-lifetime
changes, not evidence for changing `--batch` now. The current hard-coded
full-check list was compared with tracked `MODULE.bazel` paths and matches the
root plus eight nested workspaces at the audited state; the defect is the lack
of a completeness invariant, not a presently missing listed workspace.
