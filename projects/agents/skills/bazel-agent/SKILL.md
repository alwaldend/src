---
name: bazel-agent
description: >-
  Invoke Bazel safely and consistently as an agent in this monorepo through
  the bazel_agent runner. Use for every agent-executed repository Bazel
  command, for runner installation or troubleshooting, and alongside
  repo-bazel when changing Bazel files or dependencies.
---

# Use the Bazel agent runner

## Run repository commands

Run Bazel from the applicable workspace root with this form:

```sh
bazel_agent bazel test //path/to/package:all
```

`bazel` is the only Bazel entry point, and it is a validated subcommand.
The runner accepts a known Bazel command after the `bazel` keyword, rejects
arbitrary leading arguments, and places `--config=agent` in its required
position. Targets, command options, and arguments after a `bazel run`
separator pass through normally. Later options retain Bazel's normal
precedence, so use them when the task intentionally needs to override an
agent configuration setting:

```sh
bazel_agent bazel run //path/to:tool -- --tool-option
```

Do not call `bazel` directly or repeat `--config=agent`. Nested workspaces
import the same shared rc files, so the `agent` configuration resolves there
as well. `bazel_agent` adds the flag in its required position, resolves the
Bazelisk-managed `bazel` from `PATH`, and replaces itself with that process.
The replacement preserves direct signal delivery and Bazel's exit status. The
repository `.bazeliskrc` pins the Bazel version and archive hash.

For Bazel commands that support multiple targets, such as `build` and `test`,
batch compatible targets into one invocation when they use the same options:

```sh
bazel_agent bazel build //path/to:first //path/to:second
bazel_agent bazel test //path/to:first_test //path/to:second_test
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

## Install or update the runner

If `bazel_agent` is not installed, bootstrap it with the underlying command:

```sh
bazel run --config=agent //projects/bazel_agent:install
```

After every code update under `projects/bazel_agent`, reinstall the host copy:

```sh
bazel_agent bazel run //projects/bazel_agent:install
```

The install target atomically updates `~/.local/bin/bazel_agent`; confirm that
directory is on `PATH`. Installing changes host state, so do it only when the
task authorizes installation or updating the agent environment.

## Diagnose failures without bypassing the runner

- If `bazel_agent` is missing, install it rather than duplicating its flags in
  later commands.
- If it cannot find `bazel`, repair the host's `bazel` to Bazelisk provisioning
  or `PATH`; do not resolve or hard-code a separate Bazel binary in the runner.
- If Bazel cannot write its output base or Bazelisk cache, fix the managed
  sandbox writable roots. Do not create an unmanaged duplicate cache in
  `/tmp`; keep task scratch data under the repository's ignored `out/` tree.
- Treat network, credential, and remote-cache failures separately from build
  or test failures caused by the patch.
