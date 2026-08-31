---
title: Bazelrc
description: Bazelrc files
---

## Usage

```sh
bazel run //tools/bazelrc:preset.update
```

The `agent` profile must not propagate ambient `TEMP`, `TMP`, or `TMPDIR`
values into repository rules, actions, host actions, or tests. Those execution
contexts use Bazel-managed temporary storage; host tools own any explicit
task/run scratch they require.

## Links

- Rules: https://github.com/bazel-contrib/bazelrc-preset.bzl/
