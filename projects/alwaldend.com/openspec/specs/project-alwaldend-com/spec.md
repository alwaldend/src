# Alwaldend.com Specification

## Purpose

Describe the repository's main Hugo and Docsy documentation website and its
packaged build, local preview, and explicit GitHub Pages deployment behavior.

Sources: [project README](../../../README.md),
[site targets](../../../BUILD.bazel), and
[Hugo dependencies](../../../include.MODULE.bazel).
Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed 2026-09-08. This specification describes source behavior, not a live
deployment observation.

## Requirements

### Requirement: Build packaged repository documentation

The `site` target SHALL render the packaged repository documentation tree with
the configured Hugo, Docsy, Bootstrap, and local layouts and assets. Packaged
`README.md` documentation SHALL be renamed to Hugo `_index.md` pages.

#### Scenario: Build the main documentation site

- **WHEN** the `//projects/alwaldend.com:site` target is built
- **THEN** it produces rendered site output from declared documentation and
  asset inputs, with the release configuration using `https://alwaldend.com`.

### Requirement: Resolve links relative to documentation sources

The site SHALL resolve Markdown links and images relative to their source
directory, map packaged documentation pages and resources to published URLs,
and preserve unknown internal destinations rather than silently redirecting
them to GitHub.

#### Scenario: Documentation contains an unresolved internal link

- **WHEN** a Markdown destination has no matching packaged page or resource
- **THEN** it remains unresolved in the rendered output rather than being
  silently redirected to GitHub.

### Requirement: Preserve distinct anchors in combined print pages

Print rendering SHALL scope document IDs and their fragment and control
references to the source document so that combined documents retain distinct
anchors.

#### Scenario: Two printed documents contain the same heading anchor

- **WHEN** the documents are combined into a print page
- **THEN** their scoped IDs and references remain distinguishable.

### Requirement: Separate preview and publication

`site_serve` SHALL serve built content on `127.0.0.1:1313`. The deployment
archive SHALL contain the apex `CNAME`, and the explicit deployment workflow
SHALL write `.nojekyll` and push changed output to the configured Pages branch.
Building or previewing the site SHALL NOT itself publish it. The rendered site
SHALL consume the shared canvas styles provided by the reusable site shell
rather than restating them.

#### Scenario: Preview a local candidate

- **WHEN** the local `site_serve` target runs
- **THEN** the rendered site is available on the loopback preview address without
  executing the deployment target

#### Scenario: Render the shared canvas

- **WHEN** the apex site is built
- **THEN** the shared canvas style file is a declared input of its assets package
  and the rendered pages use the shared background and foreground palette

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
