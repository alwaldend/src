---
title: Alwaldend identity assets
description: Selected monogram artwork
---

`assets/monogram_s8.svg` is the canonical S8 graduated square monogram.
It preserves the approved 64 by 64 geometry, angled `l` cap, and full
clipped underline. The spiral and `l` meet at y=45.5; the underline begins
at y=50. The transparent SVG uses charcoal in light mode and off-white in
dark mode. A 2.5-unit horizontal translation centers the complete mark,
leaving equal 8.5-unit side margins and equal 8-unit vertical margins.

The SVG is available as `//data/alwaldend:monogram` for internal repository
consumers. The `projects/alwaldend.com` assets package consumes it directly
as the site's fingerprinted favicon. WebP previews are generated output,
not source.

## Reproduce the preview

The 512 by 512 dark preview uses `#ededed` ink on a `#212121` background.
Render the SVG at 448 by 448, then add 32 pixels on each side to provide
the additional preview padding while keeping the mark centered.
It was rendered using the existing host ImageMagick 7.1.2-34 with its SVG
and WebP delegates. From the repository root, derive this presentation from
the canonical source into ignored output:

```sh
mkdir -p out/s8-favicon/tmp
node <<'NODE'
const fs = require('fs');
const source = fs.readFileSync(
  'data/alwaldend/assets/monogram_s8.svg', 'utf8');
const dark = source
  .replace('viewBox="0 0 64 64"',
    'viewBox="0 0 64 64" width="448" height="448" color="#ededed"')
  .replace(/<style>[\s\S]*?<\/style>/,
    '<rect width="64" height="64" fill="#212121"/>');
fs.writeFileSync('out/s8-favicon/monogram_s8_dark.svg', dark);
NODE
TMPDIR="$PWD/out/s8-favicon/tmp" \
  MAGICK_TEMPORARY_PATH="$PWD/out/s8-favicon/tmp" \
  magick -background '#212121' out/s8-favicon/monogram_s8_dark.svg \
  -gravity center -extent 512x512 \
  -define webp:lossless=true out/s8-favicon/monogram_s8_dark.webp
```

Inspect the resulting WebP after regeneration. Rendering can vary with the
host SVG delegate; the SVG path geometry remains authoritative.
