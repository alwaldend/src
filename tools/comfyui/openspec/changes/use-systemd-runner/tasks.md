## 1. Integration

- [x] 1.1 Update server and benchmark instructions to use the generic runner with unchanged limits; inspect the documented commands.
- [x] 1.2 Start the actual Z-Image target through the documented run-under pattern, verify HTTP readiness and cgroup limits, and verify owned-process shutdown.

## 2. Validation

- [x] 2.1 Run strict change validation and the owning package's checks before publication.

## Earlier CLI-only verification

On 2026-11-01, the real `//tools/comfyui:z_image_turbo` target started through
`--run_under=//tools/systemd_run ...` and served HTTP 200 in 14.55 seconds,
including cached Bazel analysis/build. Its kernel scope reported
MemoryHigh=32 GiB, MemoryMax=36 GiB, and MemorySwapMax=0. The configured model
arguments reached the unchanged launcher. Interrupting the owned run stopped
the server, and its endpoint no longer accepted connections.

The reproduction script, startup log, and `comfyui-e2e.json` are under ignored
`out/systemd-run/`. Strict change validation, repository quality/Buildifier
checks, and affected-target lint passed. No inference benchmark was repeated
because the application and model inputs did not change.

## 3. Normal public targets

- [x] 3.1 Declare public wrappers with BUILD-owned limits and private raw binaries; inspect configuration and update commands.
- [x] 3.2 Start the public server without run_under, verify HTTP and kernel limits, then verify shutdown; validate benchmark packaging.

## Normal-target verification

On 2026-11-01, `//tools/comfyui:z_image_turbo` started without run_under,
served HTTP 200, reported MemoryHigh=32 GiB, MemoryMax=36 GiB, and zero swap,
and stopped after interruption. The public benchmark generated and decoded
one 64 by 64 PNG with one sampling step, confirming its private server launch.
This was a packaging smoke check, not a performance comparison.

Evidence is under ignored `out/systemd-rule/`: `comfyui-e2e.json`, the
benchmark run's `report.json` and PNG, and the delivery receipts. Strict
OpenSpec validation, all 28 quality/Buildifier tests, and semantic lint passed.

## 4. Frontend initialization

- [x] 4.1 Resolve the bundled frontend's physical static directory and declare its existing dependency; verify required asset responses.
- [x] 4.2 Start the exact main target on port 8188 and verify the browser reaches the editor, with backend and scope checks preserved.
- [x] 4.3 Stop the owned test instance, run validation gates, and preserve final candidate evidence for delivery.

## Frontend browser verification

On 2026-11-01, the exact command `bazel_agent bazel run
//tools/comfyui:z_image_turbo` was tested on its default port 8188 and session
directory. The browser passed the splash screen, rendered two canvases and
an enabled Run button, and loaded 685 built-in node definitions. Closing the
initial Templates dialog with Escape worked. All 214 observed JavaScript
resources returned HTTP 200, with no JavaScript resource failures. Required
HTML-referenced scripts/styles had correct content types. Offline CSP and
32/36 GiB memory limits with zero swap were preserved.

The owned test scope was stopped and port 8188 stopped accepting connections.
All 28 quality/Buildifier checks, strict change validation, and affected-target
lint passed. Evidence and reproduction checks are under ignored
`out/comfyui-frontend/`, including `browser.json`, `assets.json`, and `main.log`.
The final candidate and publication are tracked by the delivery receipts there.

## 5. Rules module integration

- [x] 5.1 Load `@rules_systemd//systemd:defs.bzl` and verify the normal main target serves the editor assets with its existing resource policy.

The standalone-module integration was verified on 2026-11-01 with the exact
main command on port 8188. The backend, required frontend scripts/styles,
and the deployed local HTTPS proxy returned 200. Kernel memory limits were
unchanged and the owned test scope stopped. Evidence is under ignored
`out/comfyui-frontend/module-startup/`. The earlier browser verification
remains applicable to the unchanged frontend launcher and wheel.
