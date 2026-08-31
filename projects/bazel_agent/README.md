---
title: Bazel agent
description: Bazel runner for repository agents
statuses:
  - active
languages:
  - go
  - sh
tags:
  - bazel
  - agent
---

`bazel_agent` is the repository Bazel entry point for agents. It keeps the
ordinary `bazel` command unchanged while consistently invoking it in batch
mode with the agent configuration:

```text
bazel_agent test //path/to/package:all
    ->
bazel --batch test --config=agent //path/to/package:all
```

The runner performs no command validation. It passes an empty argument list,
`--help`, malformed commands, and every other argument through to Bazel. When a
first argument exists, it is placed in Bazel's command position and followed
by `--config=agent`; all remaining arguments are unchanged, including arguments
after the `--` separator of a `bazel run` invocation. Later options therefore
retain Bazel's normal precedence and can override settings supplied by the
agent configuration. The runner resolves `bazel` from `PATH` and replaces
itself with that process, so signals and the final exit status are not mediated
by another wrapper process.

The runner does not create or inject a host temporary directory. Bazel actions
use declared outputs and Bazel-managed temporary storage, and tests use their
test temporary-directory contract. Repository-updating host tools that need
scratch must accept or derive an explicit task/run path under `out/<task>/`
instead of relying on ambient `TMPDIR`, `TMP`, or `TEMP` values propagated to
the whole build.

`bazel_agent doctor --workspace-root PATH --task-scratch out/<task>/<run>` is
a read-only, bounded JSON diagnostic. It reports runner and built-source
identity, Bazelisk pins, platform, rc/profile composition, task scratch
classification, and stale host-install state without dumping the environment.

## Installation

Bootstrap the host installation with the underlying Bazel command:

```sh
bazel --batch run --config=agent //projects/bazel_agent:install
```

Once installed, update it with:

```sh
bazel_agent run //projects/bazel_agent:install
```

The install target atomically replaces `~/.local/bin/bazel_agent`. After every
code update to this project, rerun the install target; the installed binary is
not updated automatically.
