---
title: Vast.ai Terraform provider
linkTitle: Vast.ai Terraform provider
description: Manage Vast.ai compute and team resources with Terraform
layout: landing
date: 2026-11-04
# Generated from project README by //tools/repo_quality:inventory_update.
statuses:
  - active
languages:
  - go
tags:
  - terraform
---

![Terraform provider request flow](/docs/projects/vastai_terraform_provider/assets/architecture.svg)

Discover available GPUs and manage Vast.ai instances with Terraform.
The provider supports on-demand rental creation, refresh, label and power-state
updates, import, and destruction through Vast.ai's REST API.
It also manages SSH public keys, account environment variables, and custom
team roles, enabling non-compute lifecycle validation without GPU rentals.
Team and API-key resources manage separate teams and sensitive project tokens.
The infrastructure bootstrap delivers the generated test token directly to
Vault, and the non-compute test root uses the project's dedicated identity.

This development provider uses HashiCorp's Terraform Plugin Framework and an
API client generated from Vast.ai's official OpenAPI schema. It has
offline lifecycle acceptance coverage and is not yet published to the Terraform
Registry. Apply records API acceptance; connection details become available
through refresh as the instance boots.

See the [documentation](/docs/projects/vastai_terraform_provider/) for setup,
supported settings, import guidance, and current limitations.
