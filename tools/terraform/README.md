---
title: Terraform
description: Terraform
---

## Links

- Docs: https://developer.hashicorp.com/terraform

## Features

- Terraform binaries use explicit operation suffixes. In particular, apply is
  `<name>.apply`; the unnamed `<name>` mutating alias is intentionally absent.
- Root `//:tf.*` commands delegate to `//tools/terraform/cmd/workspace:tf.*`.
  They retain the repository root as Terraform's working directory. Other
  packages default to their own directory; `terraform_binary` accepts an
  explicit `chdir` for relocated command wrappers.
