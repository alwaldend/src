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
scratch files are excluded.

Update all module files:

```sh
bazel run //tools/go_mod/cmd/go_mod:update
```

Check without writing:

```sh
bazel test //tools/go_mod/cmd/go_mod:test
```

The check target is part of `//:repo_quality_test`.
