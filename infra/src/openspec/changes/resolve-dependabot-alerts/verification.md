# Dependency remediation evidence

## November 4 follow-up

The live GitHub API inventory contains 63 open alerts across five manifests.
The final selections exclude every captured vulnerable range:

- Vast.ai: Pillow 13.0.0 and cryptography 50.0.2 resolve 54 alerts.
- Root and AL Go modules: Go JOSE v4.1.8 resolves six alerts.
- Promptfoo: glob 14.0.3 and minimatch 11.0.4 resolve two alerts; the existing
  Promptfoo patch remains unchanged. The user approved the newly required
  `@isaacs/nonrecursive-glob-match` package after review of its upstream owner
  and purpose. pnpm records its archive integrity; lifecycle scripts stay disabled.
- pkgsite: x/net v0.57.0 already excludes the reported `<0.55.0` range.

`out/dependabot-vastai/alert-reconciliation.json` records each alert number,
advisory, vulnerable range, and selected version. This establishes safe candidate
selections; GitHub closure still depends on default-branch integration and scanning.

Promptfoo's eight offline tests pass, including real configuration validation.
The AL package tests pass with updated Vault/JOSE dependencies. The linked
Vast.ai record describes CLI, image/PDF, and encryption compatibility evidence.
The AL module's native Go tidy passes. Root Go tidy cannot resolve existing
Bazel-generated Anki and worker-protocol packages; pinned Go get updates the
required module and its checksums without changing those unrelated packages.

Session review: `CREDENTIAL-STORE-SANDBOX` was resolved by reading the existing
host credential store through approved execution; no login or shared configuration
changed. `ROOT-GO-TIDY-GENERATED-PACKAGES` records the existing native-tidy limit
above. `NESTED-LINT-CONFIG` records that standalone Promptfoo has no lint preset;
use the root workspace's configured aspects on `@rules_promptfoo//:promptfoo`.
The initial Python compatibility probes also required the generated Bazel site
initialization; the successful probe uses that exact dependency runfiles tree.

Repository quality and Buildifier pass all 28 checks. Root semantic lint passes
for Vast.ai, its lock target, both root Vault commands, the AL package, and the
real Promptfoo binary. Formatting introduces no unrelated source changes.

## October 25 baseline

The 2026-10-25 GitHub API snapshot contains 250 open alerts across ten manifests.
One alert already selects a safe version (`tools/pkgsite/go.mod`, x/net 0.57.0).
The updates remove vulnerable selections for the other 249 alerts. Comparing
every captured range with the final candidate manifests and locks leaves zero
vulnerable selections. The baseline revision is
`58f99e2459ff83c06073b4ea0b5976d775968eac`.

## Unpatched dependency removal

Alert #324, [GHSA-wj6h-64fc-37mp](https://github.com/advisories/GHSA-wj6h-64fc-37mp),
affects all `ecdsa` releases. InfiniTime's `adafruit-nrfutil` uses this library
only for optional key signing. The pinned firmware CMake commands invoke
`dfu genpkg` without key-file arguments, and the helper already handles a
missing signing import. A uv dependency exclusion removes `ecdsa` and its unused
`six` dependency from the environment. No advisory was suppressed or dismissed.

The helper's key generation and signed DFU functionality are unavailable in
this environment. Rebuilding without ecdsa passes the firmware build test.
Whole-firmware hashes vary because SystemInfo.cpp embeds `__TIME__`. A controlled
comparison packages the same generated HEX input with and without the previous
pinned ecdsa/six libraries on PYTHONPATH. The helper reports its missing signing
import only in the excluded environment. Every BIN, DAT, and manifest member
hash matches between the two packages. The repeatable comparison is recorded
in `out/dependabot/dfu-no-ecdsa-verification.json`.

## Consumer evidence

- Cordis: all seven Bazel targets pass. The runtime E2E passes ten cases,
  including expression execution, YAML writes, scalar preservation, and restart.
  Its `yaml_roundtrip.json` undeclared output preserves the configuration and
  observed result (`answer: 42`, string `label: "true"`, boolean `enabled: false`).
- Promptfoo: six offline analysis, launcher, and configuration tests pass in its
  nested workspace; the real CLI reports 0.125.1. No model calls were needed.
- OpenSpec: its nested build test passes.
- Mermaid: renderer tests pass and generate SVG/WebP artifacts with the retained
  rendering patch and secure KaTeX dependency.
- InfiniTime: the build test passes and produces a valid
  `pinetime-mcuboot-app-dfu-1.15.0.zip` containing binary, DAT, and manifest files.
  SHA-256 for the observed build: `addf686ad11db1f33c9f807e24578d9417f4abc419175f4968500b915b64f663` (compile-time dependent).
- Python: native locks regenerate through the owning uv targets; Ansible lint
  and YubiKey Manager binaries build successfully.
- Go: the Git metadata CLI and pure Go Anki client build. The Git CLI emits verified
  commit and tag metadata from an isolated fixture repository. Al plugin transport tests pass with gRPC
  1.84.0; the repository delivery tests also pass with the updated module graph.

The broader Anki-backed client test cannot build because its unchanged Rust
`percent-encoding-iri` 2.2.0 archive lacks `Cargo.toml`. This is separate from
the dependency updates; the selected pure Go consumer and transport checks pass.

Raw inventories, range reconciliation, command logs, and delivery receipts live
in ignored `out/dependabot/`. Final candidate identity, mandatory quality and
semantic-lint results, and publication state belong to Git and delivery receipts.
