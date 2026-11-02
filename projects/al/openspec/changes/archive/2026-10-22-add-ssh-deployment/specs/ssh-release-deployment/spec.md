## Purpose

Publish release files and activate archived static websites through the
existing release tool using authenticated SSH and explicit target selection.

## ADDED Requirements

### Requirement: SSH deployment through the release command

The existing release `deploy` command SHALL use rsync over OpenSSH and accept SSH deployment metadata
alongside existing OCI deployments. An SSH-only deployment SHALL require
neither Ansible nor an OCI executable. It SHALL verify the selected host using
the configured SSH trust mechanism and reuse existing release metadata for
version and file identity.

#### Scenario: Deploy an SSH-only release

- **WHEN** a valid SSH release is selected without Soras or Ansible
- **THEN** the tool deploys over SSH to the explicitly configured environment
- **AND** an unknown or mismatched host key causes failure before upload

#### Scenario: Continue using OCI

- **WHEN** an existing OCI release is deployed
- **THEN** its existing deployment behavior remains available

### Requirement: Publish ordinary files without extraction

The tool SHALL publish each file directly beneath
`projects/<project>/releases/<version>/` on the selected host. Regular archives
SHALL remain files. The tool SHALL use rsync over SSH for transfer verification and temporary-file
publication, update existing filenames using rsync semantics, and preserve
unrelated retained files.

#### Scenario: Publish two ordinary artifacts

- **WHEN** a release contains a binary and a compressed archive
- **THEN** both are downloadable under their original filenames in the release directory
- **AND** the archive is not extracted and no public `files/` layer is added

#### Scenario: Interrupt or update an upload

- **WHEN** upload is interrupted
- **THEN** rsync reports failure and leaves previously completed files available
- **AND** a later successful run updates the selected filenames without deleting unrelated files

### Requirement: Activate a declared website archive

Only the extract step with an explicitly designated website archive SHALL
trigger extraction.
The tool SHALL accept `.tar.gz` and `.zip` archives containing a static site
at archive root, preserve the public archive, and publish extracted content
under `sites/<project>/releases/<version>/`. The link step SHALL select the
completed release by atomically changing `sites/<project>/current` after
validation.

#### Scenario: Publish a new static website

- **WHEN** upload, extract, and link are selected for a valid archive with `index.html`
- **THEN** the original archive remains downloadable
- **AND** the site's unchanged main URL serves the completed extracted release

### Requirement: Reject invalid extraction and preserve the selected site

Unsafe paths, unsupported archive entries, malformed archives, insufficient
space, or missing required site content SHALL cause activation to fail.
Concurrent deployments SHALL NOT corrupt release contents or the selected
link. Failed activation SHALL leave the previous selection intact.

#### Scenario: Reject an escaping archive

- **WHEN** a site archive contains traversal, an absolute path, or a link entry
- **THEN** extraction is rejected without writing outside its staging directory
- **AND** the selected site remains unchanged

#### Scenario: Run out of disk space

- **WHEN** upload or extraction cannot complete because storage is full
- **THEN** the tool reports failure and preserves previously published content

### Requirement: Redeploy retained sites by selecting their release

Running the link step for an existing validated extracted release SHALL
select it through the same deploy command without uploading or extracting it
again. The
version SHALL identify the retained extracted release; changed site content
requires a new version. The tool SHALL NOT introduce a separate rollback
command or delete other releases.

#### Scenario: Select an earlier retained release

- **WHEN** an operator runs link for a retained extracted version
- **THEN** only the site's selected-release link changes
- **AND** newer retained releases and their archives remain available

### Requirement: Independently selectable publication steps

SSH deployment SHALL accept an ordered subset of upload, extract, and link,
with upload as the default. Extraction SHALL use a published archive, and
linking SHALL select an existing extracted release without local payload files.

#### Scenario: Prepare a site before selecting it

- **WHEN** an administrator uploads release files and later runs extract alone
- **THEN** the archive becomes an extracted versioned site
- **AND** the current website selection remains unchanged until link is selected

#### Scenario: Select a retained release independently

- **WHEN** an administrator runs link alone with a release manifest
- **THEN** the current symlink selects that existing site version
- **AND** no archive access or rsync transfer is required

#### Scenario: Refuse invalid step selection

- **WHEN** step selection is empty, unknown, duplicated, or out of order
- **THEN** deployment fails before any remote change
