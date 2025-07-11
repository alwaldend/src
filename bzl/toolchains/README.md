---
title: Bazel rules
---

## Qt setup

- Install qt: `aqt install-qt -O /opt/qt linux desktop 6.9.0`
- Register toolchain: `register_toolchains("//bazel/qt:preinstalled-qt-toolchain")`
