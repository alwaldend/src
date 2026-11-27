---
title: flux
date: 2026-06-29
description: Flux cluster bootstrap and workload deployment
tags:
  - flux
  - dns
  - vault
---

[![Flux network showing browser access, Vault sign-in and certificate renewal, and Terraform and Forgejo deployment paths](docs/assets/network.svg)](docs/assets/network.svg)

The three Vault nodes show the same service from different network paths.
[Vertical network diagram](docs/assets/network-vertical.svg) ·
[Interactive 3D network diagram](docs/assets/network.3d.html)

Views with one Vault node are available in
[horizontal](docs/assets/network-single-vault.svg),
[vertical](docs/assets/network-single-vault-vertical.svg), and
[interactive 3D](docs/assets/network-single-vault.3d.html) layouts.

- [DNS](dns/README.md) manages the Flux host and Operator ingress records.
- [AppRole](approle/README.md) manages the Vault identity used for DNS, Xen
  Orchestra login, SSH certificates, and Terraform state.

Provision the [XCP-ng VM](tf_setup/README.md), then run the
[K3s bootstrap](ansible/README.md). Apply the separate
[Terraform workload root](tf/README.md) to install the Flux Operator and
FluxInstance. Follow the [Operator ingress runbook](tf/README.md#operator-ingress)
to make the UI available.
