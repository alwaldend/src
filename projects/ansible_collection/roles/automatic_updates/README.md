---
title: Automatic updates
description: Native operating-system security updates with packaged systemd timers
tags:
  - ansible_role
---

The shared `host` role enables automatic security updates on mutable DNF5,
DNF4, and APT hosts. It installs `dnf5-plugin-automatic`, `dnf-automatic`, or
`unattended-upgrades` and enables the distribution's existing systemd timer.
No cron job or custom timer is created. The packaged schedule is retained.

DNF applies packages with security advisories. APT retains its packaged allowed
origins and refreshes package metadata daily. Automatic reboots are disabled;
kernel updates and services requiring a restart still need maintenance.
Manually installed software such as the pinned Traefik binary is outside the
OS package updater.

Set `automatic_updates_enabled: false` to disable the managed installation
timer and update setting without removing packages. Unsupported package
managers and immutable OSTree hosts are skipped with a message. Apply only
this configuration with the `automatic_updates` tag.

Inspect the native timer and corresponding service journal to observe update
execution. Enabling the timer verifies scheduling, not that a future update
transaction has succeeded.
