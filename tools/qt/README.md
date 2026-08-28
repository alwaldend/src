---
title: Qt
description: Qt rules
languages:
  - bzl
tags:
  - bzl_rules
---

## Links

- https://bazel.build/extending/toolchains

## Hermetic Qt

The module configuration downloads SHA-256-pinned Qt 6.8.3 distributions via
`rules_qt` and registers its build tools. Bazel builds do not require Qt under
`/opt` or another host installation.
