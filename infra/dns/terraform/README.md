---
title: Shared Terraform DNS roots
description: Bazel entry point for owner-local DNS deployments
---

`infra_terraform_dns` packages a Terraform root in an owner's dedicated DNS
Bazel package. Pass the owner's AL configuration and declare the local
`:dnsconfig` filegroup for `dnsconfig.json`.

```starlark
load("//infra/dns/terraform:defs.bzl", "infra_terraform_dns")

infra_terraform_dns(
    name = "tf",
    configs = [":al"],
)
```

The default root uses Cloudflare. Set `dc1 = True` when the owner also uses
RouterOS. The AL configuration must select the matching DNS credentials and
the owner's Vault HTTP backend. The macro exposes the standard `tf.*`
commands and offline `tf_tests.*` checks.

Use `infra_terraform_dns_config` to compose the owner's AppRole with a distinct
Vault backend and DNS-only credential injection selected by `tf=dns`.
Keep that configuration separate from service commands. Existing DNS-only
roots can pass their existing selectors and state configuration. Follow the
[state transfer procedure](../README.md#transfer-from-service-state) before
operating a root extracted from service state.
