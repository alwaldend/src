## 1. Configuration

- [x] 1.1 Halve the two live user bazelrc scheduling settings through the explicitly authorized direct edit; verify four jobs and `HOST_CPUS*0.25` by reading the live file.
- [x] 1.2 Mirror the settings in checked-in Ansible desired state and update the owner README; verify the live and checked-in bazelrc files match.

## 2. Validation

- [x] 2.1 Verify Bazel parses the new build defaults and validate the OpenSpec change without running Ansible; repository quality and affected-package lint gates are recorded in the delivery receipt.

Verified on 2026-10-11: the live and checked-in bazelrc files match; `build --announce_rc --nobuild` reads `--jobs=4` and `--local_resources=cpu=HOST_CPUS*0.25`; OpenSpec strict validation passes. No Ansible playbook was run.
