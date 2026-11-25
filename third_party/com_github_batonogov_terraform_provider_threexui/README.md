---
title: Terraform provider for 3x-ui
description: Patched upstream provider built from source
---

The source pin and Go dependencies match [upstream v3.34.0](https://github.com/batonogov/terraform-provider-threexui/releases/tag/v3.34.0).
Bazel builds the provider with its own Go SDK and packages it for installation from the Terraform filesystem mirror.
The dependency manifest follows the workspace Go version convention, and the provider build selects Go 1.26.11 explicitly.

The patch lets a single observatory tagged `default` use Xray's flat JSON format.
It maps that format back to the same Terraform tag on refresh and preserves legacy tagged state during migration.
