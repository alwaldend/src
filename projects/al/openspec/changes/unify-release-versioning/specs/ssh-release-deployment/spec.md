## MODIFIED Requirements

### Requirement: SSH deployment through the release command

The release command SHALL publish with OpenSSH and rsync to an explicitly
selected SSH destination without Ansible or an OCI executable. It SHALL verify
host keys through configured OpenSSH trust. OCI metadata generation SHALL remain
available, but mutable OCI tag publication SHALL fail before any external write.

#### Scenario: Deploy an SSH-only release

- **WHEN** a valid SSH release is selected without Soras or Ansible
- **THEN** publication uses SSH with explicit environment selection
- **AND** an unknown or mismatched host key causes failure before upload

#### Scenario: Continue using OCI

- **WHEN** existing OCI deployment metadata is used to generate a manifest or page
- **THEN** its generated metadata remains available
- **AND** publication refuses mutable registry writes

#### Scenario: Refuse mutable registry publication

- **WHEN** a legacy OCI manifest is selected for publication
- **THEN** the command refuses with an immutability explanation
- **AND** no registry push is executed

### Requirement: Publish ordinary files without extraction

The tool SHALL stage and verify a complete selected project release before atomically publishing it beneath `<namespace>/<project>/releases/<version>/`. Archives SHALL remain ordinary files. An identical existing publication SHALL succeed without changes; different contents SHALL fail. No failed transfer SHALL expose a partial new release.

#### Scenario: Publish two ordinary artifacts

- **WHEN** a release contains a binary and a compressed archive
- **THEN** both become downloadable together under their original filenames
- **AND** no archive is extracted and no public files layer is added

#### Scenario: Interrupt or update an upload

- **WHEN** an upload is interrupted or a conflicting existing version is selected
- **THEN** the command fails without changing the existing public release
- **AND** an absent version becomes visible only after complete verification

#### Scenario: Retry an immutable publication

- **WHEN** the same release is published again
- **THEN** the existing identity and payload bytes are verified without changes
- **AND** an identical existing publication requires no payload staging or transfer
- **AND** conflicting or corrupted existing content causes failure

#### Scenario: Interrupted transfer

- **WHEN** a staged transfer fails
- **THEN** no new public release directory becomes visible
- **AND** earlier completed releases and activation remain unchanged

## ADDED Requirements

### Requirement: Preview artifact publication offline

The release command SHALL provide an offline preview using the same manifest,
destination selection and identity checks as artifact publication.
It SHALL report the resolved project/version, destinations, selected payloads,
and intended operations without contacting destinations or starting plugins.
The preview SHALL NOT authorize publication or claim observed remote readiness.

#### Scenario: Preview destination filtering

- **WHEN** an operator previews publication with an environment or explicit SSH override
- **THEN** the preview shows the same payload selection and destination as publication
- **AND** explicit destination selection of all items is visible without an SSH connection

#### Scenario: Refuse an invalid publication preview

- **WHEN** the selected manifest, identity, or destination is invalid
- **THEN** preview fails with actionable context before any remote operation

### Requirement: Website operations are explicit release subcommands

The tool SHALL expose optional website extraction and activation as
`al release website extract` and `al release website activate`. Artifact publish
SHALL upload only; public commands SHALL reject the retired --steps selector.
Website commands SHALL share destination selection and publication-binding guards
with artifact operations. Each SHALL provide an offline --dry-run preview without
transport, plugin startup, receipt writes, or remote-readiness claims.

#### Scenario: Publish an optional website

- **WHEN** an operator publishes a site archive, extracts it, and activates it through the dedicated commands
- **THEN** each command performs only its named operation against the selected publication
- **AND** projects without a website need only ordinary artifact publication

#### Scenario: Preview a website operation

- **WHEN** an operator requests --dry-run for website extraction or activation
- **THEN** the output reports the fixed operation, selected artifact identity, and destination without external effects
- **AND** combining --dry-run and --receipt is rejected rather than silently ignoring either option

### Requirement: Inspect current release state without remote mutation

The release command SHALL provide read-only inspection of the selected expected
publication and website state. It SHALL distinguish absent, verified, conflicting,
and invalid publication observations; report retained-site binding and current
activation; and keep local payload verification separate from remote verification.
Inspection SHALL NOT create remote directories, locks, staging, receipts, or
repairs, and SHALL NOT change activation. Its output SHALL omit credentials and
represent an observation rather than a publication lease.

#### Scenario: Recover an uncertain publication without retrying it

- **WHEN** an operator inspects a release after an interrupted attempt
- **THEN** the command verifies current publication bytes, retained-site binding, and activation
- **AND** remote contents and links are unchanged by inspection

#### Scenario: Inspect after local payload removal

- **WHEN** the expected manifest remains but local payload files are unavailable
- **THEN** inspection reports that local verification is unavailable
- **AND** it can still verify and report the selected remote publication

#### Scenario: Preserve transport uncertainty

- **WHEN** inspection cannot connect or complete remote verification
- **THEN** it fails without reporting the release as absent or verified
- **AND** it does not attempt publication or repair

### Requirement: Delete a project release explicitly

The tool SHALL support deletion at an explicitly selected SSH destination, refusing deletion while the website release is active. Deletion SHALL remove artifacts and extracted content, preserve shared Git refs, and retain no restoration state. Subsequent ordinary publication SHALL accept the absent identity.

#### Scenario: Delete and republish

- **WHEN** an inactive release is deleted
- **THEN** its artifacts and extracted content are absent
- **AND** ordinary publication can recreate that project and version

#### Scenario: Delete the selected website

- **WHEN** deletion targets the currently selected website release
- **THEN** deletion fails before removing its content

### Requirement: Use selected payloads consistently

Publication, extraction, and activation SHALL bind the same selected destination
payloads. Regular-file payload targets exposed through Bazel runfiles symlinks
SHALL be accepted and staged bytes SHALL retain size and SHA-256 verification.

#### Scenario: Activate an environment-filtered website

- **WHEN** a manifest includes files assigned to different environments
- **THEN** selecting one environment publishes and activates its selected files
- **AND** files assigned only to other environments are not required there

#### Scenario: Publish a runfiles payload

- **WHEN** a declared payload symlink resolves to matching regular-file bytes
- **THEN** publication succeeds using a verified regular staging copy

### Requirement: Retain sites and protect managed destinations

Activation SHALL retain earlier extracted sites for explicit rollback. Publication
immutability SHALL apply to both projects and third_party namespaces. Activation
and deletion SHALL refuse symlinked managed ancestry. Deletion SHALL reject mount
points in either operand, their managed ancestry below the configured content
root, or any descendants, including bind-mounted files, before removing either release tree.

#### Scenario: Roll back a retained site

- **WHEN** another complete version is activated and an earlier version is selected
- **THEN** both extracted versions and their publication archives remain available

#### Scenario: Refuse unsafe deletion

- **WHEN** a deletion operand has symlinked ancestry or contains a directory or file mount
- **THEN** deletion fails without removing either release tree or mounted state

### Requirement: Record artifact operation observations

The release commands SHALL accept an optional local receipt path for publication,
activation, extraction, and deletion. It SHALL atomically checkpoint resolved
identity, destination, publication digest and ordered steps before remote effects
and after each confirmed step, without credentials or raw error output.

#### Scenario: Retain partial progress

- **WHEN** one selected release operation completes but a later selected release does not confirm success
- **THEN** the receipt preserves the completed operation and the unconfirmed operation as unknown
- **AND** the overall receipt is incomplete and does not claim the destination was unchanged

#### Scenario: Refuse an unwritable receipt before remote effects

- **WHEN** a requested receipt cannot be checkpointed before a step
- **THEN** that step does not start
