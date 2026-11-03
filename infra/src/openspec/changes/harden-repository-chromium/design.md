## Context

See proposal.md. The dependency lock owns the browser version; Puppeteer
25.2.1 launches it. Draw.io and screenshots explicitly disable sandboxing;
Mermaid and browser tests do not. T3 is maintained separately upstream.

## Goals / Non-Goals

Retain native Chromium confinement alongside Bazel's existing sandbox, without
new dependencies, persistent host changes, or a silent unsafe fallback.
This is not a container/VM provisioning change or a claim of complete isolation
from every browser vulnerability.

## Decisions

**Proceed:** use Chrome for Testing stable 156.0.8078.80, as reported by the
[publisher manifest](https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions-with-downloads.json)
on November 3, 2026. The Linux archive was downloaded directly from Google;
its SHA-256 is recorded in the owning binary lock. No new dependency is added.

Remove `--no-sandbox` rather than introducing a parallel launcher. Puppeteer
retains sandboxing by default. Retain the existing page request and input
controls. A private browser profile or Bazel sandbox alone is not a replacement
for Chromium's renderer sandbox. Broad host sandbox relaxations are rejected.

Existing command E2Es cover successful renders and invalid input. Before edits,
the failure cases identified are: nested sandbox startup failure, silent unsafe
fallback, input/network restrictions changing, and browser-dependent layout,
shadows, labels, or highlighting changing. Reuse those E2Es and independently
observe renderer confinement rather than adding flag-string tests.

## Risks / Trade-offs

- Nested sandbox incompatibility: diagnose the actual launch failure and keep
  confinement; do not restore the disabling flag to make tests pass.
- Rendering changes: inspect representative SVG/PNG outputs and execute the
  direct consumers' existing browser tests.
- Sandbox evidence mistaken for outer isolation: record renderer ancestry,
  seccomp, privilege restrictions, and namespace separation, not just success.
- An older Electron version can carry backports: no claims or changes to T3.

## Execution state

Implementation, consumer checks, and delivery gates are complete. Task outputs
and the guarded publication receipt are under `out/chromium-security/`.

- The Bazel browser alias reports Google Chrome for Testing 156.0.8078.80.
- Draw.io and screenshot command E2Es pass with the Linux action sandbox.
- Runtime observation recorded 22 renderer processes, all with seccomp mode 2,
  NoNewPrivs set, no disabling flag, and separate PID namespaces from their
  browser ancestors. Those ancestors had seccomp mode 0, distinguishing the
  renderer filter from inherited Bazel confinement.
- Eleven additional checks pass: Mermaid renderer, native raster, aspect ratio,
  paint order, site search, site Mermaid, readability, X profile browser,
  video skill configuration, skill discovery, and Buildifier.
- Representative PNG inspection confirms literal markup and syntax formatting.
  Draw.io normal/shadow SVGs retain the Service label and positive viewboxes;
  the shadow output retains its drop-shadow filter.
- The prepared candidate passes all 40 tests in the combined repository quality
  and consumer run, affected-target semantic lint, and strict OpenSpec
  validation. Module resolution and formatting completed successfully.
- Guarded publication created the task pull request and verified the remote
  head. Final publication identity remains owned by Git and the delivery receipt.
- The initial restricted Bazel invocation could not write the existing output
  base; approved normal-cache execution succeeds. An initial process observer
  missed content-addressed executable paths; the corrected observer starts its
  own test run and retains the successful evidence in
  `out/chromium-security/logs/renderer-sandbox.json`.
