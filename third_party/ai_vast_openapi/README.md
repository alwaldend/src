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

`root_certificate` provides the official Vast instance root from
[Vast's certificate instructions](https://docs.vast.ai/guides/instances/connect/jupyter),
verified at SHA-256 `4068eae92482f045a1159b52b604142cc2c7e427a76bf8cf2070e779978e9507`.
The image-generation CLI uses this scoped trust input with hostname checking;
it is not installed into host or browser trust stores. Rotation requires a
reviewed checksum update through the same declared download.
