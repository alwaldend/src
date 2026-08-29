---
title: Agents
description: Repository-owned agent skills and evaluation assets
statuses:
  - active
languages:
  - markdown
  - bzl
tags:
  - agent
  - skills
---

This project owns the repository's Codex skills and their development-time
evaluation assets. Each skill is packaged as a `skill_library` under
`skills/<name>`.

The repository discovery path `.agents/skills` is a direct symlink to this
project's `skills/` directory. Bazel ignores that symlink and builds only the
canonical targets under `//projects/agents/skills/...`, preventing duplicate
packages while keeping the standard Codex discovery location available.

Skill evaluation data is not part of the runtime `skill_library` unless a
skill explicitly declares otherwise. Every new or updated skill includes an
offline Promptfoo configuration-validation target. That target checks the eval
harness, not answer correctness. Live behavioral evaluations are manual,
networked tests and must not be included in ordinary wildcard test runs; they
may be omitted when representative coverage requires tool calls or external
state that cannot be provided safely and reproducibly, with the gap documented
beside the eval configuration.
