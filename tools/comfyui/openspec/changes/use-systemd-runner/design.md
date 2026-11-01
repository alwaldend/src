## Context

The launcher requires a 44 GiB or smaller RAM ceiling and zero swap. The
public target should establish those limits without extra command-line flags.

## Goals / Non-Goals

Make `bazel run //tools/comfyui:z_image_turbo` sufficient to start the UI.
Keep models, inference settings, resource guards, and SD-Turbo unchanged.

## Decisions

Keep private raw server and benchmark Python binaries. Public targets use the
generic rule, with application arguments and a shared property map declared
in BUILD.bazel. The benchmark invokes its private raw server inside the
existing scope, avoiding nested scopes and duplicated model defaults.
The wrapper merges Python runfiles. Caller arguments remain available.

## Risks / Trade-offs

Commands remain in the foreground and hold the workspace's Bazel lock.
Verify model arguments, HTTP readiness, kernel limits, and shutdown through
the normal public server target. Build and inspect the public benchmark
wrapper and run a tiny one-step image check for its internal child launch;
no performance benchmark needs repeating.

## Frontend static-file serving

The index HTML can load while JavaScript assets return 404: Bazel's runfiles
tree contains file symlinks whose resolved paths are outside the served root,
and aiohttp's static handler rejects those paths. Resolve the pinned frontend
package's module file to its physical directory and pass its `static` directory
through upstream's `--front-end-root` option. Keep upstream's symlink protections
and offline Content-Security-Policy intact. Declare the existing frontend
package as a direct dependency of both launchers; no dependency pin changes.

Before implementation, cover required JavaScript/CSS status and MIME types,
browser initialization beyond the splash screen, backend readiness, existing
scope limits, and owned-process cleanup. Test the exact main target on port
8188 as requested. Preserve the browser and asset-check evidence under ignored
`out/comfyui-frontend/`.
