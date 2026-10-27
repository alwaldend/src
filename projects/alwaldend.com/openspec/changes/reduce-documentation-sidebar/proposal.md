## Why

The website archive contains 1.34 GB of expanded files. Repeating the complete
documentation tree in every page accounts for most HTML: a representative
429 KB project page carries about 380 KB of sidebar navigation.

## What Changes

- Use the pinned Docsy theme's scoped sidebar roots for documentation owners.
- Retain section navigation and parent links, without a theme fork or new dependency.
- Measure archive and representative page sizes and verify browser navigation.
- Upgrade Docsy to the latest stable release and retest first-interaction latency.
- Build the Lunr index before publication to remove synchronous browser indexing.
- Deliver the website fixes independently of infrastructure changes.

## Capabilities

### New Capabilities

- `scoped-documentation-navigation`: Documentation owners have bounded navigation trees.

## Impact

The main site's Hugo configuration and its documentation. Content, release
payloads, and theme pins retain their existing owners.
