## Design

Interpret a regular diagram according to nearby published examples. Forgejo,
Vault, Flux, DNS, download, and ingress use left-to-right flowcharts. The
site theme sets the classic look and curved connectors. Reuse that theme and
use simple semantic shapes and concise labels.

Compose the existing native SVG and WebP render rules in `mermaid_diagram`.
Use a single output stem for both formats and color schemes. The macro's
aggregate target builds all four files. Documentation that publishes SVGs
uses the SVG pair target. WebP consumers can depend on the specific variant.
Keep the existing WebP scale and aspect ratio for formats in use. Host
documentation consumes the macro's SVG target labels directly. The blog's
single dark WebP calls the native raster rule with its existing theme and
settings. Historical blog renders, renderer tests, and upstream comparisons
use their specific render rules.
