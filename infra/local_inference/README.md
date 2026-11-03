---
title: Local inference
description: Independent GPU provisioning and DNS for local inference
---

This component owns the canonical `dnsconfig.json`, consumed by Vault's
exact-hostname PKI role and the independent [DNS Terraform root](dns/README.md).
DNS has its own state and Cloudflare credentials.

The [provisioning stage](tf_setup/README.md) declares the A100 VM, encrypted
storage, KMS key, network and Yandex child DNS zone. Its AL configuration
selects only Yandex credentials and the provisioning state backend. Host
configuration is prepared separately in [PR #233](https://github.com/alwaldend/src/pull/233).

Source declarations do not establish deployment or service health. The
[provisioning change](openspec/changes/extract-local-inference-provisioning/)
records source acceptance; live deployment requires separate authorization.
