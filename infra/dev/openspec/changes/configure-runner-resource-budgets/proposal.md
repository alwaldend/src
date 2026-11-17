## Why

The dedicated Forgejo worker limits each Bazel invocation to four jobs and half
the host CPUs. Scale its scheduling defaults with the host and reserve roughly
20% of CPU and RAM capacity.

## What Changes

- Set the runner's job limit and CPU scheduling budget to `HOST_CPUS*0.8`.
- Set its memory scheduling budget to `HOST_RAM*0.8`.
- Document these per-invocation budgets and retain one Forgejo job at a time.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `bazel-resource-scheduling`: Add host-relative budgets for the dedicated runner.

## Impact

Changes `ansible/files/runner_bazelrc` and development-host documentation. The
workstation's independent scheduling defaults remain unchanged. No dependency,
checkout, workflow or live infrastructure change is included.
