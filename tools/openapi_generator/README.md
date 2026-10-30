---
title: OpenAPI Generator
description: Pinned OpenAPI Generator CLI and offline Bazel source generation
---

The wrapper runs the unmodified upstream CLI with Bazel's pinned Java runtime.
The JAR and integrity pin belong to
[`third_party/org_openapitools_openapi_generator`](../../third_party/org_openapitools_openapi_generator/README.md).

```sh
bazel_agent bazel run //tools/openapi_generator -- version
bazel_agent bazel run //tools/openapi_generator -- help generate
```

`openapi_generate` in `defs.bzl` declares the schema and configuration as inputs
and returns one directory artifact. Generation runs offline in the normal Bazel
sandbox. Its output stays in Bazel's output tree. Consumers can filter the tree
to the source language they need and compile it using their normal language
rules, as the [X SDK](../../internal/x_sdk/README.md) does for Go. Remote schema
references must be vendored before sandboxed generation.

Generation diagnostics are retained in a declared `.log` artifact rather than
printing a message for every generated file. Request `--output_groups=logs`
on the generation target to materialize the complete upstream log, for example:

```sh
bazel_agent bazel build --output_groups=logs //internal/x_sdk:generated_client
```
