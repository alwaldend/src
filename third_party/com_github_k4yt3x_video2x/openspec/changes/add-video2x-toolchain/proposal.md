## Why

The repository needs a reproducible neural frame-interpolation executable for
adding intermediate frames to existing animations. The user authorized an
established upstream project and requested a proper `third_party` integration.

## What Changes

- Pin the upstream Video2X 6.4.0 Linux x86_64 AppImage with SHA-256 integrity.
- Expose a repository-internal executable through the existing binary toolchain rules.
- Register the dependency through its owning module include and generated root module.
- Document local invocation, bundled RIFE models, hardware requirements, and verification.

## Capabilities

### New Capabilities

- `video-interpolation-toolchain`: Reproducible access to the upstream Video2X executable.

### Modified Capabilities

None.

## Impact

Adds `third_party/com_github_k4yt3x_video2x` and its root module registration.
The upstream AppImage and models remain downloaded artifacts, not tracked source.
Media outputs and runtime caches remain in ignored task scratch.
