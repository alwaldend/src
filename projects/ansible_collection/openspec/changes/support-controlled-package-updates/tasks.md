## Current implementation

- [x] Add a default-enabled immediate-upgrade flag; verify disabling it preserves native automatic configuration and skips upgrades, reboot detection, and reboot.
- [x] Keep one updater with native DNF4, DNF5, APT, and OSTree behavior.
- [x] Preserve security-update policy, no-reboot checks, and the OSTree transaction wait.
- [x] Record simplification failure cases before editing source.
- [x] Use literal automatic=true/reboot=false defaults and remove the uncalled compatibility alias.
- [x] Make automatic_updates.yaml select its backend; reduce Host Bot to one import.
- [x] Validate the exact candidate's packaging, syntax, defaults, and automatic-only selection offline.

## Evidence

Failure cases and current validation belong under
`out/host-security-updates-20261029/simplification/updater/`.
No live execution, authentication, Git operation, dependency addition, or
captured scan-runtime change is authorized to this implementation child.

On 2026-10-30, scoped packaging and Host Bot syntax passed. Pinned Ansible task
listings for Host Bot and the complete updater's automatic_updates tag selected
backend configuration without immediate upgrades or reboot. Literal defaults
parsed as automatic=true/reboot=false; strict owner OpenSpec validation passed.
The parent owns combined formatting and final runtime acceptance. No live
updater or timer reconciliation was executed for this revision.

Earlier parent maintenance exercised DNF5, APT, and OSTree on ten active hosts
with unchanged boot identities. A concurrent OSTree transaction demonstrated
the need for the retained wait. DNF4 was source-reviewed only. Earlier default
rendering and packaging receipts remain under
`out/host-security-updates-20261029/default-review/`; they predate this revision.

## Normal role configuration-only control

Review feedback is addressed by `update_all_packages_upgrade_enabled`, default
true. Host Bot declares false in owner variables and invokes the normal role.
Packaged syntax and native local check mode passed with zero changes: automatic
OSTree configuration ran while immediate upgrade, reboot-hint, and reboot tasks
were skipped, even with the fixture reboot flag true. The before/after boot
assertion passed. Evidence and source hashes are retained in
`out/host-security-updates-20261029/host-bot-updater-review/result.json`.
The check-mode verification performed no live reconciliation or reboot.

## Native simplification

- [x] Keep one public role entry point and select its backend once; retire unused direct internal task entry points.
- [x] Verify successful and failed execution, private retention, owned cleanup, and the packaged candidate.
