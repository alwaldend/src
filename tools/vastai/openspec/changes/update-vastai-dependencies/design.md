## Context

See proposal.md for scope. PyPI currently publishes Vast.ai 1.10.0, which pins
Pillow 12.2.0 and cryptography 49.0.0. The registry identifies security advisories
against both. Current stable releases are Pillow 13.0.0 and cryptography 50.0.2.

## Decisions

Proceed with uv dependency overrides for those two packages. An ordinary
constraint cannot resolve against an upstream exact pin, and merely updating
Vast.ai retains both vulnerable versions. Removing the CLI would break its
contract. Overrides preserve the upstream application and replace only existing
dependencies; the generated lock remains the exact version and hash authority.
Leave Vast.ai unconstrained in the manifest so future lock refreshes can upgrade it.

Use Python >=3.11, matching secure Pillow support and the existing Bazel runtime.
Keep the component environment isolated to avoid unrelated dependency changes.

## Risks / Trade-offs

Newer Pillow and cryptography can change APIs. Before acceptance, inspect upstream
usage and exercise CLI help, version reporting, and relevant offline behavior.
Record any compatibility failure before changing the chosen dependency versions.

The restricted sandbox cannot access the host GitHub credential store. With
host credential-store access, the API inventory confirms 54 Vast.ai alerts
among 63 open alerts. The coordinated follow-up in
`infra/src/openspec/changes/resolve-dependabot-alerts` owns the other manifests.
Raw metadata and command logs belong in ignored `out/dependabot-vastai/`.

## Verification

The generated lock selects Pillow 13.0.0 and cryptography 50.0.2. All 41
registry package selections have no active advisories in PyPI's release
metadata. All 54 captured Vast.ai alert ranges exclude the selected versions.
The CLI's help and version commands pass. The unchanged Borb 2.1.25 dependency
successfully writes an image-containing PDF using Pillow 13, and cryptography
completes an encryption/decryption round trip. Artifacts are
`out/dependabot-vastai/verification-image.pdf` and `dependency-compatibility.json`.

The legacy upstream invoice helper imports `Document` from a Borb package that
does not export it. Both the Vast.ai application wheel and Borb wheel are
unchanged from the baseline, so that existing defect is outside this update.
