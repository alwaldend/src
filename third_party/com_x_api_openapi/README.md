---
title: X API OpenAPI schema
description: Checksum-pinned mirror of the upstream X API schema
---

The [X API schema](https://api.x.com/2/openapi.json) snapshot reports API version
2.173. Its JSON snapshot, with surrounding whitespace normalized to a final
newline, is mirrored at
`https://download.alwaldend.com/third_party/com_x_api_openapi/releases/2.173-4377d1791c01/openapi.json`.
The complete SHA-256 is
`4377d1791c01632c9cf2860fc833a4cec76eff74f89d232fd7355e765e7bf0a9`.
The upstream API terms and schema descriptions remain authoritative.

`include.MODULE.bazel` owns the immutable URL and SRI pin. Bazel downloads the
snapshot anonymously and caches it as an external file; the schema body is not
stored in Git. Both local and Yandex download stores contain the same bytes.

To refresh, use `//third_party/com_x_api_openapi:fetch_schema` to download into ignored `out/`
with the uploader's existing credentials. Inspect the schema, assign a new
upstream-version-and-digest release label, generate a portable release bundle
with `//tools/release`, and publish through `//infra/download:publish.local` and
`:publish.yandex` with `--namespace third_party`. Verify both HTTPS responses and
their hashes before updating the pin. Never replace a pinned version's bytes.
