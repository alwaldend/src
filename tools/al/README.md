---
title: Repository AL configuration
description: Shared configuration for repository command wrappers
---

The reusable executable, APIs, rules, and lifecycle skills belong to
[AL](../../projects/al/README.md). This package keeps the monorepo integration:
shared Lua configuration, the explicit `//tools/al:install` entry point, and
native command tests that use the parent repository's Vault injector.

`//tools/al:config` packages the repository's shared AL configuration. The Lua source files stay
at the repository root so existing CLI configuration discovery and
`require("al_lib")` imports retain their paths.

The `//tools/al:tf.*` targets select Terraform commands declared here
with the public rules_terraform command maps and the generic AL wrapper. These
bindings use the shared config and explicitly select the repository root as
their Terraform working directory.
