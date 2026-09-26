---
name: repo-external-dependency
description: Add or upgrade external software consumed by this monorepo with reproducible pins. Use for third_party packages, Bzlmod dependencies or repositories, downloaded binaries or archives, OCI images, and language dependency manifests or locks; use bazel-nested-module for standalone projects under projects/ with their own MODULE.bazel.
---

# Manage external dependencies

## Choose the owning mechanism

1. Read the root `AGENTS.md` and the nearest owning `README.md`, `BUILD.bazel`,
   and `include.MODULE.bazel` when present. Follow the `repo-bazel` conventions
   for every Bazel command.
2. Inspect the closest dependency of the same kind before choosing a mechanism:
   - BCR modules belong in the owning `include.MODULE.bazel` as `bazel_dep`.
   - Standalone release binaries use `rules_binary_toolchain` and a checked-in
     `binary_toolchain.json`.
   - Source archives and individual files use `http_archive` or `http_file`.
   - OCI artifacts use the repository's `rules_oci` pattern and an immutable
     digest.
   - Language packages use the owning manifest and lockfile workflow.
3. Follow `project-layout` for dependency ownership and naming. Simple external
   dependency definitions belong in a reverse-domain, underscore-named package
   under `third_party/`; standalone Bazel modules keep their dependencies local.
   Keep consumer-specific configuration and first-party wrappers with the consumer.

## Check publisher reputation

Choose a reputable publisher before selecting an artifact. Prefer official
upstream releases or an established distribution or vendor. Inspect maintenance
history, project ownership, and auditable build and release provenance; do not
adopt an obscure third-party repackager merely because its binary is convenient.
A checksum establishes artifact integrity, not publisher trust. Popularity or
star counts alone neither establish nor disqualify a legitimate upstream.
When the consumer permits a host runtime prerequisite, an existing distribution
package can avoid another binary distributor; document that non-hermetic boundary.

## Add a standalone release binary

For a release binary consumed through a Bazel toolchain, create a package
such as `third_party/com_github_owner_tool/` with the following files.

1. `binary_toolchain.json` is the lock and the source of truth for archives:
   - `toolchains[].name` and every `binaries[].name` must be valid Starlark
     identifiers (letters, digits, underscores; no hyphens). The rule derives
     target names such as `{name}_binary` and `{name}_toolchain` from them, so
     a hyphenated name breaks analysis.
   - `integrity` is the SRI form `sha256-<base64>` of the SHA-256 digest, not
     the hex digest. Convert hex to base64 before pasting it into the lock.
   - Put execution-platform constraints in each archive's `toolchain` key, and
     point `binaries[].path` at the extracted file inside the archive
     `output` directory.
2. `include.MODULE.bazel` calls the `binary_toolchain_extension` with
   `lock = "//third_party/com_github_owner_tool:binary_toolchain.json"` and
   `toolchains_map = {"<name>": "com_github_owner_tool"}`, then
   registers the created repository and calls `use_repo`. Keep the underscore
   name here identical to the lock's toolchain name.
3. `BUILD.bazel` aliases the generated
   `@com_github_owner_tool//:<name>_binary` under a user-facing target
   with `visibility = ["//:__subpackages__"]`; the alias may be hyphenated
   even though the generated targets are not.
4. `README.md` links the upstream release repository.

Then regenerate the root `MODULE.bazel` with
`bazel_agent bazel run //tools/bazel_module:update`, and run
`bazel_agent bazel mod deps --lockfile_mode=update` to reconcile the module lock.
The archive contract lives in the package `binary_toolchain.json`, so adding
or changing an archive normally does not change the module lock. Verify the
pin by building and running the alias, for example
`bazel_agent bazel run //third_party/com_github_owner_tool:<alias> -- --version`,
before handing off.

## Preserve reproducibility

- Pin an immutable release, commit, or digest and use the primary publisher's
  HTTPS endpoint. Do not depend on moving branches, tags such as `latest`, or
  unverified mirrors.
- Require SRI integrity for downloads. Prefer a publisher-provided checksum or
  signature and verify it independently before adding the lock entry.
- Fetch only the platforms the repository consumes. In
  `binary_toolchain.json`, put execution constraints under `toolchain`, not the
  obsolete `rules` key.
- For a dedicated `third_party/` package, add upstream links in `README.md`, a
  least-privilege BUILD alias, an `include.MODULE.bazel` fragment, and the root
  `MODULE.bazel` include.
- Never commit downloaded binaries, package-manager caches, or generated
  credentials.

## Update and validate

Update generated locks only through their owning workflow. For root Bzlmod
changes, run:

```sh
bazel_agent bazel mod deps --lockfile_mode=update
```

Review the lockfile diff for unrelated resolution changes. Then query and build
the dependency package and its narrowest consumer, followed by:

```sh
bazel_agent bazel test //:buildifier_test
```
