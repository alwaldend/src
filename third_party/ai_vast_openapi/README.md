---
title: Vast.ai OpenAPI schema
description: Immutable official API schema for the shared Vast.ai SDK
---

The official [Vast.ai OpenAPI 3.1 schema](https://github.com/vast-ai/docs/blob/2b2d9a04e7cb542f3305c0c148c2395b0cfab86e/api-reference/openapi.yaml)
is pinned to commit `2b2d9a04e7cb542f3305c0c148c2395b0cfab86e` with SHA-256
integrity in `include.MODULE.bazel`. Its declared license is the
[Vast.ai Terms of Service](https://vast.ai/terms/).

Bazel downloads the schema anonymously from the official repository. The
shared [client library](../../internal/vastai_sdk/README.md)
owns generator configuration and adaptation; generated Go source stays in Bazel
outputs. Refresh the commit and checksum together, then rebuild and run the
provider's Terraform CLI acceptance suite. Do not edit the upstream schema.
