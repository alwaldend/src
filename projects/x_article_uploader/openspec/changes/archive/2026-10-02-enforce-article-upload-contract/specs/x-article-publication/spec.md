## MODIFIED Requirements

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

- **WHEN** the same image content has a ready upload in the current credential context and a usable lifetime
- **THEN** the publisher reuses the recorded `media_id` rather than uploading the
  same bytes again

#### Scenario: Reuse an upload in a later run

- **WHEN** draft creation runs again in a new process for an image this tool already uploaded
- **THEN** it reuses the identifier only when the cache schema, endpoint, media category,
  credential context, content digest, and conservative expiry all match
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

- **WHEN** the recorded lifetime of a cached media identifier has passed before image resolution begins
- **THEN** the cached entry is treated as a miss and the image is uploaded again
- **AND** an entry with enough lifetime for draft submission is still reused

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

## ADDED Requirements

### Requirement: Media readiness before draft creation

The publisher SHALL inspect media processing information before using an upload
identifier. A response without processing information or with succeeded state
SHALL require no status request. Pending or in-progress processing SHALL use
read-only signed status requests, respect the reported delay, and stop within a
bounded deadline. Failed, unknown, inconsistent, or incomplete processing state,
partial upload errors, and rejected status requests SHALL prevent draft creation.
Neither upload nor draft POST requests SHALL be replayed by this workflow.

#### Scenario: A static upload is immediately usable

- **WHEN** upload succeeds without processing information or with succeeded state
- **THEN** the publisher makes no processing-status request
- **AND** the media identifier can be attached to the draft

#### Scenario: Media processing remains pending

- **WHEN** the upload reports pending or in-progress processing
- **THEN** the publisher checks status after the prescribed delay within a finite budget
- **AND** it creates a draft only after observing succeeded state

#### Scenario: Processing cannot establish usable media

- **WHEN** processing fails, exceeds its deadline, returns errors, or supplies inconsistent state
- **THEN** the publisher reports the media failure without creating a draft
- **AND** it does not retry the upload POST

### Requirement: Credential-scoped media cache validity

Persistent media reuse SHALL be scoped to the upload endpoint, media category,
and OAuth credential context without persisting plaintext credentials or
requesting account identity merely to name the cache. Changing any credential
field SHALL prevent reuse from the prior context. Unscoped legacy entries SHALL
be treated as misses.

A persisted entry SHALL have a known positive lifetime, conservatively anchored
before upload and accounting for processing time and a submission reserve.
Missing or zero lifetime SHALL permit reuse only within the current resolution,
never across resolutions or processes. Negative or already unusable upload
lifetimes SHALL fail before draft creation.

#### Scenario: A token or account changes

- **WHEN** identical image bytes are submitted using a different credential context
- **THEN** the earlier context's cached identifier is not reused
- **AND** the cache contains no credential values

#### Scenario: A response provides no known lifetime

- **WHEN** an upload response omits its lifetime or reports zero
- **THEN** identical banner and body images may share that fresh identifier within one resolution
- **AND** a later resolution uploads the bytes again

#### Scenario: Uploading or processing consumes the usable lifetime

- **WHEN** a positive reported lifetime has insufficient time remaining for submission
- **THEN** draft creation stops instead of attaching the expired or near-expiry identifier

#### Scenario: An earlier image expires before its later duplicate is resolved

- **WHEN** uploading another image consumes the usable lifetime of an identifier already selected for a banner or body image
- **AND** a later body image references those same bytes
- **THEN** resolution stops without uploading the duplicate again or creating a draft
- **AND** a replacement cache entry cannot mask the expired identifier already selected for the draft
