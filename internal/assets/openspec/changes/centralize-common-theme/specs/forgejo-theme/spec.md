## ADDED Requirements

### Requirement: Forgejo uses the common palette

Forgejo SHALL apply the common theme's canvas, text, surfaces, borders, and link
colors in its Light, Dark, and Auto modes. It SHALL use the shared interface
font stack and corner radii in those modes.
Navigation weight, prose-link underlines, and keyboard focus outlines SHALL
use the same control tokens as the website.

#### Scenario: An explicit appearance mode is selected

- **WHEN** Forgejo uses Light or Dark
- **THEN** the page uses the corresponding common palette regardless of browser preference

#### Scenario: Auto appearance follows the browser

- **WHEN** Forgejo uses Auto
- **THEN** the page uses the common palette for the browser's preferred color scheme

#### Scenario: A pull request diff is opened on a phone

- **WHEN** the user views the files tab at a narrow viewport
- **THEN** diff controls can wrap and long lines can scroll within their file body

### Requirement: Forgejo uses the shared monogram

Forgejo SHALL serve the canonical Alwaldend SVG for its browser favicon and
navigation logo.

#### Scenario: A browser loads a Forgejo page

- **WHEN** the browser requests the page’s SVG favicon and navigation logo
- **THEN** both image URLs return the canonical Alwaldend monogram
