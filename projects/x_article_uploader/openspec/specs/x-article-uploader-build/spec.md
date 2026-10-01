# x-article-uploader-build Specification

## Purpose

Build the X Article uploader from pinned, hermetic inputs so conversion and its
tests run reproducibly in the repository's Bazel workflow without network access
or credentials.

## Requirements

### Requirement: Bazel targets for the uploader

The project SHALL expose Bazel targets for its converter library, its
conversion command, and its draft-creation command, with tests declared beside
the sources they cover. Targets MUST follow the repository's role-based source
layout using `internal/` for implementation and `cmd/` for entry points.

#### Scenario: Build the converter and draft command

- **WHEN** the project's conversion and draft-creation targets are built
- **THEN** both binaries build from the project's declared sources

#### Scenario: Test the converter

- **WHEN** the project's library tests run
- **THEN** they execute against the converter's declared source inputs
- **AND** they require neither credentials nor network access

### Requirement: Pinned hermetic dependencies

External build inputs, including the Markdown parser, SHALL be declared with
immutable versions and integrity information in the owning dependency
declarations, and SHALL be updated through the owning generator rather than
hand-edited.

#### Scenario: Add a Markdown parser dependency

- **WHEN** the project depends on an external Markdown parser
- **THEN** its version and integrity are pinned in the owning dependency declaration
- **AND** the applicable lock and package-check targets validate the declaration

#### Scenario: Build without network access

- **WHEN** the project is built and tested in the repository's agent configuration
- **THEN** no build or test step fetches an undeclared input

### Requirement: Behavioral tests over real posts

Tests SHALL exercise conversion against every blog post in the site content
tree and SHALL assert the mapped block types, inline ranges, and entity payloads
for the constructs those posts contain, rather than relying on a single smoke
check. The covered set MUST be derived from the content tree, so a post added
later is covered without editing the test list. Each post SHALL carry an
expected outcome: a post whose images all use a media type the upload endpoints
accept MUST satisfy the mapping contract and MAY carry the continuing
diagnostics its constructs produce, while a post that references an unacceptable
image MUST fail with the named diagnostic at the reported source position. Where
a construct is not present in any current post, a focused case SHALL cover it.

Behavior coverage SHALL be authored before the implementation it exercises: the
expected outcome for a construct, or the enumerated ways a conversion can fail,
is written first, and the implementation follows. Test cases MUST assert
observable conversion behavior and MUST NOT restate the implementation's own
constants or assert only that output changed.

#### Scenario: Convert every existing blog post

- **WHEN** the converter's tests run against the blog posts in the site content tree
- **THEN** conversion is exercised for every post in the tree, discovered from the tree itself
- **AND** each post produces the outcome its entry records: the mapping
  contract where every image is acceptable, including any continuing diagnostic
  its constructs produce, or the named failing diagnostic where an image is not
  acceptable

#### Scenario: Cover a construct absent from any current post

- **WHEN** a mapped construct does not appear in any current post
- **THEN** a focused case covers that construct
- **AND** the case asserts the resulting block type and range offsets

#### Scenario: Author behavior coverage before implementation

- **WHEN** a conversion behavior is implemented
- **THEN** the case that establishes its expected outcome already exists

### Requirement: Raster diagram renders for posts intended for a draft

A post intended for a draft SHALL reference only images whose media type the
upload endpoints accept. A Mermaid diagram in such a post SHALL be referenced as
a raster render of its `.mmd` source, and this change SHALL verify that render
with an end-to-end check. The check SHALL prove the properties this project
depends on: the produced image is in an accepted media type, it was rendered
from the diagram source through the repository's maintained Mermaid render
contract rather than a host browser, a host font, a CDN fetch, or a second
appearance, and the post that references it converts without an
unacceptable-image diagnostic.

<!-- The rule itself (a `mermaid_webp` target under tools/mermaid) and the
     authoring policy are owned by their own workspaces:
     tools/mermaid/openspec/changes/add-raster-diagram-renders and
     projects/alwaldend.com/openspec/changes/add-raster-blog-images. This change
     consumes them and verifies the properties above. -->

#### Scenario: Convert a post that references a rendered diagram

- **WHEN** a post intended for a draft references a diagram's raster render
- **THEN** the conversion check reads that rendered image from the post's package
- **AND** the image is in a media type the upload endpoints accept
- **AND** the post's recorded outcome is the mapping contract, with no
  unacceptable-image diagnostic

#### Scenario: Render a diagram hermetically as a raster image

- **WHEN** the raster render referenced by that post is produced
- **THEN** it is rendered from the diagram source through the pinned Mermaid
  inputs and the maintained appearance, not a host browser, a host font, a CDN
  fetch, or a second appearance

### Requirement: Project registration

The new project SHALL be registered everywhere the repository catalogs direct
projects, so the project builds, tests, and validates like its siblings.

#### Scenario: Register the project

- **WHEN** the project is added under `projects/`
- **THEN** it appears in the project catalog, the OpenSpec validation workspace
  sources, and its own OpenSpec workspace

#### Scenario: Validate the repository

- **WHEN** the repository's OpenSpec validation runs
- **THEN** the new owner workspace is included and its artifacts validate strictly

### Requirement: Project documentation

The project SHALL document its purpose, its two operations, its inputs and
outputs, how to run conversion and publication, and the credential and
disclosure boundary it depends on.

#### Scenario: Read the project documentation

- **WHEN** a reader opens the project README
- **THEN** it explains conversion and publication, their inputs and outputs, and
  how credentials are supplied
- **AND** it states that conversion is offline and publication is the only
  authenticated, network-performing operation
