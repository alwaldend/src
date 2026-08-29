---
title: lvg
description: Setup LVM volume groups
tags:
  - ansible_role
---

`lvg_volume_groups` describes the complete desired physical-volume list for
each volume group. Extra physical volumes are not removed unless a caller sets
`remove_extra_pvs: true`.

Creating a volume group initializes every listed device as an LVM physical
volume. Callers must resolve stable device paths and verify that each device is
the intended, unused disk before applying this role; `pvcreate` may overwrite
an existing disk signature.
