## 1. Runner

- [x] 1.1 Implement the Linux process wrapper and Bazel entry point; build and inspect its CLI.
- [x] 1.2 Verify real scope limits, argument/environment/stdin/cwd preservation, nonzero exit, interrupt cleanup, and failure cases; preserve an E2E JSON report.

## 2. Delivery

- [x] 2.1 Run repository quality and affected-target lint; preserve check logs.

## Earlier CLI-only verification

On 2026-11-01, `//tools/systemd_run` built and a real-scope E2E run verified
128 MiB maximum RAM, 96 MiB soft RAM, zero swap, literal argv, environment,
stdin, cwd, child exit 7, and interrupt exit 130 with the child stopped.
Missing executable/properties, malformed properties, missing systemd-run,
unsupported properties, and an unavailable child executable failed as expected.
The actual ComfyUI integration also served HTTP 200 with its requested limits
and stopped cleanly. Evidence and reproduction scripts are under ignored
`out/systemd-run/`, including `runner-e2e.json` and `comfyui-e2e.json`.

Strict OpenSpec validation, all 27 repository quality/Buildifier tests, and
affected-package semantic lint passed. The final candidate and publication
are tracked by the delivery receipts in the same output directory.

## 3. Executable Bazel rule

- [x] 3.1 Add the executable rule and runtime metadata loading; build and inspect its generated configuration.
- [x] 3.2 Verify literal arguments, runfiles/data, runtime environment, limits, and nonzero exit through a normal Bazel fixture; preserve JSON evidence.
- [x] 3.3 Revalidate CLI compatibility, actual ComfyUI startup, repository quality, and semantic lint; preserve delivery checks.

## Rule verification

On 2026-11-01, the normal `//tools/systemd_run/test:probe` target verified
literal caller arguments (including spaces, shell metacharacters, and an empty
argument), expanded runfiles data, wrapped runtime environment, kernel limits,
and child exit 7. The initial fixture identified a missing child runfiles
environment; forwarding the pinned Go resolver's environment corrected it.
CLI compatibility checks still passed, including interrupt cleanup.

The normal ComfyUI target served HTTP 200 under 32/36 GiB limits and zero swap,
then stopped cleanly. Its normal benchmark target generated and decoded a
64-square one-step PNG. All 28 quality/Buildifier checks and affected-target
semantic lint passed. JSON reports, logs, reproduction scripts, and final
candidate/publication receipts are under ignored `out/systemd-rule/`.

## 4. Standalone rules module

- [x] 4.1 Move the rule and launcher to `tools/rules_systemd`, add module metadata and public API, and register parent integration.
- [x] 4.2 Verify standalone build/test, real-scope fixture, parent ComfyUI startup/assets, and workspace discovery; retain evidence.

## Module verification

On 2026-11-01, the standalone module built and its public-rule build test
passed. The real-scope probe verified expanded runfiles data, runtime
environment, 96/128 MiB soft/hard memory limits, zero swap, literal spaces
and dollar expressions, empty arguments, and child exit 7. Explicitly
disabling systemd environment expansion corrected a host-version-dependent
argument transformation. The parent ComfyUI main target then served its
required frontend assets and local HTTPS proxy on port 8188 under its
existing limits and stopped cleanly. Evidence is under ignored
`out/comfyui-frontend/`, including `module-probe-final.json` and
`module-startup/comfyui-e2e.json`. Repository matrix coverage is recorded
separately by the full-repository check report and final delivery receipts.
