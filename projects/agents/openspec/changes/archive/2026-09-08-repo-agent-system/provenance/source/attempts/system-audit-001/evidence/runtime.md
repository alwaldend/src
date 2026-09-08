# Runtime, capability discovery, and workspace-control audit

## Audit binding and scope

- Goal: `repo-agent-system`
- Attempt: `system-audit-001`
- Resource version: `5`
- Goal generation: `1`
- Lifecycle generation: `4`
- Criteria revision: `2`
- Criteria digest:
  `sha256:2ac2db1242f5d3358e433b3499da5a622d06bdec49bfa690dd34cf3205e28f34`
- Goal-state digest:
  `sha256:193be5b38881faebc349f9ae1d273e24fac5d5925a9a4402b24706394ffaeb3a`

This was a read-only source audit of `.codex`, `projects/mcp_cordis`, the
repository-owned goal and Bazel-agent runtime boundaries, skill discovery, and
the tools through which an agent observes a workspace. No network, Bazel,
runtime, goal, or maintained-source mutation was performed. The only write is
this ignored audit artifact.

Codex client behavior was not externally verified because this worker was
explicitly denied network access. Statements below about `.codex/config.toml`
describe the checked-in declaration; where client consequences matter, they
are identified as inferences rather than verified product behavior.

## Verdict

The repository already has unusually good local mechanisms: exact skill-link
validation, a fixed MCP gateway, bounded structured repository reads, guarded
Git inspection, process-group cleanup, partial-scope startup recovery, and a
digest-bound durable goal store. The highest-leverage defect is not the quality
of those pieces. It is the absence of a repository-owned control-plane join
between them.

A zero-context agent must currently reconstruct five independent planes:

1. policy and routing from `AGENTS.md` and discovered skills;
2. durable intent from goal files and session bindings;
3. workspace state from Git, filesystem, and Bazel commands;
4. runtime capabilities from fixed MCP tools and a second dynamic Cordis
   catalog; and
5. optional host plugins, apps, and connectors from state outside the
   repository.

No stable interface binds those observations to one workspace/session/task
identity, reports their health and provenance, or carries authority and
resource budgets into execution and evidence. The goal architecture explicitly
leaves live topology, capability discovery, delegation, and scheduling to the
harness (`projects/goal/docs/architecture.md:166-168`), while the repository
does not define the adapter contract that would reconnect that harness state to
its durable records. This is the break in the abstraction tower.

## Current control surfaces

| Plane | Current authority | Useful guarantee | Missing join |
| --- | --- | --- | --- |
| Skill discovery | Root `BUILD.bazel:66-89` and `.agents/skills` symlinks | One canonical `skill_library` per skill; exact-state link test is documented in `projects/rules_skill/README.md:86-129` | No runtime health, version/provenance, applicability, or capability/effect index |
| Codex/MCP discovery | `.codex/config.toml:1-10` | Active worktree is resolved before the checked-in MCP launcher runs | Only Cordis is declared; no repository capability requirements or optional/fallback status |
| Dynamic runtime | `projects/mcp_cordis/cordis.yaml:1-6`, fixed gateway in `internal/mcp.mjs:54-244` | Dynamic packages remain invocable despite clients caching the initial MCP tool list | Separate package/tool catalog, process-local version, no boot epoch or desired/observed revision |
| Workspace observation | `repo_context_get/search/read` and `git_snapshot/compare` | Bounded output, explicit truncation, guarded paths, and Git commands configured to avoid prompts/lazy fetch | No one-shot coherent workspace/goal/runtime snapshot; instruction and schema content is often reloaded in full |
| Execution | Native agent tools, `bazel_agent`, and Cordis `ctx.exec()` | Bazel agent mode is consistent; Cordis tracks and kills original Linux process groups | No shared task identity, budget, capability policy, or evidence receipt across executors |
| Durable intent/evidence | `projects/goal` resources and goal skill | Optimistic concurrency, immutable criteria history, artifact digests, and one canonical writer | Runtime agents/topology and capability observations are intentionally outside the record protocol |
| Accretion | Cordis scratch packages and `cordis_promote` | Scratch code can become reusable project code without a server restart | Promotion does not require ownership, tests/evals, evidence, usage signal, or delivery review |

A tracked-file search found no repository-level Codex plugin, app, or connector
manifest. The only tracked Codex runtime registration is
`.codex/config.toml`; skill discovery is a different mechanism rooted in the
top-level `BUILD.bazel`. This is not a reason to install user-specific
connectors automatically. It is a reason to define semantic capability
requirements and let an adapter report which repository, native, or connected
provider satisfies them.

## Ranked improvements

### 1. P0: define one read-only agent-system status contract

**Evidence.** Skill identities are enumerated in `BUILD.bazel:66-89`; Cordis
starter packages are separately enumerated in
`projects/mcp_cordis/cordis.yaml:1-6`; the MCP config is separate again at
`.codex/config.toml:1-10`. `cordis_list_tools` exposes only live package
handlers (`projects/mcp_cordis/internal/mcp.mjs:137-153`). The goal design says
capability discovery and live scheduling belong to the harness
(`projects/goal/skills/goal/references/graph-organization.md:63-65`) but defines
no repository-facing harness adapter.

**Impact.** Every session pays repeated discovery calls and reasoning cost.
More importantly, an agent cannot distinguish “not installed,” “not permitted,”
“failed to initialize,” “not applicable,” and “not yet discovered” without
probing unrelated layers. Decisions are therefore made against partial,
time-skewed state.

**Change.** Add a small, versioned, read-only `agent_system status --json`
contract (also exposable as an MCP resource/tool) that *projects* existing
authorities rather than replacing them. Its context envelope should include:

- canonical workspace/worktree identity and revision;
- task/session/coordinator identity and current goal binding, if any;
- applicable instruction paths plus content digests;
- discovered skills with source labels and versions/digests;
- semantic capability requirements and current providers (native, Bazel,
  Cordis, plugin, app, or connector), including unavailable/optional reasons;
- runtime instance, package health, desired/observed revisions, and supported
  engines;
- Bazel runner and pinned Bazel provenance;
- authority, side-effect, secret, network, time, byte, and concurrency budgets;
  and
- pointers to current evidence/receipts, never embedded secret-bearing output.

Keep mutable facts in their present authorities. The status document should
carry source references and digests so duplication is detectable. It must be
safe and useful when Cordis, Bazel, Git, a connector, or a goal binding is
absent.

**Acceptance signal.** From a clean root session, one bounded call identifies
the workspace, applicable policy, current intent, capability providers,
degraded components, and next safe discovery actions. Each field names its
authority and observation time; removing any optional provider yields a
structured unavailable state rather than loss of the whole snapshot.

### 2. P0: separate the always-responsive control kernel from extension loading

**Evidence.** `mcp_cordis` awaits complete runtime initialization before it
creates the stdio server (`projects/mcp_cordis/cmd/mcp_cordis/main.mjs:83-97`).
Initialization mounts Loader, Timer, and HMR, then loads project and scratch
Includes sequentially before returning
(`projects/mcp_cordis/internal/runtime.mjs:422-487`). Package evaluation and
`apply()` may never settle, and the README explicitly says such a package can
stall lifecycle work (`projects/mcp_cordis/README.md:106-110`). The checked-in
config declares this MCP `required = true` with a 120-second startup timeout
(`.codex/config.toml:8-10`). The exact client failure UX was not verified, but
the repository clearly intends a fallible extension host to be mandatory. A
scope is also validated as one unit before its Include is mounted, so one
invalid or missing module rejects every package in that scope
(`internal/runtime.mjs:1218-1235`).

**Impact.** One bad reusable package can prevent the health/list/stop controls
needed to diagnose it from becoming reachable. Because project is attempted
before scratch, a hanging project activation also prevents the otherwise
healthy scratch scope from loading. This turns optional convenience packages
into a control-plane single point of failure and spends up to the outer startup
deadline before degradation is visible.

**Change.** Establish the fixed MCP transport and a minimal health/status kernel
first. Load each package asynchronously behind per-package admission and
startup deadlines. Health must represent `loading`, `ready`, `degraded`,
`failed`, `timed_out`, `draining`, and `disabled` without waiting for package
code. Healthy scopes and packages must remain usable. Either make the extension
host optional to the Codex session or keep only this non-extensible kernel
mandatory; do not make arbitrary plugin evaluation part of mandatory startup.

**Acceptance signal.** Inject a never-settling project package. Within a small
fixed deadline, status and management calls respond, the package is identified
as timed out, healthy packages remain callable, and the enclosing agent session
retains native/fallback capabilities.

### 3. P0: make task/session isolation real, not merely worktree-scoped

**Evidence.** Every runtime in one worktree uses the same
`out/mcp_cordis/cordis.yaml` and `out/mcp_cordis/plugins`
(`projects/mcp_cordis/internal/runtime.mjs:380-410`). Package identity has only
`project` or `scratch` scope plus name (`internal/mcp.mjs:4-6` and
`projects/mcp_cordis/README.md:10-14`). Mutation serialization is an in-memory
promise chain inside one `CordisRuntime` instance
(`internal/runtime.mjs:1311-1324`); atomic file replacement is per file and has
no cross-process compare-and-swap (`internal/runtime.mjs:290-306`). A tracked
test search found restart and in-process HMR concurrency coverage, but no two
live runtime processes mutating the same scratch catalog. The shared location
also conflicts with the repository-wide default that task-owned scratch live
under `out/<task>/` (`AGENTS.md:19-30`).

**Impact.** Concurrent sessions or delegated workers in one worktree can see,
reload, overwrite, stop, remove, or promote each other's scratch packages.
Disposable state survives without an owner or lease. The design isolates linked
worktrees, but not tasks within a worktree—the unit that the goal worker
protocol expects to have disjoint output locations
(`projects/goal/skills/goal/references/sessions-and-concurrency.md:86-100`).

**Change.** Introduce an explicit context identity hierarchy:
`workspace -> session/task -> coordinator/worker -> package invocation`.
Scratch catalogs should be namespaced by task/session and carry an owner,
creation time, lease/retention policy, and immutable input binding. Shared
mutation needs a cross-process lock plus expected revision; define, remove, and
promote should take expected source/config digests. Ordinary workers should
receive isolated scratch and never project-write authority. Project packages
should normally be read-only at runtime, with promotion expressed as a
reviewable candidate rather than an immediate shared-registry mutation.

**Acceptance signal.** Two simultaneous sessions can define the same scratch
name without observation or mutation across namespaces. A late worker cannot
publish after its goal/task binding changes. Intentional sharing uses an
explicit lease and stale writes fail deterministically.

### 4. P0: isolate scratch execution for reliability and context safety

**Evidence.** All packages share one Node process and Cordis root context
(`projects/mcp_cordis/internal/runtime.mjs:437-455`). Packages have normal Node
built-ins, and the README explicitly says the host is not a security sandbox
(`projects/mcp_cordis/README.md:106-128`). `ctx.exec()` starts with the MCP
server's entire environment and applies caller overrides
(`internal/runtime.mjs:669-701`). The base `ctx.readText()` helper only checks a
lexically resolved path and then opens it normally
(`internal/runtime.mjs:128-151,164-186`), so an internal symlink can lead outside
the workspace. The checked-in `repo_context` package adds stronger descriptor-
bound checks for its own tools (`plugins/repo_context.mjs:136-175`), but an
arbitrary scratch package need not use them.

**Impact.** A buggy or generated scratch package can block the event loop,
terminate or corrupt the shared process, retain operations indefinitely, read
ambient environment/context, and bypass the workspace helper's apparent path
boundary. Even on a dedicated machine, this is a fault-containment and result-
integrity problem, not only a hostile-code security concern.

**Change.** Preserve in-process Cordis for reviewed, trusted project packages if
its low overhead is valuable. Default generated scratch packages to a worker or
subprocess boundary with sanitized environment, explicit filesystem/network/
process capabilities, memory/CPU/time/output quotas, and revocable leases. Make
host helpers descriptor-bound and capability-enforced. Name the trust tier in
discovery and results so “scratch” is never mistaken for an isolation claim.

**Acceptance signal.** Fixtures that loop forever, call `process.exit`, mutate
globals, follow an escaping symlink, or read a withheld environment variable
cannot impair the control kernel or another package. Their termination and
denial states are observable and bounded.

### 5. P1: expose stable protocol identity and desired-versus-observed state

**Evidence.** The server advertises only a hard-coded `0.1.0` implementation
version (`projects/mcp_cordis/internal/mcp.mjs:38-41`). `catalogVersion` starts
at `1` on every process and increments when tools enter or leave the local map
(`internal/runtime.mjs:366-379,566-597`). The stale check compares only that
integer (`internal/runtime.mjs:1085-1094`), so a value can repeat after a
restart. Package inspection reports the persisted source hash and a running
boolean, not the live source revision or activation error
(`internal/runtime.mjs:1022-1040`). Updates to running entries explicitly
return `activation: "pending"`, and clients are told to poll invocation or the
tool list (`projects/mcp_cordis/README.md:89-104`); the stdio test implements
that polling loop (`projects/mcp_cordis/test/stdio_test.mjs:46-60,144-160`).

**Impact.** A caller cannot correlate a persisted update with the generation
actually serving requests. Catalog stale protection has an ABA ambiguity
across reconnects. Polling behavior spends calls/tokens and can confuse “old
generation still healthy,” “new generation loading,” and “new generation
failed.”

**Change.** Version the public control contract separately from the binary.
Return a runtime boot ID/catalog epoch and a content-derived catalog ETag with
every catalog token. Give each tool its own contract/revision hash so unrelated
package churn does not invalidate an otherwise current invocation. Track
monotonic desired and observed package revisions plus transition state,
timestamps, last bounded error, and health probe result. Expose a bounded
wait/status operation driven by Cordis lifecycle events. This does not require
a second source store or custom rollback transaction: Git and the existing
files can remain source authorities while the runtime reports what it has
observed.

**Acceptance signal.** A pre-restart catalog token is always rejected after
restart; each update reaches an exact `observedRevision == desiredRevision` or a
terminal bounded failure; no behavioral polling through the package tool is
needed to determine convergence.

### 6. P1: make capabilities machine-routable and policy-bearing

**Evidence.** Dynamic tool normalization retains only `name`, `description`,
and `inputSchema` (`projects/mcp_cordis/internal/runtime.mjs:235-287`). The
gateway therefore loses annotations, output contracts, side effects,
idempotence, retryability, cost, authority, secret/network needs, and support
status. `cordis_invoke` itself has no side-effect annotation because the
underlying operation may be anything (`internal/mcp.mjs:155-190`). Fixed
management tools have some read-only/destructive hints, but those do not
describe invoked package behavior (`internal/mcp.mjs:54-244`).

**Impact.** An agent must infer safety and expense from prose or source. The
gateway cannot admit operations against the user's authority, select the
cheapest equivalent provider, explain why a capability is unavailable, or
reuse compatible evidence reliably.

**Change.** Define a versioned semantic capability contract shared by skills,
Cordis tools, Bazel workflows, and optional connectors. At minimum include:
stable capability ID and contract version; input and output schemas; effect
class; external-state/network/secret requirements; idempotence; expected and
maximum resource budgets; timeout and cancellation semantics; error taxonomy
with retryability; provider provenance; and evidence/receipt type. Discovery
should first return compact intent tags and summaries, then fetch a selected
schema on demand.

**Acceptance signal.** Given a task and authority envelope, routing can select
or reject a provider without reading implementation source. A networked or
state-changing handler cannot run through a generic gateway unless its declared
effects fit the admission policy.

### 7. P1: remove avoidable boot cost and report launcher provenance

**Evidence.** Every MCP start invokes the checked-in shell launcher, which runs
`bazel_agent run --script_path=...` before executing the generated binary
(`projects/mcp_cordis/cmd/mcp_cordis/launch.sh:5-21`). `bazel_agent` always uses
Bazel batch mode (`projects/bazel_agent/cmd/bazel_agent/main.go:17-27`), while
its installed host binary is explicitly not updated automatically after source
changes (`projects/bazel_agent/README.md:40-56`). Bazel itself is pinned by
`.bazeliskrc:1-2`, but the wrapper has no checked startup provenance handshake.
The launcher leaves its unique generated script under `out/mcp_cordis`; it has
no cleanup path after `exec` (`cmd/mcp_cordis/launch.sh:12-21`). All three
starter packages and HMR watchers are eagerly initialized before MCP transport
(`internal/runtime.mjs:422-487`).

**Impact.** A warm session still pays Bazel process startup/analysis and eager
extension initialization; a cold session may also need repository evaluation
or dependency acquisition. A stale installed runner can silently mediate the
mandatory launcher. The fixed 120-second outer startup setting is a coarse
mask, not a latency budget or diagnosis.

**Change.** Give `bazel_agent` a versioned `doctor/provenance --json` contract
and reject a source/installed mismatch with an exact repair command. Cache or
install an MCP launcher artifact keyed by source, lockfile, toolchain, and Bazel
version; keep a source launcher fallback that reports which phase is slow.
Start the minimal status kernel directly and lazy-load packages. Record startup
phase timings locally in bounded, privacy-safe telemetry.

**Acceptance signal.** Warm offline startup meets an explicit small latency
budget without Bazel analysis; cold startup reports phase/provenance and an
actionable degraded state; a stale runner is detected before it mediates the
runtime. Measure p50/p95 rather than asserting unmeasured speedups.

### 8. P1: turn promotion into a closed learning loop

**Evidence.** `cordis_promote` copies scratch source directly into the project
scope (`projects/mcp_cordis/internal/mcp.mjs:223-243` and
`internal/runtime.mjs:931-951`). The project BUILD glob automatically includes
plugin files (`projects/mcp_cordis/BUILD.bazel:24-27`), but promotion does not
create a test, documentation, owner, compatibility declaration, evidence link,
or delivery packet. The original starter set was selected by a one-time
aggregate review of 40 sessions (`projects/mcp_cordis/goals/runtime_extensions/evidence.md:102-104`).

**Impact.** The fast path to accretion bypasses the mechanisms that make an
improvement trustworthy and maintainable. Conversely, failures and repeated
manual work do not automatically become bounded capability proposals. The goal
store preserves evidence and the delivery workflow preserves review, but no
contract connects them to runtime capability promotion.

**Change.** Make scratch-to-project promotion a candidate pipeline:

1. capture a privacy-safe recurring need or stable failure signature;
2. bind it to a goal/attempt and capability contract;
3. materialize candidate source in an isolated task area;
4. generate contract, tests/evals, ownership, fallback, and resource budgets;
5. run focused validation and compare against the prior provider;
6. publish only through repository delivery/review; and
7. update the generated capability catalog, retaining provenance and a
   deprecation/rollback path.

Do not collect transcripts or secrets. Prefer local counters, stable error
codes, repeated command shapes, and user-approved evidence.

**Acceptance signal.** Every reusable capability points to an owner, contract,
tests/evals, promotion evidence, exact delivered revision, and retirement
rule. A scratch package cannot mutate project source merely by invoking a live
runtime management tool.

### 9. P2: make context loading incremental and content-addressed

**Evidence.** `repo_context_get` defaults to embedding applicable `AGENTS.md`
contents as well as Git state (`projects/mcp_cordis/plugins/repo_context.mjs:614-776`),
with a default 128 KiB budget. A Codex session has normally already received
root instructions, so this can duplicate its largest policy input. Its quick
Git status deliberately excludes untracked files
(`plugins/repo_context.mjs:532-610`), while `git_snapshot` includes all
untracked files (`plugins/git_worktree.mjs:563-573`); neither interface marks
that semantic difference in a shared observation contract.
`cordis_list_tools` returns complete schemas for every selected live handler
(`projects/mcp_cordis/internal/runtime.mjs:1043-1070`), and server instructions
tell the client to use that catalog first (`internal/mcp.mjs:38-47`). Tool
results serialize the same value as pretty JSON text and structured content
(`internal/mcp.mjs:8-12`), which may be redundant depending on client handling.
`cordis_inspect` can include an entire package source of up to roughly 2 MiB in
that duplicated envelope (`internal/mcp.mjs:67-80,83-102` and
`internal/runtime.mjs:24,1004-1040`).

**Impact.** Correct discovery can still spend unnecessary bytes, model context,
and calls. There is no way to say “I already hold instruction digest X” or ask
for only policy deltas and relevant capability summaries.

**Change.** Return instruction paths, precedence, size, and digest first; fetch
content only for unknown digests. Let clients provide known digests. Split the
capability catalog into compact intent/effect summaries and on-demand schemas.
Offer one task-oriented bootstrap response that selects likely capabilities
without embedding every schema. Stream file line ranges instead of reading an
entire up-to-2-MiB file before slicing. Serve large source/log artifacts as
content-addressed resources or ranged reads rather than gateway JSON. Measure
bytes and calls per workflow. Extract one repository-observation contract for
path binding, Git semantics, truncation, and topology pointers; include the
applicable `AGENTS.md` chain and nearest owning `README.md`, `BUILD.bazel`, and
`include.MODULE.bazel` as paths/digests, not eagerly repeated bodies.

**Acceptance signal.** A repeated task with unchanged instructions transfers no
instruction body again. Root-to-package discovery remains under an explicit
byte/token budget, while every omitted/truncated field is detectable and
fetchable by stable digest.

### 10. P2: align deadline, degraded-engine, and error contracts

**Evidence.** The outer MCP tool timeout is 300 seconds
(`.codex/config.toml:10`), while `cordis_invoke.timeout_ms` also permits exactly
300,000 ms (`projects/mcp_cordis/internal/mcp.mjs:155-189`), leaving no cleanup
or transport margin. An invoke timeout cancels supervised child processes but
does not cancel the JavaScript handler, which may retain its Fiber lease and
delay reload/stop/shutdown (`projects/mcp_cordis/README.md:130-136`).
`repo_context_search` falls back when `rg` is absent, but directory search then
fails while explicitly selected files remain supported
(`projects/mcp_cordis/plugins/repo_context.mjs:430-473,919-940`). This engine
state is discovered only while attempting a call. Runtime errors mix lower-case
domain codes, upper-case `EXEC_*` codes, and raw host errors
(`internal/runtime.mjs:31-52`; `internal/process_supervisor.mjs:164-334`).

**Impact.** The client can abandon a call before the inner layer returns its
structured timeout/cleanup result. Agents cannot plan a fallback until after a
failed search and cannot consistently decide which failures are retryable.

**Change.** Propagate one absolute deadline through client, gateway, handler,
and child operations, reserving bounded cleanup/serialization margin. Expose
engine availability in status and allow explicit provider selection. Publish a
namespaced error schema with phase, retryability, partial-result validity,
observed state, and corrective action. Never label a still-running handler as
cancelled; report the lease/drain state exactly.

**Acceptance signal.** The outer deadline always exceeds every admitted inner
deadline plus measured cleanup margin. Forced timeout tests receive the
structured terminal result. Missing `rg`, Git, Bazel, network, and optional
connectors are visible before invocation with a valid fallback or explicit
unsupported reason.

## Runtime slice of the target abstraction tower

The following boundaries would make the existing mechanisms cohere without a
large replacement framework:

1. **Intent** — goal resources own desired outcome, criteria, and durable
   evidence identity.
2. **Context** — a short-lived envelope binds workspace, task/session,
   coordinator/worker, revisions, authority, and budgets.
3. **Capability** — a generated semantic catalog maps intent to versioned
   provider contracts across skills, native tools, Bazel, Cordis, and optional
   connectors.
4. **Policy/admission** — one decision consumes the context envelope and
   capability effects before execution; it never expands user authority.
5. **Execution** — Bazel and runtime providers execute behind uniform deadline,
   cancellation, isolation, and provenance contracts.
6. **Observation** — a bounded status/snapshot API reports desired and observed
   state, partiality, health, and exact provider identity.
7. **Evidence** — structured receipts bind operation identity and artifacts to
   goal criteria without embedding secrets. Runtime providers may write only
   bounded receipts to isolated attempt scratch; the goal coordinator alone
   selects and imports canonical evidence.
8. **Accretion** — repeated evidence proposes a candidate capability; tests,
   review, delivery, catalog update, and eventual retirement close the loop.

The essential invariant is that each layer has one authority and exposes a
versioned projection to the next. A new “agent manifest” that manually repeats
the skill list, Cordis YAML, goal state, and tool versions would worsen the
system. Generate or query the projection from those authorities and bind it by
digest.

## Strengths to preserve

- Worktree resolution and scratch roots are explicit, and the launcher releases
  Bazel's output-base lock before the long-lived server starts
  (`projects/mcp_cordis/README.md:18-31`).
- The fixed list/invoke gateway is a pragmatic compatibility boundary for MCP
  clients that cache initial schemas (`projects/mcp_cordis/README.md:44-47`).
- Package and invocation identities include explicit project/scratch scope, and
  stale in-process catalogs can be rejected (`internal/runtime.mjs:64-80,1085-1094`).
- Repository reads/searches and Git snapshots expose bounds and truncation
  rather than silently treating partial output as complete
  (`plugins/repo_context.mjs:674-1174` and
  `plugins/git_worktree.mjs:516-879`).
- Git inspection strips repository-redirecting environment and disables lazy
  fetch/prompt behavior (`plugins/git_worktree.mjs:118-151`; equivalent Git
  controls exist at `plugins/repo_context.mjs:178-210`).
- Linux child-process cleanup is explicit, bounded on repeated `/proc`
  inspection failure, and joined before results settle
  (`internal/process_supervisor.mjs:94-157,164-334`).
- Malformed configuration can degrade one scope while unfiltered listing
  preserves a healthy scope (`internal/runtime.mjs:457-487,955-1001` and
  `test/runtime_test.mjs:602-647`).
- Skill source ownership and discovery links already have a single-authority
  model with exact reconciliation tests (`projects/agents/README.md:14-31` and
  `projects/rules_skill/README.md:86-135`).
- Goal resources already supply the durable identities, digests, optimistic
  concurrency, and worker publication rules needed by the proposed context and
  evidence layers (`projects/goal/README.md:30-40` and
  `projects/goal/skills/goal/references/sessions-and-concurrency.md:19-100`).

## Audit limits

- This report infers behavior from the current working-tree source and checked-
  in tests; it does not claim deployed-host equivalence.
- No startup latency, token count, throughput, or concurrent-session race was
  measured. Recommendations deliberately ask for those measurements instead of
  asserting numerical gains.
- The absence of a repository plugin/connector manifest is bounded to tracked
  files searched in this checkout; user-level or harness-managed integrations
  may exist outside it.
- One worker performed a read-only live Cordis catalog inspection; no reported
  conclusion depends on it. No live MCP, Codex, connector, Goal, Git, or Bazel
  mutation was performed, and no Bazel command was run.
