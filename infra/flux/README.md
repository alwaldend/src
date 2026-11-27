---
title: flux
date: 2026-06-29
description: Flux cluster bootstrap and workload deployment
tags:
  - flux
  - dns
  - vault
---

- [DNS](dns/README.md) manages the existing site-local Flux records.
- [AppRole](approle/README.md) manages the Vault identity used for DNS, Xen
  Orchestra login, SSH certificates, and Terraform state.

Provision the [XCP-ng VM](tf_setup/README.md), then run the
[K3s bootstrap](ansible/README.md). Apply the separate
[Terraform workload root](tf/README.md) to install the Flux Operator and
FluxInstance.
