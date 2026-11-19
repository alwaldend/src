---
title: Source repository
description: Repository evolution specifications and workflow history
---

This project owns the [OpenSpec workspace](openspec/README.md) describing
evolution of the alwaldend/src monorepo: its shared structure, build system
and development workflows. It also retains the index of the goal-to-OpenSpec
migration. Individual projects own their own specifications and change history.

The repository root remains the Bazel workspace and owns shared executable
configuration. This directory owns repository development specifications;
adding or changing them does not provision infrastructure.

Run the [pinned CLI](../../tools/openspec/README.md) from the repository root.
It selects this project by default:

```sh
bazel_agent bazel run //tools/openspec -- list --specs
```
