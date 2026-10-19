---
title: Anki native backend and schemas
description: Pinned upstream Anki Cargo workspace, native library, and protobuf schemas
---

This package imports Anki [26.09.3](https://github.com/ankitects/anki/releases/tag/26.09.3),
commit `29bb700b951e3f0c0cb69b77c0180fc1fe33e6ba`, with its pinned submodules.
The imported source receives the maintained `wasi.patch` for offline WASI builds.
`Cargo.lock` fixes the approved workspace and server dependency versions and
registry checksums. Anki is licensed under GNU AGPL version 3 or later; the
acquired source retains its notices.

## Build

From the repository root:

```sh
bazel_agent bazel build //third_party/com_github_ankitects_anki:native
bazel_agent bazel build //third_party/com_github_ankitects_anki:protos
```

[`crate_universe`](https://bazelbuild.github.io/rules_rust/crate_universe_bzlmod.html)
reads the imported upstream workspace's `Cargo.toml` and this package's lock.
Approved server dependencies join that resolution through `crate.spec`.
The generated `all_crate_deps`, `aliases`, and `crate_edition` helpers supply
external dependencies and metadata to standard `rules_rust` targets.
`anki.BUILD.bazel` declares the required workspace-member edges; `native.bzl`
shares the ordinary library/build-script declarations and derives Cargo
environment variables from the imported manifests, including inherited
workspace fields.

The native target builds `anki`, `anki_io`, `anki_i18n`, `anki_proto`, and
`anki_proto_gen`. Anki's Python bridge and development programs participate in
workspace dependency resolution but are not built for the server. Only targets
required by a requested Bazel consumer are compiled.

Workspace feature resolution enables native TLS alongside Rustls. The
`openssl-sys` build script uses the repository's existing pinned OpenSSL headers
and static libraries through a standard crate annotation; no host OpenSSL
installation is required.

Anki's original directory structure is preserved, so source-relative `.version`,
schema, and translation paths work directly. The original Cargo build scripts
run in their member directories with declared source inputs. The protobuf
script uses the repository-pinned protoc and writes descriptors to its declared
output directory. Pure Rust generators omit unused C/C++ toolchain inputs;
native dependencies retain their own toolchain requirements.

Cargo resolves dependencies during repository acquisition using its normal
cache, or `CARGO_HOME` when supplied. Compilation uses the pinned Bazel Rust,
C/C++, and protoc toolchains in offline sandboxes. Native builds currently
support Linux x86_64. A cold build acquires the pinned Anki Git workspace,
submodules, and required registry crates.

Dependency metadata also includes `wasm32-wasip1` for the SQLite acceptance
fixtures and offline Anki guest in `projects/anki_grpc`. The patch keeps native
network dependencies behind Cargo target conditions and retains the shared sync
metadata used by local storage. Anki's descriptor generator emits explicit
unsupported handlers for the guest's network services. Its HTTP sync-server
entry points fail before acquiring filesystem or network resources.

The guest reuses the existing URL parser through Ammonia for local media URLs,
uses the host-supplied temporary directory, and performs backups and compression
synchronously without guest threads. `fsrs_wasi.patch` makes FSRS evaluation
submit work synchronously on WASI so its result channel cannot wait for absent
Rayon worker threads. Native implementations remain behind their existing
platform paths. `chrono_wasi.patch` replaces pinned Chrono's WASI UTC fallback
with host timezone callbacks, preserving historical/DST mappings for collection
scheduling, timestamps and backups. Versions and native implementations are
unchanged. These patches are applied to their pinned upstream sources during
Bazel acquisition; external cache directories are not edited.

This package exposes the upstream library, descriptors, descriptor helpers,
and Go schemas. Its native source and build dependencies contain no project
implementation inputs. First-party gRPC server source, command generation,
executable targets, and Clippy checks belong to
[projects/anki_grpc](../../projects/anki_grpc/README.md).

The `:protos` target groups all upstream Anki schemas into one `proto_library`.
A glob discovers schema files and their internal imports. Go bindings retain
separate package targets for all 25 schemas to preserve existing import paths
and generate typed gRPC clients from the original service declarations using
rules_go's standard compilers. `schemas.bzl` declares the complete package/import
catalog used by the overlay and by consumers linking all descriptors.
Application-specific descriptor registration belongs to its consumer.
Rust uses the upstream
`anki_proto` workspace member.

## Dependency metadata

`Cargo.lock` is the package's only dependency lock. `MODULE.bazel.lock` records
generated repository definitions through the normal Bzlmod workflow and can
be larger than the Cargo lock. The original backend dependency proposal and
additional workspace-only locked set were explicitly approved; the server
reuses its approved Tonic and HTTP/Tokio dependency versions.

After an intentional dependency declaration change, regenerate metadata:

```sh
bazel_agent bazel mod deps --lockfile_mode=update
```

Inspect Cargo and module lock diffs before delivery. New dependency versions
require the repository's explicit approval.
