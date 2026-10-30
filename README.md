---
title: Src
description: Projects, infrastructure, and development tools
---

## External links

- Homepage: https://alwaldend.com/
- Docs: https://alwaldend.com/docs/

## License

[AGPL-3.0](https://spdx.org/licenses/AGPL-3.0-or-later.html), see
[LICENSE.txt](data/license/LICENSE.txt).

## Start here

- [Repository setup](https://alwaldend.com/docs/misc/repo/) explains the
  development environment.
- [Agent policy](AGENTS.md) contains instructions for agents working here.
- [Repository evolution](infra/src/openspec/README.md) documents shared
  repository specifications and changes.

## Repository shape

- [`projects/`](projects/README.md): products and reusable project code.
- [`tools/`](tools/README.md): repository-wide rules, toolchains, and
  repository-internal automation.
- [`infra/`](infra/README.md): infrastructure definitions and operations.
- [`data/`](data/README.md): repository-owned data and documentation assets.
- [`third_party/`](third_party/README.md): vendored and externally sourced
  inputs.
- [`users/`](users/README.md): user-owned code and infrastructure.
- [`tools/agents/`](tools/agents/README.md): repository-wide agent
  architecture, skills, and evaluation assets.
- [`infra/src/`](infra/src/README.md): specifications and changes for evolution
  of the repository itself.

## Build documentation

- [Root module generation](tools/bazel_module/README.md): dependency include
  discovery and generated `MODULE.bazel` maintenance.
- [Repository quality](tools/repo_quality/README.md): formatting and validation.
- [Terraform rules](tools/rules_terraform/README.md): reusable Terraform rules
  and provider installation.
- [Repository documentation](projects/alwaldend.com/README.md): the website and
  repository documentation build.
