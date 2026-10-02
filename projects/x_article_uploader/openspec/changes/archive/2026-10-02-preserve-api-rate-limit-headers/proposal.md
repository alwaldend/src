## Why

Failed X API requests currently discard the response headers that distinguish
an exhausted quota from unavailable service. Preserving the limits and reset
times lets users see the known waiting boundary without spending another request.

## What Changes

- Preserve and report only the six X rate and daily usage limit headers plus
  `Retry-After` on rejected requests, including HTTP 503.
- Display valid reset timestamps in UTC alongside their original values.
- Report a conservative not-before boundary from known exhausted windows and
  applicable server advice, without promising service recovery or success.
- Keep malformed, missing, expired, or ambiguous timing evidence from producing
  a fabricated availability time; preserve the raw diagnostic values.
- Verify behavior through local HTTP E2E fixtures without real X requests.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `x-article-publication`: Preserve rate-limit metadata and report retry bounds
  as part of API failure diagnostics.

## Impact

The shared X API error model, uploader README, and E2E fixtures change. The CLI
already prints wrapped errors and needs no new retry or network-probing flow.
This change adds no external dependency and authorizes no live X request.
