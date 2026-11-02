## MODIFIED Requirements

### Requirement: Publish ordinary files without extraction

The tool SHALL publish each file directly beneath
`<namespace>/<project>/releases/<version>/` on the selected host, where namespace is projects by default or third_party
when explicitly selected or identified by the manifest. Regular archives
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

#### Scenario: Publish an upstream schema snapshot

- **WHEN** a third_party release is selected for upload
- **THEN** its payload is available beneath third_party/<name>/releases/<version>/ without affecting project releases

### Requirement: Independently selectable publication steps

Project-namespace SSH deployment SHALL accept an ordered subset of upload, extract, and link,
with upload as the default. Third-party publication SHALL support upload only. Extraction SHALL use a published archive, and
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

#### Scenario: Refuse third-party website selection

- **WHEN** extract or link is selected for a third_party artifact
- **THEN** the command fails before any remote change

## ADDED Requirements

### Requirement: Explicit SSH agent isolation

An explicit identity file SHALL preserve OpenSSH agent authentication by default.
The --ssh_disable_agent option SHALL disable agent authentication and restrict
authentication to configured identities.

#### Scenario: Agent holds the selected private key

- **WHEN** the operator selects a public identity whose private key is held by the agent
- **THEN** the generic publisher authenticates using that key
- **AND** the same request fails authentication when --ssh_disable_agent is selected
