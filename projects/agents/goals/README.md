---
title: Repository agent goals
description: Durable goals, attempts, and evidence for the repository agent system
tags:
  - agent
  - goals
---

# Repository agent goals

Each child directory is an owner-local durable goal record. Canonical resource
files and attempt artifacts are created and validated through the repository
`goal` workflow; do not edit generated records or projections directly.

Goal evidence may contain source and operational analysis, including
non-secret live or generated facts. It must never contain credentials, other
secrets, or personal information. Inspect raw artifacts before promoting their
content.

The current goal's canonical record is under `repo-agent-system/` in
repository source. The maintained, link-safe public-system synthesis is the
[current-state document](../docs/current-state.md). Context-bound raw audits
remain in the canonical goal record rather than the documentation projection.
