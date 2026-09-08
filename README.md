---
title: Src
description: Monorepo source and repository control system
---

# Source repository

This monorepo contains first-party projects, repository tooling,
infrastructure definitions, data, and pinned external sources. Facts remain
with their owning component; the repository agent system connects those facts
without replacing their authorities.

## Start here

- [Agent policy](AGENTS.md) is the repository-wide operating contract. A
  nearer `AGENTS.md` takes precedence within its subtree.
- [Agent-system current state](projects/agents/docs/current-state.md) records
  the evidence-backed baseline.
- [Agent-system architecture](projects/agents/docs/architecture.md) defines
  how intent, policy, capabilities, work, execution, evidence, delivery, and
  learning compose.
- [Agent-system roadmap](projects/agents/docs/roadmap.md) distinguishes
  proposed work from current guarantees.
- [Repository evolution](infra/src/openspec/README.md) describes changes to the monorepo
  itself. Each project and infrastructure component owns its specifications
  and maintained changes in its local `openspec/` directory. The
  [goal migration](infra/src/openspec/migration.md) preserves prior acceptance state and
  evidence provenance.
- [Repository setup](https://alwaldend.com/docs/misc/repo/) explains the human
  development environment.

The intended control flow is:

```text
request -> owner and policy -> capability -> work -> execution
        -> evidence -> delivery -> reviewed learning
```

## Repository shape

- [`projects/`](projects/README.md): products and reusable project code.
- [`tools/`](tools/README.md): repository-wide rules, toolchains, and
  repository-internal automation.
- [`infra/`](infra/README.md): infrastructure definitions and operations.
- [`data/`](data/README.md): repository-owned data and documentation assets.
- [`third_party/`](third_party/README.md): vendored and externally sourced
  inputs.
- [`users/`](users/README.md): user-owned code and infrastructure.
- [`projects/agents/`](projects/agents/README.md): repository-wide agent
  architecture, skills, and evaluation assets.
- [`infra/src/`](infra/src/README.md): specifications and changes for evolution
  of the repository itself.

Each linked tree README is authoritative for that tree. Bazel `BUILD` and
`MODULE.bazel` files remain authoritative for executable and dependency
structure; generated maps and catalogs are projections, not new sources of
truth.

## Bazel package loading

The root `MODULE.bazel` is generated with `bazel run //tools/bazel_module:update`.
Maintain dependency pins in their owning `include.MODULE.bazel` files; the
generator discovers them and the repository quality suite checks freshness.
See [root module generation](tools/bazel_module/README.md) for discovery and
ordering boundaries.

The root `BUILD.bazel` contains source exports, compatibility aliases, test
suites, and repository-wide Gazelle directives. Tool implementations live in
their owning packages so loading a root configuration file does not evaluate
unrelated language dependency generators.
`//tools/repo_quality/test/root_build:root_build_test`, included in
`//:repo_quality_test`, enforces the absence of `load()` statements.

Existing commands such as `//:gazelle`, `//:buildifier`, `//:requirements.update`,
`//:vault`, and `//:tf.plan` remain available. Shared AL configuration is owned
by `//tools/al:config`; repository documentation is assembled by
`//projects/alwaldend.com:repo_docs`.

## External links

- Homepage: https://alwaldend.com/
- Docs: https://alwaldend.com/docs/

## License

[AGPL-3.0](https://spdx.org/licenses/AGPL-3.0-or-later.html), see
[LICENSE.txt](data/license/LICENSE.txt).
