## Purpose

Provide a verified upstream ClamAV payload for controller-side malware scanning
without requiring installation during dependency resolution or build.

## ADDED Requirements

### Requirement: Pinned upstream artifacts

The dependency SHALL expose the official release DEB and its unchanged payload
tar through repository-internal Bazel targets, verified by SHA-256 integrity.

#### Scenario: Resolve the controller payload

- **WHEN** a consumer builds the dependency targets
- **THEN** the original DEB and its payload tar are available
- **AND** the tar retains the upstream `usr/local` paths and library links
- **AND** package installation scripts and services are not executed

### Requirement: Documented controller runtime

The dependency SHALL document Linux x86_64 and glibc 2.30 or newer as runtime
requirements and the extracted library path required to run its executables.

#### Scenario: Verify an isolated extraction

- **WHEN** a supported controller extracts the payload into a private directory
  and supplies its `usr/local/lib` through `LD_LIBRARY_PATH`
- **THEN** `usr/local/bin/clamscan --version` reports the pinned release
- **AND** consumer configuration and signature databases remain consumer-owned
