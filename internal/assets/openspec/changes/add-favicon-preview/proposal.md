## Why

The S8 artwork needs a standalone showcase for checking its legibility at
browser tab sizes.

## What Changes

Add a Bazel target in shared assets that generates a self-contained HTML
preview from the canonical monogram SVG.

## Capabilities

### New Capabilities

- `favicon-preview` covers local inspection of the shared favicon variant.

## Impact

The showcase builds from the shared monogram independently of the website.
