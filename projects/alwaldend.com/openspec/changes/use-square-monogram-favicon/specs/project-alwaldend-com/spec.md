## ADDED Requirements

### Requirement: Shared square monogram favicon

The homepage, blog, documentation, and project pages SHALL reference the
approved S8 square monogram as their SVG favicon. The published resource
SHALL preserve the canonical artwork in `data/alwaldend` and use a
content-derived URL so artwork changes invalidate cached icons.

#### Scenario: Build pages with the selected artwork

- **WHEN** the site is built
- **THEN** representative homepage, blog, documentation, and project pages
  reference the same fingerprinted SVG favicon
- **AND** the referenced SVG matches the canonical S8 source, including
  its centered geometry and light/dark color behavior
