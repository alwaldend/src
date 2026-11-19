## Purpose

Provide reusable Bazel rules that acquire pinned Terraform providers and run
Terraform with declared configuration, executable, and provider inputs.

## ADDED Requirements

### Requirement: Verified provider acquisition

The provider module extension SHALL download immutable HTTPS archives with
mandatory SHA256 integrity through Bazel repository fetching. It SHALL select
one version for each canonical provider source across the extension graph.

#### Scenario: Conflicting provider versions

- **WHEN** declarations select different versions of the same provider source
- **THEN** extension resolution fails before provider repositories are created.

#### Scenario: Verified archive

- **WHEN** a provider archive is fetched
- **THEN** Bazel verifies the declared integrity and exposes its canonical source,
  version, platform, and packed filesystem-mirror path.

### Requirement: Declared provider execution

Terraform rules SHALL place selected provider archives in runfiles using the
packed filesystem-mirror layout. The runner SHALL configure only that mirror
for provider installation and SHALL preserve the declared Terraform executable.

#### Scenario: Offline initialization

- **WHEN** all required providers are declared by a target
- **THEN** initialization and provider schema validation succeed without registry
  access or a host provider cache.

#### Scenario: Missing provider

- **WHEN** configuration requires a provider absent from the target's mirror
- **THEN** initialization fails without falling back to a registry download.

### Requirement: Reusable execution and repository integration

The reusable module SHALL remain independent of parent repository labels.
Repository AL wrappers SHALL preserve configuration, plugin lifecycle, backend
injection, named operations, and the saved-plan apply guard.

#### Scenario: Saved-plan apply

- **WHEN** a targeted apply wrapper receives a missing plan or extra arguments
- **THEN** the runner rejects it before Terraform initialization.

#### Scenario: Repository migration

- **WHEN** repository Terraform commands and tests are analyzed
- **THEN** they use rules_terraform, select shared provider labels, and require
  neither tools/terraform nor a checked-in .terraform.lock.hcl file.
