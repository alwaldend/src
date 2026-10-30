## Context

See proposal.md. The repository already has Java, Go, Vault injection and tree artifact copy rules.

## Goals / Non-Goals

Keep external pins in third_party and generated source in Bazel outputs. No infrastructure or secret changes are required.

## Decisions

Reuse the uploader OAuth 1.0a Vault injection through its existing al configuration. Move its signer to the SDK public oauth1 package and preserve the uploader credential alias. Fetching is an explicit network operation with no retries or redirects, bounded response size, JSON/OpenAPI validation and atomic snapshot replacement. Offline generation consumes the declared schema input and pinned generator; generated code is never written into source. rules_go expands the generated Go-only tree as compiler inputs. Keep golang.org/x/oauth2 and the explicitly approved gopkg.in/validator.v2 v2.0.1 as direct dependencies used by generated code. The user selected internal/x_sdk for repository-only use; no project landing page or release target is created.

The pinned Go generator emits invalid pointer-based identifiers for three binary/base64 media unions. Suppress generation of only those union wrappers through schema mappings and supply handwritten aliases to encoding/json.RawMessage, preserving JSON upload values without custom templates or changing the fetched schema. A generated-client HTTP fixture checks base64 media serialization. Nogo excludes only the nilness analyzer finding for the generated client.go redundant named-result nil check. Keep other binary response types unchanged.

## Risks / Trade-offs

Upstream schema and generated names may change; compile and behavior checks validate each refreshed candidate. Authentication availability may prevent live fetching; retain that failure explicitly.

The user requested removing the schema body from Git. The snapshot is now
versioned under /third_party/com_x_api_openapi on both download hosts and is
consumed through an integrity-pinned http_file dependency. Upstream refresh
authentication remains an explicit maintainer operation that writes ignored out/.

The maintenance target is owned by third_party/com_x_api_openapi:fetch_schema;
its handwritten command implementation and E2E tests belong to tools/x_api_schema.
The SDK consumes the pinned dependency and owns no infrastructure or maintenance-tool
dependency. Credential visibility is restricted to the dependency-maintenance package.

## Review corrections

XSDK-PROBLEM-DECODING-001: the default Go oneOf decoder accepts multiple Problem
variants because their type fields are plain strings. A valid partial-success
response therefore fails decoding. Use the existing schema discriminator mapping
and verify both the returned user and its resource-not-found problem through an
HTTP fixture before changing generator configuration.
