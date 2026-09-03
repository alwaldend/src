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
3. Keep dependency declarations with their consumer. Use a dedicated
   reverse-domain, underscore-named package under `third_party/` when the
   artifact is shared or needs its own repository definition.

## Add a standalone release binary

For a release binary consumed through a Bazel toolchain, create a package
such as `tools/<tool>/` with the following files.

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
   `lock = "//tools/<tool>:binary_toolchain.json"` and
   `toolchains_map = {"<name>": "com_alwaldend_src_tools_<tool>"}`, then
   registers the created repository and calls `use_repo`. Keep the underscore
   name here identical to the lock's toolchain name.
3. `BUILD.bazel` aliases the generated
   `@com_alwaldend_src_tools_<tool>//:<name>_binary` under a user-facing target
   with `visibility = ["//:__subpackages__"]`; the alias may be hyphenated
   even though the generated targets are not.
4. `README.md` links the upstream release repository.

Then add `include("//tools/<tool>:include.MODULE.bazel")` to the root
`MODULE.bazel` under the owning section, and run
`bazel_agent mod deps --lockfile_mode=update` to reconcile the module lock.
The archive contract lives in the package `binary_toolchain.json`, so adding
or changing an archive normally does not change the module lock. Verify the
pin by building and running the alias, for example
`bazel_agent bazel run //tools/<tool>:<alias> -- --version`, before handing off.

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
bazel_agent mod deps --lockfile_mode=update
```

Review the lockfile diff for unrelated resolution changes. Then query and build
the dependency package and its narrowest consumer, followed by:

```sh
bazel_agent bazel test //:buildifier_test
```
