# x-article-publication Specification

## Purpose

Create an X Article draft from a converted Markdown post, resolving referenced
images through the media upload endpoint. The uploader creates drafts only; it
never publishes an article.

## Requirements

### Requirement: Authenticated requests to X Articles

The publisher SHALL authenticate to the X API with credentials supplied at run
time through the repository's secret-injection flow. Credential values MUST NOT
be recorded in source, committed configuration, logs, or command history.

Authentication SHALL use OAuth 1.0a user context, because the existing media and
Articles endpoints accept it and its user-context token does not expire on a
fixed schedule the way an OAuth 2.0 user token does, so draft creation does not
depend on a rotating refresh token. OAuth 1.0a needs four credential fields —
the API key, the API secret, the access token, and the access token secret —
and the publisher SHALL sign each request with HMAC-SHA1 per RFC 5849 using the
repository's standard library rather than adding a signing dependency. The
publisher SHALL fail when any of the four fields is absent, rather than sending
a partially signed request.

#### Scenario: Create a draft with injected credentials

- **WHEN** the publisher runs with credentials supplied through the injection flow
- **THEN** it signs its requests to the X API with OAuth 1.0a user context
- **AND** it does not print the credential values

#### Scenario: Credentials are unavailable

- **WHEN** the publisher runs without any one of the four required credential fields
- **THEN** it fails with a diagnostic naming the missing credential reference
- **AND** it does not retry with a partial or empty credential

### Requirement: Draft-only uploads

The uploader SHALL create X Article drafts and MUST NOT publish an article. It
SHALL NOT call the Articles publish endpoint or expose a way to request
publication, so no invocation of this tool can make a post publicly visible.
Draft creation SHALL send the title that conversion parsed from the source front
matter, and MUST NOT invent, default, or omit it, because the draft endpoint
requires a title that `content_state` does not carry.

#### Scenario: Create a draft

- **WHEN** draft creation is requested for a converted document
- **THEN** the publisher reads the draft article identifier from the v2 response envelope's `data.id`
- **AND** no post becomes publicly visible
- **AND** the request carries the title parsed from the source front matter

#### Scenario: A draft response omits the identifier

- **WHEN** the draft endpoint returns a 2xx response whose `data.id` is missing or empty
- **THEN** the publisher fails rather than reporting a successful draft with no identifier to review
- **AND** it does not treat the zero-value response as a created draft

#### Scenario: Draft creation is requested without a parsed title

- **WHEN** draft creation is requested for a document with no parsed title
- **THEN** the publisher fails with a diagnostic
- **AND** it does not send a placeholder or empty title

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
the digest the locator records, rather than uploading a different image. It
SHALL enforce that containment by resolving the name through a
descriptor-anchored directory handle that rejects a symlink escaping the root,
rather than by checking a pathname and then re-opening it, so the object whose
containment is checked is the object read.

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

- **WHEN** draft creation runs again in a new process for an image this tool already uploaded
- **THEN** it reads the digest-to-`media_id` mapping from its own cache and
  reuses the identifier
- **AND** it stays correct when the cache is absent or the recorded digest no
  longer matches the artifact's, by uploading again rather than reusing a stale
  identifier

#### Scenario: Decode the media upload response

- **WHEN** the media upload endpoint returns the v2 response envelope
- **THEN** the publisher reads the media identifier from the envelope's `data.id`
- **AND** the created draft references that identifier rather than an empty one

#### Scenario: A media upload response omits the identifier

- **WHEN** the media upload endpoint returns a 2xx response whose `data.id` is missing or empty
- **THEN** the publisher fails rather than resolving the image to an empty identifier
- **AND** it does not send an entity that references a missing media identifier

#### Scenario: Expire a cached upload

- **WHEN** the recorded lifetime of a cached media identifier has passed
- **THEN** the cached entry is treated as a miss and the image is uploaded again
- **AND** an entry whose lifetime has not passed is still reused

#### Scenario: An image cannot be uploaded

- **WHEN** an image upload fails
- **THEN** the publisher fails and reports the failing image
- **AND** it does not create a draft that references a missing or placeholder media identifier

#### Scenario: An artifact names an unusable image locator

- **WHEN** an artifact's locator escapes the recorded post package or names bytes whose digest does not match the recorded one
- **THEN** the publisher fails with a diagnostic naming the locator
- **AND** it uploads no replacement image in its place

#### Scenario: An unresolved image entity has no locator

- **WHEN** an artifact retains an `image` entity without `media_items` but records no locator that resolves it
- **THEN** the publisher fails before uploading any image or creating a draft, naming the unresolved entity
- **AND** it does not send an image entity without a media item

#### Scenario: A locator does not name a distinct image entity

- **WHEN** an artifact's locator names an entity that is missing, is not an image, or is already named by another locator
- **THEN** the publisher fails before uploading any image rather than discovering the mismatch after paying for uploads

#### Scenario: Two entities share one numeric key

- **WHEN** an artifact contains more than one entity with the same numeric key
- **THEN** the publisher fails before uploading any image, because one locator could resolve only the first and leave the other without `media_items`

#### Scenario: An image entity carries a malformed media_items entry

- **WHEN** an image entity's `media_items` is present but its entry has an empty or absent `media_id`, or a missing or incorrect `media_category`
- **THEN** the entity is treated as unresolved, so the publisher requires a locator rather than sending the invalid entry to X

#### Scenario: An image entity carries a usable entry beside a malformed one

- **WHEN** an image entity's `media_items` holds a usable entry next to an empty, non-object, or incorrectly categorized one
- **THEN** the entity is treated as unresolved, because every supplied entry must be usable
- **AND** the publisher does not send the malformed array merely because one entry looked valid

#### Scenario: An image entity carries no usable key

- **WHEN** an artifact contains an image entity whose `key` is missing, non-numeric, fractional, or negative
- **THEN** the publisher fails before uploading any image, because the entity is invisible to locator coverage
- **AND** it does not send an unresolved image to X

#### Scenario: The cache cannot be written

- **WHEN** the cache cannot be recorded, such as a cache path in a read-only directory or over an existing directory
- **THEN** draft creation continues and the draft is created from the uploads that already succeeded
- **AND** the failure costs a later re-upload rather than aborting draft creation

### Requirement: Reporting and failure handling

The publisher SHALL report the identifier it obtained and SHALL fail loudly
rather than report success when a request was rejected. Failures MUST identify
which operation failed and MUST NOT be retried automatically when the response
indicates an authentication or request error.

#### Scenario: Report the created identifier

- **WHEN** a draft is created
- **THEN** the publisher reports the draft's article identifier

#### Scenario: The API rejects a request

- **WHEN** an API request is rejected
- **THEN** the publisher reports the failing operation and the API's error
- **AND** it does not report the operation as successful

#### Scenario: Validation runs without network access

- **WHEN** the project's checks validate conversion
- **THEN** they run without contacting the X API
- **AND** no test requires live credentials

### Requirement: Article banner media

Draft creation SHALL upload a selected banner through the media API and send
its identifier in the optional `cover_media` object with
`media_category: tweet_image`. The publisher SHALL share its digest-based media
cache between banner and body images, reusing an unexpired identifier for the
same bytes. It SHALL validate all banner and body image locators, bytes,
digests, and media types before any upload or draft request.

#### Scenario: Reuse one upload for banner and body

- **WHEN** the banner and a body image contain the same bytes
- **THEN** at most one media upload resolves both references
- **AND** the draft request carries that ID in `cover_media` and the body entity

#### Scenario: Upload distinct banner and body images

- **WHEN** the selected banner has different bytes from every body image
- **THEN** it receives its own media upload or cached identifier
- **AND** all media resolves before the draft request is sent

#### Scenario: Create a draft without a banner

- **WHEN** the artifact has no banner locator
- **THEN** the request omits `cover_media`

#### Scenario: Invalid image prevents every request

- **WHEN** any banner or body locator is malformed, escapes its package, names
  missing bytes, or has a mismatched digest or media type
- **THEN** draft creation fails before any media or draft HTTP request
