---
title: Terraform dependencies
description: Shared provider archive declarations for repository Terraform targets
---

`include.MODULE.bazel` owns the repository's Terraform provider versions,
immutable release URLs, and SHA256 integrity. Consumers select the generated
provider labels through the public
[rules_terraform](../../tools/rules_terraform/README.md) rules. One version is
selected per provider source; the rule module's resolution boundary can support
an explicit override policy later.

The pins were verified against the publishers' release checksums listed by the
[Terraform registry protocol](https://developer.hashicorp.com/terraform/internals/provider-registry-protocol).
Provider downloads use Bazel's verified downloader. Terraform installs only
from the selected runfiles mirror and writes provider locks in runtime workspaces.

The offline fixture in `test/` checks installation with the declared Local
provider, failure without it, and isolation from inherited registry/cache settings.
Updating a pin also requires checking consuming Terraform constraints and these
provider integration tests.
