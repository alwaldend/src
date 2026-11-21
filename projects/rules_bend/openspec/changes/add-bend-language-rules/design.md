## Context

Bend 2 runs an entry point directly and supports a separate `--check-only`
mode. Its loader resolves filesystem symlinks before looking up relative
imports.

## Goals

- Run a Bend program with `bazel run` and forward its arguments and exit status.
- Check a source and its imports with a sandboxed build action.
- Stage declared sources as regular files at their original relative paths.

## Decisions

Register a Bend runtime toolchain selected by the target platform for
`bend_binary` and a checker toolchain selected by the execution platform for
`bend_check`. The checker and its runner share an execution group. Both
modes preserve compiler runfiles and use the same runner. The executable resolves compiler and source inputs
from Bazel runfiles, while the check action receives them as action inputs.

The default Linux x86-64 toolchains share the upstream release executable.
The archive pin records the published checksum. Its Bazel executable target
retains the original `bin/bend`, `bend2/`, and `guide/` layout and declares
the standard library and guide as runtime files.

The project owns its OpenSpec validation test under `openspec/`. The test
uses the shared validation macro and the workspace sources in that package.
