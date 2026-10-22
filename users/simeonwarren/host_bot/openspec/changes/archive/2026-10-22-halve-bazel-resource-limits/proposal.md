## Why

Concurrent builds in independent workspaces each receive eight jobs and half of the host CPUs. Halving these per-invocation limits reduces contention while a shared host-wide limit remains outside this change.

## What Changes

- Set the managed user bazelrc to four jobs and `cpu=HOST_CPUS*0.25`.
- Apply those two settings directly to the live user bazelrc, as explicitly requested, without running Ansible.
- Mirror the live change in checked-in Ansible configuration and update the owner README.

## Capabilities

### New Capabilities

- `bazel-resource-scheduling`: Per-invocation scheduling defaults for host-bot builds and tests.

### Modified Capabilities

None.

## Impact

Changes `ansible/files/bazelrc`, the owner README, and `/var/home/simeonwarrenbot/.bazelrc`. No new dependency or aggregate resource enforcement is introduced. Existing running commands retain their original settings.
