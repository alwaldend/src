---
title: Readme tree
description: Tool to parse README.md files
languages:
  - go
  - bzl
tags:
  - bzl_rules
---

## Example

```sh
bazel run tools/readme_tree -- parse -g -C "${PWD}" .
```
