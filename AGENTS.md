---
title: Agents
---

## Start with the task

This is the repository-wide default; a nearer `AGENTS.md` takes precedence
within its subtree. Establish the requested outcome and existing authority,
resolve the affected path and Bazel workspace, read the nearest `README.md`,
`BUILD.bazel`, and `include.MODULE.bazel` when present, and inspect before
mutating. Load only the skills needed for the current phase.

Use [README.md](README.md) for the repository map and
[projects/agents/README.md](projects/agents/README.md) for the agent-system
documents. Each top-level tree's README owns its target visibility, allowed
build consumers, and publication boundaries; read it before changing those
relationships. Roadmaps describe intent, not current behavior or authority.

## Authority and scope

Keep one owner for each fact: the user request owns outcome and authority;
`AGENTS.md` owns agent policy; the nearest `README.md` owns component purpose
and boundaries; `CODEOWNERS` owns review accountability; BUILD and MODULE
files own executable and dependency structure; canonical skills own
procedures; goal records own durable work state; runtime providers own their
observed capabilities; Git and delivery receipts own candidate and publication
state. Derived views must identify their sources, version or digest,
observation time, unavailable fields, and truncation. Resolve conflicting
facts at their owning source; architecture does not override those owners.

- Use `answer-question` for substantive questions, including mixed requests.
  A question alone authorizes investigation, not mutation; preserve action
  authority already granted in the conversation. Quoted transformation
  payloads do not trigger this skill.
- Use `decision-review` for material design, security, operational, costly,
  irreversible, or repeatedly failing choices. The primary agent owns the
  verdict and material trade-offs. Routine reversible choices need no review.
- Reviews and audits are evidence, not authority or new acceptance criteria.
  Fix in-scope blockers and tiny defects directly related to the requested
  change; report other incidental findings separately and continue. Delegated
  workers report findings outside their assigned scope to the coordinator.
  Do not conceal correctness or safety failures behind a workaround. If a
  blocker needs material unauthorized scope expansion, ask; otherwise choose
  the smallest supported in-scope fix. Narrow unsupported claims instead of
  adding unrelated machinery.
- Delegate useful independent work when it can proceed safely in parallel.
  Keep one coordinator for shared state and bounded worker ownership. Do not
  fragment trivial tasks or invent work for unused agents. Reassess when
  dependencies or workstreams change; explain a constraint only when it
  materially affects the task.
- Never run state-changing infrastructure operations without the user's
  explicit request for the exact operation and scope. Ordinary implementation
  validation must not contact or mutate live systems.
- Treat an automatic approval rejection as a strategy signal. Diagnose it,
  choose a materially safer authorized approach, or ask; never route around
  it. Never immediately retry a rejected, failed, or rate-limited escalated
  operation. Allow at most one retry after a real delay and causal change.

## Isolate changes and scratch

Use a dedicated feature branch in its own linked worktree for every task that
can modify repository files or task-owned scratch. Verify both before the
first mutation. The default branch and checkout are read-only unless the user
explicitly authorizes that exact task there; keep them clean.

Keep every task-owned download, report, log, cache, extracted archive, and
temporary file under ignored `out/<task>/` in the applicable workspace. Point
configurable temporary and cache directories there too. Use operating-system
temporary storage only when the tool cannot use workspace scratch, and remove
only your residue before handoff. Never delete unrelated temporary files.

Secret-bearing temporary material must be task-private, access-restricted,
short-lived, and explicitly cleaned up. Never track, stage, commit, or promote
secrets as durable evidence. "Outside the repository" means outside tracked
or committable source; ignored task scratch is suitable.

## Load procedures when needed

- Load `project-layout` before creating or moving source directories or
  deciding a repository layout. Use its role-based layout for new paths;
  legacy `main/` layouts are not precedent. Project names use only ASCII
  letters, digits, and underscores.
- Load `bazel-agent` and `repo-bazel` before Bazel work. They own invocation,
  bootstrap, BUILD/Starlark/MODULE mechanics, generated files, sandboxing, and
  validation. Do not call `bazel` directly or expose raw BEP output. Run the
  owning update command instead of hand-editing generated files.
- Load `repo-external-dependency` when acquiring or updating external inputs.
  Tools used by Bazel must be pinned and checksummed, with no undeclared host
  tools or lifecycle downloads. Environment-dependent tool output is allowed
  within the requested workflow's authority and sandbox/network policy.
- Load `ast-grep` for syntax-shape searches or rewrites that plain text cannot
  express; use its repository-pinned Bazel entry point.
- Load `repo-terraform`, `repo-ansible`, and `repo-secrets` for their respective
  work. Preserve the owning `al.lua` and Bazel packaging/injection flow.
  Never commit `.terraform/`, state, plans, environment files, or local
  credentials.
- Load `goal` when work needs durable continuation or retries. Use the
  smallest record that preserves acceptance and the next action; do not
  attach a token budget unless the user explicitly requests one.
- Load `repo-delivery` when preparing publication or final handoff, before
  staging or committing delivery changes. It owns validation gates, commits,
  pushes, pull requests, reviews, and receipts. Commit and push all legitimate
  task-owned source, documentation, and configuration at delivery unless the
  user says otherwise. A conversational pause is not itself a delivery gate.
  Scratch, secrets, temporary files, and generated artifacts are excluded;
  required non-temporary binaries may be committed only through Git LFS.
- Before rewriting branch history, resolving conflicts, or recovering a
  publication, load the delivery recovery guidance and `git-rebase-remote`
  as applicable. Preserve local and remote task progress; never rewrite
  shared, human-owned, unrelated, or ambiguous history. Use exact remote
  leases. Resolve authorized minor conflicts only from unambiguous task and
  three-way evidence, stage explicit resolutions, and rerun invalidated
  checks. Abort and ask when resolution would guess intent, discard work, or
  alter unauthorized semantics.

## Use bounded tools and evidence

Prefer supported live Cordis handlers for context, reads, searches, and Git
inspection; discover the catalog before assuming one is unavailable. Otherwise
prefer purpose-built tools, MCP capabilities, and repository Bazel targets.
Use a host shell only when no suitable entry point exists. Do not build new
tooling merely to avoid a one-off command.

Use `rg`, `rg --files`, bounded `find`, or Bazel queries; never recursive
`grep`/`ls` or filesystem-root searches. Wrap shell commands in an appropriate
`timeout`, except intentionally long-running or interactive operations.
Prefer Go and Bazel-native rules for repository automation. Do not introduce
shell scripts or shell-based rules unless the task explicitly requires them.
Do not use `genrule`. Expose Go automation as Bazel `go_binary` targets; use
another language only when Go or Bazel would materially complicate the task.
Validate Go with the pinned Bazel toolchain, not host `go`.

Reuse evidence while relevant inputs are unchanged. Empty results and no-ops
are evidence with bounded coverage, not proof of absence or progress. After
two inconclusive searches, revise the approach using what is known. Before
retrying failure, identify its cause, a real change, and the predicted effect.
After two unchanged edits, diagnose the source/cache dependency; invalidate
stale generated output once deterministically and inspect the effect.

Bounded polling, mutation postconditions, and advanced refs justify fresh
observations. Record relevant revisions and observation times. After
interruption, resume from the recorded next action and recheck only mutable
dependencies, including authority, branch, candidate, and in-flight operations.
Preserve failure evidence and a stable name for recurring defects.

Start verification with `git diff --check`, then the narrowest useful package
checks under `repo-bazel`; run configured formatters and `//:buildifier_test`
for BUILD/Starlark changes. Before completion, verify source and a
representative output against the exact candidate. Command success alone is
not acceptance. Do not overwrite or stage unrelated formatter changes; record
baseline failures without waiving the delivery quality gate. Git hook
installation is optional; required checks still apply.

## Public source and protected data

Checked-in source, documentation, and fixtures are public. When the task calls
for an external service, public source and eval output derived solely from
public fixtures in an isolated public workspace may be sent without a separate
confidentiality approval. This does not authorize exposing credentials,
personal information, or secret-bearing data.

Inspect local, generated, untracked, live, and infrastructure artifacts before
disclosure because they can contain protected data. Their origin alone does
not make them confidential. Ordinary build/test/lint diagnostics and non-secret,
non-personal operational facts may be reported, including for `infra/`.

Keep source disclosure, target visibility, build consumers, artifact and
documentation publication, and secrets/personal information distinct.
"Private" tree policy means repository-internal use unless it identifies
protected content; it does not make ordinary source or operational facts
confidential.

## Communicate and maintain source

When the main agent loads a skill, say which one and why it applies. Explain
the purpose and relevant side effects of non-obvious tools or delegations;
group routine checks and report opaque tool outcomes.

Follow `.editorconfig`, `pyproject.toml`, configured formatters, and nearby
files. Preserve intentional differences between Python support and individual
tool targets. Do not use emojis in user-facing content unless the surrounding
context or tooling requires one.

Document current behavior and supported guarantees first. Keep explanations
small; retain relevant intentional plans and operational history only. Put
speculation and rejected designs in task artifacts or an explicitly maintained
decision record, not ordinary documentation. Promote lessons only when they
become stable, reviewable regressions.

Use concise, unprefixed Git subjects, a blank line before the body and trailer
footer, and `Token: value` trailers. Do not use Conventional Commit prefixes.
Goal-linked delivery commits need `Goal-Ref` and `Attempt-ID`; optional
`Learning-Proposal` precedes the generated final `LLM-disclaimer` trailer.
