# blog-section Specification

## Purpose

Give alwaldend.com a dated article section that reuses the site's Docsy shell,
packaging, and outputs instead of a separate page type.

## Requirements

### Requirement: Package an article section

The site SHALL package a `blog` content section whose index page and posts are
declared as Bazel packages and included in the site source archive, so posts
are rendered from declared inputs rather than copied into the build.

#### Scenario: Build the site with a blog post

- **WHEN** the site source archive is built while the blog section contains a
  post package
- **THEN** the rendered site contains the post page under the blog section's
  published path

#### Scenario: Build the blog section index

- **WHEN** the site is built and the blog section contains no posts
- **THEN** the section still renders its index page without failing the build

### Requirement: Publish section outputs

The blog section SHALL publish an RSS feed and a print edition in addition to
its HTML pages, and SHALL be reachable from the site's main navigation.

#### Scenario: Request the blog feed

- **WHEN** a client requests the blog section's feed
- **THEN** the site serves a feed of the section's published posts

#### Scenario: Navigate to the blog section

- **WHEN** a visitor opens any site page that renders the main navigation
- **THEN** the navigation links to the blog section index

### Requirement: Render posts with author and date metadata

Rendered posts SHALL present their publication date and, when the post declares
one, its author, and the section index SHALL order posts by date.

#### Scenario: A post declares a date and author

- **WHEN** a post declares a publication date and an author
- **THEN** the rendered post shows both, and the section index lists the post in
  date order

### Requirement: Withhold unpublished posts

A post whose front matter declares it a draft SHALL be excluded from the
release build while remaining available to local previews, so the draft state
alone keeps a post off the deployed site.

#### Scenario: Release a site containing a draft post

- **WHEN** the site is built for release while the blog section contains a post
  that declares itself a draft
- **THEN** the release output contains no page, feed entry, or sitemap entry for
  that post

#### Scenario: Preview a draft post locally

- **WHEN** the site is built for local preview while the blog section contains a
  post that declares itself a draft
- **THEN** the preview output contains the post so its author can review it

### Requirement: Reference images a syndication target can accept

A post's images SHALL use a media type the site publishes and a syndication
target can accept. Because the X Articles media upload endpoint accepts a fixed
set of image media types that excludes `image/svg+xml`, a post that is intended
for syndication SHALL NOT reference an SVG image. A Mermaid diagram in such a
post SHALL be referenced through a raster render beside `index.md`, with the
`.mmd` source remaining authoritative and the render produced by the
repository's maintained Mermaid pipeline. A post that is not intended for
syndication MAY keep an SVG diagram render.

#### Scenario: Reference a diagram from a syndicated post

- **WHEN** a post intended for syndication contains a Mermaid diagram
- **THEN** it references a raster render of that diagram
- **AND** the render is produced from the `.mmd` source rather than hand-committed
- **AND** no SVG is referenced as the post's image

#### Scenario: Reference a diagram from a site-only post

- **WHEN** a post that is not intended for syndication contains a Mermaid diagram
- **THEN** it may reference the SVG render the section already publishes
- **AND** no existing post is required to change its image references

#### Scenario: Read the authoring guidance

- **WHEN** an author follows the blog authoring guidance for a new post
- **THEN** the guidance states that a syndicated post's images use a media type
  the syndication target accepts
- **AND** it names the raster render target a diagram uses

### Requirement: Reuse native diagrams for social previews

The site SHALL provide a WebP rendering macro using the selected native light or dark palette and font owned by its SVG rendering configuration. A post MAY select a bundled image through its `images` metadata for Open Graph and large Twitter cards, without adding a duplicate image to the article body. The authoring skill SHALL encourage a representative image for every post and reuse a suitable opening image when available. Missing images SHALL NOT block delivery or cause an invented asset.

#### Scenario: Share the X uploader demonstration post

- **WHEN** the preview site renders the X uploader demonstration post
- **THEN** its diagram starts with Markdown and ends with an X Article draft
- **AND** separate stages show offline conversion, media upload or cache reuse,
  attachment of media IDs to the banner and body, and draft creation
- **AND** the diagram matches the site's dark palette with modest visible margins
- **AND** Open Graph and Twitter image metadata resolve the same bundled WebP
- **AND** offline conversion still produces exactly one image locator

#### Scenario: Add metadata to a post with an opening image

- **WHEN** a post has a suitable opening image but no selected social image
- **THEN** its metadata selects that image or its generated raster projection
- **AND** adding metadata preserves the article body and draft state

#### Scenario: Prepare a post without an image

- **WHEN** a post has no suitable image
- **THEN** guidance encourages one while permitting the metadata to remain absent
