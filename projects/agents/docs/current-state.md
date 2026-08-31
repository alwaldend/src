---
title: Repository agent system current state
description: >-
  Evidence-backed snapshot of the repository-wide agent control system,
  its present guarantees, and its unresolved joins
tags:
  - agent
  - architecture
  - repository
---

# Repository agent system current state

> Snapshot date: 2026-08-31
> Tracked source revision:
> `1423dce5fab45ce5223caeb6a24791bf1a2cc3ff`

## Scope

This document describes repository-wide control surfaces, not the behavior of
individual products. It covers how an agent discovers policy and capability,
maps ownership and workspaces, executes and validates work, retains evidence,
delivers changes, and turns feedback into reusable repository knowledge.

The snapshot distinguishes checked-in behavior from deployed or host state.
It does not claim that a configured external service is healthy, that a manual
evaluation has run, or that an ignored local artifact exists merely because a
workflow can create it.

## Evidence record

The immutable audit, criteria, methods, and evidence bindings are retained at
`projects/agents/goals/repo-agent-system/` in repository source. This document
is the maintained public-system synthesis; its direct source links are the
documentation projection's evidence surface.

## Phase 1 declared baseline

The Phase 1 candidate adds the shared `agents.alwaldend.com/v1alpha1`
contracts and a report-only registered-universe check. Registry digest
`sha256:ee6024206f87005175f0213b4d81db7be9810988d0dadcbe8c5a5347c3b745e7`
currently covers seven registration authorities, 21 skills, 28 project roots,
nine workspaces, two maintained goals, ten Cordis gateway tools, four
supported direct binaries, 24 owner-declared operations, and two generated
artifacts. The report names two legacy Terraform argument surfaces as
`requires_migration`; it does not silently call them safe or absent.

Task/run manifests and Cordis scratch now use
`out/<task>/mcp_cordis/runs/<run>/`, while Bazel action/test temporary storage
remains Bazel-owned. `bazel_agent doctor` exposes a bounded, non-environment-
dumping view of runner/source identity, Bazelisk pins, platform, profile/rc
composition, task scratch, and stale-install state.

## Direct assessment

The repository has strong specialized mechanisms and weak system-level joins.
Its build runner, skill packaging, runtime repository tools, goal integrity,
and delivery state checks are individually careful. Their relationships are
mostly conveyed by prose, copied identifiers, path conventions, and agent
judgment rather than a shared typed contract.

The system is usable by an agent that already understands it. A zero-context
agent must reconstruct policy and routing; topology and ownership; live
capability and workspace state; durable intent and evidence; and delivery,
review, version, and release state before acting.

No repository-owned interface currently joins those planes into one bounded,
versioned view of the current task, applicable policy, available capability,
allowed effects, expected cost, and reusable evidence.

## System at a glance

| Layer           | Current authority                                                                            | Present role                                                               |
| --------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| Human entry     | [`README.md`][snapshot-root-readme]                                                          | External homepage, docs, setup, and license links                          |
| Agent policy    | [`AGENTS.md`][snapshot-agents]                                                               | Repository defaults for authority, safety, layout, Bazel, and verification |
| Path policy     | Top-level tree READMEs                                                                       | Publication, visibility, and dependency constraints by tree                |
| Skill source    | Project-owned `SKILL.md` targets                                                             | Task-specific procedural knowledge                                         |
| Skill discovery | [`BUILD.bazel`][snapshot-root-build] and `.agents/skills`                                    | Exact flat projection of 20 registered skills                              |
| Build execution | [`projects/bazel_agent`](../../bazel_agent/)                                                 | Thin, consistent pass-through to pinned Bazel                              |
| Quality         | [`tools/repo_quality`](../../../tools/repo_quality/)                                         | Whole-index formatting and language-specific checks                        |
| Runtime         | [`.codex/config.toml`][snapshot-codex-config] and [`projects/mcp_cordis`](../../mcp_cordis/) | Required MCP host, dynamic packages, bounded repository and Git tools      |
| Durable intent  | [`projects/goal`](../../goal/)                                                               | Versioned goals, criteria, attempts, evidence, and projections             |
| Delivery        | [`tools/repo_delivery`](../../../tools/repo_delivery/)                                       | Exact Git/forge preparation, publication, and review mutations             |
| Versioning      | [`tools/versioning`][snapshot-versioning-skill]                                              | Global version calculation and guarded local ref creation                  |
| Documentation   | [`projects/rules_docs`](../../rules_docs/)                                                   | Recursive Markdown packaging into the repository website                   |

## Strengths

- **Policy and safety.** The root guide distinguishes questions from
  authorization, requires task-specific scratch, limits infrastructure
  mutation, warns about potentially secret-bearing diagnostics, and directs
  agents toward narrow validation
  ([`AGENTS.md`][snapshot-agents]). Bazel action defaults
  use a strict environment and deny sandbox network by default
  ([`preset.bazelrc`][snapshot-bazel-preset]).
- **Canonical sources and projections.** Skills live with their narrow owner;
  `.agents/skills` is a reconciled projection that rejects unsafe, duplicate,
  cross-repository, and stale links
  ([`projects/agents/README.md`](../) and
  [`rules_skill`](../../rules_skill/)). Goal YAML and digest-bound
  artifacts are canonical while README is generated; resource versions,
  locking, and one-writer rules protect publication
  ([`goal CLI`](../../goal/cmd/goal/)).
- **Reproducible execution.** `.bazeliskrc` pins Bazel and its archive digest.
  `bazel_agent` is a small pass-through that selects the nearest workspace and
  preserves Bazel's status. Repository quality uses Bazel-acquired tools, and
  the full checker covers the root plus eight registered nested workspaces
  while retaining task-local raw logs
  ([`bazel_agent`][snapshot-bazel-agent-main] and
  [`full-repo-check`][snapshot-full-repo-check-skill]).
- **Bounded observation.** Cordis keeps a fixed MCP gateway while dynamic tools
  change. Repository and Git tools expose path, byte, result, timeout, and
  truncation bounds; Git avoids prompts and lazy fetch; supervised process
  groups are joined before results settle
  ([`repo_context.mjs`][snapshot-repo-context] and
  [`process_supervisor.mjs`][snapshot-process-supervisor]).
- **Exact delivery state.** Repository delivery binds repository, endpoint,
  refs, OIDs, tree, forge, remote expectation, and path scope, then rechecks
  state around mutation. Versioning separately owns the global version and
  guarded local ref transactions
  ([`repo_delivery`](../../../tools/repo_delivery/) and
  [`versioning`][snapshot-versioning-skill]).

## Verified repository facts

### Entry points and topology

- The 239-line root `AGENTS.md` is the only tracked agent guide. No subtree
  currently supplies a narrower `AGENTS.md`.
- The 14-line root README links external destinations but does not identify
  the six repository trees, Bazel, agent capabilities, current goals, or a
  local system map ([`README.md`][snapshot-root-readme]).
- `projects/`, `infra/`, `tools/`, `data/`, `third_party/`, and `users/` are
  the declared content boundaries. `.agents/`, `.codex/`, root configuration,
  and `out/` form a separate control and task-state plane that the declared
  map does not present together.
- There are 182 immediate component directories under the six content trees.
  181 have both README and BUILD files. `tools/bazel_contracts` is the sole
  exception.
- There are nine tracked `MODULE.bazel` workspace roots and 73
  `include.MODULE.bazel` fragments.

### Documentation and metadata

- The snapshot contains 409 tracked README files; 407 begin with YAML
  frontmatter and all 407 declare a title.
- All 28 project roots have title, description, and exactly one lifecycle
  status from the checked Phase 1 vocabulary.
- `docs_filegroup` packages Markdown and direct documentation dependencies. It
  does not carry owner, lifecycle, policy, workspace, capability, or action
  metadata
  ([`projects/rules_docs/docs/defs.bzl`][snapshot-rules-docs-defs]).
- Root `//:docs` directly aggregates all six content trees but sets its own
  `srcs = []`, so the root README is not part of that aggregate
  ([`BUILD.bazel`][snapshot-root-build], lines 92-100).
- The website source consumes `//:docs` and renames README basenames into Hugo
  section indexes
  ([`projects/alwaldend.com/BUILD.bazel`][snapshot-site-build],
  lines 105-124).

### Capabilities

- The root registry and tracked discovery directory contain the same 21 skill
  names: 19 under `projects/agents`, the goal skill, and versioning.
- Canonical skills have strong artifact validation, but most behavioral
  targets validate configuration rather than outcomes. Descriptions drive
  routing; composition, effects, authority, cost, and observed quality are not
  one queryable contract ([`projects/agents/README.md`](../)).
- The checked-in Codex runtime declares Cordis as required. Skills, Cordis
  packages, native tools, and optional host connectors remain separate
  capability inventories.

### Intent, evidence, and delivery

- Project goals are intentionally owner-local. Goal list, graph, and catalog
  operations each start from one explicit goals root; a session binding is not
  a repository catalog
  ([`projects/goal/cmd/goal/README.md`](../../goal/cmd/goal/)).
- Goal resources prove structural identity, concurrency, criteria revision,
  and artifact integrity. Stable defect, exact external subject, regression
  set, dominant failure, and next action remain primarily in Markdown.
- Delivery proves exact prepared/current Git state but receives no attested
  validation set. Versioning creates guarded local release refs, while remote
  tag publication lies outside both owned workflows. No traversable record
  joins review finding, disposition, fix, regression, and delivered result.

## What each current proof establishes

| Mechanism              | It establishes                               | It does not establish                          |
| ---------------------- | -------------------------------------------- | ---------------------------------------------- |
| README frontmatter     | parseable page metadata                      | semantic completeness or ownership             |
| `docs_filegroup`       | packageable documentation files              | useful, reachable, current documentation       |
| skill validation       | safe artifact shape and discovery links      | correct routing or successful outcomes         |
| Bazel query/build/test | declared graph or executed target result     | action authority or user-level acceptance      |
| goal validation        | internally consistent, digest-bound records  | truth of free-form claims or external delivery |
| delivery receipt       | exact prepared and published Git/forge state | which validation commands actually passed      |
| runtime tool result    | one bounded live observation                 | a durable or cross-provider system snapshot    |

## Current system gaps

### No bounded semantic system map

Path, policy, ownership, workspace, docs, skills, runtime packages, goals,
operations, and checks have separate discovery surfaces. The latent
`readme_tree` tool parses README frontmatter into an untyped map, omits the root
README, emits checkout-specific absolute paths in JSON, and has a documented
example that omits the required output type
([`tools/readme_tree`](../../../tools/readme_tree/)). No root entry
point consumes its output.

### Ownership and path policy are not machine-resolvable

`CODEOWNERS` now uses the intentional `* @simeonwarren` catch-all. More
specific ownership remains owner-local work rather than an inferred rewrite.

Tree policy also overloads `public`, `published`, and `used in builds` across
Bazel visibility, artifact publication, documentation, source disclosure,
production dependency, and secret or personal runtime content. The public
repository statement, non-public tree rules, third-party publication command,
and website documentation graph cannot be reconciled without supplying those
meanings.

### Operations expose labels, not effect contracts

`bazel_agent` intentionally passes ordinary Bazel commands and later options
through without validating their meaning. Phase 1 adds owner-local action
declarations for the goal store, Cordis mutation surface, repository delivery,
and Terraform map. The unnamed Terraform `apply` alias is removed and the
explicit `.apply` label remains. Other Bazel targets are not thereby claimed
to be runtime-enforced; Phase 3 owns admission and enforcement.

The same separation appears in dynamic runtime tools: handlers with different
filesystem, process, and network behavior have no common capability model.
Conversation authority, host sandbox permissions, tree policy, Bazel action
isolation, runtime power, credentials, and external authorization are separate
layers with no effective-authority preflight. A Bazel network default does not
constrain repository rules, `bazel run`, or MCP. Some credentialed diagnostics
can emit secrets or personal information before summary; optional scanners are
not universal pre-execution or delivery gates.

### Validation selection and evidence handoff remain manual

Agents select affected checks from policy and package inspection. Root quality
aliases, lint aspects, focused package tests, generated-file updaters, the
full-repository matrix, and delivery validation have no shared check taxonomy
or impact plan. Full-check reports record commands, exit status, duration, and
task-local raw-log paths, but not a common candidate/configuration identity or
versioned validation receipt. Those logs need review before reuse because they
may contain secrets or personal information; their task-local retention does
not make safe operational facts confidential.

Consequently successful evidence is hard to reuse. A later agent must often
rerun checks or reinterpret prose to determine whether a changed tree, commit,
base, environment, or configuration invalidated the earlier result.

### Runtime status and isolation stop at the worktree

The required MCP server awaits Cordis project and scratch initialization before
serving. A package that never settles can delay the health controls needed to
diagnose it
([`projects/mcp_cordis/cmd/mcp_cordis/main.mjs`][snapshot-mcp-main]
and [`internal/runtime.mjs`][snapshot-mcp-runtime]).

All runtimes in one worktree share the same scratch catalog and plugin path.
Mutation serialization is process-local, and generated scratch packages share
the main Node process and ambient runtime capabilities. The README correctly
states that this is a reliability boundary, not a security sandbox
([`projects/mcp_cordis/README.md`](../../mcp_cordis/), lines 99-128).

### Durable goals are integrity-rich but system-poor

Owner-local storage preserves autonomy, but there is no repository-wide goal
identity or discovery projection. A fresh agent must know a goals root before
it can list or analyze records. The durable schema cannot directly answer all
continuation questions required by its own skill protocol, including the
stable recurring defect, stale regression evidence, exact delivered subject,
or highest-leverage next action.

### Delivery and learning lose typed causal links

Delivery and goals do not share a validation or subject receipt. Rewrite
authority is split between rebase and delivery; version, remote refs, release
snapshot, and deployed artifact do not form one typed chain. Review disposition
and its causal regression link are not queryable, so repeated friction does not
automatically become a routing case, skill revision, test, or retirement.

### Topology and context costs recur

The eight standalone module paths are repeated in `.bazelignore`, root Bzlmod
overrides, project documentation aggregation, and the full-repository checker.
README metadata, docs membership, skill links, runtime packages, and goals have
similar independently maintained joins.

There is no cheap context slice that returns only the applicable policy,
owner, workspace, capability, checks, and current evidence for one path/task.
Agents either load broad prose and inventories or pay repeated search and
Bazel analysis costs to narrow them.

## Evidence boundary

This snapshot is based on checked-in source, configuration, tests, and
read-only queries. It does not establish deployed/host equivalence, external
provider health, runtime performance or isolation under load, unrun manual
evaluations, or repository-wide build/test success.

## Handoff to the system design

The target design can use this snapshot as its baseline. It needs to define,
without duplicating current authorities:

- linked layer contracts and authoritative generated projections;
- shared identity, effect, authority, evidence, and refusal envelopes;
- a minimal offline context/status view for a zero-context agent; and
- measurable correctness, safety, latency, context, and reuse signals.

Future-state claims belong in the architecture and roadmap, not in this
current-state snapshot.

[snapshot-agents]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/AGENTS.md
[snapshot-bazel-agent-main]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/bazel_agent/cmd/bazel_agent/main.go
[snapshot-bazel-preset]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/tools/bazelrc/preset.bazelrc
[snapshot-codeowners]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/CODEOWNERS
[snapshot-codex-config]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/.codex/config.toml
[snapshot-full-repo-check-skill]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/agents/skills/full-repo-check/SKILL.md
[snapshot-mcp-main]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/mcp_cordis/cmd/mcp_cordis/main.mjs
[snapshot-mcp-runtime]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/mcp_cordis/internal/runtime.mjs
[snapshot-process-supervisor]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/mcp_cordis/internal/process_supervisor.mjs
[snapshot-repo-context]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/mcp_cordis/plugins/repo_context.mjs
[snapshot-root-build]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/BUILD.bazel
[snapshot-root-readme]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/README.md
[snapshot-rules-docs-defs]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/rules_docs/docs/defs.bzl
[snapshot-site-build]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/projects/alwaldend.com/BUILD.bazel
[snapshot-terraform-defs]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/tools/terraform/defs.bzl
[snapshot-versioning-skill]: https://github.com/alwaldend/src/blob/1423dce5fab45ce5223caeb6a24791bf1a2cc3ff/tools/versioning/skills/versioning/SKILL.md
