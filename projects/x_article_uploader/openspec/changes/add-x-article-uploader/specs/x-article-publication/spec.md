## Purpose

Publish a converted Markdown post to X as an Article, keeping draft creation
separate from publication and resolving referenced images through the media
upload endpoint.

## ADDED Requirements

### Requirement: Authenticated requests to X Articles

The publisher SHALL authenticate to the X API with credentials supplied at run
time through the repository's secret-injection flow. Credential values MUST NOT
be recorded in source, committed configuration, logs, or command history.

#### Scenario: Publish with injected credentials

- **WHEN** the publisher runs with credentials supplied through the injection flow
- **THEN** it authenticates its requests to the X API
- **AND** it does not print the credential values

#### Scenario: Credentials are unavailable

- **WHEN** the publisher runs without usable credentials
- **THEN** it fails with a diagnostic naming the missing credential reference
- **AND** it does not retry with a partial or empty credential

### Requirement: Draft creation is separate from publication

Creating a draft SHALL be a distinct operation from publishing it. Translating a
document into a draft MUST NOT publish a post, and publishing SHALL require an
explicitly requested operation so a draft can be reviewed first. Draft creation
SHALL send the title that conversion parsed from the source front matter, and
MUST NOT invent, default, or omit it, because the draft endpoint requires a
title that `content_state` does not carry.

#### Scenario: Create a draft only

- **WHEN** draft creation is requested for a converted document
- **THEN** the API returns a draft article identifier
- **AND** no post becomes publicly visible
- **AND** the request carries the title parsed from the source front matter

#### Scenario: Draft creation is requested without a parsed title

- **WHEN** draft creation is requested for a document with no parsed title
- **THEN** the publisher fails with a diagnostic
- **AND** it does not send a placeholder or empty title

#### Scenario: Publish a reviewed draft

- **WHEN** publication is explicitly requested for an existing draft identifier
- **THEN** the corresponding post is created and the resulting post identifier is reported

### Requirement: Image media upload

The publisher SHALL upload each referenced image through the X media upload
endpoint and reference the returned `media_id` in an `image` entity, because
that entity carries no document URL. Uploading MUST happen before the draft
that references the media is created.

The publisher SHALL obtain the image bytes from the locator the draft artifact
records beside the document, not by re-parsing the source Markdown, because the
API's entity `data` object rejects additional properties and so cannot carry a
path. It SHALL resolve the recorded relative path against the post package the
locator records, so the artifact stays resolvable after it is moved. It SHALL
reject a locator that escapes the post's directory or whose bytes do not match
the digest the locator records, rather than uploading a different image.

#### Scenario: Upload an image and resolve its media identifier

- **WHEN** a converted document references an unuploaded image
- **THEN** the publisher reads the image bytes from the artifact's locator and
  obtains a `media_id`
- **AND** the created draft references that `media_id` in its `image` entity

#### Scenario: Reuse an already uploaded image

- **WHEN** the same image content has already been uploaded by this tool
- **THEN** the publisher reuses the recorded `media_id` rather than uploading the
  same bytes again

#### Scenario: Reuse an upload in a later run

- **WHEN** publication runs again in a new process for an image this tool already uploaded
- **THEN** it reads the digest-to-`media_id` mapping from its own cache and
  reuses the identifier
- **AND** it stays correct when the cache is absent or the recorded digest no
  longer matches the artifact's, by uploading again rather than reusing a stale
  identifier

#### Scenario: An image cannot be uploaded

- **WHEN** an image upload fails
- **THEN** the publisher fails and reports the failing image
- **AND** it does not create a draft that references a missing or placeholder media identifier

#### Scenario: An artifact names an unusable image locator

- **WHEN** an artifact's locator escapes the recorded post package or names bytes whose digest does not match the recorded one
- **THEN** the publisher fails with a diagnostic naming the locator
- **AND** it uploads no replacement image in its place

### Requirement: Reporting and failure handling

The publisher SHALL report the identifiers it obtained and SHALL fail loudly
rather than report success when a request was rejected. Failures MUST identify
which operation failed and MUST NOT be retried automatically when the response
indicates an authentication or request error.

#### Scenario: Report created identifiers

- **WHEN** a draft is created or an article is published
- **THEN** the publisher reports the corresponding article or post identifier

#### Scenario: The API rejects a request

- **WHEN** an API request is rejected
- **THEN** the publisher reports the failing operation and the API's error
- **AND** it does not report the operation as successful

#### Scenario: Validation runs without network access

- **WHEN** the project's checks validate conversion
- **THEN** they run without contacting the X API
- **AND** no test requires live credentials
