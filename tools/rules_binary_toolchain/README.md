---
title: Rules binary toolchain
description: Bazel toolchains for packaged executable binaries
statuses:
  - active
languages:
  - bzl
tags:
  - bzl_rules
---

`rules_binary_toolchain` creates Bazel toolchains and runnable targets for
packaged executable binaries. Archive entries can include runtime files,
which the generated binary targets expose through Bazel runfiles.
