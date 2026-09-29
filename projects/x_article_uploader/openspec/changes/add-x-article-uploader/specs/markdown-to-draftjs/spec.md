## Purpose

Compile a Markdown post into the DraftJS `content_state` document that the X
Articles API accepts, so the Markdown that renders the blog also produces the
article body without hand-editing.

## ADDED Requirements

### Requirement: Source post input model

The converter SHALL treat a blog source post as YAML front matter plus a
Markdown body, and MUST NOT emit front matter as article text. It SHALL expose
the parsed metadata separately from the `content_state` document, and SHALL
report the parsed `title` as document metadata because draft creation requires
a title that `content_state` does not carry. A source whose front matter has no
non-empty `title` SHALL fail conversion with a diagnostic naming the source,
rather than producing a document that cannot be uploaded.

#### Scenario: Convert a post with front matter

- **WHEN** the converter is given a post whose front matter contains a title
- **THEN** the emitted `content_state` contains no front matter text
- **AND** the parsed title is reported as document metadata alongside it

#### Scenario: Convert a post without a title

- **WHEN** the source front matter has no non-empty title
- **THEN** conversion fails with a diagnostic naming the source
- **AND** it does not emit a document that lacks a usable title

### Requirement: Deterministic offline conversion

The converter SHALL transform a Markdown source file into a X-compatible
DraftJS `content_state` document containing a `blocks` array and an `entities`
array. Conversion MUST be deterministic and MUST NOT perform network access.

#### Scenario: Convert a post offline

- **WHEN** the converter is given a Markdown file
- **THEN** it emits a `content_state` document with `blocks` and `entities`
- **AND** it performs no network request during conversion

#### Scenario: Identical input yields identical output

- **WHEN** the same Markdown input and configuration are converted twice
- **THEN** both outputs are byte-identical

### Requirement: Heading level mapping

Markdown headings SHALL map onto the block types X exposes, and no heading MAY
produce a type outside that set. Because X exposes only three heading levels,
depth SHALL be clamped: `#` and `##` map to `header-one`, `###` maps to
`header-two`, and `####` and deeper map to `header-three`.

#### Scenario: Map a second-level heading

- **WHEN** a post contains a `##` heading
- **THEN** it becomes a `header-one` block

#### Scenario: Map a third-level heading

- **WHEN** a post contains a `###` heading
- **THEN** it becomes a `header-two` block

#### Scenario: Clamp a heading deeper than X supports

- **WHEN** a post contains a `####`, `#####`, or `######` heading
- **THEN** it becomes a `header-three` block
- **AND** no block uses a heading type outside `header-one`, `header-two`, or `header-three`

### Requirement: Block-level Markdown mapping

The converter SHALL map supported block constructs to the block type X expects:
paragraphs to `unstyled`, unordered list items to `unordered-list-item`,
ordered list items to `ordered-list-item`, and block quotes to `blockquote`.
Nested list items SHALL be emitted as list items of their own kind rather than
being merged into the parent item.

#### Scenario: Convert paragraphs and lists

- **WHEN** a post contains paragraphs and both list kinds
- **THEN** paragraphs become `unstyled` blocks and each item becomes an
  `unordered-list-item` or `ordered-list-item` block

#### Scenario: Convert a nested list

- **WHEN** a list item contains a nested list
- **THEN** the nested items are emitted as list-item blocks
- **AND** no item's text contains the nested list's raw Markdown

#### Scenario: Convert a block quote

- **WHEN** a post contains a block quote
- **THEN** it becomes a `blockquote` block containing the quoted text

### Requirement: Inline styles and link entities

Inline emphasis SHALL be expressed as `inline_style_ranges` on the enclosing
block using `bold`, `italic`, or `strikethrough`, with offsets computed against
the block's final text. Links SHALL become `link` entities referenced by
`entity_ranges` on the enclosing block. Ranges MUST be correct when a styled or
linked span is adjacent to other text in the same block.

Offsets SHALL be counted in the code units the DraftJS format uses — UTF-16 code
units of the block's final text — rather than bytes, runes, or grapheme
clusters. A block containing a supplementary character before a styled or
linked span MUST still select exactly that span, because a byte or rune offset
would select the wrong characters once the document is interpreted as a
JavaScript string.

#### Scenario: Convert emphasis

- **WHEN** a paragraph contains bold, italic, or strikethrough text
- **THEN** the enclosing block carries `inline_style_ranges` for that span
- **AND** each range's offset and length select exactly the styled text

#### Scenario: Convert a link

- **WHEN** a paragraph contains a Markdown link
- **THEN** a `link` entity carrying the destination URL is added to `entities`
- **AND** the enclosing block references it with an `entity_range` selecting the link text

#### Scenario: Keep offsets correct for adjacent spans

- **WHEN** a block contains multiple styled or linked spans next to plain text
- **THEN** every range's offset is measured against the block's final plain text
- **AND** no range extends beyond the block text

#### Scenario: Keep offsets correct after a supplementary character

- **WHEN** a block contains an emoji or other supplementary character before a styled or linked span
- **THEN** each range's offset is measured in UTF-16 code units of the block's
  final text
- **AND** the range selects exactly the styled or linked text rather than
  drifting by the character's extra code unit

### Requirement: Code blocks and tables preserve source Markdown

Code fences and tables SHALL be emitted as `atomic` blocks backed by a
`markdown` entity whose payload preserves the construct's original Markdown
source. A table MUST NOT be flattened into paragraph text, because X exposes no
table block or entity.

#### Scenario: Convert a fenced code block

- **WHEN** a post contains a fenced code block
- **THEN** it becomes an `atomic` block with a `markdown` entity
- **AND** the entity payload preserves the fenced source, including its language

#### Scenario: Convert a table

- **WHEN** a post contains a Markdown table
- **THEN** it becomes an `atomic` block with a `markdown` entity
- **AND** the entity payload preserves the table as Markdown with its rows and
  cell contents intact

#### Scenario: Report the markdown payload budget

- **WHEN** the total `markdown` entity payload for an article would exceed the
  API's per-article limit
- **THEN** conversion fails with a diagnostic naming the article and the
  measured size rather than emitting a document the API will reject

### Requirement: Inline code and thematic breaks

DraftJS exposes only `bold`, `italic`, and `strikethrough` as inline styles, so
a Markdown inline-code span has no code style to map onto. Its literal text SHALL
survive conversion unchanged, and the loss of the monospace styling SHALL be
reported as a conversion diagnostic so the author sees it rather than having it
silently dropped. A Markdown thematic break SHALL become an `atomic` block
backed by a `divider` entity, because X exposes `divider` for exactly that
purpose, rather than being dropped as an unrecognized paragraph.

#### Scenario: Convert an inline-code span

- **WHEN** a post contains a Markdown inline-code span
- **THEN** the block's final text contains the span's literal text
- **AND** conversion reports the lost code styling as a diagnostic

#### Scenario: Convert a thematic break

- **WHEN** a post contains a thematic break
- **THEN** it becomes an `atomic` block with a `divider` entity
- **AND** it does not become paragraph text containing the break's characters

### Requirement: Footnotes become a trailing section

DraftJS has no footnote construct, so footnote definitions SHALL be preserved in
a trailing `Footnotes` section rather than dropped. Each in-text footnote
reference SHALL become plain bracketed text, and each definition SHALL become an
`ordered-list-item` beneath a `Footnotes` heading.

#### Scenario: Convert a post with footnotes

- **WHEN** a post contains a footnote reference and its definition
- **THEN** the reference appears in the body as plain bracketed text
- **AND** a `Footnotes` heading is appended with each definition as an
  ordered list item

#### Scenario: Convert a post without footnotes

- **WHEN** a post contains no footnotes
- **THEN** no `Footnotes` heading is added

### Requirement: Images reference uploaded media

The `image` entity SHALL carry its media as `media_items`, whose entries require
a `media_category` and a `media_id`; `data.url` is documented for `link`
entities and MUST NOT be treated as an image source. Each image SHALL become an
`atomic` block plus an `image` entity, and unresolved images SHALL remain
explicitly unresolved so publication can resolve or reject them rather than
silently emitting a broken article. Because the media endpoints accept only a
fixed set of image media types that excludes `image/svg+xml`, the converter
SHALL report an image whose media type the upload endpoints cannot accept, with
its source position, rather than emitting an entity that publication can never
resolve.

#### Scenario: Convert an inline image

- **WHEN** a post contains an image with alt text
- **THEN** it becomes an `atomic` block with an `image` entity
- **AND** the entity preserves the alt text as its caption

#### Scenario: Publication has not yet resolved an image

- **WHEN** a converted document references an image whose media has not been uploaded
- **THEN** the unresolved image is distinguishable from a resolved one
- **AND** publication refuses to send an unresolved image as though it were resolved

#### Scenario: An image cannot be uploaded

- **WHEN** a post references an image whose media type the media upload endpoints do not accept
- **THEN** the converter reports the image with its source position
- **AND** it does not emit an entity that publication could never resolve

#### Scenario: A post follows the raster-only policy

- **WHEN** a post's images all use a media type the media upload endpoints accept
- **THEN** conversion completes with an `image` entity for each image
- **AND** no unacceptable-image diagnostic is reported for that post

### Requirement: Artifact carries image locators

Because the API's entity `data` object rejects additional properties, an image's
source cannot travel inside the entity. The emitted draft artifact SHALL
therefore carry, as metadata beside `content_state`, an entry for every
unresolved image. Each entry SHALL name the image's source path relative to the
post's directory, SHALL record the post package the path is relative to, and
SHALL record a digest of the image bytes as they were at conversion time, so the
artifact identifies both the image and the directory it resolves against.
Publication SHALL be able to obtain the image bytes from the artifact alone,
without re-parsing Markdown, and SHALL reject a locator that escapes the post's
directory or whose current bytes do not match the recorded digest.

#### Scenario: Convert an image into the artifact

- **WHEN** a post containing a raster image is converted
- **THEN** the artifact carries a locator entry for that image beside the document
- **AND** the entry names the source path relative to the post's directory, the
  post package the path resolves against, and a digest of the image bytes

#### Scenario: Associate a locator with its image entity

- **WHEN** an artifact carries two or more image locators
- **THEN** each locator entry records the entity it resolves, by the entity's
  key, so the mapping does not depend on ordering or on repeated captions
- **AND** publication attaches each returned `media_id` to the entity its locator
  names rather than by position

#### Scenario: Resolve an image from the artifact

- **WHEN** publication reads an unresolved image from the artifact
- **THEN** it obtains the image bytes from the recorded locator without
  re-parsing the source Markdown
- **AND** it resolves the relative path against the recorded post package even
  when the artifact was moved elsewhere
- **AND** it refuses a locator that escapes the post's directory or names bytes
  whose digest does not match the recorded one

### Requirement: Conversion diagnostics

The converter SHALL report a construct it cannot represent with the source
location that produced it. It SHALL distinguish two outcomes and never emit a
document that silently omits content it could not represent: a construct whose
content is preserved but whose formatting is lost is reported and conversion
continues, while a construct whose content would be lost fails conversion. The
image whose media type the upload endpoints reject and the source without a
title are failing outcomes; an inline-code span's lost styling is a continuing
one.

#### Scenario: Encounter an unsupported construct

- **WHEN** the converter cannot represent a construct and its content would be lost
- **THEN** it reports the construct with its source position
- **AND** conversion fails rather than emitting a document that discards the content

#### Scenario: Encounter a construct whose formatting is lost

- **WHEN** the converter preserves a construct's content but cannot reproduce its formatting
- **THEN** it reports the lost formatting with its source position
- **AND** conversion still completes with the content present
