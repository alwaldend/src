---
title: XCP-ng
description: XCP-ng infrastructure
tags:
  - terraform
  - xcp-ng
---

The dc1 DNS configuration maps `xcp-ng.alwaldend.com` and
`host1.xcp-ng.alwaldend.com` to `192.168.1.213`.

Terraform lives in `internal/tf` and declares the pinned
[`vatesfr/xenorchestra` provider](https://docs.xen-orchestra.com/automation/terraform-provider).
It includes empty provider and HTTP backend blocks, with no resources,
endpoints, or credentials. The provider requires Xen Orchestra connected to
the XCP-ng pool when resources are added.

`al.lua` authenticates with the `src_infra_xcp_ng` Vault AppRole and configures
the HTTP backend using
`secrets/alwaldend.com/vault1/approles/src_infra_xcp_ng/tf_backend`.
The AppRole is declared in `infra/vault/tf` and must be provisioned before
running the Terraform commands. Provider connection settings are left to the
operator.

Check Terraform formatting with:

```sh
bazel_agent bazel test //infra/xcp_ng/internal/tf:tf_tests.fmt_test
```

The formatting test does not authenticate to Vault.
