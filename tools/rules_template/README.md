---
title: Rules template
description: Bazel rules and a Go command for rendering template files
statuses:
  - active
languages:
  - bzl
  - go
tags:
  - bzl_rules
---

`rules_template` provides Bazel rules and a Go command for rendering template
files. A registered toolchain supplies the templating executable to build
actions.

## Getting started

`MODULE.bazel`:

```py
bazel_dep(name = "rules_template", version = "<VERSION>")
register_toolchains("@rules_template//main/bzl:all")
```

`BUILD.bazel`:

```py

```
