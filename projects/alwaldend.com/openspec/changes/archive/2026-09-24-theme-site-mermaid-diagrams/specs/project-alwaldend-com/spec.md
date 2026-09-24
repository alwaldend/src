## ADDED Requirements

### Requirement: Adaptive Mermaid images

Documentation Mermaid images rendered with the site preset through Markdown and the SVG
shortcode SHALL use the site's shared light/dark palette and follow the selected
color mode, including system preference. They SHALL use clean, rounded nodes and containers, straight
connectors, and embedded sans-serif typography. Color adaptation SHALL preserve
the rendered geometry, fonts, and alternative text. Blog diagrams SHALL retain
their existing publication files and original presentation.

#### Scenario: Change the selected theme

- **WHEN** a visitor selects light or dark mode while reading a documentation Mermaid diagram
- **THEN** its canvas, text, lines, labels, and actor icons adapt without reload
- **AND** an explicit selection takes precedence over the operating system

#### Scenario: Use system preference

- **WHEN** system mode is selected and the system color preference changes
- **THEN** diagrams follow the resulting site color mode

#### Scenario: Read on mobile or in print

- **WHEN** a diagram is displayed on a narrow page or combined print page
- **THEN** it remains an accessible image without document ID collisions
- **AND** its geometry and labels remain intact

#### Scenario: Read a historical blog diagram

- **WHEN** a visitor reads a blog diagram in any selected color mode
- **THEN** the image retains its existing publication URL and appearance
- **AND** the site's diagram adaptation does not rewrite it
