---
title: Tf setup
description: Terraform setup
tags:
  - terraform
  - pve
---

This stage provisions Vault and reads network settings from [dnsconfig.json](../dns/dnsconfig.json).
Manage records through the dedicated [DNS root](../dns/README.md).

Offline formatting is available through
`//infra/vault/tf_setup:tf_tests.fmt_test`.
