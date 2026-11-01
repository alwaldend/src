## Why

The local desktop needs a measured CPU image-generation trial before allocating
an XCP-ng VM. ComfyUI and model inputs must be reproducible through Bazel.

## What Changes

- Pin the approved ComfyUI source, Python dependencies, and SD-Turbo weights.
- Add local CPU launch and a repeatable generation benchmark.
- Keep `--threads` optional with no Python default and declare `--threads 12` once in the shared public-target BUILD arguments.
- Save an actual PNG, workflow, and timing report under ignored task output.

## Capabilities

### New Capabilities

- `cpu-image-generation`: Offline local generation with declared software and model inputs.

## Impact

New tools/comfyui package and third-party input declarations. No host installation,
shared Python runtime migration, infrastructure provisioning, or external API calls.

## Z-Image-Turbo extension

The user requested a local Z-Image trial. Reuse the existing ComfyUI runtime
and declare its required diffusion, text-encoder, and VAE files through Bazel.
Add an eight-step CPU benchmark with a verifiable PNG and measured timings.

Z-Image runs must use a transient 44 GiB RAM / zero-swap scope, and default
Z-Image targets refuse inference without that effective kernel boundary.
