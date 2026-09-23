## ADDED Requirements

### Requirement: Documentation titles and hierarchy

Packaged documentation without front matter that begins with a level-one
heading SHALL derive its page title from that heading. The title SHALL appear
once in the article and consistently in navigation, breadcrumbs, search,
and print output. The original heading anchor SHALL remain available.
Missing ancestor section pages SHALL preserve the documentation hierarchy
without replacing existing source metadata or leaf bundles.

#### Scenario: Browse an OpenSpec specification

- **WHEN** a visitor opens the owned-dns specification under the PVE documentation
- **THEN** the page and navigation show its title and its OpenSpec ancestry
- **AND** every breadcrumb destination resolves to a page

### Requirement: Responsive navigation states

Expanded and collapsed documentation trees SHALL preserve readable labels,
distinct expansion controls, full-row highlights, and visible keyboard focus.
Collapsed descendants SHALL be hidden from interaction. Page contents and
navigation SHALL fit narrow viewports, with wide code and tables scrolling
inside their containers. In-page TOC links SHALL reach visible headings.

#### Scenario: Expand and collapse a nested branch

- **WHEN** a visitor toggles a branch on a desktop or mobile viewport
- **THEN** its descendants appear or disappear without navigation
- **AND** activating a row opens its page while retaining the active ancestry

#### Scenario: Follow a long in-page heading

- **WHEN** a visitor follows a wrapped TOC link
- **THEN** the corresponding heading remains visible below the site header
