## MODIFIED Requirements

### Requirement: Adaptive Mermaid images

Documentation Mermaid images published through Markdown and the SVG shortcode
SHALL use the site's shared light/dark palette and follow the selected color
mode, including system preference. They SHALL use clean, rounded nodes and
containers, straight connectors, embedded sans-serif typography, and the
native Dagre layout. Edge labels and container titles SHALL use rounded,
bordered caption plates styled through Mermaid input configuration. Actor
nodes MAY use ordinary text labels.
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
