## Inputs and ownership

Use a third-party source archive with an upstream BUILD overlay, a component-local
native uv.lock consumed by rules_python, and checksum-verified SD-Turbo weights.
Use the existing Python toolchain when the complete CPU wheel set supports it;
otherwise select a component-specific Bazel Python toolchain without changing
the shared default. Upstream code retains its imports and source structure.

## Runtime

Run upstream ComfyUI as a packaged subprocess. Separate immutable runfiles from
writable input, output, user, and temporary directories. The benchmark submits
a fixed seeded SD-Turbo workflow to localhost and preserves its input JSON,
server diagnostics, generated PNG, and timing metadata. No GPU drivers or host
packages are installed. Dependency fetching is build setup; inference is offline.

## Failure cases before implementation

- Missing or incompatible CPU wheels: fail dependency resolution, not runtime installation.
- Missing/corrupt model: fail integrity checking or checkpoint loading.
- Occupied port or crashed startup: fail bounded readiness checking.
- Invalid workflow: report HTTP rejection or execution error, not a successful queue.
- Timeout: retain evidence and terminate the owned server process.
- Missing, empty, corrupt, or incorrectly sized PNG: fail output verification.
- Unexpected CUDA selection or downloads: CPU/offline flags and declared local model paths.

## Acceptance

Use the actual end-to-end generation command, inspect the PNG, and record measured
startup and generation times. Do not infer inference speed from CPU specifications.

## Command ownership

`cmd/serve` prepares writable state and executes upstream ComfyUI directly.
`cmd/benchmark` launches that same command, bounds readiness and generation,
verifies the PNG, and records timings. There is no mode-switching entry point
or shared implementation package. The package default remains the UI launcher.

## Z-Image-Turbo inputs and workflow

Use Comfy-Org's packaged Z-Image-Turbo BF16 diffusion weights, Qwen3-4B text
encoder, and VAE at one immutable revision with publisher SHA-256 checksums.
Bazel owns all downloads; the launcher only stages declared runfiles beneath
its writable model directory. Reuse ComfyUI's built-in UNET, CLIP, VAE,
AuraFlow sampling, and KSampler nodes with eight res_multistep/simple steps,
CFG 1, and shift 3. A 512-square CPU trial precedes larger resolutions.

Failure cases include a missing declared model, an invalid staging path,
missing core nodes, model-load failure or excessive memory demand, sampling
timeout, and corrupt or incorrectly sized output. Preserve logs on failure
and terminate only the owned server. No custom nodes or new Python packages.

The initial scratch download bypassed Bazel (Z_IMAGE_INPUT_OWNERSHIP). It was
stopped before inference; the maintained workflow declares every model input.

## Inference memory boundary

The initial unbounded CPU run filled host swap during model loading and was
stopped at the user's request (Z_IMAGE_MEMORY_PRESSURE). It reached two of
eight sampling steps at approximately 39.4 seconds per step, not a completed
image or end-to-end timing. Host available RAM recovered after stopping it.

Decision: proceed only inside a transient cgroup v2 scope with MemoryMax=44G,
MemoryHigh=40G, and MemorySwapMax=0. GPU offload flags do not bound CPU RAM,
and a virtual-address-space limit does not model resident model memory. The
launcher must verify an effective RAM ceiling of at most 44 GiB and disabled
swap before loading Z-Image. The bound may terminate an oversized run; do not
raise it automatically. No persistent systemd or host configuration changes.

Failure cases before implementation: missing cgroup v2 support, absent memory
controllers, unlimited RAM or swap, a ceiling above the requested bound, and
a scoped out-of-memory termination. Verify refusal outside a bounded scope,
the actual kernel limits inside it, and a complete generated PNG if it fits.

## Requested memory increase

The interactive workload reached approximately 32.8 GiB and repeatedly hit
the former 32 GiB soft limit. The user requested higher limits and stopped
the run. Raise the shared server/benchmark budget to 40 GiB soft, 44 GiB hard,
and zero swap, leaving about 18 GiB of physical RAM outside the workload on
this host. Derive the launcher guard ceiling from the same BUILD constant as
MemoryMax. Do not change model precision, worker architecture, or inference
settings. Validate both generated wrapper configurations and their resource
guard arguments; leave ComfyUI stopped and do not repeat inference.

## CPU thread configuration

Both Python entry points accept an optional positive `--threads` value with no
parser default. If omitted, preserve inherited or library-default thread
settings and omit the flag when the benchmark starts its server. The shared
`_SERVER_ARGS` list in BUILD.bazel supplies `--threads 12`
to all public server and benchmark targets. User arguments follow the declared
arguments, so an explicit caller flag overrides the BUILD value. The benchmark
forwards its parsed value to the owned server; the existing launcher maps it to
OpenMP, MKL, and OpenBLAS thread settings. Keep active user processes untouched.

Failure cases before implementation: omitted flags cause a None comparison or
forward the string None, zero/negative values are accepted, a stale eight-thread
fallback overrides the runtime, caller overrides are ignored by the wrapper,
or benchmark values do not reach its child. Build all public targets, verify
their declared arguments, and start owned server probes on separate temporary
ports to inspect the actual environment for 12 and an explicit override. Verify
both entry points accept omitted settings and reject invalid values before inference. Preserve
HTTP readiness and owned-process cleanup evidence without submitting an image.
