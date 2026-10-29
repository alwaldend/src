---
title: Host
description: Common host setup
tags:
  - ansible_role
---

Common setup invokes `update_all_packages` once for immediate package updates
and native automatic-update configuration. Set
`update_all_packages_automatic_enabled: true`
to enable native security-update timers, or staged image updates on OSTree.
Automatic reboots are disabled. Immediate reboots also default to disabled;
`update_all_packages_reboot_enabled: true` explicitly permits a required reboot.
