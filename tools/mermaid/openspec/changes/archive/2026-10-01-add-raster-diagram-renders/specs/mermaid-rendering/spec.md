## ADDED Requirements

### Requirement: Raster diagram renders

The tool SHALL expose a rule that renders a Mermaid source to a raster image in
a media type a service upload endpoint can accept, beginning with WebP. The
raster render SHALL reuse the selected SVG render contract rather than defining a
second appearance: the same theme, the same pinned browser and fonts, and
the same post-processing policy as the SVG rule. Native rendering SHALL use
the consumer's palette and font and encode the light SVG without post-processing;
the default SHALL retain the historical theme and paint-order pass.
Encoding MUST go through the pinned
browser, so the rule introduces no image-conversion dependency. The render SHALL
be a projection of the SVG render — the same maintained document, encoded to a
raster — rather than a second layout. The render SHALL be hermetic — no host
browser, host font, or network fetch — and deterministic for a given source,
theme, and configuration.

#### Scenario: Render a diagram to a raster image

- **WHEN** a diagram source is rendered through the raster rule
- **THEN** it produces an image in a service-acceptable media type
- **AND** the render uses the shared theme, the pinned browser, and the pinned
  fonts rather than host or network inputs

#### Scenario: Keep one appearance across formats

- **WHEN** the same diagram source is rendered to SVG and to the raster format
- **THEN** both renders apply the same selected theme and post-processing policy
- **AND** neither render introduces an appearance of its own
- **AND** the raster canvas reproduces the SVG render's geometry at the rule's scale

#### Scenario: Render a native site diagram

- **WHEN** the raster rule selects native mode with a palette and label font
- **THEN** it encodes the native light SVG with the same palette and font
- **AND** it performs no historical SVG post-processing
- **AND** native mode and a palette must be selected together

#### Scenario: Reproduce the natural canvas

- **WHEN** a diagram's natural canvas differs from the browser's default
  replaced-element size
- **THEN** the raster is captured at the SVG's own `viewBox` geometry rather than
  the default size
- **AND** its canvas equals the SVG render's geometry at the rule's scale

#### Scenario: Depend on no image-conversion toolchain

- **WHEN** the raster rule's inputs are inspected
- **THEN** it uses the pinned browser to encode the raster format
- **AND** it declares no image-conversion tool or library

### Requirement: Raster renders belong to their source

A raster render SHALL be reproducible from the diagram source, with the source
remaining authoritative. A consumer that must keep the image as a checked-in
asset SHALL obtain it from the rule's output rather than by hand, so the asset
can be regenerated and compared against a fresh render.

#### Scenario: Keep a checked-in raster asset current

- **WHEN** a consumer checks in a diagram's raster render
- **THEN** the asset is the rule's output for that source
- **AND** the asset can be compared against a fresh render of the same source

#### Scenario: Update the source

- **WHEN** a diagram source changes
- **THEN** its raster render changes when regenerated from the new source
- **AND** the source remains the artifact an author edits
