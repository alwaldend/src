## Purpose

Make the site's dependency structure inspectable through generated diagram
artifacts that follow the candidate Bazel graph and existing site appearance.

## ADDED Requirements

### Requirement: Build companion dependency diagrams

The site package SHALL expose Mermaid, light/dark SVG, and light/dark WebP build targets
analyzing the site root at a documented depth. These artifacts SHALL use the
shared site renderer and SHALL remain outside the site's own dependency closure.

#### Scenario: Inspect site dependencies

- **WHEN** the site's dependency WebP target is built
- **THEN** the image shows the current site's dependency relationships
- **AND** the normal site target still builds without a graph dependency cycle

#### Scenario: Inspect the dark dependency overview

- **WHEN** the site's dark dependency WebP target is built
- **THEN** it uses the site's dark palette
- **AND** its depth-eight view reaches site content and documentation owners
- **AND** documentation states the depth and excluded dependency categories

#### Scenario: Inspect ownership by package

- **WHEN** the site's mindmap companion target is built
- **THEN** it reuses the depth-eight capture and lists the same selected targets
- **AND** package branches describe ownership rather than dependency edges
- **AND** SVG and dark WebP outputs use the shared site renderer
