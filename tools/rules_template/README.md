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
actions. The `to_html_table` helper accepts up to 128 levels of nested values
and reports a rendering error when an input exceeds that depth. One invocation
accepts up to 64 MiB of template and data files combined, with up to 1,048,576
data lines. Rendering stops before buffered output exceeds 64 MiB. A size or
depth failure leaves the destination untouched.

## Getting started

`MODULE.bazel`:

```py
bazel_dep(name = "rules_template", version = "<VERSION>")
register_toolchains("@rules_template//rules:all")
```

`BUILD.bazel`:

```py

```
