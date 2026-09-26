---
title: Agents
---

## Start with the task

This is the repository-wide default; a nearer `AGENTS.md` takes precedence
within its subtree. Establish the requested outcome and authority before
mutating. Read the applicable policy chain and the owning `README.md`; inspect
`BUILD.bazel` and `MODULE.bazel` only when implementation or validation depends
on them, and load only the skills the phase needs.

## Authority and scope

Keep one owner for each fact: the user request owns outcome and authority,
`AGENTS.md` owns policy, the nearest `README.md` owns component boundaries,
`CODEOWNERS` owns review accountability, BUILD and MODULE files own executable
and dependency structure, canonical skills own procedures, OpenSpec changes own
maintained work state, runtime providers own observed capabilities, and Git and
delivery receipts own publication state.

Each fact has one authoritative source: independently maintained duplicates are
defects even when the original is hard to reach, so extend or parameterize the
owner. Generated projections, fixtures, and source-identifying summaries are
projections, not duplicates.

A question authorizes investigation, not mutation, and preserves granted
authority. The primary agent owns material trade-offs and consequential
verdicts. Reviews and audits are evidence, not new acceptance criteria: fix
in-scope blockers and tiny related defects, and report other findings.

- All infrastructure provisioning and persistent configuration MUST be defined
  in checked-in infrastructure as code and deployed through its owning
  Terraform, Ansible, or other declarative workflow. Any exception, such as an
  imperative bootstrap or emergency change, requires explicit approval for that
  exact exception, with its scope and IaC reconciliation recorded; general
  deployment authorization does not approve one, and read-only inspection is
  not one. Secret values stay in Vault; IaC carries only their references.
- Never run state-changing infrastructure operations without the user's exact
  request for that operation and scope; validation must not mutate live systems.
- Treat an automatic approval rejection as a strategy signal: diagnose it,
  choose a materially safer authorized approach, or ask. Never route around it
  or retry a rejected, failed, or rate-limited escalated operation immediately.
- Every authorized change follows delivery: validate, commit, push, and offer a
  pull request through `repo-delivery` unless the user withholds publication.

## When to load a skill

This table owns routing; each skill owns its procedure.

| When                                                                     | Load                                         |
| ------------------------------------------------------------------------ | -------------------------------------------- |
| Before the first mutation or task-scratch write                          | `repo-workspace`                             |
| Substantive question, including a mixed question-and-action request      | `answer-question`                            |
| Material, consequential, or repeatedly failing choice                    | `decision-review`                            |
| Creating or moving source, or choosing a directory layout                | `project-layout`                             |
| Any Bazel invocation, BUILD or `.bzl` mechanics, or validation scope     | `bazel-agent`, `repo-bazel`                  |
| Go implementation, refactoring, or review                                | `repo-go`                                    |
| A standalone nested workspace with its own `MODULE.bazel`                | `bazel-nested-module`                        |
| External dependencies, archives, toolchains, or lockfiles                | `repo-external-dependency`                   |
| Gazelle generation behavior or a language plugin                         | `repo-gazelle-plugin`                        |
| Structure-aware search or rewrite by syntax shape                        | `ast-grep`                                   |
| `infra/**`, `host_bot`, `al.lua`, tf/Ansible/DNS, deployment diagnostics | `repo-infra`                                 |
| Credentials, tokens, private keys, Vault policy, secret-bearing config   | `repo-secrets`                               |
| Workflow changes, reusable CI commands, or CI diagnostics                | `repo-ci`                                    |
| Hugo site or theme, or landing-site onboarding                           | `repo-hugo`                                  |
| Android app build, packaging, or publication                             | `android`                                    |
| A `.blend` asset, or Blender work judged by supplied-reference likeness  | `repo-blender`, `blender-reference-fidelity` |
| Host Codex model, provider, authentication, or shared-config migration   | `codex-migration`                            |
| Durable specifications, changes, or continuation state                   | `openspec`                                   |
| Adding or updating a repository skill                                    | `bazel-rules-skill`                          |
| Proofreading, polishing, or rewriting supplied prose                     | `spellcheck`                                 |
| Publishing an article to the blog                                        | `alwaldend-blog`                             |
| Complete repository build-and-test health check                          | `full-repo-check`                            |
| Computing a development, nightly, or weekly version                      | `versioning`                                 |
| Preparing publication or final handoff                                   | `repo-delivery`                              |
| Synchronizing an advancing remote base, or rebasing task-owned commits   | `git-rebase-remote`                          |
| Reviewing a substantial or inefficient session at task close             | `agent-ergonomics-review`                    |

## Isolate and verify changes

Use a dedicated feature branch in its own linked worktree for every task that
can modify repository files or task-owned scratch, and verify both before the
first write. The default branch and checkout are read-only unless the user
explicitly authorizes that exact task there. Preserve unrelated edits and never
discard, auto-stash, commit, or rewrite shared, human-owned, unrelated, or
ambiguous history; never lose progress on either branch side.

Verify source and a representative output against the exact candidate; command
success alone is not acceptance. Accept and commit all changes produced by the
repository's configured formatters, even outside the initial task scope, but
inspect the diff and never include semantic changes under that exception.
Exclude disposable build outputs and include required generator-maintained
files updated through their owning workflow. Commit and push legitimate
task-owned source and configuration unless the user says otherwise, with binary
artifacts only through Git LFS. `repo-delivery` owns gates, commits, push, and
requests.

## Testing

- NEVER write unit tests after writing code.
- Highly prefer end-to-end (E2E) tests as the sole testing mechanism. Use them
  to verify that complex features work. At the end of E2E tests, produce a
  verifiable and repeatable artifact.
- If you must test a system in isolation, FIRST write all the ways it could
  fail, THEN write the code.
- Tautological tests are considered harmful.
- Change-detector tests are considered harmful.
- Do not create regression tests for bug fixes without a genuine gap in
  behavior testing.

## Use bounded tools and evidence

**Adding new external dependencies is strongly discouraged and requires the
user's explicit approval for the specific dependency before it is added.**
A request to implement a feature, fix a bug, or use a capability does not
approve a new dependency. Do not add one merely because it is convenient or
an existing dependency lacks a preferred API.

Before searching for a new external dependency, inspect existing repository
code, standard libraries, manifests, lockfiles, and locally available
dependencies. Prefer reusing those options. If a new dependency is necessary,
first prepare a concrete proposal naming it, explaining why existing options
are insufficient, and describing its maintenance and supply-chain costs;
then ask the user for approval. Do not install, vendor, declare, or add it to
source or build configuration before approval.

Prefer supported live Cordis handlers, then purpose-built tools, MCP
capabilities, and repository Bazel targets; use a host shell only when no
suitable entry point exists. Use `rg`, `rg --files`, bounded `find`, or Bazel
queries rather than recursive searches, wrap commands in a `timeout`, and
extract diagnostic fields from logs before display because raw dumps can leak
protected data.

Reuse evidence while inputs are unchanged; empty results and no-ops are
evidence with bounded coverage, not proof of absence. Repeat an operation only
when changed inputs, new evidence, or expected asynchronous progress justify
it, and record revisions and observation times. After interruption, resume from
the recorded action and recheck mutable dependencies. Preserve failure evidence
and a stable defect name.

## Disclosure, naming, and communication

Checked-in source, documentation, and fixtures are public. Public source and
eval output from public fixtures in an isolated public workspace may be sent to
an external service the task calls without separate confidentiality approval;
this does not authorize exposing credentials, personal information, or
secret-bearing data. Inspect local, generated, untracked, live, and
infrastructure artifacts before disclosure because they can contain protected
data, but their origin does not make them confidential. Ordinary build, test,
and lint diagnostics and non-secret, non-personal operational facts may be
reported, including for `infra/`. Keep source disclosure, target visibility,
build consumers, artifact publication, and secrets distinct, and treat tree
"privacy" as repository-internal unless it identifies protected content.

Use precise, professional vocabulary for every name you choose, including
files, scripts, directories, functions, variables, classes, commits, and
branches, and name the purpose or behavior accurately for a technical
specification. Avoid slang, casual shorthand, chat-culture terms, jokes, and
cute wordplay, and prefer established domain terminology and standard
abbreviations. Name a dashboard-deploying script `deploy_dashboards.sh`,
never `push_dashboards.sh`.

When the main agent loads a skill, say which one and why it applies. Follow
`.editorconfig`, `pyproject.toml`, configured formatters, and nearby files,
preserving differences between Python support and individual tool targets, and
use no emojis in user-facing content unless the tooling requires one. Document
current behavior and guarantees first, keeping explanations small and retaining
only relevant plans and history; put speculation and rejected designs in task
artifacts or a decision record.
