## Decisions

The approved host-key task must work when the operator's home has no `.ssh`
directory. Native `known_hosts` creates only the file, so create its private
0700 parent in the same tagged block before reconciling keys. Verify both a
fresh home and an unchanged second application with unrelated entries retained.

Import the canonical updater once through its normal entrypoint. Owner variables
set `update_all_packages_upgrade_enabled: false`, so deployment configures native
automation without an immediate upgrade, reboot detection, or reboot task.
OSTree policy stage prepares deployments while leaving the running boot unchanged.
On-demand upgrades and audits remain in infra/security, which consumes inventory
identities without importing these owner deployment variables.
The shared role owns its task tags; the owner does not tag the whole import as
`automatic_updates`, which would also select its immediate-maintenance block.

Host Bot inherits automatic updates enabled and reboot disabled from the
consolidated role defaults. It does not duplicate those values in owner variables.
