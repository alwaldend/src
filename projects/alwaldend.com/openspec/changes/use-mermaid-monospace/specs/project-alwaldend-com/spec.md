## MODIFIED Requirements

### Requirement: Adaptive Mermaid images

Documentation Mermaid images published through Markdown and the SVG shortcode
SHALL use the site's shared light/dark palette and follow the selected color
mode, including system preference. They SHALL use clean, rounded nodes and
containers, smoothly curved connectors, embedded monospace typography, and the
native Dagre layout. Edge labels and container titles SHALL use rounded,
bordered caption plates styled through Mermaid input configuration. Actor
nodes MAY use ordinary text labels. Normal flowchart, state, class, and entity
relationship connectors SHALL have greater visual weight than thin node
outlines. Edge captions and sequence lifelines SHALL remain visually
subordinate to the content and messages they annotate. Mindmap branches and
labels SHALL remain clearly distinguishable from their backgrounds in both
modes. Thick, invisible, dashed, and marker-based relationship semantics SHALL
be preserved.
The documentation pipeline SHALL publish the renderer's SVG output without
custom SVG post-processing. Mode changes SHALL preserve geometry, fonts, and
alternative text. Blog diagrams SHALL retain their existing publication files
and original presentation.

#### Scenario: Change the selected theme

- **WHEN** a visitor selects light or dark mode while reading a documentation Mermaid diagram
- **THEN** its canvas, text, lines, and labels follow the selected mode without reload
- **AND** an explicit selection takes precedence over the operating system
- **AND** only one diagram image is displayed and exposed to accessibility tools

#### Scenario: Use system preference

- **WHEN** system mode is selected and the system color preference changes
- **THEN** diagrams follow the resulting site color mode

#### Scenario: Read on mobile or in print

- **WHEN** a diagram is displayed on a narrow page or combined print page
- **THEN** it remains an accessible image without document ID collisions
- **AND** its geometry and labels remain intact
- **AND** print uses the light variant regardless of the selected screen mode

#### Scenario: Open an image separately

- **WHEN** a visitor opens either generated SVG outside the page
- **THEN** it retains its intended palette and embedded font without site CSS

#### Scenario: Read a historical blog diagram

- **WHEN** a visitor reads a blog diagram in any selected color mode
- **THEN** the image retains its existing publication URL and appearance
- **AND** paired-image selection does not apply to blog images

#### Scenario: Read a caption

- **WHEN** a documentation diagram has an edge label or container title
- **THEN** its bordered caption follows the selected palette
- **AND** the caption fits its measured label frame without SVG post-processing

#### Scenario: Read a compact nested diagram

- **WHEN** a documentation diagram contains nested containers
- **THEN** node and rank gaps keep the overall diagram compact
- **AND** parent and child title plates have visible separation
- **AND** connectors use smooth bends without SVG post-processing

#### Scenario: Trace relationships at documentation width

- **WHEN** the representative gallery is displayed at 728px image width
- **THEN** normal graph connectors retain at least one CSS pixel of stroke
- **AND** thick edges remain stronger, invisible edges remain invisible, and
  dashed relationships and their markers remain distinct

#### Scenario: Read sequence and mindmap hierarchy

- **WHEN** a sequence or mindmap is rendered in either color mode
- **THEN** sequence guides remain visible and lighter than message strokes
- **AND** mindmap branches contrast with the canvas and root labels contrast
  with their node surface

#### Scenario: Use monospace without a host font

- **WHEN** a current site diagram is opened without external fonts
- **THEN** its embedded labels retain equal glyph advances and their measured geometry

## ADDED Requirements

### Requirement: Saved Mermaid color schemes

The project SHALL retain the Cobalt drafting, Warm editorial, Teal surfaces,
and Violet technical light/dark color schemes as named reusable definitions.
Saving these alternatives SHALL preserve the active neutral site palette.

#### Scenario: Save an alternative without activating it

- **WHEN** the saved color schemes are made available for future use
- **THEN** each retains its previously proposed light and dark colors
- **AND** ordinary site renders continue using the shared neutral palette
