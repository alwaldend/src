## ADDED Requirements

### Requirement: Imagegen rental

The implementation SHALL provide `infra/image_gen/tf_setup` with a generated
login persisted in Vault before a paid rental is created. It SHALL authenticate
with its dedicated `src_infra_image_gen` AppRole and use only that role's fresh
KV and state namespace. The existing team root SHALL issue a dedicated key and
write it to the component's `vastai` entry through an exact-path policy. The
component SHALL NOT use the personal user backend or team operator seed.
The rental SHALL use the direct approved Vast ComfyUI image tag, at least
32 GiB VRAM, verified Secure Cloud hosts and a USD 1/hour default ceiling. Only
the explicitly selected high-memory preset SHALL use its separate ceiling. Offline
validation SHALL NOT rent capacity or deploy the new role, policies or keys.

#### Scenario: Authorized configuration

- **WHEN** an operator supplies valid reviewed configuration
- **THEN** the declared behavior SHALL be available through repository-owned targets

#### Scenario: Offline validation

- **WHEN** validation runs without live credentials
- **THEN** it SHALL NOT rent capacity or write live secrets

### Requirement: Rental-owned default checkpoint

The rental root SHALL read the image, login name, hardware and model configuration
from the component config.json, with selectable chroma-hd, juggernaut, flux,
flux2-dev and qwen profiles. The selected chroma-hd profile SHALL use the
previously approved full Chroma1-HD checkpoint at revision
0e0c60ece1e82b17cb7f77342d765ba5024c40c0 with its existing FP16 T5 and FLUX VAE
pins. Juggernaut SHALL retain its previously approved all-in-one checkpoint;
the remaining profiles SHALL retain their existing publisher-pinned bundles.
The shared modules SHALL remain model-configurable.

#### Scenario: Default rental configuration

- **WHEN** the committed model profile is selected
- **THEN** the native provisioner SHALL receive the selected profile's model downloads in native ComfyUI directories before ComfyUI serves generation requests

#### Scenario: Explicit model selection

- **WHEN** the operator edits the selected profile's downloads array in config.json
- **THEN** that list SHALL supply the rental downloads, including an empty list to disable configured downloads

### Requirement: Native reference-conditioned generation

The rental SHALL supply native SDXL text-to-image, img2img and prompt-sequence
and approved RES4LYF style-reference workflows for Juggernaut, plus one-reference and separate character/style
workflows for Flux/Qwen. The Flux/Qwen profiles SHALL provide native reference
latents; Qwen SHALL also retain vision-language conditioning, index_timestep_zero,
flow shift and CFG normalization. Flux SHALL use the Base native
schedule with 50 steps, CFG 4, an encoded empty negative prompt and fresh latent noise.
All profiles SHALL enable the dependency-free first-party prompt sequence node.
Chroma1-HD and Juggernaut SHALL additionally restore the previously approved RES4LYF revision
and its exact locked Python requirements. Source and validation fixtures SHALL contain only neutral content.

#### Scenario: Default native workflows

- **WHEN** the model bundle and workflow documents finish provisioning
- **THEN** the saved workflow library SHALL contain its configured graphs with matching model filenames and model-appropriate prompts

#### Scenario: Meaningful reference and prompt behavior

- **WHEN** neutral references specify a character and style and the prompt requests a different action or setting
- **THEN** live acceptance SHALL assess the requested action/setting and recognizable reference attributes, rather than treating execution alone as success

### Requirement: Independent lifecycle

The new root SHALL use `src_infra_image_gen/tf_backend/setup` and SHALL NOT
import legacy state or credentials. Destroying the replacement's managed login
and delivered API-key entries SHALL delete all versions of those entries;
protected historical Terraform state is outside that guarantee.

#### Scenario: Legacy teardown

- **WHEN** the old stack is destroyed as authorized
- **THEN** its rental, dedicated API key, current login entry and obsolete writer grant SHALL be removed while the shared team and operator seed remain

#### Scenario: Fresh configuration delivery

- **WHEN** the replacement IaC is published
- **THEN** new AppRole, policy, key and rental deployment SHALL remain separate from source validation

### Requirement: CLI and model authentication

The component SHALL expose the approved official CLI through an owning target
that reads the current Terraform HTTPS output and its Vault web login. Gated
Flux downloads SHALL read a scoped Hugging Face token from the component's own
Vault namespace, populated manually by the operator or through external IaC.
This component SHALL only read that entry and SHALL NOT create, rotate or delete it.
No Vault or Vast account token SHALL be passed to the rental.

#### Scenario: Missing model access

- **WHEN** the publisher has not granted the token download access
- **THEN** rental creation SHALL wait for access verification rather than consume paid capacity

### Requirement: Prompt sequence generation

Each model profile SHALL supply a single-node multiline prompt workflow with
one numbered image per nonempty prompt. Flux/Qwen SHALL reuse shared
character/style references; Juggernaut SHALL use shared text instructions.
The node SHALL retain prompt order, support line and blank-line separators and
reject empty input or unsafe filename prefixes without echoing prompt content.

#### Scenario: Several panels

- **WHEN** the user enters three prompts and runs the sequence workflow once
- **THEN** native list mapping SHALL generate three ordered panels with paired numbered filenames and the same shared conditioning

### Requirement: Opt-in high-memory hardware

The component config.json SHALL provide a nominal 96 GB-or-larger hardware
profile with a 95 GiB reported-memory floor, 150 GB disk and a USD 3/hour ceiling.
The ordinary plan and rent targets SHALL consume the selected hardware profile.
The standard option SHALL retain the 32 GiB, 100 GB and USD 1/hour defaults. The
selected 48gb option SHALL require RTX PRO 5000, at least 47 GiB reported VRAM,
100 GB disk and a USD 1/hour ceiling. All host trust requirements
SHALL remain enforced. Preparing this option SHALL NOT deploy it or change the
selected model bundle.

#### Scenario: Prepare the larger option

- **WHEN** the preset and its targets are built or validated
- **THEN** no high-memory rental SHALL be created
