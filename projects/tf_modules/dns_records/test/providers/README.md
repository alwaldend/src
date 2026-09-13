---
title: DNS test providers
description: Pinned provider packages for offline Terraform tests
---

The [test target](../BUILD.bazel) uses the public
[rules_terraform](../../../../../tools/rules_terraform/README.md) wrapper with
the same Cloudflare, RouterOS, and Proxmox provider labels as infrastructure
targets. The shared extension owns provider versions, release URLs, and
checksums. Bazel fetches the declared archives before tests execute; Terraform
uses the wrapper's packaged filesystem mirror without a registry fallback.

Each regression runs Terraform in a caller-owned temporary module directory.
Tests retain an explicit environment without inherited provider credentials;
provider installation and mocked lifecycles require no network access.

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
