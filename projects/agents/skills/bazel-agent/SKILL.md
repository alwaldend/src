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
bazel_agent test //path/to/package:all
```

The runner performs no command validation. It passes an empty argument list,
`--help`, malformed commands, and every other argument through to Bazel. When a
first argument exists, it places that argument in Bazel's command position,
then passes targets, command options, and arguments after a `bazel run`
separator normally. Later options retain Bazel's normal precedence, so use
them when the task intentionally needs to override an agent configuration
setting:

```sh
bazel_agent run //path/to:tool -- --tool-option
```

Do not call `bazel` directly or repeat `--batch` or `--config=agent`.
`bazel_agent` adds those flags in their required positions, resolves the
Bazelisk-managed `bazel` from `PATH`, and replaces itself with that process.
The replacement preserves direct signal delivery and Bazel's exit status. The
repository `.bazeliskrc` pins the Bazel version and archive hash.

Use `repo-bazel` in addition to this skill when changing BUILD files, Starlark,
Bzlmod dependencies, toolchains, or the build graph.

## Install or update the runner

If `bazel_agent` is not installed, bootstrap it with the underlying command:

```sh
bazel --batch run --config=agent //projects/bazel_agent:install
```

After every code update under `projects/bazel_agent`, reinstall the host copy:

```sh
bazel_agent run //projects/bazel_agent:install
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
