---
title: Hooks
description: Git hooks
tags:
  - bzl_rules
---

The hook installer resolves the repository's effective hooks directory through
Git. This supports linked worktrees and repositories that configure
`core.hooksPath`.

The installed pre-commit hook requires the Bazelisk-managed `bazel` on `PATH`.

Install or update the checked-in hooks:

```sh
bazel run //:write_git_hooks
```

Verify that every hook is current and executable without changing it:

```sh
bazel run //:write_git_hooks -- test
```
