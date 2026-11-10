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
are deployed; the last rental was released. These remain operator instructions. The
[acceptance record](../openspec/changes/add-imagegen-rental/design.md#authorized-qwen-rental-acceptance)
owns the observed rental and workflow results.
The [module](../../../projects/tf_modules/vastai_instance_imagegen_comfyui/README.md) uses
the direct Vast ComfyUI/CUDA 12.9 image tag declared in [main.tf](main.tf),
without inheriting stock template options. The root selects an explicit published release tag. Update that one image reference to select a newer build;
there is no registry lookup or helper executable during Terraform planning.
Changing the tag produces a reviewed replacement plan for an existing rental.
The rental defaults to Qwen-Image-Edit-2511; override public model downloads using
`model_downloads` in a task-owned tfvars file. An optional reviewed `template_hash_id` may
supply provisioning environment settings; its Docker options must be environment-only.

```sh
bazel run //infra/image_gen/tf_setup:tf.plan
bazel run //infra/image_gen/tf_setup:tf.rent
```

`rent` displays the interactive plan; review it before approving. Defaults are
32 GiB VRAM, one GPU, 100 GB disk, and a USD 1/hour ceiling including requested
disk. `max_hourly_price` and `offer_id` can be supplied explicitly. The ceiling
is not a total budget and market availability is not guaranteed.

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

The default `model_downloads` in [main.tf](main.tf) selects the official
Comfy-Org Qwen-Image-Edit-2511 mixed-precision diffusion checkpoint, Qwen 2.5 VL
7B FP8-scaled encoder and Qwen VAE at fixed publisher revisions. These are the
components used by [ComfyUI's native Qwen workflow](https://docs.comfy.org/tutorials/image/qwen/qwen-image-edit-2511).
The bundle totals about 30.2 GB and fits the 100 GB disk. Vast's native
provisioner downloads it inside the rental before generation. No Hugging Face
token, custom node pack, acceleration LoRA or additional Python package is used.

Terraform stages two neutral [workflow documents](../workflows/) at startup:

- `qwen-image-edit-single-reference`: use one image for character and style.
- `qwen-image-edit-character-style`: use image 1 for the character and image 2 for the style.

Choose the reference images in Load Image and describe the new scene in the
positive prompt. The graphs feed images to both Qwen's vision-language encoder
and its native reference latents. They retain the published
`index_timestep_zero` reference method, flow shift 3.1 and CFG normalization.
The full model starts with 40 Euler/simple steps, CFG 4 and denoise 1. Output
aspect follows the scaled first reference. Denoise 1 starts fresh sampling;
references provide conditioning rather than fixing the original composition.

These defaults need live reference/prompt verification on each new image or
model bundle; native node availability and successful execution do not establish
character continuity or style fidelity. Keep test prompts and reference assets
neutral. User prompts, uploads and generated images do not belong in Git.
The existing Chroma workflow documents are historical examples and are not
staged by this profile; they need their original models and approved node pack.

Override `model_downloads` with a list of `{url, filename, kind}` objects in a
reviewed tfvars file; an explicit empty list disables configured downloads.
Pass its absolute path with `-- -var-file=/absolute/path/models.tfvars` to both
`tf.plan` and `tf.rent`. Another architecture also requires matching workflows.
Public safetensors URLs must contain no credentials. Gated tokens require
separate Vault-backed wiring.

Changing the model manifest replaces the rental and removes its local disk.
The retired personal rental was destroyed. The new component credentials and
Qwen rental were subsequently deployed under separate user authorizations.
No rental or connection URL is currently active. Future connection coordinates
come from the owning Terraform outputs. Implementation,
teardown and acceptance evidence is recorded in the
[owning change](../openspec/changes/add-imagegen-rental/design.md).
