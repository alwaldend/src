## Why

The converted artifact mixes X request fields with local upload metadata and diagnostic records. A dedicated payload makes the request boundary explicit, while strict command defaults prevent warning-bearing output from being mistaken for lossless conversion.

## What Changes

- Place title and content_state inside payload; resolved cover_media also belongs there.
- Keep image_locators and banner_locator outside payload and omit serialized diagnostics.
- Print diagnostics to stderr and fail on warnings by default. --warnings-as-errors=false permits warnings, but errors always fail.
- Reject old flat artifacts with a regeneration instruction.
- Update command and conversion-to-draft E2E coverage and the maintained contract.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- markdown-to-draftjs: nested artifact payload and stderr-only diagnostic policy.

## Impact

Converter CLI, artifact encoder/reader, publisher, project documentation and existing tests. Existing artifacts must be regenerated. No new dependencies, site changes, credentials or live X requests.
