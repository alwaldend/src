---
title: Local inference provisioning
description: Encrypted Yandex Cloud Cascade Lake V100 VM
---

This stage owns the VM, network, security group, encrypted SSD, KMS key,
and the delegated Yandex DNS zone. Providers are pinned
by the repository's shared Terraform toolchain. The Vault HTTP backend
and Yandex credentials arrive through this stage's `BUILD.bazel`.

The component AppRole is owned by `infra/vault/tf`. Prepare the Yandex folder
and service-account credentials through `infra/yandex_cloud/org1/tf` after
that identity is available. These are separately authorized prerequisites.

Required inputs are `image_id` and `admin_ipv4_cidrs`; `folder_id` is injected.
The authorized provisioning uses Fedora 43 image `fd8iku26nnkveh8s4di1`, already
used by `infra/download` and `infra/threexui`. No image family is resolved
implicitly. This standard OS image does not establish NVIDIA driver readiness;
GPU driver installation and Fedora host configuration remain host-stage work.
The current host candidate in PR #233 still requires Ubuntu and must be adapted
before it can configure this Fedora VM.

Defaults are `ru-central1-a`, one preemptible A100 with 80 GB VRAM, 28 vCPUs,
119 GiB RAM, and a 200 GiB network SSD. The disk uses AES-256 KMS encryption
and annual key rotation. Key, disk, and instance have destruction guards;
the instance does not auto-delete the disk. The boot disk explicitly uses legacy
BIOS with PCI topology V1, matching the last applied compatibility trial.

Ingress allows TCP 22 from administrator CIDRs and public TCP 80/443 for
Traefik's ACME validation, HTTPS redirect, and mTLS application route.
T3 Code and inference listen on loopback. Accounts retain normal outbound
networking, and no service account is attached to the VM. Zone and hostname
inputs come from `../dnsconfig.json`; the VM's dynamic address is maintained
in the delegated Yandex zone.

Offline formatting check:

```sh
bazel test //infra/local_inference/tf_setup:tf_tests.fmt_test
```

After separate authorization, operators can use `:tf_setup.plan` and
`:tf_setup.apply`. These commands contact Vault and Yandex Cloud and must
not be used for offline validation. Export the `ansible_host` output as
`LOCAL_INFERENCE_HOST` for the Ansible stage.

Public delegation and CNAME ownership are extracted into the independent
`infra/local_inference/dns` Terraform root and backend. This provisioning root
contains no Cloudflare resources, credentials, or DNS-only command variants.
Its Yandex child zone and VM address still consume the canonical component
DNS document. DNS has independent ownership and state. Host configuration is delivered
separately in [PR #233](https://github.com/alwaldend/src/pull/233).

The operator explicitly selected public SSH for this deployment using
`admin_ipv4_cidrs = ["0.0.0.0/0"]`; cloud-init disables password authentication
and root login. Other deployments can supply narrower IPv4 CIDRs.

Creating a VM starts it automatically. This provider does not manage a desired
power state; stopping the VM after provisioning is a separate control-plane
operation. The authorized deployment's final state must be verified explicitly.

On 2026-11-03, the approved folder/account prerequisite completed with nine
additions. The provisioning apply created the six supporting resources,
including the encrypted disk, but Yandex rejected VM creation because the
memory quota was exceeded in `ru-central1-a`, with one further unspecified
quota violation. No VM was created or booted.

The operator then requested another availability zone and authorized replacement
of unused component resources. The `ru-central1-b` attempt replaced the subnet
and deleted the old unused disk, but Yandex rejected the new disk with
`ResourceExhausted: Resource allocation is restricted`. That left five resources:
the subnet, network, security group, KMS key and child DNS zone. The boot disk
and VM were absent; disk destruction protection was restored.

The indexed incident listing described zone b as unavailable. A later live
zone-status read reports b `UP`, while the cloud's b CPU, RAM and SSD allocation
quotas remain zero. The [platform support table](https://yandex.cloud/en/docs/compute/concepts/vm-platforms)
does not support either V100 platform in zones d/e.

The operator then selected a return to `ru-central1-a` with a smaller VM.
`gpu-standard-v2` retains one 32 GB V100 and eight vCPUs while reducing host RAM
from 96 GiB to 48 GiB. Provisioning remains subject to cloud quotas and capacity.

The smaller attempt created the disk and replaced the subnet in `ru-central1-a`,
but VM creation still exceeded the memory quota, with one additional unspecified
quota violation. No VM was created or booted. A subsequent read with the cloud
management identity verified the correct target cloud and showed sufficient
published RAM, CPU, GPU and dynamic-IP headroom for the 48 GiB configuration.
Those later values do not establish the second earlier violation or prove that
a quota increase was necessary. The component identity cannot read cloud-level
quotas; the cloud-management identity can.

After the operator reported increased quota, a fresh 48 GiB retry passed the
initial quota check but failed during allocation with `Not enough resources`.
Operation `fhmvss42r60r4mrqmcp7` completed with error code 8 and no further
details; no VM was created or booted. The operator then authorized restoring
the original 96 GiB Broadwell configuration after increasing RAM quota.

That original-size retry also failed allocation with `Not enough resources`,
in operation `fhmk1pqi4vth1k6f4oj6`. Both V100 configurations have therefore
failed their tested regular-capacity requests in zone a. Six supporting
resources remain, including the encrypted disk; no VM exists. Further attempts
need a changed allocation condition or a selected alternative.

A preemptible V100 retains the requested GPU and VM size but may be stopped
by Yandex at any time, with a maximum 24-hour run. The larger Gen2 option in
zone d provides 80 GB VRAM and requires at least 144 GiB host RAM, additional
quota and a compatible UEFI image. Those alternatives require an operator
choice before changing the current regular V100 configuration.

The operator selected the preemptible V100 alternative after both regular
configurations failed allocation. This keeps the original GPU, eight vCPUs,
96 GiB RAM and encrypted disk while setting `preemptible = true`. Yandex may
stop it at any time and each run lasts at most 24 hours. The initial creation
must still be followed immediately by a verified stop. Preemptible mode does
not establish available capacity or GPU-driver readiness.

The 96 GiB preemptible attempt also failed allocation, in operation
`fhmp4c79geskjtismueq`. The next retry uses the earlier approved 48 GiB
Cascade Lake shape with the same preemptible policy, GPU capacity and disk.

The 48 GiB preemptible attempt also failed with `Not enough resources`, in
operation `fhmljpjc6622s8uo7jlk`. No VM exists. The configuration remains
preemptible Cascade Lake V100 in zone a, ready for a fresh authorized retry
after allocation conditions change. Kazakhstan's region-specific platform
limits list only the 24 GB T4i GPU, below the requested 32 GB VRAM; no
Kazakhstan VM was attempted.

The subsequent larger-GPU trials tested A100 80 GB in zone a and Gen2 80 GB
in zone d. After the operator increased zone-a CPU quota, A100 requests
reached allocation but failed with `Not enough resources` in both regular
and preemptible modes, including with legacy PCI topology V1. The Gen2
attempt was rejected by zone-d memory quota plus another unspecified quota;
its published GPU quota was also zero. These errors do not identify which
physical resource was unavailable.

The current candidate retains the last attempted preemptible A100 shape in
zone a. No VM was created or booted. The six supporting resources remain;
all three destruction guards are restored. The
[change evidence](../openspec/changes/extract-local-inference-provisioning/evidence.md)
records the exact configurations and operation IDs.
