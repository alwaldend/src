## ADDED Requirements

### Requirement: Imagegen rental

The implementation SHALL provide `infra/image_gen/tf_setup` with a generated
login persisted in Vault before a paid rental is created. It SHALL authenticate
with its dedicated `src_infra_image_gen` AppRole and use only that role's fresh
KV and state namespace. The existing team root SHALL issue a dedicated key and
write it to the component's `vastai` entry through an exact-path policy. The
component SHALL NOT use the personal user backend or team operator seed.
The rental SHALL use the direct approved Vast ComfyUI image tag, at least
32 GiB VRAM, verified Secure Cloud hosts and a USD 1/hour ceiling. Offline
validation SHALL NOT rent capacity or deploy the new role, policies or keys.

#### Scenario: Authorized configuration

- **WHEN** an operator supplies valid reviewed configuration
- **THEN** the declared behavior SHALL be available through repository-owned targets

#### Scenario: Offline validation

- **WHEN** validation runs without live credentials
- **THEN** it SHALL NOT rent capacity or write live secrets

### Requirement: Rental-owned default checkpoint

The rental root SHALL call the ComfyUI module and default to the official Comfy-Org Qwen-Image-Edit-2511 mixed-precision checkpoint and its required Qwen 2.5 VL encoder and Qwen VAE from fixed publisher revisions. The shared modules SHALL remain model-configurable.

#### Scenario: Default rental configuration

- **WHEN** model_downloads is omitted
- **THEN** the native provisioner SHALL receive the three public model downloads in native ComfyUI directories before ComfyUI serves generation requests

#### Scenario: Explicit model selection

- **WHEN** the caller supplies model_downloads
- **THEN** that list SHALL replace the rental default, including an empty list to disable configured downloads

### Requirement: Native reference-conditioned generation

The rental SHALL supply native one-reference and separate character/style
workflows using the requested Qwen bundle. Reference images SHALL reach both
the vision-language encoder and reference-latent conditioning. The graphs SHALL
retain the native index_timestep_zero method, flow shift and CFG normalization.
The default profile SHALL disable custom nodes and install no additional Python
packages. Source and validation fixtures SHALL contain only neutral content.

#### Scenario: Default native workflows

- **WHEN** the model bundle and workflow documents finish provisioning
- **THEN** the saved workflow library SHALL contain both native graphs with matching model filenames and reference-role prompts

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
