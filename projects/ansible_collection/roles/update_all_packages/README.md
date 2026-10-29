---
title: Update all packages
description: Update all system packages
tags:
  - ansible_role
---

This role owns immediate package updates and native automatic updates. The
shared `host` role imports it once.

Immediate updates refresh metadata and update installed packages on modern
systemd DNF4, DNF5 and APT hosts. APT refuses transactions requiring package
removal. Holds, exclusions, repository choices and version locks remain in
effect. Package-manager Python bindings must already be installed; automatic
binding installation is disabled for APT and DNF5. Legacy appliances such as
XCP-ng require their owning maintenance workflow.

The `/run/ostree-booted` marker takes precedence over DNF facts. OSTree upgrades
use the pinned `ansible.posix.rpm_ostree_upgrade` module to prepare a deployment
for the next boot, without live application. Prepared fixes are not necessarily
active until a later authorized reboot.
An already running native staging transaction is allowed to finish before the
immediate upgrade starts; the role waits for up to 30 minutes and never cancels it.

## Controls and tags

- `update_all_packages_upgrade_enabled` defaults to `true`. Set it to `false`
  to use the normal role for automatic-update configuration only; immediate
  upgrades, reboot detection, and the explicit reboot task are skipped.
- `update_all_packages_reboot_enabled: false` is the default. A reboot occurs
  only when explicitly enabled and the backend reports it is required.
- `update_all_packages_automatic_enabled` defaults to `true`. Set it to `false`
  to disable native automatic updates. It never enables automatic reboots.
- The `update_all_packages` tag runs the complete role.
- The `automatic_updates` tag configures native automation without an immediate
  upgrade or reboot.

Use the normal role entry point with these controls. Backend selection happens
once; task files are internal implementation and are not standalone APIs.

Reboot detection uses `dnf needs-restarting -r`, `dnf5 needs-restarting`, the
Debian `/var/run/reboot-required` marker, or the first OSTree deployment's boot
state. `update_all_packages_reboot_required` records the result. DNF exit codes
0 and 1 mean no reboot and reboot required; other failures remain failures.
The native DNF reboot-hint command must already be installed.

## Native automatic updates

The role uses native packages (`dnf-automatic`,
`dnf5-plugin-automatic`, `unattended-upgrades`), configuration paths and systemd
timers. DNF selects security advisories; APT retains its packaged allowed
origins. No custom timer or cron job is created. Other vendor settings remain
intact. Disabling automatic updates stops the managed installation timer and
does not install a missing updater package.

OSTree uses `AutomaticUpdatePolicy=stage` and
`rpm-ostreed-automatic.timer`, retaining other daemon settings. Changed policy
is reloaded before the timer is enabled. Disabling sets policy `none` and stops
the timer. This backend stages complete deployments, not security-only updates.

DNF `reboot=never` and APT `Automatic-Reboot=false` are explicit. Before enabling
timers, the role inspects native service commands and systemd actions for
separate reboot/shutdown hooks, including Cockpit overrides, and fails with an
actionable message when found. It does not replace unrelated administrator
overrides or cancel a previously scheduled reboot. Arbitrary external automation
still requires a host audit. Package scripts may restart services even when
host reboots are disabled.

Inspect the native timer state, effective configuration, update service journal,
pending packages/deployments and unchanged-run results during E2E verification.
Enabling a timer proves scheduling, not a successful future transaction.
