## ADDED Requirements

### Requirement: Render selected images above blog content

Hugo SHALL render every image declared in a blog post's front matter `images`
array, in declaration order, before the Markdown body. The print edition SHALL
include the same opening images. Rendering SHALL reuse the site's image
handling, retain supplied alternative text, and fit narrow viewports without
cropping or horizontal overflow. A post without `images` SHALL have no empty
image container. Social image metadata SHALL continue to select the declared
resources.

Authoring guidance SHALL explain that these images are displayed automatically
and do not need matching opening Markdown references. Removing existing
duplicates SHALL preserve later body illustrations and supplied prose.

#### Scenario: Read a post with selected images

- **WHEN** a post declares one or more front matter images
- **THEN** the rendered page displays them above its first body content in order
- **AND** each image retains its resource alternative text when supplied
- **AND** the same images appear in the post's print output

#### Scenario: View a post on a narrow screen

- **WHEN** the viewport is narrower than the selected image's intrinsic width
- **THEN** the image fits the content column without clipping or horizontal overflow

#### Scenario: Read a post without images

- **WHEN** the post omits `images` or declares an empty array
- **THEN** the layout emits no empty image group

#### Scenario: Migrate an existing opening illustration

- **WHEN** a post's opening Markdown image repeats its selected metadata image
- **THEN** the matching opening reference is removed from source
- **AND** its useful alternative text is retained for the automatically rendered image
- **AND** later illustrations, prose, and publication metadata are preserved

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
- **AND** Hugo displays that WebP above the post without a matching opening Markdown image
- **AND** offline X conversion records the banner without a duplicate body image locator

#### Scenario: Add metadata to a post with an opening image

- **WHEN** a post has a suitable opening image but no selected social image
- **THEN** its metadata selects that image or its generated raster projection
- **AND** guidance identifies the now-redundant opening body reference for removal
- **AND** unrelated body content and draft state are preserved

#### Scenario: Prepare a post without an image

- **WHEN** a post has no suitable image
- **THEN** guidance encourages one while permitting the metadata to remain absent
