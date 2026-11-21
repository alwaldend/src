## 1. Favicon

- [x] 1.1 Update the favicon variant of the shared S8 artwork and its documentation.
- [x] 1.2 Verify the fingerprinted resource and inspect the icon at mobile tab size.
- [x] 1.3 Verify the original S8 geometry and the 16-pixel result on a white backing.

## Evidence

The built homepage, blog, docs, project index, Anki gRPC page, and 404 page
reference the same SHA-256-fingerprinted SVG. The inline header retains its
adaptive foreground. The published paths match the original monogram, and
Chromium rendered the 16-pixel icon on a white circular backing.
