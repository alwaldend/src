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
array of `{key, value}` entries. Conversion MUST be deterministic and MUST NOT perform network access.

#### Scenario: Convert a post offline

- **WHEN** the converter is given a Markdown file
- **THEN** it emits a `content_state` document with `blocks` and `entities`
- **AND** it performs no network request during conversion

#### Scenario: Identical input yields identical output

- **WHEN** the same Markdown input and configuration are converted twice
- **THEN** both outputs are byte-identical

### Requirement: Articles wire shape

The emitted document SHALL use the exact shape the Articles draft endpoint
accepts, not the canonical DraftJS spelling. Each block SHALL name its ranges
`inline_style_ranges` and `entity_ranges`, and MUST NOT carry a `depth` field,
because the endpoint's schema sets `additionalProperties: false` and rejects
both the camelCase range names and any `depth`. Each entry of `entities` SHALL
be `{"key": "<index>", "value": {"type", "mutability", "data"}}`, matching the
endpoint's schema rather than DraftJS's inline `entityMap`.

#### Scenario: Emit the Articles field names

- **WHEN** a document with inline styles, links, and a list is converted
- **THEN** every block names its ranges `inline_style_ranges` and `entity_ranges`
  and carries no `depth`
- **AND** each entity is a `{key, value}` entry whose `value` carries `type`,
  `mutability`, and `data`
- **AND** no block or entity field falls outside the schema's allowed names

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

#### Scenario: Convert a list item with several paragraphs

- **WHEN** a list item contains more than one paragraph, as a loose list does,
  whether or not a block construct such as a nested list or a code block lies
  between them
- **THEN** the item becomes a single list-item block whose text carries all of
  the item's paragraphs
- **AND** the item's block constructs follow that block in source order
- **AND** it does not become one block per paragraph, which DraftJS would read
  as separate items and which would renumber an ordered list

#### Scenario: Convert a block construct nested in a list item

- **WHEN** a list item contains a block construct such as a fenced code block,
  an indented code block, a table, or a thematic break
- **THEN** the construct is emitted through the same block mapping it uses at
  the top level, with its source preserved
- **AND** its source is not silently dropped by walking only inline children

#### Scenario: Convert a block quote

- **WHEN** a post contains a block quote
- **THEN** it becomes a `blockquote` block containing the quoted text

#### Scenario: Convert a heading inside a block quote

- **WHEN** a block quote contains a heading, as in `> # Warning`, whether
  directly or through a nested list as in `> - # Warning`
- **THEN** the heading text is emitted as a `blockquote` block, not an ordinary
  heading outside the quote and not a list-item block that drops the quote
- **AND** the lost heading level is reported, because X exposes no quoted-heading
  block type

#### Scenario: Preserve a line break

- **WHEN** a paragraph contains a soft line break or a hard line break written
  with two trailing spaces or a trailing backslash
- **THEN** the block's text contains a newline at that break
- **AND** the two lines are not joined into one word

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

#### Scenario: Convert a link carrying a title

- **WHEN** a Markdown link carries an optional title, as in
  `[docs](https://example.test "Reference")`
- **THEN** the entity cannot carry the title, so its loss is reported as a
  continuing diagnostic naming the title and positioned at the link
- **AND** the link's text and destination are still preserved

#### Scenario: Convert an image carrying a title

- **WHEN** a Markdown image carries an optional title
- **THEN** its loss is reported as a continuing diagnostic naming the title and
  positioned at the image
- **AND** the image is still emitted as an `atomic` block

#### Scenario: A titled construct has no text of its own

- **WHEN** a titled link or image has an empty label or alt text, so the
  construct has no text node to locate it
- **THEN** the lost-title diagnostic is positioned at the construct's own
  opening delimiter rather than at the start of the enclosing block
- **AND** the reported source position still names the construct that produced
  the loss

#### Scenario: Convert an email autolink

- **WHEN** a paragraph contains an email autolink such as `<user@example.com>`
- **THEN** the `link` entity records it as `mailto:user@example.com`
- **AND** the recorded URL is an email address rather than a relative link

#### Scenario: A destination carries Markdown escapes or character references

- **WHEN** a link or image destination escapes a punctuation character or uses a
  character reference, as in `[x](https://example.test/a\(1\)?x=1&amp;y=2)`
- **THEN** the recorded destination resolves those encodings to the URL the
  author wrote
- **AND** the recorded link or image target does not contain the literal
  backslashes or unresolved entity

#### Scenario: Prose carries Markdown escapes or character references

- **WHEN** ordinary prose escapes a punctuation character or uses a character
  reference, as in `\*literal\*` or `AT&amp;T`
- **THEN** the emitted block text resolves those encodings to the characters the
  author wrote rather than keeping the backslashes or the raw entity
- **AND** a raw span such as a code span keeps its literal source, because its
  encoding is its content
- **AND** an image caption, which is flattened from its inline content, follows
  the same rule: resolution applies to each non-raw descendant while a raw span
  inside the caption keeps its literal source

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

#### Scenario: Convert an indented code block

- **WHEN** a post contains an indented code block immediately followed by a
  non-indented paragraph
- **THEN** the `markdown` entity payload contains only the code block's own lines
- **AND** the following paragraph is emitted once, as its own block, rather than
  being swallowed into the entity and repeated

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

#### Scenario: Convert a GFM task-list item

- **WHEN** a post contains a task-list item
- **THEN** the item's text keeps an indicator of its checkbox state
- **AND** the checked state is not silently dropped

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

#### Scenario: Preserve a structured footnote definition

- **WHEN** a footnote definition contains multiple paragraphs or inline
  formatting such as a link
- **THEN** each definition paragraph becomes its own ordered list item rather
  than being joined without a separator
- **AND** inline entities such as links are preserved by the normal inline
  mapper rather than flattened to plain text

#### Scenario: Preserve a construct nested in a footnote definition

- **WHEN** a footnote definition contains a nested list, block quote, fenced or
  indented code block, table, or thematic break
- **THEN** the construct is emitted through the block mapper with its content
  preserved
- **AND** it is not silently dropped

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
resolve. The declared extension SHALL be accepted as the media type only when
the file's bytes carry that format's own signature, so a file renamed to an
accepted raster extension — even one accepted format renamed to another — is
reported rather than sent to the upload endpoint.

#### Scenario: Convert an inline image

- **WHEN** a post contains an image with alt text
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
directory or whose current bytes do not match the recorded digest. Conversion
SHALL apply the same containment before it reads an image, so it neither
inspects bytes outside the post boundary nor records a locator that publication
would refuse. Both SHALL enforce containment by resolving the name through
descriptor-anchored directory handles that reject a symlink escaping the root,
rather than by checking a pathname and then re-opening it, so the object whose
containment is checked is the object read.

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

#### Scenario: Convert a source given by an absolute path

- **WHEN** the conversion command is given an absolute `--source` and no
  `--post-package`
- **THEN** the recorded locator's package is relative to the workspace named by
  `--workspace`, rather than the absolute source directory
- **AND** the artifact remains publishable, because publication refuses an
  absolute package rather than resolving it
- **AND** the workspace root is the caller-supplied `--workspace` value and is
  never re-discovered by walking the filesystem, so no repository marker found
  through a link can be mistaken for it

#### Scenario: A post directory is reached through a symlink

- **WHEN** the derived package would name a post directory whose canonical
  target, after resolving symlinks, lies outside the workspace root
- **THEN** the conversion is refused with a directive to pass `--post-package`
  explicitly, rather than recording a package the descriptor-anchored read
  would resolve past and reject
- **AND** an in-workspace symlink still records the relative package of the
  canonical directory it resolves to, so the recorded package matches the
  object publication reads
- **AND** a link that points directly at another checkout's root, a link into
  another tree beneath its root, and a post directory reached through an
  intermediate symlink-component are all refused rather than recorded relative
  to the external tree or collapsed to `.`
- **AND** containment is established by opening the directory through a
  descriptor-anchored handle on the workspace, so the no-symlink validation and
  the directory accepted for the package are the same traversal and cannot be
  separated by a writer swapping the directory for a link
- **AND** image reads for the post go through that same handle rather than
  re-opening the source pathname, so a symlink retargeted after validation
  cannot make conversion read another tree's bytes while the artifact still
  records the validated package

#### Scenario: A post package is supplied explicitly

- **WHEN** the conversion command is given `--post-package` rather than having
  one derived
- **THEN** the package is accepted only when the directory it opens to is the
  post directory, compared by file identity rather than by pathname, because
  publication resolves the locator paths against the package and a package
  naming another directory would make publication read a different object than
  conversion did
- **AND** the post's images are read through the descriptor-anchored handle the
  identity was checked on, so a package retargeted between the check and the
  read cannot supply another tree's bytes while the locator still names the
  supplied package
- **AND** an absolute package, a package that escapes the workspace, or a
  package naming a different directory than the post is refused

#### Scenario: Resolve an image from the artifact

- **WHEN** publication reads an unresolved image from the artifact
- **THEN** it obtains the image bytes from the recorded locator without
  re-parsing the source Markdown
- **AND** it resolves the relative path against the recorded post package even
  when the artifact was moved elsewhere
- **AND** it refuses a locator that escapes the post's directory or names bytes
  whose digest does not match the recorded one

#### Scenario: A symlink changes after containment validation

- **WHEN** a locator's path or its target is swapped for a symlink that leaves the
  post directory after containment was established
- **THEN** the read is refused rather than following the swapped symlink
- **AND** a swapped symlink cannot redirect the read to a different file, because
  the containment check and the read operate on the same descriptor rather than
  on a re-opened pathname

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

#### Scenario: Report continuing losses when conversion also fails

- **WHEN** a source contains both a continuing loss and a failing construct, such
  as an inline-code span and an unacceptable image
- **THEN** the continuing diagnostic is reported alongside the failing one
  rather than being discarded with the returned error
- **AND** the author sees every detected loss in one run

#### Scenario: Encounter a construct whose formatting is lost

- **WHEN** the converter preserves a construct's content but cannot reproduce its formatting
- **THEN** it reports the lost formatting with its source position
- **AND** conversion still completes with the content present

#### Scenario: Report a position after a multi-byte character

- **WHEN** a diagnostic's construct follows a multi-byte character on the same
  line
- **THEN** the reported column counts source characters rather than bytes, so it
  names the position an editor shows
- **AND** the byte offset is still recorded for slicing the source

### Requirement: Banner image metadata

The converter SHALL resolve the first entry of the optional front matter
`images` array as the article banner. It SHALL record a separate banner locator
with the post-relative path, post package, content digest, and detected media
type, using the same image validation as body images. It MUST NOT insert a body
block or remove an existing body image because the image is selected as a banner.
Malformed metadata or an unusable selected image SHALL fail conversion.

#### Scenario: Select a banner independently of body images

- **WHEN** a post selects a local image in the first `images` entry
- **THEN** conversion records its banner locator separately from body locators
- **AND** later `images` entries do not select additional banners
- **AND** body text and images retain their original order

#### Scenario: Convert without a selected banner

- **WHEN** `images` is absent or empty
- **THEN** the artifact has no banner locator and preserves existing conversion behavior

#### Scenario: Refuse an invalid selected banner

- **WHEN** `images` is not an array of strings, or its selected image is missing,
  unsupported, remote, or outside the post package
- **THEN** conversion fails with a diagnostic before any network operation
