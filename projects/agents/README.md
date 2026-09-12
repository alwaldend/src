---
title: Agents
description: Repository-wide reusable agent skills
statuses:
  - active
languages:
  - markdown
  - bzl
tags:
  - agent
  - skills
---

# Repository agent skills

This project owns the reusable cross-repository agent skills. Maintained
improvement work lives in this project's [OpenSpec workspace](openspec/). It
does not centralize component facts or runtime state: each fact remains
canonical at its natural owner.

## Start here

| Document                                                                   | Purpose                                                 |
| -------------------------------------------------------------------------- | ------------------------------------------------------- |
| [OpenSpec changes](openspec/changes/)                                      | Maintained work, acceptance, and preserved goal history |
| [Root agent guide](https://github.com/alwaldend/src/blob/master/AGENTS.md) | Current repository-wide operating policy                |

Dated audits remain in
[migrated history](../../infra/src/openspec/migration.md).

## Ownership boundaries

- `openspec/` owns agent-system requirements, changes, and evidence;
  [the migration map](../../infra/src/openspec/migration.md) locates historical records.
- Other components own their specifications and changes under their own
  `openspec/` directory. `infra/src/openspec/` describes repository evolution.
- `skills/` owns reusable cross-product, repository-wide agent procedures and
  their development-time evaluations. A product-specific procedure belongs
  with its product; a procedure that no narrower project genuinely owns lives
  here rather than in a contrived local home.
- Product-specific skills remain with their project at
  `projects/<project>/skills/<name>`. Do not split one skill across owners.

## Skill packaging and discovery

Every canonical skill is packaged as a `skill_library` in its owning
directory.

The repository discovery directory `.agents/skills/` contains one relative
symlink per skill. Each link points directly to its canonical project-owned
directory. `.agents/BUILD.bazel` declares the complete discovery set and
generates those links with `//.agents:write_skill_links`; its generated
exact-state test verifies them. Bazel ignores the discovery directory and
builds only canonical targets, preventing duplicate packages while allowing
skills from more than one owning project. Each skill grants
`//.agents:skill_discovery` read access so the owning declaration can reach it.

Skill evaluation data is not part of the runtime `skill_library` unless a
skill explicitly declares otherwise. Every new or updated skill includes an
offline Promptfoo configuration-validation target. That target checks the eval
harness, not answer correctness. Live behavioral evaluations are manual,
networked tests and must not be included in ordinary wildcard test runs; they
may be omitted when representative coverage requires tool calls or external
state that cannot be provided safely and reproducibly, with the gap documented
beside the eval configuration.
