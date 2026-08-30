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

This project owns repository-wide Codex skills and their development-time
evaluation assets. Product-specific skills live with their owning project at
`projects/<project>/skills/<name>`. Every canonical skill is packaged as a
`skill_library` in its owning directory.

The repository discovery directory `.agents/skills/` contains one relative
symlink per skill. Each link points directly to its canonical project-owned
directory. Bazel ignores the discovery directory and builds only canonical
targets, preventing duplicate packages while allowing skills from more than
one owning project.

Skill evaluation data is not part of the runtime `skill_library` unless a
skill explicitly declares otherwise. Every new or updated skill includes an
offline Promptfoo configuration-validation target. That target checks the eval
harness, not answer correctness. Live behavioral evaluations are manual,
networked tests and must not be included in ordinary wildcard test runs; they
may be omitted when representative coverage requires tool calls or external
state that cannot be provided safely and reproducibly, with the gap documented
beside the eval configuration.
