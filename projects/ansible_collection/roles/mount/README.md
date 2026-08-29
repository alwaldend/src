---
title: Mount
description: Setup posix mounts
tags:
  - ansible_role
---

Mount entries can set `opts`, `dump`, and `passno`; `state: mounted` keeps the
mount active and writes its persistent `/etc/fstab` entry.
