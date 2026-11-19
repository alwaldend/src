# Phase 2B catalog-input audit

## Binding and verdict

This read-only audit is bound to:

- source base / merged PR 39:
  `1ca06a0fc9697e8fb212b32d99d9ad3b996ea76e`;
- goal `agent-system-phase-2`, goal generation 2, lifecycle generation 4,
  criteria revision 2; and
- active attempt `goal-publication-001`.

The checked-out `HEAD` equals the bound source revision. The Phase 1 registry
digest is
`sha256:ee6024206f87005175f0213b4d81db7be9810988d0dadcbe8c5a5347c3b745e7`,
matching the resource baseline.

**Verdict: revise, then proceed in parallel with Phase 2A.** Implement the
shared deterministic envelope and the Topology, Policy, Action, Capability,
WorkspaceCheck, and AgentSystemIndex compilers without consuming continuation
state. Implement only the GoalCatalog schema and an unavailable adapter until
Phase 2A provides a pure read-only, recovery-aware validation API. The
GoalCatalog may expose a record's stable identity and `outcome`/`execution`
only after complete record validation; an invalid or interrupted record is an
eligible `unavailable` entry, never omitted or described as resumable.

## Controlling contracts

- Facts remain at natural owners; catalogs are replaceable projections with
  provenance, and conflicts are reported rather than resolved
  (`projects/agents/docs/architecture.md:30-46,53-77`).
- Deterministic catalogs carry identities, versions/digests, completeness, and
  limitations, contain no generation timestamp or absolute checkout path, and
  hash canonical schema-defined bytes excluding their own digest
  (`projects/agents/docs/architecture.md:79-107`).
- Phase 2B requires one versioned schema/compiler per concern, stable ordering,
  source-relative paths, input digests, explicit bounds, JSON plus checked
  Markdown, and a failing completeness check for omitted/invalid eligible
  skills, workspaces, goals, or components
  (`projects/agents/docs/roadmap.md:220-244`).
- Generation is offline, bounded, and non-stateful; the index contains
  references and digests rather than catalog bodies
  (`projects/agents/docs/roadmap.md:273-294`).
- Until Phase 2A passes, goal output is limited to validated identity and
  coarse status; invalid records are `unavailable`
  (`projects/agents/docs/roadmap.md:205-218`).

## Minimum shared deterministic wire contract

Each catalog JSON should have a catalog-specific schema URI and the following
common shape. The checked Markdown must render this same model and state the
JSON digest; it is not a second authority.

```text
schema                 catalog-specific .../v1alpha1 URI
kind                   stable catalog kind
id                     stable catalog identity
derivationVersion      compiler semantics version
producerRef            typed provider reference
sourceRevision         exact Git tree/commit identity used for the build
inputs[]               {path, role, digest}; sorted by (path, role)
bounds                 {eligible, emitted, unavailable, maxItems,
                        maxInputBytes, maxOutputBytes, truncated}
completeness           complete | partial | truncated | unknown
limitations[]          required whenever completeness != complete
conflicts[]             {id, code, sourcePaths[]}; sorted by id/path
items[]                 catalog-specific records in stable identity order
digest                  sha256 of canonical JSON with this field omitted
```

All paths are workspace-relative and normalized. Arrays are never `null`.
Unknown JSON fields, duplicate identities, path traversal, absolute paths,
unknown enums, unsorted set-like arrays, and a non-complete result without a
limitation fail schema validation. Do not include observation time in static
catalog bytes. A `--check` mode validates and emits the same content as report
mode, then exits nonzero for any completeness failure; invocation mode itself
must not perturb the bytes.

These rules reuse the Phase 1 reference, completeness, availability,
information, retention, and artifact-envelope vocabulary
(`tools/agents/api/v1alpha1/types.go:13-134`) and its strict/canonical JSON
behavior (`tools/agents/api/v1alpha1/canonical.go:10-61`). Phase 1 already
requires limitations for incomplete artifacts
(`tools/agents/api/v1alpha1/validation.go:69-137`).

## Catalog contracts and gates

### TopologyCatalog

Minimum records:

```text
trees[]       {id, path, readmePath, boundaryClass}
components[]  {id, path, ownerReadmePath, buildPath, title, description,
               lifecycle, docsState, docsTarget?}
workspaces[]  {id, path, modulePath, moduleName}
```

Inputs are the `repository.projects` authority, immediate `projects/*`
directories, each owner `README.md` and `BUILD.bazel`, tracked `MODULE.bazel`
roots, the six top-level boundary READMEs, and documentation membership where
declared. The owner README owns purpose/local boundaries, while MODULE/BUILD
own workspace/dependency structure (`AGENTS.md:30-40`; architecture
`projects/agents/docs/architecture.md:55-63,144-168`).

Completeness gates: every immediate registered project directory has one
README, one BUILD, a valid stable ID, title/description, and exactly one
lifecycle value; every tracked workspace root has one unique module name; each
docs claim names its BUILD source; no source path is absolute. Missing or
invalid eligible components are emitted unavailable and make `--check` fail.

Scope limitation: Phase 1 closes the maintained-component universe only at
`projects` (`tools/agents/declarations/registry.json:21-24`) and its checker
enumerates only immediate `projects/*` roots
(`tools/agents/cmd/phase1_check/main.go:434-466`). A v1 catalog may truthfully
claim complete registered-project coverage, not complete component coverage
across `infra`, `tools`, `data`, `third_party`, or `users`, until a registration
authority closes that wider universe. Tree boundary records remain safe.

### PolicyCatalog

Minimum record:

```text
policies[] {id, pathPrefix, precedence,
            agentPolicySources[], ownerBoundarySource?, reviewSource?,
            axes: {sourceDisclosure, evidenceHandling, bazelVisibility,
                   buildConsumer, artifactPublication,
                   documentationPublication, information,
                   liveEnvironmentAssociation}}
```

Every axis value carries its source path and is independently `known`,
`unknown`, or `conflict`; no axis is inferred from another. Inputs are every
tracked `AGENTS.md` in nearest-first applicability order, `CODEOWNERS`, and
owner-local boundary READMEs. The root policy explicitly separates disclosure,
target visibility, build consumers, publication, and secret/personal presence
(`AGENTS.md:42-66`); Phase 1's typed `PathPolicy` preserves disclosure,
build-consumer, publication, and information as separate fields
(`tools/agents/api/v1alpha1/types.go:70-81`).

Completeness gates: all tracked `AGENTS.md` files participate; every registered
top-level boundary has a policy record; precedence is deterministic; every
declared value cites exactly one owning source; disagreements become conflicts.
Do not claim observed target compliance from prose policy. Phase 1 registers no
policy authority, so the compiler schema must explicitly close this source
universe (or the registry must add one) before `complete` is allowed.

### ActionCatalog

Minimum records:

```text
providers[]  {id, owner, definitionPath}
actions[]    {id, providerRef, owner, sourcePath, selector, classification,
              effects[], inputs[], outputs[], information[], credentialUse,
              networkUse, environmentSelector, authorityGate, preflight,
              verification, cost, cacheability, cancellation}
aliases[]    {providerRef, selector, state, replacementRef?, reason?}
```

Inputs are the registry's four `operationFiles`
(`tools/agents/declarations/registry.json:328-332`) and their owner-local
definition paths. The operation declaration shape is already explicit
(`tools/agents/cmd/phase1_check/main.go:83-115`), and the shared atomic effect
set is closed (`tools/agents/api/v1alpha1/types.go:39-52`).

Completeness gates: strict JSON; unique action IDs, provider owners, and
provider-local selectors; known, sorted, nonempty effects; known information,
credential/network/cacheability values; all required validation fields
nonempty; definition files exist; removed aliases have an existing explicit
replacement; each supported definition surface has an exact source reconciler.
`requires_migration` remains an item and limitation, not an omission.

Important gap: the Phase 1 checker validates most operation fields only for
non-emptiness (`tools/agents/cmd/phase1_check/main.go:297-314`) and source-
reconciles Terraform selectors only
(`tools/agents/cmd/phase1_check/main.go:572-615`). Compilation from the four
manifests is safe now, but completeness must remain partial until Goal, Cordis,
and repository-delivery definition reconcilers prove the closed runnable
surface. Direct binaries are a separate registered universe and must be linked
as providers/coverage gaps rather than silently treated as action contracts.

### CapabilityCatalog

Minimum records:

```text
skills[]     {id, owner, canonicalPath, discoveryPath, layer, activation,
              exclusions[], capabilityRefs[], dependencies[], conflicts[],
              providerRequirements[], contextCost, evaluationMaturity}
providers[]  {id, owner, kind, sourcePath, classification, actionRefs[]}
```

Inputs are registry `skills`, `runtimeTools`, `directBinaries`, operation-file
providers, `.agents/skills` discovery links, and each canonical `SKILL.md` plus
its BUILD declaration. Skills own routing/procedure; providers own executable
capability (`projects/agents/docs/architecture.md:206-220`).

Completeness gates: registry/discovery sets are identical; every link is a
relative nonescaping link to one canonical skill; frontmatter identity matches;
every skill has one owner and BUILD package; dependency/conflict references
resolve, dependency graph is acyclic, and set-like fields are sorted/deduped;
every registered runtime tool/direct binary appears; action refs resolve.
External `providerRequirements` may remain explicit unresolved requirements
and must not be presented as live availability. Phase 1 currently checks the
registry/discovery name match and a small required metadata subset only
(`tools/agents/cmd/phase1_check/main.go:383-411`), so the added graph/link gates
are required before `complete`.

### WorkspaceCheckCatalog

Minimum record:

```text
workspaces[] {id, path, modulePath, moduleName,
              projections: {bazelIgnore, rootOverride, docsAggregation,
                            fullCheck},
              phases[]: {id, providerRef, commandTemplate}}
```

The authoritative universe is tracked `MODULE.bazel` roots. Reconcile it with
`.bazelignore`, root local overrides, documentation aggregation, and the full-
repository checker, as required by the roadmap
(`projects/agents/docs/roadmap.md:240-244`). The current full checker manually
lists root plus eight nested workspaces and `build`/`test`
(`projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go:19-40`),
while `.bazelignore:3-11` lists those eight nested roots and
`third_party/include.MODULE.bazel:77-127` provides their root overrides.

Completeness gates: every tracked MODULE root has a unique normalized path and
module name; all four projections have exactly the expected membership (with
explicit, schema-owned exclusions); phases use known provider/action refs;
manual counts are rejected in favor of set equality. Any drift is a conflict
and makes `--check` fail. This compiler is safe before Phase 2A.

### GoalCatalog

Minimum record before Phase 2A acceptance:

```text
goals[] {candidatePath,
         availability: available | unavailable,
         reason?,
         identity?: {name, ownerRoot, scope},
         coarseStatus?: {outcome, execution}}
```

Only a completely validated record receives `identity` and `coarseStatus`.
Never include title, criteria, relationships, generations, resource version,
attempt pointers, evidence, resume advice, or dependency state before the 2A
gate. For an invalid/interrupted record, retain only its eligible relative
directory as `candidatePath`, `availability: unavailable`, and a stable bounded
reason; do not repeat unvalidated YAML claims. The eligible universe is every
non-hidden child directory of each registered owner-local goals root.

Completeness gates: bounded root/item/file cardinality; no symlinks; strict
resource validation; directory/name match; current and immutable criteria
agreement; attempt bindings and artifact digests; exactly one active attempt
consistent with Goal status. An invalid eligible directory is emitted
unavailable and makes `--check` fail, rather than disappearing.

Existing APIs are unsafe to reuse directly:

- `List` reads only `goal.yaml` and emits generations and status without full
  record validation (`projects/goal/internal/fsstore/store.go:275-346`).
- Full `ValidateGoal` calls `loadAndValidate`
  (`projects/goal/internal/fsstore/store.go:766-827`), which begins by cleaning
  temporary residue. Cleanup removes files and directories
  (`projects/goal/internal/fsstore/integrity.go:21-79,82-122`), violating the
  offline non-stateful catalog contract.

Phase 2A must therefore expose a pure, bounded read-only inspection result that
distinguishes committed, incomplete/recoverable, and invalid states. Until
then, the GoalCatalog compiler can ship its schema, fixtures, unavailable
adapter, and negative gates, but must not ingest live goal continuation state.

### AgentSystemIndex

Minimum record:

```text
catalogs[]   {id, kind, schema, derivationVersion, digest, inputDigests[],
              completeness, queryRoutes[]}
conflicts[]  {id, catalogRefs[], sourcePaths[]}
```

No catalog payload, skill body, README text, schema body, goal state, log, or
observation is embedded. Routes are stable local labels/commands or checked
artifact paths, never checkout-specific paths. Records sort by catalog ID;
input digests sort by source identity.

Completeness gates: exactly one descriptor for every required catalog kind;
digests validate; routes resolve; byte ceiling passes; duplicate payload fields
and embedded bodies are schema-rejected. An unavailable GoalCatalog is itself
truthfully indexable: the index can be complete as an inventory while
preserving the child catalog's incomplete/unavailable state.

## Safe sequence before Phase 2A recovery completes

1. Freeze the common envelope, canonical JSON rules, bounds, report/check
   semantics, and JSON-to-Markdown parity tests.
2. Implement Topology (registered-project scope), Policy (explicitly closed
   source set), Action (initially partial until all source reconcilers exist),
   Capability, and WorkspaceCheck compilers and negative completeness fixtures.
3. Implement the GoalCatalog schema and invalid/unavailable fixtures only;
   integrate the pure goal inspector after Phase 2A.
4. Implement AgentSystemIndex over catalog descriptors, including an
   unavailable GoalCatalog descriptor.
5. Do not implement rich goal continuation, graph joins, or context resume
   advice until Phase 2A fault-injection/recovery acceptance passes.

## Decision review

The strongest case for direct reuse is that Phase 1 already reports 7
authorities, 21 skills, 10 runtime tools, 24 operations, 28 projects, 9
workspaces, and 2 goals with no missing/unclassified entries
(`projects/agents/goals/agent-system-phase-1/attempts/phase1-completion-001/evidence/phase1-validation.md:19-21`).
That path is rejected because counts are not typed catalogs, project and goal
checks are shallow, only Terraform action selectors are reconciled, no policy
authority is registered, and goal full validation mutates residue. Waiting for
all of 2A would unnecessarily block five independent static catalogs. The
least-risk alternative is the revised staged boundary above: proceed with
static owner-derived catalogs and the descriptor-only index; keep goals
unavailable until a pure inspector exists. Evidence that Phase 2A has supplied
such an inspector and interruption tests pass would change the GoalCatalog
verdict from unavailable-only to validated coarse ingestion.
