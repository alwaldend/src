---
title: Gazelle
description: Repository BUILD file generation
---

`//tools/gazelle:gazelle` runs the repository's Gazelle language plugins.
Repository-wide directives remain in root `BUILD.bazel`, where Gazelle discovers
them for all packages. Root `//:gazelle` and `//:gazelle_bin` remain compatibility
entry points.

`//tools/gazelle:gazelle_python_manifest.update` refreshes root
`gazelle_python.yaml`; `//tools/gazelle:gazelle_python_manifest.test` validates
its integrity against `tools/py/requirements.txt`. The manifest stays at the
repository root for default ancestor discovery. It can live elsewhere if the
root BUILD file sets `gazelle:python_manifest_file_name` to its relative path.
