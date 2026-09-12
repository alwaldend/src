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
- `terraform_target_binary_map` exposes `<name>.plan`, `<name>.show`, and
  `<name>.apply` for a target in an existing root. Plan includes its declared
  `-target`; apply accepts only one saved plan file, which must be reviewed
  before execution. The caller retains its root, backend, and AL configuration.
  The [runner's saved-plan guard](runner/README.md) prevents this apply target
  from starting a new plan for services elsewhere in the root.
