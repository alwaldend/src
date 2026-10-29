---
title: Appliance security maintenance
description: Private XCP-ng and XOA audits with separately selected no-reboot updates
---

The default target audits the XCP-ng and XOA groups from
`//infra/xcp_ng/ansible:inventory`. It reuses the owner's `xcpng_host` and
`xcp_ng_ansible` credential selectors. SSH uses the controller's existing
`~/.ssh/known_hosts`, requires an ED25519 match, and never adds keys. The
Host Bot configuration owns approved key provisioning and provenance.

XOA's existing firewall limits new SSH connections, so its transport reuses
one authenticated connection. The controller control socket uses a per-run
mode 0700 `/tmp/appliance-ssh.*` directory: the workspace's absolute path
exceeds Linux's 108-byte Unix-socket limit. The playbook closes the connection
and removes this temporary directory. An unreachable host can prevent Ansible's
cleanup; `ssh-control-directory.txt` records the exact controller path for
operator cleanup. Audit evidence remains in the requested workspace directory.

Both appliances use `alwaldend.main.lynis_audit` with an already installed
Ansible-compatible interpreter. XCP-ng's host-level interpreter declaration
and earlier manual bootstrap exception belong to the
[owner inventory](../../xcp_ng/ansible/README.md). Ansible reports runtime
failures directly; the audit installs no interpreter or scanner package.

## Audit

Create an existing controller-owned mode 0700 report directory. The default
audit uses the pinned Lynis archive:

```sh
umask 077
mkdir -p out/security/appliances/tmp
TMPDIR="$PWD/out/security/appliances/tmp" \
ANSIBLE_LOCAL_TEMP="$PWD/out/security/appliances/tmp" \
  bazel_agent bazel run //infra/security/appliances:ansible -- \
  --extra-vars "appliance_report_directory=$PWD/out/security/appliances"
```

Use `--limit xcpng` or `--limit xoa` for one appliance group. Use separate report
directories to compare runs. Each host gets a unique private Lynis result
directory containing scanner reports and the caller's boot identities.
Files are mode 0600. Scanner exit codes
0 and 78 mean the scan completed; 78 includes warnings. They are not a
finding that a host is uncompromised.

`:ansible.status` performs only a boot-ID and uptime
probe with one SSH connection per appliance and retains connection failures.

The [shared Lynis role](../../../projects/ansible_collection/roles/lynis_audit/README.md)
owns archive extraction, the audit profile, private evidence collection, failure
handling, and remote scratch cleanup for both appliances. Its unique result
directory is directly under `appliance_report_directory`; there is no second
appliance audit lifecycle. Set `lynis_audit_timeout_seconds` to change the
shared role's default bound. The caller records before/after boot identities
and the XOA connection recovery path there after the role returns, including
scanner failures. An interruption can prevent these final observations; it
does not establish an unchanged boot or a completed audit. Raw status and boot
checks remain available without the audit's Python runtime.

The local-only `:syntax` target checks the playbooks without AL
authentication or a host connection. Building `:ansible_bin` checks packaging;
neither is a substitute for inspecting a successful live audit's evidence.

## Package updates

Auditing never selects an update playbook. The separate
`:ansible.update_xcpng` and `:ansible.update_xoa` targets require
`appliance_report_directory` and explicit operator authorization for the
selected mutation. Both preserve and compare boot IDs.

The XCP-ng path permits only a single-host XCP-ng 8.3 pool coordinator with
HA off, no active XAPI tasks, no XOSTOR, no attached CD media, and no
unreviewed yum-cron. It records the host's enabled state, temporarily disables
new VM placement, and runs `yum update` using only the stable
`xcp-ng-base` and `xcp-ng-updates` repositories with signature/TLS verification
intact. After successful package installation it invokes the vendor
`xe-toolstack-restart`, waits for XAPI readiness, and restores the original
enabled state. Existing guests are not evacuated, stopped, or rebooted.
Restoration is also attempted after failure when XAPI responds; failure to
recover XAPI is reported rather than claiming restoration. Kernel and
hypervisor fixes can remain inactive until a separately authorized host reboot.

This vendor yum workflow remains on raw transport. A newer Ansible interpreter
does not provide the vendor Python's yum module bindings or replace its
maintenance procedure.

If the selected download mirror fails, set
`appliance_xcpng_use_origin_mirror=true` to select the vendor's documented
`https://updates.xcp-ng.org` origin for this transaction. The override applies
only to the same stable 8.3 repositories and preserves signature and TLS
verification; it does not change installed repository configuration. See the
[vendor mirror documentation](https://docs.xcp-ng.org/project/mirrors/).

The [XCP-ng vendor procedure](https://docs.xcp-ng.org/management/updates/)
documents the disable, yum update, control-plane restart, and enable sequence
without a host reboot. It requires pool coordinator ordering and HA precautions.
Its Rolling Pool
Update procedure reboots hosts and therefore cannot satisfy a no-reboot run.
The [CLI reference](https://docs.xcp-ng.org/appendix/cli_reference/)
directs package maintenance to yum instead of `xe update-*`. Multi-host and
XOSTOR maintenance require a separately reviewed coordinated workflow.

The XOA path verifies Debian and invokes the shared package-update role using
the existing interpreter and APT bindings. Missing bindings fail rather than
being installed implicitly. The role installs the
authorized standard `unattended-upgrades` component if needed and configures
native OS automatic updates with immediate and automatic reboots disabled.

The [XOA application updater](https://docs.xen-orchestra.com/updater/) is
separate from Debian OS maintenance. Its installed version and unattended
reboot semantics must be verified before application automation is changed.
This workflow does not claim or enable application automatic updates. Native
unattended XCP-ng updates meeting the pool and no-reboot constraints have not
been established; `appliance_automatic_updates=true` fails before mutation.

## ClamAV scans

The independent `:clamav` target reuses the appliance inventory and authentication
but defaults to XOA only. XCP-ng remains unscanned because the official native
ClamAV payload requires a newer glibc. This does not change either appliance's
Lynis or package-update workflow.

Follow the [ClamAV invocation instructions](../README.md#clamav-scans), selecting
`//infra/security/appliances:clamav` and including localhost in any limit, such
as `--limit 'localhost,xoa'`. The shared native role streams endpoint data to
the private controller; its signature cache and bounded transient unit are
owned by the [controller package](../clamav/README.md). `:clamav_syntax` checks
this playbook locally without authentication or live scanning.
