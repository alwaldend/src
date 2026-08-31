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
bazel_agent run //tools/agents/cmd/phase1_check -- \
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
