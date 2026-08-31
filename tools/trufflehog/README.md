---
title: Trufflehog
description: Trufflehog
languages:
  - bzl
tags:
  - bzl_rules
---

The normal `//tools/trufflehog:repo_test` scans the history reachable from the
checked-out `HEAD` and skips unrelated local refs. Run the manual
`//tools/trufflehog:repo_all_refs_test` target explicitly when an audit must
include every local Git ref, such as T3 checkpoint refs from other worktrees.

## Links

- Repo: https://github.com/trufflesecurity/trufflehog
