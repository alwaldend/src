---
title: Root Bazel module
description: Generate the root module's owning include list
---

`MODULE.bazel` is generated directly from `include.MODULE.bazel` files:

```sh
bazel run //tools/bazel_module:update
```

Dependency declarations and pins stay in their owning include files. The small
root `module(...)` header lives in the Go generator. Discovery uses the existing
`//tools/git` host Git wrapper to include tracked and untracked files while
respecting Git ignore rules. Deleted files, symlinks, `.bazelignore`
directories, and standalone workspaces beneath
`MODULE.bazel`, `WORKSPACE`, or `WORKSPACE.bazel` boundaries are excluded.

Includes are ordered by owning directory, with `third_party`, `tools`, then
`projects` before other trees. This keeps the existing precedence between those
trees, including toolchain registrations. A parent directory precedes its
children. Keep any order-sensitive declarations together in one owning include.

`//tools/bazel_module/cmd/update:freshness_test` checks the current checkout and
is included in `//:repo_quality_test`. It intentionally disables sandboxing and
result caching to detect new or removed files outside declared Bazel inputs.
The generator's unit tests remain sandboxed and cacheable.

Before deleting or moving an existing include, retain a launcher so regeneration
does not need to load the old module after its include disappears:

```sh
mkdir -p out/bazel_module
bazel run --script_path=out/bazel_module/update //tools/bazel_module:update
# Delete or move the include, then regenerate:
out/bazel_module/update
```
