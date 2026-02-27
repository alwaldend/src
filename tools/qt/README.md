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

## Qt setup

- Install qt:
  ```sh
  bazel run //tools/qt:install
  ```
- Register toolchain: `register_toolchains("//tools/qt:preinstalled_qt_toolchain")`
