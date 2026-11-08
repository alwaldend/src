---
title: Repository AL configuration
description: Shared configuration for repository command wrappers
---

The reusable executable, APIs, rules, and lifecycle skills belong to
[AL](../../projects/al/README.md). This package keeps the monorepo integration:
typed configuration targets, the explicit `//tools/al:install` entry point, and
native command tests that use AL's external injector.

`//tools/al:config` composes the shared authentication, Vault connection, and
plugin targets declared in this package. Components compose it with their own
small authentication and plugin-call targets. Component auth extends
`//tools/al:auth`; aggregate list ordering does not select overrides. Short declarations stay inline. Extract
local `defs.bzl` only for significant reuse or readability gains from complex
inputs; common graphs use reusable target macros directly. Injector calls use
`al_injector_call` with named `res_<type>` maps. Plugin executables and certificate
files are declared with Bazel labels, and their runtime dependencies propagate
through `AlConfigInfo` targets. AL loads generated JSON; no Lua configuration
sources or Lua imports are used.

The `//tools/al:tf.*` targets select Terraform commands declared here
with the public rules_terraform command maps and the generic AL wrapper. These
bindings use the shared config and explicitly select the repository root as
their Terraform working directory.

Load `al_terraform_binary_map`, `al_terraform_target_binary_map`, and
`al_terraform_test_map` from `//tools/al:al_terraform.bzl` for repository
Terraform commands. Pass `configs` and `plugin_labels` directly; this integration
selects the AL wrapper once. Provider pins and operation behavior remain owned
by standalone rules_terraform. Plain offline tests can use its maps directly.

Backend calls reuse `vault_tf_backend(name = ..., path = ..., labels = ...)` from
`//tools/vault/tf_backend:defs.bzl`. The shared configuration registers the
backend executable once and supplies its runtime files to component commands.

The injector implementation and schema belong to `@rules_al//plugins/injector`.
Its reusable `:config` target owns registration and runtime files. Component
wrappers inherit configuration targets and do not repeat plugin dependencies
in `data`. The command cleanup test likewise consumes generated configuration
targets; its startup-failure variant inherits the normal fixture and overrides
one resource. Success, startup-failure and cancellation receipts are retained
in the test's undeclared outputs.

`//tools/al:injector_dev_server_test` retains the parent integration with the pinned
Vault server; AL's standalone tests also cover a synthetic local resource graph.
