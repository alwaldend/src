---
title: Tf modules
description: Reusable Terraform modules for Vault, virtual machines, storage, and Vast.ai
statuses:
  - active
tags:
  - terraform
---

This project contains reusable Terraform modules for Vault configuration,
Proxmox virtual machines, backup storage, and related infrastructure. The
module collection is still in progress. The [Vast.ai team module](vastai_team/README.md)
composes managed or existing teams and project API keys while callers select
provider credentials and own secret storage.
The [base image-generation module](vastai_instance_imagegen/README.md) owns
constrained rentals and SSH-only or password-protected HTTPS access with at least
32 GiB VRAM. [Forge](vastai_instance_imagegen_forge/README.md) and
[ComfyUI](vastai_instance_imagegen_comfyui/README.md) modules supply native runtime
configuration and model paths. Creating a consumer does not authorize a paid rental.

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/tf_modules
- Docs: https://developer.hashicorp.com/terraform/language/modules

## Features

- Terraform modules
