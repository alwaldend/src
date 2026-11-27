---
title: go_mod
description: Keep go.mod files on one Go version
tags:
  - go
  - modules
---

# go_mod

`go_mod` keeps every tracked `go.mod` file on the Go version configured in
`tools/go_mod/cmd/go_mod/BUILD.bazel`. It discovers files with Git, so ignored
scratch files are excluded. The updater stages each changed module beside its
destination, preserves its permissions, and replaces it after a successful
write. Rerunning resumes with the remaining mismatches. The parser accepts
up to 16 MiB per module and Git's module path list may contain up to 4 MiB.

Update all module files:

```sh
bazel run //tools/go_mod/cmd/go_mod:update
```

Check without writing:

```sh
bazel test //tools/go_mod/cmd/go_mod:test
```

The check target is part of `//:repo_quality_test`.
