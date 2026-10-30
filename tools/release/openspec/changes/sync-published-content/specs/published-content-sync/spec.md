## Purpose

Publish complete versioned releases over SSH, reject duplicated local versions,
and retain only the selected extracted website while preserving archives.
Infrastructure synchronization is specified by infra/download.

## ADDED Requirements

### Requirement: Reject duplicate local release versions

Normal SSH publication SHALL reject a version already present in artifacts or
extracted sites, without modifying published content or current selection.
New uploads SHALL become public only after their full transfer completes.
Public synchronization SHALL continue to overwrite same-version destination
content from the authoritative source.

#### Scenario: Two publishers select the same new version

- **WHEN** concurrent upload attempts target the same project and version
- **THEN** exactly one can promote its complete directory
- **AND** the other fails without changing the completed release

#### Scenario: Retry extraction or select an existing version

- **WHEN** an extracted version already exists
- **THEN** extraction fails without replacing it
- **AND** an explicit link-only operation can still select it

### Requirement: Retain only the selected extracted website

Successful link publication SHALL remove other extracted versions of that same
website after atomically selecting the complete new version. It SHALL retain
all release archives, other projects, and private state. Upload-only or
extract-only publication SHALL NOT remove prior extracted versions. Cleanup
SHALL NOT follow symlinks or cross filesystem boundaries.

#### Scenario: Publish a replacement site

- **WHEN** a complete new website is selected successfully
- **THEN** current selects that version and other extracted versions are removed
- **AND** downloadable archives remain unchanged locally and after public sync

#### Scenario: Selection or cleanup fails

- **WHEN** new-site validation fails
- **THEN** the previous selection and extracted sites remain unchanged
- **WHEN** cleanup fails after selection
- **THEN** publication reports failure while keeping the new site selected
- **AND** repeating link selection can complete cleanup
