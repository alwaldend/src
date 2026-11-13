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
[acceptance record](../openspec/changes/add-imagegen-rental/design.md#requested-flux2-dev-32b-replacement)
owns the observed rental and workflow results.
The [module](../../../projects/tf_modules/vastai_instance_imagegen_comfyui/README.md) uses
the direct Vast ComfyUI/CUDA 12.9 image tag declared in [config.json](../config.json),
without inheriting stock template options. The root selects an explicit published release tag. Update that one image reference to select a newer build;
there is no registry lookup or helper executable during Terraform planning.
Changing the tag produces a reviewed replacement plan for an existing rental.
The component's [config.json](../config.json) owns the image, HTTPS username,
selected model/hardware profiles and their definitions. Terraform has no root
input variables or tfvars presets. The selected model is NoobAI XL 1.1 with the
approved SDXL Plus adapter and CLIP vision encoder, on the
existing 48 GB hardware profile. Other model profiles remain available
as `chroma-hd`, `flux2-dev`, `juggernaut`, `hidream-i1-full`, `illustrious`, `flux`,
`qwen`, `uso` and `umo-uno`. Selecting a profile does not rent until an explicitly authorized apply.

HiDream's eight workflow files are uploaded through the authenticated ComfyUI
API after startup. Its empty `startup_workflows` selection keeps those JSON
files out of the inline command; other profiles default to their full catalog.
Models, approved nodes and hardening still use the normal native startup.

```sh
bazel run //infra/image_gen/tf_setup:tf.plan
bazel run //infra/image_gen/tf_setup:tf.rent
```

`rent` displays the interactive plan; review it before approving. The selected
`48gb` hardware profile requires RTX PRO 5000 Blackwell, at least 47 GiB reported
VRAM, 100 GB disk and a USD 1/hour ceiling. The optional `96gb` profile requests
at least 95 GiB reported VRAM, 150 GB disk and a USD 3/hour ceiling.
The optional `standard` profile requests one GPU with at least 32 GiB VRAM, 100 GB
disk and a USD 1/hour ceiling including disk. The module selects the cheapest
qualifying Secure Cloud offer. The ceiling is not a total budget and market
availability is not guaranteed.

Model and hardware selections are independent. The named GPU filter applies
only to `48gb`; other profiles retain automatic GPU selection.

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
HTTPS lifecycle artifacts. Its configured-bundle check packages every model
profile through the real module against a loopback API, checking workflow
contents, model destinations and the 16 KiB startup limit. It does not download
weights or evaluate generated images. Actual image startup and TLS/authentication require
a separately scoped live acceptance run.

## Download models

`model_profile` in config.json selects matching model downloads and saved workflows.
Every saved name uses `<profile>-<purpose>`, matching its generated JSON filename:

Juggernaut also stages `juggernaut-multi-reference` and
`juggernaut-multi-reference-prompt-sequence`, with three shared character/style
references through IP-Adapter Plus. The two additional public safetensors files
add about 3.4 GB: the SDXL Plus adapter in `models/ipadapter` and its CLIP-ViT-H
encoder in `models/clip_vision`. Their publisher revisions are fixed in config.json.

| Profile             | Bundle                                                                             | Native workflows                                                                                                                                                                   |
| ------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `juggernaut`        | Juggernaut XL v9 all-in-one SDXL checkpoint                                        | `juggernaut-text-to-image`, `juggernaut-image-to-image`, `juggernaut-prompt-sequence`, `juggernaut-style-reference`, `juggernaut-style-prompt-sequence`                            |
| `noobai` (selected) | NoobAI XL 1.1 all-in-one SDXL checkpoint, SDXL Plus adapter and CLIP vision        | Seven purposes as Juggernaut, plus img2img prompt sequence; prefixed `noobai-`                                                                                                     |
| `illustrious`       | Illustrious XL v2.0-STABLE SDXL checkpoint, SDXL Plus adapter and CLIP vision      | Seven purposes as Juggernaut, prefixed `illustrious-`                                                                                                                              |
| `hidream-i1-full`   | HiDream-I1 Full 17B FP16, four local text encoders, FLUX VAE                       | Eight purposes: txt2img, img2img, their prompt sequences, style and its sequence, and both GIF paths; prefixed `hidream-i1-full-`                                                  |
| `chroma-hd`         | Full Chroma1-HD, FP16 T5 encoder, FLUX VAE                                         | `chroma-hd-text-to-image`, `chroma-hd-image-to-image`, `chroma-hd-style-reference`, `chroma-hd-prompt-sequence`, `chroma-hd-style-prompt-sequence`, `chroma-hd-image-to-image-gif` |
| `flux2-dev`         | FLUX.2 dev 32B FP8, local Mistral FP8 encoder, FLUX.2 VAE                          | `flux2-dev-text-to-image`, `flux2-dev-single-reference`, `flux2-dev-single-reference-prompt-sequence`, `flux2-dev-multi-reference`, `flux2-dev-multi-reference-prompt-sequence`    |
| `flux`              | FLUX.2 klein Base 9B FP8, Qwen3 8B mixed FP8 encoder, Flux small-decoder VAE       | `flux-single-reference`, `flux-character-style`, `flux-prompt-sequence`                                                                                                            |
| `qwen`              | Qwen-Image-Edit-2511 mixed FP8, Qwen 2.5 VL 7B FP8 encoder, Qwen VAE               | `qwen-single-reference`, `qwen-character-style`, `qwen-prompt-sequence`                                                                                                            |
| `uso`               | FLUX.1 dev FP8 all-in-one checkpoint, USO LoRA/projector and SigLIP vision encoder | Subject reference, style reference and combined subject/two-style references, each with a sequence; prefixed `uso-`                                                                |
| `umo-uno`           | Same FLUX.1 dev FP8 checkpoint and UMO-UNO LoRA                                    | Single and three-identity references, each with a sequence; prefixed `umo-uno-`                                                                                                    |

All files are fixed to publisher revisions in [config.json](../config.json).
The full Chroma bundle totals about 28 GB. The optional 32B and HiDream
bundles each total about 55 GB of model files. Model and hardware selections
are independent: choose sufficient disk/VRAM for the bundle; `48gb` provides
100 GB disk and `96gb` provides 150 GB. No new model is deployed by these additions.
Vast's native provisioner downloads them inside the rental;
no custom downloader is installed. Every profile stages the first-party
`Prompt Sequence` through native `write_files`. Full Chroma and HiDream additionally stage
the first-party chained img2img GIF node, using only ComfyUI's existing Torch,
NumPy and Pillow packages. Both first-party nodes are whitelisted for those profiles.
Chroma1-HD, HiDream and all three SDXL profiles restore
[RES4LYF](https://github.com/ClownsharkBatwing/RES4LYF/tree/6780724de0ddb98e7454d54bb80387836a911ece)
at the previously approved revision and the unchanged
[hash-locked requirements](../res4lyf-requirements.txt). Vast's native Git/pip
provisioner installs them. The SDXL profiles additionally stage the explicitly approved
ComfyUI_IPAdapter_plus pack at revision a0f451a5113cf9becb0847b92884cb10cbdec0ef;
only those two packs and Prompt Sequence are whitelisted for those profiles.
The general SDXL Plus path uses the existing image environment and requires no
FaceID/InsightFace installation. Flux, FLUX.2 dev and Qwen retain only Prompt Sequence.
The [workflow definitions](../../../projects/comfyui_tools/README.md#workflow-source-and-generation)
are compiled from Libsonnet during the Bazel build; the rental receives native
JSON and does not require Jsonnet. The archived Chroma Flash examples remain
available separately and are not staged by the full Chroma profile.

Chroma1-HD stages the already approved RES4LYF pack and its unchanged locked
requirements alongside Prompt Sequence. Native graphs start at 40 Euler/beta
steps and CFG 3; style graphs use 40 RES 2M/beta57 steps and CFG 3 with the
Chroma model patcher. Img2img preserves some reference composition, while the
style path starts from fresh noise and guides style. Neither path guarantees
character identity. Both sequence graphs share their negative prompt across
all nonempty panel lines. This full checkpoint is distinct from the archived
Flash model; additional VRAM alone does not establish better image quality.

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
for every nonempty prompt line. These are style-guided workflows,
not guarantees of style isolation or character preservation. Their sampler
starts at 30 RES 2M/beta57 steps, CFG 6.5 and full denoise; its controls differ
from the native KSampler graphs. All seven Juggernaut workflows completed with
neutral fixtures on the selected runtime; see the
[GPU acceptance record](../openspec/changes/add-imagegen-rental/design.md#replacement-gpu-acceptance).
Execution does not guarantee reference fidelity or scene adherence. The adapter
workflows give more direct character/style conditioning; lower adapter weight
when references overpower the requested scene.

For FLUX.2 dev, use three Load Image nodes as combined character/style
references. Its native sequence uses the same three reference latents for
every prompt. Start at 1024 square, batch one, 28 Euler steps and FluxGuidance 4.
It uses BasicGuider with distilled positive guidance, so a negative prompt has
no sampling effect in this family. The matching Mistral encoder runs inside the
rental; no remote text-encoding API is configured. The official ComfyUI FP8
bundle is public and does not receive the gated Klein profile's Vault HF token.

For Flux Base and Qwen, select images in Load Image: image 1 supplies character and output aspect;
image 2, when present, supplies style. Describe the new action/scene and the
reference roles in the positive prompt. Flux uses native reference latents on
both conditioning branches, an encoded empty negative prompt, fresh Flux2 noise, 50 Euler steps and CFG 4.
Qwen retains native vision-language and reference-latent conditioning,
`index_timestep_zero`, shift 3.1, CFG normalization, 40 Euler/simple steps,
CFG 4 and denoise 1. These are reference-guided generation workflows; they do
not guarantee exact style isolation or character continuity.

## Additional model options

[USO](https://docs.comfy.org/tutorials/flux/flux-1-uso) and
[UMO-UNO](https://github.com/bytedance/UMO) are optional
native ComfyUI profiles. USO downloads about 18.6 GB, including its FLUX.1 dev
FP8 checkpoint with local text encoders and VAE, a LoRA, a style projector and
SigLIP vision weights. UMO-UNO downloads about 19.2 GB, sharing that checkpoint
and adding its own LoRA. All five registry files are fixed to publisher
revisions. Both use the existing runtime image and Prompt Sequence node;
neither installs another custom-node pack or Python requirement. The shared
checkpoint retains the FLUX.1 dev non-commercial license; the USO, UMO and
SigLIP additions use Apache 2.0. These model options have only offline
packaging/graph validation, without GPU or reference-quality acceptance.

For USO, `uso-single-reference` uses one subject; `uso-style-reference` uses
one style image; `uso-character-style` uses one subject and two style images.
Use the same image in both style inputs when only one style reference is
needed. Each has a sequence counterpart. Subject references are resized with
aspect preserved, starting at a 1024-pixel longest side. Style images use the
native SigLIP crop control. UMO-UNO offers a single-reference or ordered
three-reference path; its publisher's starting reference size is 512 pixels
on the longest side. Repeating one subject in multiple inputs is not a
substitute for separate identities. It has no USO style-projector branch.
Describe reference roles and the wanted scene in natural language.

Both families start from fresh noise, with editable 1024-square output size,
Euler/simple sampling and CFG 1. USO starts at 20 steps and FluxGuidance 3.5;
UMO starts at 25 steps and FluxGuidance 4. The zeroed negative branch has no
effect at CFG 1, including the sequence node's negative output. Sequences reuse
the same original reference inputs per panel, showing progress and unique
grouped filenames. These are independent reference-conditioned panels, not
joint attention across the sequence. No reference workflow guarantees exact
style isolation or identity continuity.

The profiles below remain selectable; later rental and acceptance observations
belong to the linked deployment record. Model size alone does not predict
comic quality or reference adherence.

- [NoobAI XL 1.1](https://huggingface.co/Laxhar/noobai-XL-1.1) is an SDXL
  illustration fine-tune. Its publisher recommends Euler ancestral, 25–30
  steps and CFG 5–6. The starting preset uses 30 steps and CFG 5.5. Its custom
  license prohibits commercializing the model, derivatives and generated
  products, and imposes obligations on derivatives; review it before use.
- [Illustrious XL v2.0-STABLE](https://huggingface.co/OnomaAIResearch/Illustrious-XL-v2.0)
  is another SDXL illustration model under CreativeML OpenRAIL-M. Its matching
  starting preset is provisional, not a measured quality recommendation.
- [HiDream-I1 Full](https://huggingface.co/HiDream-ai/HiDream-I1-Full) is a 17B
  alternative. The denoiser is FP16; its published Llama 3.1 8B encoder is FP8.
  The bundle also reuses the existing FP16 T5 and FLUX VAE, plus model-specific
  CLIP-L/G. Native sampling follows the
  [official workflow](https://github.com/Comfy-Org/workflow_templates/blob/a8bc07a76478dfada740fa8fa49ac87be1bff519/templates/hidream_i1_full.json):
  50 UniPC/simple steps, CFG 5, shift 3. The transformer and CLIP weights use
  MIT licenses, T5 and the VAE use Apache 2.0, and Llama retains Meta's Llama
  3.1 Community License. The complete bundle is not licensed solely under MIT.

The 48 GB option can rely on ComfyUI's model offloading for the larger bundle;
peak VRAM, speed and style-patcher compatibility remain unmeasured. Several
models can be stored on disk and loaded as needed. Two 9B FP16 denoisers alone
need about 36 GB before encoders, activations and runtime overhead; three need
about 54 GB and cannot all remain resident on 48 GB. Quantization changes that
calculation, and ComfyUI may unload models between jobs.

The `downloads` registry owns each immutable publisher URL and destination kind;
the `custom_nodes` registry owns approved pack revisions. Each model profile's
`downloads` array lists registry filenames, and `custom_nodes` lists pack names.
HiDream's `components` map assigns its four text-encoder roles explicitly.
Add or select registry entries to change a profile's files. An
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
