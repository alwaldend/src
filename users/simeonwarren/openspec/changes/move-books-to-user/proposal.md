## Why

The book catalog and its miscellaneous documentation page belong to
`users/simeonwarren`. Keeping them together makes their ownership explicit.

## What Changes

- Move `data/books` into `users/simeonwarren/books`, preserving catalog and
  cover contents.
- Consolidate the miscellaneous books page into the moved README and package
  its data and covers through the user's documentation target.
- Serve the books page at its user-owned URL without redirects from the old
  data and miscellaneous URLs.

## Capabilities

No capability requirements change. This source and documentation relocation
uses `skip_specs: true`.

## Impact

Book source paths, Bazel documentation packaging, and the documentation URL
change. The existing shared books renderer remains in use. No new dependency
or deployment is required.
