## ADDED Requirements

### Requirement: Shared minimalist monochrome presentation

The homepage, blog, documentation, and project pages SHALL share a neutral
light/dark palette, flat interface surfaces, subtle borders, and consistent
typography. Prose links and keyboard focus SHALL remain visually identifiable.
Theme changes SHALL preserve responsive navigation and search.
Scrollbars SHALL use thin, muted thumbs, with rounded thumbs where browser
styling supports them. Taxonomy links SHALL use subtle rounded labels without
arrow shapes, including on the project index. Sidebar taxonomy items SHALL wrap compactly beneath distinct group headings,
with counts next to their labels. Documentation sidebars, headings, search
fields, code blocks, tables, and callouts SHALL use the same grayscale palette
and restrained typography. Sidebar group headings SHALL be visually distinct
from their item links, and page actions SHALL use plain rows without icons or
underlines, with visible hover and keyboard-focus states.
Project landing pages SHALL display their title.

#### Scenario: Read pages in either color mode

- **WHEN** a visitor selects light or dark mode
- **THEN** the shared canvas, text, navigation, and controls use a readable
  grayscale palette without decorative gradients or shadows

#### Scenario: Read documentation and project pages

- **WHEN** a visitor opens documentation, a code example, or a project landing
- **THEN** headings and controls share the site's typography, code follows the
  selected color mode, and sidebar links have subtle hover and active states
- **AND** project landing pages display their title before their description

#### Scenario: Navigate with a keyboard on a narrow viewport

- **WHEN** a visitor tabs through navigation and quick links on mobile
- **THEN** focused controls have a visible outline and content fits the viewport

### Requirement: Homepage navigation panel

The homepage SHALL display a bordered panel with space below the site header,
linking to GitHub, GitLab, Blog, Docs, and the project index in that order.
Individual project links SHALL follow in title order as permanently visible,
indented rows without a submenu or disclosure control. Each row SHALL display
its destination to the right of its title on desktop and below it on narrow
screens. Internal destinations SHALL follow the site's configured navigation
and generated project pages.

#### Scenario: Browse projects from the homepage

- **WHEN** a visitor opens the homepage
- **THEN** the panel offers both the project index link and all individual
  project links without requiring expansion
- **AND** each destination is visible alongside its title on desktop, while
  narrow screens keep the title and destination within the viewport
