## Context

The SDK consumes a pinned mirrored schema. Fetching upstream is a dependency-maintenance workflow.

## Decisions

Keep the canonical maintainer target beside the dependency pin, under third_party/com_x_api_openapi. Keep handwritten executable source under tools/x_api_schema, with the existing shared OAuth signer. Move the existing E2E tests unchanged. The SDK owns only runtime code and generation.

## Verification

Run the moved command E2E test, compile the new authenticated maintenance target, verify SDK and uploader E2E, and complete the existing PR delivery gates. No live upstream fetch or host deployment is needed for this relocation.
