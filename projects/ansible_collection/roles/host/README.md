---
title: Host
description: Common host setup
tags:
  - ansible_role
---

Common setup includes `automatic_updates`, which uses native distribution
packages and systemd timers for security updates without automatic reboots.
Set `automatic_updates_enabled: false` to disable that installation schedule.
