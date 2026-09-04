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
Bazel invocation behind a validated subcommand while consistently applying the
agent configuration:

```text
bazel_agent bazel test //path/to/package:all
    ->
bazel test --config=agent //path/to/package:all
```

`bazel` is the only Bazel entry point, and it is a validated subcommand. The
runner accepts a known Bazel command after the `bazel` keyword, rejects
arbitrary leading arguments, and places `--config=agent` after the command.
Targets, command options, and arguments after the `--` separator of a
`bazel run` invocation pass through unchanged. Later options therefore retain
Bazel's normal precedence and can override settings supplied by the agent
configuration. The runner resolves `bazel` from `PATH` and replaces itself
with that process, so signals and the final exit status are not mediated by
another wrapper process. Uses of the persistent Bazel server follow the host
Codex network policy; the `host-bot` profile allows loopback so the
client-server connection is not blocked.

The runner does not create or inject a host temporary directory. Bazel actions
use declared outputs and Bazel-managed temporary storage, and tests use their
test temporary-directory contract. Repository-updating host tools that need
scratch must accept or derive an explicit task/run path under `out/<task>/`
instead of relying on ambient `TMPDIR`, `TMP`, or `TEMP` values propagated to
the whole build.

## Cached repository tools

Frequently used repository control tools can bypass Bazel analysis after one
exact build:

```sh
bazel_agent tool run repo_delivery -- provider
bazel_agent tool warm mcp_cordis repo_delivery
bazel_agent tool path repo_delivery
```

The runner hashes its own binary plus each tool's declared source, dependency
locks, Bazel pins, repository rc files, optional user and host rc files, their
transitive imports, any `BAZELRC` environment files, and the platform. A cache
miss takes an exclusive per-key lock, builds the target with runfiles enabled,
and copies the complete runnable output into an atomically installed entry. A
hit executes that entry directly without starting a Bazel client or loading a
worktree graph. Concurrent misses for the same exact key share one build.
Inputs are hashed again under the lock and after the build; an in-flight source
change selects a new key instead of publishing stale output under the old one.

`BAZEL_AGENT_TOOL_CACHE` or `--cache-root` can select the cache. The default is
`/var/cache/bazel/tool_cache` when that managed directory exists, otherwise
the user's cache directory. Because entries contain executable code, the cache
root must not be writable by group or others. Source changes select a new key,
so a dirty relevant source file cannot silently reuse an older executable.

`bazel_agent doctor --workspace-root PATH --task-scratch out/<task>/<run>` is
a read-only, bounded JSON diagnostic. It reports runner and built-source
identity, Bazelisk pins, platform, rc/profile composition, task scratch
classification, and stale host-install state without dumping the environment.

## Installation

Bootstrap the host installation with the underlying Bazel command:

```sh
bazel run --config=agent //projects/bazel_agent:install
```

Once installed, update it with:

```sh
bazel_agent bazel run //projects/bazel_agent:install
```

The install target atomically replaces `~/.local/bin/bazel_agent`. After every
code update to this project, rerun the install target; the installed binary is
not updated automatically. Repository launchers probe the installed runner's
tool-cache support before using it, so an older host installation retains its
documented Bazel fallback during rollout.
