---
title: Bzl
description: Bazel code
---

## Archive binaries

Each archive `binaries` entry requires a `name` and `path`. It can also set
`runtime_files` to a list of Bazel glob patterns relative to the unpacked
archive. The generated `<name>_native_binary` exposes matching files as
runfiles; `<name>_filegroup` exposes the binary and those same files.

Omit `runtime_files` when the binary needs no extra runtime files.
