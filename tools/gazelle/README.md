---
title: Gazelle
description: Repository BUILD file generation
---

`//tools/gazelle:gazelle` runs the repository's Gazelle language plugins.
Repository-wide directives remain in root `BUILD.bazel`, where Gazelle discovers
them for all packages. Root `//:gazelle` and `//:gazelle_bin` remain compatibility
entry points.

The transitive-visibility extension is registered last so it also covers
packages created by other languages. Each tree selects a package group with
[`tools/transitive_visibility/transitive_visibility.toml`](../transitive_visibility/transitive_visibility.toml). To regenerate
only these package declarations, run:

```sh
bazel_agent bazel run //:gazelle -- -lang=enforce_transitive_visibility
```

The [plugin documentation](../../projects/bazel_gazelle_enforce_transitive_visibility/README.md)
describes inheritance and standalone use. Nested modules are separate Gazelle
workspaces; root path mappings do not cross those boundaries.

`//tools/gazelle:gazelle_python_manifest.update` refreshes root
`gazelle_python.yaml`; `//tools/gazelle:gazelle_python_manifest.test` validates
its integrity against `tools/py/uv.lock`. The manifest stays at the
repository root for default ancestor discovery. It can live elsewhere if the
root BUILD file sets `gazelle:python_manifest_file_name` to its relative path.
