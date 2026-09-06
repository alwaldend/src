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

Agent tests use Bazel's normal result cache. Tests that inspect undeclared
workspace files must opt out with the `no-cache` tag.

Agent commands use Bazel's default lockfile update mode. Review and commit
generated lockfile changes with dependency changes. CI retains strict lockfile
checking; use `--lockfile_mode=error` for an explicit reproducibility check.

In the root workspace, lint actions fail on violations. Their failure setting
is identical across normal and lint configurations to preserve the analysis
cache when switching modes. `--config=lint` enables the linter aspects and
requests only lint reports and skill validation, not ordinary target outputs.

## Links

- Rules: https://github.com/bazel-contrib/bazelrc-preset.bzl/
