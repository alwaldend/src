## Why

The approved host scanning workflow needs a reproducible ClamAV controller
payload without installing a package or starting a service on the build host.

## What Changes

- Pin the official ClamAV 1.5.4 Linux x86_64 DEB with SHA-256 integrity.
- Expose the original package and its unchanged payload tar through Bazel.
- Register the owning module include and document the runtime contract.

## Capabilities

### New Capabilities

- `clamav-payload`: Reproducible access to the upstream controller payload.

### Modified Capabilities

None.

## Impact

Adds `third_party/com_github_cisco_talos_clamav` and its generated root module
include. Consumers own configuration, database updates, deployment, and scans.
Downloaded binaries and validation output remain in ignored task scratch.
