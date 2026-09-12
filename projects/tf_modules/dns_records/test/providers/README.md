---
title: DNS test providers
description: Pinned provider packages for offline Terraform tests
---

These Linux AMD64 packages provide schemas to Terraform tests. Bazel
downloads them from the publisher before test execution and verifies the
publisher's release SHA-256 checksums. Terraform uses a local filesystem
mirror with no direct registry fallback; tests do not download providers or
contact live services.

Pins match the existing DNS, ingress, and Proxmox Terraform roots:

- [Cloudflare 5.22.0 release](https://github.com/cloudflare/terraform-provider-cloudflare/releases/tag/v5.22.0)
  and its [checksums](https://github.com/cloudflare/terraform-provider-cloudflare/releases/download/v5.22.0/terraform-provider-cloudflare_5.22.0_SHA256SUMS).
- [RouterOS 1.99.1 release](https://github.com/terraform-routeros/terraform-provider-routeros/releases/tag/v1.99.1)
  and its [checksums](https://github.com/terraform-routeros/terraform-provider-routeros/releases/download/v1.99.1/terraform-provider-routeros_1.99.1_SHA256SUMS).
- [Proxmox 3.0.2-rc07 release](https://github.com/Telmate/terraform-provider-proxmox/releases/tag/v3.0.2-rc07)
  and its [checksums](https://github.com/Telmate/terraform-provider-proxmox/releases/download/v3.0.2-rc07/terraform-provider-proxmox_3.0.2-rc07_SHA256SUMS).

The targeted import regression matches the owning Proxmox provider block by
omitting `pm_api_url`. It checks both empty state and existing Proxmox resource
state with no credentials. Without `PM_API_URL`, provider validation fails
before target pruning; supplying only a loopback endpoint permits the DNS
operation. The endpoint must receive zero requests.

The test imports a built-in Terraform resource through a root map into
`module.dns`, checks the exact no-op import, applies that saved plan, and verifies
the next targeted plan has no changes. Existing Proxmox resource attributes
remain unchanged. An untargeted control must fail authentication before
contacting the endpoint. These checks cover provider selection and state
preservation; live DNS adoption retains its separate inventory and plan checks.
