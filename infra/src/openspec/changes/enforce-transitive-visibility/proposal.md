## Why

The repository documents dependency boundaries but currently permits project
builds to consume infrastructure artifacts, including a generated certificate.
Bazel 9 supports enforcing these boundaries transitively, but declarations must
be maintained in every affected package.

## What Changes

- Add the standalone `projects/bazel_gazelle_enforce_transitive_visibility`
  Gazelle extension. It derives `package(transitive_visibility = ...)` from an
  central TOML path-to-group configuration.
- Add a package group for each repository tree and central path mappings, register the
  extension in repository Gazelle, and generate the declarations.
- Enable experimental transitive visibility in checked-in repository Bazel
  configuration and update first-party MODULE compatibility declarations.
- **BREAKING**: builds with prohibited dependencies fail during analysis.
- Verify directory-prefix selection, preservation of unrelated package attributes,
  repeat generation, direct rejection, and the documented native limitation
  for indirect consumers whose intermediary declares its own policy.
- Move public CA certificate and SSH public-key snapshots into their data
  owners. Separate Terraform roots in those directories refresh them from
  Vault data sources; infrastructure Terraform stops exporting the files.
- Repair the stale `bzlenv` references exposed by validation, using the
  existing public binary targets in `third_party` for the eight migrated tools.
- Repair the Gradle wrapper's shell-runfiles dependency so the pinned
  distribution can be analyzed and executed under the repository's Bazel version.
- Remove the retired LeetCode downloader container path and its unused Alpine
  declaration, and remove the stale root requirements alias after the uv migration.

## Capabilities

### New Capabilities

- `transitive-dependency-policy`: Repository-owned policy declarations and
  enforcement, with a reusable configuration-driven generator.

### Modified Capabilities

None.

## Impact

The change affects the new plugin module, root Gazelle integration, shared
repository Bazel configuration, first-party MODULE files, generated BUILD
declarations, and the nested-workspace validation matrix. It reuses existing
Gazelle and Go dependencies. Host configuration and live infrastructure are
outside scope. Preserve the documented website documentation and developer/test
exceptions while prohibiting the Ansible project's infrastructure dependency.
The Ansible collection may consume the data-owned public certificate snapshots.

The configuration owner is `tools/transitive_visibility/transitive_visibility.toml`.
Remove BUILD policy comments and place generated `package()` declarations after
loads, preserving existing consumer groups and exceptions.
