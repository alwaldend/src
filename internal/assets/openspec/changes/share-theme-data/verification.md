## Independent theme extraction

This change owns identity data, generated website/application Sass, the bundled
Mermaid theme provider, and native SVG/raster JSON palette inputs. It retains
the trunk Mermaid version and font embedding API. The Three.js renderer,
Mermaid upgrade, and snapshots remain in PR #211 and use the existing site
configuration and compiled Sass palette APIs. Compatibility targets derive
those inputs from the new canonical theme data.

Validate the generated Sass against the original source, the website theme and
readability suites, native SVG/raster checks, application stylesheet consumers,
repository quality, and semantic lint. Delivery receipts own exact candidate
validation and publication state.
