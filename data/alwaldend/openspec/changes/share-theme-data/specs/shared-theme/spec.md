## ADDED Requirements

### Requirement: Shared identity has one source

The website, application stylesheets, and current neutral SVG diagrams SHALL derive their
light/dark colors and diagram typography from authoritative identity data.
Consumer-specific mappings SHALL preserve the current appearance, and the
historical hand-drawn Mermaid preset SHALL remain separately selectable.

#### Scenario: Render all consumers

- **WHEN** the site, application stylesheets, and native diagrams are built
- **THEN** they consume the shared identity data without recovering colors from CSS
- **AND** their current light/dark colors and pinned diagram font are preserved
- **AND** builds and delivered diagrams require no remote theme assets
