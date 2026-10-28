## Context

Existing media tools cannot supply satisfactory neural interpolation of sparse
animation poses. The user authorized selecting an established project and adding
it to `third_party`. Video2X packages RIFE with ncnn/Vulkan in a Linux AppImage.
The repository already owns release executables through `rules_binary_toolchain`.

## Goals / Non-Goals

**Goals:** Pin Video2X 6.4.0 for Linux x86_64, expose a repository-internal Bazel
executable, and document headless execution and upstream provenance.

**Non-Goals:** Installing GPU drivers, adding a first-party media application,
supporting additional platforms, or checking private media and binary artifacts
into source control.

## Decisions

- Use upstream Video2X rather than the older standalone rife-ncnn-vulkan release.
  Video2X has a larger established project and bundles RIFE 4.6, 4.25, and 4.26.
- Pin the 6.4.0 AppImage by SHA-256 SRI through the existing binary toolchain
  extension. Its release API supplies no digest; the pin records the downloaded
  publisher artifact's locally calculated hash, not an independently signed hash.
- Expose the generated binary through the standard internal BUILD alias. Callers
  can pass AppImage's `--appimage-extract-and-run` option without a new wrapper.
- Keep version and download data in `binary_toolchain.json`; the owning module
  fragment registers it, and the root module is regenerated through its tool.
- Validate version, Vulkan device enumeration, and actual frame interpolation on
  this host. Runtime media and extracted files remain under ignored `out/`.

## Risks / Trade-offs

- The AppImage still needs a host Vulkan loader, compatible driver, render-device
  access, and AVX2. Document these requirements; do not change host configuration.
- Extract-and-run costs startup time and temporary storage. Document task-local
  `TMPDIR` and cache directories; callers may extract once for repeated work.
- Neural interpolation can distort large changes or occluded limbs. Runtime
  verification demonstrates executable behavior, not universal image quality.
- Upstream code is AGPL-3.0; retain the upstream source and license links.

## Migration Plan

This is an additive dependency with no existing consumers to migrate. Remove its
owning include and regenerate the root module to remove the integration.
