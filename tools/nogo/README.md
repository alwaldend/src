---
title: Nogo
description: Golang static analyser
---

The analyzer configuration preserves generated upstream code. `nilness` ignores
only OpenAPI Generator 7.26.0's redundant named-result nil check in the X SDK's
generated `client.go`; handwritten SDK source and other analyzers remain enabled.

## Links

- Docs: https://github.com/bazel-contrib/rules_go/blob/master/go/nogo.rst
- Bzlmod docs: https://github.com/bazel-contrib/rules_go/blob/master/docs/go/core/bzlmod.md#configuring-nogo
