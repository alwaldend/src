## Contract

Use the existing padding implementation for a default 5:2 raster canvas.
The encoder defaults to [5, 2]; the Bazel rule exposes that same default.
An explicit empty list selects natural geometry; the CLI represents it as
`--aspect-ratio natural`. Positive integer ratio overrides remain supported.

## Failure cases and verification

Before implementation, adapt the existing end-to-end renderer checks for:
divergent rule/CLI/encoder defaults, distorted or missing labels, changed content
scale or palette, broken natural-size opt-out, and invalid ratios. The fixture
comparison uses independent browser layout and captures actual/reference WebPs.
Keep native geometry checks explicitly on natural output; exercise omitted
ratios through the classic fixtures and direct encoder call. Inspect the complete
3x-ui render and verify its physical pixel dimensions have exact 5:2 proportions.

The uploader's existing SVG geometry comparison explicitly consumes a natural
WebP fixture. Its conversion behavior and published blog assets are unchanged;
the renderer's default fixture separately exercises the padded contract.
