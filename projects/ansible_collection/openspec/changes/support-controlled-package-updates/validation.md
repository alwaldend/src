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

## DNF timer opt-out (2026-11-01)

The isolated baseline left `dnf-automatic-install.timer` active and enabled
while failing on the absent generic timer. The corrected native timer tasks
passed twelve scenarios through the pinned Ansible controller and native
`systemd_service` module with an isolated systemctl executable and JSON state:
install-only DNF4, coexisting DNF4/DNF5 timers from either backend, absent timers
from either backend, check mode, enabled DNF4/DNF5 policy, discovery failure,
stop failure, and APT/OSTree opt-out. Seven repeated opt-outs converged with zero
changes. Download-only and notification timers remained unchanged.

Source review also confirmed that timer reconciliation is outside the mutable
updater package-presence gate. The fixture exercises timer task/module
orchestration; it does not exercise real systemd, package installation, or a
distribution upgrade. No live timer or package operation was performed.

Repeatable inputs, source hashes, state and command traces are under
`out/host-security-updates-20261029/dnf-timer-fix/`. The passing matrix combines
`after-stable/result.json` (the first three cases) and
`after-remaining/result.json` (the other nine). An initial fixture interpreter
failed to import the Python standard library; its replacement run later timed
out under host memory pressure. Those failures remain retained. Using the
existing system Python for fixture endpoints and releasing this worktree's
idle Bazel server allowed the remaining cases to pass. The controller and
module source remained pinned throughout. Delivery receipts own final
candidate validation and publication state.
