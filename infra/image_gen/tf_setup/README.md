---
title: Image generation provisioning
description: Generate a Vault-backed login and rent one HTTPS-protected ComfyUI instance
---

This Terraform root generates a 48-character random alphanumeric password,
stores `{username, password}` at mount `secrets`, path
`alwaldend.com/vault1/approles/src_infra_image_gen/web`, and then rents
one ComfyUI instance with that login. A failed Vault write prevents rental creation.
The backend stores state under the component AppRole's `tf_backend/setup` namespace.
The rental credential is injected as `VAST_API_KEY`; Vault provider authentication
is injected separately through AL. No account API key is passed to the container.

## Rent and connect

First deploy the `src_infra_image_gen` AppRole and exact-path team-writer
policy through [Vault Terraform](../../vault/tf/README.md#vast-team-token-delivery).
Then run [team credential setup](../../vastai/alwaldend/README.md) to generate
and store a fresh key in this component's namespace. The component credentials
are deployed; the previous Flux rental and its Base replacement were destroyed at the operator's request. These remain operator instructions. The
[acceptance record](../openspec/changes/add-imagegen-rental/design.md#base-checkpoint-replacement)
owns the observed rental and workflow results.
The [module](../../../projects/tf_modules/vastai_instance_imagegen_comfyui/README.md) uses
the direct Vast ComfyUI/CUDA 12.9 image tag declared in [config.json](../config.json),
without inheriting stock template options. The root selects an explicit published release tag. Update that one image reference to select a newer build;
there is no registry lookup or helper executable during Terraform planning.
Changing the tag produces a reviewed replacement plan for an existing rental.
The component's [config.json](../config.json) owns the image, HTTPS username,
selected model/hardware profiles and their definitions. Terraform has no root
input variables or tfvars presets. The default model is Juggernaut XL v9;
set `model_profile` to `flux` or `qwen` in that file to select FLUX.2 klein Base
9B FP8 or Qwen-Image-Edit-2511.

```sh
bazel run //infra/image_gen/tf_setup:tf.plan
bazel run //infra/image_gen/tf_setup:tf.rent
```

`rent` displays the interactive plan; review it before approving. The selected
`standard` hardware profile requests one GPU with at least 32 GiB VRAM, 100 GB
disk and a USD 1/hour ceiling including disk. The module selects the cheapest
qualifying Secure Cloud offer. The ceiling is not a total budget and market
availability is not guaranteed.

For the prepared 96 GB option, set `hardware_profile` to `96gb` in config.json,
then use the same plan and rent targets after authorizing deployment. That
profile uses a 95 GiB reported-memory floor for nominal 96 GB cards, 150 GB disk
and a USD 3/hour ceiling. Model selection is independent. The committed selection
remains `standard`; no high-memory rental has been deployed.

Creation acknowledges a contract; boot and provisioning take longer. Refresh
through the same owner after boot, reviewing its interactive plan:

```sh
bazel run //infra/image_gen/tf_setup:tf.refresh
bazel run //infra/image_gen/tf_setup:tf.output
```

`https_url` remains absent until a refresh reports the public IP and mapped port.
Verify the Vast instance certificate before entering credentials; follow the
[module's certificate guidance](../../../projects/tf_modules/vastai_instance_imagegen/README.md#certificates).
Retrieve the login through an authorized Vault client using the reference above;
do not paste it into logs or command arguments. `ssh_forward_command` provides
an alternative for account SSH keys. Export images before release:

```sh
bazel run //infra/image_gen/tf_setup:tf.release
```

Changing the password or access settings replaces the instance; export files
first. Destroy removes all versions of the managed web-login entry. Protected historical
Terraform state versions can still retain credentials. The root has
no automatic expiry or spend cutoff. See the [privacy boundary](../README.md#privacy-boundary).

## Offline validation

```sh
bazel test //infra/image_gen/tf_setup:all
```

The checks omit the live backend, validate the actual module/provider schema,
and exercise password generation and Vault persistence with a mocked Vault and
overridden rental module; the separate CLI E2E suite covers the actual rental module. They do not write Vault or rent
capacity. The provider's Terraform CLI acceptance suite emits separate SSH and
HTTPS lifecycle artifacts. Actual image startup and TLS/authentication require
a separately scoped live acceptance run.

## Download models

`model_profile` in config.json selects matching model downloads and saved workflows.
Every saved name uses `<profile>-<purpose>`, matching its generated JSON filename:

| Profile                | Bundle                                                                       | Native workflows                                                                                                                                        |
| ---------------------- | ---------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `juggernaut` (default) | Juggernaut XL v9 all-in-one SDXL checkpoint                                  | `juggernaut-text-to-image`, `juggernaut-image-to-image`, `juggernaut-prompt-sequence`, `juggernaut-style-reference`, `juggernaut-style-prompt-sequence` |
| `flux`                 | FLUX.2 klein Base 9B FP8, Qwen3 8B mixed FP8 encoder, Flux small-decoder VAE | `flux-single-reference`, `flux-character-style`, `flux-prompt-sequence`                                                                                 |
| `qwen`                 | Qwen-Image-Edit-2511 mixed FP8, Qwen 2.5 VL 7B FP8 encoder, Qwen VAE         | `qwen-single-reference`, `qwen-character-style`, `qwen-prompt-sequence`                                                                                 |

All files are fixed to publisher revisions in [config.json](../config.json). All three bundles
fit the 100 GB disk. Vast's native provisioner downloads them inside the rental;
no custom downloader is installed. Every profile stages the first-party
`Prompt Sequence` through native `write_files`. Only Juggernaut also restores
[RES4LYF](https://github.com/ClownsharkBatwing/RES4LYF/tree/6780724de0ddb98e7454d54bb80387836a911ece)
at the previously approved revision and the unchanged
[hash-locked requirements](../res4lyf-requirements.txt). Vast's native Git/pip
provisioner installs them; only RES4LYF and Prompt Sequence are whitelisted for
that profile. Flux/Qwen retain only Prompt Sequence. No additional adapter
weights are downloaded.
The [workflow definitions](../../../projects/comfyui_tools/README.md#workflow-source-and-generation)
are compiled from Libsonnet during the Bazel build; the rental receives native
JSON and does not require Jsonnet. Archived Chroma examples are generated for
documentation and are not staged by these profiles.

Flux uses the publisher's [non-commercial model](https://huggingface.co/black-forest-labs/FLUX.2-klein-base-9b-fp8).
Its checkpoint requires the user's Hugging Face license acceptance and a read-only
model token. The operator or external IaC supplies the
[Vault token entry](../README.md#identity-and-state). The rental root reads it only for a nonempty Flux manifest and
passes it as sensitive startup input to Vast's native `HF_TOKEN` downloader.
Missing access must be resolved before creating paid capacity; Juggernaut and Qwen's public
bundles need no token. All profiles retain at least 32 GiB VRAM and Secure Cloud.

Juggernaut's all-in-one checkpoint goes in `models/checkpoints` and supplies
its own CLIP encoders and VAE. Its text-to-image and prompt-sequence graphs
start with a fresh 1024-square latent. The img2img graph resizes and center-crops
the reference to 1024 square, then samples at denoise 0.65. Raise denoise for
larger changes, or lower it to retain more reference structure. All three use
30 DPM++ 2M SDE/Karras steps and CFG 6.5 as starting settings. Describe the
wanted style in the positive prompt. Img2img also carries over composition;
these native graphs do not isolate only a reference's style or character.
For a reference image's style with a new scene, use
`juggernaut-style-reference`. It routes the VAE-encoded reference through
RES4LYF's style guide and the SDXL model patcher, starting generation from
fresh noise. Describe the scene and broad medium in the positive prompt.
The guide's `weight` controls influence; zero disables it. Start at 1 and
compare at the same seed. `juggernaut-style-prompt-sequence` reuses that guide
for every panel in the growing prompt list. These are style-guided workflows,
not guarantees of style isolation or character preservation. Their sampler
starts at 30 RES 2M/beta57 steps, CFG 6.5 and full denoise; its controls differ
from the native KSampler graphs. Juggernaut GPU execution and style quality
remain unverified.

For Flux and Qwen, select images in Load Image: image 1 supplies character and output aspect;
image 2, when present, supplies style. Describe the new action/scene and the
reference roles in the positive prompt. Flux uses native reference latents on
both conditioning branches, an encoded empty negative prompt, fresh Flux2 noise, 50 Euler steps and CFG 4.
Qwen retains native vision-language and reference-latent conditioning,
`index_timestep_zero`, shift 3.1, CFG normalization, 40 Euler/simple steps,
CFG 4 and denoise 1. These are reference-guided generation workflows; they do
not guarantee exact style isolation or character continuity.

Edit the selected profile's `downloads` array to change its model files. An
empty array disables configured downloads; workflows still require matching
files. URLs must contain no credentials. Source and acceptance assets remain
neutral; private prompts, uploads and generated images do not belong in Git.

Changing the model manifest replaces the rental and removes its local disk.
The retired personal rental was destroyed. The new component credentials and
Qwen rental were subsequently deployed under separate user authorizations.
The Base replacement was created and then destroyed at the operator's request, retaining component credentials. Base GPU workflow acceptance was superseded by that release. Connection coordinates
come from the owning Terraform outputs. Implementation,
teardown and acceptance evidence is recorded in the
[owning change](../openspec/changes/add-imagegen-rental/design.md).
