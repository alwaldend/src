---
title: Local inference
description: Canonical local inference DNS declarations and independent DNS state
---

This component owns the canonical `dnsconfig.json`, consumed by Vault's
exact-hostname PKI role and the independent [DNS Terraform root](dns/README.md).
The DNS root has its own state and Cloudflare credentials; GPU provisioning and
host deployment are prepared separately in [PR #233](https://github.com/alwaldend/src/pull/233).

The [change record](openspec/changes/separate-dns-terraform-root/) records
implementation and verified DNS ownership. Source declarations alone
do not establish deployed infrastructure or service health.
