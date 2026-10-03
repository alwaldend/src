## MODIFIED Requirements

### Requirement: Conversion diagnostics

The converter SHALL report a construct it cannot represent with the source
location that produced it. It SHALL distinguish two outcomes and never emit a
document that silently omits content it could not represent: a construct whose
content is preserved but whose formatting is lost is reported as a warning.
The conversion command SHALL fail on warnings by default and SHALL permit them
only with `--warnings-as-errors=false`. Errors SHALL fail with either setting.
All diagnostics SHALL be printed to stderr and SHALL NOT appear in artifact JSON. The
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
- **AND** command conversion fails by default without writing an artifact
- **AND** `--warnings-as-errors=false` permits completion with the content present
  while still reporting the warning to stderr

#### Scenario: Report a position after a multi-byte character

- **WHEN** a diagnostic's construct follows a multi-byte character on the same
  line
- **THEN** the reported column counts source characters rather than bytes, so it
  names the position an editor shows
- **AND** the byte offset is still recorded for slicing the source

## ADDED Requirements

### Requirement: Artifact separates the X payload from local metadata

The converted artifact SHALL put title and content_state inside payload. Body
image locators and the optional banner locator SHALL remain outside payload.
Diagnostics SHALL NOT be serialized. Media resolution SHALL attach body media
IDs inside payload.content_state and the banner ID inside payload.cover_media.
Artifact reading SHALL reject the old flat format and require regeneration.

#### Scenario: Build an offline artifact

- **WHEN** conversion succeeds
- **THEN** payload contains title and content_state in the X request shape
- **AND** local media references remain alongside payload
- **AND** there is no serialized diagnostics field

#### Scenario: Resolve media for draft submission

- **WHEN** the uploader resolves an artifact's media
- **THEN** body entities inside payload receive media_items
- **AND** the banner becomes payload.cover_media
- **AND** the draft request contains only X payload fields

#### Scenario: Read a legacy flat artifact

- **WHEN** an artifact has title and content_state outside payload
- **THEN** reading fails before any API request
- **AND** the error instructs the caller to regenerate the artifact
