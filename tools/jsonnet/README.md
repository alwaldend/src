---
title: Jsonnet
description: Sandboxed generation of declared JSON outputs
---

This package exposes the approved Jsonnet compiler and formatter from the
[pinned upstream toolchain](../../third_party/com_github_google_go_jsonnet/README.md).
It owns build execution, not application configuration or generated artifacts.

Use jsonnet_to_json from defs.bzl with an entrypoint, all imported files in deps,
and explicit package-relative outs. The entrypoint returns an object keyed by
those output paths. Bazel executes the compiler without network access and puts
native JSON in the consumer's output tree. The compiler is an execution tool;
it is not installed on hosts or included in application/runtime dependencies.

```sh
bazel run //tools/jsonnet -- --version
bazel run //tools/jsonnet:format -- --version
```

The compiler's multi-output mode and directory creation preserve model-grouped
filenames. Every import must be a declared input, including configuration JSON.
No source or generated output is downloaded at execution time.
