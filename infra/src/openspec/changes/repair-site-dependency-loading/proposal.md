## Why

The full site dependency graph reaches macOS-only Python packages that the
native uv importer incorrectly resolves as Linux source builds. Bazel 8.7 also
creates non-relocatable Git strip-prefix symlinks, blocking the site when a
shared repository-content cache outlives its original output base.

## What Changes

- Partition Autoscroll's native uv resolution into macOS and its complement.
  The generated package markers work with the unmodified rules_python importer.
- Remove the initial importer patch, override, and patch-specific fixtures.
  Preserve existing package versions and artifact hashes; include the
  user-approved tomli 2.5.0 dependency selected by PySide6 for Python 3.10.
- Upgrade Bazel to 9.3.1 with the publisher's checksum. Its relative Git links
  survive isolated cross-output cache reuse. Remove the temporary Anki archive
  override, generated projections, templates, update targets, and formatting
  exception; retain the existing Cargo-locked source and crate_universe recipe.
- Select the existing transitive bazel_jar_jar dependency's current release for
  explicit Java provider imports required by Bazel 9, then validate consumers.

## Capabilities

### New Capabilities

- `python-lock-platform-reachability`: Import only packages reachable on each
  configured target platform from a native uv lockfile.

### Modified Capabilities

None.

## Impact

Autoscroll's uv configuration and generated lock, removal of the shared
rules_python override, the Bazel pin, and existing Rust source acquisition.
The user approved tomli 2.5.0. All operating systems remain in scope. No host
configuration change or deployment is involved.
This supports the analyzer and site change records in their owning projects.
