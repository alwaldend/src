---
title: ComfyUI
description: Pinned local CPU image generation through Bazel
---

This package owns local ComfyUI launch and repeatable CPU image-generation
benchmarks. Upstream source and model weights retain their upstream ownership
under third_party. Python packages are locked locally; no host Python packages
or global configuration are installed.

## Commands

```sh
bazel_agent bazel run //tools/comfyui:benchmark
bazel_agent bazel run //tools/comfyui -- --directory out/comfyui/session
```

The benchmark generates one seeded 512 by 512 SD-Turbo image with one sampling
step and eight CPU threads. Each invocation creates a separate run directory
beneath `out/comfyui/runs` containing the PNG, workflow, history, server log,
and timing report. Startup time is separate from generation time; the first
generation includes checkpoint loading. `--steps`, `--threads`, `--seed`,
`--prompt`, and `--timeout` select benchmark inputs.

The interactive command serves the local frontend at http://127.0.0.1:8188.
It uses the same CPU-only environment and checkpoint. The thin `cmd/serve` launcher prepares state and executes upstream ComfyUI
directly. `cmd/benchmark` owns readiness checks, workflow submission, output
verification, and timing; it starts that same launcher. Interrupting either
command stops its owned server. Inputs, outputs, user data, and temporary
files stay in the selected writable directory, separate from Bazel runfiles.

For interactive SD-Turbo generation, select `sd_turbo.safetensors`, 512 by 512,
Euler sampling with the normal scheduler, one to four steps, and CFG 1.0.

Custom nodes and partner API nodes are disabled. Hugging Face and Transformers
are configured for offline operation. `bazel run` is not an OS sandbox: this
configuration does not enforce a process-wide network or filesystem boundary.
Linux system libraries remain a platform prerequisite for upstream native wheels.

## Dependencies

ComfyUI v0.43.0 and SD-Turbo are declared under third_party. This package uses
Bazel Python 3.13.13 and a native uv.lock rather than the shared Python dependency
set. CPU wheel availability determines the runtime selection. Updating the lock
requires network access only in the upstream sandboxed resolution action:

```sh
bazel_agent bazel run --action_env=HTTPS_PROXY --action_env=HTTP_PROXY --action_env=NO_PROXY //tools/comfyui:uv_lock.update
```

The proxy flags forward the controller's existing configuration; omit them when
direct registry access is available. The initial build downloads several GB,
including the 5.2 GB complete model checkpoint. Subsequent inference uses the
declared local inputs without package installation or model downloads.

Model-bearing targets are tagged `manual` to keep multi-GB downloads out of
ordinary wildcard builds. The explicitly named commands above still build and run them.

## Z-Image-Turbo CPU trial

```sh
mkdir -p out/comfyui
bazel_agent bazel run --script_path=out/comfyui/z_image_benchmark //tools/comfyui:z_image_turbo_benchmark
systemd-run --user --scope -p MemoryHigh=32G -p MemoryMax=36G -p MemorySwapMax=0 "$PWD/out/comfyui/z_image_benchmark" --timeout 2400

# For the interactive UI, generate the other launch command and apply the same limits.
bazel_agent bazel run --script_path=out/comfyui/z_image_server //tools/comfyui:z_image_turbo
systemd-run --user --scope -p MemoryHigh=32G -p MemoryMax=36G -p MemorySwapMax=0 "$PWD/out/comfyui/z_image_server" --directory out/comfyui/z-image-session
```

The benchmark defaults to 512-square images, eight res_multistep/simple steps,
CFG 1, and eight CPU threads. `--size 1024` selects the upstream workflow's
normal resolution. `--steps`, `--seed`, and `--prompt` select trial inputs.
Reports include every model file's SHA-256. Generation time includes initial
model loading but excludes Bazel downloads and server startup.

The three model files are pinned under
`third_party/co_huggingface_comfy_org_z_image_turbo`: 12.3 GB diffusion weights,
8.0 GB Qwen3-4B text encoder, and 0.34 GB VAE. Bazel downloads and verifies
them; the launcher stages only declared runfiles. No additional Python
packages or custom nodes are needed. All model-bearing commands are manual.

Z-Image's default targets require an effective cgroup v2 RAM ceiling of at most
36 GiB and zero cgroup swap before starting inference. The commands above use
a transient user systemd scope; they do not install or modify host services.
The scope can terminate a run that exceeds its limit. A Linux host with a
working user systemd manager and memory controller is required. Generating
the Bazel launch script first keeps the inference scope separate from Bazel's
persistent server. SD-Turbo retains its smaller existing runtime contract.
