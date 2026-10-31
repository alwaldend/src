---
title: X SDK
description: Go client generated from a checksum-pinned X API schema mirror
---

The SDK generates a Go client from the checksum-pinned X API schema mirror in
[`third_party/com_x_api_openapi`](../../third_party/com_x_api_openapi/README.md).
Depend on `//internal/x_sdk/pkg/client` and import
`git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/client` in Bazel-built Go code.
The reusable OAuth 1.0a signer lives in `pkg/oauth1`; the Article uploader also
uses it.

## Build

```sh
al bazel build //internal/x_sdk:generate //internal/x_sdk/pkg/client
al bazel test //internal/x_sdk/...
```

Bazel fetches the pinned schema anonymously from the versioned download mirror
and caches it. Builds need no X credentials. Generation uses the pinned
[OpenAPI Generator](../../tools/openapi_generator/README.md) and `generator.json`
offline in a Bazel sandbox. Generated Go source stays under
`bazel-bin/internal/x_sdk/generate/`; neither schema bytes nor generated client
source are stored in Git. The upstream templates use the pinned
`gopkg.in/validator.v2` v2.0.1 dependency for union response decoding.

The [schema dependency](../../third_party/com_x_api_openapi/README.md) owns
upstream refresh, mirror publication, provenance and checksum updates. Its
maintenance command is implemented by [tools/x_api_schema](../../tools/x_api_schema/README.md).

## Authentication and compatibility

The schema's supported authentication methods determine the generated client
configuration. Its default server is X; override `Configuration.Servers` and
`Configuration.HTTPClient` for local verification. OAuth 1.0a callers can use
`pkg/oauth1` to sign requests in an HTTP transport. OpenAPI Generator does not
automatically sign the schema's nonstandard HTTP `OAuth` security scheme.
Response unions use the schema's discriminator mappings, including `Problem`
values in partial-success responses. The snapshot also contains discriminator
paths that the generator warns about; streaming and unmapped polymorphic
response decoding need endpoint-specific validation.

Three binary/base64 upload fields use handwritten `json.RawMessage` aliases
selected through `generator.json` because the pinned generator emits invalid Go
identifiers for those unions. Supply a base64 string encoded as a JSON string
literal for JSON uploads. Other binary download and upload types retain the
upstream generator's defaults. Nogo excludes only the generator's redundant
named-result nil check in its generated `client.go` from the `nilness` analyzer;
handwritten SDK code remains analyzed normally.
