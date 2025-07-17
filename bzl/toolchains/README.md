---
title: Toolchains
description: Bazel toolchains
cascade:
  - categories:
      - bzl
      - bzl-toolchains
---

## Links

- https://bazel.build/extending/toolchains

## Qt setup

- Install qt: `aqt install-qt -O /opt/qt linux desktop 6.9.0`
- Register toolchain: `register_toolchains("//bazel/qt:preinstalled-qt-toolchain")`
