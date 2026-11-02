---
title: AL Bazel rules
description: Configuration, command, Git, and release rules
---

AL rules and macros use the `al_` prefix; AL providers use the `Al` prefix.
Their source filenames follow the same prefix in snake case. The public API is
exported through `@rules_al//:defs.bzl`, with usage documented in
the [project README](../README.md#bazel-module-api).

The Git repository and toolchain rules support the `git` module extension.
Configured repository names, toolchain registration targets, and release target
names are chosen independently of rule and provider symbols.
