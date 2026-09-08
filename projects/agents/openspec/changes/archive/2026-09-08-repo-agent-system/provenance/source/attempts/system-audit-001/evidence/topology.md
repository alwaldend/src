# Repository topology and metadata audit

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
- Scope: repository topology and metadata only: ownership boundaries,
  README frontmatter, documentation aggregation, CODEOWNERS, README indexing,
  project lifecycle metadata, workspace boundaries, and map-query cost.
- Mutation boundary: this report only. No maintained source or canonical goal
  resource was changed.

## Executive verdict

An agent can derive a cheap *filesystem sketch* of this repository, but it
cannot cheaply derive a *trustworthy system map*.

The structural foundation is unusually good. Six named content trees have
clear prose policies, almost every immediate component has a README and BUILD
file, documentation is recursively aggregated, and standalone Bazel modules
are explicitly integrated. The failure is at the join points. Path policy,
human ownership, documentation membership, lifecycle, workspace identity,
runtime control surfaces, and executable targets live in different formats
with no typed model, common validator, or bounded projection.

The highest-consequence concrete defect is ownership: `CODEOWNERS:1` contains
the single path pattern `-`, not the catch-all pattern `*`. Normal repository
paths therefore do not resolve to the listed owner. The highest-frequency
defect is navigation: neither `README.md` nor `AGENTS.md` links a semantic
catalog, while the latent `readme_tree` entry point is broken as documented,
untyped, non-portable, and comparatively expensive to invoke.

The right correction is not another manually maintained central registry.
Keep facts with their natural authorities and generate one deterministic,
checked projection that joins them:

```text
Git paths and READMEs       -> identity, description, lifecycle, tags
top-level README policy     -> boundary rules on explicit policy axes
CODEOWNERS                  -> accountable human owner
AGENTS.md chain             -> applicable agent policy
MODULE.bazel / BUILD.bazel  -> workspace and executable/package identity
docs_filegroup              -> documentation publication membership
                                |
                                v
                 validated semantic catalog
                    /                     \
        bounded Markdown map         portable JSON map
```

That gives agents a cheap read path while preserving Bazel as the execution
graph, README files as co-located explanations, and CODEOWNERS as the review
platform's ownership authority.

## Preserved strengths

These properties should be preserved rather than replaced.

1. `AGENTS.md:124-143` defines six stable top-level content boundaries and
   explicitly delegates their detailed policy to each tree's README.
2. The six tree READMEs all have title, description, site weight, and a Hugo
   category cascade. Their simple path taxonomy is easy to infer.
3. Of 182 immediate child directories under `projects/`, `infra/`, `tools/`,
   `data/`, `third_party/`, and `users/`, 181 have both `README.md` and
   `BUILD.bazel`. The only exception is `tools/bazel_contracts`.
4. All 28 project roots have a README, BUILD file, title, and description.
5. Every tracked README directory that also has a BUILD file contains a
   textual `docs_filegroup(...)` declaration. This is a useful convention,
   even though it does not prove reachability or semantic validity.
6. A focused `bazel_agent query '//:docs' --output=build` resolves the root
   documentation aggregate to exactly the six content trees. A focused query
   of `//projects:docs` resolves all 28 project roots: 20 local packages and
   eight standalone modules imported through apparent repository names.
7. The nine tracked `MODULE.bazel` files make the root plus eight standalone
   Bazel workspaces mechanically discoverable. The 73
   `include.MODULE.bazel` files are separately recognizable as root-module
   fragments rather than standalone workspaces.
8. `projects/rules_docs_gazelle/gazelle/gazelle.go:64-80` reliably creates a
   conventional `docs` target when a Bazel package already has a README.

## Current topology and authorities

### Content and policy boundaries

The path taxonomy is clear but incomplete as a system taxonomy.

| Path | Prose role | Immediate components | Machine authority |
| --- | --- | ---: | --- |
| `projects/` | product and reusable code | 28 | path and Bazel packages |
| `infra/` | infrastructure definitions | 14 | path and Bazel packages |
| `tools/` | repository-wide tools | 97 | path; 96 documented packages |
| `data/` | private data and docs assets | 9 | path and Bazel packages |
| `third_party/` | vendored/external code | 33 | path and Bazel packages |
| `users/` | user-specific code and infra | 1 | path and Bazel packages |

The tracked control plane is outside that table:

- `.agents/skills/` contains 20 generated discovery symlinks.
- `.codex/config.toml:1-10` registers the required workspace MCP runtime.
- root configuration owns Bazel, language workspaces, quality entry points,
  and skill-link generation.
- `out/` is the ignored task-state and evidence plane.

Neither root entry point presents content and control planes together.
`README.md:6-14` has only external links and the license. `AGENTS.md` documents
the six content trees but does not map `.agents/`, `.codex/`, nested
workspaces, component ownership, or the documentation projection.

### Boundary-policy ambiguity

The top-level READMEs use the words `public`, `published`, and `used in
builds` as if each named one boolean. The graph shows that several distinct
axes are involved:

- checked-in source disclosure;
- Bazel target visibility;
- publication of binaries, modules, charts, or other artifacts;
- documentation publication;
- eligibility as a production dependency; and
- sensitivity of generated or runtime values.

Those meanings are not defined. For example:

- `AGENTS.md:10-20` says checked-in source and documentation are public
  information, while `infra/README.md:14-16` says infra must not be public or
  published.
- `third_party/README.md:14-16` says the tree must not be published, then
  lines 18-21 document a vendored Helm publication command.
- `BUILD.bazel:92-100` aggregates all six content trees, including infra,
  data, third party, tools, and users. The aggregate is included in the Hugo
  website source at `projects/alwaldend.com/BUILD.bazel:113-124`.

These may all be intentional if `public` means Bazel visibility and
`published` means a production artifact rather than documentation. An agent
cannot know that from the declared authority. A typed boundary policy must
name each axis and its exceptions.

### README metadata

The tracked README inventory is broad but its semantics are weakly typed.

| Measure | Result |
| --- | ---: |
| tracked `README.md` files | 409 |
| first-line frontmatter delimiter | 407 |
| frontmatter `title` | 407 |
| frontmatter `description` | 383 |
| frontmatter `tags` | 239 |
| frontmatter `languages` | 79 |
| frontmatter `statuses` | 17 |

The two README files without frontmatter are maintained legacy goal
projections under `projects/mcp_cordis/goals/runtime_extensions/`. Missing
frontmatter is therefore not automatically an error; document class matters.

Project-root metadata is more important for topology:

- 28 of 28 project roots have title and description.
- 17 of 28 have `statuses`; 11 have no project lifecycle value.
- 26 of 28 have `tags`; `projects/cgit` and `projects/mcp_cordis` do not.
- Existing status values are `in_progress` (5), `maintenance` (4),
  `finished` (4), `active` (2), `experimental` (1), and `abandoned` (1).

Nothing found in repository tests validates README frontmatter keys, required
fields, cardinality, or vocabulary. Hugo treats statuses and tags as
taxonomies (`projects/alwaldend.com/hugo.toml:57-61`). The documentation
Gazelle extension checks only that `README.md` exists and emits
`glob(["*.md"])`; it does not parse metadata
(`projects/rules_docs_gazelle/gazelle/gazelle.go:64-80`).

Tags should remain advisory search metadata. They are too incomplete and
uncontrolled to authorize actions or route correctness-critical work.
Project status should be typed because agents use lifecycle to decide whether
to extend, preserve, migrate, or avoid a component.

### Documentation graph

`docs_filegroup` is a sound packaging primitive, not a semantic catalog.
`projects/rules_docs/docs/defs.bzl:24-53` creates a `pkg_files` target for
local Markdown and a `pkg_filegroup` that aggregates child documentation. It
does not expose title, status, owner, workspace, policy, capability, or risk.

The graph has four material properties:

1. `//:docs` has clean direct edges to the six content boundaries.
2. The root deliberately sets `srcs = []` (`BUILD.bazel:92-100`), so the root
   README is not part of its own documentation aggregate.
3. `//projects:docs` replaces eight locally ignored standalone modules with
   apparent external repository labels (`projects/BUILD.bazel:4-22`). That is
   correct for packaging, but makes repository topology depend on Bzlmod
   resolution if an agent queries the transitive graph.
4. Generated documentation may depend on source generators and their tools.
   A transitive dependency query therefore expands far beyond documentation.

The observed query demonstrates the cost. A cold
`bazel_agent query 'deps(//:docs)' --output=label_kind` loaded approximately
1,649 packages, attempted to enumerate more than five megabytes of labels,
crossed into external code and tool dependencies, and failed after more than
100 seconds while fetching an OCI manifest from the configured registry. A
warm direct `//:docs --output=build` query completed in about 2.5 seconds.

The lesson is not that the docs graph is bad. It is that dependency closure
is the wrong repository-map API. Map generation must read small source
authorities and validate only direct documentation reachability.

### Standalone workspace topology

The eight nested workspace paths are repeated in at least four operational
surfaces:

1. `.bazelignore:4-11` excludes them from the root workspace.
2. `third_party/include.MODULE.bazel:79-127` declares their dependencies and
   local path overrides.
3. `projects/BUILD.bazel:8-21` excludes local rule packages and reinserts
   their external `:docs` targets.
4. `projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go:19-38`
   hard-codes the root and eight nested workspaces.

The `bazel-nested-module` skill correctly instructs authors to update all of
these surfaces, but prose coordination is not an invariant. A ninth nested
workspace can build locally yet be absent from full-repository validation or
documentation. `MODULE.bazel` presence should define workspace identity; the
other lists should be generated from, or validated against, that fact.

### CODEOWNERS

`CODEOWNERS` has one line:

```text
- @simeonwarren
```

In CODEOWNERS syntax the first token is a path pattern. `-` is not the
repository catch-all `*`; it describes a literal dash-named path. The file
therefore does not assign the listed owner to ordinary repository paths.

Even after correcting the catch-all, one global owner answers only review
accountability. It does not identify the narrowest semantic component. The
catalog should resolve the effective CODEOWNER for every component while
deriving component identity from path/README/BUILD structure. Do not duplicate
owners into README frontmatter unless the repository intentionally introduces
a different owner role with different semantics.

### `readme_tree`

`readme_tree` is a useful prototype, but it is not a trustworthy map surface.

- Its README example omits `--output-type`
  (`tools/readme_tree/README.md:11-15`). The default is the empty string, and
  the outputter accepts only explicit `json` or `markdown`
  (`tools/readme_tree/main/go/cmd.go:72-124` and
  `tools/readme_tree/main/go/outputter.go:29-37`). Running the documented
  command through `bazel_agent` reproduced `output type  is not supported`.
- The CLI accepts `template` as a valid enum
  (`tools/readme_tree/main/go/model.go:9-31`), but the outputter has no template
  branch and rejects it.
- The parser intentionally omits the root README
  (`tools/readme_tree/main/go/parser.go:112-120`).
- Frontmatter is split on every unanchored `---` sequence and decoded into
  `map[string]any`, with no schema
  (`tools/readme_tree/main/go/parser.go:161-188`).
- JSON contains checkout-specific absolute `dir`, `path`, and directory
  fields (`tools/readme_tree/main/go/model.go:40-53`). A corrected single-file
  invocation reproduced those absolute values.
- Markdown renders only path indentation, directory name, and title
  (`tools/readme_tree/main/go/template-markdown.txt:1-3`). It cannot answer
  owner, boundary, workspace, status, policy, docs reachability, or execution
  questions.
- Repository search found no consumer outside its own package except a bzlenv
  exposure. `tools/repo_map` is unrelated: it is a platform-dependent external
  repository downloader, so its name is a false affordance for this task.
- `bazel_agent query 'tests(//tools/readme_tree/...)'` found only shellcheck and
  shfmt tests for the legacy shell implementation. There are no behavioral
  tests for the Go parser/output contract.
- The documented run configured 46,077 targets across 308 packages and took
  about eight seconds in this checkout before failing at the CLI boundary.

This tool should either become the tested catalog generator or be retired in
favor of one. Keeping an unintegrated second map implementation increases
search noise and false confidence.

## Ranked gaps

### 1. P0: effective human ownership is absent

**Evidence:** `CODEOWNERS:1` uses literal pattern `-`.

**Consequence:** change routing, review accountability, escalation, and
component stewardship cannot be derived even though an owner name appears to
be present. This is a high-consequence false-positive affordance.

**Minimum correction:** replace the literal pattern with a real catch-all,
add narrower path rules only where responsibility differs, and validate that
every catalog component resolves to at least one CODEOWNER.

### 2. P0: boundary policy collapses distinct authority axes

**Evidence:** the public-repository statement, tree `public`/`published`
requirements, the third-party publish command, and the public docs aggregate
cannot be reconciled without guessing what each term means.

**Consequence:** an agent may either over-restrict harmless work or cross a
real publication, dependency, or sensitivity boundary.

**Minimum correction:** encode Bazel visibility, artifact publication,
documentation publication, production dependency eligibility, and sensitive
value handling as separately named fields in top-level README frontmatter.
Render prose summaries from those values or validate the prose against them.

### 3. P0: no one-hop, bounded semantic map exists

**Evidence:** root README has no local topology route; root AGENTS maps only
six content trees; docs closure is not a map; `readme_tree` is broken and
unintegrated; hidden control surfaces and workspace roots are omitted.

**Consequence:** every task spends searches, tool calls, context tokens, and
judgment to reconstruct the same topology, and two agents can reconstruct
different systems.

**Minimum correction:** generate a portable JSON catalog and bounded Markdown
map from existing authorities. Link the same model directly from `README.md`
and `AGENTS.md`. Reading the projections must require no Bazel invocation.

### 4. P1: standalone workspace membership has four mutable copies

**Evidence:** `.bazelignore`, root Bzlmod overrides, project docs aggregation,
and full-repo-check each carry the same eight paths.

**Consequence:** adding or renaming a module can silently omit documentation
or complete validation even if focused builds pass.

**Minimum correction:** discover workspace roots from tracked `MODULE.bazel`
files, record module name/path once, and make a repository test compare every
required mirror. Generate mirrors where their consumers permit it.

### 5. P1: project lifecycle metadata is partial and unvalidated

**Evidence:** 11 of 28 project roots lack `statuses`; the six observed values
have no repository-defined semantics or cardinality check.

**Consequence:** agents cannot reliably distinguish actively developed,
stable, maintenance-only, complete, experimental, and abandoned work.

**Minimum correction:** define one required status per project root with a
small documented enum and transition meanings. Keep tags optional and
advisory. Add an actionable metadata validation test.

### 6. P1: docs reachability and metadata validity are independent

**Evidence:** Gazelle creates a docs target from README presence but neither
parses frontmatter nor wires semantic parent relationships. The packaging
macro carries files and deps only.

**Consequence:** a README can be well formed but unreachable, reachable but
semantically invalid, or embedded through a custom parent without appearing
as its own package. Bazel success proves packaging, not map truth.

**Minimum correction:** keep `docs_filegroup` narrow and add catalog checks for
frontmatter schema, component coverage, direct docs reachability, and link
validity. Do not overload transitive Bazel dependency closure as metadata.

### 7. P2: one top-level tool component is outside every convention

**Evidence:** `tools/bazel_contracts/worker_protocol/go_mod_tidy_fix.go` is a
tracked source file under the tools boundary, but `tools/bazel_contracts` has
no README or BUILD file and is absent from `//tools:docs`.

**Consequence:** it has no declared purpose, owner route, queryable target,
documentation, or validation surface. The otherwise excellent 181/182
component-root coverage makes this an actionable drift signal.

**Minimum correction:** move the compatibility type to its actual owner or
make `tools/bazel_contracts` a conforming component. Make the catalog coverage
test prevent recurrence.

### 8. P2: the current README indexer is a false affordance

**Evidence:** documented invocation failure, unsupported advertised output,
absolute machine output, root omission, no schema, no consumer, and no Go
behavior tests.

**Consequence:** an agent reasonably finds a tool named for the task, spends a
Bazel analysis, and receives either an error or an incomplete map.

**Minimum correction:** use it as the seed for the catalog generator only if
its data model, portability, tests, and integration are replaced. Otherwise
deprecate it and point its README at the canonical map.

## Proposed minimal authoritative metadata model

### Principle: distribute authority, centralize projection

The repository should not introduce one hand-edited file that repeats every
component. Existing co-located facts are more accretive: adding a normal
README and BUILD file already follows local ownership. The new system should
define which source owns each fact, then generate the join.

| Fact | Canonical source | Rule |
| --- | --- | --- |
| node identity | tracked repository-relative path | never hand-assign another ID |
| title and summary | README frontmatter | required for catalog components |
| boundary policy | top-level README frontmatter | typed, separate policy axes |
| project lifecycle | project-root README `statuses` | exactly one validated value |
| tags/languages/sites | README frontmatter | optional advisory facets |
| accountable owner | `CODEOWNERS` | every component must resolve |
| agent policy | applicable `AGENTS.md` chain | paths and digests in projection |
| workspace | nearest tracked `MODULE.bazel` | root and module name derived |
| Bazel package | nearest `BUILD.bazel` | label derived, existence validated |
| documentation | direct `docs` target | direct reachability validated |
| control surface | `.agents`, `.codex`, root config | explicit catalog node class |
| hierarchy | repository path containment | deterministic, no manual child list |

Only facts that cannot be inferred should be added to frontmatter. In
particular, do not repeat path, workspace root, Bazel package, owner, or parent
in every README.

### Derived catalog record

A generated record needs no more than the following shape:

```yaml
path: projects/agents
title: Agents
description: Repository-owned agent skills and evaluation assets
nodeClass: component            # derived from path/README/BUILD convention
boundary: projects              # derived from first path segment
status: active                  # canonical README value
tags: [agent, skills]           # advisory README values
workspace:
  root: .                       # nearest MODULE.bazel
  module: com_alwaldend_src     # parsed from that module
bazel:
  package: //projects/agents    # derived
  docs: //projects/agents:docs  # validated direct target
owners: ["@simeonwarren"]       # resolved from CODEOWNERS
policyChain:
  - AGENTS.md
  - projects/README.md
```

The actual JSON projection should contain only repository-relative paths,
stable sort order, schema version, tracked-source revision or input digest,
generation command, and an explicit generated marker. It must not contain
checkout paths, environment values, credentials, inferred secrets, or the
full Bazel dependency graph.

### Node classes

A small derived classification is enough:

- `repository`: the root entry and global control plane;
- `boundary`: the six content trees;
- `component`: an owning README/BUILD root;
- `workspace`: any component with its own `MODULE.bazel`;
- `package`: deeper README/BUILD documentation nodes; and
- `document`: intentionally non-package Markdown such as goal projections.

Control-plane paths such as `.agents/` and `.codex/` should be explicit nodes
owned by the repository/agents component, not forced into a content boundary.

### Projections and query contract

Produce two checked projections from the same catalog:

1. a bounded Markdown overview for humans and zero-context agents; and
2. portable JSON for path, owner, workspace, status, and policy lookup.

The root README and AGENTS file should link the Markdown model and JSON
contract in one hop without copying mutable detail. The detailed projection
should support cheap prefix selection so an agent can read one component and
its policy chain without loading all 409 READMEs.

Regeneration may use a focused repository-owned Go target, but ordinary
orientation must read the checked projection directly. A `catalog.check`
target should regenerate in an isolated output directory and fail with the
smallest actionable diff. It must not contact the network or traverse
`deps(//:docs)`.

### Required invariants

The model becomes trustworthy only with executable checks:

1. every immediate child of a content boundary is either a valid component or
   an explicit, reasoned exception;
2. every component has title, description, README, BUILD, direct docs target,
   and effective CODEOWNER;
3. every project root has exactly one valid lifecycle status;
4. every workspace is discovered from `MODULE.bazel`, has a unique path and
   module name, and is reconciled with ignore, root integration,
   documentation, and full-check coverage;
5. every catalog path is repository-relative and resolves inside the tracked
   tree;
6. direct documentation membership is complete without using transitive
   source/tool closure;
7. generated Markdown and JSON are deterministic and match their input digest;
8. root README and AGENTS links resolve to the same canonical system model;
9. policy fields use explicit axes and do not rely on overloaded `public` or
   `published`; and
10. catalog generation has a fixed item/byte bound and a measured offline
    latency budget.

### Lowest-risk migration order

1. Correct CODEOWNERS and add coverage validation.
2. Define the small frontmatter schema and explicit boundary policy axes.
3. Complete the 11 missing project statuses and classify the one orphan tool.
4. Implement the offline catalog join and tests, initially as an observation
   tool with no build-graph rewrites.
5. Check in deterministic Markdown/JSON projections and link them from both
   root entry points.
6. Reconcile nested workspace mirrors against module discovery.
7. Retire or absorb `readme_tree` only after the replacement passes the same
   inventory and cost gates.

## Strongest counterargument

The repository already has a sensible distributed architecture: path tells an
agent the boundary, README explains the component, Bazel exposes targets,
frontmatter drives the website, and root instructions say how to search.
Adding a catalog risks becoming another stale copy.

That objection defeats a hand-maintained omnibus manifest. It does not defeat
a generated join. The current false CODEOWNERS affordance, partial lifecycle
coverage, four workspace lists, broken indexer example, and failing docs
closure are direct evidence that humans cannot keep the joins synchronized by
prose alone. The proposed model adds only a schema and validation; mutable
facts remain in their existing natural authorities.

## Evidence limits

- Counts describe tracked files at the source revision above, except the
  untracked active goal record, which was inspected only to verify immutable
  audit bindings.
- The broad docs query failed on external OCI retrieval. That failure proves
  the map query is coupled to external repository resolution; it does not
  prove the underlying documentation build is defective.
- Textual presence of `docs_filegroup(...)` does not prove that every Markdown
  file is published or that every link is valid.
- No state-changing Bazel, infrastructure, host, delivery, or remote command
  was run.
