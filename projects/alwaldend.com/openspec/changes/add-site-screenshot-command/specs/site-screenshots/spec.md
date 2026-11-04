## Purpose

Provide a repeatable local command that renders the built site's pages into
reviewable full-page screenshots without requiring a running preview server.

## ADDED Requirements

### Requirement: Capture built site pages

The screenshot command SHALL capture the homepage, project landing index, and
project documentation index by default. It SHALL support explicit site-relative
paths, viewport dimensions, color scheme, and output directory. It SHALL use
the repository's pinned browser tooling, load only built local assets, and
produce full-page PNGs plus a manifest identifying the pages and dimensions.
It SHALL fail for missing requested pages or invalid options and SHALL clean up
its temporary browser profile.

#### Scenario: Capture the standard indexes

- **WHEN** the screenshot command runs with its default pages
- **THEN** it writes screenshots of `/`, `/projects/`, and `/docs/projects/`
- **AND** its manifest identifies all three successful captures

#### Scenario: Select a page and viewport

- **WHEN** the caller supplies paths, dimensions, or a color scheme
- **THEN** the output reflects those options rather than requiring source edits

#### Scenario: Work without a preview server

- **WHEN** no site server is running
- **THEN** local assets still render and no external requests leave the browser

#### Scenario: Reject an unavailable page

- **WHEN** a requested page is missing from the built site
- **THEN** the command fails rather than reporting a successful screenshot
