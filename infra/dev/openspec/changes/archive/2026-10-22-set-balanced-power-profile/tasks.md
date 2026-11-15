## 1. Apply and record the requested setting

- [x] 1.1 Select `balanced` directly with `tuned-adm` and verify it with `tuned-adm active`.
- [x] 1.2 Record the current setting and authorized Ansible exception in the host-owned files; inspect the diff against the live profile and CPU power policy.
- [x] 1.3 Declare the desired profile in Ansible and reconcile the service, selected profile, and manual selection mode; verify the packaged playbook and affected semantic lint without deploying.
- [x] 1.4 Correct the README and change artifacts to describe the managed desired state; validate the OpenSpec change and repository quality checks.

## Verification evidence

Observed on 2026-10-22 in the host's local timezone:

- `sudo -n tuned-adm profile balanced` exited successfully.
- `tuned-adm active` reported `Current active profile: balanced`.
- `/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor` reported `powersave`.
- `/sys/devices/system/cpu/cpu0/cpufreq/energy_performance_preference` reported `balance_performance`.
- `/sys/devices/system/cpu/intel_pstate/min_perf_pct` reported `18`, down from `100`.
- `/sys/devices/system/cpu/intel_pstate/max_perf_pct` remained `100`.

No Ansible deployment was executed. Design is omitted because the change
uses the existing Tuned CLI and the component's established Ansible pattern.

Review-fix validation passed: the packaged `ansible_bin` build and its
offline `--syntax-check`, semantic lint of the Ansible package, strict
OpenSpec validation, and all 26 repository quality tests.
