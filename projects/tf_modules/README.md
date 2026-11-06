---
title: Tf modules
description: Reusable Terraform modules for Vault, virtual machines, storage, and Vast.ai
statuses:
  - in_progress
tags:
  - terraform
---

This project contains reusable Terraform modules for Vault configuration,
Proxmox virtual machines, backup storage, and related infrastructure. The
module collection is still in progress. The [Vast.ai team module](vastai_team/README.md)
composes managed or existing teams and project API keys while callers select
provider credentials and own secret storage.
The [Vast.ai image-generation module](vastai_instance_imagegen/README.md)
uses a caller-selected AI-Dock Forge template with SSH-only access and local
port forwarding. Creating a consumer does not authorize a paid rental.

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/tf_modules
- Docs: https://developer.hashicorp.com/terraform/language/modules

## Features

- Terraform modules
