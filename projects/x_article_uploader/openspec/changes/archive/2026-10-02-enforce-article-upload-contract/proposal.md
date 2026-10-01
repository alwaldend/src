## Why

Review against X's current documentation found gaps in media readiness,
cache ownership and expiry, image preflight, and Markdown entity handling.
These gaps can send unusable media or content to a quota-limited draft endpoint.

## What Changes

- Use the recommended mutable Markdown entities and an explicitly documented
  conservative weighted-length budget.
- Wait for reported media processing with bounded read-only status requests;
  stop on failure without replaying upload or draft POSTs.
- Bind persistent media reuse to its API and OAuth credential context and a
  known conservative expiry; reject unscoped legacy cache entries.
- **BREAKING**: limit simple image upload to validated static JPEG, PNG, GIF,
  and WebP at most 5 MB; reject unsupported formats and animations locally.
- Align the uploader documentation and the article's conversion showcase.
- After offline validation, make the user's explicitly authorized single
  Article draft attempt and retain its result without automatic retry.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `markdown-to-draftjs`: static image preflight and Markdown entity guarantees.
- `x-article-publication`: readiness, scoped media reuse, and lifetime handling.

## Impact

The uploader's Markdown, draft, and X API packages, their existing E2E suite,
operator documentation, and the blog conversion showcase. No new dependency,
credential mutation, Vault policy change, or article publication is required.
