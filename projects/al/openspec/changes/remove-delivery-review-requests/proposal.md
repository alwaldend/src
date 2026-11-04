## Why

The delivery review reference instructs callers to request another review after
a review-driven rewrite. Delivery should consume Codex review feedback without
initiating reviews or managing Codex review settings.

## What Changes

- Remove the automatic review-request instruction after a rewrite.
- Prohibit requesting Codex reviews or mentioning Codex to initiate a review
  during delivery; preserve inspection and guarded feedback handling.
- Keep Codex review configuration outside the delivery skill.
- Add evaluation scenarios for initial publication and review-driven updates.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `review-delivery`: Delivery consumes Codex reviews without initiating them.

## Impact

AL delivery skill, review reference, evaluation cases, and generated skill
discovery. No CLI, dependency, or live configuration changes.
