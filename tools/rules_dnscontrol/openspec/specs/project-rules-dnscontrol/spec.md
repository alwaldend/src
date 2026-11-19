# Rules DNSControl

## Purpose

Package DNS record configurations and supporting data for DNSControl's
JavaScript loader. This source baseline was observed on 2026-09-08 at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md) and
[packaging macro](../../../main/bzl/dnscontrol.bzl).
The current macro accepts a `config` parameter but does not consume it; callers
must include an entrypoint through another packaged input. The generated
manifest is `<name>.json`, so the name `requires` produces `requires.json`.

## Requirements

### Requirement: Generate a relative JSON record manifest

`dnscontrol_site` SHALL create a JSON array manifest named `<name>.json` whose
entries are relative paths to symlinks for the files selected by `srcs`.

#### Scenario: A site declares multiple record files

- **WHEN** a site supplies record files through `srcs`
- **THEN** its manifest lists each generated link with a `./` prefix so a loader can resolve it relative to the manifest directory.

### Requirement: Retain a JSON suffix on record links

Every generated record symlink SHALL end with `.json`, regardless of its source
filename, to preserve DNSControl's JSON-loader selection.

#### Scenario: A record source has a Bazel-specific path

- **WHEN** the manifest rule converts a record source path into a generated link name
- **THEN** it replaces path separators and dots in that path component and appends `.json` to the link name.

### Requirement: Aggregate the manifest and explicit supporting data

The public site target SHALL aggregate the generated manifest and record links
with files explicitly supplied in `data`, using rules_pkg packaging providers.

#### Scenario: The caller supplies a JavaScript entrypoint in data

- **WHEN** the caller includes its entrypoint in `data`
- **THEN** the site package includes that file together with the record manifest and links.
