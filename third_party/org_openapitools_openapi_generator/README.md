---
title: OpenAPI Generator dependency
description: Pinned upstream OpenAPI Generator CLI JAR
---

The Apache-2.0 [OpenAPI Generator](https://github.com/OpenAPITools/openapi-generator)
CLI is pinned to 7.26.0 from Maven Central. The self-contained JAR uses the
repository Java toolchain. Its SHA-256 SRI pin is in `include.MODULE.bazel`;
the downloaded bytes were also checked against the publisher's `.jar.sha1`.
The first-party wrapper is [`tools/openapi_generator`](../../tools/openapi_generator/README.md).
