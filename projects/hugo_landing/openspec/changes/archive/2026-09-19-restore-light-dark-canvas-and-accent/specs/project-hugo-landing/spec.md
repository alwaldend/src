## MODIFIED Requirements

### Requirement: Produce reusable source and rendered sites

`al_hugo_landing` SHALL combine a selected project README, generated Hugo
configuration, shared layouts and assets, and declared theme dependencies into
a Hugo-ready archive. `al_hugo_landing_site` SHALL additionally render that
archive into a site output. The shared assets SHALL include the canvas and
accent styles that the repository's Hugo sites render, and the archive SHALL
expose those styles for the apex site to package into its own assets tree.

The shared canvas SHALL define distinct light and dark palettes: light mode
SHALL use a pure white (`#fff`) page canvas and dark mode SHALL use a pure
black (`#000`) page canvas, each with the matching Bootstrap palette for body
text, surfaces, borders, and links. A page resolved without a color-mode
attribute SHALL use the light palette. The shared accent SHALL be the single
source for the site's primary color, and every accent-coloured surface —
including links, buttons, badges, and the dark-mode link tint — SHALL derive
from it rather than restating its value.

#### Scenario: Define a landing site for a project

- **WHEN** the caller supplies a project, title, documentation target,
  documentation URL, and repository URL
- **THEN** the macros package those inputs with the reusable shell and can
  produce rendered site files

#### Scenario: Render the shared canvas

- **WHEN** a landing site is rendered
- **THEN** it uses the shared canvas styles, and the apex site obtains the same
  styles from this shell rather than from a copy in its own tree

#### Scenario: Consume the shared canvas from the apex site

- **WHEN** the apex site packages the shared canvas styles into its assets tree
- **THEN** the shared style file is a declared input of that site's build and no
  second copy of its rules exists in the apex site tree

#### Scenario: Select a color mode

- **WHEN** a rendered page resolves to the light color mode and then to the
  dark color mode
- **THEN** the light page canvas is pure white with the light-mode palette and
  the dark page canvas is pure black with the dark-mode palette

#### Scenario: Render the shared accent

- **WHEN** a rendered page links to another page or shows a primary button
- **THEN** the link and button colors derive from the shared accent value
  rather than from a site-local copy

#### Scenario: Render a page without a color-mode attribute

- **WHEN** a rendered page has no color-mode attribute on its root element
- **THEN** it uses the pure white light palette rather than the dark palette
