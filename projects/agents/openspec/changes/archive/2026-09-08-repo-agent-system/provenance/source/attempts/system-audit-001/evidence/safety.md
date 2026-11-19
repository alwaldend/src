# Repository agent-system safety audit

## Immutable bindings

- Goal resource version: `5`
- Goal generation: `1`
- Lifecycle generation: `4`
- Attempt: `system-audit-001`
- Criteria revision: `2`
- Criteria digest:
  `sha256:2ac2db1242f5d3358e433b3499da5a622d06bdec49bfa690dd34cf3205e28f34`
- Goal-state digest:
  `sha256:193be5b38881faebc349f9ae1d273e24fac5d5925a9a4402b24706394ffaeb3a`
- Inspected tracked revision:
  `1423dce5fab45ce5223caeb6a24791bf1a2cc3ff`

## Scope and method

This was a bounded, static, read-only audit of repository-wide agent safety and
authority boundaries. It covered path classification and publication,
credentials, infrastructure mutation, destructive operations, Bazel and Codex
sandbox boundaries, network access, generated artifacts, runtime extensions,
delivery, and contract discovery. No secret-bearing value was opened or
reported. No command with live-system effects was run.

For this audit, the user's clarified phrase “outside the repo” means outside
tracked or committable source, not outside the checkout. A short-lived secret
temporary file may therefore live under ignored task-private `out/<task>/`
when it has restrictive permissions, an explicit lifetime, and cleanup. The
phrase “outside the repo” in the secrets skill
(`projects/agents/skills/repo-secrets/SKILL.md:38-40`) is wording ambiguity,
not a substantive conflict with the root scratch policy; it should be replaced
with that precise formulation.

The audit did not run Bazel. Even `bazel_agent query` can update Bazel output
and cache state, which would violate this worker's instruction to write only
this report. The apparent documentation-publication path below is therefore a
strong static inference from the rule implementation and BUILD graph, not a
recorded query result. A graph test is included as a required follow-up.

## Direct assessment

The current system has strong *local* safety mechanisms but no repository-wide,
machine-enforced action contract. Safe behavior depends on an agent remembering
which prose to read, routing to the right skill, distinguishing several kinds
of sandbox, recognizing hazardous target names, and not exposing sensitive
command output. That is too much latent state for a zero-context agent.

The highest-leverage correction is one canonical, machine-readable policy and
capability graph, compiled into a small path/task context view and enforced at
every effectful execution gateway. It must answer, before execution:

1. who owns the affected path and external resource;
2. what information classes may be read, logged, exported, built, or published;
3. whether the action is observation, hermetic computation, source mutation,
   ignored local-state mutation, credential use, network access, remote
   mutation, or destruction;
4. what exact authority, environment, target, and candidate digest are bound;
5. what bounded evidence will prove the postcondition; and
6. what expected resource cost and reusable cached evidence already exist.

The repository should build this on its existing strongest patterns:
`repo_delivery`'s inspect/prepare/receipt/publish/verify protocol, the goal
store's immutable criteria and digest-bound evidence, Bazel's strict action
environment and default network denial, and bounded runtime-tool output.

## Existing controls worth preserving

- The root guide clearly distinguishes questions from authorization, requires
  task-specific scratch, warns about infrastructure mutation and BEP leakage,
  and directs agents to narrow validation (`AGENTS.md:22-34`, `58-68`,
  `145-179`, `181-203`).
- Bazel build actions use a strict environment and deny sandbox network by
  default; the agent configuration also caps output and test output
  (`tools/bazelrc/preset.bazelrc:109-114`, `171-177`;
  `tools/bazelrc/project.bazelrc:38-60`).
- Domain skills contain valuable safety knowledge. In particular,
  `repo-secrets` prohibits printing secret material
  (`projects/agents/skills/repo-secrets/SKILL.md:8-45`), and the Terraform and
  Ansible skills distinguish validation from live mutation
  (`projects/agents/skills/repo-terraform/SKILL.md:36-54`;
  `projects/agents/skills/repo-ansible/SKILL.md:38-52`).
- `repo_delivery` is an unusually strong model for exact authority: it
  sanitizes repository identity, binds preparation to exact paths, refs,
  endpoint digests and OIDs, requires post-prepare validation, uses exact
  leases, and rechecks before and after remote mutations
  (`tools/repo_delivery/README.md:17-27`, `29-57`, `77-108`, `123-162`,
  `164-228`).
- The goal store separates canonical records from projections, uses
  compare-and-swap resource versions, digest-binds evidence, and isolates
  delegated-worker output (`projects/goal/cmd/goal/README.md:11-17`, `83-103`,
  `123-143`).
- Secret scanners exist: optional pre-commit hooks detect AWS credentials and
  private keys (`.pre-commit-config.yaml:21-25`), while TruffleHog has normal
  HEAD-history and manual all-ref targets
  (`tools/trufflehog/BUILD.bazel:10-47`).
- The MCP repository reader has strong byte, path, symlink, timeout, and output
  bounds (`projects/mcp_cordis/plugins/repo_context.mjs:13-16`, `47-65`,
  `127-175`, `1045-1175`). Scratch runtime packages are initially disabled
  unless explicitly activated
  (`projects/mcp_cordis/internal/runtime.mjs:719-752`).

These mechanisms should become implementations of shared system contracts, not
remain exceptional islands.

## Ranked findings

### 1. Critical: effectful commands are flattened into an unguarded `bazel run`

**Verified facts**

- The required `bazel_agent` runner explicitly performs no command validation;
  it passes every argument through to Bazel
  (`projects/agents/skills/bazel-agent/SKILL.md:20-25`).
- The common runner macro simply creates ordinary executable targets; it has no
  effect, authorization, credential, environment, or destructive-operation
  contract (`projects/al/rules/al/al_binary.bzl:19-27`, `46-62`).
- The Terraform map makes the *unnamed/default* target run `apply`, alongside
  `migrate`, `import`, `destroy`, `state`, `force-unlock`, and other effectful
  commands (`tools/terraform/defs.bzl:28-63`). The root package instantiates
  that map as `//:tf` without narrowing it (`BUILD.bazel:26-29`).
- The Vault map exposes secret-reading commands (`kv get`, token lookup) and
  secret mutations (`kv put`, patch, delete) as the same kind of ordinary
  executable target (`tools/vault/defs.bzl:24-50`). The root package
  instantiates that complete map (`BUILD.bazel:31-34`).
- The release executable defaults its command to `deploy`, and the project
  aggregate exposes `//projects:deploy_heads`
  (`tools/release/main/bzl/al_release_binary.bzl:57-84`;
  `projects/BUILD.bazel:24-51`).
- Source-writing workflows likewise use ordinary executables and do not share
  one reliable name: examples include `-write`, `.update`, and `write_*`
  (`tools/bzl/main/bzl/al_genquery_write_to_source_file.bzl:73-77`;
  `tools/helm/main/bzl/al_helm_chart_lock.bzl:30-34`;
  `projects/goal/docs/BUILD.bazel:11-27`).

**Why it matters**

A label typo, an assumed default, or failure to route to one domain skill can
turn a validation attempt into a credential read or remote mutation. Bazel's
action sandbox is not a defense for a program launched by `bazel run`, and
target tags are not an authorization system. Naming and prose are doing work
that an execution boundary must do.

**Architectural implication**

Introduce a mandatory `AgentAction` contract for every runnable target, with at
least:

- effect class: `observe`, `compute`, `local_state`, `source_write`,
  `network_read`, `remote_write`, or `destroy`;
- input and output information classes;
- credentials used, without values;
- permitted network destination class;
- external environment and resource selector;
- required authority granularity;
- preflight, validation, and verification contracts;
- expected cost and cacheability; and
- owning component and escalation skill.

Generate this metadata in the owning Starlark macros so it cannot drift from
the actual command. Extend the agent execution gateway—not merely a skill—to
refuse unknown runnable targets and to require exact authority for anything
beyond observation/hermetic computation. Remove mutating unnamed aliases such
as `//:tf`; a mutating label must say exactly what it does.

For remote writes and destruction, generalize the delivery protocol:

```text
inspect -> prepare/plan -> immutable candidate digest -> validate
        -> exact scoped authority -> execute -> verify -> durable receipt
```

The fast path for declared local read-only build/test/query work should remain
cheap. An emergency escape hatch should be explicit, narrow, logged, and never
treated as proof that the operation was safe.

### 2. Critical: path policy dimensions are conflated

**Verified facts**

- The root guide first says all checked-in source and fixtures are public
  information that may be sent to an explicitly requested external service
  (`AGENTS.md:10-20`). It later says every file under `infra/` and `users/`,
  and secret-bearing `data/`, is sensitive (`AGENTS.md:181-203`).
- Tree READMEs use “public” and “published” as undifferentiated binary terms:
  `infra/`, `data/`, `third_party/`, and `users/` say they must not be public
  or published (`infra/README.md:12-16`; `data/README.md:12-16`;
  `third_party/README.md:12-16`; `users/README.md:12-16`). The root guide also
  uses “public” for Bazel-like visibility and repository confidentiality in the
  same document (`AGENTS.md:10-20`, `124-143`).
- The rule implementation makes each `docs_filegroup` package its Markdown and
  recursively included documentation dependencies
  (`projects/rules_docs/docs/defs.bzl:5-53`). The root documentation aggregate
  includes every repository subpackage except `bazel-*`
  (`BUILD.bazel:92-100`), and each non-publishable top-level tree exposes its
  aggregate to the root (`infra/BUILD.bazel:42-49`; `data/BUILD.bazel:3-10`;
  `third_party/BUILD.bazel:13-20`; `users/BUILD.bazel:9-16`). The public-site
  source package consumes `//:docs`
  (`projects/alwaldend.com/BUILD.bazel:113-125`). This is an apparent direct
  publication path from non-publishable trees.
- The declared “must not be used in builds” boundary is not represented as a
  graph invariant. For example, `//infra:al` is visible to all subpackages
  (`infra/BUILD.bazel:35-40`), and a `third_party` configuration depends on it
  (`third_party/BUILD.bazel:6-11`).

**Why it matters**

An agent cannot infer whether “private” means confidential to the model,
private Bazel visibility, excluded from public artifacts, operationally
sensitive but present in a public Git repository, or forbidden as a production
dependency. Those decisions require different controls. The static
documentation graph suggests the prose boundary may already be violated.

**Architectural implication**

Create one canonical path-policy manifest with independent axes:

- source confidentiality and allowed model/external-service disclosure;
- log and evidence handling;
- Bazel target visibility;
- allowed build consumers;
- allowed publication/export destinations;
- credentials and live-environment association; and
- path owner and required reviewer.

Generate the prose tree summaries and agent context projection from it. Enforce
the build and publication axes with Bazel aspects/tests over actual dependency
graphs. At minimum, add tests proving that public release/site roots cannot
reach forbidden trees and that production targets cannot reach `tools/`,
`infra/`, or `users/` except through explicitly modeled, reviewed exemptions.

Before changing policy, run a recorded `somepath`/reverse-dependency audit to
determine whether the documentation path is intended. If publishing these
READMEs is intentional, narrow the tree policy to the exact non-publishable
artifacts; do not leave a false blanket guarantee.

### 3. Critical: ordinary output can carry credentials

**Verified facts**

- The root policy says not to paste state, plan output, inventories, decrypted
  configuration, or credentials into logs (`AGENTS.md:181-203`).
- The Terraform skill nevertheless presents `tf.plan` as a normal validation
  command and only warns afterward that it can require Vault, backend,
  provider, network, or cloud access
  (`projects/agents/skills/repo-terraform/SKILL.md:36-49`). Terraform plans can
  print operationally sensitive values before an agent gets the chance to
  summarize them.
- Vault `kv_get` and `token_lookup` are ordinary Bazel-run targets with no
  redacting output adapter (`tools/vault/defs.bzl:24-50`).
- The broad `.gitignore` covers state and variable files but contains only a
  commented suggestion for Terraform plan files, not an active plan-artifact
  rule (`.gitignore:15-49`). The root guide correctly says ignore rules are
  only a backstop (`AGENTS.md:201-203`).
- Optional pre-commit hooks and the separate TruffleHog target are not part of
  the root `repo_quality_test` or `repo_delivery` candidate gate
  (`.pre-commit-config.yaml:1-30`; `BUILD.bazel:61-64`;
  `tools/trufflehog/BUILD.bazel:10-47`).
- By contrast, the full-repo checker already demonstrates the safer pattern:
  keep raw command logs at mode `0600`, inspect only targeted excerpts, and
  never add BEP or environment dumps
  (`projects/agents/skills/full-repo-check/SKILL.md:38-50`).

**Why it matters**

“Do not quote the secret in the final answer” is too late once a command has
printed it into a model-visible tool result, transcript, terminal log, or
shared artifact. Optional commit-time scanning also cannot prevent disclosure
during diagnosis.

**Architectural implication**

- Make secret retrieval injection-only for agent workflows. Offer metadata-only
  inspection (existence, version, policy, expiry, digest) and never a generic
  agent target that emits the value.
- Run Terraform plans and similar credentialed diagnostics through a typed
  adapter that captures raw output in a restricted per-task artifact, returns a
  bounded redacted structural summary, and records additions/changes/destroys,
  replacements, environment identity, and candidate digest. Do not rely solely
  on regular-expression redaction; prefer producer-native structured output and
  allowlisted fields.
- Put generated plans under task scratch by construction and refuse source-tree
  output paths. Retain or destroy them according to an explicit classification
  and retention policy.
- Add a non-echoing, candidate-diff secret scan to the exact post-prepare
  delivery receipt. Keep the full-history/all-ref scan as a less frequent audit
  so normal delivery does not repeatedly pay that cost.
- Treat a detected credential as an incident: suppress value display, bind the
  minimal evidence, and point to rotation/revocation work without copying it.

### 4. High: runtime-extension authorities are conflated

**Verified facts**

- Trusting the workspace causes a required MCP server to be launched from
  checked-in shell code (`.codex/config.toml:1-9`;
  `projects/mcp_cordis/cmd/mcp_cordis/launch.sh:1-20`).
- Its own documentation correctly states that package code is trusted, has
  normal Node built-ins and structured process execution, and is a reliability
  boundary rather than a security sandbox
  (`projects/mcp_cordis/README.md:99-128`).
- `cordis_define` accepts either `scratch` or `project` scope and persists
  code; `cordis_promote` copies scratch code into maintained project source
  (`projects/mcp_cordis/internal/mcp.mjs:83-102`, `223-243`). Direct project
  definition bypasses the otherwise valuable promotion distinction.
- Enabling arbitrary stored code is annotated `destructiveHint: false`, reload
  is likewise non-destructive, and the generic `cordis_invoke` has no effect
  annotation even though a handler can execute processes or contact the network
  (`projects/mcp_cordis/internal/mcp.mjs:104-135`, `155-189`).
- `http_probe` accepts arbitrary destinations and request headers, can follow
  redirects, can return a one-megabyte body preview, and has an option to
  expose sensitive response headers
  (`projects/mcp_cordis/plugins/network_probe.mjs:397-438`, `441-595`). There
  is no checked address-class policy preventing access to loopback,
  link-local, metadata, or private-network endpoints.
- The repository context tools can search or read any explicitly selected
  regular file inside the workspace; their boundary is path containment, not
  information classification
  (`projects/mcp_cordis/plugins/repo_context.mjs:783-1041`, `1045-1175`).

**Why it matters**

The runtime extension mechanism is exactly the kind of accretive capability the
repository needs, but it currently makes “load code,” “trust code,” “grant
filesystem/process/network power,” “mutate maintained source,” and “invoke one
handler” nearly the same authority. Generic MCP annotations cannot accurately
describe a dynamically selected handler's effects. Scratch code also lives in
the same ignored `out/` namespace used for untrusted downloads and logs.

**Architectural implication**

- Make scratch definition inert by construction and separate definition from
  activation authority. Treat activation and reload as execution of code with
  the declared maximum effects, never as non-destructive metadata changes.
- Remove direct maintained-project writes from the live runtime API. Promotion
  should create a normal task-owned candidate that passes project layout,
  package validation, review, and repository delivery.
- Require every dynamic handler to publish the same `AgentAction` contract as a
  Bazel runnable. The broker must evaluate the selected handler contract before
  invocation; a generic invoke endpoint should advertise worst-case effects or
  issue a short-lived invocation receipt after catalog inspection.
- Run extensions in a capability sandbox with explicit filesystem roots,
  executable allowlists, inherited-environment allowlists, network destination
  classes, time/output/process limits, and no credentials unless declared and
  authorized. Current process-group cleanup is useful but is not a security
  boundary.
- Default repository reads to tracked/publicly classified source. Sensitive
  trees and ignored/untracked files need an explicit, bounded read grant.
- Harden network diagnostics against SSRF and credential exfiltration: block or
  separately authorize local/private/link-local destinations, strip credential
  headers, never return authentication headers, and bind redirect destinations
  to the same policy.

### 5. High: no bounded authoritative context slice exists

**Verified facts**

- There is exactly one `AGENTS.md` in the repository. It says more deeply
  nested agent guides would override it, but none exist (`AGENTS.md:5-8`).
- The guide says each top-level tree README is the authoritative policy and
  tells agents to read the nearest README, BUILD file, and optional module
  include
  (`AGENTS.md:70-75`, `124-143`). Those files are not automatically loaded by
  the AGENTS hierarchy.
- The repository context helper returns applicable `AGENTS.md` files only; it
  does not return the README that the root policy calls authoritative
  (`projects/mcp_cordis/plugins/repo_context.mjs:614-671`, `674-751`).
- The human root README is only a title, three external links, and a license;
  it does not link to the repository's agent/system model or safety contract
  (`README.md:1-14`).
- Skill discovery is coherent—canonical skills are linked through
  `.agents/skills/` and the root BUILD target
  (`projects/agents/README.md:14-23`; `BUILD.bazel:66-90`)—but skill metadata
  describes routing, not the precise effects of the commands it recommends.
- Several command-owning READMEs provide almost no operational or safety
  contract, including Terraform, its runner, Vault, and the core `al` project
  (`tools/terraform/README.md:1-12`; `tools/terraform/runner/README.md:1-6`;
  `tools/vault/README.md:1-8`; `projects/al/README.md:1-11`).

**Why it matters**

The agent must reconstruct policy by chasing documents whose authority and
dimensions are unclear. Missing one step changes safety, while loading all of
them wastes context and encourages stale duplication. Skills help only after
the request has already been routed correctly.

**Architectural implication**

Provide one cheap, bounded command/API such as:

```text
repo agent context --path <path> --operation <intent>
repo agent explain-target <label>
```

It should compile, rather than duplicate:

- applicable path and information policy;
- owner and task-ownership status;
- relevant skills and why they trigger;
- target effects, credentials, network and environment;
- generated/source status and canonical updater;
- narrow validation targets and cost hints;
- current goal/attempt bindings when present; and
- exact policy-source links plus truncation markers.

Keep the root `AGENTS.md` as a short policy kernel: how to obtain context, the
non-negotiable authority rules, and what the broker guarantees. Generate thin
human and nested-agent projections from the canonical manifest when automatic
path discovery needs them. The root README should link directly to the same
system model. Mutable details must have one authority.

### 6. High: ignored `out/` mixes incompatible trust classes

**Verified facts**

- All downloads, reports, logs, extracted archives, caches, and temporary files
  are required to share task-specific `out/<task>/` storage
  (`AGENTS.md:22-34`). The entire namespace is broadly ignored
  (`.gitignore:91-93`).
- Full-repository audit logs may contain sensitive output and are protected
  with `0600` files and `0700` directories
  (`projects/agents/skills/full-repo-check/SKILL.md:38-46`).
- `repo_delivery` places a trusted consistency receipt and its lock in this
  namespace and warns that a same-user writer can race it
  (`tools/repo_delivery/README.md:89-108`).
- MCP Cordis loads executable scratch modules from `out/mcp_cordis`, while the
  same root is also its launcher scratch location
  (`projects/mcp_cordis/README.md:10-14`, `27-31`;
  `projects/mcp_cordis/internal/runtime.mjs:388-409`).
- Goal session focus also lives in ignored scratch, while durable attempt
  evidence is separately digest-bound
  (`projects/goal/cmd/goal/README.md:11-17`, `133-143`).

**Why it matters**

An ignore rule prevents accidental Git inclusion; it does not establish
provenance, integrity, ownership, confidentiality, executability, freshness, or
safe cleanup. Cross-agent work and long-lived sessions make “same user” a real
concurrency boundary. A downloaded or stale artifact should never become
executable or authority-bearing merely because it lives under `out/`.

**Architectural implication**

Retain the single ignored root, but make task storage typed and manifest-bound:

```text
out/<task>/
  manifest.json       # task/owner, producer, source digest, classifications
  evidence/           # immutable selected evidence
  logs/               # restricted, never executable
  cache/              # disposable, untrusted
  candidates/         # proposed source, never live-loaded
  state/              # locked workflow receipts and exact lifetimes
```

Use restrictive creation modes for all non-public artifacts, no-execute mounts
or execution refusal for generic scratch, atomic writes, per-task locks, size
budgets, retention/cleanup state, and producer/candidate digests. A consumer
must declare which artifact type and digest it accepts. Executable extension
code belongs in an isolated sandbox store and becomes maintained code only
through ordinary source promotion. Never auto-promote logs or caches into
durable learning.

### 7. Medium-high: ownership has no common resolver

**Verified facts**

- `repo_delivery` explicitly does not infer task ownership or choose
  validation (`tools/repo_delivery/README.md:11-13`), while its safety depends
  on the caller supplying only task-owned paths and confirming that history is
  not shared or human-owned
  (`projects/agents/skills/repo-delivery/SKILL.md:84-111`).
- `CODEOWNERS` contains only `- @simeonwarren` (`CODEOWNERS:1`), which does not
  express an obvious catch-all or path-specific ownership policy under common
  CODEOWNERS syntax. Its interpretation should be verified on the actual forge.
- Goal `local-owner-root` records storage ownership, not authority over source,
  infrastructure, credentials, or remote resources
  (`projects/goal/cmd/goal/README.md:65-72`).

**Why it matters**

The strongest mutation tool begins after the hardest decision: whether the
agent owns the paths, commits, remote ref, environment, and requested effect.
No shared resolver can currently answer that question, and a receipt can only
bind a mistaken answer more precisely.

**Architectural implication**

Add a common ownership projection combining:

- code owner/reviewer policy;
- project and tree owner;
- current worktree/branch and human-vs-agent commit provenance;
- task-declared path/hunk scope;
- current goal and attempt;
- external resource/environment owner; and
- exact user-authorized operation and scope.

The broker should return `owned`, `not_owned`, or `unknown`, with evidence. It
must fail closed on `unknown` for history rewrites, remote mutation,
destruction, credential scope expansion, and publication. Correct the
CODEOWNERS root rule if the current line is not intentional, then validate the
file in repository checks.

### 8. Medium: guardrails do not compose

**Verified facts**

- The agent Bazel configuration sets good defaults, but `user.bazelrc` is
  loaded last and later command options retain normal Bazel precedence
  (`tools/bazelrc/root.bazelrc:13-17`;
  `projects/agents/skills/bazel-agent/SKILL.md:20-25`).
- Network denial is a Bazel *action sandbox* default, not a statement about
  repository rules, long-lived MCP servers, or programs executed through
  `bazel run` (`tools/bazelrc/preset.bazelrc:171-177`).
- Pre-commit installation is optional and agents must remember to run relevant
  checks themselves (`AGENTS.md:236-239`).
- There is no repository test found that compiles the prose tree policies into
  visibility, dependency, publication, effect, credential, or network
  invariants.

**Why it matters**

Each individual control is reasonable, but none describes the effective
combination of platform sandbox, Bazel action sandbox, runtime process power,
credential injection, external-system permissions, and user authority. Agents
can easily believe one layer protects another layer that it does not cover.

**Architectural implication**

Define the layers explicitly and make the preflight report their intersection:

1. conversation/user authority;
2. host-agent filesystem/network/approval sandbox;
3. repository path and information policy;
4. Bazel analysis/action isolation;
5. runtime capability and credential grant; and
6. external-system authorization and concurrency state.

Do not claim security from the weakest advisory layer. Put invariant checks in
the execution gateway and build graph, with skills explaining the result and
helping the agent recover. Add negative/fault-injection tests for unknown
effect metadata, stale receipts, wrong environment, path-policy conflicts,
secret output, private-network redirects, weakened rc flags, and concurrent
writers.

## Canonical safety contract proposed for the system architecture

One resource should own the following normalized records; every human document,
skill routing summary, Bazel provider, runtime-tool annotation, and validation
test should be a projection or consumer of it.

### Path policy

```yaml
path: infra/**
owner: infra
source_class: operational_sensitive
model_read: explicit_bounded
external_egress: deny
build_consumers: [infra_validation]
publish: deny
log_class: restricted
```

The concrete schema may differ, but the dimensions must not collapse into one
word such as `private`.

### Action contract

```yaml
label: //infra/example:tf.apply
owner: infra
effect: remote_write
environment: explicit
credentials: [vault_reference_only]
network: declared_provider_endpoints
inputs: [source, prepared_plan_digest]
outputs: [redacted_summary, operation_receipt]
authority: exact_target_environment_candidate
preflight: inspect
verify: post_state_digest
cost: high
```

Unknown fields or unknown runnable targets must fail closed for credential use,
network access, source mutation, and remote effects.

### Evidence and learning contract

Every accepted improvement should link:

- exact source/candidate identity;
- policy and action-contract versions;
- authority receipt for effects;
- bounded validation evidence and known coverage gaps;
- observed cost;
- failure classification without sensitive payloads; and
- the skill, test, contract, or architectural record changed when the lesson is
  durable.

This makes learning accretive without turning transcripts, raw logs, caches, or
secrets into permanent repository memory.

## Prioritized implementation consequences

1. **Stop adding new unclassified run targets.** Define the action/effect
   schema and annotate the common `al`, Terraform, Vault, Ansible, release,
   generator, goal, and delivery macros. Remove dangerous default aliases.
2. **Resolve path-policy contradictions.** Split confidentiality, build use,
   visibility, and publication into explicit axes; add graph checks before
   changing any claimed boundary.
3. **Protect observation channels.** Replace raw Vault/plan output with typed,
   restricted adapters and put a candidate-only secret scan in delivery.
4. **Put an enforcing broker in front of effects.** Reuse exact receipts and
   postcondition checks; expose a fast read-only path and fail closed on
   unknown or stale state.
5. **Constrain runtime accretion.** Make scratch code inert, sandbox
   activation, model handler effects, remove direct project writes, and harden
   network and sensitive-path access.
6. **Compile zero-context views.** Add `agent context` and `explain-target`,
   generated from the same canonical graph, then shorten root prose to the
   immutable kernel.
7. **Type task artifacts and ownership.** Separate cache/log/evidence/
   candidate/workflow-state semantics under `out/`; add a common ownership
   resolver.
8. **Close the loop.** Record safety refusals, near misses, and resource cost
   as bounded evidence; promote only generalized lessons into contract tests,
   skills, or maintained architecture.

## Acceptance signals

- Every runnable repository target has complete, validated effect metadata;
  no agent execution of an unknown target is possible through the supported
  gateway.
- No mutating operation is reachable through an unnamed/default label.
- A clean-root agent can obtain the applicable policy, owner, skills, action
  effects, validation, and cost hint in one bounded response without opening
  multiple prose files.
- Build-graph tests prove forbidden production dependencies and publication
  paths absent, including a recorded check of `//:docs` and the public site.
- Credential and Terraform workflows return only typed redacted summaries to
  the agent; test canaries never appear in tool output, transcripts, or
  promotable evidence.
- Every remote mutation consumes exact target/environment/candidate authority,
  detects stale state, and produces a verified receipt.
- Dynamic extensions cannot read sensitive paths, use credentials, contact
  private networks, execute processes, or write maintained source unless those
  powers are declared and explicitly granted.
- Task artifacts have type, owner, classification, source digest,
  permissions, and retention state; generic scratch is never executable or
  authority-bearing.
- Safety checks add negligible overhead to read-only/narrow build work and
  reuse candidate-bound evidence instead of rerunning expensive checks.

## Adversarial review

The strongest case for the current distributed design is that the root guide is
already explicit, domain skills contain the missing detail, Bazel denies
network to actions, the host sandbox supplies approvals, and risky labels such
as `.apply` or `.destroy` are self-describing. A universal broker could add
latency, stale metadata, and a false sense of security.

That case does not survive the direct counterexamples:

- `//:tf` means `apply`, so not every risky label is self-describing.
- `bazel_agent` intentionally validates no command.
- Bazel action network policy does not constrain an executed binary or MCP
  runtime.
- raw Vault and Terraform output can be sensitive before prose-based redaction;
- the documentation build appears to cross a declared non-publication boundary;
  and
- dynamic invocation cannot accurately inherit one static MCP annotation.

The answer is not a large policy engine on every read. It is generated
metadata, a cheap compiled context slice, a fast path for declared read-only
work, and strict two-phase gates only where effects or information boundaries
demand them. Metadata drift is addressed by generating it in the same macros
that create commands and by testing the complete runnable/publication graphs.

**Audit recommendation to the coordinator: revise.** Do not treat the current
distributed documentation as a sufficient repository-wide safety boundary.
Preserve its strongest mechanisms, but make their contracts canonical,
discoverable, compositional, and enforced before designing further agent
automation around them.
