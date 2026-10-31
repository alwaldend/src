## Inputs

- [x] Pin approved source, CPU dependency lock, and SD-Turbo model.

## Implementation

- [x] Add CPU launch and bounded end-to-end benchmark commands.
- [x] Run the benchmark locally and inspect its PNG and timing report.
- [x] Document invocation, guarantees, and observed limits.
- [x] Complete delivery validation.

## Measured acceptance

On 2026-10-31, `bazel_agent bazel run //tools/comfyui:benchmark` generated and
decoded a 512 by 512 PNG on the i5-13400 desktop. The report confirmed CPU
device type, Python 3.13.13, and torch 2.15.0+cpu. With eight threads, one
sampling step, and seed 42, startup took 10.756 seconds and first generation
took 11.032 seconds, including model loading. These are one-run observations,
not a comparison with the XCP-ng host or a steady-state throughput claim.

Evidence is under ignored
`out/comfyui/runs/20261031T222236-9d831c6f/`: workflow.json, history.json,
server.log, report.json, and output/sd_turbo_00001_.png. The PNG was visually
inspected: the red teapot/window prompt is recognizable, with distorted
handle geometry in this one-step draft. The owned server stopped afterward.

The interactive target also served HTTP 200 HTML on a task-selected localhost
port with the offline Content-Security-Policy header. Interrupting its launcher
returned exit 130 and the server endpoint stopped accepting connections.
Repeatable check evidence is in ignored out/comfyui/frontend-check.json.

The CPU lock contains 85 packages and no CUDA/NVIDIA packages. Initial
resolution failed without the controller's existing proxy forwarded to its
sandboxed network action (COMFYUI_LOCK_PROXY_REQUIRED); the documented proxy
flags resolved that failure. ComfyUI archive extraction requires an explicit
tar.gz type because its immutable codeload URL has no suffix.

A second generation after the import-policy correction produced an identical
PNG SHA-256, 4b26bc77bc476e0865620beaaa75b746c26c594f84fb70f9d0ed556324adf71a.
Its evidence is under out/comfyui/runs/20261031T222936-b3d60cc3. All 28 selected
repository/skill tests and semantic lint passed. Exact-candidate publication
checks and remote publication state are owned by the delivery receipt.

## Command separation

- [x] Separate the thin interactive launcher from benchmark orchestration.
- [x] Verify generation and interactive shutdown against the separated commands.
- [x] Validate the refactor; publication is tracked by the task delivery receipt.

On 2026-10-31, the separated benchmark produced a visually inspected 512-square
PNG in 11.495 seconds after 11.566 seconds of startup. Evidence is under
`out/comfyui-refactor/runs/20261031T224904-fefb80c8/`. The interactive command
served HTTP 200 HTML and closed its endpoint after SIGINT; evidence is in
`out/comfyui-refactor/interactive.json`. Repository quality checks passed.

The initial subprocess launch failed because Bazel target `args` apply to
`bazel run`, not direct executable invocation (COMFYUI_SUBPROCESS_ARGUMENTS).
Passing the declared server and model inputs explicitly corrected that failure.
The end-to-end rerun covers that boundary. A stale previous delivery launcher
also required rebuilding its task-local run script (DELIVERY_RUNFILES_EXPIRED).

## Z-Image-Turbo CPU trial

- [x] Declare the requested diffusion model, Qwen3 text encoder, and VAE as checksum-verified Bazel inputs.
- [x] Extend model staging and the benchmark for the Z-Image-Turbo workflow without adding Python dependencies.
- [x] Generate and inspect an image on CPU, recording timings and model identities.
- [x] Complete command validation; commit and publication state are owned by the delivery receipt.

- [x] Require and verify a transient 36 GiB RAM / zero-swap boundary for Z-Image; verify unbounded refusal and bounded execution.

## Z-Image measured acceptance

On 2026-10-31, the Bazel Z-Image-Turbo target generated and decoded a 512 by
512 PNG on CPU with eight sampling steps, eight threads, and seed 42.
Startup took 14.052 seconds and first generation took 391.319 seconds,
including model loading. This is one observed run inside a transient
MemoryHigh=32G, MemoryMax=36G, MemorySwapMax=0 scope, not a steady-state or
larger-resolution throughput claim. The report identifies the CPU device
and all three model SHA-256 values, which match the declared publisher hashes.

Evidence is under ignored `out/z-image/bounded-runs/20261031T234757-53e9154c/`:
workflow.json, history.json, server.log, report.json, and
output/benchmark_00001_.png. The PNG was visually inspected: a recognizable
red teapot beside a window; the small card text is faint at this resolution.
The owned server exited and the transient scope became inactive.

`out/z-image/memory-scope-start.txt` records the actual kernel-backed limits.
Live sampling observations stayed near 32 GiB with zero scope swap; these
observations are not a final peak measurement. The unbounded invocation was
refused before inference, with evidence in out/z-image/unbounded-refusal.json.
The stopped unbounded attempt produced no accepted image; its diagnostics
remain under out/z-image/runs/20261031T234047-bb08ee7f/.

Bazel successfully fetched and integrity-checked the model inputs; its initial
fetch/build took 1658.529 seconds. Runtime inference performed no model fetches.
The earlier scratch-download approach was stopped and its partial files removed.
The SD-Turbo compatibility command also generated a validated PNG under
out/z-image/sd-check/20261031T235556-c1689791/ in a 12 GiB / zero-swap scope.
All 28 selected repository/skill checks and semantic lint passed before final
candidate preparation. Exact-candidate checks and publication belong to the
local delivery receipt under out/z-image/.
