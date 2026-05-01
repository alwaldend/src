---
title: Docs
description: Documentation rules
languages:
  - bzl
tags:
  - bzl_rules
---

## Example

```bzl
load("//tools/docs:defs.bzl", "docs_filegroup")

docs_filegroup(
    name = "docs",
    srcs = glob(["*.md"]),
    deps = subpackages(include = ["*"], allow_empty = True),
    visibility = ["//:__pkg__"],
)
```
