---
title: Rules template
description: Rules to template files
languages:
  - bzl
  - go
tags:
  - bzl_rules
---

## Getting started

`MODULE.bazel`:
```py
bazel_dep(name = "rules_template", version = "<VERSION>")
register_toolchains("@rules_template//main/bzl:all")
```

`BUILD.bazel`:
```py
```
