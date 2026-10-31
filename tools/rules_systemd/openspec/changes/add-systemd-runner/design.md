## Context

An ordinary executable target should own the resource policy. The existing Go
launcher owns foreground execution through a transient user systemd scope.

## Goals / Non-Goals

Add a reusable rule with `binary`, `arguments`, `data`, and `properties`
attributes. Preserve the existing CLI and process contract. Builds, service
installation, and detached sessions remain outside the runtime launcher.

## Decisions

Copy the existing launcher with the repository's hermetic copy action and emit
an adjacent JSON configuration. A real copy ensures executable canonicalization
still finds the target-specific configuration. Resolve the binary through Go
runfiles; prepend configured arguments before caller arguments. Expand location
expressions, merge application/data runfiles, and forward the binary's runtime
environment provider. Arguments belong on the wrapper; raw binaries stay private.

Reuse the pinned Go runfiles library and copy toolchain. No new dependency,
shell launch script, target-building subprocess, or per-target Go compilation
is needed. Package executable and JSON together. Invalid metadata fails before
launch; direct CLI execution still requires explicit properties.

## Risks / Trade-offs

- Linux user systemd and delegated controllers remain prerequisites; failures
  never fall back to unbounded execution.
- Ordinary Bazel foreground and workspace-lock behavior remains.
- Both executable and configuration must be packaged; verify runfiles-based
  resolution and the actual ComfyUI consumer.

## Failure Cases and Verification

Before implementation, cover literal configured/caller arguments, location
expansion, runfiles data, runtime environment, memory/swap limits, nonzero exit,
and cleanup. Retain real-scope CLI checks for invalid input, unavailable
systemd, stdin/cwd preservation, and interrupts. A small Bazel fixture and the
actual ComfyUI startup produce JSON evidence under ignored `out/systemd-rule/`.
A tiny one-step inference check covers the benchmark's internal child launch;
no performance benchmark needs repeating.

## Standalone module

Use the conventional `systemd/defs.bzl` public rule surface and
`cmd/systemd_run` executable. All labels resolve within the module; the root
repository registers a local override and ComfyUI loads the external API.
The module owns its Go manifest and generated Bzlmod lock, shared Bazel rc
imports, documentation export, and existing real-scope probe. Reuse the
existing Go, Cobra, copy-action, Python fixture, and documentation pins.
Verify standalone build/test, parent runfiles and ComfyUI assets, and the
shared workspace discovery matrix before publication.

Disable systemd argument environment expansion explicitly so dollar expressions
remain literal on newer systemd releases. Verify this in the existing real-scope
fixture along with empty arguments and child exit status.
