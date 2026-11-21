---
title: Bend rules
description: Bazel rules for running and checking Bend 2 programs
statuses:
  - active
languages:
  - bzl
  - go
tags:
  - bzl_rules
---

`rules_bend` integrates the [Bend 2 compiler](https://github.com/bendlang/bend)
into Bazel. It requires Bazel 9.3.1 or later. Bazel selects the Bend executable
from a registered toolchain. Source checks select it for the execution
platform and runnable programs select it for the target platform. The module
registers Bend 2.0.40 and its standard library by default for Linux x86-64.

Run the compiler itself with `bazel run //projects/rules_bend:bend -- version`
or `bazel run @rules_bend//:bend -- version` in a standalone consumer.

Load the rules from `@rules_bend//:defs.bzl` in a standalone consumer or from
`//projects/rules_bend:defs.bzl` in this repository.

```starlark
load("@rules_bend//:defs.bzl", "bend_binary", "bend_check")

bend_binary(
    name = "hello",
    main = "hello.bend",
)

bend_check(
    name = "check_hello",
    main = "hello.bend",
)
```

For example, `hello.bend` can contain:

```text
import Base

def main() -> IO(Unit):
  IO.print("Hello from Bend")
```

Run the program with `bazel run //:hello`. Pass program arguments after
`--`, as in `bazel run //:hello -- "one argument"`. The launcher forwards
standard input, output, errors, and Bend's exit status.

Check the program with `bazel build //:check_hello`. The sandboxed action
invokes Bend with `--check-only` and emits `check_hello.checked` on success.

List local imports in `srcs`, including transitive imports. Filegroups can share
source sets between targets. The rules stage files at their workspace-relative
paths for import resolution. Include `LAWS.bend` when it participates in a
proof check. The check action runs offline.

## Register a custom Bend distribution

In a package that exposes the Bend executable and its standard library as a
Bazel executable target, declare the toolchain registrations.

```starlark
load("@rules_bend//:defs.bzl", "bend_toolchain")

bend_toolchain(
    name = "custom_bend",
    compiler = ":bend",
    compatible_with = ["@platforms//os:linux", "@platforms//cpu:x86_64"],
)
```

Register both targets in the consuming `MODULE.bazel`:

```starlark
register_toolchains("//tools/bend:custom_bend_exec", "//tools/bend:custom_bend_runtime")
```

## Validate the project specification

From the repository root, run
`bazel test //projects/rules_bend/openspec:validate_test`.
