## Single entry point (2026-10-30)

The main role selects its backend once and passes it to internal configuration
and reboot-detection tasks. Unused direct task-file interfaces are retired.
Normal role controls and the configuration-only automatic_updates tag remain;
automation defaults true and reboot defaults false.

A pinned native local check-mode invocation with immediate upgrades disabled
completed with zero changes and unchanged boot identity. Automatic-only task
selection omits immediate upgrades/reboot and selects the backend once. No
package update, timer reconciliation or reboot was performed for validation.
Collection and Host Bot packaging and packaged fleet/appliance syntax passed.
Evidence: `out/host-security-updates-20261029/native-simplification/updater/result.json`
and `package-results.json`.
