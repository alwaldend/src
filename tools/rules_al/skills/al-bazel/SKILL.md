---
name: al-bazel
description: >-
  Invoke Bazel safely and consistently as an agent in this monorepo through
  the AL runner. Use for every agent-executed repository Bazel
  command, for runner installation or troubleshooting, and alongside
  repo-bazel when changing Bazel files or dependencies.
---

# Use the Bazel agent runner

## Run repository commands

Run Bazel from the applicable workspace root with this form:

```sh
al bazel test //path/to/package:all
```

`bazel` is the only Bazel entry point, and it is a validated subcommand.
The runner accepts a known Bazel command after the `bazel` keyword, rejects
arbitrary leading arguments, and places `--config=agent` in its required
position. Targets, command options, and arguments after a `bazel run`
separator pass through normally. Later options retain Bazel's normal
precedence, so use them when the task intentionally needs to override an
agent configuration setting:

```sh
al bazel run //path/to:tool -- --tool-option
```

`bazel run` may launch the tool in its runfiles tree. Pass workspace-absolute
paths for caller-owned input and output files unless the tool explicitly
normalizes relative paths against `BUILD_WORKSPACE_DIRECTORY`. Keep generated
outputs under the task's ignored workspace `out/<task>/` directory.

Do not call `bazel` directly or repeat `--config=agent`. Nested workspaces
import the same shared rc files, so the `agent` configuration resolves there
as well. `al` adds the flag in its required position, resolves the
Bazelisk-managed `bazel` from `PATH`, and replaces itself with that process.
The replacement preserves direct signal delivery and Bazel's exit status. The
repository `.bazeliskrc` pins the Bazel version and archive hash.

For Bazel commands that support multiple targets, such as `build` and `test`,
batch compatible targets into one invocation when they use the same options:

```sh
al bazel build //path/to:first //path/to:second
al bazel test //path/to:first_test //path/to:second_test
```

The persistent Bazel server keeps startup overhead low across invocations, so
batching is about compatibility rather than cost. Do not batch a single-target
command such as `run`. Otherwise keep invocations separate only when they
require different commands or options, have a real ordering dependency, need
failure isolation for diagnosis, or would create unsafe resource contention.
Do not run separate compatible invocations merely to parallelize work that
Bazel already schedules internally.

Use `repo-bazel` in addition to this skill when changing BUILD files, Starlark,
Bzlmod dependencies, toolchains, or the build graph.

## Run repository tools

Run repository control tools from the enclosing Git worktree's root Bazel
workspace, which imports `rules_al`, even for changes in a nested
module. Nested build and test targets still use their owning workspace.

```sh
al bazel run @rules_al//:al -- delivery provider
```

For frequent execution, install the standalone tool explicitly and invoke it
directly:

```sh
al bazel run //tools/al:install
al delivery provider
```

The installed AL command is a copy of the local source-built executable. Rerun
the install target to update it; AL does not maintain an executable cache or
provide `al tool`. See the [tool documentation](../../docs/delivery/README.md).
Installing changes host state and requires installation authority.

## Install or update the runner

If `al` is not installed, bootstrap it with the underlying command:

```sh
bazel run --config=agent //tools/al:install
```

For task-only execution, pass `-- --directory "$PWD/out/task/bin"` to the
same target and invoke `out/task/bin/al` directly. This keeps the executable
inside ignored task scratch and does not change the shared host installation.

When the task authorizes updating the host runner after a code change:

```sh
al bazel run //tools/al:install
```

The install target atomically updates `~/.local/bin/al`; confirm that
directory is on `PATH`. Installing changes host state, so do it only when the
task authorizes installation or updating the agent environment.

## Diagnose failures without bypassing the runner

Allow a cold tool build enough time for analysis as well as execution; use
observed build duration to choose a bounded deadline, not a short generic
command timeout. Poll at intervals that permit useful output and user updates.
Separate commands using the same Bazel server share its lock: complete ready
CLI work before starting a long build, and wait for that build before issuing
dependent commands. Do not queue another invocation merely to check progress.

- If `al` is missing, use one of the bootstrap flows above rather than
  duplicating its flags in later commands.
- If it cannot find `bazel`, repair the host's `bazel` to Bazelisk provisioning
  or `PATH`; do not resolve or hard-code a separate Bazel binary in the runner.
- If Bazel cannot write its output base or Bazelisk cache, fix the managed
  sandbox writable roots. Do not create an unmanaged duplicate cache in
  `/tmp`; keep task scratch data under the repository's ignored `out/` tree.
- Treat network, credential, and remote-cache failures separately from build
  or test failures caused by the patch.
