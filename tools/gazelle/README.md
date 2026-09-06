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
its integrity against root `requirements.txt`. Keeping the manifest at the
repository root preserves Python dependency discovery in every package.
