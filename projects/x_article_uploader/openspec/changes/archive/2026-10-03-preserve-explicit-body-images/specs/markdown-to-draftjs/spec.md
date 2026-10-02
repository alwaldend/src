## REMOVED Requirements

### Requirement: Banner image metadata

**Reason**: The user explicitly removed opening-image suppression. Replace the
old suppression-bearing requirement with independent banner and body conversion.

**Migration**: Keep front matter banner selection, preserve every explicit body
image, and remove redundant opening references in the Hugo post sources.

## ADDED Requirements

### Requirement: Independent banner and body image conversion

The converter SHALL resolve the first entry of the optional front matter
`images` array as the article banner. It SHALL record a separate banner locator
with the post-relative path, post package, content digest, and detected media
type, using the same image validation as body images. Banner metadata MUST NOT
insert or remove any body image blocks, entities, or locators. The converter
SHALL preserve every valid explicit Markdown image in source order, including
opening images, repeated references, and images whose bytes match the banner.
Malformed metadata or an unusable selected image SHALL fail conversion. Banner
selection MUST NOT bypass body image containment or media validation.

#### Scenario: Select a banner independently of body images

- **WHEN** a post selects a local image in the first `images` entry
- **THEN** conversion records its banner locator separately from body locators
- **AND** later `images` entries do not select additional banners
- **AND** body text and images retain their original order

#### Scenario: Preserve an opening image that matches the banner

- **WHEN** a post opens with explicit Markdown images whose validated bytes
  match its selected banner, in one or several paragraphs
- **THEN** every occurrence produces its normal body block, entity, and locator
- **AND** the separate banner locator remains present

#### Scenario: Keep metadata-only images out of the body

- **WHEN** a post selects a banner in `images` and contains no Markdown images
- **THEN** conversion records the banner locator
- **AND** it creates no body image blocks, entities, or locators from metadata

#### Scenario: Preserve a secondary front matter image

- **WHEN** a body image also appears in a later `images` entry
- **THEN** the explicit Markdown image remains in the body
- **AND** that metadata entry creates no additional banner or body occurrence

#### Scenario: Convert without a selected banner

- **WHEN** `images` is absent or empty
- **THEN** the artifact has no banner locator and preserves existing conversion behavior

#### Scenario: Refuse an invalid selected banner

- **WHEN** `images` is not an array of strings, or its selected image is missing,
  unsupported, remote, or outside the post package
- **THEN** conversion fails with a diagnostic before any network operation

#### Scenario: Validate a body image that aliases the banner

- **WHEN** a body image aliases the banner's bytes through an escaping symlink
  or an unsupported media type
- **THEN** normal conversion diagnostics refuse that body reference
- **AND** matching bytes do not bypass validation

## MODIFIED Requirements

### Requirement: Images reference uploaded media

The `image` entity SHALL carry its media as `media_items`, whose entries require
a `media_category` and a `media_id`; `data.url` is documented for `link`
entities and MUST NOT be treated as an image source. Each explicit body image
SHALL become an `atomic` block plus an `image` entity, and unresolved images SHALL remain
explicitly unresolved so publication can resolve or reject them rather than
silently emitting a broken article. Because the media endpoints accept only a
fixed set of image media types that excludes `image/svg+xml`, the converter
SHALL report an image whose media type the upload endpoints cannot accept, with
its source position, rather than emitting an entity that publication can never
resolve. The declared extension SHALL be accepted as the media type only when
the file's bytes carry that format's own signature, so a file renamed to an
accepted raster extension — even one accepted format renamed to another — is
reported rather than sent to the upload endpoint.

#### Scenario: Convert an inline image

- **WHEN** a post contains an explicit Markdown image with alt text
- **THEN** it becomes an `atomic` block with an `image` entity
- **AND** the entity preserves the alt text as its caption, resolving Markdown
  escapes and character references the way prose is resolved while a raw span
  inside the alt text keeps its literal source

#### Scenario: Publication has not yet resolved an image

- **WHEN** a converted document references an image whose media has not been uploaded
- **THEN** the unresolved image is distinguishable from a resolved one
- **AND** publication refuses to send an unresolved image as though it were resolved

#### Scenario: An image cannot be uploaded

- **WHEN** a post references an image whose media type the media upload endpoints do not accept
- **THEN** the converter reports the image with its source position
- **AND** it does not emit an entity that publication could never resolve

#### Scenario: An image reference resolves outside the post directory

- **WHEN** a post's image reference is a syntactically safe relative path whose
  canonical target lies outside the post's directory, such as a symlink to a
  file elsewhere
- **THEN** the converter reports the image with its source position rather than
  reading bytes outside the post boundary
- **AND** it does not record a locator that publication would later reject, so a
  converted artifact remains publishable

#### Scenario: Image bytes do not match the declared extension

- **WHEN** a post references an image whose bytes are not the format its extension declares
- **THEN** the converter reports the image with its source position instead of treating it as the declared type
- **AND** it does not emit an entity the upload endpoint would reject

#### Scenario: An accepted format is renamed to another accepted extension

- **WHEN** bytes of one accepted raster format carry the extension of a different accepted raster format
- **THEN** the converter reports the mismatch rather than accepting the sniffed type
- **AND** it does not fall back to another accepted media type

#### Scenario: An image reference requires Markdown escaping

- **WHEN** a local image reference names a file whose name needs Markdown
  escaping, as in `![plot](plot\(final\).png)`
- **THEN** the reference is resolved before filesystem lookup, so the file it
  names is read rather than reported unreadable
- **AND** the recorded locator names the resolved path

#### Scenario: A symlink's target has a different extension than the reference

- **WHEN** an image reference such as `image.png` is an in-directory symlink to a
  file with a different raster extension, such as `actual.webp`
- **THEN** the declared-extension check classifies the image by the reference the
  author wrote, not by the canonical target's name
- **AND** bytes that do not match the reference's declared format are reported
  rather than accepted through the target's extension

#### Scenario: An inline image splits a styled span

- **WHEN** an inline image interrupts a styled or linked span
- **THEN** the text before the image keeps its style or entity range
- **AND** the text after the image keeps its style or entity range

#### Scenario: A link wraps only an image

- **WHEN** a Markdown image is the entire content of a link, as in
  `[![alt](image.png)](target)`
- **THEN** the image is still emitted as an `atomic` block with its `image`
  entity
- **AND** the lost destination is reported as a continuing diagnostic naming the
  link target, because the link has no text span to carry it

#### Scenario: A link wraps an image inside an inline wrapper

- **WHEN** the link's only content is an image nested inside a style wrapper,
  as in `[**![alt](image.png)**](target)`
- **THEN** the enclosing link is found by ancestry rather than by immediate
  parentage
- **AND** the lost destination is still reported, because the link has no text
  outside the image

#### Scenario: A link carries text beside an image

- **WHEN** a link contains non-whitespace text of its own in addition to an
  image, as in `[label ![alt](image.png)](target)`
- **THEN** the destination survives as a `link` entity over that text
- **AND** the image, which still cannot carry the link, is reported as a
  continuing diagnostic naming the destination, so the partial loss is not
  hidden merely because sibling text kept the link

#### Scenario: A post follows the raster-only policy

- **WHEN** a post's images all use a media type the media upload endpoints accept
- **THEN** conversion completes with an `image` entity for each explicit body image
- **AND** no unacceptable-image diagnostic is reported for that post
