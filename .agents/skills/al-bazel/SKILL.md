---
name: al-bazel
description: >-
  Invoke Bazel with the repository's pinned version and per-user configuration.
  Use for every agent-executed repository Bazel command, affected-test
  selection, setup or troubleshooting, and alongside
  repo-bazel when changing Bazel files or dependencies.
---

# Run Bazel

## Run repository commands

Run Bazel from the applicable workspace root with this form:

```sh
bazel test //path/to/package:all
```

Invoke the Bazelisk-managed `bazel` on `PATH` directly. The repository
`.bazeliskrc` pins its version and archive hash. Pass startup options before
the command, command options after it, and executed-tool arguments after `--`:

```sh
bazel run //path/to:tool -- --tool-option
```

`bazel run` may launch the tool in its runfiles tree. Pass workspace-absolute
paths for caller-owned input and output files unless the tool explicitly
normalizes relative paths against `BUILD_WORKSPACE_DIRECTORY`. Keep generated
outputs under the task's ignored workspace `out/<task>/` directory.

Agent output defaults belong in the account's `~/.bazelrc`, as described in
the [repository setup guide](../../../alwaldend.com/content/docs/misc/repo.md).
The bot account's configuration is managed by
[host Ansible](../../../../users/simeonwarren/host_bot/ansible/files/bazelrc).
Normal builds do not require AL or a named agent profile. Change persistent
host configuration through its owning declarative workflow; a source task
does not authorize deploying it. Command-line options can override user
defaults when a task requires different output or test behavior.

For Bazel commands that support multiple targets, such as `build` and `test`,
batch compatible targets into one invocation when they use the same options:

```sh
bazel build //path/to:first //path/to:second
bazel test //path/to:first_test //path/to:second_test
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

## Select tests affected by Git changes

Use `al bazel test` when the task calls for selecting tests affected by a Git
change range. Keep ordinary `bazel build`, `bazel test`, and semantic lint for
explicit target requests and other required checks. This command selects tests;
it is not a generic Bazel wrapper or a delivery step.

With an installed or task-local AL executable:

```sh
al bazel test --base origin/master --dry-run //projects/widget/...
al bazel test --base origin/master //projects/widget/...
```

Choose the base Git ref explicitly. Target patterns bound the requested test
scope and default to `//...`. Use `--workspace PATH` to select another Bazel
workspace, including a nested module. Selection uses Git changed paths and
Bazel dependency queries in that workspace; it does not implicitly cover every
workspace in the monorepo. `--dry-run` queries and reports the selection without
running the tests. Inspect the report's selected targets and any fallback
reason before describing its coverage.

Missing or deleted paths and build configuration changes can require selecting
the full requested test scope. That conservative fallback is not proof that
all repository checks ran. A failed Git or Bazel query is a failure, never an
empty successful selection. Preserve its diagnostics and diagnose it before
claiming that no tests were affected.

The command invokes Bazel for graph queries and test execution. No AL profile
or generic flag passthrough is involved; Bazel retains its normal configuration
and cache.

## Run repository tools

Run repository control tools from the enclosing Git worktree's root Bazel
workspace, which imports `rules_al`, even for changes in a nested
module. Nested build and test targets still use their owning workspace.

```sh
bazel run @rules_al//:al -- delivery provider
```

For frequent execution, install the standalone tool explicitly and invoke it
directly:

```sh
bazel run //tools/al:install
al delivery provider
```

The installed AL command is a copy of the local source-built executable. Rerun
the install target to update it; AL does not maintain an executable cache or
provide `al tool`. See the [tool documentation](../../docs/delivery/README.md).
Installing changes host state and requires installation authority.

## Install or update AL

Install the standalone lifecycle tool when the task authorizes it:

```sh
bazel run //tools/al:install
```

For task-only execution, pass `-- --directory "$PWD/out/task/bin"` to the
same target and invoke `out/task/bin/al` directly. This keeps the executable
inside ignored task scratch and does not change the shared host installation.

The install target atomically updates `~/.local/bin/al`; confirm that
directory is on `PATH`. Installing changes host state, so do it only when the
task authorizes installation or updating the agent environment.

## Diagnose failures

Allow a cold tool build enough time for analysis as well as execution; use
observed build duration to choose a bounded deadline, not a short generic
command timeout. Poll at intervals that permit useful output and user updates.
Separate commands using the same Bazel server share its lock: complete ready
CLI work before starting a long build, and wait for that build before issuing
dependent commands. Do not queue another invocation merely to check progress.

- If `bazel` is missing, inspect the host's Bazelisk provisioning and `PATH`.
  Installing AL does not install Bazel.
- If Bazel cannot write its output base or Bazelisk cache, fix the managed
  sandbox writable roots. Do not create an unmanaged duplicate cache in
  `/tmp`; keep task scratch data under the repository's ignored `out/` tree.
- Treat network, credential, and remote-cache failures separately from build
  or test failures caused by the patch.
