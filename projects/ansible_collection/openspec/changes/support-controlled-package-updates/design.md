## Decisions

`update_all_packages` owns immediate updates and native automation. Its single
`update_all_packages_automatic_enabled` control defaults to true; immediate
reboots default to false. `update_all_packages_upgrade_enabled` defaults to true
and gates the immediate-upgrade/reboot block. Host Bot disables that flag in
owner variables and imports the normal role once. The main entry point selects
the backend once; task files are internal implementation. The `automatic_updates`
tag remains configuration-only.

The actual PR base `a9ac58edab4a0f98aad8b2e87b442c6d248ad887` has no separate
automatic-update role or compatibility variable. Full tracked-source research
at `17e303e6d4d8e92ac3c232bc0e72ccfb7b0849a0` found no alias callers. The
separate PR202 is not a compatibility contract for this change.

Detect `/run/ostree-booted` before DNF facts. Preserve the native transaction
wait and pinned `ansible.posix.rpm_ostree_upgrade` module. Staging does not apply
fixes live. Keep native security-update scopes, existing package-manager
bindings, no-reboot settings, effective APT policy checks, and service-hook
checks before policy changes and after possible updater installation.

## Failure cases

Record before edits under
`out/host-security-updates-20261029/simplification/updater/failure-cases.md`.
Omitted controls must enable automation without reboot; an explicit canonical
false must stop managed installation without installing a missing updater.
Automatic-only controls and tags must not select immediate updates. Unsupported
backends and update/reboot-hint failures remain failures. Package scripts may
restart services; these controls do not promise interruption-free maintenance.

## Validation boundary

Use the exact packaged candidate for offline syntax and tag inspection. The
parent owns any authorized live verification, native timer/settings inspection,
unchanged boot checks, and repeated-run convergence. Earlier live evidence
predates this simplification and does not establish acceptance for it.

## Native simplification acceptance

Keep one public role entry point and select its backend once; retire unused direct internal task entry points.
Failure cases: automatic-only execution must not upgrade packages or reboot;
unsupported backends must fail, and each invocation selects its backend once.
Keep automation enabled and reboot disabled by default. Validate without live
package changes or reboot.

## DNF opt-out across native installation timers

Disabling automation discovers installed generic DNF4, install-only DNF4, and
DNF5 timer units through native systemctl. The install-only service overrides
`apply_updates=no`, so configuration alone cannot implement opt-out. Timer
reconciliation runs after backend configuration, independently of the mutable
updater package-presence gate. Missing timers are skipped; discovery and stop
errors remain failures. Download-only and notification timers and active update
services are not stopped. Enabling retains the selected backend's existing timer.
