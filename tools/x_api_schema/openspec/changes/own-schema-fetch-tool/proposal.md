## Why

The schema maintenance entry point still lives in the runtime SDK after schema ownership moved to third_party. The user requested correcting that ownership.

## What Changes

- Move the authenticated maintenance target to third_party/com_x_api_openapi.
- Move its unchanged Go implementation and E2E tests to tools/x_api_schema/cmd/fetch_schema.
- Restrict credential access and tool visibility to the dependency-maintenance package, and update current documentation.

## Capabilities

No specification delta; this relocates existing maintenance behavior without changing authentication, validation or publication.

## Impact

The schema dependency, first-party fetch tool, SDK documentation and maintenance dependency boundaries. Existing mirrored bytes and runtime generation remain unchanged.
