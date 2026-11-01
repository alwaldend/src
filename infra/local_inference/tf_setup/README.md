---
title: Local inference provisioning
description: Encrypted Yandex Cloud Broadwell V100 VM
---

This stage owns the VM, network, security group, encrypted SSD, KMS key,
Yandex address zone, and public DNS delegation/CNAME. Providers are pinned
by the repository's shared Terraform toolchain. The Vault HTTP backend,
Yandex credentials, and Cloudflare DNS credentials arrive through `../al.lua`.

Required inputs are `image_id` and `admin_ipv4_cidrs`; `folder_id` is injected.
The image must be a reviewed Ubuntu 24.04 x86_64 GPU image with a proprietary
V100-compatible NVIDIA driver at least 570.26. No image family is resolved
implicitly. This image prerequisite must be satisfied before provisioning.

Defaults are `ru-central1-a`, one non-preemptible Broadwell V100, 8 vCPUs,
96 GiB RAM, and a 200 GiB network SSD. The disk uses AES-256 KMS encryption
and annual key rotation. Key, disk, and instance have destruction guards;
the instance does not auto-delete the disk.

Ingress allows TCP 22 from administrator CIDRs and public TCP 80/443 for
Traefik's ACME validation, HTTPS redirect, and mTLS application route.
T3 Code and inference listen on loopback. Accounts retain normal outbound
networking, and no service account is attached to the VM. Public records
are consumed from `../dnsconfig.json`; the VM's dynamic address is maintained
in the delegated Yandex zone. DNS credentials need zone lookup and edit access.

Offline formatting check:

```sh
bazel_agent bazel test //infra/local_inference/tf_setup:tf_tests.fmt_test
```

After separate authorization, operators can use `:tf_setup.plan` and
`:tf_setup.apply`. These commands contact Vault and Yandex Cloud and must
not be used for offline validation. Export the `ansible_host` output as
`LOCAL_INFERENCE_HOST` for the Ansible stage.

The scoped `:dns.plan`, `:dns.show`, and `:dns.apply` commands retain this
root's backend and target `module.dns`. Apply requires a reviewed saved plan.
Required root inputs still apply; these are live operations too.
