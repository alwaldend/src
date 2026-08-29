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
bazel_agent test //:buildifier_test
```
