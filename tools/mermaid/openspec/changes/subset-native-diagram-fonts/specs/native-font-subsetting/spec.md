## ADDED Requirements

### Requirement: Subset standalone native fonts

Native diagrams SHALL embed a subset of their pinned font covering rendered
characters supported by that font, including characters decoded from entities.
Subsetting SHALL preserve glyph metrics and outlines and SHALL require no
network font access. Delivered SVG geometry SHALL remain Mermaid-generated.

#### Scenario: Render light and dark variants

- **WHEN** a native diagram is rendered through Bazel
- **THEN** both variants contain loadable reduced fonts with preserved glyph metrics
- **AND** their labels remain readable when embedded as standalone images

#### Scenario: Font subsetting fails

- **WHEN** the declared subset tool cannot process the font
- **THEN** the render fails rather than silently substituting another font
