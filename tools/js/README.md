---
title: Js
description: Js
---

Repository JavaScript tooling dependencies belong to the pnpm workspace in
`tools/`. Its `package.json` and `pnpm-lock.yaml` own their requested and resolved
versions. `//tools:node_modules` owns the Bazel package store and links; loading
the repository root package does not load the npm extension's generated rules.

Run pnpm with `//tools/pnpm`, passing `--dir` with the absolute path to `tools/`
for dependency updates. Keep lifecycle scripts disabled (`--ignore-scripts`)
when updating the lockfile.
