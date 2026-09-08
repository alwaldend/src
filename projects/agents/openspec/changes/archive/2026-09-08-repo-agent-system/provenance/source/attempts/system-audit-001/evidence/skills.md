# Skills, capability routing, discovery, and learning audit

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

## Verdict

The repository has a strong **skill artifact substrate** and a weak
**capability control plane**. Canonical ownership, Bazel packaging, validation,
and exact relative discovery links are unusually disciplined. What an agent
can do, which capabilities compose, which instruction wins, what authority is
required, how expensive a path is, and whether it has actually worked are
still encoded mostly in prose spread across `AGENTS.md`, 20 `SKILL.md` files,
BUILD declarations, and eval READMEs.

Consequently, the repository can reliably accrete more skill files, but it
cannot yet reliably accrete operational intelligence. Adding more prose to the
flat surface will make selection and conflict resolution harder. The highest
leverage change is a typed repository capability registry, consumed by a small
execution-harness resolver and projected into links, a human index, validation,
and eval selection. Preserve the current skills as the detailed knowledge
layer; make their relationships and observed quality first-class data.

### Decision review of the registry proposal

The strongest objection is that a capability schema could become a second
bureaucratic truth which the platform's native skill selector does not consume.
It could force fuzzy natural-language triggers into brittle enums, add update
work to every skill, and create false confidence in a resolver that has not
been behaviorally tested. With only 20 currently exact links, retaining the
flat prose system and adding evals is a credible lower-cost alternative. A
central mega-skill is another alternative, but it would concentrate ownership,
reload unrelated detail, and worsen the context problem.

**Verdict: revise, then proceed.** Use minimal typed metadata embedded with the
canonical skill plus generated views, not a separately authored central
manifest and not a universal meta-skill. Initially type only facts that can be
validated—identity, phase, relations, effects, authority gates, tools, and
evidence labels—while keeping nuanced trigger semantics in the concise
description. First generate the catalog and all-skill routing suite; make the
resolver authoritative only after that evidence shows it improves selection.
Abandon or narrow the schema if routine skill changes require synchronized
manual edits or if routing quality does not improve measurably.

## Current system and evidence

| Layer | Current implementation | What works | System gap |
| --- | --- | --- | --- |
| Repository policy | Root `AGENTS.md` | Clear authority, safety, layout, Bazel, and verification defaults | Routing and precedence are prose (`AGENTS.md:36-85`), not a queryable policy graph |
| Capability instructions | Canonical `SKILL.md` files owned by projects | Narrow owner placement and explicit trigger descriptions | Dependencies, exclusions, effects, authority, cost, and verification are not typed |
| Artifact envelope | `skill_library` and `SkillInfo` | Exact package-local files, stable logical paths, distinguished skill and optional UI metadata (`projects/rules_skill/skill/internal/skill_library.bzl:3-22,43-113`) | Provider exposes only files/name/root; it cannot answer capability questions |
| Validation | Global validation aspect and Go validator | Strong name/frontmatter/UI shape checks; root config enables it (`tools/bazelrc/root.bazelrc:5-7`) | Validation is not a transitive precondition of discovery/consumption; semantic fields and relationships are absent |
| Registration | Root `skill_discovery_links` list | One explicit source of truth for 20 registered skills (`BUILD.bazel:66-89`) | A valid but omitted new skill is invisible; completeness is not checked |
| Discovery projection | 20 tracked direct relative links in `.agents/skills` | Current link set exactly matches the root list; updater/checker enforce normalized, direct, same-repo links | Flat names expose no owner, scope, lifecycle, composition, or quality status |
| Routing | Runtime reads descriptions plus root trigger prose | Some boundaries are well written, such as package Bazel versus full-repo checks | Only one positive routing assertion exists; there is no all-skill confusion or composition suite |
| Behavioral evidence | Promptfoo configs/cases and a few manual live targets | Eval limitations are documented honestly; live suites isolate subject/judge state | Most targets validate config loading, not behavior; results are not bound into a durable capability status |
| Durable learning | Goal attempts record evidence and a process audit | Goal lifecycle already requires exact evidence and feedback review (`projects/goal/skills/goal/references/lifecycle-and-evidence.md:25-70`) | No path promotes repeated task friction into a reviewed routing case, fixture, or skill revision |

Observed inventory:

- The root registry and `.agents/skills` each contain the same 20 names: 18
  skills under `projects/agents`, the `goal` skill, and `versioning`. The links
  are Git-tracked mode-120000 direct relative symlinks. `.bazelignore:1`
  deliberately excludes the projection from Bazel package discovery.
- `bazel_agent test //:write_skill_links_test` passed during this audit. It took
  about 4.1 seconds and analyzed 221 packages / 6,185 configured targets to
  prove the 20-link filesystem projection. This is valuable exactness, but a
  relatively expensive everyday discovery-health query.
- The 20 canonical skill instruction files total 2,041 lines and about 13,948
  words. They are not all loaded for every task, but composition can multiply
  the cost. `repo-delivery` alone is 313 lines / 2,679 words and is mandated
  for every implementation handoff by `AGENTS.md:83-85`.
- Of the 18 `projects/agents` skills, all have an offline
  `promptfoo_validate_test`, but only `answer-question`, `decision-review`, and
  `spellcheck` declare live `promptfoo_test` targets. Their case files contain
  50 scenarios in total; 17 belong to `repo-delivery`, which has no live
  target. This is useful contract prose, not measured behavior.
- Four of the 18 core skills have no `agents/openai.yaml`; none of the present
  metadata files declares the validator-supported
  `policy.allow_implicit_invocation`. Optionality may be intentional, but the
  effective invocation policy is therefore not explicit in repository data.

### Structural strengths to retain

1. Canonical sources live with their narrow owner while discovery links are
   merely projections (`projects/agents/README.md:14-23`). This is the correct
   ownership direction.
2. `SkillInfo.files_by_path` avoids execution-path parsing, forbids files from
   other packages, and rejects duplicate logical paths
   (`skill_library.bzl:24-57`). Preserve this artifact contract.
3. Discovery rejects unsafe paths, cross-repository/generated skills,
   duplicate global names, symlink ancestors, wrong targets, extra entries,
   and non-symlink collisions
   (`skill_discovery_links.bzl:55-156,195-279,366-424`). Tracked links also
   avoid a first-clone bootstrap paradox.
4. Eval documentation consistently distinguishes configuration validation from
   behavioral proof (`projects/agents/README.md:25-32`). The answer-question
   suite additionally separates routing from quality and includes a no-skill
   control (`projects/agents/skills/answer-question/evals/README.md:27-50`).
5. The goal system already binds evidence to exact subjects and requires a
   process audit; it explicitly assigns capability discovery and scheduling to
   the execution harness
   (`projects/goal/skills/goal/references/graph-organization.md:63-73`). That is
   a clean seam for the missing control plane.

## Ranked findings and recommendations

### P0 — Resolve contradictory routes and authority before adding abstractions

Three current routes conflict or leave authority implicit:

- `git-rebase-remote` says to use it "before delivery"
  (`projects/agents/skills/git-rebase-remote/SKILL.md:4-14`), while the GitHub
  path in `repo-delivery` owns rebase and push mechanics and forbids bypassing
  its adapter (`projects/agents/skills/repo-delivery/SKILL.md:56-69`). The former
  is appropriate for standalone synchronization and the explicit Forgejo
  fallback (`repo-delivery/SKILL.md:254-274`), not every delivery.
- Root policy says every repository Bazel invocation must use `bazel_agent`
  (`AGENTS.md:145-151`), while the runner skill correctly needs one raw Bazel
  bootstrap when that runner is absent (`bazel-agent/SKILL.md:57-73`). The
  exception exists only below the rule it contradicts.
- Root policy selects `repo-delivery` for all completed implementations
  (`AGENTS.md:83-85`), and the skill says to commit and push without asking
  (`repo-delivery/SKILL.md:12-14`). Whether an implementation request grants
  remote-publication authority must be an explicit repository invariant, not
  an inference hidden inside a skill.

Narrow the rebase trigger, put the first-bootstrap exception at the
authoritative root layer, and separate or explicitly bind local completion,
local commit, remote publication, and forge mutation authority. Add these cases
to a cross-skill conflict suite before revising the prose.

### P0 — Introduce one typed capability registry and resolver

Extend the existing explicit root registration rather than creating a second
hand-authored index. Each repository skill should contribute a validated
capability record with at least:

- stable ID, display name, canonical owner, lifecycle, aliases/replacement;
- positive triggers, exclusions, path/data scopes, and lifecycle phase;
- unconditional and conditional `requires`, `composes_with`, and conflicts;
- observable side-effect classes, authorization gate, secret sensitivity, and
  network/live-state requirements;
- required tools, deterministic verification labels, eval tier/status, and an
  approximate context/runtime cost class.

Author semantic fields once in versioned, validated `SKILL.md` frontmatter
(the validator already permits an untyped `metadata` key); keep Bazel-only eval
and verification labels in the repository wrapper. Do not duplicate routing
facts in `agents/openai.yaml`, a central handwritten manifest, and prose.

The registry should validate unique IDs/names, dangling edges, cycles, scope
overlaps, explicit conflict resolution, and availability of required tools.
Generate from it:

1. the label list loaded by the root `skill_discovery_links` declaration;
2. the tracked `.agents/skills` projection;
3. a compact machine-readable catalog for the execution harness;
4. a human `CAPABILITIES.md` index showing trigger, scope, effects,
   composition, and evidence status; and
5. a task-specific activation plan that selects the smallest sufficient set
   before any full skill body is loaded.

This directly fills the seam the goal graph leaves to the execution harness.
Examples that should become edges rather than prose are
`repo-gazelle-plugin -> {bazel-nested-module, repo-bazel, bazel-agent}`
(`repo-gazelle-plugin/SKILL.md:9-13`) and
`repo-bazel -> bazel-agent` (`repo-bazel/SKILL.md:8-12`). Keep policy selection
distinct from authorization: selecting a write-capable skill must never itself
grant the write.

### P0 — Make “discoverable means validated” and registry completeness true

The root Bazel configuration applies `skill_validation_aspect`, but the aspect
does not traverse attributes (`skill_validation.bzl:39-50`). Discovery accepts
any target with `SkillInfo` and attaches no validation aspect
(`skill_discovery_links.bzl:453-499`). Therefore a narrow
`//:write_skill_links` run can consume an invalid skill; only separately
building that skill under the configured aspect validates it.

Likewise, the exact-state test proves only
`links == BUILD.skills`. Its expected set is constructed solely from the labels
the caller supplied (`skill_discovery_links.bzl:98-156`), so it cannot detect a
new canonical `skill_library` omitted from the root list. The Gazelle extension
creates local `skill_library` declarations but not global registration
(`projects/rules_skill_gazelle/README.md:10-12,61-65`).

Keep generic `rules_skill` concerned with portable artifact mechanics. Add a
repository policy wrapper/registration rule that:

- makes semantic validation a dependency of every discovery/eval consumer;
- declares explicit eligibility such as `discoverable = true` so examples and
  tests are excluded deliberately;
- checks the generated registry against every eligible skill target;
- requires an eval declaration or a typed, reviewed coverage exception; and
- validates referenced files and typed capability fields.

This also removes accretion boilerplate: all 18 core BUILD files repeat the
same `skill_library` declaration, and 14 repeat the same simple library plus
offline-eval shape. Runtime resources should become explicit. The current
`full-repo-check` broad glob packages its Go implementation and unit-test source
inside the model-facing skill bundle (`full-repo-check/BUILD.bazel:33-44`), even
though the instructions invoke the separately built executable. A repository
wrapper should keep implementation and test source out of runtime context
unless instructions actually reference it.

The Go validator currently accepts `allowed-tools`, `license`, and `metadata`
as keys without validating their contents (`validate.go:19-25,177-206`) and
emits only `Skill is valid: <name>` (`main.go:75-93`). Emit a versioned JSON
validation record containing schema version, identity/digest, parsed metadata,
and warnings so downstream catalog and evidence tools can consume facts rather
than reparse prose.

### P1 — Evaluate selection, execution, and outcome as separate layers

Only `answer-question` has a routing config, and it contains one positive
`skill-used` assertion against a workspace containing that skill alone
(`answer-question/evals/promptfooconfig.routing.yaml:27-36`). Its README notes
that “skill use” is inferred from a successful `SKILL.md` read rather than a
first-class activation event (`answer-question/evals/README.md:27-29`). This
cannot detect false positives, near-neighbor confusion, missing compositions,
wrong ordering, or a skill that is read but not followed.

The broad question trigger also needs an inert-payload boundary. A question in
text supplied for `spellcheck` is data, not a second user outcome; a material
decision question, by contrast, intentionally composes `answer-question` and
`decision-review`. Encode and test both. The shared suite should also cover
`repo-bazel` versus `full-repo-check`, ordinary versus nested modules, ordinary
Gazelle updates versus plugin work, Terraform/Ansible plus conditional secrets,
and delivery versus remote-rebase fallback.

Build a tiered, diff-aware evidence system:

1. **Static/structural:** schema, graph, references, discovery, and declared
   evaluation coverage; cheap and mandatory.
2. **Routing:** all registered skills present; table-driven positive,
   adjacent-negative, exclusion, and multi-skill cases with exact
   required/forbidden sets. Record first-class activation decisions and reasons.
3. **Deterministic fixture execution:** isolated Git remotes, small Bazel
   workspaces, fake forge/tool endpoints, synthetic Terraform/Ansible state,
   and action-trace/postcondition assertions. This is the missing middle for
   the 15 core skills whose READMEs currently explain why a read-only live model
   call is insufficient.
4. **Periodic live semantic eval:** only where a model judgment adds signal.
   Compare skill/no-skill or previous revision, sample repetitions, and bind
   results to skill digest, capability-registry digest, model/runtime version,
   fixture version, and judge. Do not put billable stochastic calls on ordinary
   `//...`.

Require stable IDs for every case and requirement assertion. Currently only 10
of 50 core shared cases have IDs, and only 16 of 55 shared assertions have an
explicit metric. Reject normalized duplicates statically: the repository
delivery corpus currently repeats both the started-review and unobservable-
review scenarios (`repo-delivery/evals/cases.yaml:179-215,261-295`). Of the 50
core cases, only 16 are exercised by any live target; report that distinction
instead of presenting configuration cases as one undifferentiated count.

Publish a compact per-capability evidence status (`structural`, `routing`,
`fixture`, `live`, `stale/unverified`) in the generated catalog. This lets an
agent calibrate confidence before acting instead of treating all discovered
skills as equally proven.

### P1 — Add a reviewed friction-to-regression learning loop

The goal lifecycle already asks each attempt to name the dominant remaining
failure and audit feedback bottlenecks
(`lifecycle-and-evidence.md:48-65`). No audited file connects those observations
to skill maintenance or eval cases.

At task close, let the harness write a sanitized structured friction record:
capabilities considered/activated, routing confidence, conflicting rules,
missing tool/knowledge, avoidable reads or commands, failed assumptions,
verification latency, and the exact public evidence identity. Aggregate by
stable signature, not free-form similarity alone. A repeated issue should
create a **proposal**, never silently rewrite instructions. Promotion should
require review, a minimized regression case or deterministic fixture, the skill
change, and fresh tiered evidence. Bind the accepted lesson to the originating
attempts and resulting skill digest.

Keep these records local by default; only sanitized public-fixture cases should
be checked in. Root visibility rules at `AGENTS.md:10-20` remain the disclosure
boundary. This creates durable compounding without turning private task traces,
secrets, or one-off model guesses into repository doctrine.

### P1 — Reduce context by making skills routers over conditional references

Because the runtime requires a selected `SKILL.md` to be read completely,
monolithic instructions are a direct recurring cost. `repo-delivery` loads both
the roughly 200-line GitHub protocol and the 47-line Forgejo fallback after
every implementation, although provider detection chooses only one. By
contrast, the `goal` and `blender-reference-fidelity` skills already use
conditional `references/` resources.

Keep each `SKILL.md` to identity, triggers/exclusions, invariant safety gates,
the high-level state machine, and explicit reference-routing rules. Move
provider-, platform-, and rare-failure detail into validated references loaded
only after the relevant observation. Start with `repo-delivery`, then examine
the broadly triggered 1,075-word `answer-question` skill. Deduplicate shared
contracts by declaring composition in the capability graph; do not copy the
same procedure into every dependent skill. Track approximate instruction bytes
selected per task and regress that metric when behavior is unchanged.

### P2 — Clarify the secret-scratch boundary

Root policy requires task-owned temporary files under ignored `out/<task>/`
(`AGENTS.md:22-34`), while `repo-secrets` says temporary secret files must be
"outside the repo" (`repo-secrets/SKILL.md:34-40`). This is wording ambiguity,
not a substantive conflict: "repository" here means tracked or committed
source, so a task-private ignored `out/<task>/` location is acceptable when it
has restrictive permissions, bounded lifetime, and reliable cleanup. State
that meaning directly so an agent does not interpret "repo" as the whole
checkout and unnecessarily move sensitive scratch elsewhere.

### P2 — Give discovery a cheap, explainable operational surface

Retain checked-in relative links for zero-step bootstrap, but add one
repository-owned `tools/agents` `check`/`sync`/`doctor` surface around the
generated registry and projections. It should support structured `--check`,
human `--explain`, explicit `--fix`, and separately authorized `--prune`, report
canonical owners and unresolved graph/tool requirements, and diagnose stale
locks. Root `AGENTS.md` and `projects/agents/README.md` currently contain no
start-here catalog/repair command.

Keep the Bazel exact-state test as the authoritative projection integration
test, but separate hermetic registry/graph checks from the necessarily local
filesystem check. The current test is tagged `local`, `no-cache`, `no-remote`,
and `no-sandbox` (`skill_discovery_links.bzl:534-543`), relies on one skill
runfile to rediscover the checkout (`skill_discovery_links.bzl:287-349,473`),
and is absent from the pre-commit hook (`tools/git_hooks/main/go/precommit.sh:11-13`).
Add the cheap completeness/catalog check to the normal hygiene path; reserve
link repair for an explicit mutation.

Longer-term, replace the duplicated generated Bash updater/checker with a small
tested portable implementation plus a generated manifest. Preserve the current
strict symlink invariants and fail-closed collision behavior. Improve
interruption diagnostics and make pruning visible rather than silently
unlinking unknown symlinks.

Close the present test gaps while doing so. The nested `rules_skill` suite
executes the updater but not the generated checker, despite the checker carrying
many exactness guarantees (`projects/rules_skill/tests/analysis/BUILD.bazel:149-197`).
Add isolated checker failure cases, schema tables/goldens, an integration test
that proves every consumer forces validation, crash/idempotence tests, and
fuzzing for Markdown/YAML/path handling.

## Recommended implementation sequence

1. Correct the three contradictory/authority routes and add regression cases
   for each.
2. Define the capability schema and stable IDs; encode the 20 current skills
   without changing their behavior. Generate the human and machine indexes.
3. Add completeness and consumer-enforced validation; make the root label list
   a generated projection while preserving root-package visibility and tracked
   symlinks.
4. Add the resolver/activation trace and the cross-skill routing confusion
   suite. Use it to settle overlaps before adding new skills.
5. Refactor `repo-delivery` into a compact router plus conditional provider
   references and measure context reduction.
6. Add deterministic fixture harnesses for the highest-risk action skills, then
   publish evidence status in the catalog.
7. Connect sanitized goal friction records to a human-reviewed
   regression-first promotion queue.
8. Clarify the nonblocking secret-scratch terminology when the owning policy
   documents are next touched.

The controlling design principle is: **one typed fact, many generated views**.
Canonical skill prose should explain how to perform a capability; the registry
should explain when and whether to activate it; the execution harness should
resolve the smallest authorized composition; evidence should say how well that
composition is known to work; and the goal loop should propose reviewed,
regression-backed improvements when reality contradicts the model.
