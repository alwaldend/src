---
title: Update all packages
description: Update all system packages
tags:
  - ansible_role
---

The role updates installed packages and detects reboot requirements with
`dnf needs-restarting -r` on Red Hat-family hosts and `/var/run/reboot-required`
on Debian-family hosts. DNF exit codes 0 and 1 mean no reboot and reboot required;
other exit codes fail the play instead of silently skipping maintenance.

`update_all_packages_reboot_required` records the result.
`update_all_packages_reboot_enabled` defaults to `true`; set it to `false` when
an approved deployment must report the requirement without rebooting. The
`reboot_required.yaml` task entry point performs detection without upgrading
packages or rebooting. The DNF reboot-hint command must be available on the host.
