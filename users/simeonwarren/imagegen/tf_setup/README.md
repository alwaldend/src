---
title: Image generation provisioning
description: Generate a Vault-backed login and rent one HTTPS-protected ComfyUI instance
---

This Terraform root generates a 48-character random alphanumeric password,
stores `{username, password}` at mount `secrets`, path
`alwaldend.com/vault1/approles/user_simeonwarren/imagegen/web`, and then rents
one ComfyUI instance with that login. A failed Vault write prevents rental creation.
The backend stores state under the user AppRole's `tf_backend/imagegen` namespace.
The rental credential is injected as `VAST_API_KEY`; Vault provider authentication
is injected separately through AL. No account API key is passed to the container.

## Rent and connect

First run [infrastructure credential setup](../../../../infra/vastai/alwaldend/README.md) to create a dedicated alwaldend
team token and store it under your existing user namespace. No new AppRole or
user credential is required. The team writer's exact-entry Vault permission is
defined by the [owning Vault configuration](../../../../infra/vault/tf/README.md#vast-team-token-delivery).
The [module](../../../../projects/tf_modules/vastai_instance_imagegen_comfyui/README.md) uses
the direct Vast ComfyUI/CUDA 12.9 image tag declared in [main.tf](main.tf),
without inheriting stock template options. The root selects an explicit published release tag. Update that one image reference to select a newer build;
there is no registry lookup or helper executable during Terraform planning.
Changing the tag produces a reviewed replacement plan for an existing rental.
The rental defaults to Chroma1-HD-Flash; override public model downloads using
`model_downloads` in a task-owned tfvars file. An optional reviewed `template_hash_id` may
supply provisioning environment settings; its Docker options must be environment-only.

```sh
bazel run //users/simeonwarren/imagegen/tf_setup:tf.plan
bazel run //users/simeonwarren/imagegen/tf_setup:tf.rent
```

`rent` displays the interactive plan; review it before approving. Defaults are
32 GiB VRAM, one GPU, 100 GB disk, and a USD 1/hour ceiling including requested
disk. `max_hourly_price` and `offer_id` can be supplied explicitly. The ceiling
is not a total budget and market availability is not guaranteed.

Creation acknowledges a contract; boot and provisioning take longer. Refresh
through the same owner after boot, reviewing its interactive plan:

```sh
bazel run //users/simeonwarren/imagegen/tf_setup:tf.refresh
bazel run //users/simeonwarren/imagegen/tf_setup:tf.output
```

`https_url` remains absent until a refresh reports the public IP and mapped port.
Verify the Vast instance certificate before entering credentials; follow the
[module's certificate guidance](../../../../projects/tf_modules/vastai_instance_imagegen/README.md#certificates).
Retrieve the login through an authorized Vault client using the reference above;
do not paste it into logs or command arguments. `ssh_forward_command` provides
an alternative for account SSH keys. Export images before release:

```sh
bazel run //users/simeonwarren/imagegen/tf_setup:tf.release
```

Changing the password or access settings replaces the instance; export files
first. Vault historical versions can retain previous credentials. The root has
no automatic expiry or spend cutoff. See the [privacy boundary](../README.md#privacy-boundary).

## Offline validation

```sh
bazel test //users/simeonwarren/imagegen/tf_setup:all
```

The checks omit the live backend, validate the actual module/provider schema,
and exercise password generation and Vault persistence with a mocked Vault and
overridden rental module; the separate CLI E2E suite covers the actual rental module. They do not write Vault or rent
capacity. The provider's Terraform CLI acceptance suite emits separate SSH and
HTTPS lifecycle artifacts. Actual image startup and TLS/authentication require
a separately scoped live acceptance run.

## Download models

The default `model_downloads` in [main.tf](main.tf) selects the full
[Chroma1-Flash](https://huggingface.co/lodestones/Chroma1-Flash) checkpoint
`Chroma1-HD-Flash.safetensors`, T5 XXL FP16 encoder and Flux VAE at fixed source
revisions. This is the standalone Flash model, not the delta weights. The three
files total about 27.9 GB and fit the 100 GB disk. Vast's provisioner downloads
them inside the rented instance before ComfyUI starts serving generation; no
Hugging Face token or additional package installation is configured by this root.

In a native ComfyUI workflow, use Load Diffusion Model for the Flash model,
CLIP Loader with type `chroma` for `t5xxl_fp16.safetensors`, and VAE Loader for
`ae.safetensors`. Flash uses distilled guidance; CFG 1 with a low-step Heun
workflow is a starting point to evaluate, not a verified output-quality setting.
ComfyUI workflows select their own sampler, steps and model components; this
root does not apply Forge presets. Custom nodes and partner API nodes are
disabled. Model license terms apply.

Override `model_downloads` with a list of `{url, filename, kind}` objects in a
reviewed tfvars file; an explicit empty list disables configured downloads.
Pass its absolute path with `-- -var-file=/absolute/path/models.tfvars` to both
`tf.plan` and `tf.rent`. See the wrapper's native model directories; another
architecture also requires a matching ComfyUI workflow. Public safetensors URLs
must contain no credentials. Gated tokens need explicit Vault-backed wiring.

The previous Forge rental was destroyed at the user's request. This source
switch retains the existing Vault login and API key and does not itself create
another instance. Image boot, inference and Flash output quality have not been
verified on a new live rental.
