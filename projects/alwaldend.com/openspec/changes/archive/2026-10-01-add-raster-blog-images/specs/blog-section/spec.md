## ADDED Requirements

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

The site SHALL provide a WebP rendering macro using the native light palette
and font owned by its SVG rendering configuration. A post MAY select a bundled
image through its `images` metadata for Open Graph and large Twitter cards,
without adding a duplicate image to the article body.

#### Scenario: Share the X uploader demonstration post

- **WHEN** the preview site renders the X uploader demonstration post
- **THEN** its diagram starts with Markdown and ends with an X Article draft
- **AND** the diagram uses the site palette and has space around its contents
- **AND** Open Graph and Twitter image metadata resolve the same bundled WebP
- **AND** offline conversion still produces exactly one image locator
