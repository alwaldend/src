---
title: Vast.ai image generation instance module
description: Reusable AI-Dock Forge rentals accessed through SSH forwarding
---

This module declares one Vast.ai image-generation instance from an existing
AI-Dock Stable Diffusion WebUI Forge template. The caller supplies the template
hash, provider authentication, and state backend. By default it selects the
cheapest eligible on-demand offer within a USD 1/hour ceiling.
The module supplies no provider credentials, account keys, or deployment root.
Applying a consumer rents paid capacity and requires separate operator approval.

## Prepare a template

Follow Vast's [image-generation guide](https://docs.vast.ai/image-generation)
and select the recommended **SD Web UI Forge** template backed by AI-Dock.
The older [Stable Diffusion guide](https://docs.vast.ai/stable-diffusion) describes
the original AUTOMATIC1111 workflow; this module targets AI-Dock Forge's
`/opt/ai-dock/bin/init.sh` startup contract.

Edit and save your own reviewed template version before renting:

1. Remove every published-port option (`-p`, `-P`, `--publish`, etc.). Keep only
   literal `-e NAME=value` or `--env NAME=value` environment assignments. The
   module rejects other Docker options, including host networking. An explicitly
   empty options string is accepted; missing options metadata fails closed.
2. Retain the template's model provisioning settings such as
   `PROVISIONING_SCRIPT`. Review the script, models, extensions, and their
   licenses. Do not place credentials in public templates or checked-in source.
3. Copy the saved template's 32-character `hash_id` from its share link. Use an
   explicit version, rather than selecting a template by its display name. A hash
   fixes template settings, not mutable image tags or provisioning URLs; pin
   those in the reviewed template where possible.

Vast [merges template environment and port options](https://docs.vast.ai/api-reference/creating-and-using-templates-with-api#precedence-rules)
into creation requests. Merely choosing SSH does not remove template web ports.
This module validates the selected template before creation, forces
`ssh_direct`, and overrides startup using AI-Dock's
[documented initialization arguments](https://github.com/ai-dock/base-image/wiki/2.0-Environment-Variables).
It disables Cloudflare quick tunnels and the Caddy, Cloudflare, Jupyter,
service portal, auxiliary SSH, and Syncthing supervisor services. Forge listens
on remote loopback port **17860**; its public proxy port 7860 is not published.
Vast's own SSH service remains available on its mapped SSH port.

The selected image and provisioning code must be trusted: validating Docker
options does not audit application code or prevent a custom script from
starting its own outbound tunnel. Other template families are unsupported.

## Compose the module

```hcl
provider "vastai" {} # Inject VAST_API_KEY through the caller's Vault configuration.

module "imagegen" {
  source = "../../../projects/tf_modules/vastai_instance_imagegen"

  template_hash_id = var.imagegen_template_hash_id
  disk_gb          = 100
  max_hourly_price = 1 # USD/hour ceiling; does not authorize a rental.
}

output "imagegen_forward" {
  value = module.imagegen.ssh_forward_command
}
```

Package the module through
`//projects/tf_modules/vastai_instance_imagegen` in the consuming Terraform
target's `data`, and select `//projects/vastai_terraform_provider:provider`.
No module-local backend, live plan, or apply target is supplied.

The guide recommends at least 8 GB VRAM and 50 GB storage, with 12–24 GB VRAM
and 100 GB storage for larger workloads. Check GPU driver/image compatibility,
availability, and the full rental/storage price before requesting approval.
Automatic selection requires a datacenter, verified host, one GPU with at least
24 GiB VRAM, 99% reported reliability, a directly mapped SSH port, and enough
disk space. Results are ordered by API-reported total hourly price with the
requested disk allocation. `min_gpu_ram_gb`, `disk_gb`, `max_hourly_price`, and
optional `gpu_name` configure the workload and price limits. Optional `offer_id`
restricts the search to a particular offer while retaining every constraint,
including the datacenter requirement. The hourly limit is not a total budget;
network transfer and continuing storage can incur additional charges.

Marketplace changes and a changed `offer_id` do not replace an existing contract.
Template, image, storage, or startup changes do require replacement and select
the currently matching offer. Label and requested power state update in place.
Each normal plan rechecks marketplace filters; an empty eligible market can
block planning even when an existing contract is retained.

## Connect and retain outputs

Configure your account's SSH public key before renting and keep the private key
in your SSH agent or local SSH configuration. Team-account key restrictions
remain the caller's responsibility. The module does not manage account keys.

Creation acknowledges the rental request; it does not wait for SSH or model
downloads. After an approved deployment, use the caller's owning refresh
workflow to observe `ssh_host` and `ssh_port`. `ssh_forward_command` is null
until both are reported. Run that output locally, then open `local_url`
(default `http://127.0.0.1:7860`). It forwards local loopback to remote loopback
over SSH, using Vast's documented [local forwarding](https://docs.vast.ai/guides/instances/connect/ssh#ssh-local-port-forwarding).

Use SCP or SFTP over the same SSH connection to upload models and download
generated images. Export wanted files before destroying an instance: destruction
deletes its local storage. Stopping an instance still incurs storage charges.

## Validation

```sh
bazel test //projects/tf_modules/vastai_instance_imagegen:all
bazel test //projects/vastai_terraform_provider/test/acceptance:acceptance_test
```

The acceptance suite exercises this actual module through the pinned Terraform
CLI and provider against a loopback API fixture, including unsafe templates,
template replacement, pending connection details, and SSH forwarding. It emits
`imagegen-lifecycle.json` in Bazel test undeclared outputs. It never rents GPU
capacity. Image startup, model provisioning, actual port exposure, and image
generation require a separately approved live acceptance run and remain unverified.
