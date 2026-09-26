---
title: XCP-ng
description: Protected local static-hosting VM
---

Provisions one Fedora 44 VM in `src_infra_download`, using named inventory
from the existing XCP-ng deployment: two CPUs, 2 GiB RAM, a 20 GiB boot disk,
and a 100 GiB content disk. The local address comes from the apex DNS owner, whose local A record is
also the target of the unique host CNAME.
The cloud template must have one boot disk and the expected `enX0` interface;
Ansible formats and mounts the second disk at `/srv/download` without force.

`prevent_destroy` rejects replacements and destruction because the pinned XO
provider owns disks with the VM. There is no automatic retained-disk replacement
workflow on XCP-ng. Do not remove that guard as a routine upgrade step.

`tf.plan` and `tf.apply` are live operations. `tf_tests.fmt_test` and the
`offline` entry point with backend-disabled initialization validate source.
