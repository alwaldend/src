## MODIFIED Requirements

### Requirement: Code blocks and tables preserve source Markdown

Code fences and tables SHALL be emitted as `atomic` blocks backed by a
`markdown` entity with `mutable` mutability whose payload preserves the construct's original Markdown
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

- **WHEN** the documented conservative weighted-length estimate for Markdown
  entities exceeds the local per-article budget
- **THEN** conversion fails with the measured estimate and budget
- **AND** the diagnostic and documentation distinguish that estimate from X's
  undocumented Articles-specific backend validation

## ADDED Requirements

### Requirement: Validate the supported static image contract

Conversion and draft preflight SHALL accept only static JPEG, PNG, GIF, and
WebP images up to 5,000,000 bytes. They SHALL reject unsupported formats,
animations, malformed images detectable by the supported decoder or container
validator, and images exceeding the documented local decoded-size guard. The
implementation SHALL distinguish X's upload limits from local resource guards
and SHALL NOT claim full WebP decoding when only container/header validation is
available. Every image SHALL pass preflight before any upload or draft request.

#### Scenario: An unsupported or animated image is referenced

- **WHEN** a banner or retained body image uses an unsupported format or animation
- **THEN** conversion reports the unusable source
- **AND** a supplied artifact containing that source is also refused before network requests

#### Scenario: A file exceeds the simple-image upload limit

- **WHEN** a selected image contains more than 5,000,000 bytes
- **THEN** conversion and draft preflight refuse it before upload

#### Scenario: A malformed image only carries a recognized signature

- **WHEN** a file's signature matches its extension but decoding or container validation fails
- **THEN** it is refused rather than sent as a usable image
