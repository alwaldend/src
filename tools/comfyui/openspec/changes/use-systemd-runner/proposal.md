## Why

Z-Image should start through an ordinary Bazel target without requiring users
to provide resource policy on the command line.

## What Changes

- Wrap the public Z-Image server and benchmark targets with `systemd_run`.
- Declare the 40 GiB soft limit, 44 GiB hard limit, and zero swap in BUILD.bazel.
- Keep raw Python binaries private and retain their resource guard.
- Document ordinary single-target launch commands.
- Resolve the bundled frontend to its physical package directory so static assets load from Bazel runfiles.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. The existing ComfyUI resource contract is unchanged.

## Impact

ComfyUI BUILD declarations, launch documentation, and bounded startup checks.
No model, inference, or persistent host configuration changes.
