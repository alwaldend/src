## MODIFIED Requirements

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
