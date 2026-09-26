## Why

The repository needs one reusable way to run Molecule scenarios on disposable
local QEMU guests through Bazel. VM lifecycle, packaged inputs, credential
isolation, and retained evidence need a tooling owner that every scenario
consumer can share.

## What Changes

- Add repository-internal Molecule rules and a Go runner under `tools/molecule`.
- Reuse `tools/ansible` executables and package mappings; add Molecule through
  the existing Python dependency workflow rather than a second Ansible pin.
- Provide independent local QEMU lifecycles with pinned images, writable
  overlays, explicit runtime prerequisites, and failure/cancellation cleanup.
- Isolate inventory and credentials and retain sanitized Bazel test artifacts.
- Document execution, cache, network, and abrupt-exit recovery behavior.

## Capabilities

### New Capabilities

- `molecule-qemu-runner`: The shared Bazel/Molecule input, VM lifecycle,
  environment, cleanup, and evidence contract.

### Modified Capabilities

None.

## Impact

Implementation affects `tools/molecule`, its pinned external tool inputs, and
the owning Ansible Python lock workflow. Reuse the existing Fedora cloud-image
pin. Initial support is Linux x86-64 with local QEMU; remote hypervisors,
persistent host provisioning, and production inventories are excluded.

The first consumer is the linked
[collection role-test change](../../../../../projects/ansible_collection/openspec/changes/add-qemu-molecule-role-tests/proposal.md),
which owns scenarios for `host`, `traefik`, and `forgejo`. This change owns
generic runner behavior and its smoke/failure acceptance; it does not define
role-specific assertions. Both changes can be implemented together while
preserving their separate specification ownership.
