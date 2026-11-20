## 1. Favicon

- [x] 1.1 Add a fixed-contrast favicon variant to the shared S8 SVG.
- [x] 1.2 Verify the fingerprinted site output and header logo in both color modes.

## Evidence

The homepage, blog, documentation, project index, and 404 page referenced
one SVG whose URL matched its SHA-256 digest. Chromium rendered the favicon
with off-white ink on charcoal in both color modes. The unselected SVG
variant kept a transparent background and switched its ink from charcoal
in light mode to off-white in dark mode.
