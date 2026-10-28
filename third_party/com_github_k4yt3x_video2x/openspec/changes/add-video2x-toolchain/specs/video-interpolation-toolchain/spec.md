## Purpose

Provide reproducible repository access to the upstream Video2X executable for
neural video frame interpolation on supported Linux hosts.

## ADDED Requirements

### Requirement: Pinned upstream executable

The package MUST obtain the Linux x86_64 Video2X 6.4.0 AppImage from its upstream
release using a checked-in SHA-256 integrity value and the existing binary
toolchain extension.

#### Scenario: Resolve the executable

- **WHEN** Bazel resolves the Video2X dependency on Linux x86_64
- **THEN** it verifies the downloaded artifact against the checked-in integrity
  value and provides the pinned upstream executable

### Requirement: Internal runtime access

The package MUST expose a repository-internal executable alias and document
running it without installing the AppImage or requiring FUSE.

#### Scenario: Invoke the packaged CLI

- **WHEN** a caller runs the alias with `--appimage-extract-and-run --version`
  on a compatible host
- **THEN** the packaged CLI reports Video2X version 6.4.0

### Requirement: Provenance and runtime requirements

The package MUST document the upstream release and license, supported platform,
Vulkan runtime requirements, and a representative RIFE interpolation invocation.

#### Scenario: Prepare a local interpolation run

- **WHEN** a caller reads the package documentation
- **THEN** they can identify the upstream artifact and the required host runtime,
  select the RIFE processor, and keep extracted files and caches in task scratch
