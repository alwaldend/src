## ADDED Requirements

### Requirement: Natural-size raster default

WebP rendering SHALL default to scale 1, using the ceiling of each natural SVG canvas dimension before optional aspect-ratio padding. Explicit positive integer scales SHALL remain supported and multiply those dimensions before padding.

#### Scenario: Omit raster scale

- **WHEN** a WebP render omits its scale
- **THEN** the diagram is rasterized at its natural size
- **AND** optional aspect-ratio padding retains the existing centering and exact-ratio behavior

#### Scenario: Request a higher raster scale

- **WHEN** a WebP render explicitly selects scale 2 or 3
- **THEN** the diagram's pixel dimensions are multiplied by that scale
- **AND** its selected appearance and aspect-ratio padding are preserved
