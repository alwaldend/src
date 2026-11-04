## MODIFIED Requirements

### Requirement: Homepage navigation panel

The homepage SHALL display a bordered panel with space below the site header,
linking to GitHub, GitLab, Blog, Docs, and the project index in that order.
Individual project links SHALL follow in descending creation date order as
permanently visible, indented rows without a submenu or disclosure control.
Each row SHALL display its destination to the right of its title on desktop
and below it on narrow screens. Internal destinations SHALL follow the site's
configured navigation and generated project pages.

#### Scenario: Browse projects from the homepage

- **WHEN** a visitor opens the homepage
- **THEN** the panel offers both the project index link and all individual
  project links without requiring expansion
- **AND** each destination is visible alongside its title on desktop, while
  narrow screens keep the title and destination within the viewport
- **AND** newer projects precede older projects using landing creation dates
