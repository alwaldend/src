## ADDED Requirements

### Requirement: Consistent workflow names

Workflow files SHALL use <profile>-<purpose>.json within their model directory.
Infrastructure and documentation SHALL reference the same names. Renaming the
existing Flux, Qwen and Chroma documents SHALL preserve their graph contents.

#### Scenario: Named workflow selection

- **WHEN** a configured profile is staged
- **THEN** its saved workflow names SHALL match the checked-in files and documented purposes

### Requirement: Juggernaut generation workflows

Juggernaut workflows SHALL load the approved all-in-one checkpoint with
CheckpointLoaderSimple, use its CLIP encoders and VAE, and provide text-to-image,
image-to-image, prompt-sequence, style-reference and style-prompt-sequence paths.
Sequence paths SHALL reuse the existing prompt-list node and paired filenames.

#### Scenario: Reference style and a new scene

- **WHEN** the style-reference workflow is selected
- **THEN** its reference latent SHALL feed the approved RES4LYF style guide, its model SHALL pass through ReSDPatcher, and sampling SHALL start from fresh latent noise

#### Scenario: Offline checks

- **WHEN** source validation runs
- **THEN** it SHALL NOT access a rental or private prompts/images and SHALL NOT claim GPU execution or reference-style quality acceptance

### Requirement: Reusable authoring and native outputs

Workflow source SHALL use shared Libsonnet construction and model-specific
composition. Bazel SHALL generate the existing native JSON names and paths.
The migration SHALL preserve parsed workflow contents, including model settings,
editor layout, IDs, widget metadata and connections. The generated JSON SHALL be
the artifact consumed by infrastructure and documentation; Jsonnet SHALL run
only in the build environment.

#### Scenario: Authoring-only migration

- **WHEN** the fourteen existing workflow definitions are generated
- **THEN** each generated native document SHALL equal its recorded pre-migration JSON value
