---
title: Qt
description: Qt rules
tags:
  - bzl_rules
---

## Links

- https://bazel.build/extending/toolchains

## Qt setup

- Install qt: `aqt install-qt -O /opt/qt linux desktop 6.9.0`
- Register toolchain: `register_toolchains("//tools/qt:preinstalled_qt_toolchain")`
