## Scope

Only the dedicated runner's Ansible-managed Bazel defaults change. The runner
retains capacity one. Workstation settings and checkout behavior are unchanged.
No live deployment was performed.

## Verification

- `git diff --check`: passed.
- `bazel run //:format`: passed; only task-owned files changed.
- `bazel build --config=lint //infra/dev/ansible:ansible_bin //infra/dev/ansible:srcs_files`: passed.
- `bazel run //infra/dev/ansible:syntax`: passed offline.
- `bazel test //:repo_quality_test`: all 27 tests passed.
- `OPENSPEC_PROJECT=infra/dev bazel run //tools/openspec -- validate configure-runner-resource-budgets --strict`: passed.

Bazel 9.3.1 accepted the exact runner bazelrc through `--bazelrc` and
`build --announce_rc --nobuild //infra/dev/ansible:srcs_files`. A representative
`test --announce_rc //tools/repo_quality/test/root_build:root_build_test` also
passed with that file, overriding only the Android SDK/NDK paths with the
workstation's existing installations. Its announced inherited build options
included the host-relative job, CPU and memory budgets.

An initial test analysis using the runner SDK paths failed because those
installations are absent on this workstation. `test --nobuild` parsed the
options but cannot execute tests; the subsequent normal test command passed.
These checks validate configuration and packaging, not an 80% hard usage cap
or observed runner performance.
