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
path, label, or task. It joins the six checked catalogs plus the applicable
`AGENTS.md` authority document and renders JSON (default) or Markdown
(`--json=false`) from the same data. Missing catalogs or documents produce
`Completeness: partial` with structured limitations rather than failure.

```sh
bazel_agent bazel run //tools/agents/cmd/agent_system:agent_system_update -- \
  --workspace-root "$PWD"
```

The capsule carries provenance (observedAt, freshness, completeness,
truncation, next discovery actions), identity (repository, worktree path,
revision, input digests, bounded byte size), and bounded component,
capability, check, and provider views. Verification:

```sh
bazel_agent bazel test //tools/agents/cmd/agent_system:agent_system_check
```

## Runtime control kernel and status

The `tools/agents/control` package owns the fixed runtime control kernel:
per-package states and deadlines, account/lock acquisition, expires-checked
leases, and expected-revision asset publication with namespace isolation. It
persists an offline-readable `packages.json` snapshot.

The `control_status` command renders the snapshot/asset status for a control
root (Markdown with `--markdown`):

```sh
bazel_agent bazel run //tools/agents/cmd/control_status:control_status -- \
  --workspace-root "$PWD" [--markdown]
```

Verification:

```sh
bazel_agent bazel test //tools/agents/control:all //tools/agents/cmd/control_status:all
```
