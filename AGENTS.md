---
title: Agents
---

## Repository guide for agents

Use this file as the repository-wide default. If a more deeply nested
`AGENTS.md` is added, its instructions take precedence for that subtree.

## Fast orientation

For every task, establish the requested outcome and authority first. Resolve
the affected path and workspace, read the nearest owner documentation, load
the smallest applicable skill set, and inspect before mutating. Use a durable
goal when work needs retries or continuation. Select the least costly check
that can answer the acceptance question, bind evidence to the exact candidate,
and use the repository delivery workflow for publication. Promote a lesson
only after it becomes a stable, reviewable regression.

Do not attach a token budget to a goal or task unless the user explicitly
tells you to do so.

## Settled results and the no-loop rule

A query is settled the moment it returns its expected-shaped result. Empty
output, zero matches, a successful no-op, and a negative search all count as
answers. A settled result is never inconclusive; do not re-run it hoping for
different output.

Exit status is irrelevant to settlement. A query that exits 0 with no
matches is exactly as settled as one that exits 0 with results.

Two invocations count as the same query when they ask the same decision
question of the same target through causally equivalent paths with the same
relevant inputs. Cosmetic changes to flags, ordering, quoting, glob patterns,
output limits, or tools do not make a new query. A changed flag, input, or path
does make a different query when it can change decision-relevant evidence or
discriminate a named failure hypothesis. Judge equivalence by causal effect,
not command bytes.

For queries, writes, and larger attempts, "the same setup" means the same
decision question, target, relevant inputs, and causal execution path. It does
not mean that a target or desired outcome may be tried only once. A retry is
allowed after the failure mechanism is identified or narrowed to a falsifiable
causal hypothesis, and the changed input or execution path could repair,
bypass, or discriminate that mechanism if the hypothesis is true. Before
retrying, reference the prior failure and record the causal change and the
observable difference it predicts. A new label, reordered command, different
tool with equivalent behavior, or parameter unrelated to the failure
hypothesis is still the same setup.

After any settled result, the next action must be one of:

- a write, patch, or other state change that acts on the settled result or
  advances a revised plan;
- a materially different query class (a different question, or a different
  target, relevant input, or causal execution path that can change the
  decision; changing tools alone does not qualify);
- a revision to the plan itself; or
- a question whose answer is needed to choose the next causal action.

Re-running a settled query is never one of these options. If you notice you
are about to issue a query you have already settled, treat that recognition
as a hard stop: name the finding, then take one of the four allowed actions.

The same hard stop applies to writes. A write is settled when it completes;
if it did not produce the intended state, diagnose the cause before writing
again. Do not issue identical or equivalent mutations because the first
mutation did not visibly change an artifact. In particular:

- Never make a no-op patch, rerun a replacement that reports no substitution,
  or otherwise count an unchanged file as progress.
- After two consecutive edits that leave the edited bytes unchanged, stop,
  identify the unchanged source or cache dependency, and repair that
  dependency.
- If cached or generated outputs appear stale, make one deterministic
  cache-invalidating action, such as changing an input name or forcing one
  target to rebuild. Do not alternate cached rebuilds with no-op writes.
- Before declaring a change complete, verify the source-level result and a
  representative output artifact, not only the output.

Variation is not progress. Repeating the same fact-finding question with a
different command, path, regex, or output format is still a repeat. After two
searches that do not produce a decision, stop searching, record what is known
and unknown, and move to the next decision or act on the best available
answer. Empty, inconclusive, and partially matching results are answers, not
invitations to reformulate. Prefer writing the change under the best current
assumption over gathering more evidence for a reversible in-scope choice.

The canonical repository-agent documents are:

- [current state](projects/agents/docs/current-state.md), an evidence snapshot;
- [architecture](projects/agents/docs/architecture.md), the composition
  contract;
- [roadmap](projects/agents/docs/roadmap.md), intentional future work rather
  than current behavior or authority; and
- [durable goals](projects/agents/goals/README.md), versioned attempts,
  acceptance state, and evidence provenance.

Use one mutation authority for each fact:

- the user request owns intended outcome and granted authority;
- the nearest `AGENTS.md` owns applicable agent policy;
- the nearest component `README.md` owns purpose and local boundaries;
- `CODEOWNERS` owns review accountability;
- `BUILD` and `MODULE.bazel` files own executable and dependency structure;
- canonical skills own reusable procedures;
- goal records own durable work state and acceptance history;
- each runtime provider owns its observed health and capabilities; and
- Git plus repository-delivery receipts own candidate and publication state.

The architecture composes these authorities; it does not override them. A
derived view must identify its sources, version or digest, observation time,
unavailable fields, and truncation. If authorities disagree, stop at the
conflict and resolve it at the owning source rather than choosing whichever
copy is convenient.

## Repository visibility

This is a public repository. Treat checked-in source, documentation, and test
fixtures as public information. When a task explicitly calls for an external
service, that checked-in public content may be sent to the service without a
separate confidentiality approval. Eval output derived only from checked-in
public fixtures and an isolated public workspace has the same classification.
This does not make credentials, personal information, or secret-bearing data
safe to disclose. Local, generated, untracked, live, or infrastructure-related
content is not confidential merely because of its origin, but inspect it
before disclosure because it can contain secrets or personal information.
Non-secret, non-personal operational facts are fair game in reports.

Keep these policy axes distinct: checked-in source disclosure, Bazel target
visibility, allowed build consumers, artifact or documentation publication,
and the presence of secrets or personal information. One axis does not imply
another. In particular, words such as “private” in tree policy describe
repository-internal use unless that policy identifies secret or personal
content; they do not make ordinary source or operational facts confidential.

## Scratch files (mandatory)

- Put every task-owned download, generated report, log, extracted archive,
  tool cache, and temporary file under an ignored, task-specific
  `out/<task>/` directory in the applicable workspace.
- Point configurable tool scratch directories at `out/<task>/` as well. Do
  not use `/tmp` for task-owned files merely because they are temporary. Set
  configurable temporary-directory and cache environment variables to the
  task's `out/<task>/` directory.
- Secret-bearing temporary material may use an ignored, task-private
  `out/<task>/` directory with restrictive permissions, minimum lifetime, and
  explicit cleanup. “Outside the repository” means outside tracked or
  committable source: never track, stage, commit, or promote secret values as
  durable evidence. Any unavoidable secret-bearing raw artifact remains
  ignored, access-controlled, and short-lived.
- Treat operating-system temporary storage, including `/tmp`, as an absolute
  last resort. Use it only when a tool cannot be directed into the applicable
  workspace's `out/<task>/` directory. Remove any task-owned residue there
  before handoff; never delete unrelated files or broad temporary directories.

## Decision-making

Use the `decision-review` skill before committing to a material design,
security, operational, costly, irreversible, or repeatedly failing choice.
Treat both the user's proposal and your current plan as hypotheses: identify
the strongest reason a domain expert would reject them, test that reason
against evidence, and compare credible alternatives. The primary agent owns
the verdict and must not delegate it to a random subagent. Optimize for the
user's actual goal, not agreement or the cheapest satisfaction of stated
constraints. State material trade-offs to the user; routine reversible
implementation choices do not need this review or narration.

Expert guidance, audits, and adversarial review are evidence, not authorization
or new acceptance criteria. Do not turn every hypothetical weakness into
required work: classify findings against the user's stated scope, fix only
in-scope blockers, and report optional hardening separately. If an overbroad
claim would require unrelated machinery, narrow the claim instead of expanding
the implementation. Do not interrupt in-scope work for routine uncertainty or
optional hardening. If completing the requested outcome requires a material
scope expansion that existing context does not authorize, ask before acting;
otherwise choose the smallest reversible in-scope approach and continue.

Treat auto-review rejections of escalated actions as authoritative strategy
signals, not obstacles to route around. When a rejection stops you, rethink
the approach before retrying or escalating further: the rejection usually
means the approach itself is wrong, such as mutating host state the sandbox
already provides or bypassing a required workflow. Retry with a materially
safer alternative, or stop and ask the user.

## Parallel work

At the start or resumption of nontrivial work, and whenever the approach or
available work changes materially, enumerate the ready independent
workstreams. When two or more can produce useful, independently reviewable
outputs concurrently, actively delegate them unless a concrete dependency,
shared-state risk, or immediate review-and-integration bottleneck makes
sequential work faster. Record that reason in the task plan or result when
remaining sequential. Do not wait for the user to request delegation again.

Prefer a few bounded workstreams such as independent research, candidate
variants, disjoint implementation modules, implementation-blind review,
verification, or artifact preparation. Keep one coordinator for canonical
state and shared artifacts. Do not fragment trivial operations, recursively
delegate without a separate benefit, or create speculative work merely to
occupy available agent slots. Recheck parallelism when work closes, stalls,
changes strategy, or exposes a new independent task.

## Questions

Use the `answer-question` skill whenever the user's requested outcome includes
a substantive question, including when the same request also authorizes an
action. Interrogative text inside quoted material or an inert transformation
payload is not itself a request for an answer. When a substantive question
asks for a material decision, use both `answer-question` and
`decision-review`. A question, including “can we,” “should we,” “why,” or “do
we need,” is a request for information and not authorization to modify files,
settings, pull requests, deployments, or other external state. Read-only
investigation is allowed when it supports a truthful answer. If a message
combines a question with an explicit request for action, act only within that
stated scope.

## Making changes

- Perform every task that can modify repository files or task-owned scratch on
  a dedicated feature branch in its own linked worktree. Treat the default
  branch and its checkout as read-only and off limits unless the user
  explicitly authorizes that exact task to use them. Read-only inspection may
  use the default checkout. Verify the branch and worktree before the first
  task-owned mutation, and keep the default checkout clean.
- When the primary agent encounters a small, bounded bug in repository tooling
  or the affected project, fix and validate it instead of silently working
  around it. Delegated workers remain within their assigned scope and report
  discovered bugs to the coordinator rather than expanding their task. If a
  bug requires substantial effort, redesign, or a rewrite, preserve the
  evidence and report it separately instead of derailing the requested work;
  any temporary workaround must be explicit and must not hide a correctness
  or safety failure.
- An agent may resolve a Git conflict without additional user confirmation
  when the branch rewrite is already authorized and the conflict is minor: its
  scope is bounded, the intended combined result is unambiguous from the task
  and repository evidence, and the resolution preserves all relevant changes
  from both sides. Inspect the exact three-way hunks, stage only explicit
  resolutions, and rerun every check invalidated by the rewrite. A conflict is
  not minor when resolving it requires guessing intent, choosing between
  unsupported behaviors, dropping work, or changing security, infrastructure,
  public-interface, or data semantics beyond the authorized task; abort and
  ask for direction in those cases. Continue to use the applicable rebase and
  delivery workflows, including exact remote leases.
- Loading the `project-layout` skill is REQUIRED before creating or moving
  any directory in the source tree, or when deciding a directory layout
  anywhere in the repository. This is a hard gate, not a suggestion: do not
  create the directory until the skill has been loaded and its layout has
  been checked against the planned path. The skill's role-based layout is
  mandatory for new projects and new subdirectories. Legacy layouts such as
  `main/` are grandfathered and must not be used as precedent for new paths.
- Read the nearest `README.md`, `BUILD.bazel`, and `include.MODULE.bazel`
  (when present) for the area being changed.
- Name projects using only ASCII letters, digits, and underscores
  (`[a-zA-Z0-9_]+`).
- Prefer a small, target-specific change. This is a large monorepo, so query,
  build, and test the affected Bazel package before considering `//...`.
- Prefer Go for repository automation and scripts. Expose them as Bazel
  `go_binary` targets and invoke them with `bazel_agent bazel run`; use
  another language only when Go or Bazel would materially complicate the task.
- Do not use `genrule`. Write a proper Bazel rule, a Go binary, or a separate
  checked-in script instead; `repo-bazel` owns the detailed guidance.
- Loading the `repo-delivery` skill is REQUIRED when making changes that will
  land in the source tree, including every implementation task. The skill owns
  staging, feature-branch commits and pushes, pull request maintenance, review
  comment handling, and the final delivery report; how to apply it is the
  agent's judgment under the skill's guidance.
- Keep all task-owned scratch under the ignored `out/` area per the Scratch
  files section below; never commit temporary files.
- Commit and push all legitimate changes (source, documentation, and
  configuration) unless the user explicitly says otherwise. Binaries,
  temporary files, task scratch, secrets, and generated artifacts are not
  covered by this default.
- Commit binaries only when they are required, are not temporary files, and
  are tracked by Git LFS. Do not commit binaries otherwise.

## Tooling

- Go code is compiled and validated by Bazel's `@rules_go//go` toolchain, not
  by a host-installed Go distribution. Use Bazel test targets and the pinned
  `bazel_agent` entry points for Go validation; do not use host `go` commands
  as the acceptance check.
- When the main agent loads a skill, tell the user which skill it loaded and
  how its context applies to the current task, so the user knows what
  instructions the agent is operating under.
- Before a non-obvious tool call or delegation, briefly tell the user its
  purpose, scope, and relevant side effects. Report the outcome when the user
  interface may otherwise show only an opaque tool event. Group routine
  repeated checks instead of narrating each mechanical call.
- Prefer live Cordis handlers over host-shell commands for repository context,
  bounded file reads and searches, Git inspection, and other operations the
  active Cordis catalog supports. Discover the catalog before assuming a
  handler is unavailable. When a recurring repository workflow lacks a safe,
  bounded handler, consider adding it to Cordis; do not expand the current
  task's authority merely to avoid a one-off shell command.
- Prefer purpose-built tools, MCP capabilities, and repository Bazel targets
  for work they support. Use direct host-shell commands only when no suitable
  tool, MCP, or Bazel target exists.
- Do not introduce shell scripts or shell-based Bazel rules unless the task
  explicitly requires them. Prefer Go programs and Bazel-native rules such as
  `go_binary`, `go_test`, and repository Bazel checks over shell wrappers.
- Acquire tools used by Bazel hermetically: use pinned, checksummed binary
  archives or Bazel-integrated package managers driven by checked-in manifests
  and lockfiles. Do not silently depend on a host-installed tool or an
  undeclared lifecycle download.
- Run every build action inside a Bazel sandbox and without network access.
  If an action genuinely needs network access or must disable sandboxing,
  document the reason in the owning project's README before enabling it and
  prefer the narrowest possible exemption.
- Run every shell command with a timeout so it cannot hang unexpectedly,
  unless there is a good reason not to (for example, a command that must
  intentionally run for a long time or an interactive session). Prefer a
  short timeout that fits the command's expected duration, and use
  `timeout <seconds> <command>` as the wrapper.
- Hermetic tool acquisition is important; hermetic tool output has lower
  priority. A tool may intentionally access the network or produce
  environment-dependent artifacts when the requested workflow permits it.
  Continue to apply the repository's authorization, secret-handling, and
  external-side-effect rules.

## Documentation

- Document current behavior and supported guarantees first.
- Include future plans only when they are intentional and relevant to users;
  include historical behavior only when it helps operate, migrate, or
  understand the current system.
- Do not add rejected ideas, speculative alternatives, or incidental design
  exploration to ordinary maintained documentation. Keep that material in the
  task's goal or research artifacts, or in an explicitly maintained decision
  record.
- Prefer the smallest durable explanation that helps a reader use or maintain
  the system.

## Searching

Do not use recursive `grep` or `ls`. Use `rg`, `rg --files`, `find` with a
bounded depth, or `bazel_agent bazel query` instead. Never search from the
filesystem root (`find /`); it can hang for minutes on large dependency
caches. Use an explicit root and `-maxdepth`, and prefer `rg --files` with a
path.

For structure-aware queries that plain-text search cannot express, use the
repository-pinned ast-grep tool and the `ast-grep` skill. Invoke it through
`bazel_agent bazel run //tools/ast-grep:ast-grep --` with the pattern and language
after the separator. Do not substitute an ast-grep binary from `PATH` or a
separate download; read `tools/ast-grep/binary_toolchain.json` for the pinned
version and follow the `ast-grep` skill for pattern validation and safe
rewrites.

Preserve failure evidence and use a stable name for recurring defects.

## Repository map and boundaries

- `projects/`: product and reusable project code. Its Bazel targets may use
  public visibility, its artifacts may be published through owned release
  workflows, and its targets may be production build dependencies.
- `infra/`: repository-internal infrastructure definitions (Terraform,
  Ansible, Flux, DNS, and host/service configuration). They must not be
  published as artifacts or consumed by production build actions.
- `tools/`: repo-wide build rules and toolchains. Except for toolchain types,
  these are repository-internal and not for production build actions; tools
  may be used by tests. Tool targets intended across the repo normally use
  `visibility = ["//:__subpackages__"]`.
- `data/`: repository-internal checked-in data and documentation assets. It
  may be used in builds, but must not be independently published.
- `third_party/`: repository-internal vendored or externally sourced code. It
  may be used in builds, but must not be published as first-party source.
- `users/`: repository-internal user-specific code and infrastructure. It must
  not be published or consumed by production builds.

The authoritative policy for each tree is its top-level `README.md`. Preserve
these visibility and publication boundaries when adding dependencies.

## Bazel and dependency management

Bazel is the primary entry point. The `bazel-agent` skill owns every
`bazel_agent` invocation and the sole bootstrap exception; the `repo-bazel`
skill owns BUILD, `.bzl`, `MODULE.bazel`, and dependency mechanics, generated
files, and the narrowest practical validation. Load both before Bazel work and
follow them instead of repeating their flags here. Do not call `bazel`
directly, use raw BEP output (it can contain secrets), or hand-edit generated
files — run the update command in their header.

## Infrastructure safety

- Ordinary read-only build, test, lint, and format-check output for targets
  under `infra/` is safe to display and must not be suppressed merely because
  the target is infrastructure-related. The restrictions below apply to
  secret-bearing values and artifacts, not normal build diagnostics.
- State, plan output, inventories, decrypted configuration, and other live or
  generated artifacts can contain credentials, other secrets, or personal
  information.
  Inspect or structurally summarize them before disclosure; never paste a
  secret or personal information into reports, logs, or commits. Non-secret,
  non-personal operational details are public and may be reported.
- `al.lua` files use the repository's custom configuration DSL and often wire
  Vault AppRole authentication into generated commands. Follow a nearby
  service's `al.lua` and Bazel target rather than invoking tools directly.
- Terraform is wrapped by `terraform_binary_map`. Typical targets are
  `:tf.fmt_check`, `:tf.plan`, and `:tf.apply`; the first two are validation,
  while `apply`, `deploy`, `destroy`, `import`, `migrate`, `state`, and
  `force_unlock` can change remote state.
- Do not run any state-changing infrastructure target unless the user
  explicitly requests that exact operation and scope. For ordinary code
  changes, limit verification to formatting, validation, queries, builds, and
  tests that do not contact or mutate live systems.
- Do not commit `.terraform/`, state files, plans, environment files, or
  machine-local credentials. Existing `.gitignore` rules are a backstop, not
  permission to create or inspect secret material unnecessarily.

## Style

Follow `.editorconfig` and the closest existing files:

- UTF-8, LF endings, a final newline, no trailing whitespace, and a preferred
  maximum line length of 79.
- Four spaces by default (including JSON); tabs only for Go and Makefiles; two
  spaces for YAML, Markdown, HTML, Proto, and TOML.
- Python formatting and analysis settings live in `pyproject.toml`. Note that
  the project metadata supports Python 3.10+, while some tool configurations
  deliberately target Python 3.9 compatibility; do not casually synchronize
  those values.
- Formatters are always right. If the repo-wide formatter changes files
  outside the focused change, commit those changes in the same candidate.
  The `repo_quality_test` gate validates the whole tree, so stale formatting
  anywhere blocks delivery until the formatter output is committed.
- Do not use emojis in any user-facing content: documentation, code comments,
  commit messages, pull requests, issues, chat or session messages, and other
  user-facing artifacts. Exceptions only when an emoji is required by the
  surrounding context or tooling.

Commit messages follow ordinary Git conventions. Use a concise, unprefixed
subject; a blank line before any body; and a blank line before the trailer
footer. Do not use Conventional Commit or other artificial subject prefixes.
Use `Token: value` trailers in Git's `git interpret-trailers` format. Goal-
linked delivery commits require both `Goal-Ref` and `Attempt-ID`; optionally
add `Learning-Proposal`. The generated `LLM-disclaimer` remains the final
trailer.

## Verification

`repo-bazel` owns the selection and ordering of build/test checks for changed
packages and repository-wide runs. Always start with `git diff --check`, then
the narrowest relevant package checks; run package formatters (buildifier,
black, mypy) for the files changed and the repository `//:buildifier_test`
for BUILD or `.bzl` changes. The checked-in pre-commit configuration supplies
repository hygiene checks; installing `//:write_git_hooks` is optional and
agents should still run the relevant checks explicitly.

## Escalated operations and settled facts

- Never immediately retry a rejected, rate-limited, or failed escalated
  operation. One retry after a real delay at most, then stop and change the
  approach or ask; the rejection usually means the approach is wrong.
- Git history, help output, and prior results are immutable within a task;
  record the conclusion in the task's `out/<task>/notes.md` (or plan) after
  the first query and move to the next decision.
- After an interrupt, resume at the exact next write; do not re-explore.
