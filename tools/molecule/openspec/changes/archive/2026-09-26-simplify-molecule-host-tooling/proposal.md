## Why

PR #112 review rejects an obscure third-party QEMU distribution and asks whether
Molecule setup needs that complexity. The supported host already provides
Fedora-packaged QEMU, so maintaining another binary toolchain and firmware
archive is unnecessary.

## What Changes

- Use host QEMU and qemu-img with their distribution-provided firmware, checked
  before convergence and identified in retained results.
- Remove the third-party QEMU package and Bazel dependency wiring.
- Put the one-command workflow and host prerequisites first in the README;
  retain the existing guest lifecycle, isolation, cleanup, and evidence contract.
- Add publisher-reputation guidance and evaluation cases to the canonical
  external-dependency skill.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `molecule-qemu-runner`: Explicit host QEMU prerequisites and runtime provenance.

## Impact

Changes `tools/molecule`, removes `third_party/com_github_hermeticbuild_qemu`,
regenerates root module wiring, and updates
`tools/agents/skills/repo-external-dependency`. No host package installation,
production configuration, or collection-role changes are required.
