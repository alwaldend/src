---
title: Image generation provisioning
description: Generate a Vault-backed login and rent one HTTPS-protected Forge instance
---

This Terraform root generates a 48-character random alphanumeric password,
stores `{username, password}` at mount `secrets`, path
`alwaldend.com/vault1/approles/user_simeonwarren/imagegen/web`, and then rents
one Forge instance with that login. A failed Vault write prevents rental creation.
The backend stores state under the user AppRole's `tf_backend/imagegen` namespace.
The rental credential is injected as `VAST_API_KEY`; Vault provider authentication
is injected separately through AL. No account API key is passed to the container.

## Rent and connect

First run [infrastructure credential setup](../../../../infra/vastai/alwaldend/README.md) to create a dedicated alwaldend
team token and store it under your existing user namespace. No new AppRole or
user credential is required. The team writer's exact-entry Vault permission is
defined by the [owning Vault configuration](../../../../infra/vault/tf/README.md#vast-team-token-delivery).
The [module](../../../../projects/tf_modules/vastai_instance_imagegen/README.md) uses
the newest published Vast Forge Neo/CUDA 12.9 tag without inheriting stock template options.
There is no literal `latest` tag upstream, so a bounded read-only Docker Hub lookup
runs at plan time. A new tag may produce a reviewed replacement plan.
The rental defaults to Juggernaut XL v9; override public model downloads using
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

The default `model_downloads` in [main.tf](main.tf) selects the publisher's
[Juggernaut XL v9](https://huggingface.co/RunDiffusion/Juggernaut-XL-v9)
standalone SDXL checkpoint at a fixed source revision. Its 7.1 GB download fits
the configured 100 GB disk. The VAE is included. The existing native provisioner
downloads it before Forge starts; no Hugging Face token is required.
The publisher describes v9 as its most popular SDXL model. Model license terms
still apply, including restrictions on paid API services.

To replace the default, supply a `model_downloads` list via Terraform variables.
For example, replace the placeholder URL with a public model file you have
reviewed and may use:

```hcl
model_downloads = [{
  url      = "https://huggingface.co/OWNER/MODEL/resolve/REVISION/model.safetensors"
  filename = "model.safetensors"
  kind     = "checkpoint"
}]
```

Pass the absolute path of your variable file with `-- -var-file=/absolute/path/models.tfvars`
to both `tf.plan` and `tf.rent`. The native provisioner downloads the configured
files on first boot before Forge starts. An explicit empty list disables the
default model download. Review model
license, disk requirements, and download charges; only public URLs are currently
supported. Gated tokens must be wired from Vault rather than embedded in a URL.
