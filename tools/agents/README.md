---
title: Agent control contracts
description: Repository-internal shared contracts for agent-system composition
statuses:
  - experimental
languages:
  - go
tags:
  - agent
  - control
---

# Agent control contracts

This project owns the repository-internal shared vocabulary used to compose
agent-system declarations and derived catalogs. It does not own component
payload schemas, user authority, runtime observations, or generated catalogs.

The initial `api/v1alpha1` package defines stable references, atomic effects,
authority and budget envelopes, independent path-policy axes, information and
retention classes, availability reasons, evidence applicability, and the
common operation and artifact envelopes. Validation fails closed for unknown
effects and malformed identities. JSON decoding is strict and canonical
encoding is deterministic.

The Phase 1 registry names the closed registration authorities, declares skill
metadata and supported direct binaries, links owner-local operation files, and
records generated-artifact ownership. Run the report-only completeness check
with:

```sh
bazel_agent bazel run //tools/agents/cmd/phase1_check -- \
  --workspace-root "$PWD" \
  --report out/<task>/phase1-report.json
```

The adjacent criteria-revision-bound resource baseline records numeric
ceilings separately from observations. Unavailable observations carry a
reason instead of an estimate.

Project roots use exactly one lifecycle value from `active`, `in_progress`,
`maintenance`, `experimental`, `finished`, or `abandoned`. These values
describe project maintenance state; they do not imply publication or
information policy.

## Bounded topology catalog

The topology compiler derives a deterministic, offline `TopologyCatalog`
from the registered `projects/*` universe, top-level boundary READMEs, and
tracked `MODULE.bazel` roots. The portable JSON at
`tools/agents/catalogs/topology.json` is authoritative; the adjacent Markdown
renders the same data and states the JSON digest.

Regenerate with:

```sh
bazel_agent bazel run //tools/agents/cmd/topology_check:topology_update
```

Verify the checked artifacts are current with:

```sh
bazel_agent bazel test //tools/agents/cmd/topology_check:topology_check_check
```

The catalog contains no timestamps, absolute checkout paths, or nested
authority; every input is bound by a content digest and the catalog digest
covers the canonical bytes. Completeness is `complete` only when every
eligible registered project, workspace, and boundary tree is emitted without
conflicts or limitations.

## Bounded context capsule

The `agent_system` command emits a bounded offline context capsule for a
workspace-relative path or local Bazel label. It joins checked catalogs with
the applicable `AGENTS.md` chain, nearest owner README, and policy source
paths. Its digest-bound document links include the nearest package
`BUILD.bazel` or `BUILD`, `include.MODULE.bazel`, and workspace `MODULE.bazel`
when present. These locate the source to inspect; they do not select targets.
JSON is the default; `--markdown` or `--json=false` renders Markdown.
Existing symlinks resolve to their canonical workspace paths before ownership
and policy selection. Escaping, dangling, or cyclic symlinks are refused;
new paths may extend a verified existing parent.

```sh
bazel_agent bazel run //tools/agents/cmd/agent_system -- \
  --workspace-root "$PWD" --path projects/agents --markdown
```

The legacy `agent_system_update` wrapper also forwards these query flags and
does not update tracked files. `--task` records a caller-supplied task ID; it
does not select a goal. Input digests bind accepted catalog and document inputs.

Relative `--workspace-root` values resolve against
`BUILD_WORKSPACE_DIRECTORY` under `bazel run`, so `--workspace-root .` means
the source workspace. Outside Bazel they resolve against the current
directory. `normalize`, `coverage`, and `aggregate` apply the same rule to
relative input, catalog, and output file flags, with optional
`--workspace-root` selecting another base. Absolute file paths retain their
meaning; output parent directories must already exist. For example:

```sh
bazel_agent bazel run //tools/agents/cmd/agent_system -- aggregate \
  --catalog tools/agents/catalogs/capability.json \
  --input tools/agents/catalogs/skill-cases.json \
  --output out/<task>/coverage.json
```

The context command observes local Git HEAD and dirty state under the selected
workspace, including untracked files and excluding ignored files. The
`identity.git` object separates these observations from optional caller
`--revision` and `--dirty-inputs` declarations. Without overrides, successful
observations populate the existing identity fields; `revisionSource` and
`dirtyInputsSource` name their basis. Failed Git reads leave `git.dirty: null`
or the revision absent with explicit reasons. Legacy identity fields then
retain an input digest and a conservative dirty default, clearly labelled.

Git commands have a combined three-second deadline and a 65,536-byte output
limit each, disable optional index writes and fsmonitor hooks, and ignore
ambient `GIT_*` variables. The report contains no changed-path list or Git
stderr. Git and catalog reads are separate observations, not an atomic
snapshot or proof that generated catalogs match HEAD.

The capsule reports the actual read time, unknown catalog freshness, and
partial completeness. It does not observe runtime health, infer authority,
resolve effective CODEOWNERS, or associate an unrelated active goal. Its
capability list is an inventory of root-workspace and selected-workspace
candidates, not task-intent routing or a complete cross-workspace dependency
closure. Workspace check phases are references, not instructions to run a
broad check.

`identity.byteSize` measures the complete canonical JSON, including its final
newline. `--max-bytes` caps the selected JSON or Markdown output at 65,536 bytes
by default; oversized output fails before writing a partial document. The
limit may be set explicitly between 1,024 and 1,048,576 bytes. Verification:

```sh
bazel_agent bazel test //tools/agents/cmd/agent_system:agent_system_check
```

## Runtime control kernel and status

The `tools/agents/control` package owns the fixed runtime control kernel:
per-package states and deadlines, account/lock acquisition, expires-checked
leases, and expected-revision asset publication with namespace isolation. It
persists an offline-readable `packages.json` snapshot.

Package registration preserves the caller's scope, desired revision, and
contract hash. Activation deadlines are fixed at registration time;
`KernelOptions.Deadlines` are fixed at kernel creation for packages introduced
through `Mark`. Reading status does not extend those deadlines. Lifecycle
transitions do not establish revision identity, so `observedRevision`
remains empty: the current API has no observed revision input.

The `control_status` command renders the snapshot/asset status for a control
root without creating directories or starting a kernel (Markdown with
`--markdown`):

```sh
bazel_agent bazel run //tools/agents/cmd/control_status:control_status -- \
  --workspace-root "$PWD" [--markdown]
```

The JSON `healthy` field describes the persisted package observations, with
`healthBasis: persisted-package-snapshot`. Ready, degraded, and disabled
packages count as healthy, matching the kernel. Missing, empty, malformed,
or rejected observations produce `healthy: null`; consumers must support
this unavailable value rather than assuming a boolean is always present.
`snapshotPath` identifies the source. The report's `observedAt` is its read
time; each package retains the writer's `observationTime`.

Freshness is `unknown` unless the caller supplies a positive duration such
as `--max-snapshot-age 5m`. That bound yields `within-age-bound` or `stale`;
stale observations make snapshot health unavailable. This is an explicit
caller policy, not a heartbeat or runtime liveness guarantee. Live runtime
health remains unavailable even when every observation is within the bound.
The snapshot format does not identify its writer, so report `runtimeId` is
empty with an unavailable reason. An asset's `runtimeId` identifies its
publisher only.

Verification:

```sh
bazel_agent bazel test //tools/agents/control:all //tools/agents/cmd/control_status:all
```
