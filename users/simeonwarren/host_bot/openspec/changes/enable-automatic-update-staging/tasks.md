## 1. Desired state

- [x] Ensure a private SSH directory exists before adding approved host keys; verify a fresh home and idempotent preservation of unrelated entries.
- [x] 1.7 Use the normal updater role with the owner upgrade flag disabled; validate native configuration-only selection and unchanged boot identity in check mode.
- [x] 1.1 Declare native automatic staging and prohibit automatic reboot.
- [x] 1.2 Validate packaged source and verify the live policy with unchanged boot identity.
- [x] 1.3 Add the explicitly approved appliance keys to known_hosts and mirror their public entries in Host Bot Ansible.
- [x] 1.4 Inherit the consolidated updater defaults instead of repeating them in host variables.
- [x] 1.5 Use the self-contained automatic-update entrypoint once; retain configuration-only deployment.
- [x] 1.6 Validate the simplified caller offline; previous live policy evidence predates this revision.

2026-10-30: Host Bot packaging, local syntax, automatic_updates task selection,
and strict owner OpenSpec validation passed. The task listing includes backend
selection and native automation without immediate upgrades or reboot. Logs
are under `out/host-security-updates-20261029/simplification/updater/`.
This revision performed no live reconciliation; the parent owns combined checks.

## Normal role configuration-only control

Review feedback is addressed by `update_all_packages_upgrade_enabled`, default
true. Host Bot declares false in owner variables and invokes the normal role.
Packaged syntax and native local check mode passed with zero changes: automatic
OSTree configuration ran while immediate upgrade, reboot-hint, and reboot tasks
were skipped, even with the fixture reboot flag true. The before/after boot
assertion passed. Evidence and source hashes are retained in
`out/host-security-updates-20261029/host-bot-updater-review/result.json`.
The check-mode verification performed no live reconciliation or reboot.

The packaged `automatic_updates` listing with immediate upgrades explicitly enabled still selects only automatic policy/timer tasks; the caller relies on native role tags. Evidence: `host-bot-updater-review/tag-selection-result.json`.

The native known_hosts module's missing-parent failure was reproduced in an
empty fixture home. The candidate tagged block then created mode 0700 `.ssh`,
added both approved keys, preserved an unrelated entry, and made zero changes
on its second application. Evidence: `ansible-review/host-keys-before/` and
`ansible-review/host-keys-after/` under the original private task root.
