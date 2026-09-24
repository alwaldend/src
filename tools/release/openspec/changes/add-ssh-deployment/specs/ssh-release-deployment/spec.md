## Purpose

Publish release files and activate archived static websites through the
existing release tool using authenticated SSH and explicit target selection.

## ADDED Requirements

### Requirement: SSH deployment through the release command

The existing release `deploy` command SHALL accept SSH deployment metadata
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
SHALL remain files. A final filename SHALL become visible only after a
completed verified upload. An identical repeat SHALL be idempotent; conflicting
bytes at the same version and filename SHALL fail without replacement.

#### Scenario: Publish two ordinary artifacts

- **WHEN** a release contains a binary and a compressed archive
- **THEN** both are downloadable under their original filenames in the release directory
- **AND** the archive is not extracted and no public `files/` layer is added

#### Scenario: Interrupt or conflict with an upload

- **WHEN** upload is interrupted or different bytes already occupy its final path
- **THEN** no incomplete replacement is exposed at that path
- **AND** the tool reports the failed artifact distinctly from completed artifacts

### Requirement: Activate a declared website archive

Only an explicitly designated website archive SHALL trigger extraction.
The tool SHALL accept `.tar.gz` and `.zip` archives containing a static site
at archive root, preserve the public archive, and publish extracted content
under `sites/<project>/releases/<version>/`. It SHALL select the completed
release by atomically changing `sites/<project>/current` after validation.

#### Scenario: Publish a new static website

- **WHEN** a valid website archive with `index.html` is deployed
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

Deploying an existing validated extracted release SHALL select it through the
same deployment operation without rebuilding or extracting it again. The
tool SHALL verify the existing release identity before selection and SHALL
NOT introduce a separate rollback command or delete other releases.

#### Scenario: Select an earlier retained release

- **WHEN** an operator deploys a retained version whose archive identity matches
- **THEN** only the site's selected-release link changes
- **AND** newer retained releases and their archives remain available
