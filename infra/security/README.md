---
title: Host security maintenance
description: Repeatable security scans and package updates without automatic reboot
---

![Security scan execution and reports](assets/security-scans.svg)

Run and copy label Ansible actions; report data returns to the controller.

The security playbook composes the existing owner inventories, runs the pinned
Lynis audit, updates all installed packages, and enables native automatic
security updates. OSTree hosts stage complete deployments for a later boot.
Immediate and automatic reboots are disabled. The playbook verifies unchanged
boot identities; it does not claim that an audit proves absence of compromise.

The default host pattern excludes the legacy Vault VM group and the XCP-ng/XOA
appliance groups. Access failures are reported and excluded from execution;
failed access is not evidence that a machine is inactive. The inventory mapping
contains source labels only. Correct host membership in the owning inventory.
Owner deployment group variables are intentionally not imported: this workflow
needs only inventory identities and explicit connection overrides.

Create a private controller report directory before running:

```sh
umask 077
mkdir -p out/security/reports out/security/tmp
TMPDIR="$PWD/out/security/tmp" bazel_agent bazel run //infra/security/ansible:ansible -- \
  --extra-vars "lynis_audit_report_directory=$PWD/out/security/reports"
```

Use `--limit` for a subset. The `security_audit` tag performs the audit only;
`update_all_packages` updates packages and reconciles automatic updates after
separately reviewed audit evidence. Full execution performs both, sequentially
per host. Audit-only runs may set `security_serial=3`. An explicit additional
inventory may supply connection overrides, including a local connection for
Host Bot when running on that exact machine. Never apply a fleet-wide local
connection override.

Authentication uses the existing operator AppRole, temporary signed SSH keys,
and the public SSH host CA. It does not inject unrelated deployment secrets.
The SSH target uses the same boundary:

```sh
bazel_agent bazel run //infra/security:ssh -- ansible@host1.ingress.alwaldend.com uname -r
```

Lynis reports and logs are private operational artifacts. Review warnings before
maintenance; the role never applies scanner hardening suggestions. Distribution
updates can restart services despite preserving the host boot. Installed kernel
fixes can remain inactive until a separately authorized reboot.

## ClamAV scans

`//infra/security/ansible:clamav` runs the shared native streaming scan role
against the existing fleet inventories. It is independent of the maintenance
playbook. `//infra/security/appliances:clamav` uses the appliance authentication
boundary and defaults to XOA only. XCP-ng remains unscanned because the official
native payload requires a newer glibc; its Lynis and update workflows are unchanged.

Both callers import one shared three-play lifecycle: prepare the controller,
scan selected hosts, and clean up before reporting accumulated failures. The
thin callers retain their separate host defaults and connection variables;
packaging retains each inventory and authentication boundary.

The controller uses the pinned official ClamAV engine in a temporary private
directory and a bounded transient systemd user unit. Endpoints stream files to
its private Unix socket through the shared role's SSH transport; they do not
load the signature database. Scans install no persistent service and perform
no reboot, deletion, or quarantine.

Startup verifies native signature metadata and daily freshness, then uses one
bounded daemon VERSION polling task to require the verified daily version and
a server response shape. The client's local-version fallback is rejected.
`signatures.json` owns signature metadata; lifecycle evidence records the
actual server version without duplicating it in a separate `engine.json`.

Remote scans require reverse Unix-socket forwarding permission for the SSH
account. The workflow preserves server SSH policy and stops before scanning
when forwarding is denied. Review any required exception through the owning
SSH configuration; the scanner does not relax login or forwarding restrictions.

Create private report and temporary directories, and retain one signature cache
across runs. For example:

```sh
umask 077
mkdir -p out/security/clamav/reports out/security/clamav/database out/security/clamav/tmp
TMPDIR="$PWD/out/security/clamav/tmp" \
ANSIBLE_LOCAL_TEMP="$PWD/out/security/clamav/tmp" \
  bazel_agent bazel run //infra/security/ansible:clamav -- \
  --extra-vars "clamav_scan_report_directory=$PWD/out/security/clamav/reports" \
  --extra-vars "clamav_database_directory=$PWD/out/security/clamav/database"
```

Run one controller workflow at a time. Preserve the cache, including cooldown
state; do not create a new cache to retry a download failure. Every `--limit`
must include localhost, for example `--limit 'localhost,host1.ingress.alwaldend.com'`
or `--limit 'localhost,xoa'`. For a local scan of Host Bot, explicitly exclude
the entire task scratch directory through `clamav_scan_exclude_paths`, including
controller files, signatures, and reports. Only Host Bot should receive a local
connection override.

The controller must have a working systemd user manager, at least 4 GiB available
RAM, and at least 5 GiB free disk. Its unit runs at Nice 19 with idle I/O priority
and a 6 GiB memory ceiling. `clamav_controller_timeout_seconds` bounds its total
lifetime (default 172800); `clamav_scan_timeout_seconds` bounds each endpoint
(default 86400). If Ansible is interrupted, the unit can remain until that
maximum runtime; the recorded unique unit can also be stopped explicitly.

Host failures are retained while later selected batches continue. The final play
attempts controller cleanup before reporting accumulated failures; host-failure
reporting still runs if cleanup fails. Inspect native scan logs and each role
receipt: completion, detections, read errors, exclusions, and
size/parser limits are distinct outcomes. Successful execution is not proof of
complete coverage or absence of compromise. The
[controller documentation](clamav/README.md) describes signature verification,
engine limits, lifetime, and private evidence.

## SSH policy

The [shared SSH role](../../projects/ansible_collection/roles/ssh/README.md)
supplies the forwarding default for `ansible`. Apply SSH configuration through
the owning host deployment with `--tags ssh` and `--limit`; this is normal SSH
setup, including certificate handling and the upstream restart handler.
The completed one-time maintenance playbook and helper are retired. Their
historical rollout evidence remains in the change's validation record; normal
owner deployment does not expose that retired reload-only or rollback workflow.
The scanner does not apply server configuration.
