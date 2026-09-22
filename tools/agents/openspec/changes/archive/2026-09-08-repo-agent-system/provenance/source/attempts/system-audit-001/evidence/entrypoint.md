# Zero-context entry-point and discoverability audit

## Audit identity

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
- Tracked source revision:
  `1423dce5fab45ce5223caeb6a24791bf1a2cc3ff`
- Scope: root entry points, repository navigation, public documentation,
  directory READMEs, and discoverability only.
- Mutation boundary: this report only; no maintained source or canonical goal
  resource was changed.

## Executive verdict

The repository has many useful local mechanisms, but it does not yet present
them as one legible system. Its highest-leverage defect is the absence of a
canonical, machine-readable semantic repository map with small human and agent
projections. Instead, topology, policy, capability, execution, documentation,
and ownership are independently encoded in root prose, tree READMEs, Bazel
files, README frontmatter, workspace manifests, generated skill links, and a
runtime MCP configuration. A zero-context agent must reconstruct the system by
searching and reconciling those surfaces.

This is not principally a documentation-volume problem. There are 409 tracked
`README.md` files, but 304 (74.3%) are at most twelve lines and predominantly
frontmatter. The human root entry point is only 299 bytes, while the sole agent
instruction file is 12,375 bytes and 1,704 words. The system is simultaneously
under-explained at its entrance and over-concentrated in one always-loaded
policy document.

The best design is therefore not a larger root manual. It is a linked control
plane:

1. one canonical system model and typed topology;
2. thin root entry points that link into it;
3. generated, bounded projections for humans and agents;
4. path- and task-specific context assembly;
5. validation that rejects drift; and
6. measured feedback that can promote recurring discoveries into the model.

No task-defeating technical premise was found. The only non-verifiable premise
is the literal claim of being "maximally" ergonomic: structural review alone
cannot prove an optimum. It must be operationalized with representative
zero-context tasks and measured correctness, calls, bytes/tokens, elapsed
time, and unsafe-action rate.

## Current entry-point graph

### Human entry

`README.md` contains only three external links and the license
(`README.md:6-14`). It does not explain the six top-level ownership trees,
identify Bazel as the execution surface, link `AGENTS.md`, point to local setup
source, or expose repository-wide health and navigation commands. The setup
link leaves the checkout for the website even though its source is available
at `projects/alwaldend.com/content/docs/misc/repo.md`.

The website is a projection of repository Markdown rather than an independent
documentation corpus:

- `BUILD.bazel:92-100` aggregates documentation from root subpackages.
- `projects/alwaldend.com/BUILD.bazel:105-124` includes `//:docs` in the site
  source and renames README basenames to Hugo section indexes.
- `projects/rules_docs/docs/defs.bzl:24-53` packages source files under
  `content/docs/<package>` without semantic validation.

That is a valuable single-source mechanism, but the root README does not tell
a local reader how this projection works or where the website pages originate.

### Agent entry

`AGENTS.md` is the only tracked `AGENTS.md`. It combines repository
classification, scratch policy, decision policy, question routing, layout,
tooling, documentation, search, tree topology, Bazel, infrastructure safety,
style, and verification in 239 lines. There is no layered instruction chain
below it, despite the file explicitly allowing nested overrides
(`AGENTS.md:5-8`).

The document is strong on invariants and safety but weak as a navigation hub:
it has no direct link to a system architecture, skill catalog, goal catalog,
runtime capability catalog, documentation source, or repository health entry
point. Its main navigation instruction asks the agent to read the nearest
README, BUILD file, and optional module fragment (`AGENTS.md:70-75`), leaving
the agent to locate, interpret, and reconcile each one.

### Capability entry

Skill ownership and discovery are comparatively well designed:

- `projects/agents/README.md:14-23` distinguishes canonical skill directories
  from `.agents/skills/` discovery symlinks.
- `BUILD.bazel:66-90` owns the explicit generated skill-link set.
- Twenty discovery symlinks currently route to eighteen repository-agent
  skills, the goal skill, and the versioning skill.

However, neither root entry point offers a capability table or task-to-skill
router. The mechanism is discoverable only after inspecting
`projects/agents/README.md`, the BUILD target, or runtime-provided skill
metadata.

The workspace also has a useful bounded context substrate:

- `.codex/config.toml:1-10` registers a required workspace MCP server.
- `projects/mcp_cordis/README.md:138-145` lists `repo_context` as a starter
  package selected from recurring task categories.
- `projects/mcp_cordis/plugins/repo_context.mjs:674-779` implements a bounded
  path description with Git state, applicable `AGENTS.md` files, and directory
  entries.

It does not yet assemble the semantic context that `AGENTS.md` asks an agent
to reconstruct: applicable ownership, nearest README and BUILD files, nested
workspace boundary, relevant skills, command risk, or verification profile.
Its default instruction inclusion can also return the 12 KB root `AGENTS.md`
again after a host has already injected it.

### Directory entry

The six tree READMEs state broad eligibility constraints, but most offer no
navigation, owner, common operation, validation entry point, or child index:

- `projects/README.md:10-16`
- `infra/README.md:10-16`
- `tools/README.md:10-24`
- `data/README.md:10-16`
- `third_party/README.md:10-21`
- `users/README.md:10-16`

The global guide says these files are authoritative for their trees
(`AGENTS.md:142-143`) while duplicating their constraints immediately above
(`AGENTS.md:124-140`). Thus the declared authorities are not the only copies.

### Generated documentation entry

The repository has a latent README indexer, but it is not an effective entry
point:

- `tools/readme_tree/main/go/parser.go:171-188` parses arbitrary frontmatter
  into an untyped map.
- `tools/readme_tree/main/go/template-markdown.txt:1-3` renders only nested
  directory links and titles; it cannot route by owner, capability, risk,
  target, or status.
- `tools/readme_tree/README.md:11-15` documents a direct `bazel` invocation
  without an output type. The output code rejects every value except explicit
  `json` or `markdown`, including the empty default
  (`tools/readme_tree/main/go/outputter.go:29-37`).
- The only repository integration found makes the binary available in bzlenv
  (`tools/bzlenv/BUILD.bazel:5-37`); no root projection consumes its output.

`tools/repo_map` is not a repository navigation map at all. Its README defines
it as a platform-dependent repository-download module extension
(`tools/repo_map/README.md:1-10`). The name is a false affordance for an agent
searching for repository topology.

## Authority and duplication matrix

| Mutable fact | Claimed or effective authority | Competing copies or projections | Consequence |
| --- | --- | --- | --- |
| Tree publication and dependency policy | Top-level tree README, by `AGENTS.md:142-143` | Repeated in `AGENTS.md:124-140`; partly encoded by Bazel visibility | An agent must decide whether prose or graph wins when they differ. |
| Agent Bazel invocation | `AGENTS.md:145-179` and the `bazel-agent` skill | Human repo, infra, third-party, and readme-tree docs embed direct `bazel` commands | Commands copied into an agent workflow bypass the required runner. |
| Documentation membership | Bazel `docs_filegroup` graph | README presence/frontmatter and Hugo layout imply a second conceptual tree | Coverage is buildable but not semantically indexed or link-checked. |
| Skills | Canonical `SKILL.md` plus root `write_skill_links` input | `.agents/skills` is a generated projection | This is mostly sound, but no human/agent catalog exposes the graph. |
| Runtime tools | `.codex/config.toml`, Cordis config, and live tool catalog | Project README describes the mechanism | Required startup capability has no root ownership, health, or fallback link. |
| Repository ownership | Normal paths are effectively unassigned: `CODEOWNERS:1` uses the literal pattern `-`; tree policy READMEs | Package names, documentation categories, and skill owners imply other notions of ownership | "Who owns this change?" has no machine-answerable semantic result. |
| Workspace topology | `MODULE.bazel`, `.bazelignore`, `go.work`, `pnpm-workspace.yaml`, docs deps | Each consumer maintains its own membership list | A path can be a root package, nested module, docs subtree, and language workspace with no unified view. |
| Plans and durable learning | Project-local goal records and goal tooling | No root goal index or system-model link exists | Prior decisions and failures are searchable only if the agent already knows where to look. |

## Strongest counterargument and response

The strongest case for the current design is that the repository already uses
good Unix-like composition: root policies, nearest READMEs, Bazel targets,
canonical skills with generated discovery links, a public docs projection, and
bounded runtime tools. A central model could become a second stale database
and add bureaucracy.

That objection is valid against a hand-maintained omnibus catalog. It does not
defend the present state: broken authoritative links, contradictory
classification language, a non-working indexer example, duplicated boundary
facts, and absent root routes are direct drift evidence. The correct response
is not central duplication. It is a small canonical contract plus generated
projections from existing authorities, with validation that fails on drift.

## Task-defeating premise

No implementation blocker or impossible dependency was found. The goal should
reject only one interpretation: no finite design-document change can prove
global maximality for all future agents and workloads. Treat "maximally
agent-intuitive" as an optimization loop with explicit workload samples and
resource budgets. Under that interpretation, the requested repository-level
system design is feasible and the existing Bazel, skill, goal, documentation,
and `repo_context` mechanisms provide strong reusable foundations.

## Ranked findings

### 1. P0: classification and publication terms are internally ambiguous

The root policy first says all checked-in source, documentation, and fixtures
are public information (`AGENTS.md:10-20`), then says all files under `infra/`
and `users/` are sensitive (`AGENTS.md:181-203`). The tree READMEs say those
trees "MUST NOT be public" or published (`infra/README.md:12-16` and
`users/README.md:12-16`), even though the repository is public and their
documentation is recursively collected into the website through the root docs
graph.

These statements may intend different axes—Git disclosure, Bazel visibility,
artifact publication, production dependency eligibility, and operational
sensitivity—but the vocabulary does not say so. This forces agents either to
over-restrict harmless checked-in material or expose something the stricter
rule meant to protect. Security ambiguity is both high consequence and common.

Improvement: define those axes once in a typed topology schema. State an
explicit precedence rule such as "path sensitivity overrides repository
disclosure defaults." Generate the human tree table and agent policy summary
from it; let tree READMEs explain rationale without restating mutable values.

### 2. P0: there is no canonical semantic map or one-hop root route

The human root gives only external destinations (`README.md:6-10`); the agent
root gives policies and ad hoc discovery commands. Neither answers the first
questions a zero-context driver has:

- What is this system for?
- Which tree owns this kind of change?
- Is this path a root package or a nested workspace?
- Which capability should execute the work?
- What may mutate local or remote state?
- What is the narrow validation and delivery path?
- Where are current plans and prior evidence?

The information exists in fragments but not as a linked abstraction tower.
Every task pays reconstruction cost and can reach a different conclusion.

Improvement: make a maintained repository-system architecture and a small,
machine-readable topology/capability manifest canonical. The architecture
should link intent -> topology -> policy -> capability -> execution ->
evidence -> delivery -> durable learning. Root README and AGENTS should each
link it in one hop and render only stable summaries.

### 3. P1: context is globally expensive and locally sparse

All scopes load the 1,704-word root guide because there are no nested agent
guides. The subsequent "read nearest README/BUILD/module" path is often low
yield: 304 of 409 tracked READMEs are twelve lines or shorter. There are 506
BUILD files and 73 `include.MODULE.bazel` files, so blind traversal is not a
cheap substitute for a semantic map.

The existing `repo_context_get` is a strong substrate but stops at filesystem,
Git, and AGENTS context. It should evolve rather than be replaced.

Improvement: add a bounded `context_for(path, task)` projection that returns
document paths and digests, not unbounded prose: applicable policy sections,
owning tree/project, nested-workspace root, canonical README/BUILD/module,
matched skills, safe discovery commands, mutation class, and validation
targets. Keep a small root policy kernel; add nested AGENTS files only where a
subtree truly changes global rules.

### 4. P1: documentation aggregation proves packaging, not usefulness

`rules_docs` copies files and dependencies but does not validate frontmatter,
local links, command labels, ownership metadata, or freshness
(`projects/rules_docs/docs/defs.bzl:5-53`). Repository quality formats Markdown
with Prettier (`tools/repo_quality/BUILD.bazel:186-269`) but its test suite
contains formatting and language checks, not link or semantic documentation
checks (`tools/repo_quality/BUILD.bazel:299-307`).

Observed drift confirms the gap. Every example link in the authoritative
`infra/README.md:18-35` points under removed `infra/dc1/...` paths; direct
existence checks found all four targets absent. Two documented apply labels
also use the removed `infra/dc1` topology.

Improvement: validate local links, required frontmatter schema, generated
indexes, referenced Bazel labels, and command risk annotations in the normal
repository quality target. A README count should not be treated as knowledge
coverage; route metadata should be measurable separately from prose presence.

### 5. P1: operational recipes mix audiences, risk levels, and stale commands

The public repo setup page mixes bootstrap, full-repository builds, Vault
login, package-lock mutation, host certificate mutation, broad mechanical
rewrites, and historical cleanup recipes in one flat page
(`projects/alwaldend.com/content/docs/misc/repo.md:8-171`). It uses direct
`bazel` throughout, including `//...` and stateful targets, while the agent
policy requires `bazel_agent` and narrow targets (`AGENTS.md:145-179`).
`infra/README.md:18-35` places unguarded `tf.apply` commands directly in the
tree entry point, and `third_party/README.md:18-21` exposes a publishing command
without a risk or authority label.

Human commands are not inherently wrong, but the pages do not declare their
audience or mutation class. An agent can reasonably copy a command from a
root-linked page and violate the agent contract.

Improvement: separate safe bootstrap and read-only discovery from operator
runbooks. Attach structured classes to every canonical workflow:
`inspect`, `build`, `test`, `generate-source`, `local-mutate`, or
`external-mutate`. Generate command snippets from target metadata where
possible and render agent/human variants explicitly.

### 6. P1: existing capability discovery is mechanically sound but semantically hidden

Generated skill symlinks avoid duplicate canonical content, and the Cordis
gateway supports bounded repository reads. Those are good foundations. The
root surfaces nevertheless omit both systems. A zero-context agent must know
to inspect `.agents/skills`, discover gateway handlers, or search project
READMEs before it can use the capabilities intended to save work.

Improvement: generate a compact capability catalog from skill frontmatter,
BUILD ownership, and runtime package descriptors. It should map task intents
to capability, authority requirements, input/output contracts, cost class, and
fallback. Root documents link the catalog; `repo_context` returns only matched
entries.

### 7. P2: ownership and durable learning have no repository-wide navigation

`CODEOWNERS:1` uses the literal pattern `-`, not a catch-all, so ordinary
repository paths have no effective CODEOWNERS assignment. The README hierarchy
encodes content categories but not accountable owners, system interfaces, or
review routes. Goal records exist in projects, yet neither root entry point nor
`projects/agents/README.md` describes where repository-wide goals live, how to
find active ones, or how accepted findings become policy, skills, or checks.

Improvement: distinguish review ownership from architectural ownership. Add
owner and interface references to the topology model, generate CODEOWNERS only
where appropriate, and publish a bounded goal/decision index. A recurring
failure should link from goal evidence to the policy, capability, or validation
change that absorbed the lesson.

### 8. P2: names and dormant tools create false navigation affordances

An agent searching for "repo map" finds a dependency-download extension, not
topology. An agent finding `readme_tree` gets an example that omits a required
output type and a result that contains only titles. Both increase search calls
while appearing to solve the problem.

Improvement: reserve names such as `repo_catalog`, `repo_context`, and
`repo_health` for system-control surfaces. Clarify or rename unrelated tools,
and either integrate `readme_tree` as a generated projection of the canonical
model or retire it from the apparent navigation path.

## Recommended dependency order

1. **Define the system contract.** Establish terminology, abstraction layers,
   authority rules, and the distinction between canonical records and derived
   projections. Do not begin by rewriting hundreds of READMEs.
2. **Define typed topology and capability records.** Include ownership,
   workspace kind, disclosure, publication, dependency eligibility,
   sensitivity, mutation classes, canonical docs, commands, and validation.
3. **Make both root files thin one-hop routers.** Keep safety invariants in
   `AGENTS.md`; put repository identity, a six-tree map, and safe quick start in
   `README.md`. Link the same architecture and generated catalog from both.
4. **Extend the existing context runtime.** Have `repo_context` query the
   canonical records and Bazel metadata to return task/path-specific bundles
   with explicit truncation and digests.
5. **Generate projections.** Produce the website tree, skill catalog, goal
   index, workflow snippets, and optional README sections from their owning
   records. Generated content must identify its source and update target.
6. **Add drift gates.** Check local links, Bazel labels, frontmatter schema,
   duplicate authorities, generated projections, risk annotations, and
   architecture links in the normal repository quality path.
7. **Measure zero-context operation.** Maintain small scenario evaluations for
   locating an owner, choosing a skill, finding a target, avoiding an unsafe
   operation, validating a change, and resuming a failed goal. Track success,
   tool calls, context bytes/tokens, wall time, and incorrect mutations.
8. **Accrete only proven lessons.** Promote repeated, evidenced patterns from
   goals into the topology, a skill, a validation gate, or a decision record;
   do not grow the root prompt with every local exception.

## Explicit contradictions and drift

1. **Disclosure:** `AGENTS.md:12-16` classifies checked-in source as public;
   `AGENTS.md:187-189` classifies all `infra/` and `users/` files as sensitive.
2. **Authority:** `AGENTS.md:142-143` declares tree READMEs authoritative but
   `AGENTS.md:124-140` repeats their mutable policies.
3. **Tool visibility:** `tools/README.md:14-18` says tools must be visible to
   the whole repo; `AGENTS.md:131-134` weakens this to "normally" for tool
   targets. The intended exceptions are not modeled.
4. **Publication:** infra/users READMEs say their trees must not be published,
   while their Markdown is included by the recursive docs graph and then by
   the public site. This may be an axis mismatch, but it is not stated as one.
5. **Bazel entry:** agents must use `bazel_agent` (`AGENTS.md:147-151`), while
   the root-linked repo page and top-level READMEs show direct `bazel` commands
   without audience labels.
6. **Infra navigation:** all four example links in `infra/README.md:20-35` are
   absent at the audited revision, and two commands reference removed paths.
7. **README indexer:** the documented invocation omits the output type that its
   implementation requires (`tools/readme_tree/README.md:11-15` versus
   `tools/readme_tree/main/go/outputter.go:29-37`).
8. **Map naming:** `tools/repo_map` sounds like topology discovery but owns an
   unrelated module-extension concern (`tools/repo_map/README.md:1-10`).
