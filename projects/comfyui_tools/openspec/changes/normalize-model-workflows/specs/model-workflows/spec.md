## ADDED Requirements

### Requirement: Grouped panel progress and unique runs

Prompt sequences SHALL remain one native ComfyUI job and media group. The
sampler SHALL display the current panel and total through native progress text.
Each execution SHALL create a unique shared filename prefix while retaining
ordered panel suffixes. Unchanged text queued again SHALL receive a new run
identity and fresh counter. Existing prompt/filename/negative slots SHALL remain.

#### Scenario: Repeated grouped generation

- **WHEN** the same three-panel workflow is queued twice
- **THEN** each run SHALL remain one job with three grouped outputs, its own filename prefix and a counter from panel 1 of 3 through panel 3 of 3

### Requirement: Chained img2img GIF

The full Chroma profile SHALL provide a chained image-to-image GIF workflow using
one prompt per frame. Each generated latent SHALL become the next frame's input.
The output SHALL be a looping GIF at the configured frame interval, with one
generated frame per prompt and no new external dependency.

#### Scenario: Ordered frame feedback

- **WHEN** a reference image and three frame prompts are supplied
- **THEN** the workflow SHALL sample three frames in order, feeding each result into the next sample, and save one three-frame GIF

### Requirement: Full Chroma1-HD workflows

The chroma-hd profile SHALL provide text-to-image, image-to-image, style-reference,
prompt-sequence and style-prompt-sequence graphs using its configured full
checkpoint, FP16 T5 encoder and matching VAE. Native graphs SHALL start at
40 Euler/beta steps and CFG 3; style graphs SHALL use the already approved
RES4LYF Chroma patcher and guide with 40 RES 2M/beta57 steps and CFG 3.
Both sequence graphs SHALL connect the shared negative prompt. Archived Flash
graphs SHALL retain their existing model names, sampling settings and links.

#### Scenario: Full checkpoint and ordered panels

- **WHEN** the selected HD profile is built and three nonempty panel lines are supplied
- **THEN** generated graphs SHALL reference the selected model files and map one ordered output filename per panel while sharing the scalar negative prompt

#### Scenario: Style reference

- **WHEN** the HD style workflow is selected
- **THEN** the reference latent SHALL feed the approved style guide and the model SHALL pass through ReChromaPatcher before fresh-noise sampling

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

### Requirement: Several character and style references

Juggernaut SHALL provide multi-reference and multi-reference-prompt-sequence
workflows using the approved SDXL Plus IP-Adapter and CLIP-ViT-H encoder. Three
reference inputs SHALL contribute to one conditioning batch; the sequence SHALL
reuse that conditioning for each nonempty prompt line. Native KSampler settings
and model loading SHALL be reused. Reference conditioning SHALL NOT be presented
as guaranteed identity preservation or isolation of every subject in an image.

#### Scenario: Several references for a panel sequence

- **WHEN** three reference images and two nonempty prompt lines are supplied
- **THEN** the graph SHALL share all references across two ordered panel generations

### Requirement: Shared sequence negative prompt

Every prompt-sequence workflow with negative conditioning SHALL expose a shared multiline negative prompt
in Prompt Sequence. Its scalar output SHALL feed the existing model-specific
negative text encoder through a converted input instead of a separate editable
negative-prompt box. New workflows SHALL retain their existing negative defaults.
The node SHALL preserve its first four widget positions and its prompt/filename
output slots, and SHALL accept an omitted negative prompt as empty text.

#### Scenario: One negative prompt for several panels

- **WHEN** a sequence contains two nonempty panel lines and a multiline shared negative prompt
- **THEN** it SHALL produce two ordered positive prompts and filenames, and one unchanged shared negative text value

#### Scenario: Existing saved input

- **WHEN** a saved node supplies only the original four inputs and uses the original two output slots
- **THEN** its panel and filename outputs SHALL retain their meaning and the additional negative output SHALL be empty

### Requirement: Native FLUX.2 dev guidance

The FLUX.2 dev 32B workflows SHALL use the requested FP8 model, local Mistral
encoder and matching VAE with native BasicGuider, FluxGuidance, Flux2Scheduler
and Euler sampling. Three reference images SHALL contribute native VAE latents
to every panel in the multi-reference sequence. The guidance-distilled path
SHALL NOT imply that an unconnected negative-prompt output affects sampling.

#### Scenario: Three references and ordered panels

- **WHEN** the operator supplies three references and several nonempty scene lines
- **THEN** the sequence SHALL generate one panel per line using all three references and the same native guidance settings

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
