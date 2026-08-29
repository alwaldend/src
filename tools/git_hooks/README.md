---
title: Hooks
description: Git hooks
tags:
  - bzl_rules
---

The hook installer resolves the repository's effective hooks directory through
Git. This supports linked worktrees and repositories that configure
`core.hooksPath`.

The installed pre-commit hook requires `bazel_agent`. Bootstrap it before
installing the hook:

```sh
bazel --batch run --config=agent //projects/bazel_agent:install
```

Install or update the checked-in hooks:

```sh
bazel_agent run //:write_git_hooks
```

Verify that every hook is current and executable without changing it:

```sh
bazel_agent run //:write_git_hooks -- test
```
